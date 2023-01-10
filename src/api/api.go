// common
// @author xiangqian
// @date 13:46 2022/12/22
package api

import (
	"auto-deploy-go/src/typ"
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
	"net/http"
	"strings"
)

var (
	zhTrans ut.Translator
	enTrans ut.Translator
)

func Init() {
	InitValidateTrans()
}

func InitValidateTrans() {
	if v, r := binding.Validator.Engine().(*validator.Validate); r {
		uni := ut.New(zh.New(), // 备用语言
			// 支持的语言
			zh.New(),
			en.New())
		if trans, r := uni.GetTranslator(typ.LocaleZh); r {
			zh_trans.RegisterDefaultTranslations(v, trans)
			zhTrans = trans
		}
		if trans, r := uni.GetTranslator(typ.LocaleEn); r {
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
		//case com.LocaleZh:
		//	validationErrTrans = errs.Translate(zhTrans)
		case typ.LocaleEn:
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
			if errMsg != "" {
				switch lang {
				case typ.LocaleEn:
					errMsg += ", "
				default:
					errMsg += "、"
				}
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

func HtmlPage(pContext *gin.Context, templateName string, pageFunc func(pContext *gin.Context, pageReq typ.PageReq) (any, gin.H, error)) {
	pageReq := typ.PageReq{Current: 1, Size: 10}
	err := ShouldBind(pContext, &pageReq)
	if err != nil {
		Html(pContext, templateName, func(pContext *gin.Context) (gin.H, error) {
			return gin.H{}, err
		})
		return
	}

	Html(pContext, templateName, func(pContext *gin.Context) (gin.H, error) {
		page, h, err := pageFunc(pContext, pageReq)
		if h == nil {
			h = gin.H{}
		}
		h["page"] = page
		return h, err
	})
}

func Html(pContext *gin.Context, templateName string, hkvFunc func(pContext *gin.Context) (gin.H, error)) {
	message := ""
	session := sessions.Default(pContext)
	smessage := session.Get("message")
	session.Delete("message")
	session.Save()
	if smessage != nil {
		message = fmt.Sprintf("%v", smessage)
	}

	h, err := hkvFunc(pContext)
	if err != nil {
		if message != "" {
			message += ", "
		}
		message += err.Error()
	}

	// user
	h["user"] = GetUser(pContext)

	// 没有消息就是最好的消息
	h["message"] = message

	pContext.HTML(http.StatusOK, templateName, h)
}

func Redirect(pContext *gin.Context, location string, message any, m map[string]any) {
	if message != nil {
		if v, r := message.(error); r {
			message = v.Error()
		}
		if m == nil {
			m = map[string]any{}
		}
		m["message"] = message
	}

	if m != nil {
		session := sessions.Default(pContext)
		for k, v := range m {
			session.Set(k, v)
		}
		session.Save()
	}

	pContext.Redirect(http.StatusMovedPermanently, location)
}
