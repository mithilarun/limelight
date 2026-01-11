package astro

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSunrise(t *testing.T) {
	testCases := []struct {
		name      string
		latitude  float64
		longitude float64
		date      time.Time
		wantHour  int
		wantMin   int
		tolerance int
	}{
		{
			name:      "san francisco winter solstice 2024",
			latitude:  37.7749,
			longitude: -122.4194,
			date:      time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  15,
			wantMin:   23,
			tolerance: 2,
		},
		{
			name:      "san francisco summer solstice 2024",
			latitude:  37.7749,
			longitude: -122.4194,
			date:      time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  12,
			wantMin:   47,
			tolerance: 2,
		},
		{
			name:      "new york winter solstice 2024",
			latitude:  40.7128,
			longitude: -74.0060,
			date:      time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  12,
			wantMin:   20,
			tolerance: 5,
		},
		{
			name:      "london summer solstice 2024",
			latitude:  51.5074,
			longitude: -0.1278,
			date:      time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  3,
			wantMin:   43,
			tolerance: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sunrise, err := CalculateSunrise(tc.latitude, tc.longitude, tc.date)
			require.NoError(t, err)

			utcSunrise := sunrise.UTC()
			wantTime := time.Date(tc.date.Year(), tc.date.Month(), tc.date.Day(), tc.wantHour, tc.wantMin, 0, 0, time.UTC)
			diff := utcSunrise.Sub(wantTime)
			diffMinutes := abs(int(diff.Minutes()))

			assert.LessOrEqual(t, diffMinutes, tc.tolerance, "time should match within %d minutes", tc.tolerance)
		})
	}
}

func TestCalculateSunset(t *testing.T) {
	testCases := []struct {
		name      string
		latitude  float64
		longitude float64
		date      time.Time
		wantHour  int
		wantMin   int
		tolerance int
	}{
		{
			name:      "san francisco winter solstice 2024",
			latitude:  37.7749,
			longitude: -122.4194,
			date:      time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  0,
			wantMin:   54,
			tolerance: 2,
		},
		{
			name:      "san francisco summer solstice 2024",
			latitude:  37.7749,
			longitude: -122.4194,
			date:      time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  3,
			wantMin:   35,
			tolerance: 2,
		},
		{
			name:      "new york winter solstice 2024",
			latitude:  40.7128,
			longitude: -74.0060,
			date:      time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  21,
			wantMin:   38,
			tolerance: 7,
		},
		{
			name:      "london summer solstice 2024",
			latitude:  51.5074,
			longitude: -0.1278,
			date:      time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
			wantHour:  20,
			wantMin:   21,
			tolerance: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sunset, err := CalculateSunset(tc.latitude, tc.longitude, tc.date)
			require.NoError(t, err)

			utcSunset := sunset.UTC()
			wantTime := time.Date(tc.date.Year(), tc.date.Month(), tc.date.Day(), tc.wantHour, tc.wantMin, 0, 0, time.UTC)
			diff := utcSunset.Sub(wantTime)
			diffMinutes := abs(int(diff.Minutes()))

			assert.LessOrEqual(t, diffMinutes, tc.tolerance, "time should match within %d minutes", tc.tolerance)
		})
	}
}

func TestCalculateSunriseInvalidLatitude(t *testing.T) {
	testCases := []struct {
		name     string
		latitude float64
	}{
		{
			name:     "latitude too high",
			latitude: 91.0,
		},
		{
			name:     "latitude too low",
			latitude: -91.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CalculateSunrise(tc.latitude, 0, time.Now())
			assert.Error(t, err)
		})
	}
}

func TestCalculateSunriseInvalidLongitude(t *testing.T) {
	testCases := []struct {
		name      string
		longitude float64
	}{
		{
			name:      "longitude too high",
			longitude: 181.0,
		},
		{
			name:      "longitude too low",
			longitude: -181.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CalculateSunrise(0, tc.longitude, time.Now())
			assert.Error(t, err)
		})
	}
}

func TestCalculateSunsetInvalidLatitude(t *testing.T) {
	testCases := []struct {
		name     string
		latitude float64
	}{
		{
			name:     "latitude too high",
			latitude: 91.0,
		},
		{
			name:     "latitude too low",
			latitude: -91.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CalculateSunset(tc.latitude, 0, time.Now())
			assert.Error(t, err)
		})
	}
}

func TestCalculateSunsetInvalidLongitude(t *testing.T) {
	testCases := []struct {
		name      string
		longitude float64
	}{
		{
			name:      "longitude too high",
			longitude: 181.0,
		},
		{
			name:      "longitude too low",
			longitude: -181.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CalculateSunset(0, tc.longitude, time.Now())
			assert.Error(t, err)
		})
	}
}

func TestPolarNight(t *testing.T) {
	_, err := CalculateSunrise(85.0, 0, time.Date(2024, 12, 21, 0, 0, 0, 0, time.UTC))
	assert.Error(t, err)
}

func TestPolarDay(t *testing.T) {
	_, err := CalculateSunset(85.0, 0, time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC))
	assert.Error(t, err)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
