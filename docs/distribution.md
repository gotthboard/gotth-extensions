# Distribution

Forgejo is the canonical development home. The public GitHub mirror exists at
`github.com/gotthboard/gotth-extensions`; retained 2026-10-23 audit evidence
found its `main` equal to Forgejo `main` at
`dfa3d975369c436ffdee98e5f6468eb2be268c54`. That historical observation is not
release evidence. No foundation tag, hosted release, or release asset exists.

The exact candidate, signed tag object, peeled commit, release metadata,
attachment inventory, owner-signed provenance/manifest/checksum evidence,
verified signer fingerprints, and every attachment byte must be proved equal on
both forges at publication. The mandatory order is published SDK, published
foundation, accepted and published webhook adopter, then GoDaddy acceptance and
publication. Neither adopter is a foundation pre-tag gate, and GoDaddy never
gates webhook acceptance or publication. Design-only repositories are not part
of this adoption path. This release-contract work does not create or synchronize
any remote object.

The plural `gotth-extensions` repository distributes only the shared
foundation. Every concrete extension is independently admitted and distributed
from its own `gotth-extension-<slug>` repository with its own artifact,
manifest digest, dependency/license inventory, provenance, security record,
tag/release history, and rollback record. A foundation mirror or release never
creates, authorizes, or proves a concrete-extension mirror or release.
