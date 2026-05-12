package api

import (
	"fmt"
	"log"
	"strings"
	"sublink/models"
	"sublink/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func currentUser(c *gin.Context) (*models.User, error) {
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	if err := user.Find(); err != nil {
		return nil, err
	}
	return user, nil
}

func ensureAdmin(c *gin.Context) (*models.User, bool) {
	user, err := currentUser(c)
	if err != nil {
		c.JSON(401, gin.H{"msg": "未登录"})
		c.Abort()
		return nil, false
	}
	if user.Role != "admin" {
		c.JSON(403, gin.H{"msg": "仅管理员可操作"})
		c.Abort()
		return nil, false
	}
	return user, true
}

func AdminOnlyMiddleware(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	c.Next()
}

func AdminUsers(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	users, err := new(models.User).All()
	if err != nil {
		c.JSON(500, gin.H{"msg": "获取用户失败"})
		return
	}
	list := []gin.H{}
	for _, item := range users {
		logs, _ := new(models.SubLogs).ListByUser(item.ID)
		list = append(list, gin.H{
			"id":                item.ID,
			"username":          item.Username,
			"nickname":          item.Nickname,
			"role":              item.Role,
			"subscriptionId":    item.SubscriptionID,
			"subscriptionToken": item.SubscriptionToken,
			"allowedRegions":    item.AllowedRegions,
			"inviteCodeId":      item.InviteCodeID,
			"disabled":          item.Disabled,
			"createdAt":         item.CreatedAt,
			"pullLogs":          logs,
		})
	}
	c.JSON(200, gin.H{"code": "00000", "data": list, "msg": "获取成功"})
}

func AdminUpdateUser(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	id := c.PostForm("id")
	if id == "" {
		c.JSON(400, gin.H{"msg": "用户ID不能为空"})
		return
	}
	var target models.User
	fmt.Sscanf(id, "%d", &target.ID)
	if err := target.FindByID(); err != nil {
		c.JSON(400, gin.H{"msg": "用户不存在"})
		return
	}
	updates := &models.User{
		AllowedRegions: c.PostForm("allowedRegions"),
	}
	if subscriptionID := c.PostForm("subscriptionId"); subscriptionID != "" {
		fmt.Sscanf(subscriptionID, "%d", &updates.SubscriptionID)
	}
	role := c.PostForm("role")
	if role != "" {
		updates.Role = role
	}
	if token := c.PostForm("subscriptionToken"); token != "" {
		updates.SubscriptionToken = token
	}
	payload := map[string]any{}
	if updates.SubscriptionID != 0 {
		payload["subscription_id"] = updates.SubscriptionID
	}
	payload["allowed_regions"] = updates.AllowedRegions
	if role != "" {
		payload["role"] = updates.Role
	}
	if updates.SubscriptionToken != "" {
		payload["subscription_token"] = updates.SubscriptionToken
	}
	if disabled := c.PostForm("disabled"); disabled != "" {
		payload["disabled"] = disabled == "true"
	}
	if target.Role == "admin" {
		if value, ok := payload["role"]; ok && value == "user" {
			c.JSON(400, gin.H{"msg": "不能直接降级当前管理员接口用户，请使用其他管理员账户处理"})
			return
		}
		if value, ok := payload["disabled"]; ok && value == true {
			c.JSON(400, gin.H{"msg": "不能禁用管理员账户"})
			return
		}
	}
	if len(payload) == 0 {
		c.JSON(400, gin.H{"msg": "没有可更新字段"})
		return
	}
	if err := models.DB.Model(&models.User{}).Where("id = ?", target.ID).Updates(payload).Error; err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"msg": "更新失败"})
		return
	}
	c.JSON(200, gin.H{"code": "00000", "msg": "更新成功"})
}

func AdminInviteList(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	invites, err := new(models.InviteCode).List()
	if err != nil {
		c.JSON(500, gin.H{"msg": "获取邀请码失败"})
		return
	}
	c.JSON(200, gin.H{"code": "00000", "data": invites, "msg": "获取成功"})
}

func AdminInviteAdd(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	code := c.PostForm("code")
	description := c.PostForm("description")
	if code == "" {
		code = strings.ToUpper(utils.RandString(10))
	}
	invite := &models.InviteCode{
		Code:        code,
		Description: description,
		Enabled:     true,
	}
	if err := invite.Add(); err != nil {
		c.JSON(400, gin.H{"msg": "创建邀请码失败，可能已存在"})
		return
	}
	c.JSON(200, gin.H{"code": "00000", "data": invite, "msg": "创建成功"})
}

