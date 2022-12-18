/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/makeAccount/controller/Library"
	"github.com/wangyi/makeAccount/controller/account"
	"github.com/wangyi/makeAccount/controller/hash"
	"github.com/wangyi/makeAccount/logger"
	"log"
	"net/http"
)

/**
注册路由
*/

func Setup() *gin.Engine {
	r := gin.New()
	//添加记录日志的中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true), Cors())
	//r.Use(Cors())
	r.Static("/static", "./static")

	r.GET("/get", account.Account)
	r.GET("/getTwo", account.Account)

	r.POST("/setWordLibrary", Library.SetWordsLibrary)
	//DownWordWords
	r.GET("/downWordWords", Library.DownWordWords)
	//UploadWords
	r.POST("/uploadWords", Library.UploadWords)
	//SetThePhrase
	r.POST("/setThePhrase", Library.SetThePhrase)

	r.POST("/setWordLibraryV2", Library.SetWordsLibraryV2)
	//DownWordWords
	r.GET("/downWordWordsV2", Library.DownWordWordsV2)
	//UploadWords
	r.POST("/uploadWordsV2", Library.UploadWordsV2)
	//SetThePhrase
	r.POST("/setThePhraseV2", Library.SetThePhraseV2)

	//GetHash
	r.POST("/GetHash", hash.GetHash)
	r.GET("/GetHash", hash.GetHash)
	r.GET("/DownHash", hash.DownHash)
	_ = r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))

	return r
}

/**
跨域设置
*/
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
