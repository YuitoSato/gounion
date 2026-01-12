package gounion

import (
	"go/ast"
	"go/types"
	"sort"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// exportUnionFacts scans the current package for union interfaces
// and exports facts for them along with their implementing types.
func exportUnionFacts(pass *analysis.Pass, inspect *inspector.Inspector) {
	// Find all interface type declarations
	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	// Map to store union interfaces: interface object -> marker method name
	unionInterfaces := make(map[*types.TypeName]string)

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Check if it's an interface type
			_, ok = typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			// Get the type object
			obj := pass.TypesInfo.Defs[typeSpec.Name]
			if obj == nil {
				continue
			}

			typeName, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}

			// Get the interface type
			iface, ok := typeName.Type().Underlying().(*types.Interface)
			if !ok {
				continue
			}

			// Check for marker methods
			markerMethod := findMarkerMethod(iface)
			if markerMethod == "" {
				continue
			}

			unionInterfaces[typeName] = markerMethod
		}
	})

	// For each union interface, find its members and export the fact
	for typeName, markerMethod := range unionInterfaces {
		members := findUnionMembers(pass, markerMethod)

		fact := &UnionInterface{
			MarkerMethod: markerMethod,
			Members:      members,
		}
		pass.ExportObjectFact(typeName, fact)
	}
}

// findMarkerMethod checks if an interface has a marker method.
// A marker method is:
// - unexported (starts with lowercase)
// - has no parameters
// - has no return values
func findMarkerMethod(iface *types.Interface) string {
	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)

		// Check if unexported
		if method.Exported() {
			continue
		}

		sig, ok := method.Type().(*types.Signature)
		if !ok {
			continue
		}

		// Check no parameters
		if sig.Params().Len() != 0 {
			continue
		}

		// Check no return values
		if sig.Results().Len() != 0 {
			continue
		}

		return method.Name()
	}
	return ""
}

// findUnionMembers finds all types in the package that implement
// the given marker method.
func findUnionMembers(pass *analysis.Pass, markerMethod string) []string {
	var members []string

	scope := pass.Pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)

		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		// Skip interfaces themselves
		if _, ok := typeName.Type().Underlying().(*types.Interface); ok {
			continue
		}

		// Check both value type and pointer type for the marker method
		if hasMarkerMethod(pass.Pkg, typeName.Type(), markerMethod) {
			members = append(members, typeName.Name())
		} else if hasMarkerMethod(pass.Pkg, types.NewPointer(typeName.Type()), markerMethod) {
			members = append(members, "*"+typeName.Name())
		}
	}

	// Sort members for consistent output
	sort.Strings(members)

	return members
}

// hasMarkerMethod checks if a type has the given marker method.
func hasMarkerMethod(pkg *types.Package, typ types.Type, markerMethod string) bool {
	mset := types.NewMethodSet(typ)
	sel := mset.Lookup(pkg, markerMethod)
	return sel != nil
}
