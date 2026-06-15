# CHANGELOG

## v1.8.0 - 2026-06-15

### Remove Leave Management Feature

- **01_PRD.md:** Removed user stories US-005 (Ajukan cuti), US-010 (Kelola cuti), US-016 (Sisa cuti); removed Leave Management from Must Have and Sisa Cuti View from Should Have; removed "Ajukan Cuti" and "Kelola Cuti" from user flow
- **02_FSD.md:** Removed §2.4 Leave Management, §2.10.4 Leave Report, §3.8 Leave Request Form, §3.16 Leave Balance, §5.3 Leave Request Flow, §6.3 Leave Validation; removed leave references from Hierarchy, Dashboard, Report, Attendance History, and Audit Log; removed leave use cases (UC6, UC7, UC16) from use case diagram; renumbered subsequent sections
- **03_Role_Matrix.md:** Removed `leave.view-all` and `leave.manage-types` permissions; removed Leave section from matrix; removed Leave Records data ownership section; updated role descriptions
- **04_TDD.md:** Removed LEAVE_TYPES, LEAVE_REQUESTS, LEAVE_BALANCES from ERD; removed §4.4 Leave Endpoints; removed leave report endpoints; removed leave indices; removed Leave Balance Initialization Strategy; renumbered subsequent sections
- **05_ITL.md:** Removed tasks T-013, T-044, T-045, T-057-T-062, T-080-T-082, T-084; removed leave test cases; updated Phase 3 description; updated summary counts

## v1.7.0 - 2026-06-12

### Multi-Session Attendance Implementation Aligned with Plan

- **04_TDD.md:** Updated ATTENDANCES ERD table — replaced `attendance_date` → `date`, `check_in_time`/`check_out_time` → `time_in`/`time_out`, `check_in_lat`/`check_in_lng` → `lat_in`/`lng_in`, `check_in_photo_url`/`check_out_photo_url` → `image_in`/`image_out` (file paths), removed `check_in_method`/`check_out_method`, `qr_code_id`, `location_id`, correction fields; added `shift_id`, `office_id`, `distance_meters`, `overtime_minutes`, `duration`; updated unique constraint from `(user_id, attendance_date)` → `(user_id, date, shift_id)`; renamed `EMPLOYEE_SHIFTS` → `USER_SHIFT_ASSIGNMENTS` with `id` PK and `is_active` field; updated `QR_CODES` with `created_by`, `revoked_at`, `updated_at` fields; updated indexing strategy to match new constraint
- **04_TDD.md:** Updated §4.2 Attendance Endpoints — removed QR check-in/checkout, history, stats, correction, late-statistics (marked as NOT YET IMPLEMENTED); added `POST /api/attendance/nearest-office`; updated permissions: `attendance.checkin`, `attendance.checkout` (new), removed `baseline` from check-in/checkout; added notes for multi-session, geolocation validation, photo validation, check-in/check-out window logic, auto-detect shift, cross-day support
- **04_TDD.md:** Updated Permission-to-API Mapping — added `attendance.checkin` → POST /api/attendance/checkin, `attendance.checkout` → POST /api/attendance/checkout; marked `attendance.view-all`, `attendance.export`, `attendance.correct` as NOT YET IMPLEMENTED
- **03_Role_Matrix.md:** Added `attendance.checkin` and `attendance.checkout` permissions (✅ for all roles); marked `attendance.view-all`, `attendance.export`, `attendance.correct` as NOT YET IMPLEMENTED
- **02_FSD.md:** Updated §2.2.1 Check-in — auto-detect shift via `findApplicableShift`, window logic `[shiftStart - 15min, shiftEnd]`, present/late status, cross-day session validation, photo validation (UUID + extension + size ≤ 5MB), file move from tmp to evidence; updated §2.2.3 Check-out — window validation (block before `shiftEnd - 15min`), overtime calculation (after `shiftEnd + 15min`), duration calculation; updated §2.9.1 Location — added `is_active` field; updated §6.2 Attendance Validation — added new error messages for window validation, session state, photo format/size
- **05_ITL.md:** Updated T-027, T-028, T-029, T-030, T-031 with implementation notes; updated T-041 (Attendance model) to [x] with notes on new fields; renamed T-043 from EmployeeShift to UserShiftAssignment; added T-046 (Multi-Session Attendance), T-047 (Nearest Office Endpoint), T-048 (Enriched Today Status); updated Testing Checklist with new test scenarios; updated Summary table (93 total tasks, 39 completed)

