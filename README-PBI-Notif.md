# Modul Notifikasi & Monitoring (Notification Engine)

## 🎯 Tujuan
Menjamin ketersediaan dan ketepatan waktu penyampaian informasi strategis kepada pengguna melalui kanal komunikasi digital (Email, WhatsApp, Telegram).

---

## ✅ Acceptance Criteria

| # | Criteria | Deskripsi |
|---|----------|----------|
| AC-1 | Kelola Template Pesan | Operasional dapat membuat, mengedit, hapus template untuk berbagai channel (Email/WA/Telegram) dengan variable placeholder |
| AC-2 | Monitor Antrian Pesan | Operasional dapat melihat daftar pesan pending dengan filter, sort, dan info status pengiriman |
| AC-3 | Lihat Status Pengiriman | Operasional dapat melihat laporan pengiriman (Sent/Failed) dengan detail error dan statistik delivery rate |

---

## 📊 Metrik Kesuksesan

| Metrik | Target | KPI |
|--------|--------|-----|
| **Reliabilitas Pengiriman** | < 5% failure | Success rate ≥ 95% |
| **Template Flexibility** | 0 code changes | Perubahan template tanpa redeploy |
| **Delivery Accuracy** | 100% | Pesan sesuai template & variables |


---

## 🏗️ Arsitektur Sistem

```
┌─────────────────────────────────────────────────────┐
│              NOTIFICATION SYSTEM                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────┐      ┌────────────────────┐  │
│  │ Application      │      │ Dashboard          │  │
│  │ (Request)        │      │ (Operasional)      │  │
│  └────────┬─────────┘      └────────┬───────────┘  │
│           │                         │              │
│           └─────────────┬───────────┘              │
│                         │                         │
│         ┌───────────────▼────────────────┐        │
│         │ Notification Service           │        │
│         │ • Validasi template            │        │
│         │ • Render variabel              │        │
│         │ • Kirim notifikasi             │        │
│         └───────────────┬────────────────┘        │
│                         │                         │
│         ┌───────────────┼───────────────┐         │
│         │               │               │         │
│    ┌────▼─────┐   ┌────▼────┐   ┌───▼──────┐   │
│    │ Email    │   │ WhatsApp │   │ Telegram │   │
│    │ Provider │   │ Provider │   │ Provider │   │
│    └──────────┘   └──────────┘   └──────────┘   │
│                                                     │
│    ┌──────────────────────────────────────────┐   │
│    │ Database                                 │   │
│    │ • templates                              │   │
│    │ • messages                               │   │
│    │ • notification_logs                      │   │
│    └──────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## 📋 Komponen Utama

| Komponen | Fungsi |
|----------|--------|
| **Notification Service** | Menerima request, validasi template, render variabel, kirim ke provider |
| **Template Engine** | Simpan & kelola template dengan variabel placeholder {{name}}, {{date}}, dll |
| **Email Provider** | Integrasi SendGrid/SMTP untuk pengiriman email |
| **WhatsApp Provider** | Integrasi Twilio untuk WhatsApp |
| **Telegram Provider** | Integrasi Telegram Bot API |
| **Dashboard** | UI untuk operasional monitor dan kelola template |

---

## 🔄 Alur Kerja

### Workflow 1: Pengiriman Notifikasi

```
┌──────────────────────────────────────┐
│ 1. Application Request Notifikasi    │
│    (User register, reset password)   │
└─────────────────┬────────────────────┘
                  │
┌─────────────────▼────────────────────┐
│ 2. Notification Service Menerima     │
│    • recipient: user@example.com      │
│    • template_id: welcome_email       │
│    • variables: {name: "John"}        │
└─────────────────┬────────────────────┘
                  │
┌─────────────────▼────────────────────┐
│ 3. Validasi Template                 │
│    • Template ada di database?        │
│    • Semua variabel tersedia?         │
│    • Format valid?                    │
└─────────────────┬────────────────────┘
                  │
          ┌───────┴────────┐
          │                │
      ┌───▼────┐      ┌───▼────────┐
      │ Valid  │      │ Invalid    │
      └───┬────┘      └────┬───────┘
          │                │
    ┌─────▼─────┐    ┌────▼──────────┐
    │ Render    │    │ Return Error   │
    │ Template  │    │ Don't send     │
    └─────┬─────┘    └────────────────┘
          │
┌─────────▼────────────────────────────┐
│ 4. Generate Pesan                     │
│    • Replace {{name}} → "John"        │
│    • Replace {{date}} → "16/2/2026"   │
│    • Save ke messages table           │
├─────────────────────────────────────┤
│ 5. Kirim via Provider                │
│    • Email → SendGrid API             │
│    • WA → Twilio API                  │
│    • Telegram → Bot API               │
└─────────────────┬────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
    ┌───▼───┐           ┌───▼──────┐
    │ Sent  │           │ Failed   │
    └───┬───┘           └───┬──────┘
        │                   │
    ┌───▼────────────┐  ┌───▼──────────────┐
    │ Save log:      │  │ Save log:        │
    │ status: SENT   │  │ status: FAILED   │
    │ timestamp      │  │ error message    │
    └────────────────┘  └──────────────────┘
```

### Workflow 2: Operasional Monitoring Antrian

```
┌────────────────────────────────┐
│ 1. Operasional Login Dashboard  │
└─────────────┬──────────────────┘
              │
┌─────────────▼──────────────────┐
│ 2. Klik Menu Monitoring         │
└─────────────┬──────────────────┘
              │
┌─────────────▼──────────────────┐
│ 3. Query Message Queue          │
│    SELECT * FROM messages       │
│    WHERE status = 'PENDING'     │
└─────────────┬──────────────────┘
              │
