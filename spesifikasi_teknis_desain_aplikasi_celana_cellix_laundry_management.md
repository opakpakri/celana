# Spesifikasi Sistem: Celana (Cellix Laundry Management)

Dokumen ini berisi rancangan teknis, desain antarmuka (_UI/UX Design System_), arsitektur **Multi-Tenant**, skema database PostgreSQL, dan spesifikasi API untuk pengembangan aplikasi **Celana (Cellix Laundry Management)** menggunakan **Go (Golang)**, **Svelte**, **Tailwind CSS**, dan **PostgreSQL**.

## 1. Identitas & Tech Stack

- **Nama Aplikasi:** Celana (_Cellix Laundry Management_)
- **Arsitektur Sistem:** Lightweight Multi-Tenant SaaS (Row-Level Tenant Isolation via `user_id`) — Ringan, modern, dan responsif untuk desktop maupun perangkat mobile/tablet kasir.
- **Backend:** Go (Golang) — Menggunakan _framework_ **Fiber** atau **Echo** dengan JWT Authentication & Tenant Isolation Middleware.
- **Frontend:** SvelteKit / Svelte + Vite — SPA/SSR yang ringan, cepat, dan reaktif.
- **Styling:** Tailwind CSS — Utility-first CSS dengan pendekatan _mobile-first responsive design_.
- **Database:** PostgreSQL — Relational database dengan indeks komposit `(user_id, ...)` untuk menjamin kecepatan query dan isolasi data.
- **Icons:** Lucide Icons (`lucide-svelte`) — Sesuai dengan gaya ikon garis tipis (_outline_) pada referensi desain.

## 2. Arsitektur Multi-Tenant & Keamanan Isolasi Data

Setiap akun pada tabel `users` merepresentasikan satu **Tenant / Pemilik Bisnis Laundry** yang independen. Untuk menjamin **User A tidak dapat melihat, mengakses, maupun mengubah data milik User B dalam kondisi apa pun**, diterapkan 4 lapis pengamanan berikut:

1. **Foreign Key `user_id` Wajib di Semua Tabel:** Setiap tabel operasional (`customers`, `services`, `orders`, `cash_flows`) wajib memiliki kolom `user_id NOT NULL REFERENCES users(id) ON DELETE CASCADE`.
2. **JWT Context Extraction (Anti-Spoofing):** Backend **dilarang** menerima `user_id` dari payload JSON request atau query parameter URL. `user_id` harus diekstrak secara eksklusif dari _claims_ token JWT yang sudah terverifikasi oleh middleware autentikasi.
3. **Strict Query Scoping:** Seluruh query database (`SELECT`, `UPDATE`, `DELETE`) wajib menyertakan klausa `WHERE user_id = $1`. Jika User A mencoba mengakses `/api/orders/99` milik User B, sistem akan mengembalikan `404 Not Found` (bukan `403`, agar tidak membocorkan eksistensi ID milik tenant lain).
4. **Cross-Tenant Reference Validation:** Saat membuat pesanan (`POST /api/orders`), backend wajib memvalidasi bahwa `customer_id` dan seluruh `service_id` yang dipilih benar-benar dimiliki oleh `user_id` yang sedang login.

## 3. Design System & Palet Warna (Berdasarkan Referensi UI)

Tampilan **Celana** menggunakan gaya **Clean SaaS / Minimalist Light Mode** dengan sudut membulat (_rounded corners_), batas tipis (_subtle borders_), aksen gradien biru-teal pada kartu utama, serta tata letak responsif (_collapsible drawer sidebar_ di layar mobile/tablet).

### Palet Warna (Tailwind CSS Mapping)

