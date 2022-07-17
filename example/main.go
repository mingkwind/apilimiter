package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mingkwind/apilimiter"
)

func main() {
	//初始化令牌桶
	bucket := &apilimiter.TokenBucket{
		Max:   100,
		Cycle: 100,
		Batch: 1,
	}
	//初始化令牌桶全局限流器
	bucket.NewTokenLimiter()
	r := gin.Default()

}
