# Epic 1 - Foundation & Core Infrastructure

**Epic Goal:** Establish the technical foundation with Go project structure, Bubble Tea TUI framework, concurrent file scanning engine, and a basic interface that validates the complete technical stack while delivering initial functionality.

### Story 1.1: Project Setup and Go Module Initialization
As a developer,
I want a properly structured Go project with all dependencies configured,
so that I can begin development with a solid, maintainable foundation.

#### Acceptance Criteria
1. Go module initialized with go.mod specifying Go 1.22+ requirement
2. Complete directory structure created matching PLAN.md specifications (cmd/, internal/, templates/)
3. All core dependencies added to go.mod (Bubble Tea v2.0.0-beta.4, Bubbles v0.21.0, Lip Gloss v1.0.0, etc.)
4. Basic Makefile with build, test, clean targets
5. .gitignore configured for Go projects with common exclusions
6. Initial README.md with project overview and build instructions

### Story 1.2: Core File Scanner Engine with Concurrency
As a user,
I want the system to efficiently scan and catalog project files in the current directory,
so that I can work with projects of any size without performance issues.

#### Acceptance Criteria
1. FileScanner struct implemented with concurrent goroutine-based directory traversal
2. Worker pool pattern using runtime.NumCPU() for optimal performance
3. Channel-based file information streaming to prevent memory bottlenecks
4. Basic binary file detection using h2non/filetype library
5. File metadata collection (size, modification time, type) for all discovered files
6. Performance requirement: scan 1000+ files in under 5 seconds
7. Graceful error handling for permission denied and other file system errors

### Story 1.3: Basic Ignore Rules Processing
As a user,
I want the file scanner to respect .gitignore patterns,
so that irrelevant files are automatically excluded from consideration.

#### Acceptance Criteria
1. .gitignore file parsing using doublestar/v4 glob pattern matching
2. Ignore rule application during file scanning phase
3. Common patterns automatically excluded (node_modules/, .git/, dist/, build/)
4. Proper handling of nested .gitignore files and directory-specific rules
5. File exclusion logged at debug level for troubleshooting
6. Performance: ignore processing adds <10% to scan time overhead

### Story 1.4: Bubble Tea Application Foundation
As a developer,
I want the TUI application framework properly initialized,
so that all future UI screens can be built on a solid reactive architecture.

#### Acceptance Criteria
1. Main Bubble Tea application struct created following Elm Architecture
2. Initial model with state management for application lifecycle
3. Root command handler with proper initialization and cleanup
4. Global key bindings implemented (ESC to exit, F1 for help)
5. Terminal size detection and responsive layout foundation
6. Error recovery and graceful shutdown on panic or interrupt
7. TEA_DEBUG support for development debugging

### Story 1.5: Health Check Interface
As a user,
I want to see a simple interface that validates the technical stack,
so that I know the application is working correctly and can provide feedback.

#### Acceptance Criteria
1. Basic TUI screen displaying "shotgun-cli" title and version
2. File scan progress indicator showing discovered files count
3. Simple status display showing scan completion and total files found
4. "Press ESC to exit" instruction for user interaction
5. Real-time updates during file scanning operation
6. Basic styling using Lip Gloss with monochrome theme foundation
7. Cross-platform terminal compatibility validated (Windows/Linux/macOS)
