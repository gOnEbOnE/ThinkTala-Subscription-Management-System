# Notification Service README (Onboarding for Live Coding)

Dokumen ini dibuat untuk kamu yang belum tahu apa-apa, lalu mendapat tugas membuat fitur di notification service.
Fokusnya adalah: mulai dari mana, paham alur, dan langkah kerja yang aman.

## 1. Gambaran singkat
Notification service mengurus:
- Broadcast notification (announcement/news) untuk user
- Template notifikasi per event_type dan channel (email/telegram)
- Log pengiriman dan retry
- Link akun Telegram ke user

Service ini bersifat event-driven: service lain mengirim event, notification service mencari template dan mengirimkan pesan.

## 2. Prasyarat minimal
Pastikan ini sudah ada di mesin kamu:
- Go 1.21+
- PostgreSQL
- Redis

Baca juga:
- QUICK_START.md (setup umum)
- README.md (overview sistem)

## 3. Lokasi file penting
Mulai baca dari sini:
- notification/main.go: entry point service
- notification/app/routes/router.go: daftar endpoint
- notification/app/modules/notification_broadcast: CRUD broadcast notif
- notification/app/modules/template_notification: template, send, logs
- notification/app/modules/telegram_link: link/unlink chat_id
- notification/core/database/migrate.go: tabel dan index
- notification/core/queue/worker.go: event queue Redis
- notification/core/utils/smtp.go: email sender (Brevo)
- notification/core/utils/telegram.go: Telegram sender
- notification/LESSONS_LEARNED.md: ringkasan desain dan alasan arsitektur

## 4. Cara menjalankan notification service
Langkah paling aman untuk lokal:

1. Siapkan environment variables di file notification/.env
   Minimal yang dibutuhkan:
   - DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME, DB_SSLMODE
   - PORT (default 5003)

   Jika pakai Redis queue:
   - REDIS_ADDR atau redis_host + redis_port + redis_pass + redis_db

   Jika pakai email:
   - BREVO_API_KEY
   - BREVO_FROM_EMAIL
   - BREVO_FROM_NAME

   Jika pakai Telegram:
   - TELEGRAM_BOT_TOKEN

2. Jalankan service:
   - cd notification
   - go run main.go

Saat startup, service akan:
- konek DB
- migrate dan seed tabel
- register routes
- menjalankan retry worker dan queue worker
- mulai Telegram listener

Jika ingin menjalankan semua service sekaligus, lihat QUICK_START.md.

## 5. Alur data yang perlu kamu pahami

A. Template-based event send
1) Service lain POST ke /api/notifications/send
2) Notification service cari template by event_type + channel
3) Placeholder diganti dari vars
4) Log disimpan (pending)
5) Kirim via email atau telegram
6) Log di-update sent atau failed
7) Retry worker menangani failure

B. Broadcast notification
1) Admin/ops membuat notification via /api/notifications
2) Data disimpan ke tabel notifications
3) Jika aktif, notif bisa dikirim ke Telegram (broadcast)

C. Telegram link
1) User mendapatkan chat_id dari bot
2) User POST /api/telegram/link dengan X-User-ID
3) chat_id disimpan ke users.telegram_chat_id

## 6. Endpoint ringkas (untuk live coding)

Broadcast notifications:
- GET /api/notifications
- GET /api/notifications/public
- GET /api/notifications/recent
- POST /api/notifications/recent/read
- POST /api/notifications/recent/read-all
- GET /api/notifications/:id
- POST /api/notifications
- PUT /api/notifications/:id
- DELETE /api/notifications/:id

Template dan send:
- GET /api/help/notification-templates
- GET /api/help/notification-templates/:id
- POST /api/help/notification-templates
- PUT /api/help/notification-templates/:id
- DELETE /api/help/notification-templates/:id
- GET /api/help/notification-templates/event-types
- POST /api/notifications/send

Telegram link:
- POST /api/telegram/link
- DELETE /api/telegram/unlink
- GET /api/telegram/status

## 7. Payload penting (supaya tidak bingung saat live)

A. Send event
POST /api/notifications/send

{
  "event_type": "otp_verification",
  "channel": "email",
  "to": "user@mail.com",
  "vars": {
    "name": "Rifqi",
    "otp": "123456"
  },
  "user_id": "uuid",
  "user_name": "Rifqi"
}

B. Template baru
POST /api/help/notification-templates

{
  "name": "OTP Verification",
  "event_type": "otp_verification",
  "channel": "email",
  "subject": "Kode OTP",
  "content": "Halo {{name}}, kode OTP kamu {{otp}}",
  "created_by": "admin"
}

C. Link Telegram
POST /api/telegram/link
Header: X-User-ID: <user_id>
Body:
{
  "chat_id": "123456789"
}

## 8. Cara mulai kalau diminta membuat fitur

Checklist yang aman:
1) Jelaskan fitur dalam 1 kalimat: input, output, dan user/role
2) Cari endpoint paling dekat di router.go
3) Tentukan modul mana yang diubah:
   - notification_broadcast untuk news/announcement
   - template_notification untuk event-based send
   - telegram_link untuk link/unlink
4) Cek types.go agar payload jelas
5) Lihat service.go untuk aturan bisnis
6) Jika perlu DB change, edit migrate.go
7) Tambahkan log dan error handling
8) Jalankan service dan test endpoint

## 9. Pola desain yang dipakai
- Controller: parsing HTTP dan validasi input
- Service: rules bisnis
- Repository: query DB

Jangan taruh logic bisnis di controller atau main.go.

## 10. Troubleshooting cepat
- Template tidak ditemukan: cek tabel notification_templates dan event_type
- Redis error: pastikan redis_host/port atau REDIS_ADDR benar
- Email gagal: cek BREVO_API_KEY
- Telegram gagal: cek TELEGRAM_BOT_TOKEN dan chat_id sudah terhubung
- users.telegram_chat_id tidak ada: migrate otomatis saat service startup

## 11. Saran alur cerita saat live coding
- Mulai dari requirement dan alur data
- Tunjukkan router, lalu turun ke module
- Jelaskan validasi di service
- Tunjukkan logging dan retry
- Tutup dengan test cepat (curl atau postman)
