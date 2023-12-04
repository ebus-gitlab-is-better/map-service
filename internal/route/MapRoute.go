package route

import (
	"encoding/json"
	"io"
	"map-service/internal/biz"
	"map-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MapRoute struct {
	uc *biz.MapUseCase
}

func NewMapRoute(uc *biz.MapUseCase) *MapRoute {
	return &MapRoute{uc: uc}
}

func (r *MapRoute) Register(router *gin.RouterGroup) {
	router.GET("/:coordinates", r.GetPath)
}

type CoordsResponse struct {
	Coords [][2]float64 `json:"coords"`
}

// @Summary	Get path
// @Accept		json
// @Produce	json
// @Tags		map
// @Param			coordinates	path	string	true	"[{longitude},{latitude};{longitude},{latitude}[;{longitude},{latitude} ...]"
//
//	@Success	200	{object}	route.CoordsResponse
//
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/maps/{coordinates}/ [get]
func (r *MapRoute) GetPath(c *gin.Context) {
	coordinates := c.Param("coordinates")
	points, err := utils.ParseCoordinates(coordinates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shape := r.uc.GetPath(points)
	coords := utils.DecodePolyline(&shape)
	c.JSON(200, CoordsResponse{
		Coords: coords,
	})
}

type PathDTO struct {
	Shape string     `json:"shape"`
	Point [2]float64 `json:"point"`
}

type PathResponse struct {
	Check bool `json:"check"`
}

// @Summary	Get path
// @Accept		json
// @Produce	json
// @Tags		map
// @Param		dto	body	route.PathDTO	true	"dto"
//
//	@Success	200	{object}	route.PathResponse
//
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/maps/check [post]
func (r *MapRoute) CheckInPath(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := PathDTO{}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, PathResponse{
		Check: r.uc.CheckPath(dto.Shape, dto.Point),
	})
}
