# `go_vibe` Weather & Analytics Backend Service

Backend service untuk aplikasi **Weather Forecast & Analytics Dashboard** dibangun dengan bahasa Go (Golang) dan menerapkan Clean Architecture.

---

## 📁 Struktur Direktori

```text
go_vibe/
├── main.go                     # Entry point aplikasi & dependency injection
├── config/
│   └── database.go             # Inisialisasi PostgreSQL & environment loader
├── domain/                     # Entities, DTOs, dan Interface Contracts
│   ├── user.go
│   ├── location.go
│   ├── weather.go
│   └── response.go
├── internal/
│   ├── delivery/
│   │   └── http/               # HTTP Handlers, Routes Gin & Tests
│   │       ├── auth_handler.go
│   │       ├── auth_and_user_handler_test.go
│   │       ├── location_handler.go
│   │       ├── location_handler_test.go
│   │       ├── user_handler.go
│   │       ├── weather_handler.go
│   │       ├── weather_handler_test.go
│   │       └── router.go
│   ├── repository/             # Implementasi Database (GORM) & Weather API
│   │   ├── location_repository.go
│   │   ├── user_repository.go
│   │   └── weather_repository.go
│   ├── usecase/                # Business Logic & Analytics
│   │   ├── location_usecase.go
│   │   ├── user_usecase.go
│   │   ├── weather_usecase.go
│   │   └── weather_usecase_test.go
│   └── middleware/             # JWT Auth Middleware & Optional Auth
│       └── auth_middleware.go
├── migrations/                 # Script Migration SQL
│   ├── 000001_create_users_table.up.sql
│   ├── 000002_create_locations_table.up.sql
│   └── 000003_create_user_locations_table.up.sql
├── .env.example
├── .env
├── go.mod
├── go.sum
└── README.md
```

---

## 🚀 Panduan Menjalankan

### 1. Download & Verifikasi Dependencies
```bash
go mod tidy
```

### 2. Jalankan Migrasi Database (Opsional jika PostgreSQL aktif)
```bash
migrate -path migrations -database "postgres://postgres:secret@localhost:5432/go_vibe_weather_db?sslmode=disable" up
```

### 3. Jalankan Pengujian Otomatis (Unit & Integration Tests)
```bash
go test -v ./...
```

### 4. Jalankan Aplikasi
```bash
go run .
```

