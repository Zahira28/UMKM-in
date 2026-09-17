# UMKM-in Frontend Client

Antarmuka web untuk platform UMKM-in, dibangun menggunakan Next.js (App Router), React 19, Tailwind CSS v4, dan TypeScript.

## Tech Stack

- Framework: Next.js (App Router)
- Library: React 19
- Styling: Tailwind CSS v4
- Language: TypeScript
- Code Quality: ESLint

## Struktur Direktori

```text
frontend/
├── app/                     # Halaman dan routing Next.js App Router
│   ├── globals.css          # Konfigurasi gaya global dan Tailwind
│   ├── layout.tsx           # Root layout aplikasi
│   └── page.tsx             # Halaman utama aplikasi
├── public/                  # Berkas statis publik (gambar, ikon, font)
├── src/                     # Komponen dan modul pendukung aplikasi
│   ├── components/          # Komponen antarmuka yang dapat digunakan kembali
│   │   ├── elements/        # Komponen elemen dasar
│   │   ├── layout/          # Komponen tata letak (navbar, sidebar, footer)
│   │   └── ui/              # Komponen antarmuka pengguna umum (button, modal, card)
│   ├── hooks/               # Custom React hooks
│   ├── lib/                 # Utilitas pustaka eksternal dan fungsi bantuan
│   └── modules/             # Modul logika berbasis fitur
├── eslint.config.mjs        # Konfigurasi linting ESLint
├── next.config.ts           # Konfigurasi Next.js
├── package.json             # Daftar pustaka dependensi dan skrip proyek
├── postcss.config.mjs       # Konfigurasi PostCSS untuk Tailwind CSS v4
└── tsconfig.json            # Konfigurasi TypeScript dan path alias (@/*)
```

## Cara Menjalankan

### 1. Prasyarat
- Node.js versi 20 ke atas.
- Paket manajer seperti npm, pnpm, atau yarn.

### 2. Instalasi Dependensi
Jalankan perintah berikut di dalam direktori frontend:
```bash
npm install
```

### 3. Konfigurasi Environment
Buat berkas .env.local di dalam direktori frontend untuk konfigurasi URL endpoint backend:
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
```

### 4. Menjalankan Server Pengembangan
```bash
npm run dev
```

Buka peramban di http://localhost:3000 untuk melihat tampilan aplikasi.

## Skrip yang Tersedia

- npm run dev : Menjalankan server pengembangan lokal dengan fitur hot reload.
- npm run build : Melakukan kompilasi dan optimasi aplikasi untuk tahap produksi.
- npm run start : Menjalankan server aplikasi hasil kompilasi produksi.
- npm run lint : Memeriksa standar penulisan kode menggunakan ESLint.
