---
title: Multi-Session Attendance
version: 1.0.0
created: 2026-06-11
last_modified: 2026-06-11
status: draft
---

# Plan: Multi-Session Attendance (Single Button, Auto-Detect Shift)

> **Problem:** User dengan 2+ shift di hari yang sama (misal Pagi 06-12 + Malam 18-24) tidak bisa check-in kedua kalinya karena model attendance hanya support 1 session per hari. Termasuk kasus cross-day: shift malam (18-24) tanggal 11/06 lembur sampai pagi 12/06, lalu user punya shift pagi di 12/06.
>
> **Solution:** Ubah attendance jadi session-based dengan unique key `(user_id, date, shift_id)`, single button yang auto-detect shift applicable, dan window logic untuk check-in/check-out.

---

## Konsep Inti

- **1 tombol** yang berubah state: `Check In` → `Check Out` → `Selesai`
- **Backend auto-detect** shift yang applicable berdasarkan waktu sekarang
- **Tidak perlu user pilih shift** — sistem cari shift yang window-nya match
- **Setiap check-in harus ada check-out** — tidak boleh ada session menggantung
- **Cross-day aware** — session aktif dari tanggal sebelumnya block check-in baru sampai di-checkout

---

## Window Logic

### Check-in Window

```
[shiftStart - buffer] ──────────── [shiftStart + buffer] ──────────── [shiftEnd]
         ↑                                      ↑                              ↑
   Bisa check-in                         Batas "tepat waktu"            Check-in masih bisa
   (status: present)                     (setelah ini = telat)           (status: late)
```

| Waktu | Status |
|-------|--------|
| `shiftStart - buffer` → `shiftStart + buffer` | Bisa check-in, status `present` |
| `shiftStart + buffer` → `shiftEnd` | Bisa check-in, status `late` |
| Di luar window | **Tidak bisa check-in** |

**Contoh:** Shift 06:00-12:00, buffer 15 menit
- 05:45 - 06:15 → check-in, status `present`
- 06:15 - 12:00 → check-in, status `late`
- < 05:45 atau > 12:00 → **error**

### Check-out Window

```
[shiftEnd - buffer] ──────────── [shiftEnd + buffer] ──────────── [∞]
         ↑                                      ↑                  ↑
   Bisa check-out                          Batas normal       Check-out = lembur
   (durasi = normal)                       (setelah ini        (overtime dihitung
                                            = lembur)           dari titik ini)
```

| Waktu | Keterangan |
|-------|------------|
| `shiftEnd - buffer` → `shiftEnd + buffer` | Check-out normal, durasi = waktu kerja |
| > `shiftEnd + buffer` | Check-out **lembur**, `overtime_minutes` dihitung dari `shiftEnd + buffer` |
| < `shiftEnd - buffer` | **Tidak bisa check-out** (belum masuk window) |

**Contoh:** Shift 06:00-12:00, buffer 15 menit
- 11:45 - 12:15 → check-out normal
- > 12:15 → check-out lembur (overtime dihitung dari 12:15)
- < 11:45 → **error**

### Validasi Session State

| Action | Validasi | Error jika gagal |
|--------|----------|------------------|
| **Check-in** | Tidak boleh ada session aktif (belum check-out) untuk user ini **(cek semua tanggal, bukan hanya hari ini)** | "Anda memiliki sesi check-in yang belum di-checkout dari tanggal [DD/MM]. Silakan check-out terlebih dahulu." |
| **Check-out** | Harus ada session aktif (sudah check-in, belum check-out) | "Belum melakukan check-in" |

### Cross-Day Scenario (Malam 11/06 → Pagi 12/06)

> **Note:** Shift malam (18-24) tanggal 11/06 yang lembur sampai pagi tanggal 12/06, lalu user punya shift pagi (06-12) di tanggal 12/06. `date` field di attendance berdasarkan **tanggal check-in**, bukan check-out.

```
=== 11 Juni ===
18:00 check-in shift malam (date = 2026-06-11, present)
— lembur melewati tengah malam —

=== 12 Juni ===
07:00 check-out shift malam (date tetap 2026-06-11, overtime dihitung dari 00:15 = 6h 45m)
07:05 check-in shift pagi (date = 2026-06-12, status: late, karena 07:05 > 06:15)
12:00 check-out shift pagi (date = 2026-06-12, durasi normal)
```

