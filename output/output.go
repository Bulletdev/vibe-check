package output

import (
	"encoding/json"
	"fmt"
	"go/token"
	"vibe-check/rules"
)

func PrintText(issues []rules.Issue) {
	if len(issues) == 0 {
		fmt.Println("Nenhum problema encontrado!")
		return
	}
	fset := token.NewFileSet()
	for _, issue := range issues {
		pos := fset.PositionFor(issue.Pos, false)
		fmt.Printf("%s [%s] Linha %d: %s\n", issue.File, issue.Severity, pos.Line, issue.Message)
	}
}

func PrintJSON(issues []rules.Issue) {
	fset := token.NewFileSet()
	type jsonIssue struct {
		File     string `json:"file"`
		Line     int    `json:"line"`
		Severity string `json:"severity"`
		Message  string `json:"message"`
	}
	var jsonIssues []jsonIssue
	for _, issue := range issues {
		pos := fset.PositionFor(issue.Pos, false)
		jsonIssues = append(jsonIssues, jsonIssue{
			File:     issue.File,
			Line:     pos.Line,
			Severity: issue.Severity,
			Message:  issue.Message,
		})
	}
	output, _ := json.MarshalIndent(jsonIssues, "", "  ")
	fmt.Println(string(output))
}
