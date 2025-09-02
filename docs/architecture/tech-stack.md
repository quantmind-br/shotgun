# Tech Stack

## Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|-----------|
| Frontend Language | Go | 1.22+ | TUI implementation | Native terminal support, cross-platform |
| Frontend Framework | Bubble Tea | v2.0.0-beta.4 | Terminal UI framework | Elm architecture, excellent keyboard support |
| UI Component Library | Bubbles | v0.21.0 | Pre-built TUI components | Accelerates development with tested components |
| State Management | Bubble Tea Models | Built-in | Immutable state management | Integrated with framework, predictable updates |
| Backend Language | Go | 1.22+ | Core business logic | Same as frontend for unified codebase |
| Backend Framework | Standard Library | 1.22+ | File operations, concurrency | Minimal dependencies, excellent performance |
| API Style | N/A | - | No API needed | Local-only application |
| Database | File System | OS-provided | Configuration and templates | No database needed for MVP |
| Cache | sync.Map | Standard library | In-memory caching | Thread-safe caching for file metadata |
| File Storage | Local File System | OS-provided | Template and output storage | Direct file system access |
| Authentication | N/A | - | No auth needed | Local-only application |
| Frontend Testing | Go testing | Standard library | Unit tests for UI components | Native Go testing support |
| Backend Testing | Go testing | Standard library | Unit tests for business logic | Consistent testing approach |
| E2E Testing | teatest | Bubble Tea testing | TUI flow testing | Framework-specific testing utilities |
| Build Tool | Go build | 1.22+ | Compilation and linking | Native Go toolchain |
| Bundler | N/A | - | Single binary output | Go compiles to single executable |
| IaC Tool | N/A | - | No infrastructure needed | Local-only application |
| CI/CD | GitHub Actions | Latest | Automated testing and releases | Free for open source, good Go support |
| Monitoring | N/A | - | No monitoring needed for MVP | Local application, no telemetry |
| Logging | log/slog | Standard library | Structured logging | Native structured logging in Go 1.21+ |
| CSS Framework | Lip Gloss | v1.0.0 | Terminal styling | Programmatic styling for TUI |
