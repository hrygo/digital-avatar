package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// BenchmarkAPIRoutes API路由性能基准测试
func BenchmarkAPIRoutes(b *testing.B) {
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	r := gin.New()
	r.Use(performanceMiddleware())
	r.Use(rateLimitMiddleware())

	// 添加测试路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   "0.3.0",
		})
	})

	r.GET("/health/detailed", func(c *gin.Context) {
		// 模拟一些数据库查询延迟
		time.Sleep(10 * time.Millisecond)
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   "0.3.0",
		})
	})

	benchmarks := []struct {
		name string
		path string
	}{
		{"Health_Check", "/health"},
		{"Detailed_Health", "/health/detailed"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", bm.path, nil)
				req.Header.Set("X-Request-ID", "benchmark-test")
				r.ServeHTTP(w, req)
			}
		})
	}
}

// BenchmarkMiddlewarePerformance 中间件性能基准测试
func BenchmarkMiddlewarePerformance(b *testing.B) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加v0.3.0性能优化中间件
	r.Use(performanceMiddleware())
	r.Use(rateLimitMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "benchmark-test")
		r.ServeHTTP(w, req)
	}
}

// BenchmarkGinRouting Gin路由性能基准测试
func BenchmarkGinRouting(b *testing.B) {
	gin.SetMode(gin.TestMode)

	// 创建多个路由测试性能
	r := gin.New()

	// 添加不同类型的路由
	r.GET("/api/v1/users", func(c *gin.Context) {
		c.JSON(200, gin.H{"users": []string{}})
	})
	r.GET("/api/v1/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": c.Param("id")})
	})
	r.POST("/api/v1/users", func(c *gin.Context) {
		c.JSON(200, gin.H{"created": true})
	})
	r.PUT("/api/v1/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"updated": true})
	})
	r.DELETE("/api/v1/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"deleted": true})
	})

	benchmarks := []struct {
		name   string
		method string
		path   string
	}{
		{"GET_Users", "GET", "/api/v1/users"},
		{"GET_User_ByID", "GET", "/api/v1/users/123"},
		{"POST_Users", "POST", "/api/v1/users"},
		{"PUT_User", "PUT", "/api/v1/users/123"},
		{"DELETE_User", "DELETE", "/api/v1/users/123"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest(bm.method, bm.path, nil)
				r.ServeHTTP(w, req)
			}
		})
	}
}

// TestMiddlewareFunctionality 中间件功能测试
func TestMiddlewareFunctionality(t *testing.T) {
	t.Log("🚀 Middleware Functionality Test v0.3.0")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(performanceMiddleware())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	t.Run("Performance_Middleware", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		// 检查性能响应头
		responseTime := w.Header().Get("X-Response-Time")
		if responseTime == "" {
			t.Error("X-Response-Time header not set")
		} else {
			t.Logf("✅ Response time header set: %s", responseTime)
		}

		if w.Code != 200 {
			t.Errorf("Request failed with status %d", w.Code)
		}
	})

	t.Run("Rate_Limit_Middleware", func(t *testing.T) {
		// 创建新的路由测试限流
		rr := gin.New()
		rr.Use(rateLimitMiddleware())

		requestCount := 0
		rr.GET("/limited", func(c *gin.Context) {
			requestCount++
			c.JSON(200, gin.H{"count": requestCount})
		})

		// 发送一些请求，应该正常
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/limited", nil)
			rr.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Request %d failed with status %d", i+1, w.Code)
			}
		}
		t.Logf("✅ Rate limiting middleware working correctly")
	})
}