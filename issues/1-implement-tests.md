---
type: task
tags: ["testing", "qa"]
---

# Implement Tests

## Context
Ensure the `sopsv` CLI and internal vault operations are thoroughly tested, including unit tests and integration tests with `sops` and `age`.

## Implementation Plan

1. **Test Framework Setup:**
   - Use standard Go testing tools (`testing`, `require` or `assert` from `testify`).
   - Create separate directories/files for unit tests and integration tests.

2. **Unit Tests (`*_test.go`):**
   - Write tests for Viper configuration loading and parsing (`config_test.go`).
   - Write tests for Cobra command flag parsing and routing.
   - Mock the `sops` and `age` interactions to isolate core vault logic (creation, retrieval, list, delete) without touching the filesystem unnecessarily.
   - Write tests for environment variable injection logic in the `exec` command.

3. **Integration Tests (`// +build integration`):**
   - Isolate these tests using the `// +build integration` directive and run them with `go test -tags=integration ./...`.
   - Setup temporary directories (e.g., `t.TempDir()`) to simulate OS-specific paths (`~/.local/share/sopsv` and `~/.config/sopsv`).
   - Verify the end-to-end file creation and true encryption/decryption flows using actual `age` key generation and `sops` operations.
   - Verify that the `exec` command correctly injects environment variables into an actual child process (e.g., executing `env` or a small shell script and validating its output).

4. **Continuous Testing:**
   - Run tests as part of the core development loop before completing any tasks.
   - Fix broken tests immediately.
   - Run `govulncheck` to ensure no vulnerabilities in dependencies are introduced.
