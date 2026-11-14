# 07 · Firebird en Docker y el alias de la base de datos

## Qué es el alias
Firebird guarda cada base de datos en **un archivo** (`.fdb`). Para conectarse, el cliente le dice al servidor *qué archivo abrir*. Hay dos formas:

| Forma | Cadena de conexión | Problema |
|---|---|---|
| Ruta real | `db:/var/lib/firebird/data/salus.fdb` | Expone cómo está organizado el servidor y obliga a cambiar la app si el archivo se mueve |
| **Alias** | `db:salus` | Ninguno: el servidor traduce el nombre a la ruta |

El alias es un **apodo** que vive en el archivo `databases.conf` **del servidor**. La aplicación solo conoce `salus`; la ruta real queda del lado de Firebird. Si mañana el archivo cambia de lugar, solo se edita `databases.conf`.

```
 Backend ──► "db:salus" ──► Firebird lee databases.conf ──► salus = /var/lib/firebird/data/salus.fdb ──► abre el archivo
```

## Cómo se ve en la configuración

**`infra/firebird/databases.conf`** (el archivo del servidor):
```ini
# Entrada original de la imagen (obligatoria)
security.db = $(dir_secDb)/security5.fdb
{
	RemoteAccess = false
	DefaultDbCachePages = 256
}

# Alias de SALUS
salus = /var/lib/firebird/data/salus.fdb
```

**`docker-compose.yml`** monta ese archivo encima del que trae la imagen:
```yaml
volumes:
  - salus-db-data:/var/lib/firebird/data                              # aquí vive salus.fdb
  - ./infra/firebird/databases.conf:/opt/firebird/databases.conf:ro   # nuestros alias
```

**`.env`** (lo que usa la app; nunca la ruta):
```ini
FIREBIRD_DB_ALIAS=salus
```

> El montaje **reemplaza** el archivo completo de la imagen, por eso se conserva la entrada `security.db`; sin ella Firebird no arranca.

## Comprobación
Conectado por alias, Firebird responde con el archivo real que abrió:
```
isql -user SALUS_APP localhost:salus
ARCHIVO_REAL                      USUARIO
/var/lib/firebird/data/salus.fdb  SALUS_APP
```
Un alias que no existe (`localhost:otra`) se rechaza con `I/O error during "open"`.

## Datos del contenedor
| Concepto | Valor |
|---|---|
| Imagen | `firebirdsql/firebird:5.0.4` |
| Servicio / contenedor | `db` / `salus-db` |
| Archivo de datos | volumen `salus-db-data` → `/var/lib/firebird/data/salus.fdb` (UTF8) |
| Usuarios | `SYSDBA` (solo mantenimiento) y `SALUS_APP` (la aplicación); contraseñas en `.env` |
| Puerto | `3050`, publicado solo en `127.0.0.1` (desarrollo). En producción no se publica: el backend entra por la red interna de Docker |
| Zona horaria | `America/Mexico_City` |

## Comandos útiles
```bash
docker compose up -d db          # levantar la BD
docker compose logs -f db        # ver el log
docker compose down              # detener (los datos se conservan en el volumen)
docker compose down -v           # detener y BORRAR los datos
```
