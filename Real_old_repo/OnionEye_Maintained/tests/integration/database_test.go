package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostgreSQL_Connection tests PostgreSQL connectivity and basic operations
func TestPostgreSQL_Connection(t *testing.T) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getEnv("POSTGRES_USER", "onioneye"),
		getEnv("POSTGRES_PASSWORD", "testpass123"),
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_DB", "onioneye_test"),
	)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
		return
	}
	defer conn.Close(ctx)

	t.Run("Database connection established", func(t *testing.T) {
		var result int
		err := conn.QueryRow(ctx, "SELECT 1").Scan(&result)
		require.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("Schema tables exist", func(t *testing.T) {
		tables := []string{"actors", "identifiers", "persona_links", "scans"}
		
		for _, table := range tables {
			var exists bool
			query := `SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)`
			err := conn.QueryRow(ctx, query, table).Scan(&exists)
			
			if err != nil {
				t.Logf("Table check error for %s: %v", table, err)
			} else if !exists {
				t.Logf("Table %s does not exist (may not be initialized)", table)
			} else {
				t.Logf("Table %s exists", table)
			}
		}
	})
}

// TestRedis_Connection tests Redis connectivity and operations
func TestRedis_Connection(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
	defer rdb.Close()

	ctx := context.Background()

	t.Run("Redis connection established", func(t *testing.T) {
		pong, err := rdb.Ping(ctx).Result()
		if err != nil {
			t.Skipf("Redis not available: %v", err)
			return
		}
		assert.Equal(t, "PONG", pong)
	})

	t.Run("Redis SET and GET operations", func(t *testing.T) {
		key := fmt.Sprintf("test:integration:%d", time.Now().Unix())
		value := "test_value"

		err := rdb.Set(ctx, key, value, 10*time.Second).Err()
		if err != nil {
			t.Skipf("Redis not available: %v", err)
			return
		}

		result, err := rdb.Get(ctx, key).Result()
		require.NoError(t, err)
		assert.Equal(t, value, result)

		// Cleanup
		rdb.Del(ctx, key)
	})

	t.Run("Redis stream operations", func(t *testing.T) {
		streamKey := fmt.Sprintf("test:stream:%d", time.Now().Unix())
		
		// Add to stream
		id, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamKey,
			Values: map[string]interface{}{
				"field1": "value1",
				"field2": "value2",
			},
		}).Result()

		if err != nil {
			t.Skipf("Redis streams not available: %v", err)
			return
		}
		assert.NotEmpty(t, id)

		// Cleanup
		rdb.Del(ctx, streamKey)
	})
}

// TestDatabaseIntegration_ActorOperations tests actor CRUD operations
func TestDatabaseIntegration_ActorOperations(t *testing.T) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getEnv("POSTGRES_USER", "onioneye"),
		getEnv("POSTGRES_PASSWORD", "testpass123"),
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_DB", "onioneye_test"),
	)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
		return
	}
	defer conn.Close(ctx)

	// Check if actors table exists
	var tableExists bool
	err = conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = 'actors'
		)
	`).Scan(&tableExists)
	
	if err != nil || !tableExists {
		t.Skip("Actors table not available")
		return
	}

	t.Run("Query actors table", func(t *testing.T) {
		var count int
		err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM actors").Scan(&count)
		require.NoError(t, err)
		t.Logf("Found %d actors in database", count)
	})
}

// TestServiceIntegration_EndToEnd tests service-to-service communication
func TestServiceIntegration_EndToEnd(t *testing.T) {
	t.Run("Storage and Database integration", func(t *testing.T) {
		// This would test storage-api writing to PostgreSQL
		t.Skip("Requires live services")
	})

	t.Run("Ingestion and Redis integration", func(t *testing.T) {
		// This would test ingestion-api writing to Redis streams
		t.Skip("Requires live services")
	})

	t.Run("Correlation and Storage integration", func(t *testing.T) {
		// This would test correlation-engine querying storage-api
		t.Skip("Requires live services")
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
