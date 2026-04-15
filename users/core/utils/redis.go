package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient  *redis.Client
	redisEnabled bool
)

// InitRedis menginisialisasi koneksi Redis jika diaktifkan di .env
func InitRedis() error {
	// 1. Cek Feature Flag di .env
	if GetEnv("redis") != "true" {
		redisEnabled = false
		log.Println("[Redis] Feature is DISABLED (Running in No-Cache mode)")
		return nil
	}

	// 2. Ambil Config
	host := GetEnv("redis_host", "localhost")
	port := GetEnv("redis_port", "6379")
	pass := GetEnv("redis_pass", "")

	// Konversi DB dari string ke int
	dbStr := GetEnv("redis_db", "0")
	dbInt, err := strconv.Atoi(dbStr)
	if err != nil {
		dbInt = 0 // Default ke 0 jika config salah
	}

	// PERBAIKAN: Gunakan 127.0.0.1 instead of localhost untuk force IPv4
	if host == "localhost" {
		host = "127.0.0.1"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	// ✅ FIX: Trim password, jangan kirim empty string sebagai AUTH
	pass = strings.TrimSpace(pass)

	// 3. Setup Client
	redisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass, // ✅ empty string = no AUTH command
		DB:           dbInt,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// 4. Test Ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		// Jangan fatal, tapi disable redis
		log.Printf("[Redis] Connection FAILED: %v (Switching to No-Cache mode)", err)
		redisEnabled = false
		redisClient = nil
		return nil // Return nil agar aplikasi tetap jalan
	}

	redisEnabled = true
	log.Printf("[Redis] Connected successfully to %s", addr)
	return nil
}

// RedisSet menyimpan data ke cache
func RedisSet(ctx context.Context, key string, value any, expiration time.Duration) error {
	if !redisEnabled || redisClient == nil {
		return nil
	}
	return redisClient.Set(ctx, key, value, expiration).Err()
}

// RedisGet mengambil data dari cache
func RedisGet(ctx context.Context, key string) (string, error) {
	if !redisEnabled || redisClient == nil {
		return "", fmt.Errorf("redis disabled")
	}
	return redisClient.Get(ctx, key).Result()
}

// RedisDel menghapus data
func RedisDel(ctx context.Context, key string) error {
	if !redisEnabled || redisClient == nil {
		return nil
	}
	return redisClient.Del(ctx, key).Err()
}

// RedisFlush menghapus semua data
func RedisFlush(ctx context.Context) error {
	if !redisEnabled || redisClient == nil {
		return nil
	}
	return redisClient.FlushDB(ctx).Err()
}

// =================================================================
// TAMBAHAN PENTING (Agar Session Manager bisa jalan)
// =================================================================

// GetRedisClient mengembalikan instance raw redis client.
// Digunakan oleh session/manager.go untuk inisialisasi redisstore.
func GetRedisClient() *redis.Client {
	return redisClient
}

// IsRedisEnabled mengecek apakah redis aktif dari luar package.
func IsRedisEnabled() bool {
	return redisEnabled
}

// PublishNotificationEvent mempublish event notifikasi ke Redis queue (RPUSH).
// Notification service akan BLPOP dan memproses event ini secara async.
// Mengembalikan error jika Redis tidak aktif atau push gagal.
func PublishNotificationEvent(eventType, channel, to string, vars map[string]string) error {
	if !redisEnabled || redisClient == nil {
		return fmt.Errorf("redis tidak aktif")
	}
	payload, err := json.Marshal(map[string]any{
		"event_type": eventType,
		"channel":    channel,
		"to":         to,
		"vars":       vars,
	})
	if err != nil {
		return fmt.Errorf("marshal payload gagal: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := redisClient.RPush(ctx, "notification:events", payload).Err(); err != nil {
		return fmt.Errorf("RPUSH gagal: %w", err)
	}
	return nil
}
