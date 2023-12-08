package schema

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yaochi-tech/lingquan-core-go/db/util"
	"strings"
)

type Field struct {
	Label        string
	Name         string
	Column       string
	Type         string
	Enum         []string
	Comment      string
	Default      string
	IsDefaultRaw bool
	IsPrimaryKey bool
	NotNull      bool
	Index        string
	Unique       string
	Length       uint
	Precision    uint
	Scale        uint
}

type Schema struct {
	Definition string
	Name       string
	TableName  string
	Comment    string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// RecordValues 获取按照模型字段顺序排列的字段数据
func (schema *Schema) RecordValues(dest map[string]interface{}) []interface{} {
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, dest[field.Name])
	}
	return fieldValues
}

// Parse 将定义的模型json转为schema
func Parse(definition string) *Schema {
	dj := gjson.Parse(definition)
	schema := &Schema{
		Definition: definition,
		Name:       dj.Get("code").String(),
		TableName:  util.ToSnake(dj.Get("code").String()),
		Comment:    dj.Get("comment").String(),
		fieldMap:   make(map[string]*Field),
	}

	fields := dj.Get("fields").Array()
	for _, f := range fields {
		// name/type/label
		field := new(Field)

		field.Label = f.Get("label").String()
		field.Name = f.Get("name").String()
		field.Column = util.ToSnake(field.Name)
		t := strings.ToLower(f.Get("type").String())
		if t == "enum" {
			field.Type = f.Get("enumType").String()
			for _, enum := range f.Get("enum").Array() {
				field.Enum = append(field.Enum, enum.String())
			}
		} else {
			field.Type = t
		}
		field.Comment = f.Get("comment").String()
		field.Default = f.Get("default_raw").String()
		field.IsDefaultRaw = field.Default != ""
		if !field.IsDefaultRaw {
			field.Default = f.Get("default").String()
		}
		field.IsPrimaryKey = t == "id"
		field.NotNull = f.Get("required").Bool()
		if idx := f.Get("index"); idx.IsBool() {
			field.Index = strings.ToUpper(fmt.Sprintf("IDX_%s_%s", schema.TableName, field.Column))
		} else {
			field.Index = idx.String()
		}
		if uni := f.Get("unique"); uni.IsBool() {
			field.Unique = strings.ToUpper(fmt.Sprintf("UNI_%s_%s", schema.TableName, field.Column))
		} else {
			field.Unique = uni.String()
		}
		field.Length = uint(f.Get("length").Uint())
		field.Precision = uint(f.Get("precision").Uint())
		field.Scale = uint(f.Get("scale").Uint())

		schema.Fields = append(schema.Fields, field)
		schema.FieldNames = append(schema.FieldNames, field.Column)
		schema.fieldMap[field.Name] = field
	}

	return schema
}
