package apilimiter

import (
	"sync"
	"time"
)

// 漏桶限流器
type LeakyBucket struct {
	Max      int64     //漏桶的最大存储上限
	Cycle    int64     //产生一块令牌的周期（每{cycle}毫秒生产一块令牌）
	lastTime time.Time //上一次加水的时间
	residue  int64     //漏桶剩余空间
	mutex    sync.Mutex
}

// NewLeakyBucket 初始化漏桶限流器
func (bucket *LeakyBucket) NewLeakyLimiter() {
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
	residue := bucket.residue + int64(duration/bucket.Cycle)
	// 漏桶最大剩余空间为bucket.Max
	if residue > bucket.Max {
		residue = bucket.Max
	}
	// 判断漏桶剩余空间是否足够
	if afterResidue := residue - num; afterResidue >= 0 {
		// 只有有剩余空间加水才更新漏桶剩余空间和加水时间
		bucket.lastTime = nowTime
		bucket.residue = afterResidue
		return true
	} else {
		return false
	}
}
