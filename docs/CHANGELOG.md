# Changelog

## Unreleased

### 2026-09-13 22:50 CDT — Require one repository per concrete extension

Record Danny's owner requirement that the plural `gotth-extensions`
repository remain the shared provider-free foundation and that every concrete
extension live in its own `gotth-extension-<slug>` repository.

Each concrete repository owns one independently deployable extension, its
manifest, conformance and verification evidence, security and dependency
records, release history, and rollback path. Provider packs are forbidden.
Consumers must continue to pin and authorize every extension independently;
repository naming grants no authority and proves no compatibility.

This contract change creates no provider repository, product integration,
credential, DNS mutation, deployment, tag, release, or GitHub mirror.

### 2026-09-13 22:05 CDT — Define the general extension foundation

Commit: `19855e4`

Affected files:

- repository policy and release documents
- product requirements, architecture, implementation specification, feature
  plan, verification contract, and workflow manifest

Explanation:

Record Danny's explicit decision to create `gotth-extensions` before concrete
providers and before `gotth-sdk`. Define a provider-free, out-of-process,
capability-scoped control and compatibility kernel grounded in GOTTH Mail's
real DNS, notification, backup, webmail, certificate, and import seams.

The host remains authoritative for policy, credentials, state, audit,
confirmation, process supervision, and product operations. The foundation
does not contain a generic invocation API or event bus.

Verification:

- prior GOTTH Extensions queue contract reconciled with the new owner ordering
- current `gotth-sdk`, `gotth-webhooks`, `gotth-authentik`, GOTTH Board, and
  GOTTH Mail contracts inspected
- implementation and runtime verification follow in later commits

Risks / non-goals:

- no provider, SDK, product change, credential, DNS mutation, deployment, tag,
  release, or public compatibility promise

### 2026-09-13 22:26 CDT — Implement and verify the V1 foundation

Source head: `95b2bc1bbe10bd1537fd5f43f5aa20b09388a29e`

Implemented strict manifest and grant parsing, deterministic canonical JSON and
SHA-256 bindings, capability/secret/interface subset negotiation, highest-
compatible-minor control selection, lifecycle validation, the deliberately
small handshake/health protobuf schema, and an external consumer compile.

Cold review repaired two real boundary defects before admission: JSON scanning
now rejects duplicate object names and caps container nesting at 32 levels;
V1 seam interfaces are explicitly extension-provided and host-called so a
grant cannot be misread as callback authority. Validation failures have an
explicit no-input-echo regression test.

Development-host verification passed full tests, race tests, vet, 100 repeated
focused runs, 20 shuffled full runs, three fuzz-smoke targets, protobuf
descriptor compilation and forbidden-field inspection, module verification,
97.1% statement coverage, public external-consumer compilation, clean-clone
verification, and diff checks.

No provider, process runner, secret store, product integration, live DNS,
deployment, tag, release, or GitHub mirror was created.
