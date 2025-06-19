package Base

import (
	"encoding/json"
	"fmt"
)

type CardInterface interface {
	GetPath() string
	GetVersion() string
	PrintCardInfo()
}

// Card is illusionCard Bset struct
type Card struct {
	Extended     map[string]*PluginData   `json:"-"`
	ExtendedList map[string]*PluginDataEx `json:"extended_list"`
	CharInfo     *ChaFileParameterEx      `json:"char_info"`
	Image1       *[]byte                  `json:"-"`
	Image2       *[]byte                  `json:"-"`
	CardType     string                   `json:"card_type"`
	LoadVersion  string                   `json:"load_version"`
	Path         string                   `json:"path"`
	MD5          string                   `json:"md5"`
}

type ChaFileParameterEx struct {
	Version   string `json:"version"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Nickname  string `json:"nickname"`
	Sex       int    `json:"sex"`
}

// TypeInt 将Card结构体中的CardType转换为int类型的函数
func (c *Card) TypeInt() int {
	switch c.CardType {
	case CT_KK:
		return CTI_Koikatu
	case CT_KKP:
		return CTI_KoikatsuParty
	case CT_KKCSP:
		return CTI_KoikatsuPartySpecialPatch
	case CT_EC:
		return CTI_EmotionCreators
	case CT_AI:
		return CTI_AiSyoujyo
	case CT_KKS:
		return CTI_KoikatsuSunshine
	default:
		return CTI_Unknown
	}
}

// Card打印zipmode信息
func (c *Card) PrintZipmodeInfo() {
	if c.ExtendedList == nil || len(c.ExtendedList) == 0 {
		return
	}

	printedGUIDs := make(map[string]bool)

	for _, pluginDataEx := range c.ExtendedList {
		if pluginDataEx.RequiredZipmodGUIDs == nil || len(pluginDataEx.RequiredZipmodGUIDs) == 0 {
			continue
		}

		for i, mod := range pluginDataEx.RequiredZipmodGUIDs {
			if _, exists := printedGUIDs[mod.GUID]; exists {
				continue // 已打印，跳过
			}
			printedGUIDs[mod.GUID] = true // 标记为已打印

			fmt.Printf("  *[mod依赖 %d]: %s (%s | LS: %d | CN: %d)\n",
				i, mod.GUID, mod.Property, mod.LocalSlot, mod.CategoryNo)
		}
	}
}

// ToJson 将Card结构体转换为json字符串的函数
func (c *Card) ToJson() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 将json字符串转换为Card结构体的函数
func (c *Card) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), c)
}
