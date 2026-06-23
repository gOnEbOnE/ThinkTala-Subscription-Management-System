# LIVE CODING CHEAT SHEET v2 — MODIFIKASI KECIL BERTARGET

> **Format ujian:** 30 menit = 1 soal C/U/D + 1 soal R.
> **Pola soal (sama kelompok lain):** BUKAN buat dari nol, BUKAN fix bug — **tambah/ubah sesuatu kecil di fitur yang sudah jalan**.
> **Cara pakai:** Ctrl+F kata kunci di **Tabel Navigasi** → lompat ke § → Ctrl+F anchor text → paste kode → verifikasi.

---

# 15 KEMUNGKINAN SOAL (URUT PROBABILITAS)

| # | Soal (kalimat asdos) | Tipe | Fitur | Modifikasi | File | Waktu |
|---|----------------------|------|-------|------------|------|-------|
| 1 | "Pada modal delete package, tampilkan berapa jumlah subscriber aktif sebelum user konfirmasi hapus" | C/U/D | Delete Package | Tambah `subscriber_count` di response list + tampil di modal | repo, types, subscriptions.html | 8 mnt |
| 2 | "Pada billing history, tambahkan filter datepicker satu tanggal saja untuk filter order pada tanggal tersebut" | R | Billing History | Tambah input date + kirim `start_date` & `end_date` sama | billing-history.html | 5 mnt |
| 3 | "Saat create/update package, validasi nama tidak boleh mengandung karakter spesial (hanya huruf, angka, spasi)" | C/U/D | Create/Update Package | Tambah regex di service | service.go | 5 mnt |
| 4 | "Pada halaman katalog client, urutkan paket dari harga termurah ke termahal" | R | Katalog | Sort array setelah fetch | packages-catalog.html | 4 mnt |
| 5 | "Pada tabel list paket operasional, tambahkan kolom Quota" | R | List Package Admin | Render field `pkg.quota` yang sudah ada di API | subscriptions.html | 4 mnt |
| 6 | "Pada halaman subscription saya, tampilkan sisa hari masa aktif berlangganan" | R | Subscription Status | Tampilkan `days_remaining` dari API (sudah dihitung backend) | subscription-me.html | 5 mnt |
| 7 | "Client tidak boleh cancel order yang sudah upload bukti transfer" | C/U/D | Cancel Order | Tambah cek `HasPaymentProof` di service | orders/service.go | 4 mnt |
| 8 | "Pada detail order client, tampilkan durasi langganan (berapa bulan)" | R | Order Detail | Tampilkan `duration_months` dari response | order-detail.html | 4 mnt |
| 9 | "Pada halaman list paket ops, tambahkan card yang menampilkan total paket berstatus ACTIVE" | R | List Package Admin | Hitung dari `allPackages` setelah load | subscriptions.html | 5 mnt |
| 10 | "Saat operasional reject pembayaran, alasan wajib minimal 10 karakter" | C/U/D | Verify Order | Validasi panjang di service + pesan di frontend | orders/service.go, orders-detail.html | 6 mnt |
| 11 | "Pada katalog, beri badge 'Harga Terbaik' pada paket dengan harga bulanan termurah" | R | Katalog | Logic badge di `buildPackageCard` | packages-catalog.html | 5 mnt |
| 12 | "Pada billing history, tambahkan kolom Durasi (bulan) di tabel" | R | Billing History | Tambah field di list API + kolom tabel | types, repository, billing-history.html | 8 mnt |
| 13 | "Saat create package, validasi harga tier 12 bulan tidak boleh lebih rendah dari tier 1 bulan" | C/U/D | Create Package | Tambah validasi di service | service.go | 5 mnt |
| 14 | "Pada list order operasional, tampilkan tanggal order dalam format Indonesia (contoh: 8 Juni 2026)" | R | Admin Order List | Format `created_at` di render | orders.html | 4 mnt |
| 15 | "Saat deactivate package, modal konfirmasi harus menyebut nama paket dan status tujuan" | C/U/D | Toggle Status | Perbaiki teks di `toggleStatus()` | subscriptions.html | 3 mnt |

---

# TABEL NAVIGASI CEPAT

| Kata kunci asdos | § |
|------------------|---|
| subscriber / pelanggan aktif / delete modal / hapus paket | **§1** |
| filter tanggal / datepicker / satu tanggal / billing history filter | **§2** |
| karakter spesial / nama paket / validasi nama | **§3** |
| urutkan / sort / termurah / termahal / katalog | **§4** |
| kolom quota / kuota / tabel paket | **§5** |
| sisa hari / days remaining / subscription aktif | **§6** |
| cancel / batalkan / bukti transfer / tidak boleh cancel | **§7** |
| durasi / bulan / duration / detail order | **§8** |
| card total / jumlah paket active / summary card | **§9** |
| reject / alasan / minimal karakter / verifikasi tolak | **§10** |
| badge / harga terbaik / termurah / label katalog | **§11** |
| kolom durasi / billing history tabel | **§12** |
| tier 12 bulan / harga tahunan / validasi harga | **§13** |
| format tanggal / Indonesia / list order ops | **§14** |
| deactivate / toggle / modal konfirmasi status | **§15** |
| deactivate blok subscriber / masih ada pelanggan | **§16** |
| upload bukti ulang / sudah ada bukti | **§17** |
| email valid / format email checkout | **§18** |
| harga tier harus > 0 / tidak boleh nol | **§19** |
| trim nama paket / spasi nama | **§20** |
| durasi harus ada tier / tidak fallback | **§21** |
| maksimal 5 tier / terlalu banyak tier | **§22** |
| nama client minimal / client_name | **§23** |
| kuota maksimal / quota max | **§24** |
| alasan reject maksimal 500 | **§25** |
| label diskon / hemat di katalog | **§26** |
| badge bukti / ada bukti billing | **§27** |
| nama pemesan / client_name detail | **§28** |
| alasan penolakan / verification_note | **§29** |
| ikon bukti / proof icon ops orders | **§30** |
| jumlah tier / berapa tier paket | **§31** |
| tanggal dibuat paket / created_at ops | **§32** |
| harga 12 bulan / annual katalog | **§33** |
| urutkan nama A-Z / sort paket ops | **§34** |
| toast invoice / unduh sukses | **§35** |
| sisa hari / dashboard subscription card | **§36** |
| badge EXPIRED / abu / subscription status | **§37** |
| kolom metode pembayaran / payment method billing | **§38** |
| label tier / ringkasan checkout | **§39** |
| sorot baris PAID / hijau billing | **§40** |
| empty state / refresh / muat ulang katalog | **§41** |
| menunggu bukti / ops order list | **§42** |
| client_name / nama pemesan API detail | **§43** |
| toast renew / dashboard submit renewal | **§44** |
| waktu verifikasi / verified_at ops detail | **§45** |
| invoice / PDF / unduh | **§R-INV** |
| upload bukti / payment proof | **§R-PROOF** |
| verify / approve pembayaran | **§R-VERIFY** |

---

# §1. DELETE MODAL — TAMPILKAN JUMLAH SUBSCRIBER

> Tipe: C/U/D | Di fitur: Delete Package | Estimasi: **8 menit**

### Apa yang diubah:
Di response list paket admin, tambahkan field `subscriber_count`, lalu tampilkan di modal delete sebelum tombol konfirmasi.

### File yang dibuka:
1. `subscription/app/modules/packages/types.go` — anchor: `type Package struct`
2. `subscription/app/modules/packages/repository.go` — anchor: `FROM subscription.packages WHERE 1=1`
3. `frontend/ops/subscriptions.html` — anchor: `function openDelete`

---

### Edit 1: `subscription/app/modules/packages/types.go`

**Ctrl+F:**
```go
type Package struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
```

**TAMBAHKAN SETELAH** field `Status`:
```go
	SubscriberCount int           `json:"subscriber_count"`
```

**Kenapa:** Field baru untuk dikirim ke frontend tanpa endpoint terpisah.

---

### Edit 2: `subscription/app/modules/packages/repository.go`

**Ctrl+F:**
```go
	query := `SELECT id, name, price, quota, status, created_at, updated_at FROM subscription.packages WHERE 1=1`
```

**UBAH MENJADI:**
```go
	query := `SELECT p.id, p.name, p.price, p.quota, p.status, p.created_at, p.updated_at,
		COALESCE((
			SELECT COUNT(*)::int FROM subscription.orders o
			WHERE o.package_id = p.id AND o.status IN ('PAID','PENDING_PAYMENT')
		), 0) AS subscriber_count
		FROM subscription.packages p WHERE 1=1`
```

**Ctrl+F:**
```go
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
```

**UBAH MENJADI:**
```go
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quota, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.SubscriberCount); err != nil {
```

**Kenapa:** Hitung subscriber saat list, sama seperti logic `CountActiveSubscribers`.

---

### Edit 3: `frontend/ops/subscriptions.html` — modal HTML

**Ctrl+F:**
```html
                    <p class="text-muted small mb-0" id="deletePkgName">This action cannot be undone.</p>
                    <input type="hidden" id="deletePkgId">
```

**TAMBAHKAN SETELAH** baris `deletePkgName`:
```html
                    <p class="text-muted small mt-2 mb-0" id="deletePkgSubscribers"></p>
```

---

### Edit 4: `frontend/ops/subscriptions.html` — function openDelete

**Ctrl+F:**
```javascript
        function openDelete(id, name) {
            document.getElementById('deletePkgId').value = id;
            document.getElementById('deletePkgName').textContent = `Delete "${name}"? This action cannot be undone.`;
            new bootstrap.Modal(document.getElementById('deletePkgModal')).show();
        }
```

