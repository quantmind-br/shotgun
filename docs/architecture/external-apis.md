# External APIs

Based on the PRD requirements and component design, **shotgun-cli-v3** is designed as a self-contained TUI application with minimal external dependencies. The application operates entirely offline and does not require external API integrations for its core functionality.

**No external APIs are required for the MVP functionality.**

The application achieves its goals through:

- **Local file system operations** for project scanning and template management
- **Embedded templates** packaged within the binary
- **Local configuration** stored in user directories
- **Offline processing** of all template generation and file operations

### Future Extension Possibilities

While not required for the initial release, potential external API integrations could enhance functionality in future versions:

#### GitHub API (Optional Future Enhancement)

- **Purpose:** Enable template sharing and community template repositories
- **Documentation:** https://docs.github.com/en/rest
- **Base URL(s):** https://api.github.com
- **Authentication:** Personal Access Token or GitHub App
- **Rate Limits:** 5,000 requests per hour (authenticated)

**Key Endpoints Used:**
- `GET /repos/{owner}/{repo}/contents/{path}` - Download template files from repositories
- `GET /search/repositories` - Discover community template repositories

**Integration Notes:** Would require network connectivity and authentication setup, potentially conflicting with the offline-first design philosophy

#### Template Registry API (Hypothetical Future Service)

- **Purpose:** Centralized template discovery and version management
- **Documentation:** Not applicable (hypothetical service)
- **Base URL(s):** https://templates.shotgun-cli.dev (example)
- **Authentication:** API key or OAuth
- **Rate Limits:** To be determined

**Key Endpoints Used:**
- `GET /api/v1/templates` - List available templates with metadata
- `GET /api/v1/templates/{id}/download` - Download specific template versions

**Integration Notes:** Would enable template ecosystem but requires infrastructure and maintenance
