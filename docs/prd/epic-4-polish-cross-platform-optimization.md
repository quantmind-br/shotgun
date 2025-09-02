# Epic 4: Polish & Cross-Platform Optimization

Add the complete keyboard navigation system with F-key shortcuts, implement visual progress indicators, ensure cross-platform compatibility, and create the distribution pipeline for releasing the application.

## Story 4.1: Global Keyboard Navigation System

As a user,  
I want consistent keyboard shortcuts across all screens,  
so that I can navigate efficiently without learning different commands.

### Acceptance Criteria
1: F1-F10 keys properly mapped and handled globally
2: Help overlay (F1) shows context-sensitive shortcuts
3: Navigation between all 5 screens works smoothly
4: ESC key handling with confirmation dialog
5: Keyboard shortcuts don't conflict during text editing
6: Tab order logical for accessibility
7: All shortcuts documented in help screen

## Story 4.2: Progress Indicators & Loading States

As a user,  
I want visual feedback during long operations,  
so that I know the application is working.

### Acceptance Criteria
1: Spinner displays during initial file scanning
2: Progress bar for file reading operations
3: Loading state for template discovery
4: Size calculation progress in confirmation screen
5: All progress indicators use consistent styling
6: Operations remain cancellable with ESC
7: Smooth animations without flickering

## Story 4.3: Cross-Platform Testing & Compatibility

As a developer,  
I want the application to work consistently across platforms,  
so that all users have the same experience.

### Acceptance Criteria
1: Application runs on Windows PowerShell and CMD
2: Application runs on macOS Terminal and iTerm2
3: Application runs on common Linux terminals
4: Unicode characters display correctly on all platforms
5: Colors degrade gracefully on limited terminals
6: Keyboard shortcuts work across all environments
7: CI pipeline tests on Windows, macOS, and Linux

## Story 4.4: Init Command & Shotgunignore Support

As a user,  
I want to create a .shotgunignore file easily,  
so that I can customize which files are excluded.

### Acceptance Criteria
1: 'shotgun init' command creates .shotgunignore file
2: Template .shotgunignore includes common patterns
3: File created only if it doesn't exist
4: Success message confirms file creation
5: Command available via Cobra CLI framework
6: Help text explains usage and purpose
7: Integration with main file scanner confirmed

## Story 4.5: Binary Distribution Pipeline

As a maintainer,  
I want automated builds for all platforms,  
so that users can easily download and use the application.

### Acceptance Criteria
1: GitHub Actions workflow builds for Linux, macOS, Windows
2: Binaries created for amd64 and arm64 architectures where applicable
3: Version information embedded in binary
4: Artifacts uploaded to GitHub Releases
5: Checksums generated for verification
6: README includes installation instructions
7: Binary size optimized with appropriate build flags
