package extensions_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gotthboard/gotth-extensions/pkg/extensions"
)

var (
	_ func([]byte) (extensions.Manifest, error)                                                       = extensions.ParseManifest
	_ func([]byte) (extensions.Grant, error)                                                          = extensions.ParseGrant
	_ func(extensions.Manifest) error                                                                 = extensions.ValidateManifest
	_ func(extensions.Grant) error                                                                    = extensions.ValidateGrant
	_ func(extensions.Manifest) ([]byte, error)                                                       = extensions.CanonicalManifest
	_ func(extensions.Manifest) (string, error)                                                       = extensions.ManifestDigest
	_ func(extensions.Grant) ([]byte, error)                                                          = extensions.CanonicalGrant
	_ func(extensions.Grant) (string, error)                                                          = extensions.GrantDigest
	_ func(extensions.Manifest, extensions.Grant, extensions.HostProfile) (extensions.Session, error) = extensions.Negotiate
	_ func(extensions.State, extensions.State) error                                                  = extensions.ValidateTransition
)

type fieldSpec struct {
	name, typ, jsonTag string
}

type sessionVector struct {
	ExtensionID      string                      `json:"extension_id"`
	ExtensionVersion string                      `json:"extension_version"`
	InstanceID       string                      `json:"instance_id"`
	Control          extensions.ProtocolVersion  `json:"control"`
	Interfaces       []extensions.InterfaceGrant `json:"interfaces"`
	Capabilities     []string                    `json:"capabilities"`
	Secrets          []string                    `json:"secrets"`
	ManifestDigest   string                      `json:"manifest_sha256"`
	GrantDigest      string                      `json:"grant_sha256"`
	Fingerprint      string                      `json:"fingerprint"`
}

type hostVector struct {
	Protocols  []extensions.VersionRange `json:"protocols"`
	Interfaces []extensions.VersionRange `json:"interfaces"`
}

func (v hostVector) profile() extensions.HostProfile {
	return extensions.HostProfile{Protocols: v.Protocols, Interfaces: v.Interfaces}
}

type compatibilityCase struct {
	Manifest          extensions.Manifest `json:"manifest"`
	CanonicalManifest string              `json:"canonical_manifest"`
	ManifestSHA256    string              `json:"manifest_sha256"`
	Grant             extensions.Grant    `json:"grant"`
	CanonicalGrant    string              `json:"canonical_grant"`
	GrantSHA256       string              `json:"grant_sha256"`
	Host              hostVector          `json:"host"`
	Session           sessionVector       `json:"session"`
	CanonicalSession  string              `json:"canonical_session"`
	SessionSHA256     string              `json:"session_sha256"`
}

type compatibilityVectors struct {
	Format                    string              `json:"format"`
	Manifest                  extensions.Manifest `json:"manifest"`
	CanonicalManifest         string              `json:"canonical_manifest"`
	ManifestSHA256            string              `json:"manifest_sha256"`
	Grant                     extensions.Grant    `json:"grant"`
	CanonicalGrant            string              `json:"canonical_grant"`
	GrantSHA256               string              `json:"grant_sha256"`
	Host                      hostVector          `json:"host"`
	Session                   sessionVector       `json:"session"`
	CanonicalSession          string              `json:"canonical_session"`
	SessionSHA256             string              `json:"session_sha256"`
	ControlOnly               compatibilityCase   `json:"control_only"`
	CanonicalDisplayNameCases []struct {
		Name              string `json:"name"`
		Document          string `json:"document"`
		ExpectedName      string `json:"expected_name"`
		CanonicalManifest string `json:"canonical_manifest"`
		ManifestSHA256    string `json:"manifest_sha256"`
	} `json:"canonical_display_name_cases"`
	RawParseCases []struct {
		Name          string `json:"name"`
		Target        string `json:"target"`
		Document      string `json:"document"`
		DocumentHex   string `json:"document_hex"`
		ExpectedError string `json:"expected_error"`
	} `json:"raw_parse_cases"`
	FixtureDecodeCases []struct {
		Name     string `json:"name"`
		Document string `json:"document"`
	} `json:"fixture_decode_cases"`
	NegotiationCases []struct {
		Name                  string `json:"name"`
		Mutation              string `json:"mutation"`
		ExpectedError         string `json:"expected_error"`
		ExpectedSessionSHA256 string `json:"expected_session_sha256"`
	} `json:"negotiation_cases"`
}

