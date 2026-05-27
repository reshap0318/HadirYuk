# Rencana Prioritas Fitur — HadirYuk

> **Status:** Pasca Data Master + Notifikasi + Peta + Manajemen Profil Lengkap
> **Terakhir Diperbarui:** 2026-05-27
> **Stack Teknologi:** Backend Go (Gin) + Frontend Vue 3 + MySQL

---

## Selesai (Fondasi)

| Area | Fitur | Status | Catatan |
|------|-------|--------|---------|
| **Autentikasi** | Login (JWT RS256), Logout, Refresh Token | Selesai | Token JWT RS256 dengan JWKS, sesi Redis, validasi token |
| **Autentikasi** | Lupa Password & Reset Password | Selesai | Backend + Frontend lengkap, token reset berkedaluwarsa 1 jam, email async |
| **UAM** | CRUD Role, CRUD Permission, Pemetaan Role-Permission, Penugasan User-Role | Selesai | RBAC 5 tabel, middleware permission, caching akses |
| **Data Master** | CRUD Shift, CRUD Lokasi Kantor, CRUD Jenis Cuti | Selesai | CRUD lengkap dengan paginasi, notifikasi otomatis saat CRUD |
| **Manajemen User** | CRUD User dengan Profil | Selesai | Profil otomatis dibuat, field: departemen, jabatan, telepon, tanggal bergabung |
| **Manajemen Profil** | Lihat & Edit Profil (/me) | Selesai | Upload avatar, ubah password, validasi form, update sesi |
| **Notifikasi** | CRUD Notifikasi, Tandai Dibaca, Hitung Belum Dibaca | Selesai | Notifikasi otomatis saat CRUD lokasi, filter berdasarkan tipe & status baca |
| **Upload File** | Upload File dengan UUID | Selesai | Penyimpanan sementara, validasi ekstensi, URL file |
| **Peta OpenStreetMap** | Komponen Peta Leaflet | Selesai | Marker seret, lingkaran radius, geolokasi browser, slider radius 50-500m |
| **Health Check** | Endpoint Kesehatan Sistem | Selesai | Status database, Redis, JWKS publik |
| **Penanganan Error** | Validasi & Error Handling | Selesai | Error field per-input, SweetAlert notifikasi, pesan error Bahasa Indonesia |

---

## Tingkatan Prioritas

### P0 — Jalur Kritis (Absensi Inti)

> Tanpa fitur ini, sistem tidak dapat memenuhi tujuan utamanya: mencatat absensi.

| # | Fitur | Kompleksitas | Dependensi | Peran Utama | Estimasi |
|---|-------|-------------|------------|-------------|----------|
| P0-1 | **Penugasan Shift Karyawan** | M | CRUD Shift (selesai), CRUD User (selesai) | HR Admin | 2 hari |
| P0-2 | **Upload Foto Wajah untuk Pengenalan Wajah** | M | Manajemen Profil (selesai) | Karyawan | 2 hari |
| P0-3 | **Absensi Check-in (Geotagging + Pengenalan Wajah)** | XL | P0-1, P0-2, Lokasi Kantor (selesai) | Karyawan | 5 hari |
| P0-4 | **Absensi Check-out (Geotagging + Pengenalan Wajah)** | L | P0-3 | Karyawan | 3 hari |
| P0-5 | **Pembuatan & Manajemen Kode QR** | M | Lokasi Kantor (selesai) | HR Admin | 2 hari |
| P0-6 | **Absensi Check-in (Kode QR)** | L | P0-5 | Karyawan | 3 hari |
| P0-7 | **Absensi Check-out (Kode QR)** | S | P0-6 | Karyawan | 1 hari |

**Total Estimasi P0: ~18 hari**

**Mengapa urutan ini:**
1. Karyawan harus ditugaskan ke shift sebelum absensi dapat divalidasi (shift menentukan jam kerja, toleransi, durasi istirahat).
2. Foto wajah harus ada sebelum pengenalan wajah dapat bekerja saat check-in.
3. Check-in geotagging adalah metode absensi utama — bangun terlebih dahulu.
4. Check-out membangun langsung dari logika check-in (catatan yang sama, status berbeda).
5. Kode QR harus dibuat sebelum dapat dipindai.
6. Check-in/out QR adalah metode sekunder — lebih sederhana dari geotagging (tanpa pengenalan wajah).

---

### P1 — Penting (Layanan Mandiri Karyawan)

