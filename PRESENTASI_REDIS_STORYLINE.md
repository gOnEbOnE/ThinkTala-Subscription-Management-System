# Presentasi Redis Thinknalyze

## 1. Tujuan Presentasi

1. Menjelaskan alasan penggunaan Redis pada arsitektur Thinknalyze.
2. Menjelaskan implementasi Redis dari sisi autentikasi, session, dan notifikasi asynchronous.
3. Menunjukkan nilai bisnis, risiko teknis, serta rekomendasi perbaikan.

---

## 2. Storyline Utama

1. Opening: tantangan microservices pada autentikasi dan pengiriman notifikasi.
2. Problem statement: tanpa storage state cepat, validasi token mahal dan notifikasi bisa menghambat flow utama.
3. Solusi: Redis sebagai token store, session backend opsional, dan queue notifikasi.
4. Implementasi nyata: alur login, validasi gateway, publish event, consume worker, retry.
5. Dampak: latency turun, decoupling meningkat, reliability notifikasi membaik.
6. Trade-off: Redis menjadi dependency kritikal sehingga perlu hardening dan observability.
7. Penutup: Redis sudah tepat, fokus berikutnya adalah standardisasi lintas service.

---

## 3. Struktur Presentasi (Slide by Slide)

## Slide 1 - Judul dan Konteks

Tujuan slide:
- Memberi konteks dan ekspektasi audiens.

Isi slide:
- Judul: Redis Implementation in Thinknalyze
- Subjudul: Auth, Session, and Async Notification
- Durasi presentasi: 15-20 menit

Speaker notes:
- Presentasi ini fokus pada implementasi Redis di codebase, bukan teori umum semata.
- Kita akan lihat alur teknis dan dampaknya ke reliability sistem.

Transisi:
- Sebelum membahas solusi, kita lihat masalah utama yang ingin diselesaikan.

## Slide 2 - Masalah Utama yang Dihadapi

Tujuan slide:
- Menjelaskan pain point yang mendorong penggunaan Redis.

Isi slide:
- Validasi auth antar service perlu cepat.
- Session state perlu bisa di-revoke dari server.
- Notifikasi tidak boleh blocking proses register dan KYC.
- Sistem harus tetap tahan saat ada gangguan dependency.

Speaker notes:
- Pada arsitektur microservices, state management yang lambat akan langsung terasa di pengalaman pengguna.
- Tanpa decoupling notifikasi, request user bisa tertahan menunggu pengiriman email.

Transisi:
- Redis dipilih untuk menutup gap ini. Selanjutnya kita lihat fundamental yang dipakai.

## Slide 3 - Fundamental Redis yang Dipakai

Tujuan slide:
- Menyatukan pemahaman konsep Redis sesuai implementasi proyek.

Isi slide:
- Redis adalah in-memory key-value store berlatensi rendah.
- Tipe data yang dominan dipakai:
  - String untuk token dan pointer session.
  - List untuk queue notifikasi.
- Operasi inti:
  - SET, GET, DEL
  - RPUSH, BLPOP
  - TTL untuk masa berlaku key

Speaker notes:
- Di Thinknalyze, Redis bukan sekadar cache tampilan, tetapi masuk ke jalur kritikal auth dan messaging.

Transisi:
- Berikut posisi Redis dalam arsitektur sistem.

## Slide 4 - Peran Redis di Arsitektur Thinknalyze

Tujuan slide:
- Memetakan domain Redis dalam sistem.

Isi slide:
- Auth token storage:
  - Browser menyimpan pointer token dalam cookie.
  - JWT raw disimpan di Redis.
- Session backend:
  - Driver session mendukung mode berbasis Redis.
- Notification queue:
  - Users service mempublish event.
  - Notification service mengonsumsi event.

Speaker notes:
- Satu Redis melayani tiga kebutuhan utama. Efisien, namun butuh disiplin key naming dan error handling.

Transisi:
- Sekarang masuk ke flow paling penting, yaitu login sampai authorization.

## Slide 5 - Alur Auth End-to-End

Tujuan slide:
- Menjelaskan pointer-token pattern yang dipakai.

Isi slide:
1. Users service membuat JWT.
2. Users service membuat UUID key, lalu encrypt key untuk cookie.
3. JWT raw disimpan di Redis dengan key UUID.
4. Gateway membaca cookie token.
5. Gateway decrypt key, GET JWT dari Redis, lalu validate JWT.
6. Gateway cek role untuk path yang diakses.

Speaker notes:
- Cookie tidak menyimpan JWT asli, hanya referensi terenkripsi.
- Keuntungan utama: revoke session bisa dilakukan cepat cukup dengan DEL key Redis.

