package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "loglinter",
		Doc:      "Checks log messages for style and sensitive data",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkCall(pass, call)
	})

	return nil, nil
}

func checkCall(pass *analysis.Pass, call *ast.CallExpr) {
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	obj, ok := pass.TypesInfo.Uses[fun.Sel]
	if !ok {
		return
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return
	}

	pkgPath := pkg.Path()
	if pkgPath != "log/slog" && pkgPath != "go.uber.org/zap" {
		return
	}

	if len(call.Args) == 0 {
		return
	}

	msgArg := call.Args[0]

	pass.Reportf(call.Pos(), "found log call: %v", msgArg)
}
