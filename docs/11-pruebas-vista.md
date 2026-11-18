# 11 · Pruebas en vista (matriz por paso)

Guía completa en este archivo. Antes de construir un paso se escribe su matriz; al terminar se recorre en el **entorno de pruebas** desde el navegador (web 1440 px y móvil 375 px) y se corrige lo que falle. Si un caso revela una regla de negocio no prevista, se pregunta al usuario antes de decidir.

- Se prueba con **`qa-admin`** (rol Admin): `bash scripts/usuarios-prueba.sh qa`. Nunca con la cuenta del admin.
- `bash scripts/usuarios-prueba.sh admin` restablece al admin con `INITIAL_ADMIN_PASSWORD` de `.env`.
- `entorno.sh pruebas probar` (integración por API) borra los usuarios al terminar: después hay que volver a correr `usuarios-prueba.sh qa`.
- Resultado: ✔ pasó · ✘ falló (con la corrección hecha) · — no aplica.

## Patrón CRUD (aplica a todos los catálogos)

| # | Caso | Esperado |
|---|---|---|
| G1 | Entrar al catálogo | Skeleton → listado ordenado por nombre |
| G2 | Guardar el formulario vacío | Error en cada campo obligatorio, sin petición |
| G3 | Cancelar en el formulario | Vuelve al listado sin guardar |
| G4 | Guardar válido | Botón "Cargando" deshabilitado, toast y vuelta al listado |
| G5 | Buscar sin acentos ni mayúsculas | Encuentra el registro |
| G6 | Buscar algo inexistente | "Sin coincidencias" con "Limpiar búsqueda" |
| G7 | Ordenar por cada columna ordenable (asc/desc) | Cambia el orden |
| G8 | Más de 10 registros | Paginación funciona; cambiar filas por página |
| G9 | Editar y recargar la página | Mantiene sesión y datos |
| G10 | Abrir `/…/<id inexistente>` | 404 dentro del layout |
| G11 | Eliminar desde la tabla (web) o desde la edición (móvil) | Diálogo rojo → toast → desaparece |
| G12 | Móvil | Tarjetas con chevron; formulario en una columna; acciones fijas abajo |
| G13 | Pestaña abierta antes de un despliegue nuevo: iniciar sesión (o navegar) | La app se recarga sola y llega a la pantalla destino; la contraseña no se borra si algo falla |

G1–G12 recorridos el 2026-11-14 en Procedimientos y Técnicos: ✔ (paginación móvil y web, orden asc/desc, recarga en edición, 404 con id inexistente o mal formado, Cancelar sin guardar).

## 1.1 Procedimientos

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| P1 | Nombre con espacios de más ("  Holep  ") | Se guarda "Holep" | ✔ |
| P2 | Nombre repetido con otras mayúsculas/acentos | Error en Nombre: ya existe | ✔ |
| P3 | Nombre con `<` o `>` | Error en el campo | ✔ |
| P4 | Eliminar y volver a dar de alta el mismo nombre | Se permite | ✔ |

## 1.2 Técnicos

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| T1 | Puesto | Solo las 5 opciones de la lista | ✔ |
| T2 | Teléfono con letras | Error en el campo | ✔ |
| T3 | Sin teléfono | Se guarda; listado muestra "—" | ✔ |
| T4 | Dos técnicos con el mismo nombre | Se permite | ✔ |
| T5 | Buscar por teléfono | Lo encuentra | ✔ |
| T6 | Editar y vaciar el teléfono | Queda vacío | ✔ |

