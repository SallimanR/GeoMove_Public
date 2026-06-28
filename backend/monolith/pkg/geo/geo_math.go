package geo

import "math"

func CalculateBearing(lat1, lon1, lat2, lon2 float64) float32 {
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δλ := (lon2 - lon1) * math.Pi / 180

	// Calculate bearing using the formula:
	// θ = atan2( sin Δλ ⋅ cos φ2 , cos φ1 ⋅ sin φ2 − sin φ1 ⋅ cos φ2 ⋅ cos Δλ )
	y := math.Sin(Δλ) * math.Cos(φ2)
	x := math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ)
	θ := math.Atan2(y, x)

	// Convert from radians to degrees and normalize to 0-360
	bearing := θ * 180 / math.Pi
	bearing = math.Mod(bearing+360, 360)
	return float32(bearing)
}

const earthRadiusKm = 6371.0

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float32 {
	x := (earthRadiusKm * math.Pi / 180) * (lat2 - lat1)
	y := (earthRadiusKm * math.Pi / 180) * (lon2 - lon1) * math.Cos(lat1)
	distance := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	return float32(distance)
}
