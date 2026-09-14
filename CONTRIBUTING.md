# Contributing

Changes proceed in this order: product requirement, architecture,
implementation specification, feature decomposition, implementation,
verification, evidence, review, and admission.

Keep patches bounded. Do not add a concrete provider, product policy, generic
RPC dispatcher, event bus, in-process plugin loader, direct database access,
Docker-socket access, secret store, or process supervisor to the foundation.

Each concrete extension must be developed in a separate repository named
`gotth-extension-<slug>`. One repository contains one independently deployable
extension. Do not create provider packs or copy shared foundation code into
extension repositories. Promote mechanism-neutral code only after at least two
real consumers prove the common contract and it receives separate review.

Every meaningful change updates `docs/CHANGELOG.md`. Run `make verify` before
requesting review. Do not commit generated scratch or coverage output.
