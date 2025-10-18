package exceptions

import (
	"bytes"
	"strings"
)

type ValidationErrors struct {
	FieldErrors []FieldError
}

func (e *ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	for i := 0; i < len(e.FieldErrors); i++ {
		buff.WriteString(e.FieldErrors[i].Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{}
}

func (e *ValidationErrors) AddError(field, tag string) {
	e.FieldErrors = append(e.FieldErrors, *NewFieldError(field, tag))
}
