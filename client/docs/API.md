# NutriTrack REST API Reference

> **Purpose:** This document is the authoritative reference for building a client-side application against the NutriTrack backend. It covers every endpoint, its required auth role, request shape, and exact response shape.

---

## Table of Contents

1. [Overview](#1-overview)
2. [Authentication & Roles](#2-authentication--roles)
3. [Common Response Envelope](#3-common-response-envelope)
4. [Error Codes](#4-error-codes)
5. [Auth Endpoints](#5-auth-endpoints)
6. [User / Me](#6-user--me)
7. [Avatar Upload](#7-avatar-upload)
8. [Admin — Platform Stats](#8-admin--platform-stats)
9. [Admin — Nutritionist Management](#9-admin--nutritionist-management)
10. [Client Management (Nutritionist)](#10-client-management-nutritionist)
11. [Foods](#11-foods)
12. [Food Categories](#12-food-categories)
13. [Medications](#13-medications)
14. [Diet Plans](#14-diet-plans)
15. [Tracking](#15-tracking)
16. [Lab Results](#16-lab-results)
17. [Messages (Chat)](#17-messages-chat)
18. [Food Requests](#18-food-requests)
19. [Push Subscriptions](#19-push-subscriptions)
20. [Notification Preferences](#20-notification-preferences)

---

## 1. Overview

| Property | Value |
|---|---|
| Base URL | `http(s)://<host>/api/v1` |
| Content-Type (JSON) | `application/json` |
| Content-Type (file upload) | `multipart/form-data` |
| Static files (avatars, uploads) | `GET /uploads/<filename>` |

All dates use `"YYYY-MM-DD"` strings.  
All timestamps are RFC 3339 (`"2025-06-01T12:00:00Z"`).  
All IDs are UUID v4 strings.

---

## 2. Authentication & Roles

The API uses **JWT Bearer tokens**.

```
Authorization: Bearer <access_token>
```

Three roles exist:

| Role | String value | How to obtain token |
|---|---|---|
| Super Admin | `super_admin` | `POST /auth/login` (email + password) |
| Nutritionist | `nutritionist` | `POST /auth/login` (email + password) |
| Client | `client` | `POST /auth/otp/verify` (mobile + OTP) |

Tokens expire. Use `POST /auth/refresh` to get a new access token.

Endpoint access matrix legend used throughout this doc:
- **Public** — no token required
- **Any Auth** — any valid token
- **Client** — role must be `client`
- **Nutritionist** — role must be `nutritionist`
- **Super Admin** — role must be `super_admin`

---

## 3. Common Response Envelope

### Success — single object
```json
{
  "data": { ... }
}
```

### Success — paginated list
```json
{
  "data": [ ... ],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}
```

### Pagination query params (on all list endpoints)
| Param | Type | Default | Max | Description |
|---|---|---|---|---|
| `page` | integer | `1` | — | 1-based page number |
| `page_size` | integer | `20` | `100` | Items per page |

### Error response
```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable message"
}
```

### HTTP status → domain error mapping
| HTTP Status | When |
|---|---|
| 200 | Success |
| 201 | Resource created |
| 204 | Success, no body |
| 302 | Redirect (lab result link) |
| 401 | Invalid / expired / revoked token, bad OTP |
| 403 | Forbidden (wrong role, not your resource) |
| 404 | Resource not found |
| 409 | Conflict (duplicate, already active plan) |
| 413 | File too large |
| 422 | Validation error |
| 429 | Rate limit exceeded |
| 500 | Internal error |

---

## 4. Error Codes

| Code | HTTP | Meaning |
|---|---|---|
| `VALIDATION_ERROR` | 422 | Missing/invalid field |
| `INVALID_MOBILE` | 422 | Mobile format invalid |
| `INVALID_FILE_TYPE` | 422 | File MIME type rejected |
| `INVALID_CREDENTIALS` | 401 | Wrong email/password |
| `INVALID_TOKEN` | 401 | JWT malformed or expired |
| `TOKEN_REVOKED` | 401 | JWT was revoked |
| `OTP_INVALID` | 401 | OTP code wrong |
| `OTP_EXPIRED` | 401 | OTP code expired |
| `OTP_MAX_ATTEMPTS` | 401 | Too many wrong OTP attempts |
| `OTP_RATE_LIMIT` | 429 | OTP send rate limit |
| `RATE_LIMIT_EXCEEDED` | 429 | General rate limit |
| `UNAUTHORIZED` | 401 | No or invalid auth |
| `FORBIDDEN` | 403 | Authenticated but not allowed |
| `FOOD_REQUEST_NOT_OWNED` | 403 | Food request belongs to another nutritionist |
| `NOT_FOUND` | 404 | Generic not found |
| `USER_NOT_FOUND` | 404 | — |
| `FOOD_NOT_FOUND` | 404 | — |
| `MEDICATION_NOT_FOUND` | 404 | — |
| `DIET_PLAN_NOT_FOUND` | 404 | — |
| `LAB_RESULT_NOT_FOUND` | 404 | — |
| `MESSAGE_NOT_FOUND` | 404 | — |
| `FOOD_REQUEST_NOT_FOUND` | 404 | — |
| `NOTIFICATION_PREFERENCE_NOT_FOUND` | 404 | — |
| `USER_ALREADY_EXISTS` | 409 | Mobile/email already registered |
| `PLAN_ALREADY_ACTIVE` | 409 | Client already has an active plan |
| `FOOD_REQUEST_ALREADY_PROCESSED` | 409 | Request was already approved/rejected |
| `FILE_TOO_LARGE` | 413 | File exceeds limit |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## 5. Auth Endpoints

### 5.1 Login (Super Admin / Nutritionist)

```
POST /auth/login
```

**Access:** Public

**Request body:**
```json
{
  "email": "admin@nutritrack.com",
  "password": "secret123"
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `email` | string | yes | valid email |
| `password` | string | yes | min 6 chars |

**Response `200`:**
```json
{
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>",
    "token_type": "Bearer",
    "user_id": "uuid",
    "role": "super_admin"
  }
}
```

---

### 5.2 Send OTP (Client)

```
POST /auth/otp/send
```

**Access:** Public  
**Rate limit:** 1 request per 60 seconds per IP

**Request body:**
```json
{
  "mobile": "09123456789"
}
```

**Response `200`:**
```json
{
  "data": { "message": "کد تأیید ارسال شد" }
}
```

---

### 5.3 Verify OTP (Client)

```
POST /auth/otp/verify
```

**Access:** Public

**Request body:**
```json
{
  "mobile": "09123456789",
  "code": "123456"
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `mobile` | string | yes | — |
| `code` | string | yes | exactly 6 chars |

**Response `200`:**
```json
{
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>",
    "token_type": "Bearer",
    "user_id": "uuid",
    "role": "client"
  }
}
```

---

### 5.4 Refresh Token

```
POST /auth/refresh
```

**Access:** Public

**Request body:**
```json
{
  "refresh_token": "<jwt>"
}
```

**Response `200`:**
```json
{
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>",
    "token_type": "Bearer",
    "user_id": "uuid",
    "role": "client"
  }
}
```

---

### 5.5 Logout

```
POST /auth/logout
```

**Access:** Any Auth

**Request body:**
```json
{
  "refresh_token": "<jwt>"
}
```

**Response `200`:**
```json
{
  "data": { "message": "با موفقیت خارج شدید" }
}
```

Both the access token (current session) and refresh token are revoked.

---

## 6. User / Me

### 6.1 Get Current User

```
GET /auth/me
```

**Access:** Any Auth

**Response `200`:**
```json
{
  "data": {
    "user_id": "uuid",
    "role": "client"
  }
}
```

> For full profile data (name, avatar, etc.) use the appropriate client or nutritionist profile endpoint.

---

## 7. Avatar Upload

### 7.1 Upload Avatar

```
PUT /users/:id/avatar
```

**Access:** Any Auth (service-level: users can only update their own; super_admin can update any)  
**Content-Type:** `multipart/form-data`  
**Max file size:** 5 MB  
**Accepted MIME types:** JPEG, PNG, WebP, GIF (validated via magic bytes)

| Form field | Type | Required |
|---|---|---|
| `avatar` | file | yes |

**Path params:**
| Param | Type | Description |
|---|---|---|
| `id` | UUID | Target user's ID |

**Response `200`:**
```json
{
  "data": {
    "id": "uuid",
    "mobile": "09123456789",
    "first_name": "Ali",
    "last_name": "Rezaei",
    "full_name": "Ali Rezaei",
    "gender": "male",
    "birth_date": "1990-05-15",
    "height": 175.0,
    "weight": 70.5,
    "bmi": 23.02,
    "avatar_url": "/uploads/avatars/uuid.jpg",
    "is_active": true,
    "nutritionist_id": "uuid-or-null",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

## 8. Admin — Platform Stats

### 8.1 Get Platform Statistics

```
GET /admin/stats
```

**Access:** Super Admin

**Response `200`:**
```json
{
  "data": {
    "total_nutritionists": 12,
    "active_nutritionists": 10,
    "inactive_nutritionists": 2,
    "total_clients": 350,
    "total_foods": 800,
    "active_diet_plans": 230
  }
}
```

---

## 9. Admin — Nutritionist Management

All endpoints require **Super Admin** role.

### 9.1 Create Nutritionist

```
POST /admin/nutritionists
```

**Request body:**
```json
{
  "email": "nutri@clinic.com",
  "password": "StrongPass8",
  "first_name": "Sara",
  "last_name": "Ahmadi",
  "mobile": "09111234567"
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `email` | string | yes | valid email |
| `password` | string | yes | min 8 chars |
| `first_name` | string | yes | — |
| `last_name` | string | yes | — |
| `mobile` | string | yes | — |

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "email": "nutri@clinic.com",
    "mobile": "09111234567",
    "first_name": "Sara",
    "last_name": "Ahmadi",
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 9.2 List Nutritionists

```
GET /admin/nutritionists?page=1&page_size=20
```

**Response `200`:** Paginated array of [nutritionist objects](#nutritionist-object)

---

### 9.3 Get Nutritionist

```
GET /admin/nutritionists/:id
```

**Path params:** `id` — UUID

**Response `200`:** Single [nutritionist object](#nutritionist-object)

---

### 9.4 Update Nutritionist

```
PATCH /admin/nutritionists/:id
```

**Request body:**
```json
{
  "first_name": "Sara",
  "last_name": "Karimi"
}
```

All fields are optional; only provided fields are updated.

**Response `200`:** Updated [nutritionist object](#nutritionist-object)

---

### 9.5 Set Nutritionist Status

```
PATCH /admin/nutritionists/:id/status
```

**Request body:**
```json
{
  "is_active": false
}
```

**Response `200`:**
```json
{
  "data": { "message": "وضعیت متخصص تغذیه با موفقیت به‌روز شد" }
}
```

---

### 9.6 Get Nutritionist's Clients

```
GET /admin/nutritionists/:id/clients?page=1&page_size=20
```

**Response `200`:** Paginated array of [client objects](#client-object)

---

#### Nutritionist Object

```json
{
  "id": "uuid",
  "email": "nutri@clinic.com",
  "mobile": "09111234567",
  "first_name": "Sara",
  "last_name": "Ahmadi",
  "is_active": true,
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

## 10. Client Management (Nutritionist)

All endpoints require **Nutritionist** role. The authenticated nutritionist ID is inferred from the JWT — never sent in the body.

### 10.1 Register Client

```
POST /clients
```

**Request body:**
```json
{
  "mobile": "09121234567",
  "first_name": "Mohammad",
  "last_name": "Hosseini",
  "gender": "male",
  "birth_date": "1992-03-10",
  "height": 178.0,
  "weight": 82.5
}
```

| Field | Type | Required | Notes |
|---|---|---|---|
| `mobile` | string | yes | — |
| `first_name` | string | yes | — |
| `last_name` | string | yes | — |
| `gender` | string | no | `"male"` or `"female"` |
| `birth_date` | string | no | `"YYYY-MM-DD"` |
| `height` | float | no | cm |
| `weight` | float | no | kg |

**Response `201`:** [Client object](#client-object)

---

### 10.2 List Clients

```
GET /clients?page=1&page_size=20
```

**Response `200`:** Paginated array of [client objects](#client-object)

---

### 10.3 Get Client Profile

```
GET /clients/:id
```

**Response `200`:** [Client object](#client-object)

---

### 10.4 Update Client

```
PATCH /clients/:id
```

**Request body:** Same fields as [Register Client](#101-register-client), all optional.

**Response `200`:** Updated [client object](#client-object)

---

### 10.5 Set Client Status

```
PATCH /clients/:id/status
```

**Request body:**
```json
{
  "is_active": false
}
```

**Response `200`:** Updated [client object](#client-object)

---

#### Client Object

```json
{
  "id": "uuid",
  "mobile": "09121234567",
  "first_name": "Mohammad",
  "last_name": "Hosseini",
  "full_name": "Mohammad Hosseini",
  "gender": "male",
  "birth_date": "1992-03-10",
  "height": 178.0,
  "weight": 82.5,
  "bmi": 26.04,
  "avatar_url": "/uploads/avatars/uuid.jpg",
  "is_active": true,
  "nutritionist_id": "uuid",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

> `bmi` is computed server-side. `avatar_url` is null until an avatar is uploaded.

---

## 11. Foods

### 11.1 Create Food

```
POST /foods
```

**Access:** Any Auth (nutritionists create public foods; clients create personal foods — role logic in service)

**Request body:**
```json
{
  "name": "برنج سفید",
  "unit": "گرم",
  "calories": 130.0,
  "protein": 2.7,
  "carbohydrate": 28.2,
  "fat": 0.3,
  "fiber": 0.4,
  "sugar": 0.0,
  "sodium": 1.0,
  "amount": 100.0,
  "category_ids": ["uuid1", "uuid2"]
}
```

| Field | Type | Required | Constraints |
|---|---|---|---|
| `name` | string | yes | — |
| `unit` | string | yes | e.g. `"گرم"`, `"عدد"`, `"ml"` |
| `calories` | float | yes | ≥ 0 |
| `protein` | float | yes | ≥ 0 |
| `carbohydrate` | float | yes | ≥ 0 |
| `fat` | float | yes | ≥ 0 |
| `fiber` | float | yes | ≥ 0 |
| `sugar` | float | yes | ≥ 0 |
| `sodium` | float | yes | ≥ 0 |
| `amount` | float | yes | ≥ 0 — portion size corresponding to the macros |
| `category_ids` | UUID[] | no | Food category UUIDs |

**Response `201`:** [Food object](#food-object)

---

### 11.2 Search Foods

```
GET /foods?q=برنج&category_id=uuid&page=1&page_size=20
```

**Access:** Any Auth

**Query params:**
| Param | Type | Required | Description |
|---|---|---|---|
| `q` | string | no | Search query (name substring) |
| `category_id` | UUID | no | Filter by category |
| `page` | int | no | default 1 |
| `page_size` | int | no | default 20, max 100 |

**Response `200`:** Paginated array of [food objects](#food-object)

---

### 11.3 Get Food

```
GET /foods/:id
```

**Access:** Any Auth

**Response `200`:** [Food object](#food-object)

---

### 11.4 Update Food

```
PATCH /foods/:id
```

**Access:** Any Auth (service enforces ownership — creator or super_admin)

**Request body:** Same as [Create Food](#111-create-food), all fields required (full replace).

**Response `200`:** Updated [food object](#food-object)

---

### 11.5 Delete Food

```
DELETE /foods/:id
```

**Access:** Any Auth (creator or super_admin)

**Response `204`:** No body

---

### 11.6 Admin — Search All Foods

```
GET /admin/foods?q=&page=1&page_size=20
```

**Access:** Super Admin  
Same query params as [Search Foods](#112-search-foods). Returns all foods regardless of creator.

---

### 11.7 Admin — Force Delete Food

```
DELETE /admin/foods/:id
```

**Access:** Super Admin  
**Response `204`:** No body

---

#### Food Object

```json
{
  "id": "uuid",
  "name": "برنج سفید",
  "unit": "گرم",
  "amount": 100.0,
  "calories": 130.0,
  "protein": 2.7,
  "carbohydrate": 28.2,
  "fat": 0.3,
  "fiber": 0.4,
  "sugar": 0.0,
  "sodium": 1.0,
  "is_active": true,
  "created_by": "uuid-or-null",
  "categories": [
    { "id": "uuid", "name": "غلات" }
  ],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## 12. Food Categories

### 12.1 List All Categories

```
GET /food-categories
```

**Access:** Any Auth

**Response `200`:**
```json
{
  "data": [
    { "id": "uuid", "name": "غلات" },
    { "id": "uuid", "name": "لبنیات" }
  ]
}
```

---

### 12.2 Admin — Create Category

```
POST /admin/food-categories
```

**Access:** Super Admin

**Request body:**
```json
{
  "name": "حبوبات"
}
```

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "name": "حبوبات",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 12.3 Admin — Delete Category

```
DELETE /admin/food-categories/:id
```

**Access:** Super Admin

**Response `200`:**
```json
{
  "data": { "message": "دسته‌بندی با موفقیت حذف شد" }
}
```

---

## 13. Medications

### 13.1 Create Medication

```
POST /medications
```

**Access:** Any Auth

**Request body:**
```json
{
  "name": "متفورمین",
  "description": "داروی دیابت",
  "unit": "mg"
}
```

| Field | Type | Required |
|---|---|---|
| `name` | string | yes |
| `description` | string | no |
| `unit` | string | yes |

**Response `201`:** [Medication object](#medication-object)

---

### 13.2 Search Medications

```
GET /medications?q=متفورمین&page=1&page_size=20
```

**Access:** Any Auth

**Response `200`:** Paginated array of [medication objects](#medication-object)

---

### 13.3 Get Medication

```
GET /medications/:id
```

**Access:** Any Auth  
**Response `200`:** [Medication object](#medication-object)

---

### 13.4 Update Medication

```
PATCH /medications/:id
```

**Access:** Any Auth (creator or super_admin)  
**Request body:** Same as create (all fields required).  
**Response `200`:** Updated [medication object](#medication-object)

---

### 13.5 Delete Medication

```
DELETE /medications/:id
```

**Access:** Any Auth (creator or super_admin)  
**Response `204`:** No body

---

### 13.6 Admin — Search All Medications

```
GET /admin/medications?q=&page=1&page_size=20
```

**Access:** Super Admin

---

### 13.7 Admin — Force Delete Medication

```
DELETE /admin/medications/:id
```

**Access:** Super Admin  
**Response `204`:** No body

---

#### Medication Object

```json
{
  "id": "uuid",
  "name": "متفورمین",
  "description": "داروی دیابت",
  "unit": "mg",
  "is_active": true,
  "created_by": "uuid-or-null",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## 14. Diet Plans

Diet plans are hierarchical: **Plan → Days → Meals → Options → Items** (foods).  
Each day can also have **Exercise Recommendations** and **Medication Prescriptions**.

### 14.1 Create Plan

```
POST /clients/:id/plans
```

**Access:** Nutritionist  
**Path params:** `id` — client UUID

**Request body:**
```json
{
  "title": "برنامه کاهش وزن",
  "start_date": "2025-07-01",
  "end_date": "2025-07-31",
  "notes": "رژیم کم‌کالری",
  "daily_water_target_ml": 2000
}
```

| Field | Type | Required |
|---|---|---|
| `title` | string | no |
| `start_date` | string | yes | `"YYYY-MM-DD"` |
| `end_date` | string | yes | `"YYYY-MM-DD"` |
| `notes` | string | no |
| `daily_water_target_ml` | int | no |

**Response `201`:** [Diet plan flat object](#diet-plan-flat-object)

> A push notification is sent to the client automatically.

---

### 14.2 List Client Plans

```
GET /clients/:id/plans?page=1&page_size=20
```

**Access:** Nutritionist or Client (client can only see their own)  
**Response `200`:** Paginated array of [diet plan flat objects](#diet-plan-flat-object)

---

### 14.3 Get Full Plan

```
GET /plans/:id
```

**Access:** Any Auth (service enforces ownership)

**Response `200`:** [Diet plan full object](#diet-plan-full-object) — includes the complete days/meals/options/items tree.

---

### 14.4 Get Active Plan (Client)

```
GET /plans/active
```

**Access:** Any Auth (client gets their own active plan; nutritionist gets the active plan for a context derived from the token)

**Response `200`:** [Diet plan full object](#diet-plan-full-object)  
**Response `404`:** `DIET_PLAN_NOT_FOUND` if no active plan exists.

---

### 14.5 Update Plan

```
PATCH /plans/:id
```

**Access:** Nutritionist (plan owner)

**Request body:** Same fields as Create Plan, all optional.

**Response `200`:** Updated [diet plan flat object](#diet-plan-flat-object)

---

### 14.6 Delete Plan

```
DELETE /plans/:id
```

**Access:** Nutritionist (plan owner)  
**Response `204`:** No body

---

### 14.7 Add Day

```
POST /plans/:id/days
```

**Access:** Nutritionist

**Request body:**
```json
{
  "day_number": 1
}
```

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "plan_id": "uuid",
    "day_number": 1,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 14.8 Delete Day

```
DELETE /plans/:id/days/:day_id
```

**Access:** Nutritionist  
**Response `204`:** No body

---

### 14.9 Add Meal

```
POST /plans/:id/days/:day_id/meals
```

**Access:** Nutritionist

**Request body:**
```json
{
  "title": "صبحانه",
  "scheduled_time": "08:00",
  "display_order": 1
}
```

| Field | Type | Required |
|---|---|---|
| `title` | string | no |
| `scheduled_time` | string | no | `"HH:MM"` |
| `display_order` | int | no |

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "day_id": "uuid",
    "title": "صبحانه",
    "scheduled_time": "08:00",
    "display_order": 1,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 14.10 Delete Meal

```
DELETE /plans/:id/days/:day_id/meals/:meal_id
```

**Access:** Nutritionist  
**Response `204`:** No body

---

### 14.11 Add Option

```
POST /plans/:id/days/:day_id/meals/:meal_id/options
```

**Access:** Nutritionist  
Each meal can have multiple options (e.g., "Option A" or "Option B").

**Request body:**
```json
{
  "option_number": 1
}
```

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "meal_id": "uuid",
    "option_number": 1,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 14.12 Delete Option

```
DELETE /plans/:id/days/:day_id/meals/:meal_id/options/:option_id
```

**Access:** Nutritionist  
**Response `204`:** No body

---

### 14.13 Add Item to Option

```
POST /plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items
```

**Access:** Nutritionist

**Request body:**
```json
{
  "food_id": "uuid",
  "quantity": 150.0,
  "unit": "گرم",
  "notes": "بدون نمک"
}
```

| Field | Type | Required |
|---|---|---|
| `food_id` | UUID | yes |
| `quantity` | float | yes |
| `unit` | string | no |
| `notes` | string | no |

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "option_id": "uuid",
    "food_id": "uuid",
    "quantity": 150.0,
    "unit": "گرم",
    "notes": "بدون نمک",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 14.14 Remove Item

```
DELETE /plans/:id/days/:day_id/meals/:meal_id/options/:option_id/items/:item_id
```

**Access:** Nutritionist  
**Response `204`:** No body

---

### 14.15 Add Exercise to Day

```
POST /plans/:id/days/:day_id/exercises
```

**Access:** Nutritionist

**Request body:**
```json
{
  "exercise_name": "پیاده‌روی",
  "duration_minutes": 30,
  "description": "پیاده‌روی سبک در هوای آزاد",
  "calories_burn_estimate": 150
}
```

| Field | Type | Required |
|---|---|---|
| `exercise_name` | string | yes |
| `duration_minutes` | int | no |
| `description` | string | no |
| `calories_burn_estimate` | int | no |

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "day_id": "uuid",
    "exercise_name": "پیاده‌روی",
    "duration_minutes": 30,
    "description": "پیاده‌روی سبک در هوای آزاد",
    "calories_burn_estimate": 150,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 14.16 Remove Exercise

```
DELETE /plans/:id/days/:day_id/exercises/:exercise_id
```

**Access:** Nutritionist  
**Response `204`:** No body

---

### 14.17 Add Prescription to Day

```
POST /plans/:id/days/:day_id/prescriptions
```

**Access:** Nutritionist

**Request body:**
```json
{
  "medication_id": "uuid",
  "dosage": "500mg",
  "frequency": "daily",
  "times": ["08:00", "20:00"],
  "instructions": "با غذا مصرف شود",
  "start_date": "2025-07-01",
  "end_date": "2025-07-31"
}
```

| Field | Type | Required |
|---|---|---|
| `medication_id` | UUID | yes |
| `dosage` | string | no |
| `frequency` | string | no | e.g. `"daily"`, `"twice_daily"` |
| `times` | string[] | no | `["HH:MM", ...]` |
| `instructions` | string | no |
| `start_date` | string | no | `"YYYY-MM-DD"` |
| `end_date` | string | no | `"YYYY-MM-DD"` |

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "day_id": "uuid",
    "medication_id": "uuid",
    "dosage": "500mg",
    "frequency": "daily",
    "times": ["08:00", "20:00"],
    "instructions": "با غذا مصرف شود",
    "start_date": "2025-07-01T00:00:00Z",
    "end_date": "2025-07-31T00:00:00Z",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

### 14.18 Remove Prescription

```
DELETE /plans/:id/days/:day_id/prescriptions/:prescription_id
```

**Access:** Nutritionist  
**Response `204`:** No body

---

#### Diet Plan Flat Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "nutritionist_id": "uuid",
  "title": "برنامه کاهش وزن",
  "start_date": "2025-07-01",
  "end_date": "2025-07-31",
  "notes": "رژیم کم‌کالری",
  "daily_water_target_ml": 2000,
  "status": "active",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

> `status` values: `"draft"` | `"active"` | `"archived"`

---

#### Diet Plan Full Object

The full object extends the flat object with a nested `days` array:

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "nutritionist_id": "uuid",
  "title": "برنامه کاهش وزن",
  "start_date": "2025-07-01",
  "end_date": "2025-07-31",
  "notes": "رژیم کم‌کالری",
  "daily_water_target_ml": 2000,
  "status": "active",
  "days": [
    {
      "id": "uuid",
      "day_number": 1,
      "total_range": {
        "min": { "calories": 1400, "protein": 80, "carbs": 160, "fat": 45, "fiber": 20 },
        "max": { "calories": 1600, "protein": 100, "carbs": 200, "fat": 55, "fiber": 30 }
      },
      "meals": [
        {
          "id": "uuid",
          "title": "صبحانه",
          "scheduled_time": "08:00",
          "display_order": 1,
          "total_range": { "min": { ... }, "max": { ... } },
          "options": [
            {
              "id": "uuid",
              "option_number": 1,
              "totals": { "calories": 350, "protein": 20, "carbs": 50, "fat": 10, "fiber": 5 },
              "items": [
                {
                  "id": "uuid",
                  "option_id": "uuid",
                  "food_id": "uuid",
                  "quantity": 150.0,
                  "unit": "گرم",
                  "notes": "",
                  "computed": { "calories": 195, "protein": 4.05, "carbs": 42.3, "fat": 0.45, "fiber": 0.6 },
                  "food": { "id": "uuid", "name": "برنج سفید", "unit": "گرم" },
                  "created_at": "2025-01-01T00:00:00Z"
                }
              ]
            }
          ]
        }
      ],
      "exercises": [
        {
          "id": "uuid",
          "exercise_name": "پیاده‌روی",
          "duration_minutes": 30,
          "description": "پیاده‌روی سبک",
          "calories_burn_estimate": 150,
          "created_at": "2025-01-01T00:00:00Z"
        }
      ],
      "prescriptions": [
        {
          "id": "uuid",
          "medication_id": "uuid",
          "dosage": "500mg",
          "frequency": "daily",
          "times": ["08:00"],
          "instructions": "با غذا مصرف شود",
          "start_date": "2025-07-01T00:00:00Z",
          "end_date": "2025-07-31T00:00:00Z",
          "medication": { "id": "uuid", "name": "متفورمین", "unit": "mg" },
          "created_at": "2025-01-01T00:00:00Z"
        }
      ]
    }
  ],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

---

## 15. Tracking

Tracking data is submitted by **clients** and viewed by clients and their nutritionists.

All write endpoints use `local_id` (a client-generated UUID string) as an idempotency key — submitting the same `local_id` twice is safe (returns the existing record).

If `logged_at` / `measured_at` is omitted, the server uses the current Tehran time.

### 15.1 Log Food

```
POST /tracking/food
```

**Access:** Client

**Request body:**
```json
{
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T08:30:00Z",
  "food_id": "uuid-or-null",
  "food_name": "برنج سفید",
  "quantity": 150.0,
  "unit": "گرم",
  "calories": 195.0,
  "protein": 4.05,
  "carbs": 42.3,
  "fat": 0.45,
  "notes": ""
}
```

| Field | Type | Required |
|---|---|---|
| `local_id` | string | yes |
| `logged_at` | RFC 3339 | no | defaults to now |
| `food_id` | UUID | no | link to food catalogue |
| `food_name` | string | yes |
| `quantity` | float | yes |
| `unit` | string | yes |
| `calories` | float | yes |
| `protein` | float | no |
| `carbs` | float | no |
| `fat` | float | no |
| `notes` | string | no |

**Response `200`:** [Food log object](#food-log-object)

---

### 15.2 Log Water

```
POST /tracking/water
```

**Access:** Client

**Request body:**
```json
{
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T09:00:00Z",
  "amount_ml": 250,
  "notes": ""
}
```

| Field | Type | Required |
|---|---|---|
| `local_id` | string | yes |
| `logged_at` | RFC 3339 | no |
| `amount_ml` | int | yes |
| `notes` | string | no |

**Response `200`:** [Water log object](#water-log-object)

---

### 15.3 Log Sleep

```
POST /tracking/sleep
```

**Access:** Client

**Request body:**
```json
{
  "local_id": "client-uuid-v4",
  "sleep_start": "2025-07-01T23:00:00Z",
  "sleep_end": "2025-07-02T07:00:00Z",
  "duration_minutes": 480,
  "quality": 4,
  "notes": ""
}
```

| Field | Type | Required | Notes |
|---|---|---|---|
| `local_id` | string | yes | — |
| `sleep_start` | RFC 3339 | yes | — |
| `sleep_end` | RFC 3339 | yes | — |
| `duration_minutes` | int | no | auto-computed if 0 |
| `quality` | int | no | 1–5 scale |
| `notes` | string | no | — |

**Response `200`:** [Sleep log object](#sleep-log-object)

---

### 15.4 Log Exercise

```
POST /tracking/exercise
```

**Access:** Client

**Request body:**
```json
{
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T07:00:00Z",
  "exercise_name": "دویدن",
  "duration_minutes": 45,
  "calories_burned": 400,
  "notes": ""
}
```

| Field | Type | Required |
|---|---|---|
| `local_id` | string | yes |
| `logged_at` | RFC 3339 | no |
| `exercise_name` | string | yes |
| `duration_minutes` | int | no |
| `calories_burned` | int | no |
| `notes` | string | no |

**Response `200`:** [Exercise log object](#exercise-log-object)

---

### 15.5 Log Medication

```
POST /tracking/medication
```

**Access:** Client

**Request body:**
```json
{
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T08:00:00Z",
  "medication_id": "uuid-or-null",
  "medication_name": "متفورمین",
  "dosage": "500mg",
  "notes": ""
}
```

| Field | Type | Required |
|---|---|---|
| `local_id` | string | yes |
| `logged_at` | RFC 3339 | no |
| `medication_id` | UUID | no |
| `medication_name` | string | yes |
| `dosage` | string | no |
| `notes` | string | no |

**Response `200`:** [Medication log object](#medication-log-object)

---

### 15.6 Log Body Measurement

```
POST /tracking/body
```

**Access:** Client

**Request body:**
```json
{
  "local_id": "client-uuid-v4",
  "measured_at": "2025-07-01T07:00:00Z",
  "weight_kg": 82.5,
  "height_cm": 178.0,
  "waist_cm": 88.0,
  "hip_cm": 96.0,
  "chest_cm": 100.0,
  "arm_cm": 32.0,
  "notes": ""
}
```

All measurement fields are optional (null means not measured today).

| Field | Type | Required |
|---|---|---|
| `local_id` | string | yes |
| `measured_at` | RFC 3339 | no |
| `weight_kg` | float | no |
| `height_cm` | float | no |
| `waist_cm` | float | no |
| `hip_cm` | float | no |
| `chest_cm` | float | no |
| `arm_cm` | float | no |
| `notes` | string | no |

**Response `200`:** [Body measurement object](#body-measurement-object)

---

### 15.7 Bulk Sync

```
POST /tracking/sync
```

**Access:** Client  
Used by mobile apps to sync offline data in one request.

**Request body:**
```json
{
  "entries": [
    {
      "type": "food",
      "payload": { /* same as LogFood body */ }
    },
    {
      "type": "water",
      "payload": { /* same as LogWater body */ }
    },
    {
      "type": "sleep",
      "payload": { /* same as LogSleep body */ }
    },
    {
      "type": "exercise",
      "payload": { /* same as LogExercise body */ }
    },
    {
      "type": "medication",
      "payload": { /* same as LogMedication body */ }
    },
    {
      "type": "body",
      "payload": { /* same as LogBody body */ }
    }
  ]
}
```

`type` values: `"food"` | `"water"` | `"sleep"` | `"exercise"` | `"medication"` | `"body"`

**Response `200`:**
```json
{
  "data": {
    "synced": 6,
    "errors": []
  }
}
```

---

### 15.8 Get Tracking (Nutritionist / Client)

```
GET /clients/:id/tracking?type=food&date=2025-07-01
```

**Access:** Nutritionist or Client (client can only read own data)  
**Path params:** `id` — client UUID

**Query params:**
| Param | Type | Required | Values |
|---|---|---|---|
| `type` | string | yes | `food` \| `water` \| `sleep` \| `exercise` \| `medication` \| `body` |
| `date` | string | no | `"YYYY-MM-DD"`, defaults to today (Tehran) |

**Response `200`:**
Returns an array of the matching log objects for the given day.

```json
{
  "data": [ /* array of log objects matching the `type` */ ]
}
```

---

#### Food Log Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T08:30:00Z",
  "logged_date": "2025-07-01",
  "food_id": "uuid-or-null",
  "food_name": "برنج سفید",
  "quantity": 150.0,
  "unit": "گرم",
  "calories": 195.0,
  "protein": 4.05,
  "carbs": 42.3,
  "fat": 0.45,
  "notes": "",
  "created_at": "2025-07-01T08:31:00Z"
}
```

#### Water Log Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T09:00:00Z",
  "logged_date": "2025-07-01",
  "amount_ml": 250,
  "notes": "",
  "created_at": "2025-07-01T09:01:00Z"
}
```

#### Sleep Log Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "local_id": "client-uuid-v4",
  "logged_date": "2025-07-02",
  "sleep_start": "2025-07-01T23:00:00Z",
  "sleep_end": "2025-07-02T07:00:00Z",
  "duration_minutes": 480,
  "quality": 4,
  "notes": "",
  "created_at": "2025-07-02T07:05:00Z"
}
```

#### Exercise Log Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T07:00:00Z",
  "logged_date": "2025-07-01",
  "exercise_name": "دویدن",
  "duration_minutes": 45,
  "calories_burned": 400,
  "notes": "",
  "created_at": "2025-07-01T07:05:00Z"
}
```

#### Medication Log Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "local_id": "client-uuid-v4",
  "logged_at": "2025-07-01T08:00:00Z",
  "logged_date": "2025-07-01",
  "medication_id": "uuid-or-null",
  "medication_name": "متفورمین",
  "dosage": "500mg",
  "notes": "",
  "created_at": "2025-07-01T08:01:00Z"
}
```

#### Body Measurement Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "local_id": "client-uuid-v4",
  "measured_at": "2025-07-01T07:00:00Z",
  "measured_date": "2025-07-01",
  "weight_kg": 82.5,
  "height_cm": 178.0,
  "waist_cm": 88.0,
  "hip_cm": 96.0,
  "chest_cm": 100.0,
  "arm_cm": 32.0,
  "notes": "",
  "created_at": "2025-07-01T07:05:00Z"
}
```

---

## 16. Lab Results

Lab results can be a **file upload** (PDF, image, etc.) or an **external link**.

**Max file size:** 10 MB  
**Content-Type:** `multipart/form-data`

### 16.1 Upload Lab Result

```
POST /clients/:id/lab-results
```

**Access:** Any Auth (client uploads own; nutritionist uploads for their client)  
**Path params:** `id` — client UUID

**Form fields:**
| Field | Type | Required | Notes |
|---|---|---|---|
| `title` | string | yes | e.g. `"آزمایش قند خون"` |
| `result_type` | string | yes | e.g. `"blood_sugar"`, `"lipid_panel"` |
| `test_date` | string | no | `"YYYY-MM-DD"` |
| `notes` | string | no | — |
| `file` | file | no* | Max 10 MB. Provide `file` OR `link`. |
| `link` | string | no* | External URL. Provide `file` OR `link`. |

**Response `201`:** [Lab result object](#lab-result-object)

---

### 16.2 List Lab Results

```
GET /clients/:id/lab-results?page=1&page_size=20
```

**Access:** Any Auth (client sees own; nutritionist sees their client's)

**Response `200`:** Paginated array of [lab result objects](#lab-result-object)

---

### 16.3 Download Lab Result

```
GET /lab-results/:id/download
```

**Access:** Any Auth (client or nutritionist of the client)

- If the result has a **file**: returns the file as an attachment (`Content-Disposition: attachment`).
- If the result has a **link**: redirects to the external URL (`302 Found`).

---

#### Lab Result Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "nutritionist_id": "uuid",
  "title": "آزمایش قند خون",
  "result_type": "blood_sugar",
  "test_date": "2025-06-15",
  "original_name": "lab_result.pdf",
  "file_type": "application/pdf",
  "file_size": 204800,
  "link": null,
  "notes": "",
  "created_at": "2025-06-15T10:00:00Z"
}
```

> `link` is `null` when a file was uploaded. `original_name`, `file_type`, `file_size` are `null`/`0` when only a link was provided.

---

## 17. Messages (Chat)

Client-nutritionist 1-on-1 messaging. Supports text and file attachments (images, PDFs, etc.).

**Max attachment size:** 10 MB  
**Content-Type:** `multipart/form-data`

### 17.1 Get Unread Count

```
GET /messages/unread-count
```

**Access:** Any Auth

**Response `200`:**
```json
{
  "data": { "unread_count": 3 }
}
```

---

### 17.2 Client — Get Conversation

```
GET /messages?page=1&page_size=20
```

**Access:** Client only  
Returns the conversation between the authenticated client and their nutritionist, newest first.

**Response `200`:** Paginated array of [message objects](#message-object)

---

### 17.3 Client — Send Message

```
POST /messages
```

**Access:** Client only  
**Content-Type:** `multipart/form-data`

| Form field | Type | Required |
|---|---|---|
| `content` | string | no* | Text content. Provide `content` and/or `file`. |
| `file` | file | no* | Max 10 MB. |

**Response `201`:** [Message object](#message-object)

> A push notification is sent to the nutritionist automatically.

---

### 17.4 Nutritionist — Get Conversation with Client

```
GET /clients/:id/messages?page=1&page_size=20
```

**Access:** Nutritionist only  
**Path params:** `id` — client UUID

**Response `200`:** Paginated array of [message objects](#message-object)

---

### 17.5 Nutritionist — Send Message to Client

```
POST /clients/:id/messages
```

**Access:** Nutritionist only  
**Content-Type:** `multipart/form-data`  
**Path params:** `id` — client UUID

| Form field | Type | Required |
|---|---|---|
| `content` | string | no* | — |
| `file` | file | no* | Max 10 MB. |

**Response `201`:** [Message object](#message-object)

> A push notification is sent to the client automatically.

---

#### Message Object

```json
{
  "id": "uuid",
  "sender_id": "uuid",
  "receiver_id": "uuid",
  "content": "سلام، برنامه‌تون رو ببینید",
  "is_mine": true,
  "read_at": null,
  "attachment": {
    "url": "/uploads/messages/file.pdf",
    "type": "application/pdf",
    "size": 102400,
    "name": "plan.pdf"
  },
  "created_at": "2025-07-01T10:00:00Z"
}
```

> `attachment` is `null` when no file was attached.  
> `is_mine` is `true` when the authenticated caller is the sender.  
> `read_at` is `null` until the receiver reads the message.

---

## 18. Food Requests

Clients can request a food to be added to the catalogue; their nutritionist approves or rejects it.

### 18.1 Submit Food Request (Client)

```
POST /food-requests
```

**Access:** Client

**Request body:**
```json
{
  "food_name": "کوفته تبریزی"
}
```

**Response `201`:** [Food request object](#food-request-object)

> A push notification is sent to the nutritionist automatically.

---

### 18.2 List Pending Food Requests (Nutritionist)

```
GET /food-requests?page=1&page_size=20
```

**Access:** Nutritionist  
Returns pending requests for this nutritionist's clients.

**Response `200`:** Paginated array of [food request objects](#food-request-object)

---

### 18.3 Approve Food Request (Nutritionist)

```
POST /food-requests/:id/approve
```

**Access:** Nutritionist

Approving a request creates a new food in the catalogue automatically.

**Request body:**
```json
{
  "name": "کوفته تبریزی",
  "unit": "عدد",
  "calories": 300.0,
  "protein": 18.0,
  "carbohydrate": 25.0,
  "fat": 12.0,
  "fiber": 2.0
}
```

| Field | Type | Required |
|---|---|---|
| `name` | string | yes |
| `unit` | string | yes |
| `calories` | float | no |
| `protein` | float | no |
| `carbohydrate` | float | no |
| `fat` | float | no |
| `fiber` | float | no |

**Response `200`:** [Food request object](#food-request-object) with `status: "approved"`

> A push notification is sent to the client automatically.

---

### 18.4 Reject Food Request (Nutritionist)

```
POST /food-requests/:id/reject
```

**Access:** Nutritionist

**Request body:**
```json
{
  "reason": "اطلاعات تغذیه‌ای کافی نیست"
}
```

**Response `200`:** [Food request object](#food-request-object) with `status: "rejected"`

> A push notification is sent to the client automatically.

---

#### Food Request Object

```json
{
  "id": "uuid",
  "client_id": "uuid",
  "nutritionist_id": "uuid",
  "food_name": "کوفته تبریزی",
  "status": "pending",
  "rejection_reason": null,
  "created_food_id": null,
  "created_at": "2025-07-01T09:00:00Z",
  "updated_at": "2025-07-01T09:00:00Z"
}
```

> `status` values: `"pending"` | `"approved"` | `"rejected"`  
> `created_food_id` is set to the new food's UUID when approved.  
> `rejection_reason` is `null` unless rejected.

---

## 19. Push Subscriptions

Web Push (VAPID) subscription management.

### 19.1 Subscribe

```
POST /push/subscribe
```

**Access:** Any Auth

**Request body:**
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/...",
  "p256dh": "base64url-encoded-key",
  "auth": "base64url-encoded-auth"
}
```

All three fields are required. These come directly from the browser's `PushSubscription` object.

**Response `201`:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "endpoint": "https://fcm.googleapis.com/fcm/send/...",
    "created_at": "2025-07-01T09:00:00Z"
  }
}
```

---

### 19.2 Unsubscribe

```
DELETE /push/subscribe
```

**Access:** Any Auth

**Request body:**
```json
{
  "endpoint": "https://fcm.googleapis.com/fcm/send/..."
}
```

**Response `204`:** No body

---

## 20. Notification Preferences

Control which push notifications a user receives.

### 20.1 Get Preferences

```
GET /notifications/preferences
```

**Access:** Any Auth

**Response `200`:** [Notification preferences object](#notification-preferences-object)

---

### 20.2 Update Preferences

```
PATCH /notifications/preferences
```

**Access:** Any Auth

**Request body:**
```json
{
  "meal_reminders": true,
  "water_reminders": true,
  "message_alerts": true,
  "diet_updates": true
}
```

All fields are boolean. This is an **upsert** — calling it for the first time creates the record.

**Response `200`:** [Notification preferences object](#notification-preferences-object)

---

#### Notification Preferences Object

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "meal_reminders": true,
  "water_reminders": true,
  "message_alerts": true,
  "diet_updates": true
}
```

---

## Appendix A — Push Notification Events

The server automatically sends Web Push notifications on these events:

| Event | Trigger endpoint | Recipient |
|---|---|---|
| New diet plan assigned | `POST /clients/:id/plans` | Client |
| New message received | `POST /messages` | Nutritionist |
| New message received | `POST /clients/:id/messages` | Client |
| Food request submitted | `POST /food-requests` | Nutritionist |
| Food request approved | `POST /food-requests/:id/approve` | Client |
| Food request rejected | `POST /food-requests/:id/reject` | Client |

---

## Appendix B — Role Access Summary

| Endpoint group | Public | Client | Nutritionist | Super Admin |
|---|---|---|---|---|
| Auth (login, OTP, refresh) | ✓ | — | — | — |
| Auth (logout, me) | — | ✓ | ✓ | ✓ |
| Avatar upload | — | ✓ (own) | ✓ (own) | ✓ (any) |
| Admin stats | — | — | — | ✓ |
| Admin nutritionist CRUD | — | — | — | ✓ |
| Client management | — | — | ✓ (own clients) | — |
| Foods CRUD | — | ✓ (own) | ✓ (own) | ✓ (any) |
| Food categories (read) | — | ✓ | ✓ | ✓ |
| Food categories (write) | — | — | — | ✓ |
| Medications CRUD | — | ✓ (own) | ✓ (own) | ✓ (any) |
| Diet plan management | — | ✓ (read own) | ✓ (full) | — |
| Tracking (write) | — | ✓ (own) | — | — |
| Tracking (read) | — | ✓ (own) | ✓ (their clients) | — |
| Lab results | — | ✓ (own) | ✓ (their clients) | — |
| Messages | — | ✓ | ✓ | — |
| Food requests (submit) | — | ✓ | — | — |
| Food requests (review) | — | — | ✓ (their clients') | — |
| Push subscriptions | — | ✓ | ✓ | ✓ |
| Notification preferences | — | ✓ | ✓ | ✓ |