> Fitur yang dibutuhkan karyawan setiap hari untuk berinteraksi dengan sistem secara bermakna.

| # | Fitur | Kompleksitas | Dependensi | Peran Utama | Estimasi |
|---|-------|-------------|------------|-------------|----------|
| P1-1 | **Riwayat & Detail Absensi** | M | P0-3, P0-4 | Karyawan | 2 hari |
| P1-2 | **Dasbor Karyawan** | M | P0-1, P0-3, P0-4 | Karyawan | 3 hari |
| P1-3 | **Pengajuan Cuti** | M | Jenis Cuti (selesai) | Karyawan | 2 hari |
| P1-4 | **Tampilan Sisa Cuti** | S | P1-3 | Karyawan | 1 hari |
| P1-5 | **Ubah Password** | S | Autentikasi (selesai) | Semua Peran | Selesai |
| P1-6 | **Lupa/Reset Password** | M | Autentikasi (selesai), Layanan email | Semua Peran | Selesai |

**Total Estimasi P1: ~8 hari** (P1-5 dan P1-6 sudah selesai)

**Mengapa urutan ini:**
1. Setelah karyawan dapat check-in/out, mereka perlu melihat catatan absensi mereka.
2. Dasbor menggabungkan data absensi + shift menjadi satu tampilan.
3. Pengajuan cuti adalah fitur layanan mandiri inti yang independen dari absensi.
4. Sisa cuti bergantung pada logika pengajuan cuti (menghitung yang terpakai vs total).
5. Manajemen password sudah selesai — berdampak tinggi untuk pengalaman pengguna.

---

### P2 — Sebaiknya Ada (Fitur Operasional HR)

> Fitur yang dibutuhkan HR Admin untuk mengelola sistem dan menangani pengecualian.

| # | Fitur | Kompleksitas | Dependensi | Peran Utama | Estimasi |
|---|-------|-------------|------------|-------------|----------|
| P2-1 | **Dasbor HR** | L | P0-3, P0-4, P1-1, P1-3 | HR Admin | 3 hari |
| P2-2 | **Koreksi Absensi** | M | P1-1 | HR Admin | 2 hari |
| P2-3 | **Tampilan Jadwal Shift** | M | P0-1 | Karyawan, HR Admin | 2 hari |
| P2-4 | **Laporan Absensi (Lihat & Ekspor Excel/PDF)** | L | P1-1 | HR Admin | 4 hari |
| P2-5 | **Laporan Cuti** | M | P1-3, P1-4 | HR Admin | 2 hari |

**Total Estimasi P2: ~13 hari**

**Mengapa urutan ini:**
1. Dasbor HR adalah pusat komando — membutuhkan data absensi + cuti agar bermakna.
2. Koreksi absensi menangani kasus tepi (lupa check-out, lokasi salah).
3. Tampilan jadwal shift membantu karyawan dan HR memvisualisasikan penugasan.
4. Laporan dibangun di atas data riwayat absensi.
5. Laporan cuti bergantung pada pengajuan/sisa cuti yang fungsional.

---

### P3 — Bagus untuk Dimiliki (Analitik & Pelaporan)

> Fitur yang memberikan wawasan tetapi tidak diperlukan untuk operasi sehari-hari.

| # | Fitur | Kompleksitas | Dependensi | Peran Utama | Estimasi |
|---|-------|-------------|------------|-------------|----------|
| P3-1 | **Statistik Keterlambatan** | M | P1-1, P0-1 | HR Admin | 2 hari |
| P3-2 | **Dasbor Admin** | S | UAM (selesai) | Super Admin | 1 hari |

**Total Estimasi P3: ~3 hari**

**Mengapa urutan ini:**
1. Statistik keterlambatan memerlukan riwayat absensi yang berfungsi penuh.
2. Dasbor Admin ringan — sebagian besar statistik sistem dan aktivitas terbaru.

---

### P4 — Bisa Ditunda (Audit & Lanjutan)

> Fitur yang menambah polesan dan kepatuhan tetapi dapat ditunda.

| # | Fitur | Kompleksitas | Dependensi | Peran Utama | Estimasi |
|---|-------|-------------|------------|-------------|----------|
| P4-1 | **Log Audit** | L | Semua operasi CRUD | Super Admin | 3 hari |

**Total Estimasi P4: ~3 hari**

