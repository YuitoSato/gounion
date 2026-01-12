package consumer

import "union"

// ===========================================
// Test Cases: Using Result from external package
// ===========================================

// ProcessResult - NG: Missing Error case
func ProcessResult(r union.Result) string {
	switch r.(type) { // want `missing cases in type switch on Result: union\.\*Error`
	case *union.Success:
		return "processed successfully"
	}
	return ""
}

// ProcessResultComplete - OK: All cases covered
func ProcessResultComplete(r union.Result) string {
	switch r.(type) {
	case *union.Success:
		return "processed successfully"
	case *union.Error:
		return "processing failed"
	}
	return ""
}

// ===========================================
// Test Cases: Using Shape from external package
// ===========================================

// DrawShape - NG: Missing Rectangle and Triangle cases
func DrawShape(s union.Shape) string {
	switch s.(type) { // want `missing cases in type switch on Shape: union\.\*Rectangle, union\.\*Triangle`
	case *union.Circle:
		return "drawing circle"
	}
	return ""
}

// DrawShapeComplete - OK: All cases covered
func DrawShapeComplete(s union.Shape) string {
	switch s.(type) {
	case *union.Circle:
		return "drawing circle"
	case *union.Rectangle:
		return "drawing rectangle"
	case *union.Triangle:
		return "drawing triangle"
	}
	return ""
}

// DrawShapeWithDefault - OK: Has default case
func DrawShapeWithDefault(s union.Shape) string {
	switch s.(type) {
	case *union.Circle:
		return "drawing circle"
	default:
		return "drawing unknown shape"
	}
}

// GetShapeName - NG: Missing Triangle case
func GetShapeName(s union.Shape) string {
	switch s.(type) { // want `missing cases in type switch on Shape: union\.\*Triangle`
	case *union.Circle:
		return "Circle"
	case *union.Rectangle:
		return "Rectangle"
	}
	return ""
}
