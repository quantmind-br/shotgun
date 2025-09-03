# Project Structure

## Current Directory Structure
```
shotgun/
├── .bmad-core/           # BMAD framework configuration and tasks
├── cmd/shotgun/          # Main application entry point
├── docs/                 # Documentation including architecture and stories
│   ├── architecture/     # Architecture documentation
│   └── stories/          # Development stories
├── internal/             # Private application packages
│   ├── core/
│   │   └── scanner/      # File scanning utilities (implemented)
│   └── models/           # Data models
├── templates/            # Built-in prompt templates
├── go.mod               # Go module definition
├── go.sum              # Dependency checksums
├── Makefile            # Build automation
└── README.md           # Project documentation
```

## Key Implemented Components
- **Scanner Package**: Core file scanning engine with concurrency support
- **Models Package**: Type definitions (internal/models/files.go)
- **Build System**: Complete Makefile with cross-platform builds

## Story-Driven Development
- Stories located in `docs/stories/`
- Following BMAD Core framework
- Currently on Story 1.3 (as per user request)
- Previous stories: 1.1 (project setup), 1.2 (file scanner engine)