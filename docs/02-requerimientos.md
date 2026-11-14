# 02 · Requerimientos — Plataforma SALUS

## 1. Objetivo

SALUS renta equipo médico y surte insumos a varios hospitales para cirugías. La plataforma reemplaza el Excel y sirve para:

1. **Ver en todo momento dónde está cada equipo**: disponible, reservado, en camino, en un hospital o en mantenimiento.
2. **Programar cirugías** con su hospital, doctor, personal, equipo e insumos, y con la hora de entrega y de devolución.
3. **Cotizar bien cada cirugía** con el material llevado y lo que se retorna.
4. **Llevar el mantenimiento del equipo** a tiempo con un historial de uso, un historial de mantenimiento y un semáforo.
5. **Controlar el inventario de insumos** con sus movimientos.

### Por qué
- No se sabe qué equipos están libres ni cuáles están en uso en qué hospital.
- El Excel une los datos con texto y comas, y por eso los conteos de uso y de stock no son confiables.
- La cotización se captura a mano y no hay tarifas.
- El mantenimiento solo es un número: no tiene fechas ni historial.

### Fuera de alcance
- Facturación, pagos y contabilidad. Lo económico se limita a una **cotización correcta**, más los gastos y la ganancia por cirugía (como en el Excel).
- Los roles User SALUS y User Doc (en espera de definición).

## 2. Conceptos

| Concepto | Definición |
|---|---|
| **Equipo** | Máquina reutilizable (torre de laparoscopia, electrocirugía, morcelador, láser). Tiene historial de uso, historial de mantenimiento, semáforo de ciclo y, si aplica, purga. |
| **Insumo** | Consumible desechable de un solo uso (sondas, bolsa recolectora, sistema de irrigación, funda, fibra). Se controla por stock. |
| **Hospital** | Destino de la renta. Son varios, en varios estados. |
| **Cirugía** | Evento en un hospital, con un doctor encargado, un procedimiento, personal, material y una ventana de renta. |
| **Ventana de renta** | De la hora de entrega a la hora de devolución del material. Define la disponibilidad y la cotización: **el equipo se cobra por hora**. |
| **Plantilla** | Lista configurable de equipo e insumos predefinidos para un procedimiento o kit. Se precarga al programar y se puede ajustar. Hoy no tiene contenido. |
| **Retorno de insumos** | Insumos llevados que no se usaron. Regresan al stock y no se cotizan. |

## 3. Roles y permisos

| Funcionalidad | Sysadmin | Admin | Doctor |
|---|:-:|:-:|:-:|
| Usuarios, roles y configuración | ✔ | — | — |
| Catálogos (equipo, insumos, doctores, hospitales, técnicos, procedimientos/plantillas, tarifas) | ✔ | ✔ | — |
| Programar, reprogramar y cancelar cirugías | ✔ | ✔ | — |
| Checklist de entrega y devolución, retorno de insumos | ✔ | ✔ | — |
| Mantenimiento y purga de equipo | ✔ | ✔ | — |
| Gastos operativos | ✔ | ✔ | — |
| Tablero de disponibilidad, agenda general y reportes | ✔ | ✔ | — |
| Su propia agenda, el material de sus cirugías y su cotización | — | — | ✔ (solo lectura) |

User SALUS y User Doc (Excel, hoja Roles): **por definir**, fuera del MVP.

## 4. Módulos y requerimientos funcionales

