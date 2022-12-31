// DB
// https://pkg.go.dev/github.com/mattn/go-sqlite3
// @author xiangqian
// @date 20:10 2022/12/21
package db

import (
	"auto-deploy-go/src/arg"
	"auto-deploy-go/src/util"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
)

func db() (*sql.DB, error) {
	return sql.Open("sqlite3", arg.Db)
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

// Add return insertId
func Add(sql string, args ...any) (int64, error) {
	pDb, err := db()
	if err != nil {
		return 0, err
	}
	defer pDb.Close()

	res, err := pDb.Exec(sql, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
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

	return res.RowsAffected()
}

// 字段集映射
// 支持 1）一个或多个属性映射；2）结构体映射；3）结构体切片映射
func rowsMapper(pRows *sql.Rows, i any) error {
	cols, err := pRows.Columns()
	if err != nil {
		return err
	}

	rflType := reflect.TypeOf(i).Elem()
	rflVal := reflect.ValueOf(i).Elem()
	switch rflType.Kind() {
	// 结构体
	case reflect.Struct:
		if pRows.Next() {
			dest := getDest(cols, rflType, rflVal)
			err = pRows.Scan(dest...)
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
				dest := getDest(cols, eRflType, eRflVal)
				err = pRows.Scan(dest...)
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

func getDest(cols []string, rflType reflect.Type, rflVal reflect.Value) []any {
	dest := make([]any, len(cols))
	for fi, fl := 0, rflType.NumField(); fi < fl; fi++ {
		typeField := rflType.Field(fi)

		// 兼容 FieldAlign() int （如果是struct字段，对齐后占用的字节数）
		if typeField.Type.Kind() == reflect.Struct {
			for sfi, sfl := 0, typeField.Type.NumField(); sfi < sfl; sfi++ {
				setDest(cols, &dest, typeField.Type.Field(sfi), rflVal)
			}
		} else {
			setDest(cols, &dest, typeField, rflVal)
		}
	}

	return dest
}

func setDest(cols []string, dest *[]any, typeField reflect.StructField, rflVal reflect.Value) {
	name := typeField.Tag.Get("sql")
	if name == "" {
		name = util.NameHumpToUnderline(typeField.Name)
	}
	for ci, col := range cols {
		if col == name {
			valField := rflVal.FieldByName(typeField.Name)
			if valField.CanAddr() {
				(*dest)[ci] = valField.Addr().Interface()
			}
			break
		}
	}
}
