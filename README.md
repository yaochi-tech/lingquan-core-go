# lingquan-core-go
瑶池-灵泉-模型核心库go语言版本

# 灵泉模型核心库
本库提供基础的模型定义转换能力，通过模型定义json数据，转换为对应的结构体，并提供结构体自动建表、增删改查及自定义查询等能力。

## 功能列表
- [x] 模型json格式规范校验
- [x] 模型json自动建表
- [x] 模型json增删改查
- [ ] 模型json自定义查询
- [ ] 模型json对应结构体代码生成

## 模块
- [x] 自定义一套类orm方式的数据库操作

# 使用
## 引用库
```go
import "github.com/yaochi-tech/lingquan-core-go"
```

## 引用对应的数据库方言
```go
go get -u github.com/yaochi-tech/lingquan-core-go/dialect/mysql
```

## 定义模型
模型json参考example目录下的模型定义文件。 
[user.json](example/user.json)

## 引擎
```go
engine, err := lingquan.StartEngine("mysql", "root:root@127.0.0.1/lingquan?charset=utf8mb4&parseTime=True&loc=Local")

defer lingquan.StopEngine(engine)
```

## 注册模型
```go
engine.RegisterModel(userJson)
```

## 迁移表
```go
engine.MigrateTable("user") // user为模型名称，即模型json中的code字段
```

## 增删改查
```go
// 获取模型的模式
schema := engine.GetSchema("user") // user为模型名称，即模型json中的code字段

// 获取所有注册的模型模式
schemas := engine.GetSchemas()

// 增加数据, 注意，ID应外部传入，不应该由数据库自动生成
count, err := engine.Insert("user", map[string]interface{}{
	"id": 1,
    "name": "张三",
    "age":  18,
})

// 删除数据
count, err := engine.Delete("user", map[string]interface{}{
    "id": 1,
})

// 更新数据
count, err := engine.Update("user", map[string]interface{}{
    "name": "张三",
    "age":  18,
}, map[string]interface{}{
    "id": 1,
})

// 查询数据
rows, err := engine.Select("user", map[string]interface{}{
    "id": 1,
}, []string{"id", "name", "age"})
```

查询条件中的特殊参数参看[where.md](db/dialect/where.md)