**UBAH MENJADI:**
```javascript
        function openDelete(id, name) {
            document.getElementById('deletePkgId').value = id;
            document.getElementById('deletePkgName').textContent = `Delete "${name}"? This action cannot be undone.`;
            const pkg = (allPackages || []).find(p => p.id === id);
            const count = pkg ? (pkg.subscriber_count || 0) : 0;
            document.getElementById('deletePkgSubscribers').textContent =
                count > 0
                    ? `⚠ ${count} pelanggan aktif masih menggunakan paket ini.`
                    : 'Tidak ada pelanggan aktif pada paket ini.';
            new bootstrap.Modal(document.getElementById('deletePkgModal')).show();
        }
```

**Kenapa:** Tampilkan info sebelum user klik Delete.

---

### Build check:
```bash
cd subscription && go build ./...
```

### Verifikasi (30 detik):
Browser: `http://localhost:2000/ops/subscriptions` → klik trash pada paket INACTIVE.
**Expected:** Modal menampilkan teks jumlah pelanggan aktif.

### Jika error:
| Error | Penyebab | Fix |
|-------|----------|-----|
| `go build` scan mismatch | Lupa ubah `rows.Scan` | Tambah `&p.SubscriberCount` di Scan |
| Modal tidak tampil angka | `allPackages` belum load | Refresh halaman dulu |
| Selalu 0 | Belum ada order PAID/PENDING | Normal jika DB kosong |

---

# §2. BILLING HISTORY — FILTER SATU TANGGAL

> Tipe: R | Di fitur: Billing History | Estimasi: **5 menit**

### Apa yang diubah:
Di filter bar billing history, tambahkan datepicker satu tanggal; saat dipilih, kirim `start_date` dan `end_date` dengan nilai sama ke API yang sudah ada.

### File yang dibuka:
1. `frontend/client/billing-history.html` — anchor: `filterStartDate`

---

### Edit 1: HTML filter bar

**Ctrl+F:**
```html
                <div>
                    <label class="form-label small fw-semibold mb-1" style="color:var(--text-muted);font-size:0.72rem;text-transform:uppercase;">From Date</label>
                    <input type="date" class="form-control form-control-sm" id="filterStartDate" style="min-width:140px;">
                </div>
```

**TAMBAHKAN SEBELUM** blok From Date:
```html
                <div>
                    <label class="form-label small fw-semibold mb-1" style="color:var(--text-muted);font-size:0.72rem;text-transform:uppercase;">On Date</label>
                    <input type="date" class="form-control form-control-sm" id="filterOnDate" style="min-width:140px;">
                </div>
```

---

### Edit 2: function `loadHistory`

**Ctrl+F:**
```javascript
            var status = document.getElementById('filterStatus').value;
            var startDate = document.getElementById('filterStartDate').value;
            var endDate = document.getElementById('filterEndDate').value;
            if (status) params.set('status', status);
            if (startDate) params.set('start_date', startDate);
            if (endDate) params.set('end_date', endDate);
```

**TAMBAHKAN SETELAH** baris `var endDate = ...`:
```javascript
            var onDate = document.getElementById('filterOnDate').value;
            if (onDate) {
                startDate = onDate;
                endDate = onDate;
            }
```

---

### Edit 3: function `resetFilter` — tambah reset field baru

**Ctrl+F:**
```javascript
            document.getElementById('filterStartDate').value = '';
            document.getElementById('filterEndDate').value = '';
```

**TAMBAHKAN SETELAH:**
```javascript
            document.getElementById('filterOnDate').value = '';
```

---

### Edit 4: event listener auto-fetch

**Ctrl+F:**
```javascript
        ['filterStatus', 'filterStartDate', 'filterEndDate'].forEach(function (id) {
```

**UBAH MENJADI:**
```javascript
        ['filterStatus', 'filterStartDate', 'filterEndDate', 'filterOnDate'].forEach(function (id) {
```

---

### Verifikasi (30 detik):
Browser: `http://localhost:2000/client/billing-history`
**Lakukan:** Pilih tanggal di "On Date".
**Expected:** Network `GET /api/orders?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` (sama), tabel hanya order hari itu.

---

# §3. VALIDASI NAMA PAKET — TANPA KARAKTER SPESIAL

> Tipe: C/U/D | Di fitur: Create/Update Package | Estimasi: **5 menit**

### Apa yang diubah:
Di service, tambahkan validasi regex nama hanya huruf, angka, spasi.

### File yang dibuka:
1. `subscription/app/modules/packages/service.go`

---

### Edit 1: tambah import

**Ctrl+F:**
```go
import (
	"context"
	"errors"
	"fmt"
)
```

**UBAH MENJADI:**
```go
import (
	"context"
	"errors"
	"fmt"
	"regexp"
)
```

---

### Edit 2: function `CreatePackage` — setelah cek nama kosong

**Ctrl+F:**
```go
	if payload.Name == "" {
		return nil, errors.New("nama paket tidak boleh kosong")
	}
	if payload.Price <= 0 {
```

**TAMBAHKAN SETELAH** cek nama kosong:
```go
	var validName = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	if !validName.MatchString(payload.Name) {
		return nil, errors.New("nama paket tidak boleh mengandung karakter spesial")
	}
```

---

### Edit 3: function `UpdatePackage` — copy blok `validName` yang sama setelah cek nama kosong (~baris 88).

---

### Build check:
```bash
cd subscription && go build ./...
```

### Verifikasi (30 detik):
Browser: `/ops/subscriptions-create` → nama `"Paket@Pro!"` → submit.
**Expected:** error `"nama paket tidak boleh mengandung karakter spesial"`.

---

# §4. KATALOG — URUTKAN HARGA TERMURAH KE TERMAHAL

> Tipe: R | Di fitur: Katalog Client | Estimasi: **4 menit**

### Apa yang diubah:
Setelah fetch katalog, sort array `data.data` berdasarkan `price` ascending sebelum render.

### File yang dibuka:
1. `frontend/client/packages-catalog.html` — anchor: `data.data.map`

---

### Edit: function `fetchCatalog`

**Ctrl+F:**
```javascript
                if (!data.success || !Array.isArray(data.data) || data.data.length === 0) {
                    emptyEl.classList.remove('d-none');
                    return;
                }

                // Render semua kartu paket
                gridEl.innerHTML = data.data.map(function (pkg, i) {
```

**TAMBAHKAN SETELAH** cek array kosong:
```javascript
                var sorted = data.data.slice().sort(function (a, b) {
                    return (Number(a.price) || 0) - (Number(b.price) || 0);
                });

```

**UBAH:**
```javascript
                gridEl.innerHTML = data.data.map(function (pkg, i) {
```
**MENJADI:**
```javascript
                gridEl.innerHTML = sorted.map(function (pkg, i) {
```

---

### Verifikasi (30 detik):
Browser: `http://localhost:2000/client/packages-catalog`
**Expected:** Paket harga terendah tampil paling kiri/atas.

---

# §5. LIST PAKET OPS — TAMBAH KOLOM QUOTA

> Tipe: R | Di fitur: List Package Admin | Estimasi: **4 menit**

### Apa yang diubah:
Di tabel subscriptions ops, tambahkan kolom Quota — data `pkg.quota` sudah ada di API.

### File yang dibuka:
1. `frontend/ops/subscriptions.html`

---

### Edit 1: thead

**Ctrl+F:**
```html
                            <th class="small fw-semibold py-3" style="text-transform:uppercase;font-size:0.72rem;">Annual
                                Price</th>
                            <th class="small fw-semibold py-3"
```

**TAMBAHKAN SETELAH** kolom Annual Price:
```html
                            <th class="small fw-semibold py-3" style="text-transform:uppercase;font-size:0.72rem;">Quota</th>
```

---

### Edit 2: tbody render di `renderTable`

**Ctrl+F:**
```javascript
                    <td class="py-3 pkg-price-yearly">${formatRp(getDurationPrice(pkg, 12))}</td>
                    <td class="py-3">${statusBadge(pkg.status)}</td>
```

**TAMBAHKAN SETELAH** annual price:
```javascript
                    <td class="py-3">${(pkg.quota || 0).toLocaleString('id-ID')}</td>
```

---

### Edit 3: colspan loading/empty — cari `colspan="5"` ubah jadi `colspan="6"`.

---

### Verifikasi (30 detik):
Browser: `http://localhost:2000/ops/subscriptions`
**Expected:** Kolom Quota muncul dengan angka (contoh: 1.000).

---

# §6. SUBSCRIPTION ME — TAMPILKAN SISA HARI

> Tipe: R | Di fitur: Subscription Status | Estimasi: **5 menit**

### Apa yang diubah:
Backend sudah hitung `days_remaining` di `GetActiveSubscriptions` — tampilkan di UI subscription-me.

### File yang dibuka:
1. `frontend/client/subscription-me.html` — anchor: `renderSubscriptionOverview`

---

### Edit 1: HTML — cari elemen end date, tambah setelahnya

**Ctrl+F di HTML** elemen `subEndDate` — tambahkan setelah parent info-item:
```html
                        <div class="info-item">
                            <div class="info-label">Sisa Hari Aktif</div>
                            <div class="info-value" id="subDaysRemaining">-</div>
                        </div>
```

*(Sesuaikan struktur HTML di sekitar `id="subEndDate"` — paste di dalam grid info yang sama.)*

---

### Edit 2: function `renderSubscriptionOverview`

**Ctrl+F:**
```javascript
            endEl.textContent = formatIdDate(display.end_date);
            defaultPackageId = display.package_id || '';
```

