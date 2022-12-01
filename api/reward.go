package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/log"
	"new-refferal/models"
)

func (api *API) UpdateReward(c *gin.Context) {
	var reward models.Reward
	if err := c.ShouldBindJSON(&reward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.services.UpdateReward(&reward)
	if err != nil {
		log.Error("[api] UpdateReward: UpdateReward", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
