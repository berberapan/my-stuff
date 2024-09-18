package validator

import (
	"regexp"
	"testing"

	"github.com/berberapan/my-stuff/internal/assert"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name     string
		value    *Validator
		expected bool
	}{
		{
			name:     "Valid",
			value:    &Validator{},
			expected: true,
		},
		{
			name:     "Not valid, non-field error populated",
			value:    &Validator{NonFieldErrors: []string{"Test Error"}},
			expected: false,
		},
		{
			name:     "Not valid, field error populated",
			value:    &Validator{FieldErrors: map[string]string{"TestError": "Error"}},
			expected: false,
		},
		{
			name:     "Not valid, both error fields populated",
			value:    &Validator{NonFieldErrors: []string{"Test Error"}, FieldErrors: map[string]string{"TestError": "Error"}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testValue := tt.value.Valid()
			assert.Equal(t, testValue, tt.expected)
		})
	}
}

func TestAddNonFieldError(t *testing.T) {
	tests := []struct {
		name     string
		initial  []string
		addition string
		expected []string
	}{
		{
			name:     "Add error to empty slice",
			initial:  []string{},
			addition: "error",
			expected: []string{"error"},
		},
		{
			name:     "Add error to slice with exisiting error",
			initial:  []string{"error1"},
			addition: "error2",
			expected: []string{"error1", "error2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testValue := &Validator{NonFieldErrors: tt.initial}
			testValue.AddNonFieldError(tt.addition)
			assert.Equal(t, len(testValue.NonFieldErrors), len(tt.expected))
			assert.Equal(t, testValue.NonFieldErrors[0], tt.expected[0])
		})
	}
}

func TestNotBlank(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "Blank",
			value:    "",
			expected: false,
		},
		{
			name:     "Blank with whitespace",
			value:    "  ",
			expected: false,
		},
		{
			name:     "Not Blank",
			value:    "Testing",
			expected: true,
		},
		{
			name:     "Not Blank with leading whitespace",
			value:    " Testing",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testValue := NotBlank(tt.value)
			assert.Equal(t, testValue, tt.expected)
		})
	}
}

func TestMinChars(t *testing.T) {
	tests := []struct {
		name     string
		minimum  int
		value    string
		expected bool
	}{
		{
			name:     "Too short",
			minimum:  5,
			value:    "abcd",
			expected: false,
		},
		{
			name:     "Same number as minimum value",
			minimum:  4,
			value:    "abcd",
			expected: true,
		},
		{
			name:     "More than minimum value",
			minimum:  3,
			value:    "abcd",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testValue := MinChars(tt.value, tt.minimum)
			assert.Equal(t, testValue, tt.expected)
		})
	}
}

func TestMatchesRequiredPattern(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		pattern  regexp.Regexp
		expected bool
	}{
		{
			name:     "Matching value and pattern",
			value:    "test@example.com",
			pattern:  *validEmailRegex,
			expected: true,
		},
		{
			name:     "None matching value and pattern",
			value:    "testexample.com",
			pattern:  *validEmailRegex,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testValue := MatchesRequiredPattern(tt.value, &tt.pattern)
			assert.Equal(t, testValue, tt.expected)
		})
	}
}
