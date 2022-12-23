// Validation
// @author xiangqian
// @date 21:52 2022/12/23
package app

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// 注册 去除字符串前后空格验证器
// E:/go/path/pkg/mod/github.com/gin-gonic/gin@v1.8.1/binding/.go:49
// Deprecated: 通过源码 github.com/gin-gonic/gin@v1.8.1/binding/default_validator.go:49 可知, 在参数校验函数中不能对结构体字段值进行修改, 因为传递的是结构体具体值, 而不是指针
func regTrimSpaceValidation() {
	if v, r := binding.Validator.Engine().(*validator.Validate); r {
		v.RegisterValidation("trim", func(fl validator.FieldLevel) bool {
			rflVal := fl.Field()
			str := rflVal.Interface().(string)
			rflVal = fl.Parent().FieldByName(fl.FieldName())
			rflVal.Set(reflect.ValueOf(strings.TrimSpace(str)))
			return true
		})
	}
}
