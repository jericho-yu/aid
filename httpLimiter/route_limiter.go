package httpLimiter

import (
	"sync"
	"time"
)

type (
	// visitor 访问者对象
	visitor struct {
		ipLimiter     *IpLimiter
		t             time.Duration
		maxVisitTimes uint16
	}

	// RouteLimiter 路由限流器
	RouteLimiter struct{ RouteSetMap *sync.Map }
)

var (
	routerLimiterOnce = sync.Once{}
	routerLimiterIns  *RouteLimiter
)

// OnceRouteLimiter 单例化：路由限流
func OnceRouteLimiter() *RouteLimiter {
	routerLimiterOnce.Do(func() { routerLimiterIns = &RouteLimiter{RouteSetMap: &sync.Map{}} })

	return routerLimiterIns
}

// Add 添加限流规则
func (my *RouteLimiter) Add(router string, t time.Duration, maxVisitTimes uint16) *RouteLimiter {
	if _, exist := my.RouteSetMap.Load(router); exist {
		my.RouteSetMap.Delete(router)
	}
	my.RouteSetMap.Store(router, &visitor{ipLimiter: NewIpLimiter(), t: t, maxVisitTimes: maxVisitTimes})
	return my
}

// Affirm 检查是否通过限流
func (my *RouteLimiter) Affirm(router, ip string) (*Visit, bool) {
	if val, exist := my.RouteSetMap.Load(router); exist {
		v := val.(*visitor)
		return v.ipLimiter.Affirm(ip, v.t, v.maxVisitTimes)
	}

	return nil, true
}