| Elemen UI                               | Warna / Efek         | Kelas Tailwind CSS / Kode Hex                                 |
| :-------------------------------------- | :------------------- | :------------------------------------------------------------ |
| **Brand Logo & Accent**                 | Biru Cerah           | `text-[#0284C7]` (`sky-600`)                                  |
| **Background Utama**                    | Abu-abu sangat muda  | `bg-[#F8FAFC]` (`slate-50`)                                   |
| **Sidebar & Card**                      | Putih bersih         | `bg-white` (`#FFFFFF`)                                        |
| **Border & Garis Pembatas**             | Abu-abu lembut       | `border-[#E2E8F0]` (`slate-200`)                              |
| **Primary Action (Tombol)**             | Biru Cerah           | `bg-[#0284C7]` hover `bg-[#0369A1]` (`sky-600` / `#007BFF`)   |
| **Active Menu Sidebar**                 | Biru Muda Lembut     | `bg-[#EFF6FF]` teks `text-[#1D4ED8]` (`blue-50` & `blue-700`) |
| **Hero Banner Dashboard**               | Gradien Biru ke Teal | `bg-gradient-to-r from-[#0284C7] via-[#0891B2] to-[#0D9488]`  |
| **Teks Utama (Heading)**                | Hitam Slate          | `text-[#0F172A]` (`slate-900`)                                |
| **Teks Sekunder (Subjudul)**            | Abu-abu Medium       | `text-[#64748B]` (`slate-500`)                                |
| **Status: Aktif / Lunas / Selesai**     | Hijau Lembut         | `bg-[#DCFCE7] text-[#15803D]` (`green-100` & `green-700`)     |
| **Status: Tipe Layanan (Per Kg/Item)**  | Biru Lembut          | `bg-[#E0F2FE] text-[#0369A1]` (`sky-100` & `sky-700`)         |
| **Status: Peringatan / Belum Lunas**    | Merah Lembut         | `bg-[#FEE2E2] text-[#B91C1C]` (`red-100` & `red-700`)         |
| **Tombol Destructive (Hapus/Nonaktif)** | Merah Solid          | `bg-[#EF4444]` hover `bg-[#DC2626]` (`red-500`)               |

## 4. Struktur Navigasi & Fitur Utama (6 Menu Sidebar)

1. **Dashboard**
   - **Branding Header Sidebar:** Logo ikon + teks **Celana** (_Cellix Laundry Management_) dan kartu info bisnis aktif milik tenant yang login.
   - **Hero Card Gradien:** Sapaan personal nama user + 4 ringkasan metrik hari ini (_Pesanan Hari Ini_, _Pendapatan_, _Dalam Proses_, _Siap Diambil_).
   - **Aksi Cepat (Quick Actions):** Kartu pintasan menuju _Buat Pesanan Baru_ dan _Lihat Semua Pesanan_.
   - **Aktivitas Terbaru:** Daftar log transaksi/perubahan status terkini milik tenant tersebut.
2. **Pesanan**
   - **Daftar Pesanan:** Filter berdasarkan status (_Semua, Antre, Proses, Siap Diambil, Selesai_) dan pencarian berdasarkan nomor nota, nama pelanggan, atau no. telepon.
   - **Form Buat Pesanan Baru:**
     - Pilih pelanggan terdaftar milik tenant (atau tambah cepat).
     - Tambah multi-item layanan (otomatis menghitung subtotal berdasarkan tipe: _Per Kg_, _Per Item_, atau _Per_ $m^2$).
     - Catatan pesanan & ringkasan total harga.
     - **Pembayaran Awal (Opsional/DP):** Input metode bayar (Tunai/Transfer/QRIS), nominal uang diterima, preset tombol nominal (_Uang Pas, Rp 50.000, Rp 100.000_), dan kalkulasi kembalian/sisa tagihan otomatis.
3. **Layanan**
   - **Statistik Katalog:** 5 kartu metrik (_Total Layanan_, _Layanan Aktif_, _Per Kg_, _Per Item_, _Per_ $m^2$).
   - **Tabel Katalog Layanan:** Menampilkan nama layanan, badge tipe satuan, harga, status (_Aktif/Nonaktif_), dan aksi (_Detail, Edit, Nonaktifkan, Hapus_).
   - **Form Tambah/Edit Layanan:** Input Nama Layanan, Tipe Layanan (`PER_KG`, `PER_ITEM`, `PER_M2`), dan Harga (Rp).
