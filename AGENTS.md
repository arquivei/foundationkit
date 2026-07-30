# Project Context

## What this service is

`foundationkit` is Arquivei/Qive's internal, opinionated Go library that standardizes the
basic needs of APIs and Workers: application lifecycle, structured errors, logging,
distributed tracing, Request ID propagation, middlewares and serialization utilities. The
goal stated in `CONTRIBUTING.md` is to provide "a single standard way of doing things" and
to make developing, deploying and maintaining services easier.

It is a **library** (there is no `main` binary), published as open source under BSD
3-Clause at `github.com/arquivei/foundationkit` and consumed by several company services.

> **Current status: the library is being broken apart.** No new functionality is being
> developed. Only **bug fixes** and **CVE fixes** are accepted.

## Commands

```bash
# Build all packages (same command as CI)
go build -v ./...

# Tests
go test ./...                     # everything
go test -v ./errors/...           # a single package
go test -run TestName ./errors/

# Vet (dedicated CI job)
go vet ./...

# Lint — required before opening a PR (config in .golangci.yml)
golangci-lint run

# Formatting (gofmt and goimports are the formatters enabled in .golangci.yml)
gofmt -l .
goimports -w .

# Vulnerabilities (same check as the Govulcheck workflow)
govulncheck ./...

# Local documentation
godoc -http=localhost:6060        # open http://localhost:6060
```

## Architecture

Single module (`go.mod` at the root, Go 1.25) with flat top-level packages — one package
per responsibility, no `internal/` and no `cmd/`:

| Package | Responsibility |
|---|---|
| `app` | Application lifecycle: admin port (9000 by default) with metrics, debug and Kubernetes probes; priority-ordered graceful shutdown; `SetupConfig` via uconfig |
| `errors` | Structured error with `Severity`, `Code`, `Op` and `KeyValue`; variadic constructor `errors.E` |
| `log` | zerolog setup; hooks and sink for Stackdriver (`log/stackdriver`), go-kit adapter (`log/kitlogger`) |
| `trace` | Legacy distributed tracing (OpenCensus + Stackdriver exporter) |
| `trace/v2` | OpenTelemetry tracing: `Setup`, `Start`, middleware for mux and go-kit, `ToMap`/`FromMap` for workers |
| `request` | Creation and propagation of Request ID (ULID) via context and HTTP |
| `contextmap` | Map of values attached to the context + endpoint middleware |
| `apiutil` | Response entities and the standard error encoder for HTTP APIs |
| `httpcomm` | HTTP/JSON client for service-to-service communication |
| `httpmiddlewares` | `net/http` middlewares: `enrichloggingmiddleware`, `trackingmiddleware` |
| `gokitmiddlewares` | go-kit endpoint middlewares: `backoff`, `dontpanic`, `logging`, `metrics` (v1–v4), `stale`, `timeout`, `tracking` |
| `metrifier` | Prometheus metrics collection per operation span |
| `retrier` | Retry with exponential backoff and a generic retry evaluator |
| `avroutil` | Avro encoder/decoder, Schema Registry wire format and mocks |
| `schemaregistry` | Schema `Repository` interface + implementation in `implschemaregistry` |
| `message` | Standard event structure (`SchemaVersion3`) for queues |
| `sefaz` | SEFAZ domain: `accesskey`, `cuf`, `nsu`, `stakeholder` |
| `ref`, `gzip`, `stringsutil` | Utilities (pointers, compression, strings) |

The expected flow in library consumers is: `app.SetupConfig` →
`log.SetupLoggerWithContext` → `app.NewDefaultApp` → registering shutdown handlers and
initializing dependencies → `app.RunAndWait(mainLoop)`.

## Conventions

**Breaking changes are versioned per subdirectory.** The module is at `v0` (current tag
`v0.10.7`), so incompatible changes create a new versioned directory instead of altering
the existing package — a pattern already applied in `trace/v2` and
`gokitmiddlewares/metricsmiddleware/v2`, `/v3`, `/v4`. Old versions stay in the
repository.

**File names by role** (repeated across several packages):

- `doc.go` — package documentation for godoc, with usage examples (`app`, `trace`, `trace/v2`, `request`)
- `config.go` — the package's `Config` struct, loaded by `app.SetupConfig` (`default:"..."` tags)
- `zerolog.go` — `MarshalZerologObject` implementation for the package's types
- `errors.go` — package-specific `Code`s and errors
- `middleware.go` — endpoint/HTTP middleware
- `mock.go` — mocks exported for consumers
- `entity_*.go` — transport entities (`apiutil`, `httpcomm`)

