# Unified Project Structure

The following structure accommodates the Go TUI application architecture with clear separation of concerns and standard Go project conventions:

```plaintext
shotgun-cli-v3/
├── .github/                    # CI/CD workflows
│   └── workflows/
│       ├── ci.yml              # Continuous integration
│       ├── release.yml         # Automated releases
│       └── codeql-analysis.yml # Security scanning
├── cmd/                        # Application entry points
│   └── shotgun/
│       └── main.go             # Main application entry
├── internal/                   # Private application packages
│   ├── ui/                     # TUI "Frontend" Layer
│   │   ├── app/                # Main application shell
│   │   │   ├── app.go          # Root Bubble Tea model
│   │   │   ├── navigation.go   # Screen navigation logic
│   │   │   └── keybindings.go  # Global key handlers
│   │   ├── screens/            # 5-screen wizard components
│   │   │   ├── filetree/       # File selection screen [1/5]
│   │   │   ├── templates/      # Template selection screen [2/5]
│   │   │   ├── taskinput/      # Task description screen [3/5]
│   │   │   ├── rulesinput/     # Rules input screen [4/5]
│   │   │   └── confirmation/   # Final confirmation screen [5/5]
│   │   ├── components/         # Reusable UI components
│   │   │   ├── filetree.go     # Hierarchical file tree widget
│   │   │   ├── templatelist.go # Template selection list
│   │   │   ├── texteditor.go   # Advanced multiline editor
│   │   │   ├── progressbar.go  # Progress indicators
│   │   │   └── statusbar.go    # Status and help display
│   │   ├── styles/             # Lip Gloss styling
│   │   │   ├── theme.go        # Monochrome theme definition
│   │   │   ├── colors.go       # Color palette constants
│   │   │   └── layout.go       # Layout and spacing utilities
│   │   └── messages/           # Bubble Tea messages
│   │       ├── navigation.go   # Screen change messages
│   │       ├── file.go         # File operation messages
│   │       └── generation.go   # Prompt generation messages
│   ├── core/                   # "Backend" Processing Engine
│   │   ├── scanner/            # File discovery and processing
│   │   │   ├── scanner.go      # Main file scanning logic
│   │   │   ├── concurrent.go   # Worker pool implementation
│   │   │   └── types.go        # File scanning data types
│   │   ├── template/           # Template processing engine
│   │   │   ├── engine.go       # Template compilation and execution
│   │   │   ├── discovery.go    # Template discovery and loading
│   │   │   ├── validation.go   # Template syntax validation
│   │   │   └── functions.go    # Custom template functions
│   │   ├── ignore/             # Ignore rules processing
│   │   │   ├── processor.go    # .gitignore/.shotgunignore logic
│   │   │   ├── patterns.go     # Glob pattern matching
│   │   │   └── precedence.go   # Rule precedence handling
│   │   └── output/             # Prompt generation and file writing
│   │       ├── generator.go    # Final prompt assembly
│   │       ├── formatter.go    # Markdown formatting
│   │       └── writer.go       # Atomic file operations
│   ├── services/               # Application services layer
│   │   ├── coordinator.go      # Service orchestration
│   │   ├── session.go          # Session persistence
│   │   └── config.go           # Configuration management
│   ├── models/                 # Shared data structures
│   │   ├── application.go      # ApplicationState and core types
│   │   ├── files.go            # FileNode and file-related types
│   │   ├── templates.go        # Template and processing types
│   │   └── config.go           # Configuration structures
│   ├── infrastructure/         # External dependencies abstraction
│   │   ├── filesystem/         # File system operations
│   │   ├── storage/            # Session storage implementations
│   │   └── platform/           # Cross-platform utilities
│   └── testutil/               # Testing utilities and helpers
│       ├── fixtures/           # Test data and fixtures
│       ├── mocks/              # Generated mocks
│       └── helpers.go          # Test helper functions
├── templates/                  # Embedded template assets
│   ├── embedded/               # Built-in templates
│   │   ├── analyze_bug.tmpl    # Bug analysis template
│   │   ├── make_diff.tmpl      # Diff creation template
│   │   ├── make_plan.tmpl      # Planning template
│   │   └── project_manager.tmpl # Project management template
│   └── examples/               # Example custom templates
│       ├── code_review.tmpl    # Code review template example
│       └── documentation.tmpl  # Documentation template example
├── configs/                    # Configuration files and examples
│   ├── default.toml            # Default configuration
│   ├── example.shotgunignore   # Example ignore rules
│   └── template-schema.json    # Template validation schema
├── scripts/                    # Build and development scripts
│   ├── build.sh               # Cross-platform build script
│   ├── test.sh                # Testing script
│   ├── lint.sh                # Linting and formatting
│   └── release.sh             # Release preparation
├── docs/                       # Documentation
│   ├── prd.md                 # Product requirements document
│   ├── architecture.md        # This architecture document
│   ├── api/                   # Internal API documentation
│   └── user-guide.md          # User documentation
├── test/                      # End-to-end tests
│   ├── integration/           # Integration tests
│   ├── e2e/                   # Full workflow tests
│   └── fixtures/              # Test project fixtures
├── .env.example               # Environment variable template
├── .gitignore                 # Git ignore rules
├── .shotgunignore             # Example shotgun ignore rules
├── .golangci.yml              # Linter configuration
├── go.mod                     # Go module definition
├── go.sum                     # Go dependency checksums
├── Makefile                   # Build automation
├── LICENSE                    # Software license
└── README.md                  # Project documentation
```