func loadV1Vectors(t *testing.T) compatibilityVectors {
	t.Helper()
	data, err := os.ReadFile("../../compat/v1/canonical-vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	vector, err := decodeV1Vectors(data)
	if err != nil {
		t.Fatal(err)
	}
	if vector.Format != "gotth.extensions.compatibility-vectors.v1" {
		t.Fatalf("unexpected vector format %q", vector.Format)
	}
	return vector
}

func decodeV1Vectors(data []byte) (compatibilityVectors, error) {
	var vector compatibilityVectors
	if err := rejectNonExactJSONFields(data, reflect.TypeOf(vector)); err != nil {
		return compatibilityVectors{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&vector); err != nil {
		return compatibilityVectors{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return compatibilityVectors{}, errors.New("compatibility vector has trailing JSON value")
	}
	return vector, nil
}

func rejectNonExactJSONFields(data []byte, target reflect.Type) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := scanExactJSONValue(decoder, target); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return errors.New("compatibility vector has trailing JSON value")
	}
	return nil
}

func scanExactJSONValue(decoder *json.Decoder, target reflect.Type) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	switch delimiter {
	case '{':
		fields, exact := exactJSONFields(target)
		seen := make(map[string]struct{})
		for decoder.More() {
			nameToken, err := decoder.Token()
			if err != nil {
				return err
			}
			name, ok := nameToken.(string)
			if !ok {
				return errors.New("compatibility vector has a non-string object name")
			}
			if _, duplicate := seen[name]; duplicate {
				return errors.New("compatibility vector has a duplicate object name")
			}
			seen[name] = struct{}{}
			child, exists := fields[name]
			if exact && !exists {
				return errors.New("compatibility vector has a non-exact object name")
			}
			if err := scanExactJSONValue(decoder, child); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim('}') {
			return errors.New("compatibility vector has an invalid object")
		}
	case '[':
		element := exactJSONElement(target)
		for decoder.More() {
			if err := scanExactJSONValue(decoder, element); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim(']') {
			return errors.New("compatibility vector has an invalid array")
		}
	default:
		return errors.New("compatibility vector has an invalid delimiter")
	}
	return nil
}

func exactJSONFields(target reflect.Type) (map[string]reflect.Type, bool) {
	for target != nil && target.Kind() == reflect.Pointer {
		target = target.Elem()
	}
	if target == nil || target.Kind() != reflect.Struct {
		return nil, false
	}
	fields := make(map[string]reflect.Type)
	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			name = strings.Split(tag, ",")[0]
		}
		if name != "-" {
			fields[name] = field.Type
		}
	}
	return fields, true
}

func exactJSONElement(target reflect.Type) reflect.Type {
	for target != nil && target.Kind() == reflect.Pointer {
		target = target.Elem()
	}
	if target != nil && (target.Kind() == reflect.Array || target.Kind() == reflect.Slice) {
		return target.Elem()
	}
	return nil
}

func TestV1VectorDecoderRejectsNonExactFields(t *testing.T) {
	tests := loadV1Vectors(t).FixtureDecodeCases
	wantNames := map[string]bool{
		"top-level-unknown": false, "nested-unknown": false,
		"top-level-case-alias": false, "top-level-alias-collision": false,
		"nested-case-alias": false, "nested-uppercase-alias": false,
		"nested-alias-collision": false, "deep-case-alias": false,
		"session-control-case-alias": false, "host-case-alias": false,
		"case-typo": false, "exact-duplicate": false,
	}
	if len(tests) != len(wantNames) {
		t.Fatalf("fixture decoder negative case count = %d, want %d", len(tests), len(wantNames))
	}
	for _, test := range tests {
		if _, ok := wantNames[test.Name]; !ok {
			t.Fatalf("unexpected fixture decoder case %q", test.Name)
		}
		wantNames[test.Name] = true
		t.Run(test.Name, func(t *testing.T) {
			if _, err := decodeV1Vectors([]byte(test.Document)); err == nil {
				t.Fatalf("non-exact fixture object name accepted in %s", test.Document)
			}
		})
	}
	for name, seen := range wantNames {
		if !seen {
			t.Fatalf("missing fixture decoder case %q", name)
		}
	}
}

