package generator

import "fmt"

// fieldTypeInfo maps an iluvatar.toml field type to the Go, proto, and
// per-driver SQL types it's generated as.
type fieldTypeInfo struct {
	goType    string
	protoType string
	sqlType   map[string]string
}

var fieldTypes = map[string]fieldTypeInfo{
	"string": {
		goType:    "string",
		protoType: "string",
		sqlType:   map[string]string{"postgres": "TEXT", "mysql": "VARCHAR(255)", "sqlite": "TEXT"},
	},
	"int": {
		goType:    "int64",
		protoType: "int64",
		sqlType:   map[string]string{"postgres": "BIGINT", "mysql": "BIGINT", "sqlite": "INTEGER"},
	},
	"float": {
		goType:    "float64",
		protoType: "double",
		sqlType:   map[string]string{"postgres": "DOUBLE PRECISION", "mysql": "DOUBLE", "sqlite": "REAL"},
	},
	"bool": {
		goType:    "bool",
		protoType: "bool",
		sqlType:   map[string]string{"postgres": "BOOLEAN", "mysql": "BOOLEAN", "sqlite": "BOOLEAN"},
	},
	"time": {
		goType:    "time.Time",
		protoType: "google.protobuf.Timestamp",
		sqlType:   map[string]string{"postgres": "TIMESTAMPTZ", "mysql": "DATETIME", "sqlite": "TIMESTAMP"},
	},
	"uuid": {
		goType:    "uuid.UUID",
		protoType: "string",
		sqlType:   map[string]string{"postgres": "UUID", "mysql": "CHAR(36)", "sqlite": "TEXT"},
	},
}

// lookupFieldType returns the type info for name, or an error listing the
// supported field types.
func lookupFieldType(name string) (fieldTypeInfo, error) {
	info, ok := fieldTypes[name]
	if !ok {
		return fieldTypeInfo{}, fmt.Errorf("unsupported field type %q (want one of: string, int, float, bool, time, uuid)", name)
	}
	return info, nil
}

// sqlTypeForDriver returns info's SQL column type for driver, or an error
// if driver isn't one iluvatar supports.
func (info fieldTypeInfo) sqlTypeForDriver(driver string) (string, error) {
	sqlType, ok := info.sqlType[driver]
	if !ok {
		return "", fmt.Errorf("unsupported database driver %q (want %q, %q, or %q)", driver, "postgres", "mysql", "sqlite")
	}
	return sqlType, nil
}
