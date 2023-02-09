package Koikatsu

//go:generate msgp
type BlockHeaderListInfo struct {
	LstInfo []*BlockHeader `msg:"lstInfo"` // Card BlockHeader Info List
}

// BlockHeader 卡片头部数据结构
type BlockHeader struct {
	Name    string `msg:"name"`    // 插件名称
	Version string `msg:"version"` // 版本
	Pos     int    `msg:"pos"`     // 数据起始位置
	Size    int    `msg:"size"`    // 数据大小
	Data    []byte // 数据内容
}
