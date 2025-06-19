package illusionCard

import (
	"errors"

	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
	"github.com/GenesisAN/illusionsCard/KKS"
	"github.com/GenesisAN/illusionsCard/util"
)

func CardTypeRead(path string) (*util.PngBuff, error) {
	return util.PngRead(path)
}

// Deprecated: 请使用 ReadCardFromPath 替代。该函数将在未来版本中移除。 ReadKK 读取KK的服装卡片,传入卡片路径
func ReadKKClothes(pgb *util.PngBuff) (*KK.KKClothesCard, error) {
	if pgb.Type != Base.CT_KKC {
		return nil, errors.New("type error:" + pgb.Type)
	}
	card, err := KK.ParesKKClothes(pgb)
	if err != nil {
		return nil, err
	}
	card.Path = pgb.FilePath
	return card, nil
}

// Deprecated: 请使用 ReadCardFromPath 替代。该函数将在未来版本中移除。 ReadKK 读取KK的卡片,传入卡片路径
func ReadKK(pgb *util.PngBuff) (*KK.KKCharaCard, error) {
	if pgb.Type != Base.CT_KK {
		return nil, errors.New("type error:" + pgb.Type)
	}
	card, err := KK.ParseKKChara(pgb)
	if err != nil {
		return nil, err
	}
	card.Path = pgb.FilePath
	return &card, nil
}

// Deprecated: 请使用 ReadCardFromPath 替代。	该函数将在未来版本中移除。
func ReadKKS(pgb *util.PngBuff) (*KKS.SunshineCharaCard, error) {
	if pgb.Type != Base.CT_KKS {
		return nil, errors.New("type error:" + pgb.Type)
	}
	card, err := KKS.ParseKKSChara(pgb)
	if err != nil {
		return nil, err
	}
	card.Path = pgb.FilePath
	return &card, nil
}
