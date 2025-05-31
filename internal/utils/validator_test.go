package utils

import (
	"testing"
)

func TestValidateTaskName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid name", "test_task", true},
		{"valid chinese name", "测试任务", true},
		{"valid mixed name", "test_任务_123", true},
		{"too short", "a", false},
		{"too long", "this_is_a_very_long_task_name_that_exceeds_the_limit", false},
		{"empty", "", false},
		{"invalid chars", "test@task", false},
		{"valid with dash", "test-task", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateTaskName(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateTaskName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateCronSpec(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid 5 fields", "0 0 12 * *", true},
		{"valid 6 fields", "0 0 0 12 * *", true},
		{"invalid 4 fields", "0 0 12 *", false},
		{"invalid 7 fields", "0 0 0 0 12 * *", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCronSpec(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateCronSpec(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid http", "http://example.com", true},
		{"valid https", "https://example.com", true},
		{"valid with path", "https://example.com/api/test", true},
		{"invalid protocol", "ftp://example.com", false},
		{"invalid format", "not-a-url", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateURL(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateURL(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid email", "test@example.com", true},
		{"valid with subdomain", "test@mail.example.com", true},
		{"invalid format", "invalid-email", false},
		{"invalid missing @", "testexample.com", false},
		{"invalid missing domain", "test@", false},
		{"empty (allowed)", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateEmail(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid strong", "Password123", true},
		{"valid with lower and digit", "password123", true},
		{"valid with upper and digit", "PASSWORD123", true},
		{"too short", "Pass1", false},
		{"too long", "ThisIsAVeryLongPasswordThatExceedsTheMaximumLengthLimitAndShouldBeRejectedByTheValidationFunctionBecauseItIsWayTooLongForAnyReasonableUse", false},
		{"only letters", "password", false},
		{"only digits", "12345678", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePassword(tt.input)
			if result != tt.expected {
				t.Errorf("ValidatePassword(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid username", "testuser", true},
		{"valid with underscore", "test_user", true},
		{"valid with numbers", "testuser123", true},
		{"too short", "ab", false},
		{"too long", "this_is_a_very_long_username_that_exceeds_limit", false},
		{"starts with number", "123user", false},
		{"starts with underscore", "_testuser", false},
		{"invalid chars", "test-user", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUsername(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateUsername(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal string", "hello world", "hello world"},
		{"with spaces", "  hello world  ", "hello world"},
		{"with tabs", "\thello\tworld\t", "hello\tworld"},
		{"with newlines", "hello\nworld", "hello\nworld"},
		{"with control chars", "hello\x00world\x01", "helloworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateJSONString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid object", `{"key": "value"}`, true},
		{"valid array", `["item1", "item2"]`, true},
		{"invalid format", `{key: value}`, false},
		{"not json", "plain text", false},
		{"empty (allowed)", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateJSONString(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateJSONString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
