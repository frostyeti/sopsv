---
tags: ["chore"]
type: chore
---

# Setup GitHub Actions

## Context
Automate CI/CD pipelines for linting, testing, and releasing the `sopsv` CLI across multiple platforms.

## Implementation Plan

1. **Continuous Integration (`.github/workflows/ci.yml`):**
   - Trigger on push to `main` and Pull Requests.
   - Setup Go environment using `actions/setup-go`.
   - Run `golangci-lint` to keep the code clean and idiomatic.
   - Run Unit tests (`go test ./...`) and Integration tests (`go test -tags=integration ./...`).
   - Run `govulncheck` to catch any vulnerabilities.
   - (Commented out) Run code coverage generation and upload step (until Codecov or another tool is configured).

2. **Continuous Deployment / Release (`.github/workflows/release.yml`):**
   - Trigger on tags (e.g., `v*`).
   - Setup Go environment using `actions/setup-go`.
   - Use non-premium features from `goreleaser/goreleaser-action` to automate releases and versioning with `args: release --clean`.
   - A release should automatically create an announcement on GitHub.

3. **GoReleaser Configuration (`.goreleaser.yaml`):**
   - Configure cross-compilation for `linux`, `windows`, and `macos` (darwin) across `amd64` and `arm64` architectures.
   - Avoid dependencies that require CGO to ensure statically linked binaries cross-compile easily.
   - Configure package managers to ensure availability in Chocolatey, Homebrew, RPM, DEB, Snap, Flatpak, and AppImage.
   - Store artifacts in the `.artifacts` folder locally if running locally.

4. **Docker Support:**
   - (Commented out) Add steps/configuration to build and push Docker images for `sopsv` until a Docker Hub account is established.
