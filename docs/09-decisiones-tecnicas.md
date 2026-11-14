# 09 · Decisiones técnicas y buenas prácticas

Registro de lo que **ya está decidido y validado**. Antes de proponer algo, revisar aquí: no se reinventa ni se cambia sin una razón nueva y el visto bueno del usuario. Cada decisión nueva se agrega a este archivo.

## 1. Proyecto y forma de trabajo
| Decisión | Detalle |
|---|---|
| Avance | Pasos cortos del plan (`06`), cada uno funcionando, probado y revisado antes del siguiente. |
| Alcance | No inventar reglas de negocio, estados ni contenido. Preguntar solo lo que sale de los docs. |
| Codificación | UTF-8 sin BOM y LF en todo. Nunca escribir archivos con PowerShell 5.1. |
| Paquetes JS | Solo pnpm. npm prohibido. |
| Git | Repo `salus-salata` en GitHub; commits pequeños por paso con mensaje en español. |

## 2. Arquitectura
| Decisión | Detalle |
|---|---|
| Contenedores | `frontend` (nginx + Vue), `backend` (Go), `db` (Firebird 5.0.4), orquestados con docker compose. Ver `08`. |
| Puerta única | Solo nginx se publica; reenvía `/api` al backend. Mismo origen → sin CORS en producción. |
| Base de datos | Firebird con **alias** `salus` en `databases.conf`; nunca la ruta real. `DatabaseAccess = Restrict` al volumen. |
| Usuario de BD | La API usa `SALUS_APP` (dueño de la base). SYSDBA solo para mantenimiento; la API se niega a arrancar con él. |
| Imagen backend | Multi-etapa, distroless, usuario no root, solo lectura, sin capacidades. Healthcheck con el propio binario. |

## 3. Backend (Go + Fiber v3)
| Decisión | Detalle |
|---|---|
| Estructura | `cmd/api` + `internal/{config, database, model, repository, service, handler, router, middleware, response, security, docs}`. Sin capas ni interfaces que no se usen. |
| Capas | `handler` (HTTP: leer, validar, responder) → `service` (reglas, sin HTTP) → `repository` (único lugar con SQL). |
| Configuración | Variables de entorno validadas al arrancar; si algo falta, no arranca y dice qué. |
| Validación | `go-playground/validator` con etiquetas `validate` y mensajes en español por campo. |
| Logs | `log/slog` (librería estándar). Texto en desarrollo, JSON en producción. Nunca contraseñas, tokens, cookies ni cuerpos. |
| Gorutinas | Solo si hacen falta, con forma clara de terminar. Hoy la única es la del servidor HTTP. Nada de procesos en segundo plano para limpiezas: se hacen en el flujo (p. ej. sesiones vencidas al iniciar sesión). |
| Zonas horarias | Fechas como `TIMESTAMP WITH TIME ZONE`, se guardan en UTC desde Go. `time/tzdata` va dentro del binario. |
| Servicios externos | Un paquete por servicio en `internal/client/`, `http.Client` con timeout y `context`, credenciales en `.env`. |

