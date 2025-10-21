package calltools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (t *Tools) HandleToolCall(ctx context.Context, funcName string, funcArgs string) []map[string]any {
	var result []map[string]any

	switch funcName {
	case "get_geocode_data":
		fmt.Println("GET_GEOCODE_DATA TOOL CALLED!")
		// Wrap in a struct with Locations field
		var payload struct {
			Locations []struct {
				Amenity string `json:"amenity"`
				Street  string `json:"street"`
				City    string `json:"city"`
				State   string `json:"state"`
				Country string `json:"country"`
			} `json:"locations"`
		}

		// Unmarshal into the wrapper
		if err := json.Unmarshal([]byte(funcArgs), &payload); err != nil {
			result = append(result, map[string]any{"error": err.Error()})
		} else {
			// Iterate payload.Locations
			for _, loc := range payload.Locations {
				geoData, err := t.GetGeoCodeData(
					ctx, loc.Amenity, loc.Street, loc.City, loc.State, loc.Country,
				)
				if err != nil {
					location := strings.Join(
						[]string{loc.Amenity, loc.Street, loc.City, loc.State, loc.Country}, " ",
					)
					result = append(result, map[string]any{
						"error":    err.Error(),
						"location": location,
					})
				} else {
					// Append each geoData entry
					result = append(result, geoData...)
				}
			}
		}
	}

	return result
}
