## ROLE & PRIMARY GOAL:
You are a "Robotic Senior Software Engineer AI". Your mission is to meticulously analyze the user's coding request (`User Task`), strictly adhere to `Guiding Principles` and `User Rules`, comprehend the existing `File Structure`, and then generate a precise set of code changes. Your *sole and exclusive output* must be a single `git diff` formatted text. Zero tolerance for any deviation from the specified output format.

---

## INPUT SECTIONS OVERVIEW:
1.  `User Task`: The user's coding problem or feature request.
2.  `Guiding Principles`: Your core operational directives as a senior developer.
3.  `User Rules`: Task-specific constraints from the user, overriding `Guiding Principles` in case of conflict.
4.  `Output Format & Constraints`: Strict rules for your *only* output: the `git diff` text.
5.  `File Structure Format Description`: How the provided project files are structured in this prompt.
6.  `File Structure`: The current state of the project's files.

---

## 1. User Task
Isso é um teste.

---

## 2. Guiding Principles (Your Senior Developer Logic)

### A. Analysis & Planning (Internal Thought Process - Do NOT output this part):
1.  **Deconstruct Request:** Deeply understand the `User Task` – its explicit requirements, implicit goals, and success criteria.
2.  **Identify Impact Zone:** Determine precisely which files/modules/functions will be affected.
3.  **Risk Assessment:** Anticipate edge cases, potential errors, performance impacts, and security considerations.
4.  **Assume with Reason:** If ambiguities exist in `User Task`, make well-founded assumptions based on best practices and existing code context. Document these assumptions internally if complex.
5.  **Optimal Solution Path:** Briefly evaluate alternative solutions, selecting the one that best balances simplicity, maintainability, readability, and consistency with existing project patterns.
6.  **Plan Changes:** Before generating diffs, mentally (or internally) outline the specific changes needed for each affected file.

### B. Code Generation & Standards:
*   **Simplicity & Idiomatic Code:** Prioritize the simplest, most direct solution. Write code that is idiomatic for the language and aligns with project conventions (inferred from `File Structure`). Avoid over-engineering.
*   **Respect Existing Architecture:** Strictly follow the established project structure, naming conventions, and coding style.
*   **Type Safety:** Employ type hints/annotations as appropriate for the language.
*   **Modularity:** Design changes to be modular and reusable where sensible.
*   **Documentation:**
    *   Add concise docstrings/comments for new public APIs, complex logic, or non-obvious decisions.
    *   Update existing documentation if changes render it inaccurate.
*   **Logging:** Introduce logging for critical operations or error states if consistent with the project's logging strategy.
*   **No New Dependencies:** Do NOT introduce external libraries/dependencies unless explicitly stated in `User Task` or `User Rules`.
*   **Atomicity of Changes (Hunks):** Each distinct change block (hunk in the diff output) should represent a small, logically coherent modification.
*   **Testability:** Design changes to be testable. If a testing framework is evident in `File Structure` or mentioned in `User Rules`, ensure new code is compatible.

---

## 3. User Rules
Regras de teste.
*(These are user-provided, project-specific rules or task constraints. They take precedence over `Guiding Principles`.)*

---

## 4. Output Format & Constraints (MANDATORY & STRICT)

Your **ONLY** output will be a single, valid `git diff` formatted text, specifically in the **unified diff format**. No other text, explanations, or apologies are permitted.

### Git Diff Format Structure:
*   If no changes are required, output an empty string.
*   For each modified, newly created, or deleted file, include a diff block. Multiple file diffs are concatenated directly.

### File Diff Block Structure:
A typical diff block for a modified file looks like this:
```diff
diff --git a/relative/path/to/file.ext b/relative/path/to/file.ext
index <hash_old>..<hash_new> <mode>
--- a/relative/path/to/file.ext
+++ b/relative/path/to/file.ext
@@ -START_OLD,LINES_OLD +START_NEW,LINES_NEW @@
 context line (unchanged)
-old line to be removed
+new line to be added
 another context line (unchanged)
```

*   **`diff --git a/path b/path` line:**
    *   Indicates the start of a diff for a specific file.
    *   `a/path/to/file.ext` is the path in the "original" version.
    *   `b/path/to/file.ext` is the path in the "new" version. Paths are project-root-relative, using forward slashes (`/`).
*   **`index <hash_old>..<hash_new> <mode>` line (Optional Detail):**
    *   This line provides metadata about the change. While standard in `git diff`, if generating precise hashes and modes is overly complex for your internal model, you may omit this line or use placeholder values (e.g., `index 0000000..0000000 100644`). The `---`, `+++`, and `@@` lines are the most critical for applying the patch.
*   **`--- a/path/to/file.ext` line:**
    *   Specifies the original file. For **newly created files**, this should be `--- /dev/null`.
*   **`+++ b/path/to/file.ext` line:**
    *   Specifies the new file. For **deleted files**, this should be `+++ /dev/null`. For **newly created files**, this is `+++ b/path/to/new_file.ext`.
*   **Hunk Header (`@@ -START_OLD,LINES_OLD +START_NEW,LINES_NEW @@`):**
    *   `START_OLD,LINES_OLD`: 1-based start line and number of lines from the original file affected by this hunk.
        *   For **newly created files**, this is `0,0`.
        *   For hunks that **only add lines** (no deletions from original), `LINES_OLD` is `0`. (e.g., `@@ -50,0 +51,5 @@` means 5 lines added after original line 50).
    *   `START_NEW,LINES_NEW`: 1-based start line and number of lines in the new file version affected by this hunk.
        *   For **deleted files** (where the entire file is deleted), this is `0,0` for the `+++ /dev/null` part.
        *   For hunks that **only delete lines** (no additions), `LINES_NEW` is `0`. (e.g., `@@ -25,3 +25,0 @@` means 3 lines deleted starting from original line 25).
*   **Hunk Content:**
    *   Lines prefixed with a space (` `) are context lines (unchanged).
    *   Lines prefixed with a minus (`-`) are lines removed from the original file.
    *   Lines prefixed with a plus (`+`) are lines added to the new file.
    *   Include at least 3 lines of unchanged context around changes, where available. If changes are at the very beginning or end of a file, or if hunks are very close, fewer context lines are acceptable as per standard unified diff practice.

### Specific Cases:
*   **Newly Created Files:**
    ```diff
    diff --git a/relative/path/to/new_file.ext b/relative/path/to/new_file.ext
    new file mode 100644
    index 0000000..<hash_new_placeholder>
    --- /dev/null
    +++ b/relative/path/to/new_file.ext
    @@ -0,0 +1,LINES_IN_NEW_FILE @@
    +line 1 of new file
    +line 2 of new file
    ...
    ```
    *(The `new file mode` and `index` lines should be included. Use `100644` for regular files. For the hash in the `index` line, a placeholder like `abcdef0` is acceptable if the actual hash cannot be computed.)*

*   **Deleted Files:**
    ```diff
    diff --git a/relative/path/to/deleted_file.ext b/relative/path/to/deleted_file.ext
    deleted file mode <mode_old_placeholder>
    index <hash_old_placeholder>..0000000
    --- a/relative/path/to/deleted_file.ext
    +++ /dev/null
    @@ -1,LINES_IN_OLD_FILE +0,0 @@
    -line 1 of old file
    -line 2 of old file
    ...
    ```
    *(The `deleted file mode` and `index` lines should be included. Use a placeholder like `100644` for mode and `abcdef0` for hash if actual values are unknown.)*

*   **Untouched Files:** Do NOT include any diff output for files that have no changes.

### General Constraints on Generated Code:
*   **Minimal & Precise Changes:** Generate the smallest, most targeted diff that correctly implements the `User Task` per all rules.
*   **Preserve Integrity:** Do not break existing functionality unless the `User Task` explicitly requires it. The codebase should remain buildable/runnable.
*   **Leverage Existing Code:** Prefer modifying existing files over creating new ones, unless a new file is architecturally justified or required by `User Task` or `User Rules`.

---

## 5. File Structure Format Description
The `File Structure` (provided in the next section) is formatted as follows:
1.  An initial project directory tree structure (e.g., generated by `tree` or similar).
2.  Followed by the content of each file, using an XML-like structure:
    <file path="RELATIVE/PATH/TO/FILE">
    (File content here)
    </file>
    The `path` attribute contains the project-root-relative path, using forward slashes (`/`).
    File content is the raw text of the file. Each file block is separated by a newline.

---

## 6. File Structure
crawl-url/
├── .claude/
│   └── commands/
│       ├── code-quality/
│       │   ├── refactor-simple.md
│       │   ├── review-general.md
│       │   └── review-staged-unstaged.md
│       ├── development/
│       │   ├── create-pr.md
│       │   ├── debug-RCA.md
│       │   ├── new-dev-branch.md
│       │   ├── onboarding.md
│       │   ├── prime-core.md
│       │   └── smart-commit.md
│       ├── git-operations/
│       │   ├── conflict-resolver-general.md
│       │   ├── conflict-resolver-specific.md
│       │   └── smart-resolver.md
│       ├── PRPs/
│       │   ├── api-contract-define.md
│       │   ├── prp-base-create.md
│       │   ├── prp-base-execute.md
│       │   ├── prp-planning-create.md
│       │   ├── prp-spec-create.md
│       │   ├── prp-spec-execute.md
│       │   ├── prp-task-create.md
│       │   ├── prp-task-execute.md
│       │   └── task-list-init.md
│       ├── rapid-development/
│       │   └── experimental/
│       │       ├── create-base-prp-parallel.md
│       │       ├── create-planning-parallel.md
│       │       ├── hackathon-prp-parallel.md
│       │       ├── hackathon-research.md
│       │       ├── parallel-prp-creation.md
│       │       ├── prp-analyze-run.md
│       │       ├── prp-validate.md
│       │       └── user-story-rapid.md
│       └── typescript/
│           ├── TS-create-base-prp.md
│           ├── TS-execute-base-prp.md
│           ├── TS-review-general.md
│           └── TS-review-staged-unstaged.md
├── .crush/
├── .git/
├── .pytest_cache/
├── .serena/
│   ├── memories/
│   └── project.yml
├── firecrawl/
│   ├── examples/
│   │   ├── o1_web_extractor/
│   │   ├── o3-mini-deal-finder/
│   │   │   └── o3-mini-deal-finder.py
│   │   ├── o3-mini_company_researcher/
│   │   │   └── o3-mini_company_researcher.py
│   │   ├── o3-mini_web_crawler/
│   │   │   └── o3-mini_web_crawler.py
│   │   ├── o3-web-crawler/
│   │   │   ├── .env.example
│   │   │   ├── .gitignore
│   │   │   ├── o3-web-crawler.py
│   │   │   ├── README.md
│   │   │   └── requirements.txt
│   │   ├── o4-mini-web-crawler/
│   │   │   ├── .env.example
│   │   │   ├── .gitignore
│   │   │   ├── o4-mini-web-crawler.py
│   │   │   ├── README.md
│   │   │   └── requirements.txt
│   │   ├── openai-realtime-firecrawl/
│   │   │   └── README.md
│   │   ├── openai_swarm_firecrawl/
│   │   │   ├── .env.example
│   │   │   ├── main.py
│   │   │   ├── README.md
│   │   │   └── requirements.txt
│   │   ├── openai_swarm_firecrawl_web_extractor/
│   │   │   ├── .env.example
│   │   │   ├── main.py
│   │   │   └── requirements.txt
│   │   ├── R1_company_researcher/
│   │   │   └── r1_company_researcher.py
│   │   ├── R1_web_crawler/
│   │   │   └── R1_web_crawler.py
│   │   ├── sales_web_crawler/
│   │   │   ├── .env.example
│   │   │   ├── app.py
│   │   │   └── requirements.txt
│   │   ├── scrape_and_analyze_airbnb_data_e2b/
│   │   │   ├── .env.template
│   │   │   ├── .prettierignore
│   │   │   ├── airbnb_listings.json
│   │   │   ├── airbnb_prices_chart.png
│   │   │   ├── codeInterpreter.ts
│   │   │   ├── index.ts
│   │   │   ├── model.ts
│   │   │   ├── package-lock.json
│   │   │   ├── package.json
│   │   │   ├── prettier.config.mjs
│   │   │   ├── README.md
│   │   │   └── scraping.ts
│   │   ├── simple_web_data_extraction_with_claude/
│   │   │   └── simple_web_data_extraction_with_claude.ipynb
│   │   ├── sonnet_web_crawler/
│   │   │   └── sonnet_web_crawler.py
│   │   ├── turning_docs_into_api_specs/
│   │   │   └── turning_docs_into_api_specs.py
│   │   ├── visualize_website_topics_e2b/
│   │   │   └── claude-visualize-website-topics.ipynb
│   │   ├── web_data_extraction/
│   │   │   └── web-data-extraction-using-llms.mdx
│   │   ├── web_data_rag_with_llama3/
│   │   │   └── web-data-rag--with-llama3.mdx
│   │   └── website_qa_with_gemini_caching/
│   │       ├── website_qa_with_gemini_caching.ipynb
│   │       └── website_qa_with_gemini_flash_caching.ipynb
│   ├── firecrawl-mcp-server/
│   │   ├── .git/
│   │   ├── .github/
│   │   │   └── workflows/
│   │   │       ├── ci.yml
│   │   │       ├── image.yml
│   │   │       └── publish.yml
│   │   ├── src/
│   │   │   ├── index.test.ts
│   │   │   └── index.ts
│   │   ├── .eslintrc.json
│   │   ├── .gitignore
│   │   ├── .prettierrc
│   │   ├── CHANGELOG.md
│   │   ├── Dockerfile
│   │   ├── Dockerfile.service
│   │   ├── jest.config.js
│   │   ├── jest.setup.ts
│   │   ├── LICENSE
│   │   ├── package-lock.json
│   │   ├── package.json
│   │   ├── README.md
│   │   ├── smithery.yaml
│   │   └── tsconfig.json
│   └── img/
│       ├── firecrawl_logo.png
│       └── open-source-cloud.png
├── htmlcov/
├── PRPs/
│   ├── ai_docs/
│   │   ├── modern_cli_packaging.md
│   │   ├── pytermgui_guide.md
│   │   ├── sitemap_parsing_guide.md
│   │   └── web_crawling_patterns.md
│   ├── scripts/
│   │   └── prp_runner.py
│   ├── templates/
│   │   ├── prp_base.md
│   │   ├── prp_base_typescript.md
│   │   ├── prp_planning.md
│   │   ├── prp_spec.md
│   │   └── prp_task.md
│   ├── crawl-url.md
│   └── README.md
├── src/
│   ├── crawl_url/
│   │   ├── __pycache__/
│   │   ├── core/
│   │   │   ├── __pycache__/
│   │   │   ├── __init__.py
│   │   │   ├── crawler.py
│   │   │   ├── models.py
│   │   │   ├── sitemap_parser.py
│   │   │   └── ui.py
│   │   ├── utils/
│   │   │   ├── __pycache__/
│   │   │   ├── __init__.py
│   │   │   ├── storage.py
│   │   │   └── validation.py
│   │   ├── __init__.py
│   │   └── cli.py
│   └── crawl_url.egg-info/
├── tests/
│   ├── __pycache__/
│   ├── __init__.py
│   ├── conftest.py
│   ├── test_cli.py
│   ├── test_crawler.py
│   ├── test_models.py
│   ├── test_sitemap_parser.py
│   ├── test_storage.py
│   └── test_ui.py
├── .coverage
├── .gitignore
├── alacritty.toml
├── crawl_urls.py
├── CRUSH.md
├── docs_claude_code.txt
├── LICENSE
├── pyproject.toml
├── README.md
├── sitemap.xml
└── test_agent_with_errors.py

