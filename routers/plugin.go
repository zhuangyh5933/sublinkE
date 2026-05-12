package routers

import (
	"sublink/api"

	"github.com/gin-gonic/gin"
)

func Plugins(r *gin.Engine) {
	pluginGroup := r.Group("/api/v1/plugins")
	pluginGroup.Use(api.AdminOnlyMiddleware)
	{
		pluginGroup.GET("/get", api.GetPlugins)               // 获取所有插件
		pluginGroup.POST("/enable/:name", api.EnablePlugin)   // 启用插件
		pluginGroup.POST("/disable/:name", api.DisablePlugin) // 禁用插件
		pluginGroup.PUT("/config", api.UpdatePluginConfig)    // 更新插件配置
		pluginGroup.GET("/config/:name", api.GetPluginConfig) // 获取插件配置
		pluginGroup.POST("/reload", api.ReloadPlugins)        // 重新加载插件
	}
}
