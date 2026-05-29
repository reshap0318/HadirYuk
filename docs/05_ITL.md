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

| Task ID | Feature/Module       | Task Description                                             | Priority | Effort | Status | FSD Ref            | TDD Ref                                          | RBAC Ref                   |
| ------- | -------------------- | ------------------------------------------------------------ | -------- | ------ | ------ | ------------------ | ------------------------------------------------ | -------------------------- |
| T-021   | Shift Assignment     | Backend: assign shift to employee endpoint                   | P0       | 3h     | [x]    | FSD §2.3.2         | POST /api/shifts/assign                          | shift.assign               |
| T-022   | Shift Assignment     | Frontend: form assign shift ke multiple employee             | P0       | 3h     | [x]    | FSD §2.3.2         | POST /api/shifts/assign                          | shift.assign               |
| T-023   | Shift Assignment     | Backend: get employee schedule endpoint                      | P0       | 2h     | [x]    | FSD §2.3.3         | GET /api/shifts/schedule                         | baseline                   |
| T-024   | Face Photo           | Backend: upload face photo + face embedding generation       | P0       | 4h     | [ ]    | FSD §2.6.2         | POST /api/users/:id/face-photo                   | baseline                   |
| T-025   | Face Photo           | Frontend: face photo capture/upload component                | P0       | 3h     | [ ]    | FSD §2.6.2         | POST /api/users/:id/face-photo                   | baseline                   |
| T-026   | Face Photo           | Backend: face recognition service (GoCV/face-recognition-go) | P0       | 6h     | [ ]    | FSD §2.2.1         | Service layer                                    | baseline                   |
| T-027   | Attendance Check-in  | Backend: check-in geotagging + Haversine validation          | P0       | 4h     | [ ]    | FSD §2.2.1         | POST /api/attendance/checkin                     | baseline                   |
| T-028   | Attendance Check-in  | Backend: check-in with face recognition match                | P0       | 4h     | [ ]    | FSD §2.2.1         | POST /api/attendance/checkin                     | baseline                   |
| T-029   | Attendance Check-in  | Frontend: check-in page (location + camera + capture)        | P0       | 5h     | [ ]    | FSD §2.2.1         | POST /api/attendance/checkin                     | baseline                   |
| T-030   | Attendance Check-out | Backend: check-out geotagging + duration calculation         | P0       | 3h     | [ ]    | FSD §2.2.3         | POST /api/attendance/checkout                    | baseline                   |
| T-031   | Attendance Check-out | Backend: check-out with optional face recognition            | P0       | 2h     | [ ]    | FSD §2.2.3         | POST /api/attendance/checkout                    | baseline                   |
| T-032   | Attendance Check-out | Frontend: check-out page (location + camera + capture)       | P0       | 3h     | [ ]    | FSD §2.2.3         | POST /api/attendance/checkout                    | baseline                   |
| T-033   | QR Code              | Backend: generate QR code with signature + expiry            | P0       | 3h     | [ ]    | FSD §2.11.1        | POST /api/qr-codes/generate                      | qrcode.generate            |
| T-034   | QR Code              | Backend: view active QR codes + revoke                       | P0       | 2h     | [ ]    | FSD §2.11.2-2.11.3 | GET /api/qr-codes, POST /api/qr-codes/:id/revoke | qrcode.view, qrcode.revoke |
| T-035   | QR Code              | Frontend: QR code management page (generate, list, revoke)   | P0       | 3h     | [ ]    | FSD §2.11          | QR Code endpoints                                | qrcode.\*                  |
| T-036   | QR Code              | Backend: QR code validation service (signature + expiry)     | P0       | 2h     | [ ]    | FSD §2.2.2         | Service layer                                    | baseline                   |
| T-037   | QR Check-in          | Backend: check-in via QR code                                | P0       | 2h     | [ ]    | FSD §2.2.2         | POST /api/attendance/checkin/qr                  | baseline                   |
| T-038   | QR Check-in          | Frontend: QR scanner component (vue-qrcode-reader)           | P0       | 3h     | [ ]    | FSD §2.2.2         | POST /api/attendance/checkin/qr                  | baseline                   |
| T-039   | QR Check-out         | Backend: check-out via QR code                               | P0       | 1h     | [ ]    | FSD §2.2.3b        | POST /api/attendance/checkout/qr                 | baseline                   |
| T-040   | QR Check-out         | Frontend: QR scanner for check-out (reuse component)         | P0       | 1h     | [ ]    | FSD §2.2.3b        | POST /api/attendance/checkout/qr                 | baseline                   |
| T-041   | Attendance Model     | Backend: Attendance GORM model + migration                   | P0       | 2h     | [ ]    | FSD §2.2           | TDD ERD                                          | -                          |
| T-042   | Attendance Model     | Backend: QRCode GORM model + migration                       | P0       | 1h     | [ ]    | FSD §2.11          | TDD ERD                                          | -                          |
| T-043   | Attendance Model     | Backend: EmployeeShift GORM model + migration                | P0       | 1h     | [ ]    | FSD §2.3.2         | TDD ERD                                          | -                          |
| T-044   | Attendance Model     | Backend: LeaveBalance GORM model + migration                 | P0       | 1h     | [ ]    | FSD §2.4.2         | TDD ERD                                          | -                          |
| T-045   | Attendance Model     | Backend: LeaveRequest GORM model + migration                 | P0       | 1h     | [ ]    | FSD §2.4.1         | TDD ERD                                          | -                          |

