-- Dummy data for portfolio / screenshot purposes.
-- Idempotent (safe to re-run), does not touch the app's own seed data.
-- Password for all employee accounts below: Karyawan#123
--
-- Embedded via go:embed (see dummy_seeder.go) and run with:
--   go run cmd/migration/main.go dummy
--
-- Requires tables already created via `go run cmd/migration/main.go up`.

-- ---------------------------------------------------------------------------
-- 1. Office location
-- ---------------------------------------------------------------------------
INSERT INTO office_locations (name, address, latitude, longitude, radius_meters, is_active, created_at, updated_at)
SELECT 'Kantor Pusat', 'Jl. Jend. Sudirman No. 1, Jakarta Pusat', -6.181788, 106.820668, 120, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM office_locations WHERE name = 'Kantor Pusat');

-- ---------------------------------------------------------------------------
-- 2. Shifts
-- ---------------------------------------------------------------------------
INSERT INTO shifts (name, start_time, end_time, break_duration, flexi_minutes, color_code, total_hours, created_at, updated_at)
SELECT 'Pagi', '08:00', '17:00', 60, 15, '#3B82F6', 8.0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM shifts WHERE name = 'Pagi');

INSERT INTO shifts (name, start_time, end_time, break_duration, flexi_minutes, color_code, total_hours, created_at, updated_at)
SELECT 'Siang', '14:00', '22:00', 60, 15, '#F59E0B', 8.0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM shifts WHERE name = 'Siang');

INSERT INTO shifts (name, start_time, end_time, break_duration, flexi_minutes, color_code, total_hours, created_at, updated_at)
SELECT 'Malam', '22:00', '06:00', 60, 15, '#8B5CF6', 8.0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM shifts WHERE name = 'Malam');

-- ---------------------------------------------------------------------------
-- 3. Roles
-- ---------------------------------------------------------------------------
INSERT INTO roles (name, description, created_at, updated_at)
SELECT 'Karyawan', 'Karyawan: hanya bisa lihat shift, riwayat shift, dan melakukan check-in/check-out', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM roles WHERE name = 'Karyawan');

-- ---------------------------------------------------------------------------
-- 4. Employee users + profiles
--    bcrypt hash below = "Karyawan#123"
-- ---------------------------------------------------------------------------
INSERT INTO users (email, password, name, avatar, created_at, updated_at)
SELECT v.email, v.password, v.name, '', NOW(), NOW()
FROM (
    SELECT 'budi.santoso@app.com'   AS email, '$2a$10$RszEG7DvvugwyKnfGL7kX.FpCdZijRvXMJqIW7g0bUXP93LnK9ZL.' AS password, 'Budi Santoso'     AS name
    UNION ALL SELECT 'siti.nurhaliza@app.com', '$2a$10$RszEG7DvvugwyKnfGL7kX.FpCdZijRvXMJqIW7g0bUXP93LnK9ZL.', 'Siti Nurhaliza'
    UNION ALL SELECT 'andi.wijaya@app.com',    '$2a$10$RszEG7DvvugwyKnfGL7kX.FpCdZijRvXMJqIW7g0bUXP93LnK9ZL.', 'Andi Wijaya'
    UNION ALL SELECT 'dewi.lestari@app.com',   '$2a$10$RszEG7DvvugwyKnfGL7kX.FpCdZijRvXMJqIW7g0bUXP93LnK9ZL.', 'Dewi Lestari'
    UNION ALL SELECT 'rizky.pratama@app.com',  '$2a$10$RszEG7DvvugwyKnfGL7kX.FpCdZijRvXMJqIW7g0bUXP93LnK9ZL.', 'Rizky Pratama'
    UNION ALL SELECT 'putri.ayu@app.com',      '$2a$10$RszEG7DvvugwyKnfGL7kX.FpCdZijRvXMJqIW7g0bUXP93LnK9ZL.', 'Putri Ayu'
    UNION ALL SELECT 'maya.kusuma@app.com',    '$2a$10$RszEG7DvvugwyKnfGL7kX.FpCdZijRvXMJqIW7g0bUXP93LnK9ZL.', 'Maya Kusuma'
) v
WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.email = v.email);