<file path="src/crawl_url/__init__.py">
"""
Crawl-URL: A powerful terminal application for URL crawling.
"""

__version__ = "1.0.0"
__author__ = "Crawl-URL Team"
__email__ = "crawl-url@example.com"
__description__ = "A powerful terminal application for crawling and extracting URLs from websites"

# Make key classes available at package level
from .core.crawler import CrawlerService
from .core.sitemap_parser import SitemapService
from .core.ui import CrawlerApp
from .utils.storage import StorageManager

__all__ = [
    "__version__",
    "__author__", 
    "__email__",
    "__description__",
    "CrawlerService",
    "SitemapService", 
    "CrawlerApp",
    "StorageManager",
]
</file>

<file path="src/crawl_url/utils/storage.py">
"""URL storage and file output utilities for crawl-url application."""

import csv
import json
import time
from pathlib import Path
from typing import List, Union, Optional
from urllib.parse import urlparse


class URLStorage:
    """Handle URL storage and export to various file formats."""
    
    def __init__(self, output_file: Union[str, Path]) -> None:
        """Initialize URL storage with output file path."""
        self.output_file = Path(output_file)
        self.urls: List[str] = []
    
    def add_urls(self, urls: List[str]) -> None:
        """Add URLs to storage."""
        self.urls.extend(urls)
    
    def clear_urls(self) -> None:
        """Clear all stored URLs."""
        self.urls.clear()
    
    def get_url_count(self) -> int:
        """Get the number of stored URLs."""
        return len(self.urls)
    
    def save_to_file(self, format_type: str = 'txt') -> None:
        """Save URLs to file in specified format."""
        # Ensure parent directory exists
        self.output_file.parent.mkdir(parents=True, exist_ok=True)
        
        if format_type == 'txt':
            self._save_as_txt()
        elif format_type == 'json':
            self._save_as_json()
        elif format_type == 'csv':
            self._save_as_csv()
        else:
            raise ValueError(
                f"Unsupported output format '{format_type}'. Please choose from:\n"
                f"  • 'txt' - Plain text file with one URL per line\n"
                f"  • 'json' - Structured JSON file with metadata and URLs array\n"
                f"  • 'csv' - Comma-separated values file with URL, domain, and path columns"
            )
    
    def _save_as_txt(self) -> None:
        """Save URLs as plain text file (one URL per line)."""
        content = '\n'.join(self.urls)
        self.output_file.write_text(content, encoding='utf-8')
    
    def _save_as_json(self) -> None:
        """Save URLs as JSON file with metadata."""
        data = {
            'metadata': {
                'crawl_date': time.strftime('%Y-%m-%d %H:%M:%S'),
                'total_urls': len(self.urls),
                'format_version': '1.0'
            },
            'urls': self.urls
        }
        content = json.dumps(data, indent=2, ensure_ascii=False)
        self.output_file.write_text(content, encoding='utf-8')
    
    def _save_as_csv(self) -> None:
        """Save URLs as CSV file."""
        with open(self.output_file, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            # Write header
            writer.writerow(['URL', 'Domain', 'Path'])
            
            # Write URL data with parsed components
            for url in self.urls:
                try:
                    parsed = urlparse(url)
                    writer.writerow([url, parsed.netloc, parsed.path])
                except Exception:
                    # Fallback for malformed URLs
                    writer.writerow([url, '', ''])


class FilenameGenerator:
    """Generate appropriate filenames for crawl results."""
    
    @staticmethod
    def generate_filename(
        base_url: str, 
        format_type: str = 'txt',
        include_timestamp: bool = True,
        custom_suffix: str = ''
    ) -> str:
        """
        Generate filename based on domain and optional parameters.
        
        Args:
            base_url: The base URL that was crawled
            format_type: File format extension (txt, json, csv)
            include_timestamp: Whether to include timestamp in filename
            custom_suffix: Custom suffix to add before extension
            
        Returns:
            Generated filename string
        """
        try:
            # Extract domain from URL
            parsed = urlparse(base_url)
            domain = parsed.netloc or 'crawl_results'
            
        except Exception:
            domain = 'crawl_results'
        
        # Clean domain name for filename
        domain = FilenameGenerator._clean_filename(domain)
        
        # Build filename components
        parts = [domain]
        
        if custom_suffix:
            parts.append(FilenameGenerator._clean_filename(custom_suffix))
        
        if include_timestamp:
            timestamp = time.strftime('%Y%m%d_%H%M%S')
            parts.append(timestamp)
        
        # Join parts and add extension
        filename = '_'.join(parts) + f'.{format_type}'
        return filename
    
    @staticmethod
    def _clean_filename(name: str) -> str:
        """Clean filename by removing invalid characters."""
        # Remove or replace invalid filename characters
        invalid_chars = '<>:"/\\|?*'
        for char in invalid_chars:
            name = name.replace(char, '_')
        
        # Remove dots except for the extension
        name = name.replace('.', '_')
        
        # Limit length and remove extra underscores
        name = name[:50]  # Reasonable filename length
        name = '_'.join(part for part in name.split('_') if part)
        
        return name or 'unnamed'


class StorageManager:
    """High-level storage management for crawl-url application."""
    
    def __init__(self) -> None:
        """Initialize storage manager."""
        self.current_storage: Optional[URLStorage] = None
    
    def create_storage(
        self, 
        base_url: str, 
        format_type: str = 'txt',
        output_path: Optional[Path] = None,
        custom_suffix: str = ''
    ) -> URLStorage:
        """
        Create a new URL storage instance.
        
        Args:
            base_url: Base URL being crawled (for filename generation)
            format_type: Output format (txt, json, csv)
            output_path: Specific output path (overrides auto-generation)
            custom_suffix: Custom suffix for auto-generated filenames
            
        Returns:
            URLStorage instance
        """
        if output_path:
            filename = output_path
        else:
            # Auto-generate filename
            filename_str = FilenameGenerator.generate_filename(
                base_url=base_url,
                format_type=format_type,
                custom_suffix=custom_suffix
            )
            filename = Path.cwd() / filename_str
        
        self.current_storage = URLStorage(filename)
        return self.current_storage
    
    def save_urls(
        self, 
        urls: List[str], 
        base_url: str,
        format_type: str = 'txt',
        output_path: Optional[Path] = None,
        custom_suffix: str = ''
    ) -> Path:
        """
        Convenience method to save URLs directly.
        
        Args:
            urls: List of URLs to save
            base_url: Base URL being crawled
            format_type: Output format
            output_path: Specific output path
            custom_suffix: Custom suffix for filename
            
        Returns:
            Path where the file was saved
        """
        storage = self.create_storage(
            base_url=base_url,
            format_type=format_type,
            output_path=output_path,
            custom_suffix=custom_suffix
        )
        
        storage.add_urls(urls)
        storage.save_to_file(format_type)
        
        return storage.output_file
    
    def get_current_storage(self) -> Optional[URLStorage]:
        """Get the current storage instance."""
        return self.current_storage
</file>

<file path="src/crawl_url/core/__init__.py">
"""Core crawling and parsing functionality."""
</file>

<file path="README.md">
# crawl-url

A powerful Python CLI tool for extracting URLs from websites with intelligent mode switching and cross-platform support.

## ✨ Features

**🕷️ URL Discovery**: Extract URLs from websites using two specialized modes
- **Sitemap Mode**: Parse `sitemap.xml` files and sitemap indexes efficiently
- **Crawl Mode**: Recursive website crawling with configurable depth and rate limiting

**🖥️ Platform Support**: Works seamlessly across Windows, Linux, and macOS
- Automatic Windows console fallback (no PyTermGUI issues)
- Full PyTermGUI interface on Linux
- Unicode-safe terminal output everywhere

**⚙️ Configuration**: 
- CLI flags for quick one-liners
- Interactive TUI mode for guided usage
- Configurable rate limiting, depth limits, and URL filtering

## 🚀 Quick Start

### Installation
```bash
# Quick setup (Windows)
python setup.bat

# Quick setup (Unix/Linux/Mac)
./setup-advanced.bat

# Manual installation
pip install crawl-url

# From source
pip install -e ".[dev]"
```

### Basic Usage
```bash
# Extract URLs from any website
crawl-url crawl https://example.com

# Use sitemap mode (faster)
crawl-url crawl https://example.com --mode sitemap

# Customize depth and rate limiting
crawl-url crawl https://example.com --depth 5 --delay 1.0

# Save results
 crawl-url crawl https://example.com --format json > urls.json
```

## 📋 Commands

### CLI Mode
```bash
# Basic crawling
crawl-url crawl <URL> [OPTIONS]

# Interactive mode
crawl-url interactive

# Help and version
crawl-url --help
crawl-url --version
```

### CLI Options
```
--mode [auto|sitemap|crawl]     # Discovery mode (default: auto)
--depth INTEGER                 # Crawl depth 1-10 (default: 3)  
--delay FLOAT                   # Request delay in seconds (default: 1.0)
--format [txt|json|csv]         # Output format (default: txt)
--filter TEXT                   # Substring to filter URLs by
--user-agent TEXT              # Custom user agent string
--timeout INTEGER               # Request timeout in seconds
```

### Interactive Mode
For guided usage with prompts and progress indicators:
```bash
crawl-url interactive
```

## 🎯 Usage Examples

### Sitemap Extraction
```bash
# Extract from sitemap.xml
crawl-url crawl https://example.com --mode sitemap --format json

# Handle sitemap indexes automatically
crawl-url crawl https://example.com/sitemap_index.xml --mode sitemap
```

### Deep Crawling
```bash
# Crawl with depth customization
crawl-url crawl https://blog.example.com --depth 5 --delay 2.0

# Respect robots.txt while crawling
crawl-url crawl https://example.com --user-agent "crawl-url/1.0"
```

### URL Filtering
```bash
# Filter by domain or path
crawl-url crawl https://example.com --filter "/blog/"
crawl-url crawl https://example.com --filter ".pdf"
```

### Output Formats
```bash
# Plain text (default)
crawl-url crawl https://example.com > urls.txt

# JSON for programmatic use
crawl-url crawl https://example.com --format json | jq '.urls[]'

# CSV for spreadsheet analysis
crawl-url crawl https://example.com --format csv > urls.csv
```

## 🛠️ Development

### Setup Development Environment
```bash
# Clone repository
git clone <repository-url>
cd crawl-url

# Install development dependencies
python -m venv venv
# Windows: venv\Scripts\activate
# Unix: source venv/bin/activate
pip install -e ".[dev]"

# Run tests
pytest                    # All tests
pytest -m unit           # Unit tests only
pytest -v                # Verbose output
pytest --cov-report=html # Coverage report
```

### Code Quality
```bash
# Format and lint
black src/ tests/
ruff check src/ tests/
mypy src/

# Run all quality checks
ruff check src/ tests/ && black --check src/ tests/ && mypy src/ && pytest

# Build package
python -m build
```

### Testing
This project uses pytest with comprehensive test coverage:
- **Unit tests**: Fast, isolated tests with mocks
- **Integration tests**: End-to-end CLI and API testing
- **Platform tests**: Cross-platform compatibility validation
- **Coverage**: HTML reports available in `htmlcov/`

Run test suites:
```bash
pytest -m unit          # 100ms fast tests
pytest -m integration   # End-to-end workflows  
pytest -m "not slow"    # Skip network tests
```

## 🏗️ Architecture

```
crawl-url/
├── src/crawl_url/
│   ├── cli.py              # Typer CLI registration
│   ├── core/
│   │   ├── models.py       # Config/result validation
│   │   ├── crawler.py      # HTTP crawling + rate limiting
│   │   ├── sitemap_parser.py  # XML parsing
│   │   └── ui.py           # PyTermGUI + console fallback
│   └── utils/
│       ├── storage.py      # File I/O (TXT/JSON/CSV)
│       └── validation.py   # URL filtering utilities
└── tests/
    ├── test_*.py          # Unit/integration tests
    └── conftest.py        # Test fixtures and mocks
```

## 🌍 Platform Compatibility

### Windows
- Automatic console mode (PyTermGUI fallback)
- PowerShell/CMD compatibility
- Unicode-safe terminal output

### Linux & Unix
- Full PyTermGUI interface with ncurses
- Rich terminal features and colors
- Bash/Zsh completion

### macOS
- Native terminal support
- Homebrew installation ready

## 📊 Performance

**Sitemap Mode**: O(n) - Direct XML parsing
**Crawl Mode**: O(n²) worst case (depth×links) with configured limits

**Tuning recommends:**
- `--depth 3` for average sites (≤1000 URLs)
- `--delay 1.0` to respect server resources
- `--timeout 30` for slower connections

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-mode`
3. Run quality checks: `ruff check src/ tests/ && black src/ tests/ && mypy src/ && pytest`
4. Submit pull request with clear commit messages

### Development Setup
```bash
# Install pre-commit hooks
pre-commit install

# Run validation test
python final_validation.py

# Package testing
pip install -e .
crawl-url --version
crawl-url crawl https://example.com --depth 1
crawl-url interactive  # Test TUI
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🐛 Troubleshooting

**Permission Issues on Windows:**
```bash
# Use PowerShell as Administrator
python setup.bat
```

**Python Version Issues:**
- Requires Python 3.8+ (see pyproject.toml)
- Tested on Python 3.8-3.12

**Network Errors:**
- Check internet connectivity
- Verify robots.txt allows crawling
- Use longer `--timeout` for slow connections
</file>

<file path="src/crawl_url/core/sitemap_parser.py">
"""Sitemap parsing implementation for crawl-url application."""

import gzip
import logging
from io import BytesIO
from typing import Dict, List, Optional, Union
from urllib.parse import urljoin, urlparse

import lxml.etree as etree
import requests

from .models import CrawlResult, SitemapEntry


class SitemapParser:
    """Memory-efficient sitemap parser using lxml."""
    
    def __init__(self, session: Optional[requests.Session] = None) -> None:
        """Initialize sitemap parser with optional session."""
        self.session = session or requests.Session()
        self.session.headers.update({
            'User-Agent': 'crawl-url/1.0 (Sitemap Parser)',
            'Accept': 'application/xml, text/xml, */*',
            'Accept-Encoding': 'gzip, deflate'
        })
        
        # Configure retry strategy for better reliability
        from requests.adapters import HTTPAdapter
        from urllib3.util.retry import Retry
        
        retry_strategy = Retry(
            total=2,
            backoff_factor=0.5,
            status_forcelist=[429, 500, 502, 503, 504],
        )
        
        adapter = HTTPAdapter(max_retries=retry_strategy)
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)
        
    def parse_sitemap_efficiently(
        self, 
        file_source: Union[str, BytesIO], 
        is_compressed: bool = False
    ) -> List[SitemapEntry]:
        """
        Memory-efficient sitemap parsing using iterparse.
        
        Args:
            file_source: File path, URL, or file-like object
            is_compressed: Whether the file is gzip compressed
            
        Returns:
            List of SitemapEntry objects
        """
        try:
            # Handle different source types
            if isinstance(file_source, str):
                if file_source.startswith(('http://', 'https://')):
                    file_obj = self._fetch_sitemap_content(file_source, timeout=10)
                else:
                    file_obj = self._open_local_file(file_source, is_compressed)
            else:
                file_obj = file_source
            
            entries = []
            context = etree.iterparse(file_obj, events=('start', 'end'))
            
            # Skip root element
            context = iter(context)
            event, root = next(context)
            
            current_url = {}
            
            for event, element in context:
                if event == 'end':
                    tag_name = element.tag.split('}')[-1]  # Remove namespace
                    
                    if tag_name == 'url':
                        # Create entry from collected data
                        if current_url.get('loc'):
                            entry = SitemapEntry(
                                loc=current_url.get('loc', ''),
                                lastmod=current_url.get('lastmod'),
                                changefreq=current_url.get('changefreq'),
                                priority=current_url.get('priority')
                            )
                            entries.append(entry)
                        current_url = {}
                        
                    elif tag_name in ('loc', 'lastmod', 'changefreq', 'priority'):
                        if element.text:
                            current_url[tag_name] = element.text.strip()
                    
                    # Memory cleanup - critical for large files
                    element.clear()
                    while element.getprevious() is not None:
                        del element.getparent()[0]
            
            if hasattr(file_obj, 'close'):
                file_obj.close()
            return entries
            
        except etree.XMLSyntaxError as e:
            logging.error(f"XML syntax error in sitemap: {e}")
            return []
        except Exception as e:
            logging.error(f"Error parsing sitemap: {e}")
            return []

    def parse_sitemap_index(self, index_source: Union[str, BytesIO]) -> List[Dict[str, str]]:
        """Parse sitemap index and return list of sitemap URLs."""
        try:
            if isinstance(index_source, str):
                if index_source.startswith(('http://', 'https://')):
                    file_obj = self._fetch_sitemap_content(index_source, timeout=10)
                else:
                    file_obj = self._open_local_file(index_source, False)
            else:
                # It's already a file-like object (BytesIO)
                file_obj = index_source
            
            sitemaps = []
            context = etree.iterparse(file_obj, events=('start', 'end'))
            current_sitemap = {}
            
            for event, element in context:
                if event == 'end':
                    tag_name = element.tag.split('}')[-1]
                    
                    if tag_name == 'sitemap':
                        if current_sitemap.get('loc'):
                            sitemaps.append(current_sitemap.copy())
                        current_sitemap = {}
                    elif tag_name in ('loc', 'lastmod'):
                        if element.text:
                            current_sitemap[tag_name] = element.text.strip()
                    
                    # Memory cleanup
                    element.clear()
                    while element.getprevious() is not None:
                        del element.getparent()[0]
            
            if hasattr(file_obj, 'close'):
                file_obj.close()
            return sitemaps
            
        except Exception as e:
            logging.error(f"Error parsing sitemap index: {e}")
            return []

    def discover_sitemaps(self, base_url: str) -> List[str]:
        """Discover sitemap URLs from robots.txt and common locations."""
        discovered_sitemaps = []
        
        # 1. Check robots.txt
        robots_sitemaps = self._get_sitemaps_from_robots(base_url)
        discovered_sitemaps.extend(robots_sitemaps)
        
        # 2. Check common locations if nothing found in robots.txt
        if not discovered_sitemaps:
            common_locations = [
                '/sitemap.xml',
                '/sitemap_index.xml',
                '/sitemaps/sitemap.xml',
                '/xml-sitemaps/sitemap.xml'
            ]
            
            for location in common_locations:
                sitemap_url = urljoin(base_url, location)
                if self._check_sitemap_exists(sitemap_url):
                    discovered_sitemaps.append(sitemap_url)
                    break  # Stop at first found
        
        return discovered_sitemaps

    def _fetch_sitemap_content(self, url: str, timeout: int = 15) -> BytesIO:
        """Fetch sitemap content from URL with compression handling."""
        response = self.session.get(url, timeout=timeout)
        response.raise_for_status()
        
        # The requests library automatically decompresses gzip content when
        # Content-Encoding: gzip is present, so we don't need to decompress again
        # Only manually decompress if the URL ends with .gz (indicating a gzipped file)
        if url.endswith('.gz'):
            try:
                # Try to decompress if it's actually gzipped
                content = gzip.decompress(response.content)
                return BytesIO(content)
            except gzip.BadGzipFile:
                # If decompression fails, use content as-is (already decompressed by requests)
                return BytesIO(response.content)
        else:
            # For regular XML files, use content directly
            return BytesIO(response.content)

    def _open_local_file(self, file_path: str, is_compressed: bool) -> Union[gzip.GzipFile, open]:
        """Open local file with optional compression handling."""
        if is_compressed or file_path.endswith('.gz'):
            return gzip.open(file_path, 'rb')
        else:
            return open(file_path, 'rb')

    def _get_sitemaps_from_robots(self, base_url: str) -> List[str]:
        """Extract sitemap URLs from robots.txt."""
        robots_url = urljoin(base_url, '/robots.txt')
        sitemaps = []
        
        try:
            response = self.session.get(robots_url, timeout=10)
            if response.status_code == 200:
                for line in response.text.split('\n'):
                    line = line.strip()
                    if line.lower().startswith('sitemap:'):
                        sitemap_url = line.split(':', 1)[1].strip()
                        sitemaps.append(sitemap_url)
        except Exception as e:
            logging.warning(f"Could not fetch robots.txt from {robots_url}: {e}")
        
        return sitemaps

    def _check_sitemap_exists(self, sitemap_url: str) -> bool:
        """Check if sitemap exists at URL."""
        try:
            response = self.session.head(sitemap_url, timeout=10)
            return response.status_code == 200
        except Exception:
            return False


