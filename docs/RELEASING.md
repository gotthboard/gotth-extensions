# Releasing

The repository begins unreleased. A technical foundation commit is not a
compatibility promise.

A first tag requires a real consumer pin, cross-language control-protocol
conformance, exact generated-SDK verification, clean security review, Forgejo
and GitHub distribution parity, immutable release artifacts, and explicit
owner release authority. Until then, use commit pins only.

This foundation is never used to publish concrete provider artifacts. Each
`gotth-extension-<slug>` repository has its own independently reviewed tags,
artifacts, manifest digest, dependency/license inventory, and rollback record.
A consumer admits and pins those releases separately; a foundation release
does not imply that any concrete extension is trusted or compatible.
