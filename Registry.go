package illusionCard

import (
	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
	"github.com/GenesisAN/illusionsCard/KKS"
)

var cardReaderMap = map[string]CardReader{
	Base.CT_KK:  KK.KKCharaReader{},
	Base.CT_KKS: KKS.SunshineCharaReader{}, // 你自己定义
	Base.CT_KKC: KK.KKClothesReader{},      // 你自己定义
}

// GetCardReaderMap returns the cardReaderMap.
func GetCardReaderMap() map[string]CardReader {
	return cardReaderMap
}
