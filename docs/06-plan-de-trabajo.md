# 06 · Plan de trabajo por fases

Cada fase entrega algo **que corre en `docker compose up` y se puede probar**. Lo esencial va primero (fases 0–4 = el sistema reemplaza al Excel en el día a día); lo demás son detalles y agregados encima de esa base.

Reglas para todas las fases:
- Stack: Vue 3 + Vuetify + TS (pnpm), Go Fiber v3, Firebird en Docker con alias. API documentada en Scalar.
- Diseño según `05-guia-diseno.md` (mobile first, Poppins, outlined, modo oscuro).
- Cada endpoint nuevo valida rol en el back. Cada fase cierra con migraciones, seed de prueba y Scalar actualizado.
- Lo que no entra en la fase se deja fuera aunque "ya que estamos" parezca fácil.

---

## Bloque A · Lo esencial

### Fase 0 · Cimientos (por pasos)
Cada paso se entrega y se revisa antes de seguir.

| Paso | Qué entrega | Listo cuando |
|---|---|---|
| **0.1 Inicializar** ✅ | `git init`, `.gitignore`, `README`, `backend/` (Go module + Fiber v3 con `/health`) y `frontend/` (`pnpm create vuetify`, TS + ESLint) | Ambos compilan y arrancan por separado |
| **0.2 Firebird** ✅ | Contenedor `firebirdsql/firebird` con alias (`databases.conf`), volumen persistente, credenciales en `.env` | Me conecto por el alias sin ver la ruta real |
| **0.3 Backend base** ✅ | Estructura idiomática de Go (`cmd/`, `internal/` con `model`, `repository`, `service`, `handler`, `router`, `config`, `middleware`), conexión a Firebird, migraciones versionadas, respuestas y errores uniformes, Scalar en `/docs` | `/health` reporta la BD conectada y `/docs` muestra la API |
| **0.4 Usuarios y seguridad** ✅ | Contraseñas con hash y sal (argon2id o bcrypt), límite de intentos por usuario y por IP en una ventana de tiempo, bloqueo temporal, roles en BD, validación de datos y de rol en el back, tokens de acceso cortos + refresh revocable, auditoría de accesos, límite estricto por IP en `/auth` y respuestas iguales para usuario inexistente o contraseña incorrecta (no revelar qué usuarios existen). Seed de un Sysadmin | Las pruebas de fuerza bruta se bloquean y cada rol solo accede a lo suyo |
| **0.5 Frontend base** ✅ | Tema, layout, login y componentes base según las instrucciones de diseño que se den en ese paso | Se ve el layout en móvil y web |
| **0.6 Integración** ✅ | Front ↔ back, guardas de ruta por rol, validación en el front, nginx con proxy de confianza (IP real), todo en `docker compose`. **Pruebas:** login/refresh/logout desde el navegador, cookie solo por mismo origen, roles en rutas y menú, cambio de contraseña obligatorio, bloqueo por IP real detrás de nginx, errores mostrados en lenguaje amigable | `docker compose up` → login → menú según rol → cerrar sesión |

**Convenciones de la API (todo el proyecto):** REST con recursos en plural (`/api/v1/hospitales`), verbos HTTP y códigos de estado correctos, JSON con `camelCase`, un solo formato de error, paginación/filtros/orden con los mismos parámetros en todos los listados, validación de entrada en cada endpoint, versión en la ruta.

### Fase 1 · Catálogos núcleo
**Entrega:** dar de alta todo lo que una cirugía necesita.
- **Equipo** (modelo maestro + unidades): clave, nombre, marca/modelo, costo por hora (tarifa con vigencia), ciclo máx., purga sí/no y cada cuántos. Unidades con clave `EQ-001-01`, serie, fecha de alta, usos previos.
- **Insumo**: clave, descripción, unidad, precio, código de barras, stock actual (entrada simple) y mínimo.
- **Hospital** (datos básicos + contactos), **Doctor**, **Técnico**, **Procedimiento** (solo nombre).
- Claves únicas entre equipo e insumo (RF-08).
- Mismo patrón CRUD en todos: listado con buscador, alta, edición, baja lógica.

**Listo cuando:** puedo capturar el catálogo real del Excel a mano y consultarlo desde el celular.

