package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/filters"
	"new-refferal/log"
	"new-refferal/models"
)

func (api *API) CreatePayout(c *gin.Context) {
	var payout models.Payout
	if err := c.ShouldBindJSON(&payout); err != nil {
		log.Error("[api] CreatePayout: ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	saved, err := api.services.CreatePayout(&payout)
	if err != nil {
		log.Error("[api] CreatePayout: CreatePayout", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, saved)
}

func (api *API) UpdatePayout(c *gin.Context) {
	var payout models.Payout
	if err := c.ShouldBindJSON(&payout); err != nil {
		log.Error("[api] UpdatePayout: ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.services.UpdatePayout(&payout)
	if err != nil {
		log.Error("[api] UpdatePayout: UpdatePayout", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payout)
}

func (api *API) GetPayouts(c *gin.Context) {
	var pagination filters.Pagination
	if err := c.Bind(&pagination); err != nil {
		log.Error("[api] GetPayouts: Bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payouts, err := api.services.GetPayouts(pagination)
	if err != nil {
		log.Error("[api] GetPayouts: GetPayouts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payouts)
}
