package calltools

import (
	"context"
	"encoding/json"
	"fmt"
)

func (t *Tools) HandleToolCall(ctx context.Context, funcName string, funcArgs string) []map[string]any {
	var result []map[string]any

	switch funcName {
	case "get_geocode_data":
		fmt.Println("GET_GEOCODE_DATA TOOL CALLED!")
		var args struct {
			Amenity string `json:"amenity"`
			Street  string `json:"street"`
			City    string `json:"city"`
			State   string `json:"state"`
			Country string `json:"country"`
		}
		if err := json.Unmarshal([]byte(funcArgs), &args); err != nil {
			result = append(result, map[string]any{"error": err.Error()})
		} else {
			geoData, err := t.GetGeoCodeData(ctx, args.Amenity, args.Street, args.City, args.State, args.Country)
			if err != nil {
				result = append(result, map[string]any{"error": err.Error()})
			} else {
				result = geoData
			}
		}
	}

	return result
}