Transisi:
- Setelah auth, kita lihat bagaimana session driver memanfaatkan Redis.

## Slide 6 - Session Driver dan Fallback

Tujuan slide:
- Menjelaskan mode session serta perilaku fallback.

Isi slide:
- Driver session mendukung:
  - cookie
  - redis
  - stateless
  - stateless_with_redis
- Pada mode stateless_with_redis, data session disimpan dengan key sess:<sid>.
- Jika driver butuh Redis tetapi Redis nonaktif, session manager fallback ke cookie.

Speaker notes:
- Fallback membantu availability, tetapi harus distandardisasi agar perilaku antar service tetap konsisten.

Transisi:
- Selanjutnya jalur asynchronous untuk notifikasi.

## Slide 7 - Queue Notifikasi dengan Redis

Tujuan slide:
- Menjelaskan decoupling antar service menggunakan Redis list.

Isi slide:
- Producer:
  - Users service RPUSH ke key notification:events.
- Consumer:
  - Notification worker BLPOP dari key yang sama.
- Worker parse payload, render template, lalu kirim notifikasi.

Speaker notes:
- Dengan pola ini, register dan KYC tidak perlu menunggu proses kirim email selesai.
- Flow utama menjadi lebih responsif dan lebih tahan beban.

Transisi:
- Asynchronous harus punya mekanisme ketahanan. Kita bahas retry dan fallback.

## Slide 8 - Reliability: Retry dan Fallback

Tujuan slide:
- Menunjukkan strategi ketahanan saat gagal kirim.

Isi slide:
- Worker queue retry koneksi Redis bila Redis belum ready saat startup.
- Gagal kirim akan dicatat sebagai failed dan dijadwalkan retry.
- Backoff retry bertahap:
  - 1 menit
  - 5 menit
  - 30 menit
- Producer-side fallback:
  - queue gagal -> HTTP direct ke notification service -> SMTP fallback

Speaker notes:
- Strategi berlapis ini mencegah notifikasi hilang begitu saja saat ada gangguan sementara.

Transisi:
- Setelah reliability, berikut sisi security dan governance.

## Slide 9 - Security dan Governance

Tujuan slide:
- Menjelaskan implikasi keamanan dari desain Redis.

Isi slide:
- Cookie menyimpan pointer terenkripsi, bukan JWT raw.
- JWT tetap diverifikasi signature sebelum dipakai.
- Role authorization dilakukan di gateway.
- Logout menghapus key Redis agar sesi tidak bisa digunakan lagi.

Speaker notes:
- Redis memegang data sensitif, jadi perlu hardening: password, network isolation, dan monitoring akses.

Transisi:
- Sekarang kondisi operasional aktual di project ini.

## Slide 10 - Kondisi Operasional Saat Ini

Tujuan slide:
- Menjelaskan dependency runtime Redis.

Isi slide:
- Orchestrator mengasumsikan Redis sudah berjalan secara eksternal.
- Docker Compose utama belum memuat service Redis.
- Konfigurasi Redis dilakukan per service melalui env.

Speaker notes:
- Ini berarti tim devops perlu SOP provisioning dan health check Redis yang jelas.

Transisi:
- Berikut temuan teknis yang perlu ditindaklanjuti.

## Slide 11 - Temuan Teknis dan Risiko

Tujuan slide:
- Menyajikan evaluasi objektif.

Isi slide:
- Ada perbedaan pola handling token antara flow modern dan flow legacy.
- Ada potensi inkonsistensi pembersihan key backup session saat logout.
- Error handling Redis belum sepenuhnya seragam di semua service.
- Risiko utama: Redis outage berdampak ke auth dan queue sekaligus.

Speaker notes:
- Sistem sudah berjalan baik, tetapi konsistensi lintas service masih menjadi area perbaikan.

Transisi:
- Berikut rekomendasi prioritas.

## Slide 12 - Rekomendasi Prioritas

Tujuan slide:
- Memberikan langkah tindak lanjut yang actionable.

Isi slide:
1. Standardisasi policy error handling Redis antar service.
2. Samakan format token handling, kurangi jalur legacy yang tidak relevan.
3. Rapikan lifecycle key saat logout agar bersih secara konsisten.
4. Tambahkan observability:
   - queue depth
   - auth Redis hit/miss
   - retry count
   - alerting saat failure meningkat

Speaker notes:
- Ini langkah praktis yang memberi dampak besar terhadap stabilitas dan troubleshooting.

Transisi:
- Kita simpulkan nilai Redis di sistem ini.

## Slide 13 - Kesimpulan

Tujuan slide:
- Menutup dengan pesan utama.

