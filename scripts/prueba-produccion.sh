#!/usr/bin/env bash
# Valida la app levantada en modo producción:
#   docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
#   bash scripts/prueba-produccion.sh
set -u
cd "$(dirname "$0")/.."
S=https://localhost:${HTTPS_PORT:-8443}
P=http://localhost:${HTTP_PORT:-8088}
T=$(mktemp -d)
ok=0; mal=0
probar() { if [ "$2" = "$3" ]; then ok=$((ok+1)); echo "  ✔ $1"; else mal=$((mal+1)); echo "  ✘ $1 (esperado $2, obtenido $3)"; fi; }
# -k: el certificado de prueba es autofirmado.
cod() { curl -sk -o "$T/r" -w '%{http_code}' "$@"; }

echo "Superficie expuesta"
probar "backend sin puertos publicados" 0 "$(docker compose ps backend --format '{{.Ports}}' | grep -c -- '->')"
probar "Firebird sin puertos publicados" 0 "$(docker compose ps db --format '{{.Ports}}' | grep -c -- '->')"
probar "HTTP redirige a HTTPS (301)" 301 "$(curl -s -o /dev/null -w '%{http_code}' $P/login)"
probar "redirección apunta a https" 1 "$(curl -s -D - -o /dev/null $P/ | grep -ci '^location: https://')"

echo "HTTPS y cabeceras"
probar "la app carga por HTTPS" 200 "$(cod $S/)"
curl -sk -D "$T/h" -o /dev/null $S/
probar "HSTS presente" 1 "$(grep -ci '^strict-transport-security: max-age=31536000' "$T/h")"
probar "CSP presente" 1 "$(grep -ci '^content-security-policy' "$T/h")"
curl -sk -D "$T/h" -o /dev/null $S/api/v1/health
probar "API por HTTPS: /health" 200 "$(cod $S/api/v1/health)"
probar "HSTS también en la API" 1 "$(grep -ci '^strict-transport-security' "$T/h")"
probar "/docs deshabilitado" 404 "$(cod $S/docs)"

echo "Cookie y logs"
curl -sk -D "$T/h" -o /dev/null -H 'Content-Type: application/json' -d '{"usuario":"no-existe","contrasena":"incorrecta-123"}' $S/api/v1/auth/login
probar "login fallido responde 401 por HTTPS" 401 "$(cod -H 'Content-Type: application/json' -d '{"usuario":"no-existe","contrasena":"incorrecta-123"}' $S/api/v1/auth/login)"
probar "logs del backend en JSON" 1 "$(docker compose logs backend --no-log-prefix --tail 5 | grep -c '^{"time"' | awk '{print ($1>0)?1:0}')"
probar "backend en modo production" 1 "$(docker compose logs backend --no-log-prefix | grep -c '"env":"production"' | awk '{print ($1>0)?1:0}')"

sql() { docker compose exec -T db sh -c 'isql -q -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD" localhost:salus'; }
# Usuario temporal para ver la cookie de una sesión real.
PASS="Prod-$(head -c 12 /dev/urandom | base64 | tr -d '/+=')"
HASH=$(printf '%s\n' "$PASS" | (cd backend && PATH="$PATH:/c/Program Files/Go/bin" go run ./cmd/hash))
printf "INSERT INTO USUARIO (USUARIO, CORREO, NOMBRE, HASH_CONTRASENA, ROL_ID, ACTIVO, DEBE_CAMBIAR_CONTRASENA, CREADO_EN, ACTUALIZADO_EN) VALUES ('prueba-prod', 'prueba-prod@salus.local', 'Prueba Prod', '%s', 3, TRUE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);\nCOMMIT;\n" "$HASH" | sql
curl -sk -D "$T/h" -o /dev/null -H 'Content-Type: application/json' -d "{\"usuario\":\"prueba-prod\",\"contrasena\":\"$PASS\"}" $S/api/v1/auth/login
probar "cookie de sesión con Secure + HttpOnly + SameSite=Strict" 1 "$(grep -i '^set-cookie: salus_refresh' "$T/h" | grep -i 'secure' | grep -i 'httponly' | grep -ci 'samesite=strict')"

# Limpia usuario y bitácora de prueba.
printf "DELETE FROM SESION WHERE USUARIO_ID IN (SELECT ID FROM USUARIO WHERE USUARIO='prueba-prod');\nDELETE FROM AUDITORIA_ACCESO WHERE IDENTIFICADOR IN ('no-existe','prueba-prod') OR USUARIO_ID IN (SELECT ID FROM USUARIO WHERE USUARIO='prueba-prod');\nDELETE FROM USUARIO WHERE USUARIO='prueba-prod';\nCOMMIT;\n" | sql
rm -rf "$T"
echo
echo "Resultado: $ok correctas, $mal fallidas"
[ "$mal" -eq 0 ]
