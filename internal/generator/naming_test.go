package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToPascalCase(t *testing.T) {
	tC := []struct {
		in, want string
	}{
		{"user", "User"},
		{"created_at", "CreatedAt"},
		{"signup_at", "SignupAt"},
		{"referrer_id", "ReferrerId"},
		{"id", "Id"},
		{"URL", "URL"},
		{"", ""},
		{"_leading", "Leading"},
		{"trailing_", "Trailing"},
		{"double__underscore", "DoubleUnderscore"},
		{"alreadyPascal", "AlreadyPascal"},
		{"a_b_c", "ABC"},
	}

	for _, tc := range tC {
		t.Run(tc.in, func(t *testing.T) {
			require.Equal(t, tc.want, toPascalCase(tc.in))
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tC := []struct {
		in, want string
	}{
		{"user", "user"},
		{"created_at", "createdAt"},
		{"signup_at", "signupAt"},
		{"referrer_id", "referrerId"},
		{"id", "id"},
		{"URL", "uRL"},
		{"", ""},
		{"alreadyCamel", "alreadyCamel"},
		{"a", "a"},
	}

	for _, tc := range tC {
		t.Run(tc.in, func(t *testing.T) {
			require.Equal(t, tc.want, toCamelCase(tc.in))
		})
	}
}
