// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package xrossover

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type NetworkMessage struct {
	_tab flatbuffers.Table
}

func GetRootAsNetworkMessage(buf []byte, offset flatbuffers.UOffsetT) *NetworkMessage {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &NetworkMessage{}
	x.Init(buf, n+offset)
	return x
}

func FinishNetworkMessageBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsNetworkMessage(buf []byte, offset flatbuffers.UOffsetT) *NetworkMessage {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &NetworkMessage{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedNetworkMessageBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *NetworkMessage) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *NetworkMessage) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *NetworkMessage) PayloadType() Payload {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return Payload(rcv._tab.GetByte(o + rcv._tab.Pos))
	}
	return 0
}

func (rcv *NetworkMessage) MutatePayloadType(n Payload) bool {
	return rcv._tab.MutateByteSlot(4, byte(n))
}

func (rcv *NetworkMessage) Payload(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func NetworkMessageStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func NetworkMessageAddPayloadType(builder *flatbuffers.Builder, payloadType Payload) {
	builder.PrependByteSlot(0, byte(payloadType), 0)
}
func NetworkMessageAddPayload(builder *flatbuffers.Builder, payload flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(payload), 0)
}
func NetworkMessageEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
