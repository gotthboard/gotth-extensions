package extensions

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestCanonicalManifestDeterministicAndNonMutating(t *testing.T) {
	first := testManifest()
	second := testManifest()
	second.Protocols = reverseRanges(second.Protocols)
	second.Interfaces = reverseRanges(second.Interfaces)
	second.Capabilities = reverseStrings(second.Capabilities)
	second.Secrets = reverseSecrets(second.Secrets)
	original := second
	original.Protocols = cloneRanges(second.Protocols)
	original.Interfaces = cloneRanges(second.Interfaces)
	original.Capabilities = cloneStrings(second.Capabilities)
	original.Secrets = append([]SecretRequirement(nil), second.Secrets...)

	left, err := CanonicalManifest(first)
	if err != nil {
		t.Fatal(err)
	}
	right, err := CanonicalManifest(second)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(left, right) {
		t.Fatalf("canonical manifests differ:\n%s\n%s", left, right)
	}
	if !reflect.DeepEqual(second, original) {
		t.Fatal("canonicalization mutated caller manifest")
	}
	leftDigest, err := ManifestDigest(first)
	if err != nil {
		t.Fatal(err)
	}
	rightDigest, err := ManifestDigest(second)
	if err != nil || leftDigest != rightDigest || len(leftDigest) != 64 {
		t.Fatalf("digest mismatch: %q %q %v", leftDigest, rightDigest, err)
	}
	if leftDigest != "53bafab4e5fcd2cf5ff97ac88a46fb5de9f2e1866c816122760fd6f13c8c6308" {
		t.Fatalf("manifest vector changed: %s", leftDigest)
	}
}

func TestCanonicalGrantDeterministicAndNonMutating(t *testing.T) {
	manifest := testManifest()
	first := testGrant(manifest)
	second := testGrant(manifest)
	second.Capabilities = reverseStrings(second.Capabilities)
	second.Interfaces = reverseInterfaceGrants(second.Interfaces)
	second.Secrets = reverseStrings(second.Secrets)
	original := second
	original.Capabilities = cloneStrings(second.Capabilities)
	original.Interfaces = append([]InterfaceGrant(nil), second.Interfaces...)
	original.Secrets = cloneStrings(second.Secrets)

	left, err := CanonicalGrant(first)
	if err != nil {
		t.Fatal(err)
	}
	right, err := CanonicalGrant(second)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(left, right) {
		t.Fatalf("canonical grants differ:\n%s\n%s", left, right)
	}
	if !reflect.DeepEqual(second, original) {
		t.Fatal("canonicalization mutated caller grant")
	}
	digest, err := GrantDigest(first)
	if err != nil {
		t.Fatal(err)
	}
	if digest != "eb46f052bc968314394806567ac345af2016eeb86e4c56f131b63a7613cf0537" {
		t.Fatalf("grant vector changed: %s", digest)
	}
}

func TestCanonicalEmptySetsDoNotDependOnNilRepresentation(t *testing.T) {
	manifestNil := testManifest()
	manifestNil.Capabilities = nil
	manifestNil.Secrets = nil
	manifestEmpty := manifestNil
	manifestEmpty.Capabilities = []string{}
	manifestEmpty.Secrets = []SecretRequirement{}

	left, err := CanonicalManifest(manifestNil)
	if err != nil {
		t.Fatal(err)
	}
	right, err := CanonicalManifest(manifestEmpty)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(left, right) {
		t.Fatalf("nil and empty manifest sets differ:\n%s\n%s", left, right)
	}

	grantNil := testGrant(testManifest())
	grantNil.Capabilities = nil
	grantNil.Interfaces = nil
	grantNil.Secrets = nil
	grantEmpty := grantNil
	grantEmpty.Capabilities = []string{}
	grantEmpty.Interfaces = []InterfaceGrant{}
	grantEmpty.Secrets = []string{}
	left, err = CanonicalGrant(grantNil)
	if err != nil {
		t.Fatal(err)
	}
	right, err = CanonicalGrant(grantEmpty)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(left, right) {
		t.Fatalf("nil and empty grant sets differ:\n%s\n%s", left, right)
	}
}

func TestCanonicalOrderingAcrossMultipleEntries(t *testing.T) {
	manifest := testManifest()
	manifest.Protocols = append(manifest.Protocols,
		VersionRange{Name: "gotth.extensions.aux", Major: 2, MinMinor: 0, MaxMinor: 1},
		VersionRange{Name: "gotth.extensions.aux", Major: 1, MinMinor: 0, MaxMinor: 1},
	)
	manifest.Interfaces = append(manifest.Interfaces,
		VersionRange{Name: "gotth.extensions.backup", Major: 2, MinMinor: 0, MaxMinor: 1},
		VersionRange{Name: "gotth.extensions.backup", Major: 1, MinMinor: 0, MaxMinor: 1},
	)
	encoded, err := CanonicalManifest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Index(encoded, []byte(`"name":"gotth.extensions.aux","major":1`)) > bytes.Index(encoded, []byte(`"name":"gotth.extensions.aux","major":2`)) {
		t.Fatalf("manifest ranges not ordered by major: %s", encoded)
	}

	grant := testGrant(manifest)
	grant.Interfaces = append(grant.Interfaces,
		InterfaceGrant{Name: "gotth.extensions.backup", Major: 1, Minor: 1},
	)
	encoded, err = CanonicalGrant(grant)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Index(encoded, []byte(`"name":"gotth.extensions.backup"`)) > bytes.Index(encoded, []byte(`"name":"gotth.extensions.dns"`)) {
		t.Fatalf("grant interfaces not ordered by name: %s", encoded)
	}
}

func TestCanonicalErrorsAndBound(t *testing.T) {
	manifest := testManifest()
	manifest.Schema = "bad"
	if _, err := CanonicalManifest(manifest); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("canonical manifest error: %v", err)
	}
	if _, err := ManifestDigest(manifest); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("manifest digest error: %v", err)
	}
	grant := testGrant(testManifest())
	grant.Schema = "bad"
	if _, err := CanonicalGrant(grant); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("canonical grant error: %v", err)
	}
	if _, err := GrantDigest(grant); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("grant digest error: %v", err)
	}
	if _, err := encodeBounded(strings.Repeat("x", MaxDocumentBytes)); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("oversized encoding error: %v", err)
	}
	if _, err := encodeBounded(make(chan int)); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("marshal failure error: %v", err)
	}
	if digestBytes([]byte("x")) != "2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881" {
		t.Fatal("SHA-256 helper changed")
	}
}

func reverseStrings(values []string) []string {
	result := cloneStrings(values)
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

func reverseRanges(values []VersionRange) []VersionRange {
	result := cloneRanges(values)
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

func reverseSecrets(values []SecretRequirement) []SecretRequirement {
	result := append([]SecretRequirement(nil), values...)
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

func reverseInterfaceGrants(values []InterfaceGrant) []InterfaceGrant {
	result := append([]InterfaceGrant(nil), values...)
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}
