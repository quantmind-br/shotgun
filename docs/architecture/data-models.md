# Data Models

### ApplicationState

**Purpose:** Central state container managing the complete wizard workflow and user session

**Key Attributes:**
- CurrentScreen: ScreenType - Active screen in the 5-screen wizard
- FileSelection: FileSelectionState - Complete file tree with selection states  
- TemplateChoice: SelectedTemplate - Chosen template with metadata
- UserInputs: UserInputData - Task description and optional rules
- GenerationConfig: GenerationSettings - Output configuration and preferences

#### TypeScript Interface
```go
type ApplicationState struct {
    CurrentScreen     ScreenType           `json:"current_screen"`
    FileSelection     *FileSelectionState  `json:"file_selection"`
    TemplateChoice    *SelectedTemplate    `json:"template_choice"`
    UserInputs        *UserInputData       `json:"user_inputs"`
    GenerationConfig  *GenerationSettings  `json:"generation_config"`
    SessionID         string               `json:"session_id"`
    CreatedAt         time.Time            `json:"created_at"`
    ModifiedAt        time.Time            `json:"modified_at"`
}
```

#### Relationships
- Contains FileSelectionState for file tree management
- References SelectedTemplate for template processing
- Includes UserInputData for text input persistence

### FileNode

**Purpose:** Represents individual files and directories in the hierarchical file tree with selection capabilities

**Key Attributes:**
- Path: string - Relative file path from project root
- Name: string - File or directory name for display
- IsDirectory: bool - Directory vs file classification
- IsSelected: bool - User selection state for prompt inclusion
- IsIgnored: bool - Excluded by .gitignore/.shotgunignore rules
- IsBinary: bool - Binary file detection result
- Size: int64 - File size in bytes for estimation calculations
- Children: []FileNode - Nested directory contents

#### TypeScript Interface
```go
type FileNode struct {
    Path        string      `json:"path"`
    Name        string      `json:"name"`
    IsDirectory bool        `json:"is_directory"`
    IsSelected  bool        `json:"is_selected"`
    IsIgnored   bool        `json:"is_ignored"`
    IsBinary    bool        `json:"is_binary"`
    Size        int64       `json:"size"`
    ModTime     time.Time   `json:"mod_time"`
    Children    []*FileNode `json:"children,omitempty"`
    Parent      *FileNode   `json:"-"` // Avoid circular references in JSON
}
```

#### Relationships
- Self-referential tree structure (Parent/Children)
- Aggregated in FileSelectionState for complete tree management

### Template

**Purpose:** Template definition containing metadata and processing instructions for prompt generation

**Key Attributes:**
- ID: string - Unique template identifier
- Name: string - Human-readable template name
- Version: string - Template version for compatibility
- Description: string - Template purpose and usage guidance
- Author: string - Template creator information
- Content: string - Raw template content with placeholders
- Variables: []string - Required template variables
- IsBuiltIn: bool - Embedded vs user-created template distinction

#### TypeScript Interface
```go
type Template struct {
    ID          string            `toml:"id" json:"id"`
    Name        string            `toml:"name" json:"name"`
    Version     string            `toml:"version" json:"version"`
    Description string            `toml:"description" json:"description"`
    Author      string            `toml:"author" json:"author"`
    Tags        []string          `toml:"tags" json:"tags"`
    Content     string            `toml:"content" json:"content"`
    Variables   []string          `toml:"variables" json:"variables"`
    IsBuiltIn   bool              `json:"is_builtin"`
    FilePath    string            `json:"file_path"`
    Metadata    map[string]string `toml:"metadata" json:"metadata"`
}
```

#### Relationships
- Referenced by SelectedTemplate in ApplicationState
- Processed by TemplateEngine for prompt generation

### PromptGeneration

**Purpose:** Contains all data and configuration needed for final prompt generation and output

**Key Attributes:**
- Template: Template - Selected template with processing instructions
- SelectedFiles: []FileNode - Files chosen for inclusion with content
- TaskDescription: string - User-provided task description
- AdditionalRules: string - Optional user rules and constraints
- OutputPath: string - Generated filename with timestamp
- EstimatedSize: int64 - Size estimation before generation
- GeneratedAt: time.Time - Generation timestamp for tracking

#### TypeScript Interface
```go
type PromptGeneration struct {
    Template        *Template   `json:"template"`
    SelectedFiles   []*FileNode `json:"selected_files"`
    TaskDescription string      `json:"task_description"`
    AdditionalRules string      `json:"additional_rules"`
    OutputPath      string      `json:"output_path"`
    EstimatedSize   int64       `json:"estimated_size"`
    ActualSize      int64       `json:"actual_size"`
    GeneratedAt     time.Time   `json:"generated_at"`
    Success         bool        `json:"success"`
    ErrorMessage    string      `json:"error_message,omitempty"`
}
```

#### Relationships
- Aggregates Template and FileNode data for processing
- Created from ApplicationState during final generation step
