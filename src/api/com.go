// common
// @author xiangqian
// @date 13:46 2022/12/22
package api

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

// 抽象实体定义
type Abs struct {
	Id         int64  `form:"id"`         // 主键id
	Rem        string `form:"rem"`        // 备注
	DelFlag    byte   `form:"delFlag"`    // 删除标识，0-正常，1-删除
	CreateTime int64  `form:"createTime"` // 创建时间（时间戳，s）
	UpdateTime int64  `form:"updateTime"` // 修改时间（时间戳，s）
}

const (
	LocalesZh = "zh"
	LocalesEn = "en"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

func InitValidateTrans() {
	if v, r := binding.Validator.Engine().(*validator.Validate); r {
		uni := ut.New(zh.New(), // 备用语言
			// 支持的语言
			zh.New(),
			en.New())
		if trans, r := uni.GetTranslator(LocalesZh); r {
			zh_trans.RegisterDefaultTranslations(v, trans)
			zhTrans = trans
		}
		if trans, r := uni.GetTranslator(LocalesEn); r {
			en_trans.RegisterDefaultTranslations(v, trans)
			enTrans = trans
		}
	}
}

func TransErr(pContext *gin.Context, err error) error {
	if errs, r := err.(validator.ValidationErrors); r {
		session := sessions.Default(pContext)
		lang := ""
		if v, r := session.Get("lang").(string); r {
			lang = v
		}
		var validationErrTrans validator.ValidationErrorsTranslations
		switch lang {
		case LocalesZh:
			validationErrTrans = errs.Translate(zhTrans)
		case LocalesEn:
			validationErrTrans = errs.Translate(enTrans)
		default:
			validationErrTrans = errs.Translate(zhTrans)
		}

		errMsg := ""
		for key, value := range validationErrTrans {
			name := key[strings.Index(key, ".")+1:]
			msg, ierr := i18n.GetMessage(fmt.Sprintf("i18n.%s", strings.ToLower(name)))
			if ierr == nil {
				value = strings.Replace(value, name, msg, 1)
			}
			errMsg += value
		}
		return errors.New(errMsg)
	}
	return err
}

func ShouldBind(pContext *gin.Context, i any) error {
	err := pContext.ShouldBind(i)
	if err != nil {
		err = TransErr(pContext, err)
	}
	return err
}
