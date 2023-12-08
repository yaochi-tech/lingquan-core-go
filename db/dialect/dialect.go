package dialect

import (
	"errors"
	"github.com/yaochi-tech/lingquan-core-go/db/schema"
)

var (
	ErrDialectNotSupported error = errors.New("dialect not supported")
)

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ string) string
	CurrentDatabaseSQL() string
	// TableExistSQL 返回查询表是否存在的sql语句，sql查询应只有一个字段，表名
	TableExistSQL(tableName, dbName string) (string, []interface{})
	CreateTableSQL(schema *schema.Schema) string
	DropTableSQL(schema *schema.Schema) string
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
