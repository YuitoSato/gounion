package gounion

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// checkTypeSwitches checks for exhaustiveness in type switch statements
// on union interfaces.
func checkTypeSwitches(pass *analysis.Pass, inspect *inspector.Inspector) {
	nodeFilter := []ast.Node{
		(*ast.TypeSwitchStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switchStmt := n.(*ast.TypeSwitchStmt)

		// Get the switch expression type
		switchType := getSwitchType(pass, switchStmt)
		if switchType == nil {
			return
		}

		// Check if it's a union interface
		namedType := extractNamedInterface(switchType)
		if namedType == nil {
			return
		}

		var unionFact UnionInterface
		if !pass.ImportObjectFact(namedType.Obj(), &unionFact) {
			return // Not a union interface
		}

		// Check for default case - if present, skip exhaustiveness check
		if hasDefaultCase(switchStmt) {
			return
		}

		// Get handled types from case clauses
		handledTypes := collectCaseTypes(pass, switchStmt)

		// Find missing types
		missing := findMissingTypes(unionFact.Members, handledTypes, namedType.Obj().Pkg())

		if len(missing) > 0 {
			pass.Reportf(switchStmt.Pos(),
				"missing cases in type switch on %s: %s",
				namedType.Obj().Name(),
				strings.Join(missing, ", "))
		}
	})
}

// getSwitchType extracts the type being switched on from a type switch statement.
func getSwitchType(pass *analysis.Pass, stmt *ast.TypeSwitchStmt) types.Type {
	// TypeSwitchStmt.Assign is either:
	// - *ast.ExprStmt containing x.(type)
	// - *ast.AssignStmt like v := x.(type)

	var typeAssert *ast.TypeAssertExpr

	switch assign := stmt.Assign.(type) {
	case *ast.ExprStmt:
		ta, ok := assign.X.(*ast.TypeAssertExpr)
		if !ok {
			return nil
		}
		typeAssert = ta
	case *ast.AssignStmt:
		if len(assign.Rhs) == 0 {
			return nil
		}
		ta, ok := assign.Rhs[0].(*ast.TypeAssertExpr)
		if !ok {
			return nil
		}
		typeAssert = ta
	default:
		return nil
	}

	if typeAssert == nil {
		return nil
	}

	tv, ok := pass.TypesInfo.Types[typeAssert.X]
	if !ok {
		return nil
	}

	return tv.Type
}

// extractNamedInterface extracts the named interface type from a type.
func extractNamedInterface(typ types.Type) *types.Named {
	named, ok := typ.(*types.Named)
	if !ok {
		return nil
	}

	// Verify it's an interface
	if _, ok := named.Underlying().(*types.Interface); !ok {
		return nil
	}

	return named
}

// hasDefaultCase checks if the type switch has a default case.
func hasDefaultCase(stmt *ast.TypeSwitchStmt) bool {
	for _, clause := range stmt.Body.List {
		caseClause, ok := clause.(*ast.CaseClause)
		if !ok {
			continue
		}
		// nil List means default case
		if caseClause.List == nil {
			return true
		}
	}
	return false
}

// collectCaseTypes collects all types mentioned in case clauses.
func collectCaseTypes(pass *analysis.Pass, stmt *ast.TypeSwitchStmt) []string {
	var handled []string

	for _, clause := range stmt.Body.List {
		caseClause, ok := clause.(*ast.CaseClause)
		if !ok {
			continue
		}

		for _, expr := range caseClause.List {
			tv, ok := pass.TypesInfo.Types[expr]
			if !ok {
				continue
			}

			// Format the type as a string for comparison
			typeStr := formatTypeForComparison(tv.Type)
			handled = append(handled, typeStr)
		}
	}

	return handled
}

// formatTypeForComparison formats a type for comparison with union members.
func formatTypeForComparison(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Pointer:
		if named, ok := t.Elem().(*types.Named); ok {
			return "*" + named.Obj().Name()
		}
	case *types.Named:
		return t.Obj().Name()
	}
	return types.TypeString(typ, nil)
}

// findMissingTypes finds union members that are not in the handled list.
func findMissingTypes(members []string, handled []string, unionPkg *types.Package) []string {
	handledSet := make(map[string]bool)
	for _, h := range handled {
		handledSet[h] = true
	}

	var missing []string
	for _, member := range members {
		if !handledSet[member] {
			// Format with package name for external references
			if unionPkg != nil {
				missing = append(missing, unionPkg.Name()+"."+member)
			} else {
				missing = append(missing, member)
			}
		}
	}

	return missing
}
