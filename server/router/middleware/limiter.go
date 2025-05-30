package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"server/pkg/bean"
	"sync"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit // 每秒允许的请求数
	b   int        // 桶大小（突发请求）
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// 初始化限流器
// 每秒生成10个令牌，桶容量10个
var ipLimiter = NewIPRateLimiter(10, 10)

func IpLimiter(c *gin.Context) {
	// 获取客户端 IP（支持代理场景）
	ip := c.ClientIP()
	// 获取该 IP 的限流器
	l := ipLimiter.GetLimiter(ip)
	// 检查是否允许请求
	if !l.Allow() {
		bean.Response.Fail(c, "Too many requests")
		c.Abort()
		return
	}
	c.Next()
}
