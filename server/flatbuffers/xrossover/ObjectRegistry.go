// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package xrossover

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ObjectRegistry struct {
	_tab flatbuffers.Table
}

func GetRootAsObjectRegistry(buf []byte, offset flatbuffers.UOffsetT) *ObjectRegistry {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ObjectRegistry{}
	x.Init(buf, n+offset)
	return x
}

func FinishObjectRegistryBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsObjectRegistry(buf []byte, offset flatbuffers.UOffsetT) *ObjectRegistry {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ObjectRegistry{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedObjectRegistryBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ObjectRegistry) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ObjectRegistry) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ObjectRegistry) Objects(obj *GameObjectWrapper, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *ObjectRegistry) ObjectsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ObjectRegistryStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ObjectRegistryAddObjects(builder *flatbuffers.Builder, objects flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(objects), 0)
}
func ObjectRegistryStartObjectsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ObjectRegistryEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
