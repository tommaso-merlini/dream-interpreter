package handler

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	requestCounts = make(map[string]int)
	countMutex    sync.Mutex
)

func RateLimiterMiddleware(maxRequests int, resetTime time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Request().URL.Path, "/public") {
				return next(c)
			}
			ip := c.RealIP()

			countMutex.Lock()
			count, exists := requestCounts[ip]
			if !exists || count == 0 {
				// Reset the count every resetTime
				go func(ip string) {
					time.Sleep(resetTime)
					countMutex.Lock()
					delete(requestCounts, ip)
					countMutex.Unlock()
				}(ip)
			}
			requestCounts[ip] = count + 1
			countMutex.Unlock()

			if count >= maxRequests {
				println("Rate limit exceeded")
				return c.JSON(
					http.StatusTooManyRequests,
					map[string]string{"error": "Rate limit exceeded. Please try again later."},
				)
			}
			println("Rate limit not exceeded")

			return next(c)
		}
	}
}
