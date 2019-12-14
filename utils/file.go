package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
)

func CopyDir(srcPath string, destPath string) error {
	if srcInfo, err := os.Stat(srcPath); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath不是一个正确的目录！")
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo不是一个正确的目录！")
			return e
		}
	}
	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			copyFile(path, destNewPath)
		}
		return nil
	})
	return err
}

//生成目录并拷贝文件
func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := PathExists(destSplitPath)
			if b == false {
				//创建目录
				_ = os.Mkdir(destSplitPath, os.ModePerm)
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

//检测文件夹路径时候存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ReadAll(path string) (map[string]string, error) {
	if fileInfo, err := os.Stat(path); err != nil {
		return nil, err
	} else {
		res := make(map[string]string, 0)
		if fileInfo.IsDir() {
			err = filepath.Walk(path, func(curPath string, info os.FileInfo, err error) error {
				if !info.IsDir() && !isHide(curPath) {
					if content, err := ioutil.ReadFile(curPath); err != nil {
						return err
					} else {
						res[curPath] = string(content)
					}
				}
				return nil
			})
		}
		return res, err
	}
}

func isHide(path string) bool{
	paths := strings.Split(path, string(filepath.Separator))
	return strings.HasPrefix(paths[len(paths) - 1], ".")
}