4. **Pelanggan**
   - **Daftar Pelanggan:** Pencarian nama/nomor telepon, tabel berisi Nama, Telepon, Alamat, Tanggal Dibuat, serta aksi (_Detail riwayat order, Edit_).
   - **Form Tambah/Edit Pelanggan:** Nama lengkap, No. Telepon (unik per tenant), Alamat, dan Catatan khusus pelanggan.
5. **Keuangan**
   - **Ringkasan Kas:** Total Pemasukan (dari pelunasan/DP pesanan), Total Pengeluaran operasional, dan Saldo Bersih.
   - **Catat Pengeluaran:** Form input pengeluaran harian (misal: beli deterjen, gas, listrik, plastik).
   - **Riwayat Arus Kas:** Tabel mutasi uang masuk dan keluar lengkap dengan filter tanggal.
6. **Pengaturan**
   - **Profil Pengguna & Bisnis:** Mengelola data akun tenant (Nama, Email, Password) sekaligus profil outlet/toko (Nama Laundry, Alamat toko, dan Nomor WhatsApp toko untuk footer struk/nota).

## 5. Skema Database Multi-Tenant (PostgreSQL)

Semua tabel operasional menggunakan `user_id` sebagai kunci isolasi tenant, dilengkapi indeks komposit dan batasan unik per-tenant (`UNIQUE (user_id, ...)`), sehingga User A dan User B bisa memiliki pelanggan dengan nomor HP yang sama atau kode nota yang sama tanpa bentrok.

```sql
-- 1. Tabel Users (Setiap baris adalah 1 Tenant / Akun Bisnis Laundry)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    business_name VARCHAR(100) NOT NULL DEFAULT 'Celana - Cellix Laundry',
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tabel Pelanggan (Terisolasi per user_id)
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, phone) -- Nomor telepon unik hanya di dalam ruang lingkup 1 tenant
);
CREATE INDEX idx_customers_user_id ON customers(user_id);

-- 3. Tabel Katalog Layanan (Terisolasi per user_id)
CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    unit_type VARCHAR(20) NOT NULL CHECK (unit_type IN ('PER_KG', 'PER_ITEM', 'PER_M2')),
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_services_user_id ON services(user_id);

-- 4. Tabel Pesanan Header (Terisolasi per user_id)
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_code VARCHAR(30) NOT NULL, -- Contoh: CLX-20261117-001
    customer_id INT NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'PROCESSING', 'READY', 'COMPLETED', 'CANCELLED')),
    payment_status VARCHAR(20) NOT NULL DEFAULT 'UNPAID'
        CHECK (payment_status IN ('UNPAID', 'PARTIAL', 'PAID')),
    total_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
    paid_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, order_code) -- Kode nota unik per tenant
);
CREATE INDEX idx_orders_user_status ON orders(user_id, status);

-- 5. Tabel Detail Item Pesanan
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    service_id INT NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
    quantity NUMERIC(10, 2) NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(12, 2) NOT NULL,
    subtotal NUMERIC(12, 2) NOT NULL
);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);

-- 6. Tabel Keuangan / Arus Kas (Terisolasi per user_id)
CREATE TABLE cash_flows (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id INT REFERENCES orders(id) ON DELETE SET NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT')),
    category VARCHAR(50) NOT NULL, -- 'ORDER_PAYMENT', 'OPERATIONAL', 'SUPPLIES'
    payment_method VARCHAR(30) DEFAULT 'TUNAI',
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    description TEXT,
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_cash_flows_user_date ON cash_flows(user_id, transaction_date DESC);
```

## 6. Rancangan Endpoint REST API (Backend Go)

Seluruh rute di bawah `/api/*` (kecuali `/api/auth/*`) dilindungi oleh **JWT & Tenant Context Middleware** yang menyuntikkan `userID` ke `c.Locals("userID")`.

