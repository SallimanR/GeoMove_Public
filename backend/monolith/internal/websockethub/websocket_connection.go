package websockethub

import "github.com/lesismal/nbio/nbhttp/websocket"

func SendDataWithRetry(conn *websocket.Conn, data []byte) error {
	err := conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		err := conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			return err
		}
	}

	return nil
}
