package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.Engine, metrics *Registry) {
	res := resource{metrics}

	r.GET("/metrics/:site", res.Get)
	r.GET("/metrics", res.All)

}

type resource struct {
	metrics *Registry
}

func (r *resource) Get(c *gin.Context) {
	name := c.Param("site")

	res, ok := r.metrics.Counters[name]
	if !ok {
		c.JSON(http.StatusNotFound, "")
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": res.Name(), "value": res.Count()})
}

func (r *resource) All(c *gin.Context) {
	c.JSON(http.StatusOK, r.metrics.Stats())
}
