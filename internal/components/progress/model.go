package progress

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/diogopedro/shotgun/internal/components/common"
	"github.com/diogopedro/shotgun/internal/styles"
)

// Model represents the progress indicator component
type Model struct {
	current    int                   // Current step (1-based)
	total      int                   // Total steps
	width      int                   // Available width for rendering
	titles     []string              // Step titles
	percentage float64               // Progress percentage (0-100)
	message    string                // Current progress message
	startTime  time.Time             // When progress started
	bytesRead  int64                 // Bytes processed so far
	totalBytes int64                 // Total bytes to process
	showETA    bool                  // Whether to show ETA
	fileCount  int                   // Current file count
	totalFiles int                   // Total file count
	fallback   styles.FallbackConfig // Terminal capability fallback configuration
}

// NewModel creates a new progress indicator
func NewModel(current, total int, titles []string) Model {
	return Model{
		current:   current,
		total:     total,
		width:     80, // Default width
		titles:    titles,
		startTime: time.Now(),
		showETA:   true,
		fallback:  styles.NewFallbackConfig(),
	}
}

// NewFileProgressModel creates a progress bar for file operations
func NewFileProgressModel(totalFiles int) Model {
	return Model{
		current:    0,
		total:      totalFiles,
		width:      80,
		startTime:  time.Now(),
		showETA:    true,
		totalFiles: totalFiles,
		fileCount:  0,
		fallback:   styles.NewFallbackConfig(),
	}
}

// NewBytesProgressModel creates a progress bar for byte-based operations
func NewBytesProgressModel(totalBytes int64) Model {
	return Model{
		current:    0,
		total:      100, // Percentage based
		width:      80,
		startTime:  time.Now(),
		showETA:    true,
		totalBytes: totalBytes,
		bytesRead:  0,
		fallback:   styles.NewFallbackConfig(),
	}
}

// SetWidth sets the available width for rendering
func (m *Model) SetWidth(width int) {
	m.width = width
}

// SetCurrent updates the current step and calculates percentage
func (m *Model) SetCurrent(current int) {
	if current >= 0 && current <= m.total {
		m.current = current
		if m.total > 0 {
			m.percentage = float64(current) / float64(m.total) * 100.0
		}
	}
}

// SetBytes updates bytes processed and calculates percentage
func (m *Model) SetBytes(read, total int64) {
	m.bytesRead = read
	m.totalBytes = total
	if total > 0 {
		m.percentage = float64(read) / float64(total) * 100.0
		m.current = int(m.percentage)
	}
}

// SetFileCount updates file processing count
func (m *Model) SetFileCount(current, total int) {
	m.fileCount = current
	m.totalFiles = total
	if total > 0 {
		m.percentage = float64(current) / float64(total) * 100.0
		m.current = current
		m.total = total
	}
}

// SetMessage updates the progress message
func (m *Model) SetMessage(message string) {
	m.message = message
}

// GetETA calculates estimated time remaining
func (m Model) GetETA() time.Duration {
	if m.percentage <= 0 || m.percentage >= 100 {
		return 0
	}

	elapsed := time.Since(m.startTime)
	if elapsed.Seconds() < 1 {
		return 0 // Too early to estimate
	}

	rate := m.percentage / elapsed.Seconds()
	if rate <= 0 {
		return 0
	}

	remaining := (100.0 - m.percentage) / rate
	return time.Duration(remaining * float64(time.Second))
}

// IncrementFile increments the file counter
func (m *Model) IncrementFile() {
	m.SetFileCount(m.fileCount+1, m.totalFiles)
}

// AddBytes adds to the bytes processed
func (m *Model) AddBytes(bytes int64) {
	m.SetBytes(m.bytesRead+bytes, m.totalBytes)
}

// Update handles bubble tea messages (for smooth animations)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case tea.Msg:
		// Handle any progress-specific messages here
		// For now, just return the model as-is
		return m, nil
	}
	return m, nil
}

// View renders the progress indicator
func (m Model) View() string {
	return m.renderProgress()
}

// ViewWithMessage renders progress with a custom message
func (m Model) ViewWithMessage(message string) string {
	oldMsg := m.message
	m.message = m.fallback.SanitizeText(message)
	result := m.renderProgress()
	m.message = oldMsg // Restore original message
	return result
}

// ViewCompact renders a compact version of the progress
func (m Model) ViewCompact() string {
	if m.totalFiles > 0 {
		return fmt.Sprintf("[%s] %d/%d files (%.1f%%)",
			m.renderSimpleBar(20), m.fileCount, m.totalFiles, m.percentage)
	}
	return fmt.Sprintf("[%s] %.1f%%", m.renderSimpleBar(20), m.percentage)
}

