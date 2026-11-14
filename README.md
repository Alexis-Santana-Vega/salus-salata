# SALUS

Plataforma de renta de equipo médico e insumos para cirugías.

```
backend/    API en Go (Fiber v3)
frontend/   Vue 3 + Vuetify + TypeScript (pnpm)
infra/      Configuración de Firebird (alias)
docs/       Análisis, requerimientos, modelo de datos, diseño, plan y arquitectura
```

## Requisitos
Git, Go 1.27+, Node 26+, pnpm 12+ y Docker. **No usar npm.**

## Levantar todo con Docker
```bash
copy .env.example .env
docker compose up -d --build
```
- App: http://localhost:8081 (nginx, puerta única)
- API directa (solo desarrollo): http://localhost:8080/api/v1/health
- Documentación (Scalar): http://localhost:8080/docs (con `DOCS_ENABLED=true`)
- Firebird: `127.0.0.1:3050`, alias `salus` ([docs/07-firebird-alias.md](docs/07-firebird-alias.md))

## Desarrollo local sin contenedor del backend
```bash
docker compose up -d db
cd backend
go run ./cmd/api
```

Frontend (http://localhost:3000):
```bash
cd frontend
pnpm install
pnpm dev
```

## Pruebas
```bash
bash scripts/prueba-integracion.sh
```
```bash
cd backend && go vet ./... && go test ./...
cd frontend && pnpm lint && pnpm build
```

Arquitectura: [docs/08-arquitectura.md](docs/08-arquitectura.md) · Plan: [docs/06-plan-de-trabajo.md](docs/06-plan-de-trabajo.md)
