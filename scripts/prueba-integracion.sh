#!/usr/bin/env bash
# Pruebas de integración del paso 0.6: todo a través de nginx (puerta única).
# Uso: bash scripts/prueba-integracion.sh   (con docker compose arriba)
# Lee la contraseña inicial de .env sin imprimirla. Al final deja la base
# como estaba: borra usuarios, sesiones y auditoría de prueba para que el
# Sysadmin inicial se vuelva a crear con la contraseña de .env.
set -u
# Estas pruebas borran usuarios y bitácora: NUNCA contra producción.
if [ "${COMPOSE_PROJECT_NAME:-salus-pruebas}" = "salus-prod" ]; then echo "Estas pruebas solo corren en el entorno de pruebas."; exit 1; fi
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
# El curl de Windows altera los acentos en la línea de comandos: los cuerpos
# con texto en español se mandan desde archivo (@archivo).
cuerpo() { printf "%s" "$1" > "$T/b.json"; echo "@$T/b.json"; }
# La API limita 30 escrituras/min por IP. Entre bloques se reinicia el backend
# (el contador vive en memoria) para que un bloque no consuma el cupo del otro.
reiniciar_limites() {
  docker compose restart backend >/dev/null 2>&1
  until [ "$(curl -s -o /dev/null -w '%{http_code}' $B/api/v1/health)" = 200 ]; do sleep 1; done
}
sql() { docker compose exec -T db sh -c 'isql -q -user "$FIREBIRD_USER" -password "$FIREBIRD_PASSWORD" localhost:salus'; }

# Parte del admin con la contraseña inicial de .env (sin importar si se cambió a mano).
bash scripts/usuarios-prueba.sh admin >/dev/null

echo "Frontend (nginx)"
probar "GET / sirve la app" 200 "$(cod $F/)"
probar "SPA: /login sirve index.html" 200 "$(cod $F/login)"
probar "SPA: ruta inexistente sirve index.html" 200 "$(cod $F/algo/que/no/existe)"
probar "cabecera CSP presente" 1 "$(curl -s -D - -o /dev/null $F/ | grep -ci '^content-security-policy')"
probar "/docs del backend no se expone" 404 "$(cod $F/docs)"
probar "API por nginx: /api/v1/health" 200 "$(cod $F/api/v1/health)"

echo "IP real detrás de nginx"
cod -H "$J" -H 'X-Real-IP: 1.2.3.4' -d '{"usuario":"ip-prueba","contrasena":"xxxxxxxxxxxx"}' $F/api/v1/auth/login >/dev/null
IP=$(echo "SELECT FIRST 1 TRIM(IP) FROM AUDITORIA_ACCESO WHERE IDENTIFICADOR='ip-prueba' ORDER BY FECHA DESC;" | sql | tail -2 | head -1 | tr -d ' ')
probar "se registró una IP" 1 "$(echo "$IP" | grep -cE '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$')"
probar "nginx reemplaza X-Real-IP falsa (no es 1.2.3.4)" 0 "$(echo "$IP" | grep -c '1.2.3.4')"
probar "no se registra la IP de nginx" 0 "$(echo "$IP" | grep -c '172.28.0.10')"
cod -H "$J" -H 'X-Real-IP: 9.9.9.9' -d '{"usuario":"ip-directa","contrasena":"xxxxxxxxxxxx"}' $B/api/v1/auth/login >/dev/null
IP2=$(echo "SELECT FIRST 1 TRIM(IP) FROM AUDITORIA_ACCESO WHERE IDENTIFICADOR='ip-directa' ORDER BY FECHA DESC;" | sql | tail -2 | head -1 | tr -d ' ')
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

