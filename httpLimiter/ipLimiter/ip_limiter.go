package ipLimiter

import (
	"time"
)

type (
	// Visit 访问记录
	Visit struct {
		// 最后一次请求时间
		lastVisit time.Time
		// 对应Time窗口内的访问次数
		visitTimes uint64
	}

	// IpLimiter ip限流器
	IpLimiter struct {
		visitMap map[string]*Visit
	}
)

var App IpLimiter

// New 实例化：Ip 限流
func (IpLimiter) New() *IpLimiter {
	return &IpLimiter{visitMap: make(map[string]*Visit)}
}

// Affirm 检查限流
func (my *IpLimiter) Affirm(ip string, t time.Duration, maxVisitTimes uint64) (*Visit, bool) {
	if maxVisitTimes == 0 || t == 0 {
		// 如果限流为0，直接通过
		return nil, true
	}

	v, ok := my.visitMap[ip]
	if !ok {
		// 若该IP是首次请求，则初始化visit
		my.visitMap[ip] = &Visit{lastVisit: time.Now(), visitTimes: 1}
	} else {
		// 若该IP非首次请求，且距离上次请求时间超过Time窗口，则重设visitTimes
		if time.Since(v.lastVisit) > t {
			v.visitTimes = 1
		} else if v.visitTimes > maxVisitTimes {
			// 若本次请求距离上次请求时间在Time窗口内，且该IP在此时间内的访问次数超过上限，则返回错误
			return v, false
		} else {
			v.visitTimes++
		}
		v.lastVisit = time.Now()
	}

	return nil, true
}

// GetLastVisiter 获取最后访问时间
func (r *Visit) GetLastVisiter() time.Time {
	return r.lastVisit
}

// GetVisitTimes 获取窗口期内访问次数
func (r *Visit) GetVisitTimes() uint64 {
	return r.visitTimes
}
