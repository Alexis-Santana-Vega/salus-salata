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
| D8 | Número repetido | Error en Número | ✔ |
| D9 | Correo repetido con otras mayúsculas | Error en Correo | ✔ |
| D10 | Buscar por nombre sin acentos, por número y por correo | Lo encuentra | ✔ |
| D11 | Editar: quitar número y cambiar teléfono | Se guarda así | ✔ |
| D12 | Eliminar y dar de alta otro con el mismo correo | Se permite | ✔ |
| D13 | Número de uno dado de baja | No se permite (los números no se reusan) | ✔ |
| D14 | C.P. con una sola colonia (42000) | Se elige sola "Centro" y se llenan municipio y estado | ✔ |
| D15 | C.P. fuera del catálogo (99999) | Aviso en el C.P.; no deja elegir colonia | ✔ |
| D16 | C.P. capturado sin elegir colonia | Error "Elige la colonia" | ✔ |
| D17 | Editar y cambiar el C.P. | Se limpia la colonia y se cargan las del nuevo C.P. | ✔ |
| D18 | Borrar el C.P. al editar | Se guarda sin dirección del catálogo | ✔ |

Primer recorrido 2026-11-14: D1–D13 ✔. Con el catálogo de C.P. (dirección normalizada) se rehicieron D4, D6, D7 y se agregaron D14–D18: todos ✔ (la lista de colonias también filtra al escribir).
