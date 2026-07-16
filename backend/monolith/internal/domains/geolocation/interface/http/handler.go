package http

import (
	"context"
	"encoding/json"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/application/query"
	"monolith/internal/websockethub"
)

type GeolocationHandler struct {
	getClosestWithinRadiusMovingDrivers query.GetClosestWithinRadiusMovingDriversHandler
	gpsChannel                          *websockethub.PubSubChannel
}

func NewGeolocationHandler(
	getClosestMovingDrivers query.GetClosestWithinRadiusMovingDriversHandler,
	gpsChannel *websockethub.PubSubChannel,
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

	resp := make(GetClosestWithinRadiusMovingDriversByIDs200JSONResponse, len(drivers))
	for i, d := range drivers {
		resp[i] = MovingDriver{
			DriverId:   int(d.DriverID),
			Lat:        d.Latitude,
			Lon:        d.Longitude,
			TravelTime: float32(d.TravelTime.Hour()*3600 + d.TravelTime.Minute()*60 + d.TravelTime.Second()),
			PathMeters: int(d.PathMeters),
			Points:     nil,
		}
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
		msgBytes, ok := h.gpsChannel.Messages.Load(int64(id))
		if !ok {
			continue
		}

		var md MovingDriver
		if err := json.Unmarshal(msgBytes, &md); err != nil {
			continue
		}

		resp = append(resp, md)
	}

	return resp, nil
}
