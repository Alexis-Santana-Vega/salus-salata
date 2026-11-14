# 05 · Guía de diseño (UI)

Reglas para los mocks y para la implementación en Vue 3 + Vuetify. La referencia visual está en el artboard **Sistema de diseño** del canvas de mocks.

## Principios
1. **Mobile first.** Cada pantalla se diseña primero para móvil (390 px) y después para web.
2. **Limpio.** En modo claro el fondo es blanco y la app no lleva color de fondo. La jerarquía sale del **peso y tamaño de la tipografía**, del uso correcto del **negro y los grises**, y de **fondos pastel** en algunos elementos.
3. **Modo oscuro.** Ahí sí hay color de fondo en la app (se define al implementar).
4. **Nunca fondos negros** en alertas, toasts ni ningún otro componente.

## Tipografía
- Familia: **Poppins** (400, 500, 600, 700).
- Escala: título de página 24/600 (móvil 20/600) · subtítulo 14/400 · título de sección 18/600 · cuerpo 14/400 · etiqueta 12/500 (en mayúsculas con tracking en cards de métricas) · número de métrica 32–36/600.

## Color

### Marca
| Token | Valor | Uso |
|---|---|---|
| primary | `#0195B7` | Rellenos, gráficas, anillo de foco, indicadores activos |
| primary-strong | `#017A96` | Fondo de botones primarios y color de enlaces/botones de texto (con blanco da 5:1; `#0195B7` da 3.5:1 y no pasa AA para texto) |
| primary-pastel / on-pastel | `#E5F4F8` / `#015F75` | Estado activo, chips informativos |
| secondary | `#4FB5AB` | Rellenos y gráficas (con blanco da 2.5:1, no se usa como fondo de texto) |
| secondary-pastel / on-pastel | `#E6F4F2` / `#1F6F67` | OK, disponible, confirmaciones |

### Texto: negro para lo principal, grises para lo demás
| Token | Valor | Uso |
|---|---|---|
| text | `#1A1D23` | Títulos, valores, celda clave de la tabla |
| text-secondary | `#4A5160` | Cuerpo y celdas normales |
| text-muted | `#6B7280` | Etiquetas, subtítulos, captions (4.9:1 sobre blanco) |
| border | `#ECEEF2` | Bordes de cards y divisores |
| outline | `#C5CBD5` | Borde de campos outlined |

### Pasteles de estado
**Regla:** sobre un fondo pastel, el texto usa el mismo tono más fuerte, nunca negro.
| Estado | Fondo | Texto fuerte | Texto secundario |
|---|---|---|---|
| Info / en hospital | `#E5F4F8` | `#015F75` | `#017A96` |
| OK / disponible | `#E6F4F2` | `#1F6F67` | `#2A857B` |
| Advertencia / mantenimiento | `#FFF4DE` | `#7A4300` | `#9A5A00` |
| Error / choque | `#FDECEC` | `#8F1717` | `#B42323` |
| En tránsito | `#F1ECFB` | `#4A2885` | `#6A43A8` |
| Neutro | `#F3F4F6` | `#4A5160` | — |

### Uso expresivo de pasteles y tipografía
- Mosaicos de resumen arriba de listas y tablas (estado de equipos, stock, gastos por cirugía): fondo pastel del estado, etiqueta 13/500 en el tono medio, número 30/600 en el tono fuerte y detalle 12/400 en el tono fuerte con 70% de opacidad.
- Eventos de agenda y tarjetas destacadas: fondo pastel con título en el tono fuerte (600) y el resto en el tono medio, nunca en negro.
- Se combinan peso (400/500/600), tamaño y opacidad del mismo tono antes que agregar colores nuevos.

