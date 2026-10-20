# HadirYuk

Sistem absensi karyawan berbasis web dengan validasi **geotagging**, **bukti foto**, dan **QR Code**. Dibangun dengan Go (Gin) di backend dan Vue 3 + TypeScript di frontend.

HadirYuk membantu perusahaan mengelola kehadiran karyawan secara lebih akurat: check-in/check-out divalidasi berdasarkan radius lokasi kantor, dilengkapi bukti foto, riwayat absensi, manajemen shift, laporan, serta pengelolaan user dan hak akses.

> Status: aplikasi web (bukan native mobile), responsive, dan mendukung PWA.

## Daftar Isi

- [Tentang Aplikasi](#tentang-aplikasi)
- [Fitur Utama](#fitur-utama)
- [Peran Pengguna](#peran-pengguna)
- [Tech Stack](#tech-stack)
- [Arsitektur](#arsitektur)
- [Struktur Repository](#struktur-repository)
- [Prasyarat](#prasyarat)
- [Instalasi Cepat dengan Docker](#instalasi-cepat-dengan-docker)
- [Development Lokal](#development-lokal)
- [Environment Variables](#environment-variables)
- [Database & Migration](#database--migration)
- [API](#api)
- [Cron Jobs](#cron-jobs)
- [Deployment & Maintenance](#deployment--maintenance)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Dokumentasi](#dokumentasi)

## Tentang Aplikasi

HadirYuk adalah aplikasi absensi untuk karyawan kantor. Setiap absensi dapat divalidasi dengan:

- **Geotagging** — koordinat user dibandingkan dengan lokasi kantor menggunakan formula Haversine dan radius yang dapat dikonfigurasi.
- **Photo evidence** — foto check-in/check-out disimpan sebagai bukti.
- **QR Code** — alternatif absensi dengan scan QR Code yang berlaku pada lokasi tertentu.
- **Shift management** — jadwal shift, assignment ke karyawan, dan tampilan jadwal.
- **Reporting** — rekap absensi, statistik keterlambatan, dan export Excel.

Aplikasi juga memiliki dashboard berbeda untuk Karyawan, HR Admin, dan Super Admin, serta fitur opsional AI Customer Service berbasis Text-to-SQL.

## Fitur Utama

### Absensi

- Check-in dan check-out dengan geotagging + foto.
- Check-in dan check-out dengan QR Code.
- Validasi radius lokasi kantor di sisi server.
- Status absensi hari ini, riwayat absensi, dan statistik bulanan.
- Koreksi absensi oleh HR Admin.
- Laporan absensi dan export Excel.
- Statistik keterlambatan.

### Manajemen

- CRUD user, profil karyawan, dan foto profil.
- CRUD role dan permission (RBAC).
- CRUD shift dan assignment shift ke karyawan.
- Jadwal shift karyawan.
- CRUD lokasi kantor beserta latitude, longitude, dan radius.
- Generate, lihat, dan revoke QR Code.
- Notifikasi in-app.

### Dashboard

- Dashboard karyawan: status hari ini, jadwal, dan ringkasan kehadiran.
- Dashboard HR: monitoring kehadiran, keterlambatan, dan rekap.
- Super Admin: akses penuh ke manajemen user, role, permission, serta seluruh dashboard.

### Fitur Opsional

- **AI Customer Service Chatbot** — Text-to-SQL ke database read-only, kompatibel dengan OpenAI, DeepSeek, Ollama, atau API OpenAI-compatible lainnya.
- **Face recognition** — opsional, memakai build tag `gocv` dan OpenCV. Secara default aplikasi memakai stub sehingga tidak memerlukan OpenCV.
- **PWA** — frontend dapat dipasang sebagai aplikasi web.

## Peran Pengguna

| Role | Akses utama |
| --- | --- |
| **Super Admin** | Semua fitur, termasuk manajemen user, role, dan permission. |
| **HR Admin** | Karyawan, shift, assignment, absensi, koreksi absensi, laporan, lokasi, QR Code, statistik keterlambatan, dan AI chat. |
| **Karyawan** | Check-in/check-out, riwayat absensi, jadwal shift, profil, dan notifikasi. |

## Tech Stack

| Bagian | Teknologi |
| --- | --- |
| Backend | Go 1.25, Gin, GORM |
| Database | MySQL 8.0 |
| Cache | Redis 7 |
| Autentikasi | JWT RS256, bcrypt |
| Validasi | go-playground/validator |
| Export | Excelize (Excel) |
| AI Chat | OpenAI-compatible API (OpenAI, DeepSeek, Ollama, dll.) |
| Face Recognition | GoCV/OpenCV (opsional) |
| Frontend | Vue 3, TypeScript, Vite 6 |
| State management | Pinia |
| Routing | Vue Router |
| UI | Tailwind CSS 4, komponen custom |
| HTTP client | Axios |
| Form validation | Vuelidate |
| Peta | Leaflet |
| QR | html5-qrcode, qrcode |
| PWA | vite-plugin-pwa |
| Container | Docker, Docker Compose |
| Web server | Nginx |

## Arsitektur

```text
┌─────────────────────────────────────┐
│ Browser — Vue 3 SPA / PWA           │
│ http://localhost:9010 (Docker)      │
│ http://localhost:5173 (dev)         │
└──────────────────┬──────────────────┘
                   │ HTTP
┌──────────────────▼──────────────────┐
│ Nginx — frontend container          │
│ /be/*      → backend:8080           │
│ /storage/* → backend:8080           │
└──────────────────┬──────────────────┘
                   │
┌──────────────────▼──────────────────┐
│ Go API — Gin :8080                  │
│ JWT · RBAC · Rate limit · Upload    │
└─────────┬─────────────────┬─────────┘
          │                 │
   ┌──────▼──────┐   ┌──────▼──────┐
   │ MySQL 8     │   │ Redis 7     │
   │ data utama  │   │ cache       │
   └─────────────┘   └─────────────┘
          │
   ┌──────▼──────────────────────────┐
   │ Local storage — be/storage/     │
   │ avatar, evidence, foto, log, key│
   └─────────────────────────────────┘
```

Pola backend:

```text
Routes → Handlers → Services → Repositories → MySQL/Redis
              ↑
        Middleware (JWT, RBAC, CORS, Rate Limit)
```

Dependency injection ada di `be/internal/di/container.go`.

## Struktur Repository

```text
.
├── be/                              # Backend Go
│   ├── cmd/
│   │   ├── api/                     # Entry point REST API
│   │   ├── migration/               # Migration + seeder
│   │   ├── genkey/                  # Generator RSA key JWT
│   │   ├── markabsent/              # Job penandaan alpha
│   │   ├── cleartmp/                # Job pembersih file tmp
│   │   └── dummy/                   # Generator absensi demo
│   ├── internal/
│   │   ├── clients/                 # Email, face recognition, OpenAI
│   │   ├── database/                # Koneksi MySQL/PostgreSQL/Redis
│   │   ├── di/                      # Dependency injection container
│   │   ├── dtos/                    # Request/response DTO
│   │   ├── handlers/                # HTTP handler
│   │   ├── helpers/                 # JWT, response, logger, dsb.
│   │   ├── middleware/              # JWT, RBAC, CORS, rate limit
│   │   ├── models/                  # Model GORM
│   │   ├── repositories/            # Akses database
│   │   ├── routes/                  # Definisi route
│   │   └── services/                # Business logic
│   ├── storage/                     # Runtime file (bind mount)
│   │   ├── attendance-evidence/     # Foto bukti absensi
│   │   ├── avatars/                 # Foto profil
│   │   ├── face-photos/             # Foto wajah (opsional)
│   │   ├── keys/                    # Private/public key JWT
│   │   ├── logs/                    # Log aplikasi
│   │   └── tmp/                     # Upload sementara
│   ├── .air.toml                    # Config hot reload (Air)
│   ├── .env.example                 # Env untuk backend lokal
│   ├── go.mod
│   └── go.sum
│
├── fe/                              # Frontend Vue 3
│   ├── src/
│   │   ├── assets/                  # Aset statis
│   │   ├── components/              # Komponen reusable
│   │   │   ├── ai-chat/
│   │   │   ├── directives/
│   │   │   ├── layouts/
│   │   │   └── utils/
│   │   ├── composables/             # useCrud, useCamera, dsb.
│   │   ├── helpers/                 # Helper tanggal, upload, dll.
│   │   ├── layouts/                 # Layout utama
│   │   ├── pages/                   # Halaman aplikasi
│   │   │   ├── auth/
│   │   │   ├── attendance/
│   │   │   ├── dashboard/
│   │   │   ├── master/
│   │   │   ├── profile/
│   │   │   ├── uam/
│   │   │   └── users/
│   │   ├── plugins/                 # Axios, SweetAlert
│   │   ├── router/                  # Vue Router + guard
│   │   ├── stores/                  # Pinia stores
│   │   ├── App.vue
│   │   └── main.ts
│   ├── .env.example
│   ├── package.json
│   └── vite.config.ts
│
├── etc/                             # Konfigurasi deployment
│   ├── Dockerfile.be
│   ├── Dockerfile.fe
│   ├── entrypoint.sh
│   ├── nginx.conf
│   ├── crontab
│   └── ai_chat_readonly_grant.sql
│
├── plans/                           # Dokumen PRD, FSD, TDD, ITL
│   ├── current/
│   ├── archive/
│   └── CHANGELOG.md
│
├── docker-compose.yml               # Backend + frontend
├── docker-compose.infra.yml         # MySQL + Redis
├── .env.example                     # Env untuk Docker Compose
└── README.md
```

## Prasyarat

Untuk menjalankan semuanya dengan Docker:

- Docker Engine 24+ atau Docker Desktop terbaru
- Docker Compose v2 (`docker compose`)
- Git

Untuk development lokal tanpa container aplikasi:

- Go 1.25+
- Node.js 22+ dan npm
- MySQL 8 dan Redis 7 (bisa memakai `docker-compose.infra.yml`)
- Git

Opsional:

- OpenCV + model FaceNet jika ingin mengaktifkan face recognition
- API key OpenAI-compatible jika ingin memakai AI chat
- SMTP server jika ingin memakai forgot/reset password lewat email

## Instalasi Cepat dengan Docker

Cara ini paling disarankan untuk mencoba aplikasi.

### 1. Clone dan masuk ke folder proyek

```bash
git clone <URL_REPOSITORY> hadiryuk
cd hadiryuk
```

### 2. Siapkan file environment

Linux/macOS/Git Bash:

```bash
cp .env.example .env
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

Buka `.env`, lalu isi minimal:

```env
DB_PASSWORD=ganti-password-mysql
REDIS_PASSWORD=ganti-password-redis

# Harus sama dengan user ACL yang dibuat infra
REDIS_USER=admin

# Host memakai nama service Docker
DB_HOST=mysql
REDIS_HOST=redis
DB_DATABASE=hadir_yuk
```

`DB_PASSWORD` dan `REDIS_PASSWORD` wajib diisi. Infra akan menolak start jika kosong.

### 3. Siapkan folder storage

```bash
mkdir -p be/storage
```

Folder ini di-bind mount ke container backend dan menyimpan foto, bukti absensi, log, serta key JWT.

### 4. Jalankan infra

```bash
docker compose -f docker-compose.infra.yml up -d
docker compose -f docker-compose.infra.yml ps
```

Perintah ini membuat:

- Container `mysql` — MySQL 8.0
- Container `redis` — Redis 7
- Network `mysql-network` dan `redis-network`
- Volume `mysql-data` dan `redis-data`

Database `hadir_yuk` dibuat otomatis oleh image MySQL saat volume masih kosong. Redis dikonfigurasi dengan ACL user `default` dan user dari `REDIS_USER` (default: `admin`), keduanya memakai password dari `REDIS_PASSWORD`.

### 5. Build dan jalankan aplikasi

```bash
docker compose up -d --build
docker compose ps
```

Backend berjalan di dalam network, frontend dipublikasikan di port **9010**.

### 6. Jalankan migration dan seed

```bash
docker exec -it hadiryuk-be /migration up
docker exec -it hadiryuk-be /migration seed
```

Opsional, untuk data demo (karyawan, shift, lokasi, riwayat absensi):

```bash
docker exec -it hadiryuk-be /migration dummy
```

> **Peringatan:** `migration up` menghapus tabel `attendances` sebelum membuat ulang schema. Jangan jalankan perintah ini pada database yang sudah berisi data absensi yang ingin dipertahankan. Lihat [Database & Migration](#database--migration).

### 7. Buka aplikasi

- Aplikasi: <http://localhost:9010>
- Health check: <http://localhost:9010/be/health>
- API base URL: `http://localhost:9010/be/api`

Akun default hasil `migration seed`:

| Role | Email | Password |
| --- | --- | --- |
| Super Admin | `suAdmin@app.com` | `@dmin#123` |
| HR Admin | `hradmin@app.com` | `HrAdmin#123` |

Data demo memakai password karyawan: `Karyawan#123`.

> Segera ganti password akun default setelah login pertama.

## Development Lokal

Mode ini menjalankan backend dan frontend langsung di host, sementara MySQL dan Redis tetap memakai Docker.

### 1. Jalankan infra saja

```bash
docker compose -f docker-compose.infra.yml up -d
```

Port yang tersedia di host:

- MySQL: `127.0.0.1:3306`
- Redis: `127.0.0.1:6379`

### 2. Jalankan backend

```bash
cd be
cp .env.example .env
```

Sesuaikan `.env`:

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=hadir_yuk
DB_USERNAME=root
DB_PASSWORD=<sama dengan DB_PASSWORD di .env root>

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_USER=admin
REDIS_PASSWORD=<sama dengan REDIS_PASSWORD di .env root>
```

Lalu:

```bash
go mod download
go run ./cmd/genkey -f
go run ./cmd/migration up
go run ./cmd/migration seed
go run ./cmd/api
```

Backend berjalan di <http://localhost:8080>. Health check: <http://localhost:8080/health>.

`go run ./cmd/genkey -f` membuat RSA key JWT di `be/storage/keys/`. Jika key sudah ada, generate ulang akan membuat token lama tidak valid.

Hot reload opsional dengan [Air](https://github.com/air-verse/air):

```bash
go install github.com/air-verse/air@latest
air
```

Konfigurasi Air ada di `be/.air.toml`.

### 3. Jalankan frontend

Buka terminal baru:

```bash
cd fe
cp .env.example .env
npm ci
npm run dev
```

Frontend berjalan di <http://localhost:5173>.

`.env` frontend:

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=HadirYuk
```

Alternatif: gunakan `VITE_API_BASE_URL=/be/api` agar request melewati proxy Vite ke `http://localhost:8080`.

### 4. Perintah development frontend

```bash
npm run dev       # jalankan dev server
npm run build     # type-check + build production
npm run preview   # preview hasil build
npm run lint      # ESLint + auto-fix
npm run format    # Prettier
```

## Environment Variables

Ada tiga file environment yang dipakai di proyek ini:

| File | Dipakai oleh | Keterangan |
| --- | --- | --- |
| `.env` | Docker Compose (`docker-compose.yml` dan `docker-compose.infra.yml`) | Konfigurasi container backend, frontend, MySQL, dan Redis. |
| `be/.env` | Backend saat dijalankan langsung di host | Gunakan `127.0.0.1` / `localhost` untuk host database dan Redis. |
| `fe/.env` | Frontend saat `npm run dev` / build | Berisi `VITE_API_BASE_URL` dan `VITE_APP_NAME`. |

### Variabel penting

| Variable | Contoh | Keterangan |
| --- | --- | --- |
| `TZ` | `Asia/Jakarta` | Timezone aplikasi dan container. |
| `DB_CONNECTION` | `mysql` | Driver database. Aplikasi menggunakan MySQL secara default. |
| `DB_HOST` | `mysql` | Nama service MySQL di Docker. Untuk backend lokal: `127.0.0.1`. |
| `DB_PORT` | `3306` | Port MySQL. |
| `DB_DATABASE` | `hadir_yuk` | Nama database yang dibuat otomatis oleh infra. |
| `DB_USERNAME` | `root` | User MySQL. |
| `DB_PASSWORD` | — | **Wajib diisi**, harus sama dengan password container MySQL. |
| `REDIS_ENABLED` | `true` | Aktifkan cache Redis. |
| `REDIS_HOST` | `redis` | Nama service Redis di Docker. Untuk backend lokal: `localhost`. |
| `REDIS_PORT` | `6379` | Port Redis. |
| `REDIS_USER` | `admin` | User ACL Redis. **Harus sama** dengan user yang dibuat infra. |
| `REDIS_PASSWORD` | — | **Wajib diisi**, harus sama dengan password Redis. |
| `REDIS_DB` | `0` | Index database Redis. |
| `APP_HOST` | `0.0.0.0` | Host bind backend. |
| `APP_PORT` | `8080` | Port internal backend. |
| `APP_URL` | `http://localhost:9010/be` | URL publik backend. |
| `APP_FE_URL` | `http://localhost:9010` | URL publik frontend. |
| `ALLOWED_ORIGINS` | `*` | Origin CORS yang diizinkan. Gunakan origin spesifik di production. |
| `TRUSTED_PROXIES` | kosong | Daftar proxy tepercaya, pisahkan dengan koma. |
| `GIN_MODE` | `debug` / `release` | Mode Gin. |
| `JWT_PRIVATE_KEY_PATH` | `storage/keys/private.pem` | Lokasi private key JWT. |
| `JWT_PUBLIC_KEY_PATH` | `storage/keys/public.pem` | Lokasi public key JWT. |
| `JWT_PASSPHRASE_PATH` | `storage/keys/passphrase` | Lokasi passphrase private key. |
| `JWT_EXPIRATION` | `24` | Masa berlaku access token (jam). |
| `JWT_REFRESH_EXPIRATION` | `168` | Masa berlaku refresh token (jam). |
| `RATE_LIMIT_REQUESTS` | `100` | Jumlah request per window. |
| `RATE_LIMIT_WINDOW` | `60` | Window rate limit (detik). |
| `SMTP_HOST` | kosong | Host SMTP untuk forgot/reset password. |
| `SMTP_PORT` | `587` | Port SMTP. |
| `SMTP_USER` | kosong | User SMTP. |
| `SMTP_PASSWORD` | kosong | Password SMTP. |
| `SMTP_FROM` | kosong | Alamat pengirim email. |

### Variabel AI Chat (opsional)

| Variable | Contoh | Keterangan |
| --- | --- | --- |
| `AI_CHAT_API_KEY` | — | API key provider. |
| `AI_CHAT_MODEL` | `gpt-4o-mini` | Model yang dipakai. |
| `AI_CHAT_BASE_URL` | kosong / `https://api.deepseek.com` | Base URL OpenAI-compatible. Kosongkan untuk OpenAI. |
| `AI_CHAT_DB_USERNAME` | `hadiryuk_ai_ro` | User MySQL read-only untuk Text-to-SQL. |
| `AI_CHAT_DB_PASSWORD` | — | Password user read-only. |
| `AI_CHAT_QUERY_TIMEOUT_SECONDS` | `10` | Timeout eksekusi query AI. |
| `AI_CHAT_MAX_ROWS` | `100` | Maksimum baris hasil query. |
| `AI_CHAT_SESSION_TTL_MINUTES` | `1440` | TTL session chat. |
| `AI_CHAT_MAX_STORED_MESSAGES` | `100` | Maksimum pesan tersimpan per session. |
| `AI_CHAT_MAX_CONTEXT_MESSAGES` | `10` | Maksimum pesan konteks ke model. |
| `AI_CHAT_RATE_LIMIT_GLOBAL_MAX` | `200` | Limit global. |
| `AI_CHAT_RATE_LIMIT_USER_MAX` | `20` | Limit per user. |

### Variabel frontend

| Variable | Contoh | Keterangan |
| --- | --- | --- |
| `VITE_API_BASE_URL` | `http://localhost:8080/api` | Base URL API untuk development lokal. Di Docker gunakan `/be/api`. |
| `VITE_APP_NAME` | `HadirYuk` | Nama aplikasi. |
| `VITE_BACKEND_URL` | `http://backend:8080` | Target proxy Vite saat menjalankan dev server. |

## Database & Migration

Migration dan seeder ada di `be/cmd/migration`.

### Membuat database

Database dibuat otomatis oleh image MySQL dari nilai `MYSQL_DATABASE` (diambil dari `DB_DATABASE`) saat container pertama kali diinisialisasi dengan volume kosong.

Jika volume `mysql-data` sudah ada tetapi database belum ada:

```bash
docker exec -it mysql mysql -uroot -p -e \
  "CREATE DATABASE IF NOT EXISTS hadir_yuk CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### Perintah migration

Di dalam container:

```bash
docker exec -it hadiryuk-be /migration up
docker exec -it hadiryuk-be /migration seed
docker exec -it hadiryuk-be /migration dummy
```

Di host (dari folder `be`):

```bash
go run ./cmd/migration up
go run ./cmd/migration seed
go run ./cmd/migration dummy
```

| Perintah | Fungsi |
| --- | --- |
| `up` | AutoMigrate semua model. **Menghapus tabel `attendances` lebih dulu.** |
| `down` | Drop semua tabel aplikasi. |
| `seed` | Seed permission, role, dan akun default. |
| `dummy` | Seed data demo: lokasi, shift, karyawan, assignment, riwayat absensi. |
| `refresh` | `down` + `up` + `seed`. Semua data aplikasi hilang. |

> **Penting:** `up` selalu drop tabel `attendances` sebelum migrate. Ini disengaja untuk perubahan schema, tetapi berarti data absensi lama akan hilang. Backup dulu jika database sudah berisi data penting.

### Membuat user MySQL read-only untuk AI Chat

Script contoh ada di `etc/ai_chat_readonly_grant.sql`. Sesuaikan nama database dan password di dalam script, lalu jalankan **setelah migration** karena grant mereferensikan tabel:

```bash
docker exec -i mysql sh -c \
  'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD"' \
  < etc/ai_chat_readonly_grant.sql
```

Setelah itu isi `AI_CHAT_DB_USERNAME` dan `AI_CHAT_DB_PASSWORD` di environment backend.

## API

Base URL:

- Docker: `http://localhost:9010/be/api`
- Local: `http://localhost:8080/api`

Endpoint publik:

- `GET /health`
- `GET /.well-known/jwks.json`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/forgot-password`
- `POST /api/auth/reset-password`

Endpoint lain memerlukan header:

```http
Authorization: Bearer <access_token>
```

### Grup endpoint utama

| Prefix | Fungsi |
| --- | --- |
| `/api/auth` | Login, refresh, logout, forgot/reset password. |
| `/api/me` | Profil sendiri, update profil, change password. |
| `/api/users` | CRUD user. |
| `/api/roles` | CRUD role. |
| `/api/permissions` | CRUD permission. |
| `/api/attendance` | Check-in/out, QR, riwayat, statistik, koreksi, laporan. |
| `/api/attendance/reports` | Laporan dan export Excel. |
| `/api/shifts` | CRUD shift dan assignment. |
| `/api/shifts/schedule` | Jadwal shift karyawan. |
| `/api/locations` | CRUD lokasi kantor. |
| `/api/qr-codes` | Generate, lihat, dan revoke QR Code. |
| `/api/dashboard` | Dashboard karyawan dan HR. |
| `/api/notifications` | Notifikasi in-app. |
| `/api/ai-chat` | Kirim pesan, riwayat, dan reset chat AI. |
| `/api/upload` | Upload file (multipart, field `file`). |

### Format response

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {}
}
```

Response paginasi menambahkan `metadata`:

```json
{
  "code": 200,
  "message": "OK",
  "data": [],
  "metadata": {
    "total": 0,
    "page": 1,
    "page_size": 10,
    "total_pages": 0
  }
}
```

Contoh login:

```bash
curl -X POST http://localhost:9010/be/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"suAdmin@app.com","password":"@dmin#123"}'
```

## Cron Jobs

Cron berjalan di dalam container backend menggunakan `crond` dari BusyBox. Konfigurasinya ada di `etc/crontab`.

| Jadwal | Job | Fungsi |
| --- | --- | --- |
| `01:00 WIB` setiap hari | `/markabsent` | Menandai karyawan yang tidak absen sebagai alpha untuk hari sebelumnya. |
| `02:00 WIB` setiap hari | `/cleartmp` | Menghapus file upload sementara yang lebih lama dari 24 jam. |

Contoh menjalankan manual:

```bash
docker exec -it hadiryuk-be /markabsent --dry-run
docker exec -it hadiryuk-be /cleartmp --dry-run
```

Log job ada di `be/storage/logs/`.

## Deployment & Maintenance

### Start

```bash
docker compose -f docker-compose.infra.yml up -d
docker compose up -d --build
```

### Stop

```bash
docker compose down
docker compose -f docker-compose.infra.yml down
```

Jangan menambahkan `-v` kecuali memang ingin menghapus volume database dan Redis.

### Update aplikasi

```bash
git pull
docker compose up -d --build
```

Migration tidak berjalan otomatis. Jalankan manual jika ada perubahan schema, dan ingat peringatan `migration up` yang menghapus tabel `attendances`.

### Backup

Database:

```bash
docker exec mysql sh -c \
  'exec mysqldump -uroot -p"$MYSQL_ROOT_PASSWORD" --single-transaction --routines --triggers "$MYSQL_DATABASE"' \
  > backup-hadiryuk.sql
```

Storage:

```bash
tar -czf backup-storage.tar.gz be/storage
```

Redis hanya menyimpan cache/session, jadi backup opsional.

### Restore database

```bash
docker exec -i mysql sh -c \
  'exec mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' \
  < backup-hadiryuk.sql
```

### Checklist keamanan production

- Ganti password akun default `suAdmin@app.com` dan `hradmin@app.com`.
- Gunakan `DB_PASSWORD` dan `REDIS_PASSWORD` yang kuat.
- Jangan commit file `.env`.
- Set `ALLOWED_ORIGINS` ke domain frontend, jangan `*`.
- Set `GIN_MODE=release`.
- Gunakan HTTPS di reverse proxy/load balancer.
- Batasi akses port MySQL dan Redis (infra sudah bind ke `127.0.0.1`).
- Backup `be/storage` dan database secara berkala.
- Pertimbangkan mengganti pasangan key JWT sebelum go-live.

## Testing

Backend:

```bash
cd be
go test ./...
go vet ./...
```

Frontend:

```bash
cd fe
npm run lint
npm run build
```

Saat ini belum ada test suite frontend; `npm run build` menjalankan type-check Vue TypeScript.

Untuk build dengan face recognition asli:

```bash
cd be
go build -tags gocv ./cmd/api
```

Build ini memerlukan OpenCV dan file model `haarcascade_frontalface_default.xml` serta `nn4.small2.v1.t7`. Tanpa tag `gocv`, aplikasi memakai stub.

## Troubleshooting

| Masalah | Kemungkinan penyebab | Solusi |
| --- | --- | --- |
| `DB_PASSWORD wajib diisi` | `DB_PASSWORD` kosong di `.env` | Isi `DB_PASSWORD` dengan password MySQL. |
| `REDIS_PASSWORD wajib diisi` | `REDIS_PASSWORD` kosong di `.env` | Isi `REDIS_PASSWORD`. |
| Backend gagal connect MySQL | Infra belum siap atau network salah | Cek `docker compose -f docker-compose.infra.yml ps`, lalu `docker logs mysql`. |
| Backend gagal connect Redis | `REDIS_USER` / password tidak cocok | Pastikan `REDIS_USER=admin` dan password sama dengan yang dipakai infra. |
| `WRONGPASS` / `NOAUTH` | ACL Redis belum cocok | Cek `docker exec -it redis redis-cli --user admin --pass <password> ping`. |
| JWT key tidak ditemukan | Key belum dibuat | Docker: restart backend agar entrypoint generate key. Lokal: `go run ./cmd/genkey -f`. |
| Port 3306/6379 sudah dipakai | Ada MySQL/Redis lain di host | Hentikan service lain atau ubah mapping port dan sesuaikan env backend lokal. |
| Halaman blank setelah rebuild | Cache Service Worker/PWA | Hard refresh, clear site data, atau buka di incognito. |
| Upload/foto gagal | Permission folder storage | Pastikan `be/storage` dapat ditulis container. Di Linux: `sudo chown -R $USER:$USER be/storage`. |
| 502 dari Nginx | Backend belum siap atau crash | Cek `docker compose logs -f backend`. |
| Data absensi hilang setelah migrate | `migration up` drop tabel `attendances` | Restore backup. Jangan jalankan `up` di database production tanpa backup. |

## Dokumentasi

Dokumen perencanaan dan desain tersedia di:

- `plans/current/01_PRD.md` — Product Requirements Document
- `plans/current/02_FSD.md` — Functional Specification Document
- `plans/current/03_Role_Matrix.md` — matriks role dan permission
- `plans/current/04_TDD.md` — Technical Design Document
- `plans/current/05_ITL.md` — Implementation Task List
- `plans/current/AI_CS_Chatbot.md` — desain AI Customer Service Chatbot
- `plans/CHANGELOG.md` — changelog proyek

## Kontribusi

- Ikuti pola layer backend: `handler → service → repository`.
- Tempatkan konfigurasi runtime di environment, bukan hardcode.
- Jalankan `go test ./...`, `npm run lint`, dan `npm run build` sebelum membuat pull request.
- Jangan menambahkan data dummy atau kredensial ke environment production.
- Untuk perubahan schema, backup database sebelum menjalankan migration.
