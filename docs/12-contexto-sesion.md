# 12 · Contexto para continuar en otra máquina (corte 2026-11-14)

Léelo completo al abrir una sesión nueva, junto con `CLAUDE.md`, `docs/06-plan-de-trabajo.md` y `docs/09-decisiones-tecnicas.md` (fuente de verdad técnica). La memoria de Claude de la otra máquina no viaja con el repo: todo lo importante está aquí.

## 1. Dónde vamos

- **Fase 0** completa (cimientos, usuarios y seguridad, front base, integración, entornos).
- **Fase 1 · Catálogos núcleo**, hechos y probados:
  - 1.0 IDs UUID v7 y transacciones
  - 1.1 Patrón CRUD + Procedimientos
  - 1.2 Técnicos
  - 1.3 Doctores (con dirección normalizada por catálogo SEPOMEX/INEGI)
  - 1.4 Hospitales + contactos
  - 1.5 Insumos + `CLAVE_CATALOGO`
  - 1.6 Equipo + tarifas con historial
  - Claves generadas por el sistema (hospital, insumo, equipo)
- **Siguiente: 1.7 Unidades de equipo**; después **1.8 Cierre de la Fase 1** (ver §6).
- Último estado: el usuario pidió "Revisar antes" de la 1.7 y no dejó observaciones. **Pregúntale al inicio si revisó y si hay ajustes** antes de construir la 1.7.

## 2. Arranque en una máquina nueva

```bash
git clone <repo salus-salata> && cd salus-salata
bash scripts/preparar.sh            # crea .env, .env.produccion y certificados (nuevos secretos)
bash scripts/entorno.sh pruebas up  # la 1.ª vez carga el catálogo de direcciones (~11 s)
bash scripts/entorno.sh pruebas probar   # integración: debe dar 126 correctas, 0 fallidas
bash scripts/usuarios-prueba.sh qa  # crea qa-admin para probar vistas (genera QA_ADMIN_PASSWORD en .env)
```

- Requisitos: Git, Go 1.27+, Node 24+ con pnpm 12, Docker Desktop encendido.
- App de pruebas: http://localhost:8081. API: :8080. Scalar: http://localhost:8080/docs.
- Usuario `admin`: contraseña = `INITIAL_ADMIN_PASSWORD` de `.env`; pide cambiarla al entrar. Si se olvida, `bash scripts/usuarios-prueba.sh admin`.
- DBeaver: `127.0.0.1:3050`, base `salus`, usuario `SALUS_APP`, contraseña `FIREBIRD_PASSWORD`. Los ID son `BINARY(16)`: leerlos con `UUID_TO_CHAR(ID)` o mostrar los binarios en hexadecimal.
- **Producción no se ha tocado** desde que se rehízo el esquema con UUID. Su base vieja no arrancará (las migraciones 0002/0003 cambiaron). Antes de desplegar hay que respaldarla, borrar su volumen y levantarla de nuevo, **solo cuando el usuario lo pida**.

## 3. Forma de trabajo acordada con el usuario

- Responder en **español**, **conciso**. Pasos cortos: cada paso se entrega funcionando, probado, con commit y push, y se revisa antes del siguiente.
- **Matriz de pruebas en vista ANTES de construir** cada paso, en `docs/11-pruebas-vista.md`. Al terminar se recorre en el navegador del entorno de pruebas con **`qa-admin`** (nunca con la cuenta del admin), en web 1440 px y en móvil 375 px. Si un caso cambia la lógica prevista, se le pregunta al usuario.
  - Leer resultados con `get_page_text`/JS en vez de muchas capturas (ahorra tokens). Una captura móvil por pantalla basta.
- **No inventar** reglas de negocio, estados ni contenido. Preguntar solo lo que sale de los docs. Si deduces una regla, dilo.
- **Claves e identificadores los genera el sistema**: el usuario nunca teclea claves, folios ni números internos. No pedir datos sensibles (p. ej. cédula profesional, que el usuario descartó).
- **Transacciones**: si una operación escribe en más de un paso o lee para escribir, va completa en `service.enTx`. Llaves foráneas siempre declaradas; nada de datos a medias ni huérfanos.
- **IDs**: UUID v7 (`model.ID`, `BINARY(16)`); nunca `BIGINT IDENTITY`. En desarrollo se puede rehacer el esquema si el usuario lo autoriza; si no, migraciones nuevas sin reiniciar la base.
- **Diseño**:
  - Mobile first y Vuetify first. Nada de `ga-*`, paddings dobles ni contenedores extra para separar: usar el espaciado propio de los componentes (`v-card-text`, `v-card-actions`, `v-row`). El usuario siente que me excedo con los espacios. Quiere rediseñar el patrón CRUD más adelante; por ahora lo dio por bueno.
  - Leer el artboard del mock antes de cada pantalla: https://claude.ai/artifact/3VR1enc7SDGVFu464sF8xg (`project/<Pantalla>.dc.html`).
  - Hay un MCP de Vuetify disponible para consultar la documentación de Vuetify 4.
