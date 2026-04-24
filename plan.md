# Plan: remove duplicate flag declarations from `model` and make `ast` the single source of truth

## Goal

Remove `model/flag.go` and keep one canonical declaration for the shared compiler/runtime flag set.

The canonical declaration should live on the AST side, since the AST version is the one we want to preserve stylistically.

## Current state

### 1. There are currently two representations of the same flag space

- `ast/ast.go`
  - exported `ast.Flag` bitmask constants such as `FlagPublic`, `FlagNative`, `FlagReadonly`, ...
  - private `nodeFlags` storage used by AST nodes
  - private aliases like `flagPublic`, `flagNative`, ...
- `model/flag.go`
  - `model.Flag` positional enum (`Flag_PUBLIC`, `Flag_PRIVATE`, ...)
  - consumers interpret it as a bit position via `1 << model.Flag_*`

### 2. The AST is not actually using `ast.Flag` as its storage contract

In `ast/ast.go`, the real node bits are defined like this:

- `flagPublic        nodeFlags = 1 << model.Flag_PUBLIC`
- `flagPrivate       nodeFlags = 1 << model.Flag_PRIVATE`
- ...

So the effective bit layout is currently owned by `model.Flag`, not by `ast.Flag`.

`ast.Flag` currently acts more like a parallel declaration than the real source of truth.

### 3. `model.Flag` has very little real usage outside this positional role

Direct usages of `model.Flag` / `model.Flag_*` are limited to:

- `ast/ast.go`
  - all AST node flag aliases are derived from `model.Flag_*`
- `runtime/internal/modules/registry_impl.go`
  - `nativeMethodFlagMask = 1 << model.Flag_NATIVE`
- `runtime/internal/exec/executor.go`
  - `hasFunctionFlag(flags int64, flag model.Flag)`
  - used with `model.Flag_ATTACHED`
- `bir/terminator.go`
  - `CalleeFlags common.Set[model.Flag]`
  - this field appears to be unused right now

That means `model.Flag` is not broadly embedded across the whole model layer. It is mostly a shared bit-position enum for AST/BIR/runtime.

### 4. BIR stores raw bitmasks, not `model.Flag`

`bir/model.go` stores flags as raw `int64` values:

- `BIRGlobalVariableDcl.Flags int64`
- `BIRFunction.Flags int64`
- `BIRParameter.Flags int64`

`bir/bir_gen.go` fills these from AST via `FlagsAsInt64()`.

So BIR does not need `model.Flag` as a type. It only needs a consistent bit layout and helper code to test bits.

### 5. There are other flag types in `model` that are **not** duplicates and should stay

These are separate domains and should not be removed as part of this work:

- `model.FuncSymbolFlags`
- `model.FieldDescriptorFlag`

Those are not the same flag set as `ast.Flag` / `model.Flag`.

## Important mismatch to resolve first

The two declarations are not just duplicates with different names. They also differ in content and ordering.

### Flags present in `ast` but not in `model`

- `FlagDeprecated`
- `FlagParameterized`
- `FlagIsolatedParam`
- `FlagInfer`
- `FlagEffectiveTypeDef`
- `FlagSourceAnnotation`

### Flags present in `model` but not in `ast`

- `Flag_PARALLEL`
- `Flag_NEVER_ALLOWED`

### Ordering is also different

The names overlap heavily, but the order does not. Today this is hidden because AST bit positions are explicitly derived from `model.Flag_*`.

This is the biggest migration risk:

- if we naively switch consumers from `model.Flag_*` to the current `ast.Flag` bit values,
- the serialized/stored BIR flag bits will change,
- and runtime checks like `attached` / `native` may start reading different bits than before.

## Recommended migration strategy

Use `ast` as the single declaration site, but preserve the **current effective bit layout** in the first refactor.

That means:

- move ownership of the bit layout to `ast`
- keep the same bit numbers currently observed by BIR/runtime
- remove `model/flag.go`
- do **not** change the meaning of existing serialized masks in the same step

This avoids mixing "deduplicate declarations" with "renumber all shared flags".

## Concrete plan

### Phase 1: make the AST package own the shared bit layout

Create a canonical flag definition in `ast` that is usable by both AST nodes and BIR/runtime consumers.

Recommended shape:

- keep `type Flag uint64` in `ast`
- make each exported AST flag constant represent the **actual mask bit** used everywhere
- assign explicit values so they preserve the current effective layout

Example direction:

- `FlagPublic = 1 << 0`
- `FlagPrivate = 1 << 1`
- `FlagRemote = 1 << 2`
- ...

using the bit positions that are currently encoded through `model.Flag_*`.

This is the key design point: the AST declaration should become the real wire/storage contract, not just a cosmetic duplicate.

#### Files affected

- `ast/ast.go`

#### Changes

- replace the current implicit `iota`-ordered `ast.Flag` constants with explicit mask values matching current effective bits
- redefine private `flagPublic`, `flagPrivate`, ... aliases in terms of `ast.Flag`, not `model.Flag`
- remove the comment claiming bit positions match `model.Flag iota values`
- add small helpers if useful, e.g. mask testing helpers that operate on `int64`/`uint64`

