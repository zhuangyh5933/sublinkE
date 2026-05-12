package api

import (
	"fmt"
	"log"
	"net/url"
	"sublink/models"
	"sublink/utils"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID                 int
	Username           string
	Nickname           string
	Avatar             string
	Mobile             string
	Email              string
	SubscriptionID     int
	SubscriptionToken  string
	SubscriptionName   string
	SubscriptionURL    string
	ClashSubscription  string
	SurgeSubscription  string
	V2raySubscription  string
	AllowedRegions     string
}

func buildSubscriptionURLs(c *gin.Context, user *models.User) (string, string, string, string, string) {
	subscriptionURL := ""
	clashURL := ""
	surgeURL := ""
	v2rayURL := ""
	subscriptionName := ""
	if user.SubscriptionToken == "" {
		return subscriptionURL, clashURL, surgeURL, v2rayURL, subscriptionName
	}
	var sub models.Subcription
	sub.ID = user.SubscriptionID
	if err := sub.Find(); err == nil {
		subscriptionName = sub.Name
	}
	baseURL := fmt.Sprintf("http://%s", c.Request.Host)
	if c.Request.TLS != nil {
		baseURL = fmt.Sprintf("https://%s", c.Request.Host)
	}
	token := url.QueryEscape(user.SubscriptionToken)
	subscriptionURL = fmt.Sprintf("%s/c/?token=%s", baseURL, token)
	clashURL = subscriptionURL + "&client=clash"
	surgeURL = subscriptionURL + "&client=surge"
	v2rayURL = subscriptionURL + "&client=v2ray"
	return subscriptionURL, clashURL, surgeURL, v2rayURL, subscriptionName
}

// 注册用户
func UserRegister(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")
	inviteCode := c.PostForm("inviteCode")
	captchaCode := c.PostForm("captchaCode")
	captchaKey := c.PostForm("captchaKey")
	if username == "" || password == "" {
		c.JSON(400, gin.H{
			"msg": "用户名或密码不能为空",
		})
		return
	}
	exists := &models.User{Username: username}
	if err := exists.Find(); err == nil {
		c.JSON(400, gin.H{
			"msg": "用户名已存在",
		})
		return
	}
	config := models.ReadConfig()
	invite := &models.InviteCode{}
	if config.InviteRequired {
		if inviteCode == "" {
			c.JSON(400, gin.H{
				"msg": "当前注册需要邀请码",
			})
			return
		}
		invite.Code = inviteCode
		if err := invite.FindByCode(); err != nil || !invite.Enabled {
			c.JSON(400, gin.H{
				"msg": "邀请码无效或已禁用",
			})
			return
		}
	}
	if !utils.VerifyCaptcha(captchaKey, captchaCode) {
		c.JSON(400, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	user := &models.User{
		Username:       username,
		Password:       password,
		Nickname:       nickname,
		Role:           "user",
		SubscriptionID: config.DefaultSubscriptionID,
	}
	if user.SubscriptionID != 0 {
		var sub models.Subcription
		sub.ID = user.SubscriptionID
		if err := sub.Find(); err != nil {
			c.JSON(400, gin.H{
				"msg": "默认订阅不存在，请先配置 default_subscription_id",
			})
			return
		}
	}
	if err := user.Create(); err != nil {
		log.Println("创建用户失败", err)
		c.JSON(400, gin.H{
			"msg": "创建用户失败",
		})
		return
	}
	if config.InviteRequired && invite.ID != 0 {
		if err := models.DB.Model(&models.InviteCode{}).Where("id = ?", invite.ID).Update("used_count", invite.UsedCount+1).Error; err == nil {
			invite.UsedCount++
		}
		_ = user.Set(&models.User{InviteCodeID: invite.ID})
		user.InviteCodeID = invite.ID
	}
	subscriptionURL, clashURL, surgeURL, v2rayURL, subscriptionName := buildSubscriptionURLs(c, user)
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"username":          user.Username,
			"nickname":          user.Nickname,
			"subscriptionId":    user.SubscriptionID,
			"subscriptionName":  subscriptionName,
			"subscriptionToken": user.SubscriptionToken,
			"subscriptionUrl":   subscriptionURL,
			"clashUrl":          clashURL,
			"surgeUrl":          surgeURL,
			"v2rayUrl":          v2rayURL,
			"allowedRegions":    user.AllowedRegions,
		},
		"msg": "注册成功",
	})
}

// 新增用户
func UserAdd(c *gin.Context) {
	user := &models.User{
		Username: "test",
		Password: "test123",
	}
	err := user.Create()
	if err != nil {
		log.Println("创建用户失败")
	}
	c.String(200, "创建用户成功")
}

// 获取用户信息
func UserMe(c *gin.Context) {
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	err := user.Find()
	if err != nil {
		c.JSON(400, gin.H{
			"code": "00000",
			"msg":  err,
		})
		return
	}
	subscriptionURL, clashURL, surgeURL, v2rayURL, subscriptionName := buildSubscriptionURLs(c, user)
	roles := []string{"USER"}
	if user.Role == "admin" {
		roles = []string{"ADMIN"}
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"avatar":            "static/avatar.gif",
			"nickname":          user.Nickname,
			"userId":            user.ID,
			"username":          user.Username,
			"roles":             roles,
			"subscriptionId":    user.SubscriptionID,
			"subscriptionName":  subscriptionName,
			"subscriptionToken": user.SubscriptionToken,
			"subscriptionUrl":   subscriptionURL,
			"clashUrl":          clashURL,
			"surgeUrl":          surgeURL,
			"v2rayUrl":          v2rayURL,
			"allowedRegions":    user.AllowedRegions,
		},
		"msg": "获取用户信息成功",
	})
}

// 获取所有用户
func UserPages(c *gin.Context) {
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	users, err := user.All()
	if err != nil {
		log.Println("获取用户信息失败")
	}
	list := []*User{}
	for i := range users {
		list = append(list, &User{
			ID:                users[i].ID,
			Username:          users[i].Username,
			Nickname:          users[i].Nickname,
			Avatar:            "static/avatar.gif",
			SubscriptionID:    users[i].SubscriptionID,
			SubscriptionToken: users[i].SubscriptionToken,
			AllowedRegions:    users[i].AllowedRegions,
		})
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": gin.H{
			"list": list,
		},
		"msg": "获取用户信息成功",
	})
}

// 更新用户信息
func UserSet(c *gin.Context) {
	NewUsername := c.PostForm("username")
	NewPassword := c.PostForm("password")
	if NewUsername == "" || NewPassword == "" {
		c.JSON(400, gin.H{
			"code": "00001",
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	err := user.Set(&models.User{
		Username: NewUsername,
		Password: NewPassword,
	})
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"code": "00000",
			"msg":  err,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "修改成功",
	})
}

func UserPullLogs(c *gin.Context) {
	username, _ := c.Get("username")
	user := &models.User{Username: username.(string)}
	if err := user.Find(); err != nil {
		c.JSON(400, gin.H{"msg": "用户不存在"})
		return
	}
	logs, err := new(models.SubLogs).ListByUser(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"msg": "获取拉取记录失败"})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": logs,
		"msg":  "获取成功",
	})
}
