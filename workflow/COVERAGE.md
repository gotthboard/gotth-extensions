# Coverage map

| Subsystem | Required harnesses | Current evidence | High-risk gaps |
| --- | --- | --- | --- |
| Manifest/grant validation | unit, table, fuzz | pass: strict/duplicate/depth/bounds/redaction tables plus `FuzzParseManifest` | none |
| Canonical JSON/digests | vectors, non-mutation, coverage | pass: fixed vectors, set ordering, nil/empty equivalence, caller preservation | none |
| Negotiation | unit, hostile, compatibility matrix, fuzz | pass: subset/stale/range/control-only/multiple-interface cases plus `FuzzNegotiate` | none |
| Lifecycle | complete transition matrix, race-safe pure calls | pass: every state pair under unit and race suites | none |
| Protobuf control schema | descriptor compile and forbidden-field scan | pass: descriptor SHA-256 `62a549b29a030f0decb46b2b1c2f382bbbc75430ba0b3d7c729bd21528b505cf` | cross-language conformance deferred until SDK/provider consumer |
| External consumption | clean temporary module compile | pass: nested independent module and final clean clone | real product pin deferred |

Statement coverage is 97.1%. The uncovered statements are defensive failures
after already-validated typed canonicalization, impossible token-shape guards
behind `encoding/json.Decoder.Token`, and byte-classification alternatives.
Every authority, compatibility, canonicalization, lifecycle, bound, duplicate,
redaction, and hostile-depth path is exercised. There is no unresolved critical
or high-risk foundation gap.

Cross-language runtime conformance and a real product/provider pin remain
required before a release; they do not block this unreleased, provider-free
technical admission.
