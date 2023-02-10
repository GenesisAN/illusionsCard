package illusionCard

import (
	"errors"
	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/KK"
	"github.com/GenesisAN/illusionsCard/KKS"
	"github.com/GenesisAN/illusionsCard/util"
)

// ReadKK 读取KK的卡片,传入卡7片路径
func ReadKK(path string) (*KK.KKCard, error) {
	pb, err := util.PngRead(path)
	if err != nil {
		return nil, err
	}
	if pb.Type != Base.CT_KK {
		return nil, errors.New("type error:" + pb.Type)
	}
	card, err := KK.ParseKKChara(pb)
	if err != nil {
		return nil, err
	}
	card.Path = path
	return &card, nil
	//版本号
}

func ReadKKS(path string) (*KKS.KKSCard, error) {
	pb, err := util.PngRead(path)
	if err != nil {
		return nil, err
	}
	if pb.Type != Base.CT_KKS {
		return nil, errors.New("type error:" + pb.Type)
	}
	card, err := KKS.ParseKKSChara(pb)
	if err != nil {
		return nil, err
	}
	card.Path = path
	return &card, nil
}