**TAMBAHKAN SETELAH:**
```javascript
            var daysEl = document.getElementById('subDaysRemaining');
            if (daysEl) {
                var days = typeof display.days_remaining === 'number' ? display.days_remaining : 0;
                daysEl.textContent = status === 'ACTIVE' ? (days + ' hari') : '-';
            }
```

---

### Verifikasi (30 detik):
Browser: `http://localhost:2000/client/subscription-me` (login CLIENT dengan sub ACTIVE)
**Expected:** "Sisa Hari Aktif" menampilkan angka hari.

---

# §7. CANCEL ORDER — BLOK JIKA SUDAH ADA BUKTI

> Tipe: C/U/D | Di fitur: Cancel Order | Estimasi: **4 menit**

### Apa yang diubah:
Di service `CancelOrder`, tambahkan cek `HasPaymentProof` — UI sudah sembunyikan tombol, backend belum enforce.

### File yang dibuka:
1. `subscription/app/modules/orders/service.go`

---

### Edit: function `CancelOrder`

**Ctrl+F:**
```go
	if rec.Status != "PENDING_PAYMENT" {
		return errors.New("Pesanan dengan status ini tidak dapat dibatalkan.")
	}

	return s.repo.CancelOrder(ctx, orderID, userID)
```

**TAMBAHKAN SEBELUM** `return s.repo.CancelOrder`:
```go
	if rec.HasPaymentProof {
		return errors.New("Pesanan yang sudah memiliki bukti transfer tidak dapat dibatalkan.")
	}

```

---

### Build check:
```bash
cd subscription && go build ./...
```

### Verifikasi (30 detik):
Curl order PENDING yang sudah punya bukti:
**Expected:** error message di atas, status tetap PENDING_PAYMENT.

---

# §8. DETAIL ORDER — TAMPILKAN DURASI (BULAN)

> Tipe: R | Di fitur: Order Detail Client | Estimasi: **4 menit**

### Apa yang diubah:
Field `duration_months` sudah ada di response `GET /api/orders/{id}` — tambahkan baris di UI.

### File yang dibuka:
1. `frontend/client/order-detail.html`

---

### Edit 1: HTML info-items

**Ctrl+F:**
```html
                        <div class="info-item">
                            <div class="info-label">Metode Pembayaran</div>
                            <div class="info-value" id="paymentMethod">-</div>
                        </div>
```

**TAMBAHKAN SEBELUM** blok Metode Pembayaran:
```html
                        <div class="info-item">
                            <div class="info-label">Durasi Langganan</div>
                            <div class="info-value" id="durationMonths">-</div>
                        </div>
```

---

### Edit 2: function `loadOrderDetail`

**Ctrl+F:**
```javascript
                document.getElementById('paymentMethod').textContent = d.payment_method || '-';
                document.getElementById('totalPrice').textContent = formatRp(d.total_price);
```

**TAMBAHKAN SETELAH:**
```javascript
                document.getElementById('durationMonths').textContent =
                    (d.duration_months ? d.duration_months + ' bulan' : '-');
```

---

### Verifikasi (30 detik):
Browser: `http://localhost:2000/client/order-detail?id=<order_id>`
**Expected:** "Durasi Langganan" menampilkan misalnya "3 bulan".

---

# §9. LIST PAKET OPS — CARD TOTAL PAKET ACTIVE

> Tipe: R | Di fitur: List Package Admin | Estimasi: **5 menit**

### Apa yang diubah:
Tambahkan card ringkasan di atas tabel yang menghitung paket ACTIVE dari data yang sudah di-fetch.

### File yang dibuka:
1. `frontend/ops/subscriptions.html`

---

### Edit 1: HTML — setelah Filter Section, sebelum Table Section

**Ctrl+F:**
```html
        <!-- Table Section -->
        <div class="ai-card">
            <div class="mb-3">
                <h6 class="fw-bold mb-0" style="color: var(--text-heading)">Package List</h6>
```

**TAMBAHKAN SEBELUM** `<!-- Table Section -->`:
```html
        <div class="row g-3 mb-3">
            <div class="col-md-4">
                <div class="ai-card py-3 px-4">
                    <div class="small text-muted text-uppercase fw-semibold">Paket Active</div>
                    <div class="fs-4 fw-bold text-main" id="activePackageCount">0</div>
                </div>
            </div>
        </div>
```

---

### Edit 2: function `loadPackages` — setelah `allPackages = json.data`

**Ctrl+F:**
```javascript
                allPackages = json.data || [];
                currentPage = 1;
                renderTable(allPackages);
```

**TAMBAHKAN SETELAH** `allPackages = ...`:
```javascript
                const activeCount = allPackages.filter(p => p.status === 'ACTIVE').length;
                const countEl = document.getElementById('activePackageCount');
                if (countEl) countEl.textContent = activeCount;
```

---

### Verifikasi (30 detik):
Browser: `http://localhost:2000/ops/subscriptions`
**Expected:** Card "Paket Active" menampilkan angka sesuai jumlah ACTIVE di DB.

---

# §10. REJECT VERIFY — ALASAN MINIMAL 10 KARAKTER

> Tipe: C/U/D | Di fitur: Verify Order | Estimasi: **6 menit**

### Apa yang diubah:
Tambahkan validasi panjang `reject_reason` minimal 10 karakter di service dan frontend.

### File yang dibuka:
1. `subscription/app/modules/orders/service.go`
2. `frontend/ops/orders-detail.html`

---

### Edit 1: `orders/service.go` — function `VerifyOrder`

**Ctrl+F:**
```go
	if action == "REJECT" && strings.TrimSpace(rejectReason) == "" {
		return nil, errors.New("alasan reject wajib diisi")
	}
```

**UBAH MENJADI:**
```go
	if action == "REJECT" {
		rejectReason = strings.TrimSpace(rejectReason)
		if rejectReason == "" {
			return nil, errors.New("alasan reject wajib diisi")
		}
		if len(rejectReason) < 10 {
			return nil, errors.New("alasan reject minimal 10 karakter")
		}
	}
```

---

### Edit 2: `orders-detail.html` — sebelum fetch verify

**Ctrl+F:**
```javascript
            if (verifyMode === 'REJECT' && !rejectReason) {
                showToast('Alasan reject wajib diisi.', 'danger');
                rejectReasonInput.focus();
```

**TAMBAHKAN SETELAH** cek kosong:
```javascript
            if (verifyMode === 'REJECT' && rejectReason.length < 10) {
                showToast('Alasan reject minimal 10 karakter.', 'danger');
                rejectReasonInput.focus();
                return;
            }
```

---

### Build check:
```bash
cd subscription && go build ./...
```

### Verifikasi (30 detik):
Browser: `/ops/orders-detail?id=<pending_with_proof>` → Reject dengan alasan "salah" (5 huruf).
**Expected:** toast minimal 10 karakter.

---

# §11. KATALOG — BADGE "HARGA TERBAIK"

> Tipe: R | Di fitur: Katalog Client | Estimasi: **5 menit**

### Apa yang diubah:
Setelah sort (atau sebelum render), tentukan paket termurah, tampilkan badge di kartunya.

### File yang dibuka:
1. `frontend/client/packages-catalog.html`

---

### Edit: function `fetchCatalog` — sebelum map render

**Ctrl+F:**
```javascript
                gridEl.innerHTML = sorted.map(function (pkg, i) {
                    return buildPackageCard(pkg, i);
                }).join('');
```

**UBAH MENJADI:**
```javascript
                var minPrice = Math.min.apply(null, sorted.map(function (p) { return Number(p.price) || 0; }));
                gridEl.innerHTML = sorted.map(function (pkg, i) {
                    var isBest = (Number(pkg.price) || 0) === minPrice;
                    return buildPackageCard(pkg, i, isBest);
                }).join('');
```

**Ctrl+F signature:**
```javascript
        function buildPackageCard(pkg, index) {
```

**UBAH MENJADI:**
```javascript
        function buildPackageCard(pkg, index, isBestPrice) {
```

**Ctrl+F di dalam buildPackageCard** setelah `<h5`:
```javascript
                '      <h5 class="fw-bold text-white mb-0">' + _escapeHtml(pkg.name) + '</h5>',
```

**TAMBAHKAN SETELAH:**
```javascript
                (isBestPrice ? '      <span class="badge bg-warning text-dark mt-2">Harga Terbaik</span>' : ''),
```

---

### Verifikasi (30 detik):
Browser: `/client/packages-catalog`
**Expected:** Satu kartu paket termurah punya badge "Harga Terbaik".

---

# §12. BILLING HISTORY — KOLOM DURASI ⚠️ KOMPLEKS

> Tipe: R | Di fitur: Billing History | Estimasi: **8 menit** | **Kerjakan:** types → repository → frontend

### Apa yang diubah:
Tambahkan `duration_months` ke list API client, lalu kolom baru di tabel.

### File yang dibuka:
1. `subscription/app/modules/orders/types.go`
2. `subscription/app/modules/orders/repository.go`
3. `frontend/client/billing-history.html`

---

### Edit 1: `types.go` — `ClientOrderListItem`

**Ctrl+F:**
```go
	PackageName            string     `json:"package_name"`
	TotalPrice             float64    `json:"total_price"`
```

**TAMBAHKAN:**
```go
	DurationMonths         int        `json:"duration_months"`
```

---

### Edit 2: `repository.go` — SELECT list client

**Ctrl+F:**
```go
			COALESCE(p.name, '-') AS package_name,
			o.total_price,
```

**UBAH MENJADI:**
```go
			COALESCE(p.name, '-') AS package_name,
			o.duration_months,
			o.total_price,
```

**Ctrl+F Scan list:**
```go
			&item.PackageName,
			&item.TotalPrice,
```

**UBAH MENJADI:**
```go
			&item.PackageName,
			&item.DurationMonths,
			&item.TotalPrice,
```

