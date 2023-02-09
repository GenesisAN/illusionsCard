package illusionCard

import (
	"bytes"
	"errors"
	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/Koikatsu"
	"github.com/GenesisAN/illusionsCard/Tools"
	"os"
	"strings"
)

// ReadKK 读取KK的卡片,传入卡7片路径
func ReadKK(path string) (*Koikatsu.KoiCard, error) {
	path = strings.Replace(path, "\\", "/", -1)
	//提取文件名
	//读取图片
	f, err := os.ReadFile(path)
	buffer := bytes.NewBuffer(f)
	if err != nil {
		return nil, err
	}
	//value := get_png(f)
	//切割图片
	pngend := bytes.Index(f, Base.PngEndChunk) + len(Base.PngEndChunk)
	if pngend == -1 {
		return nil, errors.New("not found PngEndChunk")
	}
	outpng := buffer.Next(pngend)
	//os.WriteFile("Out.png", outpng, 0776)
	fb, err := buffer.ReadByte()
	if err != nil {
		return nil, errors.New("read fail:first byte not found")
	}
	if fb == 0x7 {
		_, err = Tools.BufRead(buffer, 64, "read fail:0x7 BufRead fail")
		if err != nil {
			return nil, err
		}
	} else if fb == 0x64 {
		_, err = Tools.BufRead(buffer, 3, "read fail:0x64 BufRead fail")
		if err != nil {
			return nil, err
		}
	}
	types, err := buffer.ReadByte()
	if err != nil {
		return nil, errors.New("read fail:unknown card type string len")
	}
	cardtypebyte, err := Tools.BufRead(buffer, int(types), "read fail:card type string")
	if err != nil {
		return nil, err
	}
	cardtype := string(cardtypebyte)
	card, err := Koikatsu.ParseKoiChara(buffer)
	if err != nil {
		return nil, err
	}
	card.Image = outpng
	card.Path = path
	card.CardType = cardtype
	return &card, nil
	//版本号
}