**Validasi cross-day:**
- `FindActiveSessionByUserID` cek **semua tanggal**, bukan hanya hari ini
- Jika ada session aktif dari tanggal sebelumnya → **BLOCK** check-in baru dengan pesan: *"Anda memiliki sesi check-in yang belum di-checkout dari tanggal [DD/MM]. Silakan check-out terlebih dahulu."*
- Setelah check-out session lama → user bisa check-in untuk shift baru

**Edge case: user lupa check-out dari hari sebelumnya**
- Frontend: tampilkan alert di halaman attendance: *"Sesi dari [DD/MM] belum di-checkout. Silakan check-out untuk shift [nama shift] sebelum melanjutkan."*
- Tombol berubah ke "Check Out" untuk session yang menggantung, bukan shift hari ini
- HR bisa force close session via dashboard (future feature)

### Multi-Shift Scenario (Pagi 06-12 + Malam 18-24)

> **Note:** Shift "Siang Malam" (12-24) adalah 1 shift gabungan, bukan multi-session. Multi-session terjadi ketika user punya 2 shift terpisah di hari yang sama, contoh: Pagi (06-12) + Malam (18-24).

```
06:00 check-in shift pagi (present)
12:00 wajib check-out shift pagi (session selesai)
— break 6 jam —
18:00 check-in shift malam (present)
00:00 check-out shift malam (session selesai)

Jika lembur:
06:00 check-in shift pagi
13:00 check-out shift pagi (1 jam lembur, overtime dihitung dari 12:15)
18:00 check-in shift malam (status: late, karena 18:00 > 18:15? tidak — 18:00 masih dalam window present 17:45-18:15)
01:00 check-out shift malam (1 jam lembur, overtime dihitung dari 00:15)
```

**Aturan:** User harus check-out shift A dulu sebelum bisa check-in shift B.

---

## Perubahan yang Dibutuhkan

### Phase 1: Database + Model (2h)

| Task ID | Task | Detail | Files |
|---------|------|--------|-------|
| **MS-01** | Migration: unique index `(user_id, date, shift_id)` | Ganti unique constraint `(user_id, date)` → `(user_id, date, shift_id)` | `be/cmd/migration/` |
| **MS-02** | Model: tambah `OvertimeMinutes int` | Field baru di `Attendance` model untuk hitung menit lembur | `be/internal/models/attendance.go` |

### Phase 2: Backend — Shift Detection (2h)

| Task ID | Task | Detail | Files |
|---------|------|--------|-------|
| **MS-03** | Repository: `FindAllActiveForUserDate` | Ambil semua shift assignment aktif untuk user di tanggal tertentu | `be/internal/repositories/user_shift_assignment_repository.go` |
| **MS-04** | Service: `FindApplicableShift(userID, now)` | Loop semua shift aktif, return yang window-nya match dengan `now` (dalam `[start-buffer, end]`) | `be/internal/services/attendance_service.go` |

> **Note:** `FindNextShift` **tidak diperlukan**. `FindApplicableShift` sudah cukup — setelah user check-out shift A, call lagi `FindApplicableShift` dan akan return shift B (karena shift A sudah tidak applicable lagi untuk check-in).

### Phase 3: Backend — Check-in/Check-out Logic (4h)

