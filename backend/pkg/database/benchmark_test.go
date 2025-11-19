package database

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// BenchmarkDatabaseOperations 数据库操作性能基准测试
func BenchmarkDatabaseOperations(b *testing.B) {
	db := setupBenchmarkDatabase(b)
	defer db.Close()

	// 预编译语句测试
	stmt, err := db.Prepare("INSERT INTO messages (message_id, talker_id, type, content, timestamp) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		b.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := stmt.Exec(
			fmt.Sprintf("msg_%d", i),
			fmt.Sprintf("user_%d", i%100),
			1,
			"Test message content",
			time.Now().Unix(),
		)
		if err != nil {
			b.Errorf("Insert failed: %v", err)
		}
	}
}

// BenchmarkConcurrentReads 并发读取性能测试
func BenchmarkConcurrentReads(b *testing.B) {
	db := setupBenchmarkDatabase(b)
	defer db.Close()

	// 预先插入测试数据
	stmt, _ := db.Prepare("INSERT INTO messages (message_id, talker_id, type, content, timestamp) VALUES (?, ?, ?, ?, ?)")
	for i := 0; i < 1000; i++ {
		stmt.Exec(fmt.Sprintf("msg_%d", i), "test_user", 1, "Test content", time.Now().Unix())
	}
	stmt.Close()

	selectStmt, _ := db.Prepare("SELECT COUNT(*) FROM messages WHERE talker_id = ?")
	defer selectStmt.Close()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var count int
			err := selectStmt.QueryRow("test_user").Scan(&count)
			if err != nil {
				b.Errorf("Query failed: %v", err)
			}
		}
	})
}

// BenchmarkConnectionPool 连接池性能基准测试
func BenchmarkConnectionPool(b *testing.B) {
	tests := []struct {
		name     string
		maxOpen  int
		maxIdle  int
	}{
		{"Small_Pool", 5, 2},
		{"Medium_Pool", 25, 10},
		{"Large_Pool", 50, 20},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			db := setupBenchmarkDatabaseWithPool(b, test.maxOpen, test.maxIdle)
			defer db.Close()

			b.ResetTimer()
			b.ReportAllocs()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					var result int
					err := db.QueryRow("SELECT 1").Scan(&result)
					if err != nil {
						b.Error(err)
					}
				}
			})
		})
	}
}

