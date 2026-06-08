# Panduan Persiapan Soal Latihan Sprint 3

Dokumen ini memandu pengerjaan 10 soal latihan berdasarkan 4 fitur utama:
**KYC Resubmission · Login/Logout/Register · Superadmin Manajemen Akun Internal · Superadmin Memantau Kinerja Divisi**

---

## Peserta 1

### Soal 1 — KYC Resubmission: Tampilkan Alasan Penolakan sebagai Alert/Banner

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Sebelum form resubmission muncul, tampilkan alasan penolakan KYC secara mencolok menggunakan banner/alert agar user tahu dokumen mana yang perlu diperbaiki.

**Status kode saat ini:**
Banner penolakan sudah ada di file (`#rejectionBanner`, baris 115–123), tetapi hanya tampil di dalam form. Tugas ini meminta banner tersebut ditampilkan **secara lebih menonjol**, misalnya full-width di atas form, bukan di dalam card.

**Langkah pengerjaan:**
1. Pindahkan atau duplikat `#rejectionBanner` ke luar `<div class="ai-card">`, tepat di bawah header halaman.
2. Tambahkan class CSS agar banner lebih besar dan mencolok (misalnya `alert alert-danger` Bootstrap).
3. Pastikan teks `#rejectionReasonText` masih diisi dari JavaScript (`data.data.rejection_reason`).

**Contoh markup banner tambahan:**
```html
<!-- Di luar ai-card, sebelum resubmitFormSection -->
<div id="topRejectionAlert" class="alert d-flex align-items-start gap-3 mb-4" 
     style="background:rgba(239,68,68,0.12); border:1px solid rgba(239,68,68,0.4); border-radius:12px; display:none !important;">
    <i class="fa-solid fa-triangle-exclamation fa-xl" style="color:#ef4444; margin-top:2px;"></i>
    <div>
        <strong style="color:#ef4444;">KYC Anda Ditolak</strong>
        <p class="mb-0 mt-1 small" style="color:var(--text-muted);" id="topRejectionText"></p>
    </div>
</div>
```

**JavaScript — isi teks banner:**
```js
// Di dalam loadKYCData(), setelah baris: $('resubmitFormSection').style.display = 'block';
var topAlert = $('topRejectionAlert');
$('topRejectionText').textContent = reason || 'Pengajuan Anda sebelumnya ditolak.';
topAlert.style.display = 'flex';
```

---

### Soal 2 — Login: Tambahkan Tombol Show/Hide Password

**File:** `frontend/account/login.html`

**Apa yang harus dilakukan:**
Tambahkan tombol toggle show/hide password di field password pada halaman login.

**Status kode saat ini:**
Toggle sudah **ada** di baris 85–87 (elemen `#togglePassword`) dan JavaScript-nya di baris 116–121. **Fitur ini sudah jadi.** Pastikan kamu memverifikasi fungsinya benar-benar bekerja, lalu dokumentasikan sebagai selesai atau tingkatkan tampilannya bila perlu.

> **Tip:** Kalau sudah ada, fokuskan peningkatan ke UX — misalnya tambahkan tooltip `"Tampilkan password"` / `"Sembunyikan password"` agar lebih informatif.

**Contoh peningkatan (tooltip):**
```js
const toggleBtn = document.getElementById('togglePassword');
toggleBtn.setAttribute('title', 'Tampilkan password');

toggleBtn.addEventListener('click', function() {
    const pw = document.getElementById('password');
    const icon = this.querySelector('i');
    if (pw.type === 'password') {
        pw.type = 'text';
        icon.classList.replace('fa-eye-slash', 'fa-eye');
        this.setAttribute('title', 'Sembunyikan password');
    } else {
        pw.type = 'password';
        icon.classList.replace('fa-eye', 'fa-eye-slash');
        this.setAttribute('title', 'Tampilkan password');
    }
});
```

---

## Peserta 2

### Soal 3 — Register: Indikator Kekuatan Password (Real-time)

**File:** `frontend/account/register.html`

**Apa yang harus dilakukan:**
Tambahkan indikator kekuatan password (weak / medium / strong) yang muncul secara real-time saat user mengisi field `#password`.

**Field password yang ada** (baris 119–129):
```html
<input type="password" class="form-control" id="password" placeholder="Minimal 8 karakter" required>
<div class="validation-msg" id="passwordMsg"></div>
```

**Langkah pengerjaan:**
1. Tambahkan elemen indikator di bawah `#password` (setelah `#passwordMsg`):

```html
<div id="strengthBar" class="mt-2" style="display:none;">
    <div style="height:4px; border-radius:4px; background:#e2e8f0; overflow:hidden;">
        <div id="strengthFill" style="height:100%; width:0%; border-radius:4px; transition:all 0.3s;"></div>
    </div>
    <small id="strengthLabel" class="mt-1 d-block" style="font-size:0.75rem;"></small>
</div>
```

2. Tambahkan fungsi JavaScript:

```js
function checkPasswordStrength(val) {
    var bar = document.getElementById('strengthBar');
    var fill = document.getElementById('strengthFill');
    var label = document.getElementById('strengthLabel');
    if (!val) { bar.style.display = 'none'; return; }
    bar.style.display = 'block';

    var score = 0;
    if (val.length >= 8) score++;
    if (/[A-Z]/.test(val)) score++;
    if (/[0-9]/.test(val)) score++;
    if (/[^A-Za-z0-9]/.test(val)) score++;

    var configs = {
        1: { width:'25%', color:'#dc2626', text:'Lemah' },
        2: { width:'50%', color:'#f59e0b', text:'Sedang' },
        3: { width:'75%', color:'#3b82f6', text:'Kuat' },
        4: { width:'100%', color:'#10b981', text:'Sangat Kuat' }
    };
    var cfg = configs[score] || configs[1];
    fill.style.width = cfg.width;
    fill.style.background = cfg.color;
    label.textContent = cfg.text;
    label.style.color = cfg.color;
}

// Tambahkan event listener di input password
document.getElementById('password').addEventListener('input', function() {
    checkPasswordStrength(this.value);
});
```

