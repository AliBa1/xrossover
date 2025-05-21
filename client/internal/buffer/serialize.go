package buffer

import (
	flatbuffers "github.com/google/flatbuffers/go"
	protocol "xrossover-client/flatbuffers/xrossover"
)

func SerializeConnectionRequest(username string) []byte {
	builder := flatbuffers.NewBuilder(1024)

	user := builder.CreateString(username)

	protocol.ConnectionRequestStart(builder)
	protocol.ConnectionRequestAddUsername(builder, user)
	connReq := protocol.ConnectionRequestEnd(builder)

	protocol.NetworkMessageStart(builder)
	// protocol.NetworkMessageAddType(builder, flatbuffers.UOffsetT(protocol.PayloadConnectionRequest))
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadConnectionRequest)
	protocol.NetworkMessageAddPayload(builder, connReq)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}
