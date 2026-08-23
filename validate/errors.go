package validate

import (
	"fmt"
	"strings"
)

// FieldError describes a single validation failure.
type FieldError struct {
	Field string
	Value any
	Rule  string
}

func (e FieldError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: rule %q failed", e.Field, e.Rule)
	}
	return fmt.Sprintf("rule %q failed", e.Rule)
}

// ValidationError aggregates one or more field failures.
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	if len(e.Fields) == 0 {
		return "validation failed"
	}
	parts := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		parts[i] = f.Error()
	}
	return "validation failed: " + strings.Join(parts, "; ")
}

func (e *ValidationError) Is(target error) bool {
	_, ok := target.(*ValidationError)
	return ok
}
