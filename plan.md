# Plan: `lock` statement support

This plan implements the `lock` statement end-to-end — AST construction,
semantic validation, desugaring, BIR generation and interpretation.

## Confirmed design decisions

- **Locks are global.** Every restricted variable (module-level
  isolated variable or non-final field of an isolated class) maps to a
  single program-wide re-entrant mutex. All instances of an isolated
  class share the same lock for a given field.
- **Abrupt-exit semantics.** Normal completion, `return`, `break`,
  `continue`, and `fail` each emit an explicit `LockEnd` along their exit
  edge (rewritten in desugar). A *Ballerina panic* releases every lock
  currently held by the strand as part of unwinding (handled by the
  runtime). Go panics are not in scope.
- **Runtime `Environment`.** Introduce a new `runtime.Environment` that owns
  the existing `*modules.Registry`, the strand-id allocator, and the global
  lock table. `Runtime` holds an `*Environment`; `exec.Context` is given a
  pointer to the same `Environment`.

The work is split into self-contained phases. Each phase ends in a green test
suite (the corpus tests for the relevant stages). Phases can be reviewed/merged
incrementally.

---

## Phase 0 — Test fixtures and skeleton

Goal: make the parser-level and AST-level tests for `lock` pass before any
semantics are added, so subsequent phases get immediate corpus feedback.

1. Add `corpus/bal/subset8/08-lock/` containing the *positive* `.bal` files we
   intend to support:
   - `lock-module-var-1-v.bal` — module-level isolated int incremented inside a
     lock from an isolated function.
   - `lock-class-field-1-v.bal` — isolated class with a private non-final
     `int` field mutated under `self`-lock.
   - `lock-reentrant-1-v.bal` — recursive isolated function that re-enters its
     own lock (tests the re-entrant runtime path).
   - `lock-with-locals-1-v.bal` — lock body declaring its own locals and
     reading/writing them freely.

   And the *negative* files:
   - `lock-multiple-restricted-1-e.bal` — two distinct module-level isolated
     vars touched in one lock.
   - `lock-nested-1-e.bal` — a `lock` statement lexically inside the body of
     another `lock` statement (including via the body of a function called
     from within the outer lock? no — only *lexical* nesting; cross-function
     nesting is a runtime concern handled by re-entrancy).
   - `lock-non-isolated-assign-1-e.bal` — assignment inside lock to a captured
     mutable variable that isn't the restricted var and isn't isolated.
   - `lock-module-var-outside-1-e.bal` — read/write of a module-level isolated
     variable outside any lock.
   - `lock-public-mutable-field-1-e.bal` — non-final field of isolated class
     declared without `private`.
   - `lock-nonself-field-1-e.bal` — `obj.foo = ...` (where `obj` is not `self`)
     inside a lock.
   - `isolated-var-init-1-e.bal` — module-level `isolated` variable whose
     initializer expression captures a mutable global.

2. Add expected outputs under `corpus/parser`, `corpus/ast`, `corpus/cfg`,
   `corpus/desugared`, `corpus/bir`, `corpus/integration` for each of the above
   (using `-update` after each phase that completes the stage).

3. Phase-0 deliverable: only the parser + AST corpus expected files exist and
   are produced correctly (parser already creates `LockStatementNode`; we wire
   AST below). The other expected files are added at the phase that produces
   them.

---

## Phase 1 — AST node

Files: `ast/statements.go`, `ast/kinds.go`, `ast/walk.go`,
`ast/pretty_printer.go`, `ast/node_builder.go`.

