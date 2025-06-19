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

// // KKCharaCard 代表Koikatsu的角色卡数据结构
type KKCharaCard struct {
	*Base.Card
	CharParmeter *KKChaFileParameter
}

func (k *KKCharaCard) GetPath() string {
	return k.Path
}

func (k *KKCharaCard) GetVersion() string {
	if k.CharParmeter != nil {
		return k.CharParmeter.Version
	}
	return k.LoadVersion
}

func (k *KKCharaCard) PrintCardInfo() {
	fmt.Println("Require Plugin:")
	for _, ex := range k.ExtendedList {
		fmt.Printf("[Plugin]%s(Ver:%d)\n", ex.Name, ex.Version)
		ex.PrintMod()
	}
}

type KKCharaReader struct{}

func (r KKCharaReader) Read(pgb *util.PngBuff) (Base.CardInterface, error) {
	if pgb.Type != Base.CT_KK {
		return nil, errors.New("KKReader: invalid type " + pgb.Type)
	}
	card, err := ParseKKChara(pgb)
	if err != nil {
		return nil, err
	}
	card.Path = pgb.FilePath
	return &card, nil
}

func (card *KKCharaCard) KKChaFileParameterEx(cfp *KKChaFileParameter) {
	card.CharParmeter = cfp
	card.CharInfo = &Base.ChaFileParameterEx{}
	card.CharInfo.Lastname = cfp.Lastname
	card.CharInfo.Firstname = cfp.Firstname
	card.CharInfo.Nickname = cfp.Nickname
	card.CharInfo.Version = cfp.Version
	card.CharInfo.Sex = cfp.Sex
}

func ParseKKChara(pb *util.PngBuff) (KKCharaCard, error) {
	kc := KKCharaCard{&Base.Card{CardType: pb.Type}, &KKChaFileParameter{}}
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
	exData, exDataEx := Base.ParsePluginData(extDataO)
	kc.Extended = exData
	kc.ExtendedList = exDataEx
	kc.Image1 = pb.Png1
	kc.Image2 = pb.Png2
	return kc, nil
}
