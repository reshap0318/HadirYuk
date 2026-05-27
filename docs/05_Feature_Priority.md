# Feature Priority Plan — HadirYuk

> **Status:** Post Data Master (Shifts, Office Locations, Leave Types CRUD complete)
> **Last Updated:** 2026-05-27
> **Tech Stack:** Go (Gin) Backend + Vue 3 Frontend + MySQL

---

## Completed (Foundation)

| Area | Features | Status |
|------|----------|--------|
| **Auth** | Login (JWT RS256), Logout, Token Refresh | Done |
| **UAM** | Role CRUD, Permission CRUD, Role-Permission Mapping, User-Role Assignment | Done |
| **Data Master** | Shift CRUD, Office Location CRUD, Leave Type CRUD | Done |

---

## Priority Tiers

### P0 — Critical Path (Core Attendance)

> Without these, the system cannot fulfill its primary purpose: recording attendance.

| # | Feature | Complexity | Dependencies | Primary Role | Effort |
|---|---------|------------|--------------|--------------|--------|
| P0-1 | **Employee-Shift Assignment** | M | Shift CRUD (done), User CRUD (done) | HR Admin | 2 days |
| P0-2 | **Profile Management + Face Photo Upload** | M | Auth (done), User CRUD (done) | Karyawan | 3 days |
| P0-3 | **Attendance Check-in (Geotagging + Face Recognition)** | XL | P0-1, P0-2, Office Location (done) | Karyawan | 5 days |
| P0-4 | **Attendance Check-out (Geotagging + Face Recognition)** | L | P0-3 | Karyawan | 3 days |
| P0-5 | **QR Code Generation & Management** | M | Office Location (done) | HR Admin | 2 days |
| P0-6 | **Attendance Check-in (QR Code)** | L | P0-5 | Karyawan | 3 days |
| P0-7 | **Attendance Check-out (QR Code)** | S | P0-6 | Karyawan | 1 day |

**P0 Total Effort: ~19 days**

**Why this order:**
1. Employees must be assigned to a shift before any attendance can be validated (shift defines work hours, tolerance, break duration).
2. Face photo must exist before face recognition can work during check-in.
3. Geotagging check-in is the primary attendance method — build it first.
4. Check-out builds directly on check-in logic (same record, different state).
5. QR code must be generated before it can be scanned.
6. QR check-in/out is a secondary method — simpler than geotagging (no face recognition).

---

### P1 — Important (Employee Self-Service)

> Features that employees need daily to interact with the system meaningfully.

| # | Feature | Complexity | Dependencies | Primary Role | Effort |
|---|---------|------------|--------------|--------------|--------|
| P1-1 | **Attendance History & Detail** | M | P0-3, P0-4 | Karyawan | 2 days |
| P1-2 | **Employee Dashboard** | M | P0-1, P0-3, P0-4 | Karyawan | 3 days |
| P1-3 | **Leave Request Submission** | M | Leave Type (done) | Karyawan | 2 days |
| P1-4 | **Leave Balance View** | S | P1-3 | Karyawan | 1 day |
| P1-5 | **Password Change** | S | Auth (done) | All Roles | 1 day |
| P1-6 | **Forgot/Reset Password** | M | Auth (done), Email service | All Roles | 2 days |

**P1 Total Effort: ~11 days**

**Why this order:**
1. After employees can check in/out, they need to see their attendance records.
2. Dashboard aggregates attendance + shift data into a single view.
3. Leave request is a core self-service feature independent of attendance.
4. Leave balance depends on leave request logic (calculating used vs total).
5. Password management is low effort but high impact for user experience.

---

### P2 — Should Have (HR Operational Features)

> Features that HR Admin needs to manage the system and handle exceptions.

| # | Feature | Complexity | Dependencies | Primary Role | Effort |
|---|---------|------------|--------------|--------------|--------|
| P2-1 | **HR Dashboard** | L | P0-3, P0-4, P1-1, P1-3 | HR Admin | 3 days |
| P2-2 | **Attendance Correction** | M | P1-1 | HR Admin | 2 days |
| P2-3 | **Shift Schedule View** | M | P0-1 | Karyawan, HR Admin | 2 days |
| P2-4 | **Attendance Report (View + Export Excel/PDF)** | L | P1-1 | HR Admin | 4 days |
| P2-5 | **Leave Report** | M | P1-3, P1-4 | HR Admin | 2 days |

**P2 Total Effort: ~13 days**

