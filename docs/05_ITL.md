# Implementation Task List (ITL)

> **Project:** HadirYuk — Sistem Absensi Karyawan
> **Stack:** Go (Gin) + Vue 3 + MySQL + Redis
> **Generated:** 2026-05-29

---

## Implementation Phases

| Phase       | Name                       | Description                                                |
| ----------- | -------------------------- | ---------------------------------------------------------- |
| **Phase 1** | Foundation (DONE)          | Auth, RBAC, Data Master, Profil, Notifikasi, Upload, Peta  |
| **Phase 2** | Core Attendance (P0)       | Shift assignment, face photo, check-in/out geotagging & QR |
| **Phase 3** | Employee Self-Service (P1) | Riwayat absensi, dasbor karyawan, cuti                     |
| **Phase 4** | HR Operations (P2)         | Dasbor HR, koreksi absensi, jadwal shift, laporan          |
| **Phase 5** | Analytics & Polish (P3-P4) | Statistik keterlambatan, dasbor admin, audit log           |

---

## Phase 1: Foundation — COMPLETED ✅

| Task ID | Feature/Module  | Task Description                  | Priority | Effort | Status | FSD Ref      | TDD Ref                             | RBAC Ref               |
| ------- | --------------- | --------------------------------- | -------- | ------ | ------ | ------------ | ----------------------------------- | ---------------------- |
| T-001   | Authentication  | Login endpoint JWT RS256 + JWKS   | P0       | 3h     | [x]    | FSD §2.1.1   | POST /api/auth/login                | baseline               |
| T-002   | Authentication  | Logout endpoint + session cleanup | P0       | 1h     | [x]    | FSD §2.1.2   | POST /api/auth/logout               | baseline               |
| T-003   | Authentication  | Refresh token endpoint            | P0       | 2h     | [x]    | FSD §2.1.6   | POST /api/auth/refresh              | baseline               |
| T-004   | Authentication  | Forgot password + email async     | P0       | 3h     | [x]    | FSD §2.1.4   | POST /api/auth/forgot-password      | baseline               |
| T-005   | Authentication  | Reset password + token validation | P0       | 2h     | [x]    | FSD §2.1.5   | POST /api/auth/reset-password       | baseline               |
| T-006   | UAM             | CRUD Role endpoints               | P0       | 3h     | [x]    | FSD §2.7.1   | CRUD /api/roles                     | role.\*                |
| T-007   | UAM             | CRUD Permission endpoints         | P0       | 2h     | [x]    | FSD §2.7.2   | CRUD /api/permissions               | permission.\*          |
| T-008   | UAM             | Role-Permission mapping           | P0       | 2h     | [x]    | FSD §2.7.2   | PUT /api/roles/:id/permissions      | role.assign-permission |
| T-009   | UAM             | User-Role assignment              | P0       | 2h     | [x]    | FSD §2.7.3   | POST /api/users/:id/roles           | user.assign-role       |
| T-010   | UAM             | Permission middleware + cache     | P0       | 4h     | [x]    | FSD §2.7     | Middleware                          | All RBAC               |
| T-011   | Data Master     | CRUD Shift endpoints              | P1       | 3h     | [x]    | FSD §2.3.1   | CRUD /api/shifts                    | shift.\*               |
| T-012   | Data Master     | CRUD Office Location endpoints    | P1       | 3h     | [x]    | FSD §2.9.1   | CRUD /api/locations                 | location.\*            |
| T-013   | Data Master     | CRUD Leave Type endpoints         | P1       | 2h     | [x]    | FSD §2.4     | CRUD /api/leave/types               | leave.manage-types     |
| T-014   | User Management | CRUD User with auto-profile       | P0       | 4h     | [x]    | FSD §2.6.1   | CRUD /api/users                     | user.\*                |
| T-015   | Profile         | View/Edit profile (/me)           | P1       | 2h     | [x]    | FSD Profile  | GET/PUT /api/me                     | baseline               |
| T-016   | Profile         | Upload avatar                     | P1       | 2h     | [x]    | FSD §2.6.2   | POST /api/upload                    | baseline               |
| T-017   | Notifications   | CRUD notifications + mark read    | P2       | 3h     | [x]    | FSD -        | CRUD /api/notifications             | -                      |
| T-018   | Upload          | File upload with UUID             | P1       | 2h     | [x]    | FSD -        | POST /api/upload                    | baseline               |
| T-019   | Map             | Leaflet map component (frontend)  | P1       | 3h     | [x]    | FSD Location | Frontend                            | -                      |
| T-020   | Health          | Health check + JWKS endpoint      | P0       | 1h     | [x]    | FSD -        | GET /health, /.well-known/jwks.json | -                      |

---

## Phase 2: Core Attendance (P0)

