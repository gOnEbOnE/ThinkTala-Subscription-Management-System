const sidebarToggle = document.getElementById('sidebarToggle');
const mobileOverlay = document.getElementById('mobileOverlay');
const body = document.body;

// Load Saved Sidebar State
const savedSidebar = localStorage.getItem('sidebar_state');
if (window.innerWidth > 768 && savedSidebar === 'collapsed') {
    body.classList.add('sidebar-collapsed');
}

// Sidebar Toggle Logic
sidebarToggle.addEventListener('click', () => {
    if (window.innerWidth <= 768) { 
        body.classList.toggle('sidebar-mobile-open'); 
        body.classList.remove('sidebar-collapsed'); 
    } else { 
        body.classList.toggle('sidebar-collapsed');
        // Save State
        if (body.classList.contains('sidebar-collapsed')) {
            localStorage.setItem('sidebar_state', 'collapsed');
        } else {
            localStorage.setItem('sidebar_state', 'expanded');
        }
    }
});

mobileOverlay.addEventListener('click', () => { 
    body.classList.remove('sidebar-mobile-open'); 
});