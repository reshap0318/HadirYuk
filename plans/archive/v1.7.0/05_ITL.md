---
title: 05_ITL.md
version: 1.7.0
created: 2026-05-29
last_modified: 2026-06-12
---

# Implementation Task List (ITL)

> **Project:** HadirYuk — Sistem Absensi Karyawan
> **Stack:** Go (Gin) + Vue 3 + MySQL + Redis
> **Generated:** 2026-05-29

---

## Implementation Phases

| Phase       | Name                             | Description                                                                    |
| ----------- | -------------------------------- | ------------------------------------------------------------------------------ |
| **Phase 1** | Foundation (DONE)                | Auth, RBAC, Data Master, Profil, Notifikasi, Upload, Peta                      |
| **Phase 2** | Core Attendance (P0)             | Shift assignment, profile photo, check-in/out geotagging + photo evidence & QR |
| **Phase 3** | Employee Self-Service (P1)       | Riwayat absensi, dasbor karyawan, cuti                                         |
| **Phase 4** | HR Operations (P2)               | Dasbor HR, koreksi absensi, jadwal shift, laporan                              |
| **Phase 5** | Analytics & Polish (P3-P4)       | Statistik keterlambatan, dasbor admin, audit log                               |
| **Phase 6** | Face Recognition (Optional / P4) | Face embedding generation, face matching during check-in/out                   |

---

## Phase 1: Foundation — COMPLETED ✅

### T-001: Login Endpoint JWT RS256 + JWKS

- **Feature/Module:** Authentication
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.1.1 Functional Requirements — Login
- **TDD Ref:**
  - POST /api/auth/login

### T-002: Logout Endpoint + Session Cleanup

- **Feature/Module:** Authentication
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [x]
- **FSD Ref:**
  - §2.1.2 Functional Requirements — Logout
- **TDD Ref:**
  - POST /api/auth/logout

### T-003: Refresh Token Endpoint

- **Feature/Module:** Authentication
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.1.6 Functional Requirements — Token Refresh
- **TDD Ref:**
  - POST /api/auth/refresh

### T-004: Forgot Password + Email Async

- **Feature/Module:** Authentication
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.1.4 Functional Requirements — Forgot Password
  - §3.2 User Interaction — Forgot Password Page
- **TDD Ref:**
  - POST /api/auth/forgot-password

### T-005: Reset Password + Token Validation

- **Feature/Module:** Authentication
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.1.5 Functional Requirements — Reset Password
  - §3.3 User Interaction — Reset Password Page
- **TDD Ref:**
  - POST /api/auth/reset-password

### T-006: CRUD Role Endpoints

- **Feature/Module:** UAM
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.7.1 Functional Requirements — Create Role
- **TDD Ref:**
  - POST /api/roles
  - GET /api/roles
  - GET /api/roles/:id
  - PUT /api/roles/:id
  - DELETE /api/roles/:id

### T-007: CRUD Permission Endpoints

- **Feature/Module:** UAM
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.7.2 Functional Requirements — Assign Permissions to Role
- **TDD Ref:**
  - GET /api/permissions
  - POST /api/permissions
  - PUT /api/permissions/:id
  - DELETE /api/permissions/:id

### T-008: Role-Permission Mapping

- **Feature/Module:** UAM
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.7.2 Functional Requirements — Assign Permissions to Role
  - §3.10 User Interaction — Role Management Form
- **TDD Ref:**
  - PUT /api/roles/:id (accepts `permissions` array inline)

### T-009: User-Role Assignment

- **Feature/Module:** UAM
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.7.3 Functional Requirements — Assign Role to User
- **TDD Ref:**
  - PUT /api/users/:id (accepts `roles` array inline)

### T-010: Permission Middleware + Cache

- **Feature/Module:** UAM
- **Priority:** P0
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.7 Functional Requirements — UAM (Role & Permissions)
- **TDD Ref:**
  - Middleware (JWT RS256, CORS, RBAC)

### T-011: CRUD Shift Endpoints

- **Feature/Module:** Data Master
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.3.1 Functional Requirements — Create Shift
  - §3.7 User Interaction — Shift Management Form
- **TDD Ref:**
  - POST /api/shifts
  - GET /api/shifts
  - GET /api/shifts/:id
  - PUT /api/shifts/:id
  - DELETE /api/shifts/:id

### T-012: CRUD Office Location Endpoints

- **Feature/Module:** Data Master
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.9.1 Functional Requirements — Add Office Location
  - §3.18 User Interaction — Location Management
- **TDD Ref:**
  - POST /api/locations
  - GET /api/locations
  - GET /api/locations/:id
  - PUT /api/locations/:id
  - DELETE /api/locations/:id

### T-013: CRUD Leave Type Endpoints

- **Feature/Module:** Data Master
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.4 Leave Management — Functional Requirements
- **TDD Ref:**
  - GET /api/leave-types
  - POST /api/leave-types
  - PUT /api/leave-types/:id
  - DELETE /api/leave-types/:id

### T-014: CRUD User with Auto-Profile

- **Feature/Module:** User Management
- **Priority:** P0
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.6.1 Functional Requirements — Create Employee
  - §3.9 User Interaction — User Management Form
