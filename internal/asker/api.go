package asker

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.Engine, service Service) {
	res := resource{service}

	r.GET("/status/min", res.Min)
	r.GET("/status/max", res.Max)
	r.GET("/status/random", res.Random)

	r.GET("/status/site/:site", res.CheckStatus)
}

type resource struct {
	service Service
}

func (r *resource) Min(c *gin.Context) {
	res, err := r.service.GetMin(c)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *resource) Max(c *gin.Context) {
	res, err := r.service.GetMax(c)
	if err != nil {
		if err != nil {
			r.handleError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

func (r *resource) Random(c *gin.Context) {
	res, err := r.service.GetRandom(c)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *resource) CheckStatus(c *gin.Context) {
	name := c.Param("site")
	res, err := r.service.Get(c, name)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *resource) handleError(c *gin.Context, err error) {
	switch v := err.(type) {
	case *NotFoundError:
		c.JSON(http.StatusNotFound, v.Error())
	case *NoResponse:
		c.JSON(http.StatusNoContent, v.Error())
	default:
		c.JSON(http.StatusInternalServerError, "unknown error")
	}

	return
}
