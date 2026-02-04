package union

import "fmt"

// ===========================================
// Example 1: Result - Represents operation outcome
// ===========================================

// Result is a union type representing operation outcome.
// The isResult() marker method restricts implementations to this package.
type Result interface { // want Result:`&\{isResult \[\*Error \*Success\]\}`
	isResult()
}

// Success represents a successful result.
type Success struct {
	Value string
}

// Error represents a failure result.
type Error struct {
	Message string
	Code    int
}

func (*Success) isResult() {}
func (*Error) isResult()   {}

// ===========================================
// Example 2: Shape - Represents geometric shapes
// ===========================================

// Shape is a union type representing geometric shapes.
type Shape interface { // want Shape:`&\{isShape \[\*Circle \*Rectangle \*Triangle\]\}`
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

// ===========================================
// Test Cases: Result type
// ===========================================

// HandleResult - NG: Missing Error case
func HandleResult(r Result) string {
	switch r.(type) { // want `missing cases in type switch on Result: union\.\*Error`
	case *Success:
		return "success"
	}
	return ""
}

// HandleResultComplete - OK: All cases covered
func HandleResultComplete(r Result) string {
	switch r.(type) {
	case *Success:
		return "success"
	case *Error:
		return "error"
	}
	return ""
}

// HandleResultWithDefault - OK: Has default case
func HandleResultWithDefault(r Result) string {
	switch r.(type) {
	case *Success:
		return "success"
	default:
		return "unknown"
	}
}

// ===========================================
// Test Cases: Shape type
// ===========================================

// CalculateArea - NG: Missing Triangle case
func CalculateArea(s Shape) float64 {
	switch s := s.(type) { // want `missing cases in type switch on Shape: union\.\*Triangle`
	case *Circle:
		return 3.14 * s.Radius * s.Radius
	case *Rectangle:
		return s.Width * s.Height
	}
	return 0
}

// CalculateAreaComplete - OK: All cases covered
func CalculateAreaComplete(s Shape) float64 {
	switch s := s.(type) {
	case *Circle:
		return 3.14 * s.Radius * s.Radius
	case *Rectangle:
		return s.Width * s.Height
	case *Triangle:
		return 0.5 * s.Base * s.Height
	}
	return 0
}

// CalculateAreaWithDefault - OK: Has default case
func CalculateAreaWithDefault(s Shape) float64 {
	switch s := s.(type) {
	case *Circle:
		return 3.14 * s.Radius * s.Radius
	default:
		return 0
	}
}

// HandleResultWithDefaultPanic - NG: default only panics, missing Error case
func HandleResultWithDefaultPanic(r Result) string {
	switch r.(type) { // want `missing cases in type switch on Result: union\.\*Error`
	case *Success:
		return "success"
	default:
		panic("unreachable")
	}
}

// CalculateAreaWithDefaultPanic - NG: default only panics, missing Rectangle and Triangle
func CalculateAreaWithDefaultPanic(s Shape) float64 {
	switch s := s.(type) { // want `missing cases in type switch on Shape: union\.\*Rectangle, union\.\*Triangle`
	case *Circle:
		return 3.14 * s.Radius * s.Radius
	default:
		panic("unreachable")
	}
}

// CalculateAreaWithDefaultPanicAndLog - OK: default has multiple statements (panic + log), treated as normal default
func CalculateAreaWithDefaultPanicAndLog(s Shape) float64 {
	switch s := s.(type) {
	case *Circle:
		return 3.14 * s.Radius * s.Radius
	default:
		fmt.Println("unexpected type")
		panic("unreachable")
	}
}

// CalculateAreaWithDefaultPanicSprintf - NG: default only panics with fmt.Sprintf, missing Rectangle and Triangle
func CalculateAreaWithDefaultPanicSprintf(s Shape) float64 {
	switch s := s.(type) { // want `missing cases in type switch on Shape: union\.\*Rectangle, union\.\*Triangle`
	case *Circle:
		return 3.14 * s.Radius * s.Radius
	default:
		panic(fmt.Sprintf("unexpected type: %T", s))
	}
}

// HandleShapeWithDefaultPanicOnly - NG: default only panics, no cases at all
func HandleShapeWithDefaultPanicOnly(s Shape) {
	switch s.(type) { // want `missing cases in type switch on Shape: union\.\*Circle, union\.\*Rectangle, union\.\*Triangle`
	default:
		panic("unreachable")
	}
}

// CalculateAreaWithDefaultPanicComplete - OK: All cases covered, default panics
func CalculateAreaWithDefaultPanicComplete(s Shape) float64 {
	switch s := s.(type) {
	case *Circle:
		return 3.14 * s.Radius * s.Radius
	case *Rectangle:
		return s.Width * s.Height
	case *Triangle:
		return 0.5 * s.Base * s.Height
	default:
		panic("unreachable")
	}
}
