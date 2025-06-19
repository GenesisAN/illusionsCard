// Package KKS Package KK 用于解析Koikatsu的角色卡数据
package KKS

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"

	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/util"
)

type SunshineCharaCard struct {
	*Base.Card
	CharParmeter *KKSChaFileParameter
}

func (c *SunshineCharaCard) GetPath() string {
	return c.Path
}

func (c *SunshineCharaCard) GetVersion() string {
	if c.CharParmeter != nil {
		return c.CharParmeter.Version
	}
	return c.LoadVersion
}

func (c *SunshineCharaCard) KKSChaFileParameterEx(cfp *KKSChaFileParameter) {
	c.CharParmeter = cfp
	c.CharInfo = &Base.ChaFileParameterEx{}
	c.CharInfo.Lastname = cfp.Lastname
	c.CharInfo.Firstname = cfp.Firstname
	c.CharInfo.Version = cfp.Version
	c.CharInfo.Nickname = cfp.Nickname
	c.CharInfo.Sex = cfp.Sex
}

type SunshineCharaReader struct{}

func (r SunshineCharaReader) Read(pgb *util.PngBuff) (Base.CardInterface, error) {
	if pgb.Type != Base.CT_KKS {
		return nil, errors.New("invalid card type for SunshineCharaReader: " + pgb.Type)
	}
	card, err := ParseKKSChara(pgb)
	if err != nil {
		return nil, err
	}
	card.Path = pgb.FilePath
	return &card, nil
}

func ParseKKSChara(pb *util.PngBuff) (SunshineCharaCard, error) {
	kc := SunshineCharaCard{&Base.Card{CardType: pb.Type}, &KKSChaFileParameter{}}
	Version, err := pb.StringRead()
	if err != nil {
		return kc, err
	}
	kc.LoadVersion = Version
	FLBuf, err := pb.BuffRead(4, "BuffRead Fail:Face img len")
	if err != nil {
		return kc, err
	}

	faceLength := binary.LittleEndian.Uint32(FLBuf)
	png, err := pb.BuffRead(int(faceLength), "BuffRead Fail:Face img")
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
		cBuf, err := pb.BuffRead(bh.Size, "BuffRead Fail:"+bh.Name)
		if err != nil {
			return kc, err
		}
		bh.Data = cBuf
		bhmap[bh.Name] = bh
	}
	//var parameterBytes, extDataByte []byte
	//遍历 头部信息，获取Parameter位置
	par, ok := bhmap["Parameter"]
	if ok {
		if par.Version != "0.0.5" {
			return kc, errors.New("BlockHeaderListInfo Unmarshal Fail")
		}
	} else {
		return kc, errors.New("parameter not found")
	}
	//根据位置信息，反序列化 MsgPack 得到 KKSChaFileParameter
	var Cfp KKSChaFileParameter
	_, err = Cfp.UnmarshalMsg(par.Data)
	if err != nil {
		return kc, errors.New("KKSChaFileParameter Unmarshal Fail")
	}
	kc.KKSChaFileParameterEx(&Cfp)
	kkex, kkexok := bhmap["KKEx"]
	//根据KKEx位置信息，反序列化 得到 extDataO
	var extDataO Base.MapSArrayInterface
	if kkexok {
		_, err := extDataO.UnmarshalMsg(kkex.Data)
		if err != nil {
			return kc, errors.New("extDataO Unmarshal Fail")
		}
	}
	exData, exDataEx := Base.ParsePluginData(extDataO)
	kc.Extended = exData
	kc.ExtendedList = exDataEx
	kc.Image1 = pb.Png1
	kc.Image2 = pb.Png2
	return kc, nil
}

func (c *SunshineCharaCard) PrintCardInfo() {
	fmt.Println("Require Plugin:")
	for _, ex := range c.ExtendedList {
		fmt.Printf("[Plugin]%s(Ver:%d)\n", ex.Name, ex.Version)
		ex.PrintMod()
	}
}
