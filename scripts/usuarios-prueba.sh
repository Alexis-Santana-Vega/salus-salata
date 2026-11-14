#!/usr/bin/env bash
# Usuarios del entorno de PRUEBAS (nunca producción).
#
#   bash scripts/usuarios-prueba.sh admin   restablece al admin con INITIAL_ADMIN_PASSWORD
#                                           de .env; pide cambiarla al entrar y quita el
#                                           bloqueo por intentos fallidos.
#   bash scripts/usuarios-prueba.sh qa      crea o restablece "qa-admin" (rol Admin) para
#                                           probar las vistas sin tocar la cuenta del admin.
#                                           Su contraseña es QA_ADMIN_PASSWORD de .env (se
#                                           genera la primera vez). No pide cambiarla.
#
# Las contraseñas nunca se imprimen.
set -euo pipefail
cd "$(dirname "$0")/.."
export COMPOSE_PROJECT_NAME=salus-pruebas COMPOSE_FILE=docker-compose.yml COMPOSE_ENV_FILES=.env

leer() { { grep "^$1=" .env || true; } | cut -d= -f2- | tr -d '\r'; }
sql() { docker compose exec -T db sh -c 'isql -q -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD" localhost:salus'; }
hash() { printf '%s\n' "$1" | (cd backend && PATH="$PATH:/c/Program Files/Go/bin" go run ./cmd/hash); }

case "${1:-}" in
  admin)
    usuario=$(leer INITIAL_ADMIN_USUARIO); pass=$(leer INITIAL_ADMIN_PASSWORD)
    [ -n "$pass" ] || { echo "Falta INITIAL_ADMIN_PASSWORD en .env"; exit 1; }
    h=$(hash "$pass")
    # Contraseña, sesiones y bloqueo en una sola transacción.
    printf "UPDATE USUARIO SET HASH_CONTRASENA = '%s', DEBE_CAMBIAR_CONTRASENA = TRUE, ACTIVO = TRUE, ACTUALIZADO_EN = CURRENT_TIMESTAMP WHERE USUARIO = '%s';\nUPDATE SESION SET REVOCADA_EN = CURRENT_TIMESTAMP WHERE REVOCADA_EN IS NULL AND USUARIO_ID = (SELECT ID FROM USUARIO WHERE USUARIO = '%s');\nDELETE FROM AUDITORIA_ACCESO WHERE EVENTO IN ('LOGIN_FALLIDO', 'LOGIN_BLOQUEADO') AND IDENTIFICADOR = '%s';\nCOMMIT;\n" "$h" "$usuario" "$usuario" "$usuario" | sql
    echo "Listo: \"$usuario\" entra con INITIAL_ADMIN_PASSWORD de .env y debe cambiarla al entrar."
    ;;
  qa)
    pass=$(leer QA_ADMIN_PASSWORD)
    if [ -z "$pass" ]; then
      pass="QA-$(openssl rand -hex 12)"
      printf '\n# Usuario para probar vistas en pruebas (scripts/usuarios-prueba.sh qa).\nQA_ADMIN_PASSWORD=%s\n' "$pass" >> .env
    fi
    h=$(hash "$pass")
    printf "UPDATE OR INSERT INTO USUARIO (ID, USUARIO, CORREO, NOMBRE, HASH_CONTRASENA, ROL_ID, ACTIVO, DEBE_CAMBIAR_CONTRASENA, CREADO_EN, ACTUALIZADO_EN)\nVALUES (COALESCE((SELECT ID FROM USUARIO WHERE USUARIO = 'qa-admin'), GEN_UUID()), 'qa-admin', 'qa-admin@salus.local', 'Pruebas de vista', '%s', (SELECT ID FROM ROL WHERE CLAVE = 'ADMIN'), TRUE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)\nMATCHING (USUARIO);\nDELETE FROM AUDITORIA_ACCESO WHERE EVENTO IN ('LOGIN_FALLIDO', 'LOGIN_BLOQUEADO') AND IDENTIFICADOR = 'qa-admin';\nCOMMIT;\n" "$h" | sql
    echo "Listo: \"qa-admin\" (rol Admin) entra con QA_ADMIN_PASSWORD de .env."
    ;;
  *) sed -n '2,11p' "$0"; exit 1 ;;
esac
