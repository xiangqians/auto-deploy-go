// index admin
// @author xiangqian
// @date 13:00 2023/01/08
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var tables []typ.Table

const (
	UserTableName   = "user"
	RxTableName     = "rx"
	GitTableName    = "git"
	ServerTableName = "server"
	ItemTableName   = "item"
	RecordTableName = "record"
)

func init() {
	tables = []typ.Table{
		{Name: UserTableName, Desc: "i18n.userTable"},
		{Name: RxTableName, Desc: "i18n.rxTable"},
		{Name: GitTableName, Desc: "i18n.gitTable"},
		{Name: ServerTableName, Desc: "i18n.serverTable"},
		{Name: ItemTableName, Desc: "i18n.itemTable"},
		{Name: RecordTableName, Desc: "i18n.recordTable"},
	}
}

func Table(name string) typ.Table {
	for _, table := range tables {
		if table.Name == name {
			return table
		}
	}
	return typ.Table{}
}

func IndexAdminPage(pContext *gin.Context) {
	tableName := strings.TrimSpace(pContext.Query("tableName"))
	if tableName == "" {
		session := sessions.Default(pContext)
		sTable := session.Get("tableName")
		if v, r := sTable.(string); r {
			tableName = v
		} else {
			tableName = UserTableName
		}
	} else {
		table := Table(tableName)
		if table.Name == "" {
			tableName = UserTableName
		}
		session := sessions.Default(pContext)
		session.Set("tableName", tableName)
		session.Save()
	}

	html := func(page any, data Data, err error) {
		message := ""
		if err != nil {
			message = err.Error()
		}
		pContext.HTML(http.StatusOK, "index_admin.html", gin.H{
			"user":    GetUser(pContext),
			"table":   Table(tableName),
			"tables":  tables,
			"page":    page,
			"data":    data,
			"message": message,
		})
	}

	pageReq := typ.PageReq{Current: 1, Size: 10}
	err := ShouldBind(pContext, &pageReq)
	if err != nil {
		html(typ.Page[int]{}, Data{}, err)
		return
	}

	page, data, err := Page(pageReq, tableName)
	html(page, data, err)
	return
}

type Data struct {
	TdColspan int8     // <td></td> colspan
	Title     []string // 标题
	Data      []Data2  // 数据
}

type Data2 struct {
	typ.Abs
	Name string
	Data []any
}

func Page(pageReq typ.PageReq, tableName string) (any, Data, error) {
	switch tableName {
	case UserTableName:
		page, err := db.Page[typ.User](pageReq.Current, pageReq.Size, "SELECT u.id, u.`name`, u.nickname, u.rem, u.del_flag, u.add_time, u.upd_time FROM `user` u")
		data := Data{}
		// <td></td> colspan
		data.TdColspan = 7 + 3
		// title
		data.Title = []string{"i18n.user", "i18n.nickname"}
		// data
		if page.Data != nil && len(page.Data) > 0 {
			data2 := make([]Data2, len(page.Data))
			for i, user := range page.Data {
				data2[i] = Data2{
					Abs:  user.Abs,
					Name: user.Name,
					Data: []any{user.Name, user.Nickname},
				}
			}
			data.Data = data2
		}
		return page, data, err

	case RxTableName:
	case GitTableName:
	case ServerTableName:
	case ItemTableName:
	case RecordTableName:
	}

	return typ.Page[any]{}, Data{}, nil
}