┌─────────────▼──────────────────┐
│ 4. Tampilkan Daftar             │
│    • Message ID                 │
│    • Recipient (email/phone)    │
│    • Channel (Email/WA/TG)      │
│    • Status                     │
│    • Scheduled time             │
├────────────────────────────────┤
│ 5. Operasional Aksi:            │
│    • View Detail                │
│    • Filter by channel/date     │
│    • Cari by recipient          │
└────────────────────────────────┘
```

### Workflow 3: Kelola Template

```
┌──────────────────────────────────┐
│ 1. Operasional Login → Menu       │
│    Notifikasi → Kelola Template   │
└──────────────┬───────────────────┘
               │
┌──────────────▼───────────────────┐
│ 2. Tampilkan Daftar Template      │
│    • Template Name                │
│    • Channel (Email/WA/TG)        │
│    • Status (Active/Inactive)     │
└──────────────┬───────────────────┘
               │
         ┌─────┴─────────┬────────────┐
         │               │            │
    ┌────▼────┐  ┌──────▼────┐  ┌──▼────┐
    │ Tambah  │  │ Edit      │  │ Hapus │
    └────┬────┘  └──────┬────┘  └──┬───┘
         │               │           │
         └───────┬───────┘           │
                 │                   │
         ┌───────▼──────────────┐   │
         │ 3. Form Template    │    │
         │ • Name              │    │
         │ • Channel           │    │
         │ • Subject (Email)   │    │
         │ • Body/Content      │    │
         │ • Variables: {{}}   │    │
         └───────┬──────────────┘   │
                 │                   │
         ┌───────▼──────────────┐   │
         │ 4. Preview Template  │   │
         │ • Render sample data │   │
         │ • Cek formatting     │   │
         └───────┬──────────────┘   │
                 │                   │
         ┌───────▼──────────────┐   │
         │ 5. Save/Update/Delete    │
         │ → Database          │    │
         └───────┬──────────────┘   │
                 │                   │
         ┌───────▼──────────────┐   │
         │ 6. Success Message   │   │
         └──────────────────────┘   │
                                     │
         ┌──────────────────────────┘
         │
    ┌────▼────┐
    │ End     │
    └─────────┘
```

---

## 📡 API Notification Service

### 1. Send Notification

```
POST /api/notifications/send

Request:
{
  "recipient": "user@example.com",
  "channel": "email",          // email | whatsapp | telegram
  "template_id": "welcome",
  "variables": {"name": "John Doe"}
}

Response:
{
  "status": "success",
  "message_id": "msg_123",
  "message": "Notification sent successfully"
}
```

### 2. Get Message Status

```
GET /api/notifications/{message_id}

Response:
{
  "message_id": "msg_123",
  "status": "sent",           // sent | failed | pending
  "recipient": "user@example.com",
  "channel": "email",
  "sent_at": "2026-02-16T10:00:00Z",
  "error": null
}
```

### 3. Dashboard APIs

```
GET /api/templates
GET /api/templates/{id}
POST /api/templates
PUT /api/templates/{id}
DELETE /api/templates/{id}

GET /api/messages?status=pending
GET /api/messages/report?date=2026-02-16
```



## 🔌 Implementasi



// Usage:
// Template: "Hello {{.name}}, your OTP is {{.otp}}"
// Variables: {"name": "John", "otp": "123456"}
// Result: "Hello John, your OTP is 123456"
```

---

## 📊 Dashboard Fitur

### 1. Template Management
- **List Templates**: Tampilkan semua template dengan channel & status
- **Create Template**: Form untuk buat template baru
- **Edit Template**: Update template yang ada
- **Delete Template**: Hapus template
- **Preview**: Lihat hasil template dengan sample data

### 2. Message Monitoring
- **List Pending Messages**: Tampilkan pesan yang belum terkirim
- **Filter**: By channel, recipient, date, status
- **Sort**: By created_at, scheduled_at
- **View Detail**: Lihat isi pesan & metadata lengkap
- **Retry**: Manual retry untuk pesan gagal

### 3. Delivery Report
- **Summary Stats**: Total sent, failed, success rate
- **By Channel**: Breakdown per Email/WA/Telegram
- **By Status**: Sent vs Failed count
- **Error Analysis**: Jenis-jenis error yang terjadi
- **Export**: Export ke CSV/PDF

---

## 🎬 Contoh Use Case

### Use Case 1: User Register
```
1. User daftar akun
2. System trigger notification
3. Request: 
   - channel: email
   - template: welcome_email
   - recipient: user@email.com
   - variables: {name: "John"}
4. Notification Service:
   - Validasi template welcome_email ada
   - Render: "Welcome John!"
   - Kirim ke SendGrid
5. Response: Sent / Failed
```

### Use Case 2: Reset Password
```
1. User forget password
2. System trigger notification
3. Request:
   - channel: whatsapp
   - template: reset_otp
   - recipient: +62812xxxx
   - variables: {otp: "123456", expires: "10 min"}
4. Notification Service:
   - Render: "Kode OTP: 123456, berlaku 10 min"
   - Kirim ke Twilio
5. Save log & track status
```

---

## ✅ Checklist Implementasi

### Phase 1: Core Setup (Week 1)
- [ ] Database schema
- [ ] API endpoints
- [ ] Email provider

### Phase 2: Multi-Channel (Week 2)
- [ ] WhatsApp provider
- [ ] Telegram provider

### Phase 3: Dashboard (Week 3)
- [ ] Template UI
- [ ] Monitoring page
- [ ] Report page

### Phase 4: Testing & Deploy (Week 4)
- [ ] Testing
- [ ] Deployment

---


