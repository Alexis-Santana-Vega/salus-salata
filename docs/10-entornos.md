# 10 · Entornos: pruebas y producción

Hay **dos entornos completos e independientes**. Cada uno es un proyecto de docker compose con su propia base de datos (volumen), red, puertos y secretos. Pueden correr a la vez en la misma máquina sin tocarse.

| | Pruebas | Producción |
|---|---|---|
| Para qué | Desarrollar, probar, romper | Datos reales del negocio |
| Proyecto compose | `salus-pruebas` | `salus-prod` |
| Archivo de secretos | `.env` | `.env.produccion` |
| Compose | `docker-compose.yml` | `docker-compose.yml` + `docker-compose.prod.yml` |
| Base de datos | volumen `salus-pruebas_salus-db-data` | volumen `salus-prod_salus-db-data` |
| App | http://localhost:8081 | https://localhost:8443 (en el servidor: 443) |
| API directa / `/docs` | http://localhost:8080, Scalar activo | No expuesta, sin `/docs` |
| Firebird | `127.0.0.1:3050` (Beekeeper) | No expuesto |
| Cookies / logs | Sin `Secure` / texto | `Secure` + HSTS / JSON |
| Red interna | `172.28.0.0/24`, nginx `172.28.0.10` | `172.29.0.0/24`, nginx `172.29.0.10` |

Los secretos de producción (`.env.produccion`), los certificados (`infra/certs/`) y los respaldos (`respaldos/`) **nunca se suben a git**.

## Comandos
Todo pasa por `scripts/entorno.sh` para no equivocarse de entorno:

```bash
bash scripts/entorno.sh pruebas up            # levantar / actualizar pruebas
bash scripts/entorno.sh produccion up         # levantar / actualizar producción
bash scripts/entorno.sh <entorno> down        # detener (los datos se conservan)
bash scripts/entorno.sh <entorno> ps          # estado
bash scripts/entorno.sh <entorno> logs        # logs en vivo
bash scripts/entorno.sh <entorno> probar      # pruebas del entorno
bash scripts/entorno.sh <entorno> respaldar   # respaldo .fbk en respaldos/
bash scripts/entorno.sh produccion copiar-a-pruebas   # datos reales → pruebas (pide confirmación)
```

Para desarrollar con recarga en caliente: levantar **pruebas** y correr `pnpm dev` en `frontend/` (http://localhost:3000, usa el backend de pruebas en :8080) o `go run ./cmd/api` en `backend/` después de `docker compose stop backend`.

## Cómo trabajar uno con otro
1. **Todo se hace primero en pruebas.** Código nuevo, migraciones y pruebas (`probar`) corren en pruebas. Nunca se prueba nada en producción.
2. **Se sube a producción solo lo que pasó en pruebas**: lint, build, `go test`, `prueba-integracion.sh` y, si tocó seguridad, `prueba-seguridad.sh`.
3. **Antes de actualizar producción, respaldo**: `entorno.sh produccion respaldar`, luego `entorno.sh produccion up`. Las migraciones se aplican solas al arrancar el backend; como ya corrieron en pruebas, no hay sorpresas.
4. **Las pruebas destructivas nunca tocan producción**: `prueba-integracion.sh` y `prueba-seguridad.sh` borran usuarios y bitácora, y se niegan a correr si el proyecto es `salus-prod`. `prueba-produccion.sh` solo lee y limpia lo que crea.
5. **Datos reales en pruebas solo si hace falta** (`copiar-a-pruebas`): reproduce un problema con datos verdaderos. Recuerda que son datos reales (doctores, hospitales): trátalos con el mismo cuidado. Las pruebas de integración vuelven a dejar pruebas sin usuarios.
6. **Secretos distintos** en cada entorno: un token o contraseña de pruebas nunca sirve en producción.

## Respaldos
- `respaldar` usa `gbak` por el servicio de Firebird: la base puede estar en uso.
- El archivo queda en `respaldos/salus-<entorno>-<fecha>.fbk`. Hay que copiarlo fuera del servidor (otro disco o la nube); un respaldo en el mismo disco no protege de una falla del disco.
- Para restaurar se usa el mismo mecanismo de `copiar-a-pruebas` (`gbak -rep` como `SALUS_APP`, con el backend detenido).

## Servidor de producción (cuando exista)
1. Clonar el repo, crear `.env.produccion` desde `.env.produccion.example` con secretos nuevos.
2. Poner el certificado del dominio en `infra/certs/salus.crt` y `salus.key`.
3. En `.env.produccion`: `BIND_IP=0.0.0.0`, `HTTP_PORT=80`, `HTTPS_PORT=443`; en `frontend/nginx.prod.conf` la redirección a `https://$host$request_uri`.
4. `bash scripts/entorno.sh produccion up` y `bash scripts/entorno.sh produccion probar`.
5. Programar `respaldar` diario (cron) y copiar los respaldos fuera del servidor.
