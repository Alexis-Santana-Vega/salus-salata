/**
 * Iconos Lucide (trazo suave) como set por defecto de Vuetify.
 *
 * En la app los iconos se usan SIEMPRE por alias: <v-icon icon="$equipo" />.
 * Para cambiar un icono en toda la app basta con cambiarlo aquí.
 * Solo se importan los iconos usados (tree-shaking).
 */
import type { IconAliases, IconSet } from 'vuetify'
import {
  AlertCircle, AlertTriangle, ArrowDown, ArrowLeft, ArrowRight, ArrowUp, ArrowUpDown,
  Calendar, Check, CheckCircle2, ChevronDown, ChevronLeft, ChevronRight, ChevronsLeft,
  ChevronsRight, ChevronUp, Circle, CircleDot, Clock, Command, CornerDownLeft, Delete,
  Eye, EyeOff, File, Info, LayoutDashboard, Loader2, LogOut, Menu, Minus, Moon, Pencil,
  Pipette, Plus, Search, Square, SquareCheck, SquareMinus, Star, StarHalf, Sun, Tag,
  Trash2, Upload, User, UserPlus, Users, X, XCircle,
} from 'lucide-vue-next'
import { type Component, defineComponent, h, mergeProps } from 'vue'

// Nombre → componente. Los nombres son los alias que usa la app.
const registro: Record<string, Component> = {
  // Propios de la app
  inicio: LayoutDashboard,
  usuarios: Users,
  usuarioAlta: UserPlus,
  usuario: User,
  buscar: Search,
  modoOscuro: Moon,
  modoClaro: Sun,
  salir: LogOut,
  ver: Eye,
  ocultar: EyeOff,
  // Internos de Vuetify (botones de cerrar, checkboxes, paginación, etc.)
  collapse: ChevronUp, complete: Check, cancel: XCircle, close: X, delete: Trash2,
  clear: XCircle, success: CheckCircle2, info: Info, warning: AlertTriangle, error: AlertCircle,
  prev: ChevronLeft, next: ChevronRight, checkboxOn: SquareCheck, checkboxOff: Square,
  checkboxIndeterminate: SquareMinus, delimiter: Circle, sortAsc: ArrowUp, sortDesc: ArrowDown,
  expand: ChevronDown, menu: Menu, subgroup: ChevronDown, dropdown: ChevronDown,
  radioOn: CircleDot, radioOff: Circle, edit: Pencil, ratingEmpty: Star, ratingFull: Star,
  ratingHalf: StarHalf, loading: Loader2, first: ChevronsLeft, last: ChevronsRight,
  unfold: ArrowUpDown, file: File, plus: Plus, minus: Minus, calendar: Calendar, clock: Clock,
  treeviewCollapse: ChevronDown, treeviewExpand: ChevronRight, eyeDropper: Pipette, upload: Upload,
  color: Pipette, command: Command, ctrl: ChevronUp, space: Minus, shift: ArrowUp, alt: Command,
  enter: CornerDownLeft, arrowup: ArrowUp, arrowdown: ArrowDown, arrowleft: ArrowLeft,
  arrowright: ArrowRight, backspace: Delete, tag: Tag,
}

// Cada alias de Vuetify apunta al icono del mismo nombre: "$close" → registro.close.
export const aliases = Object.fromEntries(
  Object.keys(registro).map(k => [k, k]),
) as unknown as IconAliases

// Componente que Vuetify usa para pintar un icono del set "lucide".
const LucideIcon = defineComponent({
  name: 'LucideIcon',
  props: {
    icon: { type: [String, Object, Function], required: true },
    tag: { type: String, required: true },
  },
  setup (props, { attrs }) {
    return () => {
      const comp = typeof props.icon === 'string' ? registro[props.icon] : props.icon as Component
      return h(props.tag, mergeProps(attrs, { style: null }), [
        comp ? h(comp, { 'size': '1em', 'stroke-width': 1.75, 'aria-hidden': 'true' }) : null,
      ])
    }
  },
})

export const lucide: IconSet = { component: LucideIcon as IconSet['component'] }
