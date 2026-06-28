package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PhotonReverseGeocoderResp struct {
	Features []Properties `json:"features"`
}

type Properties struct {
	Properties struct {
		CountryCode string `json:"countrycode"`
	} `json:"properties"`
}

func ResolveCountry(ctx context.Context, lat float32, lon float32) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://localhost:2322/reverse?lon=%v&lat=%v", lon, lat), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	output := PhotonReverseGeocoderResp{}
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		return "", err
	}
	if len(output.Features) < 1 {
		return "", fmt.Errorf("no results in reverse geocoding")
	}

	return output.Features[0].Properties.CountryCode, nil
}