---

### Soal 4 — Superadmin Manajemen Akun: Konfirmasi Dialog Sebelum Hapus

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Tampilkan modal konfirmasi sebelum admin menghapus akun internal, dengan nama dan email akun ditampilkan di dalam dialog.

**Pola yang sudah ada di file ini** (baris 158–196):
File sudah punya sistem modal overlay (`#confirmModal`, `.modal-overlay`, `.modal-box`). Gunakan pola yang sama untuk aksi hapus.

**Langkah pengerjaan:**
1. Cari fungsi yang dipanggil saat tombol hapus diklik (biasanya `deleteUser(id)` atau serupa).
2. Sebelum memanggil API delete, tampilkan modal konfirmasi terlebih dahulu.
3. Isi nama & email user di dalam modal.

**Contoh pola konfirmasi (mengikuti modal yang sudah ada):**
```js
function confirmDeleteUser(userId, userName, userEmail) {
    // Isi konten modal
    document.getElementById('modalTitle').textContent = 'Hapus Akun Internal';
    document.getElementById('modalMessage').innerHTML =
        'Anda akan menghapus akun:<br><strong>' + userName + '</strong> (' + userEmail + ')<br>' +
        '<span class="text-danger small">Tindakan ini tidak dapat dibatalkan.</span>';

    // Set tombol confirm
    document.getElementById('btnConfirmAction').onclick = function() {
        closeModal();
        deleteUser(userId); // panggil fungsi hapus asli
    };

    openModal('danger');
}
```

---

## Peserta 3

### Soal 5 — Superadmin Kinerja Divisi: Filter Date Range

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan filter rentang tanggal (*date range picker*) agar superadmin dapat memfilter data kinerja berdasarkan periode tertentu.

**Langkah pengerjaan:**
1. Tambahkan dua input `type="date"` (tanggal mulai dan tanggal akhir) di area header halaman (di samping tombol Refresh):

```html
<div class="d-flex align-items-center gap-2">
    <label class="small text-muted mb-0">Dari:</label>
    <input type="date" id="filterDateFrom" class="form-control form-control-sm"
           style="width:160px; background:var(--bg-card); border-color:var(--border-color); color:var(--text-heading);">
    <label class="small text-muted mb-0">Sampai:</label>
    <input type="date" id="filterDateTo" class="form-control form-control-sm"
           style="width:160px; background:var(--bg-card); border-color:var(--border-color); color:var(--text-heading);">
    <button class="btn btn-sm btn-outline-secondary" onclick="applyDateFilter()">
        <i class="fa-solid fa-filter me-1"></i>Filter
    </button>
</div>
```

2. Modifikasi fungsi `loadData()` yang sudah ada agar mengirim parameter tanggal ke API:

```js
function applyDateFilter() {
    var from = document.getElementById('filterDateFrom').value;
    var to   = document.getElementById('filterDateTo').value;
    if (from && to && from > to) {
        alert('Tanggal mulai tidak boleh lebih besar dari tanggal akhir.');
        return;
    }
    loadData(from, to);
}

async function loadData(dateFrom, dateTo) {
    var url = '/api/management/overview';
    var params = [];
    if (dateFrom) params.push('date_from=' + dateFrom);
    if (dateTo)   params.push('date_to='   + dateTo);
    if (params.length) url += '?' + params.join('&');

    // ... lanjutkan fetch seperti sebelumnya
}
```

---

### Soal 6 — KYC Resubmission: Preview Dokumen Sebelum Submit

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Tampilkan preview dokumen yang telah diunggah **sebelum** user menekan tombol submit.

**Status kode saat ini:**
Preview sudah **ada** (`#previewContainer`, `#previewImage`) dan diisi di fungsi `handleFile()` baris 482–486. Fitur ini sudah berjalan. Peningkatan yang bisa dilakukan:

- Pastikan preview terlihat jelas (ukuran cukup besar)
- Tambahkan label "Preview KTP Baru:" di atas gambar preview
- Bedakan secara visual antara "KTP lama" (`#oldKtpContainer`) dan "KTP baru" (`#previewContainer`)

**Contoh perbaikan label:**
```html
<!-- Ganti previewContainer yang sudah ada -->
<div class="preview-container" id="previewContainer" style="display:none;">
    <small class="text-muted d-block mb-1">
        <i class="fa-solid fa-eye me-1"></i>Preview KTP Baru:
    </small>
    <img id="previewImage" src="" alt="Preview KTP"
         style="max-width:100%; max-height:300px; border-radius:8px; border:2px solid #22c55e;">
    <div class="file-info" id="fileInfo"></div>
</div>
```

---

## Peserta 4

### Soal 7 — Logout: Konfirmasi Dialog Sebelum Keluar

**File:** `frontend/assets/js/ops-layout.js` atau komponen sidebar (cek `ops-sidebar-placeholder`)

**Apa yang harus dilakukan:**
Tampilkan dialog konfirmasi dengan pilihan "Ya, Keluar" dan "Batal" sebelum user logout.

**Langkah pengerjaan:**
1. Cari tombol/link logout di sidebar layout (`ops-layout.js` atau `client-layout.js`).
2. Cegah navigasi langsung, tampilkan modal dulu.

**Contoh implementasi:**
```js
// Intercept klik logout
document.addEventListener('click', function(e) {
    var logoutBtn = e.target.closest('[data-action="logout"], a[href*="/logout"]');
    if (!logoutBtn) return;
    e.preventDefault();
    showLogoutConfirm(logoutBtn.href || '/account/logout');
});

function showLogoutConfirm(logoutUrl) {
    var confirmed = confirm('Apakah Anda yakin ingin keluar?');
    // Ganti confirm() dengan modal custom jika ada
    if (confirmed) window.location.href = logoutUrl;
}
```