class SitemapService:
    """Complete sitemap processing service for crawl-url application."""
    
    def __init__(self, progress_callback: Optional[callable] = None) -> None:
        """Initialize sitemap service with optional progress callback."""
        self.parser = SitemapParser()
        self.progress_callback = progress_callback
        
    def process_sitemap_url(self, sitemap_url: str, filter_base: Optional[str] = None) -> CrawlResult:
        """
        Process sitemap URL and return all URLs.
        Handles both regular sitemaps and sitemap indexes.
        """
        try:
            # First, try to determine if it's an index or regular sitemap
            if self._is_sitemap_index(sitemap_url):
                return self._process_sitemap_index(sitemap_url, filter_base)
            else:
                return self._process_single_sitemap(sitemap_url, filter_base)
                
        except Exception as e:
            return CrawlResult(
                success=False,
                urls=[],
                count=0,
                message=f'Error processing sitemap: {str(e)}',
                errors=[str(e)]
            )
    
    def process_base_url(self, base_url: str, filter_base: Optional[str] = None) -> CrawlResult:
        """Discover and process sitemaps from base URL."""
        try:
            if self.progress_callback:
                self.progress_callback("Discovering sitemaps...", 0)
            
            discovered_sitemaps = self.parser.discover_sitemaps(base_url)
            
            if not discovered_sitemaps:
                return CrawlResult(
                    success=False,
                    urls=[],
                    count=0,
                    message='No sitemaps found for this domain'
                )
            
            all_urls = []
            for i, sitemap_url in enumerate(discovered_sitemaps):
                if self.progress_callback:
                    self.progress_callback(
                        f"Processing sitemap {i+1}/{len(discovered_sitemaps)}", 
                        len(all_urls)
                    )
                
                result = self.process_sitemap_url(sitemap_url, filter_base)
                if result.success:
                    all_urls.extend(result.urls)
            
            return CrawlResult(
                success=True,
                urls=all_urls,
                count=len(all_urls),
                message=f'Successfully extracted {len(all_urls)} URLs from {len(discovered_sitemaps)} sitemaps'
            )
            
        except Exception as e:
            return CrawlResult(
                success=False,
                urls=[],
                count=0,
                message=f'Error processing base URL: {str(e)}',
                errors=[str(e)]
            )
    
    def _is_sitemap_index(self, sitemap_url: str) -> bool:
        """Quickly check if sitemap is an index by sampling content."""
        try:
            # Fetch first few KB to check
            response = self.parser.session.get(
                sitemap_url, 
                headers={'Range': 'bytes=0-2048'}, 
                timeout=10
            )
            content = response.content.decode('utf-8', errors='ignore')
            return '<sitemapindex' in content.lower()
        except Exception:
            return False
    
    def _process_sitemap_index(self, index_url: str, filter_base: Optional[str]) -> CrawlResult:
        """Process sitemap index and all contained sitemaps."""
        try:
            sitemaps = self.parser.parse_sitemap_index(index_url)
            if not sitemaps:
                return CrawlResult(
                    success=False,
                    urls=[],
                    count=0,
                    message='No sitemaps found in sitemap index',
                    errors=[]
                )
            
            all_urls = []
            errors = []
            successful_sitemaps = 0
            
            for i, sitemap_info in enumerate(sitemaps):
                if self.progress_callback:
                    self.progress_callback(
                        f"Processing sitemap {i+1}/{len(sitemaps)}", 
                        len(all_urls)
                    )
                
                try:
                    sitemap_url = sitemap_info['loc']
                    logging.info(f"Processing sitemap: {sitemap_url}")
                    
                    # Parse the individual sitemap with timeout handling
                    entries = self.parser.parse_sitemap_efficiently(sitemap_url)
                    urls = self._extract_and_filter_urls(entries, filter_base)
                    
                    if urls:
                        all_urls.extend(urls)
                        successful_sitemaps += 1
                        logging.info(f"Successfully processed {sitemap_url}: {len(urls)} URLs")
                    else:
                        logging.warning(f"No URLs found in sitemap: {sitemap_url}")
                        
                except Exception as e:
                    error_msg = f"Error processing sitemap {sitemap_info.get('loc', 'unknown')}: {str(e)}"
                    errors.append(error_msg)
                    logging.error(error_msg)
            
            # Consider it successful if we got at least some URLs or processed some sitemaps
            success = len(all_urls) > 0 or successful_sitemaps > 0
            
            return CrawlResult(
                success=success,
                urls=all_urls,
                count=len(all_urls),
                message=f'Processed {successful_sitemaps}/{len(sitemaps)} sitemaps from index',
                errors=errors
            )
            
        except Exception as e:
            return CrawlResult(
                success=False,
                urls=[],
                count=0,
                message=f'Error processing sitemap index: {str(e)}',
                errors=[str(e)]
            )
    
    def _process_single_sitemap(self, sitemap_url: str, filter_base: Optional[str]) -> CrawlResult:
        """Process a single sitemap file."""
        try:
            entries = self.parser.parse_sitemap_efficiently(sitemap_url)
            urls = self._extract_and_filter_urls(entries, filter_base)
            
            return CrawlResult(
                success=True,
                urls=urls,
                count=len(urls),
                message=f'Extracted {len(urls)} URLs from sitemap'
            )
        except Exception as e:
            return CrawlResult(
                success=False,
                urls=[],
                count=0,
                message=f'Error processing sitemap: {str(e)}',
                errors=[str(e)]
            )
    
    def _extract_and_filter_urls(self, entries: List[SitemapEntry], filter_base: Optional[str]) -> List[str]:
        """Extract URLs from entries and apply filtering."""
        urls = []
        for entry in entries:
            if entry.loc:
                if filter_base:
                    if entry.loc.startswith(filter_base):
                        urls.append(entry.loc)
                else:
                    urls.append(entry.loc)
        return urls
