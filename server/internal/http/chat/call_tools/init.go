package calltools

import (
	"context"

	"github.com/openai/openai-go/v3"
)

type ToolBox interface {
	// Use this function to initialize all the mentioned openai tools
	InitOpenAITools() []openai.ChatCompletionToolUnionParam

	// Use this function to handle all tool functionalities
	HandleToolCall(ctx context.Context, funcName string, funcArgs string) []map[string]any

	// Use this function to get geocodes of a place
	GetGeoCodeData(ctx context.Context, amenity string, street string, city string, state string, country string) ([]map[string]any, error)
}

type Tools struct{}

func NewCallTools() *Tools {
	return &Tools{}
}

func (t *Tools) InitOpenAITools() []openai.ChatCompletionToolUnionParam {
	tools := []openai.ChatCompletionToolUnionParam{
		// // Current weather
		// {
		// 	OfFunction: &openai.ChatCompletionFunctionToolParam{
		// 		Function: openai.FunctionDefinitionParam{
		// 			Name:        "get_current_weather",
		// 			Description: openai.String("Get current weather by latitude and longitude."),
		// 			Parameters: map[string]any{
		// 				"type":       "object",
		// 				"required":   []string{"lat", "lon"},
		// 				"properties": map[string]any{
		// 					"lat": map[string]any{"type": "number"},
		// 					"lon": map[string]any{"type": "number"},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	OfFunction: &openai.ChatCompletionFunctionToolParam{
		// 		Function: openai.FunctionDefinitionParam{
		// 			Name:        "save_itinerary",
		// 			Description: openai.String("Call this ONLY when a finalized itinerary is ready."),
		// 			Parameters:  itinerarySchema,
		// 		},
		// 	},
		// },
		// Get Geocodes from place name
		{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: openai.FunctionDefinitionParam{
					Name:        "get_geocode_data",
					Description: openai.String("Convert multiple place names to latitude/longitude in a single batch call. Pass an array of location objects."),
					Parameters: map[string]any{
						"type":     "object",
						"required": []string{"locations"},
						"properties": map[string]any{
							"locations": map[string]any{
								"type":        "array",
								"description": "Array of location objects to geocode",
								"items": map[string]any{
									"type":     "object",
									"required": []string{"city", "country"},
									"properties": map[string]any{
										"amenity": map[string]any{"type": "string", "description": "Optional: specific venue or building"},
										"street":  map[string]any{"type": "string", "description": "Optional: street address"},
										"city":    map[string]any{"type": "string"},
										"state":   map[string]any{"type": "string", "description": "Optional: state/province"},
										"country": map[string]any{"type": "string"},
									},
								},
							},
						},
					},
				},
			},
		},
		// // Hotels search
		// {
		// 	OfFunction: &openai.ChatCompletionFunctionToolParam{
		// 		Function: openai.FunctionDefinitionParam{
		// 			Name:        "find_hotels",
		// 			Description: openai.String("Search hotels given location and dates."),
		// 			Parameters: map[string]any{
		// 				"type":     "object",
		// 				"required": []string{"location","checkIn","checkOut","guests"},
		// 				"properties": map[string]any{
		// 					"location": map[string]any{"type": "string"},
		// 					"checkIn":  map[string]any{"type": "string", "format": "date"},
		// 					"checkOut": map[string]any{"type": "string", "format": "date"},
		// 					"guests":   map[string]any{"type": "integer", "minimum": 1},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// // Flights search
		// {
		// 	OfFunction: &openai.ChatCompletionFunctionToolParam{
		// 		Function: openai.FunctionDefinitionParam{
		// 			Name:        "search_flights",
		// 			Description: openai.String("Search flights between origin and destination."),
		// 			Parameters: map[string]any{
		// 				"type":     "object",
		// 				"required": []string{"origin","destination","departDate","returnDate","adults"},
		// 				"properties": map[string]any{
		// 					"origin":      map[string]any{"type": "string"},
		// 					"destination": map[string]any{"type": "string"},
		// 					"departDate":  map[string]any{"type": "string", "format": "date"},
		// 					"returnDate":  map[string]any{"type": "string", "format": "date"},
		// 					"adults":      map[string]any{"type": "integer", "minimum": 1},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// // Popular places
		// {
		// 	OfFunction: &openai.ChatCompletionFunctionToolParam{
		// 		Function: openai.FunctionDefinitionParam{
		// 			Name:        "find_most_visited_places",
		// 			Description: openai.String("Return most visited attractions in a location."),
		// 			Parameters: map[string]any{
		// 				"type":     "object",
		// 				"required": []string{"location"},
		// 				"properties": map[string]any{
		// 					"location": map[string]any{"type": "string"},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// // Latest news
		// {
		// 	OfFunction: &openai.ChatCompletionFunctionToolParam{
		// 			Function: openai.FunctionDefinitionParam{
		// 			Name:        "get_latest_news",
		// 			Description: openai.String("Fetch latest news headlines matching a query."),
		// 			Parameters: map[string]any{
		// 				"type":     "object",
		// 				"required": []string{"query"},
		// 				"properties": map[string]any{
		// 					"query": map[string]any{"type": "string"},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}

	return tools
}
