package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field represents one column information
type Field struct {
	Name string //字段名
	Type string //类型
	Tag  string //约束条件
}

// Schema
type Schema struct {
	Model      interface{}       //被映射的对象
	Name       string            //表名
	Fields     []*Field          //字段
	FieldNames []string          // 包含所有的字段（列）名
	fieldMap   map[string]*Field //字段名与Field的映射关系，方便以后直接使用
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)                     // Field returns a struct type's i'th field.
		if !p.Anonymous && ast.IsExported(p.Name) { // p.Anonymous 是否embed; ast.IsExported reports whether name starts with an upper-case letter.
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues Change struct to []interface{}
// eg: &User{Name: "Tom", Age: 18} -> ("Tom", 18)
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
