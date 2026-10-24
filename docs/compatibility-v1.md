# Foundation V1 compatibility contract

This contract defines the proposed stable surface for `v1.0.0`. It takes effect
only when an authorized `v1.0.0` tag and matching release are published. The
repository is still unreleased.

## Candidate identity and scope

The selected implementation baseline is commit
`dfa3d975369c436ffdee98e5f6468eb2be268c54`. The release candidate must be an
admitted descendant limited to this release contract, compatibility gates, and
exact-ref evidence. Its full commit ID is recorded only after this uncommitted
preparation is admitted. A tag never targets a worktree or mutable branch name.

Foundation V1 contains only the public Go package, strict manifest/grant parsing
and validation, canonical bytes and digests, subset/version negotiation,
lifecycle validation, and the `gotth.extensions.v1.ExtensionControl`
handshake/health schema. It excludes host management, business protocols,
concrete extensions/providers, process supervision, deployment, product APIs,
secret storage, and generated SDK implementation.

## Raw JSON entry boundary

Untrusted manifest and grant bytes **must** enter through `ParseManifest` and
`ParseGrant`. Ordinary `json.Unmarshal`, another permissive decoder, or
manual construction followed by `Validate*` is not an equivalent raw-input
path: permissive decoding can discard unknown members before validation and
canonical hashing. `ValidateManifest`, `ValidateGrant`, `Canonical*`, and
`Negotiate` accept only trusted in-process typed values.

Both parsers apply this sequence:

1. reject zero bytes and input longer than 65,536 bytes;
2. reject any byte sequence that is not valid UTF-8, before JSON tokenization;
3. tokenize the complete JSON value, require every object member name to match
   its declared lowercase JSON name exactly, reject duplicate object member
   names at every depth, and reject more than 32 nested object/array containers;
4. decode exactly one top-level value with unknown struct members forbidden;
5. reject every trailing non-whitespace value; and
6. apply the typed validator.

Whitespace is accepted wherever Go `encoding/json` accepts it and counts toward
the byte limit. JSON object member order is irrelevant on input. Object names
are case-sensitive: `Schema`, `SCHEMA`, nested aliases such as `Name`, casing
typos, and an alias collision such as `schema` plus `Schema` are invalid. Exact
duplicate names are forbidden even when their values are identical. This key
preflight does not replace Go `encoding/json` value decoding; number, string,
null, array, and object value semantics remain those of that decoder before
typed validation. The parsers return only stable error categories and never
echo rejected values.
In the normative raw-parser vector, each case sets exactly one of `document`
(JSON text encoded as UTF-8) or `document_hex` (lowercase hexadecimal exact
bytes). This permits invalid UTF-8 to be represented without a lossy JSON string.
The foundation test decodes the vector file itself with exact, case-sensitive
object names and duplicate names forbidden recursively at every object level,
unknown fields forbidden, and exactly one top-level JSON value required. The
fixture's own normative decoder-negative cases pin unknown names, exact
duplicates, case aliases, alias collisions, and casing typos from the top level
through nested range, host, Session, and case objects. A misspelled, recased, or
misplaced fixture member therefore cannot be silently discarded.

## Normative typed grammar and bounds

All string limits are UTF-8 byte counts, not rune counts.

- A dotted identifier is 1..128 bytes and has the required segment count stated
  below. Each segment is 1..63 ASCII bytes, begins with `a`..`z`, contains
  only `a`..`z`, `0`..`9`, or `-`, and does not end with `-`.
- Manifest IDs and grant extension IDs require at least three segments.
  Protocol, interface, capability, and secret-slot IDs require at least two.
- A display name is valid UTF-8, 1..128 bytes, has no leading or trailing Unicode
  whitespace, and contains no Unicode control rune.
- A version is 1..64 bytes and follows this bounded subset of SemVer 2.0
  syntax: three dot-separated
  decimal core components, each 0..2^31-1 with no leading zero except `0`;
  optional dot-separated prerelease identifiers after `-`; and optional
  dot-separated build identifiers after `+`. Identifiers contain only ASCII
  alphanumeric or `-`; empty identifiers are forbidden; numeric prerelease
  identifiers cannot have a leading zero. Build identifiers may.
