# API Contract UMKM-in (v1)

Dokumen ini mendefinisikan spesifikasi resmi kontrak API antara Frontend (Next.js) dan Backend (Golang Fiber) untuk platform UMKM-in.

---

## 1. Konvensi Global

### A. Base URL & Versioning
Semua endpoint berada di bawah prefix versi 1:
```text
http://localhost:8080/api/v1
```

### B. Header Standar
* `Content-Type: application/json` (untuk semua request dengan payload JSON)
* `Authorization: Bearer <token_jwt>` (untuk semua endpoint yang memerlukan autentikasi)
* `Accept: application/json`

---

## 2. Standar Format Response (Envelope)

Semua response dari backend mengembalikan format JSON standar:

### A. Response Sukses Biasa (Single Object / Status)
```json
{
  "success": true,
  "message": "Operasi berhasil dilakukan",
  "data": { ... }
}
```

### B. Response Sukses Berisi Koleksi Data dengan Pagination
```json
{
  "success": true,
  "message": "Daftar data berhasil diambil",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 45,
    "total_pages": 5
  }
}
```

### C. Response Gagal / Error
```json
{
  "success": false,
  "message": "Pesan ringkasan error",
  "errors": [
    {
      "field": "email",
      "message": "Format email tidak valid"
    }
  ]
}
```
Jika tidak ada rincian per field, nilai `errors` dapat berupa `null`.

### D. Standar Kode Status HTTP
* `200 OK` : Request berhasil dan mengembalikan data.
* `201 Created` : Entitas baru berhasil dibuat (misal pendaftaran atau pembuatan produk).
* `400 Bad Request` : Input tidak valid, JSON rusak, atau validasi gagal.
* `401 Unauthorized` : Token JWT tidak disertakan, kedaluwarsa, atau tidak valid.
* `403 Forbidden` : Pengguna tidak memiliki hak akses (misal mengedit produk milik akun lain).
* `404 Not Found` : Data yang dicari tidak ditemukan.
* `409 Conflict` : Terjadi bentrokan data unik (misal email atau username sudah terdaftar).
* `500 Internal Server Error` : Terjadi kesalahan tak terduga pada server backend.

---

## 3. Modul 1: Autentikasi dan Akun (`/api/v1/auth`)

