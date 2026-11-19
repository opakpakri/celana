<script>
	import { onMount } from 'svelte';
	import { apiFetch, formatRupiah, formatDate, orderStatusLabel, paymentStatusLabel } from '$lib/api.js';
	import { authStore } from '$lib/stores/authStore.js';
	import { getCache, setCache } from '$lib/stores/dataCache.js';
	import { maybeStartFirstVisitTour } from '$lib/tour.js';
	import Badge from '$lib/components/Badge.svelte';
	import {
		PlusCircle,
		ListOrdered,
		ShoppingBag,
		Banknote,
		Loader2,
		CheckCircle2,
		ArrowRight,
		Clock
	} from 'lucide-svelte';

	const cached = getCache('dashboard');
	let loading = !cached;
	let error = '';
	let summary = cached || {
		orders_today: 0,
		revenue_today: 0,
		processing_count: 0,
		ready_count: 0,
		recent_activity: []
	};

	onMount(async () => {
		try {
			summary = await apiFetch('/dashboard/summary');
			setCache('dashboard', summary);
		} catch (err) {
			if (!cached) error = err.message;
		} finally {
			loading = false;
			maybeStartFirstVisitTour('/');
		}
	});

	function statusVariant(status) {
		if (status === 'COMPLETED' || status === 'READY') return 'success';
		if (status === 'PROCESSING') return 'info';
		if (status === 'CANCELLED') return 'danger';
		return 'warning';
	}
</script>

<div class="space-y-6">
	<!-- Hero Card Gradien -->
	<section data-tour="dash-hero" class="rounded-2xl bg-gradient-to-r from-[#0284C7] via-[#0891B2] to-[#0D9488] p-6 text-white shadow-sm sm:p-8">
		<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
			<div>
				<p class="text-xs font-semibold uppercase tracking-wider text-sky-100">
					{$authStore.user?.business_name || 'Celana Laundry'}
				</p>
				<h1 class="mt-1 text-2xl font-extrabold sm:text-3xl">
					Halo, {$authStore.user?.name || 'Owner'}!
				</h1>
				<p class="mt-1 text-sm text-sky-100">
					Pantau ringkasan operasional outlet Anda hari ini secara real-time.
				</p>
			</div>
			<a
				href="/pesanan?new=1"
				class="inline-flex items-center justify-center gap-2 rounded-xl bg-white px-4 py-2.5 text-sm font-bold text-[#0284C7] shadow-sm transition hover:bg-sky-50"
			>
				<PlusCircle class="h-4 w-4" />
				<span>Buat Pesanan Baru</span>
			</a>
		</div>

		<!-- 4 Ringkasan Metrik Hari Ini -->
		<div class="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-4 sm:gap-4">
			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-sky-100">
					<ShoppingBag class="h-4 w-4" />
					<span>Pesanan Hari Ini</span>
				</div>
				<p class="mt-2 text-2xl font-extrabold">{summary.orders_today}</p>
			</div>

			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-sky-100">
					<Banknote class="h-4 w-4" />
					<span>Pendapatan Hari Ini</span>
				</div>
				<p class="mt-2 truncate text-xl font-extrabold sm:text-2xl">{formatRupiah(summary.revenue_today)}</p>
			</div>

			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-sky-100">
					<Loader2 class="h-4 w-4" />
					<span>Dalam Proses</span>
				</div>
				<p class="mt-2 text-2xl font-extrabold">{summary.processing_count}</p>
			</div>

			<div class="rounded-xl bg-white/15 p-4 backdrop-blur-sm">
				<div class="flex items-center gap-2 text-xs font-medium text-sky-100">
					<CheckCircle2 class="h-4 w-4" />
					<span>Siap Diambil</span>
				</div>
				<p class="mt-2 text-2xl font-extrabold">{summary.ready_count}</p>
			</div>
		</div>
	</section>

	{#if error}
		<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">
			{error}
		</div>
	{/if}

	<!-- Aksi Cepat -->
	<section data-tour="dash-quick" class="grid gap-4 sm:grid-cols-2">
		<a
			href="/pesanan?new=1"
			class="group flex items-center justify-between rounded-2xl border border-[#E2E8F0] bg-white p-5 shadow-sm transition hover:border-[#0284C7]"
		>
			<div class="flex items-center gap-4">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-[#E0F2FE] text-[#0369A1]">
					<PlusCircle class="h-6 w-6" />
				</div>
				<div>
					<h2 class="text-base font-bold text-[#0F172A]">Buat Pesanan Baru</h2>
					<p class="text-sm text-[#64748B]">Input transaksi cucian pelanggan & cetak nota</p>
				</div>
			</div>
			<ArrowRight class="h-5 w-5 text-[#64748B] transition group-hover:translate-x-1 group-hover:text-[#0284C7]" />
		</a>

		<a
			href="/pesanan"
			class="group flex items-center justify-between rounded-2xl border border-[#E2E8F0] bg-white p-5 shadow-sm transition hover:border-[#0284C7]"
		>
			<div class="flex items-center gap-4">
				<div class="flex h-12 w-12 items-center justify-center rounded-xl bg-[#EFF6FF] text-[#1D4ED8]">
					<ListOrdered class="h-6 w-6" />
				</div>
				<div>
					<h2 class="text-base font-bold text-[#0F172A]">Lihat Semua Pesanan</h2>
					<p class="text-sm text-[#64748B]">Update status antrean cucian dan pelunasan</p>
				</div>
			</div>
			<ArrowRight class="h-5 w-5 text-[#64748B] transition group-hover:translate-x-1 group-hover:text-[#0284C7]" />
		</a>
	</section>

	<!-- Aktivitas Terbaru -->
	<section data-tour="dash-activity" class="rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<div class="flex items-center justify-between border-b border-[#E2E8F0] px-6 py-4">
			<div class="flex items-center gap-2">
				<Clock class="h-5 w-5 text-[#0284C7]" />
				<h2 class="text-base font-bold text-[#0F172A]">Aktivitas Terbaru</h2>
			</div>
			<a href="/pesanan" class="text-sm font-semibold text-[#0284C7] hover:underline">Lihat semua</a>
		</div>

		{#if summary.recent_activity.length === 0}
			<div class="p-8 text-center text-sm text-[#64748B]">
				Belum ada aktivitas pesanan. Mulai dengan membuat pesanan pertama Anda!
			</div>
		{:else}
			<div class="divide-y divide-[#E2E8F0]">
				{#each summary.recent_activity as act}
					<div class="flex flex-col justify-between gap-3 px-6 py-4 sm:flex-row sm:items-center">
						<div>
							<div class="flex items-center gap-2.5">
								<span class="font-bold text-[#0F172A]">{act.order_code}</span>
								<Badge variant={statusVariant(act.status)}>{orderStatusLabel(act.status)}</Badge>
								<Badge variant={act.payment_status === 'PAID' ? 'success' : 'danger'}>
									{paymentStatusLabel(act.payment_status)}
								</Badge>
							</div>
							<p class="mt-1 text-sm text-[#64748B]">
								Pelanggan: <span class="font-medium text-[#0F172A]">{act.customer_name}</span> • {formatDate(act.updated_at)}
							</p>
						</div>
						<div class="text-left sm:text-right">
							<p class="text-base font-bold text-[#0F172A]">{formatRupiah(act.total_amount)}</p>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</section>
</div>