| Task ID | Feature/Module       | Task Description                                                                                          | Priority | Effort | Status | FSD Ref            | TDD Ref                                          | RBAC Ref                   |
| ------- | -------------------- | --------------------------------------------------------------------------------------------------------- | -------- | ------ | ------ | ------------------ | ------------------------------------------------ | -------------------------- |
| T-021   | Shift Assignment     | Fullstack: assign shift to employee, get employee schedule (BE endpoints + FE form)                       | P0       | 8h     | [x]    | FSD §2.3.2-2.3.3   | POST /api/shifts/assign, GET /api/shifts/schedule | shift.assign, baseline     |
| T-024   | Face Photo           | Fullstack: upload/capture face photo (admin only), generate face embedding, face recognition service (GoCV) | P0 | 13h | [x] | FSD §2.6.2, §2.2.1 | PUT/DELETE /api/users/:id/face-photo, Service layer | user.update |
| T-027   | Attendance Check-in  | Fullstack: check-in geotagging + Haversine validation + face recognition match (BE + FE page w/ camera)   | P0       | 13h    | [ ]    | FSD §2.2.1         | POST /api/attendance/checkin                     | baseline                   |
| T-030   | Attendance Check-out | Fullstack: check-out geotagging + duration calculation + optional face recognition (BE + FE page)         | P0       | 8h     | [ ]    | FSD §2.2.3         | POST /api/attendance/checkout                    | baseline                   |
| T-033   | QR Code Management   | Fullstack: generate QR with signature+expiry, list+revoke, validation service (BE + FE management page)   | P0       | 10h    | [ ]    | FSD §2.11          | POST /api/qr-codes/generate, GET/POST revoke      | qrcode.\*                  |
| T-037   | QR Check-in          | Fullstack: check-in via QR code validation + FE QR scanner (vue-qrcode-reader)                            | P0       | 5h     | [ ]    | FSD §2.2.2         | POST /api/attendance/checkin/qr                  | baseline                   |
| T-039   | QR Check-out         | Fullstack: check-out via QR code (reuse scanner component)                                                | P0       | 2h     | [ ]    | FSD §2.2.3b        | POST /api/attendance/checkout/qr                 | baseline                   |
| T-041   | Data Models          | Backend: GORM models + migrations for Attendance, QRCode, EmployeeShift, LeaveBalance, LeaveRequest       | P0       | 6h     | [ ]    | FSD §2.2-§2.4      | TDD ERD                                          | -                          |

---

## Phase 3: Employee Self-Service (P1)

| Task ID | Feature/Module     | Task Description                                                                             | Priority | Effort | Status | FSD Ref    | TDD Ref                     | RBAC Ref |
| ------- | ------------------ | -------------------------------------------------------------------------------------------- | -------- | ------ | ------ | ---------- | --------------------------- | -------- |
| T-050   | Attendance History | Fullstack: history endpoint (self), today status, monthly stats + FE page (filter, pagination, detail modal) | P1 | 9h     | [ ]    | FSD §2.2.4 | GET /api/attendance/history | baseline |
| T-055   | Employee Dashboard | Fullstack: dashboard endpoint + FE page (status, summary, schedule)                          | P1       | 6h     | [ ]    | FSD §2.8.1 | GET /api/dashboard/employee | baseline |
| T-057   | Leave Request      | Fullstack: submit leave request + balance deduction + FE form with validation                | P1       | 5h     | [ ]    | FSD §2.4.1 | POST /api/leave             | baseline |
| T-059   | Leave Balance      | Fullstack: leave balance endpoint, auto-init on employee creation, yearly reset cron + FE page (progress bar, history) | P1 | 7h | [ ]    | FSD §2.4.2 | GET /api/leave/balance      | baseline |

---

## Phase 4: HR Operations (P2)

