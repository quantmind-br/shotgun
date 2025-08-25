# Frontend Architecture

In the context of **shotgun-cli-v3**, the "frontend" refers to the Terminal User Interface (TUI) layer built with Bubble Tea. This section defines the TUI-specific architecture patterns and components:

### Component Architecture

The TUI follows a component-based architecture using Bubble Tea's `tea.Model` interface for each screen and reusable UI element:

#### Component Organization

```plaintext
internal/ui/
├── app/                    # Main application shell
│   ├── app.go             # Root Bubble Tea model
│   ├── navigation.go      # Screen navigation logic  
│   └── keybindings.go     # Global key handlers
├── screens/               # 5-screen wizard components
│   ├── filetree/          # File selection screen [1/5]
│   ├── templates/         # Template selection screen [2/5]
│   ├── taskinput/         # Task description screen [3/5]
│   ├── rulesinput/        # Rules input screen [4/5]
│   └── confirmation/      # Final confirmation screen [5/5]
├── components/            # Reusable UI components
│   ├── filetree.go       # Hierarchical file tree widget
│   ├── templatelist.go   # Template selection list
│   ├── texteditor.go     # Advanced multiline editor
│   ├── progressbar.go    # Progress indicators
│   └── statusbar.go      # Status and help display
├── styles/               # Lip Gloss styling
│   ├── theme.go          # Monochrome theme definition
│   ├── colors.go         # Color palette constants
│   └── layout.go         # Layout and spacing utilities
└── messages/             # Bubble Tea messages
    ├── navigation.go     # Screen change messages
    ├── file.go          # File operation messages
    └── generation.go    # Prompt generation messages
```

#### Component Template

```go
// Base component interface for all UI elements
type Component interface {
    tea.Model
    SetSize(width, height int) Component
    Focus() tea.Cmd
    Blur() tea.Cmd
    IsFocused() bool
}

// Example: File Tree Component
type FileTreeComponent struct {
    width, height int
    focused       bool
    root         *models.FileNode
    cursor       int
    viewport     viewport.Model
    styles       *styles.Theme
}

func (c FileTreeComponent) Init() tea.Cmd {
    return nil
}

func (c FileTreeComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return c.handleKeyPress(msg)
    case tea.WindowSizeMsg:
        return c.handleResize(msg)
    case messages.FileSelectionChanged:
        return c.handleSelectionChange(msg)
    }
    return c, nil
}

func (c FileTreeComponent) View() string {
    return c.styles.FileTree.Render(c.renderTree())
}
```

### State Management Architecture

The application uses a centralized state pattern with Bubble Tea's message passing for updates:

#### State Structure

```go
// Central application state
type ApplicationModel struct {
    // Navigation state
    currentScreen ScreenType
    screenStack   []ScreenType
    
    // Screen models
    screens map[ScreenType]tea.Model
    
    // Shared application data
    state *models.ApplicationState
    
    // UI state
    windowSize tea.WindowSizeMsg
    theme      *styles.Theme
    
    // Services
    fileScanner    scanner.FileScanner
    templateEngine template.Engine
    sessionManager session.Manager
}

// Screen-specific state management
type FileTreeModel struct {
    // Component state
    fileTree     *components.FileTreeComponent
    statusBar    *components.StatusBarComponent
    
    // Data state
    rootNode     *models.FileNode
    selectedFiles []*models.FileNode
    
    // UI state
    loading      bool
    error        string
}
```

#### State Management Patterns

- **Centralized Store:** ApplicationModel holds all shared state
- **Component Isolation:** Each screen manages its own UI state
- **Message Passing:** All state changes flow through Bubble Tea messages
- **Immutable Updates:** State changes create new state instances
- **Persistence:** State automatically saved to session files

### Routing Architecture

The TUI uses a screen-based routing system with validation and state preservation:

#### Route Organization

```plaintext
Screen Flow:
FileTree [1/5] → TemplateSelection [2/5] → TaskInput [3/5] → RulesInput [4/5] → Confirmation [5/5]
     ↕              ↕                        ↕                ↕                    ↕
Navigation: F2 (back), F3 (forward), ESC (exit), F1 (help)

Validation Gates:
- FileTree → TemplateSelection: At least one file selected
- TemplateSelection → TaskInput: Valid template chosen  
- TaskInput → RulesInput: Non-empty task description
- RulesInput → Confirmation: No validation (rules optional)
- Confirmation → Generation: Final user confirmation
```

#### Protected Route Pattern

```go
type ScreenRouter struct {
    current     ScreenType
    screens     map[ScreenType]tea.Model
    validators  map[ScreenType]ValidationFunc
}

type ValidationFunc func(state *models.ApplicationState) error

func (r *ScreenRouter) Navigate(direction NavigationDirection, state *models.ApplicationState) error {
    target := r.getTargetScreen(direction)
    
    // Validate transition
    if validator, exists := r.validators[target]; exists {
        if err := validator(state); err != nil {
            return fmt.Errorf("navigation blocked: %w", err)
        }
    }
    
    // Perform navigation
    r.current = target
    return nil
}

// Example validator
func validateFileSelection(state *models.ApplicationState) error {
    if len(state.FileSelection.SelectedFiles()) == 0 {
        return errors.New("at least one file must be selected")
    }
    return nil
}
```

### Frontend Services Layer

The TUI layer communicates with the processing engine through clean service interfaces:

#### API Client Setup

```go
// Service layer for TUI to backend communication
type TUIServices struct {
    fileScanner    scanner.FileScanner
    templateEngine template.Engine  
    outputGenerator output.Generator
    sessionManager session.Manager
    configManager  config.Manager
}

func NewTUIServices() *TUIServices {
    return &TUIServices{
        fileScanner:    scanner.NewConcurrentScanner(),
        templateEngine: template.NewEngine(),
        outputGenerator: output.NewGenerator(),
        sessionManager: session.NewFileManager(),
        configManager:  config.NewViperManager(),
    }
}
```

#### Service Example

```go
// File scanning service integration
type FileService struct {
    scanner scanner.FileScanner
    ignore  ignore.Processor
}

func (fs *FileService) ScanProjectFiles(rootPath string) tea.Cmd {
    return tea.Sequence(
        fs.showScanProgress(),
        tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
            result, err := fs.scanner.ScanDirectory(rootPath)
            if err != nil {
                return messages.FileScanError{Err: err}
            }
            return messages.FileScanComplete{Root: result}
        }),
    )
}

// Template processing service
func (ts *TemplateService) GeneratePrompt(data *models.PromptGeneration) tea.Cmd {
    return func() tea.Msg {
        result, err := ts.engine.ProcessTemplate(data.Template, data)
        if err != nil {
            return messages.GenerationError{Err: err}
        }
        return messages.GenerationComplete{
            Content:    result,
            OutputPath: data.OutputPath,
            Size:       int64(len(result)),
        }
    }
}
```
