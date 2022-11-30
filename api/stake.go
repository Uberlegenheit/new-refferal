package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/log"
	"new-refferal/models"
)

func (api *API) Delegate(c *gin.Context) {
	var stake models.Stake
	if err := c.ShouldBindJSON(&stake); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stk, err := api.services.SaveDelegationTx(&stake)
	if err != nil {
		log.Error("[api] Delegate: SaveDelegationTx", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stk)
}
