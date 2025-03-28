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
            function: (member_expression
                object: (identifier) @obj
                property: (property_identifier) @method)
            arguments: (arguments (string) @sql))
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
			if q.CaptureName(capture.Index) == "method" && capture.Node.Content(content) == "query" {
				sqlArg := match.Captures[2].Node.Content(content)
				if strings.Contains(sqlArg, "+") || strings.Contains(sqlArg, "${") {
					issues = append(issues, Issue{
						File:     filePath,
						Pos:      token.Pos(capture.Node.StartByte()),
						Severity: "ERRO",
						Message:  "SQL com interpolação - risco de SQL Injection",
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
