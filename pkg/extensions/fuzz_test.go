package extensions

import "testing"

func FuzzParseManifest(f *testing.F) {
	seed, _ := CanonicalManifest(testManifest())
	f.Add(seed)
	f.Add([]byte(`{"schema":"bad"}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = ParseManifest(data)
	})
}

func FuzzSemver(f *testing.F) {
	for _, seed := range []string{"1.2.3", "1.0.0-alpha+1", "01.2.3", ""} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, value string) {
		_ = validSemver(value)
	})
}

func FuzzNegotiate(f *testing.F) {
	f.Add(uint32(0), uint32(2), uint32(1), uint32(4), uint32(2))
	f.Fuzz(func(t *testing.T, extensionMin, extensionMax, hostMin, hostMax, selected uint32) {
		manifest := testManifest()
		manifest.Protocols[0].MinMinor = extensionMin
		manifest.Protocols[0].MaxMinor = extensionMax
		grant := testGrant(testManifest())
		if digest, err := ManifestDigest(manifest); err == nil {
			grant.ManifestDigest = digest
		}
		host := testHost()
		host.Protocols[0].MinMinor = hostMin
		host.Protocols[0].MaxMinor = hostMax
		grant.Interfaces[0].Minor = selected
		_, _ = Negotiate(manifest, grant, host)
	})
}
