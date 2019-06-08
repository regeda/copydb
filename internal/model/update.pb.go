// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: update.proto

/*
	Package model is a generated protocol buffer package.

	It is generated from these files:
		update.proto

	It has these top-level messages:
		Update
		Item
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
	Items []Item `protobuf:"bytes,1,rep,name=items" json:"items"`
}

func (m *Update) Reset()                    { *m = Update{} }
func (m *Update) String() string            { return proto.CompactTextString(m) }
func (*Update) ProtoMessage()               {}
func (*Update) Descriptor() ([]byte, []int) { return fileDescriptorUpdate, []int{0} }

func (m *Update) GetItems() []Item {
	if m != nil {
		return m.Items
	}
	return nil
}

type Item struct {
	ID      string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Unix    int64        `protobuf:"zigzag64,2,opt,name=unix,proto3" json:"unix,omitempty"`
	Version int64        `protobuf:"zigzag64,3,opt,name=version,proto3" json:"version,omitempty"`
	Remove  bool         `protobuf:"varint,4,opt,name=remove,proto3" json:"remove,omitempty"`
	Set     []Item_Field `protobuf:"bytes,5,rep,name=set" json:"set"`
	Unset   []Item_Field `protobuf:"bytes,6,rep,name=unset" json:"unset"`
}

func (m *Item) Reset()                    { *m = Item{} }
func (m *Item) String() string            { return proto.CompactTextString(m) }
func (*Item) ProtoMessage()               {}
func (*Item) Descriptor() ([]byte, []int) { return fileDescriptorUpdate, []int{1} }

func (m *Item) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Item) GetUnix() int64 {
	if m != nil {
		return m.Unix
	}
	return 0
}

func (m *Item) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Item) GetRemove() bool {
	if m != nil {
		return m.Remove
	}
	return false
}

func (m *Item) GetSet() []Item_Field {
	if m != nil {
		return m.Set
	}
	return nil
}

func (m *Item) GetUnset() []Item_Field {
	if m != nil {
		return m.Unset
	}
	return nil
}

type Item_Field struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Item_Field) Reset()                    { *m = Item_Field{} }
func (m *Item_Field) String() string            { return proto.CompactTextString(m) }
func (*Item_Field) ProtoMessage()               {}
func (*Item_Field) Descriptor() ([]byte, []int) { return fileDescriptorUpdate, []int{1, 0} }

func (m *Item_Field) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Item_Field) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Update)(nil), "model.Update")
	proto.RegisterType((*Item)(nil), "model.Item")
	proto.RegisterType((*Item_Field)(nil), "model.Item.Field")
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
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			dAtA[i] = 0xa
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

func (m *Item) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Item) MarshalTo(dAtA []byte) (int, error) {
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
		i = encodeVarintUpdate(dAtA, i, uint64((uint64(m.Unix)<<1)^uint64((m.Unix>>63))))
	}
	if m.Version != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintUpdate(dAtA, i, uint64((uint64(m.Version)<<1)^uint64((m.Version>>63))))
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

func (m *Item_Field) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Item_Field) MarshalTo(dAtA []byte) (int, error) {
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
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovUpdate(uint64(l))
		}
	}
	return n
}

func (m *Item) Size() (n int) {
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovUpdate(uint64(l))
	}
	if m.Unix != 0 {
		n += 1 + sozUpdate(uint64(m.Unix))
	}
	if m.Version != 0 {
		n += 1 + sozUpdate(uint64(m.Version))
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

func (m *Item_Field) Size() (n int) {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
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
			m.Items = append(m.Items, Item{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *Item) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Item: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Item: illegal tag %d (wire type %d)", fieldNum, wire)
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
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = (v >> 1) ^ uint64((int64(v&1)<<63)>>63)
			m.Unix = int64(v)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUpdate
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			v = (v >> 1) ^ uint64((int64(v&1)<<63)>>63)
			m.Version = int64(v)
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
			m.Set = append(m.Set, Item_Field{})
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
			m.Unset = append(m.Unset, Item_Field{})
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
func (m *Item_Field) Unmarshal(dAtA []byte) error {
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
	// 284 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x4d, 0x4e, 0xc3, 0x30,
	0x10, 0x85, 0xeb, 0xfc, 0x01, 0xd3, 0x2e, 0xc0, 0x8b, 0xca, 0xea, 0x22, 0x8d, 0xba, 0x21, 0x2c,
	0x9a, 0x08, 0xb8, 0x41, 0x84, 0x90, 0xba, 0xb5, 0xc4, 0x01, 0x12, 0x6c, 0x82, 0xa5, 0x3a, 0xae,
	0x12, 0xa7, 0xe2, 0x18, 0x1c, 0xab, 0x4b, 0x4e, 0x50, 0xa1, 0x70, 0x00, 0xae, 0x80, 0x3c, 0x29,
	0x12, 0x2b, 0x76, 0xef, 0xbd, 0x79, 0x1e, 0x7d, 0x1e, 0x98, 0xf5, 0x3b, 0x51, 0x5a, 0x99, 0xed,
	0x5a, 0x63, 0x0d, 0x0d, 0xb5, 0x11, 0x72, 0xbb, 0x58, 0xd7, 0xca, 0xbe, 0xf6, 0x55, 0xf6, 0x6c,
	0x74, 0x5e, 0x9b, 0xda, 0xe4, 0x38, 0xad, 0xfa, 0x17, 0x74, 0x68, 0x50, 0x8d, 0xaf, 0x56, 0xb7,
	0x10, 0x3d, 0xe1, 0x16, 0x7a, 0x0d, 0xa1, 0xb2, 0x52, 0x77, 0x8c, 0x24, 0x7e, 0x3a, 0xbd, 0x9b,
	0x66, 0xb8, 0x2f, 0xdb, 0x58, 0xa9, 0x8b, 0xe0, 0x70, 0x5c, 0x4e, 0xf8, 0x38, 0x5f, 0x7d, 0x13,
	0x08, 0x5c, 0x4a, 0xe7, 0xe0, 0x29, 0xc1, 0x48, 0x42, 0xd2, 0x8b, 0x22, 0x1a, 0x8e, 0x4b, 0x6f,
	0xf3, 0xc0, 0x3d, 0x25, 0x28, 0x85, 0xa0, 0x6f, 0xd4, 0x1b, 0xf3, 0x12, 0x92, 0x52, 0x8e, 0x9a,
	0x32, 0x38, 0xdb, 0xcb, 0xb6, 0x53, 0xa6, 0x61, 0x3e, 0xc6, 0xbf, 0x96, 0xce, 0x21, 0x6a, 0xa5,
	0x36, 0x7b, 0xc9, 0x82, 0x84, 0xa4, 0xe7, 0xfc, 0xe4, 0xe8, 0x0d, 0xf8, 0x9d, 0xb4, 0x2c, 0x44,
	0x9a, 0xab, 0x3f, 0x34, 0xd9, 0xa3, 0x92, 0x5b, 0x71, 0x62, 0x72, 0x1d, 0xba, 0x86, 0xb0, 0x6f,
	0x5c, 0x39, 0xfa, 0xbf, 0x3c, 0xb6, 0x16, 0x39, 0x84, 0x98, 0x3a, 0xd0, 0xa6, 0xd4, 0x72, 0xfc,
	0x02, 0x47, 0xed, 0x32, 0x51, 0xda, 0x12, 0xe1, 0x67, 0x1c, 0x75, 0x71, 0x79, 0x18, 0x62, 0xf2,
	0x31, 0xc4, 0xe4, 0x73, 0x88, 0xc9, 0xfb, 0x57, 0x3c, 0xa9, 0x22, 0xbc, 0xde, 0xfd, 0x4f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xbb, 0xd7, 0x0c, 0x3a, 0x83, 0x01, 0x00, 0x00,
}