- An instance ID is exactly 36 lowercase ASCII bytes in `8-4-4-4-12`
  hexadecimal form. V1 validates this lexical UUID shape, not UUID version or
  variant bits.
- A digest is exactly 64 lowercase hexadecimal ASCII bytes.
- Version major values are 1..2^31-1. Minor endpoints are 0..2^31-1 and an
  inclusive range requires `min_minor <= max_minor`.

A Manifest has exactly the JSON member names represented by its public fields.
`schema`, `id`, `name`, and `version` must validate. `protocols` has
1..16 entries and contains exactly one `gotth.extensions.control` major-1
range. `interfaces` has 0..64 entries, `capabilities` 0..128 identifiers,
and `secrets` 0..64 requirements. Protocol/interface ranges are unique by
(name, major); capabilities and secret IDs are unique by exact string.

A Grant has exactly its public JSON members. It has 0..128 capabilities, 0..64
interfaces, and 0..64 secret IDs. Grant interfaces are unique by name (therefore
a name cannot appear under two majors); capability and secret IDs are unique by
exact string. A HostProfile has the same protocol/interface range rules as a
Manifest and must contain exactly one control major-1 range.

Missing JSON list members decode as nil. For typed V1 semantics, nil and empty
optional lists are equivalent and canonicalize to `[]`. Unknown fields, extra
values, duplicate list identities, out-of-range numbers, and malformed strings
are rejected as specified above.

## Go public compatibility

The module uses semantic import versioning. While the module remains at major
version 1, exported names retain their declared kind, type/signature, constant
value, struct field order/type/tag, and documented meaning. The normative,
machine-diffable surface is `compat/v1/go-api.txt`, regenerated only by:

`go doc -all ./pkg/extensions > compat/v1/go-api.txt`

`make api-compat` fails on drift. The surface includes all constants,
sentinels, types and fields, parsers, validators, canonical/digest functions,
`Negotiate`, all lifecycle state constants, and `ValidateTransition`.
`pkg/extensions/compatibility_v1_test.go` independently compile-checks function
signatures and reflects every public struct shape.

Minor or patch releases may add separately named APIs, but cannot add fields to
existing V1 structs, change accepted V1 grammar/bounds, reinterpret a grant or
session field, weaken fail-closed behavior, or change a valid input's canonical
bytes, digest, negotiation output, or lifecycle result. A security correction
that cannot preserve this contract requires an explicit exception and a new
major/schema/protocol version.

## Canonical bytes and digests

`CanonicalManifest` and `CanonicalGrant` validate first and deep-copy every
caller slice before sorting. Manifest protocol/interface ranges sort by name,
then major; capabilities sort by string; secret requirements sort by ID. Grant
capabilities and secrets sort by string; interface grants sort by name. Ordering
is ascending bytewise UTF-8 as implemented by Go string comparison.

Manifest member order is `schema`, `id`, `name`, `version`, `protocols`,
`interfaces`, `capabilities`, `secrets`. Grant member order is `schema`,
`instance_id`, `extension_id`, `manifest_sha256`, `capabilities`,
`interfaces`, `secrets`. Nested fields retain declaration order. Encoding is
compact Go `encoding/json` with default HTML escaping, no indentation, no
trailing newline, and no insignificant whitespace. The normative display-name
cases pin Go's escaping of U+2028 and U+2029 as `\u2028` and `\u2029`, and pin
decoding of a valid UTF-16 surrogate-pair JSON escape followed by canonical
UTF-8 emission of the resulting supplementary rune.

Manifest and grant digests are lowercase SHA-256 over exactly those bytes. The
normative `compat/v1/canonical-vectors.json` fixes complete typed inputs,
canonical manifest/grant/session bytes, all three hashes, complete Session
output, strict raw-parser outcomes, negotiation outcomes, and stable error
categories. It is consumed by foundation tests; the SDK candidate copies the
three canonical byte streams to `conformance/canonical/v1/{manifest,grant,session}.json`.

## Negotiation and Session

