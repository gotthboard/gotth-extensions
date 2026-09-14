# Distribution

Forgejo is the canonical development home. GitHub is the intended public clone
and future release endpoint when explicitly created and synchronized.

This bootstrap does not create a GitHub mirror, tag, release, or deployment.

The plural `gotth-extensions` repository distributes only the shared
foundation. Every concrete extension is distributed from its own canonical
Forgejo repository named `gotth-extension-<slug>` and carries an independent
artifact, manifest digest, changelog, security record, tag/release history,
and rollback record. A public GitHub clone is created only under explicit
owner authority and exact-ref synchronization; naming this policy does not
create any concrete repository or mirror.