### 4.1 Catálogos
- **RF-01 Equipo = modelo maestro + unidades.** El **modelo** guarda clave, nombre, marca/modelo, proveedor (opcional), **costo de renta por hora**, límite de ciclo (50 por defecto) y purga sí/no con cada cuántos usos (5 por defecto). Cada **unidad** del modelo (p. ej. 3 lámparas iguales = 1 modelo con 3 unidades) tiene clave propia (`EQ-001-01`), número de serie, código QR, fecha de alta y usos previos, y lleva **su propio** uso, ciclo, purga, mantenimiento, estado y ubicación. Se pueden agregar unidades a un modelo existente.
- **RF-02 Insumo**: clave, descripción, unidad, precio, código de barras, proveedor (opcional), stock actual y stock mínimo.
- **RF-09 Proveedor** (CRUD básico): nombre, qué surte (equipo, insumos o ambos), contacto, teléfono y correo. El proveedor **no es obligatorio** al dar de alta equipo o insumo.
- **RF-03 Doctor**: número, nombre, teléfono personal, teléfono de consultorio, ciudad de residencia y correo. Puede tener un usuario con rol Doctor.
- **RF-04 Hospital**: clave (con prefijo de estado, p. ej. `GTO-001`), nombre, dirección, estado, municipio, C.P., teléfono, ubicación (latitud/longitud y liga de mapa), foto y uno o varios contactos.
- **RF-05 Técnico / personal**: nombre, puesto (ayudante, enfermería, anestesiología…) y teléfono.
- **RF-06 Procedimiento**: nombre (Holep, Artroscopia/Morcelación, Cirugía General, Traumatología…) y su plantilla opcional de equipo e insumos con cantidades.
- **RF-07 Kit**: clave y lista de equipo e insumos. Funciona como una plantilla que se puede reutilizar.
- **RF-08** Las claves son únicas y no se repiten entre equipo e insumo. Un elemento es equipo **o** insumo, nunca los dos.

### 4.2 Programación de cirugías
- **RF-10** Crear una cirugía con folio automático `CIR-###-AA`, procedimiento, hospital, doctor encargado, ayudantes/personal, fecha, hora de inicio y de fin, y **hora de entrega y de devolución**.
- **RF-11** Al elegir el procedimiento se muestra su plantilla (equipo e insumos con cantidades) y se puede ajustar. El formulario se divide en tres bloques: procedimiento, lugar y equipo humano, y fecha y horario (con la ventana de renta visible).
- **RF-12** Por cada modelo de equipo de la plantilla el sistema **asigna una unidad libre** en la ventana. *(Se eliminó "otros conceptos": lo extra se registra como gasto cobrable, ver RF-72.)*
- **RF-13** **Aviso de choque**: si ninguna unidad del modelo está libre (o el técnico ya está asignado) en una ventana que se traslapa, el sistema avisa y no deja guardar sin confirmar.
- **RF-14** Estados: Programada → En entrega → En hospital → Devuelta → Cerrada, más Cancelada. Reprogramar o cambiar de hospital guarda un historial de cambios.
- **RF-15** **Agenda semanal** por hora (como en Formatos), filtrable por hospital, doctor y equipo, más vistas de día y de mes.

### 4.3 Disponibilidad y ubicación del equipo
- **RF-20** Tablero con el estado de cada equipo: **Disponible · Reservado · En tránsito · En hospital (cuál) · En mantenimiento/purga · Fuera de servicio**, con su próxima reserva.
- **RF-21** Vista de línea de tiempo (Gantt) por equipo con sus ventanas de renta.
- **RF-22** Consultar qué equipos están libres en un rango de fecha y hora.
- **RF-23** El estado cambia con el flujo de la cirugía (entrega y devolución) y con el mantenimiento.

### 4.4 Cotización
- **RF-30** La cotización de la cirugía se calcula con: equipo (horas de renta × costo por hora) + insumos (cantidad × precio) + gastos marcados como cobrables.
- **RF-31** **El equipo se cobra solo por hora.** Cada equipo tiene un **costo de renta por hora** (con fecha de vigencia). Su importe = horas de la ventana entrega → devolución × costo por hora. Si la cirugía se reprograma y la ventana cambia, la cotización se recalcula.
- **RF-32** Hay dos momentos: la **cotización inicial** se calcula en vivo al programar (ventana × costo por hora, insumos planeados, gastos cobrables) y queda visible en la ficha y para el doctor **desde ese momento, antes de entregar el equipo**; la **cotización final** se calcula al confirmar la devolución y solo descuenta los insumos retornados.
- **RF-33** La cotización se puede imprimir o exportar a PDF con el formato de la ficha (logo, folio, hospital).

