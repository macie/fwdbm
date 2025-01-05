# fwdbm

[![Go Reference](https://pkg.go.dev/badge/github.com/macie/fwdbm.svg)](https://pkg.go.dev/github.com/macie/fwdbm)
[![Quality check status](https://github.com/macie/fwdbm/actions/workflows/check.yml/badge.svg)](https://github.com/macie/fwdbm/actions/workflows/check.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/macie/fwdbm)](https://goreportcard.com/report/github.com/macie/fwdbm)

**fwdbm** is a command-line tool for running forward-only [database migrations](https://en.wikipedia.org/wiki/Schema_migration).

It is modelled after brilliant Stack Overflow's approach (see: [Stack Overflow: How We Do Deployment - 2016 Edition](https://nickcraver.com/blog/2016/05/03/stack-overflow-how-we-do-deployment-2016-edition/#database-migrations)).

## Installation

Download [latest stable release from GitHub](https://github.com/macie/fwdbm/releases/latest) .

You can also build it manually with commands: `make && make build`.

## Development

Use `make` (GNU or BSD):

- `make` - install development dependencies
- `make check` - runs quality checks and unit tests
- `make build` - compile binary from the latest commit
- `make dist` - compile binaries from the latest commit for all supported OSes
- `make clean` - removes compilation artifacts
- `make cli-release` - tag the latest commit as a new release of CLI
- `make module-release` - tag the latest commit as a new release of Go module
- `make info` - print system info (useful for debugging).

### Versioning

The project contains CLI and Go module which can be developed with different
pace. Commits with versions are tagged with:
- `cli/vYYYY.0M.0D` (_[calendar versioning](https://calver.org/)_) - versions of command-line utility
- `vX.X.X` (_[semantic versioning](https://semver.org/)_) - versions of Go module.

## License

[MIT](./LICENSE) (see: [in Plain English](https://www.tldrlegal.com/license/mit-license)).
