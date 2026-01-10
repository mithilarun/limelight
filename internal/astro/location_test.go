package astro

import (
	"os"
	"testing"

	"github.com/mithilarun/limelight/internal/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnv(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})
}

func TestGetLocationFromConfig(t *testing.T) {
	setupTestEnv(t)

	config := &credentials.Config{
		BridgeIP:  "192.168.1.100",
		Latitude:  37.7749,
		Longitude: -122.4194,
	}
	err := credentials.SaveConfig(config)
	require.NoError(t, err)

	lat, lon, err := GetLocationFromConfig()
	require.NoError(t, err)
	assert.Equal(t, 37.7749, lat)
	assert.Equal(t, -122.4194, lon)
}

func TestGetLocationFromConfigNoFile(t *testing.T) {
	setupTestEnv(t)

	_, _, err := GetLocationFromConfig()
	assert.Error(t, err)
}

func TestGetLocationFromConfigNoLocation(t *testing.T) {
	setupTestEnv(t)

	config := &credentials.Config{
		BridgeIP: "192.168.1.100",
	}
	err := credentials.SaveConfig(config)
	require.NoError(t, err)

	_, _, err = GetLocationFromConfig()
	assert.Error(t, err)
}

func TestSetLocationInConfig(t *testing.T) {
	setupTestEnv(t)

	err := SetLocationInConfig(37.7749, -122.4194)
	require.NoError(t, err)

	config, err := credentials.LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, 37.7749, config.Latitude)
	assert.Equal(t, -122.4194, config.Longitude)
}

func TestSetLocationInConfigExistingConfig(t *testing.T) {
	setupTestEnv(t)

	config := &credentials.Config{
		BridgeIP: "192.168.1.100",
	}
	err := credentials.SaveConfig(config)
	require.NoError(t, err)

	err = SetLocationInConfig(40.7128, -74.0060)
	require.NoError(t, err)

	updated, err := credentials.LoadConfig()
	require.NoError(t, err)
	assert.Equal(t, "192.168.1.100", updated.BridgeIP)
	assert.Equal(t, 40.7128, updated.Latitude)
	assert.Equal(t, -74.0060, updated.Longitude)
}

func TestSetLocationInConfigInvalidLatitude(t *testing.T) {
	setupTestEnv(t)

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
			err := SetLocationInConfig(tc.latitude, 0)
			assert.Error(t, err)
		})
	}
}

func TestSetLocationInConfigInvalidLongitude(t *testing.T) {
	setupTestEnv(t)

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
			err := SetLocationInConfig(0, tc.longitude)
			assert.Error(t, err)
		})
	}
}
