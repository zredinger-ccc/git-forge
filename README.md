# git-forge

A lightweight, self-hostable git forge with first-class [Jujutsu](https://github.com/jj-vcs/jj) support and an agent-friendly API.

See [`docs/design/architecture.md`](docs/design/architecture.md) for the design and scope.

## Status

Pre-MVP. Foundational subissues for the implementation are tracked under issue #1.

## Quick start (development)

```sh
make build      # builds bin/git-forge
make run        # build + run on :8080
make test       # run unit tests
make check      # fmt + vet + test
```

Hit the health endpoint:

```sh
curl http://localhost:8080/healthz
```

## Configuration

All settings default to sane values. Override via environment, or point at a JSON file via `GIT_FORGE_CONFIG`.

| Env var                        | Default     | Notes                                  |
|--------------------------------|-------------|----------------------------------------|
| `GIT_FORGE_ADDR`               | `:8080`     | Listen address.                        |
| `GIT_FORGE_SHUTDOWN_TIMEOUT`   | `15`        | Graceful shutdown, seconds.            |
| `GIT_FORGE_LOG_LEVEL`          | `info`      | `debug` / `info` / `warn` / `error`.   |
| `GIT_FORGE_LOG_FORMAT`         | `json`      | `json` or `text`.                      |
| `GIT_FORGE_CONFIG`             | _(unset)_   | Path to JSON config file.              |

Environment overrides take precedence over the file.
