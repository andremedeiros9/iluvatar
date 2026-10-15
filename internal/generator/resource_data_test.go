package generator

import (
	"testing"

	"github.com/andremedeiros9/iluvatar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBuildResourceData(t *testing.T) {
	res := config.Resource{
		Name:          "user",
		Communication: "grpc",
		Fields: []config.Field{
			{Name: "email", Type: "string"},
			{Name: "signup_at", Type: "time"},
			{Name: "referrer_id", Type: "uuid"},
		},
	}

	rd, err := buildResourceData("github.com/example/widgets", "postgres", res)
	require.NoError(t, err)

	require.Equal(t, "user", rd.Name)
	require.Equal(t, "User", rd.PascalName)
	require.Equal(t, "user", rd.LowerName)
	require.Equal(t, "users", rd.PluralField)
	require.Equal(t, "Users", rd.PluralPb)
	require.Equal(t, "user", rd.PackageName)
	require.Equal(t, "user", rd.TableName)
	require.Equal(t, "postgres", rd.Driver)
	require.Equal(t, "UUID", rd.IDSQLType)
	require.Equal(t, "TIMESTAMPTZ", rd.TimestampSQLType)

	require.Len(t, rd.Fields, 3)

	email := rd.Fields[0]
	require.Equal(t, "email", email.Name)
	require.Equal(t, "Email", email.PascalName)
	require.Equal(t, "email", email.CamelName)
	require.Equal(t, "string", email.GoType)
	require.Equal(t, "string", email.ProtoType)
	require.Equal(t, "TEXT", email.SQLType)
	require.Equal(t, 2, email.NumberWithID)
	require.Equal(t, 1, email.NumberNoID)

	signup := rd.Fields[1]
	require.Equal(t, "time.Time", signup.GoType)
	require.Equal(t, "google.protobuf.Timestamp", signup.ProtoType)
	require.Equal(t, "TIMESTAMPTZ", signup.SQLType)
	require.Equal(t, 3, signup.NumberWithID)
	require.Equal(t, 2, signup.NumberNoID)

	referrer := rd.Fields[2]
	require.Equal(t, "uuid.UUID", referrer.GoType)
	require.Equal(t, "string", referrer.ProtoType)
	require.Equal(t, "UUID", referrer.SQLType)
	require.Equal(t, 4, referrer.NumberWithID)
	require.Equal(t, 3, referrer.NumberNoID)

	require.Equal(t, "id, email, signup_at, referrer_id, created_at, updated_at", rd.InsertColumns)
	require.Equal(t, "?, ?, ?, ?, ?, ?", rd.InsertPlaceholders)
	require.Equal(t, "id, email, signup_at, referrer_id, created_at, updated_at", rd.SelectColumns)
	require.Equal(t, "email = ?, signup_at = ?, referrer_id = ?, updated_at = ?", rd.UpdateSetClause)

	require.Equal(t, 5, rd.CreatedAtNumber())
	require.Equal(t, 6, rd.UpdatedAtNumber())
}

func TestBuildResourceDataUnsupportedFieldType(t *testing.T) {
	res := config.Resource{
		Name:          "user",
		Communication: "grpc",
		Fields:        []config.Field{{Name: "avatar", Type: "blob"}},
	}

	_, err := buildResourceData("github.com/example/widgets", "postgres", res)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported field type")
}

func TestBuildResourceDataUnsupportedDriver(t *testing.T) {
	res := config.Resource{
		Name:          "user",
		Communication: "grpc",
		Fields:        []config.Field{{Name: "email", Type: "string"}},
	}

	_, err := buildResourceData("github.com/example/widgets", "oracle", res)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported database driver")
}

func TestBuildResourceDataEmptyFields(t *testing.T) {
	res := config.Resource{
		Name:          "widget",
		Communication: "grpc",
	}

	rd, err := buildResourceData("github.com/example/widgets", "sqlite", res)
	require.NoError(t, err)
	require.Empty(t, rd.Fields)
	require.Equal(t, "id, created_at, updated_at", rd.InsertColumns)
	require.Equal(t, "?, ?, ?", rd.InsertPlaceholders)
	require.Equal(t, 2, rd.CreatedAtNumber())
	require.Equal(t, 3, rd.UpdatedAtNumber())
}
