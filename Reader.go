package illusionCard

import (
	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/util"
)

// CardReader 是解析各种卡片类型的策略接口
type CardReader interface {
	Read(*util.PngBuff) (Base.CardInterface, error)
}