INSERT INTO user_profiles (user_id, department, position, phone, join_date, created_at, updated_at)
SELECT u.id, v.department, v.position, v.phone, v.join_date, NOW(), NOW()
FROM (
    SELECT 'budi.santoso@app.com'    AS email, 'Engineering'      AS department, 'Backend Developer'   AS position, '081234500001' AS phone, DATE_SUB(CURDATE(), INTERVAL 2 YEAR)   AS join_date
    UNION ALL SELECT 'siti.nurhaliza@app.com', 'Engineering',       'Frontend Developer',  '081234500002', DATE_SUB(CURDATE(), INTERVAL 18 MONTH)
    UNION ALL SELECT 'andi.wijaya@app.com',    'Finance',           'Finance Staff',        '081234500003', DATE_SUB(CURDATE(), INTERVAL 3 YEAR)
    UNION ALL SELECT 'dewi.lestari@app.com',   'Marketing',         'Marketing Executive',  '081234500004', DATE_SUB(CURDATE(), INTERVAL 1 YEAR)
    UNION ALL SELECT 'rizky.pratama@app.com',  'Engineering',       'QA Engineer',          '081234500005', DATE_SUB(CURDATE(), INTERVAL 8 MONTH)
    UNION ALL SELECT 'putri.ayu@app.com',      'Human Resources',   'HR Staff',             '081234500006', DATE_SUB(CURDATE(), INTERVAL 28 MONTH)
    UNION ALL SELECT 'maya.kusuma@app.com',    'Marketing',         'Content Specialist',   '081234500008', DATE_SUB(CURDATE(), INTERVAL 5 MONTH)
) v
JOIN users u ON u.email = v.email
WHERE NOT EXISTS (SELECT 1 FROM user_profiles p WHERE p.user_id = u.id);

-- ---------------------------------------------------------------------------
-- 5. Assign "Karyawan" role to the employee users
-- ---------------------------------------------------------------------------
INSERT INTO user_has_roles (user_id, role_id, created_at)
SELECT u.id, r.id, NOW()
FROM users u
JOIN roles r ON r.name = 'Karyawan'
WHERE u.email IN (
    'budi.santoso@app.com', 'siti.nurhaliza@app.com', 'andi.wijaya@app.com', 'dewi.lestari@app.com',
    'rizky.pratama@app.com', 'putri.ayu@app.com', 'maya.kusuma@app.com'
)
AND NOT EXISTS (SELECT 1 FROM user_has_roles ur WHERE ur.user_id = u.id AND ur.role_id = r.id);

-- ---------------------------------------------------------------------------
-- 6. Assign "Pagi" shift to the employee users (current month -> +2 years)
-- ---------------------------------------------------------------------------
INSERT INTO user_shift_assignments (user_id, shift_id, start_date, end_date, is_active, created_at, updated_at)
SELECT u.id, s.id,
       DATE_FORMAT(CURDATE(), '%Y-%m-01'),
       DATE_ADD(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 2 YEAR),
       1, NOW(), NOW()
FROM users u
JOIN shifts s ON s.name = 'Pagi'
WHERE u.email IN (
    'budi.santoso@app.com', 'siti.nurhaliza@app.com', 'andi.wijaya@app.com', 'dewi.lestari@app.com',
    'rizky.pratama@app.com', 'putri.ayu@app.com', 'maya.kusuma@app.com'
)
AND NOT EXISTS (
    SELECT 1 FROM user_shift_assignments a
    WHERE a.user_id = u.id AND a.shift_id = s.id AND a.start_date = DATE_FORMAT(CURDATE(), '%Y-%m-01')
);

-- Existing assignments (from a previous run with the old +1 year window)
-- get extended too, so re-running this file after the window changes keeps
-- everyone's schedule pushed out to +2 years instead of silently stopping.
UPDATE user_shift_assignments a
JOIN users u ON u.id = a.user_id
JOIN shifts s ON s.id = a.shift_id AND s.name = 'Pagi'
SET a.end_date = DATE_ADD(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 2 YEAR)
WHERE u.email IN (
    'budi.santoso@app.com', 'siti.nurhaliza@app.com', 'andi.wijaya@app.com', 'dewi.lestari@app.com',
    'rizky.pratama@app.com', 'putri.ayu@app.com', 'maya.kusuma@app.com'
)
AND a.end_date < DATE_ADD(DATE_FORMAT(CURDATE(), '%Y-%m-01'), INTERVAL 2 YEAR);

