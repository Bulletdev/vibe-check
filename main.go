package main

import (
	"flag"
	"fmt"
	"os"
	"vibe-check/scanner"
)

func main() {
	scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
	path := scanCmd.String("path", "", "Caminho do arquivo ou pasta a ser escaneado")
	jsonOut := scanCmd.Bool("json", false, "Saída em formato JSON")
	lang := scanCmd.String("lang", "", "Linguagem específica (go, py), padrão: detectar automaticamente")

	if len(os.Args) < 2 {
		fmt.Println("Uso: vibe-check scan --path <caminho> [--json] [--lang go|py]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		scanCmd.Parse(os.Args[2:])
		if *path == "" {
			fmt.Println("Erro: informe o caminho com --path")
			os.Exit(1)
		}
		scanner.ScanPath(*path, *jsonOut, *lang)
	default:
		fmt.Println("Comando desconhecido. Use 'scan'.")
		os.Exit(1)
	}
}
