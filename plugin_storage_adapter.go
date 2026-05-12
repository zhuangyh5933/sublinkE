package main

import (
	"github.com/eun1e/sublinkE-plugins"
	"sublink/models"
)

// ModelStorageAdapter 将 models 包适配为 plugins.PluginStorage 接口
type ModelStorageAdapter struct{}

func (m *ModelStorageAdapter) GetPlugin(path string) (*plugins.PluginStorageInfo, error) {
	plugin, err := models.GetPlugin(path)
	if err != nil {
		return nil, err
	}

	if plugin == nil {
		return nil, nil
	}

	return &plugins.PluginStorageInfo{
		Name:    plugin.Name,
		Path:    plugin.Path,
		Enabled: plugin.Enabled,
		Config:  plugin.Config,
	}, nil
}

func (m *ModelStorageAdapter) SavePlugin(name, path string, enabled bool, config map[string]interface{}) error {
	return models.SavePlugin(name, path, enabled, config)
}

// 在主程序启动时调用这个函数来设置存储适配器
func initPluginStorage() {
	plugins.SetStorage(&ModelStorageAdapter{})
}
