package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sublink/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func GetIp(c *gin.Context) {
	c.Next()
	func() {
		subname, _ := c.Get("subname")
		username, _ := c.Get("username")
		userID, _ := c.Get("userId")
		ip := c.ClientIP()
		resp, err := http.Get(fmt.Sprintf("https://whois.pconline.com.cn/ipJson.jsp?ip=%s&json=true", ip))
		if err != nil {
			log.Println(err)
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
		err = json.Unmarshal(utf8Body, &ipinfo)
		if err != nil {
			log.Println(err)
			return
		}
		region := ipinfo.Addr
		if region == "" {
			region = "未知"
		}
		var sub models.Subcription
		if subname, ok := subname.(string); ok {
			sub.Name = subname
		}
		err = sub.Find()
		if err != nil {
			log.Println(err)
			return
		}
		var iplog models.SubLogs
		iplog.IP = ip
		if uid, ok := userID.(int); ok {
			iplog.UserID = uid
		}
		err = iplog.Find(sub.ID)
		if err != nil {
			iplog.SubcriptionID = sub.ID
			iplog.Date = time.Now().Format("2006-01-02 15:04:05")
			iplog.Count = 1
			iplog.Addr = ipinfo.Addr
			iplog.Region = region
			iplog.Status = "success"
			if uname, ok := username.(string); ok {
				iplog.Username = uname
			}
			iplog.Client = c.Query("client")
			if err = iplog.Add(); err != nil {
				log.Println(err)
			}
		} else {
			iplog.Count++
			iplog.Date = time.Now().Format("2006-01-02 15:04:05")
			iplog.Addr = ipinfo.Addr
			iplog.Region = region
			iplog.Client = c.Query("client")
			iplog.Status = "success"
			if err = iplog.Update(); err != nil {
				log.Println(err)
			}
		}
	}()

}
