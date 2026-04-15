# 🛡️ Auth & KYC Module - ThinkNalyze

**Branch:** `feat/PBI-KYC
**Assignee:** Christopher Matthew Hendarson (2306245592)

---

## 📖 Ringkasan Modul
Branch ini difokuskan untuk pengembangan fondasi keamanan dan legalitas akses aplikasi ThinkNalyze. Lingkup pengerjaan mencakup sistem **Autentikasi Pengguna** (Login/Register) dan sistem **Verifikasi Identitas (KYC)**. Modul ini bertindak sebagai "gerbang utama" (*Gatekeeper*) yang memastikan hanya pengguna terdaftar dan terverifikasi yang dapat mengakses fitur inti aplikasi.

---

## 🚀 Cakupan Fitur (Scope of Work)

Berikut adalah rincian fitur spesifik yang dikembangkan dalam branch ini:

### 1. 🔐 Authentication System (Auth)
Fitur untuk menangani akses masuk dan pendaftaran akun pengguna dengan standar keamanan yang baik.
- **Registration:** Formulir pendaftaran akun baru dengan validasi input (Email, Password, No. HP).
- **Login:** Mekanisme masuk pengguna menggunakan kredensial yang valid dan penyimpanan token sesi (JWT).
- **Session Management:** Pengelolaan masa aktif sesi login dan fungsi *Logout* yang aman.
- **Route Guard:** Proteksi halaman (*Private Routes*) agar tidak bisa diakses tanpa login.

### 2. 🆔 KYC Verification Center
Fitur wajib bagi pengguna baru untuk memvalidasi identitas sebelum bisa menggunakan layanan.
- **Submission Form:** Antarmuka pengunggahan foto KTP dan pengisian data NIK secara aman.
- **Status Dashboard:** Halaman khusus untuk memantau status pengajuan (*Pending, Approved, Rejected*).
- **Resubmission Logic:** Alur perbaikan data jika pengajuan sebelumnya ditolak oleh admin (menampilkan alasan penolakan dan form upload ulang).
- **Access Blocking:** Logika sistem yang membatasi akses ke menu utama jika status KYC belum *Approved*.

---

## 🔄 Alur Logika Utama (User Flow)

Berikut adalah diagram alur bagaimana Auth dan KYC bekerja sebagai satu kesatuan sistem keamanan:

```mermaid
graph TD
    Start((User)) -->|Buka Web| Login{Sudah Login?}
    
    Login -- "Belum" --> PageAuth[Halaman Login / Register]
    PageAuth -->|Submit| API_Auth[Validasi Kredensial]
    API_Auth -->|Sukses| SaveToken[Simpan Token & Redirect]
    
    Login -- "Sudah" --> CheckKYC{Cek Status KYC}
    
    CheckKYC -- "Belum / Rejected" --> PageKYC[Halaman Verifikasi]
    PageKYC -->|Upload KTP| API_KYC[Submit Data]
    PageKYC -->|Status Rejected| Resubmit[Perbaiki & Kirim Ulang]
    
    CheckKYC -- "Pending" --> PageWait[Halaman Menunggu]
    
    CheckKYC -- "Approved" --> Dashboard[Redirect ke Dashboard Utama]