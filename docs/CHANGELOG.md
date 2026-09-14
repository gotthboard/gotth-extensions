# Changelog

## Unreleased

### 2026-09-13 22:05 CDT — Define the general extension foundation

Commit: current commit; hash assigned by Git after commit

Affected files:

- repository policy and release documents
- product requirements, architecture, implementation specification, feature
  plan, verification contract, and workflow manifest

Explanation:

Record Danny's explicit decision to create `gotth-extensions` before concrete
providers and before `gotth-sdk`. Define a provider-free, out-of-process,
capability-scoped control and compatibility kernel grounded in GOTTH Mail's
real DNS, notification, backup, webmail, certificate, and import seams.

The host remains authoritative for policy, credentials, state, audit,
confirmation, process supervision, and product operations. The foundation
does not contain a generic invocation API or event bus.

Verification:

- prior GOTTH Extensions queue contract reconciled with the new owner ordering
- current `gotth-sdk`, `gotth-webhooks`, `gotth-authentik`, GOTTH Board, and
  GOTTH Mail contracts inspected
- implementation and runtime verification follow in later commits

Risks / non-goals:

- no provider, SDK, product change, credential, DNS mutation, deployment, tag,
  release, or public compatibility promise
