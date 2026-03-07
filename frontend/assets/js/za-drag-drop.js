/**
 * LAYOUT MANAGER (Drag & Drop & Resize)
 * File: assets/js/za-drag-drop.js
 */

// --- 1. INTELLIGENT DRAG & DROP ---
function initSortable() {
    if (window.innerWidth > 768) {
        const dynamicRows = document.querySelectorAll('.dynamic-row');

        dynamicRows.forEach(row => {
            // Cek jika row sudah memiliki instance Sortable, jangan buat double
            if (Sortable.get(row)) return; 

            new Sortable(row, {
                group: 'shared-dashboard',
                animation: 350,
                easing: "cubic-bezier(1, 0, 0, 1)",
                handle: '.drag-handle',
                ghostClass: 'sortable-ghost',
                dragClass: 'sortable-drag',
                
                onAdd: function (evt) {
                    updateRowLayout(evt.to); 
                },
                onRemove: function (evt) {
                    updateRowLayout(evt.from); 
                },
                onEnd: function (evt) {
                    document.body.style.cursor = 'default';
                }
            });
        });
    }
}

// --- 2. AUTO COLUMN CALCULATOR ---
function updateRowLayout(rowElement) {
    const items = rowElement.querySelectorAll('.draggable-item');
    const count = items.length;

    let newClass = 'col-lg-12'; // Default
    if (count === 2) newClass = 'col-lg-6';
    else if (count === 3) newClass = 'col-lg-4';
    else if (count >= 4) newClass = 'col-lg-3';

    items.forEach(item => {
        item.classList.remove('col-lg-12', 'col-lg-6', 'col-lg-4', 'col-lg-3', 'col-md-12');
        item.classList.add('col-12', 'col-md-6', newClass);
    });
}

// --- 3. MANUAL RESIZE LOGIC ---
function initResizable() {
    const items = document.querySelectorAll('.draggable-item');

    items.forEach(item => {
        // Mencegah duplikasi handle resize
        if (!item.querySelector('.resize-handle')) {
            const handle = document.createElement('div');
            handle.className = 'resize-handle';
            item.appendChild(handle);

            handle.addEventListener('mousedown', function(e) {
                e.preventDefault();
                const startX = e.pageX;
                const startWidth = item.offsetWidth;
                const parentRow = item.closest('.row');
                const parentWidth = parentRow.offsetWidth;
                
                item.classList.add('resizing');

                function onMouseMove(e) {
                    const currentWidth = startWidth + (e.pageX - startX);
                    const ratio = (currentWidth / parentWidth) * 12;

                    let newClass = '';
                    if (ratio > 9) newClass = 'col-lg-12';
                    else if (ratio > 5.5) newClass = 'col-lg-6';
                    else if (ratio > 3.5) newClass = 'col-lg-4';
                    else newClass = 'col-lg-3';

                    if (!item.classList.contains(newClass)) {
                        item.classList.remove('resizing'); 
                        item.classList.remove('col-lg-12', 'col-lg-8', 'col-lg-6', 'col-lg-4', 'col-lg-3');
                        item.classList.add(newClass);

                        setTimeout(() => {
                            if (e.buttons > 0) item.classList.add('resizing');
                        }, 400); 
                    }
                }

                function onMouseUp() {
                    item.classList.remove('resizing');
                    window.removeEventListener('mousemove', onMouseMove);
                    window.removeEventListener('mouseup', onMouseUp);
                }

                window.addEventListener('mousemove', onMouseMove);
                window.addEventListener('mouseup', onMouseUp);
            });
        }
    });
}