---
type: chore
tags: ["chore", "docs"]
---

# Create Docs

## Context
Create a documentation site for `sopsv` (referred to as `cast` in early drafts, but adapted to `sopsv`) that shows how to use it and provide a lot of YAML examples for common use cases, including dealing with secrets, using remote modules, and using remote tasks (if applicable).

It should also show all the available commands and how to use them, including updating remote tasks and modules, and purging Docker images (if integrating `sopsv` with such systems).

The docs site should be in the `docs` directory. (You may need to save all the docs in there to a different directory, remove it, and then add the site).

Use Astro, Starlight, ensure that search and RSS is enabled, and that it can handle new docs for new versions of the CLI.

Ideally, it should be published to Cloudflare Pages or something similar so that it can be easily accessed by users and updated with new versions.

## Implementation Plan

1. **Framework Setup:**
   - Initialize an Astro project using the Starlight template within the `docs/` directory (`npm create astro@latest -- --template starlight`).
   - Configure Starlight with search enabled and RSS feed generation.

2. **Content Strategy:**
   - Move or rewrite any existing Markdown files (e.g., from `docs/prd`) into `docs/src/content/docs`.
   - Add "Getting Started" guide explaining `sops` and `age` prerequisites and automatic setup.
   - Document the `sopsv vaults` subcommand group (new, use, list, show, delete).
   - Document the `sopsv exec` command with examples showing how multiple stores provide environment variables to a child process.
   - Document common YAML configurations (local paths vs user-profile paths).

3. **Version Support:**
   - Organize Starlight documentation structure to gracefully handle versioned documentation if necessary for new CLI versions.

4. **Deployment:**
   - Create a GitHub Actions workflow (`.github/workflows/deploy-docs.yml`) configured to deploy the `docs/` Astro site to Cloudflare Pages (using Wrangler action or Cloudflare's direct GitHub integration) on push to the `main` branch.
