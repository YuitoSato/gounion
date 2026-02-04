package union

import (
	"errors"
	"fmt"
)

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

// Sentinel error for testing
var ErrUnexpectedType = errors.New("unexpected type")

// Custom error type for testing
type UnexpectedTypeError struct {
	Type string
}

func (e *UnexpectedTypeError) Error() string {
	return "unexpected type: " + e.Type
}

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

// CalculateAreaWithDefaultPanicAndLog - NG: default ends with panic, missing Rectangle and Triangle
func CalculateAreaWithDefaultPanicAndLog(s Shape) float64 {
	switch s := s.(type) { // want `missing cases in type switch on Shape: union\.\*Rectangle, union\.\*Triangle`
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

// ===========================================
// Test Cases: Default returns error
// ===========================================

// HandleResultWithDefaultFmtErrorf - NG: default returns fmt.Errorf, missing Error case
func HandleResultWithDefaultFmtErrorf(r Result) error {
	switch r.(type) { // want `missing cases in type switch on Result: union\.\*Error`
	case *Success:
		return nil
	default:
		return fmt.Errorf("unexpected type: %T", r)
	}
}

// CalculateAreaWithDefaultErrorsNew - NG: default returns errors.New with multiple return values, missing Rectangle and Triangle
func CalculateAreaWithDefaultErrorsNew(s Shape) (float64, error) {
	switch s := s.(type) { // want `missing cases in type switch on Shape: union\.\*Rectangle, union\.\*Triangle`
	case *Circle:
		return 3.14 * s.Radius * s.Radius, nil
	default:
		return 0, errors.New("unexpected shape")
	}
}

// HandleResultWithDefaultSentinelError - NG: default returns sentinel error, missing Error case
func HandleResultWithDefaultSentinelError(r Result) error {
	switch r.(type) { // want `missing cases in type switch on Result: union\.\*Error`
	case *Success:
		return nil
	default:
		return ErrUnexpectedType
	}
}

// HandleResultWithDefaultCustomError - NG: default returns custom error type, missing Error case
func HandleResultWithDefaultCustomError(r Result) error {
	switch r.(type) { // want `missing cases in type switch on Result: union\.\*Error`
	case *Success:
		return nil
	default:
		return &UnexpectedTypeError{Type: "unknown"}
	}
}

// CalculateAreaWithDefaultErrorAndLog - NG: default ends with return error, missing Rectangle and Triangle
func CalculateAreaWithDefaultErrorAndLog(s Shape) (float64, error) {
	switch s := s.(type) { // want `missing cases in type switch on Shape: union\.\*Rectangle, union\.\*Triangle`
	case *Circle:
		return 3.14 * s.Radius * s.Radius, nil
	default:
		fmt.Println("unexpected type")
		return 0, errors.New("unexpected shape")
	}
}

// CalculateAreaWithDefaultErrorComplete - OK: All cases covered, default returns error
func CalculateAreaWithDefaultErrorComplete(s Shape) (float64, error) {
	switch s := s.(type) {
	case *Circle:
		return 3.14 * s.Radius * s.Radius, nil
	case *Rectangle:
		return s.Width * s.Height, nil
	case *Triangle:
		return 0.5 * s.Base * s.Height, nil
	default:
		return 0, errors.New("unexpected shape")
	}
}

// HandleResultWithDefaultReturnNil - OK: default returns nil (not an error), treated as normal default
func HandleResultWithDefaultReturnNil(r Result) error {
	switch r.(type) {
	case *Success:
		return nil
	default:
		return nil
	}
}
