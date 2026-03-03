package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	slogPkg = "log/slog"
	zapPkg  = "go.uber.org/zap"
)

type Config struct {
	SensitiveKeywords []string
}

func NewAnalyzer(cfg Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "loglinter",
		Doc:      "Checks log messages for style and sensitive data",
		Run:      func(p *analysis.Pass) (interface{}, error) { return run(p, cfg) },
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass, cfg Config) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		checkCall(pass, call, cfg)
	})

	return nil, nil
}

func checkCall(pass *analysis.Pass, call *ast.CallExpr, cfg Config) {
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

	msg := extractString(pass, msgArg)
	if msg == "" {
		return
	}

	checkLowercase(pass, call.Pos(), msgArg, msg)
	checkEnglishAndNoSpecials(pass, call.Pos(), msg)
	checkSensitive(pass, call.Pos(), msg, cfg)
}

func extractString(pass *analysis.Pass, expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			if err != nil {
				return ""
			}
			return s
		}
	case *ast.Ident:
		if obj, ok := pass.TypesInfo.Uses[e]; ok {
			if c, ok := obj.(*types.Const); ok {
				val := c.Val()
				if val != nil && val.Kind() == constant.String {
					return constant.StringVal(val)
				}
			}
		}
	}
	return ""
}

func checkLowercase(pass *analysis.Pass, pos token.Pos, arg ast.Expr, msg string) {
	if msg == "" {
		return
	}

	first := []rune(msg)[0]

	if unicode.IsLetter(first) && !unicode.IsLower(first) {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			newMsg := string(unicode.ToLower(first)) + msg[1:]
			newLit := strconv.Quote(newMsg)

			edit := analysis.TextEdit{
				Pos:     lit.Pos(),
				End:     lit.End(),
				NewText: []byte(newLit),
			}

			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				Message: "log message should start with a lowercase letter",
				SuggestedFixes: []analysis.SuggestedFix{{
					Message:   "lowercase first letter",
					TextEdits: []analysis.TextEdit{edit},
				}},
			})

			return
		}

		pass.Reportf(pos, "log message should start with a lowercase letter")
	}
}

func checkEnglishAndNoSpecials(pass *analysis.Pass, pos token.Pos, msg string) {
	for _, r := range msg {
		if !(unicode.IsLetter(r) && r <= 0x7F) && !unicode.IsDigit(r) && r != ' ' {
			pass.Reportf(pos, "log message should contain only English letters, digits, and spaces (no special characters or emojis)")

			return
		}
	}
}

func checkSensitive(pass *analysis.Pass, pos token.Pos, msg string, cfg Config) {
	keywords := cfg.SensitiveKeywords

	if len(keywords) == 0 {
		keywords = []string{"password", "api_key", "token", "secret", "key", "auth"}
	}

	lowerMsg := strings.ToLower(msg)

	for _, kw := range keywords {
		if strings.Contains(lowerMsg, kw) {
			pass.Reportf(pos, "log message may contain sensitive data: %q", kw)

			return
		}
	}
}
