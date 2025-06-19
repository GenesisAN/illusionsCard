package KK

import (
	"errors"
	"fmt"

	"github.com/GenesisAN/illusionsCard/Base"
	util "github.com/GenesisAN/illusionsCard/util"
)

type KKClothesCard struct {
	*Base.Card
}

func (k *KKClothesCard) CompareMissingZipMods(localGUIDs []string) []string {
	missing := []string{}
	for _, guid := range localGUIDs {
		if _, ok := k.Extended[guid]; !ok {
			missing = append(missing, guid)
		}
	}
	return missing
}
func ParesKKClothes(pb *util.PngBuff) (*KKClothesCard, error) {
	kkc := KKClothesCard{&Base.Card{}}
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
	kkex, err := pb.BuffRead(int(len), "BuffRead Fail:KKEx")
	if err != nil {
		return &kkc, err
	}
	extDataO.UnmarshalMsg(kkex)
	exData, exDataEx := Base.ParsePluginData(extDataO)
	kkc.Extended = exData
	kkc.ExtendedList = exDataEx
	kkc.Extended = exData
	kkc.Image1 = pb.Png1
	kkc.Image2 = pb.Png2
	return &kkc, nil
}

type KKClothesReader struct{}

func (r KKClothesReader) Read(pgb *util.PngBuff) (Base.CardInterface, error) {
	// 类型判断略
	card, err := ParesKKClothes(pgb)
	if err != nil {
		return nil, err
	}
	card.Path = pgb.FilePath
	return card, nil
}
