# Releasing

The repository remains unreleased. The proposed first stable version is
`v1.0.0`, with compatibility defined by `docs/compatibility-v1.md`. Drafting
that contract is not a tag, support promise, or publication authority.

## Candidate and SDK prerequisites

The implementation baseline is
`dfa3d975369c436ffdee98e5f6468eb2be268c54`. The eventual candidate is one
reviewed descendant containing the release contract, compatibility gates, and
exact-ref evidence. Record its full commit and tree IDs before admission.

Historical consumer observations do not prove current branch state, exact
candidate compatibility, or stable adoption and are not restated here as
current-head claims. A concrete extension cannot demonstrate exact stable
adoption until the foundation tag, assets, module version, and SDK are public.
The order is exact: publish the SDK, publish the foundation, accept and publish
webhook, then accept and publish GoDaddy. These adopter stages are not foundation
pre-tag or pre-publication gates, and GoDaddy does not gate webhook. Mail, other
product hosts, and design-only repositories remain outside this release and
adoption procedure.

The SDK is released independently and first. Its admitted release must bind the
final foundation candidate commit and unchanged control source/descriptor,
regenerate byte-for-byte with pinned tools, compile for Go, pass the named
non-Go conformance lane, publish dependency/license/provenance records on both
forges, and resolve from a clean external consumer without a replace directive.
The SDK version need not equal v1.0.0.

The inspected uncommitted F2 SDK worktree is based on
4b2d7c37bae3f18dafd6130f3846fc527fe01285, copies foundation source commit
dfa3d975369c436ffdee98e5f6468eb2be268c54, owns generated Go files under
gen/go/gotth/extensions/v1, and reports descriptor
62a549b29a030f0decb46b2b1c2f382bbbc75430ba0b3d7c729bd21528b505cf plus six
Go/C++ handshake-health conformance cases. Those paths and hashes coordinate
this contract, but the dirty worktree is not an SDK revision or release. F2 must
rebind source provenance to the final admitted F1 commit and be independently
committed, reviewed, verified, tagged, mirrored, and published.
The foundation tag remains blocked until all of these pre-tag gates pass at
exact refs:

1. admitted F1 commit and clean-clone full verification/security evidence;
2. published SDK-first tag/version and byte-identical Forgejo/GitHub assets;
3. Go plus at least one non-Go runtime passing the shared canonical and
   handshake/health vectors;
4. final dependency, license, and SBOM inputs plus the reviewed provenance and
   release-asset construction procedure; and
5. explicit owner authority for the signed tag, post-tag evidence signatures,
   and publication.

Final provenance/manifest/checksum bytes and their signatures necessarily follow
tag creation. Publication acceptance remains blocked until those assets exist and
exact tag, release metadata, and byte parity is proved between Forgejo and GitHub.

## Generated SDK tree hash

`release-manifest.json` binds the SDK repository, full commit, tag/version,
generator identities, and `generated_tree_sha256`. Compute that hash over the
SDK repository subtree `gen/go/gotth/extensions/v1` as follows:

1. Include every and only regular file recursively; reject symlinks, special
   files, and non-`0644` permission bits. Empty directories are ignored.
2. Encode each path relative to the SDK repository root, not the subtree root,
   using literal `/` separators and UTF-8 bytes. Reject an absolute path,
   invalid UTF-8, backslash, NUL, CR, LF, an empty segment, `.`, or `..`.
   Do not Unicode-normalize or case-fold.
3. Sort entries by unsigned bytewise path order.
4. Initialize SHA-256 with the 23 bytes `gotth-generated-tree-v1` followed by
   one NUL byte.
5. For each entry append: four-byte unsigned big-endian path-byte length; path
   bytes; the six ASCII bytes 100644; one NUL byte; eight-byte unsigned big-endian
   file-byte length; and exact file bytes.
6. Lowercase hexadecimal encoding of the final digest is
   `generated_tree_sha256`.

