# --- Etapa 1: compilar -------------------------------------------------------
FROM golang:1.27-alpine AS build
WORKDIR /src

# Dependencias primero: se cachean mientras go.mod/go.sum no cambien.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/api ./cmd/api

# --- Etapa 2: imagen final mínima --------------------------------------------
# Sin shell ni gestor de paquetes; corre como usuario no root (uid 65532).
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/api /app/api
ENV TZ=America/Mexico_City
EXPOSE 8080
USER nonroot:nonroot
HEALTHCHECK --interval=15s --timeout=5s --start-period=20s --retries=3 \
    CMD ["/app/api", "-healthcheck"]
ENTRYPOINT ["/app/api"]
