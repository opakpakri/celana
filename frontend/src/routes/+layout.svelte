<script>
	import '../app.css';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { authStore } from '$lib/stores/authStore.js';
	import '$lib/stores/themeStore.js';
	import Sidebar from '$lib/components/Sidebar.svelte';
	import { startProductTour } from '$lib/tour.js';
	import { WashingMachine, Settings, LogOut, HelpCircle } from 'lucide-svelte';

	const publicRoutes = ['/login', '/register'];

	$: isPublic = publicRoutes.includes($page.url.pathname);
	$: isSettingsActive = $page.url.pathname.startsWith('/pengaturan');
	$: if (browser) {
		if (!$authStore.token && !isPublic) {
			goto('/login');
		} else if ($authStore.token && isPublic) {
			goto('/');
		}
	}

	function handleLogout() {
		authStore.clearSession();
		goto('/login');
	}
</script>

{#if isPublic}
	<slot />
{:else if $authStore.token}
	<div class="flex min-h-screen bg-[#F8FAFC]">
		<Sidebar />
		<div class="flex min-w-0 flex-1 flex-col">
			<!-- Mobile Header -->
			<header class="no-print sticky top-0 z-30 flex items-center justify-between border-b border-[#E2E8F0] bg-white px-4 py-3 lg:hidden">
				<a href="/" class="flex items-center gap-2.5">
					<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#0284C7] text-white">
						<WashingMachine class="h-5 w-5" />
					</div>
					<div>
						<span class="block text-sm font-extrabold leading-tight text-[#0284C7]">
							{$authStore.user?.business_name || 'Celana'}
						</span>
						<span class="block text-[11px] text-[#64748B]">Halo, {$authStore.user?.name || 'Owner'}</span>
					</div>
				</a>

				<div class="flex items-center gap-2">
					<button
						type="button"
						on:click={() => startProductTour($page.url.pathname)}
						title="Panduan Halaman"
						class="rounded-xl border border-[#E2E8F0] p-2 text-[#0284C7] transition hover:bg-[#EFF6FF]"
					>
						<HelpCircle class="h-4 w-4" />
					</button>
					<a
						href="/pengaturan"
						title="Pengaturan"
						class="rounded-xl border border-[#E2E8F0] p-2 transition {isSettingsActive
							? 'bg-[#EFF6FF] text-[#1D4ED8]'
							: 'text-[#64748B] hover:bg-slate-50 hover:text-[#0F172A]'}"
					>
						<Settings class="h-4 w-4" />
					</a>
					<button
						type="button"
						on:click={handleLogout}
						title="Keluar"
						class="rounded-xl border border-[#E2E8F0] p-2 text-[#64748B] transition hover:bg-[#FEE2E2] hover:text-[#B91C1C]"
					>
						<LogOut class="h-4 w-4" />
					</button>
				</div>
			</header>

			<main class="flex-1 p-4 pb-24 sm:p-6 sm:pb-24 lg:p-8 lg:pb-8">
				<div class="mx-auto max-w-7xl">
					<slot />
				</div>
			</main>
		</div>
	</div>
{/if}
