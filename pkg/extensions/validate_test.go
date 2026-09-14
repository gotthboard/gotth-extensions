package extensions

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateManifest(t *testing.T) {
	base := testManifest()
	if err := ValidateManifest(base); err != nil {
		t.Fatalf("valid manifest rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*Manifest)
		target error
	}{
		{"schema", func(v *Manifest) { v.Schema = "v2" }, ErrUnsupportedSchema},
		{"id", func(v *Manifest) { v.ID = "Bad.ID" }, ErrInvalidInput},
		{"name", func(v *Manifest) { v.Name = " bad" }, ErrInvalidInput},
		{"name control", func(v *Manifest) { v.Name = "bad\nname" }, ErrInvalidInput},
		{"name utf8", func(v *Manifest) { v.Name = string([]byte{0xff}) }, ErrInvalidInput},
		{"version", func(v *Manifest) { v.Version = "01.0.0" }, ErrInvalidInput},
		{"protocol empty", func(v *Manifest) { v.Protocols = nil }, ErrInvalidInput},
		{"protocol missing control", func(v *Manifest) { v.Protocols[0].Name = "gotth.extensions.other" }, ErrIncompatibleVersion},
		{"protocol duplicate", func(v *Manifest) { v.Protocols = append(v.Protocols, v.Protocols[0]) }, ErrInvalidInput},
		{"protocol inverted", func(v *Manifest) { v.Protocols[0].MinMinor = 3 }, ErrInvalidInput},
		{"capability duplicate", func(v *Manifest) { v.Capabilities = append(v.Capabilities, v.Capabilities[0]) }, ErrInvalidInput},
		{"capability invalid", func(v *Manifest) { v.Capabilities[0] = "DNS.READ" }, ErrInvalidInput},
		{"secret duplicate", func(v *Manifest) { v.Secrets = append(v.Secrets, v.Secrets[0]) }, ErrInvalidInput},
		{"secret invalid", func(v *Manifest) { v.Secrets[0].ID = "secret" }, ErrInvalidInput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := base
			value.Protocols = cloneRanges(base.Protocols)
			value.Interfaces = cloneRanges(base.Interfaces)
			value.Capabilities = cloneStrings(base.Capabilities)
			value.Secrets = append([]SecretRequirement(nil), base.Secrets...)
			test.mutate(&value)
			if err := ValidateManifest(value); !errors.Is(err, test.target) {
				t.Fatalf("got %v, want %v", err, test.target)
			}
		})
	}
}

func TestEmptyInterfaceSetsAreValid(t *testing.T) {
	manifest := testManifest()
	manifest.Interfaces = nil
	if err := ValidateManifest(manifest); err != nil {
		t.Fatalf("control-only manifest rejected: %v", err)
	}
	host := testHost()
	host.Interfaces = nil
	if err := validateHostProfile(host); err != nil {
		t.Fatalf("control-only host rejected: %v", err)
	}
	grant := testGrant(manifest)
	grant.Interfaces = nil
	if _, err := Negotiate(manifest, grant, host); err != nil {
		t.Fatalf("control-only session rejected: %v", err)
	}
}

func TestValidateManifestBounds(t *testing.T) {
	checks := []func(*Manifest){
		func(v *Manifest) { v.ID = strings.Repeat("a", MaxIdentifier+1) },
		func(v *Manifest) { v.Name = strings.Repeat("a", MaxDisplayName+1) },
		func(v *Manifest) { v.Protocols = make([]VersionRange, MaxProtocols+1) },
		func(v *Manifest) { v.Interfaces = make([]VersionRange, MaxInterfaces+1) },
		func(v *Manifest) { v.Capabilities = make([]string, MaxCapabilities+1) },
		func(v *Manifest) { v.Secrets = make([]SecretRequirement, MaxSecrets+1) },
	}
	for i, mutate := range checks {
		value := testManifest()
		mutate(&value)
		if !errors.Is(ValidateManifest(value), ErrInvalidInput) {
			t.Fatalf("bound case %d accepted", i)
		}
	}
}

func TestValidateGrant(t *testing.T) {
	manifest := testManifest()
	base := testGrant(manifest)
	if err := ValidateGrant(base); err != nil {
		t.Fatalf("valid grant rejected: %v", err)
	}
	tests := []struct {
		name   string
		mutate func(*Grant)
		target error
	}{
		{"schema", func(v *Grant) { v.Schema = "v2" }, ErrUnsupportedSchema},
		{"instance", func(v *Grant) { v.InstanceID = "ABC" }, ErrInvalidInput},
		{"extension", func(v *Grant) { v.ExtensionID = "bad" }, ErrInvalidInput},
		{"digest", func(v *Grant) { v.ManifestDigest = strings.Repeat("A", 64) }, ErrInvalidInput},
		{"cap duplicate", func(v *Grant) { v.Capabilities = append(v.Capabilities, v.Capabilities[0]) }, ErrInvalidInput},
		{"interface invalid", func(v *Grant) { v.Interfaces[0].Major = 0 }, ErrInvalidInput},
		{"interface duplicate", func(v *Grant) { v.Interfaces = append(v.Interfaces, v.Interfaces[0]) }, ErrInvalidInput},
		{"secret duplicate", func(v *Grant) { v.Secrets = append(v.Secrets, v.Secrets[0]) }, ErrInvalidInput},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := base
			value.Capabilities = cloneStrings(base.Capabilities)
			value.Interfaces = append([]InterfaceGrant(nil), base.Interfaces...)
			value.Secrets = cloneStrings(base.Secrets)
			test.mutate(&value)
			if err := ValidateGrant(value); !errors.Is(err, test.target) {
				t.Fatalf("got %v, want %v", err, test.target)
			}
		})
	}
}