| Task ID | Task | Detail | Files |
|---------|------|--------|-------|
| **MS-05** | `AttendanceCheckIn`: validasi tidak ada session aktif | Cek apakah ada attendance dengan `TimeIn != nil AND TimeOut == nil` untuk user ini **(semua tanggal, bukan hanya hari ini)**. Jika ada, return error dengan tanggal session lama | `attendance_service.go` |
| **MS-06** | `AttendanceCheckIn`: auto-detect shift via `FindApplicableShift` | Tidak perlu `shift_id` dari frontend | `attendance_service.go` |
| **MS-07** | `AttendanceCheckIn`: status logic baru | `present` jika check-in dalam `[start-buffer, start+buffer]`, `late` jika setelahnya | `attendance_service.go` |
| **MS-08** | `AttendanceCheckIn`: duplicate check per shift | Cek `(user_id, date, shift_id)` bukan `(user_id, date)` | `attendance_service.go` |
| **MS-09** | `AttendanceCheckOut`: validasi ada session aktif | Cek attendance dengan `TimeIn != nil AND TimeOut == nil` | `attendance_service.go` |
| **MS-10** | `AttendanceCheckOut`: validasi window | Tidak bisa check-out sebelum `shiftEnd - buffer` | `attendance_service.go` |
| **MS-11** | `AttendanceCheckOut`: hitung `overtime_minutes` | Jika check-out > `shiftEnd + buffer`, hitung selisih | `attendance_service.go` |
| **MS-12** | `AttendanceCheckOut`: cari session berdasarkan `(user_id, date, shift_id)` | Bukan `(user_id, date)` | `attendance_service.go` |

### Phase 4: Backend — Repository + DTO + API (3h)

| Task ID | Task | Detail | Files |
|---------|------|--------|-------|
| **MS-13** | Repository: `FindActiveSessionByUserID` | Cari session aktif (TimeIn != nil AND TimeOut == nil) untuk user **(semua tanggal, bukan hanya hari ini)**. Return session + tanggalnya untuk error message | `attendance_repository.go` |
| **MS-14** | Repository: `FindByUserDateShift` | Cari attendance berdasarkan `(user_id, date, shift_id)` | `attendance_repository.go` |
| **MS-15** | Repository: `FindByUserAndDate` → return `[]Attendance` | Ubah return type dari single ke list | `attendance_repository.go` |
| **MS-16** | DTO: tambah `OvertimeMinutes` di `AttendanceDTO` | Field baru | `attendance_dto.go` |
| **MS-17** | DTO: `AttendanceTodayResponse` — return sessions + current action | Response enriched: `{ sessions: [], current_action: "checkin"|"checkout"|"done", applicable_shift: {...} }` | `attendance_dto.go` |
| **MS-18** | Endpoint: `GET /api/attendance/today` — return semua session + action yang bisa dilakukan | Ubah existing endpoint | `attendance_route.go`, `attendance_handler.go` |

### Phase 5: Frontend Store (3h)

| Task ID | Task | Detail | Files |
|---------|------|--------|-------|
| **MS-19** | Store: ganti `checkInData` → `sessions[]` | Array of session | `fe/src/stores/attendance.ts` |
| **MS-20** | Store: tambah `currentAction` | Dari endpoint `/attendance/today`, decide button state: `{ action: "checkin"|"checkout"|"done", shift: {...} }` | `fe/src/stores/attendance.ts` |
| **MS-21** | Store: tambah `todaysShifts[]` | Daftar shift yang assigned hari ini (untuk list display) | `fe/src/stores/attendance.ts` |
| **MS-22** | Store: method `executeAction()` | Call check-in atau check-out berdasarkan `currentAction` | `fe/src/stores/attendance.ts` |
| **MS-23** | Store: method `fetchTodayStatus()` → update sessions + currentAction | Refresh state | `fe/src/stores/attendance.ts` |

### Phase 6: Frontend UI (4h)

| Task ID | Task | Detail | Files |
|---------|------|--------|-------|
| **MS-24** | Attendance page: **single button** | Text & color berubah: `Check In` (green) → `Check Out` (orange) → `Selesai` (purple, disabled) | `fe/src/pages/attendance/attendance/IndexView.vue` |
| **MS-25** | Button disabled state | Disabled jika: tidak ada shift applicable, user di luar radius, atau ada session aktif tapi bukan milik shift yang applicable | `IndexView.vue` |
| **MS-26** | **Shift list** — tampilkan daftar shift hari ini | Card per shift: nama, waktu, status (Belum / Checked-in / Selesai), durasi, overtime badge | `IndexView.vue` |
| **MS-27** | Session info card | Show shift name, window time, status (present/late), overtime badge | `IndexView.vue` |
| **MS-28** | Overtime badge visual | Tampilkan jika check-out > `shiftEnd + buffer` | `IndexView.vue` |

### Phase 7: Testing (4h)

