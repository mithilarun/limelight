package astro

import (
	"math"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	degreesToRadians = math.Pi / 180.0
	radiansToDegrees = 180.0 / math.Pi
	// zenith angle for sunrise/sunset (accounts for atmospheric refraction and sun's radius)
	// Standard value from NOAA Solar Position Calculator
	sunriseZenith = 90.833
)

// CalculateSunrise calculates the sunrise time for a given location and date
func CalculateSunrise(latitude, longitude float64, date time.Time) (time.Time, error) {
	return calculateSunEvent(latitude, longitude, date, true)
}

// CalculateSunset calculates the sunset time for a given location and date
func CalculateSunset(latitude, longitude float64, date time.Time) (time.Time, error) {
	return calculateSunEvent(latitude, longitude, date, false)
}

// calculateSunEvent calculates sunrise or sunset using the standard astronomical algorithm
func calculateSunEvent(latitude, longitude float64, date time.Time, isSunrise bool) (time.Time, error) {
	if latitude < -90 || latitude > 90 {
		return time.Time{}, errors.Newf("invalid latitude: %f (must be between -90 and 90)", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return time.Time{}, errors.Newf("invalid longitude: %f (must be between -180 and 180)", longitude)
	}

	year, month, day := date.Date()
	location := date.Location()

	// Calculate day of year
	dayOfYear := date.YearDay()

	// Calculate longitude hour
	lngHour := longitude / 15.0

	// Calculate approximate time
	var t float64
	if isSunrise {
		t = float64(dayOfYear) + ((6.0 - lngHour) / 24.0)
	} else {
		t = float64(dayOfYear) + ((18.0 - lngHour) / 24.0)
	}

	// Calculate sun's mean anomaly
	M := (0.9856 * t) - 3.289

	// Calculate sun's true longitude
	L := M + (1.916 * sinDeg(M)) + (0.020 * sinDeg(2*M)) + 282.634
	L = normalizeDegrees(L)

	// Calculate sun's right ascension
	RA := radiansToDegrees * math.Atan(0.91764*tanDeg(L))
	RA = normalizeDegrees(RA)

	// Right ascension needs to be in the same quadrant as L
	Lquadrant := math.Floor(L/90.0) * 90.0
	RAquadrant := math.Floor(RA/90.0) * 90.0
	RA = RA + (Lquadrant - RAquadrant)

	// Convert RA to hours
	RA = RA / 15.0

	// Calculate sun's declination
	sinDec := 0.39782 * sinDeg(L)
	cosDec := math.Cos(math.Asin(sinDec))

	// Calculate sun's local hour angle
	cosH := (cosDeg(sunriseZenith) - (sinDec * sinDeg(latitude))) / (cosDec * cosDeg(latitude))

	// Check for polar day/night
	if cosH > 1 {
		// Sun never rises
		return time.Time{}, errors.New("sun never rises at this location on this date")
	}
	if cosH < -1 {
		// Sun never sets
		return time.Time{}, errors.New("sun never sets at this location on this date")
	}

	// Calculate hour angle
	var H float64
	if isSunrise {
		H = 360.0 - radiansToDegrees*math.Acos(cosH)
	} else {
		H = radiansToDegrees * math.Acos(cosH)
	}
	H = H / 15.0

	// Calculate local mean time
	T := H + RA - (0.06571 * t) - 6.622

	// Convert to UTC
	UT := T - lngHour
	UT = normalizeHours(UT)

	// Convert hours to time
	hours := int(UT)
	minutes := int((UT - float64(hours)) * 60)
	seconds := int(((UT-float64(hours))*60 - float64(minutes)) * 60)

	// Create time in UTC, then convert to local timezone
	utcTime := time.Date(year, month, day, hours, minutes, seconds, 0, time.UTC)
	localTime := utcTime.In(location)

	return localTime, nil
}

// sinDeg calculates sine of angle in degrees
func sinDeg(deg float64) float64 {
	return math.Sin(deg * degreesToRadians)
}

// cosDeg calculates cosine of angle in degrees
func cosDeg(deg float64) float64 {
	return math.Cos(deg * degreesToRadians)
}

// tanDeg calculates tangent of angle in degrees
func tanDeg(deg float64) float64 {
	return math.Tan(deg * degreesToRadians)
}

// normalizeDegrees normalizes an angle to 0-360 degrees
func normalizeDegrees(deg float64) float64 {
	deg = math.Mod(deg, 360.0)
	if deg < 0 {
		deg += 360.0
	}
	return deg
}

// normalizeHours normalizes hours to 0-24
func normalizeHours(hours float64) float64 {
	hours = math.Mod(hours, 24.0)
	if hours < 0 {
		hours += 24.0
	}
	return hours
}
