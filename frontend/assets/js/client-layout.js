(function () {
    // ── Route guard: redirect to login if no session ─────────────
    var guardUser = null;
    try { guardUser = JSON.parse(localStorage.getItem('user')); } catch (e) {}
    if (!guardUser || !guardUser.id) {
        window.location.href = '/account/login';
        return;
    }
    // ── Role guard: /client/* is only for CLIENT, CEO, SUPERADMIN ──────────
    var _clientRole = (guardUser.role_code || guardUser.level_code || guardUser.level || '').toString().toUpperCase();
    if (_clientRole !== 'CLIENT' && _clientRole !== 'CEO' && _clientRole !== 'SUPERADMIN') {
        var _clientRedirect = { 'OPERASIONAL': '/ops/dashboard', 'COMPLIANCE': '/compliance/dashboard' };
        window.location.href = _clientRedirect[_clientRole] || '/account/login';
        return;
    }
    // ── Prevent transition flash on load ──────────────────────────
    var css = document.createElement('style');
    css.id = 'prevent-tx';
    css.appendChild(document.createTextNode('* { transition: none !important; }'));
    document.head.appendChild(css);
    window.addEventListener('load', function () {
        setTimeout(function () { var el = document.getElementById('prevent-tx'); if (el) el.remove(); }, 50);
    });

    // ── Restore sidebar collapsed state ───────────────────────────
    if (window.innerWidth > 768 && localStorage.getItem('sidebar_state') === 'collapsed') {
        document.body.classList.add('sidebar-collapsed');
    }

    // ── Active state detection ────────────────────────────────────
    var path = window.location.pathname;
    function isActive(route) { return path === route ? ' active' : ''; }
    // KYC form page (/client/kyc) should also highlight the KYC nav item
    function isKycActive() { return (path === '/client/kyc-status' || path === '/client/kyc') ? ' active' : ''; }
    function isMembershipActive() { return path === '/client/packages-catalog' ? ' active' : ''; }
    function isTicketActive() { return (path === '/support/create' || path === '/client/support-create') ? ' active' : ''; }

    // ── Sidebar HTML ──────────────────────────────────────────────
    var sidebarHTML =
        '<div class="mobile-overlay" id="mobileOverlay"></div>' +
        '<nav class="sidebar">' +
            '<div class="sidebar-brand mb-3">' +
                '<div class="brand-wrapper">' +
                    '<i class="fa-solid fa-layer-group text-cyan brand-icon"></i>' +
                    '<div class="brand-text-content">' +
                        '<h4 class="fw-bold tracking-wider mb-0" style="color: var(--text-heading)">Think<span class="text-cyan">Tala</span></h4>' +
                        '<p class="small mb-0 text-muted" style="font-size: 0.7rem;">Client Portal</p>' +
                    '</div>' +
                '</div>' +
            '</div>' +
            '<ul class="nav flex-column flex-grow-1">' +
                '<li class="nav-item"><a class="nav-link' + isActive('/client/dashboard') + '" href="/client/dashboard"><i class="fa-solid fa-chart-pie icon-left"></i><span class="link-text">Dashboard</span></a></li>' +
                '<li class="nav-item"><a class="nav-link disabled" href="#"><i class="fa-solid fa-globe icon-left"></i><span class="link-text">Market Insight</span><span class="badge bg-secondary ms-auto" style="font-size:.55rem">Soon</span></a></li>' +
                '<li class="nav-item"><a class="nav-link disabled" href="#"><i class="fa-solid fa-satellite-dish icon-left"></i><span class="link-text">Deep Scanner</span><span class="badge bg-secondary ms-auto" style="font-size:.55rem">Soon</span></a></li>' +
                '<li class="nav-item"><a class="nav-link disabled" href="#"><i class="fa-solid fa-wand-magic-sparkles icon-left"></i><span class="link-text">Ask Nizza</span><span class="badge bg-secondary ms-auto" style="font-size:.55rem">Soon</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isKycActive() + '" href="/client/kyc-status"><i class="fa-solid fa-id-card icon-left"></i><span class="link-text">KYC Verification</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isTicketActive() + '" href="/support/create"><i class="fa-solid fa-ticket icon-left"></i><span class="link-text">Tiket</span></a></li>' +
            '</ul>' +
            '<ul class="nav flex-column mb-5">' +
                '<li class="nav-item"><a class="nav-link' + isMembershipActive() + '" href="/client/packages-catalog"><i class="fa-solid fa-crown icon-left"></i><span class="link-text">Membership</span></a></li>' +
                '<li class="nav-item"><a class="nav-link disabled" href="#"><i class="fa-solid fa-gear icon-left"></i><span class="link-text">Settings</span><span class="badge bg-secondary ms-auto" style="font-size:.55rem">Soon</span></a></li>' +
                '<li class="nav-item"><a class="nav-link text-danger" href="#" onclick="logout()"><i class="fa-solid fa-right-from-bracket icon-left"></i><span class="link-text">Logout</span></a></li>' +
            '</ul>' +
        '</nav>';

    // ── Navbar HTML ───────────────────────────────────────────────
    var navbarHTML =
        '<header class="top-navbar d-flex justify-content-between align-items-center">' +
            '<div class="d-flex align-items-center gap-2">' +
                '<button class="btn-header" id="sidebarToggle"><i class="fa-solid fa-bars fa-lg"></i></button>' +
            '</div>' +
            '<div class="nizza-wrapper mx-auto d-none d-md-block">' +
                '<div class="nizza-search">' +
                    '<i class="fa-solid fa-wand-magic-sparkles magic-icon"></i>' +
                    '<input type="text" class="form-control" placeholder="Ask Nizza AI about market trends...">' +
                    '<button class="btn btn-sm text-muted"><i class="fa-solid fa-arrow-right"></i></button>' +
                '</div>' +
            '</div>' +
            '<div class="d-flex align-items-center gap-2 gap-md-3">' +
                '<button class="btn-header" id="themeToggle"><i class="fa-solid fa-moon"></i></button>' +
                '<div class="dropdown ms-1">' +
                    '<a href="#" class="d-flex align-items-center text-decoration-none" data-bs-toggle="dropdown">' +
                        '<div class="rounded-circle p-1 position-relative" style="border: 2px solid var(--accent-cyan);">' +
                            '<img src="https://ui-avatars.com/api/?name=Client&background=0b0e17&color=fff" class="rounded-circle" width="34" height="34" id="avatarImg">' +
                        '</div>' +
                    '</a>' +
                    '<ul class="dropdown-menu dropdown-menu-end dropdown-menu-animate mt-3">' +
                        '<li class="px-3 py-2">' +
                            '<span class="d-block fw-bold text-main" id="userName">Client</span>' +
                            '<small class="text-muted" id="userEmail">client@thinktala.com</small>' +
                        '</li>' +
                        '<li><hr class="dropdown-divider border-secondary opacity-25"></li>' +
                        '<li><a class="dropdown-item text-danger" href="#" onclick="logout()"><i class="fa-solid fa-right-from-bracket me-2"></i> Logout</a></li>' +
                    '</ul>' +
                '</div>' +
            '</div>' +
        '</header>';

    // ── Inject into placeholders ──────────────────────────────────
    function inject() {
        var sp = document.getElementById('client-sidebar-placeholder');
        var np = document.getElementById('client-navbar-placeholder');

        if (sp) {
            var tmp = document.createElement('div');
            tmp.innerHTML = sidebarHTML;
            while (tmp.firstChild) { sp.parentNode.insertBefore(tmp.firstChild, sp); }
            sp.remove();
        }

        if (np) {
            var tmp2 = document.createElement('div');
            tmp2.innerHTML = navbarHTML;
            np.parentNode.replaceChild(tmp2.firstChild, np);
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', inject);
    } else {
        inject();
    }

    // ── Sidebar toggle ────────────────────────────────────────────
    document.addEventListener('DOMContentLoaded', function () {
        var btn = document.getElementById('sidebarToggle');
        if (btn) {
            btn.addEventListener('click', function () {
                document.body.classList.toggle('sidebar-collapsed');
                localStorage.setItem('sidebar_state',
                    document.body.classList.contains('sidebar-collapsed') ? 'collapsed' : 'expanded');
            });
        }

        // Populate user info from localStorage
        try {
            var user = JSON.parse(localStorage.getItem('user') || '{}');
            if (user.name) {
                var uName = document.getElementById('userName');
                var avatar = document.getElementById('avatarImg');
                if (uName) uName.textContent = user.name;
                if (avatar) avatar.src = 'https://ui-avatars.com/api/?name=' + encodeURIComponent(user.name) + '&background=0b0e17&color=fff';
            }
            if (user.email) {
                var uEmail = document.getElementById('userEmail');
                if (uEmail) uEmail.textContent = user.email;
            }
        } catch (e) { /* ignore */ }
    });

    // ── Default logout ─────────────────────────────────────────────
    window.logout = function () {
        fetch('/api/auth/logout', { method: 'POST', credentials: 'include' }).catch(function () {});
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        sessionStorage.clear();
        window.location.href = '/account/login';
    };
})();