- **Seguridad de campos**:
  - Todo campo lleva `:rules` y `maxlength`; hay `counter` global en los defaults.
  - Formatos de texto como **listas blancas** iguales en back (`handler/reglas.go`) y front (`utils/rules.ts`).
- Errores siempre amigables. pnpm siempre (npm prohibido). Scalar, nunca Swagger. Archivos UTF-8 sin BOM y LF; nunca escribir archivos con PowerShell 5.1.

## 4. Lo construido (resumen técnico)

### Backend (Go + Fiber v3)
- Patrón por catálogo:
  - migración
  - `model`
  - `repository`: lista blanca de orden con `orderBy(orden, cols, defecto, colId)`, índices únicos → campo con `integridad()`
  - `service`: `enTx` para escrituras de varios pasos
  - `handler`: `consulta()`, `idParam()`, `errorCatalogo()`
  - una línea `crud(v1, "/ruta", admin, handler)` en el router
- Roles: Sysadmin y Admin tienen acceso completo a catálogos; Doctor recibe 403.
- Baja lógica (`ACTIVO`); `DELETE` responde 204.
- **Claves generadas**:
  - Tabla `CONSECUTIVO` con `repository.Siguiente(ctx, tx, nombre)` (`UPDATE … RETURNING` dentro de la transacción; sin huecos). Nombres: `HOSPITAL:<ABR>`, `INSUMO`, `EQUIPO` y, en la 1.7, `UNIDAD:<clave equipo>`.
  - El número interno del doctor usa la secuencia `SEQ_DOCTOR_NUMERO` (`DOC-0001`; puede tener huecos).
- **`CLAVE_CATALOGO (CLAVE, TIPO)`**: equipo e insumo reservan ahí su clave (`EQUIPO`/`INSUMO`). No se libera con la baja.
- **Direcciones**:
  - Tablas `ESTADO` (32), `MUNICIPIO` y `ASENTAMIENTO` (colonia + C.P.).
  - Catálogo de Hidalgo y sus 6 vecinos (México, Querétaro, San Luis Potosí, Veracruz, Puebla, Tlaxcala) en `internal/database/catalogos/direcciones.csv.gz`, cargado al arrancar por `database.CargarDirecciones` (un estado por transacción, solo los que falten).
  - Se regenera con `go run ./cmd/catalogo-direcciones -inegi <AGEEML_utf8.csv> -sepomex <hojas.csv>`. Las hojas de SEPOMEX se exportan a CSV UTF-8 desde Excel. Los archivos fuente de INEGI y SEPOMEX **no** están en el repo.
  - Una dirección guarda calle, números y `ASENTAMIENTO_ID`. `GET /api/v1/codigos-postales/{cp}` devuelve las colonias.
- **Tarifas de equipo**: nunca se sobrescriben. Un cambio de costo agrega fila vigente desde hoy; el mismo día corrige la de hoy. "Hoy" = `model.Hoy()` (America/Mexico_City).
- Migraciones aplicadas: 0001–0015. **Una migración aplicada no se edita**: se crea otra.

### Frontend (Vue 3 + Vuetify 4)
- Piezas del patrón CRUD:
  - componentes: `ListadoCrud`, `FormularioCrud`, `SeccionFormulario`, `FormActions`, `DialogoEliminar`, `AppAvisos`, `CamposDireccion`, `CampoGenerado`
  - composables: `useListado`, `useFormulario` (con opción `enviar` para convertir texto → número), `useBaja`
  - servicio genérico: `crearRecurso`
- Menú:
  - Inventario › Equipo (Equipos y mantenimiento, Alta de equipo)
  - Inventario › Insumos (Inventario, Alta de insumo)
  - Administración › Catálogos, con pestañas Hospitales, Doctores, Técnicos y Procedimientos más sus altas
