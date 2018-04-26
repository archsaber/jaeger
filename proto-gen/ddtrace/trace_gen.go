package ddtrace

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Trace) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0002) {
		(*z) = (*z)[:zb0002]
	} else {
		(*z) = make(Trace, zb0002)
	}
	for zb0001 := range *z {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				return
			}
			(*z)[zb0001] = nil
		} else {
			if (*z)[zb0001] == nil {
				(*z)[zb0001] = new(Span)
			}
			err = (*z)[zb0001].DecodeMsg(dc)
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Trace) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0003 := range z {
		if z[zb0003] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z[zb0003].EncodeMsg(en)
			if err != nil {
				return
			}
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Trace) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0003 := range z {
		if z[zb0003] == nil {
			s += msgp.NilSize
		} else {
			s += z[zb0003].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Traces) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0003 uint32
	zb0003, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zb0003) {
		(*z) = (*z)[:zb0003]
	} else {
		(*z) = make(Traces, zb0003)
	}
	for zb0001 := range *z {
		var zb0004 uint32
		zb0004, err = dc.ReadArrayHeader()
		if err != nil {
			return
		}
		if cap((*z)[zb0001]) >= int(zb0004) {
			(*z)[zb0001] = ((*z)[zb0001])[:zb0004]
		} else {
			(*z)[zb0001] = make(Trace, zb0004)
		}
		for zb0002 := range (*z)[zb0001] {
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				(*z)[zb0001][zb0002] = nil
			} else {
				if (*z)[zb0001][zb0002] == nil {
					(*z)[zb0001][zb0002] = new(Span)
				}
				err = (*z)[zb0001][zb0002].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Traces) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zb0005 := range z {
		err = en.WriteArrayHeader(uint32(len(z[zb0005])))
		if err != nil {
			return
		}
		for zb0006 := range z[zb0005] {
			if z[zb0005][zb0006] == nil {
				err = en.WriteNil()
				if err != nil {
					return
				}
			} else {
				err = z[zb0005][zb0006].EncodeMsg(en)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Traces) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zb0005 := range z {
		s += msgp.ArrayHeaderSize
		for zb0006 := range z[zb0005] {
			if z[zb0005][zb0006] == nil {
				s += msgp.NilSize
			} else {
				s += z[zb0005][zb0006].Msgsize()
			}
		}
	}
	return
}
