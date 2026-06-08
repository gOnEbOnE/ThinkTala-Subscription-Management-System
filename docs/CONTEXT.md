# Thinknalyze Design Context (Hi-Fi Handoff)

Last updated: 2026-04-29

This document is design-only. It captures the current visual system, layout conventions, and page templates so new pages stay consistent. It intentionally avoids codebase and infrastructure details.

## 1) Design goals

- Dark-first, high-contrast, neon-accented UI.
- Dashboard-centric layout with structured grids and clear hierarchy.
- Consistent surfaces, borders, and spacing across roles and pages.
- Subtle motion and glow on interactive states; avoid heavy animation.

## 2) Information architecture (role-based)

Primary role portals (each with its own dashboard set):
- Management: analytics dashboards (Customer Retention & Churn, Package Sales) plus detail pages.
- Operational: ops dashboard and ops tools (notifications, orders, subscriptions, users).
- Client: client dashboard and account-related flows.
- Compliance: compliance dashboard and review flows.

Navigation principle: each role has a dedicated sidebar and top navbar, with two-level hierarchy where needed.

## 3) Visual system

### 3.1 Color palette (dark theme)

Backgrounds and surfaces:
- bg-dark: #0b0e17
- bg-sidebar: #0f121d
- bg-panel: #151a2d
- bg-header: rgba(15, 18, 29, 0.85)
- border-color: rgba(255, 255, 255, 0.08)

Accents:
- accent-cyan: #00f2ff
- accent-purple: #bd00ff
- accent-green: #00d26a
- accent-red: #f93e3e
- accent-orange: #ff9900

Text:
- text-heading: #ffffff
- text-main: #e2e8f0
- text-muted: #8b9bb4

Shadows:
- shadow-card: 0 4px 20px rgba(0, 0, 0, 0.3)

### 3.2 Light theme (secondary)

Light theme exists but keeps sidebar dark for contrast. Use these overrides:
- bg-dark: #f4f6f8
- bg-panel: #ffffff
- bg-header: rgba(255, 255, 255, 0.9)
- border-color: #e2e8f0
- text-heading: #0f172a
- text-main: #334155
- text-muted: #64748b
- accent-cyan: #007780
- accent-purple: #700099
- accent-green: #00a152
- accent-red: #d32f2f
- accent-orange: #e65100

### 3.3 Typography

- Primary: Inter (body, headings)
- Mono: JetBrains Mono (labels and data snippets)

Common sizes (from existing pages):
- Page title: 1.7rem, 700
- Section titles: 1.2rem to 1.25rem, 700
- Card value: 2.0rem, 700
- Meta text: 0.8rem to 0.9rem

## 4) Layout system

Global layout (desktop):
- Fixed left sidebar
- Fixed top navbar
- Main content panel with padding and grid layout

Sizing tokens:
- Sidebar width: 260px (collapsed: 80px)
- Header height: 80px
- Footer height: 50px

Spacing and radius:
- Card radius: 10px to 12px
- Buttons radius: 7px to 12px
- Inner padding: 16px to 20px

## 5) Navigation components

Sidebar:
- Brand row with icon + "ThinkTala"
- Primary nav links with icon-left
- Active state: cyan text, left border, subtle glow background
- Collapsible state for desktop; overlay for mobile

Top navbar:
- Left: sidebar toggle + role badge
- Right: theme toggle + avatar dropdown
- Dropdown uses rounded corners and soft shadow

## 6) Core components

### 6.1 KPI / Stat cards

- Label + icon + main value + meta subtext
- Uses panel background, border, and subtle shadow
- Icons in 32px rounded square

### 6.2 Chart cards

- Title + subtitle + chart wrap container
- Empty hint text below chart (muted)
- Chart wrap has border + light glass surface

### 6.3 Filter strip

- Date range inputs + period select on left
- Action buttons on right (Refresh, Export)
- Error banner and empty-note appear below

### 6.4 Tables

- Uppercase header, tight row spacing
- Empty row text centered and muted
- Action buttons with small icon buttons (square)

### 6.5 Notifications

- Toast banner for transient errors (top area)
- Inline error banner with retry action
- Empty-state note for no data

## 7) Chart styling

### 7.1 Package charts

Distinct bar/doughnut palette for packages:
- #1E88E5, #E53935, #43A047, #FB8C00, #8E24AA,
- #00897B, #6D4C41, #3949AB, #C0CA33, #F4511E

Sales trend line colors:
- Premium: #00f2ff
- Enterprise: #5ad9ff
- Starter: #86b3ff

### 7.2 Customer charts

- New customers line: cyan (accent-cyan)
- Churned customers line: #9aa3af
- Churn rate line: red (accent-red)

## 8) Page templates (existing design)

### 8.1 Management - Customer Retention & Churn

Structure:
- Page header with title and subtitle
- KPI row (Active, Loyal, Churned, Churn Rate)
- Filter strip (date range + period + actions)
- Inline error banner + empty-period note
- Chart 1: User Growth vs Churn (line chart)
- Chart 2: Churn Rate Trend (line chart)
- Table: Top Loyal Customers with search input
- Table: Recently Churned Customers

Copy conventions:
- Subtitles explain the chart intent in 1 line.
- Empty states: "Tidak ada data pada periode ini."

### 8.2 Management - Package Sales Dashboard

Structure:
- Page header with title and subtitle
- KPI row (Most Sold, Highest Revenue, Fastest Growth, Total Revenue)
- Filter strip (date range + period + actions)
- Inline error banner + empty-period note
- Chart 1: Package Sales Distribution (bar)
- Chart 2: Revenue Contribution (doughnut)
- Chart 3: Sales Trend Over Time (line)
- Table: Package Performance Details + pagination

Copy conventions:
- Use currency formatted in IDR
- Empty states: "Data tidak tersedia pada periode ini."

## 9) Interaction patterns

- Theme toggle: dark to light, with icon swap
- Sidebar collapse: remembers state on desktop
- Buttons: hover glows with cyan tint
- Dropdown menus: slide and fade in

## 10) Consistency rules for new Hi-Fi pages

- Use the same background and panel colors; avoid new accent colors.
- Keep typography scale and weights consistent with KPI and chart cards.
- Reuse filter strip and error/empty states for data pages.
- Use cyan for primary action highlights; red only for error and destructive.
- Maintain 12px rounding and 16-20px padding in cards.

If you want specific new PBIs designed, provide the PBI list and desired page outputs; the above context ensures visual continuity.