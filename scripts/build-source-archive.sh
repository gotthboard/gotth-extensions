#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

readonly archive_prefix='gotth-extensions-v1.0.0/'
readonly default_output='gotth-extensions-v1.0.0-source.tar.gz'

usage() {
  printf 'usage: %s CANDIDATE_COMMIT [OUTPUT]\n' "${0##*/}" >&2
  exit 64
}

(( $# == 1 || $# == 2 )) || usage
candidate=$1
output=${2:-$default_output}
[[ $candidate =~ ^[0-9a-f]{40}$ ]] || {
  printf 'candidate must be a full lowercase 40-hex commit ID\n' >&2
  exit 64
}

export LC_ALL=C TZ=UTC GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null
export GIT_ATTR_NOSYSTEM=1 GIT_NO_REPLACE_OBJECTS=1

root=$(git rev-parse --show-toplevel)
[[ $(pwd -P) == $(cd "$root" && pwd -P) ]] || {
  printf 'run from the repository root\n' >&2
  exit 1
}
[[ $(git rev-parse --verify 'HEAD^{commit}') == "$candidate" ]]
[[ $(git rev-parse --verify "${candidate}^{commit}") == "$candidate" ]]
[[ -z $(git status --porcelain=v1 --untracked-files=all) ]]
[[ -z $(git for-each-ref --format='%(refname)' refs/replace) ]]
info_attributes=$(git rev-parse --git-path info/attributes)
[[ ! -s $info_attributes ]]
[[ ! -e $output ]] || {
  printf 'refusing to overwrite %s\n' "$output" >&2
  exit 1
}

output_dir=$(dirname -- "$output")
[[ -d $output_dir ]]
work=$(mktemp -d)
staged=$(mktemp "$output_dir/.source-archive.XXXXXX")
cleanup() {
  rm -rf -- "$work"
  if [[ -n ${staged:-} ]]; then
    rm -f -- "$staged"
  fi
}
trap cleanup EXIT HUP INT TERM

# Any export transformation would make the inventory/tree comparison incomplete.
if git -c core.attributesFile=/dev/null grep -nE \
  '(^|[[:space:]])export-(ignore|subst)([[:space:]]|$)' \
  "$candidate" -- ':(glob)**/.gitattributes'; then
  printf 'export-ignore/export-subst attributes are forbidden\n' >&2
  exit 1
else
  status=$?
  (( status == 1 )) || exit "$status"
fi

raw_tar="$work/source.tar"
tree_records="$work/tree.records"
expected="$work/expected.inventory"
actual="$work/actual.inventory"

# A failed git archive leaves only scratch bytes; no pipeline can mask its status.
git -c core.abbrev=40 -c core.attributesFile=/dev/null -c tar.umask=0002 \
  archive --format=tar --prefix="$archive_prefix" "$candidate" > "$raw_tar"
[[ -s $raw_tar ]]
[[ $(git get-tar-commit-id < "$raw_tar") == "$candidate" ]]

# Build the expected complete path inventory independently from the Git tree.
git -c core.abbrev=40 -c core.attributesFile=/dev/null \
  ls-tree -rz -t --full-tree "$candidate" > "$tree_records"
printf '%s\n' "$archive_prefix" > "$expected"
while IFS= read -r -d '' record; do
  [[ $record == *$'\t'* ]] || {
    printf 'malformed ls-tree record\n' >&2
    exit 1
  }
  meta=${record%%$'\t'*}
  path=${record#*$'\t'}
  [[ $path != *$'\n'* && $path != *$'\r'* ]] || {
    printf 'archive paths containing CR or LF are forbidden\n' >&2
    exit 1
  }
  rest=${meta#* }
  type=${rest%% *}
  case $type in
    blob) suffix= ;;
    tree) suffix=/ ;;
    commit)
      printf 'gitlinks are forbidden in the source archive: %s\n' "$path" >&2
      exit 1
      ;;
    *)
      printf 'unexpected Git object type %s for %s\n' "$type" "$path" >&2
      exit 1
      ;;
  esac
  printf '%s%s%s\n' "$archive_prefix" "$path" "$suffix" >> "$expected"
done < "$tree_records"

# GNU tar is pinned by the release contract. Literal output is safe because CR/LF
# names were rejected above. Compare sorted multisets so duplicates also fail.
tar --quoting-style=literal --list --file="$raw_tar" > "$actual"
LC_ALL=C sort -o "$expected" "$expected"
LC_ALL=C sort -o "$actual" "$actual"
cmp --silent "$expected" "$actual"

# Compress only the validated tar, test it, and atomically expose it last.
gzip -n -9 < "$raw_tar" > "$staged"
[[ -s $staged ]]
gzip -t "$staged"
ln -- "$staged" "$output"
rm -f -- "$staged"
staged=
printf '%s\n' "$output"
