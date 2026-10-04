package generator

import "strings"

// toPascalCase converts a snake_case identifier (as used for resource and
// field names in iluvatar.toml) to PascalCase, e.g. "created_at" ->
// "CreatedAt".
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		b.WriteString(p[1:])
	}
	return b.String()
}

// toCamelCase converts a snake_case identifier to camelCase, e.g.
// "created_at" -> "createdAt". Used for local variable names in
// generated Go code.
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if pascal == "" {
		return pascal
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}
