# herodote

Git [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) historian accross repositories.

- Full-text search with Postgresql as a backend
- Light frontend (<60kb gzipped) with desktop and responsive UI
- Filters on repository, type or component
- Light shell script for loading data into backend
- GitHub Actions provided for integration

[![Build](https://github.com/ViBiOh/herodote/workflows/Build/badge.svg)](https://github.com/ViBiOh/herodote/actions)
[![codecov](https://codecov.io/gh/ViBiOh/herodote/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/herodote)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_herodote&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_herodote)

## Concepts

Herodote aims to provide a quick view of activitiy of all your repositories. It's the changelog of your organizations.

Herodote only understands conventionnal commits. Commits that don't match expectations are ignored. To ensure you have conventionnal commits, you can use [`commit-msg hooks`](https://github.com/ViBiOh/scripts/blob/main/hooks/commit-msg) and/or a [simple GitHub Action that checks it](.github/workflows/branch_clean.yaml).

Herodote loads data with its own script, which is idempotent. On cold start, with an empty index, it only loads last 50 commits.

## Getting started

### Postgres

Herodote use a Postgres database as a backend storage. You need a Postgres database for storing your datas. You can use free tier of [ElephantSQL](https://www.elephantsql.com).

Once setup, you _have to_ create schema with [Herodote DDL](sql/ddl.sql) and start the Herodote API. Configuration is done by passing `-dbHost`, `-dbName`, `-dbUser`, `-dbPass` arg or setting equivalent environment variables (cf. [API Usage](#usage) section).

### Installation

Golang binary is built with static link. You can download it directly from the [GitHub Release page](https://github.com/ViBiOh/herodote/releases) or build it by yourself by cloning this repo and running `make`.

A Docker image is available for `amd64`, `arm` and `arm64` platforms on Docker Hub: [vibioh/herodote](https://hub.docker.com/r/vibioh/herodote/tags).

You can configure app by passing CLI args or environment variables (cf. [Usage](#usage) section). CLI override environment variables.

You'll find a Kubernetes exemple in the [`infra/`](infra) folder, using my [`app chart`](https://github.com/ViBiOh/charts/tree/main/app)

### CI Integration

Herodote is fed by its own script: [herodote.sh](herodote.sh).

It automatically detects last commit's SHA in index and add only new ones of repository.

The script needs the following variables to be set (or will prompt you for):

- `GIT_HOST`: Name of your git provider (e.g. `github.com`). It's guessed from `git remote get-url --push origin` if you are in a git folder
- `GIT_REPOSITORY`: Name of your repository (e.g. `ViBiOh/herodote`). It's guessed from `git remote get-url --push origin` if you are in a git folder
- `HERODOTE_API`: URL of your Herodote API (e.g. https://herodote.vibioh.fr)
- `HERODOTE_SECRET`: `httpSecret` or your Herodote API (cf. [API Usage](#usage) section)

If you execute your script in a non-interactive environment, set the `SCRIPTS_NO_INTERACTIVE=1` for disabling prompt, guessed value will be used.

#### GitHub Actions

You can use the following GitHub Actions for pushing your commits to Herodote index on merge to `main`.

```yaml
---
name: Herodote

permissions:
  actions: none
  checks: none
  contents: none
  deployments: none
  issues: none
  packages: none
  pages: none
  pull-requests: none
  repository-projects: none
  security-events: none

on:
  push:
    branches:
      - main

jobs:
  build:
    name: Feed
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
        with:
          ref: ${{ github.event.pull_request.head.sha }}
          fetch-depth: 0
      - name: Push
        run: |
          curl --disable --silent --show-error --location --max-time 30 "https://raw.githubusercontent.com/ViBiOh/herodote/main/herodote.sh" | bash
        env:
          HERODOTE_API: https://herodote.vibioh.fr
          HERODOTE_SECRET: ${{ secrets.HERODOTE_SECRET }}
          GIT_HOST: github.com
          GIT_REPOSITORY: ${{ github.repository }}
          SCRIPTS_NO_INTERACTIVE: "1"
```

You **have to** add secrets in your repository in the repository's settings: https://github.com/YOUR_NAME/YOUR_REPOSITORY/settings/secrets

- `HERODOTE_API`: `HERODOTE_API` from [#ci-integration](#ci-integration)
- `HERODOTE_SECRET`: `HERODOTE_SECRET` from [#ci-integration](#ci-integration)

## Endpoints

- `GET /health`: healthcheck of server, always respond [`okStatus (default 204)`](#usage)
- `GET /ready`: checks external dependencies availability and then respond [`okStatus (default 204)`](#usage) or `503` during [`graceDuration`](#usage) when `SIGTERM` is received
- `GET /version`: value of `VERSION` environment variable
- `GET /metrics`: Prometheus metrics, on a dedicated port [`prometheusPort (default 9090)`](#usage)

### Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of herodote:
  -address string
        [server] Listen address {HERODOTE_ADDRESS}
  -cert string
        [server] Certificate file {HERODOTE_CERT}
  -corsCredentials
        [cors] Access-Control-Allow-Credentials {HERODOTE_CORS_CREDENTIALS}
  -corsExpose string
        [cors] Access-Control-Expose-Headers {HERODOTE_CORS_EXPOSE}
  -corsHeaders string
        [cors] Access-Control-Allow-Headers {HERODOTE_CORS_HEADERS} (default "Content-Type")
  -corsMethods string
        [cors] Access-Control-Allow-Methods {HERODOTE_CORS_METHODS} (default "GET")
  -corsOrigin string
        [cors] Access-Control-Allow-Origin {HERODOTE_CORS_ORIGIN} (default "*")
  -csp string
        [owasp] Content-Security-Policy {HERODOTE_CSP} (default "default-src 'self'; base-uri 'self'; script-src 'self' 'httputils-nonce'; style-src 'self' 'httputils-nonce'")
  -dbHost string
        [db] Host {HERODOTE_DB_HOST}
  -dbMaxConn uint
        [db] Max Open Connections {HERODOTE_DB_MAX_CONN} (default 5)
  -dbMinConn uint
        [db] Min Open Connections {HERODOTE_DB_MIN_CONN} (default 2)
  -dbName string
        [db] Name {HERODOTE_DB_NAME}
  -dbPass string
        [db] Pass {HERODOTE_DB_PASS}
  -dbPort uint
        [db] Port {HERODOTE_DB_PORT} (default 5432)
  -dbSslmode string
        [db] SSL Mode {HERODOTE_DB_SSLMODE} (default "disable")
  -dbUser string
        [db] User {HERODOTE_DB_USER}
  -frameOptions string
        [owasp] X-Frame-Options {HERODOTE_FRAME_OPTIONS} (default "deny")
  -graceDuration duration
        [http] Grace duration when SIGTERM received {HERODOTE_GRACE_DURATION} (default 30s)
  -hsts
        [owasp] Indicate Strict Transport Security {HERODOTE_HSTS} (default true)
  -httpSecret string
        [herodote] HTTP Secret Key for Update {HERODOTE_HTTP_SECRET}
  -idleTimeout duration
        [server] Idle Timeout {HERODOTE_IDLE_TIMEOUT} (default 2m0s)
  -key string
        [server] Key file {HERODOTE_KEY}
  -loggerJson
        [logger] Log format as JSON {HERODOTE_LOGGER_JSON}
  -loggerLevel string
        [logger] Logger level {HERODOTE_LOGGER_LEVEL} (default "INFO")
  -loggerLevelKey string
        [logger] Key for level in JSON {HERODOTE_LOGGER_LEVEL_KEY} (default "level")
  -loggerMessageKey string
        [logger] Key for message in JSON {HERODOTE_LOGGER_MESSAGE_KEY} (default "message")
  -loggerTimeKey string
        [logger] Key for timestamp in JSON {HERODOTE_LOGGER_TIME_KEY} (default "time")
  -minify
        Minify HTML {HERODOTE_MINIFY} (default true)
  -okStatus int
        [http] Healthy HTTP Status code {HERODOTE_OK_STATUS} (default 204)
  -pathPrefix string
        Root Path Prefix {HERODOTE_PATH_PREFIX}
  -port uint
        [server] Listen port (0 to disable) {HERODOTE_PORT} (default 1080)
  -prometheusAddress string
        [prometheus] Listen address {HERODOTE_PROMETHEUS_ADDRESS}
  -prometheusCert string
        [prometheus] Certificate file {HERODOTE_PROMETHEUS_CERT}
  -prometheusGzip
        [prometheus] Enable gzip compression of metrics output {HERODOTE_PROMETHEUS_GZIP}
  -prometheusIdleTimeout duration
        [prometheus] Idle Timeout {HERODOTE_PROMETHEUS_IDLE_TIMEOUT} (default 10s)
  -prometheusIgnore string
        [prometheus] Ignored path prefixes for metrics, comma separated {HERODOTE_PROMETHEUS_IGNORE}
  -prometheusKey string
        [prometheus] Key file {HERODOTE_PROMETHEUS_KEY}
  -prometheusPort uint
        [prometheus] Listen port (0 to disable) {HERODOTE_PROMETHEUS_PORT} (default 9090)
  -prometheusReadTimeout duration
        [prometheus] Read Timeout {HERODOTE_PROMETHEUS_READ_TIMEOUT} (default 5s)
  -prometheusShutdownTimeout duration
        [prometheus] Shutdown Timeout {HERODOTE_PROMETHEUS_SHUTDOWN_TIMEOUT} (default 5s)
  -prometheusWriteTimeout duration
        [prometheus] Write Timeout {HERODOTE_PROMETHEUS_WRITE_TIMEOUT} (default 10s)
  -publicURL string
        Public URL {HERODOTE_PUBLIC_URL} (default "https://herodote.vibioh.fr")
  -readTimeout duration
        [server] Read Timeout {HERODOTE_READ_TIMEOUT} (default 5s)
  -redisAddress string
        [redis] Redis Address fqdn:port (blank to disable) {HERODOTE_REDIS_ADDRESS} (default "localhost:6379")
  -redisAlias string
        [redis] Connection alias, for metric {HERODOTE_REDIS_ALIAS}
  -redisDatabase int
        [redis] Redis Database {HERODOTE_REDIS_DATABASE}
  -redisPassword string
        [redis] Redis Password, if any {HERODOTE_REDIS_PASSWORD}
  -redisUsername string
        [redis] Redis Username, if any {HERODOTE_REDIS_USERNAME}
  -shutdownTimeout duration
        [server] Shutdown Timeout {HERODOTE_SHUTDOWN_TIMEOUT} (default 10s)
  -title string
        Application title {HERODOTE_TITLE} (default "Herodote")
  -tracerRate string
        [tracer] OpenTracing sample rate, 'always', 'never' or a float value {HERODOTE_TRACER_RATE} (default "always")
  -tracerURL string
        [tracer] OpenTracing gRPC endpoint (e.g. otel-exporter:4317) {HERODOTE_TRACER_URL}
  -url string
        [alcotest] URL to check {HERODOTE_URL}
  -userAgent string
        [alcotest] User-Agent for check {HERODOTE_USER_AGENT} (default "Alcotest")
  -writeTimeout duration
        [server] Write Timeout {HERODOTE_WRITE_TIMEOUT} (default 10s)
```
