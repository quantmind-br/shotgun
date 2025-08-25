# Security and Performance

### Security Requirements

**Frontend Security:**
- CSP Headers: Not applicable (terminal application)
- XSS Prevention: Template output sanitization to prevent terminal escape sequence injection
- Secure Storage: Session data stored with restricted file permissions (0600) in user config directories

**Backend Security:**
- Input Validation: All file paths validated against directory traversal attacks (../ sequences)
- Rate Limiting: File scanning operations limited by configurable timeouts and resource constraints
- CORS Policy: Not applicable (no web API)

**Authentication Security:**
- Token Storage: Not applicable (single-user desktop application)
- Session Management: Local session files with appropriate file system permissions
- Password Policy: Not applicable (no authentication system)

### Performance Optimization

**Frontend Performance:**
- Bundle Size Target: Single binary under 15MB (including embedded templates)
- Loading Strategy: Progressive screen loading with lazy initialization of expensive components
- Caching Strategy: Template parsing results cached in memory, file tree state cached during navigation

**Backend Performance:**
- Response Time Target: File scanning under 5 seconds for 1000+ files, template processing under 2 seconds
- Database Optimization: File-based storage optimized with indexing and atomic operations
- Caching Strategy: Template compilation cache, file metadata cache with TTL expiration
