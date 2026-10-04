package generator

import (
	"fmt"
	"strings"

	"github.com/andremedeiros9/iluvatar/internal/config"
)

// fieldData is the per-field information resource templates render from.
type fieldData struct {
	Name       string // original snake_case name, used as the column/proto field name
	PascalName string
	CamelName  string
	GoType     string
	ProtoType  string
	SQLType    string
	// NumberWithID is this field's proto field number in messages that
	// also carry an id field (the entity message, and update requests).
	NumberWithID int
	// NumberNoID is this field's proto field number in messages that
	// don't carry an id field (create requests).
	NumberNoID int
}

// resourceData is the full set of values a resource's templates render
// from: the model/repository/service Go files, the .proto file, and the
// resource's migration.
type resourceData struct {
	Module             string
	Name               string // original snake_case name
	PascalName         string
	LowerName          string
	PluralField        string // the List response's repeated field name, e.g. "users"
	PluralPb           string // that field's Go name in generated code, e.g. "Users"
	PackageName        string
	TableName          string
	Driver             string
	IDSQLType          string
	TimestampSQLType   string
	Fields             []fieldData
	InsertColumns      string
	InsertPlaceholders string
	SelectColumns      string
	UpdateSetClause    string
}

// buildResourceData resolves res into everything its templates need to
// render, validating that every field type and the database driver are
// ones iluvatar supports.
func buildResourceData(module, driver string, res config.Resource) (resourceData, error) {
	idInfo, err := lookupFieldType("uuid")
	if err != nil {
		return resourceData{}, err
	}
	idSQLType, err := idInfo.sqlTypeForDriver(driver)
	if err != nil {
		return resourceData{}, err
	}

	timeInfo, err := lookupFieldType("time")
	if err != nil {
		return resourceData{}, err
	}
	timestampSQLType, err := timeInfo.sqlTypeForDriver(driver)
	if err != nil {
		return resourceData{}, err
	}

	name := strings.ToLower(res.Name)
	rd := resourceData{
		Module:           module,
		Name:             name,
		PascalName:       toPascalCase(name),
		LowerName:        name,
		PluralField:      name + "s",
		PluralPb:         toPascalCase(name) + "s",
		PackageName:      name,
		TableName:        name,
		Driver:           driver,
		IDSQLType:        idSQLType,
		TimestampSQLType: timestampSQLType,
	}

	columns := []string{"id"}
	placeholders := []string{"?"}
	var setClauses []string

	for i, f := range res.Fields {
		info, err := lookupFieldType(f.Type)
		if err != nil {
			return resourceData{}, fmt.Errorf("resource %q field %q: %w", res.Name, f.Name, err)
		}

		sqlType, err := info.sqlTypeForDriver(driver)
		if err != nil {
			return resourceData{}, fmt.Errorf("resource %q field %q: %w", res.Name, f.Name, err)
		}

		fieldName := strings.ToLower(f.Name)
		rd.Fields = append(rd.Fields, fieldData{
			Name:         fieldName,
			PascalName:   toPascalCase(fieldName),
			CamelName:    toCamelCase(fieldName),
			GoType:       info.goType,
			ProtoType:    info.protoType,
			SQLType:      sqlType,
			NumberWithID: i + 2,
			NumberNoID:   i + 1,
		})

		columns = append(columns, fieldName)
		placeholders = append(placeholders, "?")
		setClauses = append(setClauses, fieldName+" = ?")
	}

	columns = append(columns, "created_at", "updated_at")
	placeholders = append(placeholders, "?", "?")
	setClauses = append(setClauses, "updated_at = ?")

	rd.InsertColumns = strings.Join(columns, ", ")
	rd.InsertPlaceholders = strings.Join(placeholders, ", ")
	rd.SelectColumns = strings.Join(columns, ", ")
	rd.UpdateSetClause = strings.Join(setClauses, ", ")

	return rd, nil
}

// CreatedAtNumber and UpdatedAtNumber are the entity message's proto field
// numbers for its two fixed timestamp fields, which come after id and
// every configured field.
func (rd resourceData) CreatedAtNumber() int { return len(rd.Fields) + 2 }
func (rd resourceData) UpdatedAtNumber() int { return len(rd.Fields) + 3 }