Isi slide:
- Redis sudah menjadi komponen inti, bukan pelengkap.
- Tiga dampak utama:
  - auth state cepat
  - session management fleksibel
  - async notification lebih resilient
- Tahap selanjutnya: standardisasi dan observability untuk kesiapan production skala lebih besar.

Speaker notes:
- Pemilihan Redis sudah tepat untuk kebutuhan sistem saat ini.
- Fokus berikutnya adalah penyempurnaan operasional dan konsistensi implementasi.

Transisi:
- Lanjut demo singkat dan tanya jawab.

## Slide 14 - Demo dan Q&A

Tujuan slide:
- Membuktikan alur berjalan end-to-end.

Isi slide:
1. Demo login lalu akses route berbasis role.
2. Trigger event register atau KYC dan tunjukkan queue diproses.
3. Simulasi gangguan endpoint notifikasi dan tunjukkan fallback.
4. Q&A.

Speaker notes:
- Demo difokuskan pada ketahanan flow, bukan hanya happy path.

---

## 4. Script Narasi Ringkas (Bisa Dibaca Langsung)

Pembuka:
- Hari ini saya akan menjelaskan implementasi Redis di Thinknalyze dari sisi arsitektur dan operasional. Fokusnya adalah bagaimana Redis dipakai untuk menyelesaikan masalah autentikasi, session, dan notifikasi asynchronous.

Bagian inti:
- Redis di Thinknalyze memiliki tiga fungsi utama.
- Pertama, Redis menyimpan JWT raw yang direferensikan oleh cookie pointer terenkripsi.
- Kedua, Redis menjadi backend session pada mode tertentu.
- Ketiga, Redis menjadi queue event notifikasi dengan pola producer-consumer.

Nilai yang dihasilkan:
- Validasi auth menjadi cepat.
- Revocation session bisa dilakukan server-side.
- Register dan KYC menjadi lebih responsif karena pengiriman notifikasi dipindah ke background.

Ketahanan:
- Sistem menerapkan retry worker dan fallback bertingkat agar notifikasi tidak berhenti total saat terjadi gangguan sementara.

Penutup:
- Redis implementation ini sudah memberi fondasi kuat. Langkah berikutnya adalah menseragamkan perilaku lintas service dan menambah observability agar siap skala production.

---

## 5. Rencana Demo (5-7 Menit)

1. Login dari UI dan verifikasi akses berdasarkan role.
2. Trigger event notifikasi dari flow register atau KYC.
3. Tunjukkan worker mengonsumsi queue notification:events.
4. Simulasikan gangguan notification endpoint dan tunjukkan fallback.
5. Tutup dengan ringkasan reliability behavior.

---

## 6. Checklist Sebelum Presentasi

1. Redis server aktif.
2. Env redis=true pada service yang terkait.
3. Gateway, users, dan notification service berjalan.
4. Template notifikasi tersedia di notification service.
5. Jalur demo sudah diuji satu kali sebelum sesi.

---

## 7. Lampiran Referensi Kode Utama

### Inisialisasi Redis
- [users/main.go](users/main.go#L31)
- [account/main.go](account/main.go#L39)
- [subscription/main.go](subscription/main.go#L26)
- [gateway/main.go](gateway/main.go#L383)

### Auth pointer-token
- [users/core/token/token.go](users/core/token/token.go#L32)
- [users/core/token/token.go](users/core/token/token.go#L57)
- [gateway/auth/auth.go](gateway/auth/auth.go#L183)
- [gateway/auth/auth.go](gateway/auth/auth.go#L200)

### Queue notifikasi
- [users/core/utils/redis.go](users/core/utils/redis.go#L130)
- [users/core/utils/redis.go](users/core/utils/redis.go#L145)
- [notification/core/queue/worker.go](notification/core/queue/worker.go#L18)
- [notification/core/queue/worker.go](notification/core/queue/worker.go#L65)

### Retry dan fallback
- [notification/app/modules/template_notification/service.go](notification/app/modules/template_notification/service.go#L119)
- [notification/app/modules/template_notification/service.go](notification/app/modules/template_notification/service.go#L135)
- [users/app/modules/register/service.go](users/app/modules/register/service.go#L42)
- [users/app/modules/kyc/service.go](users/app/modules/kyc/service.go#L47)

### Catatan operasional
- [main.go](main.go#L85)
- [main.go](main.go#L180)
- [docker-compose.yml](docker-compose.yml)

### Catatan konsistensi implementasi
- [users/app/modules/login/service.go](users/app/modules/login/service.go#L83)
- [subscription/app/modules/login/service.go](subscription/app/modules/login/service.go#L59)
