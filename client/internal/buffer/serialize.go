package buffer

import (
	"net"
	protocol "xrossover-client/flatbuffers/xrossover"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SerializeConnectionRequest(username string, udpAddr *net.UDPAddr) []byte {
	builder := flatbuffers.NewBuilder(1024)

	user := builder.CreateString(username)
	udp := builder.CreateString(udpAddr.String())

	protocol.ConnectionRequestStart(builder)
	protocol.ConnectionRequestAddUsername(builder, user)
	protocol.ConnectionRequestAddUdpaddr(builder, udp)
	connReq := protocol.ConnectionRequestEnd(builder)

	protocol.NetworkMessageStart(builder)
	protocol.NetworkMessageAddPayloadType(builder, protocol.PayloadConnectionRequest)
	protocol.NetworkMessageAddPayload(builder, connReq)
	netMsg := protocol.NetworkMessageEnd(builder)

	builder.Finish(netMsg)

	return builder.FinishedBytes()
}
