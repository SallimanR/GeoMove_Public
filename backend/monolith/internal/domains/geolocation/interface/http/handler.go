package http

import (
	"context"
	"time"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/application/query"
	"monolith/internal/domains/geolocation/domain/entity"
	"monolith/internal/websockethub"
)

type GeolocationHandler struct {
	getClosestWithinRadiusMovingDrivers query.GetClosestWithinRadiusMovingDriversHandler
	gpsChannel                          *websockethub.PubSubChannel[entity.MovingDriverWithPoints]
}

func NewGeolocationHandler(
	getClosestMovingDrivers query.GetClosestWithinRadiusMovingDriversHandler,
	gpsChannel *websockethub.PubSubChannel[entity.MovingDriverWithPoints],
) GeolocationHandler {
	return GeolocationHandler{
		getClosestWithinRadiusMovingDrivers: getClosestMovingDrivers,
		gpsChannel:                          gpsChannel,
	}
}

func (h *GeolocationHandler) GetClosestWithinRadiusMovingDriversByIDs(
	ctx context.Context,
	request GetClosestWithinRadiusMovingDriversByIDsRequestObject,
) (
	GetClosestWithinRadiusMovingDriversByIDsResponseObject,
	error,
) {
	queryCmd := query.GetClosestWithinRadiusMovingDriversQuery{
		Location: dto.Location{
			Latitude:  request.Params.Lat,
			Longitude: request.Params.Lon,
		},
		RadiusMeters: uint32(request.Params.RadiusMeters),
	}

	drivers, err := h.getClosestWithinRadiusMovingDrivers.Handle(ctx, queryCmd)
	if err != nil {
		return GetClosestWithinRadiusMovingDriversByIDs500Response{}, nil
	}

	resp := make(GetClosestWithinRadiusMovingDriversByIDs200JSONResponse, 0, len(drivers))
	for _, driver := range drivers {
		data, ok := h.gpsChannel.Messages.Load(driver.DriverID)
		if !ok {
			continue
		}

		md := MovingDriver{
			DriverId:   int(data.DriverID),
			Lat:        data.Latitude,
			Lon:        data.Longitude,
			TravelTime: float32(data.TravelTime.Sub(time.Time{}).Milliseconds()),
			PathMeters: int(data.PathMeters),
			Points:     convertPointsForAPI(data.Points),
		}

		resp = append(resp, md)
	}

	return resp, nil
}

func (h *GeolocationHandler) GetMovingDriversByIDs(ctx context.Context, request GetMovingDriversByIDsRequestObject) (
	GetMovingDriversByIDsResponseObject,
	error,
) {
	if request.Body == nil {
		return GetMovingDriversByIDs400Response{}, nil
	}

	ids := request.Body.Ids
	if len(ids) == 0 {
		return GetMovingDriversByIDs200JSONResponse{}, nil
	}

	resp := make(GetMovingDriversByIDs200JSONResponse, 0, len(ids))
	for _, id := range ids {
		data, ok := h.gpsChannel.Messages.Load(int64(id))
		if !ok {
			continue
		}

		md := MovingDriver{
			DriverId:   int(data.DriverID),
			Lat:        data.Latitude,
			Lon:        data.Longitude,
			TravelTime: float32(data.TravelTime.Sub(time.Time{}).Milliseconds()),
			PathMeters: int(data.PathMeters),
			Points:     convertPointsForAPI(data.Points),
		}

		resp = append(resp, md)
	}

	return resp, nil
}

func convertPointsForAPI(points [][2]float32) [][]float32 {
	result := make([][]float32, len(points))
	for i, p := range points {
		result[i] = []float32{p[0], p[1]}
	}
	return result
}
