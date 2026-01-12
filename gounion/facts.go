package gounion

// UnionInterface is a Fact indicating that an interface is a union type
// with a private marker method and a set of implementing types.
type UnionInterface struct {
	MarkerMethod string   // e.g., "isNode"
	Members      []string // e.g., ["*BadExpr", "*Ident", "*BasicLit"]
}

// AFact implements the analysis.Fact interface.
func (*UnionInterface) AFact() {}
