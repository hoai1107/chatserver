package chat

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

const (
	SendMessageType string = "Send"
	ChangeNameType  string = "ChangeName"
	NotifyType      string = "Notify"
)

func NewMessage(messageType string, username string, msg string) Message {
	return Message{
		Type:     messageType,
		Username: username,
		Content:  msg,
	}
}
