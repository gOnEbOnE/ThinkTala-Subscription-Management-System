/**
 * App Utilities: Theme Manager & Toast Notification
 */

(function () {
    // --- 1. LOGIKA THEME (DARK/LIGHT) ---
    const initTheme = () => {
        const themeToggleBtn = document.getElementById('themeToggle');
        const html = document.documentElement;
        
        // Ambil preferensi dari Local Storage atau default ke light
        const savedTheme = localStorage.getItem('theme') || 'light';
        html.setAttribute('data-bs-theme', savedTheme);

        // Jika tombol toggle ada di halaman ini, inisialisasi ikon dan event
        if (themeToggleBtn) {
            updateThemeUI(themeToggleBtn, savedTheme);

            themeToggleBtn.addEventListener('click', () => {
                const currentTheme = html.getAttribute('data-bs-theme');
                const newTheme = currentTheme === 'light' ? 'dark' : 'light';

                html.setAttribute('data-bs-theme', newTheme);
                localStorage.setItem('theme', newTheme);
                updateThemeUI(themeToggleBtn, newTheme);
            });
        }
    };

    const updateThemeUI = (btn, theme) => {
        const icon = btn.querySelector('i');
        if (!icon) return;

        if (theme === 'dark') {
            icon.className = 'bi bi-sun-fill';
            btn.classList.replace('btn-primary', 'btn-warning');
            btn.classList.add('text-dark');
        } else {
            icon.className = 'bi bi-moon-stars-fill';
            btn.classList.replace('btn-warning', 'btn-primary');
            btn.classList.remove('text-dark');
        }
    };

    // --- 2. LOGIKA TOAST (AUTO-INJECT HTML) ---
    const injectToastHTML = () => {
        if (document.getElementById('liveToast')) return;

        const container = document.createElement('div');
        container.className = 'toast-container position-fixed top-0 end-0 p-3';
        container.style.zIndex = '1100';
        container.innerHTML = `
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
            </div>`;
        document.body.appendChild(container);
    };

    window.showToast = function (title, message, type = 'info') {
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

        const allBgClasses = Object.values(themes).map(t => t.bg);
        toastEl.classList.remove(...allBgClasses);

        const selected = themes[type] || themes.info;
        toastEl.classList.add(selected.bg);
        
        toastIcon.className = `bi fs-5 ${selected.icon}`;
        toastTitle.innerText = title;
        toastMessage.innerText = message;

        const instance = bootstrap.Toast.getOrCreateInstance(toastEl);
        instance.show();
    };

    // Jalankan inisialisasi saat DOM siap
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initTheme);
    } else {
        initTheme();
    }
})();