package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (api *API) SomeMiddleware() gin.HandlerFunc {
	fmt.Println("setting up a middleware")
	return func(c *gin.Context) {
		fmt.Println("middleware was triggered")
		c.Next()
	}
}
