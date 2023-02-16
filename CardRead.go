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

// ReadKK 读取KK的卡片,传入卡7片路径
func ReadKK(pgb *util.PngBuff) (*KK.KKCard, error) {
	if pgb.Type != Base.CT_KK {
		return nil, errors.New("type error:" + pgb.Type)
	}
	card, err := KK.ParseKKChara(pgb)
	if err != nil {
		return nil, err
	}
	card.Path = pgb.FilePath
	return &card, nil
	//版本号
}

func ReadKKS(pgb *util.PngBuff) (*KKS.KKSCard, error) {
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
