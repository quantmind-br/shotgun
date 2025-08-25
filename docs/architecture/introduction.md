# Introduction

This document outlines the complete architecture for **shotgun-cli-v3**, a high-performance Terminal User Interface (TUI) application built with Go that generates standardized LLM prompts from templates. While traditionally "fullstack" refers to web applications, this document treats the TUI interface as the "frontend" layer and the core processing engine as the "backend" layer, providing a unified architectural approach for this sophisticated command-line application.

This unified architecture combines what would traditionally be separate interface design and core system architecture, streamlining development for a modern Go application where the terminal interface and processing engine are tightly integrated through the Bubble Tea reactive framework.

### Starter Template or Existing Project

Based on the PRD review, this is a **greenfield project** that will be built from scratch using the Go ecosystem. The project follows established patterns for TUI applications:

- **Framework Foundation:** Built on Charm's Bubble Tea v2.0.0-beta.4 (Elm Architecture)
- **Component Library:** Bubbles v0.21.0 for mature UI components
- **Styling Framework:** Lip Gloss v1.0.0 for terminal styling
- **No Web Framework Dependencies:** Pure terminal application with no HTTP/web components

**Decision:** N/A - Greenfield project with carefully selected Go TUI stack

### Change Log
| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2024-03-15 | 1.0 | Initial architecture document creation | Winston (Architect) |