| Paso | Qué entrega | Listo cuando |
|---|---|---|
| **1.0 IDs UUID y transacciones** ✅ | Esquema de Fase 0 rehecho con UUID v7; escrituras de varios pasos (login, refresh, logout, cambio de contraseña, Sysadmin inicial) en una sola transacción | Login, refresh, integración y seguridad pasan igual que antes |
| **1.1 Patrón CRUD + Procedimientos** ✅ | Piezas reutilizables en back (listado con `q`/`sort`, baja lógica, 409 por duplicado) y front (listado tabla/tarjetas, formulario por secciones, acciones, diálogo de baja, estado vacío, toast). Primer catálogo: Procedimientos | Alta, edición, baja, búsqueda y orden en móvil y web; Doctor recibe 403 |
| **1.2 Técnicos** ✅ | Puesto como lista cerrada (CHECK + oneof + select), teléfono opcional validado igual en back y front | Se hace copiando el patrón |
| **1.3 Doctores** ✅ | Número opcional y único para siempre, correo único entre activos, dirección normalizada con catálogo de C.P. (SEPOMEX/INEGI) | Matriz D1–D18 en `docs/11` |
| **1.4 Hospitales + contactos** ✅ | Clave con prefijo del estado de su colonia, dirección obligatoria del catálogo, coordenadas y liga https, contactos en la misma transacción | Matriz H1–H13 en `docs/11` |
| **1.5 Insumos + `CLAVE_CATALOGO`** ✅ | Clave reservada en `CLAVE_CATALOGO` (se corrige en cascada, no se reusa tras la baja), código de barras único, precio con 2 decimales, stock ≥ 0, chip "Bajo mínimo" | Matriz I1–I12 en `docs/11` |
| **1.6 Equipo + tarifa** ✅ | Modelo con clave generada (EQ-001), purga opcional, tarifa con vigencia e historial; claves de hospital e insumo también generadas | Matriz E1–E12 en `docs/11` |
| 1.7 Unidades de equipo | Unidades en el alta y en un equipo existente, clave `EQ-001-01` automática | — |
| 1.8 Cierre | Captura real del Excel en pruebas, revisión móvil, pruebas y docs | Criterio de la fase |

**Deuda técnica (campos del mock que no entran en la Fase 1):**

| Campo del mock | Dónde | Entra en |
|---|---|---|
| Duración típica (horas) | Procedimiento | Fase 2 (precarga el horario al programar) |
| Plantilla base | Procedimiento | Fase 9 (plantillas y kits) |
| Descripción | Procedimiento | Fase 9 (junto con la plantilla) |
| Correo | Técnico | Fase 6 (contacto del personal en el seguimiento del doctor) |
| Hospitales donde opera | Doctor | Fase 2 (sugerir hospital al programar) |
| Acceso al portal | Doctor | Fase 6 (usuarios y portal Doctor) |
| Foto | Hospital | Fase 9 (ficha de hospital completa) |
| Proveedor | Equipo e insumo | Fase 7 (CRUD de proveedores) |
| Accesorios | Equipo | Fase 9 (kits) |
| Notas | Equipo | Fase 4 (ficha de equipo) |
| Tipo de material y notas | Insumo | Fase 7 (ficha de insumo e inventario) |
| Escanear series | Unidades de equipo | Fase 8 (lector de códigos) |

### Fase 2 · Programar cirugías + disponibilidad
**Entrega:** el dolor principal resuelto — saber qué está libre y apartarlo.
- Alta de cirugía: folio `CIR-###-AA`, procedimiento, hospital, doctor, personal, fecha/horas, **entrega y devolución**.
- Material: se eligen modelos de equipo e insumos con cantidad; el sistema **asigna una unidad libre** en la ventana.
- **Aviso de choque** (sin unidad libre o técnico ocupado) que pide confirmación.
- **Cotización inicial** en vivo (horas × costo por hora + insumos × precio), guardada al programar.
- Listado de cirugías con filtros + **ficha de cirugía** (datos, material, cotización inicial, estado).
- Reprogramar (recalcula cotización) y cancelar (libera unidades).
- **Tablero de disponibilidad**: cada unidad con su estado (Disponible, Reservado, En tránsito, En hospital ‹cuál›, Mantenimiento, Fuera) y su próxima reserva; consulta "¿qué está libre de A a B?".

**Listo cuando:** programo dos cirugías que se traslapan y el sistema me avisa; el tablero muestra qué unidad está reservada para cuál.