### 1.1 Register Akun Manual
* **Method:** `POST`
* **Path:** `/api/v1/auth/register`
* **Auth:** Tidak
* **Deskripsi:** Mendaftarkan pengguna baru dengan email dan kata sandi. Mengirimkan kode OTP 6 digit ke email dan mencatatnya di log terminal server.
* **Request Body:**
  ```json
  {
    "full_name": "Ayu Lestari",
    "email": "ayu@example.com",
    "password": "Password123!"
  }
  ```
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Pendaftaran berhasil. Silakan periksa email Anda untuk kode verifikasi OTP.",
    "data": {
      "email": "ayu@example.com",
      "expires_in_seconds": 600
    }
  }
  ```

### 1.2 Verifikasi Email dengan OTP
* **Method:** `POST`
* **Path:** `/api/v1/auth/verify-otp`
* **Auth:** Tidak
* **Deskripsi:** Memvalidasi 6 digit kode OTP. Jika berhasil, akun diaktifkan dan token JWT diterbitkan.
* **Request Body:**
  ```json
  {
    "email": "ayu@example.com",
    "otp_code": "749201"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Verifikasi email berhasil",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "email": "ayu@example.com",
        "full_name": "Ayu Lestari",
        "username": "",
        "is_verified": true,
        "is_onboarded": false
      }
    }
  }
  ```

### 1.3 Kirim Ulang OTP (Resend OTP)
* **Method:** `POST`
* **Path:** `/api/v1/auth/resend-otp`
* **Auth:** Tidak
* **Deskripsi:** Mengirim ulang kode verifikasi OTP baru dengan jeda minimal 60 detik.
* **Request Body:**
  ```json
  {
    "email": "ayu@example.com"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Kode OTP baru telah dikirim ke email Anda",
    "data": {
      "cooldown_seconds": 60
    }
  }
  ```

### 1.4 Login Akun Manual
* **Method:** `POST`
* **Path:** `/api/v1/auth/login`
* **Auth:** Tidak
* **Deskripsi:** Masuk menggunakan email atau username terdaftar beserta kata sandi.
* **Request Body:**
  ```json
  {
    "email_or_username": "ayu@example.com",
    "password": "Password123!"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Login berhasil",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "username": "ayu_lestari",
        "email": "ayu@example.com",
        "full_name": "Ayu Lestari",
        "avatar_url": "https://storage.umkm-in.id/avatars/ayu.jpg",
        "city": "Bandung",
        "is_verified": true,
        "ai_credits": 10
      }
    }
  }
  ```

### 1.5 Login via Google OAuth 2.0
* **Method:** `POST`
* **Path:** `/api/v1/auth/google`
* **Auth:** Tidak
* **Deskripsi:** Menerima ID token dari Google Identity Services SDK pada frontend. Akun yang login via Google otomatis berstatus terverifikasi.
* **Request Body:**
  ```json
  {
    "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjEyMy..."
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Autentikasi Google berhasil",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "email": "ayu.google@gmail.com",
        "full_name": "Ayu Lestari",
        "username": "ayu_lestari",
        "avatar_url": "https://lh3.googleusercontent.com/a/...",
        "is_verified": true,
        "is_onboarded": true
      }
    }
  }
  ```

### 1.6 Cek Sesi Akun yang Sedang Login
* **Method:** `GET`
* **Path:** `/api/v1/auth/me`
* **Auth:** Ya
* **Deskripsi:** Mendapatkan data profil akun pemilik token JWT yang aktif.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Data akun berhasil diambil",
    "data": {
      "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
      "username": "ayu_lestari",
      "email": "ayu@example.com",
      "full_name": "Ayu Lestari",
      "phone_number": "6281234567890",
      "bio": "Menjual aneka camilan rumahan tradisional khas Bandung.",
      "avatar_url": "https://storage.umkm-in.id/avatars/ayu.jpg",
      "city": "Bandung",
      "address_detail": "Jl. Cihampelas No. 12, Bandung",
      "latitude": -6.899201,
      "longitude": 107.604212,
      "ai_credits": 10,
      "created_at": "2026-09-24T10:00:00Z"
    }
  }
  ```

### 1.7 Lengkapi / Perbarui Profil (Onboarding)
* **Method:** `PUT`
* **Path:** `/api/v1/auth/profile`
* **Auth:** Ya
* **Deskripsi:** Melengkapi profil toko setelah pendaftaran atau mengubah informasi profil toko.
* **Request Body:**
  ```json
  {
    "username": "ayu_lestari",
    "full_name": "Ayu Lestari - Dapur Bunda",
    "phone_number": "6281234567890",
    "bio": "Menjual camilan rumahan segar setiap hari.",
    "avatar_url": "https://storage.umkm-in.id/avatars/ayu.jpg",
    "city": "Bandung, Jawa Barat",
    "address_detail": "Jl. Cihampelas No. 12",
    "latitude": -6.899201,
    "longitude": 107.604212
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Profil berhasil diperbarui",
    "data": { ...data user terbaru... }
  }
  ```

---

## 4. Modul 2: Profil Pengguna dan Relasi Follow (`/api/v1/users`)

### 2.1 Lihat Profil Publik Pengguna / Toko
* **Method:** `GET`
* **Path:** `/api/v1/users/:username`
* **Auth:** Opsional (jika membawa token, status `is_following` akan dievaluasi)
* **Deskripsi:** Mengambil informasi etalase toko pengguna beserta ringkasan jumlah produk, followers, dan following.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Profil toko berhasil diambil",
    "data": {
      "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
      "username": "ayu_lestari",
      "full_name": "Dapur Ayu Bandung",
      "bio": "Kue rumahan dan camilan tradisional.",
      "avatar_url": "https://storage.umkm-in.id/avatars/ayu.jpg",
      "city": "Bandung, Jawa Barat",
      "phone_number": "6281234567890",
      "products_count": 8,
      "followers_count": 142,
      "following_count": 19,
      "is_following": false
    }
  }
  ```

### 2.2 Toggle Follow / Unfollow Pengguna
* **Method:** `POST`
* **Path:** `/api/v1/users/:username/follow`
* **Auth:** Ya
* **Deskripsi:** Mengikuti akun toko atau berhenti mengikuti.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Berhasil mengikuti pengguna",
    "data": {
      "is_following": true,
      "followers_count": 143
    }
  }
  ```

