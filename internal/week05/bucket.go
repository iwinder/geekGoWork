package week05

import (
	"sync"
	"time"
)

type Bucket struct {
	sync.RWMutex
	// 总请求
	TotalCount int64
	// 失败次数
	ErrorCount int64
	// 时间
	Time int64
	// 时间 方便展示
	Timestamp time.Time
}

func NewBucket() *Bucket {
	timestamp := time.Now()
	return &Bucket{
		Timestamp: timestamp,
		Time:      timestamp.UnixNano() / 1e6}
}

func (b *Bucket) Increment(result bool) {
	b.Lock()
	defer b.Unlock()
	if !result {
		b.ErrorCount++
	}
	b.TotalCount++
}
