package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"twin-os/backend/pkg/logger"
)

// CacheItem 缓存项
type CacheItem struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	ExpiresAt  time.Time   `json:"expires_at"`
	CreatedAt  time.Time   `json:"created_at"`
	AccessCount int       `json:"access_count"`
	LastAccess time.Time   `json:"last_access"`
}

// IsExpired 检查是否过期
func (c *CacheItem) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// Cache 内存缓存实现
type Cache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
	stats CacheStats
}

// CacheStats 缓存统计
type CacheStats struct {
	Hits        int64     `json:"hits"`
	Misses      int64     `json:"misses"`
	Sets        int64     `json:"sets"`
	Evictions   int64     `json:"evictions"`
	Cleanups    int64     `json:"cleanups"`
	TotalItems  int       `json:"total_items"`
	MemoryUsage int64     `json:"memory_usage_bytes"`
	LastCleanup time.Time `json:"last_cleanup"`
}

// NewCache 创建新的缓存实例
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]*CacheItem),
		stats: CacheStats{
			LastCleanup: time.Now(),
		},
	}

	// 启动定期清理协程
	go cache.startCleanupRoutine()

	return cache
}

// Set 设置缓存项
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl <= 0 {
		ttl = 5 * time.Minute // 默认5分钟
	}

	expiresAt := time.Now().Add(ttl)

	// 检查内存使用情况
	c.checkMemoryUsage()

	item := &CacheItem{
		Key:        key,
		Value:      value,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		AccessCount: 0,
		LastAccess: time.Now(),
	}

	c.items[key] = item
	c.stats.Sets++
	c.stats.TotalItems = len(c.items)

	logger.Debug(fmt.Sprintf("Cache SET: %s (TTL: %v)", key, ttl))
	return nil
}

// Get 获取缓存项
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		c.stats.Misses++
		return nil, false
	}

	// 检查是否过期
	if item.IsExpired() {
		c.stats.Misses++
		// 异步删除过期项
		go c.Delete(key)
		return nil, false
	}

	// 更新访问统计
	item.AccessCount++
	item.LastAccess = time.Now()
	c.stats.Hits++

	logger.Debug(fmt.Sprintf("Cache HIT: %s (access_count: %d)", key, item.AccessCount))
	return item.Value, true
}

// Delete 删除缓存项
func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; exists {
		delete(c.items, key)
		c.stats.TotalItems = len(c.items)
		logger.Debug(fmt.Sprintf("Cache DELETE: %s", key))
		return true
	}
	return false
}

// Clear 清空缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheItem)
	c.stats.TotalItems = 0
	logger.Info("Cache CLEARED")
}

// Cleanup 清理过期项
func (c *Cache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	deleted := 0

	for key, item := range c.items {
		if now.After(item.ExpiresAt) {
			delete(c.items, key)
			deleted++
		}
	}

	c.stats.TotalItems = len(c.items)
	c.stats.Evictions += int64(deleted)
	c.stats.Cleanups++
	c.stats.LastCleanup = now

	if deleted > 0 {
		logger.Info(fmt.Sprintf("Cache CLEANUP: deleted %d expired items", deleted))
	}

	return deleted
}

// GetStats 获取缓存统计
func (c *Cache) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// 计算内存使用
	var memory int64
	for _, item := range c.items {
		if jsonBytes, err := json.Marshal(item); err == nil {
			memory += int64(len(jsonBytes))
		}
	}

	stats := c.stats
	stats.MemoryUsage = memory
	return stats
}

// GenerateKey 生成缓存键
func (c *Cache) GenerateKey(prefix string, params map[string]interface{}) string {
	h := md5.New()

	// 添加前缀
	h.Write([]byte(prefix))
	h.Write([]byte(":"))

	// 添加参数（按字母排序保证一致性）
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	// 排序键以确保一致性
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	for _, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte("="))
		if valueBytes, err := json.Marshal(params[key]); err == nil {
			h.Write(valueBytes)
		}
		h.Write([]byte("&"))
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// checkMemoryUsage 检查内存使用情况
func (c *Cache) checkMemoryUsage() {
	maxItems := 10000 // 最大缓存项数

	if len(c.items) > maxItems {
		// LRU清理：删除最近最少使用的项
		type accessInfo struct {
			key         string
			lastAccess  time.Time
			accessCount int
		}

		var items []accessInfo
		for key, item := range c.items {
			items = append(items, accessInfo{
				key:         key,
				lastAccess:  item.LastAccess,
				accessCount: item.AccessCount,
			})
		}

		// 按最后访问时间排序
		for i := 0; i < len(items); i++ {
			for j := i + 1; j < len(items); j++ {
				if items[i].lastAccess.After(items[j].lastAccess) {
					items[i], items[j] = items[j], items[i]
				}
			}
		}

		// 删除最旧的 25% 项
		deleteCount := len(items) / 4
		for i := 0; i < deleteCount; i++ {
			delete(c.items, items[i].key)
		}

		c.stats.Evictions += int64(deleteCount)
		c.stats.TotalItems = len(c.items)

		logger.Warn(fmt.Sprintf("Cache LRU eviction: deleted %d items", deleteCount))
	}
}

// startCleanupRoutine 启动定期清理协程
func (c *Cache) startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
		c.Cleanup()
	}
}

// CacheManager 缓存管理器
type CacheManager struct {
	instances map[string]*Cache
	mu        sync.RWMutex
}

// NewCacheManager 创建缓存管理器
func NewCacheManager() *CacheManager {
	return &CacheManager{
		instances: make(map[string]*Cache),
	}
}

// GetCache 获取指定名称的缓存实例
func (cm *CacheManager) GetCache(name string) *Cache {
	cm.mu.RLock()
	cache, exists := cm.instances[name]
	cm.mu.RUnlock()

	if !exists {
		cm.mu.Lock()
		// 双重检查
		cache, exists = cm.instances[name]
		if !exists {
			cache = NewCache()
			cm.instances[name] = cache
		}
		cm.mu.Unlock()
	}

	return cache
}

// GetAllStats 获取所有缓存统计
func (cm *CacheManager) GetAllStats() map[string]CacheStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := make(map[string]CacheStats)
	for name, cache := range cm.instances {
		stats[name] = cache.GetStats()
	}

	return stats
}

// ClearAll 清空所有缓存
func (cm *CacheManager) ClearAll() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, cache := range cm.instances {
		cache.Clear()
	}
}

// 全局缓存管理器实例
var GlobalCacheManager = NewCacheManager()

// GetGlobalCache 获取全局缓存实例
func GetGlobalCache(name string) *Cache {
	return GlobalCacheManager.GetCache(name)
}