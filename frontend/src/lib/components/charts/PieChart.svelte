<script lang="ts">
	import { cn } from "$lib/utils";

	interface DataPoint {
		x: string | number;
		y: number;
	}

	interface Props {
		data: DataPoint[];
		class?: string;
	}

	let { data = [], class: className }: Props = $props();

	let hoveredIndex = $state<number | null>(null);

	const colors = [
		"#3b82f6",
		"#0ea5e9",
		"#f59e0b",
		"#10b981",
		"#8b5cf6",
		"#ec4899",
		"#f97316"
	];

	const total = $derived(data.reduce((sum, d) => sum + d.y, 0));

	const slices = $derived.by(() => {
		let currentAngle = -Math.PI / 2;
		const cx = 100;
		const cy = 100;
		const radius = 80;

		return data.map((point, idx) => {
			const sliceAngle = (point.y / total) * 2 * Math.PI;
			const endAngle = currentAngle + sliceAngle;
			const percentage = (point.y / total) * 100;

			const x1 = cx + radius * Math.cos(currentAngle);
			const y1 = cy + radius * Math.sin(currentAngle);
			const x2 = cx + radius * Math.cos(endAngle);
			const y2 = cy + radius * Math.sin(endAngle);
			const largeArc = sliceAngle > Math.PI ? 1 : 0;

			const d = `M ${cx} ${cy} L ${x1} ${y1} A ${radius} ${radius} 0 ${largeArc} 1 ${x2} ${y2} Z`;

			const result = {
				label: String(point.x),
				value: point.y,
				percentage: Math.round(percentage * 10) / 10,
				color: colors[idx % colors.length],
				d
			};
			currentAngle = endAngle;
			return result;
		});
	});
</script>

<div class={cn("flex flex-col gap-5", className)}>
	<!-- Pie SVG -->
	<div class="flex justify-center">
		<svg
			viewBox="0 0 200 200"
			class="w-full max-w-52"
			role="img"
			aria-label="Pie chart"
		>
			{#each slices as slice, idx}
				<path
					d={slice.d}
					fill={slice.color}
					fill-opacity={hoveredIndex === idx ? "1" : "0.8"}
					stroke="white"
					stroke-width="1"
					class="transition-all cursor-pointer dark:stroke-slate-900"
					onmouseenter={() => hoveredIndex = idx}
					onmouseleave={() => hoveredIndex = null}
					role="button"
					tabindex="0"
					aria-label="{slice.label}: {slice.percentage}%"
				/>
			{/each}

			<!-- Center label on hover -->
			{#if hoveredIndex !== null}
				{@const slice = slices[hoveredIndex]}
				<text x="100" y="95" text-anchor="middle" class="text-[14px] font-black fill-slate-900 dark:fill-white" style="font-size: 14px; font-weight: 900;">
					{slice.percentage}%
				</text>
				<text x="100" y="112" text-anchor="middle" class="fill-slate-500" style="font-size: 9px;">
					{slice.label}
				</text>
			{/if}
		</svg>
	</div>

	<!-- Legend -->
	<div class="grid grid-cols-2 gap-x-6 gap-y-2.5 px-2">
		{#each slices as slice, idx}
			<div
				class="flex items-center gap-2.5 cursor-pointer group"
				onmouseenter={() => hoveredIndex = idx}
				onmouseleave={() => hoveredIndex = null}
			>
				<div
					class="w-2.5 h-2.5 rounded-full shrink-0 transition-transform group-hover:scale-125"
					style="background-color: {slice.color}"
				></div>
				<div class="min-w-0 flex-1">
					<p class="text-xs font-semibold text-slate-700 dark:text-slate-300 truncate group-hover:text-slate-900 dark:group-hover:text-white transition-colors">
						{slice.label}
					</p>
				</div>
				<span class="text-xs font-bold shrink-0" style="color: {slice.color}">
					{slice.percentage}%
				</span>
			</div>
		{/each}
	</div>
</div>
