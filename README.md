# go-template

Шаблон Go-сервиса: HTTP API, слои domain → app → adapters, PostgreSQL (pgx + sqlc + goose), Redis в compose, инструменты закреплены в `go.mod` (`go tool`).

## Новый сервис

```powershell
gonew github.com/7cout/go-template github.com/7cout/<svc>
cd <svc>
git init -b main
task build; task test; task lint
```

`gonew` меняет путь модуля в `go.mod` и импортах. Остальное — руками:

- [ ] `.golangci.yml` → `goimports.local-prefixes`
- [ ] `proto/**/*.proto` → `go_package`
- [ ] `docker-compose.yml` → `name:`
- [ ] `AGENTS.md`, `README.md` → описание сервиса
- [ ] заменить пример `Item` своей моделью

## Требования

Go (версия из `go.mod`), gcc для `-race`, [Task](https://taskfile.dev), Docker Desktop (для `task up`), `sqlc`, `goose`, `buf` — для соответствующих задач.

## Быстрый старт

```powershell
task run                         # in-memory на :8080
curl -X POST localhost:8080/items -d '{"name":"apple"}'
task up; task migrate; task run:pg   # на Postgres
```
