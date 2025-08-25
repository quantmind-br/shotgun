# Plano de Desenvolvimento — shotgun-cli (Go + Bubble Tea)

# Visão Geral
shotgun-cli é uma aplicação de terminal (TUI) escrita em Go usando Bubble Tea para gerar prompts para LLMs a partir do preenchimento de templates pré‑estabelecidos ou customizados com interface minimalista e elegante.

Problema que resolve
- Padroniza e acelera a criação de prompts complexos, evitando copiar/colar manual e variações de formatação.

Para quem é
- Desenvolvedores e times que usam LLMs em fluxos de desenvolvimento, revisão de código, planejamento e debugging.

Valor
- Reprodutibilidade, consistência e velocidade ao gerar prompts com contexto de arquivos do repositório.

Decisão aprovada nesta etapa
- Modo de interação: TUI-only (wizard full-screen, navegação 100% por teclado; campos multiline com suporte a colar).
- Linguagem: **Go 1.22+** com framework Bubble Tea v2.0.0-beta.4 para TUI minimalista e performática com compatibilidade Windows aprimorada.
- Design: Interface limpa com paleta monocromática e acentos sutis.

## Fluxo do Programa (5 Telas)
**Ao executar `shotgun-cli`, inicia-se imediatamente o mapeamento e exibição da árvore de arquivos do diretório atual.**

1. **File Tree** - Seleção de arquivos/pastas para exclusão (checkboxes iniciam marcados)
2. **Template Selection** - Escolha do template em lista vertical  
3. **Task Input** - Descrição da tarefa (editor multiline)
4. **Rules Input** - Regras opcionais (editor multiline)
5. **Confirmation** - Revisão e geração do arquivo .md

**Comportamento:** O programa trabalha exclusivamente com o diretório onde foi executado. Arquivos ignorados por `.gitignore` e `.shotgunignore` são automaticamente excluídos.

## Sistema de Navegação e Atalhos

### Atalhos Globais (disponíveis em todas as telas)
- **F1** - Ajuda contextual
- **F2** - Voltar para tela anterior  
- **F3** - Avançar para próxima tela
- **ESC** - Sair do programa (com confirmação)

### Atalhos Específicos por Tela

#### Tela 1: File Tree
- **↑↓** - Navegar entre itens
- **→/←** - Expandir/colapsar pastas
- **Space** - Marcar/desmarcar item
- **F3** - Avançar para Template Selection

#### Tela 2: Template Selection  
- **↑↓** - Navegar entre templates
- **Enter** - Selecionar template e avançar
- **F2** - Voltar para File Tree
- **F3** - Confirmar seleção e avançar

#### Tela 3: Task Input
- **Modo Edição** - Digitação normal de texto
- **Ctrl+Enter** - Finalizar edição e avançar
- **F2** - Voltar (salva conteúdo)
- **F3** - Avançar (requer conteúdo)

#### Tela 4: Rules Input  
- **Modo Edição** - Digitação normal de texto
- **Ctrl+Enter** - Finalizar edição e avançar
- **F2** - Voltar (salva conteúdo)
- **F3** - Avançar (campo opcional)
- **F4** - Pular esta etapa

#### Tela 5: Confirmation
- **F2** - Voltar para ajustes
- **F10** - Confirmar e gerar prompt
- **ESC** - Cancelar geração

**Nota:** Durante edição de texto (Task/Rules), teclas F são desabilitadas exceto F2/F3 após Ctrl+Enter para evitar conflitos.

## Comportamento de Estado e Navegação

### Persistência de Estado
- Cada tela mantém seu estado ao navegar (F2/F3)
- Checkboxes, seleções e textos são preservados
- Usuário pode voltar e ajustar qualquer tela
- Estado só é perdido ao sair (ESC) ou gerar (F10)

### Validações
- **Tela 1**: Pelo menos 1 arquivo deve estar selecionado
- **Tela 2**: Um template deve ser selecionado
- **Tela 3**: Task não pode estar vazia
- **Tela 4**: Rules é opcional (pode estar vazio)
- **Tela 5**: Confirmação explícita com F10

### Feedback Visual
- Indicador de progresso no header: [1/5], [2/5], etc
- Campos obrigatórios indicados claramente
- Warnings em amarelo (#F1FA8C) para alertas
- Erros em vermelho (#FF5555) para problemas

Premissas iniciais
- Linguagem: **Go 1.22+** (para generics otimizados, performance superior e compatibilidade Windows aprimorada).
- TUI Framework: **Bubble Tea** com Lip Gloss para estilização elegante.
- Execução global via comando: shotgun-cli (binário único).
- Suporte a Windows, Linux e macOS; terminais comuns (PowerShell, WezTerm, Bash, iTerm2 etc.).
- Templates base embarcados no binário, com possibilidade de criação de templates customizados pelo usuário em diretórios específicos do sistema.
- Exemplos e templates de referência:
  • **Templates base incluídos**:
    - [templates/prompt_analyzeBug.toml](templates/prompt_analyzeBug.toml) — Template para análise de bugs com trace de execução detalhado
    - [templates/prompt_makeDiffGitFormat.toml](templates/prompt_makeDiffGitFormat.toml) — Template para gerar diffs Git formatados a partir de código
    - [templates/prompt_makePlan.toml](templates/prompt_makePlan.toml) — Template para criação de planos arquiteturais e de refatoração
    - [templates/prompt_projectManager.toml](templates/prompt_projectManager.toml) — Template para sincronização de documentação de projetos
  • **Exemplo de prompt final**: [exemplos/ex_prompt_dev.md](exemplos/ex_prompt_dev.md) — Demonstra o formato esperado do arquivo .md gerado

# Stack Tecnológica Go

## Core Technologies
- **Go 1.22+** - Linguagem compilada com excelente performance, concorrência nativa e melhorias para Windows
- **Bubble Tea v2.0.0-beta.4** - Framework TUI elegante e reativo com keyboard enhancements refinado, debug de panics com TEA_DEBUG, e API v2 moderna
- **Bubbles v0.21.0** - Componentes prontos (filepicker, textarea, list, viewport) com horizontal scrolling
- **Lip Gloss v1.0.0** - Estilização terminal moderna com gradientes, layouts flexíveis e padding/margin customizáveis
- **Cobra** - CLI framework robusto para comandos e flags
- **text/template** - Template engine nativo do Go com segurança de tipos

## Bibliotecas Auxiliares
- **github.com/BurntSushi/toml** - Parser TOML eficiente
- **github.com/go-git/go-git/v5** - Manipulação de repositórios Git
- **github.com/bmatcuk/doublestar/v4** - Implementação de glob patterns
- **github.com/h2non/filetype** - Detecção de tipos de arquivo e binários
- **github.com/charmbracelet/bubbles** - Componentes prontos (textinput, viewport, spinner)
- **github.com/charmbracelet/glamour** - Renderização Markdown com estilo
- **github.com/spf13/viper** - Configuração estruturada e flexível
- **unicode/utf8** - Validação e manipulação de caracteres UTF-8 (biblioteca padrão)
- **golang.org/x/text** - Normalização e transformação de texto Unicode

# Funcionalidades Principais

1) Árvore de Arquivos com Exclusão (respeita .gitignore e .shotgunignore)
- O que faz: Exibe a estrutura hierárquica do diretório onde o programa foi executado; **todos os checkboxes iniciam marcados** e o usuário desmarca quais arquivos/pastas NÃO irão para o prompt. **Não há opção de navegar para outros diretórios - o programa trabalha exclusivamente com o diretório atual.** Arquivos/pastas ignorados por .gitignore e .shotgunignore são automaticamente excluídos da listagem.
- Por que é importante: Permite foco no contexto relevante e evita ruído (builds, dependências, artefatos).
- **Seleção Hierárquica**: Ao desmarcar uma pasta, todos os arquivos e subpastas dentro dela são automaticamente desmarcados. Ao marcar uma pasta, todos os itens dentro dela são marcados.
- **Feedback Visual para Binários**: Arquivos detectados como binários são visualmente distintos (ícone diferenciado, cor cinza, não selecionáveis) e automaticamente ignorados, deixando claro para o usuário por que não podem ser incluídos.
- Como funciona (alto nível): Varredura concorrente recursiva do diretório atual usando goroutines; aplicação das regras do .gitignore e .shotgunignore com `doublestar`; renderização de árvore TUI com componente tree customizado do Bubble Tea; seleção hierárquica por teclado com feedback visual suave; saída = conjunto de arquivos incluídos.

2) Suporte a `.shotgunignore` e Comando `init`
- O que faz: A CLI procura por um arquivo `.shotgunignore` na raiz do projeto para aplicar regras de exclusão adicionais, específicas do projeto. Um novo comando `shotgun-cli init` pode ser usado para criar um arquivo `.shotgunignore` de exemplo.
- Por que é importante: Oferece um controle de exclusão granular, explícito e versionável no Git, ideal para padronizar o comportamento da ferramenta em equipes e para garantir que arquivos sensíveis específicos do projeto nunca sejam incluídos.
- Como funciona: As regras no formato `.gitignore` dentro de `.shotgunignore` são processadas pelo `pathspec` e adicionadas ao conjunto de exclusões antes da varredura de arquivos.

