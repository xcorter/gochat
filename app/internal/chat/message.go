package chat

type Message struct {
	ThreadId int    `json:"threadId"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func NewMessage(threadId int, username string, message string) Message {
	return Message{
		ThreadId: threadId,
		Username: username,
		Message:  message,
	}
}