## Modo oscuro
La app sí lleva color de fondo y se usa la misma paleta, ajustada.
| Token | Valor | Uso |
|---|---|---|
| bg | `#0F1417` | Fondo de la app (gris azulado, no negro puro) |
| surface | `#171D21` | Cards, sidebar, header, diálogos |
| surface-2 | `#1E262B` | Encabezado de tabla, menús, toasts |
| border / outline | `#2A343A` / `#3A464D` | Divisores / campos outlined |
| text / secondary / muted | `#E8EDF0` / `#B3BEC5` / `#8A979F` | Texto |
| primary | `#4CC3DD` | Botones (texto `#06232B`, 7.9:1), enlaces, activo |
| primary base / secondary | `#0195B7` / `#6CCBC1` | Gráficas y rellenos |
- **Pasteles → tintes translúcidos**: fondo del color del estado al 16–22% de opacidad, texto claro del mismo tono (info `#7FD6EA`, OK `#8EDDD4`, advertencia `#F2C27A`, error `#F4A3A3`, tránsito `#C9B2F2`) y números en una versión aún más clara.
- La elevación se marca con superficies más claras, no con sombras fuertes. El botón de modo oscuro del header cambia a un icono de sol.
- En Vuetify: un tema `dark` con estos tokens y `v-theme-provider` / `useTheme()`; la preferencia se guarda por usuario.

## Layout de la app (admin)
- **Sidebar** (componente): blanco con borde derecho, logo arriba a la izquierda, secciones con título en mayúsculas pequeñas y opciones con icono.
  - Una opción **sola** lleva directo a la página principal del módulo.
  - Una opción **expandible** contiene la página principal del módulo y sus **formularios de alta** (p. ej. Equipo → Equipos y mantenimiento, Alta de equipo).
  - La opción activa lleva fondo primary-pastel y texto on-pastel en 600.
- **AppHeader** (componente): blanco con borde inferior. Botón de menú (solo icono), buscador outlined con atajo "Ctrl K", botón de modo oscuro (solo icono) y menú de usuario (avatar cuadrado redondeado + nombre + rol + chevron). **Sin notificaciones.**
- **Doctor:** encabezado propio con logo, "Mi agenda", "Mi perfil", avatar cuadrado redondeado con nombre y botón de salir con icono.
- **Móvil:** encabezado con menú a la izquierda, logo centrado y avatar a la derecha; el sidebar se abre como drawer.

## Componentes

### Encabezado de página (`PageHeader`)
Título + subtítulo, y acciones a la derecha en web. En móvil, la acción principal va abajo a todo el ancho.

### Campos de formulario
Todos con **variante outlined** (`variant="outlined"`), con la etiqueta en el borde.

### Botones
| Tipo | Estilo | Vuetify |
|---|---|---|
| Primario | Fondo `primary-strong` sólido con sombra leve | `variant="elevated"` + `elevation="2"` |
| Secundario | Solo texto, sin subrayado | `variant="text"` |
| Icono | Sin borde ni fondo | `variant="text"` + `icon` + `aria-label` |

**Carga:** el botón se deshabilita, muestra un spinner y su texto cambia a **"Cargando"** (`:loading` + `:disabled`).

### Retroalimentación
- **Toast** (`v-snackbar`): fondo blanco, borde, sombra suave y acento de color a la izquierda con icono. Para confirmaciones breves.
- **Alerta** (`v-alert` tonal): fondo pastel con texto del mismo tono. Para lo que el usuario tiene que ver o resolver.

### Cards de métricas
Fondo blanco, borde sutil y sombra mínima. Etiqueta en mayúsculas gris con un icono en un cuadro gris claro a la derecha, número grande en negro 600, una línea de contexto en gris y, cuando aporte, un gráfico:
- **Dona:** partes de un total (equipos por estado), con leyenda y valores.
- **Barras:** comparar cantidades (stock por insumo, cirugías por día), con un tono de estado si aplica.
- **Progreso:** avance hacia un límite (ciclo de 50 usos).
- **Área/sparkline:** tendencia en el tiempo.
Los colores de datos son `primary`, `secondary` y ámbar `#E0A33C`; la identidad nunca va solo por color (leyenda o etiqueta).

### Tablas
- Encabezado con más peso: 13/600 en negro sobre fondo `#F7F8FA`, con las esquinas superiores redondeadas.
- Celda clave en negro 500 y el resto en `text-secondary`; números alineados a la derecha con cifras tabulares; hover de fila muy sutil.
- Pie: "Filas por página" con un select outlined compacto, el rango "1–5 de 8" y los botones primera, anterior, siguiente y última **solo con icono, sin borde**, en gris claro cuando están deshabilitados.