// renderProgress renders the complete progress indicator
func (m Model) renderProgress() string {
	var result strings.Builder

	// Show custom message if set
	if m.message != "" {
		sanitizedMessage := m.fallback.SanitizeText(m.message)
		result.WriteString(m.styleStepTitle(sanitizedMessage))
		result.WriteString("\n")
	} else if len(m.titles) > 0 {
		// Render step counter (e.g., "Step 2 of 5")
		stepCounter := fmt.Sprintf("Step %d of %d", m.current, m.total)
		result.WriteString(m.styleStepCounter(stepCounter))
		result.WriteString("\n")

		// Render current step title if available
		if len(m.titles) >= m.current && m.current > 0 {
			title := m.titles[m.current-1]
			sanitizedTitle := m.fallback.SanitizeText(title)
			result.WriteString(m.styleStepTitle(sanitizedTitle))
			result.WriteString("\n")
		}
	}

	// Render progress bar with enhanced info
	progressBar := m.renderEnhancedProgressBar()
	result.WriteString(progressBar)
	result.WriteString("\n")

	// Show file/byte details if available
	if m.totalFiles > 0 {
		fileInfo := fmt.Sprintf("%d/%d files", m.fileCount, m.totalFiles)
		result.WriteString(m.styleProgressInfo(fileInfo))
		if m.showETA && m.percentage > 0 && m.percentage < 100 {
			eta := m.GetETA()
			if eta > 0 {
				etaText := fmt.Sprintf(" • ETA: %v", formatDuration(eta))
				result.WriteString(m.styleETA(etaText))
			}
		}
		result.WriteString("\n")
	} else if m.totalBytes > 0 {
		byteInfo := fmt.Sprintf("%s/%s", formatBytes(m.bytesRead), formatBytes(m.totalBytes))
		result.WriteString(m.styleProgressInfo(byteInfo))
		if m.showETA && m.percentage > 0 && m.percentage < 100 {
			eta := m.GetETA()
			if eta > 0 {
				etaText := fmt.Sprintf(" • ETA: %v", formatDuration(eta))
				result.WriteString(m.styleETA(etaText))
			}
		}
		result.WriteString("\n")
	}

	// Render step indicators only for traditional step-based progress
	if len(m.titles) > 0 && m.totalFiles == 0 && m.totalBytes == 0 {
		stepIndicators := m.renderStepIndicators()
		result.WriteString(stepIndicators)
	}

	return result.String()
}

// renderProgressBar creates a visual progress bar (legacy method)
func (m Model) renderProgressBar() string {
	return m.renderEnhancedProgressBar()
}

// renderEnhancedProgressBar creates a visual progress bar with percentage
func (m Model) renderEnhancedProgressBar() string {
	if m.width < 20 {
		return ""
	}

	// Use the calculated percentage or calculate from current/total
	progress := m.percentage / 100.0
	if m.percentage == 0 && m.total > 0 {
		progress = float64(m.current) / float64(m.total)
	}

	// Ensure progress is within bounds
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	// Available width for the bar (accounting for brackets and percentage display)
	percentageText := fmt.Sprintf(" %.1f%%", progress*100)
	barWidth := m.width - len(percentageText) - 2 // 2 for brackets
	if barWidth < 10 {
		barWidth = 10
	}

	// Calculate filled and empty segments
	filledWidth := int(float64(barWidth) * progress)
	emptyWidth := barWidth - filledWidth

	// Get appropriate progress characters based on terminal capabilities
	filledChar, emptyChar := m.fallback.ProgressCharacters()

	// Build the bar
	var bar strings.Builder
	bar.WriteString("[")
	bar.WriteString(strings.Repeat(filledChar, filledWidth))
	bar.WriteString(strings.Repeat(emptyChar, emptyWidth))
	bar.WriteString("]")
	bar.WriteString(percentageText)

	return m.styleProgressBar(bar.String())
}

// renderSimpleBar creates a simple progress bar for compact display
func (m Model) renderSimpleBar(width int) string {
	progress := m.percentage / 100.0
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filledWidth := int(float64(width) * progress)
	emptyWidth := width - filledWidth

	// Get appropriate progress characters based on terminal capabilities
	filledChar, emptyChar := m.fallback.ProgressCharacters()

	// Create the bar and sanitize it for terminal compatibility
	bar := strings.Repeat(filledChar, filledWidth) + strings.Repeat(emptyChar, emptyWidth)
	return m.fallback.SanitizeText(bar)
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
			arrow := m.fallback.StatusSymbols().Arrow
			if i < m.current {
				indicators.WriteString(m.styleCompletedArrow(" " + arrow + " "))
			} else {
				indicators.WriteString(m.styleFutureArrow(" " + arrow + " "))
			}
		}
	}

	return indicators.String()
}

// Styling functions using lipgloss

func (m Model) styleStepCounter(text string) string {
	return common.StyleTitle.Render(text)
}

func (m Model) styleStepTitle(text string) string {
	return common.StyleSubtitle.Render(text)
}

func (m Model) styleProgressBar(text string) string {
	return common.StyleProgressBar.Render(text)
}

func (m Model) styleCurrentStep(text string) string {
	return common.StyleHighlight.Render(text)
}

func (m Model) styleCompletedStep(text string) string {
	return common.StyleSuccess.Render(common.ProgressCheck)
}

func (m Model) styleFutureStep(text string) string {
	return common.StyleTextDim.Render(text)
}

func (m Model) styleCompletedArrow(text string) string {
	return common.StyleSuccess.Render(text)
}

func (m Model) styleFutureArrow(text string) string {
	return common.StyleTextDim.Render(text)
}

func (m Model) styleProgressInfo(text string) string {
	return common.StyleProgressInfo.Render(text)
}

func (m Model) styleETA(text string) string {
	return common.StyleProgressETA.Render(text)
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
	if m.percentage > 0 {
		return m.percentage
	}
	if m.total <= 0 {
		return 0.0
	}
	return float64(m.current) / float64(m.total) * 100.0
}

// formatBytes formats byte count into human-readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration formats duration into human-readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%02ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%02dm", int(d.Hours()), int(d.Minutes())%60)
}
