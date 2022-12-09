package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"new-refferal/conf"
)

func (api *API) Index(c *gin.Context) {
	c.String(http.StatusOK, "This is a service '%s'", conf.Service)
}

func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