- **TDD Ref:**
  - POST /api/users
  - GET /api/users
  - GET /api/users/:id
  - PUT /api/users/:id
  - DELETE /api/users/:id

### T-015: View/Edit Profile (/me)

- **Feature/Module:** Profile
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §3.20 User Interaction — Profile Page
- **TDD Ref:**
  - GET /api/me
  - PUT /api/me

### T-016: Upload Avatar

- **Feature/Module:** Profile
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.6.2 Functional Requirements — Upload Profile Photo
- **TDD Ref:**
  - POST /api/upload

### T-016b: Backend — Change Password Endpoint (Profile)

- **Feature/Module:** Profile
- **Priority:** P1
- **Estimated Effort:** 1h
- **Status:** [x]
- **FSD Ref:**
  - §2.14.1 Functional Requirements — Change Password (Profile)
- **TDD Ref:**
  - POST /api/me/change-password

### T-017: CRUD Notifications + Mark Read

- **Feature/Module:** Notifications
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - (Boilerplate — not in FSD)
- **TDD Ref:**
  - GET /api/notifications
  - GET /api/notifications/unread-count
  - GET /api/notifications/:id
  - PATCH /api/notifications/:id/read
  - PATCH /api/notifications/mark-all-read
  - DELETE /api/notifications/:id

### T-018: File Upload with UUID

- **Feature/Module:** Upload
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - (Boilerplate — not in FSD)
- **TDD Ref:**
  - POST /api/upload

### T-019: Leaflet Map Component (Frontend)

- **Feature/Module:** Map
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §3.18 User Interaction — Location Management
- **TDD Ref:**
  - Frontend component (no API endpoint)

### T-020: Health Check + JWKS Endpoint

- **Feature/Module:** Health
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [x]
- **FSD Ref:**
  - (Boilerplate — not in FSD)
- **TDD Ref:**
  - GET /health
  - GET /.well-known/jwks.json

---

## Phase 2: Core Attendance (P0)

### T-021: Backend — Assign Shift to Employee Endpoint

- **Feature/Module:** Shift Assignment
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.3.2 Functional Requirements — Assign Shift to Employee
- **TDD Ref:**
  - POST /api/shifts/assignments (request payload uses `shift`, `users` — not `shift_id`, `user_ids`)

### T-021b: Backend — Update Shift Assignment Endpoint

- **Feature/Module:** Shift Assignment
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.3.2 Functional Requirements — Assign Shift to Employee
- **TDD Ref:**
  - PUT /api/shifts/assignments/:id (request payload uses `shift`, `users`, `start_date`, `end_date`)

### T-021c: Backend — Delete Shift Assignment Endpoint

- **Feature/Module:** Shift Assignment
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [x]
- **FSD Ref:**
  - §2.3.2 Functional Requirements — Assign Shift to Employee
- **TDD Ref:**
  - DELETE /api/shifts/assignments/:id

### T-022: Frontend — Form Assign Shift ke Multiple Employee

- **Feature/Module:** Shift Assignment
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.3.2 Functional Requirements — Assign Shift to Employee
- **TDD Ref:**
  - POST /api/shifts/assignments (sends `shift`, `users` — not `shift_id`, `user_ids`)

### T-022b: Frontend — Shift Assignment List (View, Edit, Delete)

- **Feature/Module:** Shift Assignment
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.3.2 Functional Requirements — Assign Shift to Employee
- **TDD Ref:**
  - GET /api/shifts/assignments
  - PUT /api/shifts/assignments/:id
  - DELETE /api/shifts/assignments/:id

### T-023: Backend — Get Employee Schedule Endpoint

- **Feature/Module:** Shift Assignment
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.3.3 Functional Requirements — Shift Schedule View
- **TDD Ref:**
  - GET /api/shifts/schedule

### T-024: Backend — Upload Profile Photo Endpoint

- **Feature/Module:** Profile Photo
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.6.2 Functional Requirements — Upload Profile Photo
- **TDD Ref:**
  - POST /api/users/:id/face-photo

### T-025: Frontend — Profile Photo Capture/Upload Component

- **Feature/Module:** Profile Photo
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.6.2 Functional Requirements — Upload Profile Photo
  - §3.20 User Interaction — Profile Page
- **TDD Ref:**
  - POST /api/users/:id/face-photo

### T-027: Backend — Check-in Geotagging + Haversine Validation + Photo Evidence

- **Feature/Module:** Attendance Check-in
- **Priority:** P0
- **Estimated Effort:** 4h
- **Status:** [x]
- **Note:** Implemented with auto-detect shift, window logic (15min buffer), multi-session support, photo validation (UUID + extension + size ≤ 5MB), file move from tmp to evidence
- **FSD Ref:**
  - §2.2.1 Functional Requirements — Check-in (Geotagging + Photo Evidence)
  - §3.5 User Interaction — Attendance Check-in (Geotagging)
- **TDD Ref:**
  - POST /api/attendance/checkin (receives `lat`, `lng`, `image` where `image` = UUID from upload; requires `attendance.checkin` permission)
  - POST /api/upload (used by frontend before check-in)

### T-028: Backend — Photo UUID Validation for Attendance Evidence

