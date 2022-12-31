// common
// @author xiangqian
// @date 22:31 2022/12/20
package util

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/gin-contrib/i18n"
	"github.com/google/uuid"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func VerifyText(t string, maxLen int) error {
	if t == "" {
		return errors.New(i18n.MustGetMessage("i18n.anyCannotEmpty"))
	}

	if len(t) > maxLen {
		return errors.New(fmt.Sprintf(i18n.MustGetMessage("i18n.anyGtNChar"), "%v", maxLen))
	}

	return nil
}

// VerifyUserName 校验用户名
// 1-16位长度（字母，数字，下划线，减号）
func VerifyUserName(username string) error {
	if username == "" {
		return errors.New(i18n.MustGetMessage("i18n.userCannotEmpty"))
	}

	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]{1,16}$", username)
	if err == nil && matched {
		return nil
	}

	return errors.New(i18n.MustGetMessage("i18n.userNameMastNBitsLong"))
}

// VerifyPasswd 校验密码
// 8-16位长度（字母，数字，特殊字符）
func VerifyPasswd(passwd string) error {
	matched, err := regexp.MatchString("^[a-zA-Z0-9!@#$%^&*()-_=+]{8,16}$", passwd)
	if err == nil && matched {
		return nil
	}

	return errors.New(i18n.MustGetMessage("i18n.passwdMastNBitsLong"))
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

// IsExistOfPath 判断所给路径（文件/文件夹）是否存在
func IsExistOfPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// DelFile 删除文件
func DelFile(path string) error {
	return os.Remove(path)
}

// DelDir 删除文件夹
func DelDir(path string) error {
	return os.RemoveAll(path)
}

// Mkdir 创建目录
func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// Command
func Command(cmd string) (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/C", cmd), nil

	case "linux":
		return exec.Command("bash", "-c", cmd), nil

	default:
		return nil, errors.New(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
	}
}

func Cd(path string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("cd /d %s", path), nil

	case "linux":
		return fmt.Sprintf("cd %s", path), nil

	default:
		return "", errors.New(fmt.Sprintf("The current system is not supported, %v", runtime.GOOS))
	}
}

// Zip 创建zip归档文件
// prefix: zip归档前缀
// dstName: 创建zip归档文件名
// srcNames: 待归档的文件集
func Zip(prefix string, dstName string, srcNames ...string) error {
	pDstFile, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer pDstFile.Close()

	// 创建一个压缩文档
	pZipWriter := zip.NewWriter(pDstFile)
	// 关闭压缩文档
	defer pZipWriter.Close()

	for _, srcName := range srcNames {
		pSrcFile, err := os.Open(srcName)
		if err != nil {
			continue
		}

		err = zip0(pZipWriter, prefix, pSrcFile)
		pSrcFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func zip0(pZipWriter *zip.Writer, prefix string, pFile *os.File) error {
	fileInfo, err := pFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		if prefix == "" {
			prefix = fileInfo.Name()
		} else {
			prefix += "/" + fileInfo.Name()
		}

		subFileInfos, err := pFile.Readdir(-1)
		if err != nil {
			return err
		}

		for _, subFileInfo := range subFileInfos {
			pSubFile, err := os.Open(pFile.Name() + "/" + subFileInfo.Name())
			if err != nil {
				return err
			}

			err = zip0(pZipWriter, prefix, pSubFile)
			pSubFile.Close()
			if err != nil {
				return err
			}
		}

		return nil
	}

	pFileHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	if prefix != "" {
		pFileHeader.Name = prefix + "/" + pFileHeader.Name
	}
	pWriter, err := pZipWriter.CreateHeader(pFileHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(pWriter, pFile)
	return err
}