Setelah server aktif:
* **Tampilan Dashboard Web**: Buka browser ke [http://localhost:8080](http://localhost:8080)
* **Status Health Check**: `http://localhost:8080/health`
* **Dokumentasi Endpoints API**: `http://localhost:8080/api` (atau via header `Accept: application/json` pada root `/`)

---

## 📌 Daftar Lengkap Endpoint RESTful API

| Method | Endpoint | Akses | Deskripsi Fitur |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Public | Health Check Service |
| `POST` | `/api/v1/auth/register` | Public | Registrasi akun baru & dapatkan JWT |
| `POST` | `/api/v1/auth/login` | Public | Login akun & dapatkan JWT |
| `GET` | `/api/v1/users/profile` | Protected (JWT) | Ambil profil pengguna & unit suhu |
| `PUT` | `/api/v1/users/settings` | Protected (JWT) | **Pengaturan Pengguna**: Ubah satuan suhu (`Celsius` / `Fahrenheit`) |
| `GET` | `/api/v1/weather/dashboard` | Public / Auth | Ringkasan dashboard lengkap (Realtime, Highlights, Hourly, Rain, Forecast) |
| `GET` | `/api/v1/weather/hourly` | Public / Auth | **Grafik Suhu Per Jam**: Fluktuasi suhu runtut waktu, tren, min/max/avg |
| `GET` | `/api/v1/weather/chance-of-rain` | Public / Auth | **Analisis Peluang Hujan**: Probabilitas, intensitas, analisis rentang waktu harian & advisory |
| `GET` | `/api/v1/weather/forecast` | Public / Auth | **Prakiraan 3 Hari Ke Depan**: Prediksi cuaca dengan estimasi suhu maksimum ($\uparrow$) & minimum ($\downarrow$) |
| `GET` | `/api/v1/weather/calendar` | Public / Auth | **Kalender Cuaca Bulanan**: Prakiraan cuaca 1 bulan penuh (hari, tanggal, kondisi cuaca, suhu, peluang hujan, advisory aktivitas & navigasi bulan) |
| `GET` | `/api/v1/weather/compare` | Public / Auth | **Perbandingan 2 Kota**: Akses dan bandingkan cuaca 2 kota (suhu, kelembapan, angin, hujan, grafik komparatif & saran perjalanan) |
| `GET` | `/api/v1/locations/search` | Protected (JWT) | **Pencarian Wilayah**: Cari daerah di wilayah Indonesia berdasarkan nama kota atau provinsi (`?q=Jakarta`) |
| `POST` | `/api/v1/locations/saved` | Protected (JWT) | **Tambah Kota Favorit**: Simpan kota favorit pengguna (khusus wilayah Indonesia) |
| `GET` | `/api/v1/locations/saved` | Protected (JWT) | **Daftar Kota Favorit**: Ambil seluruh kota tersimpan milik pengguna |

---

## 📖 Ringkasan Fitur Utama

### 1. Grafik Suhu Per Jam (Hourly Temperature Chart)
* **Endpoint:** `GET /api/v1/weather/hourly?unit=celsius`
* Menyediakan data terformat untuk grafik: `labels` waktu (`["06 am", "09 am", ...]`), `data_points` suhu numerik (`[10, 14, 18, ...]`), `min_temp`, `max_temp`, `avg_temp`, `peak_time`, `lowest_time`, dan deskripsi tren cuaca harian.

### 2. Analisis Peluang Hujan (Chance of Rain Analysis)
* **Endpoint:** `GET /api/v1/weather/chance-of-rain`
* Menganalisis probabilitas hujan dengan klasifikasi intensitas harian:
  * Agregasi 4 rentang waktu harian: **Pagi** (`06:00 - 12:00`), **Siang** (`12:00 - 18:00`), **Sore/Malam** (`18:00 - 24:00`), **Dini Hari** (`00:00 - 06:00`).
  * Penilaian tingkat risiko: `Low`, `Moderate`, `High`, atau `Very High`.
  * Rekomendasi/saran cuaca praktis (`advisory`) untuk persiapan pengguna.

### 3. Prakiraan 3 Hari Ke Depan (3 Days Forecast)
* **Endpoint:** `GET /api/v1/weather/forecast?days=3`
* Ringkasan ramalan cuaca multi-hari dengan indikator suhu maksimum ($\uparrow$) dan minimum ($\downarrow$), nama hari, tanggal, kondisi cuaca, ikon, dan narasi ringkasan.

### 4. Pengelolaan Lokasi & Tambah Kota (Locations & Add City)
* **Search:** `GET /api/v1/locations/search?q={city}` mencari daerah di wilayah Indonesia (dengan validasi khusus daerah Indonesia dan pembersihan data luar negeri).
* **Bookmark:** `POST /api/v1/locations/saved` menyimpan relasi lokasi ke tabel `user_saved_locations`, otomatis menandai kota pertama sebagai *default*.
* **List:** `GET /api/v1/locations/saved` menampilkan daftar kota favorit pengguna.

### 5. Pengaturan Pengguna (User Settings)
* **Endpoint:** `PUT /api/v1/users/settings`
* Mengubah pengaturan unit suhu global (`Celsius` atau `Fahrenheit`). Perubahan preferensi ini secara otomatis mempengaruhi respons suhu di seluruh endpoint cuaca ketika diakses oleh pengguna terkait.

### 6. Kalender Cuaca Bulanan (Monthly Weather Calendar)
* **Endpoint:** `GET /api/v1/weather/calendar?city={city}&year={year}&month={month}&unit={unit}`
* **Fitur Utama Web:**
  * Akses langsung melalui menu navigasi sidebar **"Calendar"**, tombol header **"Kalender Cuaca"**, atau kartu quick action hero.
  * Tampilan kalender grid bulanan lengkap dengan navigasi bulan (sebelumnya, berikutnya, dan bulan ini).
  * Filter kota cepat (Jakarta, Bandung, Surabaya, Yogyakarta, Semarang, Denpasar, Medan, Makassar, dll).
  * Statistik bulanan interaktif: rata-rata suhu maks/min, jumlah hari cerah vs hari hujan, rata-rata peluang hujan, serta saran rekomendasi bulanan.
  * Inspector tanggal terpilih: klik pada tanggal manapun di kalender untuk melihat rincian suhu lengkap, persentase hujan, dan saran aktivitas outdoor yang direkomendasikan.
