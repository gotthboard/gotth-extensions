# Verification requirements

Foundation admission requires:

1. `gofmt` clean source;
2. `go vet ./...`;
3. `go test ./...`;
4. `go test -race ./...`;
5. repeated focused tests for canonicalization, negotiation, and lifecycle;
6. statement coverage with every meaningful uncovered branch recorded;
7. short fuzz smoke for public parsers/negotiation;
8. `protoc` descriptor compilation of the V1 control schema;
9. external-consumer compile against the public package;
10. clean-clone rerun and `git diff --check`;
11. cold review confirming no provider, generic invocation API, secret value,
    product policy, in-process loader, direct database access, or runtime
    side effect entered the foundation;
12. repository-boundary review confirming that concrete extensions are
    one-per-repository under valid `gotth-extension-<slug>` names, with
    independent artifact, evidence, release, and rollback ownership;
13. negative review confirming that a repository name is never accepted as
    authority, authentication, compatibility, or product admission.
14. V1 public Go documentation/API baseline, shape/signature, constant, sentinel,
    strict raw-parser (including invalid UTF-8 bytes and recursive exact-key,
    alias-collision, and duplicate rejection), strict fixture-decoder negatives
    at every object level, ordinary/control-only canonical and complete-Session
    vectors, input/Session detachment in both directions, Session
    mutation/revalidation, lifecycle-edge, and protobuf-descriptor compatibility
    gates.

Release admission additionally requires the exact-candidate, generated-SDK,
cross-language, security, signed immutable-payload, forge byte-parity, and external module
checks in `docs/RELEASING.md`. Passing the local compatibility gates does not
mean those release gates passed.

No live extension, credential, DNS record, product repository, deployment, or
remote runtime is part of this verification.

The planned host-management contract will additionally require negative proof
that metadata cannot inject markup, script, style, redirects, arbitrary form
actions, secret values, or new authority; that stale or privilege-changing
updates cannot bypass preview and confirmation; that disable revokes routing
before shutdown; and that one product cannot enumerate or mutate another
product's registry, grants, secrets, or audit state.
