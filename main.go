package main

import (
	"flag"
	"fmt"
	"os"
	"vibe-check/scanner"
	"vibe-check/web"
)

func main() {
	scanCmd := flag.NewFlagSet("scan", flag.ExitOnError)
	path := scanCmd.String("path", "", "Caminho do arquivo ou pasta a ser escaneado")
	jsonOut := scanCmd.Bool("json", false, "Saída em formato JSON")
	lang := scanCmd.String("lang", "", "Linguagem específica (go, py, java, js), padrão: detectar automaticamente")

	webCmd := flag.NewFlagSet("web", flag.ExitOnError)
	webPath := webCmd.String("path", ".", "Caminho para escanear no modo web")

	if len(os.Args) < 2 {
		fmt.Println("Uso: vibe-check [scan|web] ...")
		fmt.Println("  scan --path <caminho> [--json] [--lang go|py|java|js]")
		fmt.Println("  web [--path <caminho>]")
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
	case "web":
		webCmd.Parse(os.Args[2:])
		web.StartServer(*webPath)
	default:
		fmt.Println("Comando desconhecido. Use 'scan' ou 'web'.")
		os.Exit(1)
	}
}
