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
		log.Error("[api] Delegate: c.Get", zap.Error(fmt.Errorf("user context is empty")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context is empty"})
		return
	}
	user := val.(models.User)

	var stake models.Stake
	if err := c.ShouldBindJSON(&stake); err != nil {
		log.Error("[api] Delegate: ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stake.UserID = user.ID
	key := c.Query("key")

	ok, err := api.services.CheckDelegationKey(&user, key)
	if err != nil {
		log.Error("[api] Delegate: CheckDelegationKey", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		log.Error("[api] Delegate: CheckDelegationKey", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key"})
		return
	}

	dbStake, err := api.services.GetDelegationByTxHash(&stake)
	if err != nil {
		log.Error("[api] Delegate: GetDelegationByTxHash", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if dbStake != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hash already registered"})
		return
	}

	stk, err := api.services.SaveDelegationTx(&stake, &user)
	if err != nil {
		log.Error("[api] Delegate: SaveDelegationTx", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stk)
}

func (api *API) GetDelegationKey(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		log.Error("[api] GetDelegationKey: c.Get", zap.Error(fmt.Errorf("user context is empty")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context is empty"})
		return
	}
	user := val.(models.User)

	stk, err := api.services.GetDelegationKey(&user)
	if err != nil {
		log.Error("[api] GetDelegationKey: GetDelegationKey", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stk)
}