### 4.5 Entrega, devolución y checklist
- **RF-40** Vista propia **"Entregas y devoluciones"**, **guiada en web y en móvil**, con los movimientos del día. La vista principal es la lista del día; al elegir una entrega o devolución se abre **su revisión a pantalla completa** (sin panel lateral ni pestañas para saltar a otro movimiento) y se vuelve a la lista al terminar. Cada elemento (unidad de equipo, insumo, persona) se confirma con su propio check o con el lector; no hay captura libre en tabla. Un código que no pertenece al movimiento se rechaza con un aviso (evita errores de dedo) y uno repetido avisa que ya estaba revisado. Se captura quién recibe o entrega en el hospital.
- **RF-40a** La cirugía tiene solo 5 estados, iguales para Admin y Doctor: Programada → En entrega → En hospital → Devuelta → Cerrada. No hay estados intermedios (por ejemplo, "material preparado").
- **RF-40b** La **ficha de la cirugía enlaza a esta vista según el momento**: antes de la entrega muestra "Es momento de entregar · Iniciar entrega"; con el material en el hospital muestra "Iniciar devolución"; después muestra el resumen confirmado (hora, quién, cantidades). El botón principal de la ficha sigue la misma secuencia. La ficha no tiene checklist propio.
- **RF-41** La devolución pide, **por cada insumo**, cuántas piezas retornaron (el consumido se calcula) y que cada unidad de equipo regresó. No se confirma hasta revisar todo.
- **RF-42** Los insumos retornados vuelven al stock y se restan de la cotización final.
- **RF-43** **Lector de códigos de barras y QR**: con la cámara del celular (o un lector USB en web) para (a) marcar el checklist de entrega y devolución, (b) registrar entradas de insumos por código de barras y (c) dar de alta unidades de equipo leyendo su número de serie. Cada unidad tiene su QR imprimible. El registro guarda si fue manual o por escaneo.

### 4.6 Equipo: uso, mantenimiento y semáforo
- **RF-50** **Historial de uso**: cada cirugía en que se usó el equipo, con fecha, hospital, doctor y horas.
- **RF-51** **Historial de mantenimiento**: fecha, tipo (preventivo, correctivo, purga), uso n.º en que se hizo, proveedor/responsable, costo y notas.
- **RF-52** Usos totales, usos del ciclo actual (desde el último mantenimiento) y semáforo: **verde 0–20 · amarillo 21–40 · rojo ≥ 41**, con un máximo de 50 que se puede configurar por equipo.
- **RF-53** Purga para los equipos que la llevan (morcelador cada 5 usos): contador, alerta "Purgar ahora" y registro de purga realizada.
- **RF-54** Si un equipo llega al máximo del ciclo o le toca purga, el sistema avisa al asignarlo a una cirugía.

### 4.7 Inventario de insumos
- **RF-60** Movimientos: entrada, salida a cirugía, consumo, retorno y ajuste. Cada uno guarda fecha, usuario y cirugía.
- **RF-61** El stock se calcula a partir de los movimientos, con una alerta cuando baja del mínimo.

### 4.8 Gastos operativos
- **RF-70** Gastos por cirugía: concepto/proveedor, categoría (Personal/Honorarios, Transporte/Logística, Infraestructura, Mantenimiento…), monto y fecha.
- **RF-71** Ganancia por cirugía = cotización final − gastos ligados (como en el Excel).
- **RF-72** Cada gasto puede marcarse como **cobrable**: además de restar en la ganancia, se suma a la cotización final del cliente (p. ej. un traslado especial). Por defecto es interno.

### 4.9 Portal Doctor
- **RF-80** El doctor ve su agenda personal (semana y mes).
- **RF-81** En cada una de sus cirugías ve el hospital, el personal, el material asignado y la cotización. Solo lectura.
- **RF-82** **Seguimiento del doctor** (parte de su agenda: al elegir una cirugía se ve el detalle base y un botón "Ver seguimiento completo" abre esta vista): incluye el **listado del equipo médico y de los insumos** con cantidades y estado (entregado, en quirófano, devuelto), el avance de cada cirugía (programada, material preparado, en camino, entregado, cirugía, recolección, cotización final), con hora y quién; el estado de **su** material (confirmado, entregado o "SALUS lo está resolviendo"), avisos de cambios, contacto del personal y cotización final al cerrar. No muestra costos internos, unidades ni datos de administración.

