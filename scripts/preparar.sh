#!/usr/bin/env bash
# Prepara una máquina nueva después de clonar el repo (Git Bash en Windows).
#   bash scripts/preparar.sh
# Crea lo que git no guarda, SOLO si no existe (nunca sobrescribe):
#   .env             secretos del entorno de pruebas (aleatorios)
#   .env.produccion  secretos del entorno de producción (aleatorios, distintos)
#   infra/certs/     certificado autofirmado para probar producción en local
# e instala las dependencias del frontend. No imprime contraseñas.
set -euo pipefail
cd "$(dirname "$0")/.."

azar()   { head -c 32 /dev/urandom | base64 | tr -d '/+=\n' | head -c "$1"; }
hexazar() { head -c 48 /dev/urandom | od -An -tx1 | tr -d ' \n'; }

# Copia la plantilla cambiando los valores "cambia-..." por secretos aleatorios.
generar() { # generar <plantilla> <destino>
  if [ -f "$2" ]; then echo "• $2 ya existe: se conserva"; return; fi
  sed -e "s|^FIREBIRD_ROOT_PASSWORD=.*|FIREBIRD_ROOT_PASSWORD=$(azar 24)|" \
      -e "s|^FIREBIRD_PASSWORD=.*|FIREBIRD_PASSWORD=$(azar 24)|" \
      -e "s|^JWT_SECRET=.*|JWT_SECRET=$(hexazar)|" \
      -e "s|^INITIAL_ADMIN_PASSWORD=.*|INITIAL_ADMIN_PASSWORD=$(azar 20)|" "$1" > "$2"
  echo "✔ $2 creado (la contraseña inicial de admin está en INITIAL_ADMIN_PASSWORD)"
}

generar .env.example .env
generar .env.produccion.example .env.produccion

if [ -f infra/certs/salus.crt ]; then
  echo "• infra/certs/ ya tiene certificado: se conserva"
else
  mkdir -p infra/certs
  MSYS_NO_PATHCONV=1 openssl req -x509 -nodes -newkey rsa:2048 -days 365 -subj "/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1" \
    -keyout infra/certs/salus.key -out infra/certs/salus.crt 2>/dev/null
  chmod 644 infra/certs/salus.key
  echo "✔ certificado autofirmado de prueba en infra/certs/ (en el servidor real va el del dominio)"
fi

echo "• Instalando dependencias del frontend (pnpm)…"
(cd frontend && pnpm install --frozen-lockfile >/dev/null) && echo "✔ dependencias listas"

echo
echo "Listo. Siguiente paso:  bash scripts/entorno.sh pruebas up"
