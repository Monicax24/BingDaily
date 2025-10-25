package handler

import (
	"github.com/gin-gonic/gin"
)

// Register all the different routes
func RegisterRoutes(router *gin.Engine) {
	UserRoutes(router)
	CommunityRoutes(router)
	PostRoutes(router)
}
