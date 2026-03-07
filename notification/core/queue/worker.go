package queue

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	template "notification/app/modules/template_notification"

	"github.com/redis/go-redis/v9"
)

// QueueKey adalah nama key di Redis untuk antrian event notifikasi.
const QueueKey = "notification:events"

// StartWorker membaca event dari Redis queue (RPUSH/BLPOP) dan memprosesnya via service.
// Berjalan selamanya sampai ctx dibatalkan (graceful shutdown).
func StartWorker(ctx context.Context, svc *template.Service) {
	addr := getEnv("REDIS_ADDR", "127.0.0.1:6379")
	rdb := redis.NewClient(&redis.Options{
		Addr:        addr,
		DialTimeout: 5 * time.Second,
	})
	defer rdb.Close()

	// Test koneksi awal
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[QUEUE WORKER] Redis tidak tersedia (%v) — queue worker tidak aktif", err)
		return
	}

	log.Printf("[QUEUE WORKER] Started — listening on %s key=%s", addr, QueueKey)

	for {
		// BLPOP: blocking pop dari kiri (FIFO dengan RPUSH dari producer)
		// Timeout 5 detik agar bisa cek ctx.Done()
		result, err := rdb.BLPop(ctx, 5*time.Second, QueueKey).Result()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("[QUEUE WORKER] Stopped (context cancelled)")
				return
			default:
				// Timeout biasa — lanjut loop
				continue
			}
		}

		// result[0] = key name, result[1] = payload
		if len(result) < 2 {
			continue
		}

		var req template.SendRequest
		if err := json.Unmarshal([]byte(result[1]), &req); err != nil {
			log.Printf("[QUEUE WORKER] Payload tidak valid, dibuang: %v | raw: %s", err, result[1])
			continue
		}

		log.Printf("[QUEUE WORKER] Memproses event=%s channel=%s to=%s", req.EventType, req.Channel, req.To)
		if err := svc.Send(req); err != nil {
			// Send() sudah menyimpan log dengan status failed & jadwal retry.
			// Worker tidak perlu melakukan apa-apa lagi — retry worker akan handle.
			log.Printf("[QUEUE WORKER] Send gagal (dijadwalkan retry): %v", err)
		}
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// package utils

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// var (
// 	redisClient  *redis.Client
// 	redisEnabled bool
// )

// // InitRedis menginisialisasi koneksi Redis jika diaktifkan di .env
// func InitRedis() error {
// 	// 1. Cek Feature Flag di .env
// 	if GetEnv("redis") != "true" {
// 		redisEnabled = false
// 		log.Println("[Redis] Feature is DISABLED (Running in No-Cache mode)")
// 		return nil
// 	}

// 	// 2. Ambil Config
// 	host := GetEnv("redis_host", "localhost")
// 	port := GetEnv("redis_port", "6379")
// 	pass := GetEnv("redis_pass", "")

// 	// Konversi DB dari string ke int
// 	dbStr := GetEnv("redis_db", "0")
// 	dbInt, err := strconv.Atoi(dbStr)
// 	if err != nil {
// 		dbInt = 0 // Default ke 0 jika config salah
// 	}

// 	// PERBAIKAN: Gunakan 127.0.0.1 instead of localhost untuk force IPv4
// 	if host == "localhost" {
// 		host = "127.0.0.1"
// 	}

// 	addr := fmt.Sprintf("%s:%s", host, port)

// 	// ✅ FIX: Trim password, jangan kirim empty string sebagai AUTH
// 	pass = strings.TrimSpace(pass)

// 	// 3. Setup Client
// 	redisClient = redis.NewClient(&redis.Options{
// 		Addr:         addr,
// 		Password:     pass, // ✅ empty string = no AUTH command
// 		DB:           dbInt,
// 		DialTimeout:  5 * time.Second,
// 		ReadTimeout:  3 * time.Second,
// 		WriteTimeout: 3 * time.Second,
// 		PoolSize:     10,
// 		MinIdleConns: 5,
// 	})

// 	// 4. Test Ping
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := redisClient.Ping(ctx).Err(); err != nil {
// 		// Jangan fatal, tapi disable redis
// 		log.Printf("[Redis] Connection FAILED: %v (Switching to No-Cache mode)", err)
// 		redisEnabled = false
// 		redisClient = nil
// 		return nil // Return nil agar aplikasi tetap jalan
// 	}

// 	redisEnabled = true
// 	log.Printf("[Redis] Connected successfully to %s", addr)
// 	return nil
// }

// // RedisSet menyimpan data ke cache
// func RedisSet(ctx context.Context, key string, value any, expiration time.Duration) error {
// 	if !redisEnabled || redisClient == nil {
// 		return nil
// 	}
// 	return redisClient.Set(ctx, key, value, expiration).Err()
// }

// // RedisGet mengambil data dari cache
// func RedisGet(ctx context.Context, key string) (string, error) {
// 	if !redisEnabled || redisClient == nil {
// 		return "", fmt.Errorf("redis disabled")
// 	}
// 	return redisClient.Get(ctx, key).Result()
// }

// // RedisDel menghapus data
// func RedisDel(ctx context.Context, key string) error {
// 	if !redisEnabled || redisClient == nil {
// 		return nil
// 	}
// 	return redisClient.Del(ctx, key).Err()
// }

// // RedisFlush menghapus semua data
// func RedisFlush(ctx context.Context) error {
// 	if !redisEnabled || redisClient == nil {
// 		return nil
// 	}
// 	return redisClient.FlushDB(ctx).Err()
// }

// // =================================================================
// // TAMBAHAN PENTING (Agar Session Manager bisa jalan)
// // =================================================================

// // GetRedisClient mengembalikan instance raw redis client.
// // Digunakan oleh session/manager.go untuk inisialisasi redisstore.
// func GetRedisClient() *redis.Client {
// 	return redisClient
// }

// // IsRedisEnabled mengecek apakah redis aktif dari luar package.
// func IsRedisEnabled() bool {
// 	return redisEnabled
// }
