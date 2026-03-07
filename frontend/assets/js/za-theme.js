// Menangani Tema (Dark/Light)
const themeToggleBtn = document.getElementById('themeToggle');
const htmlElement = document.documentElement;
const themeIcon = themeToggleBtn.querySelector('i');

// Load Saved Theme
const savedTheme = localStorage.getItem('site_theme') || 'dark';
htmlElement.setAttribute('data-theme', savedTheme);
if (themeIcon) themeIcon.className = savedTheme === 'light' ? 'fa-solid fa-sun' : 'fa-solid fa-moon';

// Event Listener
themeToggleBtn.addEventListener('click', () => {
    const currentTheme = htmlElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    htmlElement.setAttribute('data-theme', newTheme);
    if (themeIcon) themeIcon.className = newTheme === 'light' ? 'fa-solid fa-sun' : 'fa-solid fa-moon';
    
    localStorage.setItem('site_theme', newTheme);
});