package dialect

import (
	"errors"
	"github.com/yaochi-tech/goqu"
	"github.com/yaochi-tech/lingquan-core-go/db/schema"
	"github.com/yaochi-tech/lingquan-core-go/util"
	"reflect"
	"strconv"
	"strings"
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

	BuildInsert(tableName string, dataList []map[string]interface{}) (string, []interface{}, error)
	BuildSelect(tableName string, selectFields []string, namedCondition map[string]interface{}) (string, []interface{}, error)
	BuildUpdate(tableName string, updateData, where map[string]interface{}) (string, []interface{}, error)
	BuildDelete(tableName string, where map[string]interface{}) (string, []interface{}, error)
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

type dialectImpl struct {
	dialect goqu.DialectWrapper
}

func (m *dialectImpl) CreateTableSQL(schema *schema.Schema) string {
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

func (m *dialectImpl) DropTableSQL(schema *schema.Schema) string {
	return "DROP TABLE IF EXISTS `" + schema.TableName + "`"
}

func (m *dialectImpl) DataTypeOf(typ string) string {
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

func (m *dialectImpl) CurrentDatabaseSQL() string {
	return "SELECT DATABASE()"
}

func (m *dialectImpl) TableExistSQL(tableName, dbName string) (string, []interface{}) {
	args := []interface{}{tableName, dbName}
	return "SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ? ", args
}

func (m *dialectImpl) columnSQL(field *schema.Field) string {
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

func (m *dialectImpl) BuildInsert(tableName string, dataList []map[string]interface{}) (string, []interface{}, error) {
	// dataList中的key转蛇形命名
	var snakeDataList []interface{}
	for _, data := range dataList {
		snakeData := make(map[string]interface{})
		for k, v := range data {
			snakeData[util.ToSnake(k)] = v
		}
		snakeDataList = append(snakeDataList, snakeData)
	}
	return m.dialect.Insert(tableName).Rows(snakeDataList...).ToSQL()
}

func (m *dialectImpl) BuildSelect(tableName string, selectFields []string, where map[string]interface{}) (string, []interface{}, error) {
	// selectFields转蛇形命名
	var snakeSelectFields []interface{}
	for _, field := range selectFields {
		snakeSelectFields = append(snakeSelectFields, util.ToSnake(field))
	}

	ds := m.dialect.From(tableName).Select(snakeSelectFields...)
	// where要处理成goqu的where语句
	var whereExList []goqu.Expression
	whereExList = whereExpression(where)

	ds = ds.Where(whereExList...)

	// 对where中的特殊操作进行处理
	for k, v := range where {
		switch k {
		case OP_LIMIT:
			// v应该是一个正整数
			if reflect.TypeOf(v).Kind() == reflect.Uint {
				i := v.(uint)
				ds = ds.Limit(i)
			}
		case OP_OFFSET:
			// v应该是一个正整数
			if reflect.TypeOf(v).Kind() == reflect.Uint {
				i := v.(uint)
				ds = ds.Offset(i)
			}
		case OP_ORDER_BY:
			// v应该是一个string或数组[]string
			if reflect.TypeOf(v).Kind() == reflect.String {
				s := v.(string)
				columnWithOrder := strings.Split(s, " ")
				column := util.ToSnake(columnWithOrder[0])
				if len(columnWithOrder) == 2 {
					if strings.ToLower(columnWithOrder[1]) == "desc" {
						ds = ds.Order(goqu.C(column).Desc())
					} else {
						ds = ds.Order(goqu.C(column).Asc())
					}
				} else {
					ds = ds.Order(goqu.C(column).Asc())
				}
			} else if reflect.TypeOf(v).Kind() == reflect.Slice {
				s := v.([]string)
				for _, co := range s {
					columnWithOrder := strings.Split(co, " ")
					column := util.ToSnake(columnWithOrder[0])
					if len(columnWithOrder) == 2 {
						if strings.ToLower(columnWithOrder[1]) == "desc" {
							ds = ds.Order(goqu.C(column).Desc())
						} else {
							ds = ds.Order(goqu.C(column).Asc())
						}
					} else {
						ds = ds.Order(goqu.C(column).Asc())
					}
				}
			}
		case OP_GROUP_BY:
			// v应该是一个string或数组[]string
			if reflect.TypeOf(v).Kind() == reflect.String {
				s := v.(string)
				column := util.ToSnake(s)
				ds = ds.GroupBy(column)
			} else if reflect.TypeOf(v).Kind() == reflect.Slice {
				s := v.([]string)
				for _, co := range s {
					column := util.ToSnake(co)
					ds = ds.GroupBy(column)
				}
			}
		case OP_HAVING:
			// v应该是一个map[string]interface{}，否则忽略
			if reflect.TypeOf(v).Kind() == reflect.Map {
				// v转为map[string]interface{}
				s := v.(map[string]interface{})
				// 递归调用whereExpression
				ds = ds.Having(whereExpression(s)...)
			}
		}
	}

	return ds.ToSQL()
}

func whereExpression(m map[string]interface{}) []goqu.Expression {
	var whereExList []goqu.Expression
	for k, v := range m {
		k = util.ToSnake(k)
		if !strings.Contains(k, " ") {
			// 如果v不是数组
			if reflect.TypeOf(v).Kind() == reflect.Slice {
				// v转为数组[]interface{}
				s := v.([]interface{})
				whereExList = append(whereExList, goqu.C(k).In(s...))
			} else {
				whereExList = append(whereExList, goqu.C(k).Eq(v))
			}
		} else {
			splited := strings.Split(k, " ")
			if len(splited) != 2 {
				continue // 忽略
			}
			op := strings.ToLower(splited[1])
			switch op {
			case OP_IN:
				// 如果v不是数组
				if reflect.TypeOf(v).Kind() == reflect.Slice {
					// v转为数组[]interface{}
					s := v.([]interface{})
					whereExList = append(whereExList, goqu.C(splited[0]).In(s...))
				} else {
					whereExList = append(whereExList, goqu.C(splited[0]).Eq(v))
				}
			case OP_LIKE:
				whereExList = append(whereExList, goqu.C(splited[0]).Like(v))
			case OP_NOT_LIKE:
				whereExList = append(whereExList, goqu.C(splited[0]).NotLike(v))
			case OP_BETWEEN:
				// 如果v不是数组
				if reflect.TypeOf(v).Kind() == reflect.Slice {
					// v转为数组[]interface{}
					s := v.([]interface{})
					if len(s) == 2 {
						whereExList = append(whereExList, goqu.C(splited[0]).Between(goqu.Range(s[0], s[1])))
					}
				}
			case OP_NOT_BETWEEN:
				// 如果v不是数组
				if reflect.TypeOf(v).Kind() == reflect.Slice {
					// v转为数组[]interface{}
					s := v.([]interface{})
					if len(s) == 2 {
						whereExList = append(whereExList, goqu.C(splited[0]).NotBetween(goqu.Range(s[0], s[1])))
					}
				}
			case OP_IS_NULL:
				whereExList = append(whereExList, goqu.C(splited[0]).IsNull())
			case OP_IS_NOT_NULL:
				whereExList = append(whereExList, goqu.C(splited[0]).IsNotNull())
			case OP_EQ:
				whereExList = append(whereExList, goqu.C(splited[0]).Eq(v))
			case OP_NEQ:
				whereExList = append(whereExList, goqu.C(splited[0]).Neq(v))
			case OP_GT:
				whereExList = append(whereExList, goqu.C(splited[0]).Gt(v))
			case OP_GTE:
				whereExList = append(whereExList, goqu.C(splited[0]).Gte(v))
			case OP_LT:
				whereExList = append(whereExList, goqu.C(splited[0]).Lt(v))
			case OP_LTE:
				whereExList = append(whereExList, goqu.C(splited[0]).Lte(v))
			case OP_IS:
				// 判断值是bool
				if reflect.TypeOf(v).Kind() == reflect.Bool {
					whereExList = append(whereExList, goqu.C(splited[0]).Is(v))
				} else {
					whereExList = append(whereExList, goqu.C(splited[0]).Eq(v))
				}
			case OP_OR:
				// v应该是一个map[string]interface{}，否则忽略
				if reflect.TypeOf(v).Kind() == reflect.Map {
					// v转为map[string]interface{}
					s := v.(map[string]interface{})
					// 递归调用whereExpression
					whereExList = append(whereExList, goqu.Or(whereExpression(s)...))
				}
			case OP_AND:
				// v应该是一个map[string]interface{}，否则忽略
				if reflect.TypeOf(v).Kind() == reflect.Map {
					// v转为map[string]interface{}
					s := v.(map[string]interface{})
					// 递归调用whereExpression
					whereExList = append(whereExList, goqu.And(whereExpression(s)...))
				}
			default:
				whereExList = append(whereExList, goqu.C(splited[0]).Eq(v))
			}
		}
	}
	return whereExList
}

func (m *dialectImpl) BuildUpdate(tableName string, updateData, where map[string]interface{}) (string, []interface{}, error) {
	// where要处理成goqu的where语句
	whereExList := whereExpression(where)

	ds := m.dialect.Update(tableName)
	// updateData中的key转蛇形命名
	var snakeUpdateData map[string]interface{}
	for k, v := range updateData {
		snakeUpdateData[util.ToSnake(k)] = v
	}

	ds = ds.Set(snakeUpdateData)

	ds = ds.Where(whereExList...)

	return ds.ToSQL()
}

func (m *dialectImpl) BuildDelete(tableName string, where map[string]interface{}) (string, []interface{}, error) {
	// where要处理成goqu的where语句
	whereExList := whereExpression(where)

	ds := m.dialect.Delete(tableName)

	ds = ds.Where(whereExList...)

	return ds.ToSQL()
}