## 1.3 Doctores

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| D1 | Vacío | Error en Nombre y Correo | ✔ |
| D2 | Correo inválido | Error en Correo | ✔ |
| D3 | Teléfonos inválidos | Error en cada teléfono | ✔ |
| D4 | C.P. con letras o distinto de 5 dígitos | Error en C.P. | ✔ |
| D5 | Alta mínima (nombre + correo) | Se guarda; listado muestra "—" en lo vacío | ✔ |
| D6 | Alta completa con dirección | Al reabrir aparecen calle, números, C.P., colonia, municipio y estado | ✔ |
| D7 | C.P. de Pachuca con varias colonias | Lista sus colonias; al elegir una se llenan Pachuca de Soto e Hidalgo | ✔ |
| D8 | Alta | Se asigna solo el número interno (DOC-0001…); en edición se ve de solo lectura | ✔ |
| D9 | Correo repetido con otras mayúsculas | Error en Correo | ✔ |
| D10 | Buscar por nombre sin acentos, por número y por correo | Lo encuentra | ✔ |
| D11 | Editar: quitar número y cambiar teléfono | Se guarda así | ✔ |
| D12 | Eliminar y dar de alta otro con el mismo correo | Se permite | ✔ |
| D13 | Número interno de uno dado de baja | No se reasigna (secuencia) | ✔ |
| D14 | C.P. con una sola colonia (42000) | Se elige sola "Centro" y se llenan municipio y estado | ✔ |
| D15 | C.P. fuera del catálogo (99999) | Aviso en el C.P.; no deja elegir colonia | ✔ |
| D16 | C.P. capturado sin elegir colonia | Error "Elige la colonia" | ✔ |
| D17 | Editar y cambiar el C.P. | Se limpia la colonia y se cargan las del nuevo C.P. | ✔ |
| D18 | Borrar el C.P. al editar | Se guarda sin dirección del catálogo | ✔ |

Primer recorrido 2026-11-14: D1–D13 ✔. Con el catálogo de C.P. (dirección normalizada) se rehicieron D4, D6, D7 y se agregaron D14–D18: todos ✔ (la lista de colonias también filtra al escribir).

## Formatos de texto (aplica a todos los formularios)

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| F1 | Escribir más del máximo | El campo no deja pasar del tope y el contador lo muestra | ✔ |
| F2 | Nombre de persona con números, `{ }`, `;`, `<` o emoji | Error de formato en el campo (y 422 en la API) | ✔ |
| F3 | Nombre con acentos, ñ, punto, apóstrofo o guion ("Dra. María-José O'Neil Peña") | Se acepta | ✔ |
| F4 | Clave o número con espacios o acentos | Error de formato | ✔ |

## 1.4 Hospitales y contactos

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| H1 | Guardar vacío | Error en Nombre, Calle y C.P./Colonia (la clave no se pide) | ✔ |
| H2 | Alta | La clave se genera con el estado de la colonia y su consecutivo (HGO-001, HGO-002…) | ✔ |
| H3 | Editar | La clave se ve de solo lectura y no cambia | ✔ |
| H4 | Alta rechazada (p. ej. contacto sin nombre) | No consume el consecutivo | ✔ (API) |
| H6 | Latitud 200 o longitud -500 | Error en el campo | ✔ |
| H7 | Liga del mapa sin https | Error en el campo | ✔ |
| H8 | Alta con 2 contactos | Se guardan los dos, en el mismo orden | ✔ |
| H9 | Contacto sin nombre | Error en el nombre de ese contacto | ✔ |
| H10 | Editar: quitar un contacto y agregar otro | Quedan exactamente los de la pantalla | ✔ |
| H11 | Listado | Clave, nombre, ciudad, teléfono y contacto principal | ✔ |
| H12 | Buscar por clave, nombre o municipio | Lo encuentra | ✔ |
| H13 | Eliminar | Diálogo, toast, sale del listado | ✔ |

Recorrido 2026-11-14 (web, qa-admin): F1–F4 y H1–H13 ✔ (H5 y H13 también por API). Hallazgos corregidos al probar: las reglas de formato validaban antes de quitar espacios de los extremos ("  Holep" se rechazaba) y el CHECK de clave de hospital necesitaba escapar el guion en `SIMILAR TO`.

## 1.5 Insumos

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| I1 | Guardar vacío | Error en Descripción y Precio (la clave no se pide; Unidad trae Pieza y existencias 0) | ✔ |
| I2 | Alta | Clave generada INS-0001, INS-0002… | ✔ |
| I3 | Alta rechazada (código de barras repetido) | No consume la clave | ✔ |
| I4 | Clave que ya usa un equipo (cuando exista el paso 1.6) | Error en Clave: ya la usa un equipo | — |
| I5 | Precio negativo, con 3 decimales o con letras | Error en Precio | ✔ |
| I6 | Stock negativo o con decimales | Error en el campo | ✔ |
| I7 | Código de barras repetido | Error en Código de barras | ✔ |
| I8 | Sin código de barras | Se guarda | ✔ |
| I9 | Stock menor al mínimo | Chip "Bajo mínimo" en el listado | ✔ |
| I10 | Editar | La clave se ve de solo lectura y no cambia | ✔ |
| I11 | Buscar por clave, descripción o código de barras | Lo encuentra | ✔ |
| I12 | Ordenar por clave, descripción, precio y stock | Cambia el orden | ✔ |

