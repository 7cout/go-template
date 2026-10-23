# go-template

Эталонный Go-сервис воркспейса `C:\7cout\ai`: HTTP API, слои domain → app → adapters, PostgreSQL через sqlc + pgx.
Новый сервис: `gonew github.com/7cout/go-template github.com/7cout/<svc>` (затем чек-лист в README).

## Команды (PowerShell, Windows)

| Команда | Что делает |
| --- | --- |
| `task build` | `go build ./...` + бинарь `bin/server` |
| `task test` | `go test -race -count=1 ./...` (нужен gcc, `CGO_ENABLED=1`) |
| `task lint` | `go tool golangci-lint run ./...` — версия закреплена в `go.mod` (`tool`) |
| `task fmt` | goimports + gofumpt |
| `task vuln` | govulncheck |
| `task all` | fmt → lint → test → build → vuln |
| `task up` / `task down` | Postgres + Redis в docker compose (тома сохраняются) |
| `task migrate` | goose up из `db/migrations` |
| `task sqlc` | перегенерировать `internal/adapters/postgres/sqlcgen` |
| `task run` / `task run:pg` | сервер in-memory / на Postgres |

Задача готова, только когда `task lint`, `task test`, `task build` зелёные (ТР-8.7). Во время итераций достаточно тестов изменённых пакетов (`go test -race ./internal/app/...`), перед сдачей — полный `task test`.

## Архитектура

- `internal/domain` — сущности, правила, доменные ошибки (`ErrNotFound`, `ErrInvalidInput`). Без зависимостей от БД и транспорта.
- `internal/app` — use cases. Интерфейсы репозиториев объявлены **здесь** (у потребителя).
- `internal/adapters/{memory,postgres,httpapi}` — реализации. postgres мапит `pgx.ErrNoRows` → `domain.ErrNotFound`, ошибки драйвера наружу не оборачивает (`%v`).
- `cmd/server` — сборка зависимостей, graceful shutdown по `signal.NotifyContext`.
- `db/migrations` — goose, они же схема для sqlc. `db/queries` — SQL для sqlc.
- `proto/` — контракт gRPC, пока только `buf lint` (генерация не подключена).

## Правила

- Конвенции Go — skill `go-conventions`. Коротко: `%w` + `errors.Is/As`; `ctx` первым аргументом и не в struct; table-driven тесты; конструктор возвращает конкретный тип.
- Код сервиса запускается на хосте (`go run`), в контейнерах только инфраструктура; bind mount исходников не использовать.
- Сгенерированный код (`sqlcgen`) руками не править — менять SQL и `task sqlc`.
- Новая зависимость → `task vuln`.
- Секреты не коммитить; `.env*` в .gitignore. Dev-пароли compose — только значения по умолчанию.
- После правки `*.go` хук сам запускает goimports + gofumpt — не форматируй вручную.
- Версия Go — из `go.mod`.

## Вики

Накопленные знания — вики воркспейса `C:\7cout\ai\knowledge\`. **Префиксы этого проекта: `go-`, `pg-`** (новому сервису — добавь свои из реестра `knowledge\AGENTS.md`).

- Вопрос по Go, PostgreSQL, библиотекам или инструментам → сначала `C:\7cout\ai\knowledge\index.md` (разделы своих префиксов), затем нужные `wiki\<префикс>-*.md`; ссылайся на страницы. Вики главнее общих знаний модели; если она расходится с кодом — скажи об этом.
- Явный запрос — `/kb-query`. Новое знание сохраняется из сессии в `C:\7cout\ai` через `/kb-ingest` или `/wrap-up`; отсюда вики не править.

## gopls MCP

Если подключён MCP `gopls` (`claude mcp add gopls -s user -- gopls mcp`), следуй инструкциям ниже (вывод `gopls mcp -instructions`).
Исключение для этого воркспейса: перед сдачей задачи полный `task test` обязателен, несмотря на шаг 7.

### Detecting a Go workspace

At the start of every session, you MUST use the `go_workspace` tool to learn about the Go workspace. ONLY if you are in a Go workspace, you MUST run `go_vulncheck` immediately afterwards to identify any existing security risks.

### Read workflow

1. **Understand the workspace layout**: use `go_workspace`.
2. **Find relevant symbols**: use `go_search` (fuzzy). EXAMPLE: `go_search({"query":"server"})`
3. **Understand a file and its intra-package dependencies**: `go_file_context` MUST be used immediately after reading any Go file for the first time. EXAMPLE: `go_file_context({"file":"/path/to/server.go"})`
4. **Understand a package's public API**: `go_package_api`. EXAMPLE: `go_package_api({"packagePaths":["example.com/internal/storage"]})`

### Editing workflow

1. **Read first**: follow the Read workflow.
2. **Find references**: before modifying any symbol, use `go_symbol_references`. EXAMPLE: `go_symbol_references({"file":"/path/to/server.go","symbol":"Server.Run"})`
3. **Make edits**, including edits to references.
4. **Check for errors**: after every modification, you MUST call `go_diagnostics` with the edited files.
5. **Fix errors** and re-run `go_diagnostics`; 'hint'/'info' may be ignored if irrelevant.
6. **Check for vulnerabilities**: if go.mod changed, run `go_vulncheck({"pattern":"./..."})`.
7. **Run tests** for the changed packages once there are no errors.
