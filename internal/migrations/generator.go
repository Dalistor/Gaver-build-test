package migrations

import (
	"fmt"
	"strings"
)

func GenerateSQLForSQLite(models []ModelInfo) string {
	var sql strings.Builder

	for _, model := range models {
		sql.WriteString(fmt.Sprintf("-- Migration for table: %s\n", model.TableName))
		sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", model.TableName))

		var fields []string
		var primaryKeys []string
		var indexes []string

		for _, field := range model.Fields {
			fieldSQL := generateSQLiteField(field)
			fields = append(fields, "    "+fieldSQL)

			if field.IsPrimaryKey {
				primaryKeys = append(primaryKeys, field.Name)
			}

			if field.IsIndex && !field.IsPrimaryKey {
				indexes = append(indexes, fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_%s ON %s (%s);", model.TableName, strings.ToLower(field.Name), model.TableName, field.Name))
			}

			if field.IsUnique && !field.IsPrimaryKey {
				indexes = append(indexes, fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS idx_%s_%s_unique ON %s (%s);", model.TableName, strings.ToLower(field.Name), model.TableName, field.Name))
			}
		}

		if len(primaryKeys) > 0 {
			fields = append(fields, fmt.Sprintf("    PRIMARY KEY (%s)", strings.Join(primaryKeys, ", ")))
		}

		sql.WriteString(strings.Join(fields, ",\n"))
		sql.WriteString("\n);\n\n")

		for _, indexSQL := range indexes {
			sql.WriteString(indexSQL + "\n")
		}

		sql.WriteString("\n")
	}

	return sql.String()
}

func generateSQLiteField(field FieldInfo) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("\"%s\"", strings.ToLower(field.Name)))

	sqlType := getSQLiteType(field)
	parts = append(parts, sqlType)

	if field.IsPrimaryKey {
		parts = append(parts, "PRIMARY KEY")
	}

	if field.IsNotNull && !field.IsPrimaryKey {
		parts = append(parts, "NOT NULL")
	}

	if field.DefaultValue != "" && !field.IsPrimaryKey {
		parts = append(parts, fmt.Sprintf("DEFAULT %s", field.DefaultValue))
	}

	return strings.Join(parts, " ")
}

func getSQLiteType(field FieldInfo) string {
	if field.SQLType != "" {
		return field.SQLType
	}

	switch field.Type {
	case "string":
		return "TEXT"
	case "int", "int32", "int64":
		return "INTEGER"
	case "uint", "uint32", "uint64":
		return "INTEGER"
	case "float32", "float64":
		return "REAL"
	case "bool", "boolean":
		return "INTEGER"
	case "time.Time":
		return "TEXT"
	case "uuid.UUID":
		return "TEXT"
	case "gorm.DeletedAt":
		return "TEXT"
	default:
		if strings.Contains(field.Type, "[]") {
			return "TEXT"
		}
		return "TEXT"
	}
}

