# Celana — Cellix Laundry Management System

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/SvelteKit-2.0+-FF3E00?style=for-the-badge&logo=svelte&logoColor=white" alt="SvelteKit" />
  <img src="https://img.shields.io/badge/Tailwind_CSS-3.4+-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white" alt="Tailwind CSS" />
  <img src="https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" alt="License" />
</p>

<p align="center">
  Sistem manajemen laundry modern, cepat, dan terisolasi berbasis <b>Lightweight Multi-Tenant SaaS</b>. Dirancang khusus untuk mempermudah operasional bisnis laundry mulai dari kasir POS, pencatatan multi-item kiloan/satuan/meteran, CRM pelanggan, arus kas, hingga notifikasi otomatis via WhatsApp.
</p>

---

## Fitur Utama

### 1. Arsitektur Multi-Tenant Terisolasi

- **Row-Level Tenant Isolation**: Setiap tenant/pemilik laundry memiliki data yang terisolasi penuh via `user_id`.
- **JWT Context Anti-Spoofing**: Autentikasi dan identitas tenant diekstrak langsung dari token JWT yang terverifikasi, bukan dari parameter request.
- **Strict Cross-Tenant Validation**: Mencegah akses silang data antar-bisnis (mengembalikan `404 Not Found` alih-alih `403` untuk menjaga keamanan privasi data).

### 2. Manajemen Pesanan & Kasir Cepat

- **Multi-Item Layanan**: Mendukung kombinasi berbagai jenis layanan dalam satu nota (Per Kg, Per Satuan/Item, dan Per Meter Persegi / $m^2$).
- **Dropdown Pencarian Cepat**: Pencarian instan untuk pelanggan dan item layanan dengan keyboard navigation.
- **Sistem Pembayaran Fleksibel**: Dukungan pembayaran Tunai, Transfer Bank, QRIS, opsi Uang Muka (DP), pelunasan bertahap, serta kalkulator uang pas & kembalian otomatis.
- **Status Siklus Pesanan**: Pelacakan status terstruktur (_QUEUE, PROCESS, READY, COMPLETED, CANCELLED_).
- **Notifikasi WhatsApp Otomatis**: Integrasi pengiriman pesan struk/nota digital otomatis ke nomor WhatsApp pelanggan saat pesanan siap diambil (_READY_).

### 3. Katalog Layanan Dinamis

- Manajemen tarif laundry fleksibel dengan tipe satuan (`PER_KG`, `PER_ITEM`, `PER_M2`).
- Fitur aktifkan/nonaktifkan layanan tanpa menghapus riwayat transaksi lama.
- Statistik katalog dan ringkasan harga per satuan.

### 4. Manajemen Pelanggan (CRM)

- Pencatatan buku kontak pelanggan lengkap dengan alamat dan catatan preferensi laundry.
- Pelacakan riwayat pesanan pelanggan dan status tagihan belum lunas.

### 5. Arus Kas & Keuangan (Cash Flow)

- Rekap otomatis pemasukan dari transaksi laundry (DP & Pelunasan).
- Pencatatan pengeluaran operasional (deterjen, listrik, parfum, plastik packaging, dll).
- Ringkasan kas real-time: Total Pemasukan, Total Pengeluaran, dan Saldo Bersih.

### 6. Pengalaman Pengguna (UI/UX) Unggulan

- **Responsive & Modern Design**: Dibangun dengan Tailwind CSS dengan layout adaptif untuk layar desktop maupun smartphone/tablet kasir.
- **Dark Mode Support**: Pilihan tema Terang (_Light_), Gelap (_Dark_), dan Otomatis (_System_).
- **Interactive Product Tour**: Panduan interaktif terintegrasi menggunakan `driver.js` untuk memandu staf/kasir baru memahami fitur aplikasi langkah demi langkah.
- **Zero-Flicker In-Memory Caching**: Navigasi antar-halaman instan tanpa jeda _loading_ berulang.

---

## Tech Stack

### Backend

