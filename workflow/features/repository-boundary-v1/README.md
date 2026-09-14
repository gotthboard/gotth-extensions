# One concrete extension per repository

ID: `repository-boundary-v1`

Canonical state lives in `workflow.toml`. This feature records the bounded
repository-topology contract: `gotth-extensions` is the shared foundation, and
each concrete extension belongs in one `gotth-extension-<slug>` repository.

This slice creates no concrete extension repository or runtime integration.

- Evidence: `evidence/2026-09-13-repository-boundary-v1.md`
- Cold review: `review/2026-09-13-cold-review.md`
