package middlewares

import (
	"bytes"
	"github.com/eun1e/sublinkE-plugins"
	"io"

	"github.com/gin-gonic/gin"
)

// responseWriter 用于捕获响应
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// PluginMiddleware 插件中间件
func PluginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取插件管理器
		manager := plugins.GetManager()

		// 触发API调用前事件
		manager.TriggerEvent(c, plugins.EventAPIBefore, c.Request.URL.Path, 0, nil, nil)

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// 复原请求体，让后续能正常读
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 创建response writer来捕获响应
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = responseWriter

		// 继续处理请求
		c.Next()

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 获取响应体
		var responseBody interface{}
		if responseWriter.body.Len() > 0 {
			responseBody = responseWriter.body.String()
		}

		// 触发API调用后事件
		manager.TriggerEvent(c, plugins.EventAPIAfter, c.Request.URL.Path, statusCode, requestBody, responseBody)

		// 根据状态码触发成功或错误事件
		if statusCode >= 200 && statusCode < 300 {
			manager.TriggerEvent(c, plugins.EventAPISuccess, c.Request.URL.Path, statusCode, requestBody, responseBody)
		} else if statusCode >= 400 {
			manager.TriggerEvent(c, plugins.EventAPIError, c.Request.URL.Path, statusCode, requestBody, responseBody)
		}
	}
}
