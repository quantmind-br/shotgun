# API Specification

Since **shotgun-cli-v3** is a standalone TUI application rather than a web service, it doesn't expose traditional REST/GraphQL APIs. However, the application does have important internal interfaces and potential extension points that function as "APIs" within the Go ecosystem:

### Internal Package APIs

The application exposes clean Go package interfaces that serve as internal APIs for modularity and testing:

#### Template Processing API
```go
// Package: internal/template
type TemplateEngine interface {
    LoadTemplates(dirs []string) error
    GetAvailableTemplates() ([]*Template, error)
    ProcessTemplate(template *Template, data *TemplateData) (string, error)
    ValidateTemplate(template *Template) error
}

type TemplateData struct {
    Task           string
    Rules          string
    FileStructure  string
    FileContents   map[string]string
    ProjectName    string
    GeneratedAt    time.Time
}
```

#### File Processing API
```go
// Package: internal/scanner
type FileScanner interface {
    ScanDirectory(rootPath string) (*FileNode, error)
    ApplyIgnoreRules(root *FileNode, rules []string) error
    GetSelectedFiles(root *FileNode) ([]*FileNode, error)
    EstimateSize(files []*FileNode) int64
}

type IgnoreProcessor interface {
    LoadRules(gitignorePath, shotgunignorePath string) error
    ShouldIgnore(path string, isDir bool) bool
    GetActiveRules() []string
}
```

#### State Management API
```go
// Package: internal/state
type StateManager interface {
    SaveSession(state *ApplicationState) error
    LoadSession(sessionID string) (*ApplicationState, error)
    ListSessions() ([]*SessionInfo, error)
    DeleteSession(sessionID string) error
}

type SessionInfo struct {
    ID          string    `json:"id"`
    CreatedAt   time.Time `json:"created_at"`
    ProjectPath string    `json:"project_path"`
    Summary     string    `json:"summary"`
}
```

### CLI Command Interface

The application exposes a command-line interface that serves as its primary external API:

```bash
# Core wizard mode (default)
shotgun

# Direct prompt generation (bypass UI)
shotgun generate --template analyze_bug --task "Debug memory leak" --output custom_prompt.md

# Template management
shotgun template list
shotgun template show analyze_bug
shotgun template create --from-file ./my-template.toml

# Session management
shotgun session list
shotgun session restore <session-id>
shotgun session clean

# Project initialization
shotgun init --create-config --create-ignore

# Configuration
shotgun config show
shotgun config set templates.directory ~/.shotgun/custom-templates
```

### Plugin Extension Points

Future extension capabilities through Go plugin interfaces:

```go
// Package: internal/plugins
type TemplateProvider interface {
    Name() string
    LoadTemplates() ([]*Template, error)
    ValidateTemplate(template *Template) error
}

type OutputProcessor interface {
    Name() string
    ProcessOutput(content string, config map[string]interface{}) (string, error)
    SupportedFormats() []string
}

type FileFilter interface {
    Name() string
    ShouldInclude(file *FileNode, context *FilterContext) bool
    Priority() int
}
```
