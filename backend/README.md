# UMKM-in Backend Service

Layanan RESTful API backend untuk platform UMKM-in, dibangun menggunakan Golang dengan framework Fiber v2 dan menerapkan pola Clean atau Layered Architecture.

## Tech Stack

- Language: Go (v1.22+)
- Web Framework: Fiber v2
- ORM and Database: GORM and PostgreSQL
- Live Reload: Air
- Configuration: godotenv

## Struktur Direktori

```text
backend/
├── cmd/
│   └── api/
│       └── main.go          # Entry point utama aplikasi backend
├── internal/                # Kode privat aplikasi
│   ├── config/              # Pembacaan konfigurasi environment (.env)
│   ├── database/            # Koneksi dan migrasi basis data
│   ├── handler/             # HTTP controller dan pendaftaran endpoint
│   ├── middleware/          # Middleware (CORS, Logger, JWT Auth)
│   ├── model/               # Definisi struct entitas tabel basis data
│   ├── repository/          # Layer query basis data (CRUD)
│   └── service/             # Layer logika bisnis (business logic)
├── pkg/                     # Helper dan utilitas publik (response formatter, dll)
├── tmp/                     # Temporary binary build dari Air (diabaikan oleh git)
├── .air.toml                # Konfigurasi hot-reload Air
├── .env.example             # Template variabel environment
├── .gitignore               # Konfigurasi pengabaian berkas git
├── go.mod                   # Definisi modul Go dan dependensi
└── go.sum                   # Checksum dependensi Go
```

## Cara Menjalankan

### 1. Prasyarat
- Go terinstal di sistem (versi 1.22 ke atas).
- PostgreSQL aktif di lokal atau melalui container Docker.
- Air terinstal jika ingin menggunakan fitur hot-reload.

### 2. Konfigurasi Environment
Salin berkas template .env.example menjadi .env:
```bash
cp .env.example .env
```
Sesuaikan konfigurasi port dan kredensial database yang digunakan.

### 3. Unduh Dependensi
```bash
go mod tidy
```

### 4. Menjalankan Server

- Mode Pengembangan dengan Hot-Reload (Air):
  ```bash
  air
  ```

- Mode Standar (tanpa Air):
  ```bash
  go run cmd/api/main.go
  ```

Server akan berjalan secara default pada port 8080 (http://localhost:8080).

## Health Check Endpoint

Untuk memverifikasi bahwa server backend telah aktif:
- Method: GET
- URL: http://localhost:8080/api/health
- Format Response:
  ```json
  {
    "status": "success",
    "message": "UMKM-in backend service is running"
  }
  ```
