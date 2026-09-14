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
  operations require their own versioned contracts.
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