3) Seleção de Template com Interface de Lista
- O que faz: Permite escolher um template de prompt em uma interface de lista simples e eficiente (inclui os 4 fornecidos e quaisquer .toml encontrados nos diretórios de templates do usuário, todos exibidos da mesma forma).
- Por que é importante: Padroniza e acelera a formatação do prompt para diferentes tarefas (debug, git diff, planejamento, PM docs).
- Como funciona: Descoberta de templates embarcados + templates dos diretórios ~/.config/shotgun-cli/templates (Linux/Mac) ou %APPDATA%/shotgun-cli/templates (Windows); **templates customizados aparecem na mesma lista sem diferenciação visual**; interface de lista usando list bubble do Bubble Tea; metadata rica (título, descrição, versão, tags) extraída do TOML; navegação por teclado com setas; seleção com Enter; animações suaves de transição.

4) Entrada Multilinha da Tarefa (TASK) com Editor Avançado e Suporte Unicode
- O que faz: Campo multiline com editor avançado para o usuário descrever a tarefa/objetivo do prompt; suporta colar texto, caracteres especiais (ç, á, ô, ñ, etc.), syntax highlighting básico, word wrap.
- Por que é importante: Aumenta a qualidade do prompt com contexto detalhado em qualquer idioma e experiência de edição superior internacional.
- Como funciona: Componente textarea do bubbles v0.21.0 com preservação de quebras de linha e suporte nativo UTF-8; syntax highlighting suave para markdown usando Lip Gloss; validação de runes com unicode/utf8; substitui o placeholder {TASK} do template via text/template do Go.

5) Entrada Multilinha de Regras (RULES) — Opcional com Validação Unicode
- O que faz: Campo multiline opcional para regras/constraints com validação; suporta colar texto, caracteres internacionais e templates pré-definidos em qualquer idioma.
- Por que é importante: Permite ajustar o comportamento do LLM às políticas do time/projeto com validação robusta e suporte internacional.
- Como funciona: Buffer multiline opcional com validação estruturada UTF-8; suporte nativo a acentos e caracteres especiais; substitui o placeholder {RULES} do template (se vazio, insere "N/A" ou mantém seção vazia conforme template text/template).

6) Montagem Assíncrona de "File Structure"
- O que faz: Insere no prompt uma seção completa com a árvore hierárquica do projeto seguida pelo conteúdo detalhado dos arquivos incluídos, processado de forma assíncrona.
- Por que é importante: Fornece ao LLM tanto a visão estrutural quanto o contexto detalhado de código/arquivos necessário para análise/geração precisa, sem bloquear a UI.
- Como funciona: 
  1. Gera árvore completa de diretórios (formato tree-like com caracteres ASCII: ├── └── │) usando processamento assíncrono
  2. Em seguida, para cada arquivo incluído, lê o conteúdo de forma assíncrona e insere bloco:
     <file path="RELATIVE/PATH/TO/FILE">
     (conteúdo do arquivo)
     </file>
  Apenas arquivos de texto são incluídos; arquivos binários são automaticamente excluídos via `filetype`. O gerenciamento do tamanho total fica sob responsabilidade do usuário. O placeholder {FILE_STRUCTURE} do template é preenchido com essa saída via text/template.

7) Estimativa de Tamanho em Tempo Real e Confirmação
- O que faz: Apresenta uma estimativa em tempo real do tamanho do arquivo final a ser gerado com progress bar e solicita confirmação do usuário antes de prosseguir.
- Por que é importante: Permite ao usuário avaliar se o prompt resultante será adequado para uso, evitando surpresas com arquivos muito grandes.
- Como funciona: Calcula tamanho estimado de forma assíncrona baseado no template + variáveis + file tree + conteúdo dos arquivos; exibe em KB/MB com indicadores visuais; progress bar durante cálculo; permite confirmar para prosseguir ou voltar para ajustar seleções.

8) Geração Assíncrona do Prompt .md na Pasta Atual
- O que faz: Gera e salva o arquivo final em Markdown após confirmação usando processamento assíncrono.
- Por que é importante: Facilita versionamento, compartilhamento e uso subsequente do prompt sem bloquear a interface.
- Como funciona: Ao confirmar no passo anterior, processa o template text/template de forma concorrente e escreve o arquivo no diretório atual com nome padrão "shotgun_prompt_YYYYMMDD_HHMM.md" para evitar sobrescritas.

9) Navegação 100% por Teclado (TUI) com Interface Moderna
- O que faz: Permite operar todo o fluxo sem mouse, incluindo navegação livre entre todas as telas do wizard com interface moderna do Bubble Tea.
- Por que é importante: Agilidade, compatibilidade com diferentes terminais/OS e flexibilidade para ajustar configurações a qualquer momento.
- Como funciona: Navegação global via teclas F (F1 ajuda; F2 tela anterior; F3 próxima tela; F4-F10 acesso direto às telas 1-7) usando key messages do Bubble Tea; teclas contextuais variam por tipo de tela (árvore: setas para navegar, espaço para marcar; multiline: Enter para nova linha, Ctrl+Enter para alternar modo edição/navegação). Estado de cada tela é preservado durante navegação via model state immutável.

10) Templates Customizados do Usuário com Validação Robusta
- O que faz: O usuário pode criar templates .toml estruturados em diretórios específicos do sistema e eles aparecem automaticamente na seleção com precedência sobre templates empacotados, com validação via structs Go.
- Por que é importante: Permite personalização global e reutilização de templates entre projetos, com metadata rica e validação automática robusta.
- Como funciona: Descoberta de templates em ~/.config/shotgun-cli/templates (Linux/Mac) ou %APPDATA%/shotgun-cli/templates (Windows); templates do usuário são listados primeiro com metadata completa (título, versão, descrição, tags) validada por structs Go; suporte a organização em subdiretórios; validação de estrutura TOML; templates embarcados servem como fallback.

