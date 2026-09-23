package gounion

import (
	"go/ast"
	"go/types"
	"sort"

	"golang.org/x/tools/go/analysis"
)

// exportUnionFacts scans the current package for union interfaces
// and exports facts for them along with their implementing types.
func exportUnionFacts(pass *analysis.Pass) {
	for typeName, markerMethod := range findUnionInterfaces(pass) {
		members := findUnionMembers(pass, markerMethod)

		fact := &UnionInterface{
			MarkerMethod: markerMethod,
			Members:      members,
		}
		pass.ExportObjectFact(typeName, fact)
	}
}

// findUnionInterfaces returns the union interfaces declared in the package,
// keyed by interface object, with their marker method name.
func findUnionInterfaces(pass *analysis.Pass) map[*types.TypeName]string {
	unionInterfaces := make(map[*types.TypeName]string)

	// Full walk: type declarations can also appear inside function bodies.
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			genDecl, ok := n.(*ast.GenDecl)
			if !ok {
				return true
			}
			for _, spec := range genDecl.Specs {
				if typeName, markerMethod := unionInterfaceOf(pass, spec); typeName != nil {
					unionInterfaces[typeName] = markerMethod
				}
			}
			return true
		})
	}

	return unionInterfaces
}

// unionInterfaceOf returns the type name and marker method of the union
// interface declared by spec, or nil if spec declares something else.
func unionInterfaceOf(pass *analysis.Pass, spec ast.Spec) (*types.TypeName, string) {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, ""
	}

	// Check if it's an interface type
	if _, ok := typeSpec.Type.(*ast.InterfaceType); !ok {
		return nil, ""
	}

	// Get the type object
	typeName, ok := pass.TypesInfo.Defs[typeSpec.Name].(*types.TypeName)
	if !ok {
		return nil, ""
	}

	// Get the interface type
	iface, ok := typeName.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, ""
	}

	// Check for marker methods
	markerMethod := findMarkerMethod(iface)
	if markerMethod == "" {
		return nil, ""
	}

	return typeName, markerMethod
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
