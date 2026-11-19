<script>
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	export let currentPage = 1;
	export let totalItems = 0;
	export let pageSize = 10;

	$: totalPages = Math.max(1, Math.ceil(totalItems / pageSize));
	$: startIdx = totalItems === 0 ? 0 : (currentPage - 1) * pageSize + 1;
	$: endIdx = Math.min(currentPage * pageSize, totalItems);
</script>

<div class="flex flex-col items-center justify-between gap-3 border-t border-[#E2E8F0] px-5 py-3 text-xs text-[#64748B] sm:flex-row">
	<span>Menampilkan {startIdx}–{endIdx} dari {totalItems} data</span>
	<div class="flex items-center gap-1.5">
		<button
			type="button"
			disabled={currentPage <= 1}
			on:click={() => (currentPage = Math.max(1, currentPage - 1))}
			class="inline-flex items-center gap-1 rounded-lg border border-[#E2E8F0] bg-white px-2.5 py-1.5 text-xs font-medium text-[#64748B] hover:bg-slate-50 disabled:opacity-40"
		>
			<ChevronLeft class="h-3.5 w-3.5" />
			<span>Sebelumnya</span>
		</button>
		<span class="flex h-7 min-w-[28px] items-center justify-center rounded-lg border border-[#0284C7] bg-[#EFF6FF] px-2 text-xs font-bold text-[#1D4ED8]">
			{currentPage}
		</span>
		<button
			type="button"
			disabled={currentPage >= totalPages}
			on:click={() => (currentPage = Math.min(totalPages, currentPage + 1))}
			class="inline-flex items-center gap-1 rounded-lg border border-[#E2E8F0] bg-white px-2.5 py-1.5 text-xs font-medium text-[#64748B] hover:bg-slate-50 disabled:opacity-40"
		>
			<span>Berikutnya</span>
			<ChevronRight class="h-3.5 w-3.5" />
		</button>
	</div>
</div>