### Key Implementation Changes vs Original Plan

| Area | Original Plan | Actual Implementation |
|------|--------------|----------------------|
| Unique constraint | `(user_id, attendance_date)` | `(user_id, date, shift_id)` — multi-session |
| Attendance fields | `check_in_method`, `qr_code_id`, `location_id`, correction fields | `shift_id`, `office_id`, `distance_meters`, `overtime_minutes`, `image_in`/`image_out` (file paths) |
| Check-in payload field | `foto` | `image` |
| Photo validation | UUID only | UUID + extension (.jpg/.jpeg/.png/.webp) + size ≤ 5MB + file move |
| Check-in logic | Simple duplicate check | Auto-detect shift, window logic, present/late status |
| Check-out logic | Allow early with warning | BLOCK before `shiftEnd - 15min`, overtime after `shiftEnd + 15min` |
| API endpoints | 8 endpoints | 4 endpoints (QR/history/stats/correct/late-stats NOT implemented) |
| Permissions | `attendance.view-all`, `attendance.export`, `attendance.correct` | `attendance.checkin`, `attendance.checkout` (new) |
| GET /attendance/today | Simple status object | Enriched: sessions[], current_action{}, todays_shifts[] |
| Nearest office | Not an endpoint | `POST /api/attendance/nearest-office` (new) |

### Validation Results

- ✅ ATTENDANCES ERD matches actual `models.Attendance` struct
- ✅ USER_SHIFT_ASSIGNMENTS ERD matches actual `models.UserShiftAssignment` struct
- ✅ API contract reflects actual implemented endpoints
- ✅ Permissions match actual route middleware (`attendance.checkin`, `attendance.checkout`)
- ✅ FSD check-in/check-out logic matches actual service implementation
- ✅ Error handling matches actual custom error messages
- ✅ ITL tasks updated with completion status and implementation notes
- ✅ All documents have YAML frontmatter with version 1.7.0

## v1.6.0 - 2026-06-08

### Change Password Endpoint Moved from Auth to Profile

- **04_TDD.md:** Removed `POST /api/auth/change-password` from §4.1 Authentication Endpoints; added `POST /api/me/change-password` to §4.11 Profile Endpoints; payload: `{ "new_password": "string", "confirm_password": "string" }` (confirm_password is frontend validation to ensure user typed new password correctly)
- **02_FSD.md:** Removed "Password Change" from Authentication hierarchy tree (§1 Functional Hierarchy); added "Change Password" under Profile hierarchy tree; removed §2.1.3 Password Change subsection; added new §2.14 Profile section with §2.14.1 Change Password (business logic: user input new password + confirm password, validate min 6 chars, validate confirm matches new password, hash and save, invalidate all sessions, user needs to re-login); updated §3.20 Profile Page — Change Password validation: "New password + confirm password, min 6 chars"
- **05_ITL.md:** Added T-016b task under Phase 1 (Profile) for `POST /api/me/change-password` endpoint

### Validation Results

- ✅ `POST /api/me/change-password` added to Profile endpoints with `{ new_password, confirm_password }` payload
- ✅ FSD hierarchy tree updated — Password Change moved from Auth to Profile
- ✅ FSD §2.1.3 removed, new §2.14 Profile section added with §2.14.1 Change Password (includes confirm password validation)
- ✅ FSD §3.20 Profile Page table updated with "New password + confirm password, min 6 chars"
- ✅ ITL T-016b added with correct FSD/TDD cross-references
- ✅ All documents have YAML frontmatter with version 1.6.0