## 4. Base de datos y consultas
| Decisión | Detalle |
|---|---|
| Migraciones | `internal/database/migrations/NNNN_descripcion.sql`, embebidas, en orden, una transacción cada una, con checksum: **una aplicada no se edita**. |
| Firebird + transacciones | Firebird no permite usar en la misma transacción una tabla creada en ella: los datos iniciales de una tabla nueva van en la **siguiente** migración. |
| SQL | Siempre parametrizado (`?`). |
| Transacciones | **Atomicidad e integridad**: toda operación que escribe en más de un paso, o que lee para luego escribir, va en una transacción con `service.enTx` (todo o nada). Padre e hijos (hospital + contactos, equipo + tarifa + unidades) se guardan juntos. Llaves foráneas siempre declaradas y con nombre; nada de datos a medias ni huérfanos. Lo lento (argon2) se calcula antes de abrir la transacción. Única excepción: la auditoría de intentos fallidos se guarda aunque la operación falle. |
| N+1 | Prohibido. Relaciones con `JOIN` o `IN` por lote. |
| Paginación | `OFFSET … ROWS FETCH NEXT … ROWS ONLY` solo en tablas pequeñas (usuarios, catálogos). Tablas que crecen (cirugías, movimientos, auditoría) → paginación por cursor (keyset). |
| Identificadores | **UUID v7** en todas las tablas (también `ROL`), generado en Go (`model.ID`) y guardado como `BINARY(16)`. En JSON viaja como texto (`"0193…-…"`). v7 lleva el tiempo al inicio: los índices crecen en orden y no se fragmentan. En DBeaver se lee con `UUID_TO_CHAR(ID)`. En desarrollo el esquema de Fase 0 se rehízo directamente con UUID (sin migraciones de conversión). Usuario y correo en minúsculas con `UNIQUE`. |
| Roles en código | El código se refiere a un rol por su `CLAVE` (`SYSADMIN`, `ADMIN`, `DOCTOR`), nunca por su id. |
| Textos que se comparan | Nombres y búsquedas en columnas `COLLATE UNICODE_CI_AI`: "holep" = "Holep" y "cirugia" = "Cirugía" al validar duplicados y al buscar. Se usa desde la Fase 1 en las columnas nuevas; los índices de texto completo siguen para la Fase 9 si hacen falta. |
| Listas cerradas | Si los valores son fijos y conocidos (p. ej. puesto del técnico, estado de la república), no se crea tabla: columna `VARCHAR` con `CHECK (… IN (…))`, la misma lista en Go (`validate:"oneof=…"`) y un `v-select` en el front. Agregar un valor = migración nueva que cambia el `CHECK`. |
| Duplicados y baja lógica | Las **claves** (equipo, insumo, hospital) son únicas para siempre, aunque el registro se dé de baja: las referencias históricas no cambian de dueño. Los **nombres** solo se validan entre registros activos (índice único parcial `WHERE ACTIVO`, Firebird 5). |
| Claves entre equipo e insumo | Tabla `CLAVE_CATALOGO (CLAVE PK, TIPO)`: equipo e insumo registran ahí su clave en la misma transacción, así la BD impide que se repita entre los dos (RF-08). |
| Tarifas | Nunca se sobrescriben: cambiar el costo por hora crea una fila nueva en `TARIFA_EQUIPO` con su `VIGENTE_DESDE`. Da el historial de precios y la cotización usa la vigente en la fecha de la cirugía. |
| Stock | Nunca negativo: `CHECK (STOCK_ACTUAL >= 0)` en la BD y validación en el back. Cuando lleguen los movimientos (Fase 7), una salida que deje el stock en negativo se rechaza con un aviso. |
| Purga | La decide el usuario por equipo (`REQUIERE_PURGA`). `PURGA_CADA` solo es obligatorio si requiere purga (`CHECK`). |
| Direcciones (México) | Desglosadas: calle, número exterior, número interior, colonia, C.P. (5 dígitos), municipio y estado (lista cerrada de las 32 entidades). Mismo bloque en Go y en el front para doctor y hospital. |

## 5. API
| Decisión | Detalle |
|---|---|
| Rutas | REST, `/api/v1/...`, recursos en plural y en español. Lo que rompa compatibilidad va a `/api/v2`. |
| Formato | JSON `camelCase`. Éxito `{data}`; listados `{data, meta:{page,pageSize,total}}`; nunca un arreglo suelto. |
| Errores | `{error:{code, message, details?, requestId}}`. `code` estable para el front; `message` en lenguaje del usuario. Un 500 nunca expone detalles. |
| Extensibilidad | Peticiones y respuestas como objeto; campos desconocidos se ignoran; los cambios son aditivos. |
| Paginación | `?page=1&pageSize=20`, máximo 100; fuera de rango → 422. |
| Búsqueda y orden | `?q=texto` (contiene, sin acentos ni mayúsculas) y `?sort=campo` o `?sort=-campo` (descendente). Cada recurso declara qué campos se pueden ordenar; otro valor → 422. |
| CRUD de catálogos | `GET /recurso` (listado, solo activos), `GET /recurso/{id}`, `POST` → 201, `PUT /recurso/{id}` → 200, `DELETE /recurso/{id}` → 204 = **baja lógica** (`ACTIVO = false`). Un id inexistente o dado de baja → 404; un duplicado → 409 con el campo en `details`. |
| Documentación | OpenAPI 3.1 escrito a mano en `internal/docs/openapi.yaml`, servido con Scalar en `/docs` solo en desarrollo. Se actualiza con cada endpoint. |