- **Feature/Module:** Attendance Check-in
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **Note:** Validates UUID existence in `storage/tmp`, file extension (.jpg/.jpeg/.png/.webp), file size ≤ 5MB, moves file to `storage/attendance-evidence`
- **FSD Ref:**
  - §2.2.1 Functional Requirements — Check-in (Geotagging + Photo Evidence)
  - §2.2.3 Functional Requirements — Check-out (Geotagging + Photo Evidence)
- **TDD Ref:**
  - POST /api/attendance/checkin (validates `image`)
  - POST /api/attendance/checkout (validates `image`)
  - POST /api/upload (returns `uuid`)

### T-029: Frontend — Check-in Page (Location + Camera + Photo Upload → image UUID)

- **Feature/Module:** Attendance Check-in
- **Priority:** P0
- **Estimated Effort:** 5h
- **Status:** [x]
- **Note:** Implemented as multi-session page with single button (Check In → Check Out → Selesai), auto-detect shift, MapCard, ActionSection, StatusCard, ClockCard, LocationCard, ShiftCarousel, CameraModal components
- **FSD Ref:**
  - §2.2.1 Functional Requirements — Check-in (Geotagging + Photo Evidence)
  - §3.5 User Interaction — Attendance Check-in (Geotagging)
  - §5.1 Feature Logic Flow — Check-in (Geotagging + Photo Evidence) Flow
- **TDD Ref:**
  - POST /api/attendance/checkin (sends `lat`, `lng`, `image`)
  - POST /api/upload (upload photo first, get `uuid`)

### T-030: Backend — Check-out Geotagging + Duration Calculation + Photo Evidence

- **Feature/Module:** Attendance Check-out
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **Note:** Implemented with window validation (block before shiftEnd-15min), overtime calculation (after shiftEnd+15min), auto-detect active session, photo validation
- **FSD Ref:**
  - §2.2.3 Functional Requirements — Check-out (Geotagging + Photo Evidence)
- **TDD Ref:**
  - POST /api/attendance/checkout (receives `lat`, `lng`, `image` where `image` = UUID from upload; requires `attendance.checkout` permission)
  - POST /api/upload (used by frontend before check-out)

### T-031: Frontend — Combined Check-in/Check-out Page (Location + Camera + Photo Upload → image UUID)

- **Feature/Module:** Attendance Check-out
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **Note:** Digabung dengan T-029 dalam satu halaman — single button berubah state: Check In (green) → Check Out (orange) → Selesai (purple). Includes: MapCard, ActionSection, StatusCard, ClockCard, LocationCard, ShiftCarousel, CameraModal. Store supports multi-session with `sessions[]`, `currentAction`, `todaysShifts[]`, `executeAction()`
- **FSD Ref:**
  - §2.2.3 Functional Requirements — Check-out (Geotagging + Photo Evidence)
- **TDD Ref:**
  - POST /api/attendance/checkout (sends `lat`, `lng`, `image`)
  - POST /api/upload (upload photo first, get `uuid`)
  - GET /api/attendance/today (enriched response with sessions, current_action, todays_shifts)
  - POST /api/attendance/nearest-office (find nearest office by coordinates)

### T-033: Backend — Generate QR Code with Signature + Expiry

- **Feature/Module:** QR Code
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **Note:** Implemented with UUID code_value, HMAC-SHA256 signature, expiry datetime, notification on generate
- **FSD Ref:**
  - §2.11.1 Functional Requirements — Generate QR Code
  - §3.12 User Interaction — QR Code Management
- **TDD Ref:**
  - POST /api/qr-codes/generate (request payload uses `office`, not `office_id`)

### T-034: Backend — View Active QR Codes + Revoke

- **Feature/Module:** QR Code
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **Note:** Implemented with filter active+non-expired, revoke sets IsActive=false+RevokedAt, QRCodeValidate service exists but not yet exposed as route
- **FSD Ref:**
  - §2.11.2 Functional Requirements — View Active QR Codes
  - §2.11.3 Functional Requirements — Revoke QR Code
- **TDD Ref:**
  - GET /api/qr-codes
  - POST /api/qr-codes/:id/revoke

### T-035: Frontend — QR Code Management Page (Generate, List, Revoke)

- **Feature/Module:** QR Code
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **Note:** Implemented with card grid layout, GenerateModal (office picker + date/time), QRCodeDisplay component (npm qrcode), preview modal, revoke with confirmation, permission-gated buttons
- **FSD Ref:**
  - §2.11 QR Code Management — Functional Requirements
  - §3.12 User Interaction — QR Code Management
- **TDD Ref:**
  - POST /api/qr-codes/generate
  - GET /api/qr-codes
  - POST /api/qr-codes/:id/revoke

### T-036: Backend — QR Code Validation Service (Signature + Expiry)

- **Feature/Module:** QR Code
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **Note:** `QRCodeValidate()` implemented in `qr_code_service.go` — validates is_active, not revoked, not expired, HMAC signature constant-time comparison
- **FSD Ref:**
  - §2.2.2 Functional Requirements — Check-in (QR Code)
- **TDD Ref:**
  - Service layer (QR Code Service)

### T-037: Backend — Check-in via QR Code

