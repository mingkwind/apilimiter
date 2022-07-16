package apilimiter

import (
	"sync"
	"time"
)

// 漏桶限流器
type LeakyBucket struct {
	Max      int64     //漏桶的最大存储上限
	Rate     int64     //水每10ms漏出的速率，亦即每10ms产生令牌的数量
	lastTime time.Time //上一次加水的时间
	residue  int64     //漏桶剩余空间
	mutex    sync.Mutex
}

// NewLeakyBucket 初始化漏桶限流器
func (bucket *LeakyBucket) NewLeakyBucket() {
	bucket.lastTime = time.Now()
	bucket.residue = 0 //bucket.Max
}

// GetToken 获取令牌
func (bucket *LeakyBucket) GetToken(num int64) bool {
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	// 获取当前时间
	nowTime := time.Now()
	duration := nowTime.Sub(bucket.lastTime).Milliseconds()
	if duration < 10 {
		// 如果时间间隔小于10ms，则不能加水
		return false
	}
	residue := bucket.residue + int64(duration/10)*bucket.Rate
	// 漏桶最大剩余空间为bucket.Max
	if residue > bucket.Max {
		residue = bucket.Max
	}
	bucket.lastTime = nowTime
	// 判断漏桶剩余空间是否足够
	if afterResidue := residue - num; afterResidue >= 0 {
		bucket.residue = afterResidue
		return true
	} else {
		bucket.residue = residue
		return false
	}
}