11) Descoberta Automática e Validação Avançada de Variáveis
- O que faz: Lê a seção [variables] do template TOML para descobrir todas as variáveis necessárias, seus tipos, obrigatoriedade e validações usando struct tags do Go.
- Por que é importante: Permite templates mais flexíveis com validação robusta e experiência de usuário superior.
- Como funciona: Parsing da seção [variables] do TOML com validação via structs Go e tags; tipos suportados (text, multiline, auto, choice, boolean, number); validação automática de obrigatoriedade e constraints; valores padrão e placeholders; geração dinâmica de campos de entrada no wizard.
- Tipos de variáveis suportados:
  • **text**: Campo de linha única (ex: título, nome)
  • **multiline**: Campo de múltiplas linhas (ex: TASK, RULES)
  • **auto**: Variáveis populadas automaticamente (ex: FILE_STRUCTURE, CURRENT_DATE)
  • **choice**: Seleção entre opções predefinidas
  • **boolean**: Sim/não para seções condicionais
  • **number**: Campos numéricos com validação de range
- Template engine: text/template do Go para lógica condicional ({{if}}, {{range}}, {{define}})

# Arquitetura Go Minimalista

## Estrutura do Projeto
```
shotgun-cli/
├── cmd/
│   └── shotgun/
│       └── main.go              # Entry point
├── internal/
│   ├── app/
│   │   ├── app.go              # Aplicação principal Bubble Tea
│   │   ├── model.go            # Model state da aplicação
│   │   └── keys.go             # Keybindings globais
│   ├── screens/                # Telas do wizard
│   │   ├── filetree/
│   │   │   ├── model.go        # File tree model
│   │   │   ├── view.go         # File tree view
│   │   │   └── update.go       # File tree update
│   │   ├── template/
│   │   │   ├── model.go        # Template selection
│   │   │   ├── view.go         # Template list view
│   │   │   └── update.go       # Template update
│   │   ├── input/
│   │   │   ├── task.go         # Task input
│   │   │   └── rules.go        # Rules input
│   │   └── confirm/
│   │       ├── model.go        # Confirmation screen
│   │       └── view.go         # Summary view
│   ├── components/             # Componentes reutilizáveis
│   │   ├── tree/               # Tree widget
│   │   ├── editor/             # Text editor
│   │   ├── list/               # List selector
│   │   ├── progress/           # Progress bar
│   │   └── statusbar/          # Status bar
│   ├── core/                   # Lógica de negócio
│   │   ├── scanner.go          # File scanner concorrente
│   │   ├── template.go         # Template engine
│   │   ├── ignore.go           # .gitignore handler
│   │   ├── builder.go          # Prompt builder
│   │   └── config.go           # Configuration
│   ├── models/                 # Data models
│   │   ├── file.go             # File tree item
│   │   ├── template.go         # Template structure
│   │   ├── variable.go         # Variable definitions
│   │   └── state.go            # App state
│   ├── styles/                 # Lip Gloss styles
│   │   ├── theme.go            # Tema minimalista
│   │   ├── colors.go           # Paleta de cores
│   │   └── components.go       # Estilos de componentes
│   └── utils/                  # Utilitários
│       ├── files.go            # File utilities
│       └── validation.go       # Validators
├── templates/                  # Templates embarcados
│   ├── prompt_analyzeBug.toml
│   ├── prompt_analyzeBug.md
│   ├── prompt_makeDiffGitFormat.toml
│   ├── prompt_makeDiffGitFormat.md
│   ├── prompt_makePlan.toml
│   ├── prompt_makePlan.md
│   ├── prompt_projectManager.toml
│   └── prompt_projectManager.md
├── go.mod                      # Dependências Go
├── go.sum
├── Makefile                    # Build automation
├── README.md
└── CHANGELOG.md
```

## Arquitetura Concorrente com Bubble Tea

### Model Principal (Elm Architecture)
```go
type Model struct {
    state       AppState
    currentView View
    scanner     *FileScanner
    template    *TemplateEngine
    width       int
    height      int
}

func InitialModel() Model {
    return Model{
        state:       NewAppState(),
        currentView: FileTreeView,
        scanner:     NewFileScanner(),
        template:    NewTemplateEngine(),
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case FileScanCompleteMsg:
        m.state.Files = msg.Files
    }
    return m, nil
}

func (m Model) View() string {
    return m.currentView.Render(m)
}
```

### Concorrência com Goroutines
- **FileScanner**: Varre diretórios em paralelo usando goroutines
- **ContentReader**: Lê arquivos concorrentemente com channels
- **TemplateProcessor**: Processa templates com worker pool
- **PromptGenerator**: Monta prompt final usando pipeline pattern

# Design da Interface Minimalista

## Paleta de Cores (Minimal Monochrome)
```go
// colors.go - Minimal Color Palette
package styles

import "github.com/charmbracelet/lipgloss"

var (
    // Base Colors - Grayscale
    ColorBackground = lipgloss.Color("#0a0a0a")  // Almost black
    ColorSurface    = lipgloss.Color("#141414")  // Dark gray
    ColorBorder     = lipgloss.Color("#2a2a2a")  // Medium gray
    ColorMuted      = lipgloss.Color("#505050")  // Muted gray
    
    // Text Colors
    ColorText       = lipgloss.Color("#e8e8e8")  // Light gray
    ColorTextDim    = lipgloss.Color("#808080")  // Dimmed text
    ColorTextBright = lipgloss.Color("#ffffff")  // Pure white
    
    // Accent Colors - Minimal
    ColorAccent     = lipgloss.Color("#6ee7b7")  // Soft mint green
    ColorHighlight  = lipgloss.Color("#fbbf24")  // Warm amber
    ColorSuccess    = lipgloss.Color("#86efac")  // Light green
    ColorWarning    = lipgloss.Color("#fde047")  // Soft yellow
    ColorError      = lipgloss.Color("#fca5a5")  // Soft red
    
    // Styles
    BaseStyle = lipgloss.NewStyle().
        Background(ColorBackground).
        Foreground(ColorText)
    
    BorderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(ColorBorder).
        Padding(1, 2)
    
    FocusedStyle = BorderStyle.Copy().
        BorderForeground(ColorAccent)
    
    SelectedStyle = lipgloss.NewStyle().
        Background(ColorSurface).
        Foreground(ColorTextBright)
    
    MutedStyle = lipgloss.NewStyle().
        Foreground(ColorTextDim)
)
```

## Componentes Visuais Minimalistas

### Header Minimalista com Progress
```
 shotgun-cli                                      
 File Selection                          [1/5]    [Texto: #e8e8e8]
 ━━━━━━────────────────────────────────          [Progress: #6ee7b7]
```

### Tela 1: File Tree Minimalista
```
 File Selection                          [1/5]    
 ━━━━━━────────────────────────────────          

 [✓] README.md                                    [Texto: #e8e8e8]
 [✓] go.mod                                       [Check: #6ee7b7]
 [✓] PLAN.md                                      
 [✓] cmd/                               ▸        [Pasta: símbolo sutil]
 [ ] tests/                             ▸        [Desmarcado: #505050]
 [✓] internal/                          ▾        [Expandido]
   └ [✓] app/                                   [Indentação limpa]
       [✓] app.go                               
       [✓] model.go                             
   └ [✓] templates/                             
       [✓] analyze_bug.toml                     
       [✓] make_diff.toml                       
   · .env                              ignored   [Ignorado: #505050]
   · .git/                             ignored   

 8 selected · 1 excluded · 12 ignored            [Status: #808080]

 Space toggle · ←→ expand · ↑↓ navigate          [Ajuda: #505050]
 F3 next · ESC exit                              
```

