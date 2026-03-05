/* FILE: za-market-insight.js 
   Fungsi: Mengambil data dari JSON, merender berita, kalender, dan sentimen.
   Status: Full Code Terintegrasi (Data + Rendering + Fallback Logic)
*/

let currentAiText = ""; 

// --- DATA CADANGAN (Agar News tidak stuck Loading jika file JSON gagal diakses) ---
const marketNewsFallback = {
    "last_updated_utc": "2026-02-11 02:05:17",
    "current_score": 0.4855996310710907,
    "articles": [
        { "title": "Stock Market LIVE: GIFT Nifty hints at positive start amid strong cues; Asian markets extend gains", "sentiment": { "label": "Positive" }, "source_url": "#" },
        { "title": "Dollar soft ahead of US data, yen holds onto its gains after election", "sentiment": { "label": "Positive" }, "source_url": "#" },
        { "title": "Global Market Today: Asian stocks extend rally to record, gold falls", "sentiment": { "label": "Positive" }, "source_url": "#" },
        { "title": "Don’t get comfortable with the global stock rally today: Goldman’s Panic Index is approaching ‘max fear’", "sentiment": { "label": "Negative" }, "source_url": "#" },
        { "title": "Japan's Nikkei 225 set to continue post election rally as Asian markets on course to open higher", "sentiment": { "label": "Positive" }, "source_url": "#" },
        { "title": "Stock Market Today: Dow Slightly Higher; Japanese Stocks Jump After Election", "sentiment": { "label": "Neutral" }, "source_url": "#" }
    ]
};

// --- 1. CORE DATA UPDATER ---
async function updateDashboard() {
    try {
        const resMacro = await fetch('macro_history.json');
        const dataMacro = await resMacro.json();
        const latest = dataMacro[dataMacro.length - 1];
        if (latest) {
            document.getElementById('lastUpdate').innerText = latest.timestamp_local.split(' ')[1];
            // Memanggil updateClock dari helper lokal di HTML
            if (typeof updateClock === "function") updateClock();
            
            updateSentiment(latest.scores.risk_score, latest.decision);
            
            if (latest.expert_analysis && latest.expert_analysis.text !== currentAiText) {
                currentAiText = latest.expert_analysis.text;
                document.getElementById('aiTime').innerText = latest.expert_analysis.generated_at;
                typeWriter(currentAiText, 'typewriterText');
            }
            const m = latest.market_metrics;
            document.getElementById('activeSessionBadge').innerText = latest.session;
            updateFlow(m.global_heatmap);
            updateDrivers(m.global_assets, m.risk_index_val, latest.session);
            updateCommodities(m.commodities);
            renderIndices(m.active_indices);
        }
    } catch (e) { console.error("Macro Fetch Error:", e); }

    try {
        const resCal = await fetch('calendar_data.json');
        const dataCal = await resCal.json();
        renderCalendar(dataCal.events);
    } catch (e) { console.error("Calendar Fetch Error:", e); }
}

// --- 2. LIVE NEWS FEED ---
async function loadMarketNews() {
    const container = document.getElementById('newsFeedContainer');
    if (!container) return;

    try {
        // Mencoba fetch file asli
        const response = await fetch('news_database.json').catch(() => null);
        
        let data;
        if (response && response.ok) {
            data = await response.json();
        } else {
            // Gunakan data cadangan jika file tidak ada atau gagal
            data = marketNewsFallback; 
        }
        renderNewsList(data);
    } catch (e) {
        console.error("News Load Error:", e);
        renderNewsList(marketNewsFallback); // Tetap render cadangan meskipun catch
    }
}

function renderNewsList(data) {
    const container = document.getElementById('newsFeedContainer');
    if(!container) return;
    
    const scoreBadge = document.getElementById('newsSentimentScore');
    const globalScore = data.current_score.toFixed(2);
    
    if (scoreBadge) {
        scoreBadge.innerText = `Sent: ${globalScore}`;
        // Atur warna badge sentimen global
        if(globalScore > 0.2) scoreBadge.className = "badge bg-success bg-opacity-25 text-green border border-success";
        else if(globalScore < -0.2) scoreBadge.className = "badge bg-danger bg-opacity-25 text-red border border-danger";
        else scoreBadge.className = "badge bg-secondary border border-secondary text-muted";
    }
    
    let html = "";
    data.articles.forEach(article => {
        const label = article.sentiment.label;
        let badgeClass = label === "Positive" ? "sent-pos" : (label === "Negative" ? "sent-neg" : "sent-neu");
        html += `
            <div class="news-item" style="border-bottom: 1px solid var(--border-color); padding: 10px 0;">
                <div class="d-flex justify-content-between mb-1">
                    <span class="sentiment-badge ${badgeClass}" style="font-size:0.6rem; padding:2px 5px; border-radius:4px; font-weight:bold;">${label}</span>
                    <small class="text-muted" style="font-size: 0.65rem;">10m ago</small>
                </div>
                <a href="${article.source_url || '#'}" target="_blank" class="text-decoration-none">
                    <p class="mb-0 text-main small fw-bold lh-sm news-title">${article.title}</p>
                </a>
            </div>`;
    });
    container.innerHTML = html;
    
    // Re-inisialisasi Resize Handle agar kartu News bisa ditarik
    if (typeof initResizable === "function") initResizable(); 
}

