package illusionCard

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/GenesisAN/illusionsCard/KK"
	"github.com/GenesisAN/illusionsCard/KKS"
)

func TestReadAllCards(t *testing.T) {
	// 合并测试文件夹路径，可根据卡类型分文件夹
	testDirs := []string{
		"./KKTest/",  // 放置 KK 角色卡/服装卡
		"./KKSTest/", // 放置 KKS 角色卡
	}

	for _, dir := range testDirs {
		files := GetAllFiles(dir, ".png")
		for _, path := range files {
			t.Run(filepath.Base(path), func(t *testing.T) {
				card, err := ReadCardFromPath(path)
				if err != nil {
					t.Errorf("读取卡片失败: %v", err)
					return
				}
				fmt.Println("解析卡片成功：", card.GetPath())

				// 根据类型断言输出额外信息
				switch c := card.(type) {
				case *KK.KKCharaCard:
					fmt.Println("→ KK角色卡:", c.CharParmeter.Nickname)
					c.PrintZipmodeInfo()

				case *KK.KKClothesCard:
					fmt.Println("→ KK服装卡:", c.Path)
					c.PrintZipmodeInfo()

				case *KKS.SunshineCharaCard:
					fmt.Println("→ KKS角色卡:", c.CharParmeter.Nickname)
					c.PrintZipmodeInfo()

				default:
					t.Errorf("未知卡片类型: %T", card)

				}
			})
		}
	}
}
