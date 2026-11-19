<script>
	import { onMount } from 'svelte';
	import { apiFetch, formatRupiah, unitLabel } from '$lib/api.js';
	import { getCache, setCache } from '$lib/stores/dataCache.js';
	import StatCard from '$lib/components/StatCard.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import Modal from '$lib/components/Modal.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import { Plus, Sparkles, CheckCircle, Scale, Package, Maximize2, Pencil, Power, Trash2, Eye } from 'lucide-svelte';

	const cachedServices = getCache('services');
	let services = cachedServices?.services || [];
	let metrics = cachedServices?.metrics || { total: 0, active: 0, per_kg: 0, per_item: 0, per_m2: 0 };
	let loading = !cachedServices;
	let error = '';
	let currentPage = 1;
	const pageSize = 10;

	$: paginatedServices = services.slice((currentPage - 1) * pageSize, currentPage * pageSize);

	let formModalOpen = false;
	let detailModalOpen = false;
	let selectedService = null;
	let editingId = null;

	let form = { name: '', unit_type: 'PER_KG', price: '' };
	let saving = false;
	let formError = '';

	async function loadServices() {
		if (services.length === 0) loading = true;
		error = '';
		try {
			const res = await apiFetch('/services');
			services = res.services || [];
			metrics = res.metrics || metrics;
			setCache('services', { services, metrics });
			currentPage = 1;
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	onMount(loadServices);

	function openAdd() {
		editingId = null;
		form = { name: '', unit_type: 'PER_KG', price: '' };
		formError = '';
		formModalOpen = true;
	}

	function openEdit(s) {
		editingId = s.id;
		form = { name: s.name, unit_type: s.unit_type, price: String(s.price) };
		formError = '';
		formModalOpen = true;
	}

	function openDetail(s) {
		selectedService = s;
		detailModalOpen = true;
	}

	async function handleSave() {
		formError = '';
		saving = true;
		try {
			const payload = {
				name: form.name,
				unit_type: form.unit_type,
				price: Number(form.price)
			};
			if (editingId) {
				await apiFetch(`/services/${editingId}`, { method: 'PUT', body: JSON.stringify(payload) });
			} else {
				await apiFetch('/services', { method: 'POST', body: JSON.stringify(payload) });
			}
			formModalOpen = false;
			await loadServices();
		} catch (err) {
			formError = err.message;
		} finally {
			saving = false;
		}
	}

	async function handleToggle(s) {
		try {
			await apiFetch(`/services/${s.id}/toggle`, { method: 'PATCH' });
			await loadServices();
		} catch (err) {
			error = err.message;
		}
	}

	async function handleDelete(s) {
		if (!confirm(`Hapus layanan "${s.name}"?`)) return;
		try {
			await apiFetch(`/services/${s.id}`, { method: 'DELETE' });
			await loadServices();
		} catch (err) {
			error = err.message;
		}
	}
</script>

<div class="space-y-6">
	<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">Katalog Layanan</h1>
			<p class="text-sm text-[#64748B]">Atur daftar layanan laundry, satuan hitung, dan tarif outlet Anda.</p>
		</div>
		<button
			type="button"
			data-tour="services-add-btn"
			on:click={openAdd}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#0284C7] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#0369A1]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah Layanan</span>
		</button>
	</div>

	<!-- 5 Kartu Metrik -->
	<div data-tour="services-stats" class="grid grid-cols-2 gap-4 lg:grid-cols-5">
		<StatCard label="Total Layanan" value={metrics.total}>
			<Sparkles slot="icon" class="h-5 w-5" />
		</StatCard>
		<StatCard label="Layanan Aktif" value={metrics.active}>
			<CheckCircle slot="icon" class="h-5 w-5" />
		</StatCard>
		<StatCard label="Per Kg" value={metrics.per_kg}>
			<Scale slot="icon" class="h-5 w-5" />
		</StatCard>
		<StatCard label="Per Item" value={metrics.per_item}>
			<Package slot="icon" class="h-5 w-5" />
		</StatCard>
		<StatCard label="Per m²" value={metrics.per_m2}>
			<Maximize2 slot="icon" class="h-5 w-5" />
		</StatCard>
	</div>

	{#if error}
		<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{error}</div>
	{/if}

	<!-- Tabel Katalog Layanan -->
	<div data-tour="services-table" class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead class="border-b border-[#E2E8F0] bg-[#F8FAFC] text-xs font-semibold text-[#64748B]">
					<tr>
						<th class="w-[35%] px-6 py-3.5">Nama Layanan</th>
						<th class="w-[22%] whitespace-nowrap px-6 py-3.5">Tipe Satuan</th>
						<th class="w-[20%] whitespace-nowrap px-6 py-3.5">Harga</th>
						<th class="w-[11%] whitespace-nowrap px-4 py-3.5">Status</th>
						<th class="w-[12%] whitespace-nowrap px-6 py-3.5 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[#E2E8F0]">
					{#if services.length === 0}
						<tr><td colspan="5" class="px-6 py-8 text-center text-[#64748B]">Belum ada layanan. Tambahkan layanan pertama Anda.</td></tr>
					{:else}
						{#each paginatedServices as s}
							<tr class="hover:bg-slate-50/70">
								<td class="px-6 py-4 font-semibold text-[#0F172A]">{s.name}</td>
								<td class="whitespace-nowrap px-6 py-4">
									<Badge variant="info">{unitLabel(s.unit_type)}</Badge>
								</td>
								<td class="whitespace-nowrap px-6 py-4 font-semibold text-[#0F172A]">{formatRupiah(s.price)}</td>
								<td class="whitespace-nowrap px-4 py-4">
									<Badge variant={s.is_active ? 'success' : 'danger'}>
										{s.is_active ? 'Aktif' : 'Nonaktif'}
									</Badge>
								</td>
								<td class="whitespace-nowrap px-6 py-4">
									<div class="flex items-center justify-end gap-1.5">
										<button
											type="button"
											on:click={() => openDetail(s)}
											title="Detail"
											class="rounded-lg p-2 text-[#64748B] hover:bg-slate-100 hover:text-[#0F172A]"
										>
											<Eye class="h-4 w-4" />
										</button>
										<button
											type="button"
											on:click={() => openEdit(s)}
											title="Edit"
											class="rounded-lg p-2 text-[#64748B] hover:bg-[#EFF6FF] hover:text-[#1D4ED8]"
										>
											<Pencil class="h-4 w-4" />
										</button>
										<button
											type="button"
											on:click={() => handleToggle(s)}
											title={s.is_active ? 'Nonaktifkan' : 'Aktifkan'}
											class="rounded-lg p-2 text-[#64748B] hover:bg-amber-50 hover:text-amber-700"
										>
											<Power class="h-4 w-4" />
										</button>
										<button
											type="button"
											on:click={() => handleDelete(s)}
											title="Hapus"
											class="rounded-lg p-2 text-[#64748B] hover:bg-[#FEE2E2] hover:text-[#B91C1C]"
										>
											<Trash2 class="h-4 w-4" />
										</button>
									</div>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
		<Pagination bind:currentPage totalItems={services.length} {pageSize} />
	</div>
</div>

<!-- Modal Tambah / Edit Layanan -->
<Modal open={formModalOpen} title={editingId ? 'Edit Layanan' : 'Tambah Layanan Baru'} on:close={() => (formModalOpen = false)}>
	<form on:submit|preventDefault={handleSave} class="space-y-4">
		{#if formError}
			<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{formError}</div>
		{/if}
		<div>
			<label for="svc-name" class="block text-sm font-semibold text-[#0F172A]">Nama Layanan</label>
			<input
				id="svc-name"
				type="text"
				bind:value={form.name}
				required
				placeholder="Nama Layanan"
				class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
			/>
		</div>
		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<label for="svc-unit" class="block text-sm font-semibold text-[#0F172A]">Tipe Satuan</label>
				<select
					id="svc-unit"
					bind:value={form.unit_type}
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				>
					<option value="PER_KG">Per Kg</option>
					<option value="PER_ITEM">Per Item</option>
					<option value="PER_M2">Per m²</option>
				</select>
			</div>
			<div>
				<label for="svc-price" class="block text-sm font-semibold text-[#0F172A]">Harga (Rp)</label>
				<input
					id="svc-price"
					type="number"
					min="0"
					step="500"
					bind:value={form.price}
					required
					placeholder="Nominal Rp"
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				/>
			</div>
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
				{saving ? 'Menyimpan...' : 'Simpan Layanan'}
			</button>
		</div>
	</form>
</Modal>

<!-- Modal Detail Layanan -->
<Modal open={detailModalOpen} title="Detail Layanan" on:close={() => (detailModalOpen = false)}>
	{#if selectedService}
		<div class="space-y-3 text-sm">
			<div class="flex justify-between border-b border-[#E2E8F0] pb-2">
				<span class="text-[#64748B]">Nama Layanan</span>
				<span class="font-bold text-[#0F172A]">{selectedService.name}</span>
			</div>
			<div class="flex justify-between border-b border-[#E2E8F0] pb-2">
				<span class="text-[#64748B]">Tipe Satuan</span>
				<Badge variant="info">{unitLabel(selectedService.unit_type)}</Badge>
			</div>
			<div class="flex justify-between border-b border-[#E2E8F0] pb-2">
				<span class="text-[#64748B]">Tarif</span>
				<span class="font-bold text-[#0F172A]">{formatRupiah(selectedService.price)}</span>
			</div>
			<div class="flex justify-between">
				<span class="text-[#64748B]">Status</span>
				<Badge variant={selectedService.is_active ? 'success' : 'danger'}>
					{selectedService.is_active ? 'Aktif' : 'Nonaktif'}
				</Badge>
			</div>
		</div>
	{/if}
</Modal>
