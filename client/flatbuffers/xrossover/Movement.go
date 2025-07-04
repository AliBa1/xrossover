// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package xrossover

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Movement struct {
	_tab flatbuffers.Table
}

func GetRootAsMovement(buf []byte, offset flatbuffers.UOffsetT) *Movement {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Movement{}
	x.Init(buf, n+offset)
	return x
}

func FinishMovementBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMovement(buf []byte, offset flatbuffers.UOffsetT) *Movement {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Movement{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMovementBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Movement) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Movement) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Movement) ObjectId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Movement) ObjectOwner() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Movement) Direction(obj *Vector3) *Vector3 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(Vector3)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func MovementStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func MovementAddObjectId(builder *flatbuffers.Builder, objectId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(objectId), 0)
}
func MovementAddObjectOwner(builder *flatbuffers.Builder, objectOwner flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(objectOwner), 0)
}
func MovementAddDirection(builder *flatbuffers.Builder, direction flatbuffers.UOffsetT) {
	builder.PrependStructSlot(2, flatbuffers.UOffsetT(direction), 0)
}
func MovementEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
