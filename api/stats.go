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
	if err := c.BindQuery(&req); err != nil {
		log.Error("[api] GetTotalStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := api.services.GetTotalStats(req)
	if err != nil {
		log.Error("[api] GetTotalStats: GetTotalStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (api *API) GetTotalStakeStats(c *gin.Context) {
	var req filters.PeriodInfoRequest
	if err := c.BindQuery(&req); err != nil {
		log.Error("[api] GetTotalStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagination filters.Pagination
	if err := c.BindQuery(&pagination); err != nil {
		log.Error("[api] GetTotalStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, length, err := api.services.GetTotalStakeStats(req, pagination)
	if err != nil {
		log.Error("[api] GetTotalStakeStats: GetTotalStakeStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":  stats,
		"length": length,
	})
}

func (api *API) GetFriendsStakeStats(c *gin.Context) {
	var req filters.PeriodInfoRequest
	if err := c.BindQuery(&req); err != nil {
		log.Error("[api] GetFriendsStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagination filters.Pagination
	if err := c.BindQuery(&pagination); err != nil {
		log.Error("[api] GetFriendsStakeStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, length, err := api.services.GetFriendsStakeStats(req, pagination)
	if err != nil {
		log.Error("[api] GetFriendsStakeStats: GetFriendsStakeStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":  stats,
		"length": length,
	})
}

func (api *API) GetBoxesStats(c *gin.Context) {
	var req filters.PeriodInfoRequest
	if err := c.BindQuery(&req); err != nil {
		log.Error("[api] GetBoxesStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pagination filters.Pagination
	if err := c.BindQuery(&pagination); err != nil {
		log.Error("[api] GetBoxesStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, length, err := api.services.GetBoxesStats(req, pagination)
	if err != nil {
		log.Error("[api] GetBoxesStats: GetBoxesStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":  stats,
		"length": length,
	})
}

func (api *API) GetInvitationsStats(c *gin.Context) {
	var pagination filters.Pagination
	if err := c.BindQuery(&pagination); err != nil {
		log.Error("[api] GetInvitationsStats: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, length, err := api.services.GetUsersInvitationsStats(pagination)
	if err != nil {
		log.Error("[api] GetInvitationsStats: GetUsersInvitationsStats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":  stats,
		"length": length,
	})
}
