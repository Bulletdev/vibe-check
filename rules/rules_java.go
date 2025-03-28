package rules

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"go/token"
	"strings"
)

func CheckJava(tree *sitter.Tree, filePath string, content []byte) []Issue {
	var issues []Issue

	// Query pra capturar chamadas de métodos
	query := `
        (method_invocation
            name: (identifier) @method
            arguments: (argument_list (string_literal) @sql))
    `
	q, _ := sitter.NewQuery([]byte(query), java.GetLanguage())
	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())
	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			if q.CaptureName(capture.Index) == "method" && capture.Node.Content(content) == "execute" {
				sqlArg := match.Captures[1].Node.Content(content)
				if strings.Contains(sqlArg, "+") {
					issues = append(issues, Issue{
						File:     filePath,
						Pos:      token.Pos(capture.Node.StartByte()),
						Severity: "ERRO",
						Message:  "SQL com concatenação - risco de SQL Injection",
					})
				}
			}
		}
	}

	// Regra: Strings hardcoded
	root := tree.RootNode()
	for i := 0; i < int(root.NamedChildCount()); i++ {
		child := root.NamedChild(i)
		if child.Type() == "string_literal" && (strings.Contains(child.Content(content), "localhost") || strings.Contains(child.Content(content), "password")) {
			issues = append(issues, Issue{
				File:     filePath,
				Pos:      token.Pos(child.StartByte()),
				Severity: "AVISO",
				Message:  "String sensível hardcoded - use variáveis de ambiente ou propriedades",
			})
		}
	}

	return issues
}
