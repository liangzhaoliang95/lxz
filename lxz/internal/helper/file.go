package helper

import (
	"os"
	"path/filepath"
)

// 递归获取文件夹下所有文件 过滤掉http-cache文件夹
func ListDefaultK8sConfigFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 过滤掉http-cache文件夹
		if info.IsDir() && (filepath.Base(path) == "http-cache" || filepath.Base(path) == "cache") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
