# API Reference – Go Feature Flag Service

Base URL (local):
http://localhost:8080

---

## Health Check

GET /healthz
curl http://localhost:8080/healthz
Response:
{
  "ok": true
}

---

## List Flags

GET /flags
curl http://localhost:8080/flags

Response:
{
  "flags": [ ... ]
}

---

## Create Flag

POST /flags
curl -X POST http://localhost:8080/flags \
  -H "Content-Type: application/json" \
  -d '{"name":"beta_banner","enabled":true,"envs":["dev","prod"],"rollout":{"type":"percentage","percentage":50}}'

Body:
{
  "name": "beta_banner",
  "enabled": true,
  "envs": ["dev", "prod"],
  "rollout": {
    "type": "percentage",
    "percentage": 50
  }
}

Response: 201 Created


## Get Flag
curl http://localhost:8080/flags/beta_banner

## Patch Flag
curl -X PATCH http://localhost:8080/flags/beta_banner \
  -H "Content-Type: application/json" \
  -d '{"enabled":false}'

## Delete Flag
curl -X DELETE http://localhost:8080/flags/beta_banner

## Evaluate
curl "http://localhost:8080/evaluate/new_checkout?env=prod&user=celeste@example.com"


