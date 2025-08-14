# Task

Convert the selected Java source into idiomatic Go **with minimal semantic drift**.

- Map packages to directories and package names.
- Map classes (including abstract) and their fields/methods to Go **interface** and **unexported implementation structs**.
- Preserve logic while adapting to Go naming, encapsulation and error handling.
- Keep diffs focused: one Java file -> one Go file (small support code if strictly necessary).

## Inputs

- Java file(s) to convert.
- Any related types the file depends on.
- Package path target within the Go module.

## Outputs

- A new Go file with idiomatic code.

## Package & File Mapping

1. **One-to-one element mapping**
    - **Files**: Map each Java source file `<Name>.java` to to one go file, named `lowercase-with-dashes.go` (e.g., `FooBar.java` -> `foo-bar.go`). Be consistent within the package.
    - **Classes/Interfaces**: For each Java class or interface define a **Go interface** named after the Java type (exported) and an **unexported implementation struct**.
    - **Methods**: Map each java method to a Go method or function. For overloads, pick distinct names (e.g., `Advance`, `AdvanceN`).

2. **Package structure**
    - Java: `package com.example.foo.something;`
    - Go: directory `foo/something` with `package something` in the file.

## Types & Encapsulation

### Classes -> Interface + Impl Struct

- Define an **interface** named exactly after the Java class, e.g. `Foo`.
- Define a constructor `NewFoo(...) Foo` returning the interface.
- Define an **unexported** struct `fooImpl` implements the interface. Prefer composition/embedding for reuse.

### Fields -> Struct Fields

- Private/encapsulated state lives in unexported struct fields (e.g., `bar int`).
- Provide getters/setters as interface methods only if the Java API requires them to be public. Avoid exporting fields directly unless they are immutable configuration.

### Constructors

For each Java public constructor:

```go
func NewFoo(params) Foo {
    return &fooImpl{/* initialize fields */}
}
```

- If multiple Java constructors exist, use either distinct names (`NewFooWithParams`) **or** the functional options pattern for optional params (avoid overloading).

### Methods & Receivers

- **Non-mutating** -> value receiver if the struct is small and the method is read-only.
- **Mutating** -> pointer receiver.
- Prefer returning `(T, error)` over panicking; translate Java exceptions to `error` values when they cross API boundaries.

### Abstract classes

```go
type Foo interface {
    // abstract methods + requires accessors
}

type fooBase struct {/* shared fields */}

func (b fooBase) GetX() T { return b.x }

type fooImpl struct {
    fooBase
}
func NewFoo(params) Foo {
    return &fooImpl{
        fooBase: fooBase{/* initialize shared fields */},
    }
}
```

### Method overloading

Java:

```java
void advance();
void advance(int n);
```

Go:

```go
func (r *readerImpl) Advance() {}
func (r *readerImpl) AdvanceN(n int) {}
```

### Equality / Hashing

- If Java only overrides `equals()`/`hashCode()`, in Go prefer:
  - Use direct `==` for comparable structs; or
  - Provide an explicit **lookup key**:

    ```go
    type FooLookupKey struct { Bar int; Baz string }
    func (f *fooImpl) FooLookupKey() FooLookupKey { return FooLookupKey{f.bar, f.baz} }
    ```

## Naming & Comments

- **Packages:** short, lowercase, single word (`something`).  
- **Exports:** capitalize to export. Keep names concise and avoid stutter (prefer `something.Reader` with type name `Reader`).

## Error Handling

- Always return `(T, error)` for fallible operations. Don't use `panic` for normal control flow.
- Wrap lower-level errors with context using `fmt.Errorf("op: %w", err)` so callers can use `errors.Is/As`.
- Match errors with `errors.Is` (sentinels) or `errors.As` (typed errors with additional context).

### Java exceptions -> Go typed errors

- Define Java-specific exceptions as typed errros in `common/errors/errors.go` (package `errors`).

```go
package errors

import "fmt"

type IndexOutOfBoundsError struct {
	index  int
	length int
}

func (e IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("Index %d out of bounds for length %d", e.index, e.length)
}

func (e IndexOutOfBoundsError) GetIndex() int  { return e.index }
func (e IndexOutOfBoundsError) GetLength() int { return e.length }
```

## Generics

- Map Java generics to Go generics when needed.

## Guardrails & Non goals

- **Do no** add file headers or license comments.
- **Do not** introduce new public APIs unless requires by the Java surface.
- **Do not** add comments unless the Java source has them.
