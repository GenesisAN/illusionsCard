// Package KK 用于解析Koikatsu的角色卡数据
package KK

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/GenesisAN/illusionsCard/Base"
	util "github.com/GenesisAN/illusionsCard/util"
)

type KKCard struct {
	*Base.Card
	CharParmeter *KKChaFileParameter
}

func (card *KKCard) KKChaFileParameterEx(cfp *KKChaFileParameter) {
	card.CharParmeter = cfp
	card.CharInfo = &Base.ChaFileParameterEx{}
	card.CharInfo.Lastname = cfp.Lastname
	card.CharInfo.Firstname = cfp.Firstname
	card.CharInfo.Version = cfp.Version
	card.CharInfo.Nickname = cfp.Nickname
}

func ParesKKClothes(pb *util.PngBuff) (*KKCard, error) {
	kkc := KKCard{&Base.Card{}, &KKChaFileParameter{}}
	Version, err := pb.StringRead()
	if err != nil {
		return &kkc, err
	}
	kkc.LoadVersion = Version
	coordinateName, err := pb.StringRead() // 读取坐标名称
	if err != nil {
		return &kkc, err
	}
	fmt.Println("Coordinate Name:", coordinateName)
	//kkc.CharInfo.Nickname = coordinateName
	num, err := pb.Int32Read()
	//load clothes and accs from the bytes array

	pb.BuffRead(int(num), "BuffRead Fail:Card BlockHeader")

	marker, err := pb.StringRead() //name
	if err != nil {
		return &kkc, err
	}
	Version, err = pb.StringRead() //version
	if err != nil {
		return &kkc, err
	}
	len, err := pb.Int32Read()
	if err != nil {
		return &kkc, err
	}
	if marker != "KKEx" {
		return &kkc, errors.New("marker not found: " + marker)
	}
	if len <= 0 {
		return &kkc, errors.New("KKEx length is zero")
	}

	var extDataO Base.MapSArrayInterface
	extDataBuf, err := pb.BuffRead(int(len), "BuffRead Fail:KKEx data")
	if err != nil {
		return &kkc, err
	}
	extDataO.UnmarshalMsg(extDataBuf)
	//exDataO处理后的数据 exData
	exData := make(map[string]*Base.PluginData)
	for S, v := range extDataO {
		if v != nil {
			//取出PluginData
			var pd Base.PluginData
			pd.Version = int(v[0].(int64))
			pd.Data = v[1]
			exData[S] = &pd
		}
	}
	// 遍历exData提取RequiredZipmodGUIDs
	exDataEx := make(map[string]*Base.PluginDataEx)
	for s, data := range exData {
		// 根据GUID找出插件对应的Data
		dex := data.DeserializeObjects()
		dex.Name = s
		dex.Version = data.Version
		exDataEx[dex.Name] = &dex
	}
	kkc.Extended = exData

	return &kkc, nil
}

func ParseKKChara(pb *util.PngBuff) (KKCard, error) {
	kc := KKCard{&Base.Card{}, &KKChaFileParameter{}}
	Version, err := pb.StringRead()
	if err != nil {
		return kc, err
	}
	kc.LoadVersion = Version
	flbui32, err := pb.UInt32Read()
	if err != nil {
		return kc, err
	}
	png, err := pb.BuffRead(int(flbui32), "BuffRead Fail:Face img")
	if err != nil {
		return kc, err
	}
	pb.Png2 = &png
	countBuf, err := pb.BuffRead(4, "BuffRead Fail:Card BlockHeader len")
	if err != nil {
		return kc, err
	}

	var count = binary.LittleEndian.Uint32(countBuf)
	bhbytes, err := pb.BuffRead(int(count), "BuffRead Fail:Card BlockHeader")
	if err != nil {
		return kc, err
	}

	var bhls BlockHeaderListInfo
	_, err = bhls.UnmarshalMsg(bhbytes)
	if err != nil {
		return kc, errors.New("BlockHeaderListInfo Unmarshal Fail")
	}

	sort.SliceStable(bhls.LstInfo, func(i, j int) bool {
		return bhls.LstInfo[i].Pos < bhls.LstInfo[j].Pos
	})
	pb.B.Next(8)
	bhmap := make(map[string]*BlockHeader)
	for _, bh := range bhls.LstInfo {
		cBuf, err := pb.BuffRead(int(bh.Size), "BuffRead Fail:"+bh.Name)
		if err != nil {
			return kc, err
		}
		bh.Data = cBuf
		bhmap[bh.Name] = bh
	}
	//遍历 头部信息，获取Parameter位置
	par, ok := bhmap["Parameter"]
	if ok {
		if par.Version != "0.0.5" {
			return kc, errors.New("BlockHeaderListInfo Unmarshal Fail")
		}
	} else {
		return kc, errors.New("parameter not found")
	}
	//根据位置信息，反序列化 MsgPack 得到 KKChaFileParameter
	var Cfp KKChaFileParameter
	_, err = Cfp.UnmarshalMsg(par.Data)
	if err != nil {
		return kc, errors.New("KKChaFileParameter Unmarshal Fail")
	}
	kc.KKChaFileParameterEx(&Cfp)
	kkex, kkexok := bhmap["KKEx"]
	//根据KKEx位置信息，反序列化 得到 extDataO
	var extDataO Base.MapSArrayInterface
	if kkexok {
		extDataO.UnmarshalMsg(kkex.Data)
	}
	//exDataO处理后的数据 exData
	exData := make(map[string]*Base.PluginData)
	kc.Extended = exData
	//exData原始数据
	for S, v := range extDataO {
		if v != nil {
			//取出PluginData
			var pd Base.PluginData
			pd.Version = int(v[0].(int64))
			pd.Data = v[1]
			exData[S] = &pd
		}
	}
	// 遍历exData提取RequiredZipmodGUIDs
	exDataEx := make(map[string]*Base.PluginDataEx)
	for s, data := range exData {
		// 根据GUID找出插件对应的Data
		dex := data.DeserializeObjects()
		dex.Name = s
		dex.Version = data.Version
		exDataEx[dex.Name] = &dex
	}
	kc.ExtendedList = exDataEx
	kc.Image1 = pb.Png1
	kc.Image2 = pb.Png2
	return kc, nil
}

func (kc *KKCard) PrintCardInfo() {
	fmt.Println("Require Plugin:")
	for _, ex := range kc.ExtendedList {
		fmt.Printf("[Plugin]%s(Ver:%d)\n", ex.Name, ex.Version)
		ex.PrintMod()
	}
}
