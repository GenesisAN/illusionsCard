package illusionCard

import (
	"os"
	"path/filepath"
	"testing"
)

// KK测试用例，需在根目录创建KKTest文件夹，并在里面放置卡片文件
func TestReadKK(t *testing.T) {
	files := GetAllFiles("./KKTest/", ".png")
	for _, v := range files {
		kk, err := ReadKK(v)
		if err != nil {
			t.Error(err)
		}
		println(kk.CharParmeter.Nickname)
	}
}

// KKS测试用例，需在根目录创建KKTest文件夹，并在里面放置卡片文件
func TestReadKKS(t *testing.T) {
	files := GetAllFiles("./KKSTest/", ".png")
	for _, v := range files {
		kks, err := ReadKKS(v)
		if err != nil {
			t.Error(err)
		}
		println(kks.CharParmeter.Nickname)
	}

}

func GetAllFiles(root, ext string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ext {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
