// DB
// @author xiangqian
// @date 20:10 2022/12/21
package db

import (
	"auto-deploy-go/src/com"
	"database/sql"
	"reflect"
	"regexp"
	"strings"
)

// 字段集映射
func RowsMapper(pRows *sql.Rows, i interface{}) error {
	cols, err := pRows.Columns()
	com.CheckErr(err)

	dest := make([]any, len(cols))
	rflType := reflect.TypeOf(i).Elem()
	rflVal := reflect.ValueOf(i).Elem()
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

// 驼峰转下划线
func NameHumpToUnderline(name string) string {
	pRegexp := regexp.MustCompile("([A-Z])")
	r := pRegexp.FindAllIndex([]byte(name), -1)
	var res = make([]string, len(r))
	resIdx := 0
	index := 0
	for i, v := range r {
		if i > 0 {
			res[resIdx] = name[index:v[0]]
			resIdx++
			index = v[0]
		}
	}
	res[resIdx] = name[index:]
	return strings.ToLower(strings.Join(res, "_"))
}
