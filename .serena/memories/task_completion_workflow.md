# Task Completion Workflow

## Standard Development Process
When completing any development task, follow this sequence:

### 1. Code Quality Checks
- `make lint` - Run formatting and static analysis
- Fix any issues reported by fmt/vet

### 2. Testing
- `make test` - Run all unit tests
- Ensure 90% minimum coverage for core packages
- `make test-coverage` - Generate coverage report if needed

### 3. Build Verification
- `make build` - Ensure clean compilation
- Test binary functionality manually if applicable

### 4. Story-Specific Requirements
- Update story file Dev Agent Record sections only
- Mark completed tasks with [x]
- Update File List in story
- Add entries to Change Log
- Set status when story complete

### 5. BMAD Framework Compliance
- Only update authorized story sections:
  - Tasks/Subtasks checkboxes
  - Dev Agent Record section
  - Debug Log References
  - Completion Notes List
  - File List
  - Change Log
  - Status

### 6. Critical Rules
- NEVER modify Story, Acceptance Criteria, Dev Notes, Testing sections
- Follow develop-story command execution order
- HALT on: unapproved dependencies, ambiguous requirements, repeated failures
- Set to "Ready for Review" only when all validations pass