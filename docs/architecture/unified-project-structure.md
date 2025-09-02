# Unified Project Structure

```plaintext
shotgun-cli/
├── .github/                    # CI/CD workflows
│   └── workflows/
│       ├── ci.yaml            # Test and build
│       ├── release.yaml       # Create releases
│       └── codeql.yaml        # Security scanning
├── cmd/
│   └── shotgun/
│       └── main.go            # Entry point
├── internal/                   # Private packages
│   ├── app/
│   │   ├── app.go            # Main app controller
│   │   ├── model.go          # App state
│   │   ├── update.go         # Message handling
│   │   ├── view.go           # View routing
│   │   └── keys.go           # Global keybindings
│   ├── screens/               # TUI screens
│   │   ├── filetree/
│   │   ├── template/
│   │   ├── input/
│   │   └── confirm/
│   ├── components/            # Reusable UI components
│   │   ├── tree/
│   │   ├── list/
│   │   ├── editor/
│   │   ├── progress/
│   │   └── statusbar/
│   ├── core/                  # Business logic
│   │   ├── scanner/
│   │   ├── template/
│   │   ├── builder/
│   │   └── config/
│   ├── models/                # Data structures
│   │   ├── file.go
│   │   ├── template.go
│   │   ├── variable.go
│   │   └── state.go
│   ├── styles/                # Lip Gloss styles
│   │   ├── theme.go
│   │   ├── colors.go
│   │   └── components.go
│   └── utils/                 # Utilities
│       ├── files.go
│       ├── validation.go
│       └── unicode.go
├── templates/                  # Embedded templates
│   ├── analyze_bug.toml
│   ├── make_diff.toml
│   ├── make_plan.toml
│   └── project_manager.toml
├── testdata/                   # Test fixtures
│   └── sample_project/
├── docs/                       # Documentation
│   ├── prd.md
│   ├── front-end-spec.md
│   └── architecture.md
├── .gitignore
├── .golangci.yml              # Linter config
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── CHANGELOG.md
```
