package dialect

import "reflect"

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    //用于将go自己的数据类型转换为数据库的数据类型
	TableExistSQL(tableName string) (string, []interface{}) // SQL语句，返回某个表是否存在
}

var dialectsMap = map[string]Dialect{}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}