func TestV1ExportedDataShapes(t *testing.T) {
	assertStruct(t, extensions.VersionRange{}, []fieldSpec{{"Name", "string", "name"}, {"Major", "uint32", "major"}, {"MinMinor", "uint32", "min_minor"}, {"MaxMinor", "uint32", "max_minor"}})
	assertStruct(t, extensions.InterfaceGrant{}, []fieldSpec{{"Name", "string", "name"}, {"Major", "uint32", "major"}, {"Minor", "uint32", "minor"}})
	assertStruct(t, extensions.SecretRequirement{}, []fieldSpec{{"ID", "string", "id"}, {"Required", "bool", "required"}})
	assertStruct(t, extensions.Manifest{}, []fieldSpec{{"Schema", "string", "schema"}, {"ID", "string", "id"}, {"Name", "string", "name"}, {"Version", "string", "version"}, {"Protocols", "[]extensions.VersionRange", "protocols"}, {"Interfaces", "[]extensions.VersionRange", "interfaces"}, {"Capabilities", "[]string", "capabilities"}, {"Secrets", "[]extensions.SecretRequirement", "secrets"}})
	assertStruct(t, extensions.Grant{}, []fieldSpec{{"Schema", "string", "schema"}, {"InstanceID", "string", "instance_id"}, {"ExtensionID", "string", "extension_id"}, {"ManifestDigest", "string", "manifest_sha256"}, {"Capabilities", "[]string", "capabilities"}, {"Interfaces", "[]extensions.InterfaceGrant", "interfaces"}, {"Secrets", "[]string", "secrets"}})
	assertStruct(t, extensions.HostProfile{}, []fieldSpec{{"Protocols", "[]extensions.VersionRange", ""}, {"Interfaces", "[]extensions.VersionRange", ""}})
	assertStruct(t, extensions.ProtocolVersion{}, []fieldSpec{{"Name", "string", "name"}, {"Major", "uint32", "major"}, {"Minor", "uint32", "minor"}})
	assertStruct(t, extensions.Session{}, []fieldSpec{{"ExtensionID", "string", ""}, {"ExtensionVersion", "string", ""}, {"InstanceID", "string", ""}, {"Control", "extensions.ProtocolVersion", ""}, {"Interfaces", "[]extensions.InterfaceGrant", ""}, {"Capabilities", "[]string", ""}, {"Secrets", "[]string", ""}, {"ManifestDigest", "string", ""}, {"GrantDigest", "string", ""}, {"Fingerprint", "string", ""}})
}

func assertStruct(t *testing.T, value any, want []fieldSpec) {
	t.Helper()
	typ := reflect.TypeOf(value)
	if typ.NumField() != len(want) {
		t.Fatalf("%s has %d fields, want %d", typ, typ.NumField(), len(want))
	}
	for i, expected := range want {
		field := typ.Field(i)
		if field.Name != expected.name || field.Type.String() != expected.typ || field.Tag.Get("json") != expected.jsonTag {
			t.Fatalf("%s field %d = (%s, %s, %q), want (%s, %s, %q)", typ, i, field.Name, field.Type, field.Tag.Get("json"), expected.name, expected.typ, expected.jsonTag)
		}
	}
}