### Phase 2: switch non-AST shared consumers from `model.Flag` to `ast.Flag`

Once `ast.Flag` owns the real masks, update the few remaining consumers.

#### `runtime/internal/modules/registry_impl.go`

Replace:

- `1 << model.Flag_NATIVE`

with:

- `int64(ast.FlagNative)`

or a helper from `ast` if one is introduced.

#### `runtime/internal/exec/executor.go`

Replace:

- `hasFunctionFlag(flags int64, flag model.Flag)`

with something like:

- `hasFunctionFlag(flags int64, flag ast.Flag)`

and update the call site to use `ast.FlagAttached`.

#### `bir/terminator.go`

Replace:

- `common.Set[model.Flag]`

with either:

- `common.Set[ast.Flag]`

or remove the field entirely if it is confirmed unused and dead.

I would prefer handling this field in the same cleanup if usage stays zero.

### Phase 3: remove `model/flag.go`

After all consumers stop depending on `model.Flag`, delete:

- `model/flag.go`

Then run a repo-wide search for:

- `model.Flag`
- `Flag_`

and make sure nothing remains except unrelated flag types such as `FuncSymbolFlags` and `FieldDescriptorFlag`.

### Phase 4: clean up AST-internal duplication

Even after `model.Flag` is gone, AST still has two layers:

- exported `ast.Flag...`
- private `flag...` aliases

That is fine if the private aliases improve readability on node code, but they should become trivial aliases over AST-owned masks.

If desired, a follow-up cleanup can simplify further by:

- keeping `nodeFlags` as the storage type
- defining `flagPublic`, `flagReadonly`, ... directly from `FlagPublic`, `FlagReadonly`, ...
- leaving the exported `ast.Flag` constants as the only canonical declaration list

## Special cases to handle carefully

### 1. Preserve current bit semantics in BIR/runtime

This is the most important requirement.

`bir/bir_gen.go` writes AST flags into BIR as raw `int64` masks. Runtime reads those masks directly.

So the first refactor should preserve these meanings:

- attached function bit
n- native function bit
- public/private/global variable bits
- parameter modifier bits such as required/defaultable/rest/included
- class/type/function bits used downstream

If those change in the same commit, the migration becomes much harder to review.

### 2. Decide what to do with AST-only and model-only names

Because the declarations do not match exactly, the canonical AST declaration needs an explicit decision for:

- `Parallel`
- `NeverAllowed`

Options:

1. add `FlagParallel` and `FlagNeverAllowed` to `ast`
2. delete them if they are truly dead

Do not silently drop them before verifying no current or near-future use depends on reserving those bit positions.

Similarly, keep the AST-only flags if they are intentionally part of the frontend representation, even if they are not currently lowered into BIR/runtime behavior.

### 3. Do not conflate this with `FuncSymbolFlags`

`model.FuncSymbolFlags` is a separate 2-bit qualifier set used on function symbols/signatures.

That should remain independent unless there is a separate design decision to unify qualifier handling.

While touching this area, it is worth reviewing `ast.(*bLangInvokableNodeBase).FuncSymbolFlags()` because it currently does:

- `return model.FuncSymbolFlags(b.flags)`

That is conceptually fragile because it relies on a cast from the full AST node bitset into a separate compact symbol-qualifier bitset.

Even if it happens to work today for current call patterns, it should not be treated as part of the shared `model.Flag` removal design.

A safer long-term shape would be to construct `FuncSymbolFlags` explicitly from AST predicates:

- set `FuncSymbolFlagIsolated` if `fn.IsIsolated()`
- set `FuncSymbolFlagTransactional` if `fn.IsTransactional()`

I would include that as a small opportunistic cleanup in the same change if tests show no behavior drift.

## Validation plan

### Repo-wide checks

Run searches for:

- `model.Flag`
- `Flag_`
- `1 << model.Flag_`

Expected result after cleanup:

- no remaining shared-flag references to `model.Flag`
- only unrelated model flag types remain

### Build/tests

At minimum:

- `go test ./...`

Because this touches AST -> BIR -> runtime flag plumbing, also run corpus coverage that exercises:

- attached methods
- extern/native functions
- public/private/global declarations
- function params with required/defaultable/rest/included modifiers
- class flags such as client/service/readonly/distinct/isolated

If there are stage-specific corpus suites, the most relevant ones are:

- AST-producing phases
- BIR generation
- interpreter/runtime execution

## Suggested implementation order

1. Make `ast.Flag` preserve current effective bit positions explicitly.
2. Rebase AST private `flag...` aliases on `ast.Flag`.
3. Update runtime/BIR shared consumers to import `ast` flags.
4. Fix or remove `bir.Call.CalleeFlags` typing.
5. Delete `model/flag.go`.
6. Run search/build/tests.
7. Optionally clean up `FuncSymbolFlags()` construction.

## End state

After the refactor:

- there is a single canonical shared flag declaration in `ast`
- AST node storage uses that declaration directly
- BIR/runtime test those same AST-defined masks directly
- `model/flag.go` is gone
- `model` keeps only non-duplicate flag types such as `FuncSymbolFlags` and `FieldDescriptorFlag`

This gets rid of the duplicate declaration without changing the effective bit contract at the same time.