1. Define
   ```go
   type BLangLock struct {
       bLangStatementBase
       Body BLangBlockStmt
       // RestrictedSymbol is filled in by the lock analyzer (Phase 4).
       // For a module-level isolated variable it is the variable's symbol;
       // for a non-final field of an isolated class it is the field's
       // SymbolRef (i.e. BLangSimpleVariable.Symbol() of the field).
       RestrictedSymbol model.SymbolRef
   }
   ```
   Add `LockNode` interface and the `_ LockNode = &BLangLock{}` /
   `_ BLangNode = &BLangLock{}` assertions; add a `NodeKind` constant.

   **On-fail clause is out of scope for this PR.** If the parser's
   `LockStatementNode` carries a non-empty `OnFailClause`, the AST builder
   raises an `unimplementedErr` ("`on fail` clause on `lock` is not yet
   supported") and produces a `BLangLock` without the clause. We will tackle
   `on fail` for `lock` in a follow-up once the basic semantics are in
   place. All later phases assume `BLangLock` has no on-fail clause.

2. `node_builder.go::TransformLockStatement` — currently panics. Build a
   `BLangLock` from the parser's `LockStatementNode`. Body is constructed
   via the existing block-statement builder. If the parser produced an
   on-fail clause, emit the unimplemented error described above.

3. `walk.go`: visit Body.

4. `pretty_printer.go`: print `lock { ... }`.

5. Run/refresh `corpus/ast` expected files for the new fixtures. Do **not**
   add fixtures that exercise `on fail` on a lock at this stage.

---

## Phase 2 — Symbol & type resolution

Files: `semantics/symbol_resolver.go`, `semantics/type_resolver.go`.

1. The lock body is just a block; reuse the block visitor. A lock does not
   introduce its own scope distinct from the body block, so no new scope is
   required here — the body's `BLangBlockStmt` already gets a block scope.

2. Type resolver: walk into Body.

3. Refresh `corpus/cfg` only after Phase 6. (Symbol/type resolution has no
   dedicated corpus output beyond the AST.)

---

## Phase 3 — Refactor isolation analysis

Prerequisite for Phase 4. Pure refactor; no new diagnostics.

In `semantics/semantic_analyzer.go`, extract the body of the existing
`isIsolatedFuncInner` into

```go
func isolationCheckNode[A analyzer](
    a A, node ast.BLangNode, extraAllowed map[model.SymbolRef]struct{},
)
```

Re-implement `isIsolatedFuncInner` as `isolationCheckNode(a, node, nil)`
so current behaviour is preserved.

Add a `BLangLock` case to the walker that does **not** descend into
`BLangLock.Body`. Reason: when `isIsolatedFuncInner` walks an isolated
function body containing a `lock`, the lock body may legitimately
reference its restricted variable, which would otherwise be flagged.
The lock analyzer (Phase 4) re-runs isolation analysis on the body
separately with the restricted variable allow-listed.

---

## Phase 4 — Lock analyzer

New file `semantics/lock_analyzer.go`.

Lock analysis is triggered from the existing `functionAnalyzer.Visit`
switch in `semantics/semantic_analyzer.go`. Add a `*ast.BLangLock` case
that invokes `analyzeLock(fa, n)` and then returns `fa` so the visitor
continues walking the body (this is what makes the nested-lock check
in step 1 fire).

For a `*ast.BLangLock` `L`, `analyzeLock` performs:

1. **Nested lock check.** If `L` is lexically inside another
   `BLangLock.Body`, emit `"lock statement cannot be nested inside
   another lock statement"` and do not descend. (Cross-function
   re-entrancy is a runtime concern and is not flagged here.)

2. **Identify the restricted variable.** Walk `L.Body` and collect every
   reference (read or write) that is one of:
   - a module-level variable with `IsIsolated() == true`;
   - `self.f` where `f` is not `final` with a `readonly` type.

   The comparison key is the variable's unnarrowed `SymbolRef` (for
   module vars) or the field's `BLangSimpleVariable.Symbol()` (for
   `self.f`). The first occurrence is the restricted variable; every
   other occurrence whose key differs emits `"more than one restricted
   variable referenced in lock statement"`. Store the chosen `SymbolRef`
   in `L.RestrictedSymbol`.

3. **Isolation analysis on the body.** Call
   `isolationCheckNode(a, L.Body, {L.RestrictedSymbol})`.

   When `analyzeLocks`'s outer walk continues, it descends into `L.Body`
   so any nested `lock` triggers the check in (1).

