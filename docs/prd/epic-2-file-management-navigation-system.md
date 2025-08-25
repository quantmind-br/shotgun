# Epic 2 - File Management & Navigation System

**Epic Goal:** Implement the complete file tree interface with hierarchical checkbox selection, advanced ignore rule processing (.gitignore and .shotgunignore), binary file handling, and full keyboard navigation to provide comprehensive file context gathering functionality.

### Story 2.1: Hierarchical File Tree Widget
As a user,
I want to see my project files displayed in a tree structure with folder expansion/collapse,
so that I can understand the project organization and navigate efficiently.

#### Acceptance Criteria
1. Tree widget component built on Bubbles filepicker with hierarchical display
2. ASCII tree characters (├── └── │) for clean visual hierarchy
3. Folder expansion/collapse with → and ← arrow keys
4. Directory icons and file type indicators for visual distinction
5. Viewport scrolling for projects with many files using Bubbles viewport
6. Performance: smooth rendering for trees with 1000+ items
7. Keyboard navigation with ↑↓ keys maintaining current selection state

### Story 2.2: Checkbox Selection System with All-Selected Default
As a user,
I want all files to start selected with checkboxes so I can deselect unwanted items,
so that I can quickly exclude irrelevant files while including everything relevant by default.

#### Acceptance Criteria
1. All discovered files and folders display with checked checkboxes initially
2. Space bar toggles individual item selection state
3. Hierarchical selection: unchecking folder unchecks all contents automatically
4. Visual distinction between checked [✓], unchecked [ ], and partially selected [◐] states
5. Selection counter showing "X selected · Y excluded · Z ignored" in status bar
6. Ctrl+A select all, Ctrl+I invert selection keyboard shortcuts
7. Selection state persists during tree navigation and expansion/collapse operations

### Story 2.3: Enhanced Ignore Rules with .shotgunignore Support
As a user,
I want to use project-specific .shotgunignore rules in addition to .gitignore,
so that I can customize file exclusions for prompt generation without affecting Git.

#### Acceptance Criteria
1. .shotgunignore file detection and parsing in project root
2. Combined rule processing: .gitignore + .shotgunignore patterns
3. Rule precedence: .shotgunignore overrides .gitignore when conflicts exist
4. Pattern support: glob patterns, directory exclusions, negation with !
5. Visual indicators for ignored files (grayed out, "ignored" label)
6. Performance: rule processing adds <15% overhead to scanning
7. `shotgun init` command creates example .shotgunignore file

### Story 2.4: Binary File Detection and Handling
As a user,
I want binary files automatically identified and excluded from selection,
so that I don't accidentally include non-text files in my prompts.

#### Acceptance Criteria
1. Automatic binary detection using filetype library during scanning
2. Binary files visually distinct (🔒 icon, grayed appearance, non-selectable)
3. Common binary extensions pre-configured (.exe, .jpg, .png, .zip, .pdf, etc.)
4. Binary files excluded from selection count but visible in tree
5. Tooltip or status indication explaining why binary files cannot be selected
6. Override capability for edge cases where binary content is needed
7. Performance: detection adds <5% overhead to file scanning

### Story 2.5: Complete File Tree Screen Integration
As a user,
I want the complete file tree interface integrated as Screen 1 of the wizard,
so that I can select files and navigate to template selection seamlessly.

#### Acceptance Criteria
1. File tree screen integrated as first screen in 5-screen wizard flow
2. Header showing "File Selection [1/5]" with progress indicator
3. Status bar showing selection counts and keyboard shortcuts
4. F3 key advances to next screen when at least one file is selected
5. F1 key shows contextual help for file tree navigation
6. Screen state preservation when navigating back from subsequent screens
7. Validation: prevents advancement if no files are selected
8. Loading state during initial file scanning with progress indication