`Negotiate` validates Manifest, Grant, and HostProfile in that order; computes
the canonical manifest digest; requires exact extension ID and digest binding;
selects the highest minor in the overlap of the single control major-1 ranges;
then verifies capabilities, secrets, and interfaces. Granted capabilities and
secret slots must be declared. Every required manifest secret must be granted.
Every interface grant must fall inside both the manifest and host range for its
(name, major). Empty interface/capability/secret sets are allowed. A normative
negative case fixes the distinct host-boundary failure where the granted minor
is valid for the manifest but outside the authoritative host range.

The returned Session owns newly allocated interface, capability, and secret
slices detached from the Grant input in both directions. Every empty Session
set is a non-nil empty slice and encodes as `[]`, never `null`. Interfaces sort
by name; capabilities and secrets sort by string. Complete ordinary and
control-only successful outputs are pinned in the normative vector. Its
Fingerprint is lowercase SHA-256 over compact JSON with
members, in order: `schema` (`gotth.extensions.session.v1`), `extension_id`,
`extension_version`, `instance_id`, `control`, `interfaces`,
`capabilities`, `secrets`, `manifest_sha256`, and `grant_sha256`.

Session exposes mutable slices and a cached Fingerprint. A caller that transfers
or stores a Session must deep-copy
`Interfaces`, `Capabilities`, and `Secrets`. Before a retained Session is
used for authorization or crosses a trust boundary, the caller must rerun
`Negotiate` from authoritative Manifest, Grant, and HostProfile values and
compare the **complete** Session, including every slice and digest. Comparing
only Fingerprint is forbidden because mutation does not refresh the cache. The
compatibility test mutates a retained copy and proves complete revalidation
detects the stale value.

## Stable error categories

Callers classify with `errors.Is`, never complete wrapped text.

| Sentinel | Stable text | Meaning |
| --- | --- | --- |
| `ErrInvalidInput` | `extensions: invalid input` | malformed, duplicate, or out-of-bound input |
| `ErrUnsupportedSchema` | `extensions: unsupported schema` | non-V1 manifest/grant schema |
| `ErrIncompatibleVersion` | `extensions: incompatible version` | missing required control declaration or overlap |
| `ErrStaleBinding` | `extensions: stale binding` | grant identity/digest mismatch |
| `ErrUngrantedCapability` | `extensions: ungranted capability` | capability outside declaration |
| `ErrUngrantedInterface` | `extensions: ungranted interface` | interface outside declaration/support |
| `ErrUngrantedSecret` | `extensions: ungranted secret` | unknown or missing secret slot |
| `ErrInvalidTransition` | `extensions: invalid transition` | unknown state or forbidden edge |

Fixed, secret-free field-category wrapping of `ErrInvalidInput` is allowed.
The parser and negotiation cases in the normative vector freeze representative
category outcomes and the validation/negotiation order above defines precedence.

## Lifecycle edges

The complete allowed V1 edge set is:

```text
discovered -> starting
starting   -> ready | failed | stopping
ready      -> degraded | stopping | failed
degraded   -> ready | stopping | failed
stopping   -> stopped | failed
stopped    -> starting
failed     -> starting | stopped
```

Every other pair, same-state pair, or unknown state returns
`ErrInvalidTransition`. Validation is pure and never starts or stops a process.

## Protobuf compatibility

The package is `gotth.extensions.v1`; the only service is
`ExtensionControl`; its only unary methods are `Handshake` then `Health`.
The V1 messages, field numbers/types/cardinality, enum names/numbers, RPC names,
paths, and streaming flags are those in `control.proto` and the complete
descriptor pinned by `compat/v1/control-descriptor.sha256`.

A successor schema reserves every removed message-field number and name and
every removed enum number and name. Protobuf syntax has no RPC-level `reserved`
declaration: removed RPC names and method paths are instead permanently
non-reusable by this contract and the retained descriptor/API baseline. V1 does
not remove or renumber anything. Any new control method or field requires
explicit compatibility review and a new protocol version; business and
management RPCs remain forbidden. Unknown protobuf fields retain protobuf's
normal wire behavior but grant no semantic compatibility or authority.
