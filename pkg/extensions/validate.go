package extensions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ParseManifest(data []byte) (Manifest, error) {
	var value Manifest
	if err := decodeStrict(data, &value); err != nil {
		return Manifest{}, err
	}
	if err := ValidateManifest(value); err != nil {
		return Manifest{}, err
	}
	return value, nil
}

func ParseGrant(data []byte) (Grant, error) {
	var value Grant
	if err := decodeStrict(data, &value); err != nil {
		return Grant{}, err
	}
	if err := ValidateGrant(value); err != nil {
		return Grant{}, err
	}
	return value, nil
}

func decodeStrict(data []byte, target any) error {
	if len(data) == 0 || len(data) > MaxDocumentBytes {
		return fieldError("document")
	}
	if err := rejectDuplicateObjectNames(data); err != nil {
		return fieldError("document")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fieldError("document")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fieldError("document")
	}
	return nil
}

func rejectDuplicateObjectNames(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := scanJSONValue(decoder, 0); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return ErrInvalidInput
	}
	return nil
}

func scanJSONValue(decoder *json.Decoder, depth int) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	if depth >= MaxJSONDepth {
		return ErrInvalidInput
	}
	switch delimiter {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			nameToken, err := decoder.Token()
			if err != nil {
				return err
			}
			name, ok := nameToken.(string)
			if !ok {
				return ErrInvalidInput
			}
			if _, exists := seen[name]; exists {
				return ErrInvalidInput
			}
			seen[name] = struct{}{}
			if err := scanJSONValue(decoder, depth+1); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim('}') {
			return ErrInvalidInput
		}
	case '[':
		for decoder.More() {
			if err := scanJSONValue(decoder, depth+1); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim(']') {
			return ErrInvalidInput
		}
	default:
		return ErrInvalidInput
	}
	return nil
}

func ValidateManifest(m Manifest) error {
	if m.Schema != ManifestSchema {
		return ErrUnsupportedSchema
	}
	if !validDottedID(m.ID, 3) {
		return fieldError("id")
	}
	if !validDisplayName(m.Name) {
		return fieldError("name")
	}
	if !validSemver(m.Version) {
		return fieldError("version")
	}
	if err := validateRanges(m.Protocols, MaxProtocols, true); err != nil {
		return err
	}
	if err := validateRanges(m.Interfaces, MaxInterfaces, false); err != nil {
		return err
	}
	if err := validateTokens(m.Capabilities, MaxCapabilities, "capabilities"); err != nil {
		return err
	}
	if len(m.Secrets) > MaxSecrets {
		return fieldError("secrets")
	}
	seenSecrets := make(map[string]struct{}, len(m.Secrets))
	for _, secret := range m.Secrets {
		if !validDottedID(secret.ID, 2) {
			return fieldError("secrets")
		}
		if _, exists := seenSecrets[secret.ID]; exists {
			return fieldError("secrets")
		}
		seenSecrets[secret.ID] = struct{}{}
	}
	return nil
}

func ValidateGrant(g Grant) error {
	if g.Schema != GrantSchema {
		return ErrUnsupportedSchema
	}
	if !validUUID(g.InstanceID) {
		return fieldError("instance_id")
	}
	if !validDottedID(g.ExtensionID, 3) {
		return fieldError("extension_id")
	}
	if !validDigest(g.ManifestDigest) {
		return fieldError("manifest_sha256")
	}
	if err := validateTokens(g.Capabilities, MaxCapabilities, "capabilities"); err != nil {
		return err
	}
	if len(g.Interfaces) > MaxInterfaces {
		return fieldError("interfaces")
	}
	seenInterfaces := make(map[string]struct{}, len(g.Interfaces))
	for _, item := range g.Interfaces {
		if !validDottedID(item.Name, 2) || item.Major == 0 || item.Major > maxVersionNumber || item.Minor > maxVersionNumber {
			return fieldError("interfaces")
		}
		if _, exists := seenInterfaces[item.Name]; exists {
			return fieldError("interfaces")
		}
		seenInterfaces[item.Name] = struct{}{}
	}
	if err := validateTokens(g.Secrets, MaxSecrets, "secrets"); err != nil {
		return err
	}
	return nil
}

func validateHostProfile(h HostProfile) error {
	if err := validateRanges(h.Protocols, MaxProtocols, true); err != nil {
		return err
	}
	return validateRanges(h.Interfaces, MaxInterfaces, false)
}

func validateRanges(ranges []VersionRange, maximum int, requireControl bool) error {
	if len(ranges) > maximum || requireControl && len(ranges) == 0 {
		return fieldError("version_ranges")
	}
	seen := make(map[string]struct{}, len(ranges))
	control := 0
	for _, item := range ranges {
		if !validDottedID(item.Name, 2) || item.Major == 0 || item.Major > maxVersionNumber || item.MinMinor > item.MaxMinor || item.MaxMinor > maxVersionNumber {
			return fieldError("version_ranges")
		}
		key := fmt.Sprintf("%s\x00%d", item.Name, item.Major)
		if _, exists := seen[key]; exists {
			return fieldError("version_ranges")
		}
		seen[key] = struct{}{}
		if item.Name == ControlName && item.Major == 1 {
			control++
		}
	}
	if requireControl && control != 1 {
		return ErrIncompatibleVersion
	}
	return nil
}

func validateTokens(values []string, maximum int, field string) error {
	if len(values) > maximum {
		return fieldError(field)
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !validDottedID(value, 2) {
			return fieldError(field)
		}
		if _, exists := seen[value]; exists {
			return fieldError(field)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validDottedID(value string, minimumSegments int) bool {
	if len(value) == 0 || len(value) > MaxIdentifier {
		return false
	}
	parts := strings.Split(value, ".")
	if len(parts) < minimumSegments {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 || part[0] < 'a' || part[0] > 'z' || part[len(part)-1] == '-' {
			return false
		}
		for i := range len(part) {
			c := part[i]
			if c < 'a' || c > 'z' {
				if c < '0' || c > '9' {
					if c != '-' {
						return false
					}
				}
			}
		}
	}
	return true
}

func validDisplayName(value string) bool {
	if len(value) == 0 || len(value) > MaxDisplayName || !utf8.ValidString(value) || strings.TrimSpace(value) != value {
		return false
	}
	for _, r := range value {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func validUUID(value string) bool {
	if len(value) != 36 {
		return false
	}
	for i := range len(value) {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if value[i] != '-' {
				return false
			}
			continue
		}
		if !isLowerHex(value[i]) {
			return false
		}
	}
	return true
}

func validDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for i := range len(value) {
		if !isLowerHex(value[i]) {
			return false
		}
	}
	return true
}

func isLowerHex(value byte) bool {
	return value >= '0' && value <= '9' || value >= 'a' && value <= 'f'
}

func fieldError(field string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, field)
}
