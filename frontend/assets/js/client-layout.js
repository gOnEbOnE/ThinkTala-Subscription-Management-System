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
    function isDiscoverActive() { return path === '/client/discover' ? ' active' : ''; }
    // KYC form page (/client/kyc) should also highlight the KYC nav item
    function isKycActive() { return (path === '/client/kyc-status' || path === '/client/kyc') ? ' active' : ''; }
    function isMembershipActive() { return path === '/client/packages-catalog' ? ' active' : ''; }
    function isBillingActive() {
        return (path === '/client/billing-history' || path === '/client/order-detail' || path === '/client/checkout') ? ' active' : '';
    }
    function isSubscriptionActive() {
        return path === '/client/subscription-me' ? ' active' : '';
    }
    function isTicketActive() { return (path === '/client/support-tickets' || path === '/client/support-ticket-detail' || path === '/support/create' || path === '/client/support-create') ? ' active' : ''; }
    function isSettingsActive() { return path === '/client/settings' ? ' active' : ''; }

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
                '<li class="sidebar-section-label"><span class="link-text"></span></li>' +
                '<li class="nav-item"><a class="nav-link' + isActive('/client/dashboard') + '" href="/client/dashboard"><i class="fa-solid fa-chart-pie icon-left"></i><span class="link-text">Dashboard</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isDiscoverActive() + '" href="/client/discover"><i class="fa-solid fa-newspaper icon-left"></i><span class="link-text">Discover</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + ((path === '/market-insight' || path === '/client/market-insight') ? ' active' : '') + '" href="/client/market-insight"><i class="fa-solid fa-globe icon-left"></i><span class="link-text">Market Insight</span></a></li>' +
                '<li class="nav-item"><a class="nav-link disabled" href="#"><i class="fa-solid fa-satellite-dish icon-left"></i><span class="link-text">Deep Scanner</span><span class="badge bg-secondary ms-auto" style="font-size:.55rem">Soon</span></a></li>' +
                '<li class="nav-item"><a class="nav-link disabled" href="#"><i class="fa-solid fa-wand-magic-sparkles icon-left"></i><span class="link-text">Ask Nizza</span><span class="badge bg-secondary ms-auto" style="font-size:.55rem">Soon</span></a></li>' +
                '<li class="sidebar-section-label"><span class="link-text"></span></li>' +
                '<li class="nav-item"><a class="nav-link' + isKycActive() + '" href="/client/kyc-status"><i class="fa-solid fa-id-card icon-left"></i><span class="link-text">KYC Verification</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isTicketActive() + '" href="/client/support-tickets"><i class="fa-solid fa-ticket icon-left"></i><span class="link-text">Ticket</span></a></li>' +
            '</ul>' +
            '<ul class="nav flex-column mb-5">' +
                '<li class="nav-item"><a class="nav-link' + isMembershipActive() + '" href="/client/packages-catalog"><i class="fa-solid fa-crown icon-left"></i><span class="link-text">Membership</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isSubscriptionActive() + '" href="/client/subscription-me"><i class="fa-solid fa-repeat icon-left"></i><span class="link-text">Langganan</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isBillingActive() + '" href="/client/billing-history"><i class="fa-solid fa-receipt icon-left"></i><span class="link-text">Billing History</span></a></li>' +
                '<li class="nav-item"><a class="nav-link' + isSettingsActive() + '" href="/client/settings"><i class="fa-solid fa-gear icon-left"></i><span class="link-text">Pengaturan</span></a></li>' +
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
                '<button class="btn-header notif-bell" id="notifBellBtn" aria-label="Notifications">' +
                    '<i class="fa-regular fa-bell"></i>' +
                    '<span class="notif-badge" id="notifBadge" style="display:none;">0</span>' +
                '</button>' +
                '<button class="btn-header" id="themeToggle"><i class="fa-solid fa-moon"></i></button>' +
                '<div class="dropdown ms-1">' +
                    '<a href="#" class="d-flex align-items-center text-decoration-none" data-bs-toggle="dropdown">' +
                        '<div class="rounded-circle p-1 position-relative" style="border: 2px solid var(--accent-cyan);">' +
                            '<img src="https://ui-avatars.com/api/?name=User&background=0b0e17&color=fff" class="rounded-circle" width="34" height="34" id="avatarImg">' +
                        '</div>' +
                    '</a>' +
                    '<ul class="dropdown-menu dropdown-menu-end dropdown-menu-animate mt-3">' +
                        '<li class="px-3 py-2">' +
                            '<span class="d-block fw-bold text-main" id="userName">User</span>' +
                            '<small class="text-muted" id="userEmail">-</small>' +
                        '</li>' +
                        '<li><hr class="dropdown-divider border-secondary opacity-25"></li>' +
                        '<li><a class="dropdown-item text-danger" href="#" onclick="logout()"><i class="fa-solid fa-right-from-bracket me-2"></i> Logout</a></li>' +
                    '</ul>' +
                '</div>' +
            '</div>' +
        '</header>';

    var drawerHTML =
        '<div class="notif-overlay" id="notifOverlay"></div>' +
        '<aside class="notif-drawer" id="notifDrawer" aria-hidden="true">' +
            '<div class="notif-drawer-header">' +
                '<div>' +
                    '<h6 class="notif-drawer-title">Notifications</h6>' +
                    '<div class="notif-drawer-subtitle">Ringkasan update terbaru</div>' +
                '</div>' +
                '<button class="notif-close-btn" id="notifCloseBtn" aria-label="Close">' +
                    '<i class="fa-solid fa-xmark"></i>' +
                '</button>' +
            '</div>' +
            '<div class="notif-drawer-actions">' +
                '<button class="notif-mark-all" id="notifMarkAll">Mark all as read</button>' +
            '</div>' +
            '<div class="notif-drawer-body" id="notifDrawerBody">' +
                '<div class="notif-empty">Memuat notifikasi terbaru...</div>' +
            '</div>' +
            '<div class="notif-drawer-footer">' +
                '<a class="notif-see-all" href="/client/discover">See all notifications <i class="fa-solid fa-arrow-right"></i></a>' +
            '</div>' +
        '</aside>';

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

        if (!document.getElementById('notifDrawer')) {
            var tmp3 = document.createElement('div');
            tmp3.innerHTML = drawerHTML;
            while (tmp3.firstChild) { document.body.appendChild(tmp3.firstChild); }
        }
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', inject);
    } else {
        inject();
    }

    // ── Notification Drawer ─────────────────────────────────────
    var notifState = {
        items: [],
        unread: 0,
        loaded: false
    };

    var notifTypeColor = {
        system: '#b8c0ff',
        promo: '#88ff47',
        warning: '#ffb347',
        info: '#7dd3fc',
        analysis: '#2be8ff',
        education: '#6ee7b7',
        event: '#7ef29a'
    };

    function escapeHtml(value) {
        return String(value || '')
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }

    function formatDrawerTime(value) {
        if (!value) return '-';
        var date = new Date(value);
        if (Number.isNaN(date.getTime())) return '-';
        return date.toLocaleString('id-ID', {
            day: '2-digit',
            month: 'short',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function formatEventType(value) {
        if (!value) return 'EVENT';
        return String(value).replace(/_/g, ' ').toUpperCase();
    }

    function formatTypeLabel(item) {
        if (!item) return 'NEWS';
        if (item.source === 'event') return formatEventType(item.event_type);
        if (item.type) return String(item.type).toUpperCase();
        return 'NEWS';
    }

    function formatChannel(channel) {
        var map = { email: 'Email', whatsapp: 'WhatsApp', telegram: 'Telegram' };
        return map[channel] || channel || '';
    }

    function channelIcon(channel) {
        var map = { email: 'fa-envelope', whatsapp: 'fa-whatsapp', telegram: 'fa-paper-plane' };
        return map[channel] || 'fa-bell';
    }

    function resolveNotificationLink(item) {
        if (!item) return '/client/discover';
        if (item.source === 'news') {
            return '/client/discover#news/' + encodeURIComponent(item.id || '');
        }
        var eventType = String(item.event_type || '').toLowerCase();
        if (eventType.indexOf('kyc') !== -1) {
            return '/client/kyc-status';
        }
        if (eventType.indexOf('payment') !== -1) {
            return '/client/billing-history';
        }
        if (eventType.indexOf('otp') !== -1) {
            return '/account/verify-otp';
        }
        if (eventType.indexOf('password_reset') !== -1) {
            return '/account/login';
        }
        return '/client/discover';
    }

    function updateNotifBadge(count) {
        var badge = document.getElementById('notifBadge');
        if (!badge) return;
        if (!count || count <= 0) {
            badge.style.display = 'none';
            return;
        }
        badge.textContent = count > 9 ? '9+' : String(count);
        badge.style.display = '';
    }

    function setDrawerOpen(isOpen) {
        var drawer = document.getElementById('notifDrawer');
        var overlay = document.getElementById('notifOverlay');
        if (drawer) {
            drawer.classList.toggle('open', isOpen);
            drawer.setAttribute('aria-hidden', isOpen ? 'false' : 'true');
        }
        if (overlay) {
            overlay.classList.toggle('open', isOpen);
        }
        document.body.classList.toggle('notif-open', isOpen);
    }

    function renderNotificationList(list) {
        var container = document.getElementById('notifDrawerBody');
        if (!container) return;
        if (!list || list.length === 0) {
            container.innerHTML = '<div class="notif-empty">Tidak ada notifikasi baru.</div>';
            return;
        }

        var html = '<div class="notif-list">' + list.map(function (item) {
            var isUnread = !item.is_read;
            var label = formatTypeLabel(item);
            var color = notifTypeColor[String(item.type || '').toLowerCase()] || '#00f2ff';
            var time = formatDrawerTime(item.created_at);
            var title = item.title || label;
            var body = item.body || 'Tidak ada ringkasan.';
            var channel = item.source === 'event' ? formatChannel(item.channel) : '';
            var channelHtml = item.source === 'event'
                ? '<span class="notif-channel"><i class="fa-solid ' + channelIcon(item.channel) + '"></i>' + escapeHtml(channel) + '</span>'
                : '';
            return '' +
                '<button type="button" class="notif-item' + (isUnread ? ' is-unread' : '') + '"' +
                ' data-id="' + escapeHtml(item.id) + '" data-source="' + escapeHtml(item.source) + '"' +
                ' data-href="' + escapeHtml(resolveNotificationLink(item)) + '">' +
                (isUnread ? '<span class="notif-unread-dot"></span>' : '') +
                '<div class="notif-item-top">' +
                '<span class="notif-type" style="--notif-color:' + color + ';">' + escapeHtml(label) + '</span>' +
                '<span class="notif-time">' + escapeHtml(time) + '</span>' +
                '</div>' +
                '<div class="notif-title">' + escapeHtml(title) + '</div>' +
                '<div class="notif-body-text">' + escapeHtml(body) + '</div>' +
                '<div class="notif-meta">' + channelHtml + '</div>' +
                '</button>';
        }).join('') + '</div>';

        container.innerHTML = html;
    }

    function fetchRecentNotifications() {
        return fetch('/api/notifications/recent', {
            credentials: 'include',
            headers: {
                'Accept': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        }).then(function (res) {
            if (!res.ok) {
                return null;
            }
            return res.json();
        }).then(function (json) {
            var list = json && Array.isArray(json.data) ? json.data : [];
            notifState.items = list;
            notifState.unread = list.filter(function (item) {
                return !item.is_read;
            }).length;
            notifState.loaded = true;
            updateNotifBadge(notifState.unread);
            renderNotificationList(list);
            return list;
        }).catch(function () {
            if (!notifState.loaded) {
                renderNotificationList([]);
            }
        });
    }

    function markNotificationRead(item) {
        if (!item || item.is_read) {
            return Promise.resolve();
        }
        return fetch('/api/notifications/recent/read', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({ id: item.id, source_type: item.source })
        }).then(function (res) {
            if (!res.ok) throw new Error('Failed to mark notification as read');
            item.is_read = true;
            notifState.unread = notifState.items.filter(function (row) {
                return !row.is_read;
            }).length;
            updateNotifBadge(notifState.unread);
            renderNotificationList(notifState.items);
        }).catch(function (err) {
            console.error('Error marking notification read:', err);
        });
    }

    function markAllNotificationsRead() {
        return fetch('/api/notifications/recent/read-all', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            }
        }).then(function (res) {
            if (!res.ok) throw new Error('Failed to mark all as read');
            notifState.items = notifState.items.map(function (item) {
                item.is_read = true;
                return item;
            });
            notifState.unread = 0;
            updateNotifBadge(0);
            renderNotificationList(notifState.items);
        }).catch(function (err) {
            console.error('Error marking all notifications read:', err);
        });
    }

    function setupNotificationDrawer() {
        var bellBtn = document.getElementById('notifBellBtn');
        var overlay = document.getElementById('notifOverlay');
        var closeBtn = document.getElementById('notifCloseBtn');
        var markAllBtn = document.getElementById('notifMarkAll');
        var drawerBody = document.getElementById('notifDrawerBody');

        if (bellBtn) {
            bellBtn.addEventListener('click', function () {
                setDrawerOpen(true);
                fetchRecentNotifications();
            });
        }

        if (overlay) {
            overlay.addEventListener('click', function () {
                setDrawerOpen(false);
            });
        }

        if (closeBtn) {
            closeBtn.addEventListener('click', function () {
                setDrawerOpen(false);
            });
        }

        if (markAllBtn) {
            markAllBtn.addEventListener('click', function () {
                markAllNotificationsRead();
            });
        }

        if (drawerBody) {
            drawerBody.addEventListener('click', function (event) {
                var target = event.target;
                while (target && target !== drawerBody) {
                    if (target.classList && target.classList.contains('notif-item')) {
                        var id = target.getAttribute('data-id');
                        var source = target.getAttribute('data-source');
                        var href = target.getAttribute('data-href');
                        var item = notifState.items.find(function (row) {
                            return row.id === id && row.source === source;
                        });
                        markNotificationRead(item).finally(function () {
                            setDrawerOpen(false);
                            if (href) {
                                window.location.href = href;
                            }
                        });
                        break;
                    }
                    target = target.parentNode;
                }
            });
        }

        document.addEventListener('keydown', function (event) {
            if (event.key === 'Escape') {
                setDrawerOpen(false);
            }
        });

        fetchRecentNotifications();
    }

    // ── Sidebar toggle ────────────────────────────────────────────
    document.addEventListener('DOMContentLoaded', function () {
        // Discover badge counter intentionally disabled.
        delete window.setClientDiscoverCount;

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

        // Simpan user_id ke localStorage agar halaman lain (settings, dll.) bisa pakai
        try {
            var _u = JSON.parse(localStorage.getItem('user') || '{}');
            if (_u.id && !localStorage.getItem('user_id')) {
                localStorage.setItem('user_id', _u.id);
            }
            // Expose ke window untuk getUserID() di settings.html
            window.__currentUserID = _u.id || '';
        } catch (e) { /* ignore */ }

        setupNotificationDrawer();
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