-- ---------------------------------------------------------------------------
-- 7. Attendance history: today + last 29 days, weekdays only, per employee.
--    Deterministic pattern from (user.id, day offset) so re-runs are stable:
--      bucket 0-6 (70%) -> present, on-time-ish check-in
--      bucket 7-8 (20%) -> late, check-in past the flexi window
--      bucket 9   (10%) -> absent, no row inserted
--    Today (n=0) only gets a check-in (time_out left NULL) — the HR dashboard
--    ("today" widgets, recent activity, not-yet-checked-in) reads TODAY's
--    attendance rows only, so it needs a live-looking in-progress record, not
--    a future checkout timestamp.
-- ---------------------------------------------------------------------------
INSERT IGNORE INTO attendances (
    user_id, shift_id, date, office_id, status, status_out,
    time_in, lat_in, lng_in, distance_meters,
    time_out, lat_out, lng_out,
    duration, duration_minutes, overtime_minutes,
    created_at, updated_at
)
SELECT
    u.id,
    s.id,
    cal.att_date,
    o.id,
    CASE WHEN cal.bucket <= 6 THEN 'present' ELSE 'late' END,
    CASE
        WHEN cal.time_out IS NULL THEN NULL
        WHEN cal.co_kind = 'early' THEN 'early_leave'
        ELSE 'on_time'
    END,
    cal.time_in,
    o.latitude  + (CAST(cal.jitter_in AS SIGNED)  - 3) * 0.00003,
    o.longitude + (CAST(cal.jitter_in AS SIGNED)  - 3) * 0.00003,
    ABS(CAST(cal.jitter_in AS SIGNED) - 3) * 3.3,
    cal.time_out,
    o.latitude  + (CAST(cal.jitter_out AS SIGNED) - 3) * 0.00003,
    o.longitude + (CAST(cal.jitter_out AS SIGNED) - 3) * 0.00003,
    CONCAT(FLOOR(TIMESTAMPDIFF(MINUTE, cal.time_in, cal.time_out) / 60), 'h ',
           MOD(TIMESTAMPDIFF(MINUTE, cal.time_in, cal.time_out), 60), 'm'),
    COALESCE(TIMESTAMPDIFF(MINUTE, cal.time_in, cal.time_out), 0),
    COALESCE(GREATEST(0, TIMESTAMPDIFF(MINUTE, TIMESTAMP(cal.att_date, '17:00:00'), cal.time_out)), 0),
    NOW(), NOW()
FROM users u
JOIN shifts s ON s.name = 'Pagi'
JOIN office_locations o ON o.name = 'Kantor Pusat'
JOIN (
    -- pre-compute everything per (user, day) once, then just SELECT from it above
    SELECT
        u5.id AS uid,
        days5.n,
        DATE_SUB(CURDATE(), INTERVAL days5.n DAY) AS att_date,
        MOD(u5.id * 13 + days5.n * 7, 10) AS bucket,
        MOD(u5.id + days5.n, 7) AS jitter_in,
        MOD(u5.id * 5 + days5.n, 7) AS jitter_out,
        CASE
            WHEN MOD(u5.id * 13 + days5.n * 7, 10) <= 6
                THEN TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL days5.n DAY), '08:00:00') - INTERVAL MOD(u5.id + days5.n, 10) MINUTE
            ELSE TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL days5.n DAY), '08:00:00') + INTERVAL (20 + MOD(u5.id + days5.n, 30)) MINUTE
        END AS time_in,
        CASE
            WHEN MOD(u5.id + days5.n, 8) = 0 THEN 'early'
            ELSE 'normal'
        END AS co_kind,
        CASE
            WHEN days5.n = 0 THEN NULL                                              -- still clocked in today
            WHEN MOD(u5.id + days5.n, 8) = 0
                THEN TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL days5.n DAY), '17:00:00') - INTERVAL 20 MINUTE
            WHEN MOD(u5.id + days5.n, 6) = 0
                THEN TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL days5.n DAY), '17:00:00') + INTERVAL (15 + MOD(u5.id * 3 + days5.n, 60)) MINUTE
            ELSE TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL days5.n DAY), '17:00:00')
        END AS time_out
    FROM users u5
    JOIN (
        SELECT 0 AS n UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5
        UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9 UNION ALL SELECT 10
        UNION ALL SELECT 11 UNION ALL SELECT 12 UNION ALL SELECT 13 UNION ALL SELECT 14 UNION ALL SELECT 15
        UNION ALL SELECT 16 UNION ALL SELECT 17 UNION ALL SELECT 18 UNION ALL SELECT 19 UNION ALL SELECT 20
        UNION ALL SELECT 21 UNION ALL SELECT 22 UNION ALL SELECT 23 UNION ALL SELECT 24 UNION ALL SELECT 25
        UNION ALL SELECT 26 UNION ALL SELECT 27 UNION ALL SELECT 28 UNION ALL SELECT 29
    ) days5 ON 1 = 1
    WHERE u5.email IN (
        'budi.santoso@app.com', 'siti.nurhaliza@app.com', 'andi.wijaya@app.com', 'dewi.lestari@app.com',
        'rizky.pratama@app.com', 'putri.ayu@app.com', 'maya.kusuma@app.com'
    )
) cal ON cal.uid = u.id
WHERE cal.bucket <> 9                          -- ~10% absent: skip row entirely
AND DAYOFWEEK(cal.att_date) NOT IN (1, 7)       -- skip Sunday(1) / Saturday(7)
-- The app's (user_id, date, shift_id) unique index isn't guaranteed to exist
-- (migrations only AutoMigrate the struct, not the model's custom GORMIndexes()),
-- so INSERT IGNORE alone can't be trusted to dedupe on re-run. Guard explicitly.
AND NOT EXISTS (
    SELECT 1 FROM attendances ax
    WHERE ax.user_id = u.id AND ax.shift_id = s.id AND ax.date = cal.att_date
);
