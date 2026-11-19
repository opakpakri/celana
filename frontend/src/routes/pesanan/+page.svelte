<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		apiFetch,
		formatRupiah,
		formatDate,
		unitLabel,
		orderStatusLabel,
		paymentStatusLabel
	} from '$lib/api.js';
	import { authStore } from '$lib/stores/authStore.js';
	import { getCache, setCache } from '$lib/stores/dataCache.js';
	import Badge from '$lib/components/Badge.svelte';
	import Modal from '$lib/components/Modal.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import {
		Plus,
		Search,
		Trash2,
		Eye,
		Printer,
		UserPlus,
		CreditCard,
		MessageCircle,
		ChevronDown,
		Check
	} from 'lucide-svelte';

	const tabs = [
		{ key: 'ALL', label: 'Semua' },
		{ key: 'PENDING', label: 'Antre' },
		{ key: 'PROCESSING', label: 'Proses' },
		{ key: 'READY', label: 'Siap Diambil' },
		{ key: 'COMPLETED', label: 'Selesai' }
	];

	const cachedOrders = getCache('orders');
	let activeStatus = 'ALL';
	let search = '';
	let orders = cachedOrders || [];
	let loading = !cachedOrders;
	let error = '';
	let currentPage = 1;
	const pageSize = 10;

	$: paginatedOrders = orders.slice((currentPage - 1) * pageSize, currentPage * pageSize);

	// Create order modal state
	let createModalOpen = false;
	let customers = [];
	let services = [];
	let selectedCustomerId = '';
	let customerDropdownOpen = false;
	let customerSearch = '';
	let quickAddOpen = false;
	let quickCustomer = { name: '', phone: '', address: '' };
	let orderItems = [{ service_id: '', quantity: 1 }];
	let orderNotes = '';
	let paymentMethod = 'TUNAI';
	let paidAmount = 0;
	let creating = false;
	let createError = '';

	// Detail & Pay modal state
	let detailModalOpen = false;
	let selectedOrder = null;
	let payAmount = 0;
	let payMethod = 'TUNAI';
	let updating = false;
	let detailError = '';

	async function loadOrders(silent = false) {
		if (!silent && orders.length === 0) loading = true;
		error = '';
		try {
			const params = new URLSearchParams();
			if (activeStatus !== 'ALL') params.set('status', activeStatus);
			if (search.trim()) params.set('search', search.trim());
			const q = params.toString() ? `?${params.toString()}` : '';
			orders = await apiFetch(`/orders${q}`);
			if (activeStatus === 'ALL' && !search.trim()) {
				setCache('orders', orders);
			}
			currentPage = 1;
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	async function loadCatalogs() {
		const [custRes, svcRes] = await Promise.all([
			apiFetch('/customers'),
			apiFetch('/services')
		]);
		customers = custRes || [];
		services = (svcRes.services || []).filter((s) => s.is_active);
	}

	onMount(async () => {
		await loadOrders();
		if ($page.url.searchParams.get('new') === '1') {
			await openCreateModal();
		}
	});

	async function openCreateModal() {
		createError = '';
		quickAddOpen = false;
		quickCustomer = { name: '', phone: '', address: '' };
		selectedCustomerId = '';
		customerDropdownOpen = false;
		customerSearch = '';
		orderItems = [{ service_id: '', quantity: 1, search: '', open: false }];
		orderNotes = '';
		paymentMethod = 'TUNAI';
		paidAmount = 0;
		createModalOpen = true;
		try {
			await loadCatalogs();
		} catch (err) {
			createError = err.message;
		}
	}

	function addItemRow() {
		orderItems = [...orderItems, { service_id: '', quantity: 1, search: '', open: false }];
	}

	function removeItemRow(idx) {
		if (orderItems.length === 1) return;
		orderItems = orderItems.filter((_, i) => i !== idx);
	}

	function filterServices(query) {
		const q = (query || '').trim().toLowerCase();
		if (!q) return services;
		return services.filter(
			(s) =>
				s.name.toLowerCase().includes(q) ||
				unitLabel(s.unit_type).toLowerCase().includes(q)
		);
	}

	function selectServiceForRow(idx, s) {
		orderItems[idx].service_id = s ? String(s.id) : '';
		orderItems[idx].search = '';
		orderItems[idx].open = false;
		orderItems = orderItems;
	}

	function toggleServiceDropdown(idx, event) {
		event.stopPropagation();
		customerDropdownOpen = false;
		orderItems = orderItems.map((r, i) => ({
			...r,
			open: i === idx ? !r.open : false,
			search: i === idx && !r.open ? '' : r.search
		}));
	}

	function toggleCustomerDropdown(event) {
		event.stopPropagation();
		if (orderItems.some((r) => r.open)) {
			orderItems = orderItems.map((r) => ({ ...r, open: false }));
		}
		customerDropdownOpen = !customerDropdownOpen;
		if (customerDropdownOpen) customerSearch = '';
	}

	function selectCustomer(c) {
		selectedCustomerId = String(c.id);
		customerSearch = '';
		customerDropdownOpen = false;
	}

	function closeAllServiceDropdowns(event) {
		if (event?.target?.closest?.('[data-dropdown]')) return;
		customerDropdownOpen = false;
		if (orderItems.some((r) => r.open)) {
			orderItems = orderItems.map((r) => ({ ...r, open: false }));
		}
	}

	function focusOnMount(node) {
		node.focus();
	}

	$: selectedCustomer = customers.find((c) => String(c.id) === String(selectedCustomerId));
	$: filteredCustomers = (() => {
		const q = customerSearch.trim().toLowerCase();
		if (!q) return customers;
		return customers.filter(
			(c) => c.name.toLowerCase().includes(q) || (c.phone || '').toLowerCase().includes(q)
		);
	})();

	function getService(id) {
		return services.find((s) => s.id === Number(id));
	}

	$: computedTotal = orderItems.reduce((sum, row) => {
		const svc = getService(row.service_id);
		const qty = Number(row.quantity) || 0;
		return sum + (svc ? svc.price * qty : 0);
	}, 0);

	$: changeAmount = Math.max(0, Number(paidAmount || 0) - computedTotal);
	$: remainingAmount = Math.max(0, computedTotal - Number(paidAmount || 0));

	async function handleQuickAddCustomer() {
		createError = '';
		try {
			const created = await apiFetch('/customers', {
				method: 'POST',
				body: JSON.stringify(quickCustomer)
			});
			customers = [created, ...customers];
			selectedCustomerId = String(created.id);
			quickAddOpen = false;
			quickCustomer = { name: '', phone: '', address: '' };
		} catch (err) {
			createError = err.message;
		}
	}

	async function handleCreateOrder() {
		createError = '';
		const validItems = orderItems
			.filter((r) => r.service_id && Number(r.quantity) > 0)
			.map((r) => ({ service_id: Number(r.service_id), quantity: Number(r.quantity) }));

		if (!selectedCustomerId || validItems.length === 0) {
			createError = 'Pilih pelanggan dan minimal satu layanan dengan jumlah > 0.';
			return;
		}

		creating = true;
		try {
			await apiFetch('/orders', {
				method: 'POST',
				body: JSON.stringify({
					customer_id: Number(selectedCustomerId),
					items: validItems,
					notes: orderNotes,
					paid_amount: Number(paidAmount || 0),
					payment_method: paymentMethod
				})
			});
			createModalOpen = false;
			await loadOrders();
		} catch (err) {
			createError = err.message;
		} finally {
			creating = false;
		}
	}

	async function openDetail(id) {
		detailError = '';
		detailModalOpen = true;
		selectedOrder = null;
		try {
			selectedOrder = await apiFetch(`/orders/${id}`);
			payAmount = Math.max(0, selectedOrder.total_amount - selectedOrder.paid_amount);
			payMethod = 'TUNAI';
		} catch (err) {
			detailError = err.message;
		}
	}

	function canTransition(currentStatus, targetStatus) {
		if (currentStatus === 'COMPLETED' || currentStatus === 'CANCELLED') return false;
		if (targetStatus === 'CANCELLED') return true;
		if (currentStatus === 'PENDING') return targetStatus === 'PROCESSING';
		if (currentStatus === 'PROCESSING') return targetStatus === 'READY';
		if (currentStatus === 'READY') return targetStatus === 'COMPLETED';
		return false;
	}

	async function handleUpdateStatus(newStatus) {
		if (!selectedOrder) return;
		if (newStatus === 'CANCELLED') {
			if (!confirm(`Batalkan pesanan ${selectedOrder.order_code}? Dana yang sudah dibayar akan dicatat sebagai pengembalian di Keuangan.`)) {
				return;
			}
		}
		updating = true;
		detailError = '';
		try {
			await apiFetch(`/orders/${selectedOrder.id}/status`, {
				method: 'PATCH',
				body: JSON.stringify({ status: newStatus })
			});
			selectedOrder = await apiFetch(`/orders/${selectedOrder.id}`);
			await loadOrders();
		} catch (err) {
			detailError = err.message;
		} finally {
			updating = false;
		}
	}

	function normalizeWhatsAppPhone(phone) {
		const digits = String(phone || '').replace(/\D/g, '');
		if (!digits) return '';
		if (digits.startsWith('0')) return '62' + digits.slice(1);
		if (digits.startsWith('8')) return '62' + digits;
		return digits;
	}

	function formatWaDateTime(iso) {
		if (!iso) return '-';
		const d = new Date(iso);
		const pad = (n) => String(n).padStart(2, '0');
		return `${pad(d.getDate())}/${pad(d.getMonth() + 1)}/${d.getFullYear()} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
	}

	function formatWaRupiah(amount) {
		return 'Rp' + new Intl.NumberFormat('id-ID', { maximumFractionDigits: 0 }).format(Number(amount || 0));
	}

	function unitShortLabel(unitType) {
		if (unitType === 'PER_KG') return 'KG';
		if (unitType === 'PER_ITEM') return 'PCS';
		if (unitType === 'PER_M2') return 'M2';
		return unitType || 'ITEM';
	}

	function sendWhatsAppNotification() {
		if (!selectedOrder || selectedOrder.status !== 'READY') return;
		const phone = normalizeWhatsAppPhone(selectedOrder.customer_phone);
		if (!phone) {
			detailError = 'Nomor telepon pelanggan tidak valid untuk WhatsApp.';
			return;
		}

		const biz = $authStore.user?.business_name || 'Celana Laundry';
		const bizAddress = $authStore.user?.address || '-';
		const bizPhone = $authStore.user?.phone || '-';
		const items = selectedOrder.items || [];
		const itemsText = items
			.map(
				(it) =>
					`\u2705 ${it.service_name}, ${it.quantity} ${unitShortLabel(it.unit_type)}\n@ ${formatWaRupiah(it.unit_price)}, Total ${formatWaRupiah(it.subtotal)}`
			)
			.join('\n');

		const remaining = Math.max(0, selectedOrder.total_amount - selectedOrder.paid_amount);
		const paymentLine =
			selectedOrder.paid_amount > 0
				? `\u{1F4B5} Dibayar ${formatWaRupiah(selectedOrder.paid_amount)}${remaining > 0 ? ` (Sisa ${formatWaRupiah(remaining)})` : ''}`
				: `\u{1F4B5} Belum ada pembayaran (Tagihan ${formatWaRupiah(selectedOrder.total_amount)})`;

		const message = [
			biz,
			bizAddress,
			bizPhone,
			'',
			'Nomor Nota :',
			selectedOrder.order_code,
			'',
			'Pelanggan Yth :',
			selectedOrder.customer_name,
			'',
			`Terima : ${formatWaDateTime(selectedOrder.created_at)}`,
			`Selesai : ${formatWaDateTime(selectedOrder.updated_at)}`,
			'',
			'======================',
			'Detail pesanan:',
			'Layanan:',
			itemsText,
			`\u2139\uFE0F ${items.length} Barang`,
			`Ket : ${selectedOrder.notes || '-'}`,
			'',
			'==============',
			'Detail biaya :',
			`Total tagihan : ${formatWaRupiah(selectedOrder.total_amount)}`,
			`Grand total : ${formatWaRupiah(selectedOrder.total_amount)}`,
			'',
			'Pembayaran:',
			paymentLine,
			'',
			`Status: ${paymentStatusLabel(selectedOrder.payment_status)}`,
			'',
			'=================',
			'Syarat dan ketentuan:',
			`S&K mengikuti ${biz}`,
			`Claim dilayani 1 x 24 jam dgn menunjukkan nota print/WA, Layanan komplain WA: ${bizPhone}`
		].join('\n');

		window.open(`https://api.whatsapp.com/send?phone=${phone}&text=${encodeURIComponent(message)}`, '_blank');
	}

	async function handlePayRemaining() {
		if (!selectedOrder || Number(payAmount) <= 0) return;
		updating = true;
		detailError = '';
		try {
			await apiFetch(`/orders/${selectedOrder.id}/pay`, {
				method: 'POST',
				body: JSON.stringify({ amount: Number(payAmount), payment_method: payMethod })
			});
			selectedOrder = await apiFetch(`/orders/${selectedOrder.id}`);
			payAmount = Math.max(0, selectedOrder.total_amount - selectedOrder.paid_amount);
			await loadOrders();
		} catch (err) {
			detailError = err.message;
		} finally {
			updating = false;
		}
	}

	function statusVariant(status) {
		if (status === 'COMPLETED' || status === 'READY') return 'success';
		if (status === 'PROCESSING') return 'info';
		if (status === 'CANCELLED') return 'danger';
		return 'warning';
	}

	function printReceipt() {
		window.print();
	}
</script>

<svelte:window on:click={closeAllServiceDropdowns} />

<div class="space-y-6">
	<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">Daftar Pesanan</h1>
			<p class="text-sm text-[#64748B]">Kelola transaksi cucian, status pengerjaan, dan pembayaran pelanggan.</p>
		</div>
		<button
			type="button"
			data-tour="orders-create-btn"
			on:click={openCreateModal}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#0284C7] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#0369A1]"
		>
			<Plus class="h-4 w-4" />
			<span>Buat Pesanan Baru</span>
		</button>
	</div>

	<!-- Filter Status & Search -->
	<div data-tour="orders-filter" class="flex flex-col justify-between gap-3 lg:flex-row lg:items-center">
		<div class="flex flex-wrap gap-2">
			{#each tabs as tab}
				<button
					type="button"
					on:click={() => {
						activeStatus = tab.key;
						loadOrders();
					}}
					class="rounded-xl px-3.5 py-2 text-xs font-semibold transition {activeStatus === tab.key
						? 'bg-[#0284C7] text-white'
						: 'border border-[#E2E8F0] bg-white text-[#64748B] hover:bg-slate-50'}"
				>
					{tab.label}
				</button>
			{/each}
		</div>

		<form on:submit|preventDefault={loadOrders} class="flex gap-2">
			<div class="relative flex-1 sm:w-72">
				<Search class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
				<input
					type="text"
					bind:value={search}
					placeholder="Cari nota, nama, atau no. HP..."
					class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2 pl-10 pr-4 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				/>
			</div>
			<button
				type="submit"
				class="rounded-xl border border-[#E2E8F0] bg-white px-4 py-2 text-sm font-semibold text-[#0F172A] hover:bg-slate-50"
			>
				Cari
			</button>
		</form>
	</div>

	{#if error}
		<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{error}</div>
	{/if}

	<!-- Tabel Pesanan -->
	<div data-tour="orders-table" class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<div class="overflow-x-auto">
			<table class="w-full text-left text-xs">
				<thead class="border-b border-[#E2E8F0] bg-[#F8FAFC] text-xs font-semibold text-[#64748B]">
					<tr>
						<th class="px-3 py-2.5">No. Nota</th>
						<th class="px-3 py-2.5">Tanggal</th>
						<th class="px-3 py-2.5">Pelanggan</th>
						<th class="px-3 py-2.5">Layanan</th>
						<th class="px-3 py-2.5">Status Cucian</th>
						<th class="px-3 py-2.5">Status Pembayaran</th>
						<th class="px-3 py-2.5">Total</th>
						<th class="px-3 py-2.5 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[#E2E8F0]">
					{#if orders.length === 0}
						<tr><td colspan="8" class="px-3 py-6 text-center text-[#64748B]">Tidak ada pesanan ditemukan.</td></tr>
					{:else}
						{#each paginatedOrders as o}
							<tr class="hover:bg-slate-50/70">
								<td class="whitespace-nowrap px-3 py-2.5 font-bold text-[#0F172A]">{o.order_code}</td>
								<td class="whitespace-nowrap px-3 py-2.5 text-[#64748B]">{formatDate(o.created_at)}</td>
								<td class="whitespace-nowrap px-3 py-2.5 font-semibold text-[#0F172A]">{o.customer_name}</td>
								<td class="max-w-[180px] truncate px-3 py-2.5 text-[#64748B]">{o.items_summary || '-'}</td>
								<td class="whitespace-nowrap px-3 py-2.5">
									<Badge variant={statusVariant(o.status)}>{orderStatusLabel(o.status)}</Badge>
								</td>
								<td class="whitespace-nowrap px-3 py-2.5">
									<Badge variant={o.payment_status === 'PAID' ? 'success' : 'danger'}>
										{paymentStatusLabel(o.payment_status)}
									</Badge>
								</td>
								<td class="whitespace-nowrap px-3 py-2.5 font-bold text-[#0F172A]">{formatRupiah(o.total_amount)}</td>
								<td class="whitespace-nowrap px-3 py-2.5 text-right">
									<button
										type="button"
										on:click={() => openDetail(o.id)}
										class="inline-flex items-center gap-1 rounded-lg border border-[#E2E8F0] px-2.5 py-1 text-xs font-semibold text-[#0284C7] hover:bg-[#EFF6FF]"
									>
										<Eye class="h-3.5 w-3.5" />
										<span>Detail</span>
									</button>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
		<Pagination bind:currentPage totalItems={orders.length} {pageSize} />
	</div>
</div>

<!-- Modal Buat Pesanan Baru -->
<Modal open={createModalOpen} title="Buat Pesanan Baru" maxWidth="max-w-2xl" on:close={() => (createModalOpen = false)}>
	<form on:submit|preventDefault={handleCreateOrder} class="space-y-5">
		{#if createError}
			<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{createError}</div>
		{/if}

		<!-- Pilih Pelanggan + Tambah Cepat -->
		<div>
			<div class="flex items-center justify-between">
				<label for="order-cust" class="block text-sm font-semibold text-[#0F172A]">Pelanggan</label>
				<button
					type="button"
					on:click={() => (quickAddOpen = !quickAddOpen)}
					class="inline-flex items-center gap-1 text-xs font-semibold text-[#0284C7] hover:underline"
				>
					<UserPlus class="h-3.5 w-3.5" />
					<span>{quickAddOpen ? 'Tutup Tambah Cepat' : 'Tambah Pelanggan Cepat'}</span>
				</button>
			</div>

			{#if quickAddOpen}
				<div class="mt-2 space-y-2 rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-3">
					<div class="grid gap-2 sm:grid-cols-2">
						<input
							type="text"
							bind:value={quickCustomer.name}
							placeholder="Nama pelanggan"
							class="rounded-lg border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none focus:border-[#0284C7]"
						/>
						<input
							type="tel"
							inputmode="numeric"
							value={quickCustomer.phone}
							on:input={(e) => {
								const val = e.target.value.replace(/\D/g, '');
								quickCustomer.phone = val;
								e.target.value = val;
							}}
							placeholder="No. Telepon"
							class="rounded-lg border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none focus:border-[#0284C7]"
						/>
					</div>
					<div class="flex gap-2">
						<input
							type="text"
							bind:value={quickCustomer.address}
							placeholder="Alamat (opsional)"
							class="flex-1 rounded-lg border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none focus:border-[#0284C7]"
						/>
						<button
							type="button"
							on:click={handleQuickAddCustomer}
							class="rounded-lg bg-[#0284C7] px-3 py-2 text-xs font-semibold text-white hover:bg-[#0369A1]"
						>
							Simpan
						</button>
					</div>
				</div>
			{/if}

			<div class="relative mt-1.5" data-dropdown>
				<button
					id="order-cust"
					type="button"
					on:click={toggleCustomerDropdown}
					class="flex w-full items-center justify-between gap-2 rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-left text-sm outline-none transition focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
				>
					<span class="truncate {selectedCustomer ? 'font-medium text-[#0F172A]' : 'text-[#64748B]'}">
						{selectedCustomer ? `${selectedCustomer.name} (${selectedCustomer.phone})` : 'Pilih pelanggan'}
					</span>
					<ChevronDown class="h-4 w-4 shrink-0 text-[#64748B] transition {customerDropdownOpen ? 'rotate-180' : ''}" />
				</button>

				{#if customerDropdownOpen}
					<div class="absolute left-0 right-0 top-full z-30 mt-1.5 rounded-2xl border border-[#E2E8F0] bg-white p-2 shadow-xl">
						<div class="pb-2">
							<div class="relative">
								<Search class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[#64748B]" />
								<input
									type="text"
									bind:value={customerSearch}
									placeholder="Cari nama atau no. telp..."
									use:focusOnMount
									class="w-full rounded-xl border border-[#0284C7] bg-white py-1.5 pl-8 pr-3 text-xs text-[#0F172A] outline-none ring-2 ring-[#E0F2FE]"
								/>
							</div>
						</div>

						<div class="max-h-48 space-y-0.5 overflow-y-auto">
							<div class="px-3 py-1.5 text-xs font-semibold text-[#64748B] opacity-70 select-none">
								Pilih pelanggan
							</div>

							{#if filteredCustomers.length === 0}
								<div class="px-3 py-2.5 text-center text-xs text-[#64748B]">Pelanggan tidak ditemukan</div>
							{:else}
								{#each filteredCustomers as c}
									{@const isSelected = String(selectedCustomerId) === String(c.id)}
									<button
										type="button"
										on:click={() => selectCustomer(c)}
										class="flex w-full items-center justify-between gap-2 rounded-xl px-3 py-2 text-left text-xs transition {isSelected
											? 'bg-[#EFF6FF] font-semibold text-[#0284C7]'
											: 'text-[#0F172A] hover:bg-slate-50'}"
									>
										<span class="truncate">{c.name} ({c.phone})</span>
										{#if isSelected}
											<Check class="h-3.5 w-3.5 shrink-0 text-[#0284C7]" />
										{/if}
									</button>
								{/each}
							{/if}
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- Multi-Item Layanan -->
		<div>
			<div class="mb-2 flex items-center justify-between">
				<span class="text-sm font-semibold text-[#0F172A]">Item Layanan</span>
				<button
					type="button"
					on:click={addItemRow}
					class="inline-flex items-center gap-1 text-xs font-semibold text-[#0284C7] hover:underline"
				>
					<Plus class="h-3.5 w-3.5" />
					<span>Tambah Item</span>
				</button>
			</div>

			<div class="space-y-2.5">
				{#each orderItems as row, idx}
					{@const svc = getService(row.service_id)}
					{@const filteredServices = filterServices(row.search)}
					<div class="flex flex-col gap-2 rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-3 sm:flex-row sm:items-center">
						<div class="relative flex-1" data-dropdown>
							<button
								type="button"
								on:click={(e) => toggleServiceDropdown(idx, e)}
								class="flex w-full items-center justify-between gap-2 rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2 text-left text-sm outline-none transition focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
							>
								<span class="truncate {svc ? 'font-medium text-[#0F172A]' : 'text-[#64748B]'}">
									{svc ? `${svc.name} — ${formatRupiah(svc.price)} (${unitLabel(svc.unit_type)})` : 'Pilih layanan'}
								</span>
								<ChevronDown class="h-4 w-4 shrink-0 text-[#64748B] transition {row.open ? 'rotate-180' : ''}" />
							</button>

							{#if row.open}
								<div class="absolute left-0 right-0 top-full z-30 mt-1.5 rounded-2xl border border-[#E2E8F0] bg-white p-2 shadow-xl">
									<div class="flex items-center gap-1.5 pb-2">
										<div class="relative flex-1">
											<Search class="pointer-events-none absolute left-3 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[#64748B]" />
											<input
												type="text"
												bind:value={row.search}
												placeholder="Cari..."
												use:focusOnMount
												class="w-full rounded-xl border border-[#0284C7] bg-white py-1.5 pl-8 pr-3 text-xs text-[#0F172A] outline-none ring-2 ring-[#E0F2FE]"
											/>
										</div>
										<button
											type="button"
											title="Tambah layanan baru"
											on:click={() => {
												createModalOpen = false;
												goto('/layanan');
											}}
											class="flex h-8 w-8 shrink-0 items-center justify-center rounded-xl border border-[#E2E8F0] bg-white text-[#0F172A] transition hover:bg-[#EFF6FF] hover:text-[#0284C7]"
										>
											<Plus class="h-4 w-4" />
										</button>
									</div>

									<div class="max-h-48 space-y-0.5 overflow-y-auto">
										<div class="px-3 py-1.5 text-xs font-semibold text-[#64748B] opacity-70 select-none">
											Pilih layanan
										</div>

										{#if filteredServices.length === 0}
											<div class="px-3 py-2.5 text-center text-xs text-[#64748B]">Layanan tidak ditemukan</div>
										{:else}
											{#each filteredServices as s}
												{@const isSelected = row.service_id === String(s.id)}
												<button
													type="button"
													on:click={() => selectServiceForRow(idx, s)}
													class="flex w-full items-center justify-between gap-2 rounded-xl px-3 py-2 text-left text-xs transition {isSelected
														? 'bg-[#EFF6FF] font-semibold text-[#0284C7]'
														: 'text-[#0F172A] hover:bg-slate-50'}"
												>
													<span class="truncate">{s.name} — {formatRupiah(s.price)} ({unitLabel(s.unit_type)})</span>
													{#if isSelected}
														<Check class="h-3.5 w-3.5 shrink-0 text-[#0284C7]" />
													{/if}
												</button>
											{/each}
										{/if}
									</div>
								</div>
							{/if}
						</div>
						<div class="flex items-center gap-2">
							<input
								type="number"
								min="0.1"
								step="0.1"
								bind:value={row.quantity}
								required
								class="w-24 rounded-lg border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none focus:border-[#0284C7]"
							/>
							<span class="w-28 text-right text-sm font-bold text-[#0F172A]">
								{formatRupiah(svc ? svc.price * (Number(row.quantity) || 0) : 0)}
							</span>
							{#if orderItems.length > 1}
								<button
									type="button"
									on:click={() => removeItemRow(idx)}
									class="rounded-lg p-1.5 text-[#64748B] hover:bg-[#FEE2E2] hover:text-[#B91C1C]"
								>
									<Trash2 class="h-4 w-4" />
								</button>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		</div>

		<!-- Catatan -->
		<div>
			<label for="order-notes" class="block text-sm font-semibold text-[#0F172A]">Catatan Pesanan</label>
			<input
				id="order-notes"
				type="text"
				bind:value={orderNotes}
				placeholder="Catatan (Opsional)"
				class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
			/>
		</div>

		<!-- Ringkasan & Pembayaran Awal -->
		<div class="rounded-2xl border border-[#E2E8F0] bg-[#F8FAFC] p-4 space-y-4">
			<div class="flex items-center justify-between border-b border-[#E2E8F0] pb-3">
				<span class="text-sm font-semibold text-[#64748B]">Total Tagihan</span>
				<span class="text-xl font-extrabold text-[#0284C7]">{formatRupiah(computedTotal)}</span>
			</div>

			<div class="grid gap-3 sm:grid-cols-2">
				<div>
					<label for="pay-method" class="block text-xs font-semibold uppercase text-[#64748B]">Metode Bayar</label>
					<select
						id="pay-method"
						bind:value={paymentMethod}
						class="mt-1 w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none focus:border-[#0284C7]"
					>
						<option value="TUNAI">Tunai</option>
						<option value="TRANSFER">Transfer Bank</option>
						<option value="QRIS">QRIS</option>
					</select>
				</div>
				<div>
					<label for="pay-nominal" class="block text-xs font-semibold uppercase text-[#64748B]">Nominal Diterima (Opsional/DP)</label>
					<input
						id="pay-nominal"
						type="number"
						min="0"
						step="500"
						bind:value={paidAmount}
						class="mt-1 w-full rounded-xl border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none focus:border-[#0284C7]"
					/>
				</div>
			</div>

			<!-- Preset Nominal -->
			<div class="flex flex-wrap gap-2">
				<button
					type="button"
					on:click={() => (paidAmount = computedTotal)}
					class="rounded-lg border border-[#E2E8F0] bg-white px-3 py-1.5 text-xs font-semibold text-[#0284C7] hover:bg-[#EFF6FF]"
				>
					Uang Pas ({formatRupiah(computedTotal)})
				</button>
				<button
					type="button"
					on:click={() => (paidAmount = 50000)}
					class="rounded-lg border border-[#E2E8F0] bg-white px-3 py-1.5 text-xs font-semibold text-[#0F172A] hover:bg-slate-50"
				>
					Rp 50.000
				</button>
				<button
					type="button"
					on:click={() => (paidAmount = 100000)}
					class="rounded-lg border border-[#E2E8F0] bg-white px-3 py-1.5 text-xs font-semibold text-[#0F172A] hover:bg-slate-50"
				>
					Rp 100.000
				</button>
			</div>

			<div class="flex justify-between text-xs font-semibold">
				<span class="text-[#64748B]">Kembalian: <strong class="text-[#15803D]">{formatRupiah(changeAmount)}</strong></span>
				<span class="text-[#64748B]">Sisa Tagihan: <strong class="text-[#B91C1C]">{formatRupiah(remainingAmount)}</strong></span>
			</div>
		</div>

		<div class="flex justify-end gap-2">
			<button
				type="button"
				on:click={() => (createModalOpen = false)}
				class="rounded-xl border border-[#E2E8F0] px-4 py-2 text-sm font-semibold text-[#64748B] hover:bg-slate-50"
			>
				Batal
			</button>
			<button
				type="submit"
				disabled={creating}
				class="rounded-xl bg-[#0284C7] px-5 py-2 text-sm font-semibold text-white hover:bg-[#0369A1] disabled:opacity-60"
			>
				{creating ? 'Menyimpan...' : 'Simpan Pesanan'}
			</button>
		</div>
	</form>
</Modal>

<!-- Modal Detail Pesanan & Nota -->
<Modal open={detailModalOpen} title="Detail Pesanan & Nota" maxWidth="max-w-2xl" on:close={() => (detailModalOpen = false)}>
	{#if detailError}
		<div class="mb-4 rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{detailError}</div>
	{/if}

	{#if !selectedOrder}
		<p class="py-8 text-center text-sm text-[#64748B]">Memuat detail pesanan...</p>
	{:else}
		<div class="space-y-5">
			<!-- Header Nota (untuk Print) -->
			<div class="rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-4">
				<div class="flex flex-wrap items-start justify-between gap-3">
					<div>
						<p class="text-xs font-bold uppercase text-[#0284C7]">{$authStore.user?.business_name || 'Celana Laundry'}</p>
						<h4 class="text-lg font-extrabold text-[#0F172A]">{selectedOrder.order_code}</h4>
						<p class="text-xs text-[#64748B]">{formatDate(selectedOrder.created_at)}</p>
					</div>
					<div class="flex items-center gap-2">
						<Badge variant={statusVariant(selectedOrder.status)}>{orderStatusLabel(selectedOrder.status)}</Badge>
						<Badge variant={selectedOrder.payment_status === 'PAID' ? 'success' : 'danger'}>
							{paymentStatusLabel(selectedOrder.payment_status)}
						</Badge>
					</div>
				</div>
				<div class="mt-3 border-t border-[#E2E8F0] pt-3 text-sm">
					<p class="font-semibold text-[#0F172A]">Pelanggan: {selectedOrder.customer_name} ({selectedOrder.customer_phone})</p>
					{#if selectedOrder.notes}
						<p class="mt-1 text-xs text-[#64748B]">Catatan: {selectedOrder.notes}</p>
					{/if}
				</div>
			</div>

			<!-- Rincian Item -->
			<div class="overflow-hidden rounded-xl border border-[#E2E8F0]">
				<table class="w-full text-left text-sm">
					<thead class="bg-[#F8FAFC] text-xs font-semibold text-[#64748B]">
						<tr>
							<th class="px-4 py-2.5">Layanan</th>
							<th class="px-4 py-2.5 text-right">Qty</th>
							<th class="px-4 py-2.5 text-right">Harga</th>
							<th class="px-4 py-2.5 text-right">Subtotal</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-[#E2E8F0]">
						{#each selectedOrder.items || [] as item}
							<tr>
								<td class="px-4 py-2.5 font-medium text-[#0F172A]">{item.service_name}</td>
								<td class="px-4 py-2.5 text-right text-[#64748B]">{item.quantity}</td>
								<td class="px-4 py-2.5 text-right text-[#64748B]">{formatRupiah(item.unit_price)}</td>
								<td class="px-4 py-2.5 text-right font-semibold text-[#0F172A]">{formatRupiah(item.subtotal)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<!-- Ringkasan Pembayaran -->
			<div class="space-y-1.5 rounded-xl border border-[#E2E8F0] p-4 text-sm">
				<div class="flex justify-between">
					<span class="text-[#64748B]">Total Tagihan</span>
					<span class="font-bold text-[#0F172A]">{formatRupiah(selectedOrder.total_amount)}</span>
				</div>
				<div class="flex justify-between">
					<span class="text-[#64748B]">Sudah Dibayar</span>
					<span class="font-bold text-[#15803D]">{formatRupiah(selectedOrder.paid_amount)}</span>
				</div>
				<div class="flex justify-between border-t border-[#E2E8F0] pt-1.5">
					<span class="font-semibold text-[#0F172A]">Sisa Tagihan</span>
					<span class="font-extrabold text-[#B91C1C]">
						{formatRupiah(Math.max(0, selectedOrder.total_amount - selectedOrder.paid_amount))}
					</span>
				</div>
			</div>

			<!-- Footer Nota Outlet -->
			<div class="text-center text-xs text-[#64748B]">
				{$authStore.user?.address || ''} • WhatsApp: {$authStore.user?.phone || '-'}
			</div>

			<!-- Kontrol Status & Pelunasan (no-print) -->
			<div class="no-print space-y-4 border-t border-[#E2E8F0] pt-4">
				<div>
					<p class="mb-2 text-xs font-semibold uppercase text-[#64748B]">Ubah Status Cucian</p>
					<div class="flex flex-wrap gap-2">
						{#each ['PENDING', 'PROCESSING', 'READY', 'COMPLETED', 'CANCELLED'] as st}
							{@const isCurrent = selectedOrder.status === st}
							{@const allowed = canTransition(selectedOrder.status, st)}
							<button
								type="button"
								disabled={updating || isCurrent || !allowed}
								on:click={() => handleUpdateStatus(st)}
								class="rounded-lg px-3 py-1.5 text-xs font-semibold transition {isCurrent
									? 'bg-[#0284C7] text-white'
									: allowed
										? st === 'CANCELLED'
											? 'border border-[#FEE2E2] bg-[#FEE2E2] text-[#B91C1C] hover:bg-red-200'
											: 'border border-[#0284C7] bg-[#EFF6FF] text-[#1D4ED8] hover:bg-sky-100'
										: 'cursor-not-allowed border border-[#E2E8F0] bg-slate-50 text-slate-400 opacity-60'}"
							>
								{orderStatusLabel(st)}
							</button>
						{/each}
					</div>
				</div>

				{#if selectedOrder.payment_status !== 'PAID'}
					<div class="rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-3">
						<p class="mb-2 flex items-center gap-1.5 text-xs font-semibold uppercase text-[#64748B]">
							<CreditCard class="h-3.5 w-3.5" />
							<span>Pelunasan / Tambah Pembayaran</span>
						</p>
						<div class="flex flex-col gap-2 sm:flex-row">
							<select
								bind:value={payMethod}
								class="rounded-lg border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none"
							>
								<option value="TUNAI">Tunai</option>
								<option value="TRANSFER">Transfer</option>
								<option value="QRIS">QRIS</option>
							</select>
							<input
								type="number"
								min="1"
								step="500"
								bind:value={payAmount}
								class="flex-1 rounded-lg border border-[#E2E8F0] bg-white px-3 py-2 text-sm outline-none"
							/>
							<button
								type="button"
								disabled={updating}
								on:click={handlePayRemaining}
								class="rounded-lg bg-[#0284C7] px-4 py-2 text-sm font-semibold text-white hover:bg-[#0369A1] disabled:opacity-60"
							>
								Bayar
							</button>
						</div>
					</div>
				{/if}

				<div class="flex flex-wrap items-center justify-end gap-2">
					<button
						type="button"
						disabled={selectedOrder.status !== 'READY'}
						on:click={sendWhatsAppNotification}
						class="inline-flex items-center gap-2 rounded-xl px-4 py-2 text-sm font-semibold transition {selectedOrder.status === 'READY'
							? 'bg-[#16A34A] text-white shadow-sm hover:bg-[#15803D]'
							: 'cursor-not-allowed border border-[#E2E8F0] bg-slate-100 text-slate-400 opacity-60'}"
					>
						<MessageCircle class="h-4 w-4" />
						<span>Kirim Notif WhatsApp</span>
					</button>
					<button
						type="button"
						on:click={printReceipt}
						class="inline-flex items-center gap-2 rounded-xl border border-[#E2E8F0] bg-white px-4 py-2 text-sm font-semibold text-[#0F172A] hover:bg-slate-50"
					>
						<Printer class="h-4 w-4" />
						<span>Cetak Nota</span>
					</button>
				</div>
			</div>
		</div>
	{/if}
</Modal>