func TestV1ConstantsAndSentinels(t *testing.T) {
	constants := map[string]any{
		"ManifestSchema": extensions.ManifestSchema, "GrantSchema": extensions.GrantSchema, "ControlName": extensions.ControlName,
		"MaxDocumentBytes": extensions.MaxDocumentBytes, "MaxIdentifier": extensions.MaxIdentifier, "MaxDisplayName": extensions.MaxDisplayName,
		"MaxVersionBytes": extensions.MaxVersionBytes, "MaxJSONDepth": extensions.MaxJSONDepth, "MaxProtocols": extensions.MaxProtocols,
		"MaxInterfaces": extensions.MaxInterfaces, "MaxCapabilities": extensions.MaxCapabilities, "MaxSecrets": extensions.MaxSecrets,
	}
	want := map[string]any{
		"ManifestSchema": "gotth.extensions.manifest.v1", "GrantSchema": "gotth.extensions.grant.v1", "ControlName": "gotth.extensions.control",
		"MaxDocumentBytes": 65536, "MaxIdentifier": 128, "MaxDisplayName": 128, "MaxVersionBytes": 64, "MaxJSONDepth": 32,
		"MaxProtocols": 16, "MaxInterfaces": 64, "MaxCapabilities": 128, "MaxSecrets": 64,
	}
	if !reflect.DeepEqual(constants, want) {
		t.Fatalf("V1 constants changed: %#v", constants)
	}
	for name, sentinel := range v1Sentinels() {
		messages := map[string]string{"ErrInvalidInput": "extensions: invalid input", "ErrUnsupportedSchema": "extensions: unsupported schema", "ErrIncompatibleVersion": "extensions: incompatible version", "ErrStaleBinding": "extensions: stale binding", "ErrUngrantedCapability": "extensions: ungranted capability", "ErrUngrantedInterface": "extensions: ungranted interface", "ErrUngrantedSecret": "extensions: ungranted secret", "ErrInvalidTransition": "extensions: invalid transition"}
		if sentinel.Error() != messages[name] || !errors.Is(sentinel, sentinel) {
			t.Fatalf("V1 sentinel changed: %s: %v", name, sentinel)
		}
	}
}

func v1Sentinels() map[string]error {
	return map[string]error{"ErrInvalidInput": extensions.ErrInvalidInput, "ErrUnsupportedSchema": extensions.ErrUnsupportedSchema, "ErrIncompatibleVersion": extensions.ErrIncompatibleVersion, "ErrStaleBinding": extensions.ErrStaleBinding, "ErrUngrantedCapability": extensions.ErrUngrantedCapability, "ErrUngrantedInterface": extensions.ErrUngrantedInterface, "ErrUngrantedSecret": extensions.ErrUngrantedSecret, "ErrInvalidTransition": extensions.ErrInvalidTransition}
}

func TestV1CanonicalAndSessionVectors(t *testing.T) {
	vector := loadV1Vectors(t)
	manifestBytes, err := extensions.CanonicalManifest(vector.Manifest)
	if err != nil || string(manifestBytes) != vector.CanonicalManifest {
		t.Fatalf("manifest vector drift: %s (%v)", manifestBytes, err)
	}
	manifestDigest, err := extensions.ManifestDigest(vector.Manifest)
	if err != nil || manifestDigest != vector.ManifestSHA256 {
		t.Fatalf("manifest digest drift: %s (%v)", manifestDigest, err)
	}
	grantBytes, err := extensions.CanonicalGrant(vector.Grant)
	if err != nil || string(grantBytes) != vector.CanonicalGrant {
		t.Fatalf("grant vector drift: %s (%v)", grantBytes, err)
	}
	grantDigest, err := extensions.GrantDigest(vector.Grant)
	if err != nil || grantDigest != vector.GrantSHA256 {
		t.Fatalf("grant digest drift: %s (%v)", grantDigest, err)
	}
	session, err := extensions.Negotiate(vector.Manifest, vector.Grant, vector.Host.profile())
	if err != nil {
		t.Fatal(err)
	}
	wantSession := extensions.Session{ExtensionID: vector.Session.ExtensionID, ExtensionVersion: vector.Session.ExtensionVersion, InstanceID: vector.Session.InstanceID, Control: vector.Session.Control, Interfaces: vector.Session.Interfaces, Capabilities: vector.Session.Capabilities, Secrets: vector.Session.Secrets, ManifestDigest: vector.Session.ManifestDigest, GrantDigest: vector.Session.GrantDigest, Fingerprint: vector.Session.Fingerprint}
	if !reflect.DeepEqual(session, wantSession) {
		t.Fatalf("complete session vector drift: got=%#v want=%#v", session, wantSession)
	}
	digest := sha256.Sum256([]byte(vector.CanonicalSession))
	if got := hex.EncodeToString(digest[:]); got != vector.SessionSHA256 || got != session.Fingerprint {
		t.Fatalf("session canonical bytes drift: %s", got)
	}
}