| Modul          | Method   | Endpoint                   | Deskripsi & Isolasi Tenant                                            |
| :------------- | :------- | :------------------------- | :-------------------------------------------------------------------- |
| **Auth**       | `POST`   | `/api/auth/register`       | Registrasi tenant/user baru + seed layanan default opsional           |
|                | `POST`   | `/api/auth/login`          | Login user & menerbitkan JWT berisi `user_id`                         |
| **Dashboard**  | `GET`    | `/api/dashboard/summary`   | Statistik hari ini & 5 aktivitas terbaru milik `user_id`              |
| **Pesanan**    | `GET`    | `/api/orders`              | List pesanan milik `user_id` (`?status=&search=`)                     |
|                | `POST`   | `/api/orders`              | Buat pesanan baru (validasi kepemilikan `customer_id` & `service_id`) |
|                | `GET`    | `/api/orders/:id`          | Detail pesanan (`WHERE id = $1 AND user_id = $2`)                     |
|                | `PATCH`  | `/api/orders/:id/status`   | Update status cucian (`WHERE id = $1 AND user_id = $2`)               |
|                | `POST`   | `/api/orders/:id/pay`      | Pelunasan sisa tagihan (`WHERE id = $1 AND user_id = $2`)             |
| **Layanan**    | `GET`    | `/api/services`            | List layanan & metrik milik `user_id`                                 |
|                | `POST`   | `/api/services`            | Tambah layanan baru untuk `user_id`                                   |
|                | `PUT`    | `/api/services/:id`        | Edit layanan (`WHERE id = $1 AND user_id = $2`)                       |
|                | `PATCH`  | `/api/services/:id/toggle` | Aktif/Nonaktifkan layanan (`WHERE id = $1 AND user_id = $2`)          |
|                | `DELETE` | `/api/services/:id`        | Hapus layanan (`WHERE id = $1 AND user_id = $2`)                      |
| **Pelanggan**  | `GET`    | `/api/customers`           | List & cari pelanggan milik `user_id`                                 |
|                | `POST`   | `/api/customers`           | Tambah pelanggan baru untuk `user_id`                                 |
|                | `PUT`    | `/api/customers/:id`       | Update pelanggan (`WHERE id = $1 AND user_id = $2`)                   |
| **Keuangan**   | `GET`    | `/api/finance/summary`     | Total masuk, keluar, & saldo kas milik `user_id`                      |
|                | `GET`    | `/api/finance/cashflows`   | Riwayat arus kas milik `user_id`                                      |
|                | `POST`   | `/api/finance/expenses`    | Catat pengeluaran baru untuk `user_id`                                |
| **Pengaturan** | `GET`    | `/api/users/profile`       | Ambil profil user & bisnis berdasarkan `user_id` JWT                  |
|                | `PUT`    | `/api/users/profile`       | Update profil user & bisnis berdasarkan `user_id` JWT                 |

## 7. Struktur Folder Proyek yang Direkomendasikan

```text
celana/
├── backend/                  # Golang (Fiber + pgx/GORM)
│   ├── cmd/api/main.go       # Entry point aplikasi Celana API
│   ├── internal/
│   │   ├── config/           # Koneksi PostgreSQL & JWT Secret
│   │   ├── middleware/       # auth_tenant.go (Validasi JWT & Ekstrak user_id)
│   │   ├── models/           # Struct Go & Request/Response DTO
│   │   ├── handlers/         # HTTP Handlers dengan strict WHERE user_id = ?
│   │   └── routes/           # Definisi endpoint API
│   ├── migrations/           # 001_init_multi_tenant.sql
│   └── go.mod
│
└── frontend/                 # Svelte + Tailwind CSS (Responsive Mobile-First)
    ├── src/
    │   ├── lib/
    │   │   ├── api.js        # Fetch wrapper otomatis menyisipkan Bearer Token
    │   │   ├── stores/       # authStore.js (State sesi user & profil bisnis)
    │   │   └── components/   # Sidebar.svelte (Responsive Drawer), StatCard.svelte, Badge.svelte
    │   ├── routes/           # Login, Register, Dashboard, Pesanan, Layanan, Pelanggan, Keuangan, Pengaturan
    │   └── app.css           # Konfigurasi Tailwind CSS
    ├── tailwind.config.js
    └── package.json
```
