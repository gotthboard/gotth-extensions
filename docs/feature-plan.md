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
- concrete providers and administrator screens.

Deferred items are not silently represented by stubs in this repository.
