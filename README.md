# SALUS

Plataforma de renta de equipo médico e insumos para cirugías.

```
backend/    API en Go (Fiber v3)
frontend/   Vue 3 + Vuetify + TypeScript (pnpm)
docs/       Análisis, requerimientos, modelo de datos, diseño y plan
```

## Requisitos
Git, Go 1.27+, Node 26+, pnpm 12+ y Docker. **No usar npm.**

## Desarrollo local

Base de datos (Firebird 5 con alias `salus`, ver [docs/07-firebird-alias.md](docs/07-firebird-alias.md)):
```bash
copy .env.example .env
docker compose up -d db
```

Backend (http://localhost:8080/health):
```bash
cd backend
go run ./cmd/api
```

Frontend (http://localhost:3000):
```bash
cd frontend
pnpm install
pnpm dev
```

El plan por fases está en [docs/06-plan-de-trabajo.md](docs/06-plan-de-trabajo.md).
