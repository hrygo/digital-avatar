package cache

import (
	"bytes"
	"bufio"
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"twin-os/backend/pkg/logger"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	TTL           time.Duration `json:"ttl"`
	SkipMethods   []string    `json:"skip_methods"`
	SkipHeaders   []string    `json:"skip_headers"`
	VaryHeaders   []string    `json:"vary_headers"`
	MaxBodySize   int64       `json:"max_body_size"`
	KeyGenerator  func(*gin.Context) string
}

// DefaultCacheConfig 默认缓存配置
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		TTL:          5 * time.Minute,
		SkipMethods:  []string{"POST", "PUT", "DELETE", "PATCH"},
		SkipHeaders:  []string{"Authorization", "Cookie"},
		VaryHeaders:  []string{"Accept-Language", "User-Agent"},
		MaxBodySize:  1024 * 1024, // 1MB
		KeyGenerator: DefaultKeyGenerator,
	}
}

// CacheMiddleware 创建缓存中间件
func CacheMiddleware(cache *Cache, config CacheConfig) gin.HandlerFunc {
	if config.KeyGenerator == nil {
		config.KeyGenerator = DefaultKeyGenerator
	}

	return func(c *gin.Context) {
		// 检查是否应该跳过缓存
		if shouldSkipCache(c, config) {
			c.Next()
			return
		}

		// 生成缓存键
		cacheKey := config.KeyGenerator(c)

		// 尝试从缓存获取响应
		if cachedResponse, found := cache.Get(cacheKey); found {
			if response, ok := cachedResponse.(*CachedResponse); ok {
				// 设置缓存相关的响应头
				c.Header("X-Cache", "HIT")
				c.Header("X-Cache-Key", cacheKey)
				c.Header("X-Cache-Age", time.Since(response.CreatedAt).String())

				// 写入缓存的响应
				for key, values := range response.Headers {
					for _, value := range values {
						c.Header(key, value)
					}
				}
				c.Data(response.StatusCode, response.ContentType, response.Body)
				c.Abort()
				return
			}
		}

		// 拦截响应
		c.Writer = &responseWriter{
			ResponseWriter: c.Writer,
			context:        c,
			cache:          cache,
			cacheKey:       cacheKey,
			config:         config,
		}

		c.Next()
	}
}

// shouldSkipCache 检查是否应该跳过缓存
func shouldSkipCache(c *gin.Context, config CacheConfig) bool {
	// 检查请求方法
	for _, method := range config.SkipMethods {
		if c.Request.Method == method {
			return true
		}
	}

	// 检查查询参数中的缓存控制
	if c.Query("no_cache") == "1" || c.Query("refresh") == "1" {
		return true
	}

	// 检查请求头中的缓存控制
	if c.GetHeader("Cache-Control") == "no-cache" || c.GetHeader("Pragma") == "no-cache" {
		return true
	}

	// 检查路径前缀
	skipPaths := []string{
		"/api/v1/wechat/sync",
		"/api/v1/error/report",
		"/api/v1/data/export",
	}

	for _, path := range skipPaths {
		if strings.HasPrefix(c.Request.URL.Path, path) {
			return true
		}
	}

	return false
}

// DefaultKeyGenerator 默认缓存键生成器
func DefaultKeyGenerator(c *gin.Context) string {
	// 基础信息
	h := md5.New()
	h.Write([]byte(c.Request.Method))
	h.Write([]byte(c.Request.URL.Path))

	// 查询参数
	if c.Request.URL.RawQuery != "" {
		h.Write([]byte(c.Request.URL.RawQuery))
	}

	// 重要的请求头
	varyHeaders := []string{"Accept-Language", "User-Agent", "Accept"}
	for _, header := range varyHeaders {
		if value := c.GetHeader(header); value != "" {
			h.Write([]byte(header))
			h.Write([]byte(":"))
			h.Write([]byte(value))
		}
	}

	return fmt.Sprintf("api:%x", h.Sum(nil))
}

