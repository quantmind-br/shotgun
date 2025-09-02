# Backend Architecture

## Service Architecture

### Function Organization
```
internal/core/
├── scanner/
│   ├── scanner.go       # Main scanner logic
│   ├── ignore.go        # Gitignore handling
│   ├── binary.go        # Binary detection
│   └── concurrent.go    # Worker pool
├── template/
│   ├── engine.go        # Template processing
│   ├── discovery.go     # Template finding
│   ├── validation.go    # Variable validation
│   └── builtin.go       # Embedded templates
├── builder/
│   ├── builder.go       # Prompt assembly
│   ├── pipeline.go      # Processing pipeline
│   └── writer.go        # Output generation
└── config/
    ├── config.go        # Configuration
    ├── paths.go         # Path resolution
    └── defaults.go      # Default values
```

### Function Template
```go
// Service function template
package scanner

type Scanner struct {
    ignorer  *ignore.Ignorer
    detector *filetype.Detector
    workers  int
}

func New(opts ...Option) *Scanner {
    s := &Scanner{
        workers: runtime.NumCPU(),
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

func (s *Scanner) ScanDirectory(root string) (<-chan FileInfo, error) {
    if err := s.validatePath(root); err != nil {
        return nil, err
    }
    
    ch := make(chan FileInfo, 100)
    go s.scan(root, ch)
    return ch, nil
}
```

## Database Architecture

### Schema Design
```sql
-- Future feature: Session history
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    timestamp INTEGER NOT NULL,
    selected_files TEXT NOT NULL, -- JSON array
    template_id TEXT NOT NULL,
    task_content TEXT,
    rules_content TEXT,
    output_path TEXT,
    output_size INTEGER
);

CREATE INDEX idx_sessions_timestamp ON sessions(timestamp DESC);
```

### Data Access Layer
```go
// Repository pattern for future persistence
type SessionRepository interface {
    Save(session *Session) error
    Load(id string) (*Session, error)
    List(limit int) ([]*Session, error)
    Delete(id string) error
}

type FileSystemRepository struct {
    basePath string
}

func (r *FileSystemRepository) Save(session *Session) error {
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }
    
    path := filepath.Join(r.basePath, session.ID+".json")
    return os.WriteFile(path, data, 0644)
}
```

## Authentication and Authorization

N/A - Shotgun CLI is a local tool with no authentication requirements. File system permissions are inherited from the operating system.
