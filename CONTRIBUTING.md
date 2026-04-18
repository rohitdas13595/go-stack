# Contributing to GoStack

Thank you for your interest in improving GoStack. This document explains how
to participate effectively.

## Code of conduct

All contributors are expected to follow the [Code of Conduct](CODE_OF_CONDUCT.md).
Be respectful, assume good intent, and keep discussion focused on the project.

## Security issues

Please **do not** open a public issue for security vulnerabilities. See
[SECURITY.md](SECURITY.md) for how to report them privately.

## Before you start

- Check [open issues](https://github.com/rohitdas13595/go-stack/issues) and
  existing pull requests to avoid duplicate work.
- For substantial changes, consider opening an issue first to discuss design
  and scope.

## Development setup

```bash
git clone https://github.com/rohitdas13595/go-stack.git
cd go-stack
go test ./...
```

Build all examples (requires the repo `go.work`):

```bash
for d in examples/*/; do (cd "$d" && go build ./...); done
```

## Pull requests

1. **Branch** from `main` or `master` (whichever is the default in this repo).
2. **Keep changes focused** — one logical concern per PR when possible.
3. **Format and test** — run `gofmt` (or `go fmt ./...`) and ensure `go test ./...` passes.
4. **Document behavior** — update `README.md`, `examples/`, or other docs when
   you change user-visible behavior or public APIs.
5. **Commit messages** — use clear, imperative summaries (e.g.
   `fix(router): handle trailing slash`). Link issues when relevant (`Fixes #123`).

## Licensing

By contributing, you agree that your contributions will be licensed under the
same terms as the project: [Apache License 2.0](LICENSE).
