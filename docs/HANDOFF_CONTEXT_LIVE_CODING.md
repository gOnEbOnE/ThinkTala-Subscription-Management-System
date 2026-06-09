# Handoff Context — Persiapan Live Coding Propen A+ (Modul Paket)

> Dokumen ini untuk **agent AI baru** (Claude/Cursor) agar langsung paham konteks tanpa baca seluruh chat history.
> Repo: `/Users/farizmuhammad/Downloads/erp/christopher` (GitLab: propensuy-thinknalyze)

---

## 1. Apa yang sedang dikerjakan?

Mahasiswa mempersiapkan **ujian live coding Propen A+** untuk modul **Paket & Subscription** di project ThinkNalyze/ThinkTala.

**Bukan** tugas implementasi fitur baru besar. **Bukan** deploy production (meski ada pekerjaan UAT sebelumnya).

Fokus utama sesi terakhir:

1. Reverse engineering codebase modul Package, Order, Subscription, Billing, Invoice PDF
2. Menyusun dokumentasi persiapan live coding
3. Menyusun **cheat sheet operasional** yang bisa diikuti kata-per-kata saat ujian 30 menit

---

## 2. Format ujian (WAJIB DIPAHAMI)

| Item | Detail |
|------|--------|
| Durasi | **30 menit** |
| Jumlah soal | **2 wajib** |
| Soal 1 | **C/U/D** — Create ATAU Update ATAU Delete |
| Soal 2 | **R** — Read / List / Detail |
| Scope soal | Hanya 4 area: Package CRUD, Pembayaran Paket, Status Subscription, Invoice PDF |

### Pola soal yang BENAR (dari kelompok lain di kelas yang sama)

Soal **BUKAN**:
- Implementasi endpoint dari nol
- Fix bug / patch gap yang belum ada

Soal **ADALAH** modifikasi kecil bertarget di fitur yang **sudah jalan**, contoh pola nyata kelompok lain:

- "Tambahkan filter datepicker satu tanggal di list news"
- "Ubah urutan audit trail dari terbaru ke terlama"
- "Saat create role, validasi jika permissions sama dengan role lain → tolak + popup error"
- "Tampilkan daftar user saat checkbox Send via Email dicentang"
- "Tambahkan card total projects published bulan ini"

**Implikasi untuk project ini:** soal kemungkinan seperti:
- Tambah kolom di tabel yang sudah ada
- Tambah filter datepicker di billing history
- Tambah validasi kecil di service (nama tanpa karakter spesial)
- Tampilkan data tambahan di modal (subscriber count sebelum delete)
- Sort/badge/format tampilan di frontend

---

## 3. File dokumentasi yang sudah dibuat

| File | Isi | Kapan dipakai |
|------|-----|---------------|
| **`LIVE_CODING_CHEATSHEET.md`** | **UTAMA saat ujian.** §1–§45 + §R-INV/PROOF/VERIFY. Anchor Ctrl+F + kode copy-paste + verifikasi 30 detik | Buka ini saat live coding |
| **`README_LIVE_CODING_PACKAGE.md`** | Referensi lengkap ~4200 baris: Part 1 (peta modul), Part 2 (implementation guide cluster A–L), Part 3 (prediksi TOP 20 C/U/D & R) | Baca sebelum ujian, jangan saat panik 30 menit |

**Aturan untuk agent baru:** kalau user minta bantu live coding → ikuti format di `LIVE_CODING_CHEATSHEET.md` (modifikasi kecil, anchor text exact, bukan penjelasan konseptual panjang).

---

## 4. Arsitektur project (ringkas)

```
Browser
  → Gateway :2000 (JWT cookie, role gate, static frontend, proxy)
      → Subscription service :5004  ← FOKUS UJIAN
          ├── modules/packages/   (CRUD paket, katalog)
          ├── modules/orders/     (order, bukti bayar, invoice, cancel, renew, verify)
          └── modules/dashboard/  (ops stats)
      → Management service :5006  (analytics dashboard — soal R mungkin, bukan C/U/D utama)
      → Users, Tickets, Notification, Account (di luar scope ujian)
```

**Stack:** Go microservices (ZaFramework) + PostgreSQL + static HTML/JS frontend (deploy Vercel) + gateway Railway.

**Pola backend Subscription:**
```
HTTP Handler (controller.go)
  → Dispatcher.DispatchAndWait("job_name")   // async job, sync wait
  → Service (validasi + business rule)
  → Repository (SQL)
  → PostgreSQL schema: subscription.*
```

---

## 5. Modul inti — file yang sering disentuh saat ujian

### Package (C/U/D utama)

