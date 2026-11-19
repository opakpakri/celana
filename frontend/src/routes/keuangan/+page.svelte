<script>
	import { onMount } from 'svelte';
	import { apiFetch, formatRupiah } from '$lib/api.js';
	import { getCache, setCache } from '$lib/stores/dataCache.js';
	import {
		Plus,
		ArrowDown,
		ArrowUp,
		Wallet,
		Trash2,
		Search,
		Tag,
		Banknote,
		ChevronLeft,
		ChevronRight
	} from 'lucide-svelte';

	const cachedFinance = getCache('finance');
	let summary = cachedFinance?.summary || { total_in: 0, total_out: 0, net_balance: 0 };
	let cashFlows = cachedFinance?.cashFlows || [];
	let loading = !cachedFinance;
	let error = '';

	// Form state (inline card matching reference)
	let txType = 'IN'; // IN | OUT
	let txDate = new Date().toISOString().slice(0, 10);
	let selectedCategory = '';
	let customCategory = '';
	let amount = '';
	let description = '';
	let saving = false;
	let formError = '';

	const inCategories = ['Order Laundry', 'Penjualan Produk', 'Modal Tambahan', 'Pemasukan Lainnya'];
	const outCategories = ['Pembelian Stok', 'Deterjen & Parfum', 'Listrik & Air', 'Gas & Operasional', 'Gaji Karyawan', 'Perawatan Mesin'];

	$: categoryOptions = txType === 'IN' ? inCategories : outCategories;

	// Filter & pagination state
	let filterType = ''; // '' | 'IN' | 'OUT'
	let search = '';
	let currentPage = 1;
	const pageSize = 10;

	async function loadFinance() {
		if (cashFlows.length === 0) loading = true;
		error = '';
		try {
			const params = new URLSearchParams();
			if (filterType) params.set('type', filterType);
			if (search.trim()) params.set('search', search.trim());
			const q = params.toString() ? `?${params.toString()}` : '';

			const [sumRes, cfRes] = await Promise.all([
				apiFetch('/finance/summary'),
				apiFetch(`/finance/cashflows${q}`)
			]);
			summary = sumRes;
			cashFlows = cfRes || [];
			if (!filterType && !search.trim()) {
				setCache('finance', { summary, cashFlows });
			}
			currentPage = 1;
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	onMount(loadFinance);

	function switchTxType(type) {
		txType = type;
		selectedCategory = '';
	}

	async function handleSaveTransaction() {
		formError = '';
		const finalCategory = customCategory.trim() || selectedCategory;
		if (!finalCategory) {
			formError = 'Pilih atau tulis kategori transaksi terlebih dahulu.';
			return;
		}
		if (!amount || Number(amount) <= 0) {
			formError = 'Masukkan nominal transaksi yang valid (> 0).';
			return;
		}

		saving = true;
		try {
			await apiFetch('/finance/transactions', {
				method: 'POST',
				body: JSON.stringify({
					type: txType,
					category: finalCategory,
					payment_method: 'TUNAI',
					amount: Number(amount),
					description: description.trim(),
					transaction_date: txDate
				})
			});
			selectedCategory = '';
			customCategory = '';
			amount = '';
			description = '';
			await loadFinance();
		} catch (err) {
			formError = err.message;
		} finally {
			saving = false;
		}
	}

	async function handleDeleteExpense(cf) {
		if (!confirm(`Hapus catatan keuangan "${cf.description || cf.category}"?`)) return;
		error = '';
		try {
			await apiFetch(`/finance/expenses/${cf.id}`, { method: 'DELETE' });
			await loadFinance();
		} catch (err) {
			error = err.message;
		}
	}

	function formatDateOnly(iso) {
		if (!iso) return '-';
		const d = new Date(iso);
		const day = String(d.getDate()).padStart(2, '0');
		const month = String(d.getMonth() + 1).padStart(2, '0');
		const year = d.getFullYear();
		return `${day}/${month}/${year}`;
	}

	$: totalItems = cashFlows.length;
	$: totalPages = Math.max(1, Math.ceil(totalItems / pageSize));
	$: paginatedFlows = cashFlows.slice((currentPage - 1) * pageSize, currentPage * pageSize);
	$: startIdx = totalItems === 0 ? 0 : (currentPage - 1) * pageSize + 1;
	$: endIdx = Math.min(currentPage * pageSize, totalItems);
</script>

<div class="space-y-5">
	<div>
		<h1 class="text-2xl font-extrabold text-[#0F172A]">Keuangan</h1>
		<p class="text-sm text-[#64748B]">Pantau pemasukan, pengeluaran operasional, dan saldo kas outlet Anda.</p>
	</div>

	<!-- 3 Kartu Ringkasan Gradien -->
	<div data-tour="finance-cards" class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4">
		<div class="rounded-2xl bg-gradient-to-br from-[#0284C7] to-[#0891B2] p-4 text-white shadow-sm sm:p-5">
			<div class="flex h-9 w-9 items-center justify-center rounded-xl bg-white/20">
				<ArrowDown class="h-4 w-4" />
			</div>
			<p class="mt-3 truncate text-lg font-extrabold tracking-tight sm:mt-4 sm:text-2xl lg:text-3xl">{formatRupiah(summary.total_in)}</p>
			<p class="mt-1 text-xs font-medium text-sky-100">Total Pemasukan</p>
		</div>

		<div class="rounded-2xl bg-gradient-to-br from-[#EF4444] to-[#E11D48] p-4 text-white shadow-sm sm:p-5">
			<div class="flex h-9 w-9 items-center justify-center rounded-xl bg-white/20">
				<ArrowUp class="h-4 w-4" />
			</div>
			<p class="mt-3 truncate text-lg font-extrabold tracking-tight sm:mt-4 sm:text-2xl lg:text-3xl">{formatRupiah(summary.total_out)}</p>
			<p class="mt-1 text-xs font-medium text-rose-100">Total Pengeluaran</p>
		</div>

		<div class="rounded-2xl bg-gradient-to-br from-[#0891B2] to-[#0D9488] p-4 text-white shadow-sm sm:p-5">
			<div class="flex h-9 w-9 items-center justify-center rounded-xl bg-white/20">
				<Wallet class="h-4 w-4" />
			</div>
			<p class="mt-3 truncate text-lg font-extrabold tracking-tight sm:mt-4 sm:text-2xl lg:text-3xl">{formatRupiah(summary.net_balance)}</p>
			<p class="mt-1 text-xs font-medium text-teal-100">Saldo Bersih</p>
		</div>
	</div>

	<!-- Kartu Form Inline: Catat Transaksi Keuangan -->
	<section data-tour="finance-form" class="rounded-2xl border border-[#E2E8F0] bg-white p-5 shadow-sm">
		<h2 class="flex items-center gap-2 text-sm font-bold text-[#0F172A]">
			<Plus class="h-4 w-4 text-[#0284C7]" />
			<span>Catat transaksi keuangan</span>
		</h2>

		<!-- Segmented Switcher Pemasukan / Pengeluaran -->
		<div class="mt-3 grid grid-cols-2 rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-1">
			<button
				type="button"
				on:click={() => switchTxType('IN')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {txType === 'IN'
					? 'bg-white text-[#0284C7] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<ArrowDown class="h-3.5 w-3.5" />
				<span>Pemasukan</span>
			</button>
			<button
				type="button"
				on:click={() => switchTxType('OUT')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {txType === 'OUT'
					? 'bg-white text-[#B91C1C] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<ArrowUp class="h-3.5 w-3.5" />
				<span>Pengeluaran</span>
			</button>
		</div>

		{#if formError}
			<div class="mt-3 rounded-xl bg-[#FEE2E2] px-4 py-2.5 text-xs font-medium text-[#B91C1C]">{formError}</div>
		{/if}

		<form on:submit|preventDefault={handleSaveTransaction} class="mt-3 space-y-3">
			<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
				<input
					type="date"
					bind:value={txDate}
					required
					class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-xs text-[#0F172A] outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				/>

				<select
					bind:value={selectedCategory}
					class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-xs text-[#0F172A] outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				>
					<option value="">Pilih kategori</option>
					{#each categoryOptions as cat}
						<option value={cat}>{cat}</option>
					{/each}
				</select>

				<div class="relative">
					<Tag class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[#64748B]" />
					<input
						type="text"
						bind:value={customCategory}
						placeholder="Atau tulis kategori"
						class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2 pl-8 pr-3 text-xs text-[#0F172A] outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
					/>
				</div>

				<div class="relative">
					<Banknote class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[#64748B]" />
					<input
						type="number"
						min="1"
						step="1"
						bind:value={amount}
						required
						placeholder="Nominal (Rp)"
						class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2 pl-8 pr-3 text-xs text-[#0F172A] outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
					/>
				</div>

				<input
					type="text"
					bind:value={description}
					placeholder="Keterangan (opsional)"
					class="w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-xs text-[#0F172A] outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				/>
			</div>

			<div>
				<button
					type="submit"
					disabled={saving}
					class="inline-flex items-center gap-1.5 rounded-xl bg-[#0284C7] px-4 py-2 text-xs font-semibold text-white shadow-sm transition hover:bg-[#0369A1] disabled:opacity-60"
				>
					<Plus class="h-3.5 w-3.5" />
					<span>{saving ? 'Menyimpan...' : 'Simpan transaksi'}</span>
				</button>
			</div>
		</form>
	</section>

	{#if error}
		<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{error}</div>
	{/if}

	<!-- Kartu Tabel Arus Kas -->
	<section data-tour="finance-table" class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<!-- Filter Bar & Search -->
		<div class="flex flex-col justify-between gap-3 border-b border-[#E2E8F0] p-4 sm:flex-row sm:items-center">
			<div class="grid grid-cols-3 rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-1 sm:w-80">
				<button
					type="button"
					on:click={() => {
						filterType = '';
						loadFinance();
					}}
					class="rounded-lg py-1.5 text-xs font-semibold transition {filterType === ''
						? 'bg-white text-[#0284C7] shadow-sm'
						: 'text-[#64748B] hover:text-[#0F172A]'}"
				>
					Semua
				</button>
				<button
					type="button"
					on:click={() => {
						filterType = 'IN';
						loadFinance();
					}}
					class="rounded-lg py-1.5 text-xs font-semibold transition {filterType === 'IN'
						? 'bg-white text-[#0284C7] shadow-sm'
						: 'text-[#64748B] hover:text-[#0F172A]'}"
				>
					Pemasukan
				</button>
				<button
					type="button"
					on:click={() => {
						filterType = 'OUT';
						loadFinance();
					}}
					class="rounded-lg py-1.5 text-xs font-semibold transition {filterType === 'OUT'
						? 'bg-white text-[#0284C7] shadow-sm'
						: 'text-[#64748B] hover:text-[#0F172A]'}"
				>
					Pengeluaran
				</button>
			</div>

			<form on:submit|preventDefault={loadFinance} class="relative sm:w-80">
				<Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
				<input
					type="text"
					bind:value={search}
					on:input={() => {
						if (!search) loadFinance();
					}}
					placeholder="Cari kategori atau keterangan"
					class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2 pl-9 pr-3 text-xs text-[#0F172A] outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				/>
			</form>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full text-left text-xs">
				<thead class="border-b border-[#E2E8F0] bg-[#F8FAFC] text-[11px] font-semibold text-[#64748B]">
					<tr>
						<th class="px-5 py-3">Tanggal</th>
						<th class="px-5 py-3">Jenis</th>
						<th class="px-5 py-3">Kategori</th>
						<th class="px-5 py-3">Keterangan</th>
						<th class="px-5 py-3 text-right">Nominal</th>
						<th class="px-5 py-3 text-right"></th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[#E2E8F0]">
					{#if paginatedFlows.length === 0}
						<tr><td colspan="6" class="px-5 py-8 text-center text-[#64748B]">Belum ada mutasi arus kas.</td></tr>
					{:else}
						{#each paginatedFlows as cf}
							<tr class="hover:bg-slate-50/70">
								<td class="whitespace-nowrap px-5 py-3.5 text-[#0F172A]">{formatDateOnly(cf.transaction_date)}</td>
								<td class="whitespace-nowrap px-5 py-3.5">
									{#if cf.type === 'IN'}
										<span class="inline-flex items-center gap-1 rounded-full bg-[#DCFCE7] px-2.5 py-0.5 text-[11px] font-semibold text-[#15803D]">
											<ArrowDown class="h-3 w-3" />
											<span>Pemasukan</span>
										</span>
									{:else}
										<span class="inline-flex items-center gap-1 rounded-full bg-[#FEE2E2] px-2.5 py-0.5 text-[11px] font-semibold text-[#B91C1C]">
											<ArrowUp class="h-3 w-3" />
											<span>Pengeluaran</span>
										</span>
									{/if}
								</td>
								<td class="whitespace-nowrap px-5 py-3.5 font-medium text-[#0F172A]">{cf.category}</td>
								<td class="px-5 py-3.5 text-[#64748B]">{cf.description || '-'}</td>
								<td class="whitespace-nowrap px-5 py-3.5 text-right font-bold {cf.type === 'IN' ? 'text-[#15803D]' : 'text-[#0F172A]'}">
									{cf.type === 'IN' ? '+' : '−'}{formatRupiah(cf.amount)}
								</td>
								<td class="whitespace-nowrap px-5 py-3.5 text-right">
									<button
										type="button"
										on:click={() => handleDeleteExpense(cf)}
										class="inline-flex items-center gap-1 rounded-lg bg-[#FEE2E2] px-2.5 py-1 text-[11px] font-semibold text-[#B91C1C] transition hover:bg-red-200"
									>
										<Trash2 class="h-3 w-3" />
										<span>Hapus</span>
									</button>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>

		<!-- Pagination Footer -->
		<div class="flex flex-col items-center justify-between gap-3 border-t border-[#E2E8F0] px-5 py-3 text-xs text-[#64748B] sm:flex-row">
			<span>Menampilkan {startIdx}–{endIdx} dari {totalItems} data</span>
			<div class="flex items-center gap-1.5">
				<button
					type="button"
					disabled={currentPage <= 1}
					on:click={() => (currentPage = Math.max(1, currentPage - 1))}
					class="inline-flex items-center gap-1 rounded-lg border border-[#E2E8F0] bg-white px-2.5 py-1.5 text-xs font-medium text-[#64748B] hover:bg-slate-50 disabled:opacity-40"
				>
					<ChevronLeft class="h-3.5 w-3.5" />
					<span>Sebelumnya</span>
				</button>
				<span class="flex h-7 min-w-[28px] items-center justify-center rounded-lg border border-[#0284C7] bg-[#EFF6FF] px-2 text-xs font-bold text-[#1D4ED8]">
					{currentPage}
				</span>
				<button
					type="button"
					disabled={currentPage >= totalPages}
					on:click={() => (currentPage = Math.min(totalPages, currentPage + 1))}
					class="inline-flex items-center gap-1 rounded-lg border border-[#E2E8F0] bg-white px-2.5 py-1.5 text-xs font-medium text-[#64748B] hover:bg-slate-50 disabled:opacity-40"
				>
					<span>Berikutnya</span>
					<ChevronRight class="h-3.5 w-3.5" />
				</button>
			</div>
		</div>
	</section>
</div>