func AdminInviteUpdate(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	var invite models.InviteCode
	id := c.PostForm("id")
	fmt.Sscanf(id, "%d", &invite.ID)
	if invite.ID == 0 {
		c.JSON(400, gin.H{"msg": "邀请码ID不能为空"})
		return
	}
	if err := invite.FindByID(); err != nil {
		c.JSON(400, gin.H{"msg": "邀请码不存在"})
		return
	}
	payload := map[string]any{}
	if description := c.PostForm("description"); description != "" {
		payload["description"] = description
	}
	if enabled := c.PostForm("enabled"); enabled != "" {
		payload["enabled"] = enabled == "true"
	}
	if len(payload) == 0 {
		c.JSON(400, gin.H{"msg": "没有可更新字段"})
		return
	}
	if err := models.DB.Model(&models.InviteCode{}).Where("id = ?", invite.ID).Updates(payload).Error; err != nil {
		c.JSON(500, gin.H{"msg": "更新失败"})
		return
	}
	c.JSON(200, gin.H{"code": "00000", "msg": "更新成功"})
}

func AdminConfigGet(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	config := models.ReadConfig()
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"jwt_secret":              config.JwtSecret,
			"api_encryption_key":      config.APIEncryptionKey,
			"expire_days":             config.ExpireDays,
			"port":                    config.Port,
			"default_subscription_id": config.DefaultSubscriptionID,
			"invite_required":         config.InviteRequired,
		},
		"msg": "获取成功",
	})
}

func AdminConfigSet(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	config := models.ReadConfig()
	if defaultSubscriptionID := c.PostForm("defaultSubscriptionId"); defaultSubscriptionID != "" {
		fmt.Sscanf(defaultSubscriptionID, "%d", &config.DefaultSubscriptionID)
	}
	if inviteRequired := c.PostForm("inviteRequired"); inviteRequired != "" {
		config.InviteRequired = inviteRequired == "true"
	}
	models.SetConfig(config)
	c.JSON(200, gin.H{"code": "00000", "msg": "保存成功"})
}

func AdminResetSubscriptionToken(c *gin.Context) {
	if _, ok := ensureAdmin(c); !ok {
		return
	}
	id := c.PostForm("id")
	var user models.User
	fmt.Sscanf(id, "%d", &user.ID)
	if user.ID == 0 {
		c.JSON(400, gin.H{"msg": "用户ID不能为空"})
		return
	}
	if err := user.FindByID(); err != nil {
		c.JSON(400, gin.H{"msg": "用户不存在"})
		return
	}
	newToken := fmt.Sprintf("sub_%d_%s", time.Now().UnixNano(), utils.RandString(24))
	if err := models.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("subscription_token", newToken).Error; err != nil {
		c.JSON(500, gin.H{"msg": "重置订阅Token失败"})
		return
	}
	c.JSON(200, gin.H{"code": "00000", "data": gin.H{"subscriptionToken": newToken}, "msg": "重置成功"})
}

func AdminDeleteUser(c *gin.Context) {
	current, ok := ensureAdmin(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var user models.User
	fmt.Sscanf(id, "%d", &user.ID)
	if user.ID == 0 {
		c.JSON(400, gin.H{"msg": "用户ID不能为空"})
		return
	}
	if err := user.FindByID(); err != nil {
		c.JSON(400, gin.H{"msg": "用户不存在"})
		return
	}
	if current.ID == user.ID {
		c.JSON(400, gin.H{"msg": "不能删除当前登录管理员"})
		return
	}
	if user.Role == "admin" {
		c.JSON(400, gin.H{"msg": "不能删除管理员账户"})
		return
	}
	_ = models.DB.Where("user_id = ?", user.ID).Delete(&models.SubLogs{}).Error
	_ = models.DB.Where("user_id = ?", user.ID).Delete(&models.AccessKey{}).Error
	if err := user.Del(); err != nil {
		c.JSON(500, gin.H{"msg": "删除用户失败"})
		return
	}
	c.JSON(200, gin.H{"code": "00000", "msg": "删除成功"})
}