**Jika menggunakan modal Bootstrap:**
```html
<!-- Tambahkan di body layout -->
<div class="modal fade" id="logoutModal" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content" style="background:var(--bg-card); border:1px solid var(--border-color);">
            <div class="modal-body text-center p-4">
                <i class="fa-solid fa-right-from-bracket fa-2x mb-3" style="color:#ef4444;"></i>
                <h5 class="fw-bold" style="color:var(--text-heading);">Keluar dari Aplikasi?</h5>
                <p class="text-muted small">Anda akan logout dari sesi ini.</p>
                <div class="d-flex gap-2 justify-content-center mt-3">
                    <button class="btn btn-outline-secondary" data-bs-dismiss="modal">Batal</button>
                    <button class="btn btn-danger" id="btnConfirmLogout">Ya, Keluar</button>
                </div>
            </div>
        </div>
    </div>
</div>
```

---

### Soal 8 — Superadmin Manajemen Akun: Kolom "Terakhir Aktif" + Sortable

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Tambahkan kolom "Terakhir Aktif" di tabel daftar akun internal, dan buat kolom tersebut bisa diurutkan (sortable).

**Langkah pengerjaan:**
1. Cari `<table id="usersTable">` dan tambahkan header kolom baru:

```html
<th style="cursor:pointer;" onclick="sortTable('last_active')">
    Terakhir Aktif
    <i class="fa-solid fa-sort ms-1 text-muted" style="font-size:0.75rem;" id="sortIcon_last_active"></i>
</th>
```

2. Saat render baris tabel dari data API, tambahkan cell:

```js
// Dalam loop render rows
var lastActive = user.last_active
    ? new Date(user.last_active).toLocaleString('id-ID', { dateStyle:'medium', timeStyle:'short' })
    : 'Belum pernah aktif';

row.innerHTML += '<td><span class="small text-muted">' + lastActive + '</span></td>';
```

3. Tambahkan logika sort:

```js
var sortDirection = {};

function sortTable(field) {
    sortDirection[field] = sortDirection[field] === 'asc' ? 'desc' : 'asc';
    usersData.sort(function(a, b) {
        var valA = a[field] || '';
        var valB = b[field] || '';
        return sortDirection[field] === 'asc'
            ? valA.localeCompare(valB)
            : valB.localeCompare(valA);
    });
    renderTable();
}
```

---

## Peserta 5

### Soal 9 — Register: Dropdown Country Code di Field Telepon

**File:** `frontend/account/register.html`

**Apa yang harus dilakukan:**
Tambahkan dropdown untuk memilih kode negara sebelum field nomor telepon.

**Field telepon yang ada** (baris 100–106):
```html
<div class="input-group">
    <span class="input-group-text"><i class="fa-solid fa-phone"></i></span>
    <input type="tel" class="form-control" id="phone" placeholder="08xxxxxxxxxx" required>
</div>
```

**Langkah pengerjaan:**
Ganti markup field telepon menjadi:

```html
<div class="input-group">
    <span class="input-group-text"><i class="fa-solid fa-phone"></i></span>
    <select class="form-select" id="countryCode" style="max-width:110px; border-radius:0; border-left:none; border-right:none;">
        <option value="+62" selected>🇮🇩 +62</option>
        <option value="+60">🇲🇾 +60</option>
        <option value="+65">🇸🇬 +65</option>
        <option value="+63">🇵🇭 +63</option>
        <option value="+66">🇹🇭 +66</option>
        <option value="+1">🇺🇸 +1</option>
    </select>
    <input type="tel" class="form-control" id="phone" placeholder="8xxxxxxxxxx" required>
</div>
<div class="validation-msg" id="phoneMsg"></div>
```

**Saat submit, gabungkan country code + nomor:**
```js
// Dalam event submit form
var code = document.getElementById('countryCode').value;  // e.g. "+62"
var num  = document.getElementById('phone').value.replace(/^0+/, ''); // hapus leading 0
var fullPhone = code + num;  // e.g. "+628123456789"
```

---

### Soal 10 — Superadmin Kinerja Divisi: Card Summary di Atas Halaman

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan card summary di bagian atas halaman yang menampilkan **total task selesai**, **total task pending**, dan **rata-rata waktu penyelesaian** untuk semua divisi.

**Langkah pengerjaan:**
1. Tambahkan section summary card sebelum `<div class="row g-4" id="cardsRow">`:

```html
<div class="row g-3 mb-4" id="summaryCards">
    <div class="col-md-4">
        <div class="division-card p-3 text-center">
            <div class="metric-val text-success" id="summaryDone">—</div>
            <div class="metric-label mt-1">
                <i class="fa-solid fa-circle-check text-success me-1"></i>Total Task Selesai
            </div>
        </div>
    </div>
    <div class="col-md-4">
        <div class="division-card p-3 text-center">
            <div class="metric-val text-warning" id="summaryPending">—</div>
            <div class="metric-label mt-1">
                <i class="fa-solid fa-clock text-warning me-1"></i>Total Task Pending
            </div>
        </div>
    </div>
    <div class="col-md-4">
        <div class="division-card p-3 text-center">
            <div class="metric-val" style="color:#60a5fa;" id="summaryAvgTime">—</div>
            <div class="metric-label mt-1">
                <i class="fa-solid fa-stopwatch me-1" style="color:#60a5fa;"></i>Rata-rata Waktu Selesai
            </div>
        </div>
    </div>
</div>
```

2. Setelah data API diterima, hitung agregat dan isi card:

```js
// Di dalam loadData(), setelah data divisions diterima
function updateSummaryCards(divisions) {
    var totalDone    = 0;
    var totalPending = 0;
    var totalAvgMs   = 0;
    var count        = 0;

    divisions.forEach(function(div) {
        totalDone    += (div.resolved || div.done || 0);
        totalPending += (div.pending  || div.open  || 0);
        if (div.avg_resolution_time) { totalAvgMs += div.avg_resolution_time; count++; }
    });

    document.getElementById('summaryDone').textContent    = totalDone;
    document.getElementById('summaryPending').textContent = totalPending;
    document.getElementById('summaryAvgTime').textContent = count
        ? (totalAvgMs / count / 3600000).toFixed(1) + ' jam'
        : 'N/A';
}
```

---

## Peserta 6

### Soal 11 — Superadmin Manajemen Akun: Tampilkan Jumlah Hasil Filter

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Tampilkan badge/teks dinamis di atas tabel yang menginformasikan jumlah akun yang sedang ditampilkan berdasarkan filter aktif, misalnya `"3 akun ditemukan"` atau `"Menampilkan semua 12 akun"`.

**Konteks kode saat ini:**
Pagination sudah menampilkan `"Menampilkan X – Y dari Z data"` di bagian bawah tabel (baris 453–461). Soal ini meminta penambahan badge informatif yang berbeda — ditaruh **di atas tabel**, lebih menonjol, dan berubah realtime mengikuti filter.

**Langkah pengerjaan:**

1. Tambahkan elemen di antara filter card dan tabel:

```html
<!-- Taruh antara .filter-card dan .table-card -->
<div class="d-flex justify-content-between align-items-center mb-2 px-1">
    <span id="resultCountBadge" class="small text-muted"></span>
    <span id="activeFilterTags" class="d-flex gap-1 flex-wrap"></span>
</div>
```

2. Di dalam fungsi `loadUsers()`, setelah baris `const total = payload.total || 0;`, tambahkan:

```js
// Update result count badge
var badge = document.getElementById('resultCountBadge');
var search = document.getElementById('searchInput').value.trim();
var role   = document.getElementById('filterRole').value;
var status = document.getElementById('filterStatus').value;
var isFiltered = search || role || status;

badge.textContent = isFiltered
    ? total + ' akun ditemukan'
    : 'Total ' + total + ' akun terdaftar';
badge.style.color = isFiltered ? 'var(--accent-cyan)' : '';

// Render active filter tags
var tagsEl = document.getElementById('activeFilterTags');
tagsEl.innerHTML = '';
if (search)  tagsEl.innerHTML += '<span class="badge bg-secondary">Cari: "' + escapeHtml(search) + '"</span>';
if (role)    tagsEl.innerHTML += '<span class="badge bg-secondary">Role: ' + escapeHtml(role) + '</span>';
if (status)  tagsEl.innerHTML += '<span class="badge bg-secondary">Status: ' + (status === 'active' ? 'Aktif' : 'Nonaktif') + '</span>';
```

3. Pastikan `resetFilters()` juga me-reset tampilan badge (akan otomatis terjadi karena `loadUsers(1)` dipanggil ulang).

---

### Soal 12 — Superadmin Kinerja Divisi: Sort Kartu Berdasarkan Task Pending

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan tombol sort di atas grid kartu untuk mengurutkan kartu divisi dari jumlah task pending **terbanyak ke tersedikit**.

**Konteks kode saat ini:**
Grid kartu dirender secara statis di HTML (baris 126–215). Data dimuat dari API di `loadData()` tetapi hasil fetch hanya mengisi angka ke elemen yang sudah ada — urutan kartu tidak berubah.

**Langkah pengerjaan:**

1. Ubah pendekatan: simpan data divisi ke array saat diterima dari API, lalu render kartu secara dinamis. Tambahkan tombol sort:

```html
<!-- Di samping tombol Refresh -->
<button class="btn btn-sm btn-outline-secondary" id="btnSortPending" onclick="toggleSortByPending()">
    <i class="fa-solid fa-arrow-down-wide-short me-1"></i>Sort: Pending Terbanyak
</button>
```

2. Simpan data divisi dan state sort:

```js
var divisionData = [];  // cache data dari API
var sortByPending = false;

function toggleSortByPending() {
    sortByPending = !sortByPending;
    var btn = document.getElementById('btnSortPending');
    btn.innerHTML = sortByPending
        ? '<i class="fa-solid fa-arrow-up-wide-short me-1"></i>Sort: Pending Terbanyak'
        : '<i class="fa-solid fa-arrow-down-wide-short me-1"></i>Sort: Default';
    renderCards();
}

function renderCards() {
    var data = divisionData.slice(); // copy
    if (sortByPending) {
        data.sort(function(a, b) { return (b.pending || 0) - (a.pending || 0); });
    }
    // render urutan kartu sesuai `data`
    var row = document.getElementById('cardsRow');
    // ubah order CSS atau re-inject kartu sesuai urutan `data`
    data.forEach(function(div, i) {
        var card = document.querySelector('[data-division="' + div.key + '"]');
        if (card) card.style.order = i;
    });
}
```

3. Tambahkan atribut `data-division="compliance"` (dst.) pada setiap `.col-lg-4` di grid agar bisa di-reorder.

---

## Peserta 7

### Soal 13 — KYC Resubmission: Lightbox Preview Gambar KTP

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Tambahkan pop-up modal *lightbox* yang menampilkan gambar KTP dalam ukuran penuh saat user mengklik thumbnail "KTP sebelumnya" (`#oldKtpPreview`).

**Konteks kode saat ini:**
Thumbnail KTP lama ditampilkan sebagai `<img id="oldKtpPreview" class="old-ktp-preview">` (baris 165). Saat ini tidak ada interaksi klik apapun pada gambar tersebut.

**Langkah pengerjaan:**

