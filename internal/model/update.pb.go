// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: update.proto

/*
	Package model is a generated protocol buffer package.

	It is generated from these files:
		update.proto

	It has these top-level messages:
		Update
*/
package model

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

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

type Update struct {
	ID      string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Unix    int64          `protobuf:"varint,2,opt,name=unix,proto3" json:"unix,omitempty"`
	Version int64          `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	Remove  bool           `protobuf:"varint,4,opt,name=remove,proto3" json:"remove,omitempty"`
	Set     []Update_Field `protobuf:"bytes,5,rep,name=set" json:"set"`
	Unset   []Update_Field `protobuf:"bytes,6,rep,name=unset" json:"unset"`
}

func (m *Update) Reset()                    { *m = Update{} }
func (m *Update) String() string            { return proto.CompactTextString(m) }
func (*Update) ProtoMessage()               {}
func (*Update) Descriptor() ([]byte, []int) { return fileDescriptorUpdate, []int{0} }

func (m *Update) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Update) GetUnix() int64 {
	if m != nil {
		return m.Unix
	}
	return 0
}

func (m *Update) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Update) GetRemove() bool {
	if m != nil {
		return m.Remove
	}
	return false
}

func (m *Update) GetSet() []Update_Field {
	if m != nil {
		return m.Set
	}
	return nil
}

func (m *Update) GetUnset() []Update_Field {
	if m != nil {
		return m.Unset
	}
	return nil
}

type Update_Field struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Update_Field) Reset()                    { *m = Update_Field{} }
func (m *Update_Field) String() string            { return proto.CompactTextString(m) }
func (*Update_Field) ProtoMessage()               {}
func (*Update_Field) Descriptor() ([]byte, []int) { return fileDescriptorUpdate, []int{0, 0} }

func (m *Update_Field) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Update_Field) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Update)(nil), "model.Update")
	proto.RegisterType((*Update_Field)(nil), "model.Update.Field")
}
func (m *Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ID) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintUpdate(dAtA, i, uint64(len(m.ID)))
		i += copy(dAtA[i:], m.ID)
	}
	if m.Unix != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintUpdate(dAtA, i, uint64(m.Unix))
	}
	if m.Version != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintUpdate(dAtA, i, uint64(m.Version))
	}
	if m.Remove {
		dAtA[i] = 0x20
		i++
		if m.Remove {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if len(m.Set) > 0 {
		for _, msg := range m.Set {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintUpdate(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Unset) > 0 {
		for _, msg := range m.Unset {
			dAtA[i] = 0x32
			i++
			i = encodeVarintUpdate(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Update_Field) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update_Field) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintUpdate(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if len(m.Data) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintUpdate(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	return i, nil
}

func encodeVarintUpdate(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Update) Size() (n int) {
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovUpdate(uint64(l))
	}
	if m.Unix != 0 {
		n += 1 + sovUpdate(uint64(m.Unix))
	}
	if m.Version != 0 {
		n += 1 + sovUpdate(uint64(m.Version))
	}
	if m.Remove {
		n += 2
	}
	if len(m.Set) > 0 {
		for _, e := range m.Set {
			l = e.Size()
			n += 1 + l + sovUpdate(uint64(l))
		}
	}
	if len(m.Unset) > 0 {
		for _, e := range m.Unset {
			l = e.Size()
			n += 1 + l + sovUpdate(uint64(l))
		}
	}
	return n
}

func (m *Update_Field) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovUpdate(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovUpdate(uint64(l))
	}
	return n
}

func sovUpdate(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozUpdate(x uint64) (n int) {
	return sovUpdate(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Update) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUpdate
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
			return fmt.Errorf("proto: Update: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Update: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
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
				return ErrInvalidLengthUpdate
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unix", wireType)
			}
			m.Unix = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Unix |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Remove", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Remove = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Set", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
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
				return ErrInvalidLengthUpdate
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Set = append(m.Set, Update_Field{})
			if err := m.Set[len(m.Set)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
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
				return ErrInvalidLengthUpdate
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Unset = append(m.Unset, Update_Field{})
			if err := m.Unset[len(m.Unset)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUpdate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthUpdate
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
func (m *Update_Field) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUpdate
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
			return fmt.Errorf("proto: Field: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Field: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
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
				return ErrInvalidLengthUpdate
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthUpdate
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUpdate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthUpdate
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
func skipUpdate(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUpdate
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
					return 0, ErrIntOverflowUpdate
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
					return 0, ErrIntOverflowUpdate
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
				return 0, ErrInvalidLengthUpdate
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowUpdate
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
				next, err := skipUpdate(dAtA[start:])
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
	ErrInvalidLengthUpdate = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUpdate   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("update.proto", fileDescriptorUpdate) }

var fileDescriptorUpdate = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x4f, 0x4a, 0xc4, 0x30,
	0x14, 0xc6, 0x27, 0xfd, 0x13, 0xf5, 0x39, 0x0b, 0x89, 0x30, 0x84, 0x59, 0x74, 0x8a, 0xab, 0x82,
	0x98, 0x82, 0xde, 0xa0, 0x88, 0xe0, 0x36, 0xe0, 0x01, 0x5a, 0x13, 0x6b, 0x60, 0xda, 0x0c, 0x9d,
	0x64, 0xf0, 0x18, 0x1e, 0x6b, 0x96, 0x9e, 0xa0, 0x48, 0x3d, 0x83, 0x7b, 0xe9, 0x8b, 0xae, 0xdd,
	0x7d, 0xdf, 0x97, 0x5f, 0x92, 0xf7, 0x3d, 0x58, 0xfa, 0x9d, 0xaa, 0x9d, 0x16, 0xbb, 0xc1, 0x3a,
	0xcb, 0xd2, 0xce, 0x2a, 0xbd, 0x5d, 0xdf, 0xb4, 0xc6, 0xbd, 0xfa, 0x46, 0x3c, 0xdb, 0xae, 0x6c,
	0x6d, 0x6b, 0x4b, 0x3c, 0x6d, 0xfc, 0x0b, 0x3a, 0x34, 0xa8, 0xc2, 0xad, 0xab, 0x6f, 0x02, 0xf4,
	0x09, 0x9f, 0x61, 0x2b, 0x88, 0x8c, 0xe2, 0x24, 0x27, 0xc5, 0x59, 0x45, 0xa7, 0x71, 0x13, 0x3d,
	0xde, 0xcb, 0xc8, 0x28, 0xc6, 0x20, 0xf1, 0xbd, 0x79, 0xe3, 0x51, 0x4e, 0x8a, 0x58, 0xa2, 0x66,
	0x1c, 0x4e, 0x0e, 0x7a, 0xd8, 0x1b, 0xdb, 0xf3, 0x18, 0xe3, 0x3f, 0xcb, 0x56, 0x40, 0x07, 0xdd,
	0xd9, 0x83, 0xe6, 0x49, 0x4e, 0x8a, 0x53, 0xf9, 0xeb, 0xd8, 0x35, 0xc4, 0x7b, 0xed, 0x78, 0x9a,
	0xc7, 0xc5, 0xf9, 0xed, 0xa5, 0xc0, 0x61, 0x45, 0xf8, 0x59, 0x3c, 0x18, 0xbd, 0x55, 0x55, 0x72,
	0x1c, 0x37, 0x0b, 0x39, 0x53, 0xac, 0x84, 0xd4, 0xf7, 0x33, 0x4e, 0xff, 0xc3, 0x03, 0xb7, 0x2e,
	0x21, 0xc5, 0x74, 0x1e, 0xb6, 0xaf, 0x3b, 0x1d, 0x6a, 0x48, 0xd4, 0x73, 0xa6, 0x6a, 0x57, 0x63,
	0x81, 0xa5, 0x44, 0x5d, 0x5d, 0x1c, 0xa7, 0x8c, 0x7c, 0x4c, 0x19, 0xf9, 0x9c, 0x32, 0xf2, 0xfe,
	0x95, 0x2d, 0x1a, 0x8a, 0x0b, 0xb9, 0xfb, 0x09, 0x00, 0x00, 0xff, 0xff, 0x31, 0x7a, 0x78, 0x75,
	0x56, 0x01, 0x00, 0x00,
}