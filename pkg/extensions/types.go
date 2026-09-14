package extensions

const (
	ManifestSchema = "gotth.extensions.manifest.v1"
	GrantSchema    = "gotth.extensions.grant.v1"
	ControlName    = "gotth.extensions.control"

	MaxDocumentBytes = 65_536
	MaxIdentifier    = 128
	MaxDisplayName   = 128
	MaxVersionBytes  = 64
	MaxProtocols     = 16
	MaxInterfaces    = 64
	MaxCapabilities  = 128
	MaxSecrets       = 64
	maxVersionNumber = 1<<31 - 1
)

type VersionRange struct {
	Name     string `json:"name"`
	Major    uint32 `json:"major"`
	MinMinor uint32 `json:"min_minor"`
	MaxMinor uint32 `json:"max_minor"`
}

type InterfaceGrant struct {
	Name  string `json:"name"`
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
}

type SecretRequirement struct {
	ID       string `json:"id"`
	Required bool   `json:"required"`
}

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

type Grant struct {
	Schema         string           `json:"schema"`
	InstanceID     string           `json:"instance_id"`
	ExtensionID    string           `json:"extension_id"`
	ManifestDigest string           `json:"manifest_sha256"`
	Capabilities   []string         `json:"capabilities"`
	Interfaces     []InterfaceGrant `json:"interfaces"`
	Secrets        []string         `json:"secrets"`
}

type HostProfile struct {
	Protocols  []VersionRange
	Interfaces []VersionRange
}

type ProtocolVersion struct {
	Name  string `json:"name"`
	Major uint32 `json:"major"`
	Minor uint32 `json:"minor"`
}

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