## 6. Seguridad
| Decisión | Detalle |
|---|---|
| Contraseñas | **argon2id** (19 MiB, t=2, p=1, OWASP) con sal aleatoria, formato PHC. Rehash automático si cambian los parámetros. |
| Política | 12–128 caracteres; no contener usuario ni correo (NIST 800-63B: largo antes que reglas de complejidad). |
| Token de acceso | JWT HS256, 15 min, en `Authorization: Bearer`. Solo se acepta HS256 (rechaza `none`). Secreto `JWT_SECRET` ≥ 32 caracteres. |
| Refresh | Aleatorio de 256 bits en cookie `salus_refresh` (`httpOnly`, `SameSite=Strict`, `Secure` en producción, ruta `/api/v1/auth`). En BD solo su SHA-256. 7 días. |
| Rotación | Cada refresh revoca el anterior y emite otro de la misma familia. Si llega uno ya usado → se revoca la familia completa (posible robo). |
| Fuerza bruta | 10 logins/min por IP (limiter) + máx. 5 fallidos por cuenta y 20 por IP en 15 min (contados en `AUDITORIA_ACCESO`) → 429. |
| Anti-enumeración | Mismo mensaje y mismo tiempo para usuario inexistente, contraseña incorrecta o cuenta inactiva (se verifica contra un hash falso). |
| Primer acceso | El Sysadmin inicial se crea desde `INITIAL_ADMIN_*` solo si no hay ninguno, y debe cambiar la contraseña: hasta entonces la API responde `PASSWORD_CHANGE_REQUIRED` en todo salvo `me`, `contrasena` y `logout`. |
| Cambio de contraseña | Cierra todas las sesiones del usuario y abre una nueva. |
| Roles | En BD (`ROL`). El token lleva el rol; `RequireRole` autoriza en cada ruta. El back es la autoridad; el front solo oculta. |
| Auditoría | `AUDITORIA_ACCESO`: login correcto, fallido, bloqueado, logout, reuso de refresh y cambio de contraseña, con IP y User-Agent. |
| Protección general | Límite por IP con ventana deslizante (120/min general, 30/min escrituras), `User-Agent` obligatorio en `/api`, escrituras solo JSON, helmet, 1 MB por petición, timeouts y tope de conexiones. |
| CORS | Solo orígenes exactos; sin `*`; https obligatorio en producción. En producción no se usa. |
| IP real | nginx manda `X-Real-IP`; el backend solo la lee si la petición viene de `TRUSTED_PROXIES` (la IP fija de nginx, `172.28.0.10` en la red `salus`). Directo al backend se ignora. |

