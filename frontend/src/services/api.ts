/**
 * Cliente HTTP de la API. Único lugar que habla con el backend.
 *
 * - Mismo origen: `/api/v1` (nginx en Docker, el proxy de Vite en desarrollo).
 * - El token de acceso viaja en `Authorization` y vive solo en memoria.
 * - Ante un 401 renueva la sesión UNA vez (una sola renovación a la vez,
 *   compartida por todas las peticiones) y reintenta.
 */
const BASE = '/api/v1'

export interface DetalleCampo { field: string, message: string }

// Error de la API ya normalizado. `mensaje` es apto para mostrar al usuario.
export class ApiError extends Error {
  constructor (
    public status: number,
    public code: string,
    public mensaje: string,
    public detalles: DetalleCampo[] = [],
  ) {
    super(code)
  }
}

interface Hooks {
  token: () => string | null
  renovar: () => Promise<boolean>
}
let hooks: Hooks = { token: () => null, renovar: async () => false }

// El store de sesión se registra aquí (evita una dependencia circular).
export function configurarApi (h: Hooks) {
  hooks = h
}

let renovando: Promise<boolean> | null = null

// Renovación compartida: si varias peticiones reciben 401 a la vez, todas
// esperan la misma promesa (dos renovaciones simultáneas parecen robo).
function renovarUnaVez () {
  renovando ??= hooks.renovar().finally(() => {
    renovando = null
  })
  return renovando
}

interface Opciones {
  metodo?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  cuerpo?: unknown
  sinRenovar?: boolean // para /auth/*: no reintentar con refresh
  completa?: boolean // devuelve {data, meta} en lugar de solo data (listados)
}

export async function api<T> (ruta: string, op: Opciones = {}): Promise<T> {
  const enviar = () => {
    const headers: Record<string, string> = { Accept: 'application/json' }
    const token = hooks.token()
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
    if (op.cuerpo !== undefined) {
      headers['Content-Type'] = 'application/json'
    }
    return fetch(BASE + ruta, {
      method: op.metodo ?? 'GET',
      headers,
      body: op.cuerpo === undefined ? undefined : JSON.stringify(op.cuerpo),
      credentials: 'same-origin', // la cookie del refresh solo viaja en /api/v1/auth
    })
  }

  let res: Response
  try {
    res = await enviar()
    if (res.status === 401 && !op.sinRenovar && hooks.token() && await renovarUnaVez()) {
      res = await enviar()
    }
  } catch {
    throw new ApiError(0, 'NETWORK_ERROR', 'No pudimos conectar con el servidor. Revisa tu conexión e intenta de nuevo.')
  }

  if (res.status === 204) {
    return undefined as T
  }
  const json = await res.json().catch(() => null)
  if (!res.ok) {
    const e = json?.error ?? {}
    throw new ApiError(res.status, e.code ?? 'INTERNAL_ERROR', mensajeAmigable(e.code, e.message), e.details ?? [])
  }
  return (op.completa ? json : json?.data) as T
}

/**
 * Traduce el código de la API a un mensaje para el usuario. Para los códigos
 * que ya traen un mensaje pensado para el usuario (validación, credenciales)
 * se usa el del servidor; para el resto, uno genérico propio.
 */
function mensajeAmigable (code?: string, delServidor?: string): string {
  switch (code) {
    case 'INVALID_CREDENTIALS':
    case 'VALIDATION_ERROR':
    case 'TOO_MANY_REQUESTS':
    case 'PASSWORD_CHANGE_REQUIRED': {
      return delServidor ?? 'Revisa los datos e intenta de nuevo.'
    }
    case 'UNAUTHORIZED': {
      return 'Tu sesión terminó. Inicia sesión de nuevo.'
    }
    case 'FORBIDDEN': {
      return 'No tienes permiso para realizar esta acción.'
    }
    case 'NOT_FOUND': {
      return 'No encontramos lo que buscas.'
    }
    case 'CONFLICT': {
      return 'Ya existe un registro con esos datos.'
    }
    case 'SERVICE_UNAVAILABLE': {
      return 'El servicio no está disponible en este momento. Intenta en unos minutos.'
    }
    default: {
      return 'Algo salió mal. Intenta de nuevo en unos minutos.'
    }
  }
}

// Listados: misma forma en toda la API ({data, meta}).
export interface Meta { page: number, pageSize: number, total: number }
export interface Pagina<T> { data: T[], meta: Meta }

export function apiPagina<T> (ruta: string): Promise<Pagina<T>> {
  return api<Pagina<T>>(ruta, { completa: true })
}
