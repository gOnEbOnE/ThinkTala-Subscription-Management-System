# Quick Start Guide - 5 Menit Setup

## ✅ Prerequisite Check

```powershell
# Pastikan sudah installed:
go version          # v1.21+
psql --version      # PostgreSQL
redis-cli --version # Redis
```

## 1️⃣ Create `.env` Files (2 menit)

Copy content dari `DOKUMENTASI_SETUP.md` **Environment Configuration** section ke:

### `gateway/.env`
```env
app_name=ZaFramework Gateway
app_env=development
app_secure=false
port=2000
timezone=Asia/Jakarta

redis=true
redis_host=localhost
redis_port=6379
redis_pass=
redis_db=0
```

### `users/.env`
```env
app_name=ZaFramework Users Service
app_env=development
app_secure=false
port=2006
timezone=Asia/Jakarta

postgres=true
read_db_host=localhost
read_db_port=5433
read_db_user=office
read_db_pass=M3l4t1cn
read_db_name=office_db
read_db_ssl_mode=disable
read_db_timezone=Asia/Jakarta

SESSION_DRIVER=cookie
SESSION_KEY=enj93VpBzHN6VULZ36sCOm5Z0wZHGw76FRvxZvNIkU3yKpGmGsI3/5KX71p8YmI/gzCQ2FSGgc1uGCD1A/tmqzIP0oVjy+7QCPHQblxa4TqYDq1wmM1V+HGnpMjnfIk3
SESSION_NAME=za_session_users
SESSION_LIFETIME=604800

redis=true
redis_host=localhost
redis_port=6379
redis_pass=
redis_db=0

JWT_EXPIRED=9600
TOKEN_ENCRYPTION_KEY=your-32-char-encryption-key-here!
```

### `account/.env`
```env
app_name=ZaFramework Account Service
app_env=development
app_secure=false
port=2001
timezone=Asia/Jakarta

postgres=true
read_db_host=localhost
read_db_port=5433
read_db_user=office
read_db_pass=M3l4t1cn
read_db_name=office_db
read_db_ssl_mode=disable
read_db_timezone=Asia/Jakarta
read_db_max_conn=10

SESSION_DRIVER=cookie
SESSION_KEY=enj93VpBzHN6VULZ36sCOm5Z0wZHGw76FRvxZvNIkU3yKpGmGsI3/5KX71p8YmI/gzCQ2FSGgc1uGCD1A/tmqzIP0oVjy+7QCPHQblxa4TqYDq1wmM1V+HGnpMjnfIk3
SESSION_NAME=za_session_account
SESSION_LIFETIME=604800

redis=true
redis_host=localhost
redis_port=6379
redis_pass=
redis_db=0

JWT_EXPIRED=9600
CSRF_AUTH_KEY=a3f9c2e4b71d6a8095f4e1c7d23b8a6c9f0e4d2b7a1c8e5f6b309d41c2e7a8f
```

## 2️⃣ Create `gateway/routes.json` (1 menit)

```json
{
  "allowedOrigins": [
    "http://localhost",
    "http://localhost:2000",
    "http://127.0.0.1:2000"
  ],
  "routes": [
    {
      "path": "/account/login/auth",
      "target": "http://localhost:2006",
      "cors": true,
      "description": "Login authentication proxy to users service"
    },
    {
      "path": "/api/auth",
      "target": "http://localhost:2006",
      "cors": true,
      "description": "Authentication API"
    },
    {
      "path": "/api/auth/register",
      "target": "http://localhost:2006",
      "cors": true,
      "description": "User registration"
    },
    {
      "path": "/api/auth/verify-otp",
      "target": "http://localhost:2006",
      "cors": true,
      "description": "OTP verification"
    },
    {
      "path": "/api/user",
      "target": "http://localhost:2001",
      "cors": true,
      "description": "User profile API"
    }
  ]
}
```

## 3️⃣ Fix Gateway `main.go` (30 detik)

Edit `gateway/main.go`, cari function `main()`, dan add ini setelah `auth.InitRedis()`:

```go
func main() {
    utils.LoadEnv(".env")
    
    if err := utils.InitJWTLoadKeys("certs/private.pem", "certs/public.pem"); err != nil {
        if err2 := utils.InitJWTLoadKeys("../users/certs/private.pem", "../users/certs/public.pem"); err2 != nil {
            log.Fatalf("[FATAL] JWT keys not found: %v / %v", err, err2)
        }
    }
    
    auth.InitRedis()

    // ⭐ ADD THIS (CRITICAL!)
    var err error
    config, err = loadConfig()
    if err != nil {
        log.Fatalf("[FATAL] Failed to load config: %v", err)
    }

    frontendDir := "../frontend"
    // ... rest of code ...
}
```

## 4️⃣ Create `certs` Package (30 detik)

Buat file: `certs/generator.go` (copy dari dokumentasi atau lihat folder reference)

## 5️⃣ Run Everything! (30 detik)

```powershell
cd C:\PROPENSUY\thinknalyze
go run main.go
```

**Expected Output:**
```
=========================================
  🚀 THINKNALYZE ORCHESTRATOR
=========================================

[*] Checking prerequisites...
  ✅ Go - OK
  ✅ PostgreSQL - OK
  ✅ Redis - OK

🔐 Generating certificates for gateway...
✅ Written to gateway/certs/private.pem
✅ Written to gateway/certs/public.pem

[REDIS] Starting Redis server...
[+] Redis started with PID xxxxx

Memulai semua service...

[+] Service account berjalan dengan PID xxxxx
[+] Service gateway berjalan dengan PID xxxxx
[+] Service users berjalan dengan PID xxxxx

=========================================
  ✅ ALL SERVICES STARTED
=========================================

  🌐 Gateway (Port 2000):
     http://localhost:2000

  🔐 Login:
     http://localhost:2000/account/login

  📊 Test Credentials:
     Email: superadmin@thinktala.com
     Pass:  Super123

=========================================
```

## ✅ Done!

Open browser: **`http://localhost:2000/account/login`**

Login dengan credentials di atas → Should redirect ke `/ops/dashboard`

---

## 🆘 Troubleshooting

### ❌ `panic: runtime error: invalid memory address`
→ Pastikan sudah add `config, err = loadConfig()` di gateway/main.go

### ❌ `redis disabled`
→ Pastikan `redis_pass=` (KOSONG) di `.env` files

### ❌ `dial tcp: no such host`
→ PostgreSQL/Redis belum running, jalankan dulu

### ❌ `Infinite redirect loop`
→ Pastikan akses via `http://localhost:2000`, BUKAN port lain

---

**More help:** See `DOKUMENTASI_SETUP.md`