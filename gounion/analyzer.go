package gounion

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the gounion analyzer that checks exhaustiveness of type switches
// on union interfaces (interfaces with private marker methods).
var Analyzer = &analysis.Analyzer{
	Name:      "gounion",
	Doc:       "checks exhaustiveness of type switches on union interfaces",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{new(UnionInterface)},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Phase 1: Detect union interfaces and export facts
	exportUnionFacts(pass, inspect)

	// Phase 2: Check type switch exhaustiveness
	checkTypeSwitches(pass, inspect)

	return nil, nil
}