### 2.3 Daftar Pengikut (Followers)
* **Method:** `GET`
* **Path:** `/api/v1/users/:username/followers?page=1&limit=20`
* **Auth:** Tidak
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Daftar followers berhasil diambil",
    "data": [
      {
        "id": "b1fe...",
        "username": "budi_santoso",
        "full_name": "Budi Santoso",
        "avatar_url": "https://..."
      }
    ],
    "meta": { "page": 1, "limit": 20, "total_items": 143, "total_pages": 8 }
  }
  ```

### 2.4 Daftar Pengguna yang Diikuti (Following)
* **Method:** `GET`
* **Path:** `/api/v1/users/:username/following?page=1&limit=20`
* **Auth:** Tidak
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Daftar following berhasil diambil",
    "data": [ ... ],
    "meta": { "page": 1, "limit": 20, "total_items": 19, "total_pages": 1 }
  }
  ```

---

## 5. Modul 3: Kategori dan Produk (`/api/v1/categories`, `/api/v1/products`)

### 3.1 Daftar Kategori
* **Method:** `GET`
* **Path:** `/api/v1/categories`
* **Auth:** Tidak
* **Deskripsi:** Mengambil seluruh kategori produk untuk filter dropdown.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Daftar kategori berhasil diambil",
    "data": [
      { "id": 1, "name": "Makanan", "slug": "makanan" },
      { "id": 2, "name": "Minuman", "slug": "minuman" },
      { "id": 3, "name": "Pakaian", "slug": "pakaian" },
      { "id": 4, "name": "Kriya", "slug": "kriya" },
      { "id": 5, "name": "Jasa", "slug": "jasa" }
    ]
  }
  ```

### 3.2 Linimasa / Feed Produk Publik (Hyperlocal Discovery)
* **Method:** `GET`
* **Path:** `/api/v1/products`
* **Auth:** Opsional (jika membawa token, status `is_liked` produk akan disesuaikan)
* **Query Parameters:**
  * `category` (string, opsional): Slug kategori (misal `makanan`).
  * `search` (string, opsional): Kata kunci pencarian judul atau deskripsi produk.
  * `sort` (string, opsional): `latest` (default), `cheapest`, `priciest`, `nearest`.
  * `latitude` (float, opsional): Titik latitude pengguna untuk sorting jarak terdekat.
  * `longitude` (float, opsional): Titik longitude pengguna.
  * `page` (int, default: 1)
  * `limit` (int, default: 10)
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Feed produk berhasil diambil",
    "data": [
      {
        "id": "e2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c",
        "name": "Kue Gula Aren Rumahan",
        "description": "Kue gula aren buatan rumahan dengan bahan segar setiap hari.",
        "price": 35000,
        "unit": "box",
        "raw_image_url": "https://storage.umkm-in.id/products/raw-1.jpg",
        "studio_image_url": "https://storage.umkm-in.id/products/studio-1.jpg",
        "is_available": true,
        "created_at": "2026-09-24T08:30:00Z",
        "likes_count": 24,
        "comments_count": 5,
        "is_liked": false,
        "category": {
          "id": 1,
          "name": "Makanan",
          "slug": "makanan"
        },
        "seller": {
          "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
          "username": "ayu_lestari",
          "full_name": "Dapur Ayu Bandung",
          "avatar_url": "https://storage.umkm-in.id/avatars/ayu.jpg",
          "city": "Bandung, Jawa Barat",
          "phone_number": "6281234567890"
        }
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 52,
      "total_pages": 6
    }
  }
  ```

