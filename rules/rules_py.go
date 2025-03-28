package rules

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
	"go/token"
	"strings"
)

func CheckPy(tree *sitter.Tree, filePath string, content []byte) []Issue {
	var issues []Issue

	// Query combinada pra capturar SQL inseguro e eval/exec
	query := `
        (call
            function: (attribute
                object: (identifier) @obj
                attribute: (identifier) @method)
            arguments: (argument_list (string) @sql))
        (call
            function: (identifier) @func
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
			captureName := q.CaptureName(capture.Index)
			text := string(node.Content(content))

			// Regra: SQL inseguro
			if captureName == "method" && text == "execute" {
				// O argumento SQL é o terceiro capture (@sql) na query
				for _, c := range match.Captures {
					if q.CaptureName(c.Index) == "sql" {
						sqlArg := c.Node.Content(content)
						if strings.Contains(sqlArg, "%") || strings.Contains(sqlArg, "+") {
							issues = append(issues, Issue{
								File:     filePath,
								Pos:      token.Pos(node.StartByte()),
								Severity: "ERRO",
								Message:  "SQL com formatação direta - risco de SQL Injection",
							})
						}
					}
				}
			}

			// Regra: eval ou exec
			if captureName == "func" && (text == "eval" || text == "exec") {
				issues = append(issues, Issue{
					File:     filePath,
					Pos:      token.Pos(node.StartByte()),
					Severity: "ERRO",
					Message:  "Uso de " + text + " - risco de execução arbitrária",
				})
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
