package graphopper

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/interface/websocket/pb"
)

type MapMatchedResponse struct {
	Paths []Path `json:"paths"`
}

type Path struct {
	Distance float32 `json:"distance"`
	Time     uint32  `json:"time"`
	Points   Points  `json:"points"`
}

// TODO: encode polyline directly? (remove "points_encoded=false" in get url)
type Points struct {
	Coordinates []dto.LocationRaw `josn:"coordinates"`
}

type MapMatchRequest struct {
	Coordinates *[]dto.Location
	Timestamps  *[]int64
}

const APIBase = "https://geomove.online/routing"

// TODO: add profile for graphopper matching
func MatchMapCoordinates(ctx context.Context, data string) (*Path, bool) {
	output := &Path{}
	// TODO: context
	url := fmt.Sprintf("%s/match?profile=car&points_encoded=false&details=time&details=distance&details=max_speed", APIBase)
	resp, err := http.Post(url, "application/gpx+xml", strings.NewReader(data))
	if err != nil {
		log.Printf("error in post req: %s", err)
		return output, false
	}
	if resp.StatusCode != http.StatusOK {
		return output, false
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return output, false
	}
	_ = resp.Body.Close()

	var matchedRes MapMatchedResponse
	err = json.Unmarshal(respBody, &matchedRes)
	if err != nil {
		log.Printf("Failed to unmarshal map matched coords: %s", err)
		return output, false
	}
	output = &matchedRes.Paths[0]

	return output, true
}

type GPX struct {
	Trk Trk `xml:"trk"`
}

type Trk struct {
	TrkSeg TrkSeg `xml:"trkseg"`
}

type TrkSeg struct {
	TrkPt []TrkPt `xml:"trkpt"`
}

type TrkPt struct {
	Lat  float32 `xml:"lat,attr"`
	Lon  float32 `xml:"lon,attr"`
	Time string  `xml:"time,omitempty"`
}

func BuildGraphhopperRequestFromWS(coords []*pb.Point, timestamps []int64) (string, error) {
	coordsNumber := len(coords)
	timestampsNumber := len(timestamps)
	points := make([]TrkPt, coordsNumber)

	for i, c := range coords {
		point := TrkPt{Lat: c.Latitude, Lon: c.Longitude}
		if timestampsNumber > i && timestamps[i] != 0 {
			point.Time = time.Unix(timestamps[i], 0).UTC().Format(time.RFC3339)
		}
		points[i] = point
	}

	gpx := GPX{
		Trk: Trk{
			TrkSeg: TrkSeg{TrkPt: points},
		},
	}
	out, err := xml.MarshalIndent(gpx, "", "")
	if err != nil {
		return "", err
	}

	return string(out), nil
}
