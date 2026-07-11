package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const APIBase = "http://localhost:2322"

type PhotonReverseGeocoderResp struct {
	Features []Properties `json:"features"`
}

type Properties struct {
	Properties struct {
		CountryCode string `json:"countrycode"`
		Name        string `json:"name"`
		Street      string `json:"street"`
		Housenumber string `json:"housenumber"`
		City        string `json:"city"`
	} `json:"properties"`
}

func ResolveCountry(ctx context.Context, lat float32, lon float32) (string, error) {
	prop, err := reverseGeocode(ctx, lat, lon)
	if err != nil {
		return "", err
	}
	return prop.CountryCode, nil
}

func ReverseGeocode(ctx context.Context, lat float32, lon float32) (string, error) {
	prop, err := reverseGeocode(ctx, lat, lon)
	if err != nil {
		return "", err
	}
	parts := []string{prop.Name, prop.Street, prop.Housenumber, prop.City}
	filtered := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	return strings.Join(filtered, ", "), nil
}

func reverseGeocode(ctx context.Context, lat float32, lon float32) (*PhotonAddressProps, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/reverse?lon=%v&lat=%v", APIBase, lon, lat), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	output := PhotonReverseGeocoderResp{}
	err = json.Unmarshal(respBody, &output)
	if err != nil {
		return nil, err
	}
	if len(output.Features) < 1 {
		return nil, fmt.Errorf("no results in reverse geocoding")
	}

	return &PhotonAddressProps{
		CountryCode: output.Features[0].Properties.CountryCode,
		Name:        output.Features[0].Properties.Name,
		Street:      output.Features[0].Properties.Street,
		Housenumber: output.Features[0].Properties.Housenumber,
		City:        output.Features[0].Properties.City,
	}, nil
}

type PhotonAddressProps struct {
	CountryCode string
	Name        string
	Street      string
	Housenumber string
	City        string
}
