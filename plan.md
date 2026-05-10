# Fixed-length array fixes — plan

## Problem summary

All `corpus/bal/subset8/08-list/fixedlength*-e.bal` tests currently panic at
AST-build time:

```
panic: array length expression handling unimplemented
ast/node_builder.go:3384  TransformArrayTypeDescriptor
```

The size expression in `int[N]` (`ArrayDimensionNode.ArrayLength()`) is never
transformed to a `BLangExpression`. As a result:

* **fixedlength1-e / fixedlength2-e** — `v[4]` on `int[4]` never gets to type
  resolution; even if it did, OOB indexing is supposed to give `never`, which
  must surface as an error in the surrounding context.
* **fixedlength3-e / fixedlength4-e** — assignment / compound-assignment to a
  fixed-length array element never reaches the LHS/RHS type compatibility check.
* **fixedlength5-e** — invocation argument `foo(v[2])` (foo expects `string`,
  `v[2]: int`) never reaches the parameter type check.
* **fixedlength6-e** — `int[3] v = [1, 2, "3"]` never reaches per-member type
  check inside the list constructor.
* **fixedlength7-e** — `int[x]` where `const x = true` should reject a non-int
  size expression.
* **fixedlength8-e** — `int[x]` where `x` is a non-constant (function-result)
  variable should reject a non-singleton size expression.
* **fixedlength9-e** — `int[3] v = [1, 2, 3, 4]` should reject too-many members.

The good news: most of the validation is already wired up correctly elsewhere
in the pipeline; the dominant blocker is the AST-stage panic plus a too-narrow
size-expression handler in the type resolver.

## Inventory of existing behavior (already correct)

* `semtypes.ListMemberTypeInnerVal` returns `NEVER` for an out-of-range index
  on a fixed-length list type.
* `ListAtomicType.MemberAtInnerVal(i)` returns `NEVER` for `i >= FixedLength`
  on a closed list type.
* `resolveListConstructorWithExpectedType` already calls `MemberAtInnerVal(i)`
  per member; if it is `NEVER` it emits `"too many members in list
  constructor"`. It also re-resolves each member with the per-index expected
  type.
* `analyzeListConstructorExpr` re-validates each member against
  `lat.MemberAtInnerVal(i)`.
* `validateResolvedType` already reports an error when an expression resolves
  to `never` and the surrounding context expects something other than `never`.
* `analyzeAssignment` already calls `validateResolvedType(rhs, lhsTy)`, so a
  type-mismatched RHS is caught.
* `bir/bir_gen.go listConstructorExpression` already fills missing fixed-length
  members with default values, so once stages 1–10 work the BIR stage is fine.

So once the AST and type-resolver issues below are fixed, the existing checks
will fire automatically for tests 1–6 and 9. Tests 7–8 need a new check.

## Fixes

### 1. AST: transform the array length expression
File: `ast/node_builder.go`, `TransformArrayTypeDescriptor`.

Replace the `panic("array length expression handling unimplemented")` branch
with:

* If `dimensionNode.ArrayLength() == nil` → keep as `nil` (open array).
* Otherwise → `n.createExpression(dimensionNode.ArrayLength())` and append the
  resulting `BLangExpression` to `sizes`.

This lets int literals, constant references, and any other expression survive
into the `BLangArrayType.Sizes` slice.

### 2. Type resolver: validate and consume the size expression
File: `semantics/type_resolver.go`, `case *ast.BLangArrayType` in
`resolveBType`.

Currently:

```go
length := int(lenExp.(*ast.BLangLiteral).Value.(int64))
```

This panics for anything other than an int literal. Replace with:

* `resolveActionOrExpression(t, nil, lenExp.(ast.BLangActionOrExpression), semtypes.INT)`
  to resolve and type-check the size expression in an `int`-expected context.
  (The chain is `nil` because the size expression appears in a type-descriptor
  position, outside any narrowing scope.)
* Read `lenExp.GetDeterminedType()`:
  * Single check: must be a subtype of `int` AND `semtypes.SingleShape(...)`
    must return a value. If either fails →
    `t.semanticError("fixed-length array size must be a singleton int",
    lenExp.GetPosition())` and return `nil, false`. Covers both
    fixedlength7-e (`const x = true`, boolean singleton) and fixedlength8-e
    (local `int x = foo();`, non-singleton int).
  * Use the singleton int value, validate `>= 0`, and pass as `length` to
    `DefineListTypeWrappedWithEnvSemTypesInt`.

This also keeps the existing happy path for `BLangLiteral` int values working
(they are singleton ints), and now additionally accepts constant references
(symbol resolution will have replaced them with their singleton type) and
rejects everything else.

Note: this resolution must happen during top-level type resolution (stage 4)
because array type descriptors can appear in module-level decls. For local
type descriptors it runs as part of the inner type-resolution pass.

