# SMS ThinkNalyze

A production-grade, microservices-based subscription management and business analytics platform built in Go. Designed to serve multiple organisational roles with purpose-built dashboards and strictly enforced access boundaries.

**Live Demo:** [propensuy-thinknalyze.vercel.app](https://propensuy-thinknalyze.vercel.app)

> Best Technical Engineering Award — Propensuy Project, 2026

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [ISO Standards Implemented](#iso-standards-implemented)
- [Demo Credentials](#demo-credentials)
- [Services](#services)

---

## Overview

SMS ThinkNalyze is a full-featured subscription management system that handles:

- **Identity & Access Management** — secure authentication, 9-role RBAC, OTP verification
- **KYC Compliance** — document verification, audit trail, regulatory-grade access gate
- **Subscription & Billing** — package catalog, order lifecycle, payment verification
- **Multi-Channel Notifications** — email, WhatsApp, and Telegram delivery with queue and retry
- **Business Analytics** — churn analysis, revenue metrics, retention tracking, CSV/PDF export
- **Support Ticketing** — ticket creation, tracking, and resolution

---

## Architecture

```
                   ┌────────────────────────────────────────┐
                   │             CLIENT BROWSER              │
                   │  (Encrypted UUID cookie — no JWT here) │
                   └─────────────────┬──────────────────────┘
                                     │ HTTPS
                   ┌─────────────────▼──────────────────────┐
                   │           GATEWAY  :2000                │
                   │   CORS · JWT Validation · RBAC Routing  │
                   └──┬──────┬──────┬──────┬──────┬──────┬──┘
                      │      │      │      │      │      │
                   Users  Acct  Notif  Sub  Mgmt  Ops  Tickets
                  :2006  :2001 :5003 :5004 :5006 :8080 :2004
                      │
               ┌──────▼──────┐     ┌───────────────────────┐
               │  PostgreSQL  │     │         Redis          │
               │  (per svc)  │     │  Sessions · JWT Store  │
               └─────────────┘     │  Notification Queue   │
                                   └───────────────────────┘
```

**Deployment:** Railway Cloud · Docker Compose · GitHub Actions CI/CD

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.25.5 · ZaFramework · Gin |
| Frontend | Vue 3 · TypeScript · Tailwind CSS |
| Database | PostgreSQL 13+ (per service) |
| Cache / Queue | Redis 6+ |
| Auth | RSA-signed JWT · Encrypted UUID cookies |
| Notifications | SMTP · Twilio (WhatsApp) · Telegram Bot API |
| Deployment | Railway Cloud · GitHub Actions |

---

## Features

### Secure Token Architecture

Rather than storing JWTs in cookies (vulnerable to XSS), ThinkNalyze uses a token-pointer pattern:

- The cookie holds only an **encrypted UUID pointer**, never the JWT itself
- The **JWT is stored server-side in Redis**, keyed by that UUID
- On every request, the Gateway decrypts the UUID, fetches the JWT from Redis, and validates it using RSA key pairs
- Logout **revokes the session globally** — not just on the client

```
Client Cookie → Encrypted UUID → Redis Lookup → JWT Validation
               (no JWT here)     (server side)   (RSA signed)
```

### Role-Based Access Control (RBAC)

Nine roles enforced at the API Gateway layer, not inside individual services:

| Role | Access |
|---|---|
| `CLIENT` | Personal dashboard, KYC, subscriptions, tickets |
| `OPERASIONAL` | Order management, notification templates, KYC queue |
| `COMPLIANCE` | KYC review and approval, compliance metrics |
| `MANAGEMENT` | Revenue analytics, churn analysis, export |
| `ADMIN` | Full administrative access |
| `ADMIN_SUPPORT` | Support and ticket management |
| `ADMIN_KYC` | KYC review and approval |
| `SUPERADMIN` | Platform-wide system management |
| `CEO` | Executive dashboards and reporting |

### KYC Compliance Module

KYC operates as a hard access gate — users cannot access platform features until verified:

```
User Submits KTP (ID Card)
        ↓
  Status: PENDING
        ↓
  Admin KYC Review
        ↓
APPROVED ←——→ REJECTED (reason logged, mandatory)
        ↓              ↓
 Access Granted  Resubmission Required
```

Every state transition is recorded with timestamps for regulatory audit.

### Multi-Channel Notification Engine

| Channel | Provider | Use Case |
|---|---|---|
| Email | SMTP | Registration, password reset, order confirmation |
| WhatsApp | Twilio API | OTP delivery, KYC status, subscription alerts |
| Telegram | Bot API | Operational alerts, admin notifications |

Redis RPUSH/BLPOP queue with priority tiers (high/normal/low) and HTTP fallback for reliability.

### Business Intelligence Dashboard

- User growth, churn rate, retention rate, loyal customer tracking
- Package sales by volume and revenue contribution
- Configurable date range filtering (monthly, yearly, custom)
- CSV and PDF export for reporting

---

## ISO Standards Implemented

| Standard | Application |
|---|---|
| ISO/IEC 27001 | ISMS: access control, cryptographic controls, session management, authentication |
| ISO/IEC 27002 | Security controls: audit logging, separation of duties, least privilege, monitoring |
| ISO 9001 | Quality management: CI/CD pipeline, modular architecture, KPI dashboards |
| ISO/IEC 25010 | Software quality: reliability fallback, async processing, maintainability |
| ISO 31000 | Risk management: identity fraud mitigation, session revocation, delivery failure handling |

---

## Demo Credentials

> These credentials are for demonstration purposes only. All roles are pre-configured with representative data.

| Role | Email | Password |
|---|---|---|
| Super Admin | superadmin@thinktala.com | Super123 |
| Operasional | ops@thinktala.com | Operas123 |
| Compliance | compliance@thinktala.com | Comply123 |
| Management | management@thinktala.com | Manage123 |
| Support | support@thinktala.com | Support123 |
| User (Client) | proaielaeu@gmail.com | User1234 |

**Access the live platform:** [propensuy-thinknalyze.vercel.app](https://propensuy-thinknalyze.vercel.app)

---

## Services

| Service | Port | Responsibility |
|---|---|---|
| Gateway | 2000 | Reverse proxy, CORS, JWT validation, RBAC routing |
| Users | 2006 | Registration, login, OTP, KYC, session management |
| Account | 2001 | User profile, account settings, personal dashboard |
| Notification | 5003 | Email/WhatsApp/Telegram delivery, template engine, queue |
| Subscription | 5004 | Package catalog, order management, payment verification |
| Operational | 8080 | Order monitoring, operational dashboards |
| Management | 5006 | Analytics, churn analysis, revenue metrics, export |
| Tickets | 2004 | Support ticket creation, tracking, resolution |

---

*Built for compliance, engineered for scale.*