# Cold review — general extension foundation

Verdict: CLEAN after repairs.

## Findings repaired

1. The first duplicate-key scanner bounded bytes but not recursive JSON depth.
   That left a stack-exhaustion path inside the manifest trust boundary.
   Commit `3a5accd` added a hard 32-container limit and boundary tests.
2. Interface direction was implicit. A general grant cannot leave callback
   authority to interpretation. Commit `f836ac1` now defines V1 interfaces as
   extension-provided and host-called, with callbacks requiring a future
   explicit contract and grant.
3. Redacted errors were structurally implemented but not directly proved.
   Commit `95b2bc1` added marker-based no-echo regression tests.

## Final hostile pass

- No generic `Invoke`, event bus, workflow engine, provider API, runner,
  supervisor, secret store, product client, direct database access, Docker
  authority, or in-process loader entered the repository.
- Manifest declarations remain requests. Only a host-issued, digest-bound
  subset can enter the effective session.
- Self-asserted extension text is not transport authentication.
- Canonicalization is bounded, deterministic, non-mutating, duplicate-safe,
  nil/empty stable, and tied to fixed vectors.
- Failure categories are stable and do not echo input.
- No userspace exists yet to regress; no compatibility promise is claimed.

The remaining deferred work is correctly outside this slice: generated SDKs,
seam protocols, providers, real consumer pins, cross-language runtime proof,
deployment, tags, and releases.