- **Bahasa**: [Go (Golang)](https://go.dev/) 1.22+
- **Framework Web**: [Fiber v2](https://gofiber.io/) — Framework web performa tinggi berbasis Fasthttp
- **Database Driver**: [pgx / lib/pq](https://github.com/lib/pq)
- **Keamanan**: JWT (_golang-jwt_), Argon2 / Bcrypt Password Hashing, CORS Protection

### Frontend

- **Framework**: [SvelteKit](https://kit.svelte.dev/) + Vite
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Icons**: [Lucide Svelte](https://lucide.dev/)
- **Tour Guide**: [driver.js](https://driverjs.com/)

### Database

- **DBMS**: [PostgreSQL](https://www.postgresql.org/) 15+

---

## Panduan Memulai (Getting Started)

### 1. Prasyarat

Pastikan Anda telah menginstal:

- [Go](https://go.dev/dl/) (versi 1.22 atau terbaru)
- [Node.js](https://nodejs.org/) (versi 18 atau 20 LTS)
- [PostgreSQL](https://www.postgresql.org/download/)

---

### 2. Konfigurasi Database & Backend

1. Buka terminal dan masuk ke folder `backend`:

   ```bash
   cd backend
   ```

2. Salin template file konfigurasi `.env`:

   ```bash
   cp .env.example .env
   ```

3. Sesuaikan konfigurasi pada file `.env`:

   ```env
   PORT=8080
   DATABASE_URL=postgres://postgres:password_anda@localhost:5432/celana_db?sslmode=disable
   JWT_SECRET=rahasia-kunci-jwt-super-aman
   CORS_ORIGIN=http://localhost:5173
   ```

4. Jalankan migrasi database PostgreSQL:

   ```bash
   # Buat database jika belum ada
   psql -U postgres -c "CREATE DATABASE celana_db;"

   # Jalankan skrip migrasi
   psql -U postgres -d celana_db -f migrations/001_init_multi_tenant.sql
   ```

5. Unduh dependencies dan jalankan backend:
   ```bash
   go mod tidy
   go run cmd/api/main.go
   ```
   _Server backend akan berjalan di `http://localhost:8080`._

---

### 3. Konfigurasi Frontend

1. Buka terminal baru dan masuk ke folder `frontend`:

   ```bash
   cd frontend
   ```

2. Pasang dependensi Node.js:

   ```bash
   npm install
   ```

3. Jalankan server pengembangan SvelteKit:
   ```bash
   npm run dev
   ```
   _Aplikasi frontend dapat diakses melalui browser di `http://localhost:5173`._

---

## Struktur Direktori

```text
celana/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go           # Entry point server Fiber
│   ├── internal/
│   │   ├── config/               # Konfigurasi environment & DB
│   │   ├── handlers/             # Controller (Auth, Orders, Customers, Finance, dll)
│   │   ├── middleware/           # Auth JWT & Tenant Isolation
│   │   ├── models/               # Definisi entitas data & DTO
│   │   └── routes/               # Routing API endpoint
│   ├── migrations/
│   │   └── 001_init_multi_tenant.sql
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/       # Komponen reusable (Sidebar, Modal, Badge, Pagination)
│   │   │   ├── stores/           # Svelte stores (Auth, Theme, In-Memory Cache)
│   │   │   ├── api.js            # Fetch wrapper terpusat
│   │   │   └── tour.js           # Konfigurasi panduan driver.js
│   │   └── routes/
│   │       ├── +page.svelte      # Dashboard utama
│   │       ├── pesanan/          # Modul transaksi & kasir
│   │       ├── layanan/          # Modul katalog paket laundry
│   │       ├── pelanggan/        # Modul CRM pelanggan
│   │       ├── keuangan/         # Modul arus kas pemasukan & pengeluaran
│   │       ├── pengaturan/       # Modul profil bisnis & outlet
│   │       ├── pengguna/         # Manajemen pengguna & role (Admin)
│   │       ├── login/            # Autentikasi masuk
│   │       └── register/         # Registrasi tenant baru
│   ├── package.json
│   ├── svelte.config.js
│   ├── tailwind.config.js
│   └── vite.config.js
│
├── .gitignore
└── README.md
```

---

## Lisensi

Proyek ini dilisensikan di bawah [MIT License](LICENSE).
Dikembangkan untuk mendukung digitalisasi dan efisiensi operasional UMKM Laundry.
