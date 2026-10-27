# Project Specification: `go_vibe` Weather & Analytics Backend Service

## 1. Overview

`go_vibe` adalah layanan backend berbasis RESTful API menggunakan Go (Golang) untuk mendukung aplikasi **Weather Forecast & Analytics Dashboard**. Backend ini bertanggung jawab menyediakan data cuaca real-time, prakiraan cuaca per jam/harian, analisis probabilitas hujan (*Chance of Rain*), pengolahan highlight cuaca harian (*humidity, wind, sunrise/sunset*), manajemen lokasi/kota favorit pengguna, serta autentikasi dan pengaturan akun.

---

## 2. Tech Stack & Dependencies

* **Language:** Go (v1.21+)
* **Web Framework:** Chi Router / Fiber / Gin
* **Database:** PostgreSQL
* **ORM / Query Builder:** GORM / `pgx` (PostgreSQL Driver)
* **Authentication:** JWT (JSON Web Tokens) & `golang.org/x/crypto/bcrypt`
* **External Integration:** OpenWeatherMap API / WeatherAPI (untuk penyedia data cuaca real-time & prakiraan)
* **Environment Management:** `godotenv`
* **Database Migration:** `golang-migrate`

---

## 3. Directory Structure

```text
go_vibe/
├── cmd/
│   └── main.go                 # Entry point aplikasi
├── config/
│   └── database.go             # Inisialisasi PostgreSQL & Env
├── internal/
│   ├── delivery/
│   │   └── http/               # HTTP Handlers & Routes
│   │       ├── auth_handler.go
│   │       ├── weather_handler.go
│   │       ├── location_handler.go
│   │       └── user_handler.go
│   ├── repository/             # Layer Interaksi Database / External APIs
│   │   ├── weather_repository.go
│   │   ├── location_repository.go
│   │   └── user_repository.go
│   ├── usecase/                # Layer Business Logic & Analytics
│   │   ├── weather_usecase.go
│   │   ├── location_usecase.go
│   │   └── user_usecase.go
│   └── middleware/             # Middleware (JWT Auth, CORS, Rate Limiter)
│       └── auth_middleware.go
├── domain/                     # Entities & Models
│   ├── user.go
│   ├── location.go
│   ├── weather.go
│   └── response.go
├── migrations/                 # Script Migration SQL
│   ├── 000001_create_users_table.up.sql
│   ├── 000002_create_locations_table.up.sql
│   └── 000003_create_user_locations_table.up.sql
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

---

## 4. Database Schema Specification

### 1. Table: `users`
| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | Primary Key | ID Unik pengguna |
| `name` | `VARCHAR(100)` | NOT NULL | Nama pengguna |
| `email` | `VARCHAR(150)` | UNIQUE, NOT NULL | Email pengguna |
| `password` | `VARCHAR(255)` | NOT NULL | Password terenkripsi (Bcrypt) |
| `temperature_unit` | `VARCHAR(10)` | DEFAULT `'Celsius'` | Unit suhu (`Celsius` / `Fahrenheit`) |
| `created_at` | `TIMESTAMP` | DEFAULT `NOW()` | Waktu pembuatan akun |
| `updated_at` | `TIMESTAMP` | DEFAULT `NOW()` | Pembaruan terakhir |

### 2. Table: `locations`
| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | Primary Key | ID Unik Lokasi |
| `city` | `VARCHAR(100)` | NOT NULL | Nama kota (misal: Jakarta) |
| `state` | `VARCHAR(100)` | NULL | Nama provinsi (misal: DKI Jakarta) |
| `country` | `VARCHAR(100)` | NOT NULL | Nama negara (Indonesia) |
| `latitude` | `DECIMAL(10, 8)` | NOT NULL | Koordinat Latitude |
| `longitude` | `DECIMAL(11, 8)` | NOT NULL | Koordinat Longitude |
| `image_url` | `TEXT` | NULL | Foto latar belakang lokasi (misal: Monas) |

### 3. Table: `user_saved_locations`
| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | Primary Key | ID Unik relasi |
| `user_id` | `UUID` | Foreign Key (`users.id`) | Reference ke tabel User |
| `location_id` | `UUID` | Foreign Key (`locations.id`) | Reference ke tabel Location |
| `is_default` | `BOOLEAN` | DEFAULT `FALSE` | Penanda lokasi utama (Current Location) |
| `created_at` | `TIMESTAMP` | DEFAULT `NOW()` | Tanggal disimpan |

---

## 5. Dashboard Data & API Endpoints Specification

### A. Authentication & Settings

#### 1. Register User
* **Endpoint:** `POST /api/v1/auth/register`
* **Access:** Public

#### 2. Login User
* **Endpoint:** `POST /api/v1/auth/login`
* **Access:** Public

#### 3. Update User Settings
* **Endpoint:** `PUT /api/v1/users/settings`
* **Access:** Authenticated (Bearer Token)
* **Request Body:**
  ```json
  {
    "temperature_unit": "Celsius"
  }
  ```

---

### B. Weather & Dashboard Analytics

#### 4. Get Dashboard Weather Overview (Current, Highlights, Hourly, Forecast)
* **Endpoint:** `GET /api/v1/weather/dashboard`
* **Access:** Authenticated / Public
* **Query Params:** `location_id` (opsional) atau default ke *Current Location*.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "data": {
      "location": {
        "id": "loc-12345",
        "city": "Jakarta",
        "state": "DKI Jakarta",
        "country": "Indonesia",
        "image_url": "https://images.unsplash.com/photo-1555899434-94d1368aa7af?auto=format&fit=crop&w=800&q=80"
      },
      "current_weather": {
        "temperature": 12,
        "condition": "Partly Cloudy",
        "condition_icon": "partly-cloudy-snow",
        "date_time": "Monday, 07:43 AM"
      },
      "highlights": {
        "precipitation": "2%",
        "humidity": "87%",
        "wind_speed": "0 km/h",
        "sunrise": "06:18 am",
        "sunset": "07:27 pm"
      },
      "hourly_temperature": [
        { "time": "06 am", "temp": 10 },
        { "time": "09 am", "temp": 14 },
        { "time": "12 pm", "temp": 18 },
        { "time": "03 pm", "temp": 22 },
        { "time": "06 pm", "temp": 19 },
        { "time": "09 pm", "temp": 14 },
        { "time": "12 am", "temp": 12 },
        { "time": "03 am", "temp": 10 }
      ],
      "chance_of_rain": [
        { "time": "09 am", "intensity": "Sunny", "percentage": 10 },
        { "time": "12 pm", "intensity": "Rainy", "percentage": 60 },
        { "time": "03 pm", "intensity": "Rainy", "percentage": 40 },
        { "time": "06 pm", "intensity": "Heavy Rain", "percentage": 90 },
        { "time": "09 pm", "intensity": "Rainy", "percentage": 30 },
        { "time": "12 am", "intensity": "Sunny", "percentage": 15 },
        { "time": "03 am", "intensity": "Sunny", "percentage": 5 }
      ],
      "three_days_forecast": [
        {
          "day": "Tuesday",
          "condition": "Cloudy",
          "temp_high": 26,
          "temp_low": 11,
          "icon": "cloudy"
        },
        {
          "day": "Wednesday",
          "condition": "Rainy",
          "temp_high": 26,
          "temp_low": 11,
          "icon": "rainy"
        },
        {
          "day": "Thursday",
          "condition": "Snowfall",
          "temp_high": 26,
          "temp_low": 11,
          "icon": "snow"
        }
      ]
    }
  }
  ```

