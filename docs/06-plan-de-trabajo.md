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
| **0.4 Usuarios y seguridad** | Contraseñas con hash y sal (argon2id o bcrypt), límite de intentos por usuario y por IP en una ventana de tiempo, bloqueo temporal, roles en BD, validación de datos y de rol en el back, tokens de acceso cortos + refresh revocable, auditoría de accesos. Seed de un Sysadmin | Las pruebas de fuerza bruta se bloquean y cada rol solo accede a lo suyo |
| **0.5 Frontend base** | Tema, layout, login y componentes base según las instrucciones de diseño que se den en ese paso | Se ve el layout en móvil y web |
| **0.6 Integración** | Front ↔ back, guardas de ruta por rol, validación en el front, todo en `docker compose` | `docker compose up` → login → menú según rol → cerrar sesión |

**Convenciones de la API (todo el proyecto):** REST con recursos en plural (`/api/v1/hospitales`), verbos HTTP y códigos de estado correctos, JSON con `camelCase`, un solo formato de error, paginación/filtros/orden con los mismos parámetros en todos los listados, validación de entrada en cada endpoint, versión en la ruta.

### Fase 1 · Catálogos núcleo
**Entrega:** dar de alta todo lo que una cirugía necesita.
- **Equipo** (modelo maestro + unidades): clave, nombre, marca/modelo, costo por hora (tarifa con vigencia), ciclo máx., purga sí/no y cada cuántos. Unidades con clave `EQ-001-01`, serie, fecha de alta, usos previos.
- **Insumo**: clave, descripción, unidad, precio, código de barras, stock actual (entrada simple) y mínimo.
- **Hospital** (datos básicos + contactos), **Doctor**, **Técnico**, **Procedimiento** (solo nombre).
- Claves únicas entre equipo e insumo (RF-08).
- Mismo patrón CRUD en todos: listado con buscador, alta, edición, baja lógica.

**Listo cuando:** puedo capturar el catálogo real del Excel a mano y consultarlo desde el celular.

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
- Páginas de error (404/403/500/503) y revisión final de estados vacíos.

---

## Notas
- RF-82 lista "material preparado" como paso del seguimiento del doctor, pero RF-40a prohíbe ese estado. En la Fase 6 se usan solo los 5 estados.
- La regla de redondeo de horas (fracciones) se fija en la Fase 2, al implementar la cotización.