# Parte de catálogos vacíos: lo que se probó a mano en la vista no interfiere.
printf 'DELETE FROM CONTACTO_HOSPITAL;\nDELETE FROM HOSPITAL;\nDELETE FROM DOCTOR;\nDELETE FROM TECNICO;\nDELETE FROM PROCEDIMIENTO;\nCOMMIT;
' | sql
echo "Catálogo de procedimientos (patrón CRUD)"
AU="Authorization: Bearer $A"
P=$F/api/v1/procedimientos
probar "alta → 201" 201 "$(cod -H "$AU" -H "$J" --data-binary "$(cuerpo '{"nombre":"  Cirugía   General "}')" $P)"
PID=$(grep -oE '"id":"[^"]+"' $T/r | cut -d'"' -f4)
probar "guarda el nombre sin espacios de más" 1 "$(grep -c '"nombre":"Cirug[^"]*a General"' $T/r)"
probar "duplicado sin acentos ni mayúsculas → 409" 409 "$(cod -H "$AU" -H "$J" -d '{"nombre":"CIRUGIA general"}' $P)"
probar "el 409 señala el campo nombre" 1 "$(grep -c '"field":"nombre"' $T/r)"
probar "nombre vacío → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"nombre":""}' $P)"
probar "nombre con HTML → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"nombre":"<b>x</b>"}' $P)"
cod -H "$AU" -H "$J" -d '{"nombre":"Holep"}' $P >/dev/null
probar "búsqueda sin acentos (?q=cirugia) → 1 resultado" 1 "$(cod -H "$AU" "$P?q=cirugia" >/dev/null; grep -o '"total":[0-9]*' $T/r | cut -d: -f2)"
probar "orden descendente (?sort=-nombre) → 200" 200 "$(cod -H "$AU" "$P?sort=-nombre")"
probar "orden por campo no permitido → 422" 422 "$(cod -H "$AU" "$P?sort=ID")"
probar "consulta por id → 200" 200 "$(cod -H "$AU" $P/$PID)"
probar "id mal formado → 404" 404 "$(cod -H "$AU" $P/123)"
probar "edición → 200" 200 "$(cod -X PUT -H "$AU" -H "$J" --data-binary "$(cuerpo '{"nombre":"Cirugía General y Laparoscopía"}')" $P/$PID)"
probar "baja lógica → 204" 204 "$(cod -X DELETE -H "$AU" $P/$PID)"
probar "dado de baja ya no se consulta → 404" 404 "$(cod -H "$AU" $P/$PID)"
probar "dado de baja sigue en la base (historial)" 1 "$(echo "SELECT COUNT(*) FROM PROCEDIMIENTO WHERE ACTIVO = FALSE;" | sql | tail -2 | head -1 | tr -d ' ')"
probar "el nombre de uno dado de baja se puede reusar → 201" 201 "$(cod -H "$AU" -H "$J" --data-binary "$(cuerpo '{"nombre":"Cirugía General y Laparoscopía"}')" $P)"
probar "sin sesión → 401" 401 "$(cod $P)"

reiniciar_limites
echo "Catálogo de técnicos"
TC=$F/api/v1/tecnicos
probar "alta con puesto y teléfono → 201" 201 "$(cod -H "$AU" -H "$J" -d '{"nombre":"Rodrigo","puesto":"AYUDANTE","telefono":"+52 477 123 4567"}' $TC)"
TID=$(grep -oE '"id":"[^"]+"' $T/r | cut -d'"' -f4)
probar "alta sin teléfono → 201" 201 "$(cod -H "$AU" -H "$J" -d '{"nombre":"Dylan","puesto":"AYUDANTE"}' $TC)"
probar "nombre repetido se permite (personas distintas) → 201" 201 "$(cod -H "$AU" -H "$J" -d '{"nombre":"Rodrigo","puesto":"LOGISTICA"}' $TC)"
probar "puesto fuera de la lista → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"nombre":"X","puesto":"JEFE"}' $TC)"
probar "teléfono inválido → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"nombre":"X","puesto":"OTRO","telefono":"abc"}' $TC)"
probar "el CHECK de la base también rechaza un puesto inválido" 1 "$(printf "INSERT INTO TECNICO (ID, NOMBRE, PUESTO, CREADO_EN, ACTUALIZADO_EN) VALUES (GEN_UUID(), 'X', 'JEFE', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);\nROLLBACK;\n" | sql 2>&1 | grep -c CK_TECNICO_PUESTO)"
probar "búsqueda por teléfono (?q=477) → 1 resultado" 1 "$(cod -H "$AU" "$TC?q=477" >/dev/null; grep -o '"total":[0-9]*' $T/r | cut -d: -f2)"
probar "orden por puesto → 200" 200 "$(cod -H "$AU" "$TC?sort=-puesto")"
probar "edición quita el teléfono → 200" 200 "$(cod -X PUT -H "$AU" -H "$J" -d '{"nombre":"Rodrigo","puesto":"ENFERMERIA","telefono":""}' $TC/$TID)"
probar "el teléfono vacío queda vacío" 1 "$(grep -c '"telefono":""' $T/r)"
probar "baja lógica → 204" 204 "$(cod -X DELETE -H "$AU" $TC/$TID)"
probar "dado de baja → 404" 404 "$(cod -H "$AU" $TC/$TID)"
reiniciar_limites
echo "Catálogo de direcciones (C.P.)"
probar "C.P. 42000 → 200" 200 "$(cod -H "$AU" $F/api/v1/codigos-postales/42000)"
probar "trae la colonia Centro de Pachuca, HGO" 1 "$(grep -c '"nombre":"Centro".*"municipio":"Pachuca de Soto","estado":"HGO"' $T/r)"
ASEN=$(grep -oE '"id":"[^"]+","cp":"42000","nombre":"Centro"' $T/r | cut -d'"' -f4)
probar "C.P. fuera del catálogo → 404" 404 "$(cod -H "$AU" $F/api/v1/codigos-postales/99999)"
probar "C.P. mal escrito → 422" 422 "$(cod -H "$AU" $F/api/v1/codigos-postales/4200)"
probar "C.P. sin sesión → 401" 401 "$(cod $F/api/v1/codigos-postales/42000)"
echo "Catálogo de doctores"
DC=$F/api/v1/doctores
probar "alta mínima (nombre + correo) → 201" 201 "$(cod -H "$AU" -H "$J" -d '{"nombre":"Dr. Aguayo","correo":"Aguayo@Mail.com"}' $DC)"
DID=$(grep -oE '"id":"[^"]+"' $T/r | cut -d'"' -f4)
probar "correo se guarda en minúsculas" 1 "$(grep -c '"correo":"aguayo@mail.com"' $T/r)"
probar "alta completa con colonia del catálogo → 201" 201 "$(cod -H "$AU" -H "$J" -d "{\"numero\":\"doc-001\",\"nombre\":\"Dr. Westendarp\",\"correo\":\"west@mail.com\",\"telPersonal\":\"477 111 2222\",\"direccion\":{\"calle\":\"Av. Juarez\",\"numExterior\":\"434\",\"asentamientoId\":\"$ASEN\"}}" $DC)"
probar "la respuesta completa C.P., municipio y estado" 1 "$(grep -c '"cp":"42000","colonia":"Centro","municipio":"Pachuca de Soto","estado":"HGO"' $T/r)"
probar "número en mayúsculas" 1 "$(grep -c '"numero":"DOC-001"' $T/r)"
probar "número repetido → 409 en numero" 1 "$(cod -H "$AU" -H "$J" -d '{"numero":"DOC-001","nombre":"X","correo":"x@mail.com"}' $DC >/dev/null; grep -c '"field":"numero"' $T/r)"
probar "correo repetido con mayúsculas → 409 en correo" 1 "$(cod -H "$AU" -H "$J" -d '{"nombre":"X","correo":"AGUAYO@mail.com"}' $DC >/dev/null; grep -c '"field":"correo"' $T/r)"
probar "colonia que no existe en el catálogo → 422 en asentamientoId" 1 "$(cod -H "$AU" -H "$J" -d '{"nombre":"X","correo":"y@mail.com","direccion":{"asentamientoId":"0193f1c2-0000-7000-8000-000000000000"}}' $DC >/dev/null; grep -c '"field":"asentamientoId"' $T/r)"
probar "correo inválido → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"nombre":"X","correo":"no-es-correo"}' $DC)"
probar "búsqueda por correo → 1" 1 "$(cod -H "$AU" "$DC?q=west@" >/dev/null; grep -o '"total":[0-9]*' $T/r | cut -d: -f2)"
probar "baja → 204" 204 "$(cod -X DELETE -H "$AU" $DC/$DID)"
probar "correo de uno dado de baja se reusa → 201" 201 "$(cod -H "$AU" -H "$J" -d '{"nombre":"Dr. Aguayo Segundo","correo":"aguayo@mail.com"}' $DC)"

reiniciar_limites
echo "Catálogo de hospitales"
HC=$F/api/v1/hospitales
CUERPO_H='{"clave":"hgo-001","nombre":"Hospital Angeles Pachuca","telefono":"771 123 4567","direccion":{"calle":"Av Revolucion","numExterior":"434","asentamientoId":"'$ASEN'"},"latitud":20.1,"longitud":-98.7,"ligaMapa":"https://maps.google.com/?q=pachuca","contactos":[{"nombre":"Juan Perez","puesto":"Quirofano","telefono":"771 000 1111"},{"nombre":"Ana Lopez","puesto":"Almacen"}]}'
probar "alta con 2 contactos → 201" 201 "$(cod -H "$AU" -H "$J" -d "$CUERPO_H" $HC)"
HID=$(grep -oE '"id":"[^"]+"' $T/r | head -1 | cut -d'"' -f4)
probar "clave en mayúsculas" 1 "$(grep -c '"clave":"HGO-001"' $T/r)"
probar "contactos en el orden capturado" 1 "$(grep -c '"contactos":\[{"nombre":"Juan Perez".*{"nombre":"Ana Lopez"' $T/r)"
probar "dirección completa desde el catálogo" 1 "$(grep -c '"municipio":"Pachuca de Soto","estado":"HGO"' $T/r)"
probar "clave repetida → 409 en clave" 1 "$(cod -H "$AU" -H "$J" -d "${CUERPO_H/Angeles/Otro}" $HC >/dev/null; grep -c '"field":"clave"' $T/r)"
probar "prefijo distinto al estado de la colonia → 422 en clave" 1 "$(cod -H "$AU" -H "$J" -d "${CUERPO_H/hgo-001/GTO-002}" $HC >/dev/null; grep -c '"field":"clave","message":"La clave debe empezar con HGO' $T/r)"
probar "clave sin formato → 422" 422 "$(cod -H "$AU" -H "$J" -d "${CUERPO_H/hgo-001/HGO001}" $HC)"
probar "latitud fuera de rango → 422" 422 "$(cod -H "$AU" -H "$J" -d "${CUERPO_H/20.1/200}" $HC)"
probar "liga sin https → 422" 422 "$(cod -H "$AU" -H "$J" -d "${CUERPO_H/https:/http:}" $HC)"
probar "contacto sin nombre → 422 en contactos[1].nombre" 1 "$(cod -H "$AU" -H "$J" -d "${CUERPO_H/Ana Lopez/}" $HC >/dev/null; grep -c '"field":"contactos\[1\].nombre"' $T/r)"
probar "sin dirección → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"clave":"HGO-009","nombre":"X"}' $HC)"
probar "edición deja solo los contactos enviados → 200" 200 "$(cod -X PUT -H "$AU" -H "$J" -d "${CUERPO_H/',{"nombre":"Ana Lopez","puesto":"Almacen"}'/}" $HC/$HID)"
probar "queda un solo contacto" 0 "$(grep -c 'Ana Lopez' $T/r)"
probar "búsqueda por municipio (?q=pachuca) → 1" 1 "$(cod -H "$AU" "$HC?q=pachuca" >/dev/null; grep -o '"total":[0-9]*' $T/r | cut -d: -f2)"
probar "baja → 204" 204 "$(cod -X DELETE -H "$AU" $HC/$HID)"
probar "la clave de uno dado de baja no se reusa → 409" 409 "$(cod -H "$AU" -H "$J" -d "$CUERPO_H" $HC)"

echo "Formatos de texto"
probar "nombre de procedimiento con ; → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"nombre":"Holep; DROP TABLE X"}' $P)"
probar "nombre de técnico con llaves → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"nombre":"{{x}}","puesto":"OTRO"}' $TC)"
probar "número de doctor con espacio → 422" 422 "$(cod -H "$AU" -H "$J" -d '{"numero":"DOC 1","nombre":"X","correo":"z@mail.com"}' $DC)"

reiniciar_limites
DPASS="Doctor-Integracion-2026"
DHASH=$(printf '%s\n' "$DPASS" | (cd backend && PATH="$PATH:/c/Program Files/Go/bin" go run ./cmd/hash))
printf "INSERT INTO USUARIO (ID, USUARIO, CORREO, NOMBRE, HASH_CONTRASENA, ROL_ID, ACTIVO, DEBE_CAMBIAR_CONTRASENA, CREADO_EN, ACTUALIZADO_EN) VALUES (GEN_UUID(), 'doctor-prueba', 'doctor-prueba@salus.local', 'Doctor Prueba', '%s', (SELECT ID FROM ROL WHERE CLAVE = 'DOCTOR'), TRUE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);\nCOMMIT;\n" "$DHASH" | sql
cod -H "$J" -d "{\"usuario\":\"doctor-prueba\",\"contrasena\":\"$DPASS\"}" $F/api/v1/auth/login >/dev/null
AD=$(grep -oE '"accessToken":"[^"]+"' $T/r | cut -d'"' -f4)
probar "rol Doctor no ve el catálogo → 403" 403 "$(cod -H "Authorization: Bearer $AD" $P)"
probar "rol Doctor no ve los técnicos → 403" 403 "$(cod -H "Authorization: Bearer $AD" $TC)"
probar "rol Doctor no ve el catálogo de doctores → 403" 403 "$(cod -H "Authorization: Bearer $AD" $DC)"
probar "rol Doctor no ve los hospitales → 403" 403 "$(cod -H "Authorization: Bearer $AD" $HC)"
probar "ruta inexistente de la API sigue en 404" 404 "$(cod -H "$AU" $F/api/v1/no-existe)"
probar "refresh con la cookie (recarga de página)" 200 "$(curl -s -b $T/c -c $T/c -o $T/r -w '%{http_code}' -X POST $F/api/v1/auth/refresh)"
probar "logout → 204" 204 "$(curl -s -b $T/c -c $T/c -o /dev/null -w '%{http_code}' -X POST $F/api/v1/auth/logout)"
probar "refresh tras logout → 401" 401 "$(curl -s -b $T/c -o /dev/null -w '%{http_code}' -X POST $F/api/v1/auth/refresh)"

echo "Limpieza"
printf 'DELETE FROM CONTACTO_HOSPITAL;\nDELETE FROM HOSPITAL;\nDELETE FROM DOCTOR;\nDELETE FROM TECNICO;\nDELETE FROM PROCEDIMIENTO;\nDELETE FROM AUDITORIA_ACCESO;\nDELETE FROM SESION;\nDELETE FROM USUARIO;\nCOMMIT;\n' | sql
docker compose restart backend >/dev/null 2>&1
rm -rf "$T"
echo
echo "Resultado: $ok correctas, $mal fallidas"
[ "$mal" -eq 0 ]
