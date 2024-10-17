package util

import "os"

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
