package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

// 生成器
var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[LIMIT] = _limit
	generators[OFFSET] = _offset

	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func _count(values ...interface{}) (string, []interface{}) {
	// SELECT COUNT(*) FROM $tableName
	return fmt.Sprintf("SELECT COUNT(*) FROM %s", values[0]), []interface{}{}
}

func _delete(values ...interface{}) (string, []interface{}) {
	// DELETE FROM $tableName
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	// UPDATE $tableName SET $fields
	tableName := values[0]
	fields := values[1].(map[string]interface{})
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString(fmt.Sprintf("UPDATE %s SET ", tableName))
	for k, v := range fields {
		sql.WriteString(fmt.Sprintf("%s = ?, ", k))
		vars = append(vars, v)
	}
	sqlStr := sql.String()
	return sqlStr[:len(sqlStr)-2], vars
}

func _offset(values ...interface{}) (string, []interface{}) {
	// OFFSET $num
	return "OFFSET ?", values
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	// ORDER BY $order
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	// 如果没有传值，直接返回空字符串
	if len(values) == 0 {
		return "", []interface{}{}
	}
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ($v3)
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func genBindVars(num int) string {
	vars := make([]string, num)
	for i := range vars {
		vars[i] = "?"
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}
