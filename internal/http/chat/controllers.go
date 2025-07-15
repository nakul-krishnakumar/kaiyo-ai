package chat

type Message struct {
	Role string
	Content string
}

type Controller struct {
	History []Message
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) SendMessage(chatID, userID, content string) (string, error) {
	reply := "Hello from kairo bot!"

	newMessage := []Message{
		{ Role: "user", Content: content }, 
		{ Role: "bot", Content: reply },
	}

	c.History = append(c.History, newMessage...)
	return reply, nil
}

func (c *Controller) GetHistory(chatID string) ([]Message, error) {
	return c.History, nil
}