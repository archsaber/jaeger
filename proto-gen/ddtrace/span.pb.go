// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: span.proto

/*
	Package ddtrace is a generated protocol buffer package.

	It is generated from these files:
		span.proto
		trace.proto

	It has these top-level messages:
		Span
		APITrace
*/
package ddtrace

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import binary "encoding/binary"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Span struct {
	Service  string             `protobuf:"bytes,1,opt,name=service,proto3" json:"service" msg:"service"`
	Name     string             `protobuf:"bytes,2,opt,name=name,proto3" json:"name" msg:"name"`
	Resource string             `protobuf:"bytes,3,opt,name=resource,proto3" json:"resource" msg:"resource"`
	TraceID  uint64             `protobuf:"varint,4,opt,name=traceID,proto3" json:"trace_id" msg:"trace_id"`
	SpanID   uint64             `protobuf:"varint,5,opt,name=spanID,proto3" json:"span_id" msg:"span_id"`
	ParentID uint64             `protobuf:"varint,6,opt,name=parentID,proto3" json:"parent_id" msg:"parent_id"`
	Start    int64              `protobuf:"varint,7,opt,name=start,proto3" json:"start" msg:"start"`
	Duration int64              `protobuf:"varint,8,opt,name=duration,proto3" json:"duration" msg:"duration"`
	Error    int32              `protobuf:"varint,9,opt,name=error,proto3" json:"error" msg:"error"`
	Meta     map[string]string  `protobuf:"bytes,10,rep,name=meta" json:"meta" msg:"meta" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Metrics  map[string]float64 `protobuf:"bytes,11,rep,name=metrics" json:"metrics" msg:"metrics" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	Type     string             `protobuf:"bytes,12,opt,name=type,proto3" json:"type" msg:"type"`
}

func (m *Span) Reset()                    { *m = Span{} }
func (m *Span) String() string            { return proto.CompactTextString(m) }
func (*Span) ProtoMessage()               {}
func (*Span) Descriptor() ([]byte, []int) { return fileDescriptorSpan, []int{0} }

func (m *Span) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Span) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Span) GetResource() string {
	if m != nil {
		return m.Resource
	}
	return ""
}

func (m *Span) GetTraceID() uint64 {
	if m != nil {
		return m.TraceID
	}
	return 0
}

func (m *Span) GetSpanID() uint64 {
	if m != nil {
		return m.SpanID
	}
	return 0
}

func (m *Span) GetParentID() uint64 {
	if m != nil {
		return m.ParentID
	}
	return 0
}

func (m *Span) GetStart() int64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *Span) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *Span) GetError() int32 {
	if m != nil {
		return m.Error
	}
	return 0
}

func (m *Span) GetMeta() map[string]string {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Span) GetMetrics() map[string]float64 {
	if m != nil {
		return m.Metrics
	}
	return nil
}