## 7. Frontend
| Decisión | Detalle |
|---|---|
| Regla de oro | **Mobile first y Vuetify first.** Componentes propios solo en casos especiales (p. ej. cards de métricas). |
| Configuración global | `src/plugins/vuetify.ts` es el único lugar para tema, colores, defaults y aliases. Sin colores, CSS ni clases sueltas en los `.vue` salvo casos especiales justificados. |
| Tema | `salus` (claro) y `salusOscuro`, con la paleta de `05`. Colores propios: `brand`, `transito`, `text-secondary`, `text-muted`. Se usan por nombre (`color="..."`, `text-text-muted`). |
| Defaults | Cada componente se predefine en `defaults` (campos outlined, botón texto, cards con borde, alertas tonales, drawer con borde, etc.). En los `.vue` solo van las props propias del caso. |
| Aliases | Componentes virtuales para variantes repetidas: `btn-primario`, `btn-icono`, `campo-busqueda`. |
| Iconos | Lucide (trazo suave) como set por defecto, registrados en `src/plugins/icons.ts` y usados **solo por alias** (`$usuarios`); se importan uno por uno (tree-shaking). |
| Tipografía | Poppins 400/500/600/700 desde `@fontsource/poppins` (sin CDN). |
| SASS | `src/styles/settings.scss`: fuente, botones sin mayúsculas y `$color-pack: false` (quita la paleta Material). |
| Tree-shaking | `vite-plugin-vuetify` con `autoImport`: solo se empaqueta lo que se usa. Páginas cargadas bajo demanda en el router. |
| Espaciado | 1) **Grid primero**: bloques de la página en `v-row` (density `comfortable` global) / `v-col`; el gutter es el espacio estándar. 2) Dentro de un `d-flex`, solo `ga-*` con la escala **2 (8 px)** entre elementos relacionados (botones, icono+texto) y **4 (16 px)** entre grupos. 3) `pa-4` solo en contenedores que no traen padding propio (p. ej. el prepend del drawer). 4) **Nunca márgenes** (`ma-*`, `mt-*`, `mb-*`…) ni valores fuera de la escala. 5) Si un espaciado debe cambiar en toda la app, se cambia en `defaults` o en SASS, no en cada archivo. |
| Validación | **Todo campo lleva `:rules`** con las reglas de `src/utils/rules.ts` (sin librerías). Mismos límites que el back. |
| Carga (skeletons) | `v-skeleton-loader` solo en el **contenido de la página**, nunca en sidebar ni header. Preferible un solo skeleton para toda la página (una animación de un jalón); si partes de la UI dependen de peticiones distintas, cada una lleva su propio skeleton para no bloquear el resto. También en listas y tablas al paginar (`loading` de `v-data-table` con skeleton de filas). |
| Transiciones | Usar las transiciones de Vuetify (`v-fade-transition`, `v-expand-transition`, `v-slide-y-transition`, prop `transition`) cuando algo aparece o desaparece, para evitar saltos. Sin animaciones CSS propias. |
| Fuentes | Títulos con **Poppins** (600/700) y texto con **Inter** (400/500), vía `$heading-font-family` y `$body-font-family`. |
| Menú | `src/config/menu.ts`: secciones, opción sola u opción con hijos (módulo + altas), filtrado por rol. |
| Diseño | `05-guia-diseno.md` y los mocks; el usuario da instrucciones al llegar a cada pantalla. |
| Errores | Siempre amigables (toast o alerta). Nunca códigos, enums, stack ni `requestId` crudo. El front traduce cada `code` de la API. |
| Patrón CRUD (front) | Un catálogo nuevo = servicio con `crearRecurso<T, D>(ruta)` + página de listado (`ListadoCrud` + `useListado` + `useBaja` + `DialogoEliminar`) + página de formulario (`FormularioCrud` + `SeccionFormulario` + `useFormulario`) + rutas, pestaña en `config/catalogos.ts` y menú. Todo el diseño del listado, formulario, acciones, baja y estados vacíos vive en esos componentes: cambiarlo ahí cambia todos los catálogos. Toasts con `useAvisosStore().exito()/error()` (`v-snackbar-queue`). Estados vacíos con `v-empty-state`. En móvil la lista abre la edición y la baja se hace desde el formulario. |
| Patrón CRUD (back) | Por catálogo: migración, `model`, `repository` (lista blanca de orden e índices únicos → campo), `service` (escrituras de varios pasos en `enTx`), `handler` con `consulta()`, `idParam()` y `errorCatalogo()`, y una línea `crud(...)` en el router. |
| Token de acceso | Solo en memoria (Pinia), nunca en `localStorage`. Al recargar la página se pide uno nuevo con `/auth/refresh`. |
| Refresh | Una sola renovación a la vez: si varias peticiones reciben 401 juntas, esperan la misma promesa (dos renovaciones simultáneas parecen robo y cierran la sesión). |
| Roles | Guardas de ruta y menús por rol, solo para la experiencia; la seguridad real está en el back. |
| Validación | Las mismas reglas que el back (largos, obligatorios) para avisar antes de enviar. |

