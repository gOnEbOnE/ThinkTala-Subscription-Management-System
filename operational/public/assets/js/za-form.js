/**
 * SecureForm Module (v2)
 * Update: Added customizable initialDelay & maxDelay
 */
class SecureForm {
    constructor(formSelector, btnSelector, options = {}) {
        this.form = document.querySelector(formSelector);
        this.submitBtn = document.querySelector(btnSelector);

        // 1. Cek apakah elemen ditemukan
        if (!this.form) {
            console.error(`SecureForm Error: Elemen "${formSelector}" tidak ditemukan di halaman.`);
            return;
        }

        // 2. Cek apakah elemen tersebut benar-benar sebuah <form>
        if (!(this.form instanceof HTMLFormElement)) {
            console.error(`SecureForm Error: Selector "${formSelector}" menunjuk ke <${this.form.tagName.toLowerCase()}>, bukan ke tag <form>.`);
            return;
        }
        
        // --- KONFIGURASI BARU ---
        this.config = {
            onSuccess: options.onSuccess || null,
            csrfSelector: options.csrfSelector || 'meta[name="csrf-token"]',
            
            // Konfigurasi Delay
            initialDelay: options.initialDelay || 2000, // Default 2 detik
            maxDelay: options.maxDelay || 30000,        // Default 30 detik (Cap)
            
            // Konfigurasi Network
            maxTimeout: 8000,
            maxRetry: 2,
            ...options
        };

        // State Management
        if (this.form) {
            this.storageKeyFail = `failCount_${this.form.id}`;
            this.storageKeyCool = `cooldownUntil_${this.form.id}`;
            
            this.state = {
                locked: false,
                failCount: Number(localStorage.getItem(this.storageKeyFail) || 0),
                cooldownUntil: Number(localStorage.getItem(this.storageKeyCool) || 0),
            };

            this.countdownTimer = null;
            this.init();
        } else {
            console.error(`SecureForm: Form dengan selector "${formSelector}" tidak ditemukan.`);
        }
    }

    init() {
        // Cek cooldown saat halaman dimuat (Persistensi)
        if (this.isCooldownActive()) {
            const sisa = Math.ceil((this.state.cooldownUntil - Date.now()) / 1000);
            this.lockButton(`Tunggu ${sisa}s`, sisa * 1000);
        }

        this.form.addEventListener("submit", (e) => this.handleSubmit(e));
    }

    // --- LOGIC UTAMA ---

