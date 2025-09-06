# Manual Cross-Platform Testing Checklist

This checklist provides comprehensive manual testing procedures for cross-platform terminal compatibility.

## Platform Test Matrix

### Windows Testing

#### Windows Command Prompt (CMD)
- [ ] **Terminal Setup**: Open Command Prompt (cmd.exe)
- [ ] **Application Launch**: `shotgun.exe --help`
- [ ] **Color Display**: Should show limited colors (16-color palette)
- [ ] **Unicode Characters**: Should show ASCII fallbacks (|/-\ instead of ⠋⠙⠹⠸)
- [ ] **Progress Bars**: Should show `[###---]` instead of `[█░░░]`
- [ ] **Keyboard Input**: 
  - [ ] ESC key works
  - [ ] Arrow keys work
  - [ ] Enter key works
  - [ ] F1-F10 keys may not work (expected)
- [ ] **Text Rendering**: No garbled Unicode characters
- [ ] **Application Exit**: Clean exit with Ctrl+C

#### Windows PowerShell 5.1
- [ ] **Terminal Setup**: Open Windows PowerShell 5.1
- [ ] **Application Launch**: `.\shotgun.exe --help`
- [ ] **Color Display**: Should show enhanced colors (256-color support)
- [ ] **Unicode Characters**: Limited Unicode, ASCII fallbacks for complex chars
- [ ] **Progress Bars**: Basic Unicode blocks may work
- [ ] **Keyboard Input**:
  - [ ] ESC key works
  - [ ] Arrow keys work
  - [ ] Enter key works
  - [ ] F1-F10 keys may conflict with PowerShell shortcuts
- [ ] **ISE Compatibility**: Test in PowerShell ISE if available
- [ ] **Application Exit**: Clean exit with Ctrl+C

#### Windows PowerShell 7.x
- [ ] **Terminal Setup**: Open PowerShell 7 (pwsh.exe)
- [ ] **Application Launch**: `.\shotgun.exe --help`
- [ ] **Color Display**: Should show true colors (16M color support)
- [ ] **Unicode Characters**: Full Unicode support including ⠋⠙⠹⠸
- [ ] **Progress Bars**: Full Unicode progress bars █░
- [ ] **Keyboard Input**:
  - [ ] ESC key works
  - [ ] Arrow keys work
  - [ ] Enter key works
  - [ ] F1-F10 keys should work (F1 may show help)
- [ ] **Cross-platform features**: All features should work
- [ ] **Application Exit**: Clean exit with Ctrl+C

#### Windows Terminal
- [ ] **Terminal Setup**: Open Windows Terminal
- [ ] **Application Launch**: `shotgun.exe --help`
- [ ] **Color Display**: Should show true colors with full palette
- [ ] **Unicode Characters**: Full Unicode support including all symbols
- [ ] **Progress Bars**: Perfect Unicode rendering
- [ ] **Keyboard Input**:
  - [ ] All keys work including F1-F10
  - [ ] ESC key works perfectly
  - [ ] Modern key combinations work
- [ ] **Font Rendering**: Test with different fonts (Cascadia Code, etc.)
- [ ] **Application Exit**: Clean exit with any method

### macOS Testing

#### macOS Terminal.app
- [ ] **Terminal Setup**: Open Terminal.app
- [ ] **Application Launch**: `./shotgun --help`
- [ ] **Color Display**: Should show true colors (16M colors)
- [ ] **Unicode Characters**: Full Unicode support
- [ ] **Progress Bars**: Perfect Unicode rendering
- [ ] **Keyboard Input**:
  - [ ] ESC key works
  - [ ] Arrow keys work
  - [ ] Enter key works
  - [ ] F1-F4 may be intercepted by system (expected)
  - [ ] F5-F10 should work
  - [ ] Cmd+C works for exit
- [ ] **Font Rendering**: Test with SF Mono and other fonts
- [ ] **Application Exit**: Clean exit with Cmd+C or ESC

#### iTerm2
- [ ] **Terminal Setup**: Open iTerm2
- [ ] **Application Launch**: `./shotgun --help`
- [ ] **Color Display**: Should show true colors with enhanced profiles
- [ ] **Unicode Characters**: Excellent Unicode support
- [ ] **Progress Bars**: Perfect rendering with custom fonts
- [ ] **Keyboard Input**:
  - [ ] All F1-F10 keys should work
  - [ ] ESC key works perfectly
  - [ ] Custom key mappings work if configured
- [ ] **Profile Testing**: Test with different iTerm2 profiles
- [ ] **Split Panes**: Test in split pane configurations
- [ ] **Application Exit**: Clean exit with any method

