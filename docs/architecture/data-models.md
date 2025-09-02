# Data Models

## FileItem

**Purpose:** Represents a file or directory in the file tree with selection state

**Key Attributes:**
- Path: string - Absolute file path
- Name: string - Display name (basename)
- IsDir: bool - Directory flag
- IsSelected: bool - Selection state
- IsBinary: bool - Binary file indicator
- IsIgnored: bool - Gitignore status
- Size: int64 - File size in bytes
- Children: []FileItem - Child items for directories

### TypeScript Interface
```typescript
interface FileItem {
  path: string;
  name: string;
  isDir: boolean;
  isSelected: boolean;
  isBinary: boolean;
  isIgnored: boolean;
  size: number;
  children: FileItem[];
}
```

### Relationships
- Hierarchical self-reference for tree structure
- Many-to-one relationship with parent directory

## Template

**Purpose:** Defines a prompt template with metadata and variables

**Key Attributes:**
- ID: string - Unique template identifier
- Name: string - Display name
- Version: string - Semantic version
- Description: string - Template purpose
- Author: string - Template creator
- Tags: []string - Categorization tags
- Variables: map[string]Variable - Template variables
- Content: string - Template content with placeholders

### TypeScript Interface
```typescript
interface Template {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  tags: string[];
  variables: Record<string, Variable>;
  content: string;
}
```

### Relationships
- One-to-many relationship with Variables
- Referenced by AppState during selection

## Variable

**Purpose:** Defines a template variable with type and validation rules

**Key Attributes:**
- Name: string - Variable name
- Type: string - Variable type (text, multiline, auto, choice, boolean, number)
- Required: bool - Required flag
- Default: string - Default value
- Placeholder: string - UI placeholder text
- MinLength: int - Minimum length for text
- MaxLength: int - Maximum length for text
- Options: []string - Choices for choice type

### TypeScript Interface
```typescript
interface Variable {
  name: string;
  type: 'text' | 'multiline' | 'auto' | 'choice' | 'boolean' | 'number';
  required: boolean;
  default?: string;
  placeholder?: string;
  minLength?: number;
  maxLength?: number;
  options?: string[];
}
```

### Relationships
- Many-to-one relationship with Template
- Referenced during prompt generation

## AppState

**Purpose:** Central application state managed by Bubble Tea

**Key Attributes:**
- CurrentScreen: ScreenType - Active screen identifier
- FileTree: []FileItem - File tree data
- SelectedFiles: []string - Selected file paths
- SelectedTemplate: *Template - Chosen template
- TaskContent: string - User-entered task
- RulesContent: string - User-entered rules
- OutputSize: int64 - Estimated output size
- Error: error - Current error state

### TypeScript Interface
```typescript
interface AppState {
  currentScreen: 'fileTree' | 'template' | 'task' | 'rules' | 'confirm';
  fileTree: FileItem[];
  selectedFiles: string[];
  selectedTemplate: Template | null;
  taskContent: string;
  rulesContent: string;
  outputSize: number;
  error: Error | null;
}
```

### Relationships
- Aggregates all other models
- Single source of truth for UI state
