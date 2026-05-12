package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sublink/models"
	"sublink/node"
	"sublink/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var SunName string

// md5加密
func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func resolveSubscriptionName(token string) (string, *models.User, error) {
	user := &models.User{SubscriptionToken: token}
	if err := user.FindBySubscriptionToken(); err == nil && user.SubscriptionID != 0 {
		var sub models.Subcription
		sub.ID = user.SubscriptionID
		if err = sub.Find(); err == nil {
			return sub.Name, user, nil
		}
	}
	Sub := new(models.Subcription)
	list, err := Sub.List()
	if err != nil {
		return "", nil, err
	}
	for _, sub := range list {
		if Md5(sub.Name) == strings.ToLower(token) {
			return sub.Name, nil, nil
		}
	}
	return "", nil, fmt.Errorf("找不到对应订阅")
}
func GetClient(c *gin.Context) {
	token := c.Query("token")
	ClientIndex := c.Query("client") // 客户端标识
	if token == "" {
		log.Println("token为空")
		c.Writer.WriteString("token为空")
		return
	}
	subName, user, err := resolveSubscriptionName(token)
	if err != nil {
		log.Println(err)
		c.Writer.WriteString(err.Error())
		return
	}
	c.Set("subname", subName)
	if user != nil {
		c.Set("userId", user.ID)
		c.Set("username", user.Username)
		c.Set("allowedRegions", user.AllowedRegions)
		if strings.TrimSpace(user.AllowedRegions) != "" {
			ip := c.ClientIP()
			resp, err := http.Get(fmt.Sprintf("https://whois.pconline.com.cn/ipJson.jsp?ip=%s&json=true", ip))
			if err != nil {
				log.Println(err)
				c.String(http.StatusForbidden, "当前地区不允许拉取订阅")
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			utf8Body, _ := simplifiedchinese.GBK.NewDecoder().Bytes(body)
			type IpInfo struct {
				Addr string `json:"addr"`
				Ip   string `json:"ip"`
			}
			ipinfo := IpInfo{}
			if err = json.Unmarshal(utf8Body, &ipinfo); err != nil {
				log.Println(err)
				c.String(http.StatusForbidden, "当前地区不允许拉取订阅")
				return
			}
			region := ipinfo.Addr
			matched := false
			for _, item := range strings.Split(user.AllowedRegions, ",") {
				if strings.Contains(region, strings.TrimSpace(item)) {
					matched = true
					break
				}
			}
			if !matched {
				var sub models.Subcription
				sub.Name = subName
				_ = sub.Find()
				failedLog := &models.SubLogs{
					IP:            ip,
					Date:          time.Now().Format("2006-01-02 15:04:05"),
					Addr:          ipinfo.Addr,
					Region:        region,
					Client:        ClientIndex,
					Status:        "blocked_region",
					Count:         1,
					SubcriptionID: sub.ID,
					UserID:        user.ID,
					Username:      user.Username,
				}
				_ = failedLog.Add()
				c.String(http.StatusForbidden, "当前地区不允许拉取订阅")
				return
			}
		}
	}
	SunName = subName
	switch ClientIndex {
	case "clash":
		GetClash(c)
		return
	case "surge":
		GetSurge(c)
		return
	case "v2ray":
		GetV2ray(c)
		return
	}
	ClientList := []string{"clash", "surge"}
	for k, v := range c.Request.Header {
		if k == "User-Agent" {
			for _, UserAgent := range v {
				if UserAgent == "" {
					fmt.Println("User-Agent为空")
				}
				for _, client := range ClientList {
					if strings.Contains(strings.ToLower(UserAgent), strings.ToLower(client)) {
						switch client {
						case "clash":
							GetClash(c)
							return
						case "surge":
							GetSurge(c)
							return
						default:
							fmt.Println("未知客户端")
						}
					}
				}
				GetV2ray(c)
				return
			}
		}
	}
	GetV2ray(c)
}
func GetV2ray(c *gin.Context) {
	var sub models.Subcription
	if SunName == "" {
		c.Writer.WriteString("订阅名为空")
		return
	}
	// subname := c.Param("subname")
	// subname := SunName
	// subname = node.Base64Decode(subname)
	sub.Name = SunName
	err := sub.Find()
	if err != nil {
		c.Writer.WriteString("找不到这个订阅:" + SunName)
		return
	}
	err = sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取错误")
		return
	}
	baselist := ""
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			baselist += strings.Join(links, "\n") + "\n"
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := utils.Base64Decode(string(body))
			baselist += nodes + "\n"
		// 默认
		default:
			baselist += v.Link + "\n"
		}
	}
	c.Set("subname", SunName)
	filename := fmt.Sprintf("%s.txt", SunName)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteString(utils.Base64Encode(baselist))
}
func GetClash(c *gin.Context) {
	var sub models.Subcription
	// subname := c.Param("subname")
	// subname := node.Base64Decode(SunName)
	sub.Name = SunName
	err := sub.Find()
	if err != nil {
		c.Writer.WriteString("找不到这个订阅:" + SunName)
		return
	}
	err = sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取错误")
		return
	}
	var urls []node.Urls
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			for _, link := range links {
				urls = append(urls, node.Urls{
					Url:             link,
					DialerProxyName: strings.TrimSpace(v.DialerProxyName),
				})
			}
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := utils.Base64Decode(string(body))
			links := strings.Split(nodes, "\n")
			for _, link := range links {
				urls = append(urls, node.Urls{
					Url:             link,
					DialerProxyName: strings.TrimSpace(v.DialerProxyName),
				})
			}
		// 默认
		default:
			urls = append(urls, node.Urls{
				Url:             v.Link,
				DialerProxyName: strings.TrimSpace(v.DialerProxyName),
			})
		}
	}

	var configs utils.SqlConfig
	err = json.Unmarshal([]byte(sub.Config), &configs)
	if err != nil {
		c.Writer.WriteString("配置读取错误")
		return
	}
	DecodeClash, err := node.EncodeClash(urls, configs)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}
	c.Set("subname", SunName)
	filename := fmt.Sprintf("%s.yaml", SunName)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteString(string(DecodeClash))
}
func GetSurge(c *gin.Context) {
	var sub models.Subcription
	// subname := c.Param("subname")
	// subname := node.Base64Decode(SunName)
	sub.Name = SunName
	err := sub.Find()
	if err != nil {
		c.Writer.WriteString("找不到这个订阅:" + SunName)
		return
	}
	err = sub.GetSub()
	if err != nil {
		c.Writer.WriteString("读取错误")
		return
	}
	urls := []string{}
	for _, v := range sub.Nodes {
		switch {
		// 如果包含多条节点
		case strings.Contains(v.Link, ","):
			links := strings.Split(v.Link, ",")
			urls = append(urls, links...)
			continue
		//如果是订阅转换
		case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
			resp, err := http.Get(v.Link)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			nodes := utils.Base64Decode(string(body))
			links := strings.Split(nodes, "\n")
			urls = append(urls, links...)
		// 默认
		default:
			urls = append(urls, v.Link)
		}
	}

	var configs utils.SqlConfig
	err = json.Unmarshal([]byte(sub.Config), &configs)
	if err != nil {
		c.Writer.WriteString("配置读取错误")
		return
	}
	// log.Println("surge路径:", configs)
	DecodeClash, err := node.EncodeSurge(urls, configs)
	if err != nil {
		c.Writer.WriteString(err.Error())
		return
	}
	c.Set("subname", SunName)
	filename := fmt.Sprintf("%s.conf", SunName)
	encodedFilename := url.QueryEscape(filename)
	c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	host := c.Request.Host
	url := c.Request.URL.String()
	// 如果包含头部更新信息
	if strings.Contains(DecodeClash, "#!MANAGED-CONFIG") {
		c.Writer.WriteString(DecodeClash)
		return
	}
	// 否则就插入头部更新信息
	interval := fmt.Sprintf("#!MANAGED-CONFIG %s interval=86400 strict=false", host+url)
	c.Writer.WriteString(string(interval + "\n" + DecodeClash))
}
