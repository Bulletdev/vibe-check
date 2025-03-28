# vibe-check

Uma ferramenta CLI e web em Go para escanear código gerado por IA (ou não) 
e identificar problemas de segurança e boas práticas. 
Suporta Go, Python, Java, JavaScript e integração com CI/CD.

## Propósito

O `vibe-check` ajuda times DevOps e de segurança a lidar com o "vibe coding", detectando vulnerabilidades como SQL Injection, chamadas inseguras e strings hardcoded.

## Funcionalidades

- **Linguagens suportadas**: Go (`.go`), Python (`.py`), Java (`.java`), JavaScript (`.js`).
- **Regras de checagem**:
    - **Go**: `http.Get` sem HTTPS, `os/exec`, SQL inseguro, strings hardcoded, erros não tratados.
    - **Python**: `os.system`, `eval/exec`, SQL inseguro, strings hardcoded.
    - **Java**: `Runtime.exec`, JDBC sem SSL, SQL inseguro, strings hardcoded.
    - **JavaScript**: `eval`, `child_process.exec`, SQL inseguro, strings hardcoded.
- **Escaneamento de pastas**: Analisa múltiplos arquivos recursivamente.
- **Saída**: Texto, JSON ou interface web.
- **Configuração**: Via `.vibe-check.yaml`.
- **CI/CD**: Integração com GitHub Actions.
- **Web**: Visualize relatórios em `http://localhost:4444`.

## Instalação

1. **Pré-requisitos**:
    - Go 1.20+
    - Git

2. **Clone**:
   ```bash
   git clone https://github.com/bulletdev/vibe-check.git
   cd vibe-check
   
3. **Dependências**

````bash
go mod init vibe-check
go get github.com/smacker/go-tree-sitter
go get github.com/smacker/go-tree-sitter/python
go get github.com/smacker/go-tree-sitter/java
go get github.com/smacker/go-tree-sitter/javascript
go get github.com/gin-gonic/gin
go get gopkg.in/yaml.v2
````

4. **Compile**

````bash
go build -o vibe-check main.go
````

5. ***Uso***

````bash
./vibe-check scan --path <caminho> [--json] [--lang go|py|java|js]
````

6. **WEB**
 - Iniciar o Servidor
`````bash
./vibe-check web [--path <caminho>]
`````` 