1. Tambahkan markup lightbox modal (di luar form, sebelum `</body>`):

```html
<div id="lightboxOverlay" onclick="closeLightbox()"
     style="display:none; position:fixed; inset:0; background:rgba(0,0,0,0.85);
            z-index:99999; cursor:zoom-out; align-items:center; justify-content:center;">
    <div onclick="event.stopPropagation()" style="position:relative; max-width:90vw; max-height:90vh;">
        <img id="lightboxImg" src="" alt="KTP Preview"
             style="max-width:100%; max-height:85vh; border-radius:12px; box-shadow:0 20px 60px rgba(0,0,0,0.5);">
        <button onclick="closeLightbox()"
                style="position:absolute; top:-14px; right:-14px; background:#ef4444; border:none;
                       color:#fff; width:32px; height:32px; border-radius:50%; cursor:pointer; font-size:1rem;">
            <i class="fa-solid fa-xmark"></i>
        </button>
    </div>
</div>
```

2. Tambahkan fungsi dan event klik:

```js
function openLightbox(src) {
    var overlay = document.getElementById('lightboxOverlay');
    document.getElementById('lightboxImg').src = src;
    overlay.style.display = 'flex';
    document.body.style.overflow = 'hidden';
}

function closeLightbox() {
    document.getElementById('lightboxOverlay').style.display = 'none';
    document.body.style.overflow = '';
}

// Tambahkan klik listener setelah loadKYCData() berhasil mengisi #oldKtpPreview
// Di dalam loadKYCData(), setelah baris: $('oldKtpContainer').style.display = 'block';
$('oldKtpPreview').style.cursor = 'zoom-in';
$('oldKtpPreview').onclick = function() { openLightbox(this.src); };
```

3. Tambahkan keyboard shortcut Escape untuk menutup:

```js
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') closeLightbox();
});
```

---

### Soal 14 — Superadmin Manajemen Akun: Sortable Column Header

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Tambahkan kemampuan sort pada header kolom tabel: klik header **"Nama"** untuk mengurutkan A→Z / Z→A, dengan ikon panah (▲/▼) yang berubah sesuai arah sort aktif.

**Konteks kode saat ini:**
Header tabel ada di baris 297–304. Data dimuat dari API via `loadUsers()` yang mengirim parameter ke `GET /api/admin/users`. Sort saat ini tidak dikirim ke API.

**Langkah pengerjaan:**

1. Ubah `<th>` untuk kolom Nama agar bisa diklik:

```html
<th style="padding:16px 20px; cursor:pointer; user-select:none;" onclick="toggleSort('full_name')" id="thName">
    Nama <i class="fa-solid fa-sort ms-1 text-muted" id="sortIcon_full_name" style="font-size:0.75rem;"></i>
</th>
```

2. Tambahkan state dan logika sort:

```js
var sortField = '';
var sortDir   = 'asc';

function toggleSort(field) {
    if (sortField === field) {
        sortDir = sortDir === 'asc' ? 'desc' : 'asc';
    } else {
        sortField = field;
        sortDir = 'asc';
    }
    updateSortIcons();
    loadUsers(1);
}

function updateSortIcons() {
    // Reset semua ikon
    document.querySelectorAll('[id^="sortIcon_"]').forEach(function(el) {
        el.className = 'fa-solid fa-sort ms-1 text-muted';
        el.style.fontSize = '0.75rem';
    });
    // Set ikon aktif
    var activeIcon = document.getElementById('sortIcon_' + sortField);
    if (activeIcon) {
        activeIcon.className = sortDir === 'asc'
            ? 'fa-solid fa-sort-up ms-1'
            : 'fa-solid fa-sort-down ms-1';
        activeIcon.style.color = 'var(--accent-cyan)';
    }
}
```

3. Tambahkan parameter sort ke URLSearchParams di `loadUsers()`:

```js
// Di dalam loadUsers(), setelah baris: if (status) params.set('status', status);
if (sortField) {
    params.set('sort_by',  sortField);
    params.set('sort_dir', sortDir);
}
```

---

## Peserta 8

### Soal 15 — Register: Modal Konfirmasi Data Sebelum Submit

**File:** `frontend/account/register.html`

**Apa yang harus dilakukan:**
Sebelum form register disubmit, tampilkan modal konfirmasi yang merangkum data yang diisi (nama, email, nomor telepon) dengan tombol **"Konfirmasi Daftar"** dan **"Ubah Data"**.

**Langkah pengerjaan:**

1. Tambahkan modal di dalam `<body>`, sebelum `</div>` penutup container:

```html
<div class="modal fade" id="confirmRegisterModal" tabindex="-1" data-bs-backdrop="static">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content" style="background:#fff; border-radius:16px;">
            <div class="modal-body p-4">
                <h5 class="fw-bold mb-1" style="color:var(--text-heading);">
                    <i class="fa-solid fa-clipboard-check me-2" style="color:#00b4d8;"></i>Konfirmasi Data
                </h5>
                <p class="small text-muted mb-3">Pastikan data Anda sudah benar sebelum mendaftar.</p>
                <table class="table table-sm table-borderless mb-3">
                    <tr><td class="text-muted small fw-semibold" style="width:40%">Nama</td>
                        <td class="fw-bold" id="confirmName"></td></tr>
                    <tr><td class="text-muted small fw-semibold">Email</td>
                        <td id="confirmEmail"></td></tr>
                    <tr><td class="text-muted small fw-semibold">No. Telepon</td>
                        <td id="confirmPhone"></td></tr>
                    <tr><td class="text-muted small fw-semibold">Tanggal Lahir</td>
                        <td id="confirmBirthdate"></td></tr>
                </table>
                <div class="d-flex gap-2">
                    <button class="btn btn-outline-secondary flex-fill" data-bs-dismiss="modal">
                        <i class="fa-solid fa-pen me-1"></i>Ubah Data
                    </button>
                    <button class="btn btn-primary-glow flex-fill" id="btnFinalSubmit">
                        <i class="fa-solid fa-check me-1"></i>Konfirmasi Daftar
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>
```

