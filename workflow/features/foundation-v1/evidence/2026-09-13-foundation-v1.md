# General extension foundation evidence — 2026-09-13

## Admitted mechanism

Source head: `95b2bc1bbe10bd1537fd5f43f5aa20b09388a29e`

- `pkg/extensions` is a pure compatibility kernel for strict manifest/grant
  parsing, canonical SHA-256 binding, explicit subset negotiation, immutable-
  by-convention effective sessions, and lifecycle transition validation.
- `proto/gotth/extensions/v1/control.proto` contains only handshake and health.
  Seam-specific product operations remain separate protocols.
- V1 seam interfaces are extension-provided and host-called. Extensions gain no
  implicit host callback, process, storage, product-policy, or secret authority.
- The package imports only Go standard-library byte, JSON, hash, formatting,
  sorting, string, and Unicode helpers. It performs no network, filesystem,
  database, process, plugin-loading, or background work.

## Requirement trace

| Requirement | Mechanism | Proof |
| --- | --- | --- |
| EXT-001–002 | `types.go`, `validate.go` | manifest/grant table and bound tests |
| EXT-003 | `Negotiate` exact intersections | hostile capability/interface/secret tests |
| EXT-004 | canonical private projections and SHA-256 | fixed vectors, ordering, nil/empty, non-mutation tests |
| EXT-005 | `negotiateControl` | overlap, mismatch, and highest-minor tests |
| EXT-006 | handshake/health-only protobuf; explicit V1 direction | descriptor compile and forbidden-field scan |
| EXT-007–009 | out-of-process/transport/deadline contract | architecture and schema review; no runtime executor exists |
| EXT-010 | `ValidateTransition` | complete state-pair matrix under unit/race tests |
| EXT-011–012 | consumer-owned runtime/secrets; IDs only | source/protobuf scan and grant/secret tests |
| EXT-013–014 | fixed sizes/counts/depth; strict schemas/digests | bounds, duplicate-key, depth, SemVer, schema, and digest tests/fuzz |
| EXT-015 | no product/database/Docker/filesystem authority | direct-import and source-scope scan |
| EXT-016 | MIT | root `LICENSE` and distribution contract |

## Development-host verification

Host: `development` (`10.0.0.97`), final source head above.

- `make verify`: pass
  - `gofmt` check, `go vet`, full tests, full race tests
  - statement coverage: 97.1%
  - protobuf descriptor compilation and forbidden-field scan
  - independent `test/consumer` module
  - `git diff --check`
- `go mod verify`: pass (`all modules verified`)
- focused canonicalization/negotiation/lifecycle tests, `-count=100`: pass
- full tests, shuffled, `-count=20`: pass
- five-second fuzz smoke for `FuzzParseManifest`, `FuzzSemver`, and
  `FuzzNegotiate`: pass; a later ten-second parser fuzz run also passed
- clean temporary clone at the exact source head, `make verify`, clean status:
  pass
- control descriptor SHA-256:
  `62a549b29a030f0decb46b2b1c2f382bbbc75430ba0b3d7c729bd21528b505cf`

The 2.9% uncovered statements are defensive error returns after validated,
bounded typed encoding; impossible token kinds behind the standard JSON token
contract; and low-level byte classifier alternatives. Consequential invalid,
duplicate, oversized, deeply nested, stale, incompatible, ungranted, mutation,
and redaction behavior is covered.

## Scope proof

Repository source/protobuf scans found no HTTP client/server, SQL, process
execution, in-process plugin loader, Docker integration, concrete DNS vendor,
generic invocation RPC, arbitrary payload/map, secret value, product mutation,
or runtime side effect. No credential was created or accepted.

This evidence admits only an unreleased foundation. SDK generation, concrete
providers, consumer integration, cross-language runtime conformance, releases,
and deployments remain future work.
