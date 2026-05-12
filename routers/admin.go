package routers

import (
	"sublink/api"

	"github.com/gin-gonic/gin"
)

func Admin(r *gin.Engine) {
	adminGroup := r.Group("/api/v1/admin")
	{
		adminGroup.GET("/users", api.AdminUsers)
		adminGroup.POST("/users/update", api.AdminUpdateUser)
		adminGroup.POST("/users/reset-subscription-token", api.AdminResetSubscriptionToken)
		adminGroup.DELETE("/users/:id", api.AdminDeleteUser)
		adminGroup.GET("/invites", api.AdminInviteList)
		adminGroup.POST("/invites/add", api.AdminInviteAdd)
		adminGroup.POST("/invites/update", api.AdminInviteUpdate)
		adminGroup.GET("/config", api.AdminConfigGet)
		adminGroup.POST("/config", api.AdminConfigSet)
	}
}
