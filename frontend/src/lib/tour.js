import { driver } from 'driver.js';
import 'driver.js/dist/driver.css';
import { browser } from '$app/environment';

const TOUR_SEEN_PREFIX = 'celana_tour_seen_';

function isVisible(selector) {
	if (!browser) return false;
	const el = document.querySelector(selector);
	if (!el) return false;
	const rect = el.getBoundingClientRect();
	return rect.width > 0 && rect.height > 0;
}

function getNavStep() {
	if (isVisible('[data-tour="sidebar-nav"]')) {
		return {
			element: '[data-tour="sidebar-nav"]',
			popover: {
				title: 'Menu Navigasi Utama',
				description:
					'Gunakan menu di samping ini untuk berpindah cepat antara Dashboard, Pesanan, Layanan, Pelanggan, Keuangan, dan Pengaturan.',
				side: 'right',
				align: 'start'
			}
		};
	}
	if (isVisible('[data-tour="mobile-nav"]')) {
		return {
			element: '[data-tour="mobile-nav"]',
			popover: {
				title: 'Navigasi Bawah (Mobile)',
				description:
					'Gunakan bar navigasi di bawah layar untuk berpindah antar menu utama dengan mudah.',
				side: 'top',
				align: 'center'
			}
		};
	}
	return null;
}

