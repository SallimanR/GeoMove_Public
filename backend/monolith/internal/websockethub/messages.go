package websockethub

import (
	"google.golang.org/protobuf/proto"

	wsPB "monolith/internal/websockethub/proto"
)

// func MarshalProtobufBatch[T proto.Message](data []byte, newMsg func() T) ([]byte, error) {
// }

func UnmarshalProtobufBatch[T proto.Message](messages []byte, newMsg func() T) ([]T, error) {
	var batch wsPB.MessageBatch
	err := proto.Unmarshal(messages, &batch)
	if err != nil {
		return nil, err
	}

	result := make([]T, 0, len(batch.Data))
	for _, msgBytes := range batch.Data {
		msg := newMsg()
		err := proto.Unmarshal(msgBytes, msg)
		if err != nil {
			return nil, err
		}

		result = append(result, msg)
	}

	return result, nil
}
