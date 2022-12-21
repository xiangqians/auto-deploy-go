// DB
// @author xiangqian
// @date 20:10 2022/12/21
package db

import (
	"auto-deploy-go/src/com"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"
)

const dataSourceName = com.DataDir + "/autodeploy.db"

func Db() (*sql.DB, error) {
	return sql.Open("sqlite3", dataSourceName)
}

func Qry(t any, sql string, args ...any) error {
	pDb, err := Db()
	if err != nil {
		return err
	}
	defer pDb.Close()

	pRows, err := pDb.Query(sql, args)
	if err != nil {
		return err
	}
	defer pRows.Close()

	for pRows.Next() {
		err = RowsMapper(pRows, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func Add(sql string, args ...any) (int64, error) {
	return exec(sql, args)
}

func Upd(sql string, args ...any) (int64, error) {
	return exec(sql, args)
}

func Del(sql string, args ...any) (int64, error) {
	return exec(sql, args)
}

// return affect
func exec(sql string, args ...any) (int64, error) {
	pDb, err := Db()
	if err != nil {
		return 0, err
	}
	defer pDb.Close()

	res, err := pDb.Exec(sql, args)
	if err != nil {
		return 0, err
	}
	res.LastInsertId()

	return res.RowsAffected()
}

// 字段集映射
func RowsMapper(pRows *sql.Rows, t any) error {
	rflType := reflect.TypeOf(t).Elem()
	switch rflType.Kind() {
	case reflect.Struct:
		return RowsMapperToStruct(pRows, t, rflType)

	case reflect.Map:
		pMap := t.(*map[string]any)
		return RowsMapperToMap(pRows, pMap)

	default:
		return errors.New("mapping is not supported")
	}
}

// 字段集映射为 Map
func RowsMapperToMap(pRows *sql.Rows, pMap *map[string]any) error {
	cols, err := pRows.Columns()
	com.CheckErr(err)

	dest := make([]any, len(cols))
	for ci, _ := range cols {
		var v any
		dest[ci] = &v
	}

	err = pRows.Scan(dest...)
	if err == nil {
		for ci, col := range cols {
			v := dest[ci]
			(*pMap)[NameUnderlineToHump(col)] = *(v.(*any))
		}
	}
	return err
}

// 字段集映射为 Struct
func RowsMapperToStruct(pRows *sql.Rows, t any, rflType reflect.Type) error {
	cols, err := pRows.Columns()
	com.CheckErr(err)

	dest := make([]any, len(cols))
	//rflType := reflect.TypeOf(t).Elem()
	rflVal := reflect.ValueOf(t).Elem()
	for fi, fl := 0, rflType.NumField(); fi < fl; fi++ {
		typeField := rflType.Field(fi)
		name := typeField.Tag.Get("sql")
		if name == "" {
			name = NameHumpToUnderline(typeField.Name)
		}
		for ci, col := range cols {
			if col == name {
				valField := rflVal.Field(fi)
				if valField.CanAddr() {
					dest[ci] = valField.Addr().Interface()
				}
				break
			}
		}
	}
	return pRows.Scan(dest...)
}

// 下划线转驼峰
func NameUnderlineToHump(name string) string {
	pRegexp := regexp.MustCompile("_")
	r := pRegexp.FindAllIndex([]byte(name), -1)
	l := len(r)

	toUpper := func(s string) string {
		if s == "" {
			return s
		}
		return strings.ToUpper(s[0:1]) + s[1:]
	}

	if l == 0 {
		return toUpper(name)
	}

	var res = make([]string, l+1)
	resIdx := 0
	index := 0
	for _, v := range r {
		s := name[index:v[0]]
		if s != "" {
			res[resIdx] = toUpper(s)
			resIdx++
		}
		index = v[0] + 1
	}
	res[resIdx] = toUpper(name[index:])

	for i, s := range res {
		if s == "" {
			res = res[0:i]
			break
		}
	}
	return strings.Join(res, "")
}

// 驼峰转下划线
func NameHumpToUnderline(name string) string {
	pRegexp := regexp.MustCompile("([A-Z])")
	r := pRegexp.FindAllIndex([]byte(name), -1)
	l := len(r)
	if l == 0 {
		return strings.ToLower(name)
	}

	var res = make([]string, l+1)
	resIdx := 0
	index := 0
	for _, v := range r {
		s := name[index:v[0]]
		if s != "" {
			res[resIdx] = s
			resIdx++
		}
		index = v[0]
	}
	res[resIdx] = name[index:]
	for i, s := range res {
		if s == "" {
			res = res[0:i]
			break
		}
	}
	return strings.ToLower(strings.Join(res, "_"))
}
