package Base

import "encoding/json"

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
