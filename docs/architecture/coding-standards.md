# Coding Standards

### Critical Fullstack Rules

- **Package Structure:** Always follow the established internal/ structure - ui/ for TUI components, core/ for business logic, services/ for orchestration
- **Error Handling:** All public functions must return error as the last return value and handle errors explicitly - never ignore errors with `_`
- **Context Propagation:** Long-running operations (file scanning, template processing) must accept context.Context as first parameter for cancellation
- **Interface Usage:** Always program to interfaces for testability - FileScanner, TemplateEngine, etc. must be interfaces, not concrete types
- **Resource Cleanup:** All file operations, goroutines, and channels must have explicit cleanup using defer or context cancellation
- **Bubble Tea Messages:** All UI state changes must flow through proper tea.Msg types - never modify model fields directly in Update methods
- **Concurrent Safety:** Shared state must be protected with sync.Mutex or use channels for communication between goroutines
- **Configuration Access:** Never access config values directly - always use the ConfigManager service with proper defaults
- **Template Security:** All template output must be sanitized before display to prevent terminal escape sequence injection
- **File Path Validation:** All file paths must be validated against directory traversal attacks using filepath.Clean and boundary checks

### Naming Conventions

| Element | Frontend | Backend | Example |
|---------|----------|---------|---------|
| TUI Models | PascalCase with Model suffix | - | `FileTreeModel`, `TemplateSelectionModel` |
| TUI Components | PascalCase with Component suffix | - | `FileTreeComponent`, `StatusBarComponent` |
| Services | PascalCase with Service suffix | PascalCase with Service suffix | `FileScanner`, `TemplateEngine` |
| Interfaces | PascalCase | PascalCase | `Scanner`, `Engine`, `Repository` |
| Methods | camelCase | camelCase | `scanDirectory()`, `processTemplate()` |
| Constants | SCREAMING_SNAKE_CASE | SCREAMING_SNAKE_CASE | `MAX_FILE_SIZE`, `DEFAULT_TIMEOUT` |
| Packages | lowercase | lowercase | `scanner`, `template`, `ignore` |
| Files | snake_case | snake_case | `file_tree.go`, `template_engine.go` |
