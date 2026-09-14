# Product requirements

## Problem

GOTTH applications already need distinct out-of-process mechanisms, including
DNS automation, notification delivery, backup storage, webmail integration,
certificate handling, and import. Leaving each product to invent extension
identity, capability grants, compatibility negotiation, lifecycle, and health
reporting produces incompatible and unsafe plugin systems.

The foundation must be general without becoming a generic execution engine.
The host knows product policy; an extension supplies one bounded mechanism.

## Concrete use cases proving the boundary

These are real consumer inputs, not promises that providers ship here:

1. A DNS provider reads a host-selected zone, previews an exact record plan,
   applies only an explicitly confirmed record set, and verifies the result.
2. A notification provider delivers an already-authorized, privacy-minimized
   alert or approval prompt and returns bounded delivery evidence.
3. A backup provider writes, retrieves, and verifies host-selected immutable
   artifacts without receiving product database access.
4. GOTTH Mail additionally identifies external webmail, certificate/ACME, and
   import-source seams. They confirm the shared control boundary but retain
   separate, seam-specific protocols.

## Requirements

- `EXT-001`: Define a language-neutral, versioned extension manifest with
  stable identity, semantic version, control-protocol ranges, seam-specific
  interface ranges, declared capabilities, and named secret requirements.
- `EXT-002`: Define a host-issued grant bound to one extension identity, one
  manifest digest, one instance, exact interface versions, exact capability
  names, and exact secret identifiers.
- `EXT-003`: Treat declarations as requests only. Negotiation admits only the
  intersection explicitly present in the grant and supported by both host and
  extension.
- `EXT-004`: Canonicalize manifests and grants deterministically and expose a
  SHA-256 digest so configuration, handshake, audit, and rollback can bind the
  same bytes.
- `EXT-005`: Negotiate one control-protocol major and the highest mutually
  supported minor version. Major mismatches fail closed.
- `EXT-006`: Keep business operations out of the control protocol. Handshake
  and health are control-plane functions; DNS, notification, backup, and other
  operations require their own versioned contracts. V1 business interfaces
  are extension-provided and host-called; no host callback is implied.
- `EXT-007`: Run extensions out of process. The foundation must not load
  arbitrary Go plugins, shared objects, scripts, or WASM modules in a host
  process.
- `EXT-008`: Require the deployment/transport boundary to authenticate the
  expected service identity before a control or business RPC is trusted.
  Extension identity text in a response is never authentication.
- `EXT-009`: Require deadlines and bounded correlation identifiers for every
  RPC. Returned messages and errors are stable, bounded, and secret-free.
- `EXT-010`: Model explicit lifecycle transitions: discovered, starting,
  ready, degraded, stopping, stopped, and failed. No state transition launches,
  kills, restarts, or retries a process implicitly.
- `EXT-011`: Keep process supervision, container orchestration, image
  admission, service credential issuance, secret storage/injection, audit
  persistence, and retry policy consumer-owned.
- `EXT-012`: Never place secret values in manifests, grants, control messages,
  digests intended for broad display, logs, or error text. A secret requirement
  names a slot only.
- `EXT-013`: Bound all strings, lists, serialized documents, version numbers,
  JSON nesting, and state-machine work. Reject duplicate and ambiguous
  identifiers and duplicate JSON object names.
- `EXT-014`: Reject unknown schema versions and noncanonical digest encodings.
  Compatibility is explicit; guessing or best-effort fallback is forbidden.
- `EXT-015`: Give extensions no direct product database, Docker socket, broad
  filesystem mount, or authority to decide product authorization/business
  policy.
- `EXT-016`: Publish owner-authored contents under MIT while retaining all
  third-party licenses.
- `EXT-017`: Keep every concrete extension in its own repository named
  `gotth-extension-<slug>`. The plural `gotth-extensions` repository remains
  the provider-free compatibility foundation and must never become a provider
  pack.
- `EXT-018`: Limit each concrete extension repository to one independently
  deployable extension with its own executable or image, exact manifest,
  tests, conformance evidence, security policy, changelog, release/tag
  lifecycle, dependency/license inventory, and rollback procedure.
- `EXT-019`: Require concrete extensions to consume an admitted,
  seam-specific contract. A provider must not create a private incompatible
  variant of an existing seam. Ownership of a new seam contract must be
  decided explicitly when its first real provider is designed; the foundation
  does not absorb speculative business RPCs.
- `EXT-020`: Require every consumer to pin and admit each extension artifact,
  version, manifest digest, grant, and transport identity independently.
  Repository naming, organization membership, or MIT licensing grants no
  runtime authority and proves no compatibility.

## Non-goals

- A concrete GoDaddy, Cloudflare, notification, backup, webmail, ACME, or import
  provider.
- Product HTTP APIs or generated product clients.
- A generic request dispatcher, event bus, workflow engine, scheduler, queue,
  service mesh, secret manager, or container supervisor.
- Direct database access or product-domain models.
- Declarative themes.
- In-process plugins or WASM.
- Live credentials, deployment, DNS mutation, tags, or releases.
- Concrete-extension source trees or creation of any `gotth-extension-*`
  repository in this foundation slice.

## Acceptance

- Requirements trace to architecture, implementation specification, source,
  tests, and retained evidence.
- Manifest/grant validation, canonicalization, digesting, version negotiation,
  capability/interface/secret subset checks, lifecycle transitions, and
  redacted errors have positive, edge, and hostile tests.
- The protobuf control schema compiles to a descriptor set without warnings.
- Format, vet, full tests, race tests, repeated focused tests, statement
  coverage, fuzz smoke, external-consumer compile, and clean-clone checks pass.
- The foundation remains provider-free and performs no external mutation.
- Documentation and distribution rules consistently enforce one concrete
  extension per `gotth-extension-<slug>` repository without claiming that any
  such repository already exists or works with a product.

## Planned host-management contract

The following requirements are planned beyond the admitted V1 foundation and
do not claim current implementation:

- `EXT-MGMT-001`: Define a bounded, versioned configuration-metadata contract
  from which a consumer can render its own native administrator form. The
  contract may describe labels, field types, validation, defaults, and named
  secret slots; it must not carry HTML, JavaScript, templates, CSS, executable
  code, or secret values.
- `EXT-MGMT-002`: Define a secret-free administrator projection for each
  installed extension: repository, immutable artifact/version pin, manifest
  digest, requested and granted capabilities/interfaces/secrets, lifecycle,
  health, enabled state, available update, and rollback pin.
- `EXT-MGMT-003`: Keep install, setup, test, enable, disable, update, rollback,
  uninstall, secret storage, process supervision, audit, and confirmation
  authority in the consuming host. The foundation may validate state and
  metadata but performs none of those actions.
- `EXT-MGMT-004`: Require disable to revoke the effective grant and stop new
  routing before process shutdown. Preserve configuration for re-enable;
  uninstall and secret deletion remain separate confirmed operations.
- `EXT-MGMT-005`: Require updates to preview artifact, manifest, interface,
  capability, configuration, and secret-slot changes before approval and to
  retain a separately admissible rollback pin.
- `EXT-MGMT-006`: Permit Mail, Board, and other consumers to share presentation
  conventions without sharing registries, grants, secrets, audit state, or
  cross-product administrator authority.
