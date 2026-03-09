/* =====================================================
   OPS PANEL — Shared Layout (Sidebar + Navbar)
   Usage: include this script in any /ops/* page.

   The page must have:
     <div id="ops-sidebar-placeholder"></div>
     <div id="ops-navbar-placeholder"></div>

   Active state is automatically set from window.location.pathname.
   ===================================================== */

(function () {
    'use strict';

    // ── Prevent transition flash on load ────────────────────────────
    const preventTx = document.createElement('style');
    preventTx.id = 'prevent-tx';
    preventTx.innerHTML = '* { transition: none !important; }';
    document.head.appendChild(preventTx);
    window.addEventListener('load', () =>
        setTimeout(() => document.getElementById('prevent-tx')?.remove(), 50)
    );

    // ── Restore collapsed state before render ───────────────────────
    if (window.innerWidth > 768 && localStorage.getItem('sidebar_state') === 'collapsed') {
        document.body.classList.add('sidebar-collapsed');
    }

    // ── Route → active keys mapping ─────────────────────────────────
    // activeKey  : nav link that gets .active
    // parentKey  : submenu parent that gets .open + .parent-active
    const ROUTES = {
        '/ops/dashboard'              : { activeKey: 'dashboard' },
        '/ops/notifications'          : { activeKey: 'notifications',          parentKey: 'notif' },
        '/ops/notification-templates' : { activeKey: 'notification-templates', parentKey: 'notif' },
        '/ops/subscriptions'          : { activeKey: 'subscriptions' },
        '/ops/subscriptions-create'   : { activeKey: 'subscriptions' },
        '/ops/subscriptions-edit'     : { activeKey: 'subscriptions' },
    };

    const currentPath = window.location.pathname;
    const route = ROUTES[currentPath] || {};
    const activeKey  = route.activeKey  || '';
    const parentKey  = route.parentKey  || '';

    // Helper: mark a link active
    function isActive(key) {
        return activeKey === key ? ' active' : '';
    }
    // Helper: open submenu parent
    function isParentActive(key) {
        return parentKey === key ? ' open parent-active' : '';
    }
    // Helper: open submenu list
    function isSubmenuOpen(key) {
        return parentKey === key ? ' open' : '';
    }

    // ── Sidebar HTML ─────────────────────────────────────────────────
    const sidebarHTML = /* html */`
<nav class="sidebar">
    <div class="sidebar-brand mb-3">
        <div class="brand-wrapper">
            <i class="fa-solid fa-layer-group text-cyan brand-icon"></i>
            <div class="brand-text-content">
                <h4 class="fw-bold tracking-wider mb-0" style="color: var(--text-heading)">
                    Think<span class="text-cyan">Tala</span>
                </h4>
                <p class="small mb-0 text-muted" style="font-size: 0.7rem;">Operations Panel</p>
            </div>
        </div>
    </div>

    <ul class="nav flex-column flex-grow-1">

        <li class="nav-item">
            <a class="nav-link${isActive('dashboard')}" href="/ops/dashboard">
                <i class="fa-solid fa-chart-pie icon-left"></i>
                <span class="link-text">Dashboard</span>
            </a>
        </li>

        <!-- Notification (expandable) -->
        <li class="nav-item nav-item-group">
            <a class="nav-link nav-link-parent${isParentActive('notif')}"
               onclick="OpsLayout.toggleSubmenu(this)">
                <i class="fa-solid fa-bell icon-left"></i>
                <span class="link-text">Notification</span>
                <i class="fa-solid fa-chevron-right caret link-text"></i>
            </a>
            <ul class="nav-submenu${isSubmenuOpen('notif')}" id="notifSubmenu">
                <li>
                    <a class="nav-link${isActive('notifications')}" href="/ops/notifications">
                        <i class="fa-solid fa-list me-2" style="font-size:.85rem"></i>
                        <span>Monitoring</span>
                    </a>
                </li>
                <li>
                    <a class="nav-link${isActive('notification-templates')}" href="/ops/notification-templates">
                        <i class="fa-solid fa-file-alt me-2" style="font-size:.85rem"></i>
                        <span>Template Management</span>
                    </a>
                </li>
            </ul>
        </li>

        <li class="nav-item">
            <a class="nav-link${isActive('subscriptions')}" href="/ops/subscriptions">
                <i class="fa-solid fa-crown icon-left"></i>
                <span class="link-text">Subscriptions</span>
            </a>
        </li>

    </ul>

    <ul class="nav flex-column mb-5">
        <li class="nav-item">
            <a class="nav-link" href="/ops/settings">
                <i class="fa-solid fa-gear icon-left"></i>
                <span class="link-text">Settings</span>
            </a>
        </li>
        <li class="nav-item">
            <a class="nav-link text-danger" href="#" onclick="OpsLayout.logout()">
                <i class="fa-solid fa-right-from-bracket icon-left"></i>
                <span class="link-text">Logout</span>
            </a>
        </li>
    </ul>
</nav>
<div class="mobile-overlay" id="mobileOverlay"></div>
`;

    // ── Navbar HTML ──────────────────────────────────────────────────
    const navbarHTML = /* html */`
<header class="top-navbar d-flex justify-content-between align-items-center">
    <div class="d-flex align-items-center gap-2">
        <button class="btn-header" id="sidebarToggle">
            <i class="fa-solid fa-bars fa-lg"></i>
        </button>
        <span class="badge bg-warning text-dark">OPERASIONAL</span>
    </div>
    <div class="d-flex align-items-center gap-2 gap-md-3">
        <button class="btn-header" id="themeToggle">
            <i class="fa-solid fa-moon"></i>
        </button>
        <div class="dropdown ms-1">
            <a href="#" class="d-flex align-items-center text-decoration-none"
               data-bs-toggle="dropdown">
                <div class="rounded-circle p-1" style="border: 2px solid var(--accent-cyan);">
                    <img src="https://ui-avatars.com/api/?name=OPS&background=0b0e17&color=fff"
                         class="rounded-circle" id="navAvatar" width="34" height="34">
                </div>
            </a>
            <ul class="dropdown-menu dropdown-menu-end dropdown-menu-animate mt-3">
                <li class="px-3 py-2">
                    <span class="d-block fw-bold text-main" id="userName">Staff Operasional</span>
                    <small class="text-muted" id="userEmail">ops@thinktala.com</small>
                </li>
                <li>
                    <hr class="dropdown-divider border-secondary opacity-25">
                </li>
                <li>
                    <a class="dropdown-item text-danger" href="#" onclick="OpsLayout.logout()">
                        <i class="fa-solid fa-right-from-bracket me-2"></i> Logout
                    </a>
                </li>
            </ul>
        </div>
    </div>
</header>
`;

    // ── Inject into placeholders ─────────────────────────────────────
    function inject() {
        const sidebarEl = document.getElementById('ops-sidebar-placeholder');
        const navbarEl  = document.getElementById('ops-navbar-placeholder');
        if (sidebarEl) sidebarEl.outerHTML = sidebarHTML;
        if (navbarEl)  navbarEl.outerHTML  = navbarHTML;
        initLayout();
    }

    // ── Init sidebar toggle & theme toggle ───────────────────────────
    function initLayout() {
        const toggle  = document.getElementById('sidebarToggle');
        const overlay = document.getElementById('mobileOverlay');
        const body    = document.body;

        if (toggle) {
            toggle.addEventListener('click', () => {
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
            overlay.addEventListener('click', () => {
                body.classList.remove('sidebar-mobile-open');
            });
        }

        // Theme toggle is handled by za-theme.js — just ensure the button exists
        // (za-theme.js binds to #themeToggle on DOMContentLoaded)
    }

    // ── Public API ───────────────────────────────────────────────────
    window.OpsLayout = {
        toggleSubmenu(el) {
            const parent  = el.closest('.nav-item-group');
            const submenu = parent.querySelector('.nav-submenu');
            el.classList.toggle('open');
            submenu.classList.toggle('open');
        },
        logout() {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/account/login';
        }
    };

    // Also expose as globals for backward-compat inline handlers
    window.toggleSubmenu = window.OpsLayout.toggleSubmenu;
    window.logout        = window.OpsLayout.logout;

    // ── Run on DOM ready ─────────────────────────────────────────────
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', inject);
    } else {
        inject();
    }

})();
