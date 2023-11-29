package route

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mojixcoder/gosrm"
)

type MapRoute struct {
	client gosrm.OSRMClient
}

func NewMapRoute(client gosrm.OSRMClient) *MapRoute {
	return &MapRoute{client: client}
}

func (r *MapRoute) Register(router *gin.RouterGroup) {
	router.GET("/:coordinates", r.GetPath)
}

// @Summary	Get path
// @Accept		json
// @Produce	json
// @Tags		map
// @Param			coordinates	path	string	true	"[{longitude},{latitude};{longitude},{latitude}[;{longitude},{latitude} ...]"
//
//	@Success	200	{object}	gosrm.LineString
//
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/maps/{coordinates}/ [get]
func (r *MapRoute) GetPath(c *gin.Context) {
	coordinates := c.Param("coordinates")
	points, err := parseCoordinates(coordinates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	routeRes, err := gosrm.Route[gosrm.LineString](context.Background(), r.client, gosrm.Request{
		Profile:     gosrm.ProfileDriving,
		Coordinates: points,
	}, gosrm.WithGeometries(gosrm.GeometryGeoJSON))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, routeRes.Routes[0].Geometry)
}

func parseCoordinates(input string) ([]gosrm.Coordinate, error) {
	var points []gosrm.Coordinate
	coords := strings.Split(input, ";")
	for _, coord := range coords {
		parts := strings.Split(coord, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid coordinate format")
		}
		lon, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, err
		}
		lat, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, err
		}
		points = append(points, gosrm.Coordinate{lon, lat})
	}
	return points, nil
}
