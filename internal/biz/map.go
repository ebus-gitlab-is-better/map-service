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

type Path struct {
	Shape   string
	Time    []float32
	Lengths []float32
}

func (uc *MapUseCase) GetPath(points []gosrm.Coordinate) *Path {
	request := valhalla.RouteRequest{}
	for _, point := range points {
		request.Locations = append(request.Locations, valhalla.Location{
			Lon: point[0],
			Lat: point[1],
		})
	}
	request.Costing = "bus"
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
		return nil
	}
	times := make([]float32, 0)
	coords := make([][2]float64, 0)
	lengths := make([]float32, 0)
	for _, leg := range route.Trip.Legs {
		coord := utils.DecodePolyline(&leg.Shape)
		coords = append(coords, coord...)
		times = append(times, float32(leg.Summary.Time))
		lengths = append(lengths, float32(leg.Summary.Length))
	}

	shape := utils.EncodePolyline(coords)
	return &Path{Shape: shape, Time: times, Lengths: lengths}
}

func (uc *MapUseCase) CheckPath(shape string, point [2]float64) bool {
	threshold := 0.001
	b, _ := utils.IsPointNearPolyline(point, shape, threshold)
	return b
}
