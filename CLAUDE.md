# CLAUDE.md · Reglas del proyecto SALUS

SALUS renta equipo médico y surte insumos a varios hospitales para cirugías. El contexto completo está en `docs/` (01 análisis, 02 requerimientos, 03 modelo de datos, 04 pendientes, 05 guía de diseño, 06 plan de trabajo, 07 alias de Firebird, 08 arquitectura). Consúltalos antes de implementar algo del negocio.

**`docs/09-decisiones-tecnicas.md` es la fuente de verdad técnica**: léelo antes de diseñar o escribir código. No contradigas ni reinventes lo que ahí está decidido; si algo debe cambiar, propónlo al usuario. Cada decisión nueva que se valide se agrega ahí.

## Forma de trabajo
- Avanzar por **pasos cortos** del plan (`docs/06-plan-de-trabajo.md`). Cada paso se entrega funcionando y se revisa antes del siguiente.
- No adelantar trabajo de otros pasos aunque parezca fácil.
- No inventar reglas de negocio, estados ni contenido (p. ej. plantillas). Si falta información, preguntar **solo lo que se desprende de los docs**, sin preguntas especulativas.
- Al terminar un paso: compila, pasa lint y se probó. Decir claramente qué quedó y qué falta.
- Responder en español.

## Codificación de archivos (obligatorio)
- Todo archivo de texto en **UTF-8 sin BOM** y fin de línea **LF** (`.editorconfig` y `.gitattributes` lo fijan).
- Crear y editar archivos con las herramientas de edición, **nunca con PowerShell 5.1** (`Set-Content`, `Out-File`, `>`): lee como ANSI, agrega BOM y rompe la ñ y los acentos.
- Si hay que transformar un archivo por consola, usar Bash (sed, heredoc).
- Después de escribir texto en español, verificar que no aparezcan `Ã`, `Â` ni BOM.

## Stack (no proponer alternativas)
- `frontend/`: Vue 3 + Vuetify + TypeScript + Pinia + Vue Router + ESLint.
- `backend/`: Go + Fiber v3.
- BD: Firebird 5 en Docker (`firebirdsql/firebird`) con **alias** en `databases.conf`. La ruta real del `.fdb` nunca va en la cadena de conexión.
- Todo se levanta con `docker compose`; cada servicio (BD, back, front) en su propio contenedor.
- **pnpm siempre; npm prohibido** (comandos, scripts, Dockerfiles, README).
- Documentación de la API con **Scalar**, nunca Swagger UI.

## Backend (Go)
- Estructura idiomática y simple: `cmd/api` + `internal/` con `config`, `database`, `model`, `repository`, `service`, `handler`, `router`, `middleware`. **Sin sobreingeniería**: nada de capas o interfaces que no se usen.
- Secretos y configuración por variables de entorno; `.env` nunca se sube (solo `.env.example`).
- SQL parametrizado siempre.
- Respuestas solo con `internal/response` (`OK`, `Created`, `List`, `NoContent` y los errores `*response.Error`); nunca `c.JSON` con formas propias.
- Migraciones en `internal/database/migrations/NNNN_descripcion.sql`; una migración aplicada **no se edita**, se crea otra.
- Logs con `slog`; nunca registrar contraseñas, tokens, cookies ni cuerpos de petición.
- La API se conecta a Firebird con `SALUS_APP`, nunca con SYSDBA.
- **Gorutinas con cuidado**: solo si hacen falta, cada una con forma clara de terminar (context) y sin fugas ni carreras. Consultar la documentación de Go antes de un patrón concurrente nuevo.
- Comentarios internos en las funciones: directos y concisos, lo justo para seguir el hilo.
- Consultas: sin N+1 (JOIN o `IN` por lote) y sin OFFSET grande; cuando un listado crezca, paginar por cursor (keyset).
- Llamadas a servicios externos: un paquete por servicio en `internal/client/`, con `http.Client` con timeout y `context`; credenciales por variables de entorno.

## API (consistente en todo el proyecto)
- REST, versión en la ruta: `/api/v1/...`; recursos en plural y en español (`/api/v1/hospitales`).
- Verbos y códigos HTTP correctos (200, 201, 204, 400, 401, 403, 404, 409, 422, 429, 500).
- JSON en `camelCase`. Un **solo formato de error** para toda la API.
- Listados con los mismos parámetros de paginación, búsqueda, filtros y orden.
- Validación de entrada en cada endpoint y autorización por rol en el back.
- **Extensible sin romper**: respuestas siempre como objeto (`data`/`meta`, nunca un arreglo suelto) y peticiones como objeto JSON, para poder agregar campos después. Los cambios son aditivos; lo que rompa compatibilidad va a `/api/v2`.
- CORS solo con orígenes explícitos (nunca `*`); en producción el front va por el mismo origen y CORS queda vacío.

