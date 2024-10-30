package routeLimiter

import (
	"sync"
	"time"

	"github.com/jericho-yu/aid/httpLimiter/ipLimiter"
)

type (
	// visitor 访问者对象
	visitor struct {
		ipLimiter     *ipLimiter.IpLimiter
		t             time.Duration
		maxVisitTimes uint64
	}

	// RouteLimiter 路由限流器
	RouteLimiter struct {
		RouteSetMap *sync.Map
	}
)

var (
	routerLimiterOnce = sync.Once{}
	routerLimiterIns  *RouteLimiter
	App               RouteLimiter
)

// Once 单例化：路由限流
func (RouteLimiter) Once() *RouteLimiter {
	routerLimiterOnce.Do(func() { routerLimiterIns = &RouteLimiter{RouteSetMap: &sync.Map{}} })
	return routerLimiterIns
}

// Add 添加限流规则
func (r *RouteLimiter) Add(router string, t time.Duration, maxVisitTimes uint64) *RouteLimiter {

	if _, exist := r.RouteSetMap.Load(router); exist {
		r.RouteSetMap.Delete(router)
	}
	r.RouteSetMap.Store(router, &visitor{ipLimiter: ipLimiter.App.New(), t: t, maxVisitTimes: maxVisitTimes})
	return r
}

// Affirm 检查是否通过限流
func (r *RouteLimiter) Affirm(router, ip string) (*ipLimiter.Visit, bool) {
	if val, exist := r.RouteSetMap.Load(router); exist {
		v := val.(*visitor)
		return v.ipLimiter.Affirm(ip, v.t, v.maxVisitTimes)
	}

	return nil, true
}