---

### Edit 3: `billing-history.html` — thead tambah kolom "Duration" setelah Package

**Ctrl+F render row:**
```javascript
                    '<td class="py-3">' + (o.package_name || '-') + '</td>' +
                    '<td class="py-3">' + formatRp(o.total_price) + '</td>' +
```

**UBAH MENJADI:**
```javascript
                    '<td class="py-3">' + (o.package_name || '-') + '</td>' +
                    '<td class="py-3">' + (o.duration_months ? o.duration_months + ' bln' : '-') + '</td>' +
                    '<td class="py-3">' + formatRp(o.total_price) + '</td>' +
```

Update `colspan` skeleton/empty dari 6 ke 7.

---

### Build check:
```bash
cd subscription && go build ./...
```

### Verifikasi (30 detik):
Browser: `/client/billing-history`
**Expected:** Kolom Duration menampilkan "1 bln", "3 bln", dll.

---

# §13. CREATE PACKAGE — VALIDASI HARGA TIER 12 ≥ TIER 1

> Tipe: C/U/D | Di fitur: Create Package | Estimasi: **5 menit**

### Apa yang diubah:
Di `CreatePackage` (dan `UpdatePackage`), bandingkan harga tier 1 bulan vs 12 bulan.

### File yang dibuka:
1. `subscription/app/modules/packages/service.go`

---

### Edit: `CreatePackage` — setelah loop validasi tier

**Ctrl+F:**
```go
	for _, t := range payload.PricingTiers {
		if t.DurationMonths <= 0 {
			return nil, errors.New("durasi harga harus lebih besar dari 0 bulan")
		}
		if t.Price < 0 {
			return nil, errors.New("harga tidak boleh negatif")
		}
	}

	existing, err := s.repo.GetPackageByName(ctx, payload.Name)
```

**TAMBAHKAN SETELAH** loop `for _, t`:
```go
	var monthlyPrice float64
	var yearlyPrice float64
	hasMonthly := false
	hasYearly := false
	for _, t := range payload.PricingTiers {
		if t.DurationMonths == 1 {
			monthlyPrice = t.Price
			hasMonthly = true
		}
		if t.DurationMonths == 12 {
			yearlyPrice = t.Price
			hasYearly = true
		}
	}
	if hasMonthly && hasYearly && yearlyPrice < monthlyPrice {
		return nil, errors.New("harga tahunan tidak boleh lebih rendah dari harga bulanan")
	}

```

Copy ke `UpdatePackage` juga (sebelum `existing, err := s.repo.GetPackageByID`).

---

### Build check:
```bash
cd subscription && go build ./...
```

### Verifikasi (30 detik):
Create paket: tier 1 bln Rp 500.000, tier 12 bln Rp 400.000 → **Expected:** error validasi.

---

# §14. LIST ORDER OPS — FORMAT TANGGAL INDONESIA

> Tipe: R | Di fitur: Admin Order List | Estimasi: **4 menit**

### Apa yang diubah:
Ubah render tanggal `created_at` ke format Indonesia di tabel ops orders.

### File yang dibuka:
1. `frontend/ops/orders.html`

---

### Edit: tambah helper + ubah render

**Ctrl+F:** `function formatRp` atau awal `<script>` — **TAMBAHKAN:**
```javascript
        function formatDateID(iso) {
            if (!iso) return '-';
            const d = new Date(iso);
            if (isNaN(d.getTime())) return '-';
            return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' });
        }
```

**Ctrl+F di renderOrders/renderTable** tempat `created_at` ditampilkan — ganti ke `formatDateID(o.created_at)`.

*(Jika belum ada kolom tanggal di tabel, tambahkan `<th>Date</th>` dan `<td>${formatDateID(o.created_at)}</td>` — field `created_at` sudah ada di response admin list.)*

---

### Verifikasi (30 detik):
Browser: `http://localhost:2000/ops/orders`
**Expected:** Tanggal format "8 Juni 2026" bukan ISO mentah.

---

# §15. TOGGLE DEACTIVATE — MODAL SEBUT NAMA + STATUS TUJUAN

> Tipe: C/U/D | Di fitur: Toggle Status | Estimasi: **3 menit**

### Apa yang diubah:
Perjelas teks modal toggle agar menyebut nama paket dan status baru.

### File yang dibuka:
1. `frontend/ops/subscriptions.html`

---

### Edit: function `toggleStatus`

**Ctrl+F:**
```javascript
            document.getElementById('toggleDesc').textContent = `Are you sure you want to ${isDeactivate ? 'deactivate' : 'activate'} the package "${name}"?`;
```

**UBAH MENJADI:**
```javascript
            const nextStatus = currentStatus === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE';
            document.getElementById('toggleDesc').textContent =
                `Paket "${name}" akan diubah dari ${currentStatus} menjadi ${nextStatus}. Lanjutkan?`;
```

---

### Verifikasi (30 detik):
Browser: `/ops/subscriptions` → klik Deactivate.
**Expected:** Modal teks menyebut nama paket dan "ACTIVE menjadi INACTIVE".

---

# §R-INV. DOWNLOAD INVOICE (SUDAH JALAN — DEMO SAJA)

> Tipe: R | Estimasi demo: **2 menit**

### Apakah perlu ubah kode?
**TIDAK**

### Verifikasi (30 detik):
Browser: `/client/billing-history` → order PAID → Download Invoice.
**Expected:** File PDF terdownload.

### Trace singkat (kalau ditanya):
`billing-history.html` → `GET /api/orders/{id}/invoice` → `GetInvoiceHandler` → `invoice_pdf.go`

---

# §R-PROOF. UPLOAD & PREVIEW BUKTI (SUDAH JALAN)

> Tipe: R | Estimasi demo: **2 menit**

### Apakah perlu ubah kode?
**TIDAK**

### Verifikasi (30 detik):
Browser: `/client/order-detail?id=<pending>` → upload JPG < 5MB.
**Expected:** Preview gambar muncul, status AWAITING VERIFICATION.

---

# §R-VERIFY. VERIFY ORDER OPS (SUDAH JALAN)

> Tipe: R | Estimasi demo: **2 menit**

### Apakah perlu ubah kode?
**TIDAK** (kecuali soal minta validasi reject → **§10**)

### Verifikasi (30 detik):
Browser: `/ops/orders-detail?id=<pending_with_proof>` → Approve.
**Expected:** Status jadi PAID, subscription aktif.

---

# TABEL ERROR OPERASIONAL

| Yang terjadi | File | Ctrl+F anchor | Fix |
|--------------|------|---------------|-----|
| `go build` scan: wrong number of arguments | `repository.go` | `rows.Scan(&p.ID` | Tambah `&p.SubscriberCount` |
| `undefined: regexp` | `service.go` | `import (` | Tambah `"regexp"` |
| Filter tanggal tidak jalan | `billing-history.html` | `filterOnDate` | Pastikan `onDate` set start & end sama |
| Kolom quota tidak muncul | `subscriptions.html` | `pkg.quota` | Update colspan thead/tbody ke 6 |
| Cancel masih bisa via API | `orders/service.go` | `HasPaymentProof` | Tambah cek sebelum `CancelOrder` |
| Badge tidak muncul di katalog | `packages-catalog.html` | `buildPackageCard` | Pastikan param `isBestPrice` diteruskan |
| Duration kolom kosong | `repository.go` | `duration_months` | Tambah di SELECT + Scan |
| 401 di browser | gateway | — | Login ulang role yang benar |
| Modal delete subscriber selalu 0 | `subscriptions.html` | `subscriber_count` | Pastikan Edit 2 repository sudah |
| Reject 5 huruf lolos | `service.go` | `len(rejectReason) < 10` | Tambah validasi §10 |
| Payment method kolom kosong | `billing-history.html` | `o.payment_method` | Field sudah di API — FE only §38 |
| client_name tidak muncul | `service.go` | `ClientName: rec.ClientName` | Tambah §43 backend |
| verified_at scan error | `repository.go` | `&rec.VerifiedAt` | Tambah di SELECT + Scan §45 |
| Dashboard days selalu `-` | `dashboard.html` | `days_remaining` | Pastikan §36 di `renderSubscriptionCard` |

---

# CHECKLIST 30 MENIT

| Waktu | Aksi |
|-------|------|
| 00:00–01:00 | Baca soal → Ctrl+F di **Tabel Navigasi** |
| 01:00–02:00 | Buka file #1 → Ctrl+F anchor → konfirmasi ketemu |
| 02:00–12:00 | **Soal C/U/D:** edit service dulu → `go build` → frontend |
| 12:00–14:00 | Test browser/curl dari bagian Verifikasi |
| 14:00–15:00 | **Soal R:** cek "TIDAK PERLU UBAH" atau edit frontend |
| 15:00–25:00 | Implement/demo soal R |
| 25:00–28:00 | Demo ke asdos — sebut file yang diubah |
| 28:00–30:00 | Buffer jika build error |

---

# IMPORT GO YANG SERING DIBUTUHKAN

```go
"strings"   // TrimSpace, ToLower
"errors"    // errors.New
"fmt"       // fmt.Errorf, fmt.Sprintf
"regexp"    // validasi nama
"sort"      // sort.Slice (jika di Go)
"time"      // filter tanggal
```

---

# POLA KODE SIAP PAKAI

### Filter tanggal satu hari (frontend → API existing)
```javascript
if (onDate) {
    params.set('start_date', onDate);
    params.set('end_date', onDate);
}
```

### Sort array frontend setelah fetch
```javascript
var sorted = data.data.slice().sort(function (a, b) {
    return (Number(a.price) || 0) - (Number(b.price) || 0);
});
```

