# Go Feature Flag Service (zero-cost)

A lightweight feature-flag service in Go (mini LaunchDarkly-style) with:
- REST API for flags CRUD
- Evaluation endpoint with env targeting + percentage rollouts
- In-memory store (no DB required)
- Optional seed file (`configs/flags.json`)
- GitHub Actions CI (fmt, vet, test -race, build)

## Run
```bash
go run ./cmd/server -addr :8080 -seed ./configs/flags.json

📘 Full API documentation: see [[API.md](https://github.com/celery-github/go-feature-flags/blob/main/API.md)](./API.md)
