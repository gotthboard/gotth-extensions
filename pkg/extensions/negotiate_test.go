package extensions

import (
	"errors"
	"reflect"
	"testing"
)

func TestNegotiate(t *testing.T) {
	manifest := testManifest()
	grant := testGrant(manifest)
	host := testHost()
	originalGrant := testGrant(manifest)

	session, err := Negotiate(manifest, grant, host)
	if err != nil {
		t.Fatal(err)
	}
	if session.ExtensionID != manifest.ID || session.ExtensionVersion != manifest.Version || session.InstanceID != grant.InstanceID {
		t.Fatalf("wrong identity: %#v", session)
	}
	if session.Control != (ProtocolVersion{Name: ControlName, Major: 1, Minor: 2}) {
		t.Fatalf("wrong control version: %#v", session.Control)
	}
	if !reflect.DeepEqual(session.Capabilities, []string{"dns.records.apply", "dns.records.read"}) || !reflect.DeepEqual(session.Secrets, []string{"dns.api.key", "dns.api.secret"}) {
		t.Fatalf("session sets not canonical: %#v", session)
	}
	if len(session.Fingerprint) != 64 || len(session.ManifestDigest) != 64 || len(session.GrantDigest) != 64 {
		t.Fatalf("missing binding digest: %#v", session)
	}
	if !reflect.DeepEqual(grant, originalGrant) {
		t.Fatal("negotiation mutated grant")
	}
	repeat, err := Negotiate(manifest, grant, host)
	if err != nil || repeat.Fingerprint != session.Fingerprint {
		t.Fatalf("non-deterministic session: %#v %v", repeat, err)
	}
}

func TestNegotiateFailures(t *testing.T) {
	manifest := testManifest()
	grant := testGrant(manifest)
	host := testHost()

	tests := []struct {
		name   string
		mutate func(*Manifest, *Grant, *HostProfile)
		target error
	}{
		{"invalid manifest", func(m *Manifest, _ *Grant, _ *HostProfile) { m.Name = "" }, ErrInvalidInput},
		{"invalid grant", func(_ *Manifest, g *Grant, _ *HostProfile) { g.InstanceID = "bad" }, ErrInvalidInput},
		{"invalid host", func(_ *Manifest, _ *Grant, h *HostProfile) { h.Protocols = nil }, ErrInvalidInput},
		{"wrong extension", func(_ *Manifest, g *Grant, _ *HostProfile) { g.ExtensionID = "gotth.dns.other" }, ErrStaleBinding},
		{"stale digest", func(_ *Manifest, g *Grant, _ *HostProfile) {
			g.ManifestDigest = string(make([]byte, 64))
		}, ErrInvalidInput},
		{"valid stale digest", func(_ *Manifest, g *Grant, _ *HostProfile) {
			g.ManifestDigest = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		}, ErrStaleBinding},
		{"control mismatch", func(_ *Manifest, _ *Grant, h *HostProfile) { h.Protocols[0].MinMinor = 3 }, ErrIncompatibleVersion},
		{"capability", func(_ *Manifest, g *Grant, _ *HostProfile) { g.Capabilities[0] = "dns.records.delete" }, ErrUngrantedCapability},
		{"secret unknown", func(_ *Manifest, g *Grant, _ *HostProfile) { g.Secrets[0] = "dns.api.unknown" }, ErrUngrantedSecret},
		{"secret required", func(_ *Manifest, g *Grant, _ *HostProfile) { g.Secrets = []string{"dns.api.key"} }, ErrUngrantedSecret},
		{"interface extension", func(_ *Manifest, g *Grant, _ *HostProfile) { g.Interfaces[0].Minor = 9 }, ErrUngrantedInterface},
		{"interface host", func(_ *Manifest, g *Grant, _ *HostProfile) { g.Interfaces[0].Minor = 1 }, ErrUngrantedInterface},
		{"interface missing", func(_ *Manifest, g *Grant, _ *HostProfile) { g.Interfaces[0].Name = "gotth.extensions.backup" }, ErrUngrantedInterface},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := manifest
			m.Protocols = cloneRanges(manifest.Protocols)
			m.Interfaces = cloneRanges(manifest.Interfaces)
			m.Capabilities = cloneStrings(manifest.Capabilities)
			m.Secrets = append([]SecretRequirement(nil), manifest.Secrets...)
			g := grant
			g.Capabilities = cloneStrings(grant.Capabilities)
			g.Interfaces = append([]InterfaceGrant(nil), grant.Interfaces...)
			g.Secrets = cloneStrings(grant.Secrets)
			h := host
			h.Protocols = cloneRanges(host.Protocols)
			h.Interfaces = cloneRanges(host.Interfaces)
			test.mutate(&m, &g, &h)
			if _, err := Negotiate(m, g, h); !errors.Is(err, test.target) {
				t.Fatalf("got %v, want %v", err, test.target)
			}
		})
	}
}

func TestNegotiationHelpers(t *testing.T) {
	if _, ok := negotiateControl(nil, nil); ok {
		t.Fatal("missing control negotiated")
	}
	left := []VersionRange{{Name: ControlName, Major: 1, MinMinor: 2, MaxMinor: 4}}
	right := []VersionRange{{Name: ControlName, Major: 2, MinMinor: 2, MaxMinor: 4}}
	if _, ok := negotiateControl(left, right); ok {
		t.Fatal("wrong major negotiated")
	}
	if uint32Decimal(0) != "0" || uint32Decimal(4294967295) != "4294967295" {
		t.Fatal("uint32 formatting wrong")
	}
	if rangeKey("x", 12) != "x\x0012" {
		t.Fatal("range key wrong")
	}
}

func TestNegotiationSortsMultipleInterfaces(t *testing.T) {
	manifest := testManifest()
	manifest.Interfaces = append(manifest.Interfaces,
		VersionRange{Name: "gotth.extensions.backup", Major: 1, MinMinor: 0, MaxMinor: 1},
	)
	grant := testGrant(manifest)
	grant.Interfaces = append(grant.Interfaces,
		InterfaceGrant{Name: "gotth.extensions.backup", Major: 1, Minor: 1},
	)
	host := testHost()
	host.Interfaces = append(host.Interfaces,
		VersionRange{Name: "gotth.extensions.backup", Major: 1, MinMinor: 1, MaxMinor: 2},
	)
	session, err := Negotiate(manifest, grant, host)
	if err != nil {
		t.Fatal(err)
	}
	want := []InterfaceGrant{
		{Name: "gotth.extensions.backup", Major: 1, Minor: 1},
		{Name: "gotth.extensions.dns", Major: 1, Minor: 2},
	}
	if !reflect.DeepEqual(session.Interfaces, want) {
		t.Fatalf("interfaces not sorted: %#v", session.Interfaces)
	}
}