## 8. Integración (paso 0.6)
| Decisión | Detalle |
|---|---|
| Puerta única | Contenedor `frontend`: `nginxinc/nginx-unprivileged` (sin root, solo lectura) sirve `dist/` y reenvía `/api/` a `backend:8080`. `/docs` no se expone por nginx. Cabeceras de seguridad y CSP en `frontend/nginx.conf`. |
| Build del front | Node 24 + corepack (instala el pnpm fijado en `package.json`); `pnpm install --frozen-lockfile`. |
| Desarrollo | Vite reenvía `/api` a `localhost:8080`: mismo origen que producción, sin CORS. |
| Cliente HTTP | `src/services/api.ts`: único punto de llamadas; normaliza errores a `ApiError` con mensaje amigable; ante 401 renueva una sola vez (promesa compartida) y reintenta. |
| Sesión en el front | `src/stores/auth.ts`: token en memoria; al abrir la app se recupera con `/auth/refresh`; si la sesión se pierde trabajando, vuelve al login con `?volver=` (solo rutas internas). |
| Guardas | `meta.publica` y `meta.roles` en las rutas; primer acceso obliga a `/cambiar-contrasena`; rol sin permiso → página 403 dentro del layout. |
| Pruebas | `scripts/prueba-integracion.sh`: SPA, CSP, IP real detrás de nginx, cookie, sesión completa y limpieza de datos de prueba. |

## 9. Ajustes posteriores
| Decisión | Detalle |
|---|---|
| Mensajes de campos | Sin `hideDetails: 'auto'`: cada campo reserva el espacio del mensaje y no salta al validar. Solo `campo-busqueda` usa `hide-details`. |
| Pantallas de acceso | `AuthLayout`: panel de marca (color primario, frase sin nombre de usuario y animación Lottie `public/hospital-lottie.json`) + formulario con la marca centrada y pie legal. En móvil solo el formulario. Lo usan login y cambio de contraseña. |
| Pestañas viejas | Tras un despliegue, una pestaña abierta pide archivos de pantallas que ya cambiaron de nombre (404). `router.onError` detecta el error de carga y recarga la app en la pantalla destino (una vez por destino cada 10 s). El login no borra la contraseña antes de confirmar la navegación. |
| Lottie | `lottie-web` build *light* (sin `eval`, compatible con la CSP), en el componente `AnimacionLottie`; respeta "reducir movimiento". |
| Producción | `docker-compose.prod.yml` encima del base: solo nginx publica puertos, HTTPS obligatorio con HSTS, `APP_ENV=production`, sin `/docs`, sin CORS, cookies `Secure`, logs JSON. Certificados en `infra/certs/` (no se suben). |
| Pruebas | `scripts/prueba-integracion.sh`, `scripts/prueba-seguridad.sh` (ataques) y `scripts/prueba-produccion.sh`. `backend/cmd/hash` genera hashes para usuarios temporales de prueba. |
| PWA | Instalable (manifest + iconos + theme-color), **sin service worker ni modo sin conexión**. Logos claro/oscuro en `src/assets`, `AppMarca` cambia según el tema. |
| Contraseñas | Componente `CampoContrasena` (ojo para mostrar/ocultar) en todos los campos de contraseña. |
| Alertas de formulario | Alias `AlertaFormulario` con `mb-5`: única excepción a "sin márgenes"; compensa el espacio de mensajes que reservan los campos. |
| Marca | `AppMarca` muestra solo el logo (sin el nombre). |
| Cambio de contraseña | Comparte el bloqueo del login: cada contraseña actual incorrecta cuenta como intento fallido de la cuenta y de la IP; además, límite de 10/min. |
| Entornos | Pruebas y producción como proyectos compose separados (`salus-pruebas` / `salus-prod`) con su propio volumen, red, puertos y secretos; `scripts/entorno.sh`; respaldos con `gbak` y copia de producción a pruebas. Ver `10-entornos.md`. |
| Sesión | Token de acceso 15 min; sesión renovable 7 días desde el último uso (`REFRESH_TOKEN_TTL`). |
