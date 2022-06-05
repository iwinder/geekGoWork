package week04

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"math/rand"
	"net/http"
	"time"
)

// InitRouter 创建 Http 路由
func InitRouter(g *gin.Engine, serverOut chan struct{}) *gin.Engine {
	//router := gin.Default()
	g.GET("/", func(c *gin.Context) {
		glog.V(2).Infoln("Entering route Handler...")
		// 获取参数
		user := c.Query("name")
		// 生成随机数
		delay := randInt(0, 2000)
		// 睡眠随机时间，模拟延时
		time.Sleep(time.Millisecond * time.Duration(delay))
		c.String(http.StatusOK, fmt.Sprintf("Welcome Gin Server [%s]\n", user))
	})

	g.GET("/shutdown", func(c *gin.Context) {
		glog.V(2).Infoln("Entering shutdown Handler...")
		serverOut <- struct{}{}
	})
	g.GET("/healthz", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	g.GET("/timeout", func(c *gin.Context) {
		// 业务开始
		glog.V(2).Infoln("Entering timeOut Handler...")
		time.Sleep(20 * time.Second)
		c.String(http.StatusOK, "Do something 20s OK\n")
	})
	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("/add", func(c *gin.Context) {
				// 如果有需要 通过 c 获取参数，调用 userServer
				// 可以单独创建文件存储 http 相关请求和返回
				// 可以根据错误不同的参数，返回不同的状态码
				c.String(http.StatusOK, "Welcome Gin Server")
			})
			userv1.GET("/:id", func(c *gin.Context) {
				// 如果有需要 通过 c 获取参数，调用 userServer
				// 可以单独创建文件存储 http 相关请求和返回
				// 可以根据错误不同的参数，返回不同的状态码
				time.Sleep(5 * time.Second)
				c.String(http.StatusOK, "Welcome Gin Server")
			})
		}
	}
	return g
}

// randInt 生成min到max的随机数
func randInt(min int, max int) int {
	// rand.Seed通过设置不同值，保证每次执行时生成的随机数不一样
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