- **Feature/Module:** QR Check-in
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **Note:** `AttendanceQRCheckIn()` — validates QR code, checks active session, auto-detects shift, window validation, creates attendance without photo/GPS, office from QR
- **FSD Ref:**
  - §2.2.2 Functional Requirements — Check-in (QR Code)
  - §3.6 User Interaction — Attendance Check-in (QR Code)
  - §5.2 Feature Logic Flow — Check-in (QR Code) Flow
- **TDD Ref:**
  - POST /api/attendance/checkin/qr

### T-038: Frontend — QR Scanner Component (html5-qrcode)

- **Feature/Module:** QR Check-in
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **Note:** `QRScanner.vue` — uses `html5-qrcode` library, start/stop/reset controls, scan result display, emits `scan` event with code_value, rear camera default
- **FSD Ref:**
  - §2.2.2 Functional Requirements — Check-in (QR Code)
  - §3.6 User Interaction — Attendance Check-in (QR Code)
- **TDD Ref:**
  - POST /api/attendance/checkin/qr

### T-039: Backend — Check-out via QR Code

- **Feature/Module:** QR Check-out
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [x]
- **Note:** `AttendanceQRCheckOut()` — validates QR code, checks active session, window validation (block before shiftEnd-15min), overtime calculation, updates attendance without photo/GPS
- **FSD Ref:**
  - §2.2.4 Functional Requirements — Check-out (QR Code)
- **TDD Ref:**
  - POST /api/attendance/checkout/qr

### T-040: Frontend — QR Scanner for Check-out (Reuse Component)

- **Feature/Module:** QR Check-out
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [x]
- **Note:** Reuses `QRScanner.vue` component — tab toggle (Geotagging / QR Code) in `IndexView.vue`, `qrCheckIn()` and `qrCheckOut()` methods in attendance store, auto-detects action from `currentAction`

### T-041: Backend — Attendance GORM Model + Migration

- **Feature/Module:** Attendance Model
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **Note:** Model includes: `ShiftID`, `OfficeID`, `DistanceMeters`, `OvertimeMinutes`, `ImageIn`/`ImageOut` (file paths). Unique constraint: `(user_id, date, shift_id)` for multi-session support. Fields `check_in_method`, `check_out_method`, `qr_code_id`, correction fields NOT included (future implementation)
- **FSD Ref:**
  - §2.2 Attendance (Absensi) — Functional Requirements
- **TDD Ref:**
  - ERD — ATTENDANCES table

### T-042: Backend — QRCode GORM Model + Migration

- **Feature/Module:** Attendance Model
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [x]
- **Note:** Model includes: `ID`, `OfficeID`, `CodeValue` (UUID, unique), `Signature` (HMAC-SHA256), `ExpiresAt`, `IsActive`, `CreatedBy`, `RevokedAt`. Table: `qr_codes`
- **FSD Ref:**
  - §2.11 QR Code Management — Functional Requirements
- **TDD Ref:**
  - ERD — QR_CODES table

### T-043: Backend — UserShiftAssignment GORM Model + Migration

- **Feature/Module:** Attendance Model
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [x]
- **Note:** Model renamed from `EmployeeShift` to `UserShiftAssignment`. Includes: `ID` (auto), `UserID`, `ShiftID`, `StartDate`, `EndDate` (nullable = ongoing), `IsActive`. Table: `user_shift_assignments`
- **FSD Ref:**
  - §2.3.2 Functional Requirements — Assign Shift to Employee
- **TDD Ref:**
  - ERD — USER_SHIFT_ASSIGNMENTS table

### T-044: Backend — LeaveBalance GORM Model + Migration

- **Feature/Module:** Attendance Model
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.2 Functional Requirements — View Leave Balance
- **TDD Ref:**
  - ERD — LEAVE_BALANCES table

### T-045: Backend — LeaveRequest GORM Model + Migration

- **Feature/Module:** Attendance Model
- **Priority:** P0
- **Estimated Effort:** 1h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.1 Functional Requirements — Submit Leave Request
- **TDD Ref:**
  - ERD — LEAVE_REQUESTS table

---

### T-046: Backend — Multi-Session Attendance (Auto-Detect Shift + Window Logic)

- **Feature/Module:** Attendance Check-in/Check-out
- **Priority:** P0
- **Estimated Effort:** 6h
- **Status:** [x]
- **Note:** Implements `findApplicableShift()` for auto-detecting shift based on current time window. Check-in window: `[shiftStart - 15min, shiftEnd]`. Check-out window: `[shiftEnd - 15min, ∞)`. Status logic: `present` if within `[shiftStart - 15min, shiftStart + 15min]`, else `late`. Overtime calculation: minutes after `shiftEnd + 15min`. Cross-day support: active session check across ALL dates.
- **FSD Ref:**
  - §2.2.1 Functional Requirements — Check-in (Geotagging + Photo Evidence)
  - §2.2.3 Functional Requirements — Check-out (Geotagging + Photo Evidence)
- **TDD Ref:**
  - POST /api/attendance/checkin
  - POST /api/attendance/checkout
  - Service: `findApplicableShift(userID, now)`

### T-047: Backend — Nearest Office Endpoint

