# health_connect_import

Ingest de datos de salud enviados por la app companion de Android (Health Connect).

## Endpoint

```
POST /api/v1/users/{userId}/import/health-connect
Authorization: Bearer <session-id>
Content-Type: application/json
```

Body: ver `ports/use_case.go` (`Payload`). Cada record incluye su `id` de Health
Connect (UUID estable), usado para deduplicar entre re-syncs.

## Persistencia

| Tipo del payload   | Destino                              | Dedup                          |
|--------------------|--------------------------------------|--------------------------------|
| `weight`           | feature `weight`                     | `external_id` (`hc:weight:<id>`) o `(user_id, date)` |
| `exercise_sessions`| feature `exercise` (type mapeado)    | `external_id` (`hc:exercise:<id>`) |
| `steps`            | tabla `external_health_records` (raw)| `external_id` (`hc:steps:<id>`) |
| `sleep`            | tabla `external_health_records` (raw)| `external_id` (`hc:sleep:<id>`) |
| `heart_rate`       | tabla `external_health_records` (raw)| `external_id` (`hc:heart_rate:<id>`) |

`steps`/`sleep`/`heart_rate` se guardan crudos (JSONB) hasta que tengan features
propias; no se pierde nada. Re-enviar el mismo payload no duplica (upsert por
external_id; weight además respeta el unique `(user_id, date)`).

La respuesta resume `created`/`skipped`/`blocked` por tipo (`blocked` = día cerrado).

## Probar localmente

```sh
task db:up                 # levanta Postgres (requiere Docker)
# correr el server con las env vars (POSTGRES_*, PORT, SEED_ADMIN_EMAIL/PASSWORD)

# 1. login para obtener un session id
curl -s -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"<admin>","password":"<pass>"}'

# 2. importar (usar el id devuelto y el userId)
curl -s -X POST localhost:8080/api/v1/users/<userId>/import/health-connect \
  -H "Authorization: Bearer <session-id>" \
  -H 'Content-Type: application/json' \
  -d @sample_payload.json

# 3. repetir el paso 2: todo debe volver como "skipped" (dedup OK)
```
