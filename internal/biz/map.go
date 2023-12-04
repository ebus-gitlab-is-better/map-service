package biz

import (
	"log"
	"map-service/internal/utils"
	"map-service/pkg/valhalla"

	"github.com/mojixcoder/gosrm"
)

type MapUseCase struct {
	client *valhalla.Client
}

func NewMapUseCase(client *valhalla.Client) *MapUseCase {
	return &MapUseCase{client: client}
}

func (uc *MapUseCase) GetPath(points []gosrm.Coordinate) string {
	request := valhalla.RouteRequest{}
	for _, point := range points {
		request.Locations = append(request.Locations, valhalla.Location{
			Lat: point[0],
			Lon: point[1],
		})
	}
	request.Costing = "bus"
	//TODO request.ExcludeLocations
	route, err := uc.client.Route(request)
	if err != nil {
		log.Fatal(err)
	}
	return route.Trip.Legs[0].Shape
}

func (uc *MapUseCase) CheckPath(shape string, point [2]float64) bool {
	threshold := 0.001
	b, _ := utils.IsPointNearPolyline(point, shape, threshold)
	return b
}
