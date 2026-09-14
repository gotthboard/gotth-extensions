# Feature plan

## Foundation V1

1. Repository and controlled documentation.
2. Manifest/grant types, validation, canonicalization, and digests.
3. Version-range negotiation and immutable effective session projection.
4. Lifecycle transition validator.
5. Language-neutral protobuf handshake/health schema.
6. Hostile tests, fuzz smoke, external-consumer compile, coverage, and cold
   review.

## Deferred consumer work

- generated clients in `gotth-sdk`;
- deployment identity and process supervision in consuming products;
- seam-specific DNS, notification, backup, webmail, ACME, and import protocols;
- concrete providers, each in its own `gotth-extension-<slug>` repository;
- administrator screens in the consuming product, not in an extension
  repository.

Deferred items are not silently represented by stubs in this repository.

## Repository-boundary sequencing

1. Admit this provider-free foundation.
2. Define and admit one real consumer seam.
3. Create one `gotth-extension-<slug>` repository for one concrete mechanism.
4. Prove that repository against the foundation, seam contract, and consumer
   without granting authority by name.
5. Repeat independently for another mechanism only when there is a real need.

No provider repository is created merely to reserve a name.

## Host management contract (planned)

1. Prove one real consumer registry and one concrete extension configuration
   need; do not invent a universal form language from screenshots.
2. Define the bounded, code-free configuration metadata and secret-free status
   projection.
3. Define preview/confirmation, enable/disable ordering, update, and rollback
   identities without moving authority out of the host.
4. Add hostile metadata, privilege-diff, stale-preview, redaction, and
   cross-product-isolation conformance tests.
5. Let each product build its own native administrator surface against the
   admitted contract.

The foundation will not contain templ, Tailwind, HTMX, product routes, a secret
store, a process supervisor, an extension catalog service, or a global GOTTH
administrator.
