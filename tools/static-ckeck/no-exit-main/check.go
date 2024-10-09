package noexitmain

import (
	"go/ast"
	"go/token"
)

// DetectExitMain detects os.Exit call in func main from main pkg
func DetectExitMain(node ast.Node) *token.Pos {
	var pos *token.Pos

	// Check if node is a file and belongs to the main package
	file, ok := node.(*ast.File)
	if !ok || file.Name.Name != "main" {
		return nil
	}

	// Look for the main function
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "main" || fn.Body == nil {
			continue
		}

		// Inspect the body for os.Exit
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if ok {
				ident, ok := selExpr.X.(*ast.Ident)
				if ok && ident.Name == "os" && selExpr.Sel.Name == "Exit" {
					pos = &ident.NamePos
					return false // Stop searching after finding os.Exit
				}
			}
			return true
		})
		break // No need to check other functions once we found the main
	}

	return pos
}
