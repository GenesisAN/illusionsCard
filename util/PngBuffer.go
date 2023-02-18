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
func PngBytesRead(f []byte) (*PngBuff, error) {
	var pb PngBuff
	pb.B = bytes.NewBuffer(f)
	//value := get_png(f)
	//切割图片
	pngend := bytes.Index(f, Base.PngEndChunk) + len(Base.PngEndChunk)
	if pngend == -1 {
		return nil, errors.New("PngRead fail:not found PngEndChunk")
	}
	png := pb.B.Next(pngend)
	pb.Png1 = &png
	//os.WriteFile("Out.png", outpng, 0776)
	fb, err := pb.B.ReadByte()
	if err != nil {
		return nil, errors.New("PngRead fail:first byte not found")
	}
	if fb == 0x7 {
		_, err = pb.BuffRead(64, "PngRead fail:0x7 BuffRead fail")
		if err != nil {
			return nil, err
		}
	} else if fb == 0x64 {
		_, err = pb.BuffRead(3, "PngRead fail:0x64 BuffRead fail")
		if err != nil {
			return nil, err
		}
	}
	pb.Type, err = pb.StringRead()
	if err != nil {
		return &pb, errors.New("PngRead fail:card type string")
	}
	return &pb, err
}

func PngRead(path string) (*PngBuff, error) {
	path = strings.Replace(path, "\\", "/", -1)
	//提取文件名
	//读取图片
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return PngBytesRead(f)
}

func (pb *PngBuff) StringRead() (string, error) {
	// Buffer中读取string的长度但这并非C#标准实现,具体参考: https://github.com/dotnet/runtime/blob/8a46f0777ef975bb3d39cfb0b477c8e5c2d02b9a/src/libraries/System.Private.CoreLib/src/System/IO/BinaryReader.cs#L544
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
