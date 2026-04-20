/* =====================================================
   MANAGEMENT PANEL — Shared Layout (Sidebar + Navbar)
   Usage: include this script in any management page.

   The page must have:
     <div id="ops-sidebar-placeholder"></div>
     <div id="ops-navbar-placeholder"></div>
   ===================================================== */

(function () {
    'use strict';

    let user = null;
    try {
        user = JSON.parse(localStorage.getItem('user') || 'null');
    } catch (_e) {
        user = null;
    }

    const hasLocalUser = !!(user && user.id);
    const roleCode = hasLocalUser
        ? String(user.role_code || user.level_code || user.level || '').toUpperCase()
        : '';

    const isCustomerPath = window.location.pathname === '/management/dashboard-customers' ||
        window.location.pathname.startsWith('/dashboard/customer/');
    const isPackagePath = window.location.pathname === '/management/dashboard-packages' ||
        window.location.pathname.startsWith('/dashboard/packages/');

    function isActive(flag) {
        return flag ? ' active' : '';
    }

    // Prevent transition flash on initial load.
    const preventTx = document.createElement('style');
    preventTx.id = 'prevent-tx';
    preventTx.innerHTML = '* { transition: none !important; }';
    document.head.appendChild(preventTx);
    window.addEventListener('load', function () {
        setTimeout(function () {
            const el = document.getElementById('prevent-tx');
            if (el) el.remove();
        }, 50);
    });

    if (window.innerWidth > 768 && localStorage.getItem('sidebar_state') === 'collapsed') {
        document.body.classList.add('sidebar-collapsed');
    }

    const sidebarHTML = '' +
'<nav class="sidebar">' +
'  <div class="sidebar-brand mb-3">' +
'    <div class="brand-wrapper">' +
'      <i class="fa-solid fa-chart-line text-cyan brand-icon"></i>' +
'      <div class="brand-text-content">' +
'        <h4 class="fw-bold tracking-wider mb-0" style="color: var(--text-heading)">Think<span class="text-cyan">Tala</span></h4>' +
'        <p class="small mb-0 text-muted" style="font-size: 0.7rem;">Management Panel</p>' +
'      </div>' +
'    </div>' +
'  </div>' +
'  <ul class="nav flex-column flex-grow-1">' +
'    <li class="nav-item">' +
'      <a class="nav-link' + isActive(isCustomerPath) + '" href="/management/dashboard-customers">' +
'        <i class="fa-solid fa-users icon-left"></i>' +
'        <span class="link-text">Customer Churn</span>' +
'      </a>' +
'    </li>' +
'    <li class="nav-item">' +
'      <a class="nav-link' + isActive(isPackagePath) + '" href="/management/dashboard-packages">' +
'        <i class="fa-solid fa-cubes icon-left"></i>' +
'        <span class="link-text">Package Sales</span>' +
'      </a>' +
'    </li>' +
'  </ul>' +
'  <ul class="nav flex-column mb-5">' +
'    <li class="nav-item">' +
'      <a class="nav-link text-danger" href="#" onclick="ManagementLayout.logout(); return false;">' +
'        <i class="fa-solid fa-right-from-bracket icon-left"></i>' +
'        <span class="link-text">Logout</span>' +
'      </a>' +
'    </li>' +
'  </ul>' +
'</nav>' +
'<div class="mobile-overlay" id="mobileOverlay"></div>';

    const navbarHTML = '' +
'<header class="top-navbar d-flex justify-content-between align-items-center">' +
'  <div class="d-flex align-items-center gap-2">' +
'    <button class="btn-header" id="sidebarToggle"><i class="fa-solid fa-bars fa-lg"></i></button>' +
'    <span class="badge bg-warning text-dark" id="roleBadge">MANAGEMENT</span>' +
'  </div>' +
'  <div class="d-flex align-items-center gap-2 gap-md-3">' +
'    <button class="btn-header" id="themeToggle"><i class="fa-solid fa-moon"></i></button>' +
'    <div class="dropdown ms-1">' +
'      <a href="#" class="d-flex align-items-center text-decoration-none" data-bs-toggle="dropdown">' +
'        <div class="rounded-circle p-1" style="border: 2px solid var(--accent-cyan);">' +
'          <img src="https://ui-avatars.com/api/?name=MGMT&background=0b0e17&color=fff" class="rounded-circle" id="navAvatar" width="34" height="34">' +
'        </div>' +
'      </a>' +
'      <ul class="dropdown-menu dropdown-menu-end dropdown-menu-animate mt-3">' +
'        <li class="px-3 py-2">' +
'          <span class="d-block fw-bold text-main" id="userName">Management User</span>' +
'          <small class="text-muted" id="userEmail">management@thinktala.com</small>' +
'        </li>' +
'        <li><hr class="dropdown-divider border-secondary opacity-25"></li>' +
'        <li>' +
'          <a class="dropdown-item text-danger" href="#" onclick="ManagementLayout.logout(); return false;">' +
'            <i class="fa-solid fa-right-from-bracket me-2"></i> Logout' +
'          </a>' +
'        </li>' +
'      </ul>' +
'    </div>' +
'  </div>' +
'</header>';

    function initLayout() {
        const toggle = document.getElementById('sidebarToggle');
        const overlay = document.getElementById('mobileOverlay');
        const body = document.body;

        if (toggle) {
            toggle.addEventListener('click', function () {
                if (window.innerWidth <= 768) {
                    body.classList.toggle('sidebar-mobile-open');
                    body.classList.remove('sidebar-collapsed');
                } else {
                    body.classList.toggle('sidebar-collapsed');
                    localStorage.setItem(
                        'sidebar_state',
                        body.classList.contains('sidebar-collapsed') ? 'collapsed' : 'expanded'
                    );
                }
            });
        }

        if (overlay) {
            overlay.addEventListener('click', function () {
                body.classList.remove('sidebar-mobile-open');
            });
        }
    }

    function initUserInfo() {
        const nameEl = document.getElementById('userName');
        const emailEl = document.getElementById('userEmail');
        const avatarEl = document.getElementById('navAvatar');
        const roleBadge = document.getElementById('roleBadge');

        if (!hasLocalUser) {
            if (roleBadge) roleBadge.textContent = 'SESSION';
            return;
        }

        if (nameEl && user.name) nameEl.textContent = user.name;
        if (emailEl && user.email) emailEl.textContent = user.email;
        if (avatarEl && user.name) {
            avatarEl.src = 'https://ui-avatars.com/api/?name=' + encodeURIComponent(user.name) + '&background=0b0e17&color=fff';
        }
        if (roleBadge) roleBadge.textContent = roleCode || 'MANAGEMENT';
    }

    function inject() {
        const sidebarEl = document.getElementById('ops-sidebar-placeholder');
        const navbarEl = document.getElementById('ops-navbar-placeholder');
        if (sidebarEl) sidebarEl.outerHTML = sidebarHTML;
        if (navbarEl) navbarEl.outerHTML = navbarHTML;
        initLayout();
        initUserInfo();
    }

    window.ManagementLayout = {
        logout: function () {
            fetch('/api/auth/logout', { method: 'POST', credentials: 'include' }).catch(function () {});
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            sessionStorage.clear();
            window.location.href = '/account/login';
        }
    };

    // Backward compatibility for any inline logout handler.
    window.logout = window.ManagementLayout.logout;

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', inject);
    } else {
        inject();
    }
})();
