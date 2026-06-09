# Panduan Persiapan Soal Latihan Sprint 3

Dokumen ini memandu pengerjaan 10 soal latihan berdasarkan 4 fitur utama:
**KYC Resubmission · Login/Logout/Register · Superadmin Manajemen Akun Internal · Superadmin Memantau Kinerja Divisi**

---

<!-- TOC START -->
## 📑 Daftar Isi

- [Peserta 1](#peserta-1)
  - [Soal 1 — KYC Resubmission: Tampilkan Alasan Penolakan sebagai Alert/Banner](#soal-1)
  - [Soal 2 — Login: Tambahkan Tombol Show/Hide Password](#soal-2)
- [Peserta 2](#peserta-2)
  - [Soal 3 — Register: Indikator Kekuatan Password (Real-time)](#soal-3)
  - [Soal 4 — Superadmin Manajemen Akun: Konfirmasi Dialog Sebelum Hapus](#soal-4)
- [Peserta 3](#peserta-3)
  - [Soal 5 — Superadmin Kinerja Divisi: Filter Date Range](#soal-5)
  - [Soal 6 — KYC Resubmission: Preview Dokumen Sebelum Submit](#soal-6)
- [Peserta 4](#peserta-4)
  - [Soal 7 — Logout: Konfirmasi Dialog Sebelum Keluar](#soal-7)
  - [Soal 8 — Superadmin Manajemen Akun: Kolom "Terakhir Aktif" + Sortable](#soal-8)
- [Peserta 5](#peserta-5)
  - [Soal 9 — Register: Dropdown Country Code di Field Telepon](#soal-9)
  - [Soal 10 — Superadmin Kinerja Divisi: Card Summary di Atas Halaman](#soal-10)
- [Peserta 6](#peserta-6)
  - [Soal 11 — Superadmin Manajemen Akun: Tampilkan Jumlah Hasil Filter](#soal-11)
  - [Soal 12 — Superadmin Kinerja Divisi: Sort Kartu Berdasarkan Task Pending](#soal-12)
- [Peserta 7](#peserta-7)
  - [Soal 13 — KYC Resubmission: Lightbox Preview Gambar KTP](#soal-13)
  - [Soal 14 — Superadmin Manajemen Akun: Sortable Column Header](#soal-14)
- [Peserta 8](#peserta-8)
  - [Soal 15 — Register: Modal Konfirmasi Data Sebelum Submit](#soal-15)
  - [Soal 16 — Superadmin Kinerja Divisi: Sort Berdasarkan Persentase Penyelesaian](#soal-16)
- [Peserta 9](#peserta-9)
  - [Soal 17 — Superadmin Manajemen Akun: Modal Quick View Detail](#soal-17)
  - [Soal 18 — KYC Resubmission: Modal Ringkasan Perubahan Data](#soal-18)
- [Peserta 10](#peserta-10)
  - [Soal 19 — Login: Modal Lockout Setelah 3 Kali Gagal](#soal-19)
  - [Soal 20 — Superadmin Kinerja Divisi: Search Bar Filter Nama Divisi](#soal-20)
- [Peserta 11](#peserta-11)
  - [Soal 21 — [CREATE] KYC Resubmission: Submit Form ke API](#soal-21)
  - [Soal 22 — [READ] KYC Resubmission: Tampilkan Status & Riwayat KYC](#soal-22)
- [Peserta 12](#peserta-12)
  - [Soal 23 — [UPDATE] KYC Resubmission: Pre-fill Form dengan Data KYC Lama](#soal-23)
  - [Soal 24 — [DELETE] KYC Resubmission: Batalkan Pengajuan Pending](#soal-24)
    - [Panduan Backend — Soal 24](#backend-soal-24)
- [Peserta 13](#peserta-13)
  - [Soal 25 — [CREATE] Register: Submit Form dengan Validasi Lengkap](#soal-25)
  - [Soal 26 — [READ] Login: Tampilkan Sambutan Nama User Setelah Berhasil Login](#soal-26)
- [Peserta 14](#peserta-14)
  - [Soal 27 — [UPDATE] Register: Form Input OTP 6-Digit dengan Auto-focus](#soal-27)
  - [Soal 28 — [DELETE] Login/Register: Logout Paksa Semua Sesi Aktif](#soal-28)
    - [Panduan Backend — Soal 28](#backend-soal-28)
- [Peserta 15](#peserta-15)
  - [Soal 29 — [CREATE] Superadmin Manajemen Akun: Form Tambah Akun Internal](#soal-29)
  - [Soal 30 — [READ] Superadmin Manajemen Akun: Export Tabel ke CSV](#soal-30)
- [Peserta 16](#peserta-16)
  - [Soal 31 — [UPDATE] Superadmin Manajemen Akun: Toggle Status Aktif/Nonaktif Inline](#soal-31)
  - [Soal 32 — [DELETE] Superadmin Manajemen Akun: Hapus Massal dengan Checkbox](#soal-32)
    - [Panduan Backend — Soal 32](#backend-soal-32)
- [Peserta 17](#peserta-17)
  - [Soal 33 — [CREATE] Superadmin Kinerja Divisi: Tambah Catatan/Note pada Divisi](#soal-33)
  - [Soal 34 — [READ] Superadmin Kinerja Divisi: Chart Bar Perbandingan Antar Divisi](#soal-34)
- [Peserta 18](#peserta-18)
  - [Soal 35 — [UPDATE] Superadmin Kinerja Divisi: Edit Target Task Divisi Inline](#soal-35)
  - [Soal 36 — [DELETE] Superadmin Kinerja Divisi: Reset Data Kinerja Divisi](#soal-36)
    - [Panduan Backend — Soal 36](#backend-soal-36)
- [Soal Tambahan — Latihan Frontend Mandiri (Soal 37–58)](#soal-tambahan-latihan-frontend-mandiri-soal-37-58)
  - [Soal 37 — Login: Peringatan Caps Lock di Field Password](#soal-37)
  - [Soal 38 — Login: Ingat Email Terakhir (localStorage)](#soal-38)
  - [Soal 39 — Register: Tampilkan Umur Otomatis dari Tanggal Lahir](#soal-39)
  - [Soal 40 — Register: Batasi Field Telepon Hanya Angka](#soal-40)
  - [Soal 41 — Verifikasi OTP: Samarkan (Mask) Alamat Email](#soal-41)
  - [Soal 42 — Manajemen Akun: Shortcut `/` Fokus Pencarian + `Esc` Bersihkan](#soal-42)
  - [Soal 43 — Manajemen Akun: Tombol "Salin Email" per Baris](#soal-43)
  - [Soal 44 — Manajemen Akun: Avatar Inisial di Kolom Nama](#soal-44)
  - [Soal 45 — Kinerja Divisi: Auto-Refresh dengan Toggle](#soal-45)
  - [Soal 46 — Kinerja Divisi: Kartu "Total Pending Semua Divisi"](#soal-46)
  - [Soal 47 — KYC Resubmission: Penghitung Karakter Field Alamat](#soal-47)
  - [Soal 48 — Pola: Menambah Tombol Beserta Fungsinya (Studi Kasus: "Salin Email")](#soal-48)
  - [Soal 49 — Pola: Menambah Field Baru di Halaman Read/Detail (Studi Kasus: "Terakhir Diperbarui")](#soal-49)
  - [Soal 50 — Pola: Menambah Pagination (Client-Side) pada List yang Belum Punya](#soal-50)
- [Tambah Field & Tombol per Fitur (Soal 51–58)](#tambah-field-tombol-per-fitur-soal-51-58)
  - [🔹 Fitur: KYC Resubmission](#fitur-kyc-resubmission)
    - [Soal 51 — Tombol "Perbarui Status" (tanpa reload)](#soal-51)
    - [Soal 52 — Field "Nomor Pengajuan" di View Read](#soal-52)
  - [🔹 Fitur: Login / Logout / Register](#fitur-login-logout-register)
    - [Soal 53 — Field "Konfirmasi Email" + Validasi Saat Submit](#soal-53)
    - [Soal 54 — Tombol "Bersihkan Formulir" (Reset)](#soal-54)
  - [🔹 Fitur: Superadmin Manajemen Akun Internal (CRUD)](#fitur-superadmin-manajemen-akun-internal-crud)
    - [Soal 55 — Field "Jumlah per Halaman" (Page Size)](#soal-55)
    - [Soal 56 — Tombol "Buat Password Kuat" (Generate)](#soal-56)
  - [🔹 Fitur: Superadmin Memantau Kinerja Divisi](#fitur-superadmin-memantau-kinerja-divisi)
    - [Soal 57 — Field Metrik "KYC Ditolak" di Kartu Compliance](#soal-57)
    - [Soal 58 — Tombol "Salin Ringkasan" ke Clipboard](#soal-58)
- [Referensi Cepat — File yang Relevan](#referensi-cepat-file-yang-relevan)
- [Pola Umum di Codebase Ini](#pola-umum-di-codebase-ini)

<!-- TOC END -->

## Peserta 1 <a id="peserta-1"></a>

### Soal 1 — KYC Resubmission: Tampilkan Alasan Penolakan sebagai Alert/Banner <a id="soal-1"></a>

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

### Soal 2 — Login: Tambahkan Tombol Show/Hide Password <a id="soal-2"></a>

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

## Peserta 2 <a id="peserta-2"></a>

### Soal 3 — Register: Indikator Kekuatan Password (Real-time) <a id="soal-3"></a>

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

### Soal 4 — Superadmin Manajemen Akun: Konfirmasi Dialog Sebelum Hapus <a id="soal-4"></a>

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

## Peserta 3 <a id="peserta-3"></a>

### Soal 5 — Superadmin Kinerja Divisi: Filter Date Range <a id="soal-5"></a>

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

### Soal 6 — KYC Resubmission: Preview Dokumen Sebelum Submit <a id="soal-6"></a>

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

## Peserta 4 <a id="peserta-4"></a>

### Soal 7 — Logout: Konfirmasi Dialog Sebelum Keluar <a id="soal-7"></a>

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

### Soal 8 — Superadmin Manajemen Akun: Kolom "Terakhir Aktif" + Sortable <a id="soal-8"></a>

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

## Peserta 5 <a id="peserta-5"></a>

### Soal 9 — Register: Dropdown Country Code di Field Telepon <a id="soal-9"></a>

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

### Soal 10 — Superadmin Kinerja Divisi: Card Summary di Atas Halaman <a id="soal-10"></a>

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

## Peserta 6 <a id="peserta-6"></a>

### Soal 11 — Superadmin Manajemen Akun: Tampilkan Jumlah Hasil Filter <a id="soal-11"></a>

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

### Soal 12 — Superadmin Kinerja Divisi: Sort Kartu Berdasarkan Task Pending <a id="soal-12"></a>

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

## Peserta 7 <a id="peserta-7"></a>

### Soal 13 — KYC Resubmission: Lightbox Preview Gambar KTP <a id="soal-13"></a>

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

### Soal 14 — Superadmin Manajemen Akun: Sortable Column Header <a id="soal-14"></a>

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

## Peserta 8 <a id="peserta-8"></a>

### Soal 15 — Register: Modal Konfirmasi Data Sebelum Submit <a id="soal-15"></a>

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

### Soal 16 — Superadmin Kinerja Divisi: Sort Berdasarkan Persentase Penyelesaian <a id="soal-16"></a>

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

## Peserta 9 <a id="peserta-9"></a>

### Soal 17 — Superadmin Manajemen Akun: Modal Quick View Detail <a id="soal-17"></a>

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

### Soal 18 — KYC Resubmission: Modal Ringkasan Perubahan Data <a id="soal-18"></a>

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

## Peserta 10 <a id="peserta-10"></a>

### Soal 19 — Login: Modal Lockout Setelah 3 Kali Gagal <a id="soal-19"></a>

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

### Soal 20 — Superadmin Kinerja Divisi: Search Bar Filter Nama Divisi <a id="soal-20"></a>

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

## Peserta 11 <a id="peserta-11"></a>

### Soal 21 — [CREATE] KYC Resubmission: Submit Form ke API <a id="soal-21"></a>

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Implementasikan tombol submit form KYC resubmission yang mengirim data (termasuk file baru) ke API `PUT /api/kyc/resubmit` menggunakan `FormData` + `fetch`, dengan loading state dan toast notifikasi sukses/gagal.

**Status kode saat ini:**
Tombol submit dan form sudah ada, tapi logika `fetch` mungkin belum lengkap atau hanya partial. Pastikan semua field (nama, NIK, tanggal lahir, dll) dan file KTP ikut terkirim.

**Langkah pengerjaan:**

1. Cari event handler submit form (biasanya `$('resubmitForm').addEventListener('submit', ...)`).
2. Gunakan `FormData` untuk mengemas semua data:

```js
async function doResubmit() {
    var btn = document.getElementById('btnSubmit');
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Mengirim...';

    var form = document.getElementById('resubmitForm');
    var fd = new FormData(form);

    // Pastikan file ikut jika user mengunggah file baru
    var fileInput = document.getElementById('ktpFile');
    if (fileInput && fileInput.files[0]) {
        fd.set('ktp_file', fileInput.files[0]);
    }

    try {
        var res = await fetch('/api/kyc/resubmit', {
            method: 'PUT',
            credentials: 'include',
            body: fd
            // JANGAN set Content-Type, biarkan browser isi boundary otomatis
        });
        var data = await res.json();
        if (res.ok) {
            showToast('success', 'KYC berhasil dikirim ulang. Menunggu review...');
            setTimeout(() => window.location.href = '/client/kyc-status', 2000);
        } else {
            showToast('error', data.message || 'Gagal mengirim KYC.');
        }
    } catch (err) {
        showToast('error', 'Koneksi gagal. Coba lagi.');
    } finally {
        btn.disabled = false;
        btn.innerHTML = '<i class="fa-solid fa-paper-plane me-2"></i>Kirim Ulang';
    }
}
```

---

### Soal 22 — [READ] KYC Resubmission: Tampilkan Status & Riwayat KYC <a id="soal-22"></a>

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Fetch status KYC terbaru user dari API `GET /api/kyc/status`, tampilkan badge status (pending/approved/rejected) di bagian atas halaman beserta tanggal pengajuan terakhir.

**Langkah pengerjaan:**

1. Tambahkan elemen status di atas form:

```html
<div id="kycStatusBar" class="d-flex align-items-center gap-3 p-3 rounded mb-4"
     style="background:rgba(255,255,255,0.04); border:1px solid var(--border-color); display:none !important;">
    <div>
        <div class="small text-muted">Status KYC Terakhir</div>
        <span id="kycStatusBadge" class="badge fs-6 mt-1"></span>
    </div>
    <div class="ms-auto text-end">
        <div class="small text-muted">Diajukan pada</div>
        <div id="kycSubmitDate" class="small" style="color:var(--text-heading);"></div>
    </div>
</div>
```

2. Fetch dan isi:

```js
async function loadKYCStatus() {
    try {
        var res = await fetch('/api/kyc/status', { credentials: 'include' });
        if (!res.ok) return;
        var data = await res.json();
        var kyc = data.data;
        if (!kyc) return;

        var statusColors = {
            pending:  { bg: '#f59e0b', text: 'Menunggu Review' },
            approved: { bg: '#22c55e', text: 'Disetujui' },
            rejected: { bg: '#ef4444', text: 'Ditolak' }
        };
        var cfg = statusColors[kyc.status] || { bg: '#94a3b8', text: kyc.status };

        var badge = document.getElementById('kycStatusBadge');
        badge.textContent = cfg.text;
        badge.style.background = cfg.bg;
        badge.style.color = '#000';

        document.getElementById('kycSubmitDate').textContent =
            new Date(kyc.created_at).toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' });

        var bar = document.getElementById('kycStatusBar');
        bar.style.display = 'flex';
    } catch (e) { /* silent */ }
}

// Panggil saat halaman load
document.addEventListener('DOMContentLoaded', loadKYCStatus);
```

---

## Peserta 12 <a id="peserta-12"></a>

### Soal 23 — [UPDATE] KYC Resubmission: Pre-fill Form dengan Data KYC Lama <a id="soal-23"></a>

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Saat halaman dimuat, otomatis isi field-field form (nama, NIK, tanggal lahir, telepon, alamat) dengan data dari KYC sebelumnya yang diambil dari API, sehingga user tinggal mengubah bagian yang perlu dikoreksi.

**Langkah pengerjaan:**

1. Di dalam fungsi `loadKYCData()` yang sudah ada, setelah data API diterima, isi field form:

```js
// Di dalam loadKYCData(), setelah: var data = await res.json();
var kyc = data.data;

// Field mapping: ID element → key dari API response
var fieldMap = {
    'fullName':  kyc.full_name  || '',
    'nik':       kyc.nik        || '',
    'birthdate': kyc.birthdate  || '',
    'phone':     kyc.phone      || '',
    'address':   kyc.address    || ''
};

Object.keys(fieldMap).forEach(function(id) {
    var el = document.getElementById(id);
    if (el && fieldMap[id]) el.value = fieldMap[id];
});

// Simpan data lama untuk perbandingan di modal ringkasan (Soal 18)
window.previousKYC = kyc;
```

2. Tambahkan visual hint bahwa data sudah di-pre-fill:

```js
// Setelah mengisi field, tambahkan kelas highlight sementara
Object.keys(fieldMap).forEach(function(id) {
    var el = document.getElementById(id);
    if (el && fieldMap[id]) {
        el.style.borderColor = 'rgba(0,212,255,0.4)';
        el.addEventListener('focus', function() {
            this.style.borderColor = '';
        }, { once: true });
    }
});
```

---

### Soal 24 — [DELETE] KYC Resubmission: Batalkan Pengajuan Pending <a id="soal-24"></a>

**File:** `frontend/client/kyc-resubmit.html`

**Apa yang harus dilakukan:**
Tambahkan tombol "Batalkan Pengajuan" yang hanya muncul jika status KYC = `'pending'`. Tombol memicu modal konfirmasi, lalu mengirim `DELETE /api/kyc/cancel` ke API.

**Langkah pengerjaan:**

1. Tambahkan tombol batalkan di dekat status bar (Soal 22):

```html
<button id="btnCancelKYC" style="display:none;"
        class="btn btn-sm btn-outline-danger"
        onclick="showCancelConfirm()">
    <i class="fa-solid fa-ban me-1"></i>Batalkan Pengajuan
</button>
```

2. Tampilkan tombol hanya jika status = pending (di `loadKYCStatus()`):

```js
// Setelah mengisi badge status
if (kyc.status === 'pending') {
    document.getElementById('btnCancelKYC').style.display = 'inline-flex';
}
```

3. Fungsi konfirmasi dan DELETE:

```js
function showCancelConfirm() {
    var ok = confirm('Yakin ingin membatalkan pengajuan KYC yang sedang pending?\nTindakan ini tidak dapat dibatalkan.');
    if (!ok) return;
    cancelKYC();
}

async function cancelKYC() {
    try {
        var res = await fetch('/api/kyc/cancel', {
            method: 'DELETE',
            credentials: 'include'
        });
        var data = await res.json();
        if (res.ok) {
            showToast('success', 'Pengajuan KYC berhasil dibatalkan.');
            setTimeout(() => location.reload(), 1500);
        } else {
            showToast('error', data.message || 'Gagal membatalkan pengajuan.');
        }
    } catch (e) {
        showToast('error', 'Koneksi gagal.');
    }
}
```

---

> **⚠️ Endpoint `DELETE /api/kyc/cancel` belum ada di backend — kerjakan langkah berikut sebelum implementasi frontend.**

#### Panduan Backend — Soal 24 <a id="backend-soal-24"></a>

**Service:** `users` · **Framework:** ZaFramework (bukan Gin)

**Langkah 1 — [users/app/routes/router.go](../users/app/routes/router.go)**

Tambah 1 baris setelah route `PUT /api/kyc/resubmit`:
```go
app.Router.HandleFunc("DELETE /api/kyc/cancel", kycController.Cancel)
```

**Langkah 2 — [users/app/modules/kyc/repository.go](../users/app/modules/kyc/repository.go)**

Di `type Repository interface { ... }`, tambah:
```go
CancelByUserID(ctx context.Context, userID string) error
```

Implementasi di bawah fungsi terakhir file:
```go
func (r *kycRepo) CancelByUserID(ctx context.Context, userID string) error {
    tag, err := r.db.Pool.Exec(ctx,
        `UPDATE kyc_submissions SET status='cancelled', updated_at=NOW()
         WHERE user_id=$1 AND status='pending'`, userID)
    if err != nil {
        return err
    }
    if tag.RowsAffected() == 0 {
        return fmt.Errorf("NOT_PENDING")
    }
    return nil
}
```

**Langkah 3 — [users/app/modules/kyc/service.go](../users/app/modules/kyc/service.go)**

> ⚠️ **JANGAN tulis `case "kyc_cancel":`.** Service ini tidak punya `switch`. Pola ZaFramework di project ini adalah **1 fungsi `Process…Job` per task** lalu didaftarkan di `main.go`. Menulis `case "kyc_cancel": m := p.(…)` akan menghasilkan error **`undefined: p`** karena variabel `p` memang tidak ada.

Tambah fungsi baru. Letakkan di antara fungsi `Process…Job` lain (mis. tepat setelah `ProcessKYCStatusJob` berakhir, sekitar baris 210):
```go
// ProcessKYCCancelJob — handler untuk dispatcher job "kyc_cancel"
func (s *Service) ProcessKYCCancelJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}
	userID, _ := data["user_id"].(string)
	if userID == "" {
		return nil, fmt.Errorf("user_id wajib diisi")
	}
	return nil, s.repo.CancelByUserID(ctx, userID)
}
```
*(Import `context` dan `fmt` sudah ada di file ini — tidak perlu menambah import.)*

**Langkah 4 — [users/main.go](../users/main.go)** — daftarkan job-nya (kalau tidak, dispatch akan error "job not found")

Tambah **1 baris** di blok `RegisterJob`, tepat setelah baris `kyc_resubmit` (sekitar baris 157):
```go
app.RegisterJob("kyc_cancel", kycService.ProcessKYCCancelJob)
```

**Langkah 5 — [users/app/modules/kyc/controller.go](../users/app/modules/kyc/controller.go)**

Tambah method `Cancel` sebagai fungsi top-level di akhir file:
```go
func (c *Controller) Cancel(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid", nil)
		return
	}
	_, err := c.Dispatcher.DispatchAndWait(r.Context(), "kyc_cancel",
		map[string]any{"user_id": userID}, concurrency.PriorityHigh)
	if err != nil {
		if err.Error() == "NOT_PENDING" {
			resp.ApiJSON(w, r, http.StatusBadRequest, false, "Tidak ada KYC pending untuk dibatalkan", nil)
			return
		}
		resp.ApiJSON(w, r, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	resp.ApiJSON(w, r, http.StatusOK, true, "Pengajuan KYC berhasil dibatalkan", nil)
}
```
*(Import `concurrency`, `resp` (alias `core/http`), dan `net/http` sudah dipakai controller lain di file ini — tidak perlu menambah import.)*

**Langkah 6 — [users/core/database/migrate.go](../users/core/database/migrate.go)** ⚠️ **WAJIB.** Tanpa ini, `DELETE /api/kyc/cancel` membalas **500** karena status `'cancelled'` ditolak CHECK constraint.

a) **DB baru** — pada `CREATE TABLE … kyc_submissions`, ubah kolom `status` (sekitar baris 126) menjadi:
```sql
status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
```

b) **DB yang sudah ada** — tambah patch idempotent ini **setelah** baris `ALTER TABLE kyc_submissions ADD COLUMN IF NOT EXISTS rejected_fields TEXT[];` (sekitar baris 144). Aman dijalankan berulang kali:
```sql
-- Patch: allow 'cancelled' status for KYC cancel feature (idempotent)
DO $$
DECLARE v TEXT;
BEGIN
	SELECT conname INTO v
	FROM pg_constraint
	WHERE conrelid = 'kyc_submissions'::regclass
	AND contype = 'c'
	AND pg_get_constraintdef(oid) LIKE '%pending%'
	AND pg_get_constraintdef(oid) NOT LIKE '%cancelled%'
	LIMIT 1;
	IF v IS NOT NULL THEN
		EXECUTE 'ALTER TABLE kyc_submissions DROP CONSTRAINT ' || v;
		EXECUTE 'ALTER TABLE kyc_submissions ADD CONSTRAINT kyc_submissions_status_check CHECK (status IN (''pending'', ''approved'', ''rejected'', ''cancelled''))';
	END IF;
END$$;
```

> **Catatan frontend Soal 24:** tombol "Batalkan Pengajuan" paling tepat dipasang di **`frontend/client/kyc-status.html`** (di blok status `pending`), bukan di `kyc-resubmit.html` — karena `kyc-resubmit.html` hanya menampilkan form untuk status `rejected`. Setelah cancel sukses, arahkan ke `/client/kyc` (form pengajuan baru), bukan `location.reload()`.

---

## Peserta 13 <a id="peserta-13"></a>

### Soal 25 — [CREATE] Register: Submit Form dengan Validasi Lengkap <a id="soal-25"></a>

**File:** `frontend/account/register.html`

**Apa yang harus dilakukan:**
Implementasikan validasi lengkap sebelum submit form register: format email, panjang minimal password (8 karakter), kesesuaian konfirmasi password, dan format nomor telepon. Kirim POST ke `/api/auth/register`, lalu redirect ke halaman OTP jika berhasil.

**Langkah pengerjaan:**

1. Fungsi validasi semua field:

```js
function validateRegisterForm() {
    var errors = [];
    var name  = document.getElementById('fullName').value.trim();
    var email = document.getElementById('email').value.trim();
    var pass  = document.getElementById('password').value;
    var conf  = document.getElementById('confirmPassword').value;
    var phone = document.getElementById('phone').value.trim();

    if (name.length < 3)
        errors.push('Nama minimal 3 karakter.');
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email))
        errors.push('Format email tidak valid.');
    if (pass.length < 8)
        errors.push('Password minimal 8 karakter.');
    if (pass !== conf)
        errors.push('Konfirmasi password tidak cocok.');
    if (!/^[0-9]{8,15}$/.test(phone.replace(/^0+/, '')))
        errors.push('Nomor telepon tidak valid (8–15 digit).');

    return errors;
}
```

2. Di event submit, jalankan validasi dulu:

```js
document.getElementById('registerForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    var errors = validateRegisterForm();
    if (errors.length) {
        showToast('error', errors[0]); // tampilkan error pertama
        return;
    }

    var btn = document.getElementById('btnRegister');
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Mendaftar...';

    try {
        var res = await fetch('/api/auth/register', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                full_name:  document.getElementById('fullName').value.trim(),
                email:      document.getElementById('email').value.trim(),
                password:   document.getElementById('password').value,
                phone:      document.getElementById('phone').value.trim(),
                birthdate:  document.getElementById('birthdate').value
            })
        });
        var data = await res.json();
        if (res.ok) {
            showToast('success', 'Registrasi berhasil! Cek email untuk OTP.');
            setTimeout(() => window.location.href = '/account/verify-otp', 1500);
        } else {
            showToast('error', data.message || 'Registrasi gagal.');
        }
    } catch (err) {
        showToast('error', 'Koneksi gagal.');
    } finally {
        btn.disabled = false;
        btn.innerHTML = 'Daftar Sekarang';
    }
});
```

---

### Soal 26 — [READ] Login: Tampilkan Sambutan Nama User Setelah Berhasil Login <a id="soal-26"></a>

**File:** `frontend/account/login.html`

**Apa yang harus dilakukan:**
Setelah login berhasil (response dari `/api/auth/login`), simpan data user ke `localStorage`, lalu tampilkan toast **"Selamat datang, [Nama]!"** selama 2 detik sebelum redirect ke dashboard.

**Status kode saat ini:**
Handler login yang ada (sekitar baris 145–191) langsung redirect tanpa menampilkan sambutan personalisasi.

**Langkah pengerjaan:**

Di dalam blok `if (res.ok)` pada handler login, tambahkan:

```js
// Simpan info user ke localStorage untuk dipakai di halaman lain
if (data.data && data.data.user) {
    localStorage.setItem('user', JSON.stringify(data.data.user));
    localStorage.setItem('user_name', data.data.user.full_name || data.data.user.name || '');
}

// Tampilkan toast sambutan dengan nama user
var name = (data.data && data.data.user)
    ? (data.data.user.full_name || data.data.user.name || 'Pengguna')
    : 'Pengguna';
showToast('success', 'Selamat datang, ' + name + '! 👋');

// Redirect setelah toast tampil
var redirectUrl = data.data.redirect || '/client/dashboard';
setTimeout(function() { window.location.href = redirectUrl; }, 1800);
```

**Baca data user di halaman lain:**
```js
// Di halaman dashboard / layout
var storedUser = localStorage.getItem('user');
if (storedUser) {
    var user = JSON.parse(storedUser);
    document.getElementById('navUserName').textContent = user.full_name || user.name;
}
```

---

## Peserta 14 <a id="peserta-14"></a>

### Soal 27 — [UPDATE] Register: Form Input OTP 6-Digit dengan Auto-focus <a id="soal-27"></a>

**File:** `frontend/account/verify-otp.html`

**Apa yang harus dilakukan:**
Implementasikan form verifikasi OTP dengan 6 kotak input terpisah (satu digit per kotak), auto-focus ke kotak berikutnya saat digit diisi, backspace pindah ke kotak sebelumnya, dan tombol Resend OTP dengan countdown 60 detik.

**Langkah pengerjaan:**

1. Markup 6 input OTP:

```html
<div class="d-flex gap-2 justify-content-center my-3" id="otpInputs">
    <input type="text" maxlength="1" class="otp-box form-control text-center fw-bold fs-4"
           style="width:52px; height:60px; border-radius:10px;" data-index="0" inputmode="numeric">
    <!-- Ulangi 5x lebih untuk index 1-5 -->
</div>
<p class="text-muted small text-center mt-2">
    Tidak menerima kode?
    <button id="btnResend" class="btn btn-link p-0 small" disabled onclick="resendOTP()">
        Kirim ulang (<span id="resendCountdown">60</span>s)
    </button>
</p>
```

2. Logika auto-focus dan backspace:

```js
document.querySelectorAll('.otp-box').forEach(function(input, idx, boxes) {
    input.addEventListener('input', function(e) {
        // Hanya izinkan digit
        this.value = this.value.replace(/\D/g, '').slice(-1);
        if (this.value && idx < boxes.length - 1) {
            boxes[idx + 1].focus();
        }
        // Auto submit saat kotak terakhir diisi
        if (idx === boxes.length - 1 && this.value) {
            submitOTP();
        }
    });

    input.addEventListener('keydown', function(e) {
        if (e.key === 'Backspace' && !this.value && idx > 0) {
            boxes[idx - 1].focus();
        }
    });
});

function getOTPValue() {
    return Array.from(document.querySelectorAll('.otp-box'))
        .map(function(el) { return el.value; }).join('');
}

// Countdown Resend OTP
var resendSeconds = 60;
var resendTimer = setInterval(function() {
    resendSeconds--;
    document.getElementById('resendCountdown').textContent = resendSeconds;
    if (resendSeconds <= 0) {
        clearInterval(resendTimer);
        var btn = document.getElementById('btnResend');
        btn.disabled = false;
        btn.textContent = 'Kirim ulang sekarang';
    }
}, 1000);

async function resendOTP() {
    var res = await fetch('/api/auth/resend-otp', { method: 'POST', credentials: 'include' });
    if (res.ok) {
        showToast('success', 'Kode OTP baru telah dikirim ke email Anda.');
        resendSeconds = 60;
        document.getElementById('btnResend').disabled = true;
        document.getElementById('btnResend').innerHTML =
            'Kirim ulang (<span id="resendCountdown">60</span>s)';
        // Restart timer
    }
}
```

---

### Soal 28 — [DELETE] Login/Register: Logout Paksa Semua Sesi Aktif <a id="soal-28"></a>

**File:** `frontend/assets/js/ops-layout.js` atau `client-layout.js`

**Apa yang harus dilakukan:**
Tambahkan opsi **"Keluar dari Semua Perangkat"** di menu dropdown profil. Fitur ini mengirim `DELETE /api/auth/sessions` untuk menghapus semua sesi aktif user sekaligus, bukan hanya sesi saat ini.

**Langkah pengerjaan:**

1. Tambahkan item menu di dropdown profil (cari area sekitar `logout` di layout):

```html
<!-- Di dalam dropdown profil, setelah tombol logout biasa -->
<li><hr class="dropdown-divider" style="border-color:var(--border-color);"></li>
<li>
    <a class="dropdown-item text-danger small" href="#" onclick="logoutAllSessions(event)">
        <i class="fa-solid fa-right-from-bracket me-2"></i>Keluar dari Semua Perangkat
    </a>
</li>
```

2. Fungsi logout semua sesi:

```js
async function logoutAllSessions(e) {
    e.preventDefault();
    var ok = confirm('Ini akan menutup semua sesi aktif di semua perangkat. Lanjutkan?');
    if (!ok) return;

    try {
        await fetch('/api/auth/sessions', {
            method: 'DELETE',
            credentials: 'include'
        });
    } catch (e) { /* biarkan lanjut */ }

    // Bersihkan storage lokal
    localStorage.clear();
    sessionStorage.clear();

    showToast('success', 'Semua sesi telah diakhiri.');
    setTimeout(() => window.location.href = '/account/login', 1500);
}
```

---

> **⚠️ Endpoint `DELETE /api/auth/sessions` belum ada di backend — kerjakan langkah berikut sebelum implementasi frontend.**

#### Panduan Backend — Soal 28 <a id="backend-soal-28"></a>

**Service:** `users` · **Framework:** ZaFramework

Sesi disimpan di Redis dengan dua key: UUID acak (dari cookie) + `session:<userID>`. Implementasi ini menghapus `session:<userID>` sehingga semua token yang bergantung pada key tersebut tidak bisa divalidasi ulang.

**Langkah 1 — [users/app/routes/router.go](../users/app/routes/router.go)**

Tambah 1 baris setelah route `POST /api/auth/logout`:
```go
app.Router.HandleFunc("DELETE /api/auth/sessions", loginController.LogoutAllSessions)
```

**Langkah 2 — [users/app/modules/login/logout.go](../users/app/modules/login/logout.go)**

Tambah fungsi baru di akhir file:
```go
func (c *Controller) LogoutAllSessions(w http.ResponseWriter, r *http.Request) {
    userID := r.Header.Get("X-User-ID")
    if userID == "" {
        resp.ApiJSON(w, r, http.StatusUnauthorized, false, "Sesi tidak valid", nil)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Hapus session key user dari Redis
    sessionKey := "session:" + userID
    if err := utils.RedisDel(ctx, sessionKey); err != nil {
        log.Printf("[LOGOUT-ALL] Redis delete session:%s error: %v", userID, err)
    }

    // Clear cookie sesi saat ini
    cookieNames := []string{"token", "_authz", "session_id",
        utils.GetEnv("SESSION_NAME", "za_session")}
    for _, name := range cookieNames {
        http.SetCookie(w, &http.Cookie{
            Name: name, Value: "", Path: "/", MaxAge: -1, HttpOnly: true,
        })
    }

    resp.ApiJSON(w, r, http.StatusOK, true, "Semua sesi berhasil diakhiri", nil)
}
```
*(Tidak perlu menambah import: `context`, `time`, `log`, `net/http`, `utils`, dan `resp` semuanya sudah di-import di `logout.go` karena dipakai fungsi `Logout` di atasnya. `utils.RedisDel` juga sudah ada — dipakai di `Logout`.)*

> ✅ **Endpoint ini stateless** — tidak perlu `RegisterJob`/dispatcher karena `LogoutAllSessions` langsung memanggil Redis & set cookie, bukan lewat worker.

---

## Peserta 15 <a id="peserta-15"></a>

### Soal 29 — [CREATE] Superadmin Manajemen Akun: Form Tambah Akun Internal <a id="soal-29"></a>

**File:** `frontend/ops/create-user.html`

**Apa yang harus dilakukan:**
Implementasikan form penuh untuk menambah akun internal baru: validasi semua field, dropdown pilihan role (Admin/Compliance/Support/Operasional), dan kirim `POST /api/admin/users` ke API.

**Langkah pengerjaan:**

1. Pastikan form memiliki field: `full_name`, `email`, `password`, `role`, `phone` (opsional).

2. Dropdown role:
```html
<select class="form-control" id="role" required>
    <option value="" disabled selected>Pilih Role</option>
    <option value="admin">Admin</option>
    <option value="compliance">Compliance</option>
    <option value="support">Customer Support</option>
    <option value="operational">Operasional</option>
</select>
```

3. Submit handler:

```js
document.getElementById('createUserForm').addEventListener('submit', async function(e) {
    e.preventDefault();

    var payload = {
        full_name: document.getElementById('fullName').value.trim(),
        email:     document.getElementById('email').value.trim(),
        password:  document.getElementById('password').value,
        role:      document.getElementById('role').value,
        phone:     document.getElementById('phone') ? document.getElementById('phone').value.trim() : ''
    };

    // Validasi dasar
    if (!payload.full_name || !payload.email || !payload.password || !payload.role) {
        showToast('error', 'Semua field wajib diisi.');
        return;
    }
    if (payload.password.length < 8) {
        showToast('error', 'Password minimal 8 karakter.');
        return;
    }

    var btn = document.getElementById('btnCreate');
    btn.disabled = true;

    try {
        var res = await fetch('/api/admin/users', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        var data = await res.json();
        if (res.ok) {
            showToast('success', 'Akun internal berhasil dibuat.');
            setTimeout(() => window.location.href = '/ops/manage-users', 1500);
        } else {
            showToast('error', data.message || 'Gagal membuat akun.');
        }
    } catch (err) {
        showToast('error', 'Koneksi gagal.');
    } finally {
        btn.disabled = false;
    }
});
```

---

### Soal 30 — [READ] Superadmin Manajemen Akun: Export Tabel ke CSV <a id="soal-30"></a>

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Tambahkan tombol **"Export CSV"** yang mengunduh data akun internal yang sedang ditampilkan (sesuai filter aktif) sebagai file `.csv`, tanpa membutuhkan backend — generate langsung dari data array di browser.

**Konteks kode saat ini:**
Data tabel tersimpan di array `usersData` setelah `loadUsers()` dipanggil. Gunakan array ini sebagai sumber data.

**Langkah pengerjaan:**

1. Tambahkan tombol di area header:

```html
<button class="btn btn-sm btn-outline-secondary" onclick="exportToCSV()">
    <i class="fa-solid fa-file-csv me-1"></i>Export CSV
</button>
```

2. Fungsi generate dan download CSV:

```js
function exportToCSV() {
    if (!usersData || usersData.length === 0) {
        showToast('error', 'Tidak ada data untuk di-export.');
        return;
    }

    var headers = ['Nama', 'Email', 'Role', 'Status', 'Tanggal Dibuat'];
    var rows = usersData.map(function(u) {
        return [
            '"' + (u.full_name || '').replace(/"/g, '""') + '"',
            '"' + (u.email    || '').replace(/"/g, '""') + '"',
            '"' + (u.role     || '').replace(/"/g, '""') + '"',
            u.status === 'active' ? 'Aktif' : 'Nonaktif',
            '"' + formatDate(u.created_at) + '"'
        ].join(',');
    });

    var csvContent = '﻿' + headers.join(',') + '\n' + rows.join('\n');
    // ﻿ = BOM agar Excel membaca UTF-8 dengan benar

    var blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    var url  = URL.createObjectURL(blob);
    var a    = document.createElement('a');
    a.href   = url;
    a.download = 'akun-internal-' + new Date().toISOString().slice(0, 10) + '.csv';
    a.click();
    URL.revokeObjectURL(url);
}
```

---

## Peserta 16 <a id="peserta-16"></a>

### Soal 31 — [UPDATE] Superadmin Manajemen Akun: Toggle Status Aktif/Nonaktif Inline <a id="soal-31"></a>

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Ganti tombol "Aktifkan/Nonaktifkan" yang ada dengan **toggle switch** langsung di kolom Status tabel. Klik toggle memunculkan konfirmasi, lalu mengirim `PATCH /api/admin/users/:id/status` untuk memperbarui status tanpa reload halaman.

**Langkah pengerjaan:**

1. Saat render baris tabel, ganti badge status statis dengan toggle:

```js
function renderStatusToggle(u) {
    var checked  = u.status === 'active' ? 'checked' : '';
    var label    = u.status === 'active' ? 'Aktif' : 'Nonaktif';
    var colorVar = u.status === 'active' ? '#22c55e' : '#94a3b8';

    return '<div class="form-check form-switch d-flex align-items-center gap-2 mb-0">' +
        '<input class="form-check-input" type="checkbox" role="switch" ' + checked +
        ' onchange="confirmToggleStatus(\'' + u.user_id + '\', \'' + u.full_name + '\', this)"' +
        ' style="width:2.2em; height:1.1em; cursor:pointer;">' +
        '<span class="small" style="color:' + colorVar + ';">' + label + '</span>' +
        '</div>';
}
```

2. Fungsi konfirmasi dan PATCH:

```js
function confirmToggleStatus(userId, userName, checkbox) {
    var willActivate = checkbox.checked;
    var actionText   = willActivate ? 'mengaktifkan' : 'menonaktifkan';

    var ok = confirm('Yakin ingin ' + actionText + ' akun "' + userName + '"?');
    if (!ok) {
        checkbox.checked = !checkbox.checked; // kembalikan toggle
        return;
    }

    toggleUserStatus(userId, willActivate, checkbox);
}

async function toggleUserStatus(userId, activate, checkbox) {
    try {
        var res = await fetch('/api/admin/users/' + userId + '/status', {
            method: 'PATCH',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ status: activate ? 'active' : 'inactive' })
        });
        var data = await res.json();
        if (!res.ok) {
            checkbox.checked = !checkbox.checked; // rollback toggle
            showToast('error', data.message || 'Gagal mengubah status.');
        } else {
            var label = checkbox.nextElementSibling;
            if (activate) {
                label.textContent = 'Aktif';
                label.style.color = '#22c55e';
            } else {
                label.textContent = 'Nonaktif';
                label.style.color = '#94a3b8';
            }
            showToast('success', 'Status akun berhasil diubah.');
        }
    } catch (e) {
        checkbox.checked = !checkbox.checked;
        showToast('error', 'Koneksi gagal.');
    }
}
```

---

> **ℹ️ Soal 31 — Tidak perlu endpoint baru.** Endpoint `PATCH /api/admin/users/{id}/deactivate` dan `…/reactivate` sudah terdaftar di [users/app/routes/router.go](../users/app/routes/router.go) baris 104–105 (controller `DeactivateUser`/`ReactivateUser`, role `SUPERADMIN`). Keduanya **tidak butuh body** (user ID diambil dari path). Ganti fetch di `toggleUserStatus()` seperti berikut:

```js
// Ganti baris: fetch('/api/admin/users/' + userId + '/status', { method: 'PATCH', ... })
// Dengan:
var endpoint = activate
    ? '/api/admin/users/' + userId + '/reactivate'
    : '/api/admin/users/' + userId + '/deactivate';

var res = await fetch(endpoint, { method: 'PATCH', credentials: 'include' });
```

---

### Soal 32 — [DELETE] Superadmin Manajemen Akun: Hapus Massal dengan Checkbox <a id="soal-32"></a>

**File:** `frontend/ops/manage-users.html`

**Apa yang harus dilakukan:**
Tambahkan checkbox di setiap baris tabel dan tombol **"Hapus Terpilih"** di header. Tombol hanya aktif saat minimal 1 checkbox dicentang, memunculkan konfirmasi jumlah akun yang akan dihapus, lalu mengirim request DELETE.

**Langkah pengerjaan:**

1. Tambahkan checkbox header dan per-baris:

```html
<!-- Header th pertama -->
<th style="padding:16px 20px; width:40px;">
    <input type="checkbox" id="checkAll" onchange="toggleCheckAll(this)">
</th>

<!-- Dalam render baris, td pertama -->
// row.innerHTML = '<td><input type="checkbox" class="row-check" value="' + u.user_id + '"></td>' + row.innerHTML;
```

2. Tombol hapus massal (di toolbar atas tabel):

```html
<button id="btnBulkDelete" class="btn btn-sm btn-outline-danger" style="display:none;"
        onclick="confirmBulkDelete()">
    <i class="fa-solid fa-trash me-1"></i>Hapus Terpilih
    (<span id="selectedCount">0</span>)
</button>
```

3. Logika checkbox dan hapus:

```js
function toggleCheckAll(master) {
    document.querySelectorAll('.row-check').forEach(function(cb) {
        cb.checked = master.checked;
    });
    updateBulkDeleteBtn();
}

document.addEventListener('change', function(e) {
    if (e.target.classList.contains('row-check')) updateBulkDeleteBtn();
});

function updateBulkDeleteBtn() {
    var checked = document.querySelectorAll('.row-check:checked');
    var btn = document.getElementById('btnBulkDelete');
    document.getElementById('selectedCount').textContent = checked.length;
    btn.style.display = checked.length > 0 ? 'inline-flex' : 'none';
}

function confirmBulkDelete() {
    var ids = Array.from(document.querySelectorAll('.row-check:checked'))
        .map(function(cb) { return cb.value; });
    if (!ids.length) return;

    var ok = confirm('Hapus ' + ids.length + ' akun yang dipilih?\nTindakan ini tidak dapat dibatalkan.');
    if (!ok) return;

    bulkDeleteUsers(ids);
}

async function bulkDeleteUsers(ids) {
    try {
        var res = await fetch('/api/admin/users/bulk-delete', {
            method: 'DELETE',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ids: ids })
        });
        var data = await res.json();
        if (res.ok) {
            showToast('success', ids.length + ' akun berhasil dihapus.');
            loadUsers(1);
        } else {
            showToast('error', data.message || 'Gagal menghapus akun.');
        }
    } catch (e) {
        showToast('error', 'Koneksi gagal.');
    }
}
```

---

> **⚠️ Endpoint `DELETE /api/admin/users/bulk-delete` belum ada di backend — kerjakan langkah berikut sebelum implementasi frontend.**

#### Panduan Backend — Soal 32 <a id="backend-soal-32"></a>

**Service:** `users` · **Framework:** ZaFramework

> ⚠️ **Cek skema dulu.** Tabel `users` di project ini berkolom: `id (UUID)`, `name`, `email`, `role_id`, `status`. **Tidak ada** kolom `role`, **tidak ada** `deleted_at`, dan `status` hanya menerima `'banned' | 'active' | 'inactive'` (CHECK constraint di `migrate.go` baris 66). Jadi query lama `SET status='deleted', deleted_at=NOW() … role != 'SUPERADMIN'` **pasti error** (3 sebab: kolom `deleted_at` tidak ada, kolom `role` tidak ada, nilai `'deleted'` melanggar CHECK). Versi di bawah memakai **hard-delete** — aman karena semua foreign key ke `users(id)` memakai `ON DELETE SET NULL`/`CASCADE`, dan baris langsung hilang dari tabel sehingga cocok dengan UX "dihapus".

**Langkah 1 — [users/app/routes/router.go](../users/app/routes/router.go)**

Tambah 1 baris setelah route `…/reactivate` (sekitar baris 105). Urutan tidak masalah: Go 1.22 `ServeMux` memilih pola paling spesifik dan **tidak ada** route `DELETE /api/admin/users/{id}`, jadi tidak bentrok:
```go
app.Router.HandleFunc("DELETE /api/admin/users/bulk-delete", adminController.BulkDeleteUsers)
```

**Langkah 2 — [users/app/modules/admin/repository.go](../users/app/modules/admin/repository.go)**

a) Di dalam `type Repository interface { … }` (setelah `InsertAuditLog`, sekitar baris 24), tambah:
```go
	BulkDeleteByIDs(ctx context.Context, ids []string) (int, error)
```

b) Implementasi di akhir file. Hanya butuh import `context` yang **sudah ada**:
```go
// BulkDeleteByIDs menghapus banyak akun sekaligus, kecuali SUPERADMIN.
func (r *adminRepo) BulkDeleteByIDs(ctx context.Context, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	tag, err := r.db.Pool.Exec(ctx, `
		DELETE FROM users
		WHERE id::text = ANY($1)
		  AND role_id NOT IN (SELECT id FROM roles WHERE UPPER(code) = 'SUPERADMIN')
	`, ids)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}
```
> `id::text = ANY($1)` membandingkan UUID (di-cast ke text) dengan `[]string` — pgx meng-encode slice jadi array Postgres, jadi tidak perlu membangun placeholder `$1,$2,…` manual (yang butuh `fmt` & `strings`).

**Langkah 3 — [users/app/modules/admin/service.go](../users/app/modules/admin/service.go)**

> ⚠️ Sama seperti Soal 24: service ini **tidak punya `switch`/`case`**. Jangan tulis `case "admin_bulk_delete_users":` — akan error **`undefined: p`**. Buat fungsi `Process…Job` baru.

Tambah fungsi (mis. setelah `ProcessReactivateUserJob`, akhir file sekitar baris 399). Import `context` & `fmt` sudah ada:
```go
// ProcessBulkDeleteUsersJob — worker untuk hapus massal akun internal
func (s *Service) ProcessBulkDeleteUsersJob(ctx context.Context, payload any) (any, error) {
	data, ok := payload.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid payload format")
	}
	ids, ok := data["ids"].([]string)
	if !ok || len(ids) == 0 {
		return nil, fmt.Errorf("ids wajib diisi")
	}
	count, err := s.repo.BulkDeleteByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return map[string]any{"deleted": count}, nil
}
```

**Langkah 4 — [users/main.go](../users/main.go)** — daftarkan job-nya (kalau tidak, dispatch error "job not found")

Tambah 1 baris di blok `RegisterJob`, setelah `admin_reactivate_user` (sekitar baris 167):
```go
app.RegisterJob("admin_bulk_delete_users", adminService.ProcessBulkDeleteUsersJob)
```

**Langkah 5 — [users/app/modules/admin/controller.go](../users/app/modules/admin/controller.go)**

Tambah method di akhir file. Semua import (`strings`, `encoding/json`, `net/http`, `concurrency`, `ehttp`) **sudah dipakai** controller lain di file ini:
```go
func (c *Controller) BulkDeleteUsers(w http.ResponseWriter, r *http.Request) {
	userRole := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Role")))
	userLevel := strings.ToUpper(strings.TrimSpace(r.Header.Get("X-User-Level")))
	if userRole != "SUPERADMIN" && userLevel != "SUPERADMIN" {
		ehttp.ApiJSON(w, r, http.StatusForbidden, false, "Akses ditolak", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var body struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.IDs) == 0 {
		ehttp.ApiJSON(w, r, http.StatusBadRequest, false, "Field 'ids' wajib diisi", nil)
		return
	}

	result, err := c.Dispatcher.DispatchAndWait(r.Context(), "admin_bulk_delete_users",
		map[string]any{"ids": body.IDs}, concurrency.PriorityHigh)
	if err != nil {
		ehttp.ApiJSON(w, r, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}
	ehttp.ApiJSON(w, r, http.StatusOK, true, "Akun berhasil dihapus", result)
}
```

> **Gateway:** path `/api/admin/users/bulk-delete` sudah ter-cover prefix `/api/admin/users/` (role `SUPERADMIN`) di `gateway/routes.json` & `routes.railway.json` — tidak perlu menambah route gateway.

---

## Peserta 17 <a id="peserta-17"></a>

### Soal 33 — [CREATE] Superadmin Kinerja Divisi: Tambah Catatan/Note pada Divisi <a id="soal-33"></a>

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan tombol **"+ Catatan"** pada setiap kartu divisi. Tombol membuka modal input teks. Catatan disimpan ke `localStorage` (tidak membutuhkan API) dan ditampilkan di bawah kartu sebagai sticky note.

**Langkah pengerjaan:**

1. Tambahkan tombol dan area catatan di setiap `.division-card`:

```html
<!-- Di dalam setiap kartu divisi, setelah metric terakhir -->
<div class="mt-3 pt-2" style="border-top:1px solid var(--border-color);">
    <div id="note-compliance" class="small text-muted fst-italic mb-1"
         style="min-height:20px;"></div>
    <button class="btn btn-link p-0 small text-muted"
            onclick="openNoteModal('compliance', 'Divisi Compliance')">
        <i class="fa-solid fa-pen-to-square me-1"></i>+ Catatan
    </button>
</div>
```

2. Modal catatan (satu modal reusable):

```html
<div class="modal fade" id="noteModal" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered modal-sm">
        <div class="modal-content" style="background:var(--bg-card); border:1px solid var(--border-color);">
            <div class="modal-body p-3">
                <h6 class="fw-bold mb-2" id="noteModalTitle"></h6>
                <textarea id="noteInput" class="form-control" rows="4"
                          placeholder="Tulis catatan untuk divisi ini..."
                          style="background:var(--bg-main); border-color:var(--border-color); color:var(--text-heading);"></textarea>
                <div class="d-flex gap-2 mt-3">
                    <button class="btn btn-sm btn-outline-secondary flex-fill"
                            data-bs-dismiss="modal">Batal</button>
                    <button class="btn btn-sm flex-fill" id="btnSaveNote"
                            style="background:var(--accent-cyan); color:#000;">
                        <i class="fa-solid fa-floppy-disk me-1"></i>Simpan
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>
```

3. Logika simpan ke localStorage:

```js
var currentNoteKey = '';

function openNoteModal(divisionKey, divisionName) {
    currentNoteKey = 'division_note_' + divisionKey;
    document.getElementById('noteModalTitle').textContent = 'Catatan: ' + divisionName;
    document.getElementById('noteInput').value = localStorage.getItem(currentNoteKey) || '';
    new bootstrap.Modal(document.getElementById('noteModal')).show();
}

document.getElementById('btnSaveNote').addEventListener('click', function() {
    var text = document.getElementById('noteInput').value.trim();
    if (text) {
        localStorage.setItem(currentNoteKey, text);
    } else {
        localStorage.removeItem(currentNoteKey);
    }
    bootstrap.Modal.getInstance(document.getElementById('noteModal')).hide();
    renderNotes(); // refresh tampilan catatan
});

function renderNotes() {
    ['compliance', 'support', 'operational'].forEach(function(key) {
        var el = document.getElementById('note-' + key);
        if (!el) return;
        var note = localStorage.getItem('division_note_' + key);
        el.textContent = note ? '"' + note + '"' : '';
    });
}

document.addEventListener('DOMContentLoaded', renderNotes);
```

---

### Soal 34 — [READ] Superadmin Kinerja Divisi: Chart Bar Perbandingan Antar Divisi <a id="soal-34"></a>

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan visualisasi chart bar sederhana (tanpa library eksternal — menggunakan CSS + flexbox) yang membandingkan persentase penyelesaian task tiap divisi secara berdampingan.

**Langkah pengerjaan:**

1. Tambahkan section chart sebelum `#cardsRow`:

```html
<div class="ai-card p-4 mb-4" id="divisionChart">
    <div class="d-flex justify-content-between align-items-center mb-3">
        <h6 class="fw-bold mb-0" style="color:var(--text-heading);">
            <i class="fa-solid fa-chart-bar me-2" style="color:var(--accent-cyan);"></i>
            Perbandingan Penyelesaian Task
        </h6>
        <span class="small text-muted" id="chartUpdated"></span>
    </div>
    <div id="chartBars" class="d-flex align-items-end gap-4 justify-content-center"
         style="height:160px;"></div>
    <div id="chartLegend" class="d-flex gap-4 justify-content-center mt-3 flex-wrap"></div>
</div>
```

2. Fungsi render chart (panggil setelah data API diterima):

```js
var DIVISION_COLORS = {
    compliance: '#00d4ff',
    support:    '#22c55e',
    operational:'#f59e0b'
};

function renderDivisionChart(divisions) {
    var barsEl   = document.getElementById('chartBars');
    var legendEl = document.getElementById('chartLegend');
    if (!barsEl) return;

    barsEl.innerHTML   = '';
    legendEl.innerHTML = '';

    document.getElementById('chartUpdated').textContent =
        'Update: ' + new Date().toLocaleTimeString('id-ID');

    divisions.forEach(function(div) {
        var resolved = div.resolved || div.done || 0;
        var pending  = div.pending  || div.open  || 0;
        var total    = resolved + pending;
        var pct      = total > 0 ? Math.round((resolved / total) * 100) : 0;
        var color    = DIVISION_COLORS[div.key] || '#94a3b8';

        // Bar
        barsEl.innerHTML +=
            '<div class="d-flex flex-column align-items-center" style="flex:1; max-width:120px;">' +
            '<div class="small fw-bold mb-1" style="color:' + color + ';">' + pct + '%</div>' +
            '<div style="width:48px; height:' + Math.max(pct, 4) + '%; background:' + color + '; ' +
            'border-radius:6px 6px 0 0; transition:height 0.6s ease; min-height:4px;"></div>' +
            '</div>';

        // Legend
        legendEl.innerHTML +=
            '<div class="d-flex align-items-center gap-1 small">' +
            '<span style="width:12px; height:12px; background:' + color + '; border-radius:3px; display:inline-block;"></span>' +
            '<span style="color:var(--text-muted);">' + (div.name || div.key) + '</span>' +
            '</div>';
    });
}
```

---

## Peserta 18 <a id="peserta-18"></a>

### Soal 35 — [UPDATE] Superadmin Kinerja Divisi: Edit Target Task Divisi Inline <a id="soal-35"></a>

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan field **"Target"** yang dapat di-edit inline pada setiap kartu divisi: double-click untuk mengaktifkan edit mode, tekan Enter atau klik luar untuk menyimpan. Target dibandingkan dengan aktual dan ditampilkan dengan warna berbeda (hijau/merah).

**Langkah pengerjaan:**

1. Tambahkan elemen target di setiap kartu:

```html
<!-- Di dalam .division-card, di bawah metric resolved -->
<div class="mt-2 small" style="color:var(--text-muted);">
    Target:
    <span class="target-display fw-semibold" data-division="compliance"
          title="Double-click untuk edit" style="cursor:pointer; text-decoration:underline dotted;">
        —
    </span>
    <span class="target-diff ms-1 small"></span>
</div>
```

2. Logika inline edit dengan localStorage:

```js
// Inisialisasi target dari localStorage
function initTargets() {
    document.querySelectorAll('.target-display').forEach(function(el) {
        var key = 'division_target_' + el.getAttribute('data-division');
        var saved = localStorage.getItem(key);
        el.textContent = saved ? saved + ' task' : 'Set target';

        el.addEventListener('dblclick', function() {
            var cur = parseInt(localStorage.getItem(key)) || 0;
            var input = document.createElement('input');
            input.type  = 'number';
            input.value = cur;
            input.min   = '0';
            input.style.cssText = 'width:70px; font-size:0.8rem; padding:1px 4px; ' +
                'background:var(--bg-main); border:1px solid var(--accent-cyan); ' +
                'color:var(--text-heading); border-radius:4px;';

            el.replaceWith(input);
            input.focus();
            input.select();

            function save() {
                var val = parseInt(input.value) || 0;
                localStorage.setItem(key, val);
                // Kembalikan ke tampilan
                el.textContent = val + ' task';
                input.replaceWith(el);
                updateTargetDiff(el.getAttribute('data-division'), val);
            }

            input.addEventListener('keydown', function(e) {
                if (e.key === 'Enter') save();
                if (e.key === 'Escape') input.replaceWith(el);
            });
            input.addEventListener('blur', save);
        });
    });
}

function updateTargetDiff(divisionKey, target) {
    var divData = window.divisionData && window.divisionData.find(function(d) {
        return d.key === divisionKey;
    });
    if (!divData) return;

    var actual = divData.resolved || divData.done || 0;
    var diff   = actual - target;
    var diffEl = document.querySelector('.target-display[data-division="' + divisionKey + '"]')
        .nextElementSibling;
    if (!diffEl) return;

    if (diff >= 0) {
        diffEl.textContent = '(+' + diff + ' ✓)';
        diffEl.style.color = '#22c55e';
    } else {
        diffEl.textContent = '(' + diff + ' ⚠)';
        diffEl.style.color = '#ef4444';
    }
}

document.addEventListener('DOMContentLoaded', initTargets);
```

---

### Soal 36 — [DELETE] Superadmin Kinerja Divisi: Reset Data Kinerja Divisi <a id="soal-36"></a>

**File:** `frontend/management/dashboard-overview.html`

**Apa yang harus dilakukan:**
Tambahkan tombol **"Reset Data"** di setiap kartu divisi yang, setelah konfirmasi, mengirim `DELETE /api/management/divisions/:key/reset` ke API untuk menghapus data kinerja periode tersebut dan me-reload kartu.

**Langkah pengerjaan:**

1. Tambahkan tombol reset di setiap kartu (di samping nama divisi):

```html
<!-- Di header kartu divisi -->
<div class="d-flex justify-content-between align-items-center mb-3">
    <div class="fw-bold division-title" style="color:var(--text-heading);">Compliance</div>
    <button class="btn btn-link p-0 text-muted" title="Reset data divisi ini"
            onclick="confirmResetDivision('compliance', 'Compliance')">
        <i class="fa-solid fa-rotate-left fa-sm"></i>
    </button>
</div>
```

2. Fungsi konfirmasi dan DELETE:

```js
function confirmResetDivision(divisionKey, divisionName) {
    var ok = confirm(
        'Reset semua data kinerja divisi "' + divisionName + '"?\n\n' +
        'Semua task (resolved & pending) akan direset ke 0.\n' +
        'Tindakan ini TIDAK DAPAT DIBATALKAN.'
    );
    if (!ok) return;

    // Double confirm untuk aksi destruktif
    var reconfirm = confirm('Konfirmasi terakhir: Lanjutkan reset divisi ' + divisionName + '?');
    if (!reconfirm) return;

    resetDivisionData(divisionKey, divisionName);
}

async function resetDivisionData(divisionKey, divisionName) {
    try {
        var res = await fetch('/api/management/divisions/' + divisionKey + '/reset', {
            method: 'DELETE',
            credentials: 'include'
        });
        var data = await res.json();
        if (res.ok) {
            showToast('success', 'Data divisi ' + divisionName + ' berhasil direset.');
            loadData(); // reload semua kartu
        } else {
            showToast('error', data.message || 'Gagal mereset data divisi.');
        }
    } catch (e) {
        showToast('error', 'Koneksi gagal.');
    }
}
```

3. Setelah reset berhasil, tambahkan visual feedback kartu:

```js
// Di dalam loadData() setelah data dimuat, highlight kartu yang baru di-reset
function highlightCard(divisionKey) {
    var card = document.querySelector('[data-division="' + divisionKey + '"] .division-card');
    if (!card) return;
    card.style.transition = 'border-color 0.5s';
    card.style.borderColor = '#22c55e';
    setTimeout(function() { card.style.borderColor = ''; }, 2000);
}
```

---

> **⚠️ Endpoint `DELETE /api/management/divisions/:key/reset` belum ada di backend — kerjakan langkah berikut sebelum implementasi frontend.**

#### Panduan Backend — Soal 36 <a id="backend-soal-36"></a>

**Service:** `management` · **Framework:** Gin (berbeda dari `users` yang pakai ZaFramework)

> Signature router: `func Register(r *gin.Engine, dashboardHandler *dashboard.Handler)`. Endpoint ini **mock** untuk menguji frontend.

**Langkah 1 — [management/app/routes/router.go](../management/app/routes/router.go)**

a) Tambah import `"time"`. Blok import saat ini hanya `net/http`, `management/app/modules/dashboard`, dan `gin`. Jadikan:
```go
import (
	"net/http"
	"time"

	"management/app/modules/dashboard"

	"github.com/gin-gonic/gin"
)
```

b) Tambah handler **di dalam** fungsi `Register(...)`, sebelum kurung tutup `}` (mis. tepat setelah baris `r.GET("/api/superadmin/dashboard/overview", …)`):
```go
	// Soal 36: reset data kinerja divisi (mock). Cek role SUPERADMIN dari header yang di-inject gateway.
	r.DELETE("/api/management/divisions/:key/reset", func(c *gin.Context) {
		if c.GetHeader("X-User-Role") != "SUPERADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Akses ditolak"})
			return
		}
		key := c.Param("key")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Data divisi " + key + " berhasil direset",
			"data":    gin.H{"key": key, "reset_at": time.Now().Format(time.RFC3339)},
		})
	})
```
*(`net/http` dan `gin` sudah ter-import; hanya `time` yang perlu ditambah.)*

**Langkah 2 — Gateway (WAJIB agar request sampai ke service)** — [gateway/routes.json](../gateway/routes.json)

Gateway **belum** punya route `/api/management/…` (service `management` selama ini hanya diakses lewat `/api/dashboard/…`). Tanpa entri ini, fetch dari frontend kena **404 di gateway** dan tidak pernah sampai ke handler. Tambah objek berikut ke array `routes` (lakukan juga di [gateway/routes.railway.json](../gateway/routes.railway.json) dengan target railway, mis. `http://management.railway.internal:8080`):
```json
{
  "path": "/api/management/",
  "target": "http://localhost:5006",
  "cors": true,
  "auth": true,
  "roles": ["MANAGEMENT", "SUPERADMIN", "ADMIN"]
}
```
> Alternatif tanpa edit gateway (untuk uji cepat): panggil langsung port service `management` lewat REST client, mis. `DELETE http://localhost:5006/api/management/divisions/compliance/reset` dengan header `X-User-Role: SUPERADMIN`, lewati gateway.

---

## Soal Tambahan — Latihan Frontend Mandiri (Soal 37–58) <a id="soal-tambahan-latihan-frontend-mandiri-soal-37-58"></a>

> Semua soal berikut **frontend-only** (tanpa sentuh backend), **self-contained** (umumnya cukup tempel 1 blok JS di dalam `<script>` file terkait), dan sudah **diverifikasi** memakai ID elemen + fungsi yang benar-benar ada di file.
>
> 📚 **Soal 48–50 adalah tutorial pola** (step-by-step) untuk 3 keterampilan inti: **(48)** menambah tombol beserta fungsinya, **(49)** menambah field di halaman read/detail, **(50)** menambah pagination pada list yang belum punya.
>
> 🧩 **Soal 51–58 dikelompokkan per fitur** (KYC · Login/Logout/Register · Manajemen Akun · Kinerja Divisi), masing-masing 1 soal **tambah field** + 1 soal **tambah tombol**.
>
> ⚠️ **Tanda tangan `showToast` BERBEDA antar file — perhatikan urutan argumen:**
>
> | File | Signature | Pemanggilan yang BENAR |
> |---|---|---|
> | `login.html`, `register.html`, `verify-otp.html`, `manage-users.html` | `showToast(msg, type)` | `showToast('Berhasil', 'success')` |
> | `kyc-resubmit.html` | `showToast(type, message)` | `showToast('success', 'Berhasil')` |
> | `kyc-status.html` | *(tidak ada)* | pakai `alert(...)` |
> | `dashboard-overview.html` | *(tidak ada)* | pakai `showError(msg)` (banner) |

### Soal 37 — Login: Peringatan Caps Lock di Field Password <a id="soal-37"></a>

**File:** `frontend/account/login.html`

Tampilkan peringatan "Caps Lock aktif" di bawah field password saat Caps Lock menyala — mencegah user salah ketik password.

**Solusi** — tempel di dalam `<script>` (mis. setelah definisi `showToast`). Elemen `#password` ada di dalam `div.mb-3`:
```js
(function () {
    var pw = document.getElementById('password');
    var warn = document.createElement('div');
    warn.className = 'small mt-1';
    warn.style.cssText = 'color:#dc2626; display:none;';
    warn.innerHTML = '<i class="fa-solid fa-triangle-exclamation me-1"></i>Caps Lock sedang aktif';
    pw.closest('.mb-3').appendChild(warn);

    function toggle(e) {
        var on = e.getModifierState && e.getModifierState('CapsLock');
        warn.style.display = on ? 'block' : 'none';
    }
    pw.addEventListener('keydown', toggle);
    pw.addEventListener('keyup', toggle);
    pw.addEventListener('blur', function () { warn.style.display = 'none'; });
})();
```

---

### Soal 38 — Login: Ingat Email Terakhir (localStorage) <a id="soal-38"></a>

**File:** `frontend/account/login.html`

Saat halaman dibuka, isi otomatis field email dengan email login terakhir agar user tidak perlu mengetik ulang.

**Solusi** — tempel di dalam `<script>`. Memakai `validateLoginForm()` & `#loginForm` yang sudah ada; listener submit kedua ini **tidak** mengganggu listener submit asli:
```js
// Prefill email terakhir saat halaman dibuka
window.addEventListener('DOMContentLoaded', function () {
    var saved = localStorage.getItem('last_login_email');
    if (saved) {
        document.getElementById('email').value = saved;
        validateLoginForm(); // aktifkan tombol Sign In jika password juga terisi
    }
});

// Simpan email setiap kali form dikirim
document.getElementById('loginForm').addEventListener('submit', function () {
    localStorage.setItem('last_login_email', document.getElementById('email').value.trim());
});
```

---

### Soal 39 — Register: Tampilkan Umur Otomatis dari Tanggal Lahir <a id="soal-39"></a>

**File:** `frontend/account/register.html`

Setelah user memilih tanggal lahir (`#birthdate`), tampilkan umurnya ("Umur: 25 tahun") di bawah field.

**Solusi** — tempel di dalam `<script>`. Variabel `elBirthdate` sudah didefinisikan global; field-nya ada di dalam `div.mb-2`:
```js
(function () {
    var info = document.createElement('small');
    info.className = 'validation-msg';
    info.style.color = '#0077b6';
    elBirthdate.closest('.mb-2').appendChild(info);

    function showAge() {
        if (!elBirthdate.value) { info.textContent = ''; return; }
        var d = new Date(elBirthdate.value), now = new Date();
        var age = now.getFullYear() - d.getFullYear();
        var m = now.getMonth() - d.getMonth();
        if (m < 0 || (m === 0 && now.getDate() < d.getDate())) age--;
        info.textContent = age >= 0 ? 'Umur: ' + age + ' tahun' : '';
    }
    elBirthdate.addEventListener('change', showAge);
    elBirthdate.addEventListener('input', showAge);
})();
```

---

### Soal 40 — Register: Batasi Field Telepon Hanya Angka <a id="soal-40"></a>

**File:** `frontend/account/register.html`

Cegah user mengetik huruf/simbol di field telepon — strip otomatis non-digit dan beri umpan balik visual (border merah sekejap).

**Solusi** — tempel di dalam `<script>`. Class `.form-control.is-invalid` (border merah) sudah ada di CSS file ini; `elPhone` & `validateAll()` sudah ada:
```js
elPhone.addEventListener('input', function () {
    var before = this.value;
    this.value = this.value.replace(/\D/g, ''); // sisakan angka saja
    if (this.value !== before) {
        this.classList.add('is-invalid');
        var el = this;
        setTimeout(function () { el.classList.remove('is-invalid'); }, 800);
    }
    validateAll(); // sinkronkan status tombol Daftar
});
```

---

### Soal 41 — Verifikasi OTP: Samarkan (Mask) Alamat Email <a id="soal-41"></a>

**File:** `frontend/account/verify-otp.html`

Demi privasi, tampilkan email yang disamarkan (`jo****@gmail.com`) di `#emailDisplay`, bukan email penuh.

**Solusi** — di dalam `<script>`, **ganti** baris `document.getElementById('emailDisplay').textContent = email;` menjadi:
```js
document.getElementById('emailDisplay').textContent = maskEmail(email);

function maskEmail(e) {
    var parts = String(e).split('@');
    if (parts.length !== 2) return e;
    var name = parts[0];
    var masked = name.length <= 2
        ? name.charAt(0) + '*'
        : name.slice(0, 2) + '*'.repeat(Math.max(1, name.length - 2));
    return masked + '@' + parts[1];
}
```

---

### Soal 42 — Manajemen Akun: Shortcut `/` Fokus Pencarian + `Esc` Bersihkan <a id="soal-42"></a>

**File:** `frontend/ops/manage-users.html`

Tekan `/` di mana saja untuk langsung fokus ke kolom pencarian; tekan `Esc` saat di kolom pencarian untuk mengosongkan + reload.

**Solusi** — tempel di dalam `<script>`. Memakai `#searchInput` & `loadUsers(1)` yang sudah ada (tanpa edit HTML):
```js
document.addEventListener('keydown', function (e) {
    var search = document.getElementById('searchInput');
    var active = document.activeElement;
    var typing = active && ['INPUT', 'TEXTAREA', 'SELECT'].indexOf(active.tagName) !== -1;

    if (e.key === '/' && !typing) {       // fokuskan pencarian
        e.preventDefault();
        search.focus();
    }
    if (e.key === 'Escape' && active === search && search.value) { // bersihkan
        search.value = '';
        loadUsers(1);
    }
});
```

---

### Soal 43 — Manajemen Akun: Tombol "Salin Email" per Baris <a id="soal-43"></a>

**File:** `frontend/ops/manage-users.html`

Tambahkan tombol salin di setiap baris untuk menyalin email user ke clipboard.

**Langkah 1** — di dalam `renderActionButtons(u)`, **sebelum** baris `return`, tambahkan:
```js
var copyBtn = '<button class="btn btn-sm detail-btn" title="Salin email" ' +
    'data-email="' + escapeHtml(u.email) + '" onclick="copyEmail(this)">' +
    '<i class="fa-solid fa-copy"></i></button>';
```
lalu ubah baris `return` menjadi (sisipkan `copyBtn`):
```js
return '<div class="d-flex gap-2 justify-content-center">' + detailBtn + copyBtn + toggleBtn + '</div>';
```

**Langkah 2** — tambahkan fungsi (memakai `showToast(msg, type)` — urutan pesan dulu):
```js
function copyEmail(btn) {
    var email = btn.getAttribute('data-email');
    navigator.clipboard.writeText(email).then(function () {
        showToast('Email disalin: ' + email, 'success');
    }).catch(function () {
        showToast('Gagal menyalin email', 'error');
    });
}
```

---

### Soal 44 — Manajemen Akun: Avatar Inisial di Kolom Nama <a id="soal-44"></a>

**File:** `frontend/ops/manage-users.html`

Tampilkan lingkaran inisial (mis. "BS" untuk "Budi Santoso") di samping nama pada setiap baris tabel.

**Langkah 1** — tambahkan helper di antara render helper lain:
```js
function getInitials(name) {
    var parts = (name || '').trim().split(/\s+/);
    if (!parts[0]) return '?';
    var first = parts[0].charAt(0);
    var last = parts.length > 1 ? parts[parts.length - 1].charAt(0) : '';
    return (first + last).toUpperCase();
}
```

**Langkah 2** — di dalam `loadUsers()` pada `list.map(...)`, **ganti** cell Nama (`'<td style="padding:16px 20px;"><strong style="color:#e2e8f0;">' + escapeHtml(u.full_name) + '</strong></td>'`) menjadi:
```js
'<td style="padding:16px 20px;">' +
    '<div class="d-flex align-items-center gap-2">' +
        '<span style="width:34px;height:34px;border-radius:50%;background:rgba(0,180,216,0.15);' +
              'color:#00b4d8;display:inline-flex;align-items:center;justify-content:center;' +
              'font-weight:700;font-size:0.8rem;flex-shrink:0;">' + getInitials(u.full_name) + '</span>' +
        '<strong style="color:#e2e8f0;">' + escapeHtml(u.full_name) + '</strong>' +
    '</div>' +
'</td>' +
```

---

### Soal 45 — Kinerja Divisi: Auto-Refresh dengan Toggle <a id="soal-45"></a>

**File:** `frontend/management/dashboard-overview.html`

Tambahkan switch "Auto" di header; saat aktif, data di-refresh otomatis tiap 30 detik via `loadData()`.

**Langkah 1** — di header (di dalam `<div class="d-flex align-items-center gap-2">`, sebelum tombol Refresh), tambahkan:
```html
<div class="form-check form-switch d-flex align-items-center gap-1 mb-0 me-2">
    <input class="form-check-input" type="checkbox" role="switch" id="autoRefreshToggle">
    <label class="form-check-label small text-muted" for="autoRefreshToggle">Auto</label>
</div>
```

**Langkah 2** — tambahkan di dalam `<script>`:
```js
var autoRefreshTimer = null;
document.getElementById('autoRefreshToggle').addEventListener('change', function () {
    if (this.checked) {
        autoRefreshTimer = setInterval(loadData, 30000); // 30 detik
        loadData();
    } else {
        clearInterval(autoRefreshTimer);
        autoRefreshTimer = null;
    }
});
```

---

### Soal 46 — Kinerja Divisi: Kartu "Total Pending Semua Divisi" <a id="soal-46"></a>

**File:** `frontend/management/dashboard-overview.html`

Tampilkan satu kartu ringkasan di atas grid yang menjumlahkan semua antrean pending (KYC + tiket terbuka + order menunggu verifikasi).

**Langkah 1** — tambahkan **sebelum** `<div class="row g-4" id="cardsRow">`:
```html
<div class="mb-4">
    <div class="division-card p-3 d-flex align-items-center justify-content-between">
        <span class="metric-label">
            <i class="fa-solid fa-layer-group me-2" style="color:#fbbf24;"></i>Total Pending Semua Divisi
        </span>
        <span class="metric-val text-warning" id="totalPending">—</span>
    </div>
</div>
```

**Langkah 2** — di dalam `loadData()`, tepat setelah blok Operational (baris yang men-set `opsActive`), tambahkan. Field `kyc.pending`, `supp.open`, `ops.pending_payment` adalah nama asli dari API:
```js
// Total pending lintas divisi
var totalPending = (kyc.pending || 0) + (supp.open || 0) + (ops.pending_payment || 0);
document.getElementById('totalPending').textContent = totalPending;
```

---

### Soal 47 — KYC Resubmission: Penghitung Karakter Field Alamat <a id="soal-47"></a>

**File:** `frontend/client/kyc-resubmit.html`

Batasi alamat maksimal 255 karakter dan tampilkan penghitung "x / 255" yang berubah merah saat mendekati batas.

**Solusi** — tempel di dalam `<script>`. Helper `$ = id => document.getElementById(id)` & textarea `#address` sudah ada:
```js
(function () {
    var addr = $('address');
    if (!addr) return;
    var MAX = 255;
    addr.setAttribute('maxlength', MAX);

    var counter = document.createElement('small');
    counter.className = 'text-muted d-block mt-1';
    counter.style.textAlign = 'right';
    addr.parentNode.appendChild(counter);

    function update() {
        var len = addr.value.length;
        counter.textContent = len + ' / ' + MAX;
        counter.style.color = len > MAX - 20 ? '#dc2626' : '';
    }
    addr.addEventListener('input', update);
    update();
})();
```

> Catatan: di file ini `showToast(type, message)` urutannya **type dulu** (mis. `showToast('error', 'Pesan')`) — beda dari halaman lain. Soal ini sendiri tidak butuh toast.

---

### Soal 48 — Pola: Menambah Tombol Beserta Fungsinya (Studi Kasus: "Salin Email") <a id="soal-48"></a>

**File:** `frontend/ops/user-detail.html`

**Tujuan belajar:** memahami **anatomi 3 langkah** sebuah tombol yang berfungsi — *markup → handler → wiring* — plus memberi umpan balik ke user. Di halaman ini sudah ada grup tombol `#actionButtons` dan objek `currentUser` (hasil fetch detail). Kita tambah tombol **Salin Email**.

**Langkah 1 — Markup tombol.** Di dalam `<div class="action-group" id="actionButtons">`, tambah (mis. setelah tombol Edit):
```html
<button id="btnCopyEmail" class="btn btn-back" onclick="copyUserEmail()">
    <i class="fa-regular fa-copy me-1"></i> Salin Email
</button>
```

**Langkah 2 — Fungsi handler.** Tambah di dalam `<script>`. Memakai `currentUser` & `showToast(msg, type)` yang sudah ada:
```js
function copyUserEmail() {
    if (!currentUser || !currentUser.email) {
        showToast('Email tidak tersedia', 'error');
        return;
    }
    navigator.clipboard.writeText(currentUser.email).then(function () {
        showToast('Email disalin: ' + currentUser.email, 'success');
    }).catch(function () {
        showToast('Gagal menyalin email', 'error');
    });
}
```

**Langkah 3 — Wiring.** Sudah otomatis lewat `onclick="copyUserEmail()"` di markup. *(Alternatif tanpa `onclick`: `document.getElementById('btnCopyEmail').addEventListener('click', copyUserEmail);` — taruh setelah elemen tersedia.)*

> **Anatomi yang bisa dipakai ulang untuk tombol APA PUN:** *markup* (tombol + class + ikon) → *handler* (fungsi yang membaca data lalu melakukan aksi) → *wiring* (`onclick` atau `addEventListener`) → *feedback* (`showToast` / ubah teks tombol jadi loading). Untuk tombol yang memanggil API, tambahkan langkah `fetch(...)` di dalam handler + state `disabled` selama proses (lihat pola di Soal 43 & `executeAction()` pada file ini).

---

### Soal 49 — Pola: Menambah Field Baru di Halaman Read/Detail (Studi Kasus: "Terakhir Diperbarui") <a id="soal-49"></a>

**File:** `frontend/ops/user-detail.html`

**Tujuan belajar:** menambah satu baris data pada tampilan detail. Field `updated_at` **sudah dikirim API** (`GET /api/admin/users/{id}` → `data.updated_at`, lihat struct `UserDetail`) tetapi belum ditampilkan.

**Pola read view di file ini:** tiap data = satu `<div class="info-row">` berisi `.info-label` (judul) + `.info-value` (isi, ber-`id`), lalu diisi di fungsi `renderDetail(u)`.

**Langkah 1 — Tambah baris konten.** Di dalam `<div id="detailContent">`, setelah baris "Terakhir Login" (`#detailLastLogin`), tambah:
```html
<div class="info-row">
    <span class="info-label">Terakhir Diperbarui</span>
    <span class="info-value" id="detailUpdatedAt">-</span>
</div>
```

**Langkah 2 — Isi datanya.** Di dalam `renderDetail(u)`, setelah baris `document.getElementById('detailLastLogin').textContent = formatDate(u.last_login_at);`, tambah:
```js
document.getElementById('detailUpdatedAt').textContent = formatDate(u.updated_at);
```

**Langkah 3 (opsional, rapi) — Skeleton loading.** Agar tampilan loading tetap selaras, di dalam `#loadingState` tambah satu baris skeleton:
```html
<div class="info-row"><span class="info-label">Terakhir Diperbarui</span><span class="skeleton" style="width:160px;"></span></div>
```

> **Resep umum menambah field read:** (1) pastikan field sudah ada di response API → (2) tambah `info-row` + `id` di area konten → (3) isi di fungsi render → (4) (opsional) tambah skeleton. **Jika field belum dikirim API**, barulah perlu ubah backend: tambah kolom di `SELECT` (repository) + field di struct response — baru lanjut 4 langkah di atas.

---

### Soal 50 — Pola: Menambah Pagination (Client-Side) pada List yang Belum Punya <a id="soal-50"></a>

**File:** `frontend/ops/notifications.html`

**Tujuan belajar:** halaman ini me-render **semua** baris sekaligus via `renderTable(list)` — berat bila data banyak. Kita tambah **pagination sisi-klien** tanpa mengubah backend, memanfaatkan cache `cachedNotifications` & fungsi `renderTable` yang sudah ada.

**Langkah 1 — Elemen kontrol halaman.** Tambah tepat setelah `</table>` (penutup tabel notifikasi):
```html
<div id="notifPagination" class="d-flex justify-content-between align-items-center mt-3 px-1"></div>
```

**Langkah 2 — State + fungsi pagination.** Tambah di dalam `<script>` (mis. tepat di bawah definisi `renderTable`):
```js
// ── Pagination sisi-klien ──
var notifPage = 1;
var notifPerPage = 10;
var notifFullList = []; // daftar (sudah difilter) yang sedang ditampilkan

// Pengganti renderTable(...) — pakai saat daftar berubah (load / search / filter)
function renderWithPagination(list) {
    notifFullList = Array.isArray(list) ? list : [];
    notifPage = 1; // selalu kembali ke halaman 1 saat daftar berubah
    renderNotifPage();
}

function renderNotifPage() {
    var total = notifFullList.length;
    var totalPages = Math.max(1, Math.ceil(total / notifPerPage));
    if (notifPage > totalPages) notifPage = totalPages;

    var start = (notifPage - 1) * notifPerPage;
    var pageItems = notifFullList.slice(start, start + notifPerPage);

    renderTable(pageItems); // pakai renderer yang sudah ada
    renderNotifPagination(total, totalPages);
}

function renderNotifPagination(total, totalPages) {
    var el = document.getElementById('notifPagination');
    if (!el) return;
    if (total === 0) { el.innerHTML = ''; return; }

    var start = (notifPage - 1) * notifPerPage + 1;
    var end = Math.min(notifPage * notifPerPage, total);

    el.innerHTML =
        '<div class="small text-muted">Menampilkan ' + start + '–' + end + ' dari ' + total + ' data</div>' +
        '<div class="d-flex gap-2 align-items-center">' +
            '<button class="btn btn-sm btn-outline-secondary" ' + (notifPage <= 1 ? 'disabled' : '') +
                ' onclick="gotoNotifPage(' + (notifPage - 1) + ')"><i class="fa-solid fa-chevron-left"></i></button>' +
            '<span class="small text-muted">Hal ' + notifPage + ' / ' + totalPages + '</span>' +
            '<button class="btn btn-sm btn-outline-secondary" ' + (notifPage >= totalPages ? 'disabled' : '') +
                ' onclick="gotoNotifPage(' + (notifPage + 1) + ')"><i class="fa-solid fa-chevron-right"></i></button>' +
        '</div>';
}

function gotoNotifPage(p) {
    notifPage = p;
    renderNotifPage();
}
```

**Langkah 3 — Sambungkan ke alur yang ada.** Ganti **kedua** pemanggilan `renderTable(applyLocalSearch(cachedNotifications))` (satu di `loadNotifications()`, satu di handler pencarian/filter) menjadi:
```js
renderWithPagination(applyLocalSearch(cachedNotifications));
```

Hasil: list tampil 10 baris per halaman, search & filter tetap jalan (otomatis balik ke halaman 1), dan tombol Prev/Next nonaktif sendiri di ujung.

> **Klien vs server?** Pagination **klien** (di atas) cocok jika API mengirim semua data sekaligus & jumlahnya wajar (≤ beberapa ratus baris). Untuk data sangat besar, pagination harus di **server** lewat query `?page=&per_page=` — lihat pola asli di `manage-users.html` (`loadUsers(page)` + `renderPagination`).

---

## Tambah Field & Tombol per Fitur (Soal 51–58) <a id="tambah-field-tombol-per-fitur-soal-51-58"></a>

### 🔹 Fitur: KYC Resubmission <a id="fitur-kyc-resubmission"></a>

#### Soal 51 — Tombol "Perbarui Status" (tanpa reload) <a id="soal-51"></a>

**File:** `frontend/client/kyc-status.html`

User pending sering me-refresh halaman untuk cek apakah sudah di-review. Tambah tombol yang memanggil ulang `fetchKYCStatus()` (sudah ada) tanpa reload.

**Langkah 1** — di dalam kartu header (yang memuat `#statusBadge`), setelah `<div id="statusMessage">…</div>`, tambah:
```html
<button class="btn btn-sm btn-outline-secondary mt-3" onclick="fetchKYCStatus()">
    <i class="fa-solid fa-arrows-rotate me-1"></i>Perbarui Status
</button>
```
Selesai — `fetchKYCStatus()` akan fetch ulang `/api/kyc/status` dan memanggil `displayKYCData()` untuk merender ulang. Tidak perlu JS tambahan.

#### Soal 52 — Field "Nomor Pengajuan" di View Read <a id="soal-52"></a>

**File:** `frontend/client/kyc-status.html`

API `/api/kyc/status` mengirim `data.kyc.id` (lihat struct `KYCStatusResult`), tapi belum ditampilkan. Tambahkan di kartu **Address & Status**.

**Langkah 1** — setelah baris "Last Updated" (`#dataReviewedAt`), tambah:
```html
<div class="info-row">
    <span class="info-label">Nomor Pengajuan</span>
    <span class="info-value" id="dataKycId" style="font-family:monospace; font-size:0.8rem;">-</span>
</div>
```

**Langkah 2** — di dalam `displayKYCData(data)`, setelah baris `document.getElementById('dataReviewedAt').textContent = ...`, tambah:
```js
document.getElementById('dataKycId').textContent = kyc.id || '-';
```

---

### 🔹 Fitur: Login / Logout / Register <a id="fitur-login-logout-register"></a>

> Untuk **Logout**, soal tombol sudah ada: Soal 7 (konfirmasi keluar) & Soal 28 (keluar semua perangkat).

#### Soal 53 — Field "Konfirmasi Email" + Validasi Saat Submit <a id="soal-53"></a>

**File:** `frontend/account/register.html`

Cegah salah ketik email dengan field konfirmasi. Karena `validateAll()` tidak tahu field baru ini, kita validasi **saat submit** memakai listener `capture` di `document` (jalan **sebelum** handler submit asli di form).

**Langkah 1** — tambah setelah blok field Email (`#email`):
```html
<div class="mb-2">
    <label class="form-label">Konfirmasi Email</label>
    <div class="input-group">
        <span class="input-group-text"><i class="fa-regular fa-envelope"></i></span>
        <input type="email" class="form-control" id="confirmEmail" placeholder="Ulangi email" required>
    </div>
    <div class="validation-msg" id="confirmEmailMsg"></div>
</div>
```

**Langkah 2** — tambah di dalam `<script>` (memakai `setMsg`/`clearMsg` yang sudah ada):
```js
// Validasi konfirmasi email saat submit — capture di document jalan lebih dulu
// daripada handler submit form, sehingga bisa membatalkannya bila tidak cocok.
document.addEventListener('submit', function (e) {
    if (!e.target || e.target.id !== 'registerForm') return;
    var email = document.getElementById('email').value.trim().toLowerCase();
    var conf  = document.getElementById('confirmEmail').value.trim().toLowerCase();
    if (email !== conf) {
        e.preventDefault();
        e.stopPropagation();   // cegah event sampai ke handler submit utama
        setMsg('confirmEmailMsg', 'Email konfirmasi tidak cocok', true);
    } else {
        clearMsg('confirmEmailMsg');
    }
}, true); // <-- true = fase capture (kunci agar jalan duluan)
```
> Kenapa `capture` di `document`, bukan di form? Karena form adalah *target* event submit-nya sendiri; dua listener di elemen yang sama dijalankan sesuai urutan pendaftaran. Listener `capture` di **ancestor** (`document`) dijamin jalan sebelum listener di target.

#### Soal 54 — Tombol "Bersihkan Formulir" (Reset) <a id="soal-54"></a>

**File:** `frontend/account/register.html`

Tambah tombol untuk mengosongkan semua field + pesan validasi sekaligus.

**Langkah 1** — tambah setelah tombol `#btnRegister`:
```html
<button type="button" class="btn btn-link w-100 mt-2 small text-muted" id="btnClearForm">
    <i class="fa-solid fa-eraser me-1"></i>Bersihkan Formulir
</button>
```

**Langkah 2** — tambah di dalam `<script>` (memakai `clearMsg`, `validateAll`, `elFullName` yang sudah ada):
```js
document.getElementById('btnClearForm').addEventListener('click', function () {
    document.getElementById('registerForm').reset();
    ['fullNameMsg', 'emailMsg', 'phoneMsg', 'birthdateMsg', 'passwordMsg', 'confirmMsg'].forEach(clearMsg);
    validateAll();        // tombol Daftar otomatis nonaktif lagi (field kosong)
    elFullName.focus();
});
```
> Jika Anda juga mengerjakan Soal 53, tambahkan `'confirmEmailMsg'` ke array agar pesannya ikut dibersihkan.

---

### 🔹 Fitur: Superadmin Manajemen Akun Internal (CRUD) <a id="fitur-superadmin-manajemen-akun-internal-crud"></a>

#### Soal 55 — Field "Jumlah per Halaman" (Page Size) <a id="soal-55"></a>

**File:** `frontend/ops/manage-users.html`

Beri user kontrol berapa baris per halaman. `loadUsers()` sudah mengirim `per_page`, tinggal buat nilainya dapat diubah.

**Langkah 1** — ubah deklarasi konstanta menjadi variabel agar bisa diganti:
```js
// Ganti baris:  const perPage = 20;
// Menjadi:
let perPage = 20;
```

**Langkah 2** — tambah selektor di area filter (di samping dropdown `#filterStatus`):
```html
<select class="form-select form-select-sm" id="filterPerPage" style="width:auto;"
        onchange="changePerPage(this.value)">
    <option value="10">10 / halaman</option>
    <option value="20" selected>20 / halaman</option>
    <option value="50">50 / halaman</option>
    <option value="100">100 / halaman</option>
</select>
```

**Langkah 3** — tambah fungsi (memakai `loadUsers` yang sudah ada):
```js
function changePerPage(val) {
    perPage = parseInt(val, 10) || 20;
    loadUsers(1); // muat ulang dari halaman 1 dengan ukuran baru
}
```

#### Soal 56 — Tombol "Buat Password Kuat" (Generate) <a id="soal-56"></a>

**File:** `frontend/ops/create-user.html`

Bantu superadmin membuat password acak yang kuat dengan sekali klik, langsung mengisi field password + konfirmasi.

**Langkah 1** — tambah di dalam blok Password (mis. setelah `<div class="password-strength" id="passwordStrength">…</div>`):
```html
<button type="button" class="btn btn-sm btn-outline-secondary mt-2" id="btnGenPassword">
    <i class="fa-solid fa-wand-magic-sparkles me-1"></i>Buat Password Kuat
</button>
```

**Langkah 2** — tambah di dalam `<script>` (memakai `passwordInput`, `confirmPasswordInput`, `showToast` yang sudah ada):
```js
document.getElementById('btnGenPassword').addEventListener('click', function () {
    var chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789!@#$%';
    var pw = '';
    for (var i = 0; i < 14; i++) pw += chars.charAt(Math.floor(Math.random() * chars.length));

    // Tampilkan sebagai teks agar admin bisa mencatatnya
    passwordInput.type = 'text';
    confirmPasswordInput.type = 'text';
    passwordInput.value = pw;
    confirmPasswordInput.value = pw;

    // Picu validasi + bar kekuatan yang sudah ada
    passwordInput.dispatchEvent(new Event('input'));
    confirmPasswordInput.dispatchEvent(new Event('input'));

    showToast('Password kuat dibuat — jangan lupa disalin', 'success');
});
```
> Untuk keamanan lebih, ganti `Math.random()` dengan `crypto.getRandomValues(new Uint32Array(1))[0]` sebagai sumber acak.

---

### 🔹 Fitur: Superadmin Memantau Kinerja Divisi <a id="fitur-superadmin-memantau-kinerja-divisi"></a>

#### Soal 57 — Field Metrik "KYC Ditolak" di Kartu Compliance <a id="soal-57"></a>

**File:** `frontend/management/dashboard-overview.html`

API overview mengirim `data.kyc.rejected` (sudah dipakai untuk menghitung "Total Processed"), tapi belum ditampilkan terpisah. Tambah sebagai metrik baru.

**Langkah 1** — di dalam body kartu **Compliance** (setelah baris metrik `#kycProcessed`), tambah:
```html
<div class="metric-row">
    <span class="metric-label"><i class="fa-solid fa-circle-xmark me-1 text-danger" style="font-size:0.75rem;"></i>KYC Ditolak</span>
    <span class="metric-val text-danger" id="kycRejected">—</span>
</div>
```

**Langkah 2** — di dalam `loadData()`, setelah baris yang men-set `kycProcessed`, tambah:
```js
document.getElementById('kycRejected').textContent = kyc.rejected != null ? kyc.rejected : '—';
```
> Opsional rapi: tambahkan `'kycRejected'` ke kedua array reset state (`['kycPending', ...]`) agar ikut menampilkan `…` saat loading dan `—` saat error.

#### Soal 58 — Tombol "Salin Ringkasan" ke Clipboard <a id="soal-58"></a>

**File:** `frontend/management/dashboard-overview.html`

Tambah tombol untuk menyalin ringkasan semua metrik (mis. untuk ditempel ke laporan/chat). Halaman ini tidak punya `showToast`, jadi umpan balik lewat teks `#lastUpdated`.

**Langkah 1** — di header (di dalam `<div class="d-flex align-items-center gap-2">`, dekat tombol Refresh), tambah:
```html
<button class="btn btn-sm btn-outline-secondary" onclick="copySummary()" title="Salin ringkasan">
    <i class="fa-regular fa-copy me-1"></i>Salin Ringkasan
</button>
```

**Langkah 2** — tambah fungsi (memakai ID metrik & `#lastUpdated` yang sudah ada):
```js
function copySummary() {
    function v(id) { return document.getElementById(id).textContent; }
    var text =
        'Ringkasan Kinerja Divisi\n' +
        '- KYC Pending      : ' + v('kycPending')   + '\n' +
        '- KYC Diproses     : ' + v('kycProcessed') + '\n' +
        '- Tiket Terbuka    : ' + v('suppOpen')     + '\n' +
        '- Tiket Selesai    : ' + v('suppResolved') + '\n' +
        '- Order Menunggu   : ' + v('opsPending')   + '\n' +
        '- Langganan Aktif  : ' + v('opsActive');

    navigator.clipboard.writeText(text).then(function () {
        var el = document.getElementById('lastUpdated');
        var old = el.textContent;
        el.textContent = 'Ringkasan disalin ✓';
        setTimeout(function () { el.textContent = old; }, 2000);
    });
}
```

---

## Referensi Cepat — File yang Relevan <a id="referensi-cepat-file-yang-relevan"></a>

| Fitur | File Utama |
|---|---|
| KYC Resubmission | `frontend/client/kyc-resubmit.html` |
| Login | `frontend/account/login.html` |
| Register | `frontend/account/register.html` |
| Logout | `frontend/assets/js/ops-layout.js` / `client-layout.js` |
| Manajemen Akun Internal | `frontend/ops/manage-users.html`, `create-user.html`, `edit-user.html` |
| Kinerja Divisi (Overview) | `frontend/management/dashboard-overview.html` |
| Kinerja Divisi (Detail) | `frontend/management/dashboard-compliance.html`, `dashboard-support.html`, `dashboard-operational.html` |

## Pola Umum di Codebase Ini <a id="pola-umum-di-codebase-ini"></a>

- **Toast notifikasi:** fungsi `showToast(type, message)` — `type` berisi `'success'` atau `'error'`
- **Modal konfirmasi:** gunakan `.modal-overlay` + `.modal-box` yang sudah ada di `manage-users.html`
- **Fetch API:** selalu pakai `credentials: 'include'`, handle status 401 dengan redirect ke `/account/login`
- **Dark theme variables:** `--bg-card`, `--border-color`, `--text-heading`, `--text-muted`, `--accent-cyan`
- **Icons:** Font Awesome 6 (`fa-solid`, `fa-regular`) — cukup tambahkan class pada `<i>`