### Validasi string di service
```go
import "regexp"
var validName = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
if !validName.MatchString(payload.Name) {
    return nil, errors.New("nama paket tidak boleh mengandung karakter spesial")
}
```

### Tampilkan field computed yang sudah ada di API
```javascript
daysEl.textContent = (display.days_remaining || 0) + ' hari';
```

### Subquery count di list query
```sql
COALESCE((SELECT COUNT(*)::int FROM subscription.orders o
  WHERE o.package_id = p.id AND o.status IN ('PAID','PENDING_PAYMENT')), 0)
```

### Format tanggal Indonesia
```javascript
new Date(iso).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
```

---

# 5 SOAL PALING MUNGKIN (PREDIKSI FINAL)

| # | Soal | Tipe | Kenapa |
|---|------|------|--------|
| 1 | Filter datepicker satu tanggal di billing history | **R** | Pola persis kelompok news di kelas |
| 2 | Tampilkan subscriber count di modal delete | **C/U/D** | Modifikasi UI + 1 field — dramatic, mudah dinilai |
| 3 | Validasi nama tanpa karakter spesial | **C/U/D** | Pola persis "validasi role permissions sama" |
| 4 | Sort katalog termurah→termahal / badge harga terbaik | **R** | Modifikasi tampilan kecil, tidak ubah backend |
| 5 | Tambah kolom Quota / card total ACTIVE | **R** | Data sudah ada — hanya render |

**Kombinasi ujian paling mungkin:** §2 (R filter tanggal) + §3 atau §1 (C/U/D validasi/tampilan data).

**Batch 2 (§16–§35):** 10 C/U/D + 10 R tambahan — tidak overlap §1–§15.

**Batch 3 (§36–§45):** 10 R/C tambahan — menutup dashboard card, subscription-me, billing, katalog, checkout, ops.

---

# BATCH 2: §16–§35 (SOAL BARU)

> §16–§25 = C/U/D | §26–§35 = R

---

## §16. DEACTIVATE DITOLAK JIKA MASIH ADA SUBSCRIBER

> Tipe: C/U/D | Fitur: Toggle Status | Estimasi: **6 menit**

### Apa yang diubah:
Di `TogglePackageStatus`, sebelum flip ke INACTIVE, cek `CountActiveSubscribers` — sama seperti delete.

### File: `subscription/app/modules/packages/service.go`

**Ctrl+F:**
```go
	newStatus := "INACTIVE"
	if existing.Status == "INACTIVE" {
		newStatus = "ACTIVE"
	}

	return s.repo.TogglePackageStatus(ctx, id, newStatus)
```

**TAMBAHKAN SEBELUM** `return s.repo.TogglePackageStatus`:
```go
	if existing.Status == "ACTIVE" {
		subscriberCount, err := s.repo.CountActiveSubscribers(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error mengecek pelanggan: %v", err)
		}
		if subscriberCount > 0 {
			return nil, errors.New("tidak dapat menonaktifkan paket yang masih memiliki pelanggan aktif")
		}
	}

```

### Build check:
```bash
cd subscription && go build ./...
```

### Verifikasi (30 detik):
Deactivate paket ACTIVE yang punya order PAID/PENDING → **Expected:** error, status tetap ACTIVE.

---

## §17. BLOK UPLOAD BUKTI TRANSFER ULANG

> Tipe: C/U/D | Fitur: Upload Payment Proof | Estimasi: **5 menit**

### File: `subscription/app/modules/orders/service.go`

**Ctrl+F:**
```go
	if rec.Status != "PENDING_PAYMENT" {
		return nil, errors.New("bukti transfer hanya dapat diunggah saat status pesanan PENDING_PAYMENT")
	}

	res, err := s.repo.SavePaymentProof(ctx, orderID, file)
```

**TAMBAHKAN SEBELUM** `res, err := s.repo.SavePaymentProof`:
```go
	if rec.HasPaymentProof {
		return nil, errors.New("bukti transfer sudah diunggah, tidak dapat mengunggah ulang")
	}

```

### Build check + Verifikasi:
Upload bukti 2x pada order yang sama → **Expected:** error pada upload kedua.

---

## §18. VALIDASI FORMAT EMAIL SAAT CREATE ORDER

> Tipe: C/U/D | Fitur: Create Order | Estimasi: **5 menit**

### File: `subscription/app/modules/orders/service.go`

**Tambah import:** `"regexp"`

**Ctrl+F:** `if userID == "" {` di dalam `CreateOrder` (setelah duration default)

**TAMBAHKAN SEBELUM** blok KYC:
```go
	dto.ClientName = strings.TrimSpace(dto.ClientName)
	dto.ClientEmail = strings.TrimSpace(dto.ClientEmail)
	if len(dto.ClientName) < 2 {
		return nil, errors.New("nama client minimal 2 karakter")
	}
	if dto.ClientEmail == "" {
		return nil, errors.New("email client wajib diisi")
	}
	var emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRe.MatchString(dto.ClientEmail) {
		return nil, errors.New("format email client tidak valid")
	}

```

### Verifikasi:
POST order dengan `client_email: "bukan-email"` → **Expected:** 400 format tidak valid.

---

## §19. HARGA TIER HARUS > 0 (BUKAN >= 0)

> Tipe: C/U/D | Fitur: Create/Update Package | Estimasi: **4 menit**

### File: `subscription/app/modules/packages/service.go`

**Ctrl+F (di CreatePackage dan UpdatePackage):**
```go
		if t.Price < 0 {
			return nil, errors.New("harga tidak boleh negatif")
		}
```

**UBAH MENJADI:**
```go
		if t.Price <= 0 {
			return nil, errors.New("harga tier harus lebih besar dari 0")
		}
```

### Verifikasi:
Create paket dengan tier price `0` → **Expected:** error harga tier.

---

## §20. TRIM NAMA PAKET SEBELUM VALIDASI

> Tipe: C/U/D | Fitur: Create/Update Package | Estimasi: **4 menit**

### File: `subscription/app/modules/packages/service.go`

**Ctrl+F di `CreatePackage`:**
```go
func (s *packageService) CreatePackage(ctx context.Context, payload CreatePackageDTO) (*Package, error) {
	if payload.Name == "" {
```

**TAMBAHKAN SETELAH** baris function signature:
```go
	payload.Name = strings.TrimSpace(payload.Name)
```

**Import:** tambah `"strings"` jika belum ada.

Ulangi di awal `UpdatePackage`:
```go
	payload.Name = strings.TrimSpace(payload.Name)
```

### Verifikasi:
Create dengan nama `"  Starter  "` → tersimpan tanpa spasi ujung.

---

## §21. DURASI >1 BULAN WAJIB PUNYA TIER (TANPA FALLBACK)

> Tipe: C/U/D | Fitur: Create Order | Estimasi: **6 menit**

### File: `subscription/app/modules/orders/service.go`

**Ctrl+F:**
```go
	price, err := s.repo.GetPricingTier(ctx, dto.PackageID, dto.DurationMonths)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil harga: %w", err)
	}
	if price == 0 {
		// Fallback: kalkulasi otomatis dari harga dasar
		price = pkg.Price * float64(dto.DurationMonths)
	}

	return s.repo.CreateOrder(ctx, userID, dto, price)
```

**UBAH MENJADI:**
```go
	price, err := s.repo.GetPricingTier(ctx, dto.PackageID, dto.DurationMonths)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil harga: %w", err)
	}
	if price <= 0 {
		return nil, errors.New("durasi yang dipilih tidak tersedia untuk paket ini")
	}

	return s.repo.CreateOrder(ctx, userID, dto, price)
```

### Verifikasi:
Checkout duration 99 bulan (tanpa tier) → **Expected:** error durasi tidak tersedia.

---

## §22. MAKSIMAL 5 PRICING TIER PER PAKET

> Tipe: C/U/D | Fitur: Create/Update Package | Estimasi: **4 menit**

### File: `subscription/app/modules/packages/service.go`

**Ctrl+F (CreatePackage, setelah cek min 1 tier):**
```go
	if len(payload.PricingTiers) == 0 {
		return nil, errors.New("minimal harus ada 1 pilihan durasi harga")
	}
```

**TAMBAHKAN SETELAH:**
```go
	if len(payload.PricingTiers) > 5 {
		return nil, errors.New("maksimal 5 pilihan durasi harga per paket")
	}
```

Copy ke `UpdatePackage`.

### Verifikasi:
Submit 6 tier → **Expected:** error maksimal 5.

---

## §23. NAMA CLIENT MINIMAL 2 KARAKTER (BACKEND)

> Tipe: C/U/D | Fitur: Create Order | Estimasi: **4 menit**

*(Gabung dengan §18 jika sudah dikerjakan — atau hanya blok nama)*

**Ctrl+F di `CreateOrder`** — tambah sebelum KYC:
```go
	if len(strings.TrimSpace(dto.ClientName)) < 2 {
		return nil, errors.New("nama client minimal 2 karakter")
	}
```

### Verifikasi:
POST `client_name: "A"` → ditolak.

---

## §24. KUOTA PAKET MAKSIMAL 1.000.000

> Tipe: C/U/D | Fitur: Create/Update Package | Estimasi: **4 menit**

**Ctrl+F:**
```go
	if payload.Quota <= 0 {
		return nil, errors.New("kuota harus lebih besar dari 0")
	}
```

**TAMBAHKAN SETELAH:**
```go
	if payload.Quota > 1000000 {
		return nil, errors.New("kuota tidak boleh lebih dari 1.000.000")
	}
```

Di `CreatePackage` dan `UpdatePackage`.

---

## §25. ALASAN REJECT MAKSIMAL 500 KARAKTER (BACKEND)