func TestV1CanonicalDisplayNameVectors(t *testing.T) {
	tests := loadV1Vectors(t).CanonicalDisplayNameCases
	if len(tests) != 2 {
		t.Fatalf("canonical display-name case count = %d, want 2", len(tests))
	}
	wantNames := map[string]bool{"line-and-paragraph-separators": false, "utf16-surrogate-pair-escape": false}
	for _, test := range tests {
		if _, ok := wantNames[test.Name]; !ok {
			t.Fatalf("unexpected canonical display-name case %q", test.Name)
		}
		wantNames[test.Name] = true
		t.Run(test.Name, func(t *testing.T) {
			manifest, err := extensions.ParseManifest([]byte(test.Document))
			if err != nil {
				t.Fatal(err)
			}
			if manifest.Name != test.ExpectedName {
				t.Fatalf("decoded display name = %q, want %q", manifest.Name, test.ExpectedName)
			}
			canonical, err := extensions.CanonicalManifest(manifest)
			if err != nil || string(canonical) != test.CanonicalManifest {
				t.Fatalf("canonical display-name vector drift: %s (%v)", canonical, err)
			}
			digest := sha256.Sum256(canonical)
			if got := hex.EncodeToString(digest[:]); got != test.ManifestSHA256 {
				t.Fatalf("display-name digest = %s, want %s", got, test.ManifestSHA256)
			}
		})
	}
	for name, seen := range wantNames {
		if !seen {
			t.Fatalf("missing canonical display-name case %q", name)
		}
	}
}

func TestV1RawParseVectors(t *testing.T) {
	vector := loadV1Vectors(t)
	exactKeyCases := map[string]bool{
		"manifest-schema-case-alias": false, "manifest-schema-uppercase-alias": false,
		"manifest-schema-alias-collision": false, "manifest-protocol-name-case-alias": false,
		"manifest-secret-required-case-typo": false, "grant-schema-case-alias": false,
		"grant-schema-alias-collision": false, "grant-interface-minor-case-alias": false,
	}
	for _, test := range vector.RawParseCases {
		_, exactKeyCase := exactKeyCases[test.Name]
		if exactKeyCase {
			exactKeyCases[test.Name] = true
		}
		t.Run(test.Name, func(t *testing.T) {
			if (test.Document == "") == (test.DocumentHex == "") {
				t.Fatal("raw vector must set exactly one of document or document_hex")
			}
			document := []byte(test.Document)
			if test.DocumentHex != "" {
				var decodeErr error
				document, decodeErr = hex.DecodeString(test.DocumentHex)
				if decodeErr != nil {
					t.Fatalf("invalid document_hex: %v", decodeErr)
				}
			}
			if exactKeyCase {
				switch test.Target {
				case "manifest":
					var permissive extensions.Manifest
					if decodeErr := json.Unmarshal(document, &permissive); decodeErr != nil || extensions.ValidateManifest(permissive) != nil {
						t.Fatalf("exact-key manifest vector has invalid permissive precondition: %v", decodeErr)
					}
				case "grant":
					var permissive extensions.Grant
					if decodeErr := json.Unmarshal(document, &permissive); decodeErr != nil || extensions.ValidateGrant(permissive) != nil {
						t.Fatalf("exact-key grant vector has invalid permissive precondition: %v", decodeErr)
					}
				}
			}
			var err error
			switch test.Target {
			case "manifest":
				_, err = extensions.ParseManifest(document)
			case "grant":
				_, err = extensions.ParseGrant(document)
			default:
				t.Fatalf("unknown target %q", test.Target)
			}
			if test.ExpectedError == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if !errors.Is(err, v1Sentinels()[test.ExpectedError]) {
				t.Fatalf("got %v, want %s", err, test.ExpectedError)
			}
		})
	}
	for name, seen := range exactKeyCases {
		if !seen {
			t.Fatalf("missing exact-key raw parser case %q", name)
		}
	}
}

