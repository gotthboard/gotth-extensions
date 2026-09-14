# gotth-extensions

`gotth-extensions` is the general, out-of-process extension control and
compatibility foundation for GOTTH applications.

It defines how a host validates an extension manifest, grants an exact subset
of requested capabilities, negotiates compatible protocol/interface versions,
and tracks lifecycle state. Concrete DNS, notification, backup, webmail,
certificate, and import mechanisms build on top of this boundary; they do not
belong in this repository.

## Boundary

- Extensions run out of process. Arbitrary in-process Go plugins are forbidden.
- The host remains authoritative for authentication, authorization, business
  policy, state admission, mutation, audit, retry, and rollback decisions.
- The control protocol provides handshake and health only. It is not a generic
  invocation API or event bus.
- Every usable capability is explicitly granted. Declaration is not authority.
- Secret values never appear in manifests, grants, handshakes, health results,
  errors, or logs.
- Extensions receive no direct product database, Docker socket, or broad host
  filesystem access.

The first unreleased foundation contains no provider, runner, supervisor,
deployment, tag, release, or compatibility promise. See the
[product requirements](docs/prd.md), [architecture](docs/architecture.md), and
[implementation specification](docs/implementation-spec.md).

Owner-authored contents are licensed under the [MIT License](LICENSE).
