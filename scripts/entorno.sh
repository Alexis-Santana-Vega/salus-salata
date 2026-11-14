#!/usr/bin/env bash
# Maneja los dos entornos de SALUS. Cada uno es un proyecto de docker compose
# con su propia base de datos (volumen), red, puertos y secretos.
#
#   bash scripts/entorno.sh <pruebas|produccion> <comando>
#
# Comandos:
#   up                 levanta (o actualiza) el entorno compilando las imágenes
#   down               detiene el entorno (los datos se conservan)
#   ps | logs          estado y logs (logs sigue en vivo; Ctrl+C para salir)
#   probar             pruebas del entorno (integración en pruebas, validación en producción)
#   respaldar          respaldo .fbk de su base en respaldos/
#   copiar-a-pruebas   (solo desde produccion) restaura el último respaldo en pruebas
set -euo pipefail
cd "$(dirname "$0")/.."

entorno=${1:-}; comando=${2:-}
export COMPOSE_PATH_SEPARATOR=:
case "$entorno" in
  pruebas)
    export COMPOSE_PROJECT_NAME=salus-pruebas COMPOSE_FILE=docker-compose.yml COMPOSE_ENV_FILES=.env
    ;;
  produccion)
    export COMPOSE_PROJECT_NAME=salus-prod COMPOSE_FILE=docker-compose.yml:docker-compose.prod.yml COMPOSE_ENV_FILES=.env.produccion
    [ -f .env.produccion ] || { echo "Falta .env.produccion (copia .env.produccion.example)."; exit 1; }
    [ -f infra/certs/salus.crt ] || { echo "Falta el certificado en infra/certs/ (salus.crt y salus.key)."; exit 1; }
    # Puertos para las pruebas (se leen sin ejecutar el archivo).
    for v in HTTP_PORT HTTPS_PORT; do
      val=$(grep -E "^$v=" .env.produccion | cut -d= -f2- | tr -d '\r')
      [ -n "$val" ] && export "$v=$val"
    done
    ;;
  *) sed -n '2,15p' "$0"; exit 1 ;;
esac

# Ejecuta isql dentro del contenedor de la base de este entorno.
sql() { docker compose exec -T db sh -c 'isql -q -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD" localhost:salus'; }

# Saca un respaldo con gbak (servicio de Firebird, con la base en uso) y lo copia a respaldos/.
respaldar() {
  mkdir -p respaldos
  local archivo="respaldos/salus-${entorno}-$(date +%Y%m%d-%H%M%S).fbk"
  docker compose exec -T db sh -c 'gbak -b -g -se localhost:service_mgr salus /tmp/respaldo.fbk -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD" && chmod 644 /tmp/respaldo.fbk'
  docker compose cp db:/tmp/respaldo.fbk "$archivo"
  docker compose exec -T db rm -f /tmp/respaldo.fbk
  echo "$archivo"
}

case "$comando" in
  up)   docker compose up -d --build ;;
  down) docker compose down ;;
  ps)   docker compose ps ;;
  logs) docker compose logs -f --tail 100 ;;
  probar)
    if [ "$entorno" = pruebas ]; then bash scripts/prueba-integracion.sh; else bash scripts/prueba-produccion.sh; fi ;;
  respaldar) echo "Respaldo creado: $(respaldar)" ;;
  copiar-a-pruebas)
    [ "$entorno" = produccion ] || { echo "Se ejecuta como: entorno.sh produccion copiar-a-pruebas"; exit 1; }
    echo "Esto REEMPLAZA la base de pruebas con los datos de producción (datos reales)."
    read -r -p "Escribe COPIAR para continuar: " ok; [ "$ok" = COPIAR ] || exit 1
    archivo=$(respaldar)
    # A partir de aquí se trabaja sobre el entorno de pruebas.
    export COMPOSE_PROJECT_NAME=salus-pruebas COMPOSE_FILE=docker-compose.yml COMPOSE_ENV_FILES=.env
    docker compose stop backend
    docker compose cp "$archivo" db:/tmp/respaldo.fbk
    # Restaura como el usuario de la app de pruebas: queda como dueño de la base.
    docker compose exec -T db sh -c 'gbak -rep -se localhost:service_mgr /tmp/respaldo.fbk salus -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD"; rm -f /tmp/respaldo.fbk'
    docker compose start backend
    echo "Pruebas ahora tiene los datos de producción del respaldo $archivo"
    ;;
  *) sed -n '2,15p' "$0"; exit 1 ;;
esac
