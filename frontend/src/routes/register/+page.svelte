<script>
	import { WashingMachine, ArrowLeft, ShieldAlert } from 'lucide-svelte';
</script>

<div class="flex min-h-screen items-center justify-center bg-[#F8FAFC] px-4 py-10">
	<div class="w-full max-w-md rounded-2xl border border-[#E2E8F0] bg-white p-8 text-center shadow-sm">
		<div class="mx-auto mb-5 flex h-12 w-12 items-center justify-center rounded-xl bg-[#0284C7] text-white shadow-sm">
			<WashingMachine class="h-6 w-6" />
		</div>
		<div class="mx-auto mb-3 inline-flex items-center gap-1.5 rounded-full bg-[#FEE2E2] px-3 py-1 text-xs font-semibold text-[#B91C1C]">
			<ShieldAlert class="h-3.5 w-3.5" />
			<span>Registrasi Dibatasi</span>
		</div>
		<h1 class="text-xl font-extrabold text-[#0F172A]">Pembuatan Akun Khusus Admin</h1>
		<p class="mt-2 text-sm text-[#64748B]">
			Pendaftaran akun outlet baru hanya dapat dilakukan oleh Admin melalui menu <strong>Pengguna</strong>.
		</p>
		<a
			href="/login"
			class="mt-6 inline-flex w-full items-center justify-center gap-2 rounded-xl bg-[#0284C7] px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-[#0369A1]"
		>
			<ArrowLeft class="h-4 w-4" />
			<span>Kembali ke Halaman Login</span>
		</a>
	</div>
</div>
