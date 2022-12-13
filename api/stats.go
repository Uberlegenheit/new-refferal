package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/filters"
	"new-refferal/log"
)

func (api *API) GetTotalRewardStats(c *gin.Context) {
	rewards, err := api.services.GetTotalRewardStats()
	if err != nil {
		log.Error("[api] GetTotalRewardStats: GetTotalRewardStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (api *API) GetTotalStats(c *gin.Context) {
	var req filters.PeriodInfoRequest
	if err := c.Bind(&req); err != nil {
		log.Error("[api] GetTotalStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetTotalStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := api.services.GetTotalStats(req, pagination)
	if err != nil {
		log.Error("[api] GetTotalStats: GetTotalStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (api *API) GetTotalStakeStats(c *gin.Context) {
	var req filters.PeriodInfoRequest
	if err := c.Bind(&req); err != nil {
		log.Error("[api] GetTotalStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetTotalStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := api.services.GetTotalStakeStats(req, pagination)
	if err != nil {
		log.Error("[api] GetTotalStakeStats: GetTotalStakeStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (api *API) GetFriendsStakeStats(c *gin.Context) {
	var req filters.PeriodInfoRequest
	if err := c.Bind(&req); err != nil {
		log.Error("[api] GetFriendsStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetFriendsStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := api.services.GetFriendsStakeStats(req, pagination)
	if err != nil {
		log.Error("[api] GetFriendsStakeStats: GetFriendsStakeStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (api *API) GetRewardPaymentStats(c *gin.Context) {
	var req filters.PeriodInfoRequest
	if err := c.Bind(&req); err != nil {
		log.Error("[api] GetRewardPaymentStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetRewardPaymentStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := api.services.GetRewardPaymentStats(req, pagination)
	if err != nil {
		log.Error("[api] GetRewardPaymentStats: GetRewardPaymentStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (api *API) GetInvitationsStats(c *gin.Context) {
	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetInvitationsStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := api.services.GetUsersInvitationsStats(pagination)
	if err != nil {
		log.Error("[api] GetInvitationsStats: GetUsersInvitationsStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
