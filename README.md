# gounion

A Go linter that checks exhaustiveness of type switches on union interfaces (sealed interfaces).

## What is a Union Type in Go?

In Go, you can simulate union types (also known as sum types or discriminated unions) by defining an interface with an unexported marker method. This pattern restricts which types can implement the interface to those defined in the same package.

```go
// Shape is a union type - only Circle, Rectangle, and Triangle can implement it
type Shape interface {
    isShape() // unexported marker method
}

type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Triangle struct{ Base, Height float64 }

func (*Circle) isShape()    {}
func (*Rectangle) isShape() {}
func (*Triangle) isShape()  {}
```

This pattern is used in Go's standard library, such as `go/ast` package for AST nodes.

## Installation

```bash
go install github.com/YuitoSato/gounion@latest
```

## Usage

```bash
gounion ./...
```

## Example

### Defining a Union Type

```go
// shape/shape.go
package shape

// Shape is a union type representing geometric shapes.
// The isShape() marker method restricts implementations to this package.
type Shape interface {
    isShape()
}

type Circle struct {
    Radius float64
}

type Rectangle struct {
    Width  float64
    Height float64
}

type Triangle struct {
    Base   float64
    Height float64
}

func (*Circle) isShape()    {}
func (*Rectangle) isShape() {}
func (*Triangle) isShape()  {}
```

### Detected Issue

```go
// main.go
package main

import "shape"

// NG: Missing Triangle case
func CalculateArea(s shape.Shape) float64 {
    switch s := s.(type) {
    case *shape.Circle:
        return 3.14 * s.Radius * s.Radius
    case *shape.Rectangle:
        return s.Width * s.Height
    // Missing *shape.Triangle!
    }
    return 0
}
```

Running gounion:

```
$ gounion ./...
main.go:7:5: missing cases in type switch on Shape: shape.*Triangle
```

### Correct Implementation

```go
// OK: All cases covered
func CalculateArea(s shape.Shape) float64 {
    switch s := s.(type) {
    case *shape.Circle:
        return 3.14 * s.Radius * s.Radius
    case *shape.Rectangle:
        return s.Width * s.Height
    case *shape.Triangle:
        return 0.5 * s.Base * s.Height
    }
    return 0
}
```

### Default Case

When a `default` case is present, no warning is issued:

```go
// OK: Has default case, no warning
func CalculateArea(s shape.Shape) float64 {
    switch s := s.(type) {
    case *shape.Circle:
        return 3.14 * s.Radius * s.Radius
    default:
        return 0
    }
}
```

However, if the `default` case ends with a `panic()` call or returns an error, the exhaustiveness check is still enforced. This is because these patterns are typically used as safety guards rather than intentional handling of unknown types:

```go
// NG: default ends with panic, missing Rectangle and Triangle
func CalculateArea(s shape.Shape) float64 {
    switch s := s.(type) {
    case *shape.Circle:
        return 3.14 * s.Radius * s.Radius
    default:
        panic("unreachable")
    }
}
```

```go
// NG: default returns error, missing Rectangle and Triangle
func CalculateArea(s shape.Shape) (float64, error) {
    switch s := s.(type) {
    case *shape.Circle:
        return 3.14 * s.Radius * s.Radius, nil
    default:
        return 0, fmt.Errorf("unexpected shape: %T", s)
    }
}
```

This also applies when additional statements (e.g., logging) precede the `panic()` or error return:

```go
// NG: default ends with panic after logging, missing Rectangle and Triangle
func CalculateArea(s shape.Shape) float64 {
    switch s := s.(type) {
    case *shape.Circle:
        return 3.14 * s.Radius * s.Radius
    default:
        log.Println("unexpected type")
        panic("unreachable")
    }
}
```

The error return detection covers all types implementing the `error` interface, including `fmt.Errorf()`, `errors.New()`, sentinel errors, and custom error types. A `default` case that returns `nil` for the error value is treated as a normal default (no exhaustiveness check).

## How It Works

1. **Detects Union Interfaces**: Finds interfaces with unexported marker methods (methods that take no parameters and return nothing)
2. **Identifies Members**: Collects all types in the package that implement the marker method
3. **Checks Exhaustiveness**: When a type switch is used on a union interface, verifies that all member types are handled
4. **Respects Default**: Skips the check if a `default` case is present, unless the default ends with a `panic()` call or returns an error

## Integration with golangci-lint

Add to your `.golangci.yml`:

```yaml
linters-settings:
  custom:
    gounion:
      path: github.com/YuitoSato/gounion
      description: checks exhaustiveness of type switches on union interfaces
```

## License

MIT
