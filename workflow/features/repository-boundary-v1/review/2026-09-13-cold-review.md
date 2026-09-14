# Cold review — concrete-extension repository boundary

## Verdict

Accept after repair.

## Finding repaired

The first implementation specification invented a 100-byte repository-name
limit without a protocol or hosting contract requiring it, and its wording
could be misread as allowing only one platform artifact or one runtime instance.
That was fake precision. Commit `fa9d2ff` removed the invented limit and made
the actual cardinality explicit: one extension identity per repository, any
number of bound platform artifacts, and any number of independently granted
instances.

## Final hostile pass

- The plural foundation cannot become a provider pack.
- One concrete repository cannot bundle unrelated provider mechanisms.
- Repository names are source-ownership labels, never credentials, grants,
  compatibility results, or product admission.
- Product policy, storage, authorization, mutation, audit, and confirmation
  remain consumer-owned.
- Seam-specific business RPCs do not leak into the handshake/health protocol.
- Independent release and rollback prevent one provider failure from forcing
  unrelated extensions to move together.
- The contract avoids premature creation of provider repositories and shared
  abstractions for imaginary consumers.
- No workflow, trust, confirmation, or runtime userspace changes occurred.

No documented runtime or hosting limit was ignored. The review is confined to
repository topology; concrete provider and product-integration contracts remain
future, separately admitted work.