---

## Phase 3: Employee Self-Service (P1)

| Task ID | Feature/Module     | Task Description                                              | Priority | Effort | Status | FSD Ref    | TDD Ref                     | RBAC Ref |
| ------- | ------------------ | ------------------------------------------------------------- | -------- | ------ | ------ | ---------- | --------------------------- | -------- |
| T-050   | Attendance History | Backend: attendance history endpoint (self)                   | P1       | 2h     | [ ]    | FSD §2.2.4 | GET /api/attendance/history | baseline |
| T-051   | Attendance History | Backend: today status endpoint                                | P1       | 1h     | [ ]    | FSD §2.2.4 | GET /api/attendance/today   | baseline |
| T-052   | Attendance History | Backend: monthly stats endpoint                               | P1       | 1h     | [ ]    | FSD §2.2.4 | GET /api/attendance/stats   | baseline |
| T-053   | Attendance History | Frontend: attendance history page (filter, pagination)        | P1       | 3h     | [ ]    | FSD §2.2.4 | GET /api/attendance/history | baseline |
| T-054   | Attendance History | Frontend: attendance detail modal                             | P1       | 2h     | [ ]    | FSD §2.2.4 | GET /api/attendance/history | baseline |
| T-055   | Employee Dashboard | Backend: employee dashboard endpoint                          | P1       | 2h     | [ ]    | FSD §2.8.1 | GET /api/dashboard/employee | baseline |
| T-056   | Employee Dashboard | Frontend: employee dashboard page (status, summary, schedule) | P1       | 4h     | [ ]    | FSD §2.8.1 | GET /api/dashboard/employee | baseline |
| T-057   | Leave Request      | Backend: submit leave request + balance deduction             | P1       | 3h     | [ ]    | FSD §2.4.1 | POST /api/leave             | baseline |
| T-058   | Leave Request      | Frontend: leave request form + validation                     | P1       | 2h     | [ ]    | FSD §2.4.1 | POST /api/leave             | baseline |
| T-059   | Leave Balance      | Backend: leave balance endpoint                               | P1       | 2h     | [ ]    | FSD §2.4.2 | GET /api/leave/balance      | baseline |
| T-060   | Leave Balance      | Frontend: leave balance page (progress bar, history)          | P1       | 2h     | [ ]    | FSD §2.4.2 | GET /api/leave/balance      | baseline |
| T-061   | Leave Balance      | Backend: auto-init leave balance on employee creation         | P1       | 1h     | [ ]    | FSD §2.4.2 | TDD Leave Balance Strategy  | -        |
| T-062   | Leave Balance      | Backend: yearly reset cron job (1 Januari)                    | P1       | 2h     | [ ]    | FSD §2.4.2 | TDD Leave Balance Strategy  | -        |

---

## Phase 4: HR Operations (P2)