> Tipe: C/U/D | Fitur: Verify Order | Estimasi: **4 menit**

### File: `subscription/app/modules/orders/service.go` — `VerifyOrder`

**Ctrl+F:**
```go
	if action == "REJECT" && rejectReason == "" {
		return nil, errors.New("alasan reject wajib diisi")
	}
```

**UBAH MENJADI:**
```go
	if action == "REJECT" {
		if rejectReason == "" {
			return nil, errors.New("alasan reject wajib diisi")
		}
		if len(rejectReason) > 500 {
			return nil, errors.New("alasan reject maksimal 500 karakter")
		}
	}
```

### Verifikasi:
Reject dengan alasan >500 karakter → ditolak backend.

---

## §26. LABEL DISKON TIER DI KARTU KATALOG

> Tipe: R | Fitur: Katalog Client | Estimasi: **5 menit**

### File: `frontend/client/packages-catalog.html`

**Ctrl+F di `buildPackageCard`:**
```javascript
            var durasiOptions = tiers.length > 0
                ? tiers.map(function (t) { return t.duration_months + ' bln'; }).join(' / ')
                : '1 bln';
```

**TAMBAHKAN SETELAH:**
```javascript
            var promoLabel = '';
            tiers.forEach(function (t) {
                if (t.label && String(t.label).trim()) {
                    promoLabel = String(t.label).trim();
                }
            });
```

**Ctrl+F:**
```javascript
                '          <span class="value">' + durasiOptions + '</span>',
```

**TAMBAHKAN SETELAH** baris durasi (jika promoLabel):
```javascript
                (promoLabel ? '        <div class="pkg-detail-row"><span class="label"><i class="fa-solid fa-percent me-2"></i>Promo</span><span class="value text-warning">' + _escapeHtml(promoLabel) + '</span></div>' : ''),
```

### Verifikasi:
`/client/packages-catalog` → paket dengan tier label "Hemat 20%" menampilkan baris Promo.

---

## §27. BADGE "ADA BUKTI" DI BILLING HISTORY

> Tipe: R | Fitur: Billing History | Estimasi: **5 menit**

### File: `frontend/client/billing-history.html`

**Ctrl+F function `statusBadge`:**
```javascript
        function statusBadge(status, hasProof) {
            if (status === 'PAID') return '<span class="badge-status badge-paid">PAID</span>';
```

**TAMBAHKAN SETELAH** baris PAID (sebelum CANCELLED):
```javascript
            if (status === 'PENDING' && hasProof) {
                return '<span class="badge-status badge-pending">AWAITING VERIFICATION</span>' +
                    ' <span class="badge bg-info text-dark ms-1" style="font-size:0.65rem;">Ada Bukti</span>';
            }
```

*(Sesuaikan jika logic badge sudah overlap — intinya tambah span "Ada Bukti" saat `hasProof`)*

### Verifikasi:
Order PENDING dengan bukti → badge extra "Ada Bukti" muncul.

---

## §28. NAMA PEMESAN DI DETAIL ORDER CLIENT

> Tipe: R | Fitur: Order Detail | Estimasi: **4 menit**

### File: `frontend/client/order-detail.html`

**Ctrl+F HTML:**
```html
                        <div class="info-item">
                            <div class="info-label">Paket</div>
                            <div class="info-value" id="packageName">-</div>
                        </div>
```

**TAMBAHKAN SEBELUM:**
```html
                        <div class="info-item">
                            <div class="info-label">Nama Pemesan</div>
                            <div class="info-value" id="clientName">-</div>
                        </div>
```

**Ctrl+F JS `loadOrderDetail`:**
```javascript
                document.getElementById('packageName').textContent = d.package_name || '-';
```

**TAMBAHKAN SEBELUM:**
```javascript
                document.getElementById('clientName').textContent = d.client_name || '-';
```

*(Pastikan `GET /api/orders/{id}` mengembalikan `client_name` — cek response; jika belum ada, tambah di service mapping detail.)*

### Verifikasi:
Detail order menampilkan nama dari checkout.

---

## §29. ALASAN PENOLAKAN JIKA CANCELLED

> Tipe: R | Fitur: Order Detail | Estimasi: **5 menit**

### File: `frontend/client/order-detail.html`

**TAMBAHKAN** di `info-items` (setelah Metode Pembayaran):
```html
                        <div class="info-item d-none" id="rejectNoteWrap">
                            <div class="info-label">Alasan Penolakan</div>
                            <div class="info-value" id="rejectNote">-</div>
                        </div>
```

**Di `loadOrderDetail` setelah set status:**
```javascript
                var rejectWrap = document.getElementById('rejectNoteWrap');
                var rejectNote = document.getElementById('rejectNote');
                if (rejectWrap && rejectNote) {
                    var note = (d.verification_note || '').trim();
                    var show = orderStatus === 'CANCELLED' && note;
                    rejectWrap.classList.toggle('d-none', !show);
                    rejectNote.textContent = show ? note : '-';
                }
```

### Verifikasi:
Order CANCELLED setelah ops reject → alasan tampil.

---

## §30. IKON BUKTI DI LIST ORDER OPS

> Tipe: R | Fitur: Admin Order List | Estimasi: **5 menit**

### File: `frontend/ops/orders.html`

**Ctrl+F di `renderTable` map:**
```javascript
                    <td class="py-3">${statusBadge(order.status || 'PENDING_PAYMENT')}</td>
```

**UBAH MENJADI:**
```javascript
                    <td class="py-3">${statusBadge(order.status || 'PENDING_PAYMENT')}${order.has_payment_proof ? ' <i class="fa-solid fa-paperclip text-info ms-1" title="Bukti transfer ada"></i>' : ''}</td>
```

### Verifikasi:
`/ops/orders` → order dengan bukti punya ikon paperclip.

---

## §31. JUMLAH TIER HARGA DI TABEL PAKET OPS

> Tipe: R | Fitur: List Package Admin | Estimasi: **4 menit**

### File: `frontend/ops/subscriptions.html`

**Ctrl+F render row:**
```javascript
                    <td class="py-3 pkg-price-yearly">${formatRp(getDurationPrice(pkg, 12))}</td>
```

**TAMBAHKAN SETELAH** (butuh kolom `<th>Tiers</th>` di thead + colspan+1):
```javascript
                    <td class="py-3">${(pkg.pricing_tiers || []).length} tier(s)</td>
```

### Verifikasi:
Tabel ops menampilkan "3 tier(s)" sesuai data API.

---

## §32. TANGGAL DIBUAT PAKET (FORMAT INDONESIA) DI OPS

> Tipe: R | Fitur: List Package Admin | Estimasi: **5 menit**

**Tambah helper di `subscriptions.html`:**
```javascript
        function formatDateID(iso) {
            if (!iso) return '—';
            const d = new Date(iso);
            return isNaN(d.getTime()) ? '—' : d.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' });
        }
```

**Tambah kolom Created** di thead + di row:
```javascript
                    <td class="py-3">${formatDateID(pkg.created_at)}</td>
```

### Verifikasi:
Kolom tanggal format "8 Juni 2026".

---

## §33. HARGA 12 BULAN DI KARTU KATALOG

> Tipe: R | Fitur: Katalog Client | Estimasi: **5 menit**

### File: `frontend/client/packages-catalog.html` — `buildPackageCard`

**Ctrl+F setelah `priceLabel`:**
```javascript
            var priceLabel = '1x ' + displayPrice;
```

**TAMBAHKAN:**
```javascript
            var yearlyTier = tiers.find(function (t) { return Number(t.duration_months) === 12; });
            var yearlyLine = yearlyTier
                ? '<div class="pkg-detail-row"><span class="label"><i class="fa-solid fa-calendar me-2"></i>12 Bulan</span><span class="value">' + formatRupiah(yearlyTier.price) + '</span></div>'
                : '';
```

**Sisipkan `yearlyLine` di array return** setelah baris Harga 1 Bulan.

### Verifikasi:
Kartu paket dengan tier 12 bln menampilkan harga tahunan.

---

## §34. URUTKAN LIST PAKET OPS A–Z

> Tipe: R | Fitur: List Package Admin | Estimasi: **4 menit**

### File: `frontend/ops/subscriptions.html` — `loadPackages`

**Ctrl+F:**
```javascript
                allPackages = json.data || [];
```

**TAMBAHKAN SETELAH:**
```javascript
                allPackages.sort((a, b) => (a.name || '').localeCompare(b.name || '', 'id'));
```

### Verifikasi:
Tabel paket terurut alfabet nama.

---

## §35. TOAST SUKSES SETELAH DOWNLOAD INVOICE

> Tipe: R | Fitur: Invoice PDF | Estimasi: **4 menit**

### File: `frontend/client/billing-history.html` — `downloadInvoice`

**Ctrl+F:**
```javascript
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
```

**TAMBAHKAN SETELAH:**
```javascript
                showToast('Invoice berhasil diunduh.', 'success');
```

Ulangi di `frontend/client/order-detail.html` function `downloadInvoice` jika ada.

### Verifikasi:
Download invoice PAID → toast hijau muncul.

---

# BATCH 3: §36–§45 (GAP CLIENT DASHBOARD / CHECKOUT)

> Menutup halaman yang belum tercakup: **dashboard subscription card**, **subscription-me**, **billing**, **katalog**, **checkout**, **ops orders**.

---

## §36. DASHBOARD — SISA HARI DI SUBSCRIPTION CARD

> Tipe: R | Fitur: Dashboard Subscription Card | Estimasi: **5 menit**

### Apa yang diubah:
Tampilkan `days_remaining` dari API `/api/subscriptions/me` di card subscription dashboard (sama seperti §6 tapi file berbeda).

