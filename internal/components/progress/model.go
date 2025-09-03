package progress

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Model represents the progress indicator component
type Model struct {
	current int    // Current step (1-based)
	total   int    // Total steps
	width   int    // Available width for rendering
	titles  []string // Step titles
}

// NewModel creates a new progress indicator
func NewModel(current, total int, titles []string) Model {
	return Model{
		current: current,
		total:   total,
		width:   80, // Default width
		titles:  titles,
	}
}

// SetWidth sets the available width for rendering
func (m *Model) SetWidth(width int) {
	m.width = width
}

// SetCurrent updates the current step
func (m *Model) SetCurrent(current int) {
	if current >= 1 && current <= m.total {
		m.current = current
	}
}

// View renders the progress indicator
func (m Model) View() string {
	return m.renderProgress()
}

// renderProgress renders the complete progress indicator
func (m Model) renderProgress() string {
	var result strings.Builder
	
	// Render step counter (e.g., "Step 2 of 5")
	stepCounter := fmt.Sprintf("Step %d of %d", m.current, m.total)
	result.WriteString(m.styleStepCounter(stepCounter))
	result.WriteString("\n")
	
	// Render current step title if available
	if len(m.titles) >= m.current {
		title := m.titles[m.current-1]
		result.WriteString(m.styleStepTitle(title))
		result.WriteString("\n")
	}
	
	// Render progress bar
	progressBar := m.renderProgressBar()
	result.WriteString(progressBar)
	result.WriteString("\n")
	
	// Render step indicators
	stepIndicators := m.renderStepIndicators()
	result.WriteString(stepIndicators)
	
	return result.String()
}

// renderProgressBar creates a visual progress bar
func (m Model) renderProgressBar() string {
	if m.width < 20 {
		return ""
	}
	
	// Calculate progress percentage
	progress := float64(m.current-1) / float64(m.total-1)
	if m.total == 1 {
		progress = 1.0
	}
	
	// Available width for the bar (accounting for brackets and padding)
	barWidth := m.width - 10
	if barWidth < 10 {
		barWidth = 10
	}
	
	// Calculate filled and empty segments
	filledWidth := int(float64(barWidth) * progress)
	emptyWidth := barWidth - filledWidth
	
	// Build the bar
	var bar strings.Builder
	bar.WriteString("[")
	bar.WriteString(strings.Repeat("█", filledWidth))
	bar.WriteString(strings.Repeat("░", emptyWidth))
	bar.WriteString("]")
	
	// Add percentage
	percentage := fmt.Sprintf(" %.0f%%", progress*100)
	bar.WriteString(percentage)
	
	return m.styleProgressBar(bar.String())
}

// renderStepIndicators creates step-by-step indicators (1 → 2 → 3 → 4 → 5)
func (m Model) renderStepIndicators() string {
	var indicators strings.Builder
	
	for i := 1; i <= m.total; i++ {
		// Add step number with appropriate styling
		stepText := fmt.Sprintf("%d", i)
		
		if i == m.current {
			// Current step - highlighted
			indicators.WriteString(m.styleCurrentStep(stepText))
		} else if i < m.current {
			// Completed step
			indicators.WriteString(m.styleCompletedStep(stepText))
		} else {
			// Future step
			indicators.WriteString(m.styleFutureStep(stepText))
		}
		
		// Add arrow separator except for last step
		if i < m.total {
			if i < m.current {
				indicators.WriteString(m.styleCompletedArrow(" → "))
			} else {
				indicators.WriteString(m.styleFutureArrow(" → "))
			}
		}
	}
	
	return indicators.String()
}

// Styling functions using lipgloss

func (m Model) styleStepCounter(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")). // Blue
		Render(text)
}

func (m Model) styleStepTitle(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("14")). // Cyan
		Render(text)
}

func (m Model) styleProgressBar(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")). // Green
		Render(text)
}

func (m Model) styleCurrentStep(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).   // Black text
		Background(lipgloss.Color("12")). // Blue background
		Padding(0, 1).
		Render(text)
}

func (m Model) styleCompletedStep(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")). // Green
		Render("✓")
}

func (m Model) styleFutureStep(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Gray
		Render(text)
}

func (m Model) styleCompletedArrow(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")). // Green
		Render(text)
}

func (m Model) styleFutureArrow(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Gray
		Render(text)
}

// GetCurrent returns the current step number
func (m Model) GetCurrent() int {
	return m.current
}

// GetTotal returns the total number of steps
func (m Model) GetTotal() int {
	return m.total
}

// IsComplete returns whether all steps are completed
func (m Model) IsComplete() bool {
	return m.current >= m.total
}

// GetProgressPercent returns the progress as a percentage (0-100)
func (m Model) GetProgressPercent() float64 {
	if m.total <= 1 {
		return 100.0
	}
	return float64(m.current-1) / float64(m.total-1) * 100.0
}