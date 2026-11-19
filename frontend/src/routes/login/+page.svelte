<script>
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api.js';
	import { authStore } from '$lib/stores/authStore.js';
	import { WashingMachine, Mail, Lock, ArrowRight } from 'lucide-svelte';

	let email = '';
	let password = '';
	let loading = false;
	let error = '';

	async function handleLogin() {
		error = '';
		loading = true;
		try {
			const res = await apiFetch('/auth/login', {
				method: 'POST',
				body: JSON.stringify({ email, password })
			});
			authStore.setSession(res.token, res.user);
			goto('/');
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-[#F8FAFC] px-4 py-10">
	<div class="w-full max-w-md rounded-2xl border border-[#E2E8F0] bg-white p-8 shadow-sm">
		<div class="mb-8 flex items-center gap-3">
			<div class="flex h-11 w-11 items-center justify-center rounded-xl bg-[#0284C7] text-white shadow-sm">
				<WashingMachine class="h-6 w-6" />
			</div>
			<div>
				<h1 class="text-xl font-extrabold text-[#0284C7]">Celana</h1>
				<p class="text-xs font-medium text-[#64748B]">Cellix Laundry Management</p>
			</div>
		</div>

		<h2 class="text-2xl font-bold text-[#0F172A]">Masuk ke Akun Outlet</h2>
		<p class="mt-1 text-sm text-[#64748B]">Kelola pesanan, layanan, pelanggan, dan arus kas outlet Anda.</p>

		{#if error}
			<div class="mt-4 rounded-xl bg-[#FEE2E2] px-4 py-3 text-sm font-medium text-[#B91C1C]">
				{error}
			</div>
		{/if}

		<form on:submit|preventDefault={handleLogin} class="mt-6 space-y-4">
			<div>
				<label for="email" class="block text-sm font-semibold text-[#0F172A]">Email</label>
				<div class="relative mt-1.5">
					<Mail class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
					<input
						id="email"
						type="email"
						bind:value={email}
						required
						placeholder="owner@laundry.com"
						class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-10 pr-4 text-sm text-[#0F172A] outline-none transition focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
					/>
				</div>
			</div>

			<div>
				<label for="password" class="block text-sm font-semibold text-[#0F172A]">Password</label>
				<div class="relative mt-1.5">
					<Lock class="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[#64748B]" />
					<input
						id="password"
						type="password"
						bind:value={password}
						required
						placeholder="••••••••"
						class="w-full rounded-xl border border-[#E2E8F0] bg-white py-2.5 pl-10 pr-4 text-sm text-[#0F172A] outline-none transition focus:border-[#0284C7] focus:ring-2 focus:ring-[#E0F2FE]"
					/>
				</div>
			</div>

			<button
				type="submit"
				disabled={loading}
				class="flex w-full items-center justify-center gap-2 rounded-xl bg-[#0284C7] px-4 py-2.5 text-sm font-semibold text-white transition hover:bg-[#0369A1] disabled:opacity-60"
			>
				<span>{loading ? 'Memproses...' : 'Masuk'}</span>
				<ArrowRight class="h-4 w-4" />
			</button>
		</form>

		<p class="mt-6 text-center text-xs text-[#64748B]">
			Pendaftaran akun baru hanya dapat dibuat oleh Admin.
		</p>
	</div>
</div>
