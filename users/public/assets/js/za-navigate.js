(function () {
    let idleTimer;
    const IDLE_TIMEOUT = 15 * 60 * 1000; // 15 Menit dalam milidetik

    function resetIdleTimer() {
        clearTimeout(idleTimer);
        // Jika dalam 15 menit tidak ada aktivitas, jalankan logout
        idleTimer = setTimeout(logoutUser, IDLE_TIMEOUT);
    }

    function logoutUser() {
        console.warn("Sesi berakhir karena tidak ada aktivitas.");
        // Arahkan ke endpoint logout backend Anda
        window.location.href = '/account/logout?reason=idle';
    }

    // Daftar aktivitas yang dianggap sebagai "User Masih Aktif"
    const activityEvents = ['mousedown', 'mousemove', 'keydown', 'scroll', 'touchstart'];

    // Pasang listener ke dokumen
    activityEvents.forEach(event => {
        document.addEventListener(event, resetIdleTimer, true);
    });

    // Jalankan timer pertama kali
    resetIdleTimer();
    // Cache sederhana agar tidak terus-terusan hit server untuk halaman yang sama
    const PageCache = new Map();

    // window.navigateTo = async function(pageKey) {
    //     // Normalisasi pageKey jika kosong
    //     if (!pageKey || pageKey === '') {
    //         pageKey = 'dashboard';
    //     }

    //     const wrapper = document.getElementById('wrapper');
    //     if (!wrapper) return;

    //     // 1. Tampilkan Spinner
    //     wrapper.innerHTML = `
    //         <div id="page-loader" class="text-center animate-page-in">
    //             <div class="spinner-border text-primary" role="status"></div>
    //         </div>`;

    //     try {
    //         let htmlContent = "";

    //         // 2. Cek Cache atau Ambil dari Backend
    //         if (PageCache.has(pageKey)) {
    //             htmlContent = PageCache.get(pageKey);
    //         } else {
    //             // Fetch ke endpoint Golang Anda
    //             const response = await fetch(`/account/${pageKey}`);
                
    //             // Cek jika sesi habis (biasanya redirect ke login)
    //             if (response.redirected && response.url.includes('/login')) {
    //                 window.location.href = '/login';
    //                 return;
    //             }

    //             if (!response.ok) throw new Error("Gagal mengambil halaman");
                
    //             htmlContent = await response.text();
    //             PageCache.set(pageKey, htmlContent); // Simpan di cache
    //         }

    //         // 3. PROSES AMAN: DOMParser + Penanganan Fragment
    //         const parser = new DOMParser();
    //         const doc = parser.parseFromString(htmlContent, 'text/html');
            
    //         // Ambil konten dari body template
    //         const newContent = doc.body.innerHTML; 

    //         // 4. Injeksi dengan Efek Transisi
    //         setTimeout(() => {
    //             // Tambahkan wrapper animasi agar transisi smooth
    //             wrapper.innerHTML = `<div class="animate-page-in">${newContent}</div>`;
                
    //             // Update Link Aktif di Sidebar
    //             updateActiveLink(pageKey);
                
    //             // Update URL tanpa reload (menggunakan Hash)
    //             if(window.location.hash !== `#${pageKey}`) {
    //                 history.pushState({page: pageKey}, "", `#${pageKey}`);
    //             }
    //         }, 300);

    //     } catch (error) {
    //         console.error(error);
    //         wrapper.innerHTML = `
    //             <div class="p-5 text-center text-danger animate-page-in">
    //                 <i class="bi bi-exclamation-triangle fs-1"></i>
    //                 <p class="mt-3">Gagal memuat halaman. Pastikan Anda masih terhubung atau silakan login kembali.</p>
    //             </div>`;
    //     }
    // };

    /**
 * Konfigurasi Keamanan:
 * 1. Mengambil Nonce dari meta tag untuk kepatuhan CSP.
 * 2. Memfilter script tanpa data-safe="true".
 */
    // const cspNonce = document.querySelector("meta[name='csp-nonce']")?.getAttribute("content");

    function runScripts(container, parsedDoc) {
        const scripts = parsedDoc.querySelectorAll('script');

        scripts.forEach((oldScript) => {
            // 1. FILTER: Cek atribut data-safe
            if (oldScript.getAttribute('data-safe') !== 'true') {
                console.warn("Blocked unsafe script:", oldScript);
                return;
            }

            const newScript = document.createElement('script');

            // 2. CLONE: Salin atribut (src, type, dll)
            Array.from(oldScript.attributes).forEach(attr => {
                newScript.setAttribute(attr.name, attr.value);
            });

            // 3. NONCE: Tambahkan nonce agar valid di mata browser (CSP)
            // if (cspNonce) {
            //     newScript.setAttribute('nonce', cspNonce);
            // }

            // 4. CONTENT: Salin isi inline script
            if (oldScript.innerHTML) {
                newScript.appendChild(document.createTextNode(oldScript.innerHTML));
            }

            // 5. EXECUTE: Masukkan ke DOM
            container.appendChild(newScript);
        });
    }

    window.navigateTo = async function(pageKey) {
        if (!pageKey || pageKey === '') pageKey = 'dashboard';

        const wrapper = document.getElementById('wrapper');
        if (!wrapper) return;

        wrapper.innerHTML = `
            <div id="page-loader" class="text-center animate-page-in">
                <div class="spinner-border text-primary" role="status"></div>
            </div>`;

        try {
            let htmlContent = "";

            if (PageCache.has(pageKey)) {
                htmlContent = PageCache.get(pageKey);
            } else {
                const response = await fetch(`/account/${pageKey}`);
                
                if (response.redirected && response.url.includes('/login')) {
                    window.location.href = '/login';
                    return;
                }

                if (!response.ok) throw new Error("Gagal mengambil halaman");
                htmlContent = await response.text();
                PageCache.set(pageKey, htmlContent);
            }

            const parser = new DOMParser();
            const doc = parser.parseFromString(htmlContent, 'text/html');
            const contentBody = doc.body;

            setTimeout(() => {
                wrapper.innerHTML = ""; 
                
                const contentDiv = document.createElement('div');
                contentDiv.className = "animate-page-in";
                
                // Masukkan HTML visual (Tanpa script tag, script ditangani runScripts)
                Array.from(contentBody.childNodes).forEach(node => {
                    if (node.nodeName !== 'SCRIPT') {
                        contentDiv.appendChild(node.cloneNode(true));
                    }
                });

                wrapper.appendChild(contentDiv);

                // Jalankan hanya script yang aman
                runScripts(contentDiv, doc);

                if (typeof updateActiveLink === 'function') updateActiveLink(pageKey);
                
                if(window.location.hash !== `#${pageKey}`) {
                    history.pushState({page: pageKey}, "", `#${pageKey}`);
                }
            }, 300);

        } catch (error) {
            console.error(error);
            wrapper.innerHTML = `<div class="p-5 text-center text-danger">Error loading page</div>`;
        }
    };

    window.updateActiveLink = function(pageKey) {
        document.querySelectorAll('.sidebar .nav-link').forEach(link => {
            link.classList.remove('active');
            
            const action = link.getAttribute('onclick');
            if (action && action.includes(`'${pageKey}'`)) {
                link.classList.add('active');
            }

            // Otomatis tutup sidebar di mobile setelah klik
            if (window.innerWidth < 768) {
                document.body.classList.remove('sidebar-open');
            }
        });
    }

    // ================================================================
    // INITIALIZER: Agar Dashboard/Halaman Aktif Terbuka Saat Start
    // ================================================================
    
    document.addEventListener('DOMContentLoaded', () => {
        // Ambil halaman dari hash (#) URL, jika tidak ada default ke dashboard
        const initialPage = window.location.hash.replace('#', '') || 'dashboard';
        window.navigateTo(initialPage);
    });

    // Menangani tombol Back/Forward di Browser
    window.onpopstate = (event) => {
        const page = (event.state && event.state.page) || 
                     window.location.hash.replace('#', '') || 
                     'dashboard';
        window.navigateTo(page);
    };

})();