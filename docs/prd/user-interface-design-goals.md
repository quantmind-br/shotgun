# User Interface Design Goals

### Overall UX Vision
Minimalist, monochrome terminal interface prioritizing speed and keyboard efficiency. The design follows a clean, distraction-free aesthetic with subtle accent colors (soft mint green #6ee7b7) for focus states and warm amber (#fbbf24) for highlights. Every interaction is optimized for professional developers who value precision and workflow efficiency over visual complexity.

### Key Interaction Paradigms
- **Keyboard-First Navigation**: 100% keyboard operation with global F-key shortcuts (F1-help, F2-back, F3-forward, ESC-exit)
- **Contextual Controls**: Screen-specific shortcuts that change based on current mode (navigation vs. editing)
- **Progressive Disclosure**: 5-screen wizard that maintains state while allowing free navigation between steps
- **Immediate Feedback**: Real-time validation, character counts, and visual indicators for all user actions

### Core Screens and Views
- **File Tree Selection** (1/5): Hierarchical file browser with checkbox states, initially all-selected with deselection workflow
- **Template Selection** (2/5): Vertical list selector showing both built-in and custom templates with metadata
- **Task Input Editor** (3/5): Multiline text editor with UTF-8 support, word wrap, and character counting
- **Rules Input Editor** (4/5): Optional multiline editor with skip functionality and visual "optional" indicators  
- **Confirmation Summary** (5/5): Complete review with size estimation, progress bars, and final generation step

### Accessibility: WCAG AA
Full keyboard navigation compliance, high contrast monochrome color scheme, clear focus indicators, and graceful degradation for limited terminals. Screen reader compatibility through semantic text structure and meaningful state descriptions.

### Branding
Early 1900s monochrome aesthetic with clean typography and minimal visual elements. Uses ASCII art borders, simple progress indicators, and restrained color palette to maintain professional, timeless appearance that works across all terminal environments.

### Target Device and Platforms: Cross-Platform
Optimized for professional terminal environments including Windows PowerShell, ConPTY, Linux bash/zsh, macOS Terminal/iTerm2, and modern terminals like WezTerm. Responsive design supporting minimum 80x24 terminal size with graceful scaling for larger displays.
