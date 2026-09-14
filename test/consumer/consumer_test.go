package consumer

import (
	"testing"

	"github.com/gotthboard/gotth-extensions/pkg/extensions"
)

func TestPublicPackageCompiles(t *testing.T) {
	if extensions.ManifestSchema == "" || extensions.ControlName == "" {
		t.Fatal("public constants missing")
	}
}
