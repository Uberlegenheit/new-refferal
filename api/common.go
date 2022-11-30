package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"new-refferal/conf"
	"strconv"
)

func (api *API) Index(c *gin.Context) {
	c.String(http.StatusOK, "This is a service '%s'", conf.Service)
}

func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (api *API) Name(c *gin.Context) {
	str := struct {
		Name string `json:"name"`
	}{}

	if err := c.ShouldBindJSON(&str); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "Specified name is %s", str.Name)
}

func (api *API) Read(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if id == 12 {
		c.JSON(http.StatusOK, gin.H{
			"name":    "Jason",
			"age":     "33",
			"country": "USA",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": false,
	})
}
