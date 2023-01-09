// index admin
// @author xiangqian
// @date 13:00 2023/01/08
package api

import (
	"auto-deploy-go/src/db"
	"auto-deploy-go/src/typ"
	"fmt"
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
	Title   []string // 标题
	Colspan int      // <td></td> colspan
	Data    []Data2  // 数据
}

type Data2 struct {
	typ.Abs
	Name string
	Data []any
}

func Page(pageReq typ.PageReq, tableName string) (any, Data, error) {
	switch tableName {
	case UserTableName:
		page, err := db.Page[typ.User](pageReq, "SELECT u.id, u.`name`, u.nickname, u.rem, u.del_flag, u.add_time, u.upd_time FROM `user` u")
		data := Data{}
		// title
		data.Title = []string{"i18n.user", "i18n.nickname", "i18n.passwd"}
		// <td></td> colspan
		data.Colspan = 7 + len(data.Title)
		// data
		if page.Data != nil && len(page.Data) > 0 {
			data2 := make([]Data2, len(page.Data))
			for i, user := range page.Data {
				data2[i] = Data2{
					Abs:  user.Abs,
					Name: user.Name,
					Data: []any{user.Name, user.Nickname, "******"},
				}
			}
			data.Data = data2
		}
		return page, data, err

	case RxTableName:
		page, err := db.Page[typ.Rx](pageReq, "SELECT r.id, r.`name`, r.owner_id, IFNULL(ou.`name`, '') AS 'owner_name',IFNULL(ou.nickname, '') AS 'owner_nickname', r.sharer_id, IFNULL(su.`name`, '') AS 'sharer_name', IFNULL(su.nickname, '') AS 'sharer_nickname', r.item_ids, COUNT(DISTINCT i.id) AS 'share_item_count', GROUP_CONCAT(i.`name`, '、') AS 'share_item_names', r.rem, r.del_flag, r.add_time, r.upd_time FROM rx r LEFT JOIN user ou ON ou.id = r.owner_id LEFT JOIN user su ON su.id = r.sharer_id LEFT JOIN item i ON r.item_ids LIKE ('%,' || i.id || ',%') GROUP BY r.id")
		data := Data{}
		// title
		data.Title = []string{
			"i18n.name",
			"i18n.owner",
			"i18n.sharer",
			"i18n.shareItemCount",
			"i18n.shareItemNames",
		}
		// <td></td> colspan
		data.Colspan = 7 + len(data.Title)
		// data
		if page.Data != nil && len(page.Data) > 0 {
			data2 := make([]Data2, len(page.Data))
			for i, rx := range page.Data {
				data2[i] = Data2{
					Abs:  rx.Abs,
					Name: rx.Name,
					Data: []any{
						rx.Name,
						fmt.Sprintf("%s, %s", rx.OwnerName, rx.OwnerNickname),
						fmt.Sprintf("%s, %s", rx.SharerName, rx.SharerNickname),
						rx.ShareItemCount,
						rx.ShareItemNames,
					},
				}
			}
			data.Data = data2
		}
		return page, data, err

	case GitTableName:
		page, err := db.Page[typ.Git](pageReq, "SELECT g.id, IFNULL(u.`name`, '') AS 'user_name', IFNULL(u.nickname, '') AS 'user_nickname', g.`name`, g.`user`, g.rem, g.del_flag, g.add_time, g.upd_time FROM git g LEFT JOIN user u ON u.id = g.user_id GROUP BY g.id")
		data := Data{}
		// title
		data.Title = []string{
			"i18n.name",
			"i18n.owner",
			"i18n.user",
			"i18n.passwd",
		}
		// <td></td> colspan
		data.Colspan = 7 + len(data.Title)
		// data
		if page.Data != nil && len(page.Data) > 0 {
			data2 := make([]Data2, len(page.Data))
			for i, git := range page.Data {
				data2[i] = Data2{
					Abs:  git.Abs,
					Name: git.Name,
					Data: []any{
						git.Name,
						fmt.Sprintf("%s, %s", git.UserName, git.UserNickname),
						git.User,
						"******",
					},
				}
			}
			data.Data = data2
		}
		return page, data, err

	case ServerTableName:
		page, err := db.Page[typ.Server](pageReq, "SELECT s.id, IFNULL(u.`name`, '') AS 'user_name', IFNULL(u.nickname, '') AS 'user_nickname', s.`name`, s.`host`, s.`port`, s.`user`, s.rem, s.del_flag, s.add_time, s.upd_time FROM server s LEFT JOIN user u ON u.id = s.user_id GROUP BY s.id")
		data := Data{}
		// title
		data.Title = []string{
			"i18n.name",
			"i18n.owner",
			"i18n.host",
			"i18n.port",
			"i18n.user",
			"i18n.passwd",
		}
		// <td></td> colspan
		data.Colspan = 7 + len(data.Title)
		// data
		if page.Data != nil && len(page.Data) > 0 {
			data2 := make([]Data2, len(page.Data))
			for i, server := range page.Data {
				data2[i] = Data2{
					Abs:  server.Abs,
					Name: server.Name,
					Data: []any{
						server.Name,
						fmt.Sprintf("%s, %s", server.UserName, server.UserNickname),
						server.Host,
						server.Port,
						server.User,
						"******",
					},
				}
			}
			data.Data = data2
		}
		return page, data, err

	case ItemTableName:
	case RecordTableName:
	}

	return typ.Page[any]{}, Data{}, nil
}
