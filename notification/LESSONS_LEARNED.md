# Lessons Learned — Notification Service Architecture

## Overview

Dalam pengembangan Notification Microservice ini, saya mempelajari implementasi arsitektur komunikasi antar-service yang terdesentralisasi, pola layering yang jelas untuk maintainability, serta strategi templating yang memungkinkan fleksibilitas delivery channel tanpa perlu deployment ulang.

---

## 1. Architecture & Design Pattern

### Microservices Communication
- **Asynchronous Events:** Notification service berfungsi sebagai *event consumer* (menerima trigger dari users, operational, subscription services) dan *broadcast publisher* (mengirim notif ke client).
- **Service Autonomy:** Setiap service (users, operational, operational) mampu beroperasi independen. Notification service hanya menyimpan template dan log; business logic tetap di service yang memanggil.
- **API Contract:** Mendefinisikan payload event yang jelas (`event_type`, `user_id`, `email`, `data`) memudahkan integrasi antar service tanpa tight coupling.

### Layered Architecture (Controller-Service-Repository)
- **Separation of Concerns:** Controller hanya parse HTTP, Service menjalankan business logic, Repository akses database → mudah test per layer dan debug.
- **Single Responsibility:** Tiap class/function punya satu alasan untuk berubah:
  - Repository berubah → query DB berubah
  - Service berubah → rule bisnis (misal: validasi subject untuk email channel) berubah
  - Controller berubah → format response HTTP berubah
- **Reusability:** Logic di Service bisa dipanggil dari berbagai handler (HTTP endpoint, event listener, background job).

---

## 2. Database Design

### Dual-Model Data
- **`notifications` table:** Broadcast announcements ke client berdasarkan role. Bersifat *published content*, tidak tied ke specific user.
- **`notification_templates` table:** Template library mapping `event_type` → template konten. Bersifat *configuration*, reusable across events.

### Indexing Strategy
- `idx_notif_tpl_event` pada `event_type` → cepat lookup template saat event triggered
- `idx_notif_tpl_channel` pada `channel` → cepat filter template by delivery method (email, WhatsApp, Telegram)
- `idx_notifications_active` pada `(is_active, target_role)` → composite index untuk query public notifications yang aktif

### Soft Configuration
- Template disimpan di database, bukan hardcoded di code → **tidak perlu deploy ulang** saat ganti konten pesan
- Placeholder dalam template (`{{name}}`, `{{otp}}`) → flexible tanpa validasi ketat di layer schema

---

## 3. Template System & Channel Abstraction

### Event-Driven Templating
- Satu `event_type` bisa punya **multiple templates per channel** (email, WhatsApp, Telegram untuk event yang sama)
- Service pemicu (users, operational) tidak perlu tahu tentang channel — hanya kirim event dengan `event_type` yang standar
- Notification service lookup template, apply placeholder, kirim ke channel yang sesuai

### Flexibility Benefits
1. **Product A/B Testing:** Bisa buat 2 template berbeda untuk event yang sama, observe mana yang convert lebih baik
2. **Localization:** Template bisa stored per bahasa tanpa perlu app code change
3. **Quick Iteration:** Ops team bisa update konten pesan dari dashboard tanpa developer intervention

---

## 4. Code Organization Best Practices

### Module Structure
```
Setiap modul (notification, template) adalah self-contained:
  ├── types.go       → Contracts (request/response)
  ├── repository.go  → Data queries
  ├── service.go     → Business rules
  └── controller.go  → HTTP handlers
```
✅ Mudah navigate — tahu mau cari apa di mana  
✅ Mudah test — bisa mock tiap layer  
✅ Mudah onboard — junior dev cepat paham alur  

### main.go Minimal
- Hanya: init DB → run migrations → register routes → start server
- **No business logic in main.go**
- Semua feature tambahan bisa di-add tanpa sentuh main.go → reduce merge conflicts, safer deployments

### Router Centralized
- Semua route registration di satu file (`app/routes/router.go`) → mudah lihat API surface
- Path consistency → avoid `/api/notification` vs `/api/notifications` kesalahan

---

## 5. Integration Patterns

