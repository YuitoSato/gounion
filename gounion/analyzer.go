package gounion

import "golang.org/x/tools/go/analysis"

// Analyzer is the gounion analyzer that checks exhaustiveness of type switches
// on union interfaces (interfaces with private marker methods).
var Analyzer = &analysis.Analyzer{
	Name: "gounion",
	Doc:  "checks exhaustiveness of type switches on union interfaces",
	Run:  run,
	// No Requires: a fact-using analyzer and its Requires run on every
	// transitive dependency, and the x/tools checker (singlechecker,
	// checker.Analyze) retains their results until the whole run ends, so
	// requiring inspect.Analyzer would keep an inspector per dependency
	// resident. The syntax is walked with ast.Inspect instead.
	FactTypes: []analysis.Fact{new(UnionInterface)},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Phase 1: Detect union interfaces and export facts
	exportUnionFacts(pass)

	// Phase 2: Check type switch exhaustiveness
	checkTypeSwitches(pass)

	return nil, nil
}