The release evidence retains the exact path inventory, per-file SHA-256 values,
and an independent implementation/check of this algorithm. A tar hash, Git tree
ID, concatenated `sha256sum` text, or filesystem traversal order is not a
substitute.

## Immutable `v1.0.0` identity and attachments

The immutable release identity is an annotated Git tag `v1.0.0`, signed by the
owner-approved OpenPGP release key and peeled to the admitted candidate. Both
forges expose the same tag object and peeled commit. If the approved key cannot
sign or verification cannot prove its exact fingerprint, release is blocked; an
unsigned tag or signature waiver is not admissible for `v1.0.0`.

Both release objects attach byte-identical files with these exact names:

- `gotth-extensions-v1.0.0-source.tar.gz`: exact output of the pinned archive
  construction below;
- `control.proto`: candidate
  `proto/gotth/extensions/v1/control.proto` bytes;
- `control.pb`: descriptor from the documented `protoc --proto_path=.
  --include_imports --descriptor_set_out=control.pb` command at the candidate;
- `canonical-vectors.json`: byte-identical copy of the exact repository path
  `compat/v1/canonical-vectors.json`;
- `go-api.txt`: byte-identical copy of `compat/v1/go-api.txt`;
- `LICENSE`;
- `dependency-licenses.json`;
- `sbom.spdx.json`;
- `provenance.json`;
- `release-manifest.json`;
- `SHA256SUMS`;
- `release-signing-key.asc`;
- `provenance.json.asc`;
- `release-manifest.json.asc`; and
- `SHA256SUMS.asc`.

### Deterministic source archive

The archive producer is exactly `scripts/build-source-archive.sh` from the
admitted candidate, run with Bash from the pinned release image, Git 2.56.0,
GNU tar 1.35, and GNU gzip 1.15. The script bytes and SHA-256, each executable
SHA-256, package coordinates, `--version` output, operating-system image digest,
and locale are release-manifest/provenance inputs. From a clean clone whose
`HEAD` equals the peeled tag commit, with no replacement refs, local/system/global
Git attributes, or working-tree input, run exactly:

```sh
scripts/build-source-archive.sh \
  "$CANDIDATE_COMMIT" gotth-extensions-v1.0.0-source.tar.gz
```

`CANDIDATE_COMMIT` is the full peeled commit ID already bound by the signed tag;
passing the commit (not its tree ID) fixes archive entry mtimes to the commit
timestamp. The script uses `set -euo pipefail`, writes `git archive` to scratch
without a pipeline, checks the embedded commit ID, and derives an independent
complete path inventory from `git ls-tree`. It rejects export transformations,
gitlinks, CR/LF path names, missing, extra, or duplicate archive entries. Only
after GNU tar inventory comparison succeeds does it compress the validated tar,
run `gzip -t`, and atomically hard-link a same-directory temporary file to the
previously absent requested output. Thus a failed `git archive`, inventory check,
compression, or no-clobber publication step cannot leave accepted output bytes.
The candidate must contain no `.gitattributes` rule that uses `export-subst` or
`export-ignore`.

Git supplies tree paths, regular file/symlink modes and bytes in Git archive
order; gzip `-n` omits the original name and timestamp, and `-9` fixes the
compression level. Do not modify the script or post-process the tar stream. Two
isolated clean clones, run by different operators or jobs with independently
installed binaries matching the pinned executable hashes, must produce
byte-identical archives. Retain both commands, validated inventories, sizes,
SHA-256 values, and `cmp` result. A mismatch blocks release; selecting either
output despite a mismatch is forbidden.

`dependency-licenses.json` is UTF-8 JSON with duplicate names forbidden. It
has separate `runtime` and `build` arrays. Every entry contains immutable
name, version, source coordinate, SPDX identifier/expression, relationship to
the release, and retained license-file SHA-256. Entries sort bytewise by name,
then version, then source coordinate. The final foundation runtime graph must
explicitly record that it contains no non-standard-library modules; Go,
protobuf, generators, scanners, and archive/SBOM tools remain build entries.
An unidentified or incompatible license blocks release.

