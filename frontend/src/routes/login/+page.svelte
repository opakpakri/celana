<script>
	import { goto } from '$app/navigation';
	import { apiFetch } from '$lib/api.js';
	import { authStore } from '$lib/stores/authStore.js';
	import {
		WashingMachine,
		Mail,
		Lock,
		ArrowRight
	} from 'lucide-svelte';

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

<div class="relative flex min-h-screen items-center justify-center overflow-hidden bg-gradient-to-br from-[#E0F2FE] via-[#F0F9FF] to-[#F8FAFC] px-4 py-10">
	<!-- Background Ambient Atmosphere (Fresh Air Glows) -->
	<div class="pointer-events-none absolute -left-20 -top-20 h-96 w-96 rounded-full bg-[#38BDF8]/20 blur-3xl animate-pulse-glow"></div>
	<div class="pointer-events-none absolute -bottom-24 -right-20 h-96 w-96 rounded-full bg-[#0284C7]/15 blur-3xl animate-pulse-glow-delayed"></div>
	<div class="pointer-events-none absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 h-[500px] w-[500px] rounded-full bg-white/40 blur-2xl"></div>

	<!-- Tali Jemuran & Pakaian Bersih (100% Menempel Presisi di Kawat Jemuran) -->
	<div class="pointer-events-none absolute inset-x-0 top-0 h-48 md:h-56 select-none overflow-hidden">
		<svg
			viewBox="0 0 1200 170"
			class="w-full h-full drop-shadow-sm"
			preserveAspectRatio="xMidYMin meet"
			fill="none"
		>
			<!-- Tali Kawat Jemuran Melengkung Alami -->
			<path
				d="M -50,30 Q 600,72 1250,30"
				stroke="#94A3B8"
				stroke-width="2"
				stroke-linecap="round"
				stroke-dasharray="6 3"
				opacity="0.8"
			/>

			<!-- 1. KAOS T-SHIRT BIRU (Tepat di Tali: X=150, Y=45) -->
			<g class="garment-1" style="transform-origin: 150px 45px;">
				<!-- Jepitan Kayu 1 (Kiri) menjepit tali dan bahu kaos -->
				<rect x="135" y="37" width="4" height="15" rx="1" fill="#D97706" />
				<!-- Jepitan Kayu 2 (Kanan) -->
				<rect x="161" y="37" width="4" height="15" rx="1" fill="#D97706" />

				<!-- Kaos T-Shirt -->
				<path
					d="M 132,45 L 140,51 C 144,48 156,48 160,51 L 168,45 L 182,54 L 175,66 L 169,63 L 169,92 L 131,92 L 131,63 L 125,66 L 118,54 Z"
					fill="#0284C7"
					class="drop-shadow-md"
				/>
				<!-- Garis Leher & Kerutan Halus -->
				<path d="M 140,51 C 144,48 156,48 160,51" stroke="#BAE6FD" stroke-width="2" stroke-linecap="round" />
				<path d="M 144,58 L 144,88 M 156,58 L 156,88" stroke="#0369A1" stroke-width="1.2" stroke-dasharray="2 2" opacity="0.6" />
			</g>

			<!-- 2. CELANA PANJANG / JEANS (Tepat di Tali: X=350, Y=56) -->
			<g class="garment-2" style="transform-origin: 350px 56px;">
				<!-- Jepitan Kayu 1 (Kiri Pinggang Celana) -->
				<rect x="333" y="48" width="4" height="15" rx="1" fill="#B45309" />
				<!-- Jepitan Kayu 2 (Kanan Pinggang Celana) -->
				<rect x="363" y="48" width="4" height="15" rx="1" fill="#B45309" />

				<!-- Celana Panjang -->
				<path
					d="M 330,56 L 370,56 L 366,118 L 354,118 L 350,78 L 346,118 L 334,118 Z"
					fill="#1D4ED8"
					class="drop-shadow-md"
				/>
				<!-- Jahitan Ban Pinggang & Saku -->
				<path d="M 330,64 L 370,64" stroke="#60A5FA" stroke-width="1.8" />
				<path d="M 335,69 Q 342,73 345,69" stroke="#93C5FD" stroke-width="1.2" />
				<path d="M 365,69 Q 358,73 355,69" stroke="#93C5FD" stroke-width="1.2" />
				<path d="M 350,64 L 350,78" stroke="#93C5FD" stroke-width="1.2" />
			</g>

			<!-- 3. HANDUK BERSIH BERGARIS (Tepat di Tali: X=850, Y=56) -->
			<g class="garment-3" style="transform-origin: 850px 56px;">
				<!-- Jepitan Kayu 1 (Kiri Handuk) -->
				<rect x="832" y="48" width="4" height="15" rx="1" fill="#D97706" />
				<!-- Jepitan Kayu 2 (Kanan Handuk) -->
				<rect x="864" y="48" width="4" height="15" rx="1" fill="#D97706" />

				<!-- Handuk Draped -->
				<rect x="828" y="56" width="44" height="58" rx="3" fill="#0D9488" class="drop-shadow-md" />
				<rect x="828" y="96" width="44" height="4" fill="#5EEAD4" opacity="0.8" />
				<rect x="828" y="103" width="44" height="2" fill="#CCFBF1" opacity="0.9" />
				<!-- Tekstur Serat Handuk -->
				<line x1="833" y1="64" x2="867" y2="64" stroke="#14B8A6" stroke-width="1.2" stroke-dasharray="2 2" />
				<line x1="833" y1="74" x2="867" y2="74" stroke="#14B8A6" stroke-width="1.2" stroke-dasharray="2 2" />
				<line x1="833" y1="84" x2="867" y2="84" stroke="#14B8A6" stroke-width="1.2" stroke-dasharray="2 2" />
			</g>

			<!-- 4. KEMEJA DI HANGER (Kait Gantungan Mengait Pas di Tali: X=1050, Y=45) -->
			<g class="garment-4" style="transform-origin: 1050px 45px;">
				<!-- Kait Hanger Logam yang mengait langsung pada kawat jemuran Y=45 -->
				<path
					d="M 1047,45 C 1045,40 1046,36 1050,36 C 1054,36 1055,39 1054,43 L 1050,53"
					stroke="#475569"
					stroke-width="2.5"
					stroke-linecap="round"
					fill="none"
				/>
				<!-- Rangka Segitiga Hanger Kayu/Kawat -->
				<path
					d="M 1050,53 L 1028,63 L 1072,63 Z"
					stroke="#475569"
					stroke-width="2"
					stroke-linejoin="round"
					fill="#F1F5F9"
				/>

				<!-- Kemeja Berkerah Rapi -->
				<path
					d="M 1032,61 L 1040,58 L 1050,62 L 1060,58 L 1068,61 L 1078,70 L 1071,78 L 1065,74 L 1065,102 L 1035,102 L 1035,74 L 1029,78 L 1022,70 Z"
					fill="#0284C7"
					class="drop-shadow-md"
				/>
				<!-- Kerah Kemeja -->
				<path d="M 1040,58 L 1050,66 L 1060,58" stroke="#E0F2FE" stroke-width="2" fill="none" stroke-linejoin="round" />
				<!-- Kancing Putih -->
				<circle cx="1050" cy="73" r="1.3" fill="#FFFFFF" />
				<circle cx="1050" cy="82" r="1.3" fill="#FFFFFF" />
				<circle cx="1050" cy="91" r="1.3" fill="#FFFFFF" />
			</g>

			<!-- Kilau Kesegaran Dekat Jemuran -->
			<path d="M 230,38 L 233,30 L 236,38 L 244,41 L 236,44 L 233,52 L 230,44 L 222,41 Z" fill="#F59E0B" opacity="0.8" class="animate-pulse" />
			<path d="M 960,34 L 962,28 L 965,34 L 971,36 L 965,38 L 962,44 L 960,38 L 954,36 Z" fill="#0EA5E9" opacity="0.8" class="animate-pulse" />
		</svg>
	</div>

	<!-- Gelembung Sabun Organik Transparan Melayang (Soap Bubbles) -->
	<div class="pointer-events-none absolute inset-0 overflow-hidden select-none">
		<div class="bubble left-[12%] top-[22%] h-20 w-20 animate-bubble-1"></div>
		<div class="bubble left-[6%] top-[42%] h-14 w-14 animate-bubble-2"></div>
		<div class="bubble bottom-[24%] left-[18%] h-16 w-16 animate-bubble-3"></div>
		<div class="bubble bottom-[38%] left-[8%] h-9 w-9 animate-bubble-4"></div>
		<div class="bubble bottom-[12%] left-[28%] h-12 w-12 animate-bubble-5"></div>

		<div class="bubble right-[14%] top-[24%] h-24 w-24 animate-bubble-6"></div>
		<div class="bubble right-[6%] top-[56%] h-14 w-14 animate-bubble-7"></div>
		<div class="bubble bottom-[18%] right-[20%] h-20 w-20 animate-bubble-8"></div>
		<div class="bubble right-[28%] top-[38%] h-10 w-10 animate-bubble-9"></div>
		<div class="bubble bottom-[32%] right-[10%] h-8 w-8 animate-bubble-10"></div>
		<div class="bubble bottom-[8%] right-[32%] h-12 w-12 animate-bubble-2"></div>
	</div>

	<!-- Login Card Container -->
	<div class="relative z-10 w-full max-w-md rounded-3xl border border-white/90 bg-white/95 p-8 shadow-2xl shadow-sky-950/10 backdrop-blur-md">
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
				class="flex w-full items-center justify-center gap-2 rounded-xl bg-[#0284C7] px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-[#0369A1] disabled:opacity-60"
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