- `router.onError` recarga la app si una pestaña vieja pide archivos de un build anterior (corrige el login atorado tras un despliegue).

### Pruebas
- `scripts/prueba-integracion.sh` (126 casos):
  - restablece al admin al empezar y parte de catálogos vacíos
  - reinicia el backend entre bloques para no chocar con el límite de 30 escrituras/min por IP
  - al terminar **borra todos los usuarios** (admin vuelve a la contraseña inicial y hay que recrear `qa-admin`)
- `scripts/prueba-seguridad.sh` (32 casos, ~5 min): correrla si se toca seguridad. Después, reiniciar el backend antes de otra prueba (deja el límite por IP lleno).
- En Git Bash el `curl` de Windows altera los acentos de la línea de comandos: los cuerpos con acentos van con `--data-binary @archivo` (función `cuerpo()` del script).

## 5. Deuda técnica (campos del mock que entran después)
Está en `docs/06-plan-de-trabajo.md` › Deuda técnica:
- **Fase 2**: duración típica del procedimiento; hospitales donde opera el doctor.
- **Fase 4**: notas de equipo; estado, ciclo y mantenimiento por unidad.
- **Fase 6**: correo del técnico; acceso del doctor al portal.
- **Fase 7**: proveedor; tipo de material y notas del insumo; mosaicos, filtros, "Registrar entrada" y movimientos de Inventario.
- **Fase 8**: escanear series.
- **Fase 9**: plantillas y kits; foto y ficha completa del hospital; accesorios de equipo.

Opcional: el usuario podría querer conservar los códigos del Excel de insumos (SON-3-18) como un campo aparte, no como clave.

## 6. Cómo hacer la 1.7 Unidades de equipo (propuesta lista)

1. **Matriz primero** en `docs/11` (casos U1…). Debe cubrir:
   - alta de equipo con unidades
   - agregar unidades a un equipo existente
   - claves `EQ-001-01`, `EQ-001-02`… generadas
   - número de serie opcional y único
   - fecha de alta (hoy por defecto, no futura)
   - usos previos (entero ≥ 0)
   - quitar una unidad en la captura antes de guardar
   - baja lógica de una unidad
   - listado de equipos con número de unidades
2. **Migración** `0016_unidades.sql`, tabla `EQUIPO_UNIDAD`:
   - columnas `ID`, `EQUIPO_ID` (FK), `CLAVE` (única), `NUM_SERIE` (única, nula), `FECHA_ALTA DATE`, `USOS_PREVIOS INTEGER >= 0`, `ACTIVO` y las fechas de creación y cambio
   - **sin** estado ni hospital actual: llegan en la Fase 2
3. **Servicio**: el alta y la edición del equipo guardan modelo, tarifa y unidades nuevas en **una** transacción. Cada clave sale de `Siguiente(ctx, tx, "UNIDAD:"+claveEquipo)` con formato `%s-%02d`.
4. **API**:
   - `POST /api/v1/equipos` acepta `unidades: [{numSerie, fechaAlta, usosPrevios}]`
   - `GET /api/v1/equipos/{id}` devuelve `unidades`
   - el listado agrega `numUnidades` (subconsulta COUNT)
   - para agregar o dar de baja una unidad en un equipo existente: `POST /api/v1/equipos/{id}/unidades` y `DELETE /api/v1/equipos/{id}/unidades/{unidadId}`
5. **Front**: sección "3 · Unidades" del mock `AltaEquipo`.
   - Lista editable con botón "Agregar unidad": serie, fecha de alta, usos previos y quitar.
   - La clave se muestra con `CampoGenerado` en las ya guardadas.
   - "Escanear series" es de la Fase 8.
6. Integración, recorrido de la matriz, docs (06 ✅, 09 si hay decisión nueva, OpenAPI), commit y push.

Luego, **1.8 Cierre de la Fase 1**:
- captura real del catálogo del Excel en pruebas
- revisión móvil de todo
- `prueba-integracion.sh` y `prueba-seguridad.sh`
- docs al día
- preguntar al usuario si quiere desplegar a producción (con `respaldar` antes)

## 7. Al cerrar cada sesión
`bash scripts/entorno.sh pruebas down`, detener `pnpm dev` si se usó, actualizar este archivo si cambió el estado, commit y push.
