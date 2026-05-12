package models

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// Plugin 插件数据模型
type Plugin struct {
	Name    string `gorm:"unique;not null"` // 插件名称，唯一
	Path    string `gorm:"not null"`        // 插件路径
	Enabled bool   `gorm:"default:false"`   // 是否启用
	Config  string `gorm:"type:text"`       // 插件配置，JSON格式存储
}

// GetPlugin 根据插件名获取单个插件
func GetPlugin(path string) (*Plugin, error) {
	var plugin Plugin
	result := DB.Where("path = ?", path).First(&plugin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("插件不存在")
		}
		return nil, result.Error
	}
	return &plugin, nil
}

// SavePlugin 保存或更新插件
// 修复数据库 SavePlugin 函数
func SavePlugin(name, filePath string, enabled bool, config map[string]interface{}) error {
	// 将配置转为JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// 检查插件是否已存在
	var plugin Plugin
	result := DB.Where("Path = ?", filePath).First(&plugin)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 创建新记录
			plugin = Plugin{
				Name:    name,
				Path:    filePath,
				Enabled: enabled,
				Config:  string(configJSON),
			}
			return DB.Create(&plugin).Error
		}
		return result.Error
	}

	// 更新现有记录 - 使用 ID 作为条件
	return DB.Model(&Plugin{}).Where("Path = ?", filePath).Updates(map[string]interface{}{
		"name":    name,
		"enabled": enabled,
		"config":  string(configJSON),
	}).Error
}