### Linux Testing

#### GNOME Terminal
- [ ] **Terminal Setup**: Open GNOME Terminal
- [ ] **Application Launch**: `./shotgun --help`
- [ ] **Color Display**: Should show true colors (16M colors)
- [ ] **Unicode Characters**: Full Unicode support
- [ ] **Progress Bars**: Perfect Unicode rendering
- [ ] **Keyboard Input**:
  - [ ] ESC key works
  - [ ] Arrow keys work
  - [ ] Enter key works
  - [ ] F1-F9 should work
  - [ ] F10 may activate menu bar (expected)
  - [ ] Ctrl+C works for exit
- [ ] **Theme Testing**: Test with different GNOME themes
- [ ] **Application Exit**: Clean exit with Ctrl+C

#### KDE Konsole
- [ ] **Terminal Setup**: Open Konsole
- [ ] **Application Launch**: `./shotgun --help`
- [ ] **Color Display**: Should show true colors with KDE integration
- [ ] **Unicode Characters**: Excellent Unicode support
- [ ] **Progress Bars**: Perfect rendering
- [ ] **Keyboard Input**:
  - [ ] All F1-F10 keys should work
  - [ ] ESC key works perfectly
  - [ ] KDE-specific shortcuts work
- [ ] **Profile Testing**: Test with different Konsole profiles
- [ ] **Split View**: Test in split view configurations
- [ ] **Application Exit**: Clean exit with any method

#### xterm
- [ ] **Terminal Setup**: Open xterm
- [ ] **Application Launch**: `./shotgun --help`  
- [ ] **Color Display**: Should show 256 colors
- [ ] **Unicode Characters**: Basic Unicode support
- [ ] **Progress Bars**: Should render correctly
- [ ] **Keyboard Input**:
  - [ ] ESC key works
  - [ ] Arrow keys work
  - [ ] Enter key works
  - [ ] F1-F10 keys should work
  - [ ] Ctrl+C works for exit
- [ ] **Configuration**: Test with different xterm configurations
- [ ] **Application Exit**: Clean exit with Ctrl+C

## Environment Variable Testing

### NO_COLOR Testing
Test on any Linux terminal:
- [ ] **Setup**: `export NO_COLOR=1`
- [ ] **Launch**: `./shotgun --help`
- [ ] **Verification**: Should show no colors, only monochrome output
- [ ] **Unicode**: Unicode should still work if terminal supports it
- [ ] **Cleanup**: `unset NO_COLOR`

### FORCE_COLOR Testing
Test on any terminal:
- [ ] **FORCE_COLOR=0**: `FORCE_COLOR=0 ./shotgun --help` (should be monochrome)
- [ ] **FORCE_COLOR=1**: `FORCE_COLOR=1 ./shotgun --help` (should be 8 colors)
- [ ] **FORCE_COLOR=2**: `FORCE_COLOR=2 ./shotgun --help` (should be 256 colors)
- [ ] **FORCE_COLOR=3**: `FORCE_COLOR=3 ./shotgun --help` (should be true color)

### TERM Variable Testing
Test on Linux/macOS:
- [ ] **vt100**: `TERM=vt100 ./shotgun --help` (should be monochrome, ASCII only)
- [ ] **xterm**: `TERM=xterm ./shotgun --help` (should be 16 colors)
- [ ] **xterm-256color**: `TERM=xterm-256color ./shotgun --help` (should be 256 colors)

## Unicode Rendering Tests

