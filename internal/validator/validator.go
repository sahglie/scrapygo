package validator

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

// Validator Define a new Validator struct which contains a map of validation error messages
// for our form fields.
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if the FieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already
	// initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) GetErrors() []string {
	fieldErrors := make([]string, 0)
	if v.FieldErrors == nil {
		return fieldErrors
	}

	for key, val := range v.FieldErrors {
		errMsg := fmt.Sprintf("%s: %s", key, val)
		fieldErrors = append(fieldErrors, errMsg)
	}

	return fieldErrors
}

func (v *Validator) FirstError() string {
	fieldErrors := v.GetErrors()
	if len(fieldErrors) > 0 {
		return fieldErrors[0]
	}

	return ""
}

// CheckField adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue returns true if a value is in a list of specific permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
