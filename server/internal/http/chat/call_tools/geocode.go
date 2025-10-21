package calltools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (t *Tools) GetGeoCodeData(ctx context.Context, amenity string, street string, city string, state string, country string) ([]map[string]any, error) {
	baseURL := "https://nominatim.openstreetmap.org/search?"

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	query := u.Query()

	query.Set("format", "json")

	if amenity != "" {
		query.Set("amenity", amenity)
	}

	if street != "" {
		query.Set("street", street)
	}

	if city != "" {
		query.Set("city", city)
	}

	if state != "" {
		query.Set("state", state)
	}

	if country != "" {
		query.Set("country", country)
	}

	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "kaiyo-ai/1.0 (contact: nakulkrishnakumar86@gmail.com)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Full URL: %s\n", req.URL.String())
	fmt.Printf("Response: %s\n", body)

	var results []map[string]any
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no results")
	}

	var data []map[string]any

	if amenity == "" {
		data = append(data, results[0])
	} else {
		data = results
	}

	return data, nil
}
