package biz

import (
	"context"
	"log"
	"map-service/internal/utils"
	"map-service/pkg/valhalla"

	accidentS "map-service/api/accident/v1"

	"github.com/mojixcoder/gosrm"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MapUseCase struct {
	client         *valhalla.Client
	accidentClient accidentS.AccidentClient
}

func NewMapUseCase(
	client *valhalla.Client,
	accidentClient accidentS.AccidentClient) *MapUseCase {
	return &MapUseCase{client: client, accidentClient: accidentClient}
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
	accidents, err := uc.accidentClient.ListAccident(context.TODO(), &emptypb.Empty{})
	if err == nil {
		excludeLocations := make([]valhalla.Location, 0)
		for _, accident := range accidents.Accidents {
			excludeLocations = append(excludeLocations, valhalla.Location{
				Lat: float64(accident.Lat),
				Lon: float64(accident.Lon),
			})
		}
		request.ExcludeLocations = excludeLocations
	}
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
