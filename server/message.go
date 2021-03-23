package server

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
}

func (message Message) String() string {
	ret := ""
	if message.Sender != "" {
		ret += message.Sender + " 说: "
	}
	ret += message.Content
	return ret
}
