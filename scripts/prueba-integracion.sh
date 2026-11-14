#!/usr/bin/env bash
# Pruebas de integración del paso 0.6: todo a través de nginx (puerta única).
# Uso: bash scripts/prueba-integracion.sh   (con docker compose arriba)
# Lee la contraseña inicial de .env sin imprimirla. Al final deja la base
# como estaba: borra usuarios, sesiones y auditoría de prueba para que el
# Sysadmin inicial se vuelva a crear con la contraseña de .env.
set -u
cd "$(dirname "$0")/.."
F=http://localhost:${FRONTEND_PORT:-8081}
B=http://localhost:8080
PASS=$(grep '^INITIAL_ADMIN_PASSWORD=' .env | cut -d= -f2- | tr -d '\r')
NUEVA="Integracion-Salus-2026-ok"
J='Content-Type: application/json'
T=$(mktemp -d)
ok=0; mal=0
probar() { # probar "descripción" esperado obtenido
  if [ "$2" = "$3" ]; then ok=$((ok+1)); echo "  ✔ $1"; else mal=$((mal+1)); echo "  ✘ $1 (esperado $2, obtenido $3)"; fi
}
cod() { curl -s -o "$T/r" -w '%{http_code}' "$@"; }
sql() { docker exec -i salus-db sh -c 'isql -q -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD" localhost:salus'; }

echo "Frontend (nginx)"
probar "GET / sirve la app" 200 "$(cod $F/)"
probar "SPA: /login sirve index.html" 200 "$(cod $F/login)"
probar "SPA: ruta inexistente sirve index.html" 200 "$(cod $F/algo/que/no/existe)"
probar "cabecera CSP presente" 1 "$(curl -s -D - -o /dev/null $F/ | grep -ci '^content-security-policy')"
probar "/docs del backend no se expone" 404 "$(cod $F/docs)"
probar "API por nginx: /api/v1/health" 200 "$(cod $F/api/v1/health)"

echo "IP real detrás de nginx"
cod -H "$J" -H 'X-Real-IP: 1.2.3.4' -d '{"usuario":"ip-prueba","contrasena":"xxxxxxxxxxxx"}' $F/api/v1/auth/login >/dev/null
IP=$(echo "SELECT FIRST 1 TRIM(IP) FROM AUDITORIA_ACCESO WHERE IDENTIFICADOR='ip-prueba' ORDER BY ID DESC;" | sql | tail -2 | head -1 | tr -d ' ')
probar "se registró una IP" 1 "$(echo "$IP" | grep -cE '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$')"
probar "nginx reemplaza X-Real-IP falsa (no es 1.2.3.4)" 0 "$(echo "$IP" | grep -c '1.2.3.4')"
probar "no se registra la IP de nginx" 0 "$(echo "$IP" | grep -c '172.28.0.10')"
cod -H "$J" -H 'X-Real-IP: 9.9.9.9' -d '{"usuario":"ip-directa","contrasena":"xxxxxxxxxxxx"}' $B/api/v1/auth/login >/dev/null
IP2=$(echo "SELECT FIRST 1 TRIM(IP) FROM AUDITORIA_ACCESO WHERE IDENTIFICADOR='ip-directa' ORDER BY ID DESC;" | sql | tail -2 | head -1 | tr -d ' ')
probar "directo al backend se ignora X-Real-IP" 0 "$(echo "$IP2" | grep -c '9.9.9.9')"

echo "Sesión por el mismo origen"
probar "login incorrecto → 401" 401 "$(cod -H "$J" -d '{"usuario":"admin","contrasena":"incorrecta-123"}' $F/api/v1/auth/login)"
probar "login correcto → 200" 200 "$(curl -s -c $T/c -D $T/h -o $T/r -w '%{http_code}' -H "$J" -d "{\"usuario\":\"admin\",\"contrasena\":\"$PASS\"}" $F/api/v1/auth/login)"
probar "cookie httpOnly + SameSite=Strict + ruta /api/v1/auth" 1 "$(grep -i '^set-cookie: salus_refresh' $T/h | grep -i 'httponly' | grep -i 'samesite=strict' | grep -ci 'path=/api/v1/auth')"
A=$(grep -oE '"accessToken":"[^"]+"' $T/r | cut -d'"' -f4)
probar "cambio de contraseña obligatorio bloquea /usuarios" 403 "$(cod -H "Authorization: Bearer $A" $F/api/v1/usuarios)"
probar "cambio de contraseña → 200" 200 "$(curl -s -c $T/c -o $T/r -w '%{http_code}' -X PUT -H "$J" -H "Authorization: Bearer $A" -d "{\"contrasenaActual\":\"$PASS\",\"contrasenaNueva\":\"$NUEVA\"}" $F/api/v1/auth/contrasena)"
A=$(grep -oE '"accessToken":"[^"]+"' $T/r | cut -d'"' -f4)
probar "Sysadmin entra a /usuarios" 200 "$(cod -H "Authorization: Bearer $A" $F/api/v1/usuarios)"
probar "refresh con la cookie (recarga de página)" 200 "$(curl -s -b $T/c -c $T/c -o $T/r -w '%{http_code}' -X POST $F/api/v1/auth/refresh)"
probar "logout → 204" 204 "$(curl -s -b $T/c -c $T/c -o /dev/null -w '%{http_code}' -X POST $F/api/v1/auth/logout)"
probar "refresh tras logout → 401" 401 "$(curl -s -b $T/c -o /dev/null -w '%{http_code}' -X POST $F/api/v1/auth/refresh)"

echo "Limpieza"
printf 'DELETE FROM AUDITORIA_ACCESO;\nDELETE FROM SESION;\nDELETE FROM USUARIO;\nCOMMIT;\n' | sql
docker compose restart backend >/dev/null 2>&1
rm -rf "$T"
echo
echo "Resultado: $ok correctas, $mal fallidas"
[ "$mal" -eq 0 ]
