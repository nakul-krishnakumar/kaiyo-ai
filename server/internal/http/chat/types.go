package chat

import (
	"time"

	"github.com/openai/openai-go"
)

type Model struct {
	Name         string `mapstructure:"model_name"`
	Type         string `mapstructure:"model_type"`
	SystemPrompt string `mapstructure:"system_prompt"`
}

type Message struct {
	Role      string
	Content   string
	CreatedAt time.Time
}

type Controller struct {
	Client *openai.Client // OpenAI Client
	Model
	History []openai.ChatCompletionMessageParamUnion // context memory to store messages
}

type Handler struct {
	Controller *Controller
}

type UserInput struct {
	ChatID  string
	UserID  string
	Content string
}
