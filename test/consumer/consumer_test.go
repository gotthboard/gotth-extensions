package consumer

import (
	"errors"
	"testing"

	"github.com/gotthboard/gotth-extensions/pkg/extensions"
)

func TestPublicPackageCompiles(t *testing.T) {
	if extensions.ManifestSchema == "" || extensions.ControlName == "" {
		t.Fatal("public constants missing")
	}
}

func TestUntrustedDocumentsUseStrictParsers(t *testing.T) {
	manifest := []byte(`{"schema":"gotth.extensions.manifest.v1","id":"gotth.dns.example","name":"Example","version":"1.0.0","protocols":[{"name":"gotth.extensions.control","major":1,"min_minor":0,"max_minor":0}],"interfaces":[],"capabilities":[],"secrets":[],"unknown":true}`)
	if _, err := extensions.ParseManifest(manifest); !errors.Is(err, extensions.ErrInvalidInput) {
		t.Fatalf("unknown manifest field was not rejected: %v", err)
	}
	grant := []byte(`{"schema":"gotth.extensions.grant.v1","instance_id":"11111111-1111-4111-8111-111111111111","extension_id":"gotth.dns.example","manifest_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","capabilities":[],"interfaces":[],"secrets":[],"secrets":[]}`)
	if _, err := extensions.ParseGrant(grant); !errors.Is(err, extensions.ErrInvalidInput) {
		t.Fatalf("duplicate grant field was not rejected: %v", err)
	}
}
