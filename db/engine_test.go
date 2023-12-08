package db

import (
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const def = `{
  "code": "user",
  "name": "用户",
  "description": "用户模型",
  "comment": "用户表",
  "fields": [
    {
      "label": "主键",
      "name": "id",
      "type": "ID"
    },
    {
      "label": "用户名",
      "name": "username",
      "type": "String",
      "comment": "用户名",
      "showInList": true,
      "showable": true,
      "queryCondition": true,
      "orderable": true,
      "required": true,
      "unique": true,
      "index": true,
      "length": 20,
      "validations": [
        {
          "type": "required",
          "message": "用户名不能为空"
        },
        {
          "minlength": 6,
          "message": "用户名不能少于6个字符"
        },
        {
          "maxlength": 20,
          "message": "用户名不能多于20个字符"
        }
      ]
    },
    {
      "label": "密码",
      "name": "password",
      "type": "string",
      "comment": "密码",
      "showInList": false,
      "showable": false,
      "queryCondition": false,
      "orderable": false,
      "required": true,
      "length": 50,
      "crypt": "BCRYPT",
      "validations": [
        {
          "type": "required",
          "message": "密码不能为空"
        },
        {
          "minlength": 6,
          "message": "密码不能少于6个字符"
        },
        {
          "maxlength": 50,
          "message": "密码不能多于50个字符"
        }
      ]
    },
    {
      "label": "昵称",
      "name": "nickname",
      "type": "string",
      "comment": "昵称",
      "showInList": true,
      "showable": true,
      "queryCondition": true,
      "orderable": true,
      "required": true,
      "length": 20,
      "validations": [
        {
          "type": "required",
          "message": "昵称不能为空"
        },
        {
          "minlength": 2,
          "message": "昵称不能少于2个字符"
        },
        {
          "maxlength": 20,
          "message": "昵称不能多于20个字符"
        }
      ]
    },
    {
      "label": "邮箱",
      "name": "email",
      "type": "string",
      "comment": "邮箱",
      "showInList": true,
      "showable": true,
      "queryCondition": true,
      "orderable": true,
      "required": true,
      "unique": true,
      "index": true,
      "length": 50,
      "validations": [
        {
          "type": "required",
          "message": "邮箱不能为空"
        },
        {
          "type": "email",
          "message": "邮箱格式不正确"
        },
        {
          "minlength": 6,
          "message": "邮箱不能少于6个字符"
        },
        {
          "maxlength": 50,
          "message": "邮箱不能多于50个字符"
        }
      ]
    },
    {
      "label": "手机号",
      "name": "mobile",
      "type": "string",
      "comment": "手机号",
      "showInList": true,
      "showable": true,
      "queryCondition": true,
      "orderable": true,
      "required": true,
      "unique": true,
      "index": true,
      "length": 11,
      "validations": [
        {
          "type": "required",
          "message": "手机号不能为空"
        },
        {
          "type": "mobile",
          "message": "手机号格式不正确"
        },
        {
          "minlength": 11,
          "message": "手机号不能少于11个字符"
        },
        {
          "maxlength": 11,
          "message": "手机号不能多于11个字符"
        }
      ]
    },
    {
      "label": "头像",
      "name": "avatar",
      "type": "string",
      "comment": "头像",
      "showInList": true,
      "showable": true,
      "queryCondition": false,
      "orderable": false,
      "required": false,
      "length": 500
    },
    {
      "label": "性别",
      "name": "gender",
      "type": "enum",
      "enumType": "string",
      "comment": "性别",
      "default": "保密",
      "showInList": true,
      "showable": true,
      "queryCondition": true,
      "orderable": false,
      "required": false,
      "length": 2,
      "enum": [
        "男",
        "女",
        "保密"
      ]
    }
  ],
  "options": {
    "timestamps": true,
    "softDelete": true
  },
  "values": [
    {
      "username": "admin",
      "password": "123456",
      "nickname": "管理员",
      "email": "support@yaochi.tech",
      "gender": "保密"
    }
  ]
}`

func TestEngine(t *testing.T) {
	Convey("数据库引擎集成测试", t, func() {
		engine, err := NewEngine("mysql", "root:root@/lowcode?charset=utf8mb4&parseTime=True&loc=Local")
		So(err, ShouldBeNil)
		So(engine, ShouldNotBeNil)

		// 未注册模型
		schema := engine.GetSchema("notfound")
		So(schema, ShouldBeNil)

		// 注册模型
		name, err := engine.Register(def)
		So(err, ShouldBeNil)
		So(name, ShouldEqual, "user")

		// 已注册模型
		schema = engine.GetSchema("user")
		So(schema, ShouldNotBeNil)
		So(schema.Name, ShouldEqual, "user")

		// 获取所有模型
		schemas := engine.GetSchemas()
		So(schemas, ShouldNotBeNil)
		So(len(schemas), ShouldEqual, 1)

		// 检查模型对应的表格不存在
		var exists bool
		exists, err = engine.SchemaTableExists("user")
		So(err, ShouldBeNil)
		So(exists, ShouldBeFalse)

		// 创建表
		err = engine.MigrateTable("user")
		So(err, ShouldBeNil)

		// 检查模型对应的表格存在
		exists, err = engine.SchemaTableExists("user")
		So(err, ShouldBeNil)
		So(exists, ShouldBeTrue)

		// 删除表
		err = engine.DropTable("user")
		So(err, ShouldBeNil)

		// 检查模型对应的表格不存在
		exists, err = engine.SchemaTableExists("user")
		So(err, ShouldBeNil)
		So(exists, ShouldBeFalse)
	})

}