- **Feature/Module:** Attendance Geolocation
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **Note:** `POST /api/attendance/nearest-office` — finds nearest active office by coordinates using Haversine formula, returns distance in meters
- **FSD Ref:**
  - §2.2.1 Functional Requirements — Check-in (Geotagging + Photo Evidence)
- **TDD Ref:**
  - POST /api/attendance/nearest-office

### T-048: Backend — Enriched Today Status Endpoint

- **Feature/Module:** Attendance Status
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **Note:** `GET /api/attendance/today` returns enriched response: `{ sessions[], current_action{action, shift, cross_day_session}, todays_shifts[] }`. Supports multi-session display and cross-day scenario detection.
- **FSD Ref:**
  - §2.2.5 Functional Requirements — Attendance History
- **TDD Ref:**
  - GET /api/attendance/today

### T-049: Frontend — Multi-Session Attendance Store

- **Feature/Module:** Attendance Store
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **Note:** Store supports `sessions[]`, `currentAction`, `todaysShifts[]`, `executeAction()` method. Legacy fields (`checkInData`, `checkOutData`, `todayStatus`) kept for backward compatibility. Haversine distance calculation on frontend for proximity check.
- **FSD Ref:**
  - §2.2.1 Functional Requirements — Check-in (Geotagging + Photo Evidence)
  - §2.2.3 Functional Requirements — Check-out (Geotagging + Photo Evidence)
- **TDD Ref:**
  - GET /api/attendance/today
  - POST /api/attendance/checkin
  - POST /api/attendance/checkout
  - POST /api/attendance/nearest-office

### T-049b: Frontend — Attendance Page Components

