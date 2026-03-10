(function () {
    // ── Route guard: redirect to login if no session ─────────────
    var user = null;
    try { user = JSON.parse(localStorage.getItem('user')); } catch (e) {}
    if (!user || !user.id) {
        window.location.href = '/account/login';
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

    // ── Sidebar HTML ──────────────────────────────────────────────
    var sidebarHTML =
        '<div class="mobile-overlay" id="mobileOverlay"></div>' +
        '<nav class="sidebar">' +
            '<div class="sidebar-brand mb-3">' +
                '<div class="brand-wrapper">' +
                    '<i class="fa-solid fa-layer-group text-cyan brand-icon"></i>' +
                    '<div class="brand-text-content">' +
                        '<h4 class="fw-bold tracking-wider mb-0" style="color: var(--text-heading)">Think<span class="text-cyan">Tala</span></h4>' +
                        '<p class="small mb-0 text-muted" style="font-size: 0.7rem;">Compliance Panel</p>' +
                    '</div>' +
                '</div>' +
            '</div>' +
            '<ul class="nav flex-column flex-grow-1">' +
                '<li class="nav-item"><a class="nav-link' + isActive('/compliance/dashboard') + '" href="/compliance/dashboard"><i class="fa-solid fa-shield-halved icon-left"></i><span class="link-text">KYC Review</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isActive('/compliance/reports') + '" href="/compliance/reports"><i class="fa-solid fa-chart-bar icon-left"></i><span class="link-text">Reports</span></a></li>' +
            '</ul>' +
            '<ul class="nav flex-column mb-5">' +
                '<li class="nav-item"><a class="nav-link text-danger" href="#" onclick="logout()"><i class="fa-solid fa-right-from-bracket icon-left"></i><span class="link-text">Logout</span></a></li>' +
            '</ul>' +
        '</nav>';

    // ── Navbar HTML ───────────────────────────────────────────────
    var navbarHTML =
        '<header class="top-navbar d-flex justify-content-between align-items-center">' +
            '<div class="d-flex align-items-center gap-2">' +
                '<button class="btn-header" id="sidebarToggle"><i class="fa-solid fa-bars fa-lg"></i></button>' +
                '<span class="badge bg-danger" style="font-size: 0.7rem;">COMPLIANCE</span>' +
                '<span class="badge bg-info text-dark" id="assumedRoleBadge" style="display:none;"><i class="fa-solid fa-user-secret me-1"></i>Sedang sebagai: <strong id="assumedRoleName"></strong></span>' +
            '</div>' +
            '<div class="d-flex align-items-center gap-2">' +
                '<div class="dropdown" id="assumeRoleSection" style="display:none;">' +
                    '<button class="btn btn-sm btn-outline-light dropdown-toggle" data-bs-toggle="dropdown" style="font-size:0.8rem;"><i class="fa-solid fa-user-secret me-1"></i>Simulasi Peran</button>' +
                    '<ul class="dropdown-menu dropdown-menu-end dropdown-menu-animate mt-2">' +
                        '<li><h6 class="dropdown-header">Pilih Simulasi Peran</h6></li>' +
                        '<li><a class="dropdown-item" href="#" onclick="assumeRole(\"OPERASIONAL\")"><i class="fa-solid fa-cogs me-2"></i>Operasional</a></li>' +
                        '<li><a class="dropdown-item" href="#" onclick="assumeRole(\"COMPLIANCE\")"><i class="fa-solid fa-shield-halved me-2"></i>Compliance</a></li>' +
                        '<li><a class="dropdown-item" href="#" onclick="assumeRole(\"CEO\")"><i class="fa-solid fa-briefcase me-2"></i>CEO</a></li>' +
                        '<li><a class="dropdown-item" href="#" onclick="assumeRole(\"CLIENT\")"><i class="fa-solid fa-user me-2"></i>Client</a></li>' +
                    '</ul>' +
                '</div>' +
                '<button class="btn-header" id="themeToggle"><i class="fa-solid fa-moon"></i></button>' +
                '<div class="dropdown ms-1">' +
                    '<a href="#" class="d-flex align-items-center text-decoration-none" data-bs-toggle="dropdown">' +
                        '<div class="rounded-circle p-1" style="border: 2px solid var(--accent-cyan);">' +
                            '<img src="https://ui-avatars.com/api/?name=CM&background=0b0e17&color=fff" class="rounded-circle" width="34" height="34" id="avatarImg">' +
                        '</div>' +
                    '</a>' +
                    '<ul class="dropdown-menu dropdown-menu-end dropdown-menu-animate mt-3">' +
                        '<li class="px-3 py-2">' +
                            '<span class="d-block fw-bold text-main" id="userName">Compliance</span>' +
                            '<small class="text-muted" id="userEmail">compliance@thinktala.com</small>' +
                        '</li>' +
                        '<li><hr class="dropdown-divider border-secondary opacity-25"></li>' +
                        '<li><a class="dropdown-item text-danger" href="#" onclick="logout()"><i class="fa-solid fa-right-from-bracket me-2"></i>Logout</a></li>' +
                    '</ul>' +
                '</div>' +
            '</div>' +
        '</header>';

    // ── Inject into placeholders ──────────────────────────────────
    function inject() {
        var sp = document.getElementById('compliance-sidebar-placeholder');
        var np = document.getElementById('compliance-navbar-placeholder');

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
        // Populate user info from localStorage
        try {
            if (user && user.name) {
                var uName = document.getElementById('userName');
                var avatar = document.getElementById('avatarImg');
                if (uName) uName.textContent = user.name;
                if (avatar) avatar.src = 'https://ui-avatars.com/api/?name=' + encodeURIComponent(user.name) + '&background=0b0e17&color=fff';
            }
            if (user && user.email) {
                var uEmail = document.getElementById('userEmail');
                if (uEmail) uEmail.textContent = user.email;
            }
        } catch (e) { /* ignore */ }

        var btn = document.getElementById('sidebarToggle');
        if (btn) {
            btn.addEventListener('click', function () {
                document.body.classList.toggle('sidebar-collapsed');
                localStorage.setItem('sidebar_state',
                    document.body.classList.contains('sidebar-collapsed') ? 'collapsed' : 'expanded');
            });
        }

        // Show assume role section for SUPERADMIN
        if ((user.level_code || '').toUpperCase() === 'SUPERADMIN') {
            var section = document.getElementById('assumeRoleSection');
            if (section) section.style.display = '';
        }
        // Show assumed role badge
        if (user.assumed_role) {
            var badge = document.getElementById('assumedRoleBadge');
            var nameSpan = document.getElementById('assumedRoleName');
            if (badge && nameSpan) { nameSpan.textContent = user.role_code || ''; badge.style.display = ''; }
        }
    });

    // ── Assume role ───────────────────────────────────────────────
    window.assumeRole = async function (targetRoleCode) {
        try {
            var res = await fetch('/api/auth/assume-role', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ target_role_code: targetRoleCode })
            });
            var json = await res.json();
            if (json.success || json.status) {
                var u = JSON.parse(localStorage.getItem('user') || '{}');
                u.assumed_role = true; u.role_code = targetRoleCode;
                localStorage.setItem('user', JSON.stringify(u));
                window.location.href = (json.data && json.data.redirect_url) || '/ops/dashboard';
            } else { alert(json.message || json.msg || 'Gagal simulasi role'); }
        } catch (e) { alert('Gagal menghubungi server'); }
    };

    // ── Default logout — same across all compliance pages ─────────
    window.logout = function () {
        fetch('/api/auth/logout', { method: 'POST', credentials: 'include' }).catch(function () {});
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        sessionStorage.clear();
        window.location.href = '/account/login';
    };
})();
