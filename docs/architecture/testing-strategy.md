# Testing Strategy

### Testing Pyramid

```plaintext
                E2E Tests
               /         \
          Integration Tests
         /                 \
    TUI Unit Tests    Core Unit Tests
```

### Test Organization

#### Frontend Tests

```plaintext
internal/ui/
├── app/
│   ├── app_test.go              # Application shell tests
│   ├── navigation_test.go       # Screen navigation tests
│   └── keybindings_test.go      # Key handler tests
├── screens/
│   ├── filetree/
│   │   ├── screen_test.go       # File tree screen tests
│   │   └── integration_test.go  # File tree integration tests
│   ├── templates/
│   │   ├── screen_test.go       # Template selection tests
│   │   └── mock_templates_test.go # Template mocking tests
│   └── confirmation/
│       ├── screen_test.go       # Confirmation screen tests
│       └── generation_test.go   # Generation workflow tests
├── components/
│   ├── filetree_test.go         # File tree component tests
│   ├── texteditor_test.go       # Text editor component tests
│   └── progressbar_test.go      # Progress bar component tests
└── styles/
    └── theme_test.go            # Theme and styling tests
```

#### Backend Tests

```plaintext
internal/core/
├── scanner/
│   ├── scanner_test.go          # File scanning unit tests
│   ├── concurrent_test.go       # Concurrency tests
│   └── benchmark_test.go        # Performance benchmarks
├── template/
│   ├── engine_test.go           # Template processing tests
│   ├── validation_test.go       # Template validation tests
│   └── functions_test.go        # Custom function tests
├── ignore/
│   ├── processor_test.go        # Ignore rules tests
│   └── patterns_test.go         # Pattern matching tests
└── output/
    ├── generator_test.go        # Output generation tests
    └── writer_test.go           # File writing tests
```

#### E2E Tests

```plaintext
test/e2e/
├── workflows/
│   ├── complete_generation_test.go    # Full workflow tests
│   ├── session_management_test.go     # Session persistence tests
│   └── error_recovery_test.go         # Error handling tests
├── fixtures/
│   ├── sample_project/                # Test project structures
│   ├── custom_templates/              # Test template files
│   └── expected_outputs/              # Golden file outputs
└── helpers/
    ├── tui_simulator.go               # TUI interaction simulation
    └── project_builder.go             # Test project creation
```

### Test Examples

#### Frontend Component Test

```go
func TestFileTreeComponent_Selection(t *testing.T) {
    // Setup test data
    rootNode := &models.FileNode{
        Path:        "test-project",
        Name:        "test-project",
        IsDirectory: true,
        IsSelected:  true,
        Children: []*models.FileNode{
            {Path: "main.go", Name: "main.go", IsSelected: true},
            {Path: "README.md", Name: "README.md", IsSelected: true},
        },
    }
    
    // Create component
    component := components.NewFileTreeComponent()
    component.SetRoot(rootNode)
    
    tests := []struct {
        name           string
        keypress       string
        expectedCursor int
        expectedSelected bool
    }{
        {
            name:           "arrow down moves cursor",
            keypress:       "down",
            expectedCursor: 1,
            expectedSelected: true,
        },
        {
            name:           "space toggles selection",
            keypress:       " ",
            expectedCursor: 0,
            expectedSelected: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Simulate key press
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.keypress[0]}}
            model, _ := component.Update(msg)
            
            // Verify expectations
            updatedComponent := model.(components.FileTreeComponent)
            assert.Equal(t, tt.expectedCursor, updatedComponent.Cursor())
            assert.Equal(t, tt.expectedSelected, updatedComponent.CurrentNode().IsSelected)
        })
    }
}
```

#### Backend API Test

```go
func TestTemplateEngine_ProcessTemplate(t *testing.T) {
    // Setup template engine
    engine := template.NewEngine()
    
    // Test template
    tmpl := &models.Template{
        ID:      "test-template",
        Name:    "Test Template",
        Content: "Task: {{.Task}}\n\nFiles:\n{{range .FileContents}}{{.}}\n{{end}}",
    }
    
    // Test data
    data := &template.TemplateData{
        Task: "Debug memory leak",
        FileContents: map[string]string{
            "main.go": "package main\n\nfunc main() {}\n",
            "README.md": "# Test Project\n",
        },
    }
    
    tests := []struct {
        name     string
        template *models.Template
        data     *template.TemplateData
        want     string
        wantErr  bool
    }{
        {
            name:     "valid template processing",
            template: tmpl,
            data:     data,
            want:     "Task: Debug memory leak\n\nFiles:\npackage main\n\nfunc main() {}\n\n# Test Project\n\n",
            wantErr:  false,
        },
        {
            name: "invalid template syntax",
            template: &models.Template{
                ID:      "invalid",
                Content: "{{.InvalidSyntax",
            },
            data:    data,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := engine.ProcessTemplate(tt.template, tt.data)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, result)
        })
    }
}
```

#### E2E Test

```go
func TestCompleteWorkflow_E2E(t *testing.T) {
    // Create test project
    projectDir := test.CreateSampleProject(t, "test-project", []string{
        "main.go",
        "README.md",
        ".gitignore",
        "internal/core.go",
    })
    defer os.RemoveAll(projectDir)
    
    // Change to project directory
    originalDir, _ := os.Getwd()
    os.Chdir(projectDir)
    defer os.Chdir(originalDir)
    
    // Create test application
    app := ui.NewApp()
    app.SetTestMode(true)
    
    // Run complete workflow simulation
    tm := teatest.NewTestModel(t, app, teatest.WithInitialTermSize(120, 40))
    
    // Step 1: File selection
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "File Selection [1/5]")
    })
    
    // Deselect .gitignore file
    tm.Send(tea.KeyMsg{Type: tea.KeyDown})    // Navigate to .gitignore
    tm.Send(tea.KeyMsg{Type: tea.KeyDown})
    tm.Send(tea.KeyMsg{Type: tea.KeySpace})   // Deselect
    tm.Send(tea.KeyMsg{Type: tea.KeyF3})      // Next screen
    
    // Step 2: Template selection
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Template Selection [2/5]")
    })
    
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})   // Select default template
    
    // Step 3: Task input
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Task Description [3/5]")
    })
    
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Debug memory leak in core module")})
    tm.Send(tea.KeyMsg{Type: tea.KeyF3})      // Next screen
    
    // Step 4: Rules input (skip)
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Rules")
    })
    
    tm.Send(tea.KeyMsg{Type: tea.KeyF4})      // Skip rules
    
    // Step 5: Confirmation and generation
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Confirm Generation [5/5]")
    })
    
    tm.Send(tea.KeyMsg{Type: tea.KeyF10})     // Generate prompt
    
    // Verify output file was created
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return strings.Contains(string(bts), "Successfully generated")
    })
    
    // Verify file contents
    files, _ := filepath.Glob("shotgun_prompt_*.md")
    require.Len(t, files, 1, "Expected one generated prompt file")
    
    content, err := ioutil.ReadFile(files[0])
    require.NoError(t, err)
    
    // Verify prompt contains expected elements
    assert.Contains(t, string(content), "Debug memory leak in core module")
    assert.Contains(t, string(content), "main.go")
    assert.Contains(t, string(content), "README.md")
    assert.NotContains(t, string(content), ".gitignore") // Should be excluded
}
```