func TestV1ControlOnlyVector(t *testing.T) {
	vector := loadV1Vectors(t).ControlOnly
	manifestBytes, err := extensions.CanonicalManifest(vector.Manifest)
	if err != nil || string(manifestBytes) != vector.CanonicalManifest {
		t.Fatalf("control-only manifest vector drift: %s (%v)", manifestBytes, err)
	}
	if digest := sha256.Sum256(manifestBytes); hex.EncodeToString(digest[:]) != vector.ManifestSHA256 {
		t.Fatal("control-only manifest digest drift")
	}
	grantBytes, err := extensions.CanonicalGrant(vector.Grant)
	if err != nil || string(grantBytes) != vector.CanonicalGrant {
		t.Fatalf("control-only grant vector drift: %s (%v)", grantBytes, err)
	}
	if digest := sha256.Sum256(grantBytes); hex.EncodeToString(digest[:]) != vector.GrantSHA256 {
		t.Fatal("control-only grant digest drift")
	}
	session, err := extensions.Negotiate(vector.Manifest, vector.Grant, vector.Host.profile())
	if err != nil {
		t.Fatal(err)
	}
	wantSession := extensions.Session{ExtensionID: vector.Session.ExtensionID, ExtensionVersion: vector.Session.ExtensionVersion, InstanceID: vector.Session.InstanceID, Control: vector.Session.Control, Interfaces: vector.Session.Interfaces, Capabilities: vector.Session.Capabilities, Secrets: vector.Session.Secrets, ManifestDigest: vector.Session.ManifestDigest, GrantDigest: vector.Session.GrantDigest, Fingerprint: vector.Session.Fingerprint}
	if !reflect.DeepEqual(session, wantSession) {
		t.Fatalf("control-only complete Session drift: got=%#v want=%#v", session, wantSession)
	}
	if session.Interfaces == nil || session.Capabilities == nil || session.Secrets == nil {
		t.Fatalf("empty Session sets must be non-nil []: %#v", session)
	}
	digest := sha256.Sum256([]byte(vector.CanonicalSession))
	if got := hex.EncodeToString(digest[:]); got != vector.SessionSHA256 || got != session.Fingerprint {
		t.Fatalf("control-only Session canonical bytes drift: %s", got)
	}
}

func TestV1NegotiationVectors(t *testing.T) {
	vector := loadV1Vectors(t)
	sawHostBoundary := false
	for _, test := range vector.NegotiationCases {
		t.Run(test.Name, func(t *testing.T) {
			manifest := cloneManifest(vector.Manifest)
			grant := cloneGrant(vector.Grant)
			host := cloneHost(vector.Host.profile())
			switch test.Mutation {
			case "none":
			case "grant.extension_id":
				grant.ExtensionID = "gotth.dns.other"
			case "grant.manifest_sha256":
				grant.ManifestDigest = strings.Repeat("a", 64)
			case "host.control.no_overlap":
				host.Protocols[0].MinMinor = 3
			case "grant.capability.unknown":
				grant.Capabilities[0] = "dns.records.delete"
			case "grant.interface.out_of_range":
				grant.Interfaces[0].Minor = 4
			case "grant.interface.outside_host_range":
				grant.Interfaces[0].Minor = 1
			case "grant.secret.unknown":
				grant.Secrets[0] = "dns.api.unknown"
			case "grant.secret.missing_required":
				grant.Secrets = []string{"dns.api.key"}
			default:
				t.Fatalf("unknown mutation %q", test.Mutation)
			}
			if test.Mutation == "grant.interface.outside_host_range" {
				sawHostBoundary = true
				minor := grant.Interfaces[0].Minor
				manifestRange := manifest.Interfaces[0]
				hostRange := host.Interfaces[0]
				if minor < manifestRange.MinMinor || minor > manifestRange.MaxMinor {
					t.Fatalf("vector precondition failed: grant minor %d is outside manifest range", minor)
				}
				if minor >= hostRange.MinMinor && minor <= hostRange.MaxMinor {
					t.Fatalf("vector precondition failed: grant minor %d is inside host range", minor)
				}
			}
			session, err := extensions.Negotiate(manifest, grant, host)
			if test.ExpectedError == "" {
				if err != nil || session.Fingerprint != test.ExpectedSessionSHA256 {
					t.Fatalf("session=%#v err=%v", session, err)
				}
				return
			}
			if !errors.Is(err, v1Sentinels()[test.ExpectedError]) {
				t.Fatalf("got %v, want %s", err, test.ExpectedError)
			}
		})
	}
	if !sawHostBoundary {
		t.Fatal("missing manifest-valid but host-range-invalid negotiation vector")
	}
}