| Layer | Path |
|-------|------|
| Route | `subscription/app/routes/router.go` |
| Controller | `subscription/app/modules/packages/controller.go` |
| Service | `subscription/app/modules/packages/service.go` |
| Repository | `subscription/app/modules/packages/repository.go` |
| Types/DTO | `subscription/app/modules/packages/types.go` |
| Job register | `subscription/main.go` |
| UI CRUD ops | `frontend/ops/subscriptions.html`, `subscriptions-create.html`, `subscriptions-edit.html` |

**PENTING:** CRUD paket UI ada di **`frontend/ops/subscriptions*.html`**, BUKAN `management/dashboard-packages.html` (itu analytics read-only).

### Order & Billing

| Layer | Path |
|-------|------|
| Controller/Service/Repo | `subscription/app/modules/orders/*.go` |
| Invoice PDF | `subscription/app/modules/orders/invoice_pdf.go` |
| Client UI | `frontend/client/checkout.html`, `billing-history.html`, `order-detail.html` |
| Ops UI | `frontend/ops/orders.html`, `orders-detail.html` |

### Subscription status

| Path |
|------|
| `frontend/client/subscription-me.html` |
| `GET /api/subscriptions/me` di orders controller |

### Katalog client (soal R)

| Path |
|------|
| `frontend/client/packages-catalog.html` |
| `GET /api/subscription/catalog` |

---

## 6. Path canonical frontend (jangan pakai path UAT lama)

| Fungsi | URL benar |
|--------|-----------|
| Riwayat order client | `/client/billing-history` |
| Status subscription | `/client/subscription-me` |
| Katalog paket | `/client/packages-catalog` |
| Checkout | `/client/checkout?package_id=` |
| CRUD paket ops | `/ops/subscriptions` |
| List order ops | `/ops/orders` |

Gateway sudah redirect `/orders/history` → `/client/billing-history`.

---

## 7. Database tables relevan

Schema: `subscription`

| Tabel | Fungsi |
|-------|--------|
| `packages` | id, name, price, quota, status (ACTIVE/INACTIVE/DELETED) |
| `package_pricing` | duration_months, price, label per package |
| `orders` | invoice, user_id, package_id, duration_months, payment_proof (BYTEA), status |
| `subscriptions` | order_id, user_id, start_date, end_date, status |
| `invoice_download_logs` | audit download PDF |

Migration: `subscription/core/database/migrate.go`

---

## 8. Fitur yang SUDAH implemented (baseline — jangan buat ulang)

### Package
- Create/Update/Delete/Toggle status paket + pricing tiers
- Validasi: nama unik, harga/kuota > 0, min 1 tier
- Soft delete: hanya INACTIVE, cek subscriber PAID/PENDING
- Admin list dengan query param `status`, `min_price`, `max_price` (backend ada; UI ops filter masih client-side — kandidat soal)
- Katalog ACTIVE only untuk client

### Order flow
- Checkout → `POST /api/orders` → PENDING_PAYMENT (cek KYC, paket ACTIVE)
- Upload payment proof (max 5MB, JPG/PNG/WebP)
- Ops verify APPROVE/REJECT → PAID + aktivasi subscription
- Client cancel PENDING (UI blok jika ada bukti; backend belum enforce — kandidat soal §7)
- Renew order `POST /api/orders/renew`
- Invoice PDF client & admin (PAID only) + audit log

### Subscription
- `GET /api/subscriptions/me` — active_subscriptions, latest_subscription, can_renew, days_remaining (backend hitung, UI belum tampilkan semua — kandidat soal §6)

### UAT yang sudah dikerjakan di sesi sebelumnya (PBI-63, 64, 66, 67)
- Billing history: filter auto-fetch, skeleton, cancel real-time
- Invoice PDF dark theme, Content-Disposition filename
- Renew flow di dashboard/subscription-me
- Cancel modal di order-detail
- File baru: `subscription-me.html`, `invoice_pdf.go`

---

## 9. Gap / kandidat soal (modifikasi kecil — ada di cheat sheet)

Bukan "bug wajib diperbaiki" — ini **prediksi** apa yang mungkin diminta asdos:

| § | Modifikasi | File utama |
|---|------------|------------|
| §1 | Subscriber count di modal delete | repo, types, subscriptions.html |
| §2 | Filter datepicker satu tanggal billing history | billing-history.html |
| §3 | Validasi nama tanpa karakter spesial | packages/service.go |
| §4 | Sort katalog termurah→termahal | packages-catalog.html |
| §5 | Kolom Quota di tabel ops | subscriptions.html |
| §6 | Tampilkan days_remaining | subscription-me.html |
| §7 | Cancel blok jika ada bukti | orders/service.go |
| §8 | Tampilkan duration_months di detail order | order-detail.html |
| §9 | Card total paket ACTIVE | subscriptions.html |
| §10 | Reject reason min 10 karakter | orders/service.go, orders-detail.html |

