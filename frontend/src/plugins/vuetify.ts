/**
 * Configuración global de Vuetify: ÚNICO lugar para colores, tema, estilos
 * por defecto de componentes (defaults) y componentes virtuales (aliases).
 *
 * Regla del proyecto: nada de colores ni estilos sueltos en los archivos .vue.
 * Si un componente se ve distinto en toda la app, se cambia aquí.
 */
import { createVuetify } from 'vuetify'
import { VAlert, VBtn, VTextField } from 'vuetify/components'
import { es } from 'vuetify/locale'
import { aliases, lucide } from './icons'
import 'vuetify/styles'

// Paleta SALUS (docs/05-guia-diseno.md). Los nombres de color se usan con
// la prop `color` o las clases text-* / bg-* de Vuetify.
const claro = {
  'background': '#FFFFFF', // sin color de fondo en modo claro
  'surface': '#FFFFFF',
  'surface-bright': '#F7F8FA',
  'surface-light': '#F3F4F6',
  'surface-variant': '#4A5160',
  'on-surface-variant': '#F3F4F6',
  'on-background': '#1A1D23',
  'on-surface': '#1A1D23',
  'primary': '#017A96', // primary-strong: botones y enlaces (contraste AA)
  'primary-darken-1': '#015F75',
  'brand': '#0195B7', // primary de marca: gráficas, rellenos
  'secondary': '#4FB5AB',
  'secondary-darken-1': '#1F6F67',
  'info': '#015F75',
  'success': '#1F6F67',
  'warning': '#9A5A00',
  'error': '#B42323',
  'transito': '#4A2885',
  'text-secondary': '#4A5160',
  'text-muted': '#6B7280',
}

const oscuro = {
  'background': '#0F1417',
  'surface': '#171D21',
  'surface-bright': '#1E262B',
  'surface-light': '#1E262B',
  'surface-variant': '#B3BEC5',
  'on-surface-variant': '#1E262B',
  'on-background': '#E8EDF0',
  'on-surface': '#E8EDF0',
  'primary': '#4CC3DD',
  'on-primary': '#06232B',
  'primary-darken-1': '#7FD6EA',
  'brand': '#0195B7',
  'secondary': '#6CCBC1',
  'on-secondary': '#06232B',
  'secondary-darken-1': '#8EDDD4',
  'info': '#7FD6EA',
  'success': '#8EDDD4',
  'warning': '#F2C27A',
  'error': '#F4A3A3',
  'transito': '#C9B2F2',
  'text-secondary': '#B3BEC5',
  'text-muted': '#8A979F',
}

export default createVuetify({
  locale: { locale: 'es', messages: { es } },

  icons: { defaultSet: 'lucide', aliases, sets: { lucide } },

  theme: {
    defaultTheme: 'salus',
    themes: {
      salus: {
        dark: false,
        colors: claro,
        variables: { 'border-color': '#1A1D23', 'border-opacity': 0.1 },
      },
      salusOscuro: {
        dark: true,
        colors: oscuro,
        variables: { 'border-color': '#E8EDF0', 'border-opacity': 0.12 },
      },
    },
  },

  // Componentes virtuales: variantes que se repiten en toda la app.
  aliases: {
    BtnPrimario: VBtn, // acción principal: sólido con sombra leve
    BtnIcono: VBtn, // solo icono, sin borde ni fondo
    CampoBusqueda: VTextField, // buscador junto a tablas y en el header
    AlertaFormulario: VAlert, // alerta arriba de los campos de un formulario
  },

  // Estilo por defecto de cada componente. Aquí viven las decisiones de diseño.
  defaults: {
    // Densidad comfortable en todos los componentes que la admiten.
    global: { ripple: true, density: 'comfortable' },

    // Botón secundario por defecto: solo texto.
    VBtn: { variant: 'text', color: 'primary', rounded: 'lg' },
    BtnPrimario: { variant: 'elevated', color: 'primary', elevation: 2, rounded: 'lg' },
    BtnIcono: { variant: 'text', icon: true, color: 'text-secondary', rounded: 'lg' },

    // Campos: siempre outlined con la etiqueta en el borde.
    VTextField: { variant: 'outlined', color: 'primary', rounded: 'lg' },
    VTextarea: { variant: 'outlined', color: 'primary', rounded: 'lg' },
    VSelect: { variant: 'outlined', color: 'primary', rounded: 'lg' },
    VAutocomplete: { variant: 'outlined', color: 'primary', rounded: 'lg' },
    VCombobox: { variant: 'outlined', color: 'primary', rounded: 'lg' },
    CampoBusqueda: {
      variant: 'outlined', rounded: 'lg', hideDetails: true, singleLine: true,
      prependInnerIcon: '$buscar', placeholder: 'Buscar', color: 'primary',
    },

    // Grid: los espacios entre bloques salen de aquí (ver docs/09 › Espaciado).
    VRow: { density: 'comfortable' },

    VCard: { elevation: 0, border: true, rounded: 'lg' },
    VAlert: { variant: 'tonal', color: 'info', rounded: 'lg', border: false },
    // Margen inferior = espacio de mensajes que reservan los campos, para que
    // la alerta no quede pegada al primer campo (regla de formularios, docs/09).
    AlertaFormulario: { variant: 'tonal', color: 'error', icon: '$error', rounded: 'lg', class: 'mb-5' },
    VChip: { variant: 'tonal', rounded: 'lg', size: 'small' },
    VSnackbar: { color: 'surface', elevation: 3, rounded: 'lg', location: 'top right', timeout: 4000 },
    VDialog: { maxWidth: 560 },
    VAvatar: { rounded: 'lg' },
    VMenu: { offset: 8 },

    // Navegación: drawer blanco con borde, lista anidada.
    VNavigationDrawer: { elevation: 0, border: 'e', width: 264 },
    VAppBar: { flat: true, border: 'b', color: 'surface', density: 'comfortable' },
    VList: { slim: true, color: 'primary' },
    VListItem: { rounded: 'lg' },
    VListSubheader: { class: 'text-label-small text-uppercase text-text-muted' },
  },
})
