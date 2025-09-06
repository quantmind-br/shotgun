package generate

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/diogopedro/shotgun/internal/components/common"
	"github.com/diogopedro/shotgun/internal/core/builder"
)

// StartGenerationCmd starts the async prompt generation process
func StartGenerationCmd(config builder.GenerationConfig) tea.Cmd {
	return func() tea.Msg {
		generator := builder.NewPromptGenerator()

		// Create progress callback that sends progress messages
		progressCallback := func(stage string, progress float64) {
			// This would ideally send progress updates, but we'll handle it in the async command
		}

		return generator.GenerateAsync(config, progressCallback)
	}
}

// WritePromptToFileCmd writes the generated prompt to a file
func WritePromptToFileCmd(result *builder.GeneratedPrompt) tea.Cmd {
	return func() tea.Msg {
		if result == nil {
			return FileWriteCompleteMsg{
				Result:     nil,
				OutputFile: "",
				Error:      fmt.Errorf("no result to write"),
			}
		}

		writer := builder.NewFileWriter()
		outputFile, err := writer.WritePromptFile(result.Content, "")

		return FileWriteCompleteMsg{
			Result:     result,
			OutputFile: outputFile,
			Error:      err,
		}
	}
}

// CancelGenerationCmd cancels the current generation process
func CancelGenerationCmd() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual cancellation mechanism
		return GenerationCancelledMsg{}
	}
}

// NavigateBackCmd navigates back to the previous screen
func NavigateBackCmd() tea.Cmd {
	return func() tea.Msg {
		return NavigateBackMsg{}
	}
}

// NavigateToFileTreeCmd navigates to the file tree screen
func NavigateToFileTreeCmd() tea.Cmd {
	return func() tea.Msg {
		return NavigateToFileTreeMsg{}
	}
}

// OpenFileCmd opens a file using the system's default application
func OpenFileCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", "", filePath)
		case "darwin":
			cmd = exec.Command("open", filePath)
		case "linux":
			cmd = exec.Command("xdg-open", filePath)
		default:
			return fmt.Errorf("unsupported operating system for opening files")
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		return nil
	}
}

// RetryGenerationCmd retries the generation process
func RetryGenerationCmd() tea.Cmd {
	return func() tea.Msg {
		return RetryGenerationMsg{}
	}
}

// ProgressTickCmd provides periodic progress updates during generation
func ProgressTickCmd() tea.Cmd {
	return tea.Tick(common.ProgressUpdateRate, func(time.Time) tea.Msg {
		return struct{}{} // Generic tick message for spinner updates
	})
}