### Tela 2: Template Selection Minimalista
```
 Template Selection                      [2/5]    
 ━━━━━━━━━━────────────────────────────          

 › analyze_bug                          v2.1.0    [Seleção: #6ee7b7]
   Bug analysis with execution traces             [Descrição: #808080]

   make_plan                            v1.3.0    [Normal: #e8e8e8]
   Create detailed project plans                  

   make_diff                            v1.0.0    
   Generate formatted git diffs                   

   project_manager                      v1.2.0    
   Sync project documentation                     

   my_custom_template                   v1.0.0    
   User custom template for reviews               


 ↑↓ navigate · Enter select                       [Ajuda: #505050]
 F2 back · F3 next · ESC exit                     
```

### Tela 3: Task Input Minimalista
```
 Task Description                        [3/5]    
 ━━━━━━━━━━━━━━────────────────────────          
 editing                                          [Modo: #6ee7b7]

 ┌────────────────────────────────────────┐      [Borda sutil: #2a2a2a]
 │ Analyze the authentication bug in the  │      [Texto: #e8e8e8]
 │ login system. The error occurs when    │      [Fundo: #141414]
 │ users try to login with special chars  │
 │ in their password.                     │
 │                                        │
 │ Steps to reproduce:                    │
 │ 1. Go to login page                    │
 │ 2. Enter email: test@example.com       │
 │ 3. Enter password with symbols: P@ss!  │
 │ 4. Click login                         │
 │ 5. Error appears: "Invalid credentials"│
 │                                        │
 │ ▌                                      │      [Cursor minimalista]
 └────────────────────────────────────────┘

 11 lines · 456 chars                            [Info: #505050]

 Ctrl+Enter finish                               
 F2 back · F3 next · ESC exit                    
```

### Tela 4: Rules Input Minimalista
```
 Rules · optional                        [4/5]    
 ━━━━━━━━━━━━━━━━──────────────────────          
 editing                                          

 ┌────────────────────────────────────────┐      
 │ Use Go best practices                  │      [Texto: #e8e8e8]
 │ Follow clean architecture              │      
 │ Include comprehensive tests            │      
 │                                        │      
 │ ▌                                      │      
 └────────────────────────────────────────┘      

 3 lines · 89 chars                              
 optional field                                  [Info: #505050]

 Ctrl+Enter finish · F4 skip                     
 F2 back · F3 next · ESC exit                    
```

### Tela 5: Confirmation Minimalista
```
 Confirm Generation                      [5/5]    
 ━━━━━━━━━━━━━━━━━━━━──────────────────          

 Summary                                          [Título: #e8e8e8]
 ─────────────────────────────────────           [Linha: #2a2a2a]

 Template        analyze_bug v2.1.0               
 Selected        8 files                          
 Excluded        1 folder (tests/)                
 Lines           ~2,340                           

 Output                                           
 shotgun_prompt_20240315_1430.md                 [Nome: #6ee7b7]

 Size estimate                                    
 ████████████████████░░░░░░  423 KB              [Bar: #6ee7b7/#2a2a2a]

 ⚠ Large file - consider excluding more          [Aviso: #fde047]


 F2 back · F10 generate · ESC cancel             
```

# Implementação das 5 Telas

