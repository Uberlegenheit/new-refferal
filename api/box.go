package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/log"
	"new-refferal/models"
)

func (api *API) OpenBox(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		log.Error("[api] OpenBox: c.Get", zap.Error(fmt.Errorf("user context is empty")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context is empty"})
		return
	}
	user := val.(models.User)

	err := api.services.OpenBox(&user)
	if err != nil {
		log.Error("[api] OpenBox: OpenBox", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