`sbom.spdx.json` describes the exact source archive and final module/build
inputs using an admitted SPDX JSON version. Its namespace and creation metadata
are immutable release values, not mutable branch URLs. Package checksums and
relationships agree with `dependency-licenses.json` and the final module graph.

`provenance.json` records the candidate commit/tree, clean-clone source,
builder/tool versions and executable hashes, exact commands, start/end times,
source inputs, SDK identity/tree hash, and hashes of every pre-manifest
attachment (source archive,
proto, descriptor, vectors, API baseline, license, dependency inventory, and
SBOM). It records retained verification/security evidence IDs and contains no
credentials or private/mutable endpoint as an identity. It deliberately does
not hash `release-manifest.json` or `SHA256SUMS`.

`release-manifest.json` binds version; signed tag object and peeled commit/tree;
all pre-manifest attachment hashes; module path/version and module sum; Go and
protobuf toolchain identities; source/descriptor/vector/API hashes; SDK
repository, commit, tag/version, module sum, generators and generated-tree hash;
license/inventory/SBOM hashes; provenance hash; and the owner-approved OpenPGP
primary-key fingerprint, signing subkey fingerprint, exact signer UID, and
detached-signature filenames. All object IDs and fingerprints are full length.
Missing values are blockers, never null placeholders.

All JSON attachments use UTF-8, one trailing LF, duplicate names forbidden, and
stable bytewise key/order rules documented by their generator. `SHA256SUMS`
contains lowercase SHA-256, two ASCII spaces, and the exact filename for every
attachment then present except itself and `SHA256SUMS.asc`, one LF-terminated
line per file sorted by unsigned bytewise filename. Thus it includes
`release-signing-key.asc`, `provenance.json.asc`, and
`release-manifest.json.asc`. This ordering avoids checksum cycles: provenance
binds only pre-manifest files; manifest binds provenance and pre-manifest files;
the first two detached signatures bind their exact JSON bytes; checksums bind
the manifest, those signatures, and every earlier attachment; and the final
detached signature binds the exact checksum file.

The same owner-approved OpenPGP key that signs the annotated tag must make
binary-input, ASCII-armored detached signatures over the exact published bytes
of `provenance.json`, `release-manifest.json`, and `SHA256SUMS`, after the tag
exists. Export its minimal public key as `release-signing-key.asc`. The owner
authorization record fixes the primary-key fingerprint, allowed signing-subkey
fingerprint, signer UID, and release version before signing; a key bundled only
inside the release is not self-authenticating. Record GnuPG version, executable
hash, signature creation time, and exact commands. Use the full approved signing
subkey fingerprint with `!`; abbreviated key IDs are forbidden:

```sh
gpg --batch --armor --export-options export-minimal \
  --export "$PRIMARY_FPR!" > release-signing-key.asc
gpg --batch --yes --armor --local-user "$SIGNING_SUBKEY_FPR!" \
  --output provenance.json.asc --detach-sign provenance.json
gpg --batch --yes --armor --local-user "$SIGNING_SUBKEY_FPR!" \
  --output release-manifest.json.asc --detach-sign release-manifest.json
# Construct SHA256SUMS only after the preceding signatures exist.
gpg --batch --yes --armor --local-user "$SIGNING_SUBKEY_FPR!" \
  --output SHA256SUMS.asc --detach-sign SHA256SUMS
```

Verification imports `release-signing-key.asc` into a new empty keyring, proves
its full primary and signing-subkey fingerprints equal the separately approved
record, and runs `gpg --batch --status-fd=1 --verify SIGNATURE FILE` for all
three pairs. Each verification must emit `VALIDSIG` for the approved signing
subkey and expected primary fingerprint; expiration, revocation, bad signature,
wrong fingerprint/UID, or unverifiable status blocks release. Download both
forges independently and require byte-for-byte equality of the key, three signed
files, all three signatures, every checksum-covered attachment, and release
metadata before accepting publication.

