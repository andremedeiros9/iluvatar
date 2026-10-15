package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupFieldType(t *testing.T) {
	tC := []struct {
		name string
		want fieldTypeInfo
	}{
		{"string", fieldTypeInfo{goType: "string", protoType: "string", sqlType: map[string]string{"postgres": "TEXT", "mysql": "VARCHAR(255)", "sqlite": "TEXT"}}},
		{"int", fieldTypeInfo{goType: "int64", protoType: "int64", sqlType: map[string]string{"postgres": "BIGINT", "mysql": "BIGINT", "sqlite": "INTEGER"}}},
		{"float", fieldTypeInfo{goType: "float64", protoType: "double", sqlType: map[string]string{"postgres": "DOUBLE PRECISION", "mysql": "DOUBLE", "sqlite": "REAL"}}},
		{"bool", fieldTypeInfo{goType: "bool", protoType: "bool", sqlType: map[string]string{"postgres": "BOOLEAN", "mysql": "BOOLEAN", "sqlite": "BOOLEAN"}}},
		{"time", fieldTypeInfo{goType: "time.Time", protoType: "google.protobuf.Timestamp", sqlType: map[string]string{"postgres": "TIMESTAMPTZ", "mysql": "DATETIME", "sqlite": "TIMESTAMP"}}},
		{"uuid", fieldTypeInfo{goType: "uuid.UUID", protoType: "string", sqlType: map[string]string{"postgres": "UUID", "mysql": "CHAR(36)", "sqlite": "TEXT"}}},
	}

	for _, tc := range tC {
		t.Run(tc.name, func(t *testing.T) {
			got, err := lookupFieldType(tc.name)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	_, err := lookupFieldType("blob")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported field type")
	require.Contains(t, err.Error(), "string, int, float, bool, time, uuid")
}

func TestSQLTypeForDriver(t *testing.T) {
	info, err := lookupFieldType("string")
	require.NoError(t, err)

	tC := []struct {
		driver, want string
	}{
		{"postgres", "TEXT"},
		{"mysql", "VARCHAR(255)"},
		{"sqlite", "TEXT"},
	}

	for _, tc := range tC {
		t.Run(tc.driver, func(t *testing.T) {
			got, err := info.sqlTypeForDriver(tc.driver)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	_, err = info.sqlTypeForDriver("oracle")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported database driver")
}