### 3.3 Detail Produk
* **Method:** `GET`
* **Path:** `/api/v1/products/:id`
* **Auth:** Opsional
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Detail produk berhasil diambil",
    "data": { ...detail lengkap objek produk seperti pada feed... }
  }
  ```

### 3.4 Buat Produk Baru
* **Method:** `POST`
* **Path:** `/api/v1/products`
* **Auth:** Ya
* **Request Body:**
  ```json
  {
    "category_id": 1,
    "name": "Kue Gula Aren Rumahan",
    "description": "Kue gula aren buatan rumahan dengan bahan pilihan segar setiap hari.",
    "price": 35000,
    "unit": "box",
    "raw_image_url": "https://storage.umkm-in.id/products/raw-1.jpg",
    "studio_image_url": "https://storage.umkm-in.id/products/studio-1.jpg",
    "is_available": true
  }
  ```
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Produk berhasil ditambahkan ke etalase",
    "data": {
      "id": "e2a3b4c5-d6e7-4f8a-9b0c-1d2e3f4a5b6c",
      "name": "Kue Gula Aren Rumahan",
      "price": 35000,
      "is_available": true,
      "created_at": "2026-09-24T10:30:00Z"
    }
  }
  ```

### 3.5 Perbarui Data Produk
* **Method:** `PUT`
* **Path:** `/api/v1/products/:id`
* **Auth:** Ya (hanya pemilik produk)
* **Request Body:**
  ```json
  {
    "category_id": 1,
    "name": "Kue Gula Aren Premium",
    "description": "Deskripsi baru produk yang diperbarui.",
    "price": 38000,
    "unit": "box",
    "studio_image_url": "https://...",
    "is_available": true
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Produk berhasil diperbarui",
    "data": { ...data produk terbaru... }
  }
  ```

### 3.6 Hapus Produk
* **Method:** `DELETE`
* **Path:** `/api/v1/products/:id`
* **Auth:** Ya (hanya pemilik produk)
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Produk berhasil dihapus dari etalase",
    "data": null
  }
  ```

---

## 6. Modul 4: Komentar dan Reaksi / Likes (`/api/v1/products/:id/...`)

### 4.1 Daftar Komentar Produk (Mendukung Balasan Bertingkat)
* **Method:** `GET`
* **Path:** `/api/v1/products/:id/comments`
* **Auth:** Tidak
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Daftar komentar berhasil diambil",
    "data": [
      {
        "id": 101,
        "content": "Bisa kirim ke area Buah Batu Bandung nggak kak?",
        "created_at": "2026-09-24T09:00:00Z",
        "user": {
          "id": "c1...",
          "username": "siti_bandung",
          "full_name": "Siti Rahma",
          "avatar_url": "https://..."
        },
        "replies": [
          {
            "id": 102,
            "parent_id": 101,
            "content": "Bisa kak, kami pakai kurir instan!",
            "created_at": "2026-09-24T09:15:00Z",
            "user": {
              "id": "a0eebc99...",
              "username": "ayu_lestari",
              "full_name": "Dapur Ayu Bandung",
              "avatar_url": "https://..."
            }
          }
        ]
      }
    ]
  }
  ```

### 4.2 Kirim Komentar Baru atau Balasan
* **Method:** `POST`
* **Path:** `/api/v1/products/:id/comments`
* **Auth:** Ya
* **Request Body:**
  ```json
  {
    "content": "Bisa dikirim hari ini kak?",
    "parent_id": null
  }
  ```
  *(Isi `parent_id` dengan id komentar utama jika ingin membalas komentar).*
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Komentar berhasil dikirim",
    "data": {
      "id": 103,
      "product_id": "e2a3...",
      "content": "Bisa dikirim hari ini kak?",
      "parent_id": null,
      "created_at": "2026-09-24T10:45:00Z"
    }
  }
  ```

### 4.3 Hapus Komentar
* **Method:** `DELETE`
* **Path:** `/api/v1/comments/:comment_id`
* **Auth:** Ya (pembuat komentar atau pemilik produk)
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Komentar berhasil dihapus",
    "data": null
  }
  ```