## Tela 1: File Tree Selection
- Baseado no componente **filepicker v0.21.0** do Bubbles com melhorias de navegação e performance
- Customização do filepicker para seleção múltipla com checkboxes e suporte aprimorado a keyboard navigation
- **Todos os checkboxes iniciam marcados por padrão**
- **Pastas iniciam colapsadas (exceto raiz)**
- Navegação hierárquica com teclas → ou ← para expandir/colapsar
- **Seleção hierárquica**: desmarcar pasta desmarca todo conteúdo
- Visual feedback para binários com cor accent (ícone 🔒, cor #4a4a5a)
- Indicadores de .gitignore com estilo monocromático (#4a4a5a)
- Viewport virtual para listas grandes usando viewport bubble
- **Sem navegação para diretórios externos - trabalha apenas com o diretório de execução**

## Tela 2: Template Selection
- Baseado no componente **list** do Bubbles para seleção
- Interface de lista vertical com paleta monocromática
- Navegação com setas para cima/baixo
- Seleção com Enter ou F3
- Exibição unificada de templates built-in e customizados
- Nome do template com versão inline (cor #bd93f9 para versão)
- Descrição em linha secundária (cor #6272a4)
- Destaque visual no item selecionado (fundo #44475a, borda #bd93f9)

## Tela 3: Task Input
- Baseado no componente **textarea v0.21.0** do Bubbles com horizontal scrolling e melhor handling de text wrap
- Editor multiline com syntax highlighting suave (paleta monocromática) e suporte aprimorado para linhas longas
- Line numbers em #6272a4 e indicador de posição
- Word wrap inteligente
- Clipboard support (Ctrl+V/Ctrl+C)
- Character/word count em tempo real (cor #8be9fd)
- Ctrl+Enter para finalizar edição
- F-keys desabilitadas durante edição
- Borda em #bd93f9 quando focado

## Tela 4: Rules Input
- Campo opcional com indicação clara (#6272a4 para "optional")
- Baseado no componente **textarea v0.21.0** do Bubbles com as mesmas melhorias de horizontal scrolling
- Mesmos recursos do editor de Task com paleta monocromática e performance aprimorada
- F4 para pular esta etapa
- Auto-save ao navegar
- Indicador visual de campo opcional com ícone ℹ

## Tela 5: Confirmation & Generation
- Resumo completo antes de gerar (styled com Lip Gloss)
- Cálculo concorrente de tamanho estimado
- Progress bar visual durante estimativa (#50fa7b para progresso, #44475a para fundo) usando viewport com melhor responsividade
- Warnings para arquivos muito grandes (#f1fa8c) com estimativas mais precisas
- F10 para confirmar geração
- Nome do arquivo com timestamp automático (#50fa7b)
- Geração concorrente com goroutines sem bloquear UI, otimizada para Go 1.22+
- Bordas arredondadas estilo minimalista (#4a4a5a) renderizadas com Lip Gloss v1.0.0

# Sistema de Templates Go text/template

## Estrutura TOML Completa
```toml
[meta]
name = "analyze_bug"
version = "2.1.0"
description = "Comprehensive bug analysis template"
author = "Shotgun Team"
tags = ["debug", "analysis", "bug-fix"]
category = "development"

[variables]
# Campo de texto simples
title = { type = "text", required = true, placeholder = "Bug title", max_length = 100 }

# Campo multiline com validação
task = { 
    type = "multiline", 
    required = true, 
    placeholder = "Describe the bug in detail...",
    min_lines = 3,
    max_lines = 50
}

# Campo opcional com valor padrão
rules = { 
    type = "multiline", 
    required = false, 
    default = "Follow standard debugging practices",
    placeholder = "Additional rules or constraints..."
}

# Seleção de opções
context = { 
    type = "choice", 
    options = ["frontend", "backend", "fullstack", "mobile"],
    default = "fullstack",
    required = true
}

# Campo numérico com range
priority = { 
    type = "number", 
    min = 1, 
    max = 5, 
    default = 3,
    description = "Bug priority (1=low, 5=critical)"
}

# Campo booleano para seções condicionais
include_logs = {
    type = "boolean",
    default = true,
    description = "Include system logs in analysis"
}

# Variáveis automáticas
file_structure = { type = "auto", source = "files" }
current_date = { type = "auto", source = "date" }
project_name = { type = "auto", source = "dirname" }

[template]
content = """
# Bug Analysis: {{.Title}}

**Context**: {{.Context | title}}  
**Priority**: {{.PriorityStars}} ({{.Priority}}/5)  
**Date**: {{.CurrentDate}}  
**Project**: {{.ProjectName}}

## Task Description
{{.Task}}

{{if and (ne .Rules "N/A") (ne .Rules "")}}
## Rules & Constraints
{{.Rules}}
{{end}}

{{if .IncludeLogs}}
## System Context
*Analysis includes system logs and traces*
{{end}}

## Project Structure & Files
{{.FileStructure}}

{{if ge .Priority 4}}
---
⚠ **HIGH PRIORITY ISSUE** - Requires immediate attention
{{end}}
"""

[validation]
# Validações customizadas
custom_validators = [
    "task_should_mention_steps",
    "priority_justification_if_high"
]

[ui]
# Configurações de UI específicas
step_order = ["title", "context", "priority", "task", "rules", "include_logs"]
group_layout = {
    "Basic Info" = ["title", "context", "priority"],
    "Details" = ["task", "rules"],
    "Options" = ["include_logs"]
}
```

## Validação e Suporte UTF-8
```go
package utils

import (
    "unicode/utf8"
    "golang.org/x/text/unicode/norm"
    "strings"
)

// ValidateUTF8Input valida se a string contém UTF-8 válido
func ValidateUTF8Input(s string) bool {
    return utf8.ValidString(s)
}

// NormalizeText normaliza texto Unicode (NFD -> NFC)
func NormalizeText(text string) string {
    return norm.NFC.String(text)
}

// SafeStringLength conta corretamente caracteres Unicode
func SafeStringLength(s string) int {
    return utf8.RuneCountInString(s)
}

// ContainsSpecialChars verifica se contém caracteres não-ASCII
func ContainsSpecialChars(s string) bool {
    for _, r := range s {
        if r > 127 {
            return true
        }
    }
    return false
}

// ValidateInputRunes valida caracteres de entrada
func ValidateInputRunes(input string) error {
    if !utf8.ValidString(input) {
        return fmt.Errorf("input contains invalid UTF-8 sequences")
    }
    
    for i, r := range input {
        if r == utf8.RuneError {
            return fmt.Errorf("invalid rune at position %d", i)
        }
        
        // Permitir caracteres printáveis e espaços em branco
        if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
            return fmt.Errorf("non-printable character at position %d", i)
        }
    }
    
    return nil
}

// TextInputModel with UTF-8 support
type TextInputModel struct {
    Value       string
    Placeholder string
    MaxRunes    int
}

func (m *TextInputModel) InsertRune(r rune) bool {
    if m.MaxRunes > 0 && utf8.RuneCountInString(m.Value) >= m.MaxRunes {
        return false
    }
    
    if unicode.IsPrint(r) || unicode.IsSpace(r) {
        m.Value += string(r)
        return true
    }
    
    return false
}
```

## Template Engine Go
```go
package core

import (
    "bytes"
    "strings"
    "text/template"
    "time"
    "path/filepath"
)

type TemplateEngine struct {
    funcMap template.FuncMap
}

func NewTemplateEngine() *TemplateEngine {
    return &TemplateEngine{
        funcMap: template.FuncMap{
            "title":     strings.Title,
            "lower":     strings.ToLower,
            "upper":     strings.ToUpper,
            "trim":      strings.TrimSpace,
            "wordCount": wordCount,
            "lineCount": lineCount,
            "now":       time.Now,
            "date":      formatDate,
        },
    }
}

func (e *TemplateEngine) Render(content string, data interface{}) (string, error) {
    tmpl, err := template.New("prompt").
        Funcs(e.funcMap).
        Parse(content)
    if err != nil {
        return "", err
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    
    return buf.String(), nil
}

func wordCount(s string) int {
    return len(strings.Fields(s))
}

func lineCount(s string) int {
    return len(strings.Split(s, "\n"))
}

func formatDate(t time.Time) string {
    return t.Format("2006-01-02")
}
```

# Otimizações de Performance

## Viewport Virtual com Bubble Tea
```go
type FileTreeModel struct {
    items        []FileItem
    viewport     viewport.Model
    cursor       int
    windowHeight int
}

func NewFileTreeModel() FileTreeModel {
    vp := viewport.New(80, 20)
    vp.Style = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder())
    
    return FileTreeModel{
        items:    []FileItem{},
        viewport: vp,
    }
}

func (m FileTreeModel) View() string {
    // Renderiza apenas itens visíveis
    visible := m.getVisibleItems()
    content := m.renderItems(visible)
    m.viewport.SetContent(content)
    return m.viewport.View()
}
```

## Cache Concorrente
```go
type FileScanner struct {
    cache     sync.Map
    mu        sync.RWMutex
    ignorer   *ignore.Ignorer
}

func (s *FileScanner) GetFileInfo(path string) (*FileInfo, error) {
    // Check cache first
    if cached, ok := s.cache.Load(path); ok {
        return cached.(*FileInfo), nil
    }
    
    // Scan file
    info, err := s.scanFile(path)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    s.cache.Store(path, info)
    return info, nil
}

func (s *FileScanner) ScanDirectory(root string) <-chan FileInfo {
    ch := make(chan FileInfo, 100)
    
    go func() {
        defer close(ch)
        
        // Worker pool for parallel scanning
        var wg sync.WaitGroup
        sem := make(chan struct{}, runtime.NumCPU())
        
        filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            if err != nil || s.ignorer.Ignore(path) {
                return nil
            }
            
            wg.Add(1)
            sem <- struct{}{} // Acquire semaphore
            
            go func(p string, i os.FileInfo) {
                defer wg.Done()
                defer func() { <-sem }() // Release semaphore
                
                if fileInfo := s.processFile(p, i); fileInfo != nil {
                    ch <- *fileInfo
                }
            }(path, info)
            
            return nil
        })
        
        wg.Wait()
    }()
    
    return ch
}
```

## Pipeline Pattern
```go
type PromptBuilder struct {
    scanner  *FileScanner
    template *TemplateEngine
}

func (b *PromptBuilder) Build(config BuildConfig) tea.Cmd {
    return func() tea.Msg {
        // Pipeline stages
        files := b.scanner.ScanDirectory(config.Root)
        filtered := b.filterFiles(files, config.Selected)
        content := b.readContents(filtered)
        prompt := b.assemblePrompt(content, config.Template)
        
        return PromptReadyMsg{Content: prompt}
    }
}

func (b *PromptBuilder) filterFiles(in <-chan FileInfo, selected map[string]bool) <-chan FileInfo {
    out := make(chan FileInfo)
    go func() {
        defer close(out)
        for file := range in {
            if selected[file.Path] {
                out <- file
            }
        }
    }()
    return out
}
```

# Features Adicionais Go

## CLI Interface com Cobra
```go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/shotgun/internal/app"
)

var rootCmd = &cobra.Command{
    Use:   "shotgun",
    Short: "Generate LLM prompts from templates",
    Run: func(cmd *cobra.Command, args []string) {
        // Start TUI application
        app := app.New()
        if err := app.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    },
}

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a prompt using the TUI wizard",
    Run: func(cmd *cobra.Command, args []string) {
        template, _ := cmd.Flags().GetString("template")
        task, _ := cmd.Flags().GetString("task")
        output, _ := cmd.Flags().GetString("output")
        quick, _ := cmd.Flags().GetBool("quick")
        
        app := app.New()
        if quick {
            app.LoadLastSession()
        }
        if template != "" {
            app.SetTemplate(template)
        }
        if task != "" {
            app.SetTask(task)
        }
        
        app.Run()
    },
}

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize .shotgunignore file",
    Run: func(cmd *cobra.Command, args []string) {
        CreateShotgunIgnore()
    },
}

func init() {
    rootCmd.AddCommand(generateCmd)
    rootCmd.AddCommand(initCmd)
    
    generateCmd.Flags().StringP("template", "t", "", "Template name to use")
    generateCmd.Flags().String("task", "", "Task description")
    generateCmd.Flags().StringP("output", "o", "", "Output file path")
    generateCmd.Flags().Bool("quick", false, "Use last configuration")
}

func Execute() error {
    return rootCmd.Execute()
}
```

## Configuração com Viper
```go
package core

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"
    
    "github.com/spf13/viper"
)

type Config struct {
    TemplatesDir    string `mapstructure:"templates_dir"`
    HistoryFile     string `mapstructure:"history_file"`
    MaxHistory      int    `mapstructure:"max_history"`
    AutoSave        bool   `mapstructure:"auto_save"`
    DefaultTemplate string `mapstructure:"default_template"`
}

func LoadConfig() (*Config, error) {
    viper.SetDefault("templates_dir", getConfigDir()+"/templates")
    viper.SetDefault("history_file", getConfigDir()+"/history.json")
    viper.SetDefault("max_history", 50)
    viper.SetDefault("auto_save", true)
    viper.SetDefault("default_template", "")
    
    viper.SetEnvPrefix("SHOTGUN")
    viper.AutomaticEnv()
    
    viper.SetConfigName("config")
    viper.SetConfigType("toml")
    viper.AddConfigPath(getConfigDir())
    
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

type Session struct {
    Timestamp     time.Time         `json:"timestamp"`
    SelectedFiles []string          `json:"selected_files"`
    Template      string            `json:"template"`
    Variables     map[string]string `json:"variables"`
    Task          string            `json:"task"`
    Rules         string            `json:"rules"`
}

type SessionManager struct {
    config *Config
}

func (sm *SessionManager) SaveSession(state *AppState) error {
    session := Session{
        Timestamp:     time.Now(),
        SelectedFiles: state.SelectedFiles,
        Template:      state.SelectedTemplate,
        Variables:     state.TemplateVariables,
        Task:          state.TaskContent,
        Rules:         state.RulesContent,
    }
    
    // Load history
    history, err := sm.LoadHistory()
    if err != nil {
        history = []Session{}
    }
    
    // Add new session and trim
    history = append([]Session{session}, history...)
    if len(history) > sm.config.MaxHistory {
        history = history[:sm.config.MaxHistory]
    }
    
    // Save to file
    data, err := json.MarshalIndent(history, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(sm.config.HistoryFile, data, 0644)
}

func getConfigDir() string {
    if runtime.GOOS == "windows" {
        return os.Getenv("APPDATA") + "/shotgun-cli"
    }
    home, _ := os.UserHomeDir()
    return home + "/.config/shotgun-cli"
}
```

## Histórico com Bubble Tea
```go
type HistoryModel struct {
    history  []Session
    selected int
    list     list.Model
}

func NewHistoryModel() HistoryModel {
    items := []list.Item{}
    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Recent Sessions"
    l.Styles.Title = titleStyle
    
    return HistoryModel{
        list: l,
    }
}

func (m HistoryModel) LoadHistory() tea.Cmd {
    return func() tea.Msg {
        sm := &SessionManager{config: loadConfig()}
        history, err := sm.LoadHistory()
        if err != nil {
            return errMsg{err}
        }
        return historyLoadedMsg{history}
    }
}

func (m HistoryModel) View() string {
    if len(m.history) == 0 {
        return "No recent sessions"
    }
    
    var s strings.Builder
    for i, session := range m.history[:min(10, len(m.history))] {
        cursor := " "
        if i == m.selected {
            cursor = ">"
        }
        
        s.WriteString(fmt.Sprintf(
            "%s %s  %s  %d files\n",
            cursor,
            session.Timestamp.Format("01/02"),
            session.Template,
            len(session.SelectedFiles),
        ))
    }
    
    s.WriteString("\nEnter restore · c copy · d delete")
    return s.String()
}
```

# Testes e Qualidade

## Estratégia de Testes Abrangente

### Níveis de Teste
1. **Unit Tests**: Cobertura mínima 90% para lógica de negócio
2. **Integration Tests**: Testes de fluxo completo TUI
3. **E2E Tests**: Simulação de uso real com expect/golden files
4. **Performance Tests**: Benchmarks para operações críticas
5. **Cross-platform Tests**: CI em Windows, Linux, macOS

### Ferramentas de Teste
- **testify**: Assertions e mocks
- **teatest**: Testing helper para Bubble Tea
- **golden**: Golden file testing para output
- **fuzzing**: Go native fuzzing para inputs
- **race detector**: Detecção de condições de corrida

## Testing Strategy Go
```go
// internal/core/scanner_test.go
package core

import (
    "testing"
    "os"
    "path/filepath"
)

func TestFileScannerRespectsGitignore(t *testing.T) {
    scanner := NewFileScanner()
    files := make([]FileInfo, 0)
    
    for file := range scanner.ScanDirectory("./testdata/project") {
        files = append(files, file)
    }
    
    // Verify node_modules is not included
    for _, f := range files {
        if filepath.Base(f.Path) == "node_modules" {
            t.Errorf("node_modules should be ignored")
        }
    }
}

func TestBinaryDetection(t *testing.T) {
    scanner := NewFileScanner()
    
    tests := []struct {
        file     string
        isBinary bool
    }{
        {"test.jpg", true},
        {"test.txt", false},
        {"test.go", false},
        {"test.exe", true},
    }
    
    for _, tt := range tests {
        result := scanner.IsBinary(tt.file)
        if result != tt.isBinary {
            t.Errorf("%s: expected %v, got %v", tt.file, tt.isBinary, result)
        }
    }
}

// internal/app/app_test.go
package app

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

func TestAppNavigation(t *testing.T) {
    model := InitialModel()
    
    // Test navigation
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
    model, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
    
    if model.state.CurrentScreen != FileTreeView {
        t.Errorf("Expected FileTreeView, got %v", model.state.CurrentScreen)
    }
}

func BenchmarkFileScanning(b *testing.B) {
    scanner := NewFileScanner()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        files := make([]FileInfo, 0)
        for file := range scanner.ScanDirectory(".") {
            files = append(files, file)
        }
    }
}
```

## CI/CD Pipeline
```yaml
# .github/workflows/ci.yml
name: CI/CD

on: [push, pull_request]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        go-version: ['1.21', '1.22']
        
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
        
    - name: Install dependencies
      run: go mod download
      
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
      
    - name: Run benchmarks
      run: go test -bench=. ./...
      
    - name: Run linting
      run: |
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        golangci-lint run
        
  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
        
    - name: Build binaries
      run: |
        GOOS=linux GOARCH=amd64 go build -o dist/shotgun-linux-amd64 ./cmd/shotgun
        GOOS=windows GOARCH=amd64 go build -o dist/shotgun-windows-amd64.exe ./cmd/shotgun
        GOOS=darwin GOARCH=amd64 go build -o dist/shotgun-darwin-amd64 ./cmd/shotgun
        GOOS=darwin GOARCH=arm64 go build -o dist/shotgun-darwin-arm64 ./cmd/shotgun
        
    - name: Create Release
      if: startsWith(github.ref, 'refs/tags/')
      uses: softprops/action-gh-release@v1
      with:
        files: dist/shotgun*
```

# Instalação e Distribuição

## Via go install
```bash
# Instalação direta do source
go install github.com/user/shotgun-cli/cmd/shotgun@latest
```

## Via Homebrew (macOS/Linux)
```bash
# Tap e instalação
brew tap user/shotgun
brew install shotgun-cli
```

## Binário Pré-compilado
```bash
# Linux
wget https://github.com/user/shotgun-cli/releases/latest/download/shotgun-linux-amd64
chmod +x shotgun-linux-amd64
sudo mv shotgun-linux-amd64 /usr/local/bin/shotgun

# macOS (Intel)
curl -L https://github.com/user/shotgun-cli/releases/latest/download/shotgun-darwin-amd64 -o shotgun
chmod +x shotgun
sudo mv shotgun /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/user/shotgun-cli/releases/latest/download/shotgun-darwin-arm64 -o shotgun
chmod +x shotgun
sudo mv shotgun /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/user/shotgun-cli/releases/latest/download/shotgun-windows-amd64.exe" -OutFile "shotgun.exe"
Move-Item shotgun.exe C:\Windows\System32\
```

## Build from Source
```bash
git clone https://github.com/user/shotgun-cli.git
cd shotgun-cli
go build -o shotgun ./cmd/shotgun
./shotgun
```

## Configuração go.mod
```go
module github.com/user/shotgun-cli

go 1.22

require (
    github.com/charmbracelet/bubbletea/v2 v2.0.0-beta.4
    github.com/charmbracelet/bubbles v0.21.0
    github.com/charmbracelet/lipgloss v1.0.0
    github.com/charmbracelet/glamour v0.6.0
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.2
    github.com/BurntSushi/toml v1.3.2
    github.com/bmatcuk/doublestar/v4 v4.6.1
    github.com/h2non/filetype v1.1.3
    github.com/go-git/go-git/v5 v5.11.0
    golang.org/x/text v0.18.0
)

require (
    // indirect dependencies
    github.com/atotto/clipboard v0.1.4 // indirect
    github.com/containerd/console v1.0.4 // indirect
    github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
    github.com/mattn/go-isatty v0.0.20 // indirect
    github.com/mattn/go-runewidth v0.0.15 // indirect
    github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
    github.com/muesli/cancelreader v0.2.2 // indirect
    github.com/muesli/reflow v0.3.0 // indirect
    github.com/muesli/termenv v0.15.2 // indirect
    golang.org/x/sync v0.5.0 // indirect
    golang.org/x/sys v0.15.0 // indirect
    golang.org/x/term v0.15.0 // indirect
    golang.org/x/text v0.14.0 // indirect
)
```

## Makefile
```makefile
.PHONY: build test clean install

BINARY_NAME=shotgun
BUILD_DIR=dist

build:
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/shotgun

build-all:
	@echo "Building for all platforms..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/shotgun
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/shotgun
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/shotgun
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/shotgun

test:
	@echo "Running tests..."
	go test -v -race ./...

bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

lint:
	@echo "Running linter..."
	golangci-lint run

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

install: build
	@echo "Installing..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
```

# Riscos e Mitigações

1) **Compatibilidade de Terminal (Windows/Linux)**
- Risco: Diferenças de rendering entre terminais, suporte a cores inconsistente, teclas especiais.
- Mitigação: Bubble Tea v2.0.0-beta.4 tem excelente suporte cross-platform com API v2 moderna; fallback gracioso para terminais limitados; testes em múltiplos terminais (PowerShell, CMD, bash, wezterm).
- **Melhorias Windows 2025**: Bubble Tea v2.0.0-beta.4 resolve problemas críticos:
  - Keyboard enhancements refinado com funções separadas (WithKeyReleases, WithUniformKeyLayout, RequestKeyDisambiguation)
  - Debug de panics melhorado com TEA_DEBUG=true
  - Melhor reconhecimento de combinações de teclas (Ctrl+C, Alt, etc.)
  - Suporte aprimorado para diferentes codepages do Windows
  - Detecção automática de capacidades do terminal (cores, Unicode)
  - Compatibilidade com ConPTY (Windows 10 1903+)

2) **Performance em Repositórios Grandes**
- Risco: UI pode ficar lenta durante scan de milhares de arquivos.
- Mitigação: **Concorrência nativa do Go** com goroutines e channels; worker pools; viewport virtual do Bubble Tea; pipeline pattern para processamento streaming; zero-allocation onde possível.

3) **Simplicidade do Bubble Tea v2**
- Risco: Menos features out-of-the-box comparado a frameworks maiores; quebras de compatibilidade na API v2.
- Mitigação: Bubble Tea v2 segue Elm Architecture (simples e previsível) com melhorias modernas; excelente documentação; bubbles v0.21.0 library com componentes prontos e melhorados (filepicker, textarea com horizontal scroll); Lip Gloss v1.0.0 para estilização poderosa; comunidade Charm muito ativa com releases beta frequentes; TEA_DEBUG para facilitar desenvolvimento.

4) **Gestão de Estado Imutável**
- Risco: State management com imutabilidade em Go.
- Mitigação: Model único centralizado seguindo Elm Architecture; updates funcionais puros; state transitions explícitas; serialização simples com structs Go.

5) **Validação de Templates**
- Risco: Templates TOML malformados, variáveis inconsistentes.
- Mitigação: Validação com struct tags Go; parsing type-safe com BurntSushi/toml; error handling idiomático Go; templates embarcados como fallback.

6) **Memory Usage com Arquivos Grandes**
- Risco: Carregar muitos arquivos grandes pode consumir muita RAM.
- Mitigação: Streaming file reading; chunk processing; content caching inteligente; memory monitoring; garbage collection proativo.

7) **Encoding Issues e Suporte a Caracteres Especiais**
- Risco: Arquivos com encodings diferentes, caracteres especiais como "ç", "á", "ô", etc.
- Mitigação: 
  - **Go UTF-8 nativo**: Go trabalha nativamente com UTF-8, todos os strings são UTF-8 por padrão
  - **Bubble Tea Unicode**: Framework suporta completamente caracteres Unicode/UTF-8 em input e output
  - **Validação de runes**: Usar `unicode/utf8` package para validação de caracteres válidos
  - **Input handling**: textinput e textarea do Bubbles processam automaticamente caracteres internacionais
  - **Terminal encoding**: Detecção automática de capacidades do terminal para Unicode
  - **Fallback gracioso**: Para terminais limitados, degradação suave mantendo funcionalidade

8) **Template Security**
- Risco: Injeção em templates text/template.
- Mitigação: text/template do Go é seguro por padrão (escape automático); sem execução de código arbitrário; validação de entrada; sanitização de output.

# Segurança e Proteção

## Segurança de Dados

### 1. Proteção de Informações Sensíveis
- **Detecção automática**: Scanner identifica arquivos com potenciais secrets (`.env`, `.pem`, `*_key`, etc.)
- **Exclusão padrão**: Arquivos sensíveis excluídos automaticamente do output
- **Aviso visual**: Alerta quando arquivo potencialmente sensível é selecionado
- **Sanitização**: Opção de mascarar valores em variáveis de ambiente

### 2. Validação de Entrada
- **Template validation**: TOML parsing com validação de estrutura
- **Path traversal protection**: Previne acesso a diretórios fora do escopo
- **Input sanitization**: Limpeza de caracteres de controle em inputs do usuário
- **Size limits**: Limites configuráveis para tamanho de arquivos

### 3. Segurança de Execução
- **No code execution**: Templates não executam código arbitrário
- **Memory safety**: Go garante memory safety e previne buffer overflows
- **Concurrent safety**: Uso correto de channels e sync primitives
- **Panic recovery**: Handlers para recuperação graciosa de panics

## Privacidade

### 1. Dados Locais
- **No telemetry**: Nenhum dado enviado para servidores externos
- **Local storage only**: Configurações e histórico apenas local
- **User control**: Usuário tem controle total sobre dados salvos
- **Clean uninstall**: Remoção completa sem deixar rastros

### 2. Configurações de Segurança
```yaml
# security.yaml
security:
  exclude_sensitive: true
  mask_env_vars: true
  max_file_size: 1MB
  allowed_extensions:
    - .go
    - .js
    - .py
  blocked_patterns:
    - "*_secret*"
    - "*.key"
    - "*.pem"
```

## Auditoria e Compliance

### 1. Logging de Segurança
- **Audit trail**: Log de arquivos acessados (opcional)
- **Error logging**: Registro de tentativas de acesso negadas
- **Session logging**: Histórico de sessões com timestamps

### 2. Best Practices
- **Principle of least privilege**: Acesso mínimo necessário
- **Defense in depth**: Múltiplas camadas de proteção
- **Fail secure**: Em caso de erro, falha de forma segura
- **Regular updates**: Dependências sempre atualizadas

# Cronograma de Desenvolvimento

## Fase 1: Core Infrastructure (Semanas 1-2)
- Estrutura básica do projeto
- Models com structs Go e validação
- File scanner concorrente com goroutines
- Template engine text/template
- Testes básicos

## Fase 2: TUI Foundation (Semanas 3-4)
- App Bubble Tea v2.0.0-beta.4 base
- Sistema de navegação
- Componentes básicos
- Lip Gloss styling
- State management com Elm Architecture

## Fase 3: Core Screens (Semanas 5-7)
- File tree widget
- Template selection
- Variable inputs
- Task/Rules editors
- Navigation flow

## Fase 4: Advanced Features (Semanas 8-9)
- Size estimation
- Progress indicators
- History system
- Configuration
- Error handling

## Fase 5: Polish & Testing (Semanas 10-11)
- UI refinements
- Performance optimization
- Comprehensive testing
- Documentation
- CLI interface

## Fase 6: Distribution (Semana 12)
- Go build para múltiplas plataformas
- CI/CD pipeline
- GitHub releases
- Release artifacts
- User documentation

# Métricas de Sucesso

## Performance
- Tempo de startup < 2s
- File scan de 1000+ arquivos < 5s
- UI responsiva (< 16ms frame time)
- Memory usage < 100MB para repos médios

## Usabilidade
- Zero travamentos durante operações I/O
- Navegação fluida entre telas
- Prompts gerados em < 30s
- Suporte a terminais 80x24

## Qualidade
- Test coverage > 90%
- Zero memory leaks
- Graceful error handling
- Cross-platform compatibility

## Funcionalidade
- Suporte a repos com 10k+ arquivos
- Templates customizados funcionais
- Histórico e sessions persistentes
- Integração com .gitignore/.shotgunignore

# Documentação do Usuário

## Instalação

### Binários Pré-compilados
```bash
# Linux/Mac
curl -sSL https://github.com/user/shotgun-cli/releases/latest/download/shotgun-cli-$(uname -s)-$(uname -m) -o shotgun
chmod +x shotgun

# Windows
# Baixar shotgun-cli-windows-amd64.exe do GitHub Releases
```

### Via Go Install
```bash
go install github.com/user/shotgun-cli@latest
```

## Uso Básico

### Modo Interativo (TUI)
```bash
shotgun          # Abre interface interativa
shotgun -t bug   # Abre com template específico
```

### Modo CLI Direto
```bash
shotgun generate --template bug --task "Fix login issue" --output prompt.md
shotgun list-templates  # Lista templates disponíveis
```

## Configuração

### Arquivo de Configuração
- Linux/Mac: `~/.config/shotgun-cli/config.yaml`
- Windows: `%APPDATA%\shotgun-cli\config.yaml`

```yaml
# config.yaml
default_template: "dev"
max_file_size: 1000000
exclude_patterns:
  - "*.log"
  - "node_modules/**"
color_scheme: "monochrome"
```

### Templates Customizados
- Linux/Mac: `~/.config/shotgun-cli/templates/`
- Windows: `%APPDATA%\shotgun-cli\templates\`

## Atalhos de Teclado

### Navegação Global
- `F1` - Ajuda contextual
- `F2` - Tela anterior
- `F3` - Próxima tela
- `F4-F10` - Acesso direto às telas
- `Ctrl+C` - Sair
- `Ctrl+S` - Salvar sessão

### Árvore de Arquivos
- `↑/↓` - Navegar
- `Espaço` - Marcar/desmarcar
- `Enter` - Expandir/recolher
- `Ctrl+A` - Marcar todos
- `Ctrl+I` - Inverter seleção

### Editores de Texto
- `Enter` - Nova linha
- `Ctrl+Enter` - Alternar modo edição/navegação
- `Tab` - Indentação
- `Shift+Tab` - Remove indentação

## Troubleshooting

### Problemas Comuns
1. **Terminal não suporta cores**: Use `TERM=xterm-256color shotgun`
2. **Caracteres Unicode quebrados**: Verifique encoding UTF-8 do terminal
3. **Performance lenta**: Ajuste `max_file_size` na configuração
4. **Bubble Tea panic**: Use `TEA_DEBUG=true shotgun` para debug

# Plano de Contingência para API Beta

## Riscos da API v2.0.0-beta.4

### 1. Breaking Changes
- **Risco**: API pode mudar entre versões beta
- **Mitigação**: 
  - Pin exato da versão no go.mod
  - Testes de regressão abrangentes
  - Camada de abstração sobre APIs críticas
  - Monitoramento de releases do Bubble Tea

### 2. Bugs não Descobertos
- **Risco**: Comportamentos inesperados em edge cases
- **Mitigação**:
  - Testes em múltiplos ambientes
  - Fallback handlers para panics
  - Log detalhado de erros
  - Comunicação ativa com comunidade Charm

### 3. Features Incompletas
- **Risco**: Funcionalidades podem não estar totalmente implementadas
- **Mitigação**:
  - Implementação defensiva
  - Feature flags para recursos experimentais
  - Graceful degradation quando possível

## Estratégia de Migração

1. **Monitoramento**: Check semanal de novos releases
2. **Testes**: Suite automatizada para cada atualização
3. **Rollback**: Manter versão estável anterior como fallback
4. **Documentação**: Changelog detalhado de mudanças

Este plano estabelece uma base sólida para criar uma versão Go moderna e eficiente do shotgun-cli, aproveitando as melhores práticas e ferramentas do ecossistema Go para TUI development com Bubble Tea v2.