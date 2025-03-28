package rules

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
	"go/token"
	"strings"
)

func CheckPy(tree *sitter.Tree, filePath string, content []byte) []Issue {
	var issues []Issue

	query := `
        (call
            function: (attribute
                object: (identifier) @module
                attribute: (identifier) @func)
            arguments: (argument_list) @args)
    `
	q, _ := sitter.NewQuery([]byte(query), python.GetLanguage())
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
			switch q.CaptureName(capture.Index) {
			case "module":
				if text == "os" && node.NextSibling().Content(content) == "system" {
					issues = append(issues, Issue{
						File:     filePath,
						Pos:      token.Pos(node.StartByte()),
						Severity: "ERRO",
						Message:  "Uso de os.system - risco de injeção de comando",
					})
				}
			case "func":
				if text == "eval" || text == "exec" {
					issues = append(issues, Issue{
						File:     filePath,
						Pos:      token.Pos(node.StartByte()),
						Severity: "ERRO",
						Message:  "Uso de " + text + " - risco de execução arbitrária",
					})
				}
			}
		}
	}

	// Regra: Strings hardcoded
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
