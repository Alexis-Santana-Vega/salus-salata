# 01 · Análisis del Excel actual

**Fuente base:** `Registro_Cirugias_Semaforo_Completo_V4_V2.xlsx` (versión más nueva).
**Versión anterior:** `Registro_Cirugias_Semaforo_Completo_V4.xlsx`. Al final se listan las diferencias.

## 1. Contexto

SALUS (Salus Salata) renta equipo médico y surte insumos a varios hospitales para cirugías. El rentador lleva el material al hospital, lo entrega a una hora y lo recoge más tarde. Hoy el control completo vive en un Excel de 9 hojas que relaciona los datos con texto separado por comas y con fórmulas de conteo.

## 2. Hojas y contenido

| # | Hoja | Qué contiene | Estado |
|---|------|--------------|--------|
| 1 | **Formatos** | Agenda semanal por hora y ficha de una cirugía (información, cotización, checklist, retorno de insumos) | Es un boceto con un ejemplo lleno |
| 2 | **Roles** | Sys Admin, Admin, Usuario SALUS, Usuario Doc | Descripción breve |
| 3 | **Registro de Cirugías** | Una fila por cirugía con personal, kit, equipos, consumibles, cobro, gastos y ganancia | 2 cirugías de ejemplo |
| 4 | **Catálogo de Equipos** | Clave, nombre, modelo/marca, veces utilizado | 11 renglones; mezcla equipo, instrumental y consumibles |
| 5 | **Catálogo de Hospitales** *(nueva en V2)* | Clave, nombre, dirección, contacto, GPS, teléfono, foto, estado, municipio, C.P. | 1 hospital completo y 1 incompleto |
| 6 | **Doctores** | Número, nombre, teléfono personal, teléfono de consultorio, ciudad, correo | Vacía (el título dice "Catálogo de equipo médico") |
| 7 | **Inventario de Insumos** | Clave, descripción, stock inicial, unidades usadas, stock actual | 8 insumos |
| 8 | **Desglose de Gastos** | Cirugía, concepto/proveedor, categoría, monto, fecha | 4 gastos |
| 9 | **Control de Mantenimiento** | Usos totales, último mantenimiento, ciclo actual, semáforo, purga | 4 equipos |

### 2.1 Formatos: agenda semanal
- Una cuadrícula de lunes a domingo con renglones por hora, de 6:00 a 24:00.
- Las cirugías se ponen como bloques que abarcan varias horas (CIR-001 de 6:00 a 9:00, CIR-002 de 11:00 a 15:00).

### 2.2 Formatos: ficha de cirugía (ejemplo de V2)
| Sección | Contenido del ejemplo |
|---|---|
| Encabezado | Logotipo · "Salus Salata" · folio `CIR-001-26` (consecutivo y año) |
| Información | Tipo: **Holep** · Doctor: **Aguayo** · Hospital: **Angeles** · Ayudante: **Rodrigo / Dylan** |
| Cotización | Equipo: Láser 3, Torre 4 · Insumos 7 · Otros: XXXX 3, YYYY 4 · **total = SUM = 21** |
| Hospital | Imagen del hospital · Liga · Contacto |
| Checklist | Equipos: Láser OK, Torre OK · Insumos: XXX OK, YYY OK · Personal |
| Retorno de insumos | Insumos llevados 7 − retornados 2 = **5** (`R22 = N16 − Q22`) |
| Nota | "Catálogo de procedimientos — Dar de alta la base para procedimiento" |

> **Hallazgo:** los renglones de la cotización mezclan equipo (Láser, Torre), insumos y conceptos sin nombre ("Otros"). En el sistema cada renglón se tendrá que clasificar como **equipo** o como **insumo**.

### 2.3 Registro de Cirugías
| ID | Fecha | Tipo de servicio | Personal | Kit | Equipos | Consumibles | Cobro | Gastos directos | Gastos adicionales | Ganancia |
|---|---|---|---|---|---|---|---|---|---|---|
| CIR-001 | 2026-08-06 | Cirugía General | Dr. Westendarp, Enf. Quirúrgica | KIT-001 | EQ-001, EQ-002 | INS-001, INS-001, INS-002 | 25,000 | 5,000 | 3,300 | 16,700 |
| CIR-002 | 2026-08-07 | Artroscopia / Morcelación | Dr. Westendarp, Anestesiólogo | KIT-002 | EQ-001, EQ-004 | INS-001, INS-003, INS-003 | 35,000 | 7,000 | 4,700 | 23,300 |

- La columna de tipo de servicio también trae "Traumatologia" como tercer valor.
- En V2 hay un borrador de lo que lleva `KIT-001` (x, y, z, a, b) debajo de la tabla.

## 3. Fórmulas y reglas actuales

