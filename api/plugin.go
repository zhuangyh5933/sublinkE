package api

import (
	"github.com/eun1e/sublinkE-plugins"
	"github.com/gin-gonic/gin"
)

// PluginListResponse 插件列表响应
type PluginListResponse struct {
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	FilePath    string                 `json:"filePath"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// PluginConfigRequest 插件配置请求
type PluginConfigRequest struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// GetPlugins 获取所有插件列表
func GetPlugins(c *gin.Context) {
	manager := plugins.GetManager()
	allPlugins := manager.GetAllPlugins()

	var response []PluginListResponse
	for _, plugin := range allPlugins {
		response = append(response, PluginListResponse{
			Name:        plugin.Name,
			Version:     plugin.Version,
			Description: plugin.Description,
			Enabled:     plugin.Enabled,
			FilePath:    plugin.FilePath,
			Config:      plugin.Config,
		})
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"data": response,
		"msg":  "获取插件列表成功",
	})
}

// EnablePlugin 启用插件
func EnablePlugin(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(400, gin.H{
			"code": "40000",
			"msg":  "插件名称不能为空",
		})
		return
	}

	manager := plugins.GetManager()
	if err := manager.EnablePlugin(pluginName); err != nil {
		c.JSON(400, gin.H{
			"code": "40000",
			"msg":  "启用插件失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "插件启用成功",
	})
}

// DisablePlugin 禁用插件
func DisablePlugin(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(400, gin.H{
			"code": "40000",
			"msg":  "插件名称不能为空",
		})
		return
	}

	manager := plugins.GetManager()
	if err := manager.DisablePlugin(pluginName); err != nil {
		c.JSON(400, gin.H{
			"code": "40000",
			"msg":  "禁用插件失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "插件禁用成功",
	})
}

// UpdatePluginConfig 更新插件配置
func UpdatePluginConfig(c *gin.Context) {
	var req PluginConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": "40000",
			"msg":  "请求参数错误: " + err.Error(),
		})
		return
	}

	manager := plugins.GetManager()
	if err := manager.UpdatePluginConfig(req.Name, req.Config); err != nil {
		c.JSON(400, gin.H{
			"code": "40000",
			"msg":  "更新插件配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "插件配置更新成功",
	})
}

// GetPluginConfig 获取插件配置
func GetPluginConfig(c *gin.Context) {
	pluginName := c.Param("name")
	if pluginName == "" {
		c.JSON(400, gin.H{
			"code": "40000",
			"msg":  "插件名称不能为空",
		})
		return
	}

	manager := plugins.GetManager()
	plugin, exists := manager.GetPlugin(pluginName)
	if !exists {
		c.JSON(404, gin.H{
			"code": "40400",
			"msg":  "插件不存在",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"data": plugin.Config,
		"msg":  "获取插件配置成功",
	})
}

// ReloadPlugins 重新加载插件
func ReloadPlugins(c *gin.Context) {
	manager := plugins.GetManager()

	// 关闭所有插件
	manager.Shutdown()

	// 重新加载插件
	if err := manager.LoadPlugins(); err != nil {
		c.JSON(500, gin.H{
			"code": "50000",
			"msg":  "重新加载插件失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "插件重新加载成功",
	})
}