<style>
	/* Animasi ayunan jemuran tertiup angin sejuk - pivot tepat di tali jemuran */
	@keyframes sway1 {
		0%, 100% {
			transform: rotate(-2deg);
		}
		50% {
			transform: rotate(2.5deg);
		}
	}
	@keyframes sway2 {
		0%, 100% {
			transform: rotate(2.2deg);
		}
		50% {
			transform: rotate(-1.8deg);
		}
	}
	@keyframes sway3 {
		0%, 100% {
			transform: rotate(-1.5deg);
		}
		50% {
			transform: rotate(2deg);
		}
	}
	@keyframes sway4 {
		0%, 100% {
			transform: rotate(2deg);
		}
		50% {
			transform: rotate(-2.2deg);
		}
	}

	.garment-1 {
		animation: sway1 5.5s ease-in-out infinite;
	}
	.garment-2 {
		animation: sway2 6.2s ease-in-out infinite 0.4s;
	}
	.garment-3 {
		animation: sway3 5.8s ease-in-out infinite 1s;
	}
	.garment-4 {
		animation: sway4 6.8s ease-in-out infinite 0.7s;
	}

	/* Animasi gelembung melayang alami */
	@keyframes floatBubble {
		0%, 100% {
			transform: translateY(0px) scale(1);
		}
		50% {
			transform: translateY(-26px) scale(1.06);
		}
	}

	@keyframes pulseGlow {
		0%, 100% {
			opacity: 0.3;
			transform: scale(1);
		}
		50% {
			opacity: 0.6;
			transform: scale(1.1);
		}
	}

	.animate-pulse-glow {
		animation: pulseGlow 8s ease-in-out infinite;
	}
	.animate-pulse-glow-delayed {
		animation: pulseGlow 10s ease-in-out infinite 2.5s;
	}

	/* Bubble animations with varied durations */
	.animate-bubble-1 { animation: floatBubble 5.2s ease-in-out infinite; }
	.animate-bubble-2 { animation: floatBubble 4.4s ease-in-out infinite 0.8s; }
	.animate-bubble-3 { animation: floatBubble 6.1s ease-in-out infinite 1.6s; }
	.animate-bubble-4 { animation: floatBubble 3.9s ease-in-out infinite 0.4s; }
	.animate-bubble-5 { animation: floatBubble 5.8s ease-in-out infinite 1.2s; }
	.animate-bubble-6 { animation: floatBubble 4.6s ease-in-out infinite 2.1s; }
	.animate-bubble-7 { animation: floatBubble 6.4s ease-in-out infinite 0.6s; }
	.animate-bubble-8 { animation: floatBubble 4.2s ease-in-out infinite 1.5s; }
	.animate-bubble-9 { animation: floatBubble 3.7s ease-in-out infinite 2.7s; }
	.animate-bubble-10 { animation: floatBubble 5.0s ease-in-out infinite 1.9s; }

	/* Soap bubble 3D realistic glass style */
	.bubble {
		position: absolute;
		border-radius: 9999px;
		background: radial-gradient(circle at 35% 30%, rgba(255, 255, 255, 0.85), rgba(186, 230, 253, 0.4) 55%, rgba(56, 189, 248, 0.25) 100%);
		border: 1.5px solid rgba(255, 255, 255, 0.75);
		box-shadow: inset -4px -4px 10px rgba(14, 165, 233, 0.25), 0 8px 24px rgba(2, 132, 199, 0.08);
		pointer-events: none;
	}
	.bubble::after {
		content: '';
		position: absolute;
		top: 15%;
		left: 20%;
		width: 25%;
		height: 25%;
		border-radius: 9999px;
		background: radial-gradient(circle, rgba(255, 255, 255, 0.95), rgba(255, 255, 255, 0));
	}
</style>
