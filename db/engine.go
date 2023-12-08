package db

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/yaochi-tech/lingquan-core-go/db/clause"
	"github.com/yaochi-tech/lingquan-core-go/db/dialect"
	"github.com/yaochi-tech/lingquan-core-go/db/schema"
	"sync"
)

var (
	ErrSchemaNotRegistered error = errors.New("schema not registered")
)

// Engine 数据库引擎, 该引擎通过解析模型json文件, 生成对应的数据库表，并对表进行增删改查操作
type Engine struct {
	DB              *sqlx.DB
	dialect         dialect.Dialect
	currentDatabase string
	schemas         map[string]*schema.Schema // 模型名(code) => 模型
	lock            sync.RWMutex
}

func NewEngine(driverName, dataSourceName string) (*Engine, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	engine := new(Engine)
	engine.DB = db
	engine.schemas = make(map[string]*schema.Schema)
	var ok bool
	engine.dialect, ok = dialect.GetDialect(driverName)
	if !ok {
		err = dialect.ErrDialectNotSupported
		return nil, err
	}

	// 查询当前数据库
	sql := engine.dialect.CurrentDatabaseSQL()
	// 使用sqlx的Get方法，将查询结果赋值给engine.currentDatabase
	err = engine.DB.Get(&engine.currentDatabase, sql)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// Close 关闭数据库连接
func (engine *Engine) Close() {
	_ = engine.DB.Close()
}

// Register 注册模型
func (engine *Engine) Register(definition string) (string, error) {
	// 加锁
	engine.lock.Lock()
	defer engine.lock.Unlock()
	s := schema.Parse(definition)
	engine.schemas[s.Name] = s // 这里的Name是模型名(code)
	return s.Name, nil
}

// GetSchema 获取模型
func (engine *Engine) GetSchema(name string) *schema.Schema {
	engine.lock.RLock()
	defer engine.lock.RUnlock()
	return engine.schemas[name]
}

// GetSchemas 获取所有模型
func (engine *Engine) GetSchemas() map[string]*schema.Schema {
	engine.lock.RLock()
	defer engine.lock.RUnlock()
	return engine.schemas
}

// SchemaTableExists 检查模型对应的表格是否存在
func (engine *Engine) SchemaTableExists(name string, tx ...*sqlx.Tx) (bool, error) {
	s := engine.GetSchema(name)
	if s == nil {
		// 抛出异常，模型不存在
		return false, ErrSchemaNotRegistered
	}
	sql, args := engine.dialect.TableExistSQL(s.Name, engine.currentDatabase)

	var row *sqlx.Row
	if len(tx) > 0 {
		row = tx[0].QueryRowx(sql, args...)
	} else {
		row = engine.DB.QueryRowx(sql, args...)
	}
	if row.Err() != nil {
		return false, row.Err()
	}

	var tableName string
	err := row.Scan(&tableName)
	if err != nil {
		return false, nil
	}
	return tableName != "", nil
}

// MigrateTable 迁移表
func (engine *Engine) MigrateTable(name string, tx ...*sqlx.Tx) error {
	s := engine.GetSchema(name)
	if s == nil {
		return nil
	}
	// 先查看是否存在表
	tableExists, err := engine.SchemaTableExists(name, tx...)
	if err != nil {
		return err
	}

	if !tableExists {
		// 如果不存在表，则创建表
		sql := engine.dialect.CreateTableSQL(s)
		if len(tx) > 0 {
			_, err = tx[0].Exec(sql)
		} else {
			_, err = engine.DB.Exec(sql)
		}
		if err != nil {
			return err
		}
	} else {
		// 如果表已经存在，则检查字段、索引等是否有变化
		// TODO
	}

	return err
}

// DropTable 删除表
func (engine *Engine) DropTable(name string) error {
	s := engine.GetSchema(name)
	if s == nil {
		return nil
	}
	sql := engine.dialect.DropTableSQL(s)
	_, err := engine.DB.Exec(sql)
	return err
}

// Migrate 迁移所有注册的模型
func (engine *Engine) Migrate() error {
	if len(engine.schemas) == 0 {
		return nil
	}
	tx, err := engine.DB.Beginx()
	if err != nil {
		return err
	}
	for name := range engine.schemas {
		err = engine.MigrateTable(name, tx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Insert 插入数据
func (engine *Engine) Insert(name string, datas ...map[string]interface{}) (int64, error) {
	s := engine.GetSchema(name)
	if s == nil {
		return 0, nil
	}

	c := new(clause.Clause)
	c.Set(clause.INSERT, s.Name, s.Fields)
	for _, data := range datas {
		c.Set(clause.VALUES, data)
	}
	sql, args := c.Build()
	res, err := engine.DB.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Find 查询数据, where中的条件使用命名参数，如：where = "id = :id", namedCondition = map[string]interface{}{"id": 1}
func (engine *Engine) Find(name, where string, namedCondition map[string]interface{}, limit, offset int, orderBy string) ([]map[string]interface{}, error) {
	s := engine.GetSchema(name)
	if s == nil {
		return nil, nil
	}

	c := new(clause.Clause)
	c.Set(clause.SELECT, s.Name, s.Fields)
	c.Set(clause.WHERE, where, namedCondition)
	c.Set(clause.LIMIT, limit)
	c.Set(clause.OFFSET, offset)
	c.Set(clause.ORDERBY, orderBy)
	sql, args := c.Build()
	rows, err := engine.DB.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// Update 更新数据, where中的条件使用命名参数，如：where = "id = :id", namedCondition = map[string]interface{}{"id": 1}
func (engine *Engine) Update(name, where string, namedCondition, data map[string]interface{}) (int64, error) {
	s := engine.GetSchema(name)
	if s == nil {
		return 0, nil
	}

	c := new(clause.Clause)
	c.Set(clause.UPDATE, s.Name, data)
	c.Set(clause.WHERE, where, namedCondition)
	sql, args := c.Build()
	res, err := engine.DB.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete 删除数据, where中的条件使用命名参数，如：where = "id = :id", namedCondition = map[string]interface{}{"id": 1}，注意，delete方法必须有where条件
func (engine *Engine) Delete(name, where string, namedCondition map[string]interface{}) (int64, error) {
	s := engine.GetSchema(name)
	if s == nil {
		return 0, nil
	}

	// 判断where条件是否存在
	if where == "" {
		return 0, errors.New("delete method must have where condition")
	}

	c := new(clause.Clause)
	c.Set(clause.DELETE, s.Name)
	c.Set(clause.WHERE, where, namedCondition)
	sql, args := c.Build()
	res, err := engine.DB.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Count 统计数据, where中的条件使用命名参数，如：where = "id = :id", namedCondition = map[string]interface{}{"id": 1}
func (engine *Engine) Count(name, where string, namedCondition map[string]interface{}) (int64, error) {
	s := engine.GetSchema(name)
	if s == nil {
		return 0, nil
	}

	c := new(clause.Clause)
	c.Set(clause.COUNT, s.Name)
	c.Set(clause.WHERE, where, namedCondition)
	sql, args := c.Build()
	var count int64
	err := engine.DB.Get(&count, sql, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// First 查询第一条数据, where中的条件使用命名参数，如：where = "id = :id", namedCondition = map[string]interface{}{"id": 1}
func (engine *Engine) First(name, where string, namedCondition map[string]interface{}) (map[string]interface{}, error) {
	s := engine.GetSchema(name)
	if s == nil {
		return nil, nil
	}

	// 可以直接调用Find方法
	results, err := engine.Find(name, where, namedCondition, 1, 0, "")

	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}