**Mengapa urutan ini:**
1. Log audit adalah perhatian lintas modul — paling baik diimplementasikan setelah semua operasi CRUD stabil sehingga setiap aksi dapat dicatat secara konsisten.

---

## Peta Jalan Implementasi (Bertahap)

```
Tahap 1 — P0: Mesin Absensi Inti          (~18 hari)
  Minggu 1-2: Penugasan Shift, Foto Wajah, Check-in/out Geotagging
  Minggu 3:   Pembuatan Kode QR, Check-in/out QR

Tahap 2 — P1: Layanan Mandiri Karyawan    (~8 hari)
  Minggu 4:   Riwayat Absensi, Dasbor Karyawan
  Minggu 5:   Pengajuan Cuti, Sisa Cuti

Tahap 3 — P2: Operasional HR              (~13 hari)
  Minggu 6-7: Dasbor HR, Koreksi Absensi, Jadwal Shift
  Minggu 8:   Laporan Absensi, Laporan Cuti

Tahap 4 — P3: Analitik                    (~3 hari)
  Minggu 9:   Statistik Keterlambatan, Dasbor Admin

Tahap 5 — P4: Audit                       (~3 hari)
  Minggu 10:  Implementasi Log Audit di semua modul
```

**Total Estimasi Usaha: ~45 hari kerja**

---

## Grafik Dependensi

```
Data Master (selesai) ──┬── P0-1 Penugasan Shift Karyawan ──┬── P0-3 Check-in (Geotagging)
                         │                                     ├── P0-4 Check-out (Geotagging)
                         │                                     ├── P1-2 Dasbor Karyawan
                         │                                     └── P2-3 Tampilan Jadwal Shift
                         │
                         ├── P0-2 Foto Wajah ───────────────── P0-3 Check-in (Geotagging)
                         │
                         ├── P0-5 Pembuatan Kode QR ───────────┬── P0-6 Check-in (QR)
                         │                                      └── P0-7 Check-out (QR)
                         │
                         └── Jenis Cuti (selesai) ───────────── P1-3 Pengajuan Cuti ── P1-4 Sisa Cuti
                                                                                   └── P2-5 Laporan Cuti

P0-3 + P0-4 ─────────────────────────────┬── P1-1 Riwayat Absensi ──┬── P2-2 Koreksi Absensi
                                         │                          ├── P2-1 Dasbor HR
                                         │                          ├── P2-4 Laporan Absensi
                                         │                          └── P3-1 Statistik Keterlambatan
                                         │
                                         └── P1-2 Dasbor Karyawan

Autentikasi (selesai) ───────────────────┬── P1-5 Ubah Password (SELESAI)
                                         └── P1-6 Lupa/Reset Password (SELESAI)

UAM (selesai) ─────────────────────────── P3-2 Dasbor Admin

Semua operasi CRUD ────────────────────── P4-1 Log Audit
```

---

## Risiko & Mitigasi

| Risiko | Dampak | Mitigasi |
|--------|--------|----------|
| Integrasi pustaka pengenalan wajah | Tinggi | Prototipe pengenalan wajah lebih awal di P0-2; siapkan fallback ke check-in hanya foto |
| Akurasi GPS di browser mobile | Sedang | Implementasi perhitungan jarak haversine dengan toleransi yang dapat dikonfigurasi; catat koordinat mentah untuk debugging |
| Waktu kedaluwarsa Kode QR | Rendah | Gunakan validasi waktu sisi server; tambahkan periode toleransi 30 detik |
| Akurasi perhitungan sisa cuti | Sedang | Tulis uji unit untuk perhitungan hari kerja (kecuali akhir pekan) |
| Performa pembuatan file ekspor | Rendah | Streaming dataset besar; tambahkan paginasi untuk laporan dengan >1000 baris |

---

## Catatan

- **Estimasi usaha** mengasumsikan satu pengembang full-stack yang bekerja di backend (Go/Gin) dan frontend (Vue 3) secara bersamaan.
- **Pekerjaan paralel:** Frontend dan backend untuk setiap fitur dapat dikembangkan secara paralel setelah kontrak API didefinisikan.
- **Pengujian:** Setiap tingkatan harus mencakup uji unit + uji integrasi sebelum berpindah ke tingkatan berikutnya.
- **Milestone:** Setelah P0 selesai, sistem dapat mencatat absensi — ini adalah produk minimum yang layak (MVP).

---

## Fitur yang Sudah Selesai — Rincian Teknis

