package dialect

import (
	"github.com/yaochi-tech/lingquan-core-go/db/schema"
	"strconv"
	"strings"
)

type mysql struct{}

func (m mysql) CreateTableSQL(schema *schema.Schema) string {
	var columns []string
	var primaryKeys []string
	for _, field := range schema.Fields {
		columns = append(columns, m.columnSQL(field))
		if field.IsPrimaryKey {
			primaryKeys = append(primaryKeys, field.Name)
		}
	}
	var sql strings.Builder
	sql.WriteString("CREATE TABLE IF NOT EXISTS ")
	sql.WriteString("`" + schema.TableName + "`")
	sql.WriteString(" (")
	sql.WriteString(strings.Join(columns, ","))
	if len(primaryKeys) > 0 {
		sql.WriteString(", PRIMARY KEY(")
		sql.WriteString(strings.Join(primaryKeys, ","))
		sql.WriteString(")")
	}
	sql.WriteString(")")
	return sql.String()
}

func (m mysql) DropTableSQL(schema *schema.Schema) string {
	return "DROP TABLE IF EXISTS `" + schema.TableName + "`"
}

func (m mysql) DataTypeOf(typ string) string {
	t := strings.ToLower(typ)
	switch t {
	case "string":
		return "varchar"
	case "text":
		return "text"
	case "id", "int64", "uint64":
		return "bigint"
	case "int", "uint":
		return "int"
	case "bool":
		return "bool"
	case "float32", "float64":
		return "float"
	case "date", "datetime":
		return "datetime"
	case "json":
		return "json"
	}
	panic("invalid sql type " + typ)
}

func (m mysql) CurrentDatabaseSQL() string {
	return "SELECT DATABASE()"
}

func (m mysql) TableExistSQL(tableName, dbName string) (string, []interface{}) {
	args := []interface{}{tableName, dbName}
	return "SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ? ", args
}

func (m mysql) columnSQL(field *schema.Field) string {
	var sql strings.Builder
	sql.WriteString("`" + field.Name + "`")

	if field.Type == "string" {
		length := "255"
		if field.Length > 0 {
			length = strconv.Itoa(int(field.Length))
		}
		sql.WriteString(" varchar(" + length + ")")
	} else if field.Type == "id" {
		sql.WriteString(" bigint")
	} else if field.Type == "int" {
		sql.WriteString(" int")
	} else if field.Type == "bool" {
		sql.WriteString(" bool")
	} else if field.Type == "float32" {
		sql.WriteString(" float")
	} else if field.Type == "float64" {
		sql.WriteString(" double")
	} else if field.Type == "date" {
		sql.WriteString(" datetime")
	} else if field.Type == "datetime" {
		sql.WriteString(" datetime")
	} else if field.Type == "json" {
		sql.WriteString(" json")
	} else if field.Type == "text" {
		sql.WriteString(" text")
	}

	if field.IsPrimaryKey || field.NotNull {
		sql.WriteString(" NOT NULL")
	}
	if field.Default != "" {
		if field.IsDefaultRaw || (field.Type != "string" && field.Type != "json") {
			sql.WriteString(" DEFAULT " + field.Default)
		} else {
			sql.WriteString(" DEFAULT '" + field.Default + "'")
		}
	}
	return sql.String()
}

var _ Dialect = (*mysql)(nil)

func init() {
	RegisterDialect("mysql", &mysql{})
}