// --- 3. CURRENCY STRENGTH METER ---
function loadCSM() {
    const currencies = [
        { name: "USD", val: 85, color: "bg-success" }, 
        { name: "JPY", val: 20, color: "bg-danger" }, 
        { name: "EUR", val: 45, color: "bg-danger" }, 
        { name: "GBP", val: 60, color: "bg-warning" }
    ];
    const container = document.getElementById('csmContainer');
    if(!container) return;
    
    let html = "";
    currencies.forEach(curr => {
        html += `
            <div>
                <div class="d-flex justify-content-between align-items-end mb-1">
                    <span class="fw-bold text-main font-mono small">${curr.name}</span>
                    <span class="small">${curr.val/10}</span>
                </div>
                <div class="csm-bar-bg" style="height: 6px; background: rgba(255,255,255,0.1); border-radius: 10px; overflow: hidden;">
                    <div class="csm-bar-fill ${curr.color}" style="width: ${curr.val}%; height: 100%; border-radius: 10px;"></div>
                </div>
            </div>`;
    });
    container.innerHTML = html;
    
    if (typeof initResizable === "function") initResizable();
}

// --- 4. DATA HELPERS ---
function updateSentiment(score, decision) {
    const lbl = document.getElementById('riskLabel');
    if(!lbl) return;
    lbl.innerText = decision.status.replace(/[^a-zA-Z ]/g, "");
    lbl.className = score > 0.2 ? "fw-bold text-green" : (score < -0.2 ? "fw-bold text-red" : "fw-bold text-warning");
    
    const pointer = document.getElementById('riskPointer');
    if(pointer) pointer.style.left = ((score + 1) / 2 * 100) + "%";
    
    const moodText = document.getElementById('aiMoodText');
    if(moodText) moodText.innerText = `"${decision.strategy}"`;
}

function updateFlow(map) {
    const updateBox = (idBox, idVal, val) => {
        const el = document.getElementById(idVal);
        const box = document.getElementById(idBox);
        if(!el || !box) return;
        el.innerText = val.toFixed(2) + "%";
        box.style.background = val > 0.1 ? "rgba(0, 210, 106, 0.1)" : (val < -0.1 ? "rgba(249, 62, 62, 0.1)" : "");
        el.className = val > 0.1 ? "fw-bold text-green" : (val < -0.1 ? "fw-bold text-red" : "fw-bold text-main");
    };
    updateBox('boxAsia', 'valAsia', map.ASIA);
    updateBox('boxEurope', 'valEurope', map.EUROGE || map.EUROPE); // Toleransi jika ada typo key
    updateBox('boxUS', 'valUS', map.US);
}

function updateDrivers(assets, riskIndex, sessionName) {
    const set = (idVal, val, inv=false) => {
        const el = document.getElementById(idVal);
        if(!el) return;
        el.innerText = (val>0?"+":"") + val.toFixed(2) + "%";
        el.className = "fw-bold " + (val > 0.01 ? (inv ? "text-red" : "text-green") : (val < -0.01 ? (inv ? "text-green" : "text-red") : "text-main"));
    };
    set('d-index', riskIndex); 
    set('d-vix', assets.vix, true); 
    set('d-dxy', assets.dxy, true); 
    set('d-yield', assets.bond, true);
}

function updateCommodities(comm) {
    const set = (id, val) => { 
        const el = document.getElementById(id);
        if(el) el.innerText = (val>0?"+":"") + val.toFixed(2) + "%"; 
    };
    if(comm) {
        set('comm-ind', comm.INDUSTRIAL); 
        set('comm-ene', comm.ENERGY); 
        set('comm-pre', comm.PRECIOUS); 
        set('comm-agr', comm.AGRICULTURE);
    }
}

function renderIndices(indices) {
    const c = document.getElementById('indicesContainer');
    if(!c) return;
    c.innerHTML = "";
    for (const [key, val] of Object.entries(indices)) {
        c.innerHTML += `
            <div class="col-6">
                <div class="p-2 border border-secondary rounded">
                    <div class="d-flex justify-content-between mb-1">
                        <small class="text-muted fw-bold font-mono" style="font-size:0.7rem;">${key}</small>
                        <small class="${val >= 0 ? 'text-green' : 'text-red'}">${val.toFixed(2)}%</small>
                    </div>
                    <div class="progress" style="height: 3px; background: rgba(255,255,255,0.1);">
                        <div class="progress-bar ${val>=0?'bg-success':'bg-danger'}" style="width: ${Math.min(Math.abs(val)*40, 100)}%"></div>
                    </div>
                </div>
            </div>`;
    }
}

function renderCalendar(events) {
    const container = document.getElementById('calendarList');
    if(!container) return;
    let html = "";
    if(events) {
        events.forEach(evt => {
            const color = evt.impact === "High" ? "text-red" : (evt.impact === "Medium" ? "text-warning" : "text-green");
            html += `
                <div class="cal-item d-flex align-items-center justify-content-between" style="border-bottom: 1px solid rgba(255,255,255,0.05); padding: 8px 0;">
                    <div class="d-flex align-items-center">
                        <span class="badge cal-badge me-2" style="font-size:0.65rem;">${evt.country}</span>
                        <div class="small text-main fw-bold text-truncate" style="max-width:140px;">${evt.name}</div>
                    </div>
                    <small class="${color} fw-bold" style="font-size:0.65rem;">${evt.impact}</small>
                </div>`;
        });
    }
    container.innerHTML = html || "<div class='text-muted small text-center'>No events found.</div>";
}

function typeWriter(text, elementId) {
    const el = document.getElementById(elementId);
    if(!el) return;
    el.innerHTML = ""; let i = 0;
    function type() { 
        if (i < text.length) { 
            el.innerHTML += text.charAt(i); 
            i++; 
            setTimeout(type, 25); 
        } 
    }
    type();
}