func TestV1GrammarAndBounds(t *testing.T) {
	vector := loadV1Vectors(t)
	manifest := cloneManifest(vector.Manifest)
	manifest.ID = strings.Repeat("a", 63) + "." + strings.Repeat("b", 62) + ".c"
	manifest.Name = strings.Repeat("n", extensions.MaxDisplayName)
	manifest.Version = "1.0.0-" + strings.Repeat("a", extensions.MaxVersionBytes-len("1.0.0-"))
	manifest.Protocols = []extensions.VersionRange{{Name: extensions.ControlName, Major: 1, MinMinor: 0, MaxMinor: 1<<31 - 1}}
	for i := 1; i < extensions.MaxProtocols; i++ {
		manifest.Protocols = append(manifest.Protocols, extensions.VersionRange{Name: "gotth.protocol.p" + string(rune('a'+i)), Major: 1})
	}
	manifest.Interfaces = nil
	for i := 0; i < extensions.MaxInterfaces; i++ {
		manifest.Interfaces = append(manifest.Interfaces, extensions.VersionRange{Name: "gotth.interface.i" + decimal(i), Major: 1})
	}
	manifest.Capabilities = nil
	for i := 0; i < extensions.MaxCapabilities; i++ {
		manifest.Capabilities = append(manifest.Capabilities, "capability.item."+decimal(i))
	}
	manifest.Secrets = nil
	for i := 0; i < extensions.MaxSecrets; i++ {
		manifest.Secrets = append(manifest.Secrets, extensions.SecretRequirement{ID: "secret.slot." + decimal(i)})
	}
	if err := extensions.ValidateManifest(manifest); err != nil {
		t.Fatalf("maximum valid manifest rejected: %v", err)
	}
	manifest.ID += "x"
	if !errors.Is(extensions.ValidateManifest(manifest), extensions.ErrInvalidInput) {
		t.Fatal("overlong identifier accepted")
	}

	document := []byte(vector.CanonicalManifest)
	document = append(document, []byte(strings.Repeat(" ", extensions.MaxDocumentBytes-len(document)))...)
	if _, err := extensions.ParseManifest(document); err != nil {
		t.Fatalf("maximum-size document rejected: %v", err)
	}
	document = append(document, ' ')
	if _, err := extensions.ParseManifest(document); !errors.Is(err, extensions.ErrInvalidInput) {
		t.Fatal("oversized document accepted")
	}
	tooDeep := strings.Repeat("[", extensions.MaxJSONDepth+1) + "0" + strings.Repeat("]", extensions.MaxJSONDepth+1)
	if _, err := extensions.ParseManifest([]byte(tooDeep)); !errors.Is(err, extensions.ErrInvalidInput) {
		t.Fatal("over-depth document accepted")
	}
}

func decimal(value int) string {
	if value == 0 {
		return "zero"
	}
	const digits = "abcdefghijklmnopqrstuvwxyz"
	result := ""
	for value > 0 {
		result = string(digits[value%len(digits)]) + result
		value /= len(digits)
	}
	return result
}

