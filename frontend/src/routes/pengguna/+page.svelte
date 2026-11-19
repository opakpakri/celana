<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { apiFetch, formatShortDate } from '$lib/api.js';
	import { authStore } from '$lib/stores/authStore.js';
	import { getCache, setCache } from '$lib/stores/dataCache.js';
	import Badge from '$lib/components/Badge.svelte';
	import Modal from '$lib/components/Modal.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import { Plus, Search, Pencil, Trash2, Shield, User as UserIcon } from 'lucide-svelte';

	const cachedUsers = getCache('admin_users');
	let users = cachedUsers || [];
	let search = '';
	let loading = !cachedUsers;
	let error = '';
	let currentPage = 1;
	const pageSize = 10;

	let formModalOpen = false;
	let editingId = null;
	let form = {
		name: '',
		email: '',
		password: '',
		role: 'USER',
		business_name: '',
		phone: '',
		address: ''
	};
	let saving = false;
	let formError = '';

	$: paginatedUsers = users.slice((currentPage - 1) * pageSize, currentPage * pageSize);

	async function loadUsers() {
		if (users.length === 0) loading = true;
		error = '';
		try {
			const q = search.trim() ? `?search=${encodeURIComponent(search.trim())}` : '';
			users = await apiFetch(`/admin/users${q}`);
			if (!search.trim()) setCache('admin_users', users);
			currentPage = 1;
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		if ($authStore.user?.role !== 'ADMIN') {
			goto('/');
			return;
		}
		loadUsers();
	});

	function openAdd() {
		editingId = null;
		form = {
			name: '',
			email: '',
			password: '',
			role: 'USER',
			business_name: '',
			phone: '',
			address: ''
		};
		formError = '';
		formModalOpen = true;
	}

	function openEdit(u) {
		editingId = u.id;
		form = {
			name: u.name,
			email: u.email,
			password: '',
			role: u.role || 'USER',
			business_name: u.business_name || '',
			phone: u.phone || '',
			address: u.address || ''
		};
		formError = '';
		formModalOpen = true;
	}

	async function handleSave() {
		formError = '';
		saving = true;
		try {
			const payload = { ...form };
			if (editingId && !payload.password) {
				delete payload.password;
			}
			if (editingId) {
				await apiFetch(`/admin/users/${editingId}`, {
					method: 'PUT',
					body: JSON.stringify(payload)
				});
			} else {
				await apiFetch('/admin/users', {
					method: 'POST',
					body: JSON.stringify(payload)
				});
			}
			formModalOpen = false;
			await loadUsers();
		} catch (err) {
			formError = err.message;
		} finally {
			saving = false;
		}
	}

	async function handleDelete(u) {
		if (!confirm(`Hapus akun "${u.name}" (${u.email})? Seluruh data outlet milik akun ini akan ikut terhapus.`)) return;
		error = '';
		try {
			await apiFetch(`/admin/users/${u.id}`, { method: 'DELETE' });
			await loadUsers();
		} catch (err) {
			error = err.message;
		}
	}
</script>

