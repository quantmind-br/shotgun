package template

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/diogopedro/shotgun/internal/models"
)

// Messages for template screen
type TemplatesLoadedMsg struct {
	Templates []models.Template
}

type TemplateLoadErrorMsg struct {
	Error error
}

type TemplateSelectedMsg struct {
	Template *models.Template
}

// TemplateModel represents the template selection screen state
type TemplateModel struct {
	// Template data
	templates []models.Template
	list      list.Model
	selected  *models.Template

	// Layout components
	viewport viewport.Model
	width    int
	height   int

	// Screen state
	loading bool
	err     error
	ready   bool

	// Detail panel state
	showDetails bool

	// Key mappings
	keyMap KeyMap
}

// NewTemplateModel creates a new template selection model
func NewTemplateModel() TemplateModel {
	// Create list with custom delegate for template display
	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(3) // Allow 3 lines per item for name, version, description
	delegate.SetSpacing(1)

	templateList := list.New([]list.Item{}, delegate, 0, 0)
	templateList.SetShowHelp(false)
	templateList.SetShowStatusBar(false)
	templateList.SetFilteringEnabled(false)
	templateList.Title = "Select Template"

	return TemplateModel{
		templates:   []models.Template{},
		list:        templateList,
		selected:    nil,
		viewport:    viewport.New(0, 0),
		loading:     true,
		err:         nil,
		ready:       false,
		showDetails: true,
		keyMap:      DefaultKeyMap(),
	}
}

// Init initializes the template model
func (m TemplateModel) Init() tea.Cmd {
	return nil
}

// UpdateSize updates the model dimensions for responsive layout
func (m *TemplateModel) UpdateSize(width, height int) {
	m.width = width
	m.height = height

	// Reserve space for title, help text, and borders
	listHeight := height - 6
	listWidth := width

	// If showing details panel, split the width
	if m.showDetails && width > 80 {
		listWidth = width * 2 / 3
		detailWidth := width - listWidth - 2 // 2 for border
		m.viewport.Width = detailWidth
		m.viewport.Height = listHeight
	}

	m.list.SetSize(listWidth, listHeight)
}

// SetTemplates updates the template list
func (m *TemplateModel) SetTemplates(templates []models.Template) {
	m.templates = templates
	items := make([]list.Item, len(templates))

	for i, template := range templates {
		items[i] = TemplateItem{Template: template}
	}

	m.list.SetItems(items)
	m.loading = false
	m.ready = true

	// Auto-select first template if available
	if len(templates) > 0 {
		m.selected = &templates[0]
	}
}

// GetSelected returns the currently selected template
func (m TemplateModel) GetSelected() *models.Template {
	return m.selected
}

// CanAdvance returns true if a template is selected
func (m TemplateModel) CanAdvance() bool {
	return m.selected != nil
}

// SetError sets an error state
func (m *TemplateModel) SetError(err error) {
	m.err = err
	m.loading = false
}

// IsLoading returns true if templates are currently loading
func (m TemplateModel) IsLoading() bool {
	return m.loading
}

// TemplateItem implements list.Item for template display
type TemplateItem struct {
	Template models.Template
}

// FilterValue returns the value to filter on
func (t TemplateItem) FilterValue() string {
	return t.Template.Name
}

// Title returns the template name with version
func (t TemplateItem) Title() string {
	if t.Template.Version != "" {
		return t.Template.Name + " v" + t.Template.Version
	}
	return t.Template.Name
}

// Description returns the template description
func (t TemplateItem) Description() string {
	return t.Template.Description
}
