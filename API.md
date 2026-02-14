# API Reference – Go Feature Flag Service

Base URL (local):
http://localhost:8080

---

## Health Check

GET /healthz

Response:
{
  "ok": true
}

---

## List Flags

GET /flags

Response:
{
  "flags": [ ... ]
}

---

## Create Flag

POST /flags

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
