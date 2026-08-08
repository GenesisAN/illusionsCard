package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/GenesisAN/illusionsCard/Base"
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
	pngEndIndex := bytes.Index(f, Base.PngEndChunk)
	if pngEndIndex < 0 {
		return nil, errors.New("PngRead fail:not found PngEndChunk")
	}
	pngEnd := pngEndIndex + len(Base.PngEndChunk)
	png := f[:pngEnd]
	pb.Png1 = &png
	payload := f[pngEnd:]
	if len(payload) == 0 {
		return nil, errors.New("PngRead fail:card data not found")
	}

	pb.B = bytes.NewBuffer(payload)
	fb, err := pb.B.ReadByte()
	if err != nil {
		return nil, errors.New("PngRead fail:first byte not found")
	}
	if fb == 0x7 {
		pb.B = bytes.NewBuffer(payload)
		pb.Type = Base.CT_KKSC
		return &pb, nil
	} else if fb == 0x64 {
		_, err = pb.BuffRead(3, "PngRead fail:0x64 BuffRead fail")
		if err != nil {
			return nil, err
		}
	} else {
		pb.B = bytes.NewBuffer(payload)
		version, versionErr := pb.StringRead()
		if versionErr != nil || !looksLikeVersion(version) {
			return nil, errors.New("PngRead fail:unknown card prefix")
		}
		pb.B = bytes.NewBuffer(payload)
		pb.Type = Base.CT_KKSC
		return &pb, nil
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
	pb, err := PngBytesRead(f)
	if err != nil {
		return nil, err
	}
	pb.FilePath = path
	return pb, nil
}

func (pb *PngBuff) StringRead() (string, error) {
	length, err := read7BitEncodedInt(pb.B)
	if err != nil {
		return "", fmt.Errorf("StringRead fail: %w", err)
	}
	cardtypebyte, err := pb.BuffRead(length, "StringRead fail: truncated string")
	if err != nil {
		return "", err
	}
	return string(cardtypebyte), nil
}

// read7BitEncodedInt implements the length prefix used by .NET BinaryReader.
func read7BitEncodedInt(r io.ByteReader) (int, error) {
	var result uint32
	for shift := uint(0); shift < 35; shift += 7 {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if shift == 28 && b > 0x0f {
			return 0, errors.New("invalid 7-bit encoded integer")
		}
		result |= uint32(b&0x7f) << shift
		if b&0x80 == 0 {
			return int(result), nil
		}
	}
	return 0, errors.New("invalid 7-bit encoded integer")
}

func looksLikeVersion(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}
	return true
}

func (pb *PngBuff) UInt32Read() (uint32, error) {
	types, err := pb.BuffRead(4, "UInt32Read fail:unknown Int32 len")
	if err != nil {
		return 0, errors.New("UInt32Read fail:unknown Int32 len")
	}
	return binary.LittleEndian.Uint32(types), nil
}

func (pb *PngBuff) Int32Read() (int32, error) {
	types, err := pb.BuffRead(4, "Int32Read fail:unknown Int32 len")
	if err != nil {
		return 0, errors.New("Int32Read fail:unknown Int32 len")
	}
	return int32(binary.LittleEndian.Uint32(types)), nil
}