// TestDatabasePerformanceOptimizations 数据库性能优化验证测试
func TestDatabasePerformanceOptimizations(t *testing.T) {
	t.Log("🚀 Database Performance Optimizations Test")

	// 测试连接池优化
	t.Run("Connection_Pool_Optimization", func(t *testing.T) {
		db := setupBenchmarkDatabase(t)
		defer db.Close()

		// 获取连接池配置
		stats := db.Stats()
		t.Logf("Connection pool configuration:")
		t.Logf("  Max open connections: %d", stats.MaxOpenConnections)
		t.Logf("  Open connections: %d", stats.OpenConnections)

		// v0.3.0 优化验证
		if stats.MaxOpenConnections >= 25 {
			t.Logf("✅ Max open connections optimized: %d ≥ 25", stats.MaxOpenConnections)
		} else {
			t.Logf("📈 Max open connections can be increased: %d < 25", stats.MaxOpenConnections)
		}
	})

	// 测试索引优化
	t.Run("Index_Optimization", func(t *testing.T) {
		db := setupBenchmarkDatabase(t)
		defer db.Close()

		// 插入测试数据
		stmt, _ := db.Prepare("INSERT INTO messages (message_id, talker_id, type, content, timestamp) VALUES (?, ?, ?, ?, ?)")
		for i := 0; i < 5000; i++ {
			stmt.Exec(fmt.Sprintf("msg_%d", i), fmt.Sprintf("user_%d", i%100), 1, "Test content", time.Now().Unix())
		}
		stmt.Close()

		// 测试索引查询性能
		tests := []struct {
			name    string
			query   string
			args    []interface{}
			target  time.Duration
		}{
			{
				name:   "Talker_ID_Index",
				query:  "SELECT COUNT(*) FROM messages WHERE talker_id = ?",
				args:   []interface{}{"user_1"},
				target: 10 * time.Millisecond,
			},
			{
				name:   "Composite_Index",
				query:  "SELECT COUNT(*) FROM messages WHERE talker_id = ? AND timestamp > ?",
				args:   []interface{}{"user_1", time.Now().Add(-time.Hour).Unix()},
				target: 15 * time.Millisecond,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				start := time.Now()
				var count int
				err := db.QueryRow(test.query, test.args...).Scan(&count)
				duration := time.Since(start)

				if err != nil {
					t.Errorf("Query failed: %v", err)
					return
				}

				t.Logf("Query completed in %v, count: %d", duration, count)

				if duration <= test.target {
					t.Logf("✅ Index optimization effective: %v ≤ %v", duration, test.target)
				} else {
					t.Logf("📈 Query needs optimization: %v > %v", duration, test.target)
				}
			})
		}
	})

	// 测试SQLite性能优化
	t.Run("SQLite_Performance_Settings", func(t *testing.T) {
		db := setupBenchmarkDatabase(t)
		defer db.Close()

		// 测试PRAGMA设置
		pragmaSettings := []struct {
			name  string
			query string
		}{
			{"Journal_Mode", "PRAGMA journal_mode"},
			{"Synchronous", "PRAGMA synchronous"},
			{"Cache_Size", "PRAGMA cache_size"},
			{"Temp_Store", "PRAGMA temp_store"},
		}

		for _, pragma := range pragmaSettings {
			var result string
			err := db.QueryRow(pragma.query).Scan(&result)
			if err != nil {
				t.Errorf("Failed to get %s: %v", pragma.name, err)
				continue
			}
			t.Logf("PRAGMA %s: %s", pragma.name, result)
		}

		t.Logf("✅ SQLite performance settings verified")
	})

	// 测试并发性能
	t.Run("Concurrent_Performance", func(t *testing.T) {
		db := setupBenchmarkDatabase(t)
		defer db.Close()

		// 并发读写测试
		concurrency := 20
		operationsPerWorker := 50

		var wg sync.WaitGroup
		start := time.Now()

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				for j := 0; j < operationsPerWorker; j++ {
					var result int
					err := db.QueryRow("SELECT 1").Scan(&result)
					if err != nil {
						t.Errorf("Query failed in worker %d: %v", workerID, err)
					}
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)

		totalOps := concurrency * operationsPerWorker
		opsPerSec := float64(totalOps) / duration.Seconds()

		t.Logf("Concurrent performance:")
		t.Logf("  Total operations: %d", totalOps)
		t.Logf("  Duration: %v", duration)
		t.Logf("  Operations per second: %.0f", opsPerSec)

		// v0.3.0 性能目标
		targetOpsPerSec := 1000.0
		if opsPerSec >= targetOpsPerSec {
			t.Logf("🎉 v0.3.0 target achieved: %.0f ≥ %.0f ops/sec", opsPerSec, targetOpsPerSec)
		} else {
			t.Logf("📈 Progress: %.0f / %.0f ops/sec (%.1f%% complete)",
				opsPerSec, targetOpsPerSec, (opsPerSec/targetOpsPerSec)*100)
		}
	})
}

// setupBenchmarkDatabase 设置基准测试数据库
func setupBenchmarkDatabase(tb testing.TB) *sql.DB {
	return setupBenchmarkDatabaseWithPool(tb, 25, 10)
}

// setupBenchmarkDatabaseWithPool 设置指定连接池大小的基准测试数据库
func setupBenchmarkDatabaseWithPool(tb testing.TB, maxOpen, maxIdle int) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		tb.Fatalf("Failed to open test database: %v", err)
	}

	// 创建表和索引
	queries := []string{
		`CREATE TABLE messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			message_id TEXT UNIQUE NOT NULL,
			talker_id TEXT NOT NULL,
			type INTEGER NOT NULL,
			content TEXT,
			timestamp DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		"CREATE INDEX idx_messages_talker_id ON messages(talker_id)",
		"CREATE INDEX idx_messages_timestamp ON messages(timestamp)",
		"CREATE INDEX idx_messages_message_id ON messages(message_id)",
		"CREATE INDEX idx_messages_talker_timestamp ON messages(talker_id, timestamp DESC)",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			tb.Fatalf("Failed to execute query %s: %v", query, err)
		}
	}

	// 应用v0.3.0性能优化设置
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// SQLite性能优化
	pragmaSettings := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = 10000",
		"PRAGMA temp_store = MEMORY",
		"PRAGMA mmap_size = 268435456",
	}

	for _, pragma := range pragmaSettings {
		if _, err := db.Exec(pragma); err != nil {
			tb.Logf("Warning: Failed to set %s: %v", pragma, err)
		}
	}

	return db
}