## v1.5.0 - 2026-06-08

### Permission Cleanup & flexi_minutes Addition

- **03_Role_Matrix.md:** Removed `user.assign-role` permission (role assignment now inline via `PUT /api/users/:id`); removed `role.assign-permission` permission (permission assignment now inline via `PUT /api/roles/:id`); updated matrix table to remove both permission rows
- **04_TDD.md:** Added `flexi_minutes` field to SHIFTS ERD table; added `flexi_minutes` to POST /api/shifts and PUT /api/shifts/:id request payloads; updated Password Policy from "Minimum 8 characters" → "Minimum 6 characters"
- **02_FSD.md:** Updated §2.1.3 Password Change — changed min 8 char → min 6 char; updated §3.1 Login Page — changed min 8 → min 6; updated §3.3 Reset Password Page — changed min 8 → min 6; added "Flexi Minutes" row to §3.7 Shift Management Form

### Validation Results

- ✅ Permission list reduced from 34 to 32 entries
- ✅ `user.assign-role` and `role.assign-permission` removed from all documents
- ✅ `flexi_minutes` added to Shift ERD, API contract, and FSD form
- ✅ Password minLength consistently set to 6 across all documents
- ✅ All documents have YAML frontmatter with version 1.5.0

## v1.4.0 - 2026-06-08

### Endpoint Consolidation & Field Naming Alignment

- **04_TDD.md:** Removed `PUT /api/roles/:id/permissions` from §4.6 — permissions now handled inline via `PUT /api/roles/:id` which accepts `permissions` array; removed `POST /api/users/:id/roles` from §4.5 — roles now handled inline via `PUT /api/users/:id` which accepts `roles` array; updated Permission-to-API Mapping to reflect inline handling; updated upload response from `document_id` → `uuid`; updated check-in/check-out payloads from `{ "latitude", "longitude", "document_id" }` → `{ "lat", "lng", "foto" }`; removed `confirm_password` from reset-password payload; removed `join_date` from POST /api/users payload
- **02_FSD.md:** Updated §2.1.5 Reset Password — removed confirm_password (frontend handles confirmation validation); updated §2.2.1 Check-in — changed `document_id` to `foto` (UUID), `latitude`/`longitude` to `lat`/`lng`; updated §2.2.3 Check-out — same field changes; updated §2.6.1 Create Employee — removed `join_date` from input, noted it's auto-set to `created_at`; updated §3.9 User Management Form — removed "Join Date" row; updated §3.18 Location Management — removed "Coordinate Input" row (coordinates only via map click); updated §5.1 flow diagram — changed `document_id` to `foto` (uuid); updated §6.2 validation — changed "document_id tidak valid" to "foto tidak valid"
- **05_ITL.md:** Updated T-008 — changed TDD Ref to `PUT /api/roles/:id` (inline permissions); updated T-009 — changed TDD Ref to `PUT /api/users/:id` (inline roles); updated T-027, T-028, T-029, T-030, T-031 — changed all `document_id` references to `foto` (UUID from upload), `latitude`/`longitude` to `lat`/`lng`; updated Testing Checklist T-028 description
- **01_PRD.md:** Version bump only
- **03_Role_Matrix.md:** Version bump only

### Field Naming Convention Summary (Updated)

| Before | After | Used In |
|--------|-------|---------|
| `latitude` | `lat` | POST /api/attendance/checkin, POST /api/attendance/checkout |
| `longitude` | `lng` | POST /api/attendance/checkin, POST /api/attendance/checkout |
| `document_id` | `foto` | POST /api/attendance/checkin, POST /api/attendance/checkout |
| Upload response `document_id` | Upload response `uuid` | POST /api/upload |
| `confirm_password` | (removed) | POST /api/auth/reset-password |
| `join_date` | (auto-set) | POST /api/users |

### Validation Results

