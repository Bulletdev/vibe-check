package rules

import (
	"go/ast"
	"go/token"
	"strings"
)

type Issue struct {
	File     string
	Pos      token.Pos
	Severity string
	Message  string
}

func CheckGo(file *ast.File, filePath string, fset *token.FileSet) []Issue {
	var issues []Issue

	ast.Inspect(file, func(n ast.Node) bool {
		// Regra 1: http.Get inseguro
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := fun.X.(*ast.Ident); ok {
					if pkg.Name == "http" && fun.Sel.Name == "Get" {
						issues = append(issues, Issue{
							File:     filePath,
							Pos:      call.Pos(),
							Severity: "ERRO",
							Message:  "Uso de http.Get sem HTTPS - potencial risco de segurança",
						})
					}
				}
			}
		}

		// Regra 2: Strings hardcoded
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if strings.Contains(lit.Value, "localhost") || strings.Contains(lit.Value, "password") {
				issues = append(issues, Issue{
					File:     filePath,
					Pos:      lit.Pos(),
					Severity: "AVISO",
					Message:  "String sensível hardcoded - use variáveis de ambiente",
				})
			}
		}

		// Regra 3: os/exec inseguro
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := fun.X.(*ast.Ident); ok {
					if pkg.Name == "exec" && (fun.Sel.Name == "Command" || fun.Sel.Name == "Run") {
						issues = append(issues, Issue{
							File:     filePath,
							Pos:      call.Pos(),
							Severity: "ERRO",
							Message:  "Uso de os/exec sem sanitização - risco de injeção de comando",
						})
					}
				}
			}
		}

		// Regra 4: Erros não tratados
		if assign, ok := n.(*ast.AssignStmt); ok {
			for _, rhs := range assign.Rhs {
				if call, ok := rhs.(*ast.CallExpr); ok {
					if !isErrorChecked(call, n) {
						issues = append(issues, Issue{
							File:     filePath,
							Pos:      call.Pos(),
							Severity: "AVISO",
							Message:  "Erro potencialmente não tratado",
						})
					}
				}
			}
		}

		return true
	})

	return issues
}

func isErrorChecked(call *ast.CallExpr, parent ast.Node) bool {
	if stmt, ok := parent.(*ast.AssignStmt); ok {
		for _, lhs := range stmt.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok && ident.Name == "err" {
				return true
			}
		}
	}
	return false
}