Open question (Q3 in chat): if reaching `resolveActionOrExpression` from
`resolveBType` is awkward (chain plumbing, recursion concerns), an alternative
is a small dedicated `resolveArraySizeExpr` that handles only:

* `BLangLiteral` int → its `IntConst` type
* `BLangSimpleVarRef` → look up symbol; if it's a constant with an int
  singleton type, use it; otherwise error.

The first option (full expression resolution) is simpler and more uniform; the
second is more restricted and matches the spec, which only allows int literal
or constant ref. We prefer option 1 unless there is a layering concern.

### 3. Index-based access: rely on `never` projection (no code change)
File: `semantics/type_resolver.go`, `resolveIndexBasedAccess`.

Already calls `ListMemberTypeInnerVal(... containerExprTy, keyExprTy)`. For an
out-of-range constant index this returns `NEVER`. In the surrounding context
`validateResolvedType` will see `IsNever(resolvedTy) && !IsNever(expectedType)`
and emit an "incompatible types" diagnostic, satisfying fixedlength1-e and
fixedlength2-e.

If a dedicated "index out of bounds" message is preferred, we can add an
optional check here when `keyExprTy` is a singleton int constant and the
projection is `NEVER` while `containerExprTy` is non-empty list.

### 4. Assignment / compound assignment (no code change)
`analyzeAssignment` → `validateResolvedType(rhs, lhsTy)` already enforces
RHS subtype LHS, so fixedlength3-e (`v[2] = 1.5`) and fixedlength4-e
(`v[2] -= 1.5`) become errors automatically.

### 5. Function call argument (no code change)
Function-call argument checking already calls `IsSubtype(argTy, paramTy)`, so
fixedlength5-e (`foo(v[2])`, `foo(string)`) becomes an error automatically.

### 6. List constructor (no code change)
Both `resolveListConstructorWithExpectedType` and `analyzeListConstructorExpr`
already use `lat.MemberAtInnerVal(i)`, which returns `NEVER` past the fixed
length; this gives:

* fixedlength6-e: `int[3] v = [1, 2, "3"]` — element 2 expected `int`, gets
  `"3"` → mismatch.
* fixedlength9-e: `int[3] v = [1, 2, 3, 4]` — 4th member's expected type is
  `NEVER` → "too many members" error.

## Test housekeeping

Adding the AST-level fix unblocks the valid tests
`fixedlength1-v.bal` and `fixedlength2-v.bal`, which are currently skipped
in:

* `ast/corpus_ast_test.go`
* `semantics/corpus_symbol_resolver_test.go`
* `semantics/corpus_type_resolver_test.go`
* `semantics/corpus_semantic_analysis_test.go`
* `semantics/corpus_cfg_test.go`
* `desugar/corpus_desugar_test.go`

Also `corpus/integration_test.go:144-161` lists all nine `-e.bal` files in the
"currently failing" block.

Pending Q4 in chat. Default plan:

1. Run each affected corpus test with `-update` to regenerate expected outputs
   under `corpus/ast`, `corpus/symbolresolver`, `corpus/typeresolver`,
   `corpus/semantic`, `corpus/cfg`, `corpus/desugar` for the two `-v` files.
2. Remove the two `-v` entries from each skip list.
3. Remove the nine `-e` entries from `corpus/integration_test.go` failing
   block.
4. Run `go test ./...` and the targeted corpus tests; fix any fallout.

## Verification

For each `-e` test, after the fixes:

| Test            | Expected diagnostic                                                         |
| --------------- | --------------------------------------------------------------------------- |
| fixedlength1-e  | incompatible type (got `never`) at `v[4]`                                   |
| fixedlength2-e  | incompatible type (got `never`) at `v[4] = 4`                               |
| fixedlength3-e  | incompatible type (got `float`, expected `int`) at `1.5`                    |
| fixedlength4-e  | same as 3 (compound assignment routes through analyzeAssignment)            |
| fixedlength5-e  | incompatible argument type (got `int`, expected `string`)                   |
| fixedlength6-e  | incompatible type (got `"3"`, expected `int`) at member 2                   |
| fixedlength7-e  | fixed-length array size must be of type int (at `x`)                        |
| fixedlength8-e  | fixed-length array size must be a compile-time constant (at `x`)            |
| fixedlength9-e  | too many members in list constructor                                        |

For each `-v` test, the existing `@output` comments must match interpreter
output once stage-11 execution is reached.

## Out of scope

* Open arrays (`int[]`) and inferred-length arrays (`int[*]`) are unchanged.
* Multi-dimensional fixed-length arrays beyond what the resolver loop already
  supports.
* Runtime bounds checking for non-constant indices (that is enforced by the
  interpreter's array-access opcode and is unrelated to these tests).
