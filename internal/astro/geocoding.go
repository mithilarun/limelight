package astro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	nominatimURL = "https://nominatim.openstreetmap.org/search"
	userAgent    = "limelight/1.0"
)

// GeocodingResult represents a location result from geocoding
type GeocodingResult struct {
	DisplayName string
	Latitude    float64
	Longitude   float64
}

// nominatimResponse represents the JSON response from Nominatim API
type nominatimResponse struct {
	PlaceID     int     `json:"place_id"`
	Lat         string  `json:"lat"`
	Lon         string  `json:"lon"`
	DisplayName string  `json:"display_name"`
	Type        string  `json:"type"`
	Importance  float64 `json:"importance"`
}

// GeocodeCity looks up the coordinates for a city using Nominatim (OpenStreetMap)
func GeocodeCity(cityName string) (*GeocodingResult, error) {
	if cityName == "" {
		return nil, errors.New("city name cannot be empty")
	}

	params := url.Values{}
	params.Set("q", cityName)
	params.Set("format", "json")
	params.Set("limit", "1")
	params.Set("addressdetails", "0")

	reqURL := fmt.Sprintf("%s?%s", nominatimURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make geocoding request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Newf("geocoding request failed with status %d", resp.StatusCode)
	}

	var results []nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, errors.Wrap(err, "failed to decode geocoding response")
	}

	if len(results) == 0 {
		return nil, errors.Newf("no results found for city: %s", cityName)
	}

	result := results[0]

	var lat, lon float64
	if _, err := fmt.Sscanf(result.Lat, "%f", &lat); err != nil {
		return nil, errors.Wrap(err, "failed to parse latitude")
	}
	if _, err := fmt.Sscanf(result.Lon, "%f", &lon); err != nil {
		return nil, errors.Wrap(err, "failed to parse longitude")
	}

	return &GeocodingResult{
		DisplayName: result.DisplayName,
		Latitude:    lat,
		Longitude:   lon,
	}, nil
}
