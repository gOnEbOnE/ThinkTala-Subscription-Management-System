# Personal Report – Implementation/Coding
**Proyek Pengembangan Sistem Informasi | Genap 2025/2026**

## 👤 Identitas
- **Nama:** Christopher Matthew Hendarson
- **NIM:** 2306245592
- **Kelas:** C
- **Nama Tim / Proyek:** Propensuy/ThinkNalyze
- **Role dalam Tim:** Implementation / Coding

---

## 🎯 Fokus Kontribusi Utama
Sprint 1 — Integrasi antar modul via API Gateway dan implementasi sistem akun (registrasi & autentikasi)

**PBI / Fitur yang Dikerjakan:**
PBI01-PBI06
Registrasi, login, logout, dashboard masing masing role

---

## 🏗️ Arsitektur & Hubungan Antar Modul

### Gambaran Umum
Project ini menggunakan pendekatan modular microservices di dalam satu monorepo. Dari `docker-compose.yml`, service utama yang berjalan adalah Gateway, Users, Notification, Operational, Subscription, dan Tickets, ditambah PostgreSQL dan Redis.

Untuk kontribusi Sprint 1, alur kritis yang saya kerjakan berfokus pada 3 komponen:
- **Gateway** sebagai pintu masuk tunggal HTTP (`:2000`) dan reverse proxy antar service.
- **Users module** sebagai implementasi inti registrasi, verifikasi OTP, login, dan pembuatan JWT.
- **Account module** sebagai lapisan halaman akun/dashboard dan mekanisme session-authz di sisi aplikasi account.

### Peran Gateway
Gateway memuat konfigurasi route dari `gateway/routes.json`, lalu mendaftarkan handler proxy di `gateway/main.go`.

Peran utamanya:
- Menjadi **single entry point** untuk frontend dan API.
- Melakukan **reverse proxy** ke service target berdasarkan prefix route.
- Menjalankan **role-based authorization** sebelum meneruskan request ke service backend tertentu.
- Mengelola **session cookie `token`** untuk logout/redirect, dan meneruskan identitas user sebagai header downstream (`X-User-Role`, `X-User-ID`, `X-User-Email`).

Contoh routing yang ditemukan:
- `/api/auth/*` -> Users service (`localhost:2006` pada config gateway saat local mapping)
- `/account/login/auth` -> Users service
- `/api/kyc/*`, `/api/admin/kyc*` -> endpoint KYC pada service Users (dengan role check di gateway)

### Hubungan Account -> Gateway -> User
Alur autentikasi lintas modul berjalan seperti ini:
1. Client membuka halaman akun melalui gateway (`/account/login`, `/account/register`).
2. Untuk login, client mengirim `POST /account/login/auth` ke gateway.
3. Gateway mem-proxy request ke Users service (`loginController.Auth`).
4. Users service validasi kredensial:
   - baca user dari PostgreSQL,
   - verifikasi password hash,
   - cek status akun aktif,
   - generate JWT (RS256).
5. Users menyimpan JWT mentah di Redis menggunakan key UUID, lalu mengirim cookie `token` berisi key terenkripsi.
6. Untuk request berikutnya ke route protected (`/client/*`, `/ops/*`, `/compliance/*`, `/api/*` tertentu), gateway:
   - ambil cookie `token`,
   - decrypt key,
   - ambil JWT dari Redis,
   - validasi JWT,
   - cek role terhadap path,
   - forward request ke service target jika lolos.

Di sisi Account module, middleware authz juga membaca session token model yang sama (encrypted key -> Redis -> JWT) untuk menjaga konsistensi akses halaman account/dashboard.

### Diagram Alur (Text)
```text
Client
  -> Gateway: GET /account/login
  -> Gateway: POST /account/login/auth
      -> Proxy to Users: Auth()
          -> Verify password + user status in PostgreSQL
          -> Create JWT (RS256)
          -> Store raw JWT in Redis (key=UUID)
          -> Return cookie token (encrypted UUID key)

Client + cookie token
  -> Gateway: GET /client/dashboard or /api/kyc/status
      -> Decrypt token key
      -> Fetch JWT from Redis
      -> Validate JWT + role authorization
      -> Proxy to target service (Users/Operational/etc)
```

---

## 🛠️ Penjelasan Implementasi per Modul

### 1. Gateway
- **File utama:** `gateway/main.go`, `gateway/routes.json`, `gateway/auth/auth.go`
- **Cara kerja:** Gateway dibuat dengan `net/http`, dynamic route config dari JSON, lalu reverse proxy ke target service menggunakan `httputil.NewSingleHostReverseProxy`.
- **Middleware yang digunakan:**
  - CORS middleware (`withCORS`)
  - request logging (`withLogging`)
  - role auth (`withRoleAuth`, `withRolesAuth`) yang validasi token dari cookie + Redis + JWT
- **Route yang dikonfigurasi:**
  - Auth/account routes: `/api/auth/*`, `/account/login/auth`, `/api/auth/logout`, `/api/auth/assume-role`, `/api/auth/roles`
  - User/KYC routes: `/api/kyc/*`, `/api/admin/kyc*`
  - Service lain: `/api/notifications*`, `/api/help/*`, `/api/admin/packages*`, `/api/subscriptions*`, `/api/admin/orders*`, `/api/operational/*`

### 2. Account Module
- **File utama:** `account/main.go`, `account/app/routes/router.go`, `account/app/modules/account/controller.go`, `account/core/token/token.go`, `account/core/http/middleware.go`
- **Fitur yang diimplementasikan:**
  - Registrasi akun: Alur registrasi endpoint API dijalankan oleh Users module, sedangkan Account module menyediakan konteks halaman akun dan proteksi akses halaman berbasis sesi.
  - Login / penggunaan akun: Endpoint login diproxy gateway ke Users module; Account module menggunakan token/session yang sama untuk authorize akses halaman account/dashboard dan melakukan logout yang membersihkan sesi Redis + cookie/session.
- **Logic utama:**
  - Implementasi dipisah agar **gateway tetap stateless di sisi request handling**, sementara state sesi disimpan di Redis.
  - Cookie tidak menyimpan JWT langsung, tetapi menyimpan **encrypted session key**, sehingga jika cookie bocor, token asli tidak langsung terekspos.
  - Verifikasi role dilakukan di gateway sebelum request ke modul tujuan, sehingga akses dashboard per role lebih konsisten.
- **Code snippet utama:**

```go
// users/app/modules/login/controller.go
if loginRes.Token != "" {
    tokenKey, err := token.SetUserAuthz(w, r, loginRes.Token)
    if err != nil {
        c.Response.JSON(w, r, map[string]any{"status": false, "msg": err.Error()})
        return
    }

    // Cookie yang dibaca gateway pada request berikutnya
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    tokenKey,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })
}

// gateway/auth/auth.go
tokenCookie, _ := r.Cookie("token")
decryptedKey, _ := utils.Decrypt(tokenCookie.Value)
jwtRaw, _ := utils.RedisGet(ctx, string(decryptedKey))
claims, _ := utils.ValidateJWT(jwtRaw)
```