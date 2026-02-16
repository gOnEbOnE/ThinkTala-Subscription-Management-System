# Thinknalyze - Microservices Ecosystem

**Thinknalyze** adalah ekosistem microservices berbasis Golang yang dirancang dengan arsitektur **API Gateway**. Proyek ini menggunakan Docker Compose untuk orkestrasi kontainer, memungkinkan skalabilitas tinggi dan manajemen layanan yang terisolasi.

---


## 🏗️ Arsitektur Layanan

Sistem ini terdiri dari satu gerbang utama (Gateway) dan lima layanan mikro internal:

| Service | Port Internal | Deskripsi |
| :--- | :--- | :--- |
| **Gateway** | **2000** | Entry point (HTTPS), Reverse Proxy, SSL, & CORS |
| **Notification** | 2001 | Layanan pesan dan notifikasi |
| **Operational** | 2002 | Layanan manajemen operasional & dashboard |
| **Subscription** | 2003 | Layanan manajemen paket dan aplikasi |
| **Tickets** | 2004 | Layanan dukungan dan ticketing |
| **Users** | 2005 | Layanan autentikasi dan profil pengguna |

---

## 📂 Struktur Folder

```text
thinknalyze/
├── gateway/           # Core API Gateway
│   ├── .env           # Gateway configuration (Required: port, timezone)
│   ├── routes.json    # Dynamic proxy configuration
│   └── Dockerfile
├── notification/      # Microservice Notification
├── operational/       # Microservice Operational
├── subscription/      # Microservice Subscription
├── tickets/           # Microservice Tickets
├── users/             # Microservice Users
├── docker-compose.yml # Docker orchestration
└── README.md          # Project documentation
```

---

## 🚀 Persiapan Instalasi

### 1. Prasyarat
* **Docker** & **Docker Compose V2** terinstal di server.

### 2. Konfigurasi Lingkungan (.env)
Buat file `.env` di dalam folder `gateway/`:
```env
port=2000
timezone=Asia/Jakarta
```

### 3. Inisialisasi Go Module (Pertama kali)
Jalankan perintah ini di setiap folder (termasuk gateway):
```bash
go mod init thinknalyze/nama-folder
go mod tidy
```

---

## 🛠️ Menjalankan Project

Jalankan seluruh layanan menggunakan Docker Compose dari root folder:

```bash
# Membangun image dan menjalankan kontainer di background
docker compose up -d --build
```

Untuk mematikan layanan:
```bash
docker compose down
```

---

## 🧪 Pengujian & Verifikasi

Setelah kontainer berjalan, verifikasi koneksi menggunakan `curl` atau browser (abaikan peringatan SSL jika menggunakan self-signed):

* **Gateway Health Check:**
  ```bash
  curl -k https://localhost:2000/health
  ```
* **Verify Route (Proxy to Users):**
  ```bash
  curl -k https://localhost:2000/account/register
  ```

---

## 🔄 Konfigurasi Dinamis (Hot Reload)

Anda dapat mengubah rute di `gateway/routes.json` tanpa perlu membangun ulang kontainer. Setelah mengubah file tersebut, cukup panggil endpoint reload:

**Endpoint:** `POST https://localhost:2000/admin/config/reload`

---

## 📝 Cheat Sheet Docker

| Tujuan | Perintah |
| :--- | :--- |
| Lihat semua kontainer | `docker compose ps` |
| Lihat log Gateway | `docker compose logs -f gateway` |
| Restart satu service | `docker compose restart <nama_service>` |
| Build ulang tanpa cache | `docker compose build --no-cache` |

---

## 📄 Lisensi
Copyright © 2026. **Thinknalyze**.  
*Optimized for efficiency and high-concurrency microservices.*