| Task ID | Feature/Module        | Task Description                                                                                | Priority | Effort | Status | FSD Ref     | TDD Ref                                  | RBAC Ref            |
| ------- | --------------------- | ----------------------------------------------------------------------------------------------- | -------- | ------ | ------ | ----------- | ---------------------------------------- | ------------------- |
| T-070   | HR Dashboard          | Fullstack: HR dashboard endpoint (stats, chart, not-attended, leaves) + FE page (cards, chart, tables) | P2    | 7h     | [ ]    | FSD §2.8.2  | GET /api/dashboard/hr                    | dashboard.view-hr   |
| T-072   | Attendance Correction | Fullstack: correct attendance endpoint + FE correction form (datetime picker, reason)           | P2       | 4h     | [ ]    | FSD §2.2.5  | PUT /api/attendance/:id/correct          | attendance.correct  |
| T-074   | Shift Schedule        | Fullstack: employee schedule endpoint (all employees) + FE calendar view (month, color-coded)   | P2       | 5h     | [ ]    | FSD §2.3.3  | GET /api/shifts/schedule                 | baseline            |
| T-076   | Attendance Report     | Fullstack: report endpoint with filters, export Excel (.xlsx) + PDF + FE page (filters, table, export buttons) | P2 | 12h    | [ ]    | FSD §2.10   | GET /api/reports/attendance              | report.\*           |
| T-080   | Leave Report          | Fullstack: leave report endpoint with filters, export Excel/PDF + FE page                       | P2       | 6h     | [ ]    | FSD §2.10.4 | GET /api/reports/leave                   | report.view         |
| T-083   | Attendance View-All   | Backend: attendance history with user_id=all filter                                              | P2       | 1h     | [ ]    | FSD §2.2.4  | GET /api/attendance/history?user_id=all  | attendance.view-all |
| T-084   | Leave View-All        | Backend: leave history with user_id=all filter                                                   | P2       | 1h     | [ ]    | FSD §2.4    | GET /api/leave?user_id=all               | leave.view-all      |

---

## Phase 5: Analytics & Polish (P3-P4)

| Task ID | Feature/Module  | Task Description                                                                        | Priority | Effort | Status | FSD Ref     | TDD Ref                             | RBAC Ref             |
| ------- | --------------- | --------------------------------------------------------------------------------------- | -------- | ------ | ------ | ----------- | ----------------------------------- | -------------------- |
| T-090   | Late Statistics | Fullstack: late statistics endpoint + FE page (chart, table, trend)                     | P3       | 6h     | [ ]    | FSD §2.5.1  | GET /api/attendance/late-statistics | late-statistic.view  |
| T-092   | Admin Dashboard | Fullstack: admin dashboard endpoint (system stats, activity, health) + FE page          | P3       | 4h     | [ ]    | FSD Admin   | GET /api/dashboard/admin            | dashboard.view-admin |
| T-094   | Audit Log       | Fullstack: audit log middleware, GORM model + migration, endpoint with filters + FE page (filters, detail modal) | P4 | 10h    | [ ]    | FSD §2.12.1 | GET /api/audit-logs                 | audit.view           |

---

## Testing Checklist

| Task ID | Test Scenario                 | Expected Result                                | Status |
| ------- | ----------------------------- | ---------------------------------------------- | ------ |
| T-021   | Assign shift to employee      | Employee ter-assign, record di employee_shifts | [ ]    |
| T-027   | Check-in dalam radius         | Absensi tersimpan dengan status present        | [ ]    |
| T-027   | Check-in luar radius          | Error "Anda berada di luar area kantor"        | [ ]    |
| T-028   | Face match                    | Absensi tersimpan                              | [ ]    |
| T-028   | Face tidak match              | Error "Wajah tidak dikenali"                   | [ ]    |
| T-030   | Check-out setelah check-in    | Record updated dengan duration                 | [ ]    |
| T-030   | Check-out sebelum check-in    | Error "Belum check-in"                         | [ ]    |
| T-033   | Generate QR code              | QR code tersimpan dengan signature + expiry    | [ ]    |
| T-037   | Check-in QR valid             | Absensi tersimpan                              | [ ]    |
| T-037   | Check-in QR expired           | Error "QR Code sudah expired"                  | [ ]    |
| T-057   | Submit cuti valid             | Leave request tersimpan, balance berkurang     | [ ]    |
| T-057   | Submit cuti saldo tidak cukup | Error "Sisa cuti tidak mencukupi"              | [ ]    |
| T-057   | Submit cuti overlap           | Error "Tanggal cuti overlap"                   | [ ]    |
| T-072   | Koreksi absensi               | Record updated, audit log tercatat             | [ ]    |
| T-077   | Export Excel                  | File .xlsx terdownload dengan data benar       | [ ]    |
| T-078   | Export PDF                    | File .pdf terdownload dengan data benar        | [ ]    |
| T-090   | Late statistics               | Data keterlambatan akurat per employee         | [ ]    |
| T-094   | Audit log CRUD                | Setiap perubahan tercatat di audit log         | [ ]    |

---

## Summary

| Phase                     | Total Tasks | Completed | Remaining | Est. Hours |
| ------------------------- | ----------- | --------- | --------- | ---------- |
| Phase 1 (Foundation)      | 20          | 20        | 0         | ~48h       |
| Phase 2 (Core Attendance) | 8           | 2         | 6         | ~57h       |
| Phase 3 (Self-Service)    | 4           | 0         | 4         | ~27h       |
| Phase 4 (HR Operations)   | 7           | 0         | 7         | ~36h       |
| Phase 5 (Analytics)       | 3           | 0         | 3         | ~20h       |
| **Total**                 | **42**      | **21**    | **21**    | **~188h**  |
