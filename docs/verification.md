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
    side effect entered the foundation.

No live extension, credential, DNS record, product repository, deployment, or
remote runtime is part of this verification.
