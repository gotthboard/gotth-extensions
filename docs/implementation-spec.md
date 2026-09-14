# Implementation specification

## Public package

The sole initial Go package is `github.com/gotthboard/gotth-extensions/pkg/extensions`.

## Concrete repository contract

A concrete extension repository name must match:

```text
^gotth-extension-[a-z][a-z0-9]*(?:-[a-z0-9]+)*$
```

The repository name is lowercase ASCII and contains exactly one independently
deployable extension identity. Hosting-service length limits remain a
distribution concern and are not invented by this protocol specification. The
manifest `ID` remains the protocol identity defined below; the repository must
document the mapping, but hosts must not derive or trust one from the other.
One repository may publish platform-specific artifacts for that same identity,
and a consumer may run multiple independently granted instances.

Every concrete repository must retain:

- the build definition for its executable or immutable image artifacts;
- one exact extension identity, with every released artifact bound to its
  exact manifest and digest;
- seam-contract and foundation-conformance tests;
- security, dependency/license, verification, changelog, release, and rollback
  records.

Provider bundles are forbidden. Product source, product policy, product
database clients, broad credentials, and copied foundation implementations are
forbidden. A new seam-specific protocol requires an explicit ownership and
versioning decision before provider implementation; it is not added to the
handshake/health control schema for convenience.

### Manifest

```go
type Manifest struct {
    Schema       string
    ID           string
    Name         string
    Version      string
    Protocols    []VersionRange
    Interfaces   []VersionRange
    Capabilities []string
    Secrets      []SecretRequirement
}
```

`Schema` is exactly `gotth.extensions.manifest.v1`. `ID` is a lowercase,
dot-separated identifier with at least three segments. `Name` is bounded UTF-8
display text without control characters. `Version` is an exact SemVer 2.0
version.

`VersionRange` has `Name`, `Major`, `MinMinor`, and `MaxMinor`. Range names are
lowercase dot-separated identifiers. Major is positive, `MinMinor <= MaxMinor`,
and each `(name, major)` pair is unique. A manifest must advertise exactly one
`gotth.extensions.control` major-1 range.

Capabilities use the same token grammar and are unique. A secret requirement
contains only a lowercase dot-separated ID and `Required`; values are
impossible to represent.

### Grant

```go
type Grant struct {
    Schema         string
    InstanceID     string
    ExtensionID    string
    ManifestDigest string
    Capabilities   []string
    Interfaces     []InterfaceGrant
    Secrets        []string
}
```

`Schema` is exactly `gotth.extensions.grant.v1`. Instance IDs use lowercase
UUID text. Manifest digests are exactly 64 lowercase hex characters.
`InterfaceGrant` selects one exact `(name, major, minor)`.
Interface names are unique within a grant. V1 interfaces are extension-
provided and host-called; a grant does not create a callback into the host.

### Host support and negotiation

```go
type HostProfile struct {
    Protocols  []VersionRange
    Interfaces []VersionRange
}

func Negotiate(Manifest, Grant, HostProfile) (Session, error)
```

Negotiation validates all inputs, recomputes the manifest digest, checks exact
identity/digest binding, verifies every granted capability/secret/interface is
declared, verifies every granted interface version is inside both manifest and
host ranges, and selects the highest mutually supported control-protocol minor.

The returned session contains the extension and instance IDs, selected control
version, sorted exact interfaces/capabilities/secrets, both canonical digests,
and a session fingerprint. The fingerprint is SHA-256 over the canonical
session projection. It is configuration identity, not a credential.

### Canonicalization and digest

```go
func CanonicalManifest(Manifest) ([]byte, error)
func ManifestDigest(Manifest) (string, error)
func CanonicalGrant(Grant) ([]byte, error)
func GrantDigest(Grant) (string, error)
```

Functions validate and copy inputs, sort set-like lists, and use compact
`encoding/json` output from private wire structs. They do not mutate caller
data. Encoded output must not exceed 65,536 bytes.

### Lifecycle

```go
type State string
func ValidateTransition(from, to State) error
```

Only the transitions documented in `docs/architecture.md` are accepted. The
function is pure and performs no runtime action.

### Stable errors

Export sentinels for invalid input, unsupported schema, incompatible version,
stale binding, ungranted capability/interface/secret, and invalid transition.
Returned errors may wrap a sentinel with a fixed field category but must not
echo values.

## Bounds

- canonical document: 65,536 bytes;
- identifier/token: 128 bytes;
- display name: 128 UTF-8 bytes;
- semantic version: 64 bytes;
- JSON nesting: 32 container levels;
- protocols: 16;
- interfaces: 64;
- capabilities: 128;
- secret requirements/grants: 64;
- challenge: 32 bytes exactly;
- correlation ID: lowercase UUID text;
- health code: 64-byte lowercase dot token.

Unsigned integers in the public types are bounded to signed 31-bit values to
avoid language interop surprises.

## Protobuf control schema

Package: `gotth.extensions.v1`.

`ExtensionControl.Handshake` request contains correlation ID, 32-byte
challenge, expected extension ID, expected manifest/grant SHA-256 digests, and
expected session fingerprint. Response echoes the challenge and bindings plus
the extension version and negotiated control version.

`ExtensionControl.Health` request contains correlation ID and session
fingerprint. Response contains one enum state and one bounded stable code.

The schema contains no generic payload, arbitrary metadata map, credential,
secret value, configuration mutation, process-control method, or business RPC.
Unknown protobuf fields follow normal protobuf forward-compatibility rules;
semantic admission still depends on negotiated versions.

## Verification contract

- table tests for every validation rule and lifecycle edge;
- canonical order/digest vectors and caller-slice non-mutation tests;
- grant subset, stale binding, major mismatch, minor overlap, and highest-minor
  negotiation tests;
- fuzz targets for semantic version, manifest validation/canonicalization, and
  negotiation panic resistance;
- descriptor compilation and inspection for forbidden generic/secret fields;
- external module compile using only the public package;
- format, vet, unit, race, repeat, coverage, clean-clone, and diff checks.
- repository-policy review proving that the plural foundation remains
  provider-free and that every documented concrete extension uses a separate
  valid `gotth-extension-<slug>` repository identity;
- negative review proving that repository naming is never treated as a grant,
  transport credential, compatibility result, or product admission.
