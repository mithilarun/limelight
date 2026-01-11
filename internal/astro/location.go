package astro

import (
	"github.com/cockroachdb/errors"
	"github.com/mithilarun/limelight/internal/credentials"
)

// GetLocationFromConfig retrieves the latitude and longitude from the config file
func GetLocationFromConfig() (latitude, longitude float64, err error) {
	config, err := credentials.LoadConfig()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to load config")
	}

	if config == nil {
		return 0, 0, errors.New("config file does not exist")
	}

	if config.Latitude == 0 && config.Longitude == 0 {
		return 0, 0, errors.New("latitude and longitude not set in config")
	}

	if config.Latitude < -90 || config.Latitude > 90 {
		return 0, 0, errors.Newf("invalid latitude in config: %f (must be between -90 and 90)", config.Latitude)
	}
	if config.Longitude < -180 || config.Longitude > 180 {
		return 0, 0, errors.Newf("invalid longitude in config: %f (must be between -180 and 180)", config.Longitude)
	}

	return config.Latitude, config.Longitude, nil
}

// SetLocationInConfig updates the latitude and longitude in the config file
func SetLocationInConfig(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return errors.Newf("invalid latitude: %f (must be between -90 and 90)", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return errors.Newf("invalid longitude: %f (must be between -180 and 180)", longitude)
	}

	config, err := credentials.LoadConfig()
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	if config == nil {
		config = &credentials.Config{}
	}

	config.Latitude = latitude
	config.Longitude = longitude

	if err := credentials.SaveConfig(config); err != nil {
		return errors.Wrap(err, "failed to save config")
	}

	return nil
}

// SetLocationByCity looks up a city's coordinates using geocoding and saves them to config
func SetLocationByCity(cityName string) (*GeocodingResult, error) {
	result, err := GeocodeCity(cityName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to geocode city: %s", cityName)
	}

	if err := SetLocationInConfig(result.Latitude, result.Longitude); err != nil {
		return nil, errors.Wrap(err, "failed to save location to config")
	}

	return result, nil
}
