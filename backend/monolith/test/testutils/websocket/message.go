package testws

import (
	"testing"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	wsPB "monolith/internal/websockethub/proto"
)

func WriteMessage(t testing.TB, conn *websocket.Conn, data []byte) {
	err := conn.WriteMessage(websocket.BinaryMessage, data)
	require.NoErrorf(t, err, "failed to write msg pub")
}

func UnmarshalWSResponseMessage(t testing.TB, data []byte) *wsPB.ResponseMessage {
	var resp wsPB.Response
	err := proto.Unmarshal(data, &resp)
	require.NoErrorf(t, err, "failed to unmarshal response")
	respMsg := resp.Payload.(*wsPB.Response_Response).Response
	return respMsg
}
