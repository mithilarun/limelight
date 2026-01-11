package astro

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeocodeCity(t *testing.T) {
	// sleep for 2 seconds to respect Nominatim usage policy
	time.Sleep(time.Second * 2)

	result, err := GeocodeCity("San Francisco")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Greater(t, result.Latitude, 37.0)
	assert.Less(t, result.Latitude, 38.0)
	assert.Greater(t, result.Longitude, -123.0)
	assert.Less(t, result.Longitude, -122.0)
	assert.Contains(t, result.DisplayName, "San Francisco")
}

func TestGeocodeCityEmptyName(t *testing.T) {
	_, err := GeocodeCity("")
	assert.Error(t, err)
}

func TestGeocodeCityNotFound(t *testing.T) {
	// sleep for 2 seconds to respect Nominatim usage policy
	time.Sleep(time.Second * 2)

	_, err := GeocodeCity("ThisCityDoesNotExist123456")
	assert.Error(t, err)
}

func TestGeocodeCityMultipleWords(t *testing.T) {
	// sleep for 2 seconds to respect Nominatim usage policy
	time.Sleep(time.Second * 2)

	result, err := GeocodeCity("New York")
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Greater(t, result.Latitude, 40.0)
	assert.Less(t, result.Latitude, 41.0)
	assert.Contains(t, result.DisplayName, "New York")
}
