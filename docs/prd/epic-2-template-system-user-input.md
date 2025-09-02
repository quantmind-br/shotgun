# Epic 2: Template System & User Input

Create the template discovery and selection system, implement the template parsing engine, and build the multiline text input components that allow users to provide task descriptions and optional rules for their prompts.

## Story 2.1: Template Discovery & Loading System

As a user,  
I want the application to find and load available templates,  
so that I can choose from built-in or custom templates.

### Acceptance Criteria
1: Built-in templates embedded in binary from templates/ directory
2: User templates discovered from ~/.config/shotgun-cli/templates (or Windows equivalent)
3: TOML parser correctly reads template metadata (name, version, description, variables)
4: Template validation ensures required fields are present
5: Templates with parsing errors are skipped with warning logs
6: Both built-in and user templates appear in unified list
7: Unit tests cover template discovery and parsing logic

## Story 2.2: Template Selection UI Screen

As a user,  
I want to browse and select a template from available options,  
so that I can use the appropriate format for my use case.

### Acceptance Criteria
1: List view displays all discovered templates with name and description
2: Version number shows inline with template name
3: Arrow keys navigate up/down through template list
4: Enter or F3 selects template and advances to next screen
5: Template metadata (author, tags) displays in detail panel
6: Visual distinction for currently selected template
7: F2 returns to file tree screen preserving selection

## Story 2.3: Task Input Screen with Multiline Editor

As a user,  
I want to describe my task in detail with proper formatting,  
so that the LLM receives clear context about what I need.

### Acceptance Criteria
1: Multiline text editor using Bubbles textarea component
2: Support for copy/paste operations (Ctrl+C, Ctrl+V)
3: UTF-8 character support for international text
4: Line and character count displays in real-time
5: Ctrl+Enter finalizes input and advances
6: F3 advances only if content is non-empty
7: Text persists when navigating back via F2
8: Word wrap functions correctly for long lines

## Story 2.4: Rules Input Screen (Optional)

As a user,  
I want to optionally specify additional rules or constraints,  
so that I can guide the LLM's response style or requirements.

### Acceptance Criteria
1: Multiline text editor similar to task input screen
2: Clear indication that this field is optional
3: F4 key skips this screen entirely
4: F3 advances regardless of content (can be empty)
5: Content persists when navigating between screens
6: Placeholder text suggests example rules
7: Same UTF-8 and clipboard support as task input
