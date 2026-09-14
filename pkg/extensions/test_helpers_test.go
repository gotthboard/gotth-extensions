package extensions

func testManifest() Manifest {
	return Manifest{
		Schema:  ManifestSchema,
		ID:      "gotth.dns.godaddy",
		Name:    "GoDaddy DNS",
		Version: "1.2.3-alpha.1+build.5",
		Protocols: []VersionRange{
			{Name: ControlName, Major: 1, MinMinor: 0, MaxMinor: 2},
		},
		Interfaces: []VersionRange{
			{Name: "gotth.extensions.dns", Major: 1, MinMinor: 0, MaxMinor: 3},
		},
		Capabilities: []string{"dns.records.read", "dns.records.plan", "dns.records.apply"},
		Secrets: []SecretRequirement{
			{ID: "dns.api.key", Required: true},
			{ID: "dns.api.secret", Required: true},
			{ID: "dns.account.id", Required: false},
		},
	}
}

func testGrant(manifest Manifest) Grant {
	digest, err := ManifestDigest(manifest)
	if err != nil {
		panic(err)
	}
	return Grant{
		Schema:         GrantSchema,
		InstanceID:     "11111111-1111-4111-8111-111111111111",
		ExtensionID:    manifest.ID,
		ManifestDigest: digest,
		Capabilities:   []string{"dns.records.read", "dns.records.apply"},
		Interfaces:     []InterfaceGrant{{Name: "gotth.extensions.dns", Major: 1, Minor: 2}},
		Secrets:        []string{"dns.api.key", "dns.api.secret"},
	}
}

func testHost() HostProfile {
	return HostProfile{
		Protocols:  []VersionRange{{Name: ControlName, Major: 1, MinMinor: 1, MaxMinor: 5}},
		Interfaces: []VersionRange{{Name: "gotth.extensions.dns", Major: 1, MinMinor: 2, MaxMinor: 4}},
	}
}
