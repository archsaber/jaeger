package ddtrace

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ServicesMetadata) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0005 uint32
	zb0005, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	if (*z) == nil && zb0005 > 0 {
		(*z) = make(ServicesMetadata, zb0005)
	} else if len((*z)) > 0 {
		for key, _ := range *z {
			delete((*z), key)
		}
	}
	for zb0005 > 0 {
		zb0005--
		var zb0001 string
		var zb0002 map[string]string
		zb0001, err = dc.ReadString()
		if err != nil {
			return
		}
		var zb0006 uint32
		zb0006, err = dc.ReadMapHeader()
		if err != nil {
			return
		}
		if zb0002 == nil && zb0006 > 0 {
			zb0002 = make(map[string]string, zb0006)
		} else if len(zb0002) > 0 {
			for key, _ := range zb0002 {
				delete(zb0002, key)
			}
		}
		for zb0006 > 0 {
			zb0006--
			var zb0003 string
			var zb0004 string
			zb0003, err = dc.ReadString()
			if err != nil {
				return
			}
			zb0004, err = dc.ReadString()
			if err != nil {
				return
			}
			zb0002[zb0003] = zb0004
		}
		(*z)[zb0001] = zb0002
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ServicesMetadata) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0007, zb0008 := range z {
		err = en.WriteString(zb0007)
		if err != nil {
			return
		}
		err = en.WriteMapHeader(uint32(len(zb0008)))
		if err != nil {
			return
		}
		for zb0009, zb0010 := range zb0008 {
			err = en.WriteString(zb0009)
			if err != nil {
				return
			}
			err = en.WriteString(zb0010)
			if err != nil {
				return
			}
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ServicesMetadata) Msgsize() (s int) {
	s = msgp.MapHeaderSize
	if z != nil {
		for zb0007, zb0008 := range z {
			_ = zb0008
			s += msgp.StringPrefixSize + len(zb0007) + msgp.MapHeaderSize
			if zb0008 != nil {
				for zb0009, zb0010 := range zb0008 {
					_ = zb0010
					s += msgp.StringPrefixSize + len(zb0009) + msgp.StringPrefixSize + len(zb0010)
				}
			}
		}
	}
	return
}