### Event Forwarding Flow
```
users/register → OTP generated
       │
       └─→ (async) POST /api/notifications/event (to notification service)
             event_type: "otp_verification"
             user_id: "uuid"
             email: "user@mail.com"
             data: { otp: "123456", expires_in: "5m" }
             │
             ├─→ notification service lookup template (event_type → database)
             ├─→ apply placeholder ({{otp}} → 123456)
             └─→ log matched template for audit
```

### Why Async?
- Notification service slow/down → **tidak block** user registration flow
- Retry logic di notification service independent dari caller
- Scalability → notification queue bisa diganti dengan RabbitMQ/Kafka tanpa perlu ubah users service

---

## 6. Data Consistency & Audit Trail

### Immutable Event Records
- Setiap event stored dengan `created_at`, `created_by`
- Template history preserved (update template tidak menghapus old versions)
- Enable audit trail: "Kapa template berubah? Apa isinya dulu?" → traceable

### Validation Layering
- **Repository level:** Database constraint (NOT NULL, FK, UNIQUE)
- **Service level:** Business rule (subject required for email, event_type valid, channel supported)
- **Controller level:** Input validation (binding:"required", field type coercion)
- Benefit: Jika ada bug di controller, service masih catch. Jika bug di service, DB constraint masih guard.

---

## 7. Operational Considerations

### Monitoring Points
- Notification service availability (status endpoint)
- Template lookup success/failure ratio
- Event processing latency (log event received → template matched)
- Placeholder substitution errors (if template expect {{kyc_reason}} tapi caller tidak kirim)

### Future Scalability
- Notification service stateless → easy horizontal scale (multiple instances behind load balancer)
- Database as bottleneck → can add read replica untuk query template
- Event queue pattern ready → replace sync HTTP POST dengan async queue (RabbitMQ, Kafka, Google Pub/Sub)

---

## 8. Collaborative Development Insights

### Cross-Service Coordination
| Service | Role | Integration Point |
|---------|------|-------------------|
| **users** | Event producer | Trigger `otp_verification`, `user_register` events |
| **operational** | Event producer | Trigger `user_kyc_approved`, `user_kyc_rejected` events |
| **notification** | Event consumer + publisher | Lookup template, store broadcast notifications |
| **subscription** | Event producer (future) | Trigger `subscription_renewed`, `payment_failed` |

### API Contract Clarity
```json
// Standard Event Payload (shared across all producers)
{
  "event_type": "string",      // otp_verification | user_kyc_approved | ...
  "user_id": "string",         // Who this event is about
  "email": "string",           // Where to send notification
  "name": "string",            // Recipient display name
  "data": { /* ... */ }        // Event-specific payload (otp, kyc_reason, etc)
}
```
**Benefit:** Notification service build once, support many producers without renegotiation

### Ownership Boundaries
- **Users service:** Owns user credential, OTP generation, user account lifecycle
- **Operational service:** Owns KYC review flow, document validation
- **Notification service:** Owns template library, broadcast content, delivery metadata
- No duplicated logic, clear responsibility ✓

---

## 9. Key Takeaways

1. **Template as Configuration:** Database-driven templates → dynamic content without redeployment
2. **Event Standardization:** Common event payload contract → easy onboarding of new service producers
3. **Async by Default:** Keep service calls non-blocking for resilience
4. **Layered Code:** Controller ≠ Service ≠ Repository → clean testing, debugging, scaling
5. **Audit Everything:** timestamps, user actions, template versions → enable operational visibility
6. **Microservice Autonomy:** Each service independent, communicate via well-defined events, not shared databases
7. **Documentation First:** Lessons learned, API contract, module README → easier handoff, onboarding, maintenance

---

## 10. What Would Improve Further

- [ ] Add **event logging table** (`notification_events`) to track every event received + matched template
- [ ] Implement **rate limiting** on event endpoint to prevent spam/abuse
- [ ] Add **placeholder validation** in service layer: verify template placeholders are provided in event.data
- [ ] Implement **retry logic** for failed deliveries (WhatsApp throttle, email bounce, etc)
- [ ] Create **template versioning** (keep history, rollback capability)
- [ ] Add **delivery status tracking** table to know which notification was actually sent to whom

---

**Session Outcome:** Notification Microservice berjalan sebagai event-driven, template-based centralized hub komunikasi, dengan clean layered architecture, allowing other services to send structured events without worrying about "how to notify"—notification service handle everything.
