package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mingkwind/apilimiter"
)

func main() {
	//初始化令牌桶
	bucket := &apilimiter.LeakyBucket{
		Max:   100,
		Cycle: 10,
	}
	//初始化令牌桶全局限流器
	bucket.NewLeakyBucket()
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		//获取令牌
		token := bucket.GetToken(1)
		if token {
			c.JSON(200, gin.H{
				"message": "ok",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "no",
			})
		}
	})
	r.Run(":8080")
}
