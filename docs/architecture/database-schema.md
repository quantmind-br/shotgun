# Database Schema

**shotgun-cli-v3** is designed as a standalone TUI application that operates entirely through file system operations without requiring a traditional database. However, the application does have structured data storage requirements that are implemented through file-based persistence:

### File-Based Data Storage

Instead of a database, the application uses structured file storage with the following schema:

#### Session Storage Schema (JSON)

**File Location:** `~/.config/shotgun-cli/sessions/{session-id}.json`

```json
{
  "session_id": "20240315-143022-abc123",
  "created_at": "2024-03-15T14:30:22Z",
  "modified_at": "2024-03-15T14:35:10Z",
  "project_path": "/Users/dev/my-project",
  "current_screen": "confirmation",
  "file_selection": {
    "root_path": "/Users/dev/my-project",
    "total_files": 156,
    "selected_count": 23,
    "ignored_count": 89,
    "tree": {
      "path": ".",
      "name": "my-project",
      "is_directory": true,
      "is_selected": true,
      "is_ignored": false,
      "is_binary": false,
      "size": 0,
      "children": [...]
    }
  },
  "template_choice": {
    "template_id": "analyze_bug",
    "template_name": "Bug Analysis",
    "is_builtin": true,
    "file_path": "embedded://templates/analyze_bug.tmpl"
  },
  "user_inputs": {
    "task_description": "Debug memory leak in file processing",
    "additional_rules": "Focus on goroutine usage and channel handling",
    "character_count": 89,
    "word_count": 12
  },
  "generation_config": {
    "output_filename": "shotgun_prompt_20240315_1430.md",
    "estimated_size": 45678,
    "include_binary_warning": true,
    "max_file_size": 1048576
  }
}
```

#### Template Metadata Schema (TOML)

**File Location:** `~/.config/shotgun-cli/templates/{template-name}.toml`

```toml
[metadata]
id = "custom_code_review"
name = "Code Review Assistant"
version = "1.2.0"
description = "Comprehensive code review with security and performance focus"
author = "Development Team"
tags = ["code-review", "security", "performance"]
created_at = "2024-03-10T09:00:00Z"
updated_at = "2024-03-15T11:30:00Z"

[variables]
required = ["TASK", "FILE_STRUCTURE", "FILE_CONTENTS"]
optional = ["RULES", "PROJECT_NAME"]

[settings]
max_file_size = 1048576
exclude_binary = true
preserve_structure = true

[content]
template = """
# Code Review Request

## Task
{{.TASK}}

{{if .RULES}}
## Additional Guidelines
{{.RULES}}
{{end}}

## Project Structure
{{.FILE_STRUCTURE}}

## File Contents
{{range $path, $content := .FILE_CONTENTS}}
### {{$path}}
```
{{$content}}
```
{{end}}
"""
```

#### Configuration Schema (TOML)

**File Location:** `~/.config/shotgun-cli/config.toml`

```toml
[general]
default_template = "make_plan"
auto_save_sessions = true
max_sessions = 50
session_cleanup_days = 30

[templates]
directory = "~/.config/shotgun-cli/templates"
auto_discover = true
validate_on_load = true

[file_processing]
max_file_size = 1048576  # 1MB
max_total_size = 10485760  # 10MB  
worker_count = 0  # 0 = auto-detect based on CPU cores
scan_timeout = 30  # seconds

[ui]
theme = "monochrome"
accent_color = "#6ee7b7"
highlight_color = "#fbbf24"
animation_speed = "normal"

[output]
default_filename = "shotgun_prompt_{{.Timestamp}}.md"
output_directory = "."
backup_prompts = false
```

### Data Relationships and Constraints

#### File System Constraints

```plaintext
Session Files:
- Filename format: YYYYMMDD-HHMMSS-{random}.json
- Maximum size: 10MB (prevent excessive session data)
- Retention: Configurable cleanup after N days
- Validation: JSON schema validation on load

Template Files:
- Filename format: {template-id}.toml
- TOML metadata validation required
- Template content must pass Go text/template parsing
- No circular dependencies in template references

Configuration:
- Single config.toml file per user
- Environment variable overrides supported
- Default values for all required settings
- Backward compatibility for schema versions
```

#### Index Files for Performance

**Session Index:** `~/.config/shotgun-cli/sessions/index.json`

```json
{
  "last_updated": "2024-03-15T14:35:10Z",
  "sessions": [
    {
      "id": "20240315-143022-abc123",
      "created_at": "2024-03-15T14:30:22Z",
      "project_path": "/Users/dev/my-project",
      "summary": "Bug Analysis - my-project",
      "file_count": 23,
      "estimated_size": 45678
    }
  ]
}
```