**Why this order:**
1. HR Dashboard is the command center — needs attendance + leave data to be meaningful.
2. Attendance correction handles edge cases (forgot check-out, wrong location).
3. Shift schedule view helps both employees and HR visualize assignments.
4. Reports are built on top of attendance history data.
5. Leave report depends on leave request/balance being functional.

---

### P3 — Nice to Have (Analytics & Reporting)

> Features that provide insights but are not required for day-to-day operations.

| # | Feature | Complexity | Dependencies | Primary Role | Effort |
|---|---------|------------|--------------|--------------|--------|
| P3-1 | **Late Statistics** | M | P1-1, P0-1 | HR Admin | 2 days |
| P3-2 | **Admin Dashboard** | S | UAM (done) | Super Admin | 1 day |

**P3 Total Effort: ~3 days**

**Why this order:**
1. Late statistics require attendance history to be fully functional.
2. Admin Dashboard is lightweight — mostly system stats and recent activity.

---

### P4 — Could Have (Audit & Advanced)

> Features that add polish and compliance but can be deferred.

| # | Feature | Complexity | Dependencies | Primary Role | Effort |
|---|---------|------------|--------------|--------------|--------|
| P4-1 | **Audit Log** | L | All CRUD operations | Super Admin | 3 days |

**P4 Total Effort: ~3 days**

**Why this order:**
1. Audit log is a cross-cutting concern — best implemented after all CRUD operations are stable so every action can be logged consistently.

---

## Implementation Roadmap (Phased)

```
Phase 1 — P0: Core Attendance Engine     (~19 days)
  Week 1-2: Employee-Shift, Profile/Face Photo, Geotagging Check-in/out
  Week 3:   QR Code Generation, QR Check-in/out

Phase 2 — P1: Employee Self-Service       (~11 days)
  Week 4:   Attendance History, Employee Dashboard
  Week 5:   Leave Request, Leave Balance, Password Management

Phase 3 — P2: HR Operations               (~13 days)
  Week 6-7: HR Dashboard, Attendance Correction, Shift Schedule
  Week 8:   Attendance Report, Leave Report

Phase 4 — P3: Analytics                   (~3 days)
  Week 9:   Late Statistics, Admin Dashboard

Phase 5 — P4: Audit                       (~3 days)
  Week 10:  Audit Log implementation across all modules
```

**Total Estimated Effort: ~49 working days**

---

## Dependency Graph

```
Data Master (done) ──┬── P0-1 Employee-Shift Assignment ──┬── P0-3 Check-in (Geotagging)
                      │                                     ├── P0-4 Check-out (Geotagging)
                      │                                     ├── P1-2 Employee Dashboard
                      │                                     └── P2-3 Shift Schedule View
                      │
                      ├── P0-2 Profile + Face Photo ──────── P0-3 Check-in (Geotagging)
                      │
                      ├── P0-5 QR Code Generation ──────────┬── P0-6 Check-in (QR)
                      │                                      └── P0-7 Check-out (QR)
                      │
                      └── Leave Type (done) ──────────────── P1-3 Leave Request ── P1-4 Leave Balance
                                                                                 └── P2-5 Leave Report

P0-3 + P0-4 ────────────────────────────┬── P1-1 Attendance History ──┬── P2-2 Attendance Correction
                                        │                             ├── P2-1 HR Dashboard
                                        │                             ├── P2-4 Attendance Report
                                        │                             └── P3-1 Late Statistics
                                        │
                                        └── P1-2 Employee Dashboard

Auth (done) ────────────────────────────┬── P1-5 Password Change
                                        └── P1-6 Forgot/Reset Password

UAM (done) ───────────────────────────── P3-2 Admin Dashboard

All CRUD operations ──────────────────── P4-1 Audit Log
```

---

## Risk & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| Face recognition library integration | High | Prototype face recognition early in P0-2; have fallback to photo-only check-in |
| GPS accuracy on mobile browsers | Medium | Implement haversine distance calculation with configurable tolerance; log raw coordinates for debugging |
| QR Code expiry timing | Low | Use server-side time validation; add 30-second grace period |
| Leave balance calculation accuracy | Medium | Write unit tests for business day calculation (exclude weekends) |
| Export file generation performance | Low | Stream large datasets; add pagination for reports with >1000 rows |

---

## Notes

- **Effort estimates** assume a single full-stack developer working on backend (Go/Gin) and frontend (Vue 3) simultaneously.
- **Parallel work:** Frontend and backend for each feature can be developed in parallel once API contracts are defined.
- **Testing:** Each tier should include unit tests + integration tests before moving to the next tier.
- **Milestone:** After P0 is complete, the system can record attendance — this is the minimum viable product (MVP).
