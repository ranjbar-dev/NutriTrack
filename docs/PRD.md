# Product Requirements Document (PRD)

## NutriTrack — Nutrition Management PWA

**Version:** 1.0  
**Date:** April 19, 2026  
**Status:** Draft  
**Author:** Joe (Product Owner)

---

## Table of Contents

1. [Overview](#1-overview)
2. [Goals & Non-Goals](#2-goals--non-goals)
3. [User Roles & Permissions](#3-user-roles--permissions)
4. [Authentication & Registration](#4-authentication--registration)
5. [Feature Specifications](#5-feature-specifications)
   - 5.1 [Food Database Management](#51-food-database-management)
   - 5.2 [Medication Database Management](#52-medication-database-management)
   - 5.3 [Diet Plan Management](#53-diet-plan-management)
   - 5.4 [Daily Tracking (Client Side)](#54-daily-tracking-client-side)
   - 5.5 [Body Measurements](#55-body-measurements)
   - 5.6 [Exercise Management](#56-exercise-management)
   - 5.7 [Water Intake Tracking](#57-water-intake-tracking)
   - 5.8 [Sleep Tracking](#58-sleep-tracking)
   - 5.9 [Medication Tracking](#59-medication-tracking)
   - 5.10 [Lab Results Upload](#510-lab-results-upload)
   - 5.11 [Messaging System](#511-messaging-system)
   - 5.12 [Food Request System](#512-food-request-system)
   - 5.13 [Client Management (Nutritionist)](#513-client-management-nutritionist)
   - 5.14 [Super Admin Panel](#514-super-admin-panel)
6. [Offline & Sync Strategy](#6-offline--sync-strategy)
7. [Notifications](#7-notifications)
8. [Tech Stack & Architecture](#8-tech-stack--architecture)
9. [Data Model](#9-data-model)
10. [Non-Functional Requirements](#10-non-functional-requirements)
11. [Decision Log](#11-decision-log)
12. [Assumptions](#12-assumptions)

---

## 1. Overview

NutriTrack is a mobile-first Progressive Web Application (PWA) for managing the relationship between nutritionists and their clients. Nutritionists create personalized diet plans, prescribe medications, and monitor client progress. Clients view their plans, track daily intake, exercise, sleep, water consumption, and body measurements — with full offline support for viewing and data entry.

The platform is **Persian-only (RTL)** and designed exclusively for **mobile viewport**.

**Tech Stack:**
- **Backend:** Go (Golang)
- **Database:** PostgreSQL
- **Frontend:** Nuxt.js (Vue 3 / Nuxt 4)
- **Infrastructure:** Self-hosted on Hetzner with Docker/Docker Compose

---

## 2. Goals & Non-Goals

### Goals

- Provide a structured, digital workflow for nutritionists to manage clients and diet plans
- Enable clients to track daily food intake, exercise, water, sleep, weight, and body measurements
- Support full offline capability for clients (view plans + record data, sync when online)
- Deliver a minimal, clean, Persian RTL mobile experience
- Build a shared food and medication database accessible to all nutritionists
- Enable secure messaging between clients and their assigned nutritionist
- Push notifications via PWA Web Push for reminders and messages

### Non-Goals

- Desktop-optimized UI (mobile-only design)
- Multi-language support (Persian only, no i18n infrastructure needed)
- Real-time video/voice consultation
- Payment processing or subscription billing
- Integration with external health devices or wearables
- AI-powered diet recommendations
- Calorie auto-detection from food photos

---

## 3. User Roles & Permissions

### 3.1 Super Admin

- Create, edit, activate/deactivate nutritionist accounts
- Full CRUD on shared food database
- Full CRUD on shared medication database
- View platform-wide statistics
- Cannot manage individual clients (that's the nutritionist's domain)

### 3.2 Nutritionist (کارشناس تغذیه)

- Manage own clients: register, activate, deactivate, view records
- Create and manage diet plans for own clients
- Prescribe medications to own clients
- CRUD on shared food database (additions visible to all)
- CRUD on shared medication database (additions visible to all)
- View client tracking data: food logs, weight, measurements, exercise, sleep, water, medication intake
- Chat with own clients (text, image, file)
- Set required daily water intake per client (optional)
- Record client weight and body measurements
- Approve/reject food addition requests from clients

### 3.3 Client (کاربر عادی)

- View active diet plan and meal options
- Log daily food intake (select which option was eaten per meal, optional)
- Track daily water intake with timestamps
- Track sleep times
- Track exercise with optional calorie burn
- Log medication intake with timestamps
- Record daily weight
- Record body measurements (waist, hip, abdomen, thigh, chest, wrist)
- Upload lab results (PDF or link)
- Chat with assigned nutritionist (text, image, file)
- Submit food addition requests
- Receive push notifications

---

## 4. Authentication & Registration

### 4.1 Super Admin Login

- Email + password authentication
- Seeded/created via backend migration or CLI command
- No self-registration

### 4.2 Nutritionist Login

- Email + password authentication
- Account created exclusively by Super Admin
- No self-registration

### 4.3 Client Registration & Login

- **Registration:** Nutritionist registers the client by entering: full name, mobile number, date of birth, height, gender
- **Login:** Client receives OTP via SMS to their registered mobile number
- **Session:** JWT-based with refresh tokens
- **No self-registration** — clients cannot sign up without a nutritionist

### 4.4 OTP Flow

1. Client enters mobile number
2. Backend generates 6-digit OTP, sends via SMS gateway
3. OTP valid for 2 minutes, max 3 attempts
4. On success, issue JWT access token (15 min) + refresh token (30 days)
5. Rate limit: max 3 OTP requests per phone per 10 minutes

---

## 5. Feature Specifications

### 5.1 Food Database Management

The food database is a **shared, platform-wide** resource. Both Super Admin and Nutritionists can manage it.

#### Food Item Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Food name in Persian |
| `category` | enum | Yes | See categories below |
| `calories` | float | Yes | Per measurement unit |
| `protein` | float | Yes | Grams |
| `carbohydrates` | float | Yes | Grams |
| `fat` | float | Yes | Grams |
| `fiber` | float | No | Grams |
| `sugar` | float | No | Grams |
| `sodium` | float | No | Milligrams |
| `measurement_unit` | enum | Yes | See units below |
| `measurement_amount` | float | Yes | Base amount for nutritional values |
| `description` | text | No | Additional notes |
| `is_active` | boolean | Yes | Soft delete |
| `created_by` | FK → User | Yes | Who added this item |

#### Food Categories

- `breakfast` — صبحانه
- `lunch` — ناهار
- `dinner` — شام
- `snack` — میان‌وعده
- `fruit` — میوه
- `beverage` — نوشیدنی
- `supplement` — مکمل
- `other` — سایر

> **Note:** A food item can belong to multiple categories (many-to-many relationship), as some items (e.g., bread, eggs) can appear in any meal type.

#### Measurement Units

- `gram` — گرم
- `kilogram` — کیلوگرم
- `tablespoon` — قاشق غذاخوری
- `teaspoon` — قاشق چایخوری
- `cup` — لیوان
- `piece` — عدد
- `slice` — برش
- `palm` — کف دست
- `matchbox` — قوطی کبریت
- `bowl` — کاسه
- `milliliter` — میلی‌لیتر
- `liter` — لیتر

#### Search & Filter

- Full-text search on name (Persian-aware)
- Filter by category
- Filter by active/inactive status
- Pagination with 20 items per page

---

### 5.2 Medication Database Management

Shared platform-wide resource, same access rules as food database.

#### Medication Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Medication name |
| `generic_name` | string | No | Generic/scientific name |
| `form` | enum | Yes | Tablet, capsule, syrup, injection, drop, powder, other |
| `dosage_unit` | string | Yes | e.g., mg, ml, IU |
| `description` | text | No | Usage notes, side effects |
| `is_active` | boolean | Yes | Soft delete |
| `created_by` | FK → User | Yes | Who added this item |

---

### 5.3 Diet Plan Management

This is the core feature. A nutritionist creates a periodic diet plan for a client.

#### Diet Plan Structure

```
Diet Plan (دوره رژیم)
├── period: start_date → end_date
├── notes: text (general notes from nutritionist)
├── daily_water_target_ml: integer (optional)
├── status: active | archived
│
├── Plan Day (روز رژیم) — can be "day template" applied to multiple dates
│   ├── day_number: 1, 2, 3, ...
│   │
│   ├── Meal (وعده غذایی)
│   │   ├── title: string (e.g., "صبحانه", "میان‌وعده عصر")
│   │   ├── scheduled_time: time (e.g., 09:00)
│   │   ├── order: integer (display order)
│   │   │
│   │   ├── Meal Option (گزینه غذایی) — client picks ONE per meal
│   │   │   ├── option_number: 1, 2, 3, ...
│   │   │   │
│   │   │   ├── Meal Option Item (آیتم غذایی)
│   │   │   │   ├── food_item: FK → Food
│   │   │   │   ├── quantity: float
│   │   │   │   ├── measurement_unit: enum
│   │   │   │   ├── notes: string (optional, e.g., "آب‌پز")
│   │   │   │   └── computed: calories, protein, carbs, fat, fiber
│   │   │   │
│   │   │   └── Computed Totals: sum of all items
│   │   │
│   │   └── Computed Totals: based on all options (min/max range)
│   │
│   ├── Exercise Recommendation (پیشنهاد ورزش)
│   │   ├── exercise_name: string
│   │   ├── duration_minutes: integer
│   │   ├── description: text (optional)
│   │   └── calories_burn_estimate: integer (optional)
│   │
│   └── Computed Daily Totals: sum of meals (min/max range based on options)
│
└── Prescribed Medications
    ├── medication: FK → Medication
    ├── dosage: string (e.g., "1 tablet")
    ├── frequency: string (e.g., "twice daily")
    ├── times: array of time (e.g., ["08:00", "20:00"])
    ├── instructions: text (e.g., "after meal")
    └── duration: start_date → end_date (optional, may extend beyond diet plan)
```

#### Diet Plan Rules

- Only **one active** diet plan per client at any time
- When a new plan is created, the previous plan is automatically **archived**
- Archived plans remain viewable for both nutritionist and client (history)
- Nutritionist can see **computed nutritional totals** for each option, meal, and day while building the plan
- Diet plan days can be:
  - **Unique per date:** Each date has its own meals
  - **Repeating pattern:** e.g., a 7-day cycle that repeats throughout the period

#### Nutritional Computation Display

When building a diet plan, the nutritionist sees real-time computed values:

```
وعده صبحانه — ساعت ۹:۰۰
├── گزینه ۱: ۳۵۰ کالری | پروتئین ۲۰g | کربوهیدرات ۳۵g | چربی ۱۵g
│   ├── نان (۳ کف دست) — ۲۱۰ کالری
│   └── تخم مرغ آب‌پز (۲ عدد) — ۱۴۰ کالری
│
├── گزینه ۲: ۳۲۰ کالری | پروتئین ۱۸g | کربوهیدرات ۴۰g | چربی ۱۲g
│   ├── نان (۴ کف دست) — ۲۸۰ کالری
│   └── پنیر (۳ قوطی کبریت) — ۴۰ کالری

جمع کل روزانه: ۱,۸۰۰ — ۲,۱۰۰ کالری (بسته به انتخاب‌ها)
```

---

### 5.4 Daily Tracking (Client Side)

Clients can optionally log what they actually ate each day.

#### Food Log Entry

| Field | Type | Required |
|-------|------|----------|
| `date` | date | Yes |
| `meal_id` | FK → Meal | Yes |
| `selected_option_id` | FK → MealOption | No (null = skipped meal) |
| `logged_at` | timestamp | Auto |
| `notes` | text | No |

- Client sees their daily plan with all meals and options
- For each meal, they can tap to select which option they ate, or mark as skipped
- This is **optional** — no enforcement
- Nutritionist can view the client's food log history

---

### 5.5 Body Measurements

Both client and nutritionist can record measurements. Each entry is date-stamped.

#### Measurement Fields

| Field | Type | Unit |
|-------|------|------|
| `weight` | float | kg |
| `waist` | float | cm |
| `hip` | float | cm |
| `abdomen` | float | cm |
| `thigh` | float | cm |
| `chest` | float | cm |
| `wrist` | float | cm |

#### Rules

- One record per date per field (last write wins, or update existing)
- Both client and nutritionist can create entries
- `recorded_by` field tracks who entered the data
- History is viewable as a list and as a simple chart (weight over time, measurements over time)
- Height and date of birth are set at registration and editable only by nutritionist

---

### 5.6 Exercise Management

#### Nutritionist Side (Recommendation)

- Part of the diet plan (see 5.3)
- Each plan day can have exercise recommendations with name, duration, description, and estimated calorie burn

#### Client Side (Tracking)

| Field | Type | Required |
|-------|------|----------|
| `date` | date | Yes |
| `exercise_name` | string | Yes |
| `duration_minutes` | integer | Yes |
| `calories_burned` | float | No (optional) |
| `notes` | text | No |
| `logged_at` | timestamp | Auto |

- Multiple exercise entries per day allowed
- Nutritionist can view client's exercise history

---

### 5.7 Water Intake Tracking

#### Nutritionist Side

- Set `daily_water_target_ml` in the diet plan (optional)
- View client's daily water logs

#### Client Side

| Field | Type | Required |
|-------|------|----------|
| `date` | date | Yes |
| `amount_ml` | integer | Yes |
| `time` | time | No (optional) |
| `logged_at` | timestamp | Auto |

- Multiple entries per day (each glass/cup is a separate entry)
- Daily summary shows total vs. target (if set)
- Simple visual progress indicator (e.g., water glass filling up)

---

### 5.8 Sleep Tracking

| Field | Type | Required |
|-------|------|----------|
| `date` | date | Yes |
| `sleep_time` | datetime | Yes |
| `wake_time` | datetime | Yes |
| `quality` | enum | No (good, fair, poor) |
| `notes` | text | No |

- One entry per date (can be updated)
- Nutritionist can view client's sleep history
- Duration is computed from sleep_time and wake_time

---

### 5.9 Medication Tracking

#### Nutritionist Side (Prescription)

- Part of the diet plan (see 5.3)
- Prescribe from shared medication database with dosage, frequency, and timing

#### Client Side (Intake Logging)

| Field | Type | Required |
|-------|------|----------|
| `date` | date | Yes |
| `prescribed_medication_id` | FK | No (null for self-reported meds) |
| `medication_name` | string | Yes (auto-filled if prescribed) |
| `dosage` | string | Yes |
| `taken_at` | time | Yes |
| `notes` | text | No |

- Client sees a list of prescribed medications with scheduled times
- Client taps to mark each dose as taken (with timestamp)
- Client can also log non-prescribed medications/supplements manually
- Nutritionist can view the medication intake history

---

### 5.10 Lab Results Upload

| Field | Type | Required |
|-------|------|----------|
| `title` | string | Yes (e.g., "آزمایش خون — فروردین ۱۴۰۵") |
| `type` | enum | Yes (blood_test, urine_test, thyroid, hormone, allergy, other) |
| `date` | date | Yes (date of the test) |
| `file` | file (PDF) | No |
| `link` | URL | No |
| `notes` | text | No |

- At least one of `file` or `link` must be provided
- Max file size: 10 MB
- Accepted formats: PDF, JPG, PNG
- Nutritionist can view and download all uploaded lab results
- Files stored on Hetzner server filesystem

---

### 5.11 Messaging System

Chat-style messaging between client and their assigned nutritionist.

#### Message Fields

| Field | Type | Required |
|-------|------|----------|
| `sender_id` | FK → User | Yes |
| `receiver_id` | FK → User | Yes |
| `content` | text | No (if attachment present) |
| `attachment_type` | enum | No (image, file) |
| `attachment_path` | string | No |
| `sent_at` | timestamp | Auto |
| `read_at` | timestamp | No |

#### Rules

- Client can only message their assigned nutritionist
- Nutritionist can message any of their own clients
- Supported attachments: images (JPG, PNG, max 5 MB), files (PDF, max 10 MB)
- Not real-time WebSocket — use **polling** (every 10 seconds when chat is open)
- Unread message count shown as badge
- Push notification sent on new message (via Web Push)
- Messages are ordered chronologically
- No message editing or deletion
- Files stored on Hetzner server filesystem

---

### 5.12 Food Request System

When a client cannot find a food item in the database, they can request its addition.

#### Request Fields

| Field | Type | Required |
|-------|------|----------|
| `food_name` | string | Yes |
| `description` | text | No (e.g., "نان سنگک محلی") |
| `status` | enum | Auto (pending → approved / rejected) |
| `requested_by` | FK → Client | Auto |
| `reviewed_by` | FK → Nutritionist | No |
| `created_at` | timestamp | Auto |

#### Flow

1. Client submits food request with name and optional description
2. Request goes to the client's assigned nutritionist
3. Nutritionist reviews: approves (creates the food item in shared database) or rejects (with optional reason)
4. Client receives notification of the result

---

### 5.13 Client Management (Nutritionist)

#### Client List View

- List of all clients with: name, mobile, status (active/inactive), current plan status, last activity date
- Search by name or mobile
- Filter by active/inactive
- Sort by name or last activity

#### Client Profile View

- Personal info: name, mobile, date of birth, height, gender
- Current active diet plan (summary)
- History tabs:
  - Weight & measurements chart
  - Food logs
  - Exercise logs
  - Water intake
  - Sleep logs
  - Medication logs
  - Lab results
  - Archived diet plans
- Quick actions: new diet plan, send message, deactivate client

#### Client Registration

Nutritionist fills in:

| Field | Type | Required |
|-------|------|----------|
| `full_name` | string | Yes |
| `mobile` | string | Yes (unique, validated Iranian format) |
| `date_of_birth` | date | Yes |
| `height_cm` | float | Yes |
| `gender` | enum | Yes (male, female) |
| `notes` | text | No |

---

### 5.14 Super Admin Panel

#### Nutritionist Management

- List all nutritionists with: name, email, status, client count, created date
- Create new nutritionist: name, email, password
- Activate / deactivate nutritionist accounts
- View nutritionist's client list (read-only)

#### Food & Medication Database

- Full CRUD on food items (same interface as nutritionist, but with ability to edit/delete items created by others)
- Full CRUD on medications
- View audit log of who created/modified items

#### Platform Statistics (Basic)

- Total nutritionists (active/inactive)
- Total clients (active/inactive)
- Total food items
- Total active diet plans

---

## 6. Offline & Sync Strategy

Offline capability is **client-side only** (not nutritionist or super admin).

### What Works Offline

| Feature | Offline Capability |
|---------|-------------------|
| View active diet plan | ✅ Full read access |
| Log food intake | ✅ Queue for sync |
| Log water intake | ✅ Queue for sync |
| Log exercise | ✅ Queue for sync |
| Log sleep | ✅ Queue for sync |
| Log weight & measurements | ✅ Queue for sync |
| Log medication intake | ✅ Queue for sync |
| View messages | ✅ Cached messages only |
| Send messages | ✅ Queue for sync |
| Upload lab results | ❌ Requires internet |
| View food database | ✅ Cached subset (plan-related foods) |

### Technical Approach

- **Service Worker** for caching static assets and API responses
- **IndexedDB** for storing:
  - Active diet plan data
  - Pending log entries (food, water, exercise, sleep, weight, measurements, medications)
  - Cached messages
  - Queued outgoing messages
- **Sync Manager:**
  - On network reconnect, push all queued entries to backend API
  - Conflict resolution: **last-write-wins** with server timestamp
  - Each queued entry has a `local_id` (UUID) and `synced_at` (null until synced)
  - If sync fails for an entry, retry with exponential backoff (max 3 retries, then flag for manual retry)
- **Background Sync API** where supported, fallback to polling on app open

### Data Freshness

- Diet plan: cached on first load, refreshed on app open if online
- Tracking data: always write-through (local + API when online)
- Messages: cached last 50 messages per conversation, fetch new on open

---

## 7. Notifications

Push notifications via **PWA Web Push** (using VAPID keys).

### Notification Triggers

| Event | Recipient | Priority |
|-------|-----------|----------|
| New message received | Client / Nutritionist | High |
| New diet plan assigned | Client | High |
| Food request approved/rejected | Client | Medium |
| Meal time reminder | Client | Medium |
| Medication reminder | Client | Medium |
| Water intake reminder | Client | Low |
| Client registered (OTP sent) | Client | High |

### Implementation Notes

- Backend sends push notifications via Web Push protocol
- Client subscribes to push on first login (permission prompt)
- Notification payload includes: title, body, action URL, icon
- Meal and medication reminders are scheduled based on diet plan times
- Reminder preferences: client can enable/disable each reminder type in settings

---

## 8. Tech Stack & Architecture

### Backend — Go (Golang)

- **Framework:** Fiber or Echo (HTTP router)
- **ORM / Query Builder:** sqlc or GORM
- **Authentication:** JWT (access + refresh tokens) with OTP for clients
- **SMS Gateway:** Configurable adapter (e.g., Kavenegar, Melipayamak for Iranian SMS)
- **File Storage:** Local filesystem on Hetzner (`/data/uploads/`) with path stored in DB
- **Push Notifications:** Web Push library (github.com/SherClockHolmes/webpush-go)
- **API Design:** RESTful JSON API
- **Migration:** Atlas or golang-migrate

### Database — PostgreSQL

- Primary data store
- Full-text search for food/medication names (with Persian text search configuration)
- Indexes on: user lookups, diet plan active status, date-based queries
- Soft deletes where applicable (`deleted_at` timestamp)

### Frontend — Nuxt 4 (Vue 3)

- **PWA Plugin:** @vite-pwa/nuxt for service worker and manifest
- **UI Framework:** Minimal custom components (Tailwind CSS with RTL plugin)
- **State Management:** Pinia
- **Offline Storage:** IndexedDB via Dexie.js
- **HTTP Client:** ofetch / useFetch with offline queue wrapper
- **Charts:** Chart.js or lightweight alternative for weight/measurement trends
- **RTL:** Full RTL layout, Persian numeral display, Jalali (Shamsi) calendar

### Infrastructure

- **Hosting:** Hetzner dedicated/VPS
- **Containerization:** Docker + Docker Compose
- **Reverse Proxy:** Traefik or Nginx
- **CI/CD:** GitLab CI/CD
- **Monitoring:** Grafana + Loki (existing stack)

### Architecture Diagram

```
┌─────────────────────────────────────────────────┐
│                   Client (PWA)                   │
│        Nuxt 4 + Service Worker + IndexedDB       │
└──────────────────────┬──────────────────────────┘
                       │ HTTPS (REST JSON)
                       ▼
┌──────────────────────────────────────────────────┐
│              Traefik (Reverse Proxy)              │
└──────────────────────┬───────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────┐
│               Go API Server                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐  │
│  │   Auth   │ │   API    │ │  Push Notifier   │  │
│  │ (JWT+OTP)│ │ Handlers │ │  (Web Push)      │  │
│  └──────────┘ └──────────┘ └──────────────────┘  │
│  ┌──────────────────────────────────────────────┐ │
│  │           Business Logic Layer               │ │
│  └──────────────────────────────────────────────┘ │
└────────┬────────────────────────────┬────────────┘
         │                            │
         ▼                            ▼
┌────────────────┐          ┌─────────────────┐
│  PostgreSQL    │          │  File Storage   │
│  (Data Store)  │          │  (Hetzner Disk) │
└────────────────┘          └─────────────────┘
```

---

## 9. Data Model

### Core Tables

#### `users`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| role | enum | super_admin, nutritionist, client |
| full_name | varchar(255) | |
| email | varchar(255) | Nullable, for admin/nutritionist |
| password_hash | varchar(255) | Nullable, for admin/nutritionist |
| mobile | varchar(15) | Nullable, for client (unique) |
| date_of_birth | date | For clients |
| height_cm | float | For clients |
| gender | enum | male, female |
| nutritionist_id | UUID FK | For clients — assigned nutritionist |
| is_active | boolean | Default true |
| created_at | timestamp | |
| updated_at | timestamp | |

#### `foods`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| name | varchar(255) | Persian, indexed for full-text search |
| calories | float | Per measurement unit + amount |
| protein | float | Grams |
| carbohydrates | float | Grams |
| fat | float | Grams |
| fiber | float | Nullable |
| sugar | float | Nullable |
| sodium | float | Nullable (mg) |
| measurement_unit | enum | See units list |
| measurement_amount | float | Base amount |
| description | text | Nullable |
| is_active | boolean | Default true |
| created_by | UUID FK → users | |
| created_at | timestamp | |
| updated_at | timestamp | |

#### `food_categories` (junction table)
| Column | Type | Notes |
|--------|------|-------|
| food_id | UUID FK | |
| category | enum | breakfast, lunch, dinner, snack, fruit, beverage, supplement, other |

#### `medications`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| name | varchar(255) | |
| generic_name | varchar(255) | Nullable |
| form | enum | tablet, capsule, syrup, injection, drop, powder, other |
| dosage_unit | varchar(50) | |
| description | text | Nullable |
| is_active | boolean | Default true |
| created_by | UUID FK → users | |
| created_at | timestamp | |
| updated_at | timestamp | |

#### `diet_plans`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK → users | |
| nutritionist_id | UUID FK → users | |
| start_date | date | |
| end_date | date | |
| daily_water_target_ml | integer | Nullable |
| notes | text | Nullable |
| status | enum | active, archived |
| created_at | timestamp | |
| updated_at | timestamp | |

#### `plan_days`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| diet_plan_id | UUID FK | |
| day_number | integer | 1-indexed |
| created_at | timestamp | |

#### `meals`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| plan_day_id | UUID FK | |
| title | varchar(255) | e.g., "صبحانه" |
| scheduled_time | time | |
| display_order | integer | |
| created_at | timestamp | |

#### `meal_options`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| meal_id | UUID FK | |
| option_number | integer | |
| created_at | timestamp | |

#### `meal_option_items`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| meal_option_id | UUID FK | |
| food_id | UUID FK → foods | |
| quantity | float | |
| measurement_unit | enum | |
| notes | varchar(255) | Nullable |
| created_at | timestamp | |

#### `prescribed_medications`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| diet_plan_id | UUID FK | |
| medication_id | UUID FK → medications | |
| dosage | varchar(100) | |
| frequency | varchar(100) | |
| times | jsonb | Array of time strings |
| instructions | text | Nullable |
| start_date | date | |
| end_date | date | Nullable |
| created_at | timestamp | |

#### `exercise_recommendations`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| plan_day_id | UUID FK | |
| exercise_name | varchar(255) | |
| duration_minutes | integer | |
| description | text | Nullable |
| calories_burn_estimate | integer | Nullable |
| created_at | timestamp | |

### Tracking Tables

#### `food_logs`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK | |
| date | date | |
| meal_id | UUID FK | |
| selected_option_id | UUID FK | Nullable (null = skipped) |
| notes | text | Nullable |
| local_id | UUID | For offline sync deduplication |
| created_at | timestamp | |

#### `water_logs`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK | |
| date | date | |
| amount_ml | integer | |
| time | time | Nullable |
| local_id | UUID | |
| created_at | timestamp | |

#### `sleep_logs`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK | |
| date | date | Unique per client |
| sleep_time | timestamp | |
| wake_time | timestamp | |
| quality | enum | good, fair, poor (nullable) |
| notes | text | Nullable |
| local_id | UUID | |
| created_at | timestamp | |

#### `exercise_logs`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK | |
| date | date | |
| exercise_name | varchar(255) | |
| duration_minutes | integer | |
| calories_burned | float | Nullable |
| notes | text | Nullable |
| local_id | UUID | |
| created_at | timestamp | |

#### `medication_logs`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK | |
| date | date | |
| prescribed_medication_id | UUID FK | Nullable (null for self-reported) |
| medication_name | varchar(255) | |
| dosage | varchar(100) | |
| taken_at | time | |
| notes | text | Nullable |
| local_id | UUID | |
| created_at | timestamp | |

#### `body_measurements`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK | |
| date | date | |
| weight_kg | float | Nullable |
| waist_cm | float | Nullable |
| hip_cm | float | Nullable |
| abdomen_cm | float | Nullable |
| thigh_cm | float | Nullable |
| chest_cm | float | Nullable |
| wrist_cm | float | Nullable |
| recorded_by | UUID FK → users | Client or Nutritionist |
| local_id | UUID | |
| created_at | timestamp | |

#### `lab_results`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| client_id | UUID FK | |
| title | varchar(255) | |
| type | enum | blood_test, urine_test, thyroid, hormone, allergy, other |
| test_date | date | |
| file_path | varchar(500) | Nullable |
| link | varchar(500) | Nullable |
| notes | text | Nullable |
| created_at | timestamp | |

#### `messages`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| sender_id | UUID FK | |
| receiver_id | UUID FK | |
| content | text | Nullable |
| attachment_type | enum | Nullable (image, file) |
| attachment_path | varchar(500) | Nullable |
| sent_at | timestamp | |
| read_at | timestamp | Nullable |

#### `food_requests`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| food_name | varchar(255) | |
| description | text | Nullable |
| status | enum | pending, approved, rejected |
| rejection_reason | text | Nullable |
| requested_by | UUID FK → users (client) | |
| reviewed_by | UUID FK → users (nutritionist) | Nullable |
| created_at | timestamp | |
| updated_at | timestamp | |

#### `push_subscriptions`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| user_id | UUID FK | |
| endpoint | text | Web Push endpoint |
| p256dh | text | Public key |
| auth | text | Auth secret |
| created_at | timestamp | |

#### `notification_preferences`
| Column | Type | Notes |
|--------|------|-------|
| id | UUID | PK |
| user_id | UUID FK | Unique |
| meal_reminders | boolean | Default true |
| medication_reminders | boolean | Default true |
| water_reminders | boolean | Default true |
| message_notifications | boolean | Default true |
| created_at | timestamp | |
| updated_at | timestamp | |

---

## 10. Non-Functional Requirements

### Performance

- API response time: < 200ms for standard CRUD operations
- Diet plan load time: < 500ms (full plan with all days, meals, options)
- PWA initial load: < 3 seconds on 3G
- Offline mode activation: < 1 second after network loss

### Scale

- Support up to 50 nutritionists and 10,000 clients
- Concurrent active users: ~500
- Database size estimate: ~5 GB in first year (excluding uploaded files)
- File storage estimate: ~50 GB in first year

### Security

- All traffic over HTTPS (TLS 1.2+)
- JWT tokens with short expiry (15 min access, 30 day refresh)
- OTP rate limiting and brute force protection
- Input validation and sanitization on all endpoints
- SQL injection prevention via parameterized queries
- File upload validation (type, size, content sniffing)
- CORS restricted to app domain only
- Nutritionist can only access their own clients' data (row-level authorization)
- Passwords hashed with bcrypt (cost factor 12)

### Reliability

- Target uptime: 99.5%
- Daily automated PostgreSQL backups
- File storage backups (weekly)
- Health check endpoint for monitoring

### Maintenance

- Structured logging (JSON) to stdout, collected by Loki
- Database migrations via version-controlled migration files
- Docker-based deployment for reproducibility
- GitLab CI/CD pipeline for automated testing and deployment

---

## 11. Decision Log

| # | Decision | Alternatives Considered | Reason |
|---|----------|------------------------|--------|
| 1 | Shared food/medication database across all nutritionists | Per-nutritionist isolated databases | Reduces duplication, builds a richer database over time, simpler to maintain |
| 2 | Client-only offline support | Full offline for all roles | Keeps complexity manageable; nutritionists typically work with stable internet |
| 3 | OTP via SMS for client login | Email/password, magic link | Most Iranian clients prefer SMS; no email dependency; simpler UX on mobile |
| 4 | Email/password for nutritionist login | OTP for all roles | Nutritionists need stable sessions for long workflows; email is more professional |
| 5 | Super admin required for nutritionist creation | Self-registration with approval | More control over who joins the platform; prevents spam |
| 6 | Single active diet plan per client | Multiple concurrent plans | Avoids confusion; one plan with all aspects (food, exercise, medication) is cleaner |
| 7 | Polling for chat (not WebSocket) | WebSocket, SSE | Simpler to implement; acceptable latency for this use case; works better offline |
| 8 | Files stored on Hetzner local disk | MinIO, S3 | Simpler setup; sufficient for current scale; can migrate later if needed |
| 9 | Food request goes to nutritionist first | Direct to super admin | Nutritionist knows the client's context and can add it immediately |
| 10 | PWA (not native app) | React Native, Flutter | Single codebase, no app store dependency, easier updates, good offline via Service Worker |
| 11 | Persian-only, no i18n | Multi-language from start | YAGNI — target audience is Iranian; add i18n later if needed |

---

## 12. Assumptions

1. An SMS gateway with Persian support (e.g., Kavenegar) will be available and configured
2. The Hetzner server has sufficient disk space for file uploads (~50 GB first year)
3. Users have modern mobile browsers supporting Service Workers and Web Push (Chrome, Safari 16+)
4. Persian full-text search in PostgreSQL will be configured using a custom text search configuration or pg_trgm extension
5. Jalali (Shamsi) calendar conversion will be handled on the frontend using a library like `jalaali-js`
6. All nutritional values in the food database are entered manually by nutritionists (no external API integration)
7. The platform will not handle any payment or financial transactions
8. Nutritionists are verified professionals — the platform does not validate credentials
9. SMS costs are borne by the platform operator
10. No GDPR or specific data protection compliance is required initially (Iranian domestic use)
