// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package xrossover

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Move struct {
	_tab flatbuffers.Table
}

func GetRootAsMove(buf []byte, offset flatbuffers.UOffsetT) *Move {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Move{}
	x.Init(buf, n+offset)
	return x
}

func FinishMoveBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMove(buf []byte, offset flatbuffers.UOffsetT) *Move {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Move{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMoveBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Move) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Move) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Move) Direction(obj *Vector3) *Vector3 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
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

func MoveStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func MoveAddDirection(builder *flatbuffers.Builder, direction flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(direction), 0)
}
func MoveEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