### 4.10 Usuarios (Sysadmin)
- **RF-85** Listado de usuarios con su rol, estado y último acceso, con alta, edición y baja (la baja se confirma y conserva el historial).
- **RF-86** Al dar de alta un usuario se elige el rol (Sysadmin, Admin o Doctor). Si es Doctor, se vincula a un registro del catálogo de doctores. El acceso se activa por invitación al correo.

### 4.11 Fichas de detalle
- **RF-87** Cada registro principal tiene una ficha de consulta: **cirugía** (estado, información, hospital, equipo, insumos con retorno, checklist de entrega y devolución, cotización inicial y final, gastos, resumen económico e historial de cambios), **equipo** (ubicación, reservas, ciclo, purga, ficha técnica, historial de uso y de mantenimiento), **insumo** (datos, stock contra mínimo y movimientos) y **hospital** (datos, contactos, equipo presente y cirugías).
- **RF-88** Toda entidad de catálogo tiene alta, edición y baja con el mismo patrón (ver `05-guia-diseno.md`, sección CRUD).

### 4.12 Reportes / dashboard
- **RF-90** Dashboard con equipo por estado, cirugías de la semana, semáforos en amarillo o rojo, purgas pendientes e insumos bajo mínimo.
- **RF-91** Reportes de uso por equipo, cirugías por hospital y por doctor, y cotización, gastos y ganancia por periodo.

## 5. Requerimientos no funcionales

| ID | Requerimiento |
|---|---|
| RNF-01 | Monorepo con `frontend/` (Vue 3 + Vuetify + TypeScript + ESLint) y `backend/` (Go Fiber v3). |
| RNF-02 | Paquetes del front con **pnpm**; npm no está permitido. |
| RNF-03 | Firebird en Docker (`firebirdsql/firebird`) con un **alias** de base de datos (`databases.conf`). La ruta real no va en la cadena de conexión y las credenciales van en variables de entorno. |
| RNF-04 | Todo se levanta con `docker compose` (front, back y BD). |
| RNF-05 | API REST documentada con **Scalar** (sin Swagger UI). |
| RNF-06 | Autenticación con usuario y contraseña, y autorización por rol en el back. |
| RNF-07 | Auditoría: quién creó o cambió cirugías, movimientos de inventario y mantenimientos. |
| RNF-08 | Zona horaria de México y moneda MXN. |
| RNF-09 | Interfaz responsiva: el checklist se debe poder usar desde el celular en el hospital. |

## 6. MVP y fases

| Fase | Contenido |
|---|---|
| **MVP** | Catálogos (RF-01…08) · Programación y agenda (RF-10…15) · Tablero de disponibilidad (RF-20…23) · Cotización por hora de renta (RF-30…32) · Checklist y retorno (RF-40…42) · Uso, mantenimiento y semáforo (RF-50…54) · Roles Sysadmin, Admin y Doctor |
| **Fase 2** | Inventario con movimientos completo (RF-60…61) · Gastos y ganancia (RF-70…71) · Portal Doctor (RF-80…81) · Exportar cotización a PDF (RF-33) |
| **Fase 3** | Reportes (RF-91) · Definir User SALUS y User Doc si aplican |

## 7. Supuestos
- El ciclo de 50 usos y la purga cada 5 son valores por defecto que se configuran por equipo.
- Los insumos retornados regresan al stock y no se cotizan.
- Un "uso" de equipo es una cirugía en la que se usó.
- Las plantillas existen como mecanismo, pero se llenan después; el sistema funciona sin ellas.
- Durante la migración, los renglones del Excel que mezclan equipo e insumo los clasifica el usuario.
