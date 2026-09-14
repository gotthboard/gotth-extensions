# Contributing

Changes proceed in this order: product requirement, architecture,
implementation specification, feature decomposition, implementation,
verification, evidence, review, and admission.

Keep patches bounded. Do not add a concrete provider, product policy, generic
RPC dispatcher, event bus, in-process plugin loader, direct database access,
Docker-socket access, secret store, or process supervisor to the foundation.

Every meaningful change updates `docs/CHANGELOG.md`. Run `make verify` before
requesting review. Do not commit generated scratch or coverage output.