</file>

<file path="src/crawl_url/core/crawler.py">
"""Web crawling implementation for crawl-url application."""

import hashlib
import logging
import time
import urllib.robotparser
from collections import defaultdict, deque
from typing import Dict, List, Optional, Set
from urllib.parse import urljoin, urlparse, urlunparse

import requests
from bs4 import BeautifulSoup
from requests.adapters import HTTPAdapter
from requests.exceptions import ConnectionError, RequestException, Timeout
from urllib3.util.retry import Retry

from .models import CrawlResult


class URLExtractor:
    """Extract and validate URLs from web pages."""
    
    def __init__(self, allowed_domains: Optional[Set[str]] = None, url_filter_base: Optional[str] = None) -> None:
        """Initialize URL extractor with filtering options."""
        self.allowed_domains = set(allowed_domains or [])
        self.url_filter_base = url_filter_base
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'crawl-url/1.0 (Compatible Web Crawler)'
        })
        
        # Configure retry strategy
        retry_strategy = Retry(
            total=3,
            backoff_factor=1,
            status_forcelist=[429, 500, 502, 503, 504],
        )
        
        adapter = HTTPAdapter(max_retries=retry_strategy)
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)
    
    def extract_urls_from_page(self, url: str, timeout: int = 10) -> Set[str]:
        """Extract all URLs from a single page."""
        try:
            response = self.session.get(url, timeout=timeout)
            response.raise_for_status()
            
            if not self._is_html_content(response):
                return set()
            
            soup = BeautifulSoup(response.content, 'html.parser')
            urls = set()
            
            # Extract from anchor tags
            for link in soup.find_all('a', href=True):
                absolute_url = urljoin(url, link['href'])
                if self._is_valid_url(absolute_url):
                    urls.add(self._normalize_url(absolute_url))
            
            return urls
            
        except requests.RequestException as e:
            logging.warning(f"Error fetching {url}: {e}")
            return set()
    
    def _is_html_content(self, response: requests.Response) -> bool:
        """Check if response contains HTML content."""
        content_type = response.headers.get('content-type', '').lower()
        return 'text/html' in content_type
    
    def _is_valid_url(self, url: str) -> bool:
        """Validate and filter URLs."""
        try:
            parsed = urlparse(url)
            
            # Basic validation
            if not (parsed.scheme and parsed.netloc):
                return False
            
            # Domain filtering
            if self.allowed_domains and parsed.netloc not in self.allowed_domains:
                return False
            
            # Base URL filtering
            if self.url_filter_base:
                if not url.startswith(self.url_filter_base):
                    return False
            
            return True
            
        except Exception:
            return False
    
    def _normalize_url(self, url: str) -> str:
        """Normalize URL by removing fragments."""
        parsed = urlparse(url)
        # Remove fragment (anchor)
        normalized = urlunparse((
            parsed.scheme,
            parsed.netloc, 
            parsed.path,
            parsed.params,
            parsed.query,
            ''  # Remove fragment
        ))
        return normalized


class RateLimiter:
    """Implement respectful rate limiting for web crawling."""
    
    def __init__(self, default_delay: float = 1.0) -> None:
        """Initialize rate limiter with default delay."""
        self.default_delay = default_delay
        self.domain_delays: Dict[str, float] = defaultdict(lambda: default_delay)
        self.last_request_time: Dict[str, float] = defaultdict(float)
    
    def wait_if_needed(self, url: str) -> None:
        """Implement respectful delay based on domain."""
        domain = urlparse(url).netloc
        current_time = time.time()
        time_since_last = current_time - self.last_request_time[domain]
        
        delay_needed = self.domain_delays[domain] - time_since_last
        if delay_needed > 0:
            time.sleep(delay_needed)
        
        self.last_request_time[domain] = time.time()
    
    def set_domain_delay(self, domain: str, delay: float) -> None:
        """Set custom delay for specific domain."""
        self.domain_delays[domain] = delay


class RobotsTxtChecker:
    """Check robots.txt compliance for URLs."""
    
    def __init__(self, user_agent: str = '*') -> None:
        """Initialize robots.txt checker."""
        self.user_agent = user_agent
        self.robots_cache: Dict[str, Optional[urllib.robotparser.RobotFileParser]] = {}
    
    def can_fetch(self, url: str) -> bool:
        """Check if URL can be fetched according to robots.txt."""
        try:
            parsed_url = urlparse(url)
            base_url = f"{parsed_url.scheme}://{parsed_url.netloc}"
            
            if base_url not in self.robots_cache:
                robots_url = urljoin(base_url, '/robots.txt')
                rp = urllib.robotparser.RobotFileParser()
                rp.set_url(robots_url)
                try:
                    rp.read()
                    self.robots_cache[base_url] = rp
                except Exception:
                    # If robots.txt can't be read, allow crawling
                    self.robots_cache[base_url] = None
            
            robots_parser = self.robots_cache[base_url]
            if robots_parser is None:
                return True
            
            # Fix for Python robotparser bug: when there are no Disallow rules
            # for a user-agent, it should default to allowing all URLs
            try:
                result = robots_parser.can_fetch(self.user_agent, url)
                
                # Check if this is the robotparser bug where no rules exist
                # but it returns False anyway
                if not result:
                    # Check if there are any disallow rules for this user agent
                    has_disallow_rules = False
                    if robots_parser.default_entry:
                        for rule in robots_parser.default_entry.rulelines:
                            if not rule.allowance:  # Disallow rule
                                has_disallow_rules = True
                                break
                    
                    # If no disallow rules exist, allow access
                    if not has_disallow_rules:
                        return True
                
                return result
            except Exception:
                # If there's any error with robots parsing, be permissive
                return True
                
        except Exception:
            return True  # Allow if robots.txt can't be checked


