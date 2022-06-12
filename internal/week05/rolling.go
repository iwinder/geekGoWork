package week05

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	opt "github.com/iwinder/geekGoWork/internal/pkg/options"
	"github.com/pkg/errors"
	"sync"
	"time"
)

//var rl sync.RWMutex

type RollingWindow struct {
	sync.RWMutex
	// 触发熔断的请求总数阈值
	limitCount int64
	// 触发熔断的失败率阈值
	errorPercentage float64

	// 滑动窗口的长度（时间间隔） ms  1000 = 1s
	timeInMilliseconds int
	// 滑动窗口中桶的个数
	limitBucket int
	// 熔断间隔
	brokenTimePeriod int64

	// 每个桶对应的窗口长度
	bucketSizeInMs int
	// 桶
	buckets []*Bucket
	// 当前总值
	hc *HealthCounts
	// 熔断状态
	brokenFlag bool
	// 上次熔断发生时间
	lastBreakTime      int64
	lastInitBucketTIme int64

	done chan int
}

type HealthCounts struct {
	TotalRequests   int64
	ErrorRequests   int64
	ErrorPercentage float64
}

func NewHealthCounts() *HealthCounts {
	return &HealthCounts{
		TotalRequests:   0,
		ErrorRequests:   0,
		ErrorPercentage: 0,
	}
}

// NewRollingWindow 新建实例
func NewRollingWindow(opt *opt.RollingOption) *RollingWindow {
	if opt.TimeInMilliseconds%opt.LimitBucket != 0 {
		fmt.Errorf("the timeInMilliseconds must divide equally into limitBucket. For example 1000/10 is ok, 1000/11 is not")
	}
	r := &RollingWindow{
		timeInMilliseconds: opt.TimeInMilliseconds,
		limitBucket:        opt.LimitBucket,
		bucketSizeInMs:     opt.TimeInMilliseconds / opt.LimitBucket,
		buckets:            make([]*Bucket, 0, opt.LimitBucket),
		lastInitBucketTIme: getNowTimeInMs(),
		limitCount:         opt.LimitCount,
		errorPercentage:    opt.ErrorPercentage,
		brokenTimePeriod:   opt.BrokenTimePeriod,
		hc:                 NewHealthCounts(),
	}
	r.done = make(chan int)
	return r
}

func (r *RollingWindow) RunWindow() {
	//go func() {
	for {
		select {
		case <-r.done:
			glog.V(2).Infoln("RollingWindow ShutDone RunWindow...")
			return
		default:
			r.appendBucket()
			glog.V(2).Infoln("RollingWindow  RunWindow created...")
			// 每隔 r.bucketSizeInMs 创建一个新的桶
			time.Sleep(time.Millisecond * time.Duration(r.bucketSizeInMs))
		}
	}
	glog.V(2).Infoln("RollingWindow   RunWindow ZZZ21...")
	//}()
}

// CheckBroken 检测是否阻塞
func (r *RollingWindow) CheckBroken() bool {
	r.Lock()
	defer r.Unlock()
	if r.brokenFlag {
		if r.getBrokenTimeFlag() {
			r.brokenFlag = false
		}
		return r.brokenFlag
	}
	if r.getBreakJudgementState() {
		r.brokenFlag = true
		r.lastBreakTime = getNowTimeInMs()
		glog.V(2).Infoln("RollingWindow need brokenFlag...")
	}
	return r.brokenFlag
}

// RecordReqResult 在桶中记录当次结果
func (r *RollingWindow) RecordReqResult(result bool) {
	r.getCurrentBucket().Increment(result)
}

// ShowAllBucket 展示当前滑动窗口的所有桶状态
func (r *RollingWindow) ShowAllBucket() {
	for _, v := range r.buckets {
		fmt.Printf("id: [%v] | total: [%d] | failed: [%d]\n", v.Timestamp, v.TotalCount, v.ErrorCount)
	}
}

// ShutDone 退出
func (r *RollingWindow) ShutDone(ctx context.Context) error {
	close(r.done)
	select {
	case <-ctx.Done():
		return errors.New("timeout")
	}
}

// 增加新桶
func (r *RollingWindow) appendBucket() {
	r.Lock()
	defer r.Unlock()
	currentTime := getNowTimeInMs()
	// 当超过滑动窗口的长度，重新创建新桶
	if currentTime-r.lastInitBucketTIme >= int64(r.timeInMilliseconds) {
		r.buckets = make([]*Bucket, 0, r.limitBucket)
	}
	r.buckets = append(r.buckets, NewBucket())
	if !(len(r.buckets) < r.limitBucket+1) {
		r.buckets = r.buckets[1:]
	}
}

// 获取当前可用桶
func (r *RollingWindow) getCurrentBucket() *Bucket {
	currentTime := getNowTimeInMs()
	currentBucket := r.peekLast()

	// 当前桶不为空且当前时间在当前桶的时间范围内
	if currentBucket != nil && currentBucket.Time+int64(r.bucketSizeInMs) > currentTime {
		return currentBucket
	}
	r.appendBucket()
	return r.peekLast()
}

func (r *RollingWindow) peekLast() *Bucket {
	return r.buckets[len(r.buckets)-1]
}

// 获取是否需要熔断
func (r *RollingWindow) getBreakJudgementState() bool {
	r.hc.TotalRequests = 0
	r.hc.ErrorRequests = 0
	for _, v := range r.buckets {
		r.hc.TotalRequests += v.TotalCount
		r.hc.ErrorRequests += v.ErrorCount
	}
	if r.hc.ErrorRequests == 0 {
		r.hc.ErrorPercentage = 0
	} else {
		r.hc.ErrorPercentage = float64(r.hc.TotalRequests) / float64(r.hc.ErrorRequests)
	}

	glog.V(2).Infoln("RollingWindow getBreakJudgementState...", r.hc.TotalRequests, r.hc.ErrorRequests, r.hc.ErrorPercentage)
	// 如果超过失败率或者超过总请求数 则触发熔断
	if r.hc.ErrorPercentage >= r.errorPercentage || r.hc.TotalRequests >= r.limitCount {
		return true
	}
	return false
}

// 获取当前时间是否已过熔断间隔时间
func (r *RollingWindow) getBrokenTimeFlag() bool {
	currentTime := getNowTimeInMs()
	return r.lastBreakTime+r.brokenTimePeriod > currentTime
}

func getNowTimeInMs() int64 {
	return time.Now().UnixNano() / 1e6
}
