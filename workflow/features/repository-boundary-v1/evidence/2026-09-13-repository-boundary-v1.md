# Concrete-extension repository-boundary evidence — 2026-09-13

## Admitted contract

Contract source head: `fa9d2ffb9971b0ee751023d40806674599793cdb`

- `gotth-extensions` remains the plural, provider-free control and
  compatibility foundation.
- Every concrete extension belongs in one canonical Forgejo repository named
  `gotth-extension-<slug>`.
- A concrete repository owns one independently deployable extension identity,
  its artifacts, exact manifest bindings, tests, conformance evidence,
  security and dependency records, release history, and rollback path.
- Multiple platform artifacts and multiple independently granted runtime
  instances may represent that same extension identity.
- Repository naming grants no capability, transport identity, compatibility,
  product admission, or release trust.
- No concrete extension repository was created by this feature.

## Requirement trace

| Requirement | Contract | Proof |
| --- | --- | --- |
| EXT-017 | one `gotth-extension-<slug>` repository per concrete extension; plural foundation remains provider-free | PRD, README, architecture, contribution and distribution rules |
| EXT-018 | one independently releasable extension identity owns artifacts, manifest bindings, evidence, security, dependency, release, and rollback records | PRD, architecture, implementation and release specifications |
| EXT-019 | providers consume an explicitly owned seam contract; business RPCs stay outside control | PRD, architecture, implementation specification |
| EXT-020 | consumers independently pin artifact/version/digest/grant/transport identity | PRD, README, architecture, release and verification contracts |

## Development-host verification

Host: `development` (`10.0.0.97`), exact contract source head above.

- `make verify`: pass
  - `go vet`, full unit tests, and full race tests
  - statement coverage: 97.1%
  - protobuf descriptor compilation and forbidden-field inspection
  - independent `test/consumer` module
  - `git diff --check`
- `go mod verify`: pass (`all modules verified`)
- workflow event JSON parsing: pass
- required policy present in README, contribution rules, PRD, architecture,
  implementation, feature plan, verification, distribution, and release docs:
  pass
- no top-level concrete `gotth-extension-*` source tree: pass
- clean temporary clone and clean final status: pass

## Scope proof

The source diff contains documentation and workflow records only. It changes no
Go package, protobuf schema, consumer, product, provider, credential, DNS
record, deployment, tag, release, or mirror. The future GoDaddy DNS extension
is merely an example of the naming contract; it was not created.
