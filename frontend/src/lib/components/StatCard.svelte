<script>
	export let label = '';
	export let value = '';
	export let sublabel = '';
</script>

<div class="rounded-2xl border border-[#E2E8F0] bg-white p-5 shadow-sm transition hover:border-slate-300">
	<div class="flex items-start justify-between gap-3">
		<div class="min-w-0">
			<p class="text-xs font-medium uppercase tracking-wider text-[#64748B]">{label}</p>
			<p class="mt-2 truncate text-2xl font-bold text-[#0F172A]">{value}</p>
			{#if sublabel}
				<p class="mt-1 text-xs text-[#64748B]">{sublabel}</p>
			{/if}
		</div>
		<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#EFF6FF] text-[#0284C7]">
			<slot name="icon" />
		</div>
	</div>
</div>