Use an empty mode-0700 `GNUPGHOME`; import only `release-signing-key.asc`; capture
`gpg --batch --with-colons --fingerprint --fingerprint` and prove its exact
`fpr` and `uid` records equal the separately approved identity; then run exactly:

```sh
gpg --batch --status-fd=1 --verify provenance.json.asc provenance.json
gpg --batch --status-fd=1 --verify release-manifest.json.asc release-manifest.json
gpg --batch --status-fd=1 --verify SHA256SUMS.asc SHA256SUMS
```

## SDK-first construction and publication sequence

The order is normative; constructing the final manifest before the tag exists
is impossible because the manifest binds the tag object.

1. Admit and commit F1; record exact commit/tree. Do not tag it yet.
2. Rebind F2 SDK source provenance to that commit, admit and commit the SDK,
   then—with separate owner authority—create its immutable tag and publish the
   SDK release on both forges. Verify tag/assets and a clean external import.
3. Run foundation exact-candidate clean-clone, full/race/repeat/shuffle/fuzz,
   descriptor/API, consumer, vulnerability, secret-history, dependency,
   license, and provenance gates. Seal logs with command, tool version,
   timestamps, exit status, candidate ID, and output hash.
4. Obtain explicit owner authority for foundation `v1.0.0`. Create the signed
   annotated tag exactly once at the admitted candidate, verify it locally, push
   the same tag object to both forges, and verify tag-object and peeled-commit
   equality. Never recreate or move it.
5. Make two isolated fresh public clones of that exact mirrored tag. Independently
   construct the source archive with the pinned commands and prove exact byte
   equality before selecting those bytes for the release. Construct each other
   pre-manifest attachment once. Resolve the public module version without a
   replace directive and record its module sum before constructing the manifest.
6. Construct `provenance.json`, then `release-manifest.json` using the now-known
   tag object and already-published SDK identity. Owner-sign those two exact files,
   construct `SHA256SUMS` including their detached signatures, and owner-sign the
   checksum file. Verify every schema, inventory, hash, owner identity/signature,
   license, and secret scan in a separate staging directory.
7. Upload the same attachment bytes and release metadata to Forgejo and GitHub;
   do not rebuild per forge.
8. Download both public release objects and every attachment independently.
   Verify tag object, peeled commit, metadata, inventory, size, SHA-256,
   signature/provenance, and byte equality.
9. Independently verify a clean external module resolves v1.0.0 without a
   replace directive and produces the same module sum recorded in the manifest.
   Only then may current-status documents say released.
10. After step 9, update `gotth-extension-webhook` in separately authorized work
    to pin the public foundation `v1.0.0` and published SDK identity, remove any
    release-test replacement, and retain clean module, canonical,
    handshake/health, and hostile-path evidence at the exact webhook commit.
    Accept that revision independently, publish it, and verify the public
    identity before beginning GoDaddy adoption. Webhook acceptance and
    publication do not wait for GoDaddy.
11. Only after step 10 is accepted and published, update
    `gotth-extension-godaddy-dns` in its own separately authorized work to pin
    the same public foundation and SDK identities. Retain the same categories of
    exact-revision evidence, accept the GoDaddy revision independently, publish
    it, and verify its public identity. Neither adopter result can retroactively
    serve as a pre-tag gate for the already published foundation. Design-only
    repositories are not part of steps 10 or 11.

## Failure and rollback

Before step 4, discard staging output and repair or replace the candidate. Once
the tag is created, never delete, overwrite, recreate, or retarget it. A failure
after tag creation leaves an immutable incomplete/affected release identity:
stop publication, record the failure, and use a new patch version after the same
gates. Do not repair assets in place under the old version.

Foundation rollback is consumer-first, not repository deletion: an adopter
that has completed step 10 or 11 pins its last independently admitted foundation and
SDK identities, regenerates or selects matching bindings, reruns its own
admission gates, and revokes stale sessions/grants. The `v1.0.0` release record
names a tested commit-pin fallback before tagging. Rollback of foundation
identity never implies rollback of every concrete extension artifact, and a
foundation release never publishes concrete provider artifacts.
