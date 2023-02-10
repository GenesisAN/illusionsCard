package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/GenesisAN/illusionsCard/Base"
	"os"
	"strings"
)

type PngBuff struct {
	B        *bytes.Buffer
	FilePath string
	Png1     *[]byte
	Png2     *[]byte
	Type     string
}

func get_png(file []byte) int {
	res1 := bytes.Index(file, Base.PngEndChunk)
	return res1
}
func (pb *PngBuff) BuffRead(n int, errMsg string) ([]byte, error) {
	if pb.B.Len() < n {
		return nil, errors.New(errMsg)
	}
	return pb.B.Next(n), nil
}
func PngRead(path string) (*PngBuff, error) {
	var pb PngBuff
	path = strings.Replace(path, "\\", "/", -1)
	//提取文件名
	//读取图片
	f, err := os.ReadFile(path)
	if err != nil {
		return &pb, err
	}
	pb.B = bytes.NewBuffer(f)

	//value := get_png(f)
	//切割图片
	pngend := bytes.Index(f, Base.PngEndChunk) + len(Base.PngEndChunk)
	if pngend == -1 {
		return &pb, errors.New("PngRead fail:not found PngEndChunk")
	}
	png := pb.B.Next(pngend)
	pb.Png1 = &png
	//os.WriteFile("Out.png", outpng, 0776)
	fb, err := pb.B.ReadByte()
	if err != nil {
		return &pb, errors.New("PngRead fail:first byte not found")
	}
	if fb == 0x7 {
		_, err = pb.BuffRead(64, "PngRead fail:0x7 BuffRead fail")
		if err != nil {
			return &pb, err
		}
	} else if fb == 0x64 {
		_, err = pb.BuffRead(3, "PngRead fail:0x64 BuffRead fail")
		if err != nil {
			return &pb, err
		}
	}
	pb.Type, err = pb.StringRead()
	if err != nil {
		return &pb, errors.New("PngRead fail:card type string")
	}
	return &pb, err
}

func (pb *PngBuff) StringRead() (string, error) {
	types, err := pb.B.ReadByte()
	if err != nil {
		return "", errors.New("StringRead fail:unknown string len")
	}
	cardtypebyte, err := pb.BuffRead(int(types), "StringRead fail")
	if err != nil {
		return "", err
	}
	return string(cardtypebyte), nil
}
func (pb *PngBuff) UInt32Read() (uint32, error) {
	types, err := pb.BuffRead(4, "UInt32Read fail:unknown Int32 len")
	if err != nil {
		return 0, errors.New("UInt32Read fail:unknown Int32 len")
	}
	return binary.LittleEndian.Uint32(types), nil
}