#!/usr/bin/env bash
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

MODE="no-coverage"
MAX_PARALLEL=3
SKIP_PATTERN='TestParseCorpusFiles|TestJBalUnitTests|TestJBalUnitBIRTests'
TIMEOUT="30m"
GO_TEST_PACKAGES_PARALLEL="$(getconf _NPROCESSORS_ONLN 2>/dev/null || echo 4)"

while [ $# -gt 0 ]; do
  case "$1" in
    --with-coverage)
      MODE="with-coverage"
      shift
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

discover_modules() {
  local modules
  modules=(".")
  while IFS= read -r mod; do
    [ -n "$mod" ] && modules+=("$mod")
  done < <(git ls-files '**/go.mod' | grep -Ev '^go\.mod$|/vendor/' | sed 's#/go.mod$##' | sort -u)
  printf '%s\n' "${modules[@]}"
}

normalize_profiles() {
  local root f line path rest tmp
  root="$(git rev-parse --show-toplevel)"
  root="${root%/}"
  shopt -s nullglob
  for f in .artifacts/coverage/*.out; do
    tmp="${f}.tmp"
    : > "$tmp"
    while IFS= read -r line || [ -n "$line" ]; do
      if [[ "$line" == mode:* ]] || [[ -z "$line" ]]; then
        printf '%s\n' "$line" >> "$tmp"
        continue
      fi
      if [[ "$line" =~ ^(.+):([0-9]+\.[0-9]+,[0-9]+\.[0-9]+[[:space:]]+[0-9]+[[:space:]]+[0-9]+)$ ]]; then
        path="${BASH_REMATCH[1]}"
        rest="${BASH_REMATCH[2]}"
        if [[ "$path" == "$root"/* ]]; then
          path="${path#"$root"/}"
        fi
        printf '%s:%s\n' "$path" "$rest" >> "$tmp"
      else
        printf '%s\n' "$line" >> "$tmp"
      fi
    done < "$f"
    mv "$tmp" "$f"
  done
  shopt -u nullglob
}

run_no_coverage() {
  local module="$1"
  (
    cd "$module"
    if [ "$module" = "." ]; then
      go test -race -count=1 -timeout "$TIMEOUT" \
        -p "$GO_TEST_PACKAGES_PARALLEL" \
        -skip "$SKIP_PATTERN" ./...
    else
      go test -count=1 -timeout "$TIMEOUT" \
        -p "$GO_TEST_PACKAGES_PARALLEL" \
        -skip "$SKIP_PATTERN" ./...
    fi
  )
}

run_with_coverage() {
  local module="$1"
  local safe_module cov_dir out_profile repo_root
  repo_root="$(pwd)"
  safe_module="${module//\//-}"
  if [ "$module" = "." ]; then
    safe_module="root"
  fi
  cov_dir="${repo_root}/.cover/${safe_module}_codecov"
  out_profile="${repo_root}/.artifacts/coverage/${safe_module}.out"
  mkdir -p "$cov_dir" "$(dirname "$out_profile")"

  (
    cd "$module"
    if [ "$module" = "." ]; then
      BAL_GOCOVERDIR="$cov_dir" go test -race -count=1 -timeout "$TIMEOUT" \
        -p "$GO_TEST_PACKAGES_PARALLEL" \
        -skip "$SKIP_PATTERN" \
        -coverpkg=./... -coverprofile="$out_profile" -covermode=atomic ./...
    else
      CODECOV_INTEGRATION_COVERDIR="$cov_dir" go test -count=1 -timeout "$TIMEOUT" \
        -p "$GO_TEST_PACKAGES_PARALLEL" \
        -skip "$SKIP_PATTERN" \
        -coverpkg=./... -coverprofile="$out_profile" -covermode=atomic ./...
    fi
  )

  if [ -n "$(ls -A "$cov_dir" 2>/dev/null)" ]; then
    go tool covdata textfmt -i="$cov_dir" \
      -o="${repo_root}/.artifacts/coverage/${safe_module}-executable.out"
  fi
}

normalize_codecov_paths() {
  local repo_root module safe_module module_path f
  repo_root="$(git rev-parse --show-toplevel)"
  while IFS= read -r module; do
    safe_module="${module//\//-}"
    [ "$module" = "." ] && safe_module="root"
    module_path="$(cd "$repo_root/$module" && go list -m -f '{{.Path}}')"
    for f in "$repo_root/.artifacts/coverage/${safe_module}.out" \
      "$repo_root/.artifacts/coverage/${safe_module}-executable.out"; do
      [ -f "$f" ] || continue
      python3 "$repo_root/.github/scripts/normalize_coverage_paths.py" \
        "$f" "$module_path" "$module"
    done
  done < <(discover_modules)
}

main() {
  local modules=()
  while IFS= read -r module; do
    modules+=("$module")
  done < <(discover_modules)

  if [ "$MODE" = "with-coverage" ]; then
    mkdir -p .artifacts/coverage
    echo "Running tests with coverage"
  else
    echo "Running tests"
  fi

  local failed=0
  local pids=()
  for module in "${modules[@]}"; do
    if [ "$MODE" = "with-coverage" ]; then
      run_with_coverage "$module" &
    else
      run_no_coverage "$module" &
    fi
    pids+=("$!")

    while [ "${#pids[@]}" -ge "$MAX_PARALLEL" ]; do
      if ! wait "${pids[0]}"; then
        failed=1
      fi
      pids=("${pids[@]:1}")
    done
  done

  for pid in "${pids[@]}"; do
    if ! wait "$pid"; then
      failed=1
    fi
  done

  if [ "$MODE" = "with-coverage" ]; then
    normalize_codecov_paths
    normalize_profiles
  fi

  return "$failed"
}

main
