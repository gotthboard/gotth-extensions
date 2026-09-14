# Coverage map

| Subsystem | Required harnesses | Current evidence | High-risk gaps |
| --- | --- | --- | --- |
| Manifest/grant validation | unit, table, fuzz | pending | none accepted |
| Canonical JSON/digests | vectors, non-mutation, coverage | pending | none accepted |
| Negotiation | unit, hostile, compatibility matrix, fuzz | pending | none accepted |
| Lifecycle | complete transition matrix, race-safe pure calls | pending | none accepted |
| Protobuf control schema | descriptor compile and forbidden-field scan | pending | cross-language conformance deferred until SDK/provider consumer |
| External consumption | clean temporary module compile | pending | real product pin deferred |

Unresolved critical/high-risk gaps block foundation completion. Cross-language
runtime and real-product validation block a release, not this unreleased
provider-free technical admission.
