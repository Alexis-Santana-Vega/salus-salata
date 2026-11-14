#!/usr/bin/env bash
# Pruebas de ataque contra la app levantada con docker compose.
# Uso: bash scripts/prueba-seguridad.sh
# Crea un usuario temporal, ataca por nginx (puerta única) y al final borra
# usuario, sesiones y bitácora de prueba (también quita los bloqueos).
# Tarda ~4 min: respeta el límite de 10 logins/min por IP entre bloques.
set -u
cd "$(dirname "$0")/.."
F=http://localhost:${FRONTEND_PORT:-8081}
A=$F/api/v1
J='Content-Type: application/json'
T=$(mktemp -d)
ok=0; mal=0
probar() { if [ "$2" = "$3" ]; then ok=$((ok+1)); echo "  ✔ $1"; else mal=$((mal+1)); echo "  ✘ $1 (esperado $2, obtenido $3)"; fi; }
cod() { curl -s -o "$T/r" -w '%{http_code}' "$@"; }
codigo() { grep -oE '"code":"[A-Z_]+"' "$T/r" | head -1 | cut -d'"' -f4; }
sql() { docker exec -i salus-db sh -c 'isql -q -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD" localhost:salus'; }
esperar() { echo "  … esperando 61 s (límite de login por minuto)"; sleep 61; }
b64url() { base64 -w0 | tr '+/' '-_' | tr -d '='; }

# Usuario temporal con contraseña aleatoria (no se imprime).
PASS="Seg-$(head -c 12 /dev/urandom | base64 | tr -d '/+=')"
HASH=$(printf '%s\n' "$PASS" | (cd backend && PATH="$PATH:/c/Program Files/Go/bin" go run ./cmd/hash))
printf "INSERT INTO USUARIO (USUARIO, CORREO, NOMBRE, HASH_CONTRASENA, ROL_ID, ACTIVO, DEBE_CAMBIAR_CONTRASENA, CREADO_EN, ACTUALIZADO_EN) VALUES ('atacante', 'atacante@salus.local', 'Prueba Seguridad', '%s', 3, TRUE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);\nCOMMIT;\n" "$HASH" | sql

echo "1. Inyección SQL y validación de entrada"
probar "SQLi ' OR '1'='1 en usuario → credenciales inválidas" INVALID_CREDENTIALS "$(cod -H "$J" -d "{\"usuario\":\"' OR '1'='1\",\"contrasena\":\"' OR '1'='1\"}" $A/auth/login >/dev/null; codigo)"
probar "SQLi admin'-- → credenciales inválidas" INVALID_CREDENTIALS "$(cod -H "$J" -d "{\"usuario\":\"admin'--\",\"contrasena\":\"x\"}" $A/auth/login >/dev/null; codigo)"
probar "campo de 300 caracteres → 422" 422 "$(cod -H "$J" -d "{\"usuario\":\"$(printf 'a%.0s' {1..300})\",\"contrasena\":\"x\"}" $A/auth/login)"
probar "JSON malformado → 400" 400 "$(cod -H "$J" -d '{"usuario":' $A/auth/login)"
probar "texto plano en escritura → 415" 415 "$(cod -H 'Content-Type: text/plain' -d 'usuario=admin' $A/auth/login)"
head -c 1100000 /dev/zero | tr '\0' 'a' > "$T/grande"
probar "cuerpo de 1.1 MB → 413" 413 "$(cod -H "$J" --data-binary @"$T/grande" $A/auth/login)"
probar "sin User-Agent → 403" 403 "$(cod -A '' $A/health)"

esperar
echo "2. Tokens falsificados y escalada de rol"
cod -c "$T/c" -H "$J" -d "{\"usuario\":\"atacante\",\"contrasena\":\"$PASS\"}" $A/auth/login >/dev/null
TOK=$(grep -oE '"accessToken":"[^"]+"' "$T/r" | cut -d'"' -f4)
probar "sesión de Doctor válida" 200 "$(cod -H "Authorization: Bearer $TOK" $A/auth/me)"
probar "Doctor no entra a /usuarios (solo Sysadmin)" 403 "$(cod -H "Authorization: Bearer $TOK" $A/usuarios)"
H=$(printf '{"alg":"none","typ":"JWT"}' | b64url)
P=$(printf '{"sub":"1","rol":"SYSADMIN","iss":"salus-api","exp":4102444800}' | b64url)
probar "JWT alg=none con rol SYSADMIN → 401" 401 "$(cod -H "Authorization: Bearer $H.$P." $A/usuarios)"
FIRMA=$(echo "$TOK" | cut -d. -f3)
H2=$(echo "$TOK" | cut -d. -f1)
probar "JWT con rol cambiado y firma original → 401" 401 "$(cod -H "Authorization: Bearer $H2.$P.$FIRMA" $A/usuarios)"
HS=$(printf '{"alg":"HS256","typ":"JWT"}' | b64url)
FALSA=$(printf '%s' "$HS.$P" | openssl dgst -sha256 -hmac "secreto-adivinado" -binary | b64url)
probar "JWT firmado con otro secreto → 401" 401 "$(cod -H "Authorization: Bearer $HS.$P.$FALSA" $A/usuarios)"
probar "token basura → 401" 401 "$(cod -H "Authorization: Bearer abc.def.ghi" $A/auth/me)"
probar "SQL en parámetro con rol sin permiso → 403 (no llega a la consulta)" 403 "$(cod -H "Authorization: Bearer $TOK" "$A/usuarios?pageSize=1;DROP%20TABLE%20USUARIO")"

