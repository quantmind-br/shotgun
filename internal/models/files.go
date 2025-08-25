package models

import (
	"time"
)

// FileNode represents a file or directory in the project structure
type FileNode struct {
	Path        string      `json:"path"`
	Name        string      `json:"name"`
	IsDirectory bool        `json:"is_directory"`
	IsSelected  bool        `json:"is_selected"`
	IsIgnored   bool        `json:"is_ignored"`
	IsBinary    bool        `json:"is_binary"`
	Size        int64       `json:"size"`
	ModTime     time.Time   `json:"mod_time"`
	Children    []*FileNode `json:"children,omitempty"`
	Parent      *FileNode   `json:"-"`
}