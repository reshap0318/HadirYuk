-- AI Customer Service Chatbot — readonly MySQL user (plan §6.1)
-- Ganti hadir_yuk dan alpha0318 sesuai server, lalu jalankan sebagai root/admin.
--
-- Catatan: sengaja grant per-tabel (bukan `GRANT SELECT ON hadir_yuk.*`), karena
-- MySQL gak bisa REVOKE per-tabel di atas grant level database — begitu SELECT
-- di-grant db-wide, gak bisa dipersempit lagi khusus buat tabel `users`.

CREATE USER IF NOT EXISTS 'hadiryuk_ai_ro'@'%' IDENTIFIED BY 'alpha0318';

GRANT SELECT ON hadir_yuk.attendances             TO 'hadiryuk_ai_ro'@'%';
GRANT SELECT ON hadir_yuk.notifications           TO 'hadiryuk_ai_ro'@'%';
GRANT SELECT ON hadir_yuk.office_locations        TO 'hadiryuk_ai_ro'@'%';
GRANT SELECT ON hadir_yuk.qr_codes                TO 'hadiryuk_ai_ro'@'%';
GRANT SELECT ON hadir_yuk.shifts                  TO 'hadiryuk_ai_ro'@'%';
GRANT SELECT ON hadir_yuk.user_profiles           TO 'hadiryuk_ai_ro'@'%';
GRANT SELECT ON hadir_yuk.user_shift_assignments  TO 'hadiryuk_ai_ro'@'%';

-- users: HANYA id + name (password/email/avatar TIDAK BOLEH keluar lewat AI chat, §6.2)
GRANT SELECT (id, name) ON hadir_yuk.users TO 'hadiryuk_ai_ro'@'%';

-- password_resets, roles, permissions, role_has_permissions, user_has_roles:
-- sengaja TIDAK di-grant sama sekali (full blacklist, §6.2).

FLUSH PRIVILEGES;