    async handleSubmit(e) {
        e.preventDefault();
        if (this.state.locked || this.isCooldownActive()) return;

        const originalBtnText = this.submitBtn.innerHTML;
        this.lockButton();

        try {
            const formData = new FormData(this.form);
            const payload = Object.fromEntries(formData.entries());
            const csrfToken = this.getCsrfToken();

            const res = await this.fetchWithRetry(this.form.action, {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                    "X-CSRF-Token": csrfToken,
                    "X-Requested-With": "XMLHttpRequest"
                },
                body: JSON.stringify(payload)
            });

            const json = await res.json();

            if (!json.status) {
                // Gagal Login / Validasi
                this.activateCooldown(json.msg || "Validasi gagal.");
                
                // Kembalikan tombol ke teks asli setelah cooldown selesai (handled by unlockButton inside timer)
                // Tapi karena activateCooldown memicu lockButton dengan timer, kita tidak perlu setHTML manual disini.
                return;
            }

            // --- SUKSES ---
            this.resetState();
            this.handleSuccessAction(json);

        } catch (err) {
            console.error(err);
            this.activateCooldown("Gagal terhubung ke server.");
        }
    }

    handleSuccessAction(responseJson) {
        if (window.showToast) {
            window.showToast("Berhasil", responseJson.msg || "Sukses.", "success");
        }

        const action = this.config.onSuccess;

        if (typeof action === 'string') {
            setTimeout(() => window.location.href = action, 1000);
        } else if (typeof action === 'function') {
            this.submitBtn.innerHTML = "Selesai";
            action(responseJson);
            this.state.locked = false;
            this.submitBtn.disabled = false;
        } else {
            this.submitBtn.innerHTML = "Kirim Data";
            this.state.locked = false;
            this.submitBtn.disabled = false;
            this.form.reset();
        }
    }

    // --- HELPER & SECURITY ---

    getCsrfToken() {
        const meta = document.querySelector(this.config.csrfSelector);
        return meta ? meta.content : (document.querySelector(this.config.csrfSelector)?.value || '');
    }

    isCooldownActive() {
        return this.state.cooldownUntil > Date.now();
    }

    /**
     * LOGIKA BARU UNTUK DELAY
     * Menggunakan initialDelay dari config
     */
    activateCooldown(msg) {
        this.state.failCount++;
        localStorage.setItem(this.storageKeyFail, this.state.failCount);

        // RUMUS: initialDelay * (2 pangkat (failCount - 1))
        // Contoh jika initialDelay = 2000 (2 detik):
        // Gagal ke-1: 2000 * 1 = 2 detik
        // Gagal ke-2: 2000 * 2 = 4 detik
        // Gagal ke-3: 2000 * 4 = 8 detik
        let calculatedDelay = this.config.initialDelay * (2 ** (this.state.failCount - 1));
        
        // Pastikan tidak melebihi batas maksimal (Max Cap)
        const delay = Math.min(calculatedDelay, this.config.maxDelay);
        
        this.state.cooldownUntil = Date.now() + delay;
        localStorage.setItem(this.storageKeyCool, this.state.cooldownUntil);

        if (window.showToast) window.showToast("Gagal", msg, 'warning');
        
        // Kunci tombol selama durasi delay
        this.lockButton(null, delay);
    }

    resetState() {
        this.state.failCount = 0;
        this.state.cooldownUntil = 0;
        localStorage.removeItem(this.storageKeyFail);
        localStorage.removeItem(this.storageKeyCool);
    }

    lockButton(text = "Memproses...", timerDelay = 0) {
        this.submitBtn.disabled = true;
        this.state.locked = true;

        if (timerDelay > 0) {
            if (this.countdownTimer) clearInterval(this.countdownTimer);
            
            let remaining = Math.ceil(timerDelay / 1000);
            this.submitBtn.innerHTML = `Tunggu ${remaining}s`;

            this.countdownTimer = setInterval(() => {
                remaining--;
                if (remaining <= 0) {
                    this.unlockButton();
                    return;
                }
                this.submitBtn.innerHTML = `Tunggu ${remaining}s`;
            }, 1000);
        } else {
            this.submitBtn.innerHTML = `<span class="spinner-border spinner-border-sm me-2"></span>${text}`;
        }
    }

    unlockButton() {
        if (this.countdownTimer) {
            clearInterval(this.countdownTimer);
            this.countdownTimer = null;
        }
        this.submitBtn.disabled = false;
        // Kita bisa kembalikan ke teks default, atau biarkan user mengatur ulang lewat callback
        // Untuk aman, kita set default umum:
        this.submitBtn.innerHTML = "Kirim Ulang"; 
        this.state.locked = false;
    }

    async fetchWithTimeout(url, options) {
        const controller = new AbortController();
        const id = setTimeout(() => controller.abort(), this.config.maxTimeout);
        try {
            const response = await fetch(url, { ...options, signal: controller.signal });
            clearTimeout(id);
            return response;
        } catch (error) {
            clearTimeout(id);
            throw error;
        }
    }

    async fetchWithRetry(url, options, retry = this.config.maxRetry) {
        try {
            const res = await this.fetchWithTimeout(url, options);
            if (!res.ok) throw new Error(`HTTP ${res.status}`);
            return res;
        } catch (err) {
            if (retry <= 0) throw err;
            return await this.fetchWithRetry(url, options, retry - 1);
        }
    }
}

window.SecureForm = SecureForm;