| Dato | Fórmula en el Excel | Regla de negocio |
|---|---|---|
| Gastos adicionales | `SUMIF('Desglose de Gastos'!A:A, ID, D:D)` | Suma de los gastos ligados a la cirugía |
| Ganancia neta | `Cobro − Gastos directos − Gastos adicionales` | Los gastos directos se capturan a mano |
| Veces que se usó un equipo | `COUNTIF(Equipos utilizados, "*EQ-001*")` | Cuenta las cirugías cuyo texto contiene la clave |
| Unidades de insumo usadas | `SUMPRODUCT((LEN(txt) − LEN(SUBSTITUTE(txt, clave, ""))) / LEN(clave))` | Cuenta cuántas veces aparece la clave dentro del texto |
| Stock actual | `Stock inicial − Unidades usadas` | No hay entradas ni retornos |
| Usos del ciclo actual | `Usos totales − Último mantenimiento (n.º de uso)` | El ciclo se reinicia con el mantenimiento |
| Semáforo del ciclo | Formato condicional: **0–20 verde · 21–40 amarillo · ≥41 rojo** (máximo 50) | |
| Usos para purga | `MOD(usos del ciclo, 5)` | Solo para el morcelador |
| Estatus de purga | `IF(MOD(usos, 5) = 0 y usos > 0, "¡PURGAR AHORA!", "OK")` | Rojo si toca purgar |
| Total de la ficha | `SUM(N14:N18)` | Suma de los renglones de la cotización |
| Insumos cobrados | `Llevados − Retornados` | Lo que regresa no se consume |

## 4. Áreas de oportunidad

1. **No se sabe dónde está cada equipo.** No hay estado (disponible, en hospital, en mantenimiento) ni horas de entrega y devolución. Es el dolor principal: no se ve qué equipo está libre ni en qué hospital está el que se usa.
2. **No hay control de choques.** Nada impide asignar el mismo equipo a dos cirugías que se traslapan.
3. **Las relaciones son texto con comas.** Los conteos fallan con facilidad: `"*EQ-001*"` también cuenta `EQ-0010`, y un espacio o un error de dedo cambia el resultado.
4. **Equipo e insumo están mezclados.** El catálogo de equipos incluye instrumental (`INS-001`, `INS-002`), varillas (`VAR-*`) y la fibra láser (`FIB-550-001`), que además aparece en Insumos. El registro usa `INS-001/002/003` como consumibles, aunque `INS-003` no existe en ningún catálogo.
5. **Los nombres no coinciden entre hojas.** EQ-002 es "Electrocirugia" en el catálogo y "Electrocauterio" en mantenimiento. EQ-003 es "Electrocirugia Erbe VIO3" en el catálogo y "Artroscopio" en mantenimiento.
6. **Faltan catálogos.** No hay técnicos/personal (el personal es texto libre), procedimientos con su plantilla, kits ni tarifas. Doctores está vacía.
7. **No hay precios.** El cobro total se captura a mano y no hay tarifa por equipo ni precio por insumo, así que la cotización no se puede calcular.
8. **El mantenimiento no tiene bitácora.** "Último mantenimiento" es solo un número de uso: no hay fecha, tipo, proveedor ni costo. La purga no registra cuándo se hizo.
9. **El inventario es limitado.** No hay stock mínimo ni movimientos (entradas, salidas a hospital, consumos, retornos). El retorno de insumos solo existe en la ficha y no regresa al stock.
10. **Las cirugías no tienen estado ni historial.** Si una cirugía cambia de fecha o de hospital, no queda rastro.
11. **El mismo folio trae datos distintos.** En la ficha, CIR-001 es *Holep* con el Dr. Aguayo; en el Registro de Cirugías, CIR-001 es *Cirugía General* con el Dr. Westendarp. Tampoco coinciden las fechas: la agenda la pone en lunes y el registro el 06-ago (jueves). Sin una sola fuente por cirugía, los datos se contradicen.
12. **Los hospitales están incompletos.** Solo hay un registro completo, y el GPS y la foto son texto ("Latitud - Longuitud", "Pegar foto").

## 5. Cambios de V4 a V4_V2

| Cambio | Impacto |
|---|---|
| Nueva hoja **Catálogo de Hospitales** con dirección, GPS, foto, contactos, estado, municipio y C.P. Las claves llevan prefijo de estado (`GTO-001`, `SLP-001`) | Confirma que hay varios hospitales en varios estados |
| La ficha de Formatos trae un ejemplo real (Holep, Aguayo, Angeles, Rodrigo / Dylan) | Muestra cómo se usa: varios ayudantes por cirugía |
| Cotización en la ficha con renglones y total | Primer esbozo de la cotización por cirugía |
| Checklist con OK por renglón | El checklist se marca renglón por renglón |
| Retorno de insumos con resta | Lo retornado ajusta la cotización y el inventario |
| Campos Liga y Contacto en la ficha | Datos del hospital para la logística |
| Nota "Catálogo de procedimientos" | Plantilla base de material por procedimiento |
| Borrador de composición de `KIT-001` | Los kits son listas de piezas |
| Las demás hojas no cambian | — |
