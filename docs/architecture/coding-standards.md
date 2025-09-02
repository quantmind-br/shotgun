# Coding Standards

## Critical Fullstack Rules

- **Type Definitions:** Always define types in internal/models and import from there
- **Error Handling:** All functions that can fail must return error as last value
- **Concurrent Access:** Use channels or sync primitives for shared state
- **Resource Cleanup:** Always use defer for cleanup operations
- **Context Propagation:** Pass context.Context for cancellable operations
- **Testing Coverage:** Minimum 90% coverage for core packages
- **Panic Recovery:** Never panic in library code, only in main()

## Naming Conventions

| Element | Frontend | Backend | Example |
|---------|----------|---------|---------|
| Components | PascalCase | - | `FileTreeModel` |
| Functions | camelCase | camelCase | `scanDirectory()` |
| Interfaces | PascalCase+er | PascalCase+er | `Scanner` |
| Packages | lowercase | lowercase | `scanner` |
