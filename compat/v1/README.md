# V1 compatibility fixtures

These files pin the compatibility surface proposed for foundation `v1.0.0`.
They are gates, not evidence that a release exists.

- `canonical-vectors.json` is the normative vector set for complete ordinary and
  control-only manifest, grant, host, and Session values; exact canonical
  manifest/grant/session bytes and SHA-256 values; canonical Unicode display
  names including U+2028/U+2029 and a valid UTF-16 surrogate-pair escape; text or
  hexadecimal raw-byte parser inputs (including invalid UTF-8 and recursive
  exact-key negatives); manifest-valid but host-invalid interface negotiation;
  strict fixture-decoder negatives; and stable error categories.
- `go-api.txt` is the machine-diffable output of
  `go doc -all ./pkg/extensions` and includes every exported declaration and
  its documentation.
- `control-descriptor.sha256` pins the descriptor emitted by the documented
  `protoc` command for `gotth.extensions.v1`.

`pkg/extensions/compatibility_v1_test.go` strictly decodes this fixture with
exact case-sensitive object names and exact duplicates forbidden recursively,
unknown fields forbidden, and one top-level value required. It compile-checks
functions, reflects public struct shapes, executes vectors and limits, checks
every lifecycle edge,
and proves input/Session slice detachment in both directions plus the Session
deep-copy/revalidation rule. `make api-compat` and
`make proto-compat` reject API/descriptor drift. Updating a baseline does not
make a breaking change compatible; `docs/compatibility-v1.md` decides whether
a new major/schema/protocol version is required.
