package schema

import (
	"github.com/yaochi-tech/lingquan-core-go/db/dialect"
	"testing"
)

func TestParse(t *testing.T) {

	type args struct {
		dest string
		d    dialect.Dialect
	}

	mysqlDialect, _ := dialect.GetDialect("mysql")

	fields := []*Field{
		{
			Label:        "主键",
			Name:         "id",
			Column:       "id",
			Type:         "ID",
			Comment:      "",
			Default:      "",
			IsDefaultRaw: false,
			IsPrimaryKey: true,
			NotNull:      false,
			Index:        "",
			Unique:       "",
			Length:       0,
			Precision:    0,
			Scale:        0,
		},
		{
			Label:        "用户名",
			Name:         "username",
			Column:       "username",
			Type:         "String",
			Comment:      "用户名",
			Default:      "",
			IsDefaultRaw: false,
			IsPrimaryKey: false,
			NotNull:      true,
			Index:        "IDX_USER_USERNAME",
			Unique:       "UNI_USER_USERNAME",
			Length:       20,
			Precision:    0,
			Scale:        0,
		},
		{
			Label:        "密码",
			Name:         "password",
			Column:       "password",
			Type:         "String",
			Comment:      "密码",
			Default:      "",
			IsDefaultRaw: false,
			IsPrimaryKey: false,
			NotNull:      true,
			Index:        "",
			Unique:       "",
			Length:       50,
			Precision:    0,
			Scale:        0,
		},
		{
			Label:        "昵称",
			Name:         "nickname",
			Column:       "nickname",
			Type:         "String",
			Comment:      "昵称",
			Default:      "",
			IsDefaultRaw: false,
			IsPrimaryKey: false,
			NotNull:      true,
			Index:        "",
			Unique:       "",
			Length:       20,
			Precision:    0,
			Scale:        0,
		},
		{
			Label:        "邮箱",
			Name:         "email",
			Column:       "email",
			Type:         "String",
			Comment:      "邮箱",
			Default:      "",
			IsDefaultRaw: false,
			IsPrimaryKey: false,
			NotNull:      true,
			Index:        "IDX_USER_EMAIL",
			Unique:       "UNI_USER_EMAIL",
			Length:       50,
			Precision:    0,
			Scale:        0,
		},
		{
			Label:        "手机号",
			Name:         "mobile",
			Column:       "mobile",
			Type:         "string",
			Comment:      "手机号",
			Default:      "",
			IsDefaultRaw: false,
			IsPrimaryKey: false,
			NotNull:      true,
			Index:        "IDX_USER_MOBILE",
			Unique:       "UNI_USER_MOBILE",
			Length:       11,
			Precision:    0,
			Scale:        0,
		},
		{
			Label:        "头像",
			Name:         "avatar",
			Column:       "avatar",
			Type:         "string",
			Comment:      "头像",
			Default:      "",
			IsDefaultRaw: false,
			IsPrimaryKey: false,
			NotNull:      false,
			Index:        "",
			Unique:       "",
			Length:       500,
			Precision:    0,
			Scale:        0,
		},
		{
			Label:        "性别",
			Name:         "gender",
			Column:       "gender",
			Type:         "string",
			Enum:         []string{"男", "女", "保密"},
			Comment:      "性别",
			Default:      "保密",
			IsDefaultRaw: false,
			IsPrimaryKey: false,
			NotNull:      false,
			Index:        "",
			Unique:       "",
			Length:       2,
			Precision:    0,
			Scale:        0,
		},
	}

	tests := []struct {
		name string
		args args
		want *Schema
	}{
		{
			name: "测试结构体解析",
			args: args{
				dest: `{
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
}`,
				d: mysqlDialect,
			},
			want: &Schema{
				Name:       "user",
				TableName:  "user",
				Fields:     fields,
				FieldNames: []string{"id", "username", "password", "nickname", "email", "mobile", "avatar", "gender"},
				fieldMap: map[string]*Field{
					"id":       fields[0],
					"username": fields[1],
					"password": fields[2],
					"nickname": fields[3],
					"email":    fields[4],
					"mobile":   fields[5],
					"avatar":   fields[6],
					"gender":   fields[7],
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.dest); got == nil || got.Name != tt.want.Name || len(got.Fields) != len(tt.want.Fields) {
				t.Errorf("Parse() = %v failed", got)
			}
		})
	}
}
