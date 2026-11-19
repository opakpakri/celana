<script>
	import { X } from 'lucide-svelte';
	import { createEventDispatcher } from 'svelte';

	export let open = false;
	export let title = '';
	export let maxWidth = 'max-w-lg';

	const dispatch = createEventDispatcher();

	function close() {
		dispatch('close');
	}
</script>

{#if open}
	<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
		<button
			type="button"
			aria-label="Tutup modal"
			class="fixed inset-0 bg-slate-900/40 backdrop-blur-[1px]"
			on:click={close}
		></button>
		<div class="relative z-10 w-full {maxWidth} max-h-[90vh] overflow-y-auto rounded-2xl border border-[#E2E8F0] bg-white shadow-xl">
			<div class="sticky top-0 z-10 flex items-center justify-between border-b border-[#E2E8F0] bg-white px-6 py-4">
				<h3 class="text-lg font-bold text-[#0F172A]">{title}</h3>
				<button
					type="button"
					on:click={close}
					class="rounded-lg p-1.5 text-[#64748B] transition hover:bg-slate-100 hover:text-[#0F172A]"
				>
					<X class="h-5 w-5" />
				</button>
			</div>
			<div class="p-6">
				<slot />
			</div>
		</div>
	</div>
{/if}
