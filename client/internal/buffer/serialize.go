package buffer

import (
	flatbuffers "github.com/google/flatbuffers/go"
	schema "xrossover-client/flatbuffers/xrossover"
)

func SerializeConnectionRequest(username string) []byte {
	builder := flatbuffers.NewBuilder(1024)

	user := builder.CreateString(username)

	schema.ConnectionRequestStart(builder)
	schema.ConnectionRequestAddUsername(builder, user)
	client := schema.ConnectionRequestEnd(builder)

	builder.Finish(client)
	return builder.FinishedBytes()
}
