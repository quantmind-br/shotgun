package spinner

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/diogopedro/shotgun/internal/components/common"
	"github.com/diogopedro/shotgun/internal/styles"
)

// SpinnerStyle represents different spinner animation styles
type SpinnerStyle string

const (
	SpinnerDots   SpinnerStyle = "dots"
	SpinnerLine   SpinnerStyle = "line"
	SpinnerCircle SpinnerStyle = "circle"
)

// Model represents the spinner component for loading states
type Model struct {
	spinner     spinner.Model
	loading     bool
	message     string
	style       lipgloss.Style
	started     time.Time
	minDuration time.Duration         // Anti-flicker: minimum display duration
	fallback    styles.FallbackConfig // Terminal capability fallback configuration
}

// MinSpinnerDuration is the minimum time a spinner should be displayed to avoid flicker
const MinSpinnerDuration = common.MinDisplayDuration

// New creates a new spinner component with the specified style
func New(style SpinnerStyle) Model {
	s := spinner.New()
	fallback := styles.NewFallbackConfig()

	// Configure spinner based on terminal capabilities and requested style
	spinnerChars := fallback.SpinnerCharacters()

	// Create custom spinner with appropriate characters
	s.Spinner = spinner.Spinner{
		Frames: spinnerChars,
		FPS:    time.Second / 80, // ~12.5 FPS as defined in common constants
	}

	// Set spinner colors to match app theme
	s.Style = common.StyleSpinner

	return Model{
		spinner:     s,
		loading:     false,
		message:     "",
		style:       common.StyleSpinnerMessage,
		minDuration: MinSpinnerDuration,
		fallback:    fallback,
	}
}

// Start begins the spinner animation
func (m *Model) Start() tea.Cmd {
	m.loading = true
	m.started = time.Now()
	return m.spinner.Tick
}

// Stop ends the spinner animation
func (m *Model) Stop() {
	m.loading = false
}

// SetMessage updates the loading message
func (m *Model) SetMessage(msg string) {
	m.message = msg
}

// IsLoading returns whether the spinner is currently active
func (m Model) IsLoading() bool {
	return m.loading
}

// ShouldHide determines if spinner can be hidden (anti-flicker)
func (m Model) ShouldHide() bool {
	if !m.loading {
		elapsed := time.Since(m.started)
		return elapsed >= m.minDuration
	}
	return false
}

// Update handles spinner tick messages and animations
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.loading {
		return m, nil
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the spinner with optional message
func (m Model) View() string {
	if !m.loading {
		return ""
	}

	var result string

	// Get the spinner view and sanitize it for terminal compatibility
	spinnerView := m.fallback.SanitizeText(m.spinner.View())

	// Render spinner with message
	if m.message != "" {
		// Sanitize message for terminal compatibility
		sanitizedMessage := m.fallback.SanitizeText(m.message)
		result = spinnerView + " " + m.style.Render(sanitizedMessage)
	} else {
		result = spinnerView
	}

	return result
}

// ViewWithCancel renders spinner with cancellation hint
func (m Model) ViewWithCancel() string {
	if !m.loading {
		return ""
	}

	view := m.View()
	if view == "" {
		return ""
	}

	cancelHint := common.StyleCancelHint.
		Render("\n\n[Press ESC to cancel]")

	return view + cancelHint
}

// LoadingTracker helps manage loading state timing for anti-flicker
type LoadingTracker struct {
	startTime time.Time
	minShown  bool
}

// NewLoadingTracker creates a new loading state tracker
func NewLoadingTracker() LoadingTracker {
	return LoadingTracker{
		startTime: time.Now(),
		minShown:  false,
	}
}

// ShouldHide determines if loading state can be hidden based on minimum duration
func (l *LoadingTracker) ShouldHide() bool {
	elapsed := time.Since(l.startTime)
	if elapsed >= MinSpinnerDuration {
		l.minShown = true
		return true
	}
	return false
}

// HasShownMinimum returns whether minimum display time has elapsed
func (l LoadingTracker) HasShownMinimum() bool {
	return l.minShown
}