### 4.4 Toggle Like / Reaction Produk
* **Method:** `POST`
* **Path:** `/api/v1/products/:id/like`
* **Auth:** Ya
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Reaksi berhasil diperbarui",
    "data": {
      "is_liked": true,
      "likes_count": 25
    }
  }
  ```

---

## 7. Modul 5: Direct-to-WhatsApp Order Generator (`/api/v1/products/:id/whatsapp-order`)

### 7.1 Generate Link WhatsApp Terformat Otomatis
* **Method:** `GET`
* **Path:** `/api/v1/products/:id/whatsapp-order`
* **Auth:** Tidak
* **Deskripsi:** Mengambil data produk dan nomor WhatsApp toko untuk menyusun URL WhatsApp API (`https://wa.me/...`) dengan pesan pembuka otomatis yang sopan dan jelas.
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Link WhatsApp order berhasil dibuat",
    "data": {
      "phone_number": "6281234567890",
      "seller_name": "Dapur Ayu Bandung",
      "product_name": "Kue Gula Aren Rumahan",
      "price": 35000,
      "whatsapp_url": "https://wa.me/6281234567890?text=Halo%20kak%2C%20saya%20tertarik%20dengan%20produk%20*Kue%20Gula%20Aren%20Rumahan*%20seharga%20Rp%2035.000%20di%20UMKM-in.%20Apakah%20stok%20masih%20tersedia%3F",
      "draft_message": "Halo kak, saya tertarik dengan produk *Kue Gula Aren Rumahan* seharga Rp 35.000 di UMKM-in. Apakah stok masih tersedia?"
    }
  }
  ```

---

## 8. Modul 6: Media Upload dan Object Storage (`/api/v1/upload`)

### 8.1 Upload Berkas Gambar (Foto Profil / Produk)
* **Method:** `POST`
* **Path:** `/api/v1/upload`
* **Auth:** Ya
* **Content-Type:** `multipart/form-data`
* **Form Field:** `file` (berkas gambar, format: `.jpg`, `.jpeg`, `.png`, `.webp`, ukuran maksimal 5 MB)
* **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Berkas gambar berhasil diunggah",
    "data": {
      "url": "https://storage.umkm-in.id/uploads/prod_1727170800_abc123.webp",
      "filename": "prod_1727170800_abc123.webp",
      "file_size": 348120,
      "content_type": "image/webp"
    }
  }
  ```

---

## 9. Modul 7: AI Creative Studio dan Copywriting (`/api/v1/ai`)

### 9.1 Smart Copywriter Generator
* **Method:** `POST`
* **Path:** `/api/v1/ai/generate-copywriting`
* **Auth:** Ya
* **Deskripsi:** Membuat deskripsi produk persuasif, ringkasan keunggulan, dan rekomendasi tagar otomatis. Mengurangi kuota `ai_credits` pengguna sebanyak 1 poin.
* **Request Body:**
  ```json
  {
    "title": "Kue Gula Aren Rumahan",
    "category": "Makanan",
    "vibe": "Rustic Wood",
    "custom_notes": "Tanpa pengawet, dibuat dari gula aren murni petani lokal"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Copywriting berhasil dibuat oleh AI",
    "data": {
      "caption": "Kue Gula Aren buatan rumahan dengan bahan pilihan segar setiap hari. Dipanggang pas, manisnya gula aren alami tidak bikin enek! Pas untuk camilan santai keluarga.",
      "selling_points": [
        "100 persen gula aren murni tanpa pemanis buatan",
        "Dipanggang segar setiap pagi",
        "Kemasan ramah lingkungan"
      ],
      "hashtags": [
        "#makananlokal",
        "#kuegulaaren",
        "#umkmjuara",
        "#kulinerbandung"
      ],
      "remaining_ai_credits": 9
    }
  }
  ```

### 9.2 AI Background Inpainting / Studio Photo
* **Method:** `POST`
* **Path:** `/api/v1/ai/studio-inpainting`
* **Auth:** Ya
* **Deskripsi:** Menghapus latar belakang kotor dari foto kamera ponsel mentah dan menggantinya dengan latar studio komersial sesuai preset vibe.
* **Request Body:**
  ```json
  {
    "raw_image_url": "https://storage.umkm-in.id/products/raw-1.jpg",
    "vibe_preset": "Rustic Wood"
  }
  ```
* **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Foto produk studio AI berhasil digenerate",
    "data": {
      "studio_image_url": "https://storage.umkm-in.id/products/studio-1.jpg",
      "remaining_ai_credits": 8
    }
  }
  ```
