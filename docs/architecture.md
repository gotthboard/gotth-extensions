# Architecture

## Decision

`gotth-extensions` owns a small compatibility kernel and a language-neutral
control schema. It does not own product operations or process execution.

```text
consumer policy / operator confirmation / audit
                     |
                     v
       manifest + host grant + host support
                     |
                     v
        validate -> canonical digest -> negotiate
                     |
                     v
        immutable effective session contract
                     |
        authenticated transport identity
                     |
                     v
      control RPCs       seam-specific RPCs
   (handshake/health)   (defined elsewhere)
```

## Repository boundary

- `pkg/extensions` contains pure manifest, grant, negotiation, digest, and
  lifecycle mechanics.
- `proto/gotth/extensions/v1/control.proto` defines handshake and health wire
  messages. Generated clients belong in a language SDK such as `gotth-sdk`.
- `gotth-sdk` may later consume this stable protocol. It is not pulled into the
  foundation and the foundation does not generate product clients.
- Concrete provider repositories define seam-specific RPCs and mechanisms.
- Consumer applications own policy, storage, audit, administration, and
  orchestration.

### Concrete-extension topology

The repository topology is deliberately one-to-one:

```text
gotth-extensions                 shared control/compatibility foundation
gotth-extension-<slug>           exactly one deployable extension
gotth-extension-<another-slug>   exactly one other deployable extension
```

`<slug>` is a stable lowercase ASCII, hyphen-separated mechanism identifier.
The repository name identifies source ownership only. It does not authenticate
the process, grant capabilities, or establish protocol compatibility.

Each concrete repository owns its executable or image, manifest, tests,
conformance evidence, security and release records, dependencies, third-party
licenses, and rollback instructions. It must not contain a bundle of providers
or product application code. A fault, revocation, upgrade, or rollback must be
isolatable to that extension rather than forcing unrelated providers to move
together.

One extension per repository does not mean one runtime instance or one target.
A consumer may create multiple independently granted instances of the same
extension artifact—for example, separate DNS zones—without duplicating the
repository. The manifest identity names the implementation; the host-issued
instance ID and grant isolate each deployment.

Seam-specific wire contracts remain separate from the generic control
protocol. The first real implementation of a new seam must explicitly settle
where that seam contract is owned. Providers consume that admitted contract;
they do not each invent a nearly identical RPC. Shared mechanism-neutral code
is promoted only after two real consumers demonstrate it. Copying foundation
code between extension repositories or creating abstractions for imaginary
future providers is rejected.

This reverses the former placeholder ordering deliberately: Danny required the
general extension boundary first so clients and providers can be built on top.
It does not create a compatibility promise or declare `gotth-sdk` complete.

## Authority and trust

An extension manifest is untrusted input. A grant is trusted only after the
host obtains it from its own authorized configuration path. The effective
session is the intersection of:

1. a validated manifest;
2. a grant bound to the manifest digest and extension identity;
3. host-supported control and seam-specific interface ranges.

No capability exists merely because an extension declares it. A handshake
response cannot widen the negotiated session. The transport must authenticate
the expected service identity independently (for example, pinned mTLS identity
or an equivalently scoped local credential). The control schema does not carry
credentials and does not pretend self-asserted IDs are authentication.

## Canonical documents

Manifest and grant structs contain no maps. Canonicalization validates a copy,
sorts set-like slices by their full identity, encodes compact JSON with stable
field order, and appends no insignificant whitespace. SHA-256 over those exact
bytes is rendered as 64 lowercase hexadecimal characters.

Canonicalization never mutates caller-owned slices. Work is bounded by fixed
maximum counts, 32 JSON container levels, and a 64 KiB encoded-document
limit. Duplicate JSON object names are rejected before typed decoding so a
text document has only one interpretation.

## Compatibility

A protocol or interface range is `(name, major, min_minor, max_minor)`. Names
and majors are unique within a document. Negotiation requires equal names and
majors and selects the highest minor in the overlap. A grant selects one exact
version for each admitted seam-specific interface.

In V1, every seam-specific interface is a service provided by the extension
and called by the host. Extension-initiated callbacks into a host are not
silently implied by an interface grant; they require a future explicit
contract and host grant.

Control protocol and business interfaces are separate. Adding a DNS method
cannot silently change the control protocol. Unknown schemas and incompatible
major versions fail closed. Minor compatibility is admitted only through an
explicit overlapping range.

## Lifecycle

The package validates transitions only:

```text
discovered -> starting -> ready <-> degraded
                   |        |           |
                   v        +----+------+ 
                 failed          v
                              stopping -> stopped

failed  -> starting | stopped
stopped -> starting
```

`ready` and `degraded` may also enter `failed`. `starting` may enter
`stopping`. Same-state and every unlisted transition are rejected. A consumer
must record why it transitions and owns restart/backoff policy.

## Control transport

The V1 protobuf service contains only `Handshake` and `Health`.

- Handshake binds correlation ID, random challenge, expected identity,
  expected manifest/grant digests, and the negotiated session fingerprint.
- Health reports a bounded state/code pair. It carries no arbitrary detail,
  secret, raw dependency error, or product data.
- Every call requires a caller deadline. Deployments reject unauthenticated
  peers before invoking application handlers.

There is deliberately no generic `Invoke`, `Execute`, event subscription,
configuration mutation, secret retrieval, shutdown, or restart RPC.

## Secret boundary

The manifest can name secret slots and the grant can admit an exact subset.
Values are supplied through a consumer-owned, provider-scoped mechanism after
deployment identity is authenticated. Values never traverse the control RPCs
or enter foundation storage. Rotation and revocation remain host operations.

## Failure model

- Invalid documents, stale digests, unsupported versions, ungranted requests,
  wrong extension identity, and invalid lifecycle transitions fail before any
  extension operation.
- Returned errors use stable categories and never interpolate untrusted input.
- A failed extension degrades only its mechanism. It cannot partially commit
  host state because the host owns mutation and confirmation.
- Network completion can be ambiguous. Seam-specific protocols must define
  idempotency and reconciliation; the control kernel does not invent them.
- A concrete-extension release fails independently. Consumers pin it
  independently and can revoke or roll it back without changing unrelated
  extensions.

## Cost model

Validation is linear in input bytes and list counts. Canonicalization copies
bounded slices and sorts them, costing `O(n log n)` comparisons and `O(n)`
auxiliary elements. Negotiation uses bounded indexed lookups and returns a new
immutable-by-convention session value. No background goroutine, network call,
filesystem access, database access, or process launch occurs.