### Fase 3 · Entregas y devoluciones
**Entrega:** el flujo completo de la cirugía con sus 5 estados.
- Vista **Entregas y devoluciones**: lista del día → al elegir una, su revisión a pantalla completa → regreso a la lista.
- Check por elemento (unidad, insumo, persona) **manual** en esta fase; se captura quién recibe/entrega.
- Devolución: piezas retornadas por insumo (consumido calculado) y cada unidad regresada; no confirma sin revisar todo.
- Estados Programada → En entrega → En hospital → Devuelta → Cerrada; la ubicación de cada unidad cambia con ellos.
- Retornados vuelven al stock; **cotización final** = inicial − retornados.
- La ficha de cirugía muestra el botón según el momento ("Iniciar entrega" / "Iniciar devolución" / resumen).
- Cada cirugía cerrada registra un **uso** por unidad.

**Listo cuando:** llevo una cirugía de Programada a Cerrada desde el celular y veo la cotización final y el stock ajustado.

### Fase 4 · Mantenimiento y semáforo
**Entrega:** el equipo se cuida solo con avisos.
- Ficha de equipo/unidad: ubicación, reservas, historial de uso, historial de mantenimiento.
- Registrar mantenimiento (preventivo, correctivo, purga) con fecha, uso n.º, responsable, costo y notas.
- Usos totales, usos del ciclo y **semáforo** (verde/amarillo/rojo según `ciclo_max`); contador de **purga** con "Purgar ahora".
- Aviso al asignar una unidad en rojo o con purga pendiente.
- Mandar una unidad a mantenimiento la saca de la disponibilidad.

**Listo cuando:** una unidad que llega a 41 usos sale en rojo y el sistema avisa al querer asignarla.

> **Al terminar la Fase 4 el sistema ya reemplaza al Excel en la operación diaria.**

---

## Bloque B · Visibilidad y usuarios

### Fase 5 · Dashboard y agenda
- Dashboard: equipo por estado, cirugías de la semana, semáforos amarillo/rojo, purgas pendientes, insumos bajo mínimo.
- Agenda semanal por hora, más día y mes; filtros por hospital, doctor y equipo.

### Fase 6 · Usuarios y portal Doctor
- Sysadmin: listado de usuarios, alta con rol, vinculación a doctor, baja que conserva historial.
- Doctor: su agenda (semana/mes), detalle de su cirugía con material y cotización, y **"Ver seguimiento completo"** (equipo e insumos con estado, avance con hora y quién, contacto del personal, cotización final). Sin costos internos ni unidades.
- Invitación por correo: requiere un servidor SMTP; mientras no exista, el Sysadmin da una contraseña temporal.

---

## Bloque C · Agregados

### Fase 7 · Inventario, gastos y proveedores
- Movimientos de insumo completos (entrada, salida, consumo, retorno, ajuste) y stock calculado desde ellos; ficha de insumo.
- Gastos por cirugía con categoría y switch **cobrable** (suma a la cotización final); ganancia por cirugía.
- CRUD de proveedores (opcional en altas de equipo e insumo).

### Fase 8 · Lector de códigos y QR
- Escaneo con cámara (móvil) y lector USB (web) en el checklist; rechaza códigos ajenos y avisa repetidos; guarda manual vs. escaneo.
- Entradas de insumo por código de barras; alta de unidades leyendo la serie.
- QR imprimible por unidad.

### Fase 9 · Detalles y pulido
- Cotización a PDF con formato de la ficha (logo, folio, hospital).
- Línea de tiempo (Gantt) por equipo.
- Plantillas por procedimiento y kits (precargan material al programar).
- Historial de cambios de la cirugía y auditoría (RNF-07).
- Ficha de hospital completa (ubicación, foto, equipo presente).
- Reportes por periodo: uso por equipo, cirugías por hospital/doctor, cotización, gastos y ganancia.
- **Búsqueda eficiente sin importar acentos ni mayúsculas** (`cirugia` encuentra "Cirugía", `pena` encuentra "Peña"): columnas de búsqueda con colación `UNICODE_CI_AI` de Firebird, índices sobre ellas y búsqueda por prefijo; si un catálogo crece mucho, evaluar el índice de texto completo.
- Páginas de error (404/403/500/503) y revisión final de estados vacíos.

---

## Notas
- RF-82 lista "material preparado" como paso del seguimiento del doctor, pero RF-40a prohíbe ese estado. En la Fase 6 se usan solo los 5 estados.
- La regla de redondeo de horas (fracciones) se fija en la Fase 2, al implementar la cotización.