A lock body with zero candidates is not an error: `L.RestrictedSymbol`
is left unset and step 3 calls `isolationCheckNode(a, L.Body, nil)`.

The nested-lock check requires `analyzeLock` to know whether the
visitor is already inside a lock body. Add an `inLock bool` (or
equivalent) field to `functionAnalyzer`, set to true while descending
into `L.Body` and restored on the way out.

---

## Phase 5 — Adjacent isolation rules

Three small semantic-analyzer additions that are in scope for this PR
but share no code with the lock analyzer. Each is independent.

### 5.1 Module-isolated-var access must be inside a `lock`

Add a `*ast.BLangSimpleVarRef` case to `functionAnalyzer.Visit`: if the
unnarrowed symbol is a module-level variable with `IsIsolated() == true`
and `fa.inLock` is false, emit `"access of an isolated variable must be
inside a lock statement"`. Reuses the `inLock` state already maintained
by the analyzer for Phase 4.

Test: `lock-module-var-outside-1-e.bal`.

### 5.2 Module-isolated-var initializer must be isolated

Where module-level variable declarations are validated, if the variable
`IsIsolated()` and has an initializer expression `e`, run
`isolationCheckNode(a, e, nil)` on `e`.

Test: `isolated-var-init-1-e.bal`.

### 5.3 Non-final field of isolated class must be `private`

Tighten `validateIsolatedClassFields` so a non-final field of an
isolated class must be declared `private`.

Test: `lock-public-mutable-field-1-e.bal`.

---

## Phase 6 — CFG

Files: `semantics/cfg_analyzer.go`, `semantics/control_flow_analyzer.go`.

1. Treat the lock body as a regular block for reachability and explicit-return
   purposes. The lock itself doesn't terminate; control flows through it.
2. Update `corpus/cfg` expected outputs.

No interesting changes here yet — the lock-release-on-abrupt-exit problem is
solved in desugar (Phase 7), not in CFG, because the CFG operates over the AST
shape.

---

## Phase 7 — Desugar

File: `desugar/statement.go`.

Walk `BLangLock.Body` recursively (as today for any block). Nothing else
is required from desugar; lock-release placement is handled in BIR-gen
(Phase 8).

Update `corpus/desugared`.

---

## Phase 8 — BIR

Files: `bir/terminator.go`, `bir/model.go`, `bir/pretty_print.go`,
`bir/codec/*`, `bir/bir_gen.go`.

1. Define two new terminator instructions:

   ```go
   type LockStart struct {
       BIRTerminatorBase
       LockKey string
   }

   type LockEnd struct {
       BIRTerminatorBase
       LockKey string
   }
   ```

   Both have `ThenBB` (entry of body / continuation) like every other
   terminator. The string key is stable across BIR (de)serialization —
   no per-compilation id allocation, no compiler-context state.

2. Pretty-print and serialize/deserialize them (codec). New
   `InstructionKind` constants.

3. **Lock-key construction.** Add a helper in `bir_gen.go`:

   ```go
   func buildLockKey(ctx *Context, restricted model.SymbolRef) string {
       sym := ctx.CompilerContext.GetSymbol(restricted)
       if owner := /* owning class, if any */; owner != nil {
           return buildLookupKey(restricted.Package,
               owner.Name() + "." + sym.Name())
       }
       return buildLookupKey(restricted.Package, sym.Name())
   }
   ```

   This reuses `buildLookupKey` (already used for global vars,
   functions, classes, methods) so the resulting strings live in the
   same namespace as the rest of the BIR's symbolic identifiers.

4. `bir_gen.go`'s per-function context gets an `activeLockKey *string`
   field (nil when not inside a lock). Phase 4 has already rejected
   nested locks, so at most one lock is active at any point.

5. `bir_gen.go` handling of `*ast.BLangLock` `L`:
   - `key := buildLockKey(ctx, L.RestrictedSymbol)`
   - Close the current BB with `LockStart{LockKey: key,
     ThenBB=bodyEntry}`.
   - Set `activeLockKey = &key` and emit Body BBs as usual.
   - Clear `activeLockKey` and close the body's tail BB with
     `LockEnd{LockKey: key, ThenBB=afterLock}`.