## Seguridad de usuarios
- Contraseñas con hash y sal (argon2id o bcrypt); nunca en texto plano ni en logs.
- Límite de intentos por usuario y por IP en una ventana de tiempo, con bloqueo temporal (anti fuerza bruta).
- Roles guardados en BD (Sysadmin, Admin, Doctor). Validación de rol y de datos en **back y front**; el back es la autoridad.
- Tokens de acceso de vida corta + refresh revocable. Auditoría de accesos.

## Frontend y diseño
- **Regla de oro: mobile first y Vuetify first.** Componentes propios solo en casos especiales (cards de métricas).
- Tema, colores, defaults de componentes y aliases **solo** en `src/plugins/vuetify.ts`. Sin colores, CSS ni clases sueltas en los `.vue`; sin props repetidas que deberían ir en `defaults`.
- Iconos Lucide solo por alias (`$nombre`) registrados en `src/plugins/icons.ts`.
- **Todo campo de formulario lleva `:rules`** de `src/utils/rules.ts` (sin librerías externas).
- **Skeletons** (`v-skeleton-loader`) solo en el contenido de la página, nunca en sidebar ni header; un solo skeleton por página si se puede, y skeletons propios solo para partes con peticiones independientes. También al paginar listas y tablas.
- **Transiciones de Vuetify** (`v-fade-transition`, `v-expand-transition`, etc.) cuando algo aparece o desaparece; sin animaciones CSS propias.
- Fuentes: Poppins para títulos, Inter para texto.
- **Antes de construir o ajustar una pantalla, leer su artboard del mock** (https://claude.ai/artifact/3VR1enc7SDGVFu464sF8xg, `project/<Pantalla>.dc.html`, y `Componentes.dc.html` para elementos de UI) y seguirlo: estructura, textos y jerarquía. No inventar diseño.
- Componentes de UI del mock ya hechos: `PageHeader` (título + subtítulo + acciones), `AppMarca`, `ErrorPagina` (404/403/500/503). Reusarlos, no repetir su diseño a mano.
- **Sin espaciados manuales** (`pa-*`, `ma-*`, `mt-*`…): los dan el grid (`v-row`/`v-col`) y los componentes; si hay que cambiarlos, en la config global. Densidad `comfortable` global; `v-list` con `slim` y sin `nav`.
- Tipografía con las clases MD3 de Vuetify 4 (`text-headline-small`, `text-title-large`, `text-body-medium`, `text-label-large`…), no `text-h*`.
- Páginas de error dentro del dashboard se pintan en el área de contenido, centradas (nunca reemplazan sidebar ni header).
- Seguir `docs/05-guia-diseno.md`. Resumen: mobile first, Poppins, primario #0195B7 / secundario #4FB5AB, campos outlined, modo oscuro requerido, sin fondos negros, en pasteles el texto lleva el mismo tono más fuerte.
- Botón primario sólido; secundario solo texto; de icono sin borde ni fondo. Si dispara una carga: se deshabilita, muestra spinner y dice "Cargando".
- El usuario da instrucciones de diseño al llegar al paso del front; no adelantarse.
- **Errores al usuario siempre amigables** (toast o alerta según el caso): nunca mostrar códigos, enums, stack, `requestId` crudo ni mensajes técnicos. El front traduce cada `code` de la API a lenguaje del usuario; lo técnico queda en los logs.

## Negocio (lo que más se confunde)
- **Equipo** (reutilizable, modelo maestro + unidades, se cobra **solo por hora**) ≠ **Insumo** (desechable, por stock, los no usados regresan).
- La cirugía tiene **solo 5 estados**: Programada → En entrega → En hospital → Devuelta → Cerrada (+ Cancelada). Sin estados intermedios.
- Cotización inicial al programar; la final solo descuenta lo retornado y suma gastos cobrables.
- **Fuera de alcance**: facturación, pagos y contabilidad. Roles User SALUS y User Doc en espera.

## Comandos
```bash
# Backend
cd backend && go run ./cmd/api
cd backend && go vet ./... && go test ./...

# Frontend
cd frontend && pnpm dev
cd frontend && pnpm lint && pnpm build

# Todo
docker compose up -d
```
