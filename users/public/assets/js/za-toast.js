/**
 * Auto-Inject Toast Utility
 * Menangani HTML & Logic dalam satu file.
 */
(function() {
    // 1. Fungsi untuk Inject HTML ke Body
    const injectToastHTML = () => {
        if (document.getElementById('liveToast')) return; // Jangan inject dua kali

        const toastContainer = document.createElement('div');
        toastContainer.className = 'toast-container position-fixed top-0 end-0 p-3';
        toastContainer.style.zIndex = '1100';
        toastContainer.innerHTML = `
            <div id="liveToast" class="toast align-items-center border-0 rounded-4 shadow-lg" role="alert" aria-live="assertive" aria-atomic="true">
                <div class="d-flex">
                    <div class="toast-body d-flex align-items-center gap-2">
                        <i id="toastIcon" class="bi fs-5"></i>
                        <div>
                            <strong class="d-block" id="toastTitle"></strong>
                            <small id="toastMessage"></small>
                        </div>
                    </div>
                    <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
                </div>
            </div>
        `;
        document.body.appendChild(toastContainer);
    };

    // 2. Fungsi Utama showToast
    window.showToast = function(title, message, type = 'info') {
        // Pastikan HTML sudah ada di DOM
        injectToastHTML();

        const toastEl = document.getElementById('liveToast');
        const toastIcon = document.getElementById('toastIcon');
        const toastTitle = document.getElementById('toastTitle');
        const toastMessage = document.getElementById('toastMessage');

        const themes = {
            success: { bg: 'text-bg-success', icon: 'bi-check-circle-fill' },
            danger:  { bg: 'text-bg-danger',  icon: 'bi-exclamation-octagon-fill' },
            warning: { bg: 'text-bg-warning', icon: 'bi-exclamation-triangle-fill' },
            info:    { bg: 'text-bg-info',    icon: 'bi-info-circle-fill' },
            primary: { bg: 'text-bg-primary', icon: 'bi-megaphone-fill' }
        };

        // Reset & Set Tema
        const themeClasses = Object.values(themes).map(t => t.bg);
        toastEl.classList.remove(...themeClasses);
        
        const selected = themes[type] || themes.info;
        toastEl.classList.add(selected.bg);
        
        // Update Icon & Konten
        toastIcon.className = `bi fs-5 ${selected.icon}`;
        toastTitle.innerText = title;
        toastMessage.innerText = message;

        // Tampilkan
        const instance = bootstrap.Toast.getOrCreateInstance(toastEl);
        instance.show();
    };

    // Inject HTML segera setelah script dimuat jika DOM sudah siap
    if (document.body) {
        injectToastHTML();
    } else {
        document.addEventListener('DOMContentLoaded', injectToastHTML);
    }
})();