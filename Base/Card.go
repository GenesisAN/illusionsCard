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

func (c *Card) ToJson() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
