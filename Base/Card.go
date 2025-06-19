package Base

import (
	"encoding/json"
	"fmt"
)

type CardInterface interface {
	CompareMissingZipMods(localGUIDs []string) []string
	GetZipmodsDependencies() []string
	GetPath() string
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
func (k *Card) GetPath() string {
	return k.Path
}

// CompareMissingMods 对比卡片所需的 zipmod 插件与本地 mods，返回缺失插件的 GUID 映射
// CompareMissingMods compares the card's required zipmod GUIDs with local GUID list.
// Returns a map of missing GUID -> ResolveInfo.
func (c *Card) CompareMissingMods(localGUIDs []string) map[string]ResolveInfo {
	missing := make(map[string]ResolveInfo)
	if c.ExtendedList == nil {
		return missing
	}

	// 构造本地 GUID 快速查询 map
	localModSet := make(map[string]struct{}, len(localGUIDs))
	for _, guid := range localGUIDs {
		localModSet[guid] = struct{}{}
	}

	// 提取卡片中 universalautoresolver 的插件依赖
	if v, ok := c.ExtendedList["com.bepis.sideloader.universalautoresolver"]; ok {
		for _, mod := range v.RequiredZipmodGUIDs {
			if _, found := localModSet[mod.GUID]; !found {
				missing[mod.GUID] = mod
			}
		}
	}
	return missing
}

// 返回卡牌依赖的DLL信息
func (c *Card) GetDLLDependencies() []string {
	if c.ExtendedList == nil || len(c.ExtendedList) == 0 {
		return nil
	}
	var dependencies []string
	for i, _ := range c.ExtendedList {
		dependencies = append(dependencies, i)
	}
	return dependencies
}

// 返回卡片的zipmod依赖信息
func (c *Card) GetZipmodsDependencies() []string {
	if c.ExtendedList == nil || len(c.ExtendedList) == 0 {
		return nil
	}

	var dependencies []string
	for _, pluginDataEx := range c.ExtendedList {
		if pluginDataEx.RequiredZipmodGUIDs == nil || len(pluginDataEx.RequiredZipmodGUIDs) == 0 {
			continue
		}

		for _, mod := range pluginDataEx.RequiredZipmodGUIDs {
			dependencies = append(dependencies, mod.GUID)
		}
	}
	return dependencies
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
			fmt.Printf("  *[mod依赖 %d]: %s\n",
				i, mod.GUID)
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
