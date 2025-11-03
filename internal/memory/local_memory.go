package memory

import "github.com/editor3-astera-editora/chat-rag/internal/model"

const MaxMessages = 10

var memoryStore = make(map[string][]model.Message)

func AddMessage(userID string, msg model.Message) {
	msgs := memoryStore[userID]

	if len(msgs) >= MaxMessages {
		msgs = msgs[1:]
	}
	memoryStore[userID] = append(msgs, msg)
}

func GetMemory(userID string) []model.Message {
	return memoryStore[userID]
}

func ClearMemory(userID string) {
	delete(memoryStore, userID)
}
