package KK

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *BlockHeader) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "version":
			z.Version, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "pos":
			z.Pos, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Pos")
				return
			}
		case "size":
			z.Size, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Size")
				return
			}
		case "Data":
			z.Data, err = dc.ReadBytes(z.Data)
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *BlockHeader) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "name"
	err = en.Append(0x85, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	// write "version"
	err = en.Append(0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Version)
	if err != nil {
		err = msgp.WrapError(err, "Version")
		return
	}
	// write "pos"
	err = en.Append(0xa3, 0x70, 0x6f, 0x73)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Pos)
	if err != nil {
		err = msgp.WrapError(err, "Pos")
		return
	}
	// write "size"
	err = en.Append(0xa4, 0x73, 0x69, 0x7a, 0x65)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Size)
	if err != nil {
		err = msgp.WrapError(err, "Size")
		return
	}
	// write "Data"
	err = en.Append(0xa4, 0x44, 0x61, 0x74, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.Data)
	if err != nil {
		err = msgp.WrapError(err, "Data")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *BlockHeader) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "name"
	o = append(o, 0x85, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "version"
	o = append(o, 0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Version)
	// string "pos"
	o = append(o, 0xa3, 0x70, 0x6f, 0x73)
	o = msgp.AppendInt(o, z.Pos)
	// string "size"
	o = append(o, 0xa4, 0x73, 0x69, 0x7a, 0x65)
	o = msgp.AppendInt(o, z.Size)
	// string "Data"
	o = append(o, 0xa4, 0x44, 0x61, 0x74, 0x61)
	o = msgp.AppendBytes(o, z.Data)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlockHeader) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "version":
			z.Version, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "pos":
			z.Pos, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Pos")
				return
			}
		case "size":
			z.Size, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Size")
				return
			}
		case "Data":
			z.Data, bts, err = msgp.ReadBytesBytes(bts, z.Data)
			if err != nil {
				err = msgp.WrapError(err, "Data")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *BlockHeader) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 8 + msgp.StringPrefixSize + len(z.Version) + 4 + msgp.IntSize + 5 + msgp.IntSize + 5 + msgp.BytesPrefixSize + len(z.Data)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *BlockHeaderListInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "lstInfo":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "LstInfo")
				return
			}
			if cap(z.LstInfo) >= int(zb0002) {
				z.LstInfo = (z.LstInfo)[:zb0002]
			} else {
				z.LstInfo = make([]*BlockHeader, zb0002)
			}
			for za0001 := range z.LstInfo {
				if dc.IsNil() {
					err = dc.ReadNil()
					if err != nil {
						err = msgp.WrapError(err, "LstInfo", za0001)
						return
					}
					z.LstInfo[za0001] = nil
				} else {
					if z.LstInfo[za0001] == nil {
						z.LstInfo[za0001] = new(BlockHeader)
					}
					err = z.LstInfo[za0001].DecodeMsg(dc)
					if err != nil {
						err = msgp.WrapError(err, "LstInfo", za0001)
						return
					}
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *BlockHeaderListInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "lstInfo"
	err = en.Append(0x81, 0xa7, 0x6c, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.LstInfo)))
	if err != nil {
		err = msgp.WrapError(err, "LstInfo")
		return
	}
	for za0001 := range z.LstInfo {
		if z.LstInfo[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.LstInfo[za0001].EncodeMsg(en)
			if err != nil {
				err = msgp.WrapError(err, "LstInfo", za0001)
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *BlockHeaderListInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "lstInfo"
	o = append(o, 0x81, 0xa7, 0x6c, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f)
	o = msgp.AppendArrayHeader(o, uint32(len(z.LstInfo)))
	for za0001 := range z.LstInfo {
		if z.LstInfo[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.LstInfo[za0001].MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "LstInfo", za0001)
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *BlockHeaderListInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "lstInfo":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "LstInfo")
				return
			}
			if cap(z.LstInfo) >= int(zb0002) {
				z.LstInfo = (z.LstInfo)[:zb0002]
			} else {
				z.LstInfo = make([]*BlockHeader, zb0002)
			}
			for za0001 := range z.LstInfo {
				if msgp.IsNil(bts) {
					bts, err = msgp.ReadNilBytes(bts)
					if err != nil {
						return
					}
					z.LstInfo[za0001] = nil
				} else {
					if z.LstInfo[za0001] == nil {
						z.LstInfo[za0001] = new(BlockHeader)
					}
					bts, err = z.LstInfo[za0001].UnmarshalMsg(bts)
					if err != nil {
						err = msgp.WrapError(err, "LstInfo", za0001)
						return
					}
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *BlockHeaderListInfo) Msgsize() (s int) {
	s = 1 + 8 + msgp.ArrayHeaderSize
	for za0001 := range z.LstInfo {
		if z.LstInfo[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.LstInfo[za0001].Msgsize()
		}
	}
	return
}