- ✅ Removed dedicated endpoints consolidated into their parent PUT endpoints
- ✅ Upload response field name matches implementation (`uuid` not `document_id`)
- ✅ Attendance payload field names match implementation (`lat`, `lng`, `foto`)
- ✅ Reset password payload simplified (no `confirm_password` — frontend-only validation)
- ✅ User creation no longer requires manual `join_date` input
- ✅ Location form no longer shows manual coordinate input fields
- ✅ All documents have YAML frontmatter with version 1.4.0

## v1.3.1 - 2026-06-08

### Mermaid Diagram Fixes

- **04_TDD.md:** Fixed System Architecture diagram — removed parentheses from node label `Face Recognition Service (Optional)` → `Face Recognition Service`; renamed subgraph `Data Layer` → `DataLayer` (spaces in subgraph IDs cause parse errors); updated style reference to match new subgraph name
- **04_TDD.md:** Fixed ERD diagram — removed inline comment `// OPTIONAL: only used if face recognition enabled` from `face_embedding` field (Mermaid ERD does not support `//` comments inside table definitions); cleaned up trailing whitespace throughout diagram

### Validation Results

- ✅ All mermaid diagrams in 04_TDD.md now parse correctly
- ✅ System Architecture diagram renders without errors
- ✅ ERD diagram renders without errors
- ✅ All documents have YAML frontmatter with version 1.3.1

## v1.3.0 - 2026-06-08

### Photo Upload Flow & Field Naming Convention

- **01_PRD.md:** Version bump only (no content changes)
- **02_FSD.md:** Updated §2.2.1 Check-in (Geotagging + Photo Evidence) — added two-step upload flow: frontend uploads photo via `POST /api/upload` first, receives `document_id`, then sends `document_id` to check-in endpoint; updated §2.2.3 Check-out with same flow; updated §3.5 screen elements — added Submit Button with validation (enabled only if photo uploaded/document_id exists); updated §5.1 flow diagram — added upload step and document_id flow; updated §6.2 validation — added `document_id` tidak valid error
- **03_Role_Matrix.md:** Version bump only (no permission changes needed)
- **04_TDD.md:** Updated §4.2 Attendance Endpoints — changed `photo: "base64"` to `document_id: "uuid-string"` for check-in/check-out payloads; updated §4.3 Shift Endpoints — changed `user_ids` → `users`, `shift_id` → `shift`, query param `user_id` → `user`; updated §4.4 Leave Endpoints — changed `leave_type_id` → `leave_type`; updated §4.5 User Management — changed `role_ids` → `roles`; updated §4.6 UAM — changed `permission_ids` → `permissions`; updated §4.8 QR Code — changed `office_id` → `office` in payload and query param; updated §4.10 Report — changed query param `user_id` → `user`; updated §4.13 System — updated `/api/upload` response to return `document_id` instead of `uuid`; updated Permission-to-API Mapping — changed `user_id=all` → `user=all` query params; added note to §4.2 explaining the two-step upload flow
- **05_ITL.md:** Updated T-021, T-022 (shift assign — `shift`, `users` fields); updated T-008 (role-permission — `permissions` field); updated T-009 (user-role — `roles` field); updated T-027, T-028, T-029, T-030, T-031 (attendance check-in/out — `document_id` flow, two-step upload); updated T-033 (QR code — `office` field); updated T-057 (leave request — `leave_type` field); updated all TDD Ref notes to reflect new field names

### Field Naming Convention Summary

All request payload fields that previously used `_id` suffix have been renamed for cleaner validation error messages:

| Before | After | Used In |
|--------|-------|---------|
| `shift_id` | `shift` | POST /api/shifts/assign |
| `user_ids` | `users` | POST /api/shifts/assign |
| `leave_type_id` | `leave_type` | POST /api/leave |
| `role_ids` | `roles` | POST /api/users/:id/roles |
| `permission_ids` | `permissions` | PUT /api/roles/:id/permissions |
| `office_id` | `office` | POST /api/qr-codes/generate, GET /api/qr-codes |
| `photo` (base64) | `document_id` (uuid) | POST /api/attendance/checkin, POST /api/attendance/checkout |
| `user_id` (query) | `user` (query) | GET /api/shifts/schedule, GET /api/reports/attendance, GET /api/attendance/late-statistics |

