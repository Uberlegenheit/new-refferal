package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/filters"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (api *API) GetAllRewards(c *gin.Context) {
	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetAllRewards: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rewards, err := api.services.GetAllRewards(pagination)
	if err != nil {
		log.Error("[api] GetAllRewards: GetAllRewards", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (api *API) GetMyRewards(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		log.Error("[api] GetMyRewards: c.Get", zap.Error(fmt.Errorf("user context is empty")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context is empty"})
		return
	}
	user := val.(models.User)

	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetMyRewards: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rewards, err := api.services.GetUserRewardsByID(&user, pagination)
	if err != nil {
		log.Error("[api] GetAllRewards: GetAllRewards", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (api *API) GetMyLink(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		log.Error("[api] GetMyLink: c.Get", zap.Error(fmt.Errorf("user context is empty")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context is empty"})
		return
	}
	user := val.(models.User)

	link, err := api.services.GetLinkByUserID(&user)
	if err != nil {
		log.Error("[api] GetMyLink: GetLinkByUserID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, link)
}

func (api *API) GetInvitedFriends(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		log.Error("[api] GetInvitedFriends: c.Get", zap.Error(fmt.Errorf("user context is empty")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "user context is empty"})
		return
	}
	user := val.(models.User)

	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetInvitedFriends: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invited, err := api.services.GetInvitedUsersStakes(&user, pagination)
	if err != nil {
		log.Error("[api] GetInvitedFriends: GetInvitedUsersStakes", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invited)
}