### File: `frontend/client/dashboard.html`

---

### Edit 1: HTML — setelah baris Nearest Due Date

**Ctrl+F:**
```html
                        <div class="col-md-4">
                            <small class="text-muted text-uppercase" style="font-size:0.68rem;">Nearest Due Date</small>
                            <div class="fw-semibold" id="subEndDate" style="color:var(--text-main);">-</div>
                        </div>
                    </div>
```

**TAMBAHKAN SEBELUM** `</div>` penutup row (setelah blok di atas):
```html
                        <div class="col-md-4">
                            <small class="text-muted text-uppercase" style="font-size:0.68rem;">Sisa Hari Aktif</small>
                            <div class="fw-semibold" id="subDaysRemaining" style="color:var(--accent-cyan);">-</div>
                        </div>
```

---

### Edit 2: function `renderSubscriptionCard`

**Ctrl+F:**
```javascript
            startEl.textContent = formatIdDate(sortedByStartDate[0].start_date);
            endEl.textContent = formatIdDate(sortedByEndDate[0].end_date);

            subListEl.innerHTML = list.map(item => {
```

**TAMBAHKAN SETELAH** `endEl.textContent = ...`:
```javascript
            const nearest = sortedByEndDate[0];
            const daysEl = document.getElementById('subDaysRemaining');
            if (daysEl) {
                const days = typeof nearest.days_remaining === 'number' ? nearest.days_remaining : null;
                daysEl.textContent = days !== null ? (days + ' hari') : '-';
            }
```

**Ctrl+F** blok empty state di `renderSubscriptionCard`:
```javascript
                subListEl.innerHTML = '<div class="text-muted small">No active packages at the moment.</div>';
```

**TAMBAHKAN SEBELUM** `return;` di blok empty:
```javascript
                const daysElEmpty = document.getElementById('subDaysRemaining');
                if (daysElEmpty) daysElEmpty.textContent = '-';
```

### Verifikasi:
Browser: `/client/dashboard` (client dengan sub ACTIVE) → card menampilkan "X hari".

---

## §37. SUBSCRIPTION ME — BADGE EXPIRED WARNA ABU

> Tipe: R | Fitur: Subscription Status | Estimasi: **4 menit**

### File: `frontend/client/subscription-me.html` — `renderSubscriptionOverview`

**Ctrl+F:**
```javascript
            badge.textContent = status === 'ACTIVE' ? 'ACTIVE' : (status === 'EXPIRED' ? 'EXPIRED' : status || '-');
            badge.style.background = status === 'ACTIVE' ? 'rgba(16,185,129,0.18)' : 'rgba(255,255,255,0.1)';
            badge.style.color = status === 'ACTIVE' ? '#34d399' : 'var(--text-muted)';
```

**UBAH MENJADI:**
```javascript
            badge.textContent = status === 'ACTIVE' ? 'ACTIVE' : (status === 'EXPIRED' ? 'EXPIRED' : status || '-');
            if (status === 'ACTIVE') {
                badge.style.background = 'rgba(16,185,129,0.18)';
                badge.style.color = '#34d399';
            } else if (status === 'EXPIRED') {
                badge.style.background = 'rgba(156,163,175,0.22)';
                badge.style.color = '#9ca3af';
            } else {
                badge.style.background = 'rgba(255,255,255,0.1)';
                badge.style.color = 'var(--text-muted)';
            }
```

### Verifikasi:
`/client/subscription-me` dengan sub EXPIRED → badge abu, bukan hijau.

---

## §38. BILLING HISTORY — KOLOM METODE PEMBAYARAN

> Tipe: R | Fitur: Billing History | Estimasi: **5 menit** | **FE saja** — `payment_method` sudah ada di API.

### File: `frontend/client/billing-history.html`

---

### Edit 1: thead — setelah kolom Total Amount

**Ctrl+F:**
```html
                            <th class="small fw-semibold py-3" style="text-transform: uppercase; font-size: 0.72rem;">Total Amount</th>
                            <th class="small fw-semibold py-3" style="text-transform: uppercase; font-size: 0.72rem;">Status</th>
```

**TAMBAHKAN SETELAH** Total Amount:
```html
                            <th class="small fw-semibold py-3" style="text-transform: uppercase; font-size: 0.72rem;">Payment Method</th>
```

---

### Edit 2: `renderHistory` — baris tabel desktop

**Ctrl+F:**
```javascript
                    '<td class="py-3">' + formatRp(o.total_price) + '</td>' +
                    '<td class="py-3 order-status-cell">' + statusBadge(normalized, hasProof) + '</td>' +
```

**UBAH MENJADI:**
```javascript
                    '<td class="py-3">' + formatRp(o.total_price) + '</td>' +
                    '<td class="py-3">' + (o.payment_method || '-') + '</td>' +
                    '<td class="py-3 order-status-cell">' + statusBadge(normalized, hasProof) + '</td>' +
```

---

### Edit 3: mobile card — tambah baris Payment

**Ctrl+F:**
```javascript
                    '<div class="mobile-item-row"><span class="mobile-item-label">Total</span><span class="mobile-item-value">' + formatRp(o.total_price) + '</span></div>' +
```

**TAMBAHKAN SETELAH:**
```javascript
                    '<div class="mobile-item-row"><span class="mobile-item-label">Payment</span><span class="mobile-item-value">' + (o.payment_method || '-') + '</span></div>' +
```

---

### Edit 4: update `colspan` skeleton/empty dari **6** ke **7**

**Ctrl+F:** `colspan="6"` → ganti semua jadi `colspan="7"`.

### Verifikasi:
`/client/billing-history` → kolom "Transfer Bank" muncul.

---

## §39. CHECKOUT — LABEL TIER DI RINGKASAN

> Tipe: R | Fitur: Checkout | Estimasi: **6 menit**

### File: `frontend/client/checkout.html`

---

### Edit 1: HTML ringkasan — setelah Periode

**Ctrl+F:**
```html
                                <div class="sum-label">Periode</div>
                                <div class="sum-value" id="sumPeriod">—</div>

                                <hr class="sum-divider">
```

**TAMBAHKAN SETELAH** `sumPeriod`:
```html
                                <div class="sum-label" id="sumTierLabelWrap">Label Promo</div>
                                <div class="sum-value" id="sumTierLabel">—</div>
```

---

### Edit 2: function `updateDurationPricing`

**Ctrl+F:**
```javascript
            document.getElementById('sumPeriod').textContent = months + ' bulan';
            document.getElementById('sumSubtotal').textContent = formatRp(selectedTier.price);
```

**TAMBAHKAN SETELAH:**
```javascript
            var tierLabel = (selectedTier.label || '').trim();
            var tierLabelEl = document.getElementById('sumTierLabel');
            var tierWrap = document.getElementById('sumTierLabelWrap');
            if (tierLabelEl) tierLabelEl.textContent = tierLabel || '-';
            if (tierWrap) tierWrap.classList.toggle('d-none', !tierLabel);
```

### Verifikasi:
Checkout paket dengan tier label "Hemat 20%" → ringkasan menampilkan label.

---

## §40. BILLING HISTORY — SOROT BARIS PAID HIJAU

> Tipe: R | Fitur: Billing History | Estimasi: **4 menit**

### File: `frontend/client/billing-history.html`

---

### Edit 1: CSS — di blok `<style>`

**Ctrl+F:**
```css
        .badge-status {
```

**TAMBAHKAN SEBELUM:**
```css
        .history-table tbody tr.row-paid td {
            background: rgba(16, 185, 129, 0.08) !important;
        }

```

---

### Edit 2: `renderHistory` — class pada `<tr>`

**Ctrl+F:**
```javascript
                return '<tr data-order-id="' + o.order_id + '">' +
```

**UBAH MENJADI:**
```javascript
                var rowPaidClass = status === 'PAID' ? ' row-paid' : '';
                return '<tr class="' + rowPaidClass.trim() + '" data-order-id="' + o.order_id + '">' +
```

*(Alternatif lebih simpel: `return '<tr' + (status === 'PAID' ? ' class="row-paid"' : '') + ' data-order-id="' + o.order_id + '">'`)*

### Verifikasi:
Baris order PAID punya background hijau transparan.

---

## §41. KATALOG — EMPTY STATE + TOMBOL MUAT ULANG

> Tipe: R | Fitur: Packages Catalogue | Estimasi: **4 menit**

### File: `frontend/client/packages-catalog.html`

**Ctrl+F:**
```html
                <h6 class="fw-semibold mb-1" style="color: var(--text-heading)">Belum Ada Paket Aktif</h6>
                <p class="text-muted small mb-0">Saat ini belum ada paket yang tersedia. Coba lagi nanti.</p>
            </div>
        </div>
```

**UBAH `mb-0` pada `<p>` dan TAMBAHKAN tombol:**
```html
                <h6 class="fw-semibold mb-1" style="color: var(--text-heading)">Belum Ada Paket Aktif</h6>
                <p class="text-muted small mb-3">Saat ini belum ada paket yang tersedia. Coba lagi nanti.</p>
                <button type="button" class="btn btn-sm btn-outline-light" onclick="fetchCatalog()">
                    <i class="fa-solid fa-rotate-right me-2"></i>Muat Ulang
                </button>
            </div>
        </div>
```

### Verifikasi:
Katalog kosong → tombol Muat Ulang memanggil `fetchCatalog()` (sama pola error state).

---

## §42. OPS ORDERS — TEKS "MENUNGGU BUKTI" UNTUK PENDING

> Tipe: R | Fitur: Admin Order List | Estimasi: **5 menit**

### File: `frontend/ops/orders.html`

**Ctrl+F di `renderTable`:**
```javascript
                    <td class="py-3">${statusBadge(order.status || 'PENDING_PAYMENT')}</td>
```

