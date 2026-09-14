package extensions

import "sort"

func Negotiate(manifest Manifest, grant Grant, host HostProfile) (Session, error) {
	if err := ValidateManifest(manifest); err != nil {
		return Session{}, err
	}
	if err := ValidateGrant(grant); err != nil {
		return Session{}, err
	}
	if err := validateHostProfile(host); err != nil {
		return Session{}, err
	}
	manifestDigest, err := ManifestDigest(manifest)
	if err != nil {
		return Session{}, err
	}
	if manifest.ID != grant.ExtensionID || manifestDigest != grant.ManifestDigest {
		return Session{}, ErrStaleBinding
	}

	control, ok := negotiateControl(manifest.Protocols, host.Protocols)
	if !ok {
		return Session{}, ErrIncompatibleVersion
	}
	if err := verifyCapabilities(manifest.Capabilities, grant.Capabilities); err != nil {
		return Session{}, err
	}
	if err := verifySecrets(manifest.Secrets, grant.Secrets); err != nil {
		return Session{}, err
	}
	if err := verifyInterfaces(manifest.Interfaces, host.Interfaces, grant.Interfaces); err != nil {
		return Session{}, err
	}

	grantDigest, err := GrantDigest(grant)
	if err != nil {
		return Session{}, err
	}
	session := Session{
		ExtensionID:      manifest.ID,
		ExtensionVersion: manifest.Version,
		InstanceID:       grant.InstanceID,
		Control:          control,
		Interfaces:       append([]InterfaceGrant(nil), grant.Interfaces...),
		Capabilities:     cloneStrings(grant.Capabilities),
		Secrets:          cloneStrings(grant.Secrets),
		ManifestDigest:   manifestDigest,
		GrantDigest:      grantDigest,
	}
	sort.Slice(session.Interfaces, func(i, j int) bool { return session.Interfaces[i].Name < session.Interfaces[j].Name })
	sort.Strings(session.Capabilities)
	sort.Strings(session.Secrets)
	fingerprint, err := sessionFingerprint(session)
	if err != nil {
		return Session{}, err
	}
	session.Fingerprint = fingerprint
	return session, nil
}

func negotiateControl(extension, host []VersionRange) (ProtocolVersion, bool) {
	var extensionRange VersionRange
	found := false
	for _, item := range extension {
		if item.Name == ControlName && item.Major == 1 {
			extensionRange = item
			found = true
			break
		}
	}
	if !found {
		return ProtocolVersion{}, false
	}
	for _, item := range host {
		if item.Name != ControlName || item.Major != extensionRange.Major {
			continue
		}
		minimum := max(extensionRange.MinMinor, item.MinMinor)
		maximum := min(extensionRange.MaxMinor, item.MaxMinor)
		if minimum <= maximum {
			return ProtocolVersion{Name: ControlName, Major: 1, Minor: maximum}, true
		}
	}
	return ProtocolVersion{}, false
}

func verifyCapabilities(declared, granted []string) error {
	available := stringSet(declared)
	for _, item := range granted {
		if _, ok := available[item]; !ok {
			return ErrUngrantedCapability
		}
	}
	return nil
}

func verifySecrets(declared []SecretRequirement, granted []string) error {
	available := make(map[string]SecretRequirement, len(declared))
	for _, item := range declared {
		available[item.ID] = item
	}
	selected := stringSet(granted)
	for _, item := range granted {
		if _, ok := available[item]; !ok {
			return ErrUngrantedSecret
		}
	}
	for _, item := range declared {
		if item.Required {
			if _, ok := selected[item.ID]; !ok {
				return ErrUngrantedSecret
			}
		}
	}
	return nil
}

func verifyInterfaces(declared, supported []VersionRange, granted []InterfaceGrant) error {
	declaredIndex := rangeIndex(declared)
	supportedIndex := rangeIndex(supported)
	for _, item := range granted {
		key := rangeKey(item.Name, item.Major)
		left, leftOK := declaredIndex[key]
		right, rightOK := supportedIndex[key]
		if !leftOK || !rightOK || item.Minor < left.MinMinor || item.Minor > left.MaxMinor || item.Minor < right.MinMinor || item.Minor > right.MaxMinor {
			return ErrUngrantedInterface
		}
	}
	return nil
}

func rangeIndex(values []VersionRange) map[string]VersionRange {
	result := make(map[string]VersionRange, len(values))
	for _, item := range values {
		result[rangeKey(item.Name, item.Major)] = item
	}
	return result
}

func rangeKey(name string, major uint32) string {
	return name + "\x00" + uint32Decimal(major)
}

func uint32Decimal(value uint32) string {
	if value == 0 {
		return "0"
	}
	var scratch [10]byte
	position := len(scratch)
	for value > 0 {
		position--
		scratch[position] = byte('0' + value%10)
		value /= 10
	}
	return string(scratch[position:])
}

func stringSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func sessionFingerprint(session Session) (string, error) {
	type sessionWire struct {
		Schema           string           `json:"schema"`
		ExtensionID      string           `json:"extension_id"`
		ExtensionVersion string           `json:"extension_version"`
		InstanceID       string           `json:"instance_id"`
		Control          ProtocolVersion  `json:"control"`
		Interfaces       []InterfaceGrant `json:"interfaces"`
		Capabilities     []string         `json:"capabilities"`
		Secrets          []string         `json:"secrets"`
		ManifestDigest   string           `json:"manifest_sha256"`
		GrantDigest      string           `json:"grant_sha256"`
	}
	encoded, err := encodeBounded(sessionWire{
		Schema:           "gotth.extensions.session.v1",
		ExtensionID:      session.ExtensionID,
		ExtensionVersion: session.ExtensionVersion,
		InstanceID:       session.InstanceID,
		Control:          session.Control,
		Interfaces:       session.Interfaces,
		Capabilities:     session.Capabilities,
		Secrets:          session.Secrets,
		ManifestDigest:   session.ManifestDigest,
		GrantDigest:      session.GrantDigest,
	})
	if err != nil {
		return "", err
	}
	return digestBytes(encoded), nil
}