class URLDeduplicator:
    """Memory-efficient URL deduplication using hashes."""
    
    def __init__(self, max_memory_hashes: int = 100000) -> None:
        """Initialize URL deduplicator."""
        self.url_hashes: Set[str] = set()
        self.max_memory_hashes = max_memory_hashes
        
    def is_duplicate(self, url: str) -> bool:
        """Check if URL has been seen before using memory-efficient hashing."""
        url_hash = hashlib.md5(url.encode('utf-8')).hexdigest()
        
        if url_hash in self.url_hashes:
            return True
            
        # Prevent unlimited memory growth
        if len(self.url_hashes) >= self.max_memory_hashes:
            # Remove oldest 20% of hashes (simple cleanup)
            hashes_list = list(self.url_hashes)
            self.url_hashes = set(hashes_list[len(hashes_list)//5:])
        
        self.url_hashes.add(url_hash)
        return False


class WebCrawler:
    """Recursive web crawler with respectful crawling policies."""
    
    def __init__(
        self, 
        max_depth: int = 3, 
        delay: float = 1.0, 
        max_urls: int = 1000,
        progress_callback: Optional[callable] = None
    ) -> None:
        """Initialize web crawler with configuration."""
        self.max_depth = max_depth
        self.delay = delay
        self.max_urls = max_urls
        self.progress_callback = progress_callback
        
        self.visited_urls: Set[str] = set()
        self.url_queue: deque = deque()
        self.rate_limiter = RateLimiter(delay)
        self.extractor = URLExtractor()
        self.robots_checker = RobotsTxtChecker()
        self.deduplicator = URLDeduplicator()
        
        # Setup logging
        self.logger = logging.getLogger(__name__)
        
    def crawl_website(self, start_url: str, filter_base: Optional[str] = None) -> CrawlResult:
        """Main crawling method with progress reporting."""
        try:
            # Configure extractor with filter
            if filter_base:
                self.extractor.url_filter_base = filter_base
            
            self.url_queue.append((start_url, 0))  # (url, depth)
            all_urls = []
            errors = []
            
            while self.url_queue and len(all_urls) < self.max_urls:
                current_url, depth = self.url_queue.popleft()
                
                if (current_url in self.visited_urls or 
                    depth > self.max_depth or
                    self.deduplicator.is_duplicate(current_url)):
                    continue
                
                # Progress callback for UI updates
                if self.progress_callback:
                    self.progress_callback(f"Crawling: {self._shorten_url(current_url)}", len(all_urls))
                
                # Respectful crawling checks
                if not self.robots_checker.can_fetch(current_url):
                    self.logger.info(f"Robots.txt disallows crawling: {current_url}")
                    continue
                    
                self.rate_limiter.wait_if_needed(current_url)
                
                # Extract URLs from current page
                try:
                    urls = self.extractor.extract_urls_from_page(current_url)
                    all_urls.extend(list(urls))
                    self.visited_urls.add(current_url)
                    
                    # Add new URLs to queue for next depth level
                    if depth < self.max_depth:
                        for url in list(urls)[:50]:  # Limit to prevent explosion
                            if url not in self.visited_urls:
                                self.url_queue.append((url, depth + 1))
                                
                except Exception as e:
                    error_msg = f"Error crawling {current_url}: {e}"
                    errors.append(error_msg)
                    self.logger.error(error_msg)
            
            return CrawlResult(
                success=True,
                urls=all_urls,
                count=len(all_urls),
                message=f'Successfully crawled {len(all_urls)} URLs',
                errors=errors
            )
            
        except Exception as e:
            return CrawlResult(
                success=False,
                urls=[],
                count=0,
                message=f'Crawling failed: {str(e)}',
                errors=[str(e)]
            )
    
    def _shorten_url(self, url: str, max_length: int = 50) -> str:
        """Shorten URL for display purposes."""
        if len(url) <= max_length:
            return url
        return url[:max_length-3] + "..."


class CrawlerService:
    """Service class to integrate web crawler with terminal app."""
    
    def __init__(self, progress_callback: Optional[callable] = None) -> None:
        """Initialize crawler service."""
        self.progress_callback = progress_callback
        
    def crawl_url(
        self, 
        url: str, 
        max_depth: int = 3, 
        delay: float = 1.0, 
        filter_base: Optional[str] = None
    ) -> CrawlResult:
        """Main crawling service method."""
        try:
            crawler = WebCrawler(
                max_depth=max_depth,
                delay=delay,
                progress_callback=self.progress_callback
            )
            
            return crawler.crawl_website(url, filter_base)
            
        except Exception as e:
            return CrawlResult(
                success=False,
                urls=[],
                count=0,
                message=f'Crawling service failed: {str(e)}',
                errors=[str(e)]
            )
</file>

<file path="src/crawl_url/utils/validation.py">
"""Validation utilities with user-friendly error messages."""

import re
from typing import Tuple
from urllib.parse import urlparse


def validate_url(url: str) -> Tuple[bool, str]:
    """
    Validate URL format and provide helpful error messages.
    
    Args:
        url: URL string to validate
        
    Returns:
        Tuple of (is_valid, error_message). error_message is empty string if valid.
    """
    if not url or not url.strip():
        return False, "🌐 URL cannot be empty. Please enter a website URL (e.g., https://example.com)"
    
    url = url.strip()
    
    # Check for common protocol mistakes
    if not url.startswith(('http://', 'https://')):
        if url.startswith(('www.', 'ftp.', 'mail.')):
            return False, (
                f"🔗 Missing protocol. Did you mean 'https://{url}'?\n"
                f"URLs must start with 'http://' or 'https://'"
            )
        elif '.' in url and not url.startswith(('mailto:', 'javascript:', 'file:')):
            return False, (
                f"🔗 Missing protocol. Did you mean 'https://{url}'?\n"
                f"URLs must start with 'http://' or 'https://'"
            )
        else:
            return False, (
                "🔗 Invalid URL format. URLs must start with 'http://' or 'https://'\n"
                "Examples: https://example.com, http://subdomain.site.org/path"
            )
    
    # Parse URL to validate structure
    try:
        parsed = urlparse(url)
        
        if not parsed.netloc:
            return False, (
                "🌐 Invalid URL: Missing domain name\n"
                "Example of valid URL: https://example.com/path"
            )
        
        # Check for common mistakes
        if parsed.netloc.startswith('.') or parsed.netloc.endswith('.'):
            return False, (
                f"🌐 Invalid domain '{parsed.netloc}': Domain cannot start or end with a dot\n"
                "Example: https://example.com (not https://.example.com.)"
            )
        
        if '..' in parsed.netloc:
            return False, (
                f"🌐 Invalid domain '{parsed.netloc}': Domain cannot contain consecutive dots\n"
                "Example: https://sub.example.com (not https://sub..example.com)"
            )
        
        # Check for suspicious protocols in parsed URL
        if parsed.scheme.lower() not in ('http', 'https'):
            return False, (
                f"🔒 Unsupported protocol '{parsed.scheme}'. Only HTTP and HTTPS are supported\n"
                "Please use URLs starting with 'http://' or 'https://'"
            )
        
        return True, ""
        
    except Exception as e:
        return False, (
            f"🔗 Invalid URL format: {str(e)}\n"
            "Please check your URL and try again. Example: https://example.com"
        )


def validate_crawl_depth(depth_str: str) -> Tuple[bool, str, int]:
    """
    Validate crawl depth input with helpful error messages.
    
    Args:
        depth_str: String representation of depth
        
    Returns:
        Tuple of (is_valid, error_message, parsed_depth)
    """
    if not depth_str or not depth_str.strip():
        return False, "🔍 Crawl depth cannot be empty. Please enter a number between 1 and 10", 0
    
    depth_str = depth_str.strip()
    
    try:
        depth = int(depth_str)
        
        if depth < 1:
            return False, (
                f"🔍 Crawl depth '{depth}' is too low. Minimum depth is 1\n"
                "Depth 1 means crawling only the starting page"
            ), 0
        
        if depth > 10:
            return False, (
                f"🔍 Crawl depth '{depth}' is too high. Maximum depth is 10\n"
                "High depths can take very long and may overwhelm websites"
            ), 0
        
        # Provide guidance for different depth values
        guidance = ""
        if depth == 1:
            guidance = " (Quick: crawls only the starting page)"
        elif depth <= 3:
            guidance = " (Recommended: good balance of speed and coverage)"
        elif depth <= 5:
            guidance = " (Thorough: may take longer but finds more URLs)"
        else:
            guidance = " (Deep crawl: may be very slow, use with caution)"
        
        return True, guidance, depth
        
    except ValueError:
        if '.' in depth_str:
            return False, (
                f"🔍 Crawl depth '{depth_str}' must be a whole number, not a decimal\n"
                "Examples: 1, 2, 3 (not 1.5, 2.0)"
            ), 0
        else:
            return False, (
                f"🔍 Invalid crawl depth '{depth_str}'. Please enter a whole number between 1 and 10\n"
                "Examples: 1 (surface), 3 (recommended), 5 (thorough)"
            ), 0


def validate_delay(delay_str: str) -> Tuple[bool, str, float]:
    """
    Validate delay input with helpful error messages.
    
    Args:
        delay_str: String representation of delay in seconds
        
    Returns:
        Tuple of (is_valid, error_message, parsed_delay)
    """
    if not delay_str or not delay_str.strip():
        return False, "⏱️ Delay cannot be empty. Please enter a number (e.g., 1.0 for 1 second)", 0.0
    
    delay_str = delay_str.strip()
    
    try:
        delay = float(delay_str)
        
        if delay < 0:
            return False, (
                f"⏱️ Delay '{delay}' cannot be negative\n"
                "Use 0 for no delay, or positive numbers like 1.0 for delays"
            ), 0.0
        
        if delay < 0.5 and delay > 0:
            return False, (
                f"⏱️ Delay '{delay}' is too short. Minimum recommended delay is 0.5 seconds\n"
                "Short delays may overwhelm websites and get you blocked"
            ), 0.0
        
        if delay > 10:
            return False, (
                f"⏱️ Delay '{delay}' is very long. Consider using a shorter delay\n"
                "Delays over 10 seconds will make crawling extremely slow"
            ), 0.0
        
        # Provide guidance for different delay values
        guidance = ""
        if delay == 0:
            guidance = " (No delay: fastest but may overwhelm servers)"
        elif delay < 1:
            guidance = " (Fast: minimal delay)"
        elif delay <= 2:
            guidance = " (Recommended: respectful crawling speed)"
        else:
            guidance = " (Conservative: very respectful but slower)"
        
        return True, guidance, delay
        
    except ValueError:
        return False, (
            f"⏱️ Invalid delay '{delay_str}'. Please enter a number\n"
            "Examples: 1.0 (1 second), 0.5 (half second), 2.5 (2.5 seconds)"
        ), 0.0


def validate_filter_url(filter_url: str, base_url: str) -> Tuple[bool, str]:
    """
    Validate URL filter with helpful error messages.
    
    Args:
        filter_url: Filter URL string
        base_url: Base URL being crawled
        
    Returns:
        Tuple of (is_valid, error_message). error_message is empty string if valid.
    """
    if not filter_url or not filter_url.strip():
        return True, ""  # Empty filter is valid (no filtering)
    
    filter_url = filter_url.strip()
    
    # Validate filter URL format
    is_valid, error = validate_url(filter_url)
    if not is_valid:
        return False, f"🎯 Filter URL is invalid: {error}"
    
    # Check if filter is related to base URL (helpful warning)
    try:
        base_domain = urlparse(base_url).netloc
        filter_domain = urlparse(filter_url).netloc
        
        if base_domain and filter_domain and base_domain != filter_domain:
            return False, (
                f"🎯 Filter domain '{filter_domain}' doesn't match crawl domain '{base_domain}'\n"
                f"This will likely result in no URLs being found.\n"
                f"Example valid filter for {base_url}: {base_url}/docs/"
            )
        
        return True, ""
        
    except Exception:
        # If parsing fails, still allow the filter (let the crawler handle it)
        return True, ""


def suggest_url_fix(invalid_url: str) -> str:
    """
    Suggest a corrected URL for common user mistakes.
    
    Args:
        invalid_url: The invalid URL string
        
    Returns:
        Suggested corrected URL or empty string if no obvious fix
    """
    if not invalid_url:
        return ""
    
    invalid_url = invalid_url.strip()
    
    # Common case: missing protocol
    if not invalid_url.startswith(('http://', 'https://', 'ftp://', 'mailto:')):
        if '.' in invalid_url and not invalid_url.startswith(('.', '/')):
            return f"https://{invalid_url}"
    
    # Common case: mixed up protocols
    if invalid_url.startswith('htp://') or invalid_url.startswith('htps://'):
        return invalid_url.replace('htp://', 'http://').replace('htps://', 'https://')
    
    if invalid_url.startswith('http:///') or invalid_url.startswith('https:///'):
        return invalid_url.replace(':///', '://')
    
    return ""


def get_crawl_mode_explanation(mode: str) -> str:
    """
    Get detailed explanation of a crawling mode.
    
    Args:
        mode: Crawling mode ('auto', 'sitemap', 'crawl')
        
    Returns:
        Detailed explanation string
    """
    explanations = {
        'auto': (
            "🔍 Auto mode automatically detects the best crawling method:\n"
            "  • If URL ends with .xml → Uses sitemap mode\n"
            "  • Otherwise → Uses website crawling mode\n"
            "  • Best choice when you're unsure"
        ),
        'sitemap': (
            "🗺️ Sitemap mode extracts URLs from XML sitemaps:\n"
            "  • Fastest method for supported websites\n"
            "  • Finds URLs listed in sitemap.xml files\n"
            "  • Automatically discovers sitemaps from robots.txt\n"
            "  • Use when you know the site has sitemaps"
        ),
        'crawl': (
            "🕷️ Crawl mode recursively explores website pages:\n"
            "  • Follows links from page to page\n"
            "  • Can find URLs not in sitemaps\n"
            "  • Respects robots.txt and rate limits\n"
            "  • Use when sitemaps aren't available"
        )
    }
    
    return explanations.get(mode, f"Unknown mode: {mode}")


def get_output_format_explanation(format_type: str) -> str:
    """
    Get detailed explanation of an output format.
    
    Args:
        format_type: Output format ('txt', 'json', 'csv')
        
    Returns:
        Detailed explanation string
    """
    explanations = {
        'txt': (
            "📄 Plain text format (.txt):\n"
            "  • One URL per line\n"
            "  • Simple and lightweight\n"
            "  • Easy to process with scripts\n"
            "  • Best for: Simple URL lists"
        ),
        'json': (
            "📊 JSON format (.json):\n"
            "  • Structured data with metadata\n"
            "  • Includes crawl date, count, base URL\n"
            "  • Machine-readable format\n"
            "  • Best for: API integration, data analysis"
        ),
        'csv': (
            "📈 CSV format (.csv):\n"
            "  • Comma-separated values\n"
            "  • Columns: URL, Domain, Path\n"
            "  • Opens in Excel/Google Sheets\n"
            "  • Best for: Spreadsheet analysis, sorting"
        )
    }
    
    return explanations.get(format_type, f"Unknown format: {format_type}")
</file>

<file path="src/crawl_url/cli.py">
"""Main CLI entry point for crawl-url application."""

import sys
from pathlib import Path
from typing import List, Optional
from urllib.parse import urlparse

import typer
from rich.console import Console
from rich.progress import Progress, SpinnerColumn, TextColumn, BarColumn, TaskProgressColumn
from rich.table import Table

from . import __version__, __description__
from .core.crawler import CrawlerService
from .core.models import CrawlConfig
from .core.sitemap_parser import SitemapService
from .core.ui import CrawlerApp
from .utils.storage import StorageManager

# Create main Typer app
app = typer.Typer(
    name="crawl-url",
    help=f"🕷️ {__description__}",
    epilog="Made with ❤️ using Typer and PyTermGUI",
    add_completion=True,
    rich_markup_mode="rich"
)

console = Console()


def version_callback(value: bool) -> None:
    """Show version and exit."""
    if value:
        console.print(f"[bold blue]crawl-url[/bold blue] version [green]{__version__}[/green]")
        raise typer.Exit()


@app.callback()
def main_callback(
    version: Optional[bool] = typer.Option(
        None, 
        "--version", 
        "-V", 
        callback=version_callback,
        help="Show version and exit"
    )
) -> None:
    """🕷️ Crawl-URL: A powerful terminal application for URL crawling."""
    pass


@app.command()
def interactive() -> None:
    """🖥️ Launch interactive terminal UI mode."""
    try:
        crawler_app = CrawlerApp()
        crawler_app.run()
    except KeyboardInterrupt:
        console.print("\n[yellow]Operation cancelled by user[/yellow]")
        raise typer.Exit(0)
    except Exception as e:
        console.print(f"[red]Error launching interactive mode: {e}[/red]")
        raise typer.Exit(1)


@app.command()
def crawl(
    url: str = typer.Argument(
        ..., 
        help="🌐 URL to crawl (website URL or sitemap.xml URL)"
    ),
    mode: str = typer.Option(
        "auto", 
        "--mode", 
        "-m", 
        help="🔍 Crawling mode: auto (detect), sitemap (XML only), crawl (recursive)",
        case_sensitive=False
    ),
    output: Optional[Path] = typer.Option(
        None, 
        "--output", 
        "-o", 
        help="💾 Output file path (auto-generated if not specified)"
    ),
    format_type: str = typer.Option(
        "txt", 
        "--format", 
        "-f", 
        help="📄 Output format",
        case_sensitive=False
    ),
    depth: int = typer.Option(
        3, 
        "--depth", 
        "-d", 
        help="🔍 Maximum crawling depth (crawl mode only)",
        min=1,
        max=10
    ),
    filter_base: Optional[str] = typer.Option(
        None,
        "--filter",
        "-fb",
        help="🎯 Filter URLs by base URL (only URLs starting with this will be included)"
    ),
    delay: float = typer.Option(
        1.0, 
        "--delay", 
        help="⏱️ Delay between requests in seconds (minimum 0.5 for respectful crawling)",
        min=0.5
    ),
    max_urls: int = typer.Option(
        1000,
        "--max-urls",
        help="📊 Maximum number of URLs to extract",
        min=1,
        max=10000
    ),
    verbose: bool = typer.Option(
        False, 
        "--verbose", 
        "-v", 
        help="🔊 Enable verbose output with progress information"
    ),
) -> None:
    """🕷️ Crawl a URL and extract URLs (command-line mode)."""
    
    # Create configuration
    try:
        config = CrawlConfig(
            url=url,
            mode=mode,
            max_depth=depth,
            delay=delay,
            filter_base=filter_base,
            output_path=output,
            output_format=format_type,
            verbose=verbose
        )
    except ValueError as e:
        console.print(f"[red]Configuration error: {e}[/red]")
        raise typer.Exit(1)
    
    # Show initial information
    if verbose:
        _display_crawl_info(config, max_urls)
    
    # Determine output filename if not provided
    if output is None:
        domain = urlparse(url).netloc or "crawl_results"
        timestamp = __import__('time').strftime('%Y%m%d_%H%M%S')
        output = Path(f"{domain}_{timestamp}.{format_type}")
    
    try:
        # Perform crawling based on mode
        if mode == "sitemap" or (mode == "auto" and url.endswith('.xml')):
            result = _crawl_sitemap_mode(config, verbose)
        else:
            result = _crawl_website_mode(config, max_urls, verbose)
        
        if result.success:
            # Save results
            storage_manager = StorageManager()
            final_output_path = storage_manager.save_urls(
                urls=result.urls,
                base_url=config.url,
                format_type=config.output_format,
                output_path=output
            )
            
            # Display success information
            console.print(f"\n[green]Success![/green]")
            console.print(f"Found [bold cyan]{result.count}[/bold cyan] URLs")
            console.print(f"Saved to: [bold]{final_output_path}[/bold]")
            
            if verbose and result.urls:
                _display_summary_table(result.urls[:10])  # Show first 10
            
            # Show any warnings/errors
            if result.errors:
                console.print(f"\n[yellow]{len(result.errors)} warnings/errors occurred:[/yellow]")
                for error in result.errors[:5]:  # Show first 5 errors
                    console.print(f"  • {error}")
                if len(result.errors) > 5:
                    console.print(f"  • ... and {len(result.errors) - 5} more")
        else:
            console.print(f"[red]{result.message}[/red]")
            if result.errors and verbose:
                console.print("[red]Errors:[/red]")
                for error in result.errors:
                    console.print(f"  • {error}")
            raise typer.Exit(1)
            
    except KeyboardInterrupt:
        console.print("\n[yellow]Operation cancelled by user[/yellow]")
        raise typer.Exit(0)
    except Exception as e:
        console.print(f"[red]Unexpected error: {e}[/red]")
        if verbose:
            import traceback
            console.print("[dim]Traceback:[/dim]")
            console.print(traceback.format_exc())
        raise typer.Exit(1)


def _display_crawl_info(config: CrawlConfig, max_urls: int) -> None:
    """Display crawl configuration information."""
    table = Table(title="Crawl Configuration")
    table.add_column("Setting", style="cyan")
    table.add_column("Value", style="white")
    
    table.add_row("URL", config.url)
    table.add_row("Mode", config.mode)
    table.add_row("Max Depth", str(config.max_depth))
    table.add_row("Delay", f"{config.delay}s")
    table.add_row("Max URLs", str(max_urls))
    table.add_row("Output Format", config.output_format.upper())
    
    if config.filter_base:
        table.add_row("URL Filter", config.filter_base)
    
    console.print(table)
    console.print()


def _crawl_sitemap_mode(config: CrawlConfig, verbose: bool):
    """Perform sitemap crawling with progress display."""
    console.print("[blue]Sitemap crawling mode[/blue]")
    
    if verbose:
        # Create progress display
        with Progress(
            SpinnerColumn(),
            TextColumn("[progress.description]{task.description}"),
            TextColumn("•"),
            TextColumn("[cyan]{task.fields[urls_found]}[/cyan] URLs found"),
            console=console,
            transient=True
        ) as progress:
            
            progress_task = progress.add_task("Processing sitemap...", urls_found=0)
            
            def progress_callback(message: str, count: int):
                progress.update(progress_task, description=message, urls_found=count)
            
            service = SitemapService(progress_callback=progress_callback)
            
            if config.url.endswith('.xml'):
                result = service.process_sitemap_url(config.url, config.filter_base)
            else:
                result = service.process_base_url(config.url, config.filter_base)
    else:
        # Simple mode without progress
        service = SitemapService()
        
        if config.url.endswith('.xml'):
            result = service.process_sitemap_url(config.url, config.filter_base)
        else:
            result = service.process_base_url(config.url, config.filter_base)
    
    return result


def _crawl_website_mode(config: CrawlConfig, max_urls: int, verbose: bool):
    """Perform website crawling with progress display."""
    console.print("[blue]Website crawling mode[/blue]")
    
    if verbose:
        # Create progress display
        with Progress(
            SpinnerColumn(),
            TextColumn("[progress.description]{task.description}"),
            BarColumn(),
            TaskProgressColumn(),
            TextColumn("•"),
            TextColumn("[cyan]{task.fields[urls_found]}[/cyan] URLs"),
            console=console,
            transient=True
        ) as progress:
            
            progress_task = progress.add_task(
                "Crawling website...", 
                total=max_urls,
                urls_found=0
            )
            
            def progress_callback(message: str, count: int):
                progress.update(
                    progress_task, 
                    description=message,
                    completed=min(count, max_urls),
                    urls_found=count
                )
            
            service = CrawlerService(progress_callback=progress_callback)
            result = service.crawl_url(
                url=config.url,
                max_depth=config.max_depth,
                delay=config.delay,
                filter_base=config.filter_base
            )
    else:
        # Simple mode without progress
        service = CrawlerService()
        result = service.crawl_url(
            url=config.url,
            max_depth=config.max_depth,
            delay=config.delay,
            filter_base=config.filter_base
        )
    
    return result


def _display_summary_table(urls: List[str]) -> None:
    """Display summary table of URLs."""
    table = Table(title="Sample URLs Found")
    table.add_column("#", justify="right", style="cyan", width=4)
    table.add_column("URL", style="blue")
    
    for i, url in enumerate(urls, 1):
        # Truncate very long URLs for display
        display_url = url if len(url) <= 80 else url[:77] + "..."
        table.add_row(str(i), display_url)
    
    console.print(table)


def main() -> None:
    """Main entry point for the CLI application."""
    try:
        app()
    except KeyboardInterrupt:
        console.print("\n[yellow]Operation cancelled by user[/yellow]")
        sys.exit(0)
    except Exception as e:
        console.print(f"[red]Unexpected error: {e}[/red]")
        sys.exit(1)


if __name__ == "__main__":
    main()
</file>

<file path=".gitignore">
# Byte-compiled / optimized / DLL files
__pycache__/
*.py[cod]
*$py.class

# C extensions
*.so

# Distribution / packaging
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
share/python-wheels/
*.egg-info/
.installed.cfg
*.egg
MANIFEST

# PyInstaller
#  Usually these files are written by a python script from a template
#  before PyInstaller builds the exe, so as to inject date/other infos into it.
*.manifest
*.spec

# Installer logs
pip-log.txt
pip-delete-this-directory.txt

# Unit test / coverage reports
htmlcov/
.tox/
.nox/
.coverage
.coverage.*
.cache
nosetests.xml
coverage.xml
*.cover
*.py,cover
.hypothesis/
.pytest_cache/
cover/

# Translations
*.mo
*.pot

# Django stuff:
*.log
local_settings.py
db.sqlite3
db.sqlite3-journal

# Flask stuff:
instance/
.webassets-cache

# Scrapy stuff:
.scrapy

# Sphinx documentation
docs/_build/

# PyBuilder
.pybuilder/
target/

# Jupyter Notebook
.ipynb_checkpoints

# IPython
profile_default/
ipython_config.py

# pyenv
#   For a library or package, you might want to ignore these files since the code is
#   intended to run in multiple environments; otherwise, check them in:
# .python-version

# pipenv
#   According to pypa/pipenv#598, it is recommended to include Pipfile.lock in version control.
#   However, in case of collaboration, if having platform-specific dependencies or dependencies
#   having no cross-platform support, pipenv may install dependencies that don't work, or not
#   install all needed dependencies.
#Pipfile.lock

# poetry
#   Similar to Pipfile.lock, it is generally recommended to include poetry.lock in version control.
#   This is especially recommended for binary packages to ensure reproducibility, and is more
#   commonly ignored for libraries.
#   https://python-poetry.org/docs/basic-usage/#commit-your-poetrylock-file-to-version-control
#poetry.lock

# pdm
#   Similar to Pipfile.lock, it is generally recommended to include pdm.lock in version control.
#pdm.lock
#   pdm stores project-wide configurations in .pdm.toml, but it is recommended to not include it
#   in version control.
#   https://pdm.fming.dev/#use-with-ide
.pdm.toml

# PEP 582; used by e.g. github.com/David-OConnor/pyflow and github.com/pdm-project/pdm
__pypackages__/

# Celery stuff
celerybeat-schedule
celerybeat.pid

# SageMath parsed files
*.sage.py

# Environments
.env
.venv
env/
venv/
ENV/
env.bak/
venv.bak/

# Spyder project settings
.spyderproject
.spyproject

# Rope project settings
.ropeproject

# mkdocs documentation
/site

# mypy
.mypy_cache/
.dmypy.json
dmypy.json

# Pyre type checker
.pyre/

# pytype static type analyzer
.pytype/

# Cython debug symbols
cython_debug/

# PyCharm
#  JetBrains specific template is maintained in a separate JetBrains.gitignore that can
#  be added to the global gitignore or merged into this project gitignore.  For a PyCharm
#  project, it is preferred to specify the gitignore entries in the project directory.
#  See https://github.com/github/gitignore/blob/main/Global/JetBrains.gitignore
.idea/

# VS Code
.vscode/

# Local output files
*.txt
*.json
*.csv
crawl.log

# OS-specific files
.DS_Store
Thumbs.db

# Crush agent local data
.crush/
</file>

<file path="src/crawl_url/core/models.py">
"""Core data models for crawl-url application."""

from dataclasses import dataclass, field
from typing import List, Optional
from pathlib import Path


@dataclass
class CrawlConfig:
    """Configuration for crawling operations."""
    
    url: str
    mode: str = "auto"  # 'auto', 'sitemap', or 'crawl'
    max_depth: int = 3
    delay: float = 1.0
    filter_base: Optional[str] = None
    output_path: Optional[Path] = None
    output_format: str = "txt"
    verbose: bool = False
    
    def __post_init__(self) -> None:
        """Validate configuration after initialization."""
        if self.mode not in ("auto", "sitemap", "crawl"):
            raise ValueError(
                f"Invalid crawling mode '{self.mode}'. Please choose from:\n"
                f"  • 'auto' - Automatically detect sitemap.xml URLs vs website URLs\n"
                f"  • 'sitemap' - Extract URLs from sitemap.xml files only\n"
                f"  • 'crawl' - Recursively crawl website pages"
            )
        
        if self.max_depth < 1 or self.max_depth > 10:
            raise ValueError(
                f"Invalid crawl depth '{self.max_depth}'. Depth must be between 1 and 10.\n"
                f"  • Use depth 1-2 for quick surface-level crawling\n"
                f"  • Use depth 3-5 for thorough crawling (recommended)\n"
                f"  • Use depth 6-10 for deep crawling (may be slow)"
            )
        
        if self.delay < 0:
            raise ValueError(
                f"Invalid delay '{self.delay}'. Delay must be non-negative (0 or greater).\n"
                f"  • Recommended: 1.0 seconds for respectful crawling\n"
                f"  • Minimum: 0.5 seconds to avoid overwhelming servers\n"
                f"  • Use higher values (2-5 seconds) for conservative crawling"
            )
        
        if self.output_format not in ("txt", "json", "csv"):
            raise ValueError(
                f"Invalid output format '{self.output_format}'. Please choose from:\n"
                f"  • 'txt' - Simple text file with one URL per line\n"
                f"  • 'json' - Structured JSON with metadata and URLs array\n"
                f"  • 'csv' - Comma-separated values with URL, domain, and path columns"
            )


@dataclass
class CrawlResult:
    """Result of a crawling operation."""
    
    success: bool
    urls: List[str]
    count: int
    message: str
    errors: List[str] = field(default_factory=list)
    
    def __post_init__(self) -> None:
        """Ensure count matches URL list length."""
        if self.count != len(self.urls):
            self.count = len(self.urls)


@dataclass
class SitemapEntry:
    """Sitemap URL entry with optional metadata."""
    
    loc: str
    lastmod: Optional[str] = None
    changefreq: Optional[str] = None
    priority: Optional[str] = None
    
    def __post_init__(self) -> None:
        """Validate sitemap entry data."""
        if not self.loc:
            raise ValueError(
                "Sitemap entry is missing required location URL. "
                "Each sitemap entry must have a valid <loc> element with a URL."
            )
        
        # Validate changefreq if provided
        valid_changefreq = {
            "always", "hourly", "daily", "weekly", 
            "monthly", "yearly", "never"
        }
        if self.changefreq and self.changefreq not in valid_changefreq:
            valid_options = "', '".join(sorted(valid_changefreq))
            raise ValueError(
                f"Invalid changefreq value '{self.changefreq}' in sitemap entry.\n"
                f"Valid changefreq values are: '{valid_options}'"
            )
        
        # Validate priority if provided
        if self.priority:
            try:
                priority_val = float(self.priority)
                if not (0.0 <= priority_val <= 1.0):
                    raise ValueError(
                        f"Invalid priority value '{self.priority}' in sitemap entry.\n"
                        f"Priority must be a decimal number between 0.0 and 1.0 (e.g., '0.5', '0.8', '1.0')"
                    )
            except ValueError as e:
                if "Priority must be" in str(e):
                    raise
                raise ValueError(
                    f"Invalid priority format '{self.priority}' in sitemap entry.\n"
                    f"Priority must be a valid decimal number between 0.0 and 1.0 (e.g., '0.5', '0.8', '1.0')"
                )
</file>

<file path="crawl_urls.py">
import asyncio
from mcp import get_mcp_client

async def crawl_urls_from_file(file_path):
    # Inicializa o cliente MCP
    client = get_mcp_client()
    
    # Lê as URLs do arquivo
    with open(file_path, 'r') as file:
        urls = [line.strip() for line in file if line.strip()]
    
    # Processa cada URL com crawl_single_page
    for i, url in enumerate(urls, 1):
        print(f"Processando ({i}/{len(urls)}): {url}")
        try:
            result = await client.call_tool("crawl_single_page", {"url": url})
            print(f"Resultado para {url}: {result}")
        except Exception as e:
            print(f"Erro ao processar {url}: {e}")

if __name__ == "__main__":
    asyncio.run(crawl_urls_from_file("docs_claude_code.txt"))
</file>

<file path="src/crawl_url/utils/__init__.py">
"""Utility functions and helpers."""
</file>

<file path="LICENSE">
MIT License

Copyright (c) 2024 Crawl-URL Team

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
</file>

<file path="pyproject.toml">
[build-system]
requires = ["setuptools>=61.0", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "crawl-url"
version = "1.0.0"
description = "A powerful terminal application for crawling and extracting URLs from websites"
readme = "README.md"
requires-python = ">=3.8"
license = {text = "MIT"}
authors = [
    {name = "Crawl-URL Team", email = "crawl-url@example.com"}
]
maintainers = [
    {name = "Crawl-URL Team", email = "crawl-url@example.com"}
]
keywords = ["cli", "web-crawling", "url", "scraping", "terminal", "sitemap"]
classifiers = [
    "Development Status :: 4 - Beta",
    "Environment :: Console",
    "Intended Audience :: Developers",
    "Intended Audience :: System Administrators",
    "License :: OSI Approved :: MIT License",
    "Operating System :: OS Independent",
    "Programming Language :: Python :: 3",
    "Programming Language :: Python :: 3.8",
    "Programming Language :: Python :: 3.9",
    "Programming Language :: Python :: 3.10",
    "Programming Language :: Python :: 3.11",
    "Programming Language :: Python :: 3.12",
    "Topic :: Internet :: WWW/HTTP",
    "Topic :: System :: Systems Administration",
    "Topic :: Utilities",
]

# Core dependencies
dependencies = [
    "requests>=2.28.0",
    "beautifulsoup4>=4.11.0",
    "lxml>=4.9.0",
    "pytermgui>=7.4.0",
    "typer>=0.9.0",
    "rich>=13.0.0",
]

# Optional dependencies
[project.optional-dependencies]
dev = [
    "pytest>=7.0.0",
    "pytest-cov>=4.0.0",
    "pytest-mock>=3.10.0",
    "black>=23.0.0",
    "ruff>=0.1.0",
    "mypy>=1.0.0",
    "pre-commit>=3.0.0",
    "twine>=4.0.0",
    "build>=0.10.0",
]
test = [
    "pytest>=7.0.0",
    "pytest-cov>=4.0.0",
    "pytest-mock>=3.10.0",
]

# Console script entry points
[project.scripts]
crawl-url = "crawl_url.cli:main"

[project.urls]
Homepage = "https://github.com/crawl-url/crawl-url"
Repository = "https://github.com/crawl-url/crawl-url"
Issues = "https://github.com/crawl-url/crawl-url/issues"
Documentation = "https://github.com/crawl-url/crawl-url#readme"

# Setuptools-specific configuration
[tool.setuptools.packages.find]
where = ["src"]
include = ["crawl_url*"]

# Black formatting configuration
[tool.black]
line-length = 88
target-version = ["py38", "py39", "py310", "py311", "py312"]
include = '\.pyi?$'
extend-exclude = '''
/(
  # directories
  \.eggs
  | \.git
  | \.hg
  | \.mypy_cache
  | \.tox
  | \.venv
  | build
  | dist
)/
'''

# Ruff linting configuration  
[tool.ruff]
target-version = "py38"
line-length = 88
select = [
    "E",  # pycodestyle errors
    "W",  # pycodestyle warnings
    "F",  # pyflakes
    "I",  # isort
    "B",  # flake8-bugbear
    "C4", # flake8-comprehensions
    "UP", # pyupgrade
]
ignore = [
    "E501",  # line too long, handled by black
    "B008",  # do not perform function calls in argument defaults
]

[tool.ruff.per-file-ignores]
"__init__.py" = ["F401"]

# MyPy type checking configuration
[tool.mypy]
python_version = "3.8"
check_untyped_defs = true
disallow_any_generics = true
disallow_incomplete_defs = true
disallow_untyped_defs = true
no_implicit_optional = true
warn_redundant_casts = true
warn_unused_ignores = true

# Pytest configuration
[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py", "*_test.py"]
python_classes = ["Test*"]
python_functions = ["test_*"]
addopts = [
    "--strict-markers",
    "--strict-config",
    "--cov=crawl_url",
    "--cov-report=html",
    "--cov-report=term-missing",
]
markers = [
    "slow: marks tests as slow (deselect with '-m \"not slow\"')",
    "integration: marks tests as integration tests",
    "unit: marks tests as unit tests",
]

# Coverage configuration
[tool.coverage.run]
source = ["src"]
branch = true

[tool.coverage.report]
exclude_lines = [
    "pragma: no cover",
    "def __repr__",
    "raise AssertionError",
    "raise NotImplementedError",
    "if __name__ == .__main__.:",
]
</file>

<file path="alacritty.toml">
# =============================================================================
# Configuração do Alacritty para PowerShell 7
# =============================================================================

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO GERAL
# -----------------------------------------------------------------------------
[general]
# Recarregamento automático da configuração
live_config_reload = true

# Diretório de trabalho inicial (opcional - descomente se necessário)
# working_directory = "C:\\Users\\%USERNAME%"

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DO SHELL - PowerShell 7
# -----------------------------------------------------------------------------
[shell]
program = "pwsh"  # PowerShell 7 (pwsh.exe)
args = [
    "-NoLogo",                    # Remove o logo de inicialização
    "-WorkingDirectory", ".",     # Define diretório de trabalho atual
]

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DA JANELA
# -----------------------------------------------------------------------------
[window]
# Dimensões da janela (colunas x linhas)
dimensions = { columns = 120, lines = 35 }

# Posição inicial da janela (opcional)
# position = { x = 100, y = 100 }

# Padding interno da janela
padding = { x = 15, y = 15 }

# Padding dinâmico (distribui espaço extra uniformemente)
dynamic_padding = true

# Decorações da janela
decorations = "Full"

# Opacidade da janela (0.0 = transparente, 1.0 = opaco)
opacity = 0.95

# Desfoque de fundo (funciona no Windows 11)
blur = true

# Título da janela
title = "Alacritty - PowerShell 7"

# Permitir que aplicações alterem o título
dynamic_title = true

# Modo de inicialização
startup_mode = "Windowed"

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE FONTE
# -----------------------------------------------------------------------------
[font]
# Fonte normal
normal = { family = "JetBrains Mono", style = "Regular" }

# Fonte negrito
bold = { family = "JetBrains Mono", style = "Bold" }

# Fonte itálico
italic = { family = "JetBrains Mono", style = "Italic" }

# Fonte negrito itálico
bold_italic = { family = "JetBrains Mono", style = "Bold Italic" }

# Tamanho da fonte
size = 11.0

# Offset da fonte (espaçamento)
offset = { x = 0, y = 1 }

# Desenho de caixas integrado
builtin_box_drawing = true

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE CORES - Tema Dracula Adaptado
# -----------------------------------------------------------------------------
[colors]
# Cores primárias
[colors.primary]
background = "#282a36"
foreground = "#f8f8f2"
dim_foreground = "#6272a4"

# Cores do cursor
[colors.cursor]
text = "#282a36"
cursor = "#f8f8f2"

# Cores do cursor no modo vi
[colors.vi_mode_cursor]
text = "#282a36"
cursor = "#ffb86c"

# Cores de seleção
[colors.selection]
text = "#282a36"
background = "#44475a"

# Cores de busca
[colors.search.matches]
foreground = "#282a36"
background = "#ffb86c"

[colors.search.focused_match]
foreground = "#282a36"
background = "#ff79c6"

# Cores normais
[colors.normal]
black = "#000000"
red = "#ff5555"
green = "#50fa7b"
yellow = "#f1fa8c"
blue = "#bd93f9"
magenta = "#ff79c6"
cyan = "#8be9fd"
white = "#bfbfbf"

# Cores brilhantes
[colors.bright]
black = "#4d4d4d"
red = "#ff6e67"
green = "#5af78e"
yellow = "#f4f99d"
blue = "#caa9fa"
magenta = "#ff92d0"
cyan = "#9aedfe"
white = "#e6e6e6"

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DO CURSOR
# -----------------------------------------------------------------------------
[cursor]
# Estilo do cursor
style = { shape = "Block", blinking = "On" }

# Intervalo de piscada (em milissegundos)
blink_interval = 500

# Timeout de piscada (em segundos)
blink_timeout = 0

# Cursor oco quando a janela não está focada
unfocused_hollow = true

# Espessura do cursor (0.0 a 1.0)
thickness = 0.15

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE ROLAGEM
# -----------------------------------------------------------------------------
[scrolling]
# Histórico de linhas
history = 50000

# Multiplicador de rolagem
multiplier = 3

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE SELEÇÃO
# -----------------------------------------------------------------------------
[selection]
# Caracteres de escape semântico
semantic_escape_chars = ",│`|:\"' ()[]{}<>\t"

# Salvar seleção no clipboard automaticamente
save_to_clipboard = true

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE BELL (Notificações)
# -----------------------------------------------------------------------------
[bell]
# Duração do bell visual (0 = desabilitado)
duration = 100

# Cor do bell visual
color = "#ffb86c"

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE TERMINAL
# -----------------------------------------------------------------------------
[terminal]
# Suporte a OSC 52 (clipboard)
osc52 = "CopyPaste"

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE MOUSE
# -----------------------------------------------------------------------------
[mouse]
# Ocultar cursor quando digitando
hide_when_typing = true

# Bindings do mouse
[[mouse.bindings]]
mouse = "Right"
action = "PasteSelection"

# -----------------------------------------------------------------------------
# CONFIGURAÇÃO DE TECLADO
# -----------------------------------------------------------------------------
[keyboard]
# Bindings personalizados
bindings = [
    # Navegação de abas (se usando tmux ou similar)
    { key = "T", mods = "Control|Shift", action = "CreateNewWindow" },
    { key = "W", mods = "Control|Shift", action = "Quit" },
    
    # Controle de fonte
    { key = "Plus", mods = "Control", action = "IncreaseFontSize" },
    { key = "Minus", mods = "Control", action = "DecreaseFontSize" },
    { key = "Key0", mods = "Control", action = "ResetFontSize" },
    
    # Clipboard
    { key = "C", mods = "Control|Shift", action = "Copy" },
    { key = "V", mods = "Control|Shift", action = "Paste" },
    
    # Busca
    { key = "F", mods = "Control|Shift", action = "SearchForward" },
    { key = "B", mods = "Control|Shift", action = "SearchBackward" },
    
    # Limpeza do terminal
    { key = "L", mods = "Control|Shift", chars = "clear\n" },
    
    # PowerShell específicos
    { key = "R", mods = "Control", chars = "\u0012" },  # Ctrl+R para histórico
]

# -----------------------------------------------------------------------------
# CONFIGURAÇÕES DE DEBUG (opcional)
# -----------------------------------------------------------------------------
[debug]
# Nível de log
log_level = "Warn"

# Log persistente
persistent_logging = false

# Timer de renderização
render_timer = false

# -----------------------------------------------------------------------------
# VARIÁVEIS DE AMBIENTE
# -----------------------------------------------------------------------------
[env]
# Força o PowerShell a usar UTF-8
POWERSHELL_TELEMETRY_OPTOUT = "1"

# Terminal colorido
TERM = "xterm-256color"
</file>

<file path="src/crawl_url/core/ui.py">
"""PyTermGUI interactive interface for crawl-url application."""

import shutil
import sys
from pathlib import Path
from typing import List, Optional

from ..utils.storage import StorageManager
from .crawler import CrawlerService
from .models import CrawlConfig
from .sitemap_parser import SitemapService

# Try to import PyTermGUI with fallback handling
try:
    import pytermgui as ptg
    PTG_AVAILABLE = True
except ImportError:
    PTG_AVAILABLE = False


class CrawlerApp:
    """Main PyTermGUI application for interactive URL crawling."""
    
    def __init__(self) -> None:
        """Initialize the crawler application."""
        self.manager: Optional[ptg.WindowManager] = None
        self.results: List[str] = []
        self.current_config: Optional[CrawlConfig] = None
        self.storage_manager = StorageManager()
        
        # Input field references
        self.url_input: Optional[ptg.InputField] = None
        self.mode_input: Optional[ptg.InputField] = None
        self.filter_input: Optional[ptg.InputField] = None
        self.depth_input: Optional[ptg.InputField] = None
        self.delay_input: Optional[ptg.InputField] = None
        self.format_input: Optional[ptg.InputField] = None
        
        # Status and progress
        self.status_label: Optional[ptg.Label] = None
        self.progress_label: Optional[ptg.Label] = None
        
    def run(self) -> None:
        """Run the interactive TUI application with fallback."""
        import platform
        
        # On Windows, prefer console mode by default due to PyTermGUI compatibility issues
        if platform.system() == "Windows":
            print("Windows detected - using console mode for better compatibility")
            self._run_console_fallback()
            return
            
        if not PTG_AVAILABLE:
            self._run_console_fallback()
            return
        
        try:
            self._test_tui_compatibility()
            self._run_tui_mode()
        except KeyboardInterrupt:
            print("\nOperation cancelled by user")
            return
        except Exception as e:
            print(f"TUI mode not available: {str(e)}")
            print("Falling back to console mode...")
            self._run_console_fallback()
    
    def _test_tui_compatibility(self) -> None:
        """Test PyTermGUI compatibility before launching full interface."""
        try:
            # Simple, fast compatibility test - just try to create basic objects
            test_label = ptg.Label("Test")
            test_window = ptg.Window(test_label, width=20, height=3)
            
            # If we can create these objects without exception, TUI should work
            # No need to actually display anything
            
        except Exception as e:
            raise Exception(f"PyTermGUI compatibility test failed: {e}")
    
    def _run_tui_mode(self) -> None:
        """Run the full TUI interface."""
        self.manager = ptg.WindowManager()
        
        try:
            with self.manager:
                main_window = self._create_main_window()
                self.manager.add(main_window)
                self.manager.run()
        except KeyboardInterrupt:
            self._safe_exit("Operation cancelled by user")
        except Exception as e:
            self._safe_exit(f"TUI error: {e}")
    
    def _create_main_window(self) -> ptg.Window:
        """Create the main configuration window."""
        terminal_width, terminal_height = shutil.get_terminal_size()
        
        # Responsive design
        if terminal_width < 80:
            width = terminal_width - 4
            box = ptg.boxes.SINGLE
        else:
            width = 80
            box = ptg.boxes.DOUBLE
        
        # Create input fields
        self.url_input = ptg.InputField("https://", prompt="URL to crawl: ")
        self.mode_input = ptg.InputField("auto", prompt="Mode (auto/sitemap/crawl): ")
        self.filter_input = ptg.InputField("", prompt="Filter base URL (optional): ")
        self.depth_input = ptg.InputField("3", prompt="Max depth (crawl mode): ")
        self.delay_input = ptg.InputField("1.0", prompt="Delay between requests: ")
        self.format_input = ptg.InputField("txt", prompt="Output format (txt/json/csv): ")
        
        # Status label
        self.status_label = ptg.Label("[blue]ℹ️ Ready to crawl[/blue]")
        
        return ptg.Window(
            ptg.Label("[bold blue]🕷️ Crawl-URL Interactive Interface[/bold blue]", centered=True),
            "",
            ptg.Label("[bold]Configuration[/bold]"),
            self.url_input,
            self.mode_input,
            self.filter_input,
            self.depth_input,
            self.delay_input,
            self.format_input,
            "",
            self.status_label,
            "",
            ptg.Container(
                ptg.Button("Start Crawling", self._start_crawl),
                ptg.Button("View Results", self._show_results),
                ptg.Button("Clear Results", self._clear_results),
                ptg.Button("Help", self._show_help),
                ptg.Button("Quit", self._quit_app),
                box=ptg.boxes.EMPTY
            ),
            width=width,
            box=box
        ).center()
    
    def _create_progress_window(self) -> ptg.Window:
        """Create progress display window."""
        self.progress_label = ptg.Label("Initializing...")
        
        return ptg.Window(
            ptg.Label("[bold]Crawling Progress[/bold]", centered=True),
            "",
            self.progress_label,
            ptg.Label(""),
            self.status_label,
            "",
            ptg.Button("Cancel", self._cancel_crawl),
            width=60,
            box=ptg.boxes.SINGLE
        ).center()
    
    def _create_results_window(self) -> ptg.Window:
        """Create results display window with scrollable content."""
        # Create scrollable container for URLs
        results_container = ptg.Container()
        
        if not self.results:
            results_container.add(ptg.Label("[yellow]No results to display[/yellow]"))
        else:
            # Limit display for performance (first 100 URLs)
            display_urls = self.results[:100]
            
            for i, url in enumerate(display_urls, 1):
                # Truncate very long URLs
                display_url = url if len(url) <= 70 else url[:67] + "..."
                label = ptg.Label(f"{i:3d}. {display_url}")
                results_container.add(label)
            
            if len(self.results) > 100:
                results_container.add(ptg.Label(f"[dim]... and {len(self.results) - 100} more URLs[/dim]"))
        
        return ptg.Window(
            ptg.Label(f"[bold]Found {len(self.results)} URLs[/bold]", centered=True),
            "",
            results_container,
            "",
            ptg.Container(
                ptg.Button("Export Results", self._export_results),
                ptg.Button("Back to Main", self._show_main),
                box=ptg.boxes.EMPTY
            ),
            width=min(100, shutil.get_terminal_size()[0] - 4),
            height=min(25, shutil.get_terminal_size()[1] - 4),
            box=ptg.boxes.DOUBLE
        ).center()
    
    def _create_help_window(self) -> ptg.Window:
        """Create help information window."""
        help_text = [
            "[bold]Crawl-URL Help[/bold]",
            "",
            "[bold]Modes:[/bold]",
            "• auto: Automatically detect sitemap or crawl mode",
            "• sitemap: Parse sitemap.xml files only", 
            "• crawl: Recursively crawl website pages",
            "",
            "[bold]URL Filter:[/bold]",
            "Only include URLs starting with the filter base",
            "Example: https://docs.anthropic.com/en/docs/",
            "",
            "[bold]Output Formats:[/bold]",
            "• txt: Plain text, one URL per line",
            "• json: JSON format with metadata",
            "• csv: CSV format with URL components",
            "",
            "[bold]Tips:[/bold]",
            "• Use delay ≥1.0 for respectful crawling",
            "• Filter helps focus on specific site sections",
            "• Check robots.txt compliance automatically"
        ]
        
        help_container = ptg.Container()
        for line in help_text:
            help_container.add(ptg.Label(line))
        
        return ptg.Window(
            help_container,
            "",
            ptg.Button("Close", self._show_main),
            width=60,
            height=20,
            box=ptg.boxes.SINGLE
        ).center()
    
    def _start_crawl(self) -> None:
        """Start the crawling process."""
        try:
            # Validate and create configuration
            config = self._create_config_from_inputs()
            if not config:
                return
            
            self.current_config = config
            
            # Show progress window
            progress_window = self._create_progress_window()
            self.manager.remove_window(self.manager.focused_window)
            self.manager.add(progress_window)
            
            # Update status
            self._update_status("Starting crawl...", "info")
            
            # Perform crawling based on mode
            if config.mode in ("sitemap", "auto") and (
                config.url.endswith('.xml') or config.mode == "sitemap"
            ):
                self._crawl_sitemap(config)
            else:
                self._crawl_website(config)
            
        except Exception as e:
            self._update_status(f"Error: {e}", "error")
    
    def _create_config_from_inputs(self) -> Optional[CrawlConfig]:
        """Create crawl configuration from input fields."""
        try:
            url = self.url_input.value.strip()
            if not url:
                self._update_status("🌐 Please enter a website URL or sitemap.xml URL to crawl", "error")
                return None
            
            mode = self.mode_input.value.strip().lower()
            if mode not in ("auto", "sitemap", "crawl"):
                self._update_status("🔍 Mode must be: 'auto' (detect), 'sitemap' (XML only), or 'crawl' (recursive)", "error")
                return None
            
            try:
                max_depth = int(self.depth_input.value.strip())
                delay = float(self.delay_input.value.strip())
            except ValueError:
                self._update_status("⚙️ Invalid input: Depth must be a whole number (1-10), delay must be a decimal number (e.g., 1.0)", "error")
                return None
            
            filter_base = self.filter_input.value.strip() or None
            output_format = self.format_input.value.strip().lower()
            
            if output_format not in ("txt", "json", "csv"):
                self._update_status("📄 Output format must be: 'txt' (plain text), 'json' (structured), or 'csv' (spreadsheet)", "error")
                return None
            
            return CrawlConfig(
                url=url,
                mode=mode,
                max_depth=max_depth,
                delay=delay,
                filter_base=filter_base,
                output_format=output_format
            )
            
        except Exception as e:
            self._update_status(f"⚠️ Configuration problem: {e}", "error")
            return None
    
    def _crawl_sitemap(self, config: CrawlConfig) -> None:
        """Perform sitemap crawling."""
        service = SitemapService(progress_callback=self._progress_callback)
        
        if config.url.endswith('.xml'):
            result = service.process_sitemap_url(config.url, config.filter_base)
        else:
            result = service.process_base_url(config.url, config.filter_base)
        
        self._handle_crawl_result(result, config)
    
    def _crawl_website(self, config: CrawlConfig) -> None:
        """Perform website crawling."""
        service = CrawlerService(progress_callback=self._progress_callback)
        result = service.crawl_url(
            url=config.url,
            max_depth=config.max_depth,
            delay=config.delay,
            filter_base=config.filter_base
        )
        
        self._handle_crawl_result(result, config)
    
    def _handle_crawl_result(self, result, config: CrawlConfig) -> None:
        """Handle the result of crawling operation."""
        if result.success:
            self.results = result.urls
            self._update_status(f"✅ Found {result.count} URLs", "success")
            
            # Auto-save results
            try:
                output_path = self.storage_manager.save_urls(
                    urls=result.urls,
                    base_url=config.url,
                    format_type=config.output_format
                )
                self._update_status(f"📁 Saved to: {output_path.name}", "info")
            except Exception as e:
                self._update_status(f"⚠️ Crawl completed but save failed: {e}", "warning")
        else:
            self._update_status(f"❌ {result.message}", "error")
    
    def _progress_callback(self, message: str, count: int) -> None:
        """Update progress display."""
        if self.progress_label:
            self.progress_label.value = f"[cyan]🔄 {message}[/cyan]"
        if self.status_label and count > 0:
            self.status_label.value = f"[green]{count} URLs found so far...[/green]"
    
    def _update_status(self, message: str, status_type: str = "info") -> None:
        """Update status display with colored message."""
        colors = {
            "info": "blue",
            "success": "green", 
            "warning": "yellow",
            "error": "red"
        }
        
        color = colors.get(status_type, "white")
        if self.status_label:
            self.status_label.value = f"[{color}]{message}[/{color}]"
    
    def _show_results(self) -> None:
        """Show results window."""
        results_window = self._create_results_window()
        self.manager.remove_window(self.manager.focused_window)
        self.manager.add(results_window)
    
    def _show_help(self) -> None:
        """Show help window."""
        help_window = self._create_help_window()
        self.manager.remove_window(self.manager.focused_window)
        self.manager.add(help_window)
    
    def _show_main(self) -> None:
        """Return to main window."""
        main_window = self._create_main_window()
        self.manager.remove_window(self.manager.focused_window)
        self.manager.add(main_window)
    
    def _export_results(self) -> None:
        """Export results in different format."""
        if not self.results or not self.current_config:
            self._update_status("No results to export", "warning")
            return
        
        try:
            formats = ["txt", "json", "csv"]
            current_format = self.current_config.output_format
            
            # Try next format in sequence
            next_format_idx = (formats.index(current_format) + 1) % len(formats)
            next_format = formats[next_format_idx]
            
            output_path = self.storage_manager.save_urls(
                urls=self.results,
                base_url=self.current_config.url,
                format_type=next_format,
                custom_suffix="export"
            )
            
            self._update_status(f"📁 Exported as {next_format.upper()}: {output_path.name}", "success")
            
        except Exception as e:
            self._update_status(f"Export failed: {e}", "error")
    
    def _clear_results(self) -> None:
        """Clear current results."""
        self.results.clear()
        self._update_status("Results cleared", "info")
    
    def _cancel_crawl(self) -> None:
        """Cancel current crawl operation."""
        self._show_main()
        self._update_status("Crawl cancelled", "warning")
    
    def _quit_app(self) -> None:
        """Quit the application."""
        self._safe_exit("Goodbye!")
    
    def _safe_exit(self, message: str = "") -> None:
        """Safely exit the application."""
        if message:
            print(message)
        if self.manager:
            self.manager.stop()
        sys.exit(0)
    
    def _run_console_fallback(self) -> None:
        """Run simplified console-based interface when TUI is not available."""
        print("Crawl-URL Console Mode")
        print("=" * 50)
        print("Running in console mode for Windows compatibility.")
        print("For full interactive experience on Linux, PyTermGUI is available.")
        print()
        
        try:
            url = input("Enter URL to crawl: ").strip()
            if not url:
                print("No URL provided. Exiting.")
                return
            
            mode = input("Mode (auto/sitemap/crawl) [auto]: ").strip().lower() or "auto"
            filter_base = input("Filter base URL (optional): ").strip() or None
            
            print("\nStarting crawl...")
            
            if mode in ("sitemap", "auto") and (url.endswith('.xml') or mode == "sitemap"):
                service = SitemapService()
                if url.endswith('.xml'):
                    result = service.process_sitemap_url(url, filter_base)
                else:
                    result = service.process_base_url(url, filter_base)
            else:
                service = CrawlerService()
                result = service.crawl_url(url=url, filter_base=filter_base)
            
            if result.success:
                print(f"\n[SUCCESS] Found {result.count} URLs")
                
                # Save results
                storage_manager = StorageManager()
                output_path = storage_manager.save_urls(
                    urls=result.urls,
                    base_url=url,
                    format_type="txt"
                )
                print(f"Results saved to: {output_path}")
                
                # Show first few URLs
                print("\nFirst 10 URLs:")
                for i, result_url in enumerate(result.urls[:10], 1):
                    print(f"{i:2d}. {result_url}")
                
                if len(result.urls) > 10:
                    print(f"... and {len(result.urls) - 10} more URLs")
            else:
                print(f"\n[FAILED] Crawl failed: {result.message}")
                
        except KeyboardInterrupt:
            print("\n\nOperation cancelled by user.")
        except Exception as e:
            print(f"\nError: {e}")


# Mock classes for testing when PyTermGUI is not available
class MockWindowManager:
    """Mock WindowManager for testing."""
    
    def __enter__(self):
        return self
    
    def __exit__(self, *args):
        pass
    
    def add(self, window):
        pass
    
    def run(self):
        pass
</file>

