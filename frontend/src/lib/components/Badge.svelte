<script>
	export let variant = 'neutral'; // success | info | danger | warning | neutral

	const styles = {
		success: 'bg-[#DCFCE7] text-[#15803D]',
		info: 'bg-[#E0F2FE] text-[#0369A1]',
		danger: 'bg-[#FEE2E2] text-[#B91C1C]',
		warning: 'bg-amber-100 text-amber-800',
		neutral: 'bg-slate-100 text-slate-600'
	};
</script>

<span class="inline-flex items-center whitespace-nowrap rounded-full px-2.5 py-0.5 text-xs font-semibold {styles[variant] || styles.neutral}">
	<slot />
</span>
