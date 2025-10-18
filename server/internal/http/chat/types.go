package chat

import (
	"time"

	"github.com/openai/openai-go/v3"
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

type Itinerary struct {
	Destination string    `json:"destination"`         // e.g., "Coorg"
	StartDate   string    `json:"startDate,omitempty"` // ISO date string, e.g., "2025-11-01"
	EndDate     string    `json:"endDate,omitempty"`   // ISO date string, e.g., "2025-11-03"
	Currency    string    `json:"currency,omitempty"`  // e.g., "INR", "USD"
	Days        []DayPlan `json:"days"`                // Day-wise plan
}

type DayPlan struct {
	Day   int       `json:"day"`             // 1-based day index
	Label string    `json:"label,omitempty"` // Optional label, e.g., "Arrival", "Trek day"
	Items []DayItem `json:"items"`           // Stops/activities for the day
}

type DayItem struct {
	Title     string  `json:"title"`               // Activity/title, e.g., "Tadiandamol Trek"
	City      string  `json:"city,omitempty"`      // City/town context, e.g., "Madikeri"
	Place     string  `json:"place,omitempty"`     // POI name, e.g., "Abbey Falls"
	Category  string  `json:"category,omitempty"`  // e.g., "sightseeing", "food", "trek"
	StartTime string  `json:"startTime,omitempty"` // "09:00" (24h) or ISO time
	EndTime   string  `json:"endTime,omitempty"`   // "11:30"
	Notes     string  `json:"notes,omitempty"`     // Free-form notes
	Lat       float64 `json:"lat,omitempty"`       // Geocoded latitude
	Lon       float64 `json:"lon,omitempty"`       // Geocoded longitude
}

type Controller struct {
	Client *openai.Client // OpenAI Client
	Model
	History   []openai.ChatCompletionMessageParamUnion // context memory to store messages
	Itinerary *Itinerary
}

type Handler struct {
	Controller *Controller
}

type UserInput struct {
	ChatID  string
	UserID  string
	Content string
}
