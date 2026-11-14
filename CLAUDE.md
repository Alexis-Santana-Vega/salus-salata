# CLAUDE.md · Reglas del proyecto SALUS

SALUS renta equipo médico y surte insumos a varios hospitales para cirugías. El contexto completo está en `docs/` (01 análisis, 02 requerimientos, 03 modelo de datos, 04 pendientes, 05 guía de diseño, 06 plan de trabajo). Consúltalos antes de implementar algo del negocio.

## Forma de trabajo
- Avanzar por **pasos cortos** del plan (`docs/06-plan-de-trabajo.md`). Cada paso se entrega funcionando y se revisa antes del siguiente.
- No adelantar trabajo de otros pasos aunque parezca fácil.
- No inventar reglas de negocio, estados ni contenido (p. ej. plantillas). Si falta información, preguntar **solo lo que se desprende de los docs**, sin preguntas especulativas.
- Al terminar un paso: compila, pasa lint y se probó. Decir claramente qué quedó y qué falta.
- Responder en español.

## Codificación de archivos (obligatorio)
- Todo archivo de texto en **UTF-8 sin BOM** y fin de línea **LF** (`.editorconfig` y `.gitattributes` lo fijan).
- Crear y editar archivos con las herramientas de edición, **nunca con PowerShell 5.1** (`Set-Content`, `Out-File`, `>`): lee como ANSI, agrega BOM y rompe la ñ y los acentos.
- Si hay que transformar un archivo por consola, usar Bash (sed, heredoc).
- Después de escribir texto en español, verificar que no aparezcan `Ã`, `Â` ni BOM.

## Stack (no proponer alternativas)
- `frontend/`: Vue 3 + Vuetify + TypeScript + Pinia + Vue Router + ESLint.
- `backend/`: Go + Fiber v3.
- BD: Firebird 5 en Docker (`firebirdsql/firebird`) con **alias** en `databases.conf`. La ruta real del `.fdb` nunca va en la cadena de conexión.
- Todo se levanta con `docker compose`; cada servicio (BD, back, front) en su propio contenedor.
- **pnpm siempre; npm prohibido** (comandos, scripts, Dockerfiles, README).
- Documentación de la API con **Scalar**, nunca Swagger UI.

## Backend (Go)
- Estructura idiomática y simple: `cmd/api` + `internal/` con `config`, `database`, `model`, `repository`, `service`, `handler`, `router`, `middleware`. **Sin sobreingeniería**: nada de capas o interfaces que no se usen.
- Secretos y configuración por variables de entorno; `.env` nunca se sube (solo `.env.example`).
- SQL parametrizado siempre.

## API (consistente en todo el proyecto)
- REST, versión en la ruta: `/api/v1/...`; recursos en plural y en español (`/api/v1/hospitales`).
- Verbos y códigos HTTP correctos (200, 201, 204, 400, 401, 403, 404, 409, 422, 429, 500).
- JSON en `camelCase`. Un **solo formato de error** para toda la API.
- Listados con los mismos parámetros de paginación, búsqueda, filtros y orden.
- Validación de entrada en cada endpoint y autorización por rol en el back.

## Seguridad de usuarios
- Contraseñas con hash y sal (argon2id o bcrypt); nunca en texto plano ni en logs.
- Límite de intentos por usuario y por IP en una ventana de tiempo, con bloqueo temporal (anti fuerza bruta).
- Roles guardados en BD (Sysadmin, Admin, Doctor). Validación de rol y de datos en **back y front**; el back es la autoridad.
- Tokens de acceso de vida corta + refresh revocable. Auditoría de accesos.

## Frontend y diseño
- Seguir `docs/05-guia-diseno.md`. Resumen: mobile first, Poppins, primario #0195B7 / secundario #4FB5AB, campos outlined, modo oscuro requerido, sin fondos negros, en pasteles el texto lleva el mismo tono más fuerte.
- Botón primario sólido; secundario solo texto; de icono sin borde ni fondo. Si dispara una carga: se deshabilita, muestra spinner y dice "Cargando".
- El usuario da instrucciones de diseño al llegar al paso del front; no adelantarse.

## Negocio (lo que más se confunde)
- **Equipo** (reutilizable, modelo maestro + unidades, se cobra **solo por hora**) ≠ **Insumo** (desechable, por stock, los no usados regresan).
- La cirugía tiene **solo 5 estados**: Programada → En entrega → En hospital → Devuelta → Cerrada (+ Cancelada). Sin estados intermedios.
- Cotización inicial al programar; la final solo descuenta lo retornado y suma gastos cobrables.
- **Fuera de alcance**: facturación, pagos y contabilidad. Roles User SALUS y User Doc en espera.

## Comandos
```bash
# Backend
cd backend && go run ./cmd/api
cd backend && go vet ./... && go build ./...

# Frontend
cd frontend && pnpm dev
cd frontend && pnpm lint && pnpm build

# Todo
docker compose up -d
```
