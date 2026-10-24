# Coverage map

| Subsystem | Required harnesses | Current evidence | Release gap |
| --- | --- | --- | --- |
| Manifest/grant validation | unit, table, fuzz | retained baseline plus normative text/hex raw-byte vectors, invalid UTF-8 rejection, recursive exact-case key/alias-collision rejection, and duplicate/depth/bounds/redaction tables | exact-candidate full/fuzz rerun pending |
| Canonical JSON/digests | vectors, non-mutation | ordinary and control-only manifest, grant, and complete Session bytes/outputs/digests plus U+2028/U+2029 and surrogate-pair display-name behavior pinned in `compat/v1/canonical-vectors.json`; fixture decoding recursively rejects unknown, duplicate, recased, collided-alias, and casing-typo keys | SDK copy/parity must bind final F1 ref |
| Negotiation | unit, hostile, matrix, fuzz | retained subset/stale/range/control-only/multiple-interface tests plus a manifest-valid but host-range-invalid grant vector, stable-error vectors, actual input/Session slice detachment in both directions, and Session mutation/revalidation | exact-candidate repeat/shuffle/fuzz pending |
| Lifecycle | complete transition matrix | every known-state pair and unknown states pinned | exact-candidate full/race rerun pending |
| Go API | docs, signature/shape baseline | `compat/v1/go-api.txt`, `make api-compat`, compile signatures, reflected fields, and exported docs | independent release review pending |
| Protobuf control schema | descriptor plus scope scan | descriptor SHA-256 `62a549b29a030f0decb46b2b1c2f382bbbc75430ba0b3d7c729bd21528b505cf` pinned | exact-ref SDK admission and structural review pending |
| Cross-language SDK | Go/C++ shared vectors | uncommitted F2 candidate reports Go and protobuf C++ 36.1 passing six handshake/health cases | not admitted, committed, tagged, or published |
| Concrete consumption | clean module plus executable adopter | no current stable-adoption claim; historical observations are not exact release evidence | after SDK and foundation publication, webhook must pin them, pass, be accepted and published; only then may GoDaddy do the same; design-only repositories are off path |
| Release/security | clean clone, dependency/license/SBOM, owner-signed provenance/manifest/checksums, parity | release contract plus fail-fast inventory-validating source-archive script defined | all exact-candidate evidence, signer authority, independent archive equality, and publication pending |

The retained foundation admission at `dfa3d975369c436ffdee98e5f6468eb2be268c54`
reports 97.1% statement coverage and reviewed defensive uncovered branches. That
number was not rerun for this uncommitted F1 patch and is not release-candidate
coverage. F1's fresh local checks are limited to the lightweight commands named
in its evidence.

There is no known critical/high-risk defect in the admitted provider-free
kernel. That statement does **not** mean release risk is clear: exact-candidate
heavy verification, admitted/published SDK and non-Go conformance,
vulnerability/secret scans, dependency/license/SBOM records, signed-tag and
owner-signed post-tag evidence, and forge/release parity remain unresolved
release blockers. Webhook acceptance/publication after foundation publication,
then GoDaddy adoption after published webhook, remain required post-publication
stages, not release blockers. GoDaddy does not gate webhook.