**UBAH MENJADI:**
```javascript
                    <td class="py-3">${statusBadge(order.status || 'PENDING_PAYMENT')}${normalizeStatus(order.status) === 'PENDING_PAYMENT' && !order.has_payment_proof ? ' <small class="text-muted ms-1">Menunggu Bukti</small>' : ''}</td>
```

### Verifikasi:
`/ops/orders` → order PENDING tanpa bukti menampilkan teks "Menunggu Bukti".

---

## §43. DETAIL ORDER CLIENT — `client_name` DI API ⚠️ BACKEND

> Tipe: R | Fitur: Order Detail | Estimasi: **6 menit** | Lengkapi §28 jika field belum di response.

### File:
1. `subscription/app/modules/orders/types.go`
2. `subscription/app/modules/orders/service.go`
3. `frontend/client/order-detail.html` (sudah di §28)

---

### Edit 1: `types.go` — `ClientOrderDetail`

**Ctrl+F:**
```go
	PackageName            string     `json:"package_name"`
	DurationMonths         int        `json:"duration_months"`
```

**TAMBAHKAN:**
```go
	ClientName             string     `json:"client_name"`
```

---

### Edit 2: `service.go` — `GetOrderDetailForClient`

**Ctrl+F:**
```go
	return &ClientOrderDetail{
		OrderID:                rec.OrderID,
		InvoiceNumber:          rec.InvoiceNumber,
		PackageName:            rec.PackageName,
```

**TAMBAHKAN** setelah `PackageName`:
```go
		ClientName:             rec.ClientName,
```

### Build:
```bash
cd subscription && go build ./...
```

### Verifikasi:
`GET /api/orders/{id}` → JSON ada `client_name` → UI §28 menampilkan nama.

---

## §44. DASHBOARD — TOAST SUKSES RENEW (BUKAN ALERT)

> Tipe: R | Fitur: Dashboard Renew | Estimasi: **5 menit**

### File: `frontend/client/dashboard.html`

---

### Edit 1: HTML toast — sebelum `<footer class="fixed-footer">`

**Ctrl+F:**
```html
    <footer class="fixed-footer">
```

**TAMBAHKAN SEBELUM:**
```html
    <div class="position-fixed bottom-0 end-0 p-3" style="z-index:9999">
        <div id="toastEl" class="toast align-items-center text-white border-0" role="alert">
            <div class="d-flex">
                <div class="toast-body fw-semibold" id="toastMsg">-</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        </div>
    </div>

```

---

### Edit 2: helper `showToast` — di dalam `<script>`, sebelum `submitRenew`

```javascript
        function showToast(msg, type) {
            const t = document.getElementById('toastEl');
            const m = document.getElementById('toastMsg');
            if (!t || !m) return;
            t.classList.remove('bg-success', 'bg-danger', 'bg-warning');
            t.classList.add(type === 'success' ? 'bg-success' : type === 'warning' ? 'bg-warning' : 'bg-danger');
            m.textContent = msg;
            bootstrap.Toast.getOrCreateInstance(t, { delay: 3500 }).show();
        }
```

---

### Edit 3: `submitRenew` — ganti `alert` error

**Ctrl+F:**
```javascript
                    alert(json.error_message || 'Failed to submit renewal. Please try again.');
```

**UBAH MENJADI:**
```javascript
                    showToast(json.error_message || 'Failed to submit renewal. Please try again.', 'danger');
```

**Ctrl+F:**
```javascript
                alert('Failed to connect to the server. Please try again.');
```

**UBAH MENJADI:**
```javascript
                showToast('Failed to connect to the server. Please try again.', 'danger');
```

### Verifikasi:
Renew gagal → toast merah (bukan alert browser). Sukses tetap redirect ke billing.

---

## §45. OPS ORDER DETAIL — WAKTU VERIFIKASI (`verified_at`) ⚠️ KOMPLEKS

> Tipe: R | Fitur: Ops Order Detail | Estimasi: **8 menit** | types → repository → service → frontend

### File:
1. `subscription/app/modules/orders/types.go`
2. `subscription/app/modules/orders/repository.go`
3. `subscription/app/modules/orders/service.go`
4. `frontend/ops/orders-detail.html`

---

### Edit 1: `types.go` — `OrderRecord` dan `AdminOrderDetail`

**Ctrl+F `OrderRecord`:**
```go
	CreatedAt               time.Time
}
```

**UBAH MENJADI:**
```go
	VerifiedAt              *time.Time
	CreatedAt               time.Time
}
```

**Ctrl+F `AdminOrderDetail` — setelah `CreatedAt`:**
```go
	CreatedAt              time.Time  `json:"created_at"`
}
```

**UBAH MENJADI:**
```go
	CreatedAt              time.Time  `json:"created_at"`
	VerifiedAt             *time.Time `json:"verified_at,omitempty"`
}
```

---

### Edit 2: `repository.go` — SELECT detail (kedua query `getOrderByIDWithFallback`)

**Ctrl+F** (di `baseWithUsers` dan `baseWithoutUsers` — ubah keduanya):
```go
		o.created_at
	FROM subscription.orders o
```

**UBAH MENJADI:**
```go
		o.verified_at,
		o.created_at
	FROM subscription.orders o
```

**Ctrl+F Scan:**
```go
		&rec.CreatedAt,
	)
```

**UBAH MENJADI:**
```go
		&rec.VerifiedAt,
		&rec.CreatedAt,
	)
```

---

### Edit 3: `service.go` — `GetOrderDetailForAdmin`

**Ctrl+F:**
```go
		CreatedAt: rec.CreatedAt,
	}, nil
}
```

*(di function `GetOrderDetailForAdmin`)*

**UBAH `CreatedAt` line menjadi:**
```go
		CreatedAt:  rec.CreatedAt,
		VerifiedAt: rec.VerifiedAt,
```

---

### Edit 4: `orders-detail.html` — tampilkan jika PAID

**Ctrl+F** setelah `dTotal`:
```html
                            <div class="info-item">
                                <div class="info-label">Total Pembayaran</div>
                                <div class="info-value" id="dTotal">-</div>
                            </div>
                        </div>
```

**TAMBAHKAN SEBELUM** `</div>` penutup info-items:
```html
                            <div class="info-item d-none" id="verifiedAtWrap">
                                <div class="info-label">Waktu Verifikasi</div>
                                <div class="info-value" id="dVerifiedAt">-</div>
                            </div>
```

**Ctrl+F di `renderDetail` setelah set `dTotal`:**
```javascript
            document.getElementById('dTotal').textContent = formatRp(order.total_price);

            const note = String(order.verification_note || '').trim();
```

**TAMBAHKAN SETELAH:**
```javascript
            const verifiedWrap = document.getElementById('verifiedAtWrap');
            const verifiedEl = document.getElementById('dVerifiedAt');
            if (verifiedWrap && verifiedEl) {
                const show = status === 'PAID' && order.verified_at;
                verifiedWrap.classList.toggle('d-none', !show);
                verifiedEl.textContent = show
                    ? new Date(order.verified_at).toLocaleString('id-ID', { day: '2-digit', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' })
                    : '-';
            }
```

### Build + Verifikasi:
```bash
cd subscription && go build ./...
```
`/ops/orders-detail?id=<paid>` → "Waktu Verifikasi" tampil setelah approve.

---

# PETA COVERAGE LENGKAP (SEMUA HALAMAN PAKET)

| Halaman | URL | § yang cover | Status |
|---------|-----|--------------|--------|
| Katalog client | `/client/packages-catalog` | §4, §11, §26, §33, §41 | ✅ lengkap |
| Checkout | `/client/checkout` | §18, §21, §39 + KYC sudah ada | ✅ cukup |
| Billing history | `/client/billing-history` | §2, §12, §27, §35, §38, §40, §R-INV | ✅ lengkap |
| Order detail client | `/client/order-detail` | §7, §8, §28, §29, §35, §43, §R-PROOF | ✅ lengkap |
| Subscription status | `/client/subscription-me` | §6, §37 + renew sudah ada | ✅ lengkap |
| Dashboard card | `/client/dashboard` | §36, §44 | ✅ baru dilengkapi |
| List paket ops | `/ops/subscriptions` | §1, §5, §9, §15, §16, §31–§34 | ✅ lengkap |
| Create/edit paket ops | `/ops/subscriptions-create/edit` | §3, §13, §19–§22, §24 | ✅ via service |
| List order ops | `/ops/orders` | §14, §30, §42 | ✅ lengkap |
| Detail order ops | `/ops/orders-detail` | §10, §45, §R-VERIFY | ✅ lengkap |

**Total skenario operasional: §1–§45 + §R-INV/PROOF/VERIFY = 48**

---

# 5 SOAL BATCH 2 PALING MUNGKIN

| # | Soal | Tipe |
|---|------|------|
| 1 | §16 deactivate blok subscriber | C/U/D |
| 2 | §17 blok upload bukti ulang | C/U/D |
| 3 | §18 validasi email | C/U/D |
| 4 | §27 badge Ada Bukti billing | R |
| 5 | §30 ikon bukti ops orders | R |

---

# 5 SOAL BATCH 3 PALING MUNGKIN (CLIENT AREA)

| # | Soal | Tipe |
|---|------|------|
| 1 | §36 sisa hari di dashboard subscription card | R |
| 2 | §38 kolom payment method billing | R |
| 3 | §39 label tier checkout | R |
| 4 | §40 sorot baris PAID hijau | R |
| 5 | §42 Menunggu Bukti di ops orders | R |

---

*Cheat sheet v2.2 — §1–§45 + §R. Scan kode: Juni 2026.*