- **Columnas ordenables** donde aporte (clave, nombre, cantidades): flecha doble gris cuando la columna no está ordenada y flecha de color hacia arriba o abajo en la columna activa (`sortable` en `v-data-table`, con `aria-sort`).

### Búsqueda y filtros (decidido: dentro de la card de la tabla)
- El **buscador va junto a la tabla**, dentro de su card y no separado: un solo campo con icono de lupa y placeholder, sin etiqueta flotante.
- A la derecha, el botón **"Filtros"** (texto con icono y un contador de filtros activos) abre un **popover** con selects outlined y las acciones "Restablecer" (texto) y "Aplicar" (primario).
- Si hay muchos filtros: panel lateral en web. En móvil siempre es una hoja a pantalla completa.
- Los filtros activos se muestran como chips que se pueden quitar, junto a "Limpiar filtros".

### Listas de datos
- **Lista descriptiva** (detalle de un registro): filas etiqueta–valor con la etiqueta en gris a la izquierda (ancho fijo), el valor en negro a la derecha y un divisor sutil entre filas (`v-list` o `dl`). En móvil, la etiqueta va arriba del valor.
- **Card resumen**: datos apilados (etiqueta gris arriba, valor en negrita abajo) con divisores, y un pie separado con barra de progreso (p. ej. ciclo de 50 usos o purga).

### Estados vacíos (`EmptyState`)
Ilustración pastel en tonos de la marca, título en 600, descripción en gris y hasta dos acciones. Casos:
- **Sin coincidencias** (búsqueda o filtros): acción "Limpiar filtros".
- **Sin registros**: acción principal para crear y una secundaria (p. ej. importar).
- **Conexión perdida**: "Reintentar" con estado de carga. Si se pierde mientras se trabaja, aviso ámbar arriba sin bloquear la pantalla.
- **Agenda vacía** (doctor).
- Dentro de tablas se usa una versión compacta.

### Diálogos
- Encabezado limpio: título (y subtítulo opcional) a la izquierda y botón cerrar **solo icono, sin fondo**.
- Cuerpo con scroll; campos outlined en dos columnas en web.
- Acciones en el **pie** (o al final del scroll): "Cancelar" (texto) y la acción principal (primario, con estado de carga). Al guardar se cierra y aparece un toast.
- En móvil el diálogo ocupa toda la pantalla y las acciones van fijas abajo.
- El fondo detrás del diálogo es un velo gris semitransparente.

### CRUD (patrón para todos los módulos)
- **Listado:** card con buscador y filtros junto a la tabla, columnas ordenables y acciones por fila (ver, editar con icono de lápiz, eliminar con icono de bote). En móvil se usa una lista de tarjetas con chevron.
- **Ficha (detalle):** encabezado con breadcrumb, título y subtítulo, más "Editar" (texto) y la acción principal del módulo (primario). Datos en lista etiqueta–valor, cards resumen con progreso y tablas relacionadas (cirugías, movimientos, historial). En móvil las acciones van fijas abajo.
- **Alta y edición:** página propia (también accesible desde el sidebar) con secciones (título + ayuda) y campos outlined en rejilla. El componente `FormActions` va al pie, con "Cancelar" (texto) y "Guardar" (con carga y toast). Editar abre el mismo formulario con los datos cargados. En móvil, campos en una columna y acciones fijas abajo.
- **Eliminar:** diálogo de confirmación con icono en pastel rojo, texto que explica la consecuencia, "Cancelar" (texto) y **"Eliminar" sólido en rojo** (`#B42323`). Es la única excepción al color primario.
- **Acciones rápidas** sobre un registro existente (registrar mantenimiento, registrar gasto, registrar entrada de insumo): diálogo en web y pantalla completa en móvil.

### Páginas de error
Logo (esquina en web, centrado en móvil), código grande en primario, título, explicación breve y acciones: **404** no encontrada, **403** sin permiso, **500** error del servidor y **503** mantenimiento. Siempre "Volver al inicio" y, si aplica, "Reintentar" o "Cambiar de cuenta".

### Login
Usuario o correo + contraseña y un solo botón **Entrar** (en el mock hay dos, uno por rol, solo para navegar). Sin acceso con Google.
