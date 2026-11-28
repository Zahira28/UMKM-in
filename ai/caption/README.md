# AI Caption: Improve Deskripsi Produk

Prompt dan skrip pengujian untuk fitur **"Improve dengan AI"** di UMKM-in. Penjual menulis deskripsi produk seadanya, lalu AI merapikannya dan menyarankan hashtag, tanpa menambah fakta yang tidak ditulis penjual.

- **Model:** Gemini 3.5 Flash Lite (`gemini-3.5-flash-lite`), tier gratis Google AI Studio
- **Versi prompt aktif:** v2 (`prompts/system_instructions_v2.md`)
- **Issue terkait:** #26

> Folder ini berisi prompt dan alat uji, bukan model. Model dipanggil lewat API Google.

## Struktur Folder

```text
ai/caption/
├── prompts/
│   ├── system_instructions_v1.md   # versi awal (arsip)
│   ├── system_instructions_v2.md   # versi aktif
│   └── prompt_config.md            # konfigurasi model, aturan, dan riwayat versi
├── dataset/
│   └── test_cases.json             # 20 kasus uji
├── scripts/
│   └── run_tests.py                # menjalankan semua kasus uji ke Gemini API
├── evaluation/
│   └── model_testing_results.xlsx  # hasil evaluasi
├── requirements.txt
├── .env.example
└── README.md
```

## Kontrak Input dan Output

**Input** (JSON):

| Field | Tipe | Keterangan |
|---|---|---|
| `name` | string | Nama produk |
| `category` | string | Kategori produk |
| `price` | number | Harga produk |
| `description` | string | Deskripsi awal dari penjual |

**Output** (JSON, hanya dua field):

```json
{
  "description": "Cookies dengan gula aren yang renyah dan enak banget.",
  "hashtags": ["#CookiesGulaAren", "#CookiesRenyah", "#GulaArenRenyah"]
}
```

- `hashtags` berisi 3-7 item dan **sudah berawalan `#`**, jadi bisa langsung dipakai di frontend.
- Output selalu JSON valid, tanpa Markdown atau teks di luar JSON.

## Aturan Utama Prompt
Daftar lengkap ada di `prompts/prompt_config.md`

## Cara Menjalankan Tes

1. **Pasang dependensi** (disarankan di dalam virtual environment):

   ```powershell
   py -m venv .venv
   .\.venv\Scripts\Activate.ps1
   python -m pip install -r requirements.txt
   ```

2. **Siapkan API key.** Salin `.env.example` menjadi `.env`, lalu isi:

   ```text
   GEMINI_API_KEY=isi_key
   ```

3. **Jalankan pengujian:**

   ```powershell
   python scripts\run_tests.py
   ```

   Hasil pengujian berupa file Excel (lihat folder `evaluation/`).

## Hasil Pengujian

Pengujian dilakukan pada 20 kasus yang mencakup bahasa informal, deskripsi kosong, klaim berlebihan, informasi bertentangan, dan deskripsi yang sudah bagus.

| Prompt | File |
|---|---|---|
| 1 | v2 | `evaluation/model_testing_results.xlsx` |

## Keterbatasan yang Diketahui

- Klaim kesehatan kadang tidak dihapus, hanya dilunakkan (mis. "membantu meningkatkan imun").
- Saat dua informasi bertentangan, model terkadang memilih salah satu sisi.
- Hashtag kadang menempelkan kata kategori ke nama produk sehingga terdengar mekanis.
- Hasil bisa sedikit berbeda tiap pemanggilan karena temperature 1.0.
- Aturan prompt hanya mengurangi risiko, bukan jaminan. Karena platform tidak memverifikasi klaim penjual, sebaiknya form produk memuat pernyataan bahwa penjual bertanggung jawab atas kebenaran informasi.

## Catatan Integrasi

**Backend**

- Panggil Gemini dari backend, bukan dari frontend. API key disimpan di environment variable `GEMINI_API_KEY`.
- Kirim ke AI hanya `name`, `category`, `price`, dan `description`. Jangan kirim nomor telepon atau alamat penjual.
- Kurangi `AICredits` hanya jika pemanggilan berhasil.
- Tangani error `429` (batas request) dengan retry atau pesan yang jelas ke pengguna.

**Frontend**

- Tampilkan hasil AI sebagai pratinjau dengan pilihan **Pakai / Buat ulang / Batal**, agar tulisan penjual tidak langsung tertimpa.
- Aktifkan tombol "Improve dengan AI" hanya jika deskripsi sudah terisi.

## Riwayat Versi

| Versi | Tanggal | Perubahan |
|---|---|---|
| 1.0 | 10 November 2026 | Prompt awal |
| 2.0 | 20 November 2026 | Klaim penjual dipertahankan tanpa diperkuat, aturan informasi bertentangan, pembatasan hashtag, dan gaya bahasa dibatasi sesuai input |