# Testing Strategy

## Testing Pyramid
```
        E2E Tests
       /        \
    Integration Tests
    /              \
Frontend Unit  Backend Unit
```

## Test Organization

### Frontend Tests
```
internal/screens/filetree/
├── model_test.go
├── update_test.go
└── view_test.go

internal/components/tree/
├── tree_test.go
└── testdata/
    └── sample_tree.json
```

### Backend Tests
```
internal/core/scanner/
├── scanner_test.go
├── ignore_test.go
├── binary_test.go
└── testdata/
    ├── .gitignore
    └── sample_files/

internal/core/template/
├── engine_test.go
├── discovery_test.go
└── testdata/
    └── templates/
```

### E2E Tests
```
e2e/
├── full_flow_test.go
├── navigation_test.go
└── fixtures/
    └── test_project/
```

## Test Examples

### Frontend Component Test
```go
func TestFileTreeToggleSelection(t *testing.T) {
    model := NewFileTreeModel()
    model.items = []FileItem{
        {Path: "/test.go", Name: "test.go", IsSelected: false},
    }
    model.cursor = 0
    
    // Toggle selection
    model.toggleSelection()
    
    assert.True(t, model.items[0].IsSelected)
    assert.Contains(t, model.selected, "/test.go")
}
```

### Backend API Test
```go
func TestScannerRespectsGitignore(t *testing.T) {
    scanner := NewScanner()
    scanner.LoadGitignore([]string{"*.log", "node_modules/"})
    
    files := []string{}
    for file := range scanner.ScanDirectory("./testdata") {
        files = append(files, file.Path)
    }
    
    assert.NotContains(t, files, "test.log")
    assert.NotContains(t, files, "node_modules")
}
```

### E2E Test
```go
func TestCompletePromptGeneration(t *testing.T) {
    app := teatest.NewTestModel(t, app.New())
    
    // Navigate through screens
    app.Send(tea.KeyMsg{Type: tea.KeySpace}) // Select file
    app.Send(tea.KeyMsg{Type: tea.KeyF3})    // Next screen
    app.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Select template
    app.Type("Fix authentication bug")        // Enter task
    app.Send(tea.KeyMsg{Type: tea.KeyCtrlEnter})
    app.Send(tea.KeyMsg{Type: tea.KeyF4})    // Skip rules
    app.Send(tea.KeyMsg{Type: tea.KeyF10})   // Generate
    
    // Verify output file created
    _, err := os.Stat("shotgun_prompt_*.md")
    assert.NoError(t, err)
}
```
