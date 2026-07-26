package pricing

const (
	basePickupFee        = 2500.0
	perBlockedWheelFee   = 500.0
	minPrice             = 3000.0
	cityDistanceKm       = 50.0
	intercityRatePercent = 0.5
)

type EstimateRequest struct {
	DistanceMeters       int32
	CarWeightKg          int32
	HowManyWheelsBlocked int16
}

func ratePerKm(weightKg int32) float64 {
	switch {
	case weightKg <= 1500:
		return 60
	case weightKg <= 2500:
		return 80
	case weightKg <= 3500:
		return 95
	default:
		return 120
	}
}

func CalculatePrice(req EstimateRequest) int32 {
	distanceKm := float64(req.DistanceMeters) / 1000.0
	baseRate := ratePerKm(req.CarWeightKg)

	cityKm := min(distanceKm, cityDistanceKm)
	intercityKm := max(distanceKm-cityDistanceKm, 0)

	transportFee := cityKm*baseRate + intercityKm*baseRate*intercityRatePercent
	blockedFee := float64(req.HowManyWheelsBlocked) * perBlockedWheelFee

	total := basePickupFee + transportFee + blockedFee
	if total < minPrice {
		total = minPrice
	}

	return int32(total)
}
