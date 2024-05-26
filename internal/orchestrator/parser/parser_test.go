package parser

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	testCases := []struct {
		expr     string
		expected float64
		errMsg   string
	}{
		{"2 + 3", 5.0, ""},
		{"2 - 3", -1.0, ""},
		{"2 * 3", 6.0, ""},
		{"6 / 3", 2.0, ""},
		{"2 + 3 * 4", 14.0, ""},
		{"(2 + 3) * 4", 20.0, ""},
		{"2 + 3 * (4 + 5)", 29.0, ""},
		{"2.5 + 3.7", 6.2, ""},
		{"2 * (3 + 4) / 2", 7.0, ""},
		//{"2 / 0", 0.0, "division by zero"},
		{"2 + ", 0.0, "unexpected end of expression"},
		{"2 + 3 * ", 0.0, "unexpected end of expression"},
		{"(2 + 3", 0.0, "missing closing parenthesis"},
		{"2 + 3 ) ", 0.0, "unexpected token in expression"},
		{"2 % 3", 0.0, "unexpected character: %"},
	}
	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			root, err := Parse(tc.expr)
			if err != nil {
				if tc.errMsg == "" {
					t.Errorf("unexpected error: %v", err)
				} else if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("expected error '%s', got '%v'", tc.errMsg, err)
				}
				return
			}

			if tc.errMsg != "" {
				t.Errorf("expected error '%s', but got no error", tc.errMsg)
				return
			}

			result, err := root.Evaluate()
			if err != nil {
				t.Errorf("unexpected error during evaluation: %v", err)
				return
			}

			if result != tc.expected {
				t.Errorf("expected %f, got %f", tc.expected, result)
			}
		})
	}
}
