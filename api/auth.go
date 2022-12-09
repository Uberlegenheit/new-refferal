package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"new-refferal/log"
	"new-refferal/models"
)

const adminRole = "admin"

func (api *API) SomeMiddleware() gin.HandlerFunc {
	fmt.Println("setting up a middleware")
	return func(c *gin.Context) {
		fmt.Println("middleware was triggered")
		c.Next()
	}
}

func (api *API) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBind(&user); err != nil {
			log.Error("[api] AdminMiddleware: fetch user", zap.Error(err))
			return
		}

		if user.Role != adminRole {
			log.Error("[api] AdminMiddleware: ", zap.Error(fmt.Errorf("user isn't admin")))
			return
		}

		c.Next()
	}
}
