package test

import (
	"fmt"
	"github/dpwgc/apilimiter"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

//令牌桶全局限流器测试
func TestTokenLimiter(t *testing.T) {

	fmt.Println("--- tokenLimiter ---")

	//设置令牌桶最大容量为100，每500毫秒生产50个令牌，相当于每1000毫秒最多只能取出100个令牌
	bucket := apilimiter.Bucket{
		Max:   100,
		Cycle: 500,
		Batch: 50,
	}

	//初始化令牌桶限流器
	bucket.NewTokenLimiter()
	sucNum := new(int64) //成功请求数
	*sucNum = 0
	//模拟200次循环请求
	for i := 0; i < 200; i++ {
		//每次访问至取出1个令牌
		isOk := bucket.GetToken(1)
		if isOk {
			*sucNum++
			// fmt.Println(i, "Access successful", "[Time]:", time.Now().Unix())
		} else {
			// fmt.Println(i, "Access failed.", "Token bucket is empty", "[Time]:", time.Now().Unix())
		}
	}
	if *sucNum > 100 {
		t.Errorf("loop request sucNum expected <= 100, got %d", *sucNum)
	}
	time.Sleep(time.Second * 1)

	//模拟200次并发请求
	wg := &sync.WaitGroup{}
	reqChan := make(chan int, 200)
	*sucNum = 0
	for i := 0; i < 200; i++ {
		go func(i int, sucNum *int64, wg *sync.WaitGroup) {
			<-reqChan
			isOk := bucket.GetToken(1)
			if isOk {
				atomic.AddInt64(sucNum, 1)
				// fmt.Println(i, "Access successful", "[Time]:", time.Now().Unix())
			} else {
				//fmt.Println(i, "Access failed.", "Token bucket is empty", "[Time]:", time.Now().Unix())
			}
			wg.Done()
		}(i, sucNum, wg)
	}
	for i := 0; i < 200; i++ {
		wg.Add(1)
		reqChan <- i
	}
	wg.Wait()
	if *sucNum > 100 {
		t.Errorf("Concurrency request sucNum expected <= 100, got %d", *sucNum)
	}
}
