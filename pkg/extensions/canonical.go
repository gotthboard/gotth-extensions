package extensions

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

type manifestWire struct {
	Schema       string              `json:"schema"`
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Protocols    []VersionRange      `json:"protocols"`
	Interfaces   []VersionRange      `json:"interfaces"`
	Capabilities []string            `json:"capabilities"`
	Secrets      []SecretRequirement `json:"secrets"`
}

type grantWire struct {
	Schema         string           `json:"schema"`
	InstanceID     string           `json:"instance_id"`
	ExtensionID    string           `json:"extension_id"`
	ManifestDigest string           `json:"manifest_sha256"`
	Capabilities   []string         `json:"capabilities"`
	Interfaces     []InterfaceGrant `json:"interfaces"`
	Secrets        []string         `json:"secrets"`
}

func CanonicalManifest(manifest Manifest) ([]byte, error) {
	if err := ValidateManifest(manifest); err != nil {
		return nil, err
	}
	copyValue := manifestWire{
		Schema:       manifest.Schema,
		ID:           manifest.ID,
		Name:         manifest.Name,
		Version:      manifest.Version,
		Protocols:    cloneRanges(manifest.Protocols),
		Interfaces:   cloneRanges(manifest.Interfaces),
		Capabilities: cloneStrings(manifest.Capabilities),
		Secrets:      cloneSecrets(manifest.Secrets),
	}
	sortRanges(copyValue.Protocols)
	sortRanges(copyValue.Interfaces)
	sort.Strings(copyValue.Capabilities)
	sort.Slice(copyValue.Secrets, func(i, j int) bool { return copyValue.Secrets[i].ID < copyValue.Secrets[j].ID })
	return encodeBounded(copyValue)
}

func ManifestDigest(manifest Manifest) (string, error) {
	canonical, err := CanonicalManifest(manifest)
	if err != nil {
		return "", err
	}
	return digestBytes(canonical), nil
}

func CanonicalGrant(grant Grant) ([]byte, error) {
	if err := ValidateGrant(grant); err != nil {
		return nil, err
	}
	copyValue := grantWire{
		Schema:         grant.Schema,
		InstanceID:     grant.InstanceID,
		ExtensionID:    grant.ExtensionID,
		ManifestDigest: grant.ManifestDigest,
		Capabilities:   cloneStrings(grant.Capabilities),
		Interfaces:     cloneInterfaceGrants(grant.Interfaces),
		Secrets:        cloneStrings(grant.Secrets),
	}
	sort.Strings(copyValue.Capabilities)
	sort.Slice(copyValue.Interfaces, func(i, j int) bool { return copyValue.Interfaces[i].Name < copyValue.Interfaces[j].Name })
	sort.Strings(copyValue.Secrets)
	return encodeBounded(copyValue)
}

func GrantDigest(grant Grant) (string, error) {
	canonical, err := CanonicalGrant(grant)
	if err != nil {
		return "", err
	}
	return digestBytes(canonical), nil
}

func encodeBounded(value any) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil || len(encoded) > MaxDocumentBytes {
		return nil, fieldError("document")
	}
	return encoded, nil
}

func digestBytes(value []byte) string {
	digest := sha256.Sum256(value)
	return hex.EncodeToString(digest[:])
}

func cloneStrings(values []string) []string {
	result := make([]string, len(values))
	copy(result, values)
	return result
}

func cloneRanges(values []VersionRange) []VersionRange {
	result := make([]VersionRange, len(values))
	copy(result, values)
	return result
}

func cloneInterfaceGrants(values []InterfaceGrant) []InterfaceGrant {
	result := make([]InterfaceGrant, len(values))
	copy(result, values)
	return result
}

func cloneSecrets(values []SecretRequirement) []SecretRequirement {
	result := make([]SecretRequirement, len(values))
	copy(result, values)
	return result
}

func sortRanges(values []VersionRange) {
	sort.Slice(values, func(i, j int) bool {
		if values[i].Name != values[j].Name {
			return values[i].Name < values[j].Name
		}
		return values[i].Major < values[j].Major
	})
}
