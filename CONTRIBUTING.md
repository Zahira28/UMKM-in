# Panduan Kontribusi (Contributing Guidelines) UMKM-in

Selamat datang di proyek pengembangan **UMKM-in**! Untuk menjaga kerapian, kualitas kode, dan kelancaran kolaborasi tim, seluruh anggota tim diwajibkan mengikuti panduan kerja berikut.

---

## 1. Aturan Penamaan Branch (Branching Strategy)

Dilarang membuat branch dengan nama yang ambigu. Seluruh branch kerja harus mengikuti konvensi berikut:

```text
<tipe>/<deskripsi-fitur-singkat>
```

### Kategori Branch:
* **`feat/`** : Untuk penambahan fitur baru (contoh: `feat/github-actions`, `feat/auth-login`, `feat/erd-schema`).
* **`fix/`** : Untuk perbaikan bug atau error (contoh: `fix/whatsapp-button-link`, `fix/feed-sorting`).
* **`docs/`** : Untuk pembaruan berkas dokumentasi (contoh: `docs/update-readme`, `docs/lab-2.4`).
* **`chore/`** : Untuk pemeliharaan dependensi atau konfigurasi (contoh: `chore/setup-linter`).

---

## 2. Aturan Pesan Commit (Conventional Commits)

Format pesan commit wajib diawali dengan jenis perubahan dalam huruf kecil:

```text
<tipe>: <deskripsi singkat perubahan>
```

### Contoh:
* `feat: add toggle availability on product`
* `fix: handle null gps coordinates in feed`
* `docs: update lab 2.4 assignment report`
* `ci: configure github actions workflow for lab 2.5`

---

## 3. Alur Kerja Pull Request (PR)

1. **Dilarang Push Langsung ke Branch `main`**:
   Branch `main` harus selalu bersih, stabil, dan siap pakai.
2. **Buat Branch Baru dari `main`**:
   ```bash
   git checkout main
   git pull origin main
   git checkout -b feat/nama-fitur
   ```
3. **Commit dan Push ke Branch Baru**:
   ```bash
   git add .
   git commit -m "feat: deskripsi perubahan"
   git push -u origin feat/nama-fitur
   ```
4. **Buka Pull Request (PR)**:
   * Target branch tujuan selalu mengarah ke `main`.
   * Tuliskan ringkasan perubahan yang jelas pada deskripsi PR.
   * Tautkan ke Issue terkait jika ada (contoh: `Closes #1`).
5. **Verifikasi CI Pipeline**:
   Pastikan pengujian otomatis GitHub Actions berhasil (status centang hijau) sebelum branch di-*merge*.