Recorrido 2026-11-14 (web, qa-admin): I1–I12 ✔ (I11 e I12 también por API). En la vista se comprobó la atomicidad: un alta rechazada por código de barras repetido no dejó reservada su clave. I4 se prueba al construir equipo (1.6).

Claves generadas (2026-11-14): el usuario ya no teclea claves de hospital ni de insumo; H2–H4 e I1–I3, I10 se rehicieron.

## 1.6 Equipo (modelo) y tarifas

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| E1 | Guardar vacío | Error en Nombre, Marca/modelo y Costo por hora (Ciclo trae 50) | ✔ |
| E2 | Alta | Clave generada EQ-001, EQ-002…; tarifa vigente desde hoy | ✔ |
| E3 | Clave de equipo frente a insumo (I4) | Nunca chocan: las dos se reservan en el mismo catálogo de claves | ✔ |
| E4 | Costo por hora negativo o con 3 decimales | Error en el campo | ✔ |
| E5 | Ciclo 0 o con letras | Error en el campo | ✔ |
| E6 | Activar "Requiere purga" | Aparece "Purga cada (usos)" y es obligatorio; al apagarlo desaparece y no se guarda | ✔ |
| E7 | Editar sin cambiar el costo | No se crea tarifa nueva | ✔ |
| E8 | Editar con otro costo | Se agrega una tarifa al historial (vigente desde hoy); el costo vigente cambia | ✔ |
| E9 | Cambiar el costo dos veces el mismo día | Queda una sola tarifa de hoy con el último costo | ✔ |
| E10 | Historial de tarifas en la edición | Fecha y costo, la más reciente primero | ✔ |
| E11 | Listado | Clave, nombre, marca/modelo, costo por hora vigente, ciclo y purga | ✔ |
| E12 | Buscar por clave, nombre o marca | Lo encuentra | ✔ |

Recorrido 2026-11-14 (web, qa-admin): E1–E12 ✔ (E3, E8, E10 con varias tarifas y E12 por API). Claves generadas verificadas en vista: el alta de hospital e insumo ya no pide clave y la edición la muestra de solo lectura.

## 1.7 Unidades de equipo

| # | Caso | Esperado | Resultado |
|---|---|---|---|
| U1 | Alta de equipo con 3 unidades | Claves EQ-00N-01, -02, -03 generadas; el chip dice "3 unidades" | ✔ API |
| U2 | Alta sin unidades | Se guarda (las unidades se agregan después) | ✔ API |
| U3 | Agregar unidades a un equipo existente | Siguen el consecutivo del equipo (-04, -05…), nunca se repite una clave | ✔ API |
| U4 | Quitar una fila nueva antes de guardar | No se guarda ni consume clave | |
| U5 | Número de serie vacío | Se guarda sin serie | ✔ API |
| U6 | Serie repetida (en la misma captura o con otra unidad activa) | Error en la Serie de esa fila; no se guarda nada | ✔ API |
| U7 | Fecha de alta: hoy por defecto; una fecha futura | Error en la fecha | ✔ API |
| U8 | Usos previos negativos, con decimales o vacíos | Error en el campo | ✔ API |
| U9 | Editar la serie o los usos de una unidad guardada | Se actualiza; su clave no cambia | ✔ API |
| U10 | Quitar una unidad guardada y guardar | Baja lógica: deja de verse y su clave no se reutiliza | ✔ API |
| U11 | Listado de equipos | Columna Unidades con el número de unidades activas | ✔ API |
| U12 | Error en una unidad al guardar | Ni el modelo ni las demás unidades se guardan (atomicidad) | ✔ API |
| U13 | Móvil | Las filas de unidad se apilan sin scroll horizontal | |

Recorrido 2026-11-17 (API, `prueba-integracion.sh`, 142/142): U1–U3, U5–U12 ✔. U4 y U13 (vista y móvil) quedan para recorrer en el navegador con `qa-admin` al iniciar la siguiente sesión: en esta no se pudo iniciar sesión en el navegador.
