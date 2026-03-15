# AGENTS.md

## Cursor Cloud specific instructions

**drone-zip** is a Go CLI tool (Drone CI / Woodpecker CI plugin) that compresses files into ZIP archives using glob patterns. Single Go module, no external services required.

### Lint / Test / Build

Standard commands from CI (`.github/workflows/ci.yml`):

- **Format check:** `gofmt -s -l .` (expect 0 output)
- **Vet:** `go vet ./...`
- **Test:** `go test -v --count=1`
- **Build:** `go build -o drone-zip-release`

### Non-obvious caveats

- **Test fixtures must be generated before running tests.** Run `sh scripts/generate-test-files.sh` to create the `test/` and `test-output/` directories with fixture files. Without these, all tests will fail. The update script handles this automatically.
- The built binary is a CLI invoked via `--input` and `--output` flags (or `PLUGIN_INPUT` / `PLUGIN_OUTPUT` env vars in CI).