**5 soal paling mungkin:** §2 (R filter tanggal) + §1 atau §3 (C/U/D).

---

## 10. Dead code traps (controller expect error, service belum return)

Jika asdos minta "lengkapi validasi" — cek ini dulu di `packages/controller.go`:

- `"harga tahunan tidak boleh lebih rendah dari harga bulanan"` → belum di service (§13)
- `"tidak dapat mengubah paket yang sedang aktif"` → belum di service (bisa diminta mirip §2 pola validasi)

---

## 11. Urutan pengerjaan saat live coding

### Soal C/U/D
```
service.go → repository.go (jika perlu SQL) → controller.go (jika perlu) → frontend HTML → go build ./...
```

### Soal R
```
frontend HTML (fetch/render) → trace ke controller → service → repository (jika perlu field baru di API)
```

### Build check
```bash
cd subscription && go build ./...
cd gateway && go build .
```

### Test lokal
- Gateway: `localhost:2000`
- Login role: CLIENT (billing/katalog), OPERASIONAL (ops/subscriptions, ops/orders)

---

## 12. Role & route penting

| Role | Halaman / API |
|------|---------------|
| CLIENT | `/client/*`, `GET/POST /api/orders`, `/api/subscriptions/me` |
| OPERASIONAL | `/ops/subscriptions*`, `/api/admin/packages`, `/api/admin/orders`, verify |
| MANAGEMENT | `/management/dashboard-packages` → `GET /api/dashboard/packages` (service terpisah) |

Header dari gateway: `X-User-ID`, `X-User-Role`

---

## 13. Yang TIDAK perlu dikerjakan agent kecuali diminta explicit

- Commit/push git
- Deploy Vercel/Railway
- Refactor besar di luar scope soal
- Menulis dokumentasi panjang baru (sudah ada README + cheat sheet)
- Fix semua gap sekaligus — hanya implement skenario yang diminta user/asdos

---

## 14. Instruksi untuk agent baru

Ketika user minta bantuan live coding:

1. **Tanya atau identifikasi** kata kunci soal → cocokkan ke § di `LIVE_CODING_CHEATSHEET.md`
2. **Baca file aktual** di repo — jangan pakai line number dari memory, pakai **anchor text Ctrl+F**
3. **Tulis perubahan minimal** — 1 skenario = modifikasi kecil, bukan rewrite modul
4. **Format output:** File dibuka → anchor → TAMBAHKAN/UBAH → build check → verifikasi browser 1 kalimat
5. **Jangan** tulis penjelasan arsitektur panjang kecuali user minta

Ketika user minta prediksi soal baru:
- Ikuti pola kelompok lain (modifikasi kecil di fitur existing)
- Base pada file yang benar-benar ada di repo christopher
- Estimasi waktu 3–8 menit per soal

---

## 15. Status repo (perkiraan)

- Branch aktif: `PBI-Dashboard` (banyak commit ahead of remote)
- Ada perubahan UAT belum tentu semua di-commit
- File penting yang pernah untracked: `frontend/client/subscription-me.html`, `subscription/app/modules/orders/invoice_pdf.go`
- Remote: GitLab (bukan GitHub)

---

## 16. Prompt singkat untuk paste ke Claude

```
Konteks: Saya persiapan live coding Propen A+ modul Paket di repo ThinkNalyze (christopher).
Format ujian: 30 menit, 1 soal C/U/D + 1 soal R.
Pola soal: modifikasi kecil di fitur yang sudah jalan (filter, kolom, validasi, badge) — BUKAN buat dari nol.
Dokumen utama: LIVE_CODING_CHEATSHEET.md (§1-§45 + demo §R). Coverage lengkap: katalog, checkout, billing, subscription-me, dashboard card, ops.
Backend fokus: subscription service (packages + orders modules).
Frontend CRUD paket: frontend/ops/subscriptions*.html (bukan management dashboard).
Path client: /client/billing-history, /client/packages-catalog, /client/subscription-me.
Baca docs/HANDOFF_CONTEXT_LIVE_CODING.md dan LIVE_CODING_CHEATSHEET.md sebelum membantu.
Saat implementasi: anchor text Ctrl+F, kode minimal, go build, verifikasi browser.
```

---

*Terakhir diupdate: Juni 2026 — sesi persiapan live coding modul Paket.*
