package scanner

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"vibe-check/config"
	"vibe-check/output"
	"vibe-check/rules"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
)

func ScanPath(path string, jsonOut bool, lang string) {
	var allIssues []rules.Issue

	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Erro ao acessar o caminho: %v\n", err)
		return
	}

	if !info.IsDir() {
		issues := scanFile(path, lang)
		allIssues = append(allIssues, issues...)
	} else {
		err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(filePath, ".go") || strings.HasSuffix(filePath, ".py") || strings.HasSuffix(filePath, ".java") || strings.HasSuffix(filePath, ".js")) {
				issues := scanFile(filePath, lang)
				allIssues = append(allIssues, issues...)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Erro ao escanear pasta: %v\n", err)
			return
		}
	}

	if jsonOut {
		output.PrintJSON(allIssues)
	} else {
		output.PrintText(allIssues)
	}
}

func scanFile(filePath string, lang string) []rules.Issue {
	if lang == "" {
		if strings.HasSuffix(filePath, ".go") {
			lang = "go"
		} else if strings.HasSuffix(filePath, ".py") {
			lang = "py"
		} else if strings.HasSuffix(filePath, ".java") {
			lang = "java"
		} else if strings.HasSuffix(filePath, ".js") {
			lang = "js"
		} else {
			return nil
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Erro ao ler %s: %v\n", filePath, err)
		return nil
	}

	switch lang {
	case "go":
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, filePath, nil, 0)
		if err != nil {
			fmt.Printf("Erro ao parsear %s: %v\n", filePath, err)
			return nil
		}
		return rules.CheckGo(file, filePath, fset)
	case "py":
		parser := sitter.NewParser()
		parser.SetLanguage(python.GetLanguage())
		tree := parser.Parse(nil, content)
		return rules.CheckPy(tree, filePath, content)
	case "java":
		parser := sitter.NewParser()
		parser.SetLanguage(java.GetLanguage())
		tree := parser.Parse(nil, content)
		return rules.CheckJava(tree, filePath, content)
	case "js":
		parser := sitter.NewParser()
		parser.SetLanguage(javascript.GetLanguage())
		tree := parser.Parse(nil, content)
		return rules.CheckJs(tree, filePath, content)
	default:
		fmt.Printf("Linguagem não suportada: %s\n", lang)
		return nil
	}
}
