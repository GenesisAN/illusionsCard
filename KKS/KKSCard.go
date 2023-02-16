// Package KKS Package KK 用于解析Koikatsu的角色卡数据
package KKS

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/GenesisAN/illusionsCard/Base"
	"github.com/GenesisAN/illusionsCard/util"
	"sort"
)

type KKSCard struct {
	*Base.Card
	CharParmeter *KKSChaFileParameter
}

func (c *KKSCard) KKChaFileParameterEx(cfp *KKSChaFileParameter) {
	c.CharParmeter = cfp
	c.CharInfo = &Base.ChaFileParameterEx{}
	c.CharInfo.Lastname = cfp.Lastname
	c.CharInfo.Firstname = cfp.Firstname
	c.CharInfo.Version = cfp.Version
	c.CharInfo.Nickname = cfp.Nickname
	c.CharInfo.Sex = cfp.Sex
}

func ParseKKSChara(pb *util.PngBuff) (KKSCard, error) {
	kc := KKSCard{&Base.Card{}, &KKSChaFileParameter{}}
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
	kc.KKChaFileParameterEx(&Cfp)
	kkex, kkexok := bhmap["KKEx"]
	//根据KKEx位置信息，反序列化 得到 extDataO
	var extDataO Base.MapSArrayInterface
	if kkexok {
		_, err := extDataO.UnmarshalMsg(kkex.Data)
		if err != nil {
			return kc, errors.New("extDataO Unmarshal Fail")
		}
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

func (c *KKSCard) PrintCardInfo() {
	fmt.Println("Require Plugin:")
	for _, ex := range c.ExtendedList {
		fmt.Printf("[Plugin]%s(Ver:%d)\n", ex.Name, ex.Version)
		ex.PrintMod()
	}
}
