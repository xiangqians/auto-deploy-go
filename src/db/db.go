// DB
// https://pkg.go.dev/github.com/mattn/go-sqlite3
// @author xiangqian
// @date 20:10 2022/12/21
package db

import (
	"auto-deploy-go/src/com"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
)

const dataSourceName = com.DataDir + "/autodeploy.db"

func db() (*sql.DB, error) {
	return sql.Open("sqlite3", dataSourceName)
}

func Qry(i any, sql string, args ...any) error {
	pDb, err := db()
	if err != nil {
		return err
	}
	defer pDb.Close()

	pRows, err := pDb.Query(sql, args...)
	if err != nil {
		return err
	}
	defer pRows.Close()

	err = rowsMapper(pRows, i)
	if err != nil {
		return err
	}

	return nil
}

func Add(sql string, args ...any) (int64, error) {
	return exec(sql, args...)
}

func Upd(sql string, args ...any) (int64, error) {
	return exec(sql, args...)
}

func Del(sql string, args ...any) (int64, error) {
	return exec(sql, args...)
}

// return affect
func exec(sql string, args ...any) (int64, error) {
	pDb, err := db()
	if err != nil {
		return 0, err
	}
	defer pDb.Close()

	res, err := pDb.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	res.LastInsertId()

	return res.RowsAffected()
}

// 字段集映射
// 支持 1）一个或多个属性映射；2）结构体映射；3）结构体切片映射
func rowsMapper(pRows *sql.Rows, i any) error {
	cols, err := pRows.Columns()
	if err != nil {
		return err
	}

	getDest := func(rflType reflect.Type, rflVal reflect.Value) []any {
		dest := make([]any, len(cols))
		for fi, fl := 0, rflType.NumField(); fi < fl; fi++ {
			typeField := rflType.Field(fi)
			name := typeField.Tag.Get("sql")
			if name == "" {
				name = com.NameHumpToUnderline(typeField.Name)
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
		return dest
	}

	rflType := reflect.TypeOf(i).Elem()
	rflVal := reflect.ValueOf(i).Elem()
	switch rflType.Kind() {
	// 结构体
	case reflect.Struct:
		if pRows.Next() {
			err = pRows.Scan(getDest(rflType, rflVal)...)
		}

	// 切片
	case reflect.Slice:
		eVal := rflVal.Index(0)
		l := rflVal.Len()
		switch eVal.Kind() {
		// 结构体切片
		case reflect.Struct:
			e := eVal.Interface()
			eRflType := reflect.TypeOf(e)
			var eRflVal reflect.Value
			idx := 0
			for pRows.Next() {
				if idx < l {
					eRflVal = rflVal.Index(idx).Addr().Elem()
				} else {
					pE := reflect.New(eRflType).Interface()
					eRflVal = reflect.ValueOf(pE).Elem()
				}
				err = pRows.Scan(getDest(eRflType, eRflVal)...)
				if err != nil {
					return err
				}

				// 切片（slice）扩容
				if idx >= l {
					rflVal0 := rflVal
					rflVal = reflect.Append(rflVal, eRflVal)
					rflVal0.Set(rflVal)
					rflVal = rflVal0
				}
				idx++
			}

		// 普通指针类型数组
		default:
			if pRows.Next() {
				dest := make([]any, l)
				for ei := 0; ei < l; ei++ {
					e := rflVal.Index(ei).Interface()
					dest[ei] = e
				}
				err = pRows.Scan(dest...)
			}
		}

	// 普通指针类型
	default:
		if pRows.Next() {
			err = pRows.Scan(i)
		}
	}

	return err
}