echo "3. Robo de refresh token"
cp "$T/c" "$T/robada"
probar "refresh legítimo rota la cookie" 200 "$(cod -b "$T/c" -c "$T/c" -X POST $A/auth/refresh)"
probar "atacante reusa la cookie robada → 401" 401 "$(cod -b "$T/robada" -X POST $A/auth/refresh)"
probar "tras el reuso, la sesión legítima también se cierra" 401 "$(cod -b "$T/c" -X POST $A/auth/refresh)"

echo "4. nginx: archivos y cabeceras"
for ruta in "/.env" "/.git/config" "/../../etc/passwd" "/%2e%2e/%2e%2e/etc/passwd" "/api/../.env"; do
  curl -s --path-as-is "$F$ruta" -o "$T/r"
  probar "no expone $ruta" 0 "$(grep -cE 'JWT_SECRET|FIREBIRD_PASSWORD|root:x:0|\[core\]' "$T/r")"
done
curl -s -D "$T/h" -o /dev/null $F/
probar "X-Frame-Options DENY (anti clickjacking)" 1 "$(grep -ci '^x-frame-options: deny' "$T/h")"
probar "sin versión de nginx en Server" 0 "$(grep -ciE '^server: nginx/[0-9]' "$T/h")"
probar "sin X-Powered-By" 0 "$(grep -ci '^x-powered-by' "$T/h")"
probar "CORS: origen ajeno sin permiso" 0 "$(curl -s -D - -o /dev/null -X OPTIONS -H 'Origin: https://malo.com' -H 'Access-Control-Request-Method: POST' $A/auth/login | grep -ci '^access-control-allow-origin')"
probar "TRACE no permitido" 1 "$(cod -X TRACE $F/ | grep -cE '^(405|404|403)$')"

esperar
echo "5. Fuerza bruta y enumeración de usuarios"
for i in 1 2 3 4 5; do cod -H "$J" -d '{"usuario":"atacante","contrasena":"incorrecta-123"}' $A/auth/login >/dev/null; done
probar "6.º intento sobre la cuenta → 429 aunque la contraseña sea correcta" 429 "$(cod -H "$J" -d "{\"usuario\":\"atacante\",\"contrasena\":\"$PASS\"}" $A/auth/login)"
t1=$(curl -s -o /dev/null -w '%{time_total}' -H "$J" -d '{"usuario":"no-existe-xyz","contrasena":"incorrecta-123"}' $A/auth/login)
t2=$(curl -s -o /dev/null -w '%{time_total}' -H "$J" -d '{"usuario":"admin","contrasena":"incorrecta-123"}' $A/auth/login)
probar "tiempos parecidos con usuario existente e inexistente (<2x)" 1 "$(awk -v a="$t1" -v b="$t2" 'BEGIN{r=(a>b)?a/b:b/a; print (r<2)?1:0}')"

esperar
echo "6. Bloqueo por IP (20 fallos en 15 min, rotando cuentas)"
for i in $(seq 1 9); do cod -H "$J" -d "{\"usuario\":\"rota-$i\",\"contrasena\":\"x\"}" $A/auth/login >/dev/null; done
esperar
for i in $(seq 10 18); do cod -H "$J" -d "{\"usuario\":\"rota-$i\",\"contrasena\":\"x\"}" $A/auth/login >/dev/null; done
probar "cuenta nueva desde la misma IP → 429" 429 "$(cod -H "$J" -d '{"usuario":"cuenta-limpia","contrasena":"x"}' $A/auth/login)"

echo "7. Ráfaga de peticiones (bots)"
n429=$(for i in $(seq 1 140); do curl -s -o /dev/null -w '%{http_code}\n' $A/health; done | grep -c 429)
probar "más de 120/min por IP recibe 429" 1 "$([ "$n429" -gt 0 ] && echo 1 || echo 0)"

echo "Limpieza"
printf "DELETE FROM SESION WHERE USUARIO_ID IN (SELECT ID FROM USUARIO WHERE USUARIO='atacante');\nDELETE FROM AUDITORIA_ACCESO;\nDELETE FROM USUARIO WHERE USUARIO='atacante';\nCOMMIT;\n" | sql
rm -rf "$T"
echo
echo "Resultado: $ok correctas, $mal fallidas"
[ "$mal" -eq 0 ]
