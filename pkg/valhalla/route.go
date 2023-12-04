package valhalla

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type DirectionsOptions struct {
	Units     string `json:"units,omitempty"`
	Narrative bool   `json:"narrative,omitempty"`
	Language  string `json:"language,omitempty"`
}
type CostingOptions struct {
	Bus CostingOption `json:"bus,omitempty"`
}

type CostingOption struct {
	CountryCrossingPenalty int `json:"country_crossing_penalty,omitempty"`
	TollBoothPenalty       int `json:"toll_booth_penalty,omitempty"`
	CountryCrossingCost    int `json:"country_crossing_cost,omitempty"`

	Weight float32 `json:"weight,omitempty"`
	Height float32 `json:"height,omitempty"`
	Width  float32 `json:"width,omitempty"`
	Length float32 `json:"length,omitempty"`
}
type RouteRequest struct {
	Locations         []Location        `json:"locations,omitempty"`
	Costing           string            `json:"costing,omitempty"`
	DirectionsOptions DirectionsOptions `json:"directions_options,omitempty"`
	ID                string            `json:"id,omitempty"`
	CostingOptions    CostingOptions    `json:"costing_options,omitempty"`
	ExcludeLocations  []Location        `json:"exclude_locations,omitempty"`
	Verbose           bool              `json:"verbose,omitempty"`
}

type Contour struct {
	Time  int    `json:"time,omitempty"`
	Color string `json:"color,omitempty"`
}

type RouteResponse struct {
	ID   string `json:"id"`
	Trip struct {
		Language      string `json:"language"`
		Status        int    `json:"status"`
		Units         string `json:"units"`
		StatusMessage string `json:"status_message"`
		Legs          []struct {
			Shape   string `json:"shape"`
			Summary struct {
				MaxLon float64 `json:"max_lon"`
				MaxLat float64 `json:"max_lat"`
				Time   int     `json:"time"`
				Length float64 `json:"length"`
				MinLat float64 `json:"min_lat"`
				MinLon float64 `json:"min_lon"`
			} `json:"summary"`
			Maneuvers []Maneuver `json:"maneuvers"`
		} `json:"legs"`
		Summary struct {
			MaxLon float64 `json:"max_lon"`
			MaxLat float64 `json:"max_lat"`
			Time   int     `json:"time"`
			Length float64 `json:"length"`
			MinLat float64 `json:"min_lat"`
			MinLon float64 `json:"min_lon"`
		} `json:"summary"`
		Locations []struct {
			OriginalIndex int     `json:"original_index"`
			Type          string  `json:"type"`
			Lon           float64 `json:"lon"`
			Lat           float64 `json:"lat"`
			SideOfStreet  string  `json:"side_of_street"`
		} `json:"locations"`
	} `json:"trip"`
}

type Maneuver struct {
	TravelMode                       string   `json:"travel_mode"`
	BeginShapeIndex                  int      `json:"begin_shape_index"`
	Length                           float64  `json:"length"`
	Time                             int      `json:"time"`
	Type                             int      `json:"type"`
	EndShapeIndex                    int      `json:"end_shape_index"`
	Instruction                      string   `json:"instruction"`
	VerbalPreTransitionInstruction   string   `json:"verbal_pre_transition_instruction"`
	TravelType                       string   `json:"travel_type"`
	StreetNames                      []string `json:"street_names,omitempty"`
	VerbalTransitionAlertInstruction string   `json:"verbal_transition_alert_instruction,omitempty"`
	VerbalPostTransitionInstruction  string   `json:"verbal_post_transition_instruction,omitempty"`
}

type ValhallaError struct {
	ErrorCode  int    `json:"error_code"`
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
}

func (c *Client) Route(request RouteRequest) (RouteResponse, error) {
	r, err := json.Marshal(request)
	if err != nil {
		return RouteResponse{}, err
	}
	fmt.Println(request)
	response, err := c.request("GET", "route", bytes.NewReader(r))
	if err != nil {
		return RouteResponse{}, err
	}

	result := RouteResponse{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return RouteResponse{}, err
	}

	return result, nil
}