**Errors.** Always `errors.E(op, err, severity, code)` from this library's own `errors`
package. Every function declares `const op = errors.Op("package.Func")`. Severities:
`SeverityInput` (expected caller error → 400), `SeverityRuntime` (retryable, e.g.
timeout), `SeverityFatal` (not retryable).

**Logging.** zerolog everywhere. Loggable types implement `MarshalZerologObject` and are
used with `log.Ctx(ctx).Info().EmbedObject(x)`.

**Tests.** `*_test.go` in the same package, table-driven with `stretchr/testify`
(`assert`/`require`). Fixtures in `testdata/` (e.g. `avroutil/testdata/schemas`,
`stringsutil/testdata/fuzz`). Runnable examples in `examples/` — excluded from lint.

**Commits.** Semantic commits (`renovate.json` uses `:semanticCommits`); history follows
`fix(...)`, `chore(deps): ...`.

**Style.** `.editorconfig` (LF, trailing newline). Lint with 20 linters enabled in
`.golangci.yml`, including `gosec`, `gocritic`, `gocyclo`, `dupl`, `prealloc`, `unparam`.

> Note: `mise.toml` at the root is **not versioned** and today pins only the `claude` CLI
> — it does not manage the project's Go toolchain.

## Relevant dependencies

**Main frameworks and libraries** (direct, from `go.mod`):

- `github.com/go-kit/kit` + `go-kit/log` — endpoints and middlewares
- `github.com/rs/zerolog` — structured logging (the base of the whole `log` package)
- `github.com/gorilla/mux` — HTTP routing (used in `trace/v2/http.go`)
- `github.com/omeid/uconfig` — configuration via env var + `config.json` (`app.SetupConfig`)
- `go.opentelemetry.io/otel` (+ `sdk`, `otlptracehttp`, `otelmux`, `go-logr/zerologr`) — `trace/v2`
- `go.opencensus.io` + `contrib.go.opencensus.io/exporter/stackdriver` — `trace` (legacy)
- `cloud.google.com/go/logging` — `log/stackdriver`
- `github.com/prometheus/client_golang` — metrics in `app`, `metrifier` and `metricsmiddleware`
- `github.com/arquivei/avro/v2` — `avroutil` and `schemaregistry`. Internally maintained
  fork of the discontinued `github.com/hamba/avro/v2`; it carries CVE fixes that upstream
  never released, so do not switch back to `hamba/avro`
- `github.com/oklog/ulid/v2` — IDs in `request` and `message`
- `github.com/stretchr/testify` — tests

**External services assumed by consumers:** Kubernetes (probes at `/healthy` and
`/ready`), Prometheus (admin port scrape), Google Cloud Logging / Stackdriver, OTLP
collector and Confluent Schema Registry.

## What NOT to do

- **Do not add new functionality.** The library is being broken apart; accept only bug
  and CVE fixes.
- **Do not break the public API.** The library is consumed by several services. If an
  incompatible change is unavoidable, create a new versioned directory (`.../v2`) instead
  of altering the existing package.
- **Do not update dependencies manually** and do not run `go mod tidy` on your own —
  bumps belong to Renovate (grouped, Mondays before 8am).
  `github.com/omeid/uconfig@v1.2.1` is pinned in `ignoreDeps` and must not be updated.
- **Do not use the stdlib `errors` directly** in library code — use
  `github.com/arquivei/foundationkit/errors`.
- **Do not remove old package versions** (`trace`, `metricsmiddleware` v1–v3); there are
  still consumers on them.
- **Do not open a PR without running `golangci-lint run`** — an explicit `README.md`
  requirement.

## Extra information

- Public godoc: https://pkg.go.dev/github.com/arquivei/foundationkit
- Go Report Card: https://goreportcard.com/report/github.com/arquivei/foundationkit
- PR criteria: `CONTRIBUTING.md` · Code of conduct: `CODE_OF_CONDUCT.md`
- Versioning: SemVer via repository tags (currently `v0.10.7`)
- OTel SDK environment variables used by `trace/v2`:
  https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/
- CI (GitHub Actions, in `.github/workflows/`): `go.yml` (build/test/vet),
  `golangci-lint.yml`, `govulcheck.yml`, `stale.yml`
- CODEOWNERS: `@victormn` and `@rjfonseca` review the whole repository