function buildStepsForPath(pathname) {
	const steps = [];

	if (pathname === '/') {
		const navStep = getNavStep();
		if (navStep) steps.push(navStep);

		steps.push(
			{
				element: '[data-tour="dash-hero"]',
				popover: {
					title: 'Ringkasan Operasional Hari Ini',
					description:
						'Pantau 4 metrik utama outlet Anda hari ini: jumlah pesanan masuk, pendapatan hari ini, cucian dalam proses, dan cucian siap diambil.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="dash-quick"]',
				popover: {
					title: 'Pintasan Aksi Cepat',
					description:
						'Klik "Buat Pesanan Baru" untuk langsung menginput transaksi kasir, atau "Lihat Semua Pesanan" untuk mengelola antrean cucian.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="dash-activity"]',
				popover: {
					title: 'Aktivitas Terbaru',
					description:
						'Menampilkan 5 pesanan terakhir beserta status pengerjaan dan status pembayarannya secara langsung.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/pesanan')) {
		steps.push(
			{
				element: '[data-tour="orders-create-btn"]',
				popover: {
					title: 'Buat Pesanan Baru',
					description:
						'Klik tombol ini untuk membuat nota baru: pilih/tambah pelanggan cepat, cari multi-item layanan, dan catat pembayaran awal (DP/Lunas).',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="orders-filter"]',
				popover: {
					title: 'Filter Status & Pencarian',
					description:
						'Saring daftar pesanan berdasarkan tahap pengerjaan (Antre, Proses, Siap Diambil, Selesai) atau cari berdasarkan nomor nota dan nama pelanggan.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="orders-table"]',
				popover: {
					title: 'Tabel Pesanan & Detail Nota',
					description:
						'Klik tombol "Detail" pada pesanan untuk mengubah status cucian secara berurutan, menerima pelunasan, mengirim notifikasi WhatsApp saat Siap Diambil, atau mencetak nota.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/layanan')) {
		steps.push(
			{
				element: '[data-tour="services-add-btn"]',
				popover: {
					title: 'Tambah Layanan Laundry',
					description:
						'Tambahkan katalog layanan baru beserta tipe satuan hitung (Per Kg, Per Item, atau Per m²) dan tarifnya.',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="services-stats"]',
				popover: {
					title: 'Statistik Katalog Layanan',
					description:
						'Ringkasan jumlah total layanan, layanan yang sedang aktif, serta pembagian berdasarkan tipe satuan.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="services-table"]',
				popover: {
					title: 'Kelola Daftar Layanan',
					description:
						'Gunakan tombol aksi di sebelah kanan untuk melihat detail, mengubah harga, menonaktifkan sementara, atau menghapus layanan.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/pelanggan')) {
		steps.push(
			{
				element: '[data-tour="customers-add-btn"]',
				popover: {
					title: 'Tambah Pelanggan Baru',
					description:
						'Simpan data pelanggan tetap beserta nomor WhatsApp, alamat, dan catatan khusus (misal: alergi parfum tertentu).',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="customers-search"]',
				popover: {
					title: 'Cari Data Pelanggan',
					description:
						'Ketik nama atau nomor telepon pelanggan untuk menemukan data dengan cepat.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="customers-table"]',
				popover: {
					title: 'Riwayat & Edit Pelanggan',
					description:
						'Klik tombol "Riwayat" untuk melihat seluruh histori pesanan pelanggan tersebut, atau ikon pensil untuk memperbarui datanya.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/keuangan')) {
		steps.push(
			{
				element: '[data-tour="finance-cards"]',
				popover: {
					title: 'Ringkasan Arus Kas',
					description:
						'Pantau Total Pemasukan (otomatis dari pesanan maupun input manual), Total Pengeluaran, dan Saldo Bersih outlet.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="finance-form"]',
				popover: {
					title: 'Catat Transaksi Keuangan',
					description:
						'Pilih tab Pemasukan atau Pengeluaran, tentukan tanggal, pilih/tulis kategori, masukkan nominal, lalu klik "Simpan transaksi".',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="finance-table"]',
				popover: {
					title: 'Tabel Mutasi & Filter',
					description:
						'Filter riwayat berdasarkan jenis transaksi, cari berdasarkan kategori/keterangan, atau hapus catatan transaksi jika diperlukan.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/pengguna')) {
		steps.push(
			{
				element: '[data-tour="users-add-btn"]',
				popover: {
					title: 'Tambah Akun Pengguna / Outlet',
					description:
						'Sebagai Admin, Anda dapat mendaftarkan akun pemilik outlet baru atau akun Admin tambahan melalui tombol ini.',
					side: 'bottom',
					align: 'end'
				}
			},
			{
				element: '[data-tour="users-table"]',
				popover: {
					title: 'Daftar Akun Terdaftar',
					description:
						'Kelola seluruh akun tenant: ubah informasi outlet, ganti role, reset password, atau hapus akun.',
					side: 'top',
					align: 'start'
				}
			}
		);
	} else if (pathname.startsWith('/pengaturan')) {
		steps.push(
			{
				element: '[data-tour="settings-theme"]',
				popover: {
					title: 'Pengaturan Tema Tampilan',
					description:
						'Pilih tema Terang, Gelap (Dark Mode), atau Sistem yang otomatis mengikuti pengaturan perangkat Anda.',
					side: 'bottom',
					align: 'start'
				}
			},
			{
				element: '[data-tour="settings-form"]',
				popover: {
					title: 'Profil Akun & Identitas Outlet',
					description:
						'Perbarui Nama Laundry, Nomor WhatsApp Toko, dan Alamat Toko. Informasi ini akan otomatis tampil pada nota cetak dan pesan WhatsApp pelanggan.',
					side: 'top',
					align: 'start'
				}
			}
		);
	}

	return steps.filter((s) => !s.element || isVisible(s.element));
}

export function startProductTour(pathname = '/') {
	if (!browser) return;
	const steps = buildStepsForPath(pathname);
	if (steps.length === 0) return;

	const drv = driver({
		showProgress: true,
		animate: true,
		allowClose: true,
		nextBtnText: 'Lanjut →',
		prevBtnText: '← Kembali',
		doneBtnText: 'Selesai',
		progressText: '{{current}} dari {{total}}',
		popoverClass: 'celana-driver-popover',
		steps
	});

	drv.drive();
}

export function maybeStartFirstVisitTour(pathname = '/') {
	if (!browser || pathname !== '/') return;
	const key = `${TOUR_SEEN_PREFIX}dashboard`;
	if (localStorage.getItem(key)) return;
	localStorage.setItem(key, '1');
	setTimeout(() => {
		startProductTour('/');
	}, 500);
}
