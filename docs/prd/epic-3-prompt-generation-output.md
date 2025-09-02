# Epic 3: Prompt Generation & Output

Implement the prompt assembly system that combines templates with file contents, create the confirmation screen with size estimation, and generate the final Markdown output file with proper formatting.

## Story 3.1: Template Processing Engine

As a developer,  
I want templates to be processed with variable substitution,  
so that dynamic content is properly inserted into the output.

### Acceptance Criteria
1: Go text/template engine processes template content
2: Variable substitution works for all defined variables (TASK, RULES, FILE_STRUCTURE)
3: Automatic variables populated (CURRENT_DATE, PROJECT_NAME)
4: Conditional sections ({{if}}) process correctly
5: Template functions available (upper, lower, trim)
6: Error handling for missing variables or template errors
7: Unit tests cover various template scenarios

## Story 3.2: File Structure Assembly

As a user,  
I want selected files to be included with their content in the prompt,  
so that the LLM has full context for analysis.

### Acceptance Criteria
1: Generate tree-like structure showing directory hierarchy
2: Use ASCII characters for tree visualization (├── └── │)
3: Read file contents asynchronously using goroutines
4: Wrap file contents in <file path="...">content</file> tags
5: Skip binary files with appropriate message
6: Handle large files gracefully (streaming read)
7: Preserve exact file content including whitespace

## Story 3.3: Confirmation Screen with Size Estimation

As a user,  
I want to review prompt details before generation,  
so that I can ensure the output will be appropriate.

### Acceptance Criteria
1: Display summary of selections (template name, file count, excluded items)
2: Calculate and display estimated output size in KB/MB
3: Show output filename with timestamp
4: Progress bar displays during size calculation
5: Warning appears for very large outputs (>500KB)
6: F10 confirms and triggers generation
7: F2 allows returning to make adjustments

## Story 3.4: Prompt Generation & File Writing

As a user,  
I want the final prompt saved to a file,  
so that I can use it with my preferred LLM tool.

### Acceptance Criteria
1: Combine template + variables + file structure into final output
2: Save to current directory with timestamp filename (shotgun_prompt_YYYYMMDD_HHMM.md)
3: Ensure no filename collisions with incrementing counter if needed
4: Display success message with full file path
5: Handle write errors gracefully with clear error message
6: Non-blocking generation using goroutines
7: Progress indicator during file writing for large outputs
