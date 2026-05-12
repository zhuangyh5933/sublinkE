package api

import (
	"sublink/models"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Title     string   `json:"title"`
	Icon      string   `json:"icon"`
	Hidden    bool     `json:"hidden"`
	Roles     []string `json:"roles"`
	KeepAlive bool     `json:"keepAlive,omitempty"`
}

type Child struct {
	Path      string `json:"path"`
	Component string `json:"component"`
	Name      string `json:"name"`
	Meta      Meta   `json:"meta"`
}

type Menu struct {
	Path      string  `json:"path"`
	Component string  `json:"component"`
	Redirect  string  `json:"redirect"`
	Name      string  `json:"name"`
	Meta      Meta    `json:"meta"`
	Children  []Child `json:"children"`
}

func GetMenus(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"msg": "未登录"})
		return
	}
	user := &models.User{Username: username.(string)}
	if err := user.Find(); err != nil {
		c.JSON(401, gin.H{"msg": "用户不存在"})
		return
	}
	if user.Role == "admin" {
		menus := []Menu{{
			Path:      "/system",
			Component: "Layout",
			Name:      "system",
			Meta: Meta{Title: "system", Icon: "system", Hidden: false, Roles: []string{"ADMIN"}},
			Children: []Child{
				{Path: "user/set", Component: "system/user/set", Name: "Userset", Meta: Meta{Title: "userset", Icon: "role", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}},
				{Path: "user/admin", Component: "system/user/admin", Name: "UserAdmin", Meta: Meta{Title: "用户管理", Icon: "user", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}},
				{Path: "invite/index", Component: "system/invite/index", Name: "InviteAdmin", Meta: Meta{Title: "邀请码", Icon: "key", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}},
				{Path: "config/index", Component: "system/config/index", Name: "SystemConfig", Meta: Meta{Title: "注册配置", Icon: "setting", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}},
			},
		},
			{Path: "/apikey", Component: "Layout", Name: "apikey", Meta: Meta{Title: "apikey", Icon: "key", Hidden: true, Roles: []string{"ADMIN"}}, Children: []Child{{Path: "index", Component: "apikey/index", Name: "ApikeyIndex", Meta: Meta{Title: "apikey", Icon: "key", Hidden: true, Roles: []string{"ADMIN"}, KeepAlive: true}}}},
			{Path: "/subcription", Component: "Layout", Redirect: "/subcription/subs", Name: "subcription", Meta: Meta{Title: "subcription", Icon: "client", Hidden: false, Roles: []string{"ADMIN"}}, Children: []Child{{Path: "subs", Component: "subcription/subs", Name: "Subs", Meta: Meta{Title: "sublist", Icon: "link", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}}, {Path: "nodes", Component: "subcription/nodes", Name: "Nodes", Meta: Meta{Title: "nodelist", Icon: "publish", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}}, {Path: "template", Component: "subcription/template", Name: "Template", Meta: Meta{Title: "templatelist", Icon: "document", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}}}},
			{Path: "/plugin", Component: "Layout", Redirect: "/plugin/index", Name: "plugin", Meta: Meta{Title: "plugin", Icon: "system", Hidden: false, Roles: []string{"ADMIN"}}, Children: []Child{{Path: "index", Component: "plugin/index", Name: "PluginList", Meta: Meta{Title: "PluginList", Icon: "system", Hidden: false, Roles: []string{"ADMIN"}, KeepAlive: true}}}},
		}
		c.JSON(200, gin.H{"code": "00000", "data": menus, "msg": "获取成功"})
		return
	}
	menus := []Menu{{
		Path:      "/user",
		Component: "Layout",
		Redirect:  "/user/subscription",
		Name:      "user-panel",
		Meta: Meta{Title: "用户中心", Icon: "user", Hidden: false, Roles: []string{"USER"}},
		Children: []Child{
			{Path: "subscription", Component: "system/user/set", Name: "Userset", Meta: Meta{Title: "我的订阅", Icon: "link", Hidden: false, Roles: []string{"USER"}, KeepAlive: true}},
			{Path: "apikey", Component: "apikey/index", Name: "ApikeyIndex", Meta: Meta{Title: "APIKey", Icon: "key", Hidden: false, Roles: []string{"USER"}, KeepAlive: true}},
		},
	}}
	c.JSON(200, gin.H{"code": "00000", "data": menus, "msg": "获取成功"})
}
