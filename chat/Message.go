package chat

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

// TODO:
// Change username
