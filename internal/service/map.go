package service

import (
	"context"

	pb "map-service/api/map/v1"
	"map-service/internal/biz"

	"github.com/mojixcoder/gosrm"
)

type MapService struct {
	uc *biz.MapUseCase
	pb.UnimplementedMapServer
}

func NewMapService(uc *biz.MapUseCase) *MapService {
	return &MapService{
		uc: uc,
	}
}

func (s *MapService) GetPath(ctx context.Context, req *pb.GetPathRequest) (*pb.PathResponse, error) {
	points := make([]gosrm.Coordinate, 0)
	for _, point := range req.Points {
		points = append(points, gosrm.Coordinate{
			float64(point.Lon), float64(point.Lat),
		})
	}
	return &pb.PathResponse{
		Shape: s.uc.GetPath(points),
	}, nil
}
func (s *MapService) CheckPath(ctx context.Context, req *pb.CheckPathRequest) (*pb.CheckPathResponse, error) {
	return &pb.CheckPathResponse{
		IsValid: s.uc.CheckPath(req.Shape, [2]float64{float64(req.Point.Lon), float64(req.Point.Lat)}),
	}, nil
}
