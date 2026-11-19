<script>
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import { apiFetch } from '$lib/api.js';
	import { authStore } from '$lib/stores/authStore.js';
	import { themeStore } from '$lib/stores/themeStore.js';
	import { Save, Sun, Moon, Monitor, UserCog, ArrowRight } from 'lucide-svelte';

	const initialUser = get(authStore).user;
	let form = {
		name: initialUser?.name || '',
		email: initialUser?.email || '',
		business_name: initialUser?.business_name || '',
		phone: initialUser?.phone || '',
		address: initialUser?.address || '',
		password: ''
	};
	let loading = !initialUser;
	let saving = false;
	let error = '';
	let success = '';

	onMount(async () => {
		try {
			const user = await apiFetch('/users/profile');
			form = {
				name: user.name || '',
				email: user.email || '',
				business_name: user.business_name || '',
				phone: user.phone || '',
				address: user.address || '',
				password: ''
			};
			authStore.updateUser(user);
		} catch (err) {
			if (!initialUser) error = err.message;
		} finally {
			loading = false;
		}
	});

	async function handleSave() {
		error = '';
		success = '';
		saving = true;
		try {
			const payload = { ...form };
			if (!payload.password) delete payload.password;
			const updated = await apiFetch('/users/profile', {
				method: 'PUT',
				body: JSON.stringify(payload)
			});
			authStore.updateUser(updated);
			form.password = '';
			success = 'Profil akun dan bisnis berhasil diperbarui.';
		} catch (err) {
			error = err.message;
		} finally {
			saving = false;
		}
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-extrabold text-[#0F172A]">Pengaturan Akun & Outlet</h1>
		<p class="text-sm text-[#64748B]">Kelola informasi pemilik, tampilan aplikasi, dan detail outlet yang tampil pada nota pesanan.</p>
	</div>

	<!-- Kartu Pengaturan Tema (Light / Dark Mode) -->
	<div data-tour="settings-theme" class="flex flex-col justify-between gap-4 rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-sm sm:flex-row sm:items-center">
		<div>
			<h2 class="text-base font-bold text-[#0F172A]">Tampilan Aplikasi</h2>
			<p class="text-sm text-[#64748B]">Pilih tema terang atau gelap sesuai kenyamanan Anda.</p>
		</div>
		<div class="grid grid-cols-3 rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] p-1 sm:w-80">
			<button
				type="button"
				on:click={() => themeStore.setTheme('light')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {$themeStore === 'light'
					? 'bg-white text-[#0284C7] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<Sun class="h-4 w-4" />
				<span>Terang</span>
			</button>
			<button
				type="button"
				on:click={() => themeStore.setTheme('dark')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {$themeStore === 'dark'
					? 'bg-white text-[#0284C7] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<Moon class="h-4 w-4" />
				<span>Gelap</span>
			</button>
			<button
				type="button"
				on:click={() => themeStore.setTheme('system')}
				class="flex items-center justify-center gap-1.5 rounded-lg py-2 text-xs font-semibold transition {$themeStore === 'system'
					? 'bg-white text-[#0284C7] shadow-sm'
					: 'text-[#64748B] hover:text-[#0F172A]'}"
			>
				<Monitor class="h-4 w-4" />
				<span>Sistem</span>
			</button>
		</div>
	</div>

	{#if $authStore.user?.role === 'ADMIN'}
		<a
			href="/pengguna"
			class="group flex items-center justify-between rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-sm transition hover:border-[#0284C7]"
		>
			<div class="flex items-center gap-4">
				<div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#E0F2FE] text-[#0369A1]">
					<UserCog class="h-5 w-5" />
				</div>
				<div>
					<h2 class="text-base font-bold text-[#0F172A]">Manajemen Pengguna</h2>
					<p class="text-sm text-[#64748B]">Buat, edit, dan kelola akun outlet pengguna (khusus Admin).</p>
				</div>
			</div>
			<ArrowRight class="h-5 w-5 text-[#64748B] transition group-hover:translate-x-1 group-hover:text-[#0284C7]" />
		</a>
	{/if}

	<div data-tour="settings-form" class="rounded-2xl border border-[#E2E8F0] bg-white p-6 shadow-sm sm:p-8">
		<form on:submit|preventDefault={handleSave} class="space-y-5">
				<div>
					<h2 class="text-base font-bold text-[#0F172A]">Pengaturan Akun</h2>
					<p class="text-sm text-[#64748B]">Informasi identitas pemilik dan email login akun Anda.</p>
				</div>

				{#if error}
					<div class="rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">{error}</div>
				{/if}
				{#if success}
					<div class="rounded-xl bg-[#DCFCE7] px-4 py-3 text-sm font-medium text-[#15803D]">{success}</div>
				{/if}

				<div class="grid gap-4 sm:grid-cols-2">
					<div>
						<label for="set-name" class="block text-sm font-semibold text-[#0F172A]">Nama Pemilik</label>
						<input
							id="set-name"
							type="text"
							bind:value={form.name}
							readonly
							class="mt-1.5 w-full cursor-not-allowed rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] px-3.5 py-2.5 text-sm text-[#64748B] outline-none"
						/>
					</div>
					<div>
						<label for="set-email" class="block text-sm font-semibold text-[#0F172A]">Email Login</label>
						<input
							id="set-email"
							type="email"
							bind:value={form.email}
							readonly
							class="mt-1.5 w-full cursor-not-allowed rounded-xl border border-[#E2E8F0] bg-[#F8FAFC] px-3.5 py-2.5 text-sm text-[#64748B] outline-none"
						/>
					</div>
				</div>

				<div class="border-t border-[#E2E8F0] pt-5">
					<h2 class="mb-4 text-base font-bold text-[#0F172A]">Profil Outlet / Toko (Footer Nota)</h2>
					<div class="grid gap-4 sm:grid-cols-2">
						<div>
							<label for="set-biz" class="block text-sm font-semibold text-[#0F172A]">Nama Laundry</label>
							<input
								id="set-biz"
								type="text"
								bind:value={form.business_name}
								required
								class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
							/>
						</div>
						<div>
							<label for="set-phone" class="block text-sm font-semibold text-[#0F172A]">Nomor WhatsApp Toko</label>
							<input
								id="set-phone"
								type="text"
								bind:value={form.phone}
								placeholder="081234567890"
								class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
							/>
						</div>
					</div>

					<div class="mt-4">
						<label for="set-addr" class="block text-sm font-semibold text-[#0F172A]">Alamat Toko</label>
						<textarea
							id="set-addr"
							rows="2"
							bind:value={form.address}
							class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
						></textarea>
					</div>
				</div>

				<div class="border-t border-[#E2E8F0] pt-5">
					<label for="set-pass" class="block text-sm font-semibold text-[#0F172A]">Password Baru (Opsional)</label>
					<input
						id="set-pass"
						type="password"
						bind:value={form.password}
						minlength="6"
						placeholder="Kosongkan jika tidak ingin mengubah password"
						class="mt-1.5 w-full rounded-xl border border-[#E2E8F0] px-3.5 py-2.5 text-sm outline-none focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
					/>
				</div>

				<div class="flex justify-end pt-2">
					<button
						type="submit"
						disabled={saving}
						class="inline-flex items-center gap-2 rounded-xl bg-[#0284C7] px-5 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#0369A1] disabled:opacity-60"
					>
						<Save class="h-4 w-4" />
					<span>{saving ? 'Menyimpan...' : 'Simpan Perubahan'}</span>
				</button>
			</div>
		</form>
	</div>
</div>