- **Feature/Module:** Attendance UI
- **Priority:** P0
- **Estimated Effort:** 4h
- **Status:** [x]
- **Note:** Components: `IndexView.vue` (main page), `MapCard.vue` (Leaflet map), `ActionSection.vue` (single button), `StatusCard.vue` (session status), `ClockCard.vue` (real-time clock), `LocationCard.vue` (distance info), `ShiftCarousel.vue` (today's shifts), `CameraModal.vue` (photo capture). Single button changes state: Check In (green) → Check Out (orange) → Selesai (purple).
- **FSD Ref:**
  - §3.5 User Interaction — Attendance Check-in (Geotagging)
  - §3.4 User Interaction — Karyawan Dashboard
- **TDD Ref:**
  - Frontend components (no API endpoint)

---

## Phase 3: Employee Self-Service (P1)

### T-050: Backend — Attendance History Endpoint (Self)

- **Feature/Module:** Attendance History
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.5 Functional Requirements — Attendance History
  - §3.15 User Interaction — Attendance History
- **TDD Ref:**
  - GET /api/attendance/history

### T-051: Backend — Today Status Endpoint

- **Feature/Module:** Attendance History
- **Priority:** P1
- **Estimated Effort:** 1h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.5 Functional Requirements — Attendance History
- **TDD Ref:**
  - GET /api/attendance/today

### T-052: Backend — Monthly Stats Endpoint

- **Feature/Module:** Attendance History
- **Priority:** P1
- **Estimated Effort:** 1h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.5 Functional Requirements — Attendance History
- **TDD Ref:**
  - GET /api/attendance/stats

### T-053: Frontend — Attendance History Page (Filter, Pagination)

- **Feature/Module:** Attendance History
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.5 Functional Requirements — Attendance History
  - §3.15 User Interaction — Attendance History
- **TDD Ref:**
  - GET /api/attendance/history

### T-054: Frontend — Attendance Detail Modal

- **Feature/Module:** Attendance History
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.5 Functional Requirements — Attendance History
- **TDD Ref:**
  - GET /api/attendance/history

### T-055: Backend — Employee Dashboard Endpoint

- **Feature/Module:** Employee Dashboard
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.8.1 Functional Requirements — Karyawan Dashboard
- **TDD Ref:**
  - GET /api/dashboard/employee

### T-056: Frontend — Employee Dashboard Page (Status, Summary, Schedule)

- **Feature/Module:** Employee Dashboard
- **Priority:** P1
- **Estimated Effort:** 4h
- **Status:** [ ]
- **FSD Ref:**
  - §2.8.1 Functional Requirements — Karyawan Dashboard
  - §3.4 User Interaction — Karyawan Dashboard
- **TDD Ref:**
  - GET /api/dashboard/employee

### T-057: Backend — Submit Leave Request + Balance Deduction

- **Feature/Module:** Leave Request
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.1 Functional Requirements — Submit Leave Request
  - §3.8 User Interaction — Leave Request Form
  - §5.3 Feature Logic Flow — Leave Request Flow
- **TDD Ref:**
  - POST /api/leave (request payload uses `leave_type`, not `leave_type_id`)

### T-058: Frontend — Leave Request Form + Validation

- **Feature/Module:** Leave Request
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.1 Functional Requirements — Submit Leave Request
  - §3.8 User Interaction — Leave Request Form
- **TDD Ref:**
  - POST /api/leave

### T-059: Backend — Leave Balance Endpoint

- **Feature/Module:** Leave Balance
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.2 Functional Requirements — View Leave Balance
  - §3.16 User Interaction — Leave Balance
- **TDD Ref:**
  - GET /api/leave/balance

### T-060: Frontend — Leave Balance Page (Progress Bar, History)

- **Feature/Module:** Leave Balance
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.2 Functional Requirements — View Leave Balance
  - §3.16 User Interaction — Leave Balance
- **TDD Ref:**
  - GET /api/leave/balance

### T-061: Backend — Auto-Init Leave Balance on Employee Creation

- **Feature/Module:** Leave Balance
- **Priority:** P1
- **Estimated Effort:** 1h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.2 Functional Requirements — View Leave Balance
- **TDD Ref:**
  - Leave Balance Initialization Strategy (TDD §3)

### T-062: Backend — Yearly Reset Cron Job (1 Januari)

- **Feature/Module:** Leave Balance
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4.2 Functional Requirements — View Leave Balance
- **TDD Ref:**
  - Leave Balance Initialization Strategy (TDD §3)

---

## Phase 4: HR Operations (P2)

### T-070: Backend — HR Dashboard Endpoint (Stats, Chart, Not-Attended, Leaves)

- **Feature/Module:** HR Dashboard
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [ ]
- **FSD Ref:**
  - §2.8.2 Functional Requirements — HR Dashboard
  - §3.13 User Interaction — HR Dashboard
- **TDD Ref:**
  - GET /api/dashboard/hr

### T-071: Frontend — HR Dashboard Page (Cards, Chart, Tables)

- **Feature/Module:** HR Dashboard
- **Priority:** P2
- **Estimated Effort:** 4h
- **Status:** [ ]
- **FSD Ref:**
  - §2.8.2 Functional Requirements — HR Dashboard
  - §3.13 User Interaction — HR Dashboard
- **TDD Ref:**
  - GET /api/dashboard/hr

### T-072: Backend — Correct Attendance Endpoint

- **Feature/Module:** Attendance Correction
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.6 Functional Requirements — Attendance Correction
  - §3.11 User Interaction — Attendance Correction Form
- **TDD Ref:**
  - PUT /api/attendance/:id/correct

### T-073: Frontend — Correction Form (DateTime Picker, Reason)

- **Feature/Module:** Attendance Correction
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.6 Functional Requirements — Attendance Correction
  - §3.11 User Interaction — Attendance Correction Form
- **TDD Ref:**
  - PUT /api/attendance/:id/correct

### T-074: Backend — Employee Schedule Endpoint (All Employees)

- **Feature/Module:** Shift Schedule
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.3.3 Functional Requirements — Shift Schedule View
  - §3.17 User Interaction — Shift Schedule View
- **TDD Ref:**
  - GET /api/shifts/schedule

### T-075: Frontend — Shift Schedule Calendar View (Month, Color-Coded)

- **Feature/Module:** Shift Schedule
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [ ]
- **FSD Ref:**
  - §2.3.3 Functional Requirements — Shift Schedule View
  - §3.17 User Interaction — Shift Schedule View
- **TDD Ref:**
  - GET /api/shifts/schedule

### T-076: Backend — Attendance Report Endpoint with Filters

- **Feature/Module:** Attendance Report
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.10.1 Functional Requirements — Attendance Report
  - §3.15 User Interaction — Attendance History
- **TDD Ref:**
  - GET /api/reports/attendance

### T-077: Backend — Export to Excel (.xlsx)

- **Feature/Module:** Attendance Report
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [ ]
- **FSD Ref:**
  - §2.10.2 Functional Requirements — Export to Excel
- **TDD Ref:**
  - GET /api/reports/attendance/export/excel

### T-078: Backend — Export to PDF

- **Feature/Module:** Attendance Report
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [ ]
- **FSD Ref:**
  - §2.10.3 Functional Requirements — Export to PDF
- **TDD Ref:**
  - GET /api/reports/attendance/export/pdf

### T-079: Frontend — Report Page (Filters, Table, Export Buttons)

- **Feature/Module:** Attendance Report
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.10 Reporting — Functional Requirements
  - §3.15 User Interaction — Attendance History
- **TDD Ref:**
  - GET /api/reports/attendance
  - GET /api/reports/attendance/export/excel
  - GET /api/reports/attendance/export/pdf

### T-080: Backend — Leave Report Endpoint with Filters

- **Feature/Module:** Leave Report
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.10.4 Functional Requirements — Leave Report
- **TDD Ref:**
  - GET /api/reports/leave

### T-081: Backend — Leave Report Export Excel/PDF

- **Feature/Module:** Leave Report
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.10.4 Functional Requirements — Leave Report
- **TDD Ref:**
  - GET /api/reports/leave/export/excel
  - GET /api/reports/leave/export/pdf

### T-082: Frontend — Leave Report Page

- **Feature/Module:** Leave Report
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.10.4 Functional Requirements — Leave Report
- **TDD Ref:**
  - GET /api/reports/leave

### T-083: Backend — Attendance History with user_id=all Filter

- **Feature/Module:** Attendance View-All
- **Priority:** P2
- **Estimated Effort:** 1h
- **Status:** [x]
- **FSD Ref:**
  - §2.2.5 Functional Requirements — Attendance History
- **TDD Ref:**
  - GET /api/attendance/history

### T-084: Backend — Leave History with user_id=all Filter

- **Feature/Module:** Leave View-All
- **Priority:** P2
- **Estimated Effort:** 1h
- **Status:** [ ]
- **FSD Ref:**
  - §2.4 Leave Management — Functional Requirements
- **TDD Ref:**
  - GET /api/leave

---

## Phase 5: Analytics & Polish (P3-P4)

### T-090: Backend — Late Statistics Endpoint

- **Feature/Module:** Late Statistics
- **Priority:** P3
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.5.1 Functional Requirements — View Late Statistics
- **TDD Ref:**
  - GET /api/attendance/late-statistics

### T-091: Frontend — Late Statistics Page (Chart, Table, Trend)

- **Feature/Module:** Late Statistics
- **Priority:** P3
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.5.1 Functional Requirements — View Late Statistics
- **TDD Ref:**
  - GET /api/attendance/late-statistics

### T-092: Backend — Admin Dashboard Endpoint (System Stats, Activity, Health)

- **Feature/Module:** Admin Dashboard
- **Priority:** P3
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §3.14 User Interaction — Admin Dashboard
- **TDD Ref:**
  - GET /api/dashboard/admin

### T-093: Frontend — Admin Dashboard Page

- **Feature/Module:** Admin Dashboard
- **Priority:** P3
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §3.14 User Interaction — Admin Dashboard
- **TDD Ref:**
  - GET /api/dashboard/admin

### T-094: Backend — Audit Log Middleware (Record CRUD Changes)

- **Feature/Module:** Audit Log
- **Priority:** P4
- **Estimated Effort:** 4h
- **Status:** [ ]
- **FSD Ref:**
  - §2.12.1 Functional Requirements — View Audit Log
  - §3.19 User Interaction — Audit Log View
- **TDD Ref:**
  - GET /api/audit-logs

### T-095: Backend — Audit Log Endpoint with Filters

- **Feature/Module:** Audit Log
- **Priority:** P4
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.12.1 Functional Requirements — View Audit Log
  - §3.19 User Interaction — Audit Log View
- **TDD Ref:**
  - GET /api/audit-logs

### T-096: Frontend — Audit Log Page (Filters, Detail Modal)

- **Feature/Module:** Audit Log
- **Priority:** P4
- **Estimated Effort:** 3h
- **Status:** [ ]
- **FSD Ref:**
  - §2.12.1 Functional Requirements — View Audit Log
  - §3.19 User Interaction — Audit Log View
- **TDD Ref:**
  - GET /api/audit-logs

### T-097: Backend — AuditLog GORM Model + Migration

- **Feature/Module:** Audit Log
- **Priority:** P4
- **Estimated Effort:** 1h
- **Status:** [ ]
- **FSD Ref:**
  - §2.12.1 Functional Requirements — View Audit Log
- **TDD Ref:**
  - ERD — AUDIT_LOGS table (new)

---

## Phase 6: Face Recognition (Optional / P4)

> **Note:** All tasks in this phase are OPTIONAL. Core check-in/out already works with photo evidence. Face recognition can be added later as an enhancement.

### T-100: Backend — Face Recognition Service (GoCV/face-recognition-go)

- **Feature/Module:** Face Recognition
- **Priority:** P4
- **Estimated Effort:** 6h
- **Status:** [ ]
- **FSD Ref:**
  - §2.13.2 Functional Requirements — Face Recognition during Check-in (Optional Enhancement)
- **TDD Ref:**
  - Service layer (Face Recognition Service — Optional)

### T-101: Backend — Face Embedding Generation on Photo Upload

- **Feature/Module:** Face Recognition
- **Priority:** P4
- **Estimated Effort:** 4h
- **Status:** [ ]
- **FSD Ref:**
  - §2.13.1 Functional Requirements — Upload Face Photo for Recognition
  - §2.6.2 Functional Requirements — Upload Profile Photo
- **TDD Ref:**
  - POST /api/users/:id/face-photo

### T-102: Backend — Face Matching during Check-in

- **Feature/Module:** Face Recognition
- **Priority:** P4
- **Estimated Effort:** 4h
- **Status:** [ ]
- **FSD Ref:**
  - §2.13.2 Functional Requirements — Face Recognition during Check-in (Optional Enhancement)
- **TDD Ref:**
  - POST /api/attendance/checkin

### T-103: Backend — Face Matching during Check-out

- **Feature/Module:** Face Recognition
- **Priority:** P4
- **Estimated Effort:** 2h
- **Status:** [ ]
- **FSD Ref:**
  - §2.13.3 Functional Requirements — Face Recognition during Check-out (Optional Enhancement)
- **TDD Ref:**
  - POST /api/attendance/checkout

### T-104: Frontend — Face Photo Capture Component with Face Detection

- **Feature/Module:** Face Recognition
- **Priority:** P4
- **Estimated Effort:** 3h
- **Status:** [ ]
- **FSD Ref:**
  - §2.13.1 Functional Requirements — Upload Face Photo for Recognition
  - §3.20 User Interaction — Profile Page
- **TDD Ref:**
  - POST /api/users/:id/face-photo

### T-105: Backend — Face Recognition Feature Toggle Endpoint

- **Feature/Module:** Face Recognition
- **Priority:** P4
- **Estimated Effort:** 1h
- **Status:** [ ]
- **FSD Ref:**
  - §2.13 Face Recognition (Optional / Could Have)
- **TDD Ref:**
  - POST /api/face-recognition/enable
  - GET /api/face-recognition/status

---

## Testing Checklist

| Task ID | Test Scenario                 | Expected Result                                | Status |
| ------- | ----------------------------- | ---------------------------------------------- | ------ |
| T-021   | Assign shift to employee      | Employee ter-assign, record di user_shift_assignments | [ ]    |
| T-027   | Check-in dalam radius + foto  | Absensi tersimpan dengan status present + foto | [ ]    |
| T-027   | Check-in luar radius          | Error "Anda berada di luar area kantor"        | [ ]    |
| T-027   | Check-in di luar window       | Error "Waktu check-in di luar jendela shift"   | [ ]    |
| T-027   | Check-in saat ada session aktif | Error "Anda memiliki sesi check-in yang belum di-checkout" | [ ]    |
| T-028   | Upload foto evidence (UUID)   | Foto tersimpan sebagai bukti kehadiran         | [ ]    |
| T-028   | Foto format tidak valid       | Error "Format foto tidak didukung"             | [ ]    |
| T-028   | Foto terlalu besar (>5MB)     | Error "Ukuran foto terlalu besar"              | [ ]    |
| T-030   | Check-out setelah check-in    | Record updated dengan duration                 | [ ]    |
| T-030   | Check-out sebelum window      | Error "Belum waktunya check-out"               | [ ]    |
| T-030   | Check-out lembur (> shiftEnd + 15min) | overtime_minutes terhitung              | [ ]    |
| T-030   | Check-out tanpa check-in      | Error "Belum melakukan check-in"               | [ ]    |
| T-046   | 2 shift berurutan (06-12 + 18-24) | Check-in/out keduanya berhasil          | [ ]    |
| T-046   | Cross-day (malam 11/06 → pagi 12/06) | Session date tetap 11/06, overtime terhitung | [ ]    |
| T-046   | Auto-detect shift             | Sistem pilih shift yang applicable berdasarkan waktu | [ ]    |
| T-047   | Nearest office endpoint       | Return office terdekat + distance              | [ ]    |
| T-048   | Today status enriched         | Return sessions[], current_action, todays_shifts[] | [ ]    |
| T-048   | Cross-day session detection   | current_action.cross_day_session terisi        | [ ]    |
| T-033   | Generate QR code              | QR code tersimpan dengan signature + expiry    | [x]    |
| T-034   | View active QR codes          | List QR aktif (non-expired, non-revoked)       | [x]    |
| T-034   | Revoke QR code                | QR code tidak aktif setelah revoke             | [x]    |
| T-035   | QR management page            | Generate, view, revoke berjalan di frontend    | [x]    |
| T-036   | QR validation service         | Validasi signature, expiry, revoke             | [x]    |
| T-037   | Check-in QR valid             | Absensi tersimpan tanpa foto/GPS               | [x]    |
| T-037   | Check-in QR expired           | Error "QR Code sudah expired"                  | [x]    |
| T-038   | QR scanner component          | Scan QR dari kamera, emit code_value           | [x]    |
| T-039   | Check-out QR valid            | Record updated dengan duration + overtime      | [x]    |
| T-040   | QR scanner reuse check-out    | Tab toggle Geotagging/QR di IndexView          | [x]    |
| T-037   | Check-in QR expired           | Error "QR Code sudah expired"                  | [ ]    |
| T-057   | Submit cuti valid             | Leave request tersimpan, balance berkurang     | [ ]    |
| T-057   | Submit cuti saldo tidak cukup | Error "Sisa cuti tidak mencukupi"              | [ ]    |
| T-057   | Submit cuti overlap           | Error "Tanggal cuti overlap"                   | [ ]    |
| T-072   | Koreksi absensi               | Record updated, audit log tercatat             | [ ]    |
| T-077   | Export Excel                  | File .xlsx terdownload dengan data benar       | [ ]    |
| T-078   | Export PDF                    | File .pdf terdownload dengan data benar        | [ ]    |
| T-090   | Late statistics               | Data keterlambatan akurat per employee         | [ ]    |
| T-094   | Audit log CRUD                | Setiap perubahan tercatat di audit log         | [ ]    |
| T-100   | Face recognition enabled      | Face matching berjalan saat check-in           | [ ]    |
| T-100   | Face recognition disabled     | Check-in berjalan tanpa face matching          | [ ]    |

---

## Summary

| Phase                             | Total Tasks | Completed | Remaining | Est. Hours |
| --------------------------------- | ----------- | --------- | --------- | ---------- |
| Phase 1 (Foundation)              | 20          | 20        | 0         | ~48h       |
| Phase 2 (Core Attendance)         | 31          | 28        | 3         | ~83h       |
| Phase 3 (Self-Service)            | 13          | 4         | 9         | ~28h       |
| Phase 4 (HR Operations)           | 15          | 5         | 10        | ~38h       |
| Phase 5 (Analytics)               | 8           | 2         | 6         | ~20h       |
| Phase 6 (Face Recognition — Opt.) | 6           | 0         | 6         | ~20h       |
| **Total**                         | **93**      | **59**    | **34**    | **~237h**  |