2. Ubah event submit form agar menampilkan modal dulu, bukan langsung submit:

```js
document.getElementById('registerForm').addEventListener('submit', function(e) {
    e.preventDefault();
    if (!isFormValid()) return; // validasi tetap jalan

    // Isi data ke modal
    document.getElementById('confirmName').textContent      = document.getElementById('fullName').value;
    document.getElementById('confirmEmail').textContent     = document.getElementById('email').value;
    document.getElementById('confirmPhone').textContent     = document.getElementById('phone').value;
    document.getElementById('confirmBirthdate').textContent = document.getElementById('birthdate').value;

    new bootstrap.Modal(document.getElementById('confirmRegisterModal')).show();
});

// Tombol "Konfirmasi Daftar" yang memanggil fetch API
document.getElementById('btnFinalSubmit').addEventListener('click', function() {
    bootstrap.Modal.getInstance(document.getElementById('confirmRegisterModal')).hide();
    submitRegister(); // pindahkan logika fetch ke fungsi terpisah
});
```

---

### Soal 16 — Superadmin Kinerja Divisi: Sort Berdasarkan Persentase Penyelesaian

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan tombol sort untuk mengurutkan kartu divisi berdasarkan **persentase penyelesaian** (resolved ÷ total × 100%), dari tertinggi ke terendah.

**Perbedaan dari Soal 12:**
Soal 12 sort berdasarkan nilai absolut pending. Soal ini menggunakan rasio penyelesaian — divisi dengan total task sedikit tapi semua selesai akan muncul lebih atas dibandingkan divisi besar yang banyak pending-nya.

**Langkah pengerjaan:**

1. Tambahkan tombol sort kedua (atau dropdown sort):

```html
<select id="sortSelect" class="form-select form-select-sm" style="width:auto;"
        onchange="applySortDivisions(this.value)">
    <option value="">Urutan Default</option>
    <option value="pending_desc">Pending Terbanyak</option>
    <option value="completion_desc">Penyelesaian Tertinggi (%)</option>
    <option value="completion_asc">Penyelesaian Terendah (%)</option>
</select>
```

2. Fungsi sort dengan perhitungan persentase:

```js
function getCompletionRate(div) {
    var resolved = div.resolved || div.done || 0;
    var pending  = div.pending  || div.open  || 0;
    var total    = resolved + pending;
    return total > 0 ? (resolved / total) * 100 : 0;
}

function applySortDivisions(mode) {
    var data = divisionData.slice();
    if (mode === 'pending_desc') {
        data.sort(function(a, b) { return (b.pending || 0) - (a.pending || 0); });
    } else if (mode === 'completion_desc') {
        data.sort(function(a, b) { return getCompletionRate(b) - getCompletionRate(a); });
    } else if (mode === 'completion_asc') {
        data.sort(function(a, b) { return getCompletionRate(a) - getCompletionRate(b); });
    }
    // Re-order kartu berdasarkan `data`
    var row = document.getElementById('cardsRow');
    data.forEach(function(div, i) {
        var card = row.querySelector('[data-division="' + div.key + '"]');
        if (card) card.style.order = i;
    });
}
```

3. Tambahkan atribut `style="display:flex; flex-wrap:wrap;"` pada `#cardsRow` agar CSS `order` bekerja.

---

## Peserta 9

### Soal 17 — Superadmin Manajemen Akun: Modal Quick View Detail

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Tambahkan tombol **"Quick View"** di setiap baris tabel yang menampilkan detail lengkap akun dalam pop-up modal, tanpa perlu berpindah ke halaman `/ops/user-detail`.

**Konteks kode saat ini:**
Tombol "Detail" saat ini navigate ke `/ops/user-detail?id=...` (baris 374). Soal ini menambahkan tombol baru di samping tombol tersebut yang membuka modal inline.

**Langkah pengerjaan:**

1. Tambahkan markup modal detail (setelah `#confirmModal`):

```html
<div class="modal-overlay" id="detailModal">
    <div class="modal-box" style="max-width:520px; text-align:left;">
        <div class="d-flex justify-content-between align-items-start mb-3">
            <h5 class="modal-title mb-0">
                <i class="fa-solid fa-id-card me-2" style="color:var(--accent-cyan);"></i>Detail Akun
            </h5>
            <button onclick="closeDetailModal()"
                    style="background:none; border:none; color:#94a3b8; font-size:1.2rem; cursor:pointer;">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
        <div id="detailModalContent"></div>
    </div>
</div>
```

2. Tambahkan tombol Quick View di fungsi `renderActionButtons()`:

```js
// Tambahkan setelah var detailBtn = ...
var quickBtn = '<button class="btn btn-sm btn-outline-secondary" ' +
    'onclick=\'showQuickView(' + JSON.stringify(u) + ')\'>' +
    '<i class="fa-solid fa-bolt me-1"></i>Quick</button>';
// Tambahkan quickBtn ke dalam return string
```

3. Fungsi menampilkan modal dengan data user:

```js
function showQuickView(u) {
    var html =
        '<div class="mb-2"><small class="text-muted">Nama</small>' +
        '<div class="fw-bold" style="color:#e2e8f0;">' + escapeHtml(u.full_name) + '</div></div>' +
        '<div class="mb-2"><small class="text-muted">Email</small>' +
        '<div>' + escapeHtml(u.email) + '</div></div>' +
        '<div class="mb-2"><small class="text-muted">Role</small>' +
        '<div><span class="role-badge">' + escapeHtml(u.role || '-') + '</span></div></div>' +
        '<div class="mb-2"><small class="text-muted">Status</small>' +
        '<div>' + renderStatus(u.status) + '</div></div>' +
        '<div class="mb-3"><small class="text-muted">Dibuat</small>' +
        '<div class="small">' + formatDate(u.created_at) + '</div></div>' +
        '<a href="/ops/user-detail?id=' + u.user_id + '" class="btn btn-sm w-100" ' +
        'style="background:var(--accent-cyan); color:#000; font-weight:600;">' +
        '<i class="fa-solid fa-arrow-up-right-from-square me-1"></i>Buka Halaman Detail</a>';

    document.getElementById('detailModalContent').innerHTML = html;
    document.getElementById('detailModal').classList.add('show');
}

function closeDetailModal() {
    document.getElementById('detailModal').classList.remove('show');
}
```

---

### Soal 18 — KYC Resubmission: Modal Ringkasan Perubahan Data

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Sebelum tombol "Kirim Ulang" memproses request, tampilkan modal yang membandingkan **nilai lama vs nilai baru** untuk setiap field yang berubah.

**Langkah pengerjaan:**

1. Tambahkan modal sebelum `</body>`:

```html
<div id="changesSummaryOverlay" onclick="closeChangesSummary()"
     style="display:none; position:fixed; inset:0; background:rgba(0,0,0,0.7);
            z-index:99998; align-items:center; justify-content:center;">
    <div onclick="event.stopPropagation()"
         style="background:var(--bg-card,#1e2030); border:1px solid var(--border-color);
                border-radius:16px; padding:2rem; max-width:500px; width:90%; max-height:80vh; overflow-y:auto;">
        <h5 class="fw-bold mb-1" style="color:var(--text-heading);">
            <i class="fa-solid fa-arrows-rotate me-2 text-cyan"></i>Ringkasan Perubahan
        </h5>
        <p class="text-muted small mb-3">Pastikan perubahan data berikut sudah benar.</p>
        <div id="changesList" class="mb-3"></div>
        <div class="d-flex gap-2">
            <button class="btn btn-outline-secondary flex-fill" onclick="closeChangesSummary()">
                <i class="fa-solid fa-pen me-1"></i>Kembali Edit
            </button>
            <button class="btn btn-submit flex-fill" id="btnConfirmResubmit">
                <i class="fa-solid fa-paper-plane me-1"></i>Kirim Ulang
            </button>
        </div>
    </div>
</div>
```

2. Bangun fungsi pembanding dan tampilkan sebelum submit:

```js
var FIELD_LABELS = {
    full_name: 'Nama Lengkap', nik: 'NIK',
    birthdate: 'Tanggal Lahir', phone: 'No. Telepon', address: 'Alamat'
};

function buildChangesList() {
    var fields = { full_name: 'fullName', nik: 'nik', birthdate: 'birthdate', phone: 'phone', address: 'address' };
    var rows = '';
    for (var key in fields) {
        var oldVal = (previousKYC && previousKYC[key]) || '—';
        var newVal = document.getElementById(fields[key]).value.trim() || '—';
        var changed = oldVal !== newVal;
        rows += '<div class="mb-2 p-2 rounded" style="background:rgba(255,255,255,0.03);">' +
            '<div class="small fw-semibold mb-1" style="color:' + (changed ? '#00d4ff' : '#94a3b8') + ';">' +
            FIELD_LABELS[key] + (changed ? ' <span class="badge bg-warning text-dark ms-1" style="font-size:0.6rem;">Berubah</span>' : '') +
            '</div>' +
            '<div class="d-flex gap-2 small">' +
            '<span class="text-muted text-decoration-line-through">' + escapeHtml(oldVal) + '</span>' +
            (changed ? '<i class="fa-solid fa-arrow-right text-muted"></i><span style="color:#22c55e;">' + escapeHtml(newVal) + '</span>' : '') +
            '</div></div>';
    }
    return rows;
}

// Di event submit form, ganti langsung proses dengan tampilkan modal dulu:
$('resubmitForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    if (!validateForm()) return;
    // Tampilkan modal ringkasan
    document.getElementById('changesList').innerHTML = buildChangesList();
    document.getElementById('changesSummaryOverlay').style.display = 'flex';
});

function closeChangesSummary() {
    document.getElementById('changesSummaryOverlay').style.display = 'none';
}

// Pindahkan logika fetch ke fungsi terpisah, panggil dari tombol konfirmasi
document.getElementById('btnConfirmResubmit').addEventListener('click', function() {
    closeChangesSummary();
    doResubmit(); // fungsi yang berisi logika fetch PUT /api/kyc/resubmit
});
```

---

## Peserta 10

### Soal 19 — Login: Modal Lockout Setelah 3 Kali Gagal

**File:** `frontend/account/login.html`

**Apa yang harus dilakukan:**
Tampilkan modal pop-up **"Terlalu Banyak Percobaan"** apabila pengguna gagal login 3 kali berturut-turut, dengan countdown timer 30 detik sebelum bisa mencoba lagi.

**Langkah pengerjaan:**

1. Tambahkan markup modal sebelum `</body>`:

```html
<div class="modal fade" id="lockoutModal" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered modal-sm">
        <div class="modal-content" style="background:#fff; border-radius:16px; text-align:center;">
            <div class="modal-body p-4">
                <div style="width:64px; height:64px; border-radius:50%; background:rgba(220,38,38,0.1);
                            margin:0 auto 1rem; display:flex; align-items:center; justify-content:center;">
                    <i class="fa-solid fa-lock fa-xl" style="color:#dc2626;"></i>
                </div>
                <h5 class="fw-bold" style="color:var(--text-heading);">Terlalu Banyak Percobaan</h5>
                <p class="small text-muted mb-2">Silakan tunggu sebelum mencoba lagi.</p>
                <div style="font-size:2rem; font-weight:800; color:#dc2626;" id="lockCountdown">30</div>
                <p class="small text-muted">detik</p>
            </div>
        </div>
    </div>
</div>
```