| Task ID | Task | Detail |
|---------|------|--------|
| **MS-29** | Test: 1 shift normal (06-12, check-in 06:00, check-out 12:00) | Regression — status present, durasi 6h |
| **MS-30** | Test: 1 shift telat (check-in 06:30 → status `late`) | Status logic |
| **MS-31** | Test: 1 shift lembur (check-out 13:00 → overtime 45 menit) | Overtime calc (dari 12:15) |
| **MS-32** | Test: 2 shift berurutan (06-12 + 18-24) — check-in/out keduanya | Multi-session |
| **MS-33** | Test: 2 shift + lembur (06-12 check-out 13:00, 18-24 check-out 01:00) | Multi-session + overtime |
| **MS-34** | Test: check-out sebelum window (11:30 → error) | Validation |
| **MS-35** | Test: check-in di luar window (05:30 → error) | Validation |
| **MS-36** | Test: check-in saat ada session aktif → error | Session state validation |
| **MS-37** | Test: check-out tanpa check-in → error | Session state validation |
| **MS-38** | Test: cross-day — check-in malam 11/06, check-out 07:00 12/06, check-in pagi 12/06 | Cross-day session + overtime |
| **MS-39** | Test: cross-day + lupa check-out → block check-in baru, show alert | Session state validation |

---

## Estimasi Total

| Phase | Tasks | Effort |
|-------|-------|--------|
| Database + Model | MS-01, MS-02 | 2h |
| Backend — Shift Detection | MS-03, MS-04 | 2h |
| Backend — Check-in/Check-out Logic | MS-05 s/d MS-12 | 4h |
| Backend — Repository + DTO + API | MS-13 s/d MS-18 | 3h |
| Frontend Store | MS-19 s/d MS-23 | 3h |
| Frontend UI | MS-24 s/d MS-28 | 4h |
| Testing | MS-29 s/d MS-39 | 5h |
| **Total** | **39 tasks** | **~23h** |

---

## UX Wireframe

```
┌─────────────────────────────────────────────────┐
│  Absensi Hari Ini — Rabu, 11 Jun 2026           │
├─────────────────────────────────────────────────┤
│                                                 │
│  📍 Kantor: Head Office (150m)                  │
│  📏 Jarak: 45m — Dalam area ✅                  │
│                                                 │
├─────────────────────────────────────────────────┤
│                                                 │
│        ┌─────────────────────────────┐          │
│        │                             │          │
│        │      [  Check Out  ]        │          │
│        │                             │          │
│        │   Shift Pagi • 06:00-12:00  │          │
│        │   Checked-in: 06:05 (present)          │
│        │   Durasi: 5h 50m            │          │
│        │                             │          │
│        └─────────────────────────────┘          │
│                                                 │
├─────────────────────────────────────────────────┤
│  Shift Hari Ini                                 │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │ ☀️  Shift Pagi    06:00 - 12:00           │  │
│  │     🔄 Checked-in (Active)                │  │
│  │     Check-in:  06:05 (present)            │  │
│  │     Durasi: 5h 50m ...                    │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │ 🌙  Shift Malam   18:00 - 24:00           │  │
│  │     ⏳ Belum Mulai                        │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
├─────────────────────────────────────────────────┤
│  Total Hari Ini: 5h 50m (1 session active)      │
└─────────────────────────────────────────────────┘
```

---

## Backward Compatibility

- User dengan **1 shift per hari** → UI tetap sama, tidak ada perubahan UX
- Endpoint `/attendance/today` tetap ada, response enriched dengan sessions array
- Migration tidak menghapus data existing, hanya tambah unique constraint
- Data existing otomatis dapat `shift_id` dari record yang sudah ada

---

## Risk & Mitigation

| Risk | Mitigation |
|------|------------|
| Data existing conflict dengan unique constraint baru | Migration script handle existing data — set default shift_id dari record yang ada |
| User check-out lupa, session menggantung | Frontend tampilkan alert + tombol check-out untuk session lama; Admin bisa force close via HR dashboard (future feature) |
| Shift overlap (12-18 dan 12-20) | `FindApplicableShift` ambil shift dengan `endTime` terdekat dari `now` |
| Cross-day lembur (malam → pagi besok) | `FindActiveSessionByUserID` cek semua tanggal, block check-in baru sampai session lama di-checkout |