func (m *Span) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func init() {
	proto.RegisterType((*Span)(nil), "ddtrace.Span")
}
func (m *Span) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Span) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Service) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintSpan(dAtA, i, uint64(len(m.Service)))
		i += copy(dAtA[i:], m.Service)
	}
	if len(m.Name) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintSpan(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Resource) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintSpan(dAtA, i, uint64(len(m.Resource)))
		i += copy(dAtA[i:], m.Resource)
	}
	if m.TraceID != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintSpan(dAtA, i, uint64(m.TraceID))
	}
	if m.SpanID != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintSpan(dAtA, i, uint64(m.SpanID))
	}
	if m.ParentID != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintSpan(dAtA, i, uint64(m.ParentID))
	}
	if m.Start != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintSpan(dAtA, i, uint64(m.Start))
	}
	if m.Duration != 0 {
		dAtA[i] = 0x40
		i++
		i = encodeVarintSpan(dAtA, i, uint64(m.Duration))
	}
	if m.Error != 0 {
		dAtA[i] = 0x48
		i++
		i = encodeVarintSpan(dAtA, i, uint64(m.Error))
	}
	if len(m.Meta) > 0 {
		for k, _ := range m.Meta {
			dAtA[i] = 0x52
			i++
			v := m.Meta[k]
			mapSize := 1 + len(k) + sovSpan(uint64(len(k))) + 1 + len(v) + sovSpan(uint64(len(v)))
			i = encodeVarintSpan(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintSpan(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintSpan(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if len(m.Metrics) > 0 {
		for k, _ := range m.Metrics {
			dAtA[i] = 0x5a
			i++
			v := m.Metrics[k]
			mapSize := 1 + len(k) + sovSpan(uint64(len(k))) + 1 + 8
			i = encodeVarintSpan(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintSpan(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x11
			i++
			binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(v))))
			i += 8
		}
	}
	if len(m.Type) > 0 {
		dAtA[i] = 0x62
		i++
		i = encodeVarintSpan(dAtA, i, uint64(len(m.Type)))
		i += copy(dAtA[i:], m.Type)
	}
	return i, nil
}

func encodeVarintSpan(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Span) Size() (n int) {
	var l int
	_ = l
	l = len(m.Service)
	if l > 0 {
		n += 1 + l + sovSpan(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovSpan(uint64(l))
	}
	l = len(m.Resource)
	if l > 0 {
		n += 1 + l + sovSpan(uint64(l))
	}
	if m.TraceID != 0 {
		n += 1 + sovSpan(uint64(m.TraceID))
	}
	if m.SpanID != 0 {
		n += 1 + sovSpan(uint64(m.SpanID))
	}
	if m.ParentID != 0 {
		n += 1 + sovSpan(uint64(m.ParentID))
	}
	if m.Start != 0 {
		n += 1 + sovSpan(uint64(m.Start))
	}
	if m.Duration != 0 {
		n += 1 + sovSpan(uint64(m.Duration))
	}
	if m.Error != 0 {
		n += 1 + sovSpan(uint64(m.Error))
	}
	if len(m.Meta) > 0 {
		for k, v := range m.Meta {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovSpan(uint64(len(k))) + 1 + len(v) + sovSpan(uint64(len(v)))
			n += mapEntrySize + 1 + sovSpan(uint64(mapEntrySize))
		}
	}
	if len(m.Metrics) > 0 {
		for k, v := range m.Metrics {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovSpan(uint64(len(k))) + 1 + 8
			n += mapEntrySize + 1 + sovSpan(uint64(mapEntrySize))
		}
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovSpan(uint64(l))
	}
	return n
}

func sovSpan(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSpan(x uint64) (n int) {
	return sovSpan(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Span) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSpan
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Span: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Span: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Service", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpan
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Service = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpan
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Resource", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpan
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Resource = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TraceID", wireType)
			}
			m.TraceID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TraceID |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpanID", wireType)
			}
			m.SpanID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SpanID |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParentID", wireType)
			}
			m.ParentID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ParentID |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Start", wireType)
			}
			m.Start = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Start |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			m.Duration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Duration |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			m.Error = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Error |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSpan
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Meta == nil {
				m.Meta = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSpan
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSpan
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthSpan
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSpan
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthSpan
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipSpan(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthSpan
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Meta[mapkey] = mapvalue
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metrics", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSpan
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Metrics == nil {
				m.Metrics = make(map[string]float64)
			}
			var mapkey string
			var mapvalue float64
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowSpan
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowSpan
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthSpan
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapvaluetemp uint64
					if (iNdEx + 8) > l {
						return io.ErrUnexpectedEOF
					}
					mapvaluetemp = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
					iNdEx += 8
					mapvalue = math.Float64frombits(mapvaluetemp)
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipSpan(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthSpan
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Metrics[mapkey] = mapvalue
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSpan
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSpan(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSpan
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSpan(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSpan
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSpan
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthSpan
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSpan
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipSpan(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthSpan = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSpan   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("span.proto", fileDescriptorSpan) }

var fileDescriptorSpan = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xcf, 0x8e, 0xd3, 0x30,
	0x10, 0xc6, 0xf1, 0x36, 0x6d, 0x5a, 0x77, 0x81, 0x95, 0x85, 0xc0, 0xaa, 0x50, 0x12, 0xf9, 0x14,
	0x21, 0x91, 0x95, 0x00, 0xc1, 0xaa, 0xe2, 0x54, 0xca, 0xa1, 0x07, 0x2e, 0xe6, 0x01, 0x90, 0x9b,
	0x9a, 0x12, 0x41, 0xfe, 0xc8, 0x71, 0x56, 0xea, 0x5b, 0xf0, 0x1c, 0x3c, 0x09, 0x47, 0x9e, 0x20,
	0x42, 0xe5, 0x96, 0x63, 0x9f, 0x00, 0x79, 0x9c, 0x98, 0x4a, 0x1c, 0xf6, 0x96, 0xef, 0x37, 0xf3,
	0x79, 0x3c, 0xe3, 0x09, 0xc6, 0x75, 0x25, 0x8a, 0xa4, 0x52, 0xa5, 0x2e, 0x89, 0xbf, 0xdb, 0x69,
	0x25, 0x52, 0xb9, 0x78, 0xbe, 0xcf, 0xf4, 0x97, 0x66, 0x9b, 0xa4, 0x65, 0x7e, 0xbd, 0x2f, 0xf7,
	0xe5, 0x35, 0xc4, 0xb7, 0xcd, 0x67, 0x50, 0x20, 0xe0, 0xcb, 0xfa, 0xd8, 0x8f, 0x09, 0xf6, 0x3e,
	0x56, 0xa2, 0x20, 0xaf, 0xb1, 0x5f, 0x4b, 0x75, 0x9b, 0xa5, 0x92, 0xa2, 0x08, 0xc5, 0xb3, 0xd5,
	0xd3, 0xae, 0x0d, 0x07, 0x74, 0x6a, 0xc3, 0xfb, 0x79, 0xbd, 0x5f, 0xb2, 0x5e, 0x33, 0x3e, 0x44,
	0xc8, 0x33, 0xec, 0x15, 0x22, 0x97, 0xf4, 0x02, 0x4c, 0x8f, 0xbb, 0x36, 0x04, 0x7d, 0x6a, 0x43,
	0x0c, 0x0e, 0x23, 0x18, 0x07, 0x46, 0x96, 0x78, 0xaa, 0x64, 0x5d, 0x36, 0x2a, 0x95, 0x74, 0x04,
	0xf9, 0x41, 0xd7, 0x86, 0x8e, 0x9d, 0xda, 0xf0, 0x01, 0x78, 0x06, 0xc0, 0xb8, 0x8b, 0x91, 0x1b,
	0xec, 0x43, 0x83, 0x9b, 0x35, 0xf5, 0x22, 0x14, 0x7b, 0xd6, 0x0a, 0xe8, 0x53, 0xb6, 0x73, 0xd6,
	0x01, 0x30, 0x3e, 0xa4, 0x93, 0x57, 0x78, 0x62, 0x06, 0xb5, 0x59, 0xd3, 0x31, 0x18, 0x6d, 0x63,
	0x95, 0x28, 0xac, 0xaf, 0x6f, 0xcc, 0x6a, 0xc6, 0xfb, 0x5c, 0xf2, 0x16, 0x4f, 0x2b, 0xa1, 0x64,
	0xa1, 0x37, 0x6b, 0x3a, 0x01, 0x5f, 0xd4, 0xb5, 0xe1, 0xcc, 0x32, 0xeb, 0x7c, 0x08, 0x4e, 0x47,
	0x18, 0x77, 0x0e, 0x92, 0xe0, 0x71, 0xad, 0x85, 0xd2, 0xd4, 0x8f, 0x50, 0x3c, 0x5a, 0xd1, 0xae,
	0x0d, 0x2d, 0x38, 0xb5, 0xe1, 0xdc, 0x16, 0x34, 0x8a, 0x71, 0x4b, 0xcd, 0x64, 0x76, 0x8d, 0x12,
	0x3a, 0x2b, 0x0b, 0x3a, 0x05, 0x0b, 0xb4, 0x37, 0x30, 0xd7, 0xde, 0x00, 0x18, 0x77, 0x31, 0x53,
	0x4b, 0x2a, 0x55, 0x2a, 0x3a, 0x8b, 0x50, 0x3c, 0xb6, 0xb5, 0x00, 0xb8, 0x5a, 0xa0, 0x18, 0xb7,
	0x94, 0xbc, 0xc3, 0x5e, 0x2e, 0xb5, 0xa0, 0x38, 0x1a, 0xc5, 0xf3, 0x17, 0x4f, 0x92, 0x7e, 0x73,
	0x12, 0xb3, 0x06, 0xc9, 0x07, 0xa9, 0xc5, 0xfb, 0x42, 0xab, 0x83, 0x7d, 0x4a, 0x93, 0xe8, 0x9e,
	0xd2, 0x08, 0xc6, 0x81, 0x11, 0x8e, 0xfd, 0x5c, 0x6a, 0x95, 0xa5, 0x35, 0x9d, 0xc3, 0x39, 0x8b,
	0xff, 0xce, 0x31, 0x41, 0x7b, 0x14, 0x4c, 0xbc, 0x4f, 0x77, 0x13, 0xef, 0x35, 0xe3, 0x43, 0xc4,
	0xac, 0x92, 0x3e, 0x54, 0x92, 0x5e, 0xfe, 0x5b, 0x25, 0xa3, 0x5d, 0x7d, 0x23, 0x18, 0x07, 0xb6,
	0x78, 0x83, 0x67, 0xee, 0xaa, 0xe4, 0x0a, 0x8f, 0xbe, 0xca, 0x83, 0xdd, 0x5b, 0x6e, 0x3e, 0xc9,
	0x23, 0x3c, 0xbe, 0x15, 0xdf, 0x9a, 0x7e, 0x2d, 0xb9, 0x15, 0xcb, 0x8b, 0x1b, 0xb4, 0x58, 0xe2,
	0xcb, 0xf3, 0xbb, 0xdd, 0xe5, 0x45, 0x67, 0xde, 0xd5, 0xd5, 0xcf, 0x63, 0x80, 0x7e, 0x1d, 0x03,
	0xf4, 0xfb, 0x18, 0xa0, 0xef, 0x7f, 0x82, 0x7b, 0xdb, 0x09, 0xfc, 0x45, 0x2f, 0xff, 0x06, 0x00,
	0x00, 0xff, 0xff, 0xa5, 0x7f, 0x70, 0xa3, 0x8b, 0x03, 0x00, 0x00,
}