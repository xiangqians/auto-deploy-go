// common
// @author xiangqian
// @date 22:31 2022/12/20
package com

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

// DataDir
// const DataDir = "./data"
const DataDir = "C:\\Users\\xiangqian\\Desktop\\tmp\\auto-deploy\\data"

// VerifyUsername 校验用户名
// 1-16位长度（字母，数字，下划线，减号）
func VerifyUsername(username string) error {
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,16}$", username)
	if err == nil && matched {
		return nil
	}

	return errors.New("用户名1-16位长度（字母，数字，下划线，减号）")
}

// VerifyPasswd 校验密码
// 8-16位长度（字母，数字，特殊字符）
func VerifyPasswd(passwd string) error {
	matched, err := regexp.MatchString("^[a-zA-Z0-9!@#$%^&*()-_=+]{8,16}$", passwd)
	if err == nil && matched {
		return nil
	}

	return errors.New("密码8-16位长度（字母，数字，特殊字符）")
}

// Uuid https://github.com/google/uuid
func Uuid() string {
	return uuid.New().String()
}

// NameHumpToUnderline 驼峰转下划线
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