2. Tambahkan logika counter gagal login:

```js
var failedAttempts = 0;
var lockoutTimer   = null;
var lockModal      = null;

function startLockout() {
    lockModal = new bootstrap.Modal(document.getElementById('lockoutModal'));
    lockModal.show();
    document.getElementById('btnSubmit').disabled = true;

    var seconds = 30;
    document.getElementById('lockCountdown').textContent = seconds;

    lockoutTimer = setInterval(function() {
        seconds--;
        document.getElementById('lockCountdown').textContent = seconds;
        if (seconds <= 0) {
            clearInterval(lockoutTimer);
            lockModal.hide();
            failedAttempts = 0;
            document.getElementById('btnSubmit').disabled = false;
        }
    }, 1000);
}

// Di dalam event handler submit, setelah login gagal:
// else {
//     failedAttempts++;
//     if (failedAttempts >= 3) { startLockout(); return; }
//     showToast(data.msg || 'Email atau Password salah');
// }
```

3. Integrasikan ke handler `loginForm` yang sudah ada (baris 145–191) dengan menambah pengecekan `if (failedAttempts >= 3)` dan `failedAttempts++` pada blok `else`.

---

### Soal 20 — Superadmin Kinerja Divisi: Search Bar Filter Nama Divisi

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan search bar di atas grid kartu untuk memfilter divisi berdasarkan nama secara real-time — kartu yang tidak cocok tersembunyi, kartu yang cocok tetap tampil.

**Langkah pengerjaan:**

1. Tambahkan search bar di area header (di samping tombol Refresh):

```html
<div class="input-group input-group-sm" style="width:220px;">
    <span class="input-group-text" style="background:var(--bg-card); border-color:var(--border-color);">
        <i class="fa-solid fa-magnifying-glass text-muted"></i>
    </span>
    <input type="text" id="divisionSearch" class="form-control"
           placeholder="Cari divisi..."
           style="background:var(--bg-card); border-color:var(--border-color); color:var(--text-heading);"
           oninput="filterDivisionCards(this.value)">
</div>
```

2. Tambahkan atribut `data-name` pada setiap wrapper kartu divisi:

```html
<!-- Contoh untuk Compliance -->
<div class="col-lg-4 col-md-6" data-name="compliance">
<!-- Contoh untuk Customer Support -->
<div class="col-lg-4 col-md-6" data-name="customer support">
<!-- Contoh untuk Operasional -->
<div class="col-lg-4 col-md-6" data-name="operasional">
```

3. Fungsi filter — sembunyikan kartu yang tidak cocok dan highlight teks yang sesuai:

```js
function filterDivisionCards(query) {
    var q = query.trim().toLowerCase();
    document.querySelectorAll('#cardsRow [data-name]').forEach(function(col) {
        var name = col.getAttribute('data-name');
        var match = !q || name.includes(q);
        col.style.display = match ? '' : 'none';

        // Highlight nama divisi yang cocok
        var titleEl = col.querySelector('.division-title');
        if (titleEl) {
            var original = titleEl.getAttribute('data-original') || titleEl.textContent;
            titleEl.setAttribute('data-original', original);
            if (q && original.toLowerCase().includes(q)) {
                var regex = new RegExp('(' + q.replace(/[.*+?^${}()|[\]\\]/g,'\\$&') + ')', 'gi');
                titleEl.innerHTML = original.replace(regex,
                    '<mark style="background:#fbbf24;color:#000;border-radius:3px;padding:0 2px;">$1</mark>');
            } else {
                titleEl.textContent = original;
            }
        }
    });

    // Tampilkan pesan jika tidak ada yang cocok
    var visible = document.querySelectorAll('#cardsRow [data-name]:not([style*="none"])').length;
    var emptyEl = document.getElementById('noResultsDivision');
    if (emptyEl) emptyEl.style.display = visible === 0 ? 'block' : 'none';
}
```

4. Tambahkan elemen "tidak ada hasil" setelah `#cardsRow`:

```html
<div id="noResultsDivision" style="display:none;" class="text-center text-muted py-5">
    <i class="fa-solid fa-magnifying-glass d-block mb-2" style="font-size:2rem;"></i>
    Tidak ada divisi yang cocok dengan "<span id="noResultsQuery"></span>"
</div>
```

---

## Referensi Cepat — File yang Relevan

| Fitur | File Utama |
|---|---|
| KYC Resubmission | `frontend/client/kyc-resubmit.html` |
| Login | `frontend/account/login.html` |
| Register | `frontend/account/register.html` |
| Logout | `frontend/assets/js/ops-layout.js` / `client-layout.js` |
| Manajemen Akun Internal | `frontend/ops/manage-users.html`, `create-user.html`, `edit-user.html` |
| Kinerja Divisi (Overview) | `frontend/management/dashboard-overview.html` |
| Kinerja Divisi (Detail) | `frontend/management/dashboard-compliance.html`, `dashboard-support.html`, `dashboard-operational.html` |

## Pola Umum di Codebase Ini

- **Toast notifikasi:** fungsi `showToast(type, message)` — `type` berisi `'success'` atau `'error'`
- **Modal konfirmasi:** gunakan `.modal-overlay` + `.modal-box` yang sudah ada di `manage-users.html`
- **Fetch API:** selalu pakai `credentials: 'include'`, handle status 401 dengan redirect ke `/account/login`
- **Dark theme variables:** `--bg-card`, `--border-color`, `--text-heading`, `--text-muted`, `--accent-cyan`
- **Icons:** Font Awesome 6 (`fa-solid`, `fa-regular`) — cukup tambahkan class pada `<i>`
