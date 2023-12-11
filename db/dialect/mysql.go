package dialect

import (
	"github.com/yaochi-tech/goqu"
	_ "github.com/yaochi-tech/goqu/dialect/mysql"
)

var _ Dialect = (*dialectImpl)(nil)

func init() {
	RegisterDialect("mysql", &dialectImpl{
		dialect: goqu.Dialect("mysql"),
	})
}