### Autentikasi
- ✅ **Login** — JWT RS256 dengan JWKS, penyimpanan sesi Redis, respons dengan token + refresh token + data user + permission
- ✅ **Logout** — Pembersihan token sisi klien, pemanggilan endpoint logout
- ✅ **Refresh Token** — Validasi refresh token, pembuatan token baru, pembaruan sesi
- ✅ **Lupa Password** — Pembuatan token reset acak, hashing token, penyimpanan ke database, pengiriman email async
- ✅ **Reset Password** — Validasi token (kedaluwarsa, sudah dipakai), transaksi database, hash password baru

### Manajemen Akses Pengguna (UAM)
- ✅ **CRUD Role** — Buat, baca, perbarui, hapus role dengan paginasi
- ✅ **CRUD Permission** — Buat, baca, perbarui, hapus permission dengan format `{module}.{action}`
- ✅ **Pemetaan Role-Permission** — Tambah/hapus permission dari role
- ✅ **Penugasan User-Role** — Tambah/hapus role dari user saat membuat/memperbarui user
- ✅ **Middleware Permission** — `RequirePermission` untuk proteksi endpoint
- ✅ **Caching Akses** — Cache permission per user dengan invalidasi otomatis

### Data Master
- ✅ **CRUD Shift** — Nama unik, jam mulai/selesai, durasi istirahat, kode warna, total jam
- ✅ **CRUD Lokasi Kantor** — Nama, alamat, latitude/longitude, radius meter, status aktif
- ✅ **CRUD Jenis Cuti** — Nama unik, deskripsi, hari default, status berbayar

### Manajemen User & Profil
- ✅ **CRUD User** — Email unik, password ter-hash, avatar, profil otomatis, penugasan role
- ✅ **Profil Saya (/me)** — Lihat profil, perbarui profil (nama, telepon, departemen, jabatan, avatar)
- ✅ **Upload Avatar** — Upload file dengan UUID, pemindahan dari tmp ke avatars, hapus file lama
- ✅ **Ubah Password via Profil** — Validasi password minimal 6 karakter, konfirmasi password cocok
- ✅ **Field Profil** — Departemen, jabatan, telepon, tanggal bergabung, foto wajah, embedding wajah

### Notifikasi
- ✅ **CRUD Notifikasi** — Buat, baca (paginasi), hapus notifikasi per user
- ✅ **Tandai Dibaca** — Tandai satu notifikasi atau semua notifikasi sebagai dibaca
- ✅ **Hitung Belum Dibaca** — Endpoint untuk jumlah notifikasi belum dibaca
- ✅ **Filter** — Filter berdasarkan status baca dan tipe notifikasi
- ✅ **Notifikasi Otomatis** — Dibuat otomatis saat CRUD lokasi kantor

### Peta OpenStreetMap
- ✅ **Komponen FormMap** — Peta Leaflet dengan tile OpenStreetMap
- ✅ **Marker Seret** — Marker dapat diseret untuk memilih lokasi
- ✅ **Lingkaran Radius** — Visualisasi radius validasi absensi
- ✅ **Geolokasi Browser** — Tombol "Gunakan Lokasi Saya" dengan akurasi tinggi
- ✅ **Slider Radius** — Pengaturan radius 50-500 meter
- ✅ **Koordinat Real-time** — Tampilan latitude/longitude dengan 6 desimal
- ✅ **Integrasi Form Lokasi** — FormModal lokasi menggunakan FormMap untuk pemilihan lokasi

### Upload File
- ✅ **Endpoint Upload** — Upload file dengan validasi ekstensi, penyimpanan dengan UUID
- ✅ **Penyimpanan Sementara** — File disimpan di `storage/tmp` sebelum dipindahkan
- ✅ **URL File** — Generate URL file yang dapat diakses

### Kesehatan Sistem
- ✅ **Endpoint /health** — Status database (ping), status Redis (jika tersedia), status keseluruhan
- ✅ **Endpoint JWKS** — Publik key JSON Web Key Set untuk validasi token

### Penanganan Error
- ✅ **Error Validasi** — Pesan error per field, format 422 Unprocessable Entity
- ✅ **Error Bidang** — FieldError untuk validasi bisnis (email sudah ada, dll)
- ✅ **SweetAlert** — Notifikasi sukses/error di frontend dengan pesan Bahasa Indonesia
- ✅ **Error Autentikasi** — Pesan error untuk token tidak valid, kredensial salah, token kedaluwarsa
