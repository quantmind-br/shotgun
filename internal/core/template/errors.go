package template

import (
	"fmt"
)

// Error types for different failure categories
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeParsing
	ErrorTypeValidation
	ErrorTypeFileAccess
	ErrorTypePathTraversal
	ErrorTypeContentSize
	ErrorTypeEmbedding
)

// String returns a string representation of the error type
func (e ErrorType) String() string {
	switch e {
	case ErrorTypeParsing:
		return "parsing"
	case ErrorTypeValidation:
		return "validation"
	case ErrorTypeFileAccess:
		return "file_access"
	case ErrorTypePathTraversal:
		return "path_traversal"
	case ErrorTypeContentSize:
		return "content_size"
	case ErrorTypeEmbedding:
		return "embedding"
	default:
		return "unknown"
	}
}

// TemplateError represents an error with additional context
type TemplateError struct {
	Type         ErrorType
	TemplatePath string
	Message      string
	Cause        error
}

// Error implements the error interface
func (e *TemplateError) Error() string {
	if e.TemplatePath != "" {
		return fmt.Sprintf("template error [%s] in '%s': %s", e.Type, e.TemplatePath, e.Message)
	}
	return fmt.Sprintf("template error [%s]: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *TemplateError) Unwrap() error {
	return e.Cause
}

// NewTemplateError creates a new template error
func NewTemplateError(errorType ErrorType, templatePath, message string, cause error) *TemplateError {
	return &TemplateError{
		Type:         errorType,
		TemplatePath: templatePath,
		Message:      message,
		Cause:        cause,
	}
}

// NewParsingError creates a parsing error
func NewParsingError(templatePath string, cause error) *TemplateError {
	return NewTemplateError(ErrorTypeParsing, templatePath,
		"failed to parse TOML template", cause)
}

// NewValidationError creates a validation error
func NewValidationError(templatePath, message string) *TemplateError {
	return NewTemplateError(ErrorTypeValidation, templatePath, message, nil)
}

// NewFileAccessError creates a file access error
func NewFileAccessError(templatePath string, cause error) *TemplateError {
	return NewTemplateError(ErrorTypeFileAccess, templatePath,
		"failed to access template file", cause)
}

// NewPathTraversalError creates a path traversal security error
func NewPathTraversalError(templatePath string) *TemplateError {
	return NewTemplateError(ErrorTypePathTraversal, templatePath,
		"path traversal attempt detected", nil)
}

// NewContentSizeError creates a content size error
func NewContentSizeError(templatePath string, size int) *TemplateError {
	return NewTemplateError(ErrorTypeContentSize, templatePath,
		fmt.Sprintf("template file too large: %d bytes", size), nil)
}

// NewEmbeddingError creates an embedding error
func NewEmbeddingError(message string, cause error) *TemplateError {
	return NewTemplateError(ErrorTypeEmbedding, "",
		fmt.Sprintf("embedded template system error: %s", message), cause)
}

// ErrorAggregator collects multiple errors during batch operations
type ErrorAggregator struct {
	errors []error
}

// Add adds an error to the aggregator
func (ea *ErrorAggregator) Add(err error) {
	if err != nil {
		ea.errors = append(ea.errors, err)
	}
}

// HasErrors returns true if any errors were collected
func (ea *ErrorAggregator) HasErrors() bool {
	return len(ea.errors) > 0
}

// Count returns the number of collected errors
func (ea *ErrorAggregator) Count() int {
	return len(ea.errors)
}

// Errors returns all collected errors
func (ea *ErrorAggregator) Errors() []error {
	return ea.errors
}

// Error returns a combined error message, implementing the error interface
func (ea *ErrorAggregator) Error() string {
	if len(ea.errors) == 0 {
		return "no errors"
	}

	if len(ea.errors) == 1 {
		return ea.errors[0].Error()
	}

	return fmt.Sprintf("multiple errors occurred (%d): %v", len(ea.errors), ea.errors)
}

// NewErrorAggregator creates a new error aggregator
func NewErrorAggregator() *ErrorAggregator {
	return &ErrorAggregator{
		errors: make([]error, 0),
	}
}
