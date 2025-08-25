# Epic 3 - Template System & User Input

**Epic Goal:** Create comprehensive template discovery and management system for built-in and custom user templates, implement advanced multiline text editors with full UTF-8 support, and complete the 5-screen wizard workflow with seamless navigation and state management.

### Story 3.1: Template Discovery and Metadata System
As a user,
I want to see all available templates (built-in and custom) with rich metadata,
so that I can choose the most appropriate template for my current task.

#### Acceptance Criteria
1. Template discovery system scanning embedded templates and user config directories
2. Cross-platform config directory support (~/.config/shotgun-cli/templates, %APPDATA%/shotgun-cli/templates)
3. TOML metadata parsing for template name, version, description, author, and tags
4. Template validation ensuring required sections and variables are present
5. Unified template list showing both built-in and custom templates without visual distinction
6. Template precedence: user templates override built-in templates with same name
7. Error handling for malformed templates with user-friendly error messages

### Story 3.2: Template Selection Interface
As a user,
I want to select templates from a clean, navigable list interface,
so that I can quickly choose the right template and understand its purpose.

#### Acceptance Criteria
1. Vertical list interface using Bubbles list component with template metadata display
2. Template entries showing name, version, and description in structured format
3. Keyboard navigation with ↑↓ arrows and Enter/F3 for selection
4. Visual highlighting of currently selected template with focus indicators
5. Screen integration as "Template Selection [2/5]" in wizard flow
6. F2 navigation back to file tree with state preservation
7. Template preview or expanded description available via additional key press
8. Validation: prevents advancement without template selection

### Story 3.3: Advanced Multiline Text Editor for Task Input
As a user,
I want a sophisticated text editor for describing my task with full UTF-8 support,
so that I can provide detailed context in any language with special characters.

#### Acceptance Criteria
1. Multiline text editor built on Bubbles textarea v0.21.0 with horizontal scrolling
2. Full UTF-8 character support including accented characters (ç, á, ô, ñ, etc.)
3. Text input features: word wrap, line numbers, character/word count display
4. Clipboard integration (Ctrl+V paste, Ctrl+C copy) with proper UTF-8 handling
5. Editor modes: editing mode (normal text input) and navigation mode (F-key shortcuts)
6. Ctrl+Enter toggles between editing and navigation modes
7. Screen integration as "Task Description [3/5]" with state preservation
8. Validation: prevents advancement with empty task content

### Story 3.4: Optional Rules Input Editor
As a user,
I want an optional editor for additional rules and constraints,
so that I can provide specific guidance for the LLM when needed.

#### Acceptance Criteria
1. Optional multiline editor with same advanced features as task editor
2. Clear visual indication that field is optional with "optional" label
3. F4 key to skip this step entirely and advance to next screen
4. Auto-save functionality preserving content when navigating between screens
5. Screen integration as "Rules · optional [4/5]" with appropriate styling
6. Same UTF-8 and clipboard support as task editor
7. Graceful handling of empty content (treated as no additional rules)

### Story 3.5: Complete Wizard Flow Integration
As a user,
I want seamless navigation between all screens with preserved state,
so that I can review and modify any step without losing my work.

#### Acceptance Criteria
1. Complete 5-screen wizard flow: File Tree → Templates → Task → Rules → Confirmation
2. Global navigation: F2 (back), F3 (forward), F1 (help), ESC (exit with confirmation)
3. State preservation: all selections, text content, and UI states maintained during navigation
4. Progress indicator showing current screen (e.g., [3/5]) in all screen headers
5. Screen-specific validation preventing invalid forward navigation
6. Contextual help (F1) showing relevant shortcuts and instructions per screen
7. Consistent styling and theme across all screens using Lip Gloss monochrome palette
8. Error states and user feedback for invalid operations or missing required content
