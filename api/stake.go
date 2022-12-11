package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/log"
	"new-refferal/models"
)

func (api *API) Delegate(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		log.Error("[api] GetMyLink: c.Get", zap.Error(fmt.Errorf("user context is empty")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context is empty"})
		return
	}
	user := val.(models.User)

	var stake models.Stake
	if err := c.ShouldBindJSON(&stake); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stake.UserID = user.ID

	stk, err := api.services.SaveDelegationTx(&stake)
	if err != nil {
		log.Error("[api] Delegate: SaveDelegationTx", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stk)
}