| Task ID | Feature/Module        | Task Description                                                    | Priority | Effort | Status | FSD Ref     | TDD Ref                                  | RBAC Ref            |
| ------- | --------------------- | ------------------------------------------------------------------- | -------- | ------ | ------ | ----------- | ---------------------------------------- | ------------------- |
| T-070   | HR Dashboard          | Backend: HR dashboard endpoint (stats, chart, not-attended, leaves) | P2       | 3h     | [ ]    | FSD §2.8.2  | GET /api/dashboard/hr                    | dashboard.view-hr   |
| T-071   | HR Dashboard          | Frontend: HR dashboard page (cards, chart, tables)                  | P2       | 4h     | [ ]    | FSD §2.8.2  | GET /api/dashboard/hr                    | dashboard.view-hr   |
| T-072   | Attendance Correction | Backend: correct attendance endpoint                                | P2       | 2h     | [ ]    | FSD §2.2.5  | PUT /api/attendance/:id/correct          | attendance.correct  |
| T-073   | Attendance Correction | Frontend: correction form (datetime picker, reason)                 | P2       | 2h     | [ ]    | FSD §2.2.5  | PUT /api/attendance/:id/correct          | attendance.correct  |
| T-074   | Shift Schedule        | Backend: employee schedule endpoint (all employees)                 | P2       | 2h     | [ ]    | FSD §2.3.3  | GET /api/shifts/schedule                 | baseline            |
| T-075   | Shift Schedule        | Frontend: shift schedule calendar view (month, color-coded)         | P2       | 3h     | [ ]    | FSD §2.3.3  | GET /api/shifts/schedule                 | baseline            |
| T-076   | Attendance Report     | Backend: attendance report endpoint with filters                    | P2       | 3h     | [ ]    | FSD §2.10.1 | GET /api/reports/attendance              | report.view         |
| T-077   | Attendance Report     | Backend: export to Excel (.xlsx)                                    | P2       | 3h     | [ ]    | FSD §2.10.2 | GET /api/reports/attendance/export/excel | report.export-excel |
| T-078   | Attendance Report     | Backend: export to PDF                                              | P2       | 3h     | [ ]    | FSD §2.10.3 | GET /api/reports/attendance/export/pdf   | report.export-pdf   |
| T-079   | Attendance Report     | Frontend: report page (filters, table, export buttons)              | P2       | 3h     | [ ]    | FSD §2.10   | Report endpoints                         | report.\*           |
| T-080   | Leave Report          | Backend: leave report endpoint with filters                         | P2       | 2h     | [ ]    | FSD §2.10.4 | GET /api/reports/leave                   | report.view         |
| T-081   | Leave Report          | Backend: leave report export Excel/PDF                              | P2       | 2h     | [ ]    | FSD §2.10.4 | GET /api/reports/leave/export/\*         | report.export-\*    |
| T-082   | Leave Report          | Frontend: leave report page                                         | P2       | 2h     | [ ]    | FSD §2.10.4 | GET /api/reports/leave                   | report.view         |
| T-083   | Attendance View-All   | Backend: attendance history with user_id=all filter                 | P2       | 1h     | [ ]    | FSD §2.2.4  | GET /api/attendance/history?user_id=all  | attendance.view-all |
| T-084   | Leave View-All        | Backend: leave history with user_id=all filter                      | P2       | 1h     | [ ]    | FSD §2.4    | GET /api/leave?user_id=all               | leave.view-all      |

---

## Phase 5: Analytics & Polish (P3-P4)

| Task ID | Feature/Module  | Task Description                                                   | Priority | Effort | Status | FSD Ref     | TDD Ref                             | RBAC Ref             |
| ------- | --------------- | ------------------------------------------------------------------ | -------- | ------ | ------ | ----------- | ----------------------------------- | -------------------- |
| T-090   | Late Statistics | Backend: late statistics endpoint                                  | P3       | 3h     | [ ]    | FSD §2.5.1  | GET /api/attendance/late-statistics | late-statistic.view  |
| T-091   | Late Statistics | Frontend: late statistics page (chart, table, trend)               | P3       | 3h     | [ ]    | FSD §2.5.1  | GET /api/attendance/late-statistics | late-statistic.view  |
| T-092   | Admin Dashboard | Backend: admin dashboard endpoint (system stats, activity, health) | P3       | 2h     | [ ]    | FSD Admin   | GET /api/dashboard/admin            | dashboard.view-admin |
| T-093   | Admin Dashboard | Frontend: admin dashboard page                                     | P3       | 2h     | [ ]    | FSD Admin   | GET /api/dashboard/admin            | dashboard.view-admin |
| T-094   | Audit Log       | Backend: audit log middleware (record CRUD changes)                | P4       | 4h     | [ ]    | FSD §2.12.1 | GET /api/audit-logs                 | audit.view           |
| T-095   | Audit Log       | Backend: audit log endpoint with filters                           | P4       | 2h     | [ ]    | FSD §2.12.1 | GET /api/audit-logs                 | audit.view           |
| T-096   | Audit Log       | Frontend: audit log page (filters, detail modal)                   | P4       | 3h     | [ ]    | FSD §2.12.1 | GET /api/audit-logs                 | audit.view           |
| T-097   | Audit Log       | Backend: AuditLog GORM model + migration                           | P4       | 1h     | [ ]    | FSD §2.12.1 | TDD ERD (new)                       | -                    |

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
| Phase 2 (Core Attendance) | 25          | 0         | 25        | ~68h       |
| Phase 3 (Self-Service)    | 13          | 0         | 13        | ~28h       |
| Phase 4 (HR Operations)   | 15          | 0         | 15        | ~38h       |
| Phase 5 (Analytics)       | 8           | 0         | 8         | ~20h       |
| **Total**                 | **81**      | **20**    | **61**    | **~202h**  |
