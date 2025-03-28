package rules

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"strings"
)

func CheckJs(tree *sitter.Tree, filePath string, content []byte) []Issue {
	var issues []Issue

	query := `
        (call_expression
            function: (identifier) @func
            arguments: (arguments) @args)
    `
	q, _ := sitter.NewQuery([]byte(query), javascript.GetLanguage())
	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			text := string(node.Content(content))
			if q.CaptureName(capture.Index) == "func" {
				if text == "eval" {
					issues = append(issues, Issue{
						File:     filePath,
						Pos:      token.Pos(node.StartByte()),
						Severity: "ERRO",
						Message:  "Uso de eval - risco de execução arbitrária",
					})
				} else if text == "exec" && strings.Contains(node.Parent().Content(content), "child_process") {
					issues = append(issues, Issue{
						File:     filePath,
						Pos:      token.Pos(node.StartByte()),
						Severity: "ERRO",
						Message:  "Uso de child_process.exec - risco de injeção de comando",
					})
				}
			}
		}
	}

	// Strings hardcoded
	root := tree.RootNode()
	for i := 0; i < int(root.NamedChildCount()); i++ {
		child := root.NamedChild(i)
		if child.Type() == "string" && (strings.Contains(child.Content(content), "localhost") || strings.Contains(child.Content(content), "password")) {
			issues = append(issues, Issue{
				File:     filePath,
				Pos:      token.Pos(child.StartByte()),
				Severity: "AVISO",
				Message:  "String sensível hardcoded - use variáveis de ambiente",
			})
		}
	}

	return issues
}
