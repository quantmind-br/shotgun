# Project Brief: Shotgun CLI

## Executive Summary

Shotgun CLI is a terminal-based user interface (TUI) application written in Go using Bubble Tea that generates standardized LLM prompts from pre-established or custom templates with a minimalist and elegant interface.

**Primary Problem**: Developers and teams using LLMs in development workflows face inconsistent prompt formatting, manual copy-paste processes, and formatting variations that reduce efficiency and reproducibility.

**Target Market**: Developers and teams who use LLMs for development workflows, code review, planning, and debugging tasks.

**Key Value Proposition**: Provides reproducibility, consistency, and speed when generating prompts with repository file context through a keyboard-only TUI interface.

## Problem Statement

**Current State**: Developers increasingly rely on LLMs for code analysis, debugging, planning, and review tasks. However, the current process involves:
- Manual copying and pasting of code files into prompts
- Inconsistent formatting across team members
- Time-consuming context gathering for each LLM interaction
- Lack of standardized templates for common development tasks
- Difficulty reproducing successful prompt patterns

**Impact**: This inefficiency costs development teams significant time daily, reduces the quality and consistency of LLM interactions, and creates friction in AI-assisted development workflows.

**Why Existing Solutions Fall Short**: Current tools are either web-based (breaking terminal workflow), too generic (not optimized for code contexts), or require complex setup processes that don't integrate well with existing development environments.

**Urgency**: As LLM adoption in development workflows accelerates, the need for efficient, reproducible prompt generation becomes critical for maintaining team productivity and code quality standards.

## Proposed Solution

**Core Concept**: A single-binary Go TUI application that provides a 5-screen wizard for generating structured LLM prompts from repository context and predefined templates.

**Key Differentiators**:
- **100% keyboard navigation** optimized for developer workflows
- **Repository-aware file selection** with .gitignore/.shotgunignore support  
- **Template system** with built-in templates for common dev tasks (bug analysis, code review, planning)
- **Minimalist monochrome design** that works in any terminal environment
- **Zero configuration startup** - works immediately in any project directory

**Why This Will Succeed**: Leverages Go's excellent cross-platform terminal capabilities, Bubble Tea's mature TUI framework, and focuses specifically on the developer use case rather than trying to be a general-purpose tool.

**High-Level Vision**: The standard tool for generating consistent, high-quality LLM prompts in development workflows, similar to how tools like `git` or `docker` became essential development utilities.

## Target Users

### Primary User Segment: Software Development Teams

**Profile**: 
- 2-20 person development teams using modern toolchains
- Already using LLMs for code assistance (GitHub Copilot, ChatGPT, Claude)
- Command-line comfortable developers
- Working in Git repositories with standard project structures

**Current Behaviors**:
- Manually copying code files into LLM interfaces
- Using screenshots or partial code snippets for context
- Inconsistent prompt formats across team members
- Spending 5-15 minutes per LLM session gathering context

**Pain Points**:
- Time lost in context gathering and formatting
- Inconsistent results due to prompt variations
- Difficulty sharing effective prompt patterns
- Breaking flow to switch to web interfaces

**Goals**: 
- Reduce time spent on prompt preparation
- Increase consistency and quality of LLM interactions
- Maintain terminal-based workflow
- Share and standardize effective prompting patterns

### Secondary User Segment: Solo Developers & Consultants

**Profile**: Independent developers, freelancers, and technical consultants working across multiple client projects

**Needs**: Quick context switching between projects, professional prompt formatting, reproducible analysis patterns

## Goals & Success Metrics

### Business Objectives
- **Adoption**: 1,000+ GitHub stars within 6 months of release
- **Usage**: 10,000+ downloads across all platforms
- **Community**: Active template sharing and contribution ecosystem
- **Integration**: Mentioned in developer productivity blog posts and workflows

### User Success Metrics  
- **Time Savings**: Average 5-10 minutes saved per LLM interaction
- **Consistency**: 95% of team members using standardized templates
- **Frequency**: Daily usage by primary user segments
- **Satisfaction**: 4.5+ stars on package managers and GitHub

### Key Performance Indicators (KPIs)
- **Performance**: File scanning <5s for repos with 1000+ files
- **Reliability**: <1% crash rate across all platforms
- **Usability**: <2 minutes to complete first successful prompt generation
- **Compatibility**: Support for 95% of common terminal environments

## MVP Scope

### Core Features (Must Have)
- **File Tree Navigation**: Interactive selection with checkboxes, hierarchical selection, .gitignore/.shotgunignore support
- **Template Selection**: Built-in templates (bug analysis, diff generation, planning, project management) with metadata display  
- **Text Input Fields**: Multiline task description and optional rules input with UTF-8 support
- **Prompt Generation**: Combine template + selected files + user input into formatted .md output
- **Keyboard Navigation**: Complete 5-screen wizard navigable entirely by keyboard (F1-F10, arrows, space, enter)

