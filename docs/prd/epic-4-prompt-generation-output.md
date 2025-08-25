# Epic 4 - Prompt Generation & Output

**Epic Goal:** Develop the comprehensive template processing engine, concurrent file content aggregation system, real-time size estimation with progress visualization, and final Markdown output generation with robust error handling to deliver the complete MVP functionality.

### Story 4.1: Template Processing Engine
As a user,
I want the system to process templates with my inputs and selected files,
so that I can generate customized prompts based on the chosen template structure.

#### Acceptance Criteria
1. Template processing engine using Go's text/template package with custom function maps
2. Variable substitution system supporting all template placeholders (TASK, RULES, FILE_STRUCTURE, etc.)
3. Template function library including string manipulation (title, upper, lower, trim, wordCount, lineCount)
4. Conditional logic support in templates ({{if}}, {{range}}, {{with}}) for dynamic content
5. Template validation during processing with meaningful error messages for syntax issues
6. Safety measures preventing template injection or code execution vulnerabilities
7. Performance: template processing completes in <2 seconds for standard templates

### Story 4.2: Concurrent File Content Aggregation
As a user,
I want the system to efficiently read and organize selected file contents,
so that I can include comprehensive project context in my prompts.

#### Acceptance Criteria
1. Concurrent file reading using worker pool pattern with configurable concurrency
2. File content aggregation in structured format with file path headers
3. Directory tree generation using ASCII characters (├── └── │) for visual hierarchy
4. Content format: `<file path="RELATIVE/PATH">` blocks with actual file contents
5. Binary file exclusion with automatic detection and user notification
6. Large file handling with size limits and user warnings for oversized content
7. Performance: aggregate 100+ files in <10 seconds with progress indication
8. Error handling for file permission issues, missing files, or encoding problems

### Story 4.3: Real-Time Size Estimation System
As a user,
I want to see an accurate estimate of the final prompt size before generation,
so that I can make informed decisions about file inclusion and template complexity.

#### Acceptance Criteria
1. Real-time size calculation as user modifies selections or input content
2. Size estimation in human-readable units (KB, MB) with byte precision available
3. Progress bar visualization during size calculation with smooth animations
4. Warning thresholds for large prompts (>500KB warning, >1MB strong warning)
5. Breakdown showing contribution from different sections (template, files, inputs)
6. Performance impact estimation and recommendations for large file sets
7. Size recalculation triggers on file selection changes or input modifications

### Story 4.4: Confirmation Screen with Summary
As a user,
I want to review all my selections and see the complete summary before generation,
so that I can verify everything is correct before creating the final prompt.

#### Acceptance Criteria
1. Comprehensive summary screen showing template, selected files count, excluded items
2. Size estimate with visual progress bar and warning indicators
3. Generated filename preview with timestamp format (shotgun_prompt_YYYYMMDD_HHMM.md)
4. Final review of task description and rules (truncated preview if very long)
5. Clear action buttons: F2 (back for changes), F10 (generate), ESC (cancel)
6. Screen integration as "Confirm Generation [5/5]" completing the wizard
7. Warning messages for potential issues (very large size, no files selected, etc.)

### Story 4.5: Markdown Output Generation
As a user,
I want the system to generate a complete Markdown prompt file in the current directory,
so that I can immediately use the prompt with my preferred LLM tools.

#### Acceptance Criteria
1. Final prompt assembly combining template structure with processed content
2. Markdown file generation with timestamped filename in current working directory
3. File structure section with complete directory tree and file contents
4. Proper Markdown formatting with headers, code blocks, and structured sections
5. UTF-8 encoding support maintaining all special characters from inputs
6. Atomic file writing preventing partial files on interruption or errors
7. Success confirmation with file path display and file size information
8. Error handling for disk space issues, permission problems, or write failures
9. Generation progress indication during file writing process
