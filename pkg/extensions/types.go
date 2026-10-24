package extensions

const (
	// ManifestSchema is the only manifest schema accepted by V1 parsers and validators.
	ManifestSchema = "gotth.extensions.manifest.v1"
	// GrantSchema is the only grant schema accepted by V1 parsers and validators.
	GrantSchema = "gotth.extensions.grant.v1"
	// ControlName identifies the provider-free handshake and health protocol.
	ControlName = "gotth.extensions.control"

	// MaxDocumentBytes is the maximum accepted raw or canonical JSON size in
	// bytes.
	MaxDocumentBytes = 65_536
	// MaxIdentifier is the maximum byte length of a dotted identifier.
	MaxIdentifier = 128
	// MaxDisplayName is the maximum UTF-8 byte length of a display name.
	MaxDisplayName = 128
	// MaxVersionBytes is the maximum byte length of a semantic version.
	MaxVersionBytes = 64
	// MaxJSONDepth is the maximum number of nested JSON containers.
	MaxJSONDepth = 32
	// MaxProtocols is the maximum number of manifest or host protocol ranges.
	MaxProtocols = 16
	// MaxInterfaces is the maximum number of interface ranges or grants.
	MaxInterfaces = 64
	// MaxCapabilities is the maximum number of capability identifiers.
	MaxCapabilities = 128
	// MaxSecrets is the maximum number of secret-slot identifiers.
	MaxSecrets       = 64
	maxVersionNumber = 1<<31 - 1
)

// VersionRange declares an inclusive minor-version range for one name and
// major version.
type VersionRange struct {
	Name     string `json:"name"`
	Major    uint32 `json:"major"`
	MinMinor uint32 `json:"min_minor"`
	MaxMinor uint32 `json:"max_minor"`
}

// InterfaceGrant selects one exact interface major and minor version.
type InterfaceGrant struct {
	Name  string `json:"name"`
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
}

// SecretRequirement declares a secret slot by identifier, never a secret value.
type SecretRequirement struct {
	ID       string `json:"id"`
	Required bool   `json:"required"`
}

// Manifest declares an extension identity and its requested protocol, interface,
// capability, and secret-slot surface.
type Manifest struct {
	Schema       string              `json:"schema"`
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Protocols    []VersionRange      `json:"protocols"`
	Interfaces   []VersionRange      `json:"interfaces"`
	Capabilities []string            `json:"capabilities"`
	Secrets      []SecretRequirement `json:"secrets"`
}

// Grant binds one extension instance to a manifest digest and an admitted subset.
type Grant struct {
	Schema         string           `json:"schema"`
	InstanceID     string           `json:"instance_id"`
	ExtensionID    string           `json:"extension_id"`
	ManifestDigest string           `json:"manifest_sha256"`
	Capabilities   []string         `json:"capabilities"`
	Interfaces     []InterfaceGrant `json:"interfaces"`
	Secrets        []string         `json:"secrets"`
}

// HostProfile declares the protocol and interface ranges supported by a host.
type HostProfile struct {
	Protocols  []VersionRange
	Interfaces []VersionRange
}

// ProtocolVersion is one negotiated protocol name, major, and minor version.
type ProtocolVersion struct {
	Name  string `json:"name"`
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
}

// Session is the complete negotiated identity and authority projection.
//
// Session contains caller-owned mutable slices. Negotiate returns those slices
// detached from its inputs. Callers crossing a trust or ownership boundary must
// deep-copy all slices. Before using a retained Session for authority, rerun
// Negotiate from the authoritative Manifest, Grant, and HostProfile and compare
// the complete Session value; the cached Fingerprint alone cannot detect caller
// mutation.
type Session struct {
	ExtensionID      string
	ExtensionVersion string
	InstanceID       string
	Control          ProtocolVersion
	Interfaces       []InterfaceGrant
	Capabilities     []string
	Secrets          []string
	ManifestDigest   string
	GrantDigest      string
	Fingerprint      string
}
