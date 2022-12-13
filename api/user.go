package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/log"
	"new-refferal/models"
)

func (api *API) Me(c *gin.Context) {
	val, ok := c.Get("user")
	if !ok {
		log.Error("[api] Me: Get", zap.Error(fmt.Errorf("cannot fetch active user")))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot fetch active user"})
		return
	}
	user, ok := val.(models.User)
	if !ok {
		log.Error("[api] Me: cast", zap.Error(fmt.Errorf("cannot cast active user")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot cast active user"})
		return
	}

	box, err := api.services.GetAvailableBoxesByUserID(user.ID)
	if err != nil {
		log.Error("[api] Me: GetAvailableBoxesByUserID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	link, err := api.services.GetLinkByUserID(&user)
	if err != nil {
		log.Error("[api] Me: GetLinkByUserID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stake, err := api.services.GetMyStakeSum(user.ID)
	if err != nil {
		log.Error("[api] Me: GetMyStakeSum", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"box":   box,
		"link":  link,
		"stake": stake,
	})
}

func (api *API) SignIn(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("[api] SignIn: ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usr, err := api.services.LogInOrRegister(&user)
	if err != nil {
		log.Error("[api] SignIn: LogInOrRegister", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	td, err := api.services.CreateToken(usr.WalletAddress)
	if err != nil {
		log.Error("[api] SignIn: CreateToken", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = api.services.CreateAuth(usr.WalletAddress, td)
	if err != nil {
		log.Error("[api] SignIn: CreateAuth", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"access_token":    td.AccessToken,
		"refresh_token":   td.RefreshToken,
		"access_expired":  td.AtExpires,
		"refresh_expired": td.AtExpires,
	})
}

func (api *API) SignOut(c *gin.Context) {
	au, err := api.services.ExtractTokenMetadata(c)
	if err != nil {
		log.Error("[api] SignOut: ExtractTokenMetadata", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	refreshUuid, _ := api.services.FetchAuth(&models.AccessDetails{
		AccessUuid: fmt.Sprintf("%s_refresh", au.AccessUuid),
	})

	delErr := api.services.DeleteAuth(au.AccessUuid,
		refreshUuid,
		fmt.Sprintf("%s_refresh", au.AccessUuid),
		fmt.Sprintf("%s_access", refreshUuid))
	if delErr != nil {
		log.Error("[api] SignOut: DeleteAuth", zap.Error(delErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": delErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (api *API) Refresh(c *gin.Context) {
	td, err := api.services.Refresh(c.Request)
	if err != nil {
		log.Error("[api] Refresh: Refresh", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"access_token":    td.AccessToken,
		"refresh_token":   td.RefreshToken,
		"access_expired":  td.AtExpires,
		"refresh_expired": td.AtExpires,
	})
}