<div class="space-y-6">
	<div class="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-2xl font-extrabold text-[#0F172A]">Manajemen Pengguna</h1>
			<p class="text-sm text-[#64748B]">Menu khusus Admin untuk membuat dan mengelola akun outlet pengguna.</p>
		</div>
		<button
			type="button"
			data-tour="users-add-btn"
			on:click={openAdd}
			class="inline-flex items-center justify-center gap-2 rounded-xl bg-[#0284C7] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#0369A1]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah Pengguna</span>
		</button>
	</div>

	<form on:submit|preventDefault={loadUsers} class="flex gap-2">
		<div class="relative flex-1 max-w-md">
			<Search class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
			<input
				type="text"
				bind:value={search}
				placeholder="Cari nama, email, atau nama outlet..."
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

	<div data-tour="users-table" class="overflow-hidden rounded-2xl border border-[#E2E8F0] bg-white shadow-sm">
		<div class="overflow-x-auto">
			<table class="w-full text-left text-sm">
				<thead class="border-b border-[#E2E8F0] bg-[#F8FAFC] text-xs font-semibold uppercase tracking-wider text-[#64748B]">
					<tr>
						<th class="px-5 py-3.5">Nama</th>
						<th class="px-5 py-3.5">Email</th>
						<th class="px-5 py-3.5">Nama Outlet</th>
						<th class="px-5 py-3.5">Role</th>
						<th class="px-5 py-3.5">Dibuat</th>
						<th class="px-5 py-3.5 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[#E2E8F0]">
					{#if paginatedUsers.length === 0}
						<tr><td colspan="6" class="px-5 py-8 text-center text-[#64748B]">Tidak ada pengguna ditemukan.</td></tr>
					{:else}
						{#each paginatedUsers as u}
							<tr class="hover:bg-slate-50/70">
								<td class="px-5 py-3.5 font-semibold text-[#0F172A]">{u.name}</td>
								<td class="px-5 py-3.5 text-[#64748B]">{u.email}</td>
								<td class="px-5 py-3.5 text-[#0F172A]">{u.business_name}</td>
								<td class="px-5 py-3.5">
									{#if u.role === 'ADMIN'}
										<Badge variant="info">
											<Shield class="mr-1 h-3 w-3" />
											ADMIN
										</Badge>
									{:else}
										<Badge variant="neutral">
											<UserIcon class="mr-1 h-3 w-3" />
											USER
										</Badge>
									{/if}
								</td>
								<td class="px-5 py-3.5 text-[#64748B]">{formatShortDate(u.created_at)}</td>
								<td class="px-5 py-3.5 text-right">
									<div class="flex items-center justify-end gap-1.5">
										<button
											type="button"
											on:click={() => openEdit(u)}
											title="Edit"
											class="rounded-lg p-2 text-[#64748B] hover:bg-[#EFF6FF] hover:text-[#1D4ED8]"
										>
											<Pencil class="h-4 w-4" />
										</button>
										{#if u.id !== $authStore.user?.id}
											<button
												type="button"
												on:click={() => handleDelete(u)}
												title="Hapus"
												class="rounded-lg p-2 text-[#64748B] hover:bg-[#FEE2E2] hover:text-[#B91C1C]"
											>
												<Trash2 class="h-4 w-4" />
											</button>
										{/if}
									</div>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
		<Pagination bind:currentPage totalItems={users.length} {pageSize} />
	</div>
</div>

<!-- Modal Tambah / Edit Pengguna -->
<Modal open={formModalOpen} title={editingId ? 'Edit Pengguna' : 'Tambah Pengguna Baru'} maxWidth="max-w-lg" on:close={() => (formModalOpen = false)}>
	<form on:submit|preventDefault={handleSave} class="space-y-4">
		{#if formError}
			<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{formError}</div>
		{/if}
		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<label for="u-name" class="block text-sm font-semibold text-[#0F172A]">Nama Lengkap</label>
				<input
					id="u-name"
					type="text"
					bind:value={form.name}
					required
					placeholder="Nama User"
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7]"
				/>
			</div>
			<div>
				<label for="u-role" class="block text-sm font-semibold text-[#0F172A]">Role</label>
				<select
					id="u-role"
					bind:value={form.role}
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] bg-white px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7]"
				>
					<option value="USER">User (Pemilik Outlet)</option>
					<option value="ADMIN">Admin</option>
				</select>
			</div>
		</div>

		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<label for="u-email" class="block text-sm font-semibold text-[#0F172A]">Email Login</label>
				<input
					id="u-email"
					type="email"
					bind:value={form.email}
					required
					placeholder="Email User"
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7]"
				/>
			</div>
			<div>
				<label for="u-biz" class="block text-sm font-semibold text-[#0F172A]">Nama Outlet</label>
				<input
					id="u-biz"
					type="text"
					bind:value={form.business_name}
					required
					placeholder="Nama Outlet"
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7]"
				/>
			</div>
		</div>

		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<label for="u-phone" class="block text-sm font-semibold text-[#0F172A]">No. WhatsApp</label>
				<input
					id="u-phone"
					type="text"
					bind:value={form.phone}
					placeholder="No Telp User"
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7]"
				/>
			</div>
			<div>
				<label for="u-pass" class="block text-sm font-semibold text-[#0F172A]">
					Password {editingId ? '(Opsional)' : '(min. 6)'}
				</label>
				<input
					id="u-pass"
					type="password"
					bind:value={form.password}
					required={!editingId}
					minlength="6"
					placeholder={editingId ? 'Kosongkan jika tidak diubah' : '••••••••'}
					class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7]"
				/>
			</div>
		</div>

		<div>
			<label for="u-addr" class="block text-sm font-semibold text-[#0F172A]">Alamat Outlet</label>
			<input
				id="u-addr"
				type="text"
				bind:value={form.address}
				placeholder="Alamat User"
				class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7]"
			/>
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
				{saving ? 'Menyimpan...' : 'Simpan Pengguna'}
			</button>
		</div>
	</form>
</Modal>