### Validation Results

- ✅ All ITL tasks have valid FSD Ref matching FSD section numbers
- ✅ All ITL tasks have valid TDD Ref matching TDD API endpoints
- ✅ Photo upload flow is consistent across FSD (§2.2.1, §2.2.3, §5.1) and TDD (§4.2, §4.13)
- ✅ Field naming convention applied consistently to all request payloads
- ✅ All documents have YAML frontmatter with version 1.3.0

## v1.2.0 - 2026-06-08

### Face Recognition Separated from Core Check-in

- **01_PRD.md:** Moved Face Recognition from Must Have to Could Have; added Photo Evidence as Must Have; updated user stories US-001 (geotagging + photo evidence), US-019 (profile photo instead of face photo); updated goals, success metrics, constraints, and assumptions to reflect face recognition as optional
- **02_FSD.md:** Renamed §2.2.1 to "Check-in (Geotagging + Photo Evidence)" — removed face recognition step, changed to capture/upload photo as evidence; renamed §2.2.3 to "Check-out (Geotagging + Photo Evidence)" — same change; added new §2.13 "Face Recognition (Optional / Could Have)" with 3 sub-sections (§2.13.1 Upload Face Photo, §2.13.2 Face Recognition during Check-in, §2.13.3 Face Recognition during Check-out); updated §2.6.1 and §2.6.2 to reflect profile photo (face embedding only if face recognition enabled); updated §3.5 screen elements (added Upload Button, Photo Preview); updated §3.9 and §3.20 to use Profile Photo instead of Face Photo; updated §5.1 flow diagram to use photo evidence instead of face recognition; updated §6.2 and §6.4 validation tables; updated use case diagram with UC28 (optional face recognition)
- **03_Role_Matrix.md:** Version bump only (no permission changes needed)
- **04_TDD.md:** Marked Face Recognition tech stack entry as optional; updated ERD to mark `face_embedding` as optional; added note to §4.2 Attendance Endpoints that photo is for evidence only; renamed face-photo endpoint response message; added new §4.14 "Face Recognition Endpoints (Optional)" with enable/status endpoints; updated environment variables with `FACE_RECOGNITION_ENABLED` flag
- **05_ITL.md:** Moved face recognition tasks (T-026, T-028, T-031) to new Phase 6 (Face Recognition — Optional, P4); created new T-024 (Upload Profile Photo), T-025 (Profile Photo Component), T-027 (Check-in + Photo Evidence), T-028 (Photo Upload Handler), T-030 (Check-out + Photo Evidence), T-031 (Check-out Page); added Phase 6 with 6 new tasks (T-100 through T-105); updated all FSD/TDD cross-references to match new section numbers; updated testing checklist; updated summary table (87 total tasks across 6 phases)

### Validation Results

- ✅ All ITL tasks have valid FSD Ref matching FSD section numbers (§2.2.1, §2.2.3, §2.6.2, §2.13.x, §3.5, §3.20, §5.1)
- ✅ All ITL tasks have valid TDD Ref matching TDD API endpoints (POST /api/attendance/checkin, POST /api/attendance/checkout, POST /api/users/:id/face-photo, etc.)
- ✅ Face recognition is completely separated from check-in flow — check-in works with photo evidence only
- ✅ Photo upload during check-in documented as evidence (not face matching) in all documents
- ✅ Phase 6 tasks use correct FSD §2.13.x references and TDD optional endpoint references
- ✅ All documents have YAML frontmatter with version 1.2.0

## v1.1.0 - 2026-06-08

### Format & Structure Update