### Out of Scope for MVP
- Web interface or GUI version
- Real-time collaboration features
- Integration with specific LLM APIs
- Advanced template editor within the TUI
- Plugin system or extensions
- Cloud synchronization of templates or settings

### MVP Success Criteria
- **Functional**: Successfully generates usable prompts for all built-in templates
- **Performance**: Handles typical development repositories (100-1000 files) smoothly
- **Compatibility**: Works on Windows PowerShell, macOS Terminal, and Linux terminals
- **Usability**: New users can generate first prompt within 5 minutes

## Post-MVP Vision

### Phase 2 Features
- **Custom Template Creation**: In-app TOML template editor
- **Session History**: Save and restore previous configurations
- **Quick Mode**: CLI flags for non-interactive usage
- **Advanced Filtering**: Complex file selection patterns and rules
- **Template Marketplace**: Community template sharing

### Long-term Vision
- **Industry Standard**: The de-facto tool for structured LLM prompt generation in development
- **Ecosystem Integration**: Official templates for popular frameworks and tools
- **Team Features**: Shared template repositories and team-wide standardization
- **Performance Optimization**: Handle enterprise-scale monorepos efficiently

### Expansion Opportunities
- **IDE Integrations**: VS Code extension, Vim plugin, JetBrains integration
- **CI/CD Integration**: Automated prompt generation for code review workflows  
- **Template Analytics**: Usage patterns and effectiveness metrics
- **Professional Services**: Custom template development for enterprise clients

## Technical Considerations

### Platform Requirements
- **Target Platforms**: Windows, macOS, Linux
- **Browser/OS Support**: Any terminal with basic ANSI support; optimized for modern terminals
- **Performance Requirements**: <2s startup, <5s file scanning for 1000+ files, responsive UI (<16ms frame time)

### Technology Preferences  
- **Frontend**: Bubble Tea v2.0.0-beta.4 TUI framework with Lip Gloss styling
- **Backend**: Pure Go with standard library, no external dependencies for core functionality
- **Database**: File-based configuration and templates (TOML format)
- **Hosting/Infrastructure**: Distributed as single binary via GitHub releases

### Architecture Considerations
- **Repository Structure**: Standard Go project layout with internal packages
- **Service Architecture**: Single-binary application with modular internal architecture
- **Integration Requirements**: Git integration for .gitignore parsing, file type detection
- **Security/Compliance**: No network access required, all processing local, respect .gitignore patterns

## Constraints & Assumptions

### Constraints
- **Budget**: Open source development with no initial funding
- **Timeline**: 12-week development timeline from start to first stable release
- **Resources**: Single developer initially, with potential community contributions
- **Technical**: Must work in resource-constrained environments, no external service dependencies

### Key Assumptions
- Developers are comfortable with command-line tools and keyboard navigation
- Go single-binary distribution model is acceptable for target audience
- Terminal-based interface is preferred over GUI for integration with developer workflows
- Built-in templates cover 80% of common use cases without customization
- Cross-platform compatibility is essential for team adoption
- Bubble Tea framework will remain stable and well-maintained

## Risks & Open Questions

### Key Risks
- **Terminal Compatibility**: Rendering inconsistencies across different terminal environments and operating systems
- **Performance Degradation**: UI responsiveness with very large repositories (10k+ files)
- **Framework Stability**: Dependency on beta version of Bubble Tea v2 may introduce breaking changes
- **User Adoption**: Developers may prefer existing web-based or IDE-integrated solutions

### Open Questions
- Should the tool support non-Git repositories, or focus exclusively on Git-based projects?
- How complex should the built-in template system be versus relying on user customization?
- What is the optimal balance between features and simplicity for the MVP release?
- Should there be any integration with popular Git hosting services (GitHub, GitLab)?

### Areas Needing Further Research
- **Competitive Analysis**: Detailed comparison with existing prompt generation tools
- **User Interviews**: Validation of assumptions about developer workflows and pain points
- **Technical Feasibility**: Performance benchmarking with large, real-world repositories
- **Template Design**: User research on most valuable built-in template types

## Next Steps

### Immediate Actions
1. **Validate Core Assumptions**: Conduct interviews with 10-15 developers about current LLM prompt workflows
2. **Technical Proof of Concept**: Build basic file scanner and Bubble Tea navigation prototype  
3. **Template Research**: Analyze common LLM prompting patterns in developer workflows
4. **Competitive Analysis**: Research existing tools and identify differentiation opportunities
5. **Architecture Design**: Finalize Go project structure and component interfaces

### PM Handoff
This Project Brief provides the full context for **shotgun-cli**. Please start in 'PRD Generation Mode', review the brief thoroughly to work with the user to create the PRD section by section as the template indicates, asking for any necessary clarification or suggesting improvements.