// CachedResponse 缓存的响应
type CachedResponse struct {
	StatusCode int               `json:"status_code"`
	ContentType string            `json:"content_type"`
	Headers     map[string][]string `json:"headers"`
	Body        []byte            `json:"body"`
	CreatedAt   time.Time         `json:"created_at"`
}

// responseWriter 拦截响应的写入器
type responseWriter struct {
	gin.ResponseWriter
	context  *gin.Context
	cache    *Cache
	cacheKey string
	config   CacheConfig

	buffer    *bytes.Buffer
	status    int
	written   bool
}

// Write 实现io.Writer接口
func (w *responseWriter) Write(data []byte) (int, error) {
	if !w.written {
		// 首次写入，初始化缓冲区
		w.buffer = bytes.NewBuffer(data)
		w.written = true
	} else {
		w.buffer.Write(data)
	}

	return w.ResponseWriter.Write(data)
}

// WriteHeader 实现http.ResponseWriter接口
func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Hijack 实现http.Hijacker接口
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("response writer does not implement http.Hijacker")
}

// Flush 实现http.Flusher接口
func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// shouldCacheResponse 检查响应是否应该被缓存
func (w *responseWriter) shouldCacheResponse() bool {
	// 检查状态码
	cacheableStatusCodes := []int{200, 201, 202, 204, 301, 302, 304, 404}
	statusCode := w.status
	if statusCode == 0 {
		statusCode = 200 // 默认状态码
	}

	cacheable := false
	for _, code := range cacheableStatusCodes {
		if statusCode == code {
			cacheable = true
			break
		}
	}

	if !cacheable {
		return false
	}

	// 检查内容类型
	contentType := w.context.GetHeader("Content-Type")
	if contentType != "" {
		// 只缓存特定类型的响应
		cacheableTypes := []string{
			"application/json",
			"text/html",
			"text/plain",
			"application/xml",
		}

		for _, ct := range cacheableTypes {
			if strings.Contains(contentType, ct) {
				cacheable = true
				break
			}
		}

		if !cacheable {
			return false
		}
	}

	// 检查响应大小
	if w.buffer != nil && int64(w.buffer.Len()) > w.config.MaxBodySize {
		return false
	}

	return true
}

// saveToCache 保存响应到缓存
func (w *responseWriter) saveToCache() {
	if !w.shouldCacheResponse() {
		return
	}

	// 收集响应头
	headers := make(map[string][]string)
	for key, values := range w.context.Writer.Header() {
		// 过滤掉不应该缓存的头
		shouldSkip := false
		for _, skipHeader := range w.config.SkipHeaders {
			if strings.EqualFold(key, skipHeader) {
				shouldSkip = true
				break
			}
		}

		if !shouldSkip {
			headers[key] = values
		}
	}

	// 获取响应体
	var body []byte
	if w.buffer != nil {
		body = w.buffer.Bytes()
	} else {
		body = []byte{}
	}

	// 创建缓存响应对象
	cachedResponse := &CachedResponse{
		StatusCode: w.status,
		ContentType: w.context.GetHeader("Content-Type"),
		Headers:     headers,
		Body:        body,
		CreatedAt:   time.Now(),
	}

	// 保存到缓存
	if err := w.cache.Set(w.cacheKey, cachedResponse, w.config.TTL); err != nil {
		logger.Warn("Failed to cache response: " + err.Error())
	} else {
		w.context.Header("X-Cache", "MISS")
		w.context.Header("X-Cache-Key", w.cacheKey)
		logger.Debug(fmt.Sprintf("Cached response: %s (size: %d bytes)", w.cacheKey, len(body)))
	}
}

// 清理缓存响应写入器的内存
func (w *responseWriter) cleanup() {
	if w.buffer != nil {
		w.buffer.Reset()
		w.buffer = nil
	}
}