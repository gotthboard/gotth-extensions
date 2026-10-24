# Foundation V1 compatibility and release contract

ID: `foundation-v1-release-contract`

This active feature prepares the proposed stable `v1.0.0` compatibility and
release contract. It adds mechanical API, grammar, raw-byte, control-only
Session, slice-detachment, and descriptor drift gates. It defines the signed
post-tag payload, deterministic archive, generated-tree hash,
dependency/license/SBOM records, provenance, the exact SDK publication →
foundation publication → webhook acceptance/publication → GoDaddy adoption
order, and consumer-first rollback procedure. Design-only repositories are not
on that adoption path.

It does not create a candidate commit, SDK, management API, provider, tag,
release, deployment, or support promise. Canonical state is in `workflow.toml`.
Task-local check logs are noncanonical evidence and remain outside this source
tree.