func TestStrictDocumentParsing(t *testing.T) {
	manifest := testManifest()
	encoded, err := CanonicalManifest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseManifest(encoded)
	if err != nil || parsed.ID != manifest.ID {
		t.Fatalf("parse manifest: %#v %v", parsed, err)
	}
	grant := testGrant(manifest)
	grantJSON, err := CanonicalGrant(grant)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseGrant(grantJSON); err != nil {
		t.Fatalf("parse grant: %v", err)
	}

	bad := [][]byte{
		nil,
		[]byte(`{"schema":"gotth.extensions.manifest.v1","unknown":true}`),
		[]byte(`{"schema":"gotth.extensions.manifest.v1","schema":"gotth.extensions.manifest.v1"}`),
		[]byte(`{"schema":"gotth.extensions.manifest.v1","id":"gotth.dns.test","name":"Test","version":"1.0.0","protocols":[{"name":"gotth.extensions.control","name":"gotth.extensions.control","major":1,"min_minor":0,"max_minor":0}],"interfaces":[],"capabilities":[],"secrets":[]}`),
		append(append([]byte(nil), encoded...), []byte(` {}`)...),
		[]byte(`{`),
		make([]byte, MaxDocumentBytes+1),
	}
	for i, input := range bad {
		if _, err := ParseManifest(input); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("bad document %d returned %v", i, err)
		}
	}
	if _, err := ParseGrant(encoded); err == nil {
		t.Fatal("manifest decoded as grant")
	}
	if _, err := ParseGrant([]byte(`null`)); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("null grant returned %v", err)
	}
}

func TestJSONDuplicateScanner(t *testing.T) {
	valid := [][]byte{
		[]byte(`null`),
		[]byte(`true`),
		[]byte(`123`),
		[]byte(`"text"`),
		[]byte(`[1,{"a":2}]`),
		[]byte(`{"a":{"b":1},"c":[2,3]}`),
		[]byte(strings.Repeat("[", MaxJSONDepth) + "0" + strings.Repeat("]", MaxJSONDepth)),
	}
	for _, input := range valid {
		if err := rejectDuplicateObjectNames(input); err != nil {
			t.Fatalf("valid JSON rejected: %s: %v", input, err)
		}
	}
	invalid := [][]byte{
		[]byte(``),
		[]byte(`{"a":1,"a":2}`),
		[]byte(`{"a":{"b":1,"b":2}}`),
		[]byte(`[1,2`),
		[]byte(`{"a":1} {"b":2}`),
		[]byte(strings.Repeat("[", MaxJSONDepth+1) + "0" + strings.Repeat("]", MaxJSONDepth+1)),
	}
	for _, input := range invalid {
		if err := rejectDuplicateObjectNames(input); err == nil {
			t.Fatalf("invalid JSON accepted: %s", input)
		}
	}
}

func TestIdentifierUUIDDigestAndDisplayHelpers(t *testing.T) {
	validIDs := []string{"gotth.dns.godaddy", "dns.records.read", "a-b.c2.d3"}
	for _, value := range validIDs {
		if !validDottedID(value, 2) {
			t.Fatalf("valid ID rejected: %q", value)
		}
	}
	invalidIDs := []string{"", "one", ".a.b", "a..b", "a-.b", "A.b", "1a.b", "a._b", strings.Repeat("a", 64) + ".b"}
	for _, value := range invalidIDs {
		if validDottedID(value, 2) {
			t.Fatalf("invalid ID accepted: %q", value)
		}
	}
	if !validUUID("11111111-1111-4111-8111-111111111111") || validUUID("11111111-1111-4111-8111-11111111111A") || validUUID("111111111111-4111-8111-111111111111") {
		t.Fatal("UUID policy wrong")
	}
	if !validDigest(strings.Repeat("a", 64)) || validDigest(strings.Repeat("A", 64)) || validDigest("abc") {
		t.Fatal("digest policy wrong")
	}
	if !validDisplayName("DNS Provider") || validDisplayName("") || validDisplayName("name ") {
		t.Fatal("display-name policy wrong")
	}
}

func TestRangeValidationEdges(t *testing.T) {
	valid := []VersionRange{{Name: ControlName, Major: 1, MinMinor: 0, MaxMinor: 0}}
	if err := validateRanges(valid, 1, true); err != nil {
		t.Fatal(err)
	}
	tests := [][]VersionRange{
		{{Name: "bad", Major: 1}},
		{{Name: ControlName, Major: 0}},
		{{Name: ControlName, Major: maxVersionNumber + 1}},
		{{Name: ControlName, Major: 1, MaxMinor: maxVersionNumber + 1}},
	}
	for i, input := range tests {
		if !errors.Is(validateRanges(input, 2, true), ErrInvalidInput) {
			t.Fatalf("invalid range case %d accepted", i)
		}
	}
}

func TestValidationErrorsDoNotEchoInput(t *testing.T) {
	const marker = "do-not-echo-this-value"
	manifest := testManifest()
	manifest.Name = marker + "\n"
	if err := ValidateManifest(manifest); err == nil || strings.Contains(err.Error(), marker) {
		t.Fatalf("manifest error leaked input: %v", err)
	}
	input := []byte(`{"schema":"gotth.extensions.manifest.v1","` + marker + `":true}`)
	if _, err := ParseManifest(input); err == nil || strings.Contains(err.Error(), marker) {
		t.Fatalf("parse error leaked input: %v", err)
	}
	grant := testGrant(testManifest())
	grant.ExtensionID = marker
	if err := ValidateGrant(grant); err == nil || strings.Contains(err.Error(), marker) {
		t.Fatalf("grant error leaked input: %v", err)
	}
}
