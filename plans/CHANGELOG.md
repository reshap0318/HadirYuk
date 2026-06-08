# CHANGELOG

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