---

### C. Location Management (*Add City & Locations Menu*)

#### 5. Search Location / City
* **Endpoint:** `GET /api/v1/locations/search?q={city_name}`
* **Access:** Authenticated

#### 6. Add Saved Location (*Add City*)
* **Endpoint:** `POST /api/v1/locations/saved`
* **Access:** Authenticated (Bearer Token)
* **Request Body:**
  ```json
  {
    "location_id": "loc-67890"
  }
  ```

#### 7. Get All Saved Locations
* **Endpoint:** `GET /api/v1/locations/saved`
* **Access:** Authenticated (Bearer Token)

---

## 6. Environment Variables (`.env`)

```env
APP_PORT=8080
APP_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=go_vibe_weather_db
DB_SSLMODE=disable

# External Weather API Key
WEATHER_API_KEY=your_openweather_or_weatherapi_key
WEATHER_API_BASE_URL=https://api.weatherapi.com/v1

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key
JWT_EXPIRATION_HOURS=24
```

---

## 7. Setup & Run Instructions

1. **Clone Repositori & Install Dependencies:**
   ```bash
   git clone https://github.com/user/go_vibe.git
   cd go_vibe
   go mod download
   ```

2. **Konfigurasi Environment:**
   Salin `.env.example` ke `.env` lalu sesuaikan kredensial PostgreSQL dan API Key cuaca Anda.
   ```bash
   cp .env.example .env
   ```

3. **Jalankan Migration Database:**
   ```bash
   migrate -path migrations -database "postgres://postgres:secret@localhost:5432/go_vibe_weather_db?sslmode=disable" up
   ```

4. **Jalankan Aplikasi:**
   ```bash
   go run cmd/main.go
   ```