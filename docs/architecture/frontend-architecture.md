# Frontend Architecture

## Component Architecture

### Component Organization
```
internal/screens/
├── filetree/
│   ├── model.go      # FileTree state
│   ├── update.go     # Message handling
│   ├── view.go       # Rendering logic
│   └── keys.go       # Keybindings
├── template/
│   ├── model.go
│   ├── update.go
│   ├── view.go
│   └── list.go       # List component
├── input/
│   ├── task.go       # Task input screen
│   ├── rules.go      # Rules input screen
│   └── editor.go     # Shared editor
└── confirm/
    ├── model.go
    ├── view.go
    └── progress.go   # Progress bar
```

### Component Template
```go
// Standard Bubble Tea component structure
type FileTreeModel struct {
    items    []FileItem
    cursor   int
    selected map[string]bool
    viewport viewport.Model
}

func NewFileTreeModel() FileTreeModel {
    return FileTreeModel{
        items:    []FileItem{},
        selected: make(map[string]bool),
        viewport: viewport.New(80, 20),
    }
}

func (m FileTreeModel) Update(msg tea.Msg) (FileTreeModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            m.cursor--
        case "down", "j":
            m.cursor++
        case " ":
            m.toggleSelection()
        }
    }
    return m, nil
}

func (m FileTreeModel) View() string {
    return m.viewport.View()
}
```

## State Management Architecture

### State Structure
```go
type AppState struct {
    CurrentScreen Screen
    FileTree      FileTreeModel
    Template      TemplateModel
    TaskInput     InputModel
    RulesInput    InputModel
    Confirmation  ConfirmModel
    
    // Shared state
    SelectedFiles    []string
    SelectedTemplate *Template
    TaskContent      string
    RulesContent     string
    
    // UI state
    WindowSize tea.WindowSizeMsg
    Error      error
}
```

### State Management Patterns
- Immutable state updates through Bubble Tea
- Message-based communication between components
- Command pattern for async operations
- Single source of truth in AppState

## Routing Architecture

### Route Organization
```
Screen Flow:
1. FileTree    (mandatory)
2. Template    (mandatory)  
3. TaskInput   (mandatory)
4. RulesInput  (optional - F4 to skip)
5. Confirm     (mandatory)
```

### Protected Route Pattern
```go
// Screen transition validation
func (m Model) canAdvance() bool {
    switch m.CurrentScreen {
    case FileTreeScreen:
        return len(m.SelectedFiles) > 0
    case TemplateScreen:
        return m.SelectedTemplate != nil
    case TaskScreen:
        return m.TaskContent != ""
    case RulesScreen:
        return true // Optional
    default:
        return false
    }
}

func (m Model) advance() (Model, tea.Cmd) {
    if !m.canAdvance() {
        m.Error = ErrValidationFailed
        return m, nil
    }
    m.CurrentScreen++
    return m, m.initScreen()
}
```

## Frontend Services Layer

### API Client Setup
N/A - No API client needed for local file operations

### Service Example
```go
// File service for frontend screens
type FileService struct {
    scanner *FileScanner
    cache   *sync.Map
}

func (s *FileService) GetFileTree(root string) tea.Cmd {
    return func() tea.Msg {
        files := s.scanner.ScanDirectory(root)
        return FileTreeLoadedMsg{Files: files}
    }
}

func (s *FileService) ReadFileContent(path string) tea.Cmd {
    return func() tea.Msg {
        content, err := os.ReadFile(path)
        if err != nil {
            return ErrorMsg{err}
        }
        return FileContentMsg{Path: path, Content: string(content)}
    }
}
```