6. **Abrupt-exit handling.** When `bir_gen.go` is about to close a BB
   with a `Return` / `Break` / `Continue` / `Fail` terminator while
   `activeLockKey != nil`, it first closes the current BB with
   `LockEnd{LockKey: *activeLockKey, ThenBB=newBB}`, opens `newBB` as
   the current BB, and emits the original terminator there.

   Panics are not handled here; the runtime drains held locks on panic
   unwinding (Phase 9).

7. Update `corpus/bir`.

---

## Phase 9 — Runtime: Environment, strand id, lock table, re-entrant mutex

Files: new `runtime/environment.go`, modified `runtime/runtime.go`,
`runtime/internal/exec/context.go`, new `runtime/internal/exec/locks.go`,
`runtime/internal/exec/terminators.go`, `runtime/internal/exec/errors.go`,
`runtime/internal/exec/frame.go`.

### 7.1 `runtime.Environment`

New type that owns runtime-wide state:

```go
// runtime/environment.go
type Environment struct {
    registry      *modules.Registry
    nextStrandID  atomic.Uint64 // initialised to 1; wraps on overflow
    locks         *locks.Mutexes // see 7.3
}

func NewEnvironment() *Environment {
    env := &Environment{
        registry: modules.NewRegistry(),
        locks:    locks.NewMutexes(),
    }
    env.nextStrandID.Store(1)
    return env
}

func (e *Environment) Registry() *modules.Registry { return e.registry }
func (e *Environment) Locks() *locks.Mutexes      { return e.locks }
func (e *Environment) AllocateStrandID() uint64 {
    for {
        v := e.nextStrandID.Add(1)
        if v != 0 { // skip the unowned sentinel on wrap
            return v
        }
    }
}
```

`Runtime` is updated:

```go
type Runtime struct {
    env      *Environment
    platform pal.Platform
}

func NewRuntime(platform pal.Platform) *Runtime {
    rt := &Runtime{env: NewEnvironment(), platform: platform}
    for _, init := range moduleInitializers { init(rt) }
    return rt
}
```

All existing `rt.registry.X(...)` call-sites become `rt.env.Registry().X(...)`.
`Runtime.Interpret` passes `rt.env` to `exec.Interpret`.

Introduce a small `runtime/internal/locks` package containing
`ReentrantMutex` and `Mutexes` (id → mutex map). `Environment` holds a
`*locks.Mutexes`; `exec` imports `runtime/internal/locks`. This avoids
an import cycle between `runtime` and `runtime/internal/exec`.

### 7.2 Strand id on the context

`exec.Context`:

```go
type Context struct {
    env       *Environment // replaces registry field
    strandID  uint64
    callStack callStack
    typeCtx   semtypes.Context
}

func NewContext(env *Environment) *Context {
    return &Context{
        env:       env,
        strandID:  env.AllocateStrandID(),
        callStack: callStack{elements: make([]*Frame, 0, 32)},
        typeCtx:   semtypes.TypeCheckContext(env.Registry().GetTypeEnv()),
    }
}

func (ctx *Context) StrandID() uint64       { return ctx.strandID }
func (ctx *Context) Env() *Environment      { return ctx.env }
func (ctx *Context) Registry() *modules.Registry { return ctx.env.Registry() }
```

Existing accessors that touched `ctx.registry` are routed via
`ctx.env.Registry()`. Today only one `Context` is created per `Interpret`,
so strand allocation is a no-op semantically but is plumbed for future
worker support.

### 7.3 Reentrant mutex and lock table

`runtime/internal/locks/locks.go`:

```go
type ReentrantMutex struct {
    mu    sync.Mutex
    cond  *sync.Cond
    owner uint64 // 0 == unowned
    count int
}

func NewReentrantMutex() *ReentrantMutex {
    m := &ReentrantMutex{}
    m.cond = sync.NewCond(&m.mu)
    return m
}

func (r *ReentrantMutex) Lock(strandID uint64)   { /* per DESIGN sketch */ }
func (r *ReentrantMutex) Unlock(strandID uint64) { /* per DESIGN sketch */ }

type Mutexes struct {
    mu sync.Mutex
    m  map[string]*ReentrantMutex
}

func NewMutexes() *Mutexes { return &Mutexes{m: map[string]*ReentrantMutex{}} }

func (t *Mutexes) Get(key string) *ReentrantMutex {
    t.mu.Lock(); defer t.mu.Unlock()
    if m, ok := t.m[id]; ok { return m }
    m := NewReentrantMutex()
    t.m[id] = m
    return m
}
```

`owner == 0` means free; strand ids start at 1 (see 7.1), so 0 is a valid
sentinel. `Unlock` panics if the caller is not the owner — that is a
runtime invariant violation, not a Ballerina panic.

### 7.4 Held-lock list on the strand

Locks are held by the *strand*, not by a specific frame. Add a held-lock
stack to `Context`:

```go
type heldLock struct {
    key   string
    mutex *locks.ReentrantMutex
}

type Context struct {
    ...
    heldLocks []heldLock
}
```

- `LockStart` pushes `heldLock{id, env.Locks().Get(id)}` after acquiring.
- `LockEnd` pops the top entry, asserts the id matches, and releases.
- On *Ballerina panic* unwinding (existing mechanism in
  `runtime/internal/exec/errors.go`): before unwinding past the entry
  point of the top-most lock-holding frame, drain `heldLocks` in LIFO
  order, calling `Unlock(strandID)` for each. We do this once at the
  panic-handling site rather than per-frame, because locks are
  strand-scoped.

  Concretely: when a Ballerina panic propagates and is about to be
  observed at the top-level (or caught by an enclosing trap/check), call
  `ctx.releaseAllHeldLocks()` before resuming normal control flow.

### 7.5 Terminator handlers

In `runtime/internal/exec/terminators.go`:

```go
case *bir.LockStart:
    m := ctx.Env().Locks().Get(n.LockKey)
    m.Lock(ctx.StrandID())
    ctx.heldLocks = append(ctx.heldLocks, heldLock{n.LockKey, m})
    return n.ThenBB

case *bir.LockEnd:
    top := ctx.heldLocks[len(ctx.heldLocks)-1]
    if top.key != n.LockKey { /* internal-error panic */ }
    ctx.heldLocks = ctx.heldLocks[:len(ctx.heldLocks)-1]
    top.mutex.Unlock(ctx.StrandID())
    return n.ThenBB
```

### 7.6 Tests

Add an integration test (`lock-reentrant-1-v.bal`) that exercises a
recursive isolated function calling itself inside `lock { ... }` to
validate re-entrancy.

Add a panic test (`lock-panic-release-1-p.bal`) under `corpus/bal/...`
where a Ballerina panic inside a lock body propagates out and a subsequent
`lock` on the same id from the same strand still succeeds (i.e. the lock
was in fact released).

Update `corpus/integration` and any expected files.

---

## Phase 10 — Final wiring & docs

1. Run `go build ./...` and full corpus test suite.
2. Update any existing isolation-analysis fixtures that change behaviour due
   to the "non-final field of isolated class must be private" tightening.
3. Make sure `corpus/integration/*.txtar` references for any renamed/added
   files are consistent.
4. Add a `TODO.md` entry tracking that locks are coarse-grained (one mutex
   per restricted symbol, shared across all instances of an isolated
   class). Revisit when worker/`start` support lands and contention becomes
   measurable.

---

## Out of scope for this PR

- `start`/worker syntax. We add strand-id plumbing but no concurrent strand
  spawning; everything runs in strand 1.
- BIR-on-disk format changes are minimal (two new instructions). No format
  migration utility.
- Fairness in `ReentrantMutex` (`Signal` may starve; acceptable for now,
  documented as a follow-up).
- Per-instance / per-field-region fine-grained locking. The current design
  is the "naive" implementation from the spec quote: one global recursive
  mutex per restricted symbol.
