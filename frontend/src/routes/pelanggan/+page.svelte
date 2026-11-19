<script>
	import { onMount } from 'svelte';
	import { apiFetch, formatShortDate, formatRupiah, orderStatusLabel, paymentStatusLabel } from '$lib/api.js';
	import { getCache, setCache } from '$lib/stores/dataCache.js';
	import Modal from '$lib/components/Modal.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import { Plus, Search, Pencil, History, Phone, MapPin } from 'lucide-svelte';

	const cachedCustomers = getCache('customers');
	let customers = cachedCustomers || [];
	let search = '';
	let loading = !cachedCustomers;
	let error = '';
	let currentPage = 1;
	const pageSize = 10;

	$: paginatedCustomers = customers.slice((currentPage - 1) * pageSize, currentPage * pageSize);

	let formModalOpen = false;
	let editingId = null;
	let form = { name: '', phone: '', address: '', notes: '' };
	let saving = false;
	let formError = '';

	let historyModalOpen = false;
	let historyData = null;
	let historyLoading = false;

	async function loadCustomers() {
		if (customers.length === 0) loading = true;
		error = '';
		try {
			const q = search ? `?search=${encodeURIComponent(search)}` : '';
			customers = await apiFetch(`/customers${q}`);
			if (!search.trim()) setCache('customers', customers);
			currentPage = 1;
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	onMount(loadCustomers);

	function openAdd() {
		editingId = null;
		form = { name: '', phone: '', address: '', notes: '' };
		formError = '';
		formModalOpen = true;
	}

	function openEdit(c) {
		editingId = c.id;
		form = { name: c.name, phone: c.phone, address: c.address || '', notes: c.notes || '' };
		formError = '';
		formModalOpen = true;
	}

	async function handleSave() {
		formError = '';
		saving = true;
		try {
			if (editingId) {
				await apiFetch(`/customers/${editingId}`, { method: 'PUT', body: JSON.stringify(form) });
			} else {
				await apiFetch('/customers', { method: 'POST', body: JSON.stringify(form) });
			}
			formModalOpen = false;
			await loadCustomers();
		} catch (err) {
			formError = err.message;
		} finally {
			saving = false;
		}
	}

	async function openHistory(c) {
		historyModalOpen = true;
		historyLoading = true;
		historyData = null;
		try {
			historyData = await apiFetch(`/customers/${c.id}`);
		} catch (err) {
			error = err.message;
			historyModalOpen = false;
		} finally {
			historyLoading = false;
		}
	}
</script>

<div class="space-y-6">
	<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">Daftar Pelanggan</h1>
			<p class="text-sm text-[#64748B]">Kelola data pelanggan dan riwayat pesanan masing-masing pelanggan.</p>
		</div>
		<button
			type="button"
			data-tour="customers-add-btn"
			on:click={openAdd}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#0284C7] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#0369A1]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah Pelanggan</span>
		</button>
	</div>

	<!-- Pencarian -->
	<form data-tour="customers-search" on:submit|preventDefault={loadCustomers} class="flex gap-2">
		<div class="relative flex-1 max-w-md">
			<Search class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
			<input
				type="text"
				bind:value={search}
				placeholder="Cari nama atau nomor telepon..."
				class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-10 pr-4 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
			/>
		</div>
		<button
			type="submit"
			class="rounded-xl border border-[#E2E8F0] bg-white px-4 py-2.5 text-sm font-semibold text-[#0F172A] hover:bg-slate-50"
		>
			Cari
		</button>
	</form>

	{#if error}
		<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{error}</div>
	{/if}

	<div data-tour="customers-table" class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead class="border-b border-[#E2E8F0] bg-[#F8FAFC] text-xs font-semibold text-[#64748B]">
					<tr>
						<th class="px-5 py-3.5">Nama</th>
						<th class="whitespace-nowrap px-5 py-3.5">Telepon</th>
						<th class="px-5 py-3.5">Alamat</th>
						<th class="min-w-[120px] whitespace-nowrap px-5 py-3.5">Total Order</th>
						<th class="whitespace-nowrap px-5 py-3.5">Tanggal Dibuat</th>
						<th class="whitespace-nowrap px-5 py-3.5 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[#E2E8F0]">
					{#if customers.length === 0}
						<tr><td colspan="6" class="px-5 py-8 text-center text-[#64748B]">Belum ada pelanggan ditemukan.</td></tr>
					{:else}
						{#each paginatedCustomers as c}
							<tr class="hover:bg-slate-50/70">
								<td class="px-5 py-4">
									<p class="font-semibold text-[#0F172A]">{c.name}</p>
									{#if c.notes}
										<p class="text-xs text-[#64748B]">{c.notes}</p>
									{/if}
								</td>
								<td class="whitespace-nowrap px-5 py-4 text-[#0F172A]">
									<span class="inline-flex items-center gap-1.5">
										<Phone class="h-3.5 w-3.5 text-[#64748B]" />
										{c.phone}
									</span>
								</td>
								<td class="px-5 py-4 text-[#64748B]">{c.address || '-'}</td>
								<td class="min-w-[120px] whitespace-nowrap px-5 py-4">
									<Badge variant="info">{c.order_count || 0} order</Badge>
								</td>
								<td class="whitespace-nowrap px-5 py-4 text-[#64748B]">{formatShortDate(c.created_at)}</td>
								<td class="whitespace-nowrap px-5 py-4">
									<div class="flex items-center justify-end gap-1.5">
										<button
											type="button"
											on:click={() => openHistory(c)}
											title="Riwayat Order"
											class="inline-flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-xs font-semibold text-[#0284C7] hover:bg-[#EFF6FF]"
										>
											<History class="h-4 w-4" />
											<span>Riwayat</span>
										</button>
										<button
											type="button"
											on:click={() => openEdit(c)}
											title="Edit"
											class="rounded-lg p-2 text-[#64748B] hover:bg-[#EFF6FF] hover:text-[#1D4ED8]"
										>
											<Pencil class="h-4 w-4" />
										</button>
									</div>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
		<Pagination bind:currentPage totalItems={customers.length} {pageSize} />
	</div>
</div>

<!-- Modal Tambah / Edit Pelanggan -->
<Modal open={formModalOpen} title={editingId ? 'Edit Pelanggan' : 'Tambah Pelanggan Baru'} on:close={() => (formModalOpen = false)}>
	<form on:submit|preventDefault={handleSave} class="space-y-4">
		{#if formError}
			<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{formError}</div>
		{/if}
		<div>
			<label for="cust-name" class="block text-sm font-semibold text-[#0F172A]">Nama Lengkap</label>
			<input
				id="cust-name"
				type="text"
				bind:value={form.name}
				required
				placeholder="Nama Pelanggan"
				class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
			/>
		</div>
		<div>
			<label for="cust-phone" class="block text-sm font-semibold text-[#0F172A]">No. Telepon / WhatsApp</label>
			<input
				id="cust-phone"
				type="text"
				bind:value={form.phone}
				required
				placeholder="No Telp Pelanggan"
				class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
			/>
		</div>
		<div>
			<label for="cust-addr" class="block text-sm font-semibold text-[#0F172A]">Alamat</label>
			<input
				id="cust-addr"
				type="text"
				bind:value={form.address}
				placeholder="Alamat Pelanggan"
				class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
			/>
		</div>
		<div>
			<label for="cust-notes" class="block text-sm font-semibold text-[#0F172A]">Catatan Khusus</label>
			<textarea
				id="cust-notes"
				rows="2"
				bind:value={form.notes}
				placeholder="Catatan (Opsional)"
				class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
			></textarea>
		</div>
		<div class="flex justify-end gap-2 pt-2">
			<button
				type="button"
				on:click={() => (formModalOpen = false)}
				class="rounded-xl border border-[#E2E8F0] px-4 py-2 text-sm font-semibold text-[#64748B] hover:bg-slate-50"
			>
				Batal
			</button>
			<button
				type="submit"
				disabled={saving}
				class="rounded-xl bg-[#0284C7] px-4 py-2 text-sm font-semibold text-white hover:bg-[#0369A1] disabled:opacity-60"
			>
				{saving ? 'Menyimpan...' : 'Simpan Pelanggan'}
			</button>
		</div>
	</form>
</Modal>

<!-- Modal Riwayat Order Pelanggan -->
<Modal open={historyModalOpen} title="Detail & Riwayat Pesanan" maxWidth="max-w-2xl" on:close={() => (historyModalOpen = false)}>
	{#if historyLoading}
		<p class="py-6 text-center text-sm text-[#64748B]">Memuat riwayat...</p>
	{:else if historyData}
		<div class="space-y-4">
			<div class="rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-4">
				<p class="text-base font-bold text-[#0F172A]">{historyData.customer.name}</p>
				<p class="mt-1 flex items-center gap-2 text-xs text-[#64748B]">
					<Phone class="h-3.5 w-3.5" /> {historyData.customer.phone}
					{#if historyData.customer.address}
						<span>•</span>
						<MapPin class="h-3.5 w-3.5" /> {historyData.customer.address}
					{/if}
				</p>
			</div>

			{#if historyData.orders.length === 0}
				<p class="py-4 text-center text-sm text-[#64748B]">Pelanggan ini belum memiliki pesanan.</p>
			{:else}
				<div class="divide-y divide-[#E2E8F0] rounded-xl border border-[#E2E8F0]">
					{#each historyData.orders as o}
						<div class="flex items-center justify-between px-4 py-3 text-sm">
							<div>
								<p class="font-bold text-[#0F172A]">{o.order_code}</p>
								<p class="text-xs text-[#64748B]">{formatShortDate(o.created_at)}</p>
							</div>
							<div class="flex items-center gap-2">
								<Badge variant={o.status === 'COMPLETED' || o.status === 'READY' ? 'success' : 'info'}>
									{orderStatusLabel(o.status)}
								</Badge>
								<Badge variant={o.payment_status === 'PAID' ? 'success' : 'danger'}>
									{paymentStatusLabel(o.payment_status)}
								</Badge>
								<span class="ml-2 font-bold text-[#0F172A]">{formatRupiah(o.total_amount)}</span>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</Modal>