### Basic Unicode Test
In each terminal, verify these characters render correctly:
- [ ] **Progress Bars**: `█▉▊▋▌▍▎▏░`
- [ ] **Spinner**: `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
- [ ] **Box Drawing**: `┌─┐│└─┘`
- [ ] **Arrows**: `←↑→↓`
- [ ] **Status Symbols**: `✓✗⚠ⓘ`

### Fallback Verification
In terminals with limited Unicode (CMD, PowerShell 5.1):
- [ ] **Progress Bars**: Should show `[###---]` instead of `[█░░]`
- [ ] **Spinner**: Should show `|/-\` instead of `⠋⠙⠹⠸`
- [ ] **Box Drawing**: Should show `+--+|` instead of `┌─┐│`
- [ ] **Status Symbols**: Should show `+x!i` instead of `✓✗⚠ⓘ`

## Color Rendering Tests

### True Color Test (16M colors)
In modern terminals (Windows Terminal, iTerm2, GNOME Terminal):
- [ ] **Gradients**: Should show smooth color transitions
- [ ] **Syntax Highlighting**: Rich colors for different elements
- [ ] **Status Colors**: Distinct colors for success/warning/error

### Limited Color Test
In basic terminals (CMD, xterm):
- [ ] **Distinct Colors**: Should show clearly different colors for different states
- [ ] **Readability**: Text should be readable against backgrounds
- [ ] **No Color Bleeding**: Colors should not interfere with text

## Keyboard Input Tests

### Universal Keys (should work everywhere)
- [ ] **ESC**: Exits or cancels current operation
- [ ] **Enter**: Confirms selections
- [ ] **Arrow Keys**: Navigate menus/lists
- [ ] **Tab**: Navigate between elements
- [ ] **Ctrl+C**: Force quit/interrupt

### Function Keys (F1-F10)
Test where expected to work (Windows Terminal, iTerm2, Linux terminals):
- [ ] **F1**: Shows help (if implemented)
- [ ] **F2-F10**: Perform assigned functions (if implemented)

Test where may not work (CMD, PowerShell 5.1):
- [ ] **Graceful Handling**: Application doesn't crash with F-key input
- [ ] **Alternative Methods**: Alternative ways to access F-key functions

## Performance Testing

### Startup Time
- [ ] **Cold Start**: `time ./shotgun --help` (should be < 1 second)
- [ ] **Warm Start**: Run multiple times, verify consistent performance

### Memory Usage
- [ ] **Memory Footprint**: Monitor memory usage during operation
- [ ] **Memory Leaks**: Run for extended periods, verify stable memory

### Large Output Handling
- [ ] **Long Lists**: Test with applications that generate long output
- [ ] **Rapid Updates**: Test with fast-updating progress indicators

## Error Handling Tests

### Invalid Terminal Configurations
- [ ] **TERM=invalid**: `TERM=invalid ./shotgun --help`
- [ ] **Missing Terminal**: Test in minimal environments

### Signal Handling
- [ ] **Ctrl+C**: Should exit cleanly
- [ ] **SIGTERM**: Should handle termination gracefully
- [ ] **Window Resize**: Should handle terminal resize events

## Integration Testing

### SSH Testing
- [ ] **Local SSH**: `ssh localhost` then run application
- [ ] **Remote SSH**: Test over actual SSH connections
- [ ] **Terminal Forwarding**: Verify X11/terminal forwarding works

### Screen/Tmux Testing  
- [ ] **GNU Screen**: Run application inside screen session
- [ ] **tmux**: Run application inside tmux session
- [ ] **Nested Sessions**: Test nested screen/tmux sessions

### Container Testing
- [ ] **Docker**: Run in Docker container
- [ ] **Minimal Container**: Test in Alpine/minimal containers
- [ ] **Different Base Images**: Test Ubuntu, CentOS, etc.

## Documentation Verification

### Help System
- [ ] **--help flag**: Shows comprehensive help
- [ ] **Error Messages**: Clear, helpful error messages
- [ ] **Version Info**: `--version` shows correct information

### Platform-Specific Notes
- [ ] **Installation Instructions**: Work for each platform
- [ ] **Known Issues**: Documented issues match actual behavior
- [ ] **Workarounds**: Provided workarounds actually work

## Automated Test Verification

### CI Pipeline Results
- [ ] **All Platforms Pass**: GitHub Actions shows green for all platforms
- [ ] **Test Coverage**: Coverage reports show adequate coverage
- [ ] **Build Artifacts**: All platform binaries build successfully

### Local Test Suite
- [ ] **Unit Tests**: `go test ./...` passes on local platform
- [ ] **E2E Tests**: `go test ./e2e/...` passes on local platform
- [ ] **Integration Tests**: All integration tests pass

## Sign-off

### Platform Maintainers
- [ ] **Windows**: Tested and approved by Windows user
- [ ] **macOS**: Tested and approved by macOS user  
- [ ] **Linux**: Tested and approved by Linux user

### Test Evidence
- [ ] **Screenshots**: Captured for each platform/terminal combination
- [ ] **Screen Recordings**: Video evidence of key functionality
- [ ] **Test Logs**: Detailed logs of all test executions

### Release Approval
- [ ] **All Tests Pass**: Every checklist item completed successfully
- [ ] **Known Issues Documented**: Any issues properly documented
- [ ] **Performance Acceptable**: Performance meets requirements
- [ ] **Ready for Release**: Application ready for cross-platform distribution

---

## Testing Notes

Use this space to record any issues, observations, or notes during testing:

```
Date: ___________
Platform: ___________
Terminal: ___________
Tester: ___________

Issues Found:
- 

Observations:
- 

Performance Notes:
- 
```