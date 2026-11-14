# --- Etapa 1: compilar la app con pnpm --------------------------------------
# Node 24 trae corepack, que instala la versión de pnpm fijada en package.json.
FROM node:24-alpine AS build
WORKDIR /app
RUN corepack enable

# Dependencias primero: se cachean mientras el lockfile no cambie.
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm build-only

# --- Etapa 2: nginx sin root sirviendo los archivos estáticos ---------------
FROM nginxinc/nginx-unprivileged:1.28-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 8080
