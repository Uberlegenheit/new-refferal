package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"new-refferal/log"
	"new-refferal/models"
)

const adminRole = "admin"

func (api *API) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, err := api.services.ExtractTokenMetadata(c)
		if err != nil {
			log.Error("[api] AuthMiddleware: ExtractTokenMetadata", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		walletAddr, err := api.services.FetchAuth(ad)
		if err != nil || walletAddr == "" {
			log.Error("[api] AuthMiddleware: FetchAuth", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot fetch active user"})
			c.Abort()
			return
		}

		user, err := api.services.GetUserByWalletAddress(walletAddr)
		if err != nil {
			log.Error("[api] AuthMiddleware: GetUserByWalletAddress", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", *user)
		c.Next()
	}
}

func (api *API) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user")
		if !ok {
			log.Error("[api] AdminMiddleware: c.Get", zap.Error(fmt.Errorf("cannot get active user")))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get active user"})
			c.Abort()
			return
		}

		user, ok := val.(models.User)
		if !ok {
			log.Error("[api] AdminMiddleware: cast", zap.Error(fmt.Errorf("cannot cast val to user")))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot cast val to user"})
			c.Abort()
			return
		}
		if user.Role != adminRole {
			log.Error("[api] AdminMiddleware: ", zap.Error(fmt.Errorf("user isn't admin")))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user isn't admin"})
			c.Abort()
			return
		}

		c.Next()
	}
}
