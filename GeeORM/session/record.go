package session

import (
	"errors"
	"geeorm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find: 传入一个切片的指针，查询的结果保存在切片中
// eg: var users []User
// s.Find(&Users)
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()                                   // 获取切片的单个元素类型 destType
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable() // 映射出表结构 RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT) // clause 构建 select 语句
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() { //enumerate rows
		dest := reflect.New(destType).Elem() // 创建 destType 的实例 dest，
		var values []interface{}
		for _, name := range table.FieldNames { // 将 dest 的所有字段平铺开到 values
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil { // 用 rows.Scan() 将该行记录的每一列的值依次赋值给 values 中的每一个字段
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// Update support map[string]interface{}
// or k-v list : "Name", "Tom", "Age", 18
func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

/* 链式调用是一种简化代码的编程方式，能够使代码更简洁、易读。
链式调用的原理也非常简单，某个对象调用某个方法后，将该对象的引用/指针返回，即可以继续调用该对象的其他方法。
通常来说，当某个对象需要一次调用多个方法来设置其属性时，就非常适合改造为链式调用了。
s.Where("Age > 18").Limit(3).Find(&users)
*/

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

//OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
