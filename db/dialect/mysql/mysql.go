package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/yaochi-tech/goqu"
	"github.com/yaochi-tech/lingquan-core-go/db/dialect"
)

var _ dialect.Dialect = (*dialect.DialectWrapper)(nil)

func init() {
	dialect.RegisterDialect("mysql", &dialect.DialectWrapper{
		Dialect: goqu.Dialect("mysql"),
	})
}
