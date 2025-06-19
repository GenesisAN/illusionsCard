package illusionCard

import (
	"errors"

	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/util"
)

// ReadCardFromPath 是最终用户调用的统一工厂入口：通过路径读取并解析卡片
func ReadCardFromPath(path string) (Base.CardInterface, error) {
	pgb, err := util.PngRead(path)
	if err != nil {
		return nil, err
	}
	return ReadCard(pgb)
}

// ReadCard 是工厂方法的核心逻辑，使用策略模式读取正确类型的卡片
func ReadCard(pgb *util.PngBuff) (Base.CardInterface, error) {
	readerMap := GetCardReaderMap()
	reader, ok := readerMap[pgb.Type]
	if !ok {
		return nil, errors.New("unsupported card type: " + pgb.Type)
	}
	return reader.Read(pgb)
}
