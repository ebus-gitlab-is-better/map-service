package route

import (
	"fmt"
	"log"
	"map-service/pkg/valhalla"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mojixcoder/gosrm"
	"github.com/twpayne/go-polyline"
)

type MapRoute struct {
	client *valhalla.Client
}

func NewMapRoute(client *valhalla.Client) *MapRoute {
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
	request := valhalla.RouteRequest{}
	for _, point := range points {
		request.Locations = append(request.Locations, valhalla.Location{
			Lat: point[0],
			Lon: point[1],
		})
	}
	request.Costing = "bus"
	//TODO request.ExcludeLocations
	route, err := r.client.Route(request)
	if err != nil {
		log.Fatal(err)
	}
	shape := route.Trip.Legs[0].Shape
	coords, _, _ := polyline.DecodeCoords([]byte(shape))
	fmt.Println(coords)
	c.JSON(200, coords)
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
