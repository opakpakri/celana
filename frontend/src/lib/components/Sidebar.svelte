<script>
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/authStore.js';
	import { startProductTour } from '$lib/tour.js';
	import {
		LayoutDashboard,
		ShoppingBag,
		Sparkles,
		Users,
		Wallet,
		Settings,
		LogOut,
		Store,
		WashingMachine,
		UserCog,
		HelpCircle
	} from 'lucide-svelte';

	const navItems = [
		{ href: '/', label: 'Dashboard', icon: LayoutDashboard },
		{ href: '/pesanan', label: 'Pesanan', icon: ShoppingBag },
		{ href: '/layanan', label: 'Layanan', icon: Sparkles },
		{ href: '/pelanggan', label: 'Pelanggan', icon: Users },
		{ href: '/keuangan', label: 'Keuangan', icon: Wallet },
		{ href: '/pengguna', label: 'Pengguna', icon: UserCog, adminOnly: true },
		{ href: '/pengaturan', label: 'Pengaturan', icon: Settings }
	];

	$: isAdmin = $authStore.user?.role === 'ADMIN';
	$: visibleNavItems = navItems.filter((item) => !item.adminOnly || isAdmin);
	const mobileNavItems = navItems.slice(0, 5);

	$: currentPath = $page.url.pathname;

	function isActive(href, path) {
		if (href === '/') return path === '/';
		return path.startsWith(href);
	}

	function handleLogout() {
		authStore.clearSession();
		goto('/login');
	}
</script>

<!-- Desktop Sticky Sidebar -->
<aside
	class="no-print hidden w-60 shrink-0 flex-col bg-gradient-to-b from-[#0284C7] via-[#0369A1] to-[#0891B2] text-white shadow-sm lg:sticky lg:top-0 lg:flex lg:h-screen"
>
	<!-- Branding Header -->
	<div class="flex items-center justify-between border-b border-white/15 px-4 py-4">
		<a href="/" class="flex items-center gap-2.5">
			<div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-white/20 text-white backdrop-blur-sm">
				<WashingMachine class="h-5 w-5" />
			</div>
			<div class="min-w-0">
				<span class="block text-base font-extrabold tracking-tight text-white">Celana</span>
				<span class="block truncate text-[10px] font-medium text-sky-100">Cellix Laundry</span>
			</div>
		</a>
	</div>

	<!-- Active Tenant Business Card -->
	<div class="px-3 pt-3">
		<div class="rounded-xl border border-white/15 bg-white/10 p-3 backdrop-blur-sm">
			<div class="flex items-center gap-2">
				<div class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-white/20 text-white">
					<Store class="h-3.5 w-3.5" />
				</div>
				<div class="min-w-0">
					<p class="truncate text-xs font-bold text-white">
						{$authStore.user?.business_name || 'Celana Laundry'}
					</p>
					<p class="truncate text-[11px] text-sky-100">
						{$authStore.user?.name || 'Owner'}
					</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Main Navigation Items -->
	<nav data-tour="sidebar-nav" class="flex-1 space-y-1 overflow-y-auto px-3 py-3">
		{#each visibleNavItems as item}
			{@const active = isActive(item.href, currentPath)}
			<a
				href={item.href}
				class="flex items-center gap-2.5 rounded-xl px-3 py-2 text-sm font-semibold transition {active
					? 'bg-white text-[#0284C7] shadow-sm'
					: 'text-sky-100 hover:bg-white/10 hover:text-white'}"
			>
				<svelte:component this={item.icon} class="h-4 w-4 shrink-0" />
				<span>{item.label}</span>
			</a>
		{/each}
	</nav>

	<!-- Panduan & Logout Footer -->
	<div class="space-y-1 border-t border-white/15 p-3">
		<button
			type="button"
			on:click={() => startProductTour(currentPath)}
			class="flex w-full items-center gap-2.5 rounded-xl px-3 py-2 text-sm font-semibold text-sky-100 transition hover:bg-white/10 hover:text-white"
		>
			<HelpCircle class="h-4 w-4 shrink-0" />
			<span>Panduan Halaman</span>
		</button>
		<button
			type="button"
			on:click={handleLogout}
			class="flex w-full items-center gap-2.5 rounded-xl px-3 py-2 text-sm font-semibold text-sky-100 transition hover:bg-red-500/20 hover:text-white"
		>
			<LogOut class="h-4 w-4 shrink-0" />
			<span>Keluar</span>
		</button>
	</div>
</aside>

<!-- Mobile Bottom Navbar -->
<nav
	data-tour="mobile-nav"
	class="no-print fixed inset-x-0 bottom-0 z-40 grid grid-cols-5 gap-1.5 rounded-t-2xl bg-[#0284C7] px-2.5 py-2.5 shadow-[0_-4px_20px_rgba(2,132,199,0.22)] lg:hidden"
>
	{#each mobileNavItems as item}
		{@const active = isActive(item.href, currentPath)}
		<a
			href={item.href}
			class="flex flex-col items-center justify-center gap-1.5 rounded-xl py-2.5 text-xs font-semibold transition {active
				? 'bg-white text-[#0284C7] shadow-sm'
				: 'text-sky-100 hover:bg-white/10 hover:text-white'}"
		>
			<svelte:component this={item.icon} class="h-5 w-5 shrink-0" />
			<span class="truncate">{item.label}</span>
		</a>
	{/each}
</nav>