- **All documents:** Added YAML frontmatter with `title`, `version`, `created`, and `last_modified` fields
- **All documents:** Updated version from `1.0.0` to `1.1.0`
- **Folder structure:** Restructured from `docs/` to `plans/` with `current/`, `archive/`, and `CHANGELOG.md`
- **Archive:** Moved all v1.0.0 documents to `plans/archive/v1.0.0/`

### 02_FSD.md

- **§2.x Functional Requirements:** Added `§` prefix to ALL section numbers (e.g., `2.1` → `§2.1`, `2.2.1` → `§2.2.1`)
- **§3.x User Interaction & Screen Elements:** Added sub-numbering to ALL screen sections (`§3.1` through `§3.20`)
  - Login Page → §3.1
  - Forgot Password Page → §3.2
  - Reset Password Page → §3.3
  - Karyawan Dashboard → §3.4
  - Attendance Check-in (Geotagging) → §3.5
  - Attendance Check-in (QR Code) → §3.6
  - Shift Management Form → §3.7
  - Leave Request Form → §3.8
  - User Management Form → §3.9
  - Role Management Form → §3.10
  - Attendance Correction Form → §3.11
  - QR Code Management → §3.12
  - HR Dashboard → §3.13
  - Admin Dashboard → §3.14
  - Attendance History → §3.15
  - Leave Balance → §3.16
  - Shift Schedule View → §3.17
  - Location Management → §3.18
  - Audit Log View → §3.19
  - Profile Page → §3.20
- **§5.x Feature Logic Flow:** Added sub-numbering to ALL flow diagrams
  - Check-in (Geotagging + Face Recognition) Flow → §5.1
  - Check-in (QR Code) Flow → §5.2
  - Leave Request Flow → §5.3
- **§6.x Error Handling & Validation:** Added sub-numbering and grouped validation tables by feature
  - Authentication Validation → §6.1
  - Attendance Validation → §6.2
  - Leave Validation → §6.3
  - User & Upload Validation → §6.4
  - General Validation → §6.5

### 04_TDD.md

- **§4.x API Contract:** Added sub-numbering to ALL endpoint groups
  - Authentication Endpoints → §4.1
  - Attendance Endpoints → §4.2
  - Shift Endpoints → §4.3
  - Leave Endpoints → §4.4
  - User Management Endpoints → §4.5
  - UAM (Role & Permissions) Endpoints → §4.6
  - Location Endpoints → §4.7
  - QR Code Management Endpoints → §4.8
  - Dashboard Endpoints → §4.9
  - Report Endpoints → §4.10
  - Profile Endpoints → §4.11
  - Notification Endpoints → §4.12
  - System Endpoints → §4.13

### 05_ITL.md

- **Format:** Converted ALL tasks from table format to list format
- **Structure:** Each task now follows the standard block format:
  ```
  ### T-XXX: Task Name
  - **Feature/Module:** ...
  - **Priority:** ...
  - **Estimated Effort:** ...
  - **Status:** ...
  - **FSD Ref:** ...
  - **TDD Ref:** ...
  ```
- **Cross-references:** Updated ALL FSD Ref to use new `§`-prefixed section numbers
- **Cross-references:** Updated ALL TDD Ref to use new `§`-prefixed endpoint group numbers
- **Total tasks:** 81 tasks (T-001 through T-097, with gaps for future expansion)

### Validation Results

- ✅ All 81 ITL tasks have valid FSD Ref matching FSD section numbers
- ✅ All 81 ITL tasks have valid TDD Ref matching TDD API endpoint groups
- ✅ All FSD sub-numbering is consistent (§2.x, §3.x, §5.x, §6.x)
- ✅ All TDD API groups are numbered (§4.1 through §4.13)
- ✅ All documents have YAML frontmatter with correct version metadata

## v1.0.0 - 2026-05-29

- Initial generation of HadirYuk project planner documents
- 01_PRD.md — Product Requirements Document
- 02_FSD.md — Functional Specification Document
- 03_Role_Matrix.md — Role & Permissions Matrix
- 04_TDD.md — Technical Design Document
- 05_ITL.md — Implementation Task List (81 tasks across 5 phases)