func TestV1SessionCopyAndRevalidationRule(t *testing.T) {
	vector := loadV1Vectors(t)
	authoritative, err := extensions.Negotiate(vector.Manifest, vector.Grant, vector.Host.profile())
	if err != nil {
		t.Fatal(err)
	}
	retained := cloneSession(authoritative)
	retained.Capabilities[0] = "dns.records.tampered"
	if reflect.DeepEqual(retained, authoritative) || retained.Fingerprint != authoritative.Fingerprint {
		t.Fatal("mutation test did not preserve the stale cached fingerprint")
	}
	revalidated, err := extensions.Negotiate(vector.Manifest, vector.Grant, vector.Host.profile())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(revalidated, authoritative) || reflect.DeepEqual(revalidated, retained) {
		t.Fatal("complete-session revalidation did not detect mutation")
	}
	copied := cloneSession(authoritative)
	copied.Secrets[0] = "secret.slot.changed"
	if authoritative.Secrets[0] == copied.Secrets[0] {
		t.Fatal("session deep copy aliases slices")
	}
}

func cloneManifest(value extensions.Manifest) extensions.Manifest {
	value.Protocols = append([]extensions.VersionRange(nil), value.Protocols...)
	value.Interfaces = append([]extensions.VersionRange(nil), value.Interfaces...)
	value.Capabilities = append([]string(nil), value.Capabilities...)
	value.Secrets = append([]extensions.SecretRequirement(nil), value.Secrets...)
	return value
}
func cloneGrant(value extensions.Grant) extensions.Grant {
	value.Capabilities = append([]string(nil), value.Capabilities...)
	value.Interfaces = append([]extensions.InterfaceGrant(nil), value.Interfaces...)
	value.Secrets = append([]string(nil), value.Secrets...)
	return value
}
func cloneHost(value extensions.HostProfile) extensions.HostProfile {
	value.Protocols = append([]extensions.VersionRange(nil), value.Protocols...)
	value.Interfaces = append([]extensions.VersionRange(nil), value.Interfaces...)
	return value
}
func cloneSession(value extensions.Session) extensions.Session {
	value.Interfaces = append([]extensions.InterfaceGrant(nil), value.Interfaces...)
	value.Capabilities = append([]string(nil), value.Capabilities...)
	value.Secrets = append([]string(nil), value.Secrets...)
	return value
}

func TestV1LifecycleEdges(t *testing.T) {
	states := []extensions.State{extensions.StateDiscovered, extensions.StateStarting, extensions.StateReady, extensions.StateDegraded, extensions.StateStopping, extensions.StateStopped, extensions.StateFailed}
	wantValues := []extensions.State{"discovered", "starting", "ready", "degraded", "stopping", "stopped", "failed"}
	if !reflect.DeepEqual(states, wantValues) {
		t.Fatalf("V1 state strings changed: %v", states)
	}
	allowed := map[[2]extensions.State]bool{{extensions.StateDiscovered, extensions.StateStarting}: true, {extensions.StateStarting, extensions.StateReady}: true, {extensions.StateStarting, extensions.StateFailed}: true, {extensions.StateStarting, extensions.StateStopping}: true, {extensions.StateReady, extensions.StateDegraded}: true, {extensions.StateReady, extensions.StateStopping}: true, {extensions.StateReady, extensions.StateFailed}: true, {extensions.StateDegraded, extensions.StateReady}: true, {extensions.StateDegraded, extensions.StateStopping}: true, {extensions.StateDegraded, extensions.StateFailed}: true, {extensions.StateStopping, extensions.StateStopped}: true, {extensions.StateStopping, extensions.StateFailed}: true, {extensions.StateStopped, extensions.StateStarting}: true, {extensions.StateFailed, extensions.StateStarting}: true, {extensions.StateFailed, extensions.StateStopped}: true}
	for _, from := range states {
		for _, to := range states {
			err := extensions.ValidateTransition(from, to)
			if allowed[[2]extensions.State{from, to}] != (err == nil) {
				t.Fatalf("V1 lifecycle edge changed: %s -> %s: %v", from, to, err)
			}
		}
	}
	for _, edge := range [][2]extensions.State{{"unknown", extensions.StateReady}, {extensions.StateReady, "unknown"}} {
		if !errors.Is(extensions.ValidateTransition(edge[0], edge[1]), extensions.ErrInvalidTransition) {
			t.Fatalf("V1 unknown-state behavior changed: %s -> %s", edge[0], edge[1])
		}
	}
}
