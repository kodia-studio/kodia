<script lang="ts">
	import { scaleBand, scaleLinear } from "d3-scale";
	import { cn } from "$lib/utils";

	interface DataPoint {
		x: string | number;
		y: number;
	}

	interface Props {
		data: DataPoint[];
		title?: string;
		height?: number;
		class?: string;
	}

	let { data = [], title, height = 300, class: className }: Props = $props();

	const padding = { left: 50, right: 20, top: 20, bottom: 40 };
	let hoveredIndex = $state<number | null>(null);

	const xScale = $derived.by(() => {
		return scaleBand()
			.domain(data.map(d => String(d.x)))
			.range([padding.left, 100 - padding.right])
			.padding(0.4);
	});

	const yScale = $derived.by(() => {
		const max = data.length > 0 ? Math.max(...data.map(d => d.y)) : 100;
		return scaleLinear()
			.domain([0, max])
			.range([height - padding.bottom, padding.top]);
	});

	const yTicks = $derived.by(() => {
		if (data.length === 0) return [];
		const max = Math.max(...data.map(d => d.y));
		const step = Math.pow(10, Math.floor(Math.log10(max))) / 2;
		const ticks = [];
		for (let i = 0; i <= max; i += step) {
			ticks.push(i);
		}
		return ticks.slice(0, 5);
	});

	function handleBarHover(index: number) {
		hoveredIndex = index;
	}

	function handleBarLeave() {
		hoveredIndex = null;
	}
</script>

<div class={cn("space-y-4", className)}>
	{#if title}
		<h3 class="text-xs font-black uppercase tracking-widest text-slate-500 ml-1">
			{title}
		</h3>
	{/if}

	<div class="card-premium p-6 relative" style="height: {height}px">
		<svg
			{height}
			viewBox="0 0 100 {height}"
			preserveAspectRatio="none"
			class="w-full absolute inset-0"
			role="img"
			aria-label="Bar chart"
		>
			<!-- Y Axis -->
			<line
				x1={padding.left}
				y1={padding.top}
				x2={padding.left}
				y2={height - padding.bottom}
				stroke="currentColor"
				stroke-width="0.5"
				opacity="0.2"
			/>

			<!-- Y Axis Ticks & Labels -->
			{#each yTicks as tick}
				{@const y = yScale(tick)}
				<line
					x1={padding.left - 2}
					y1={y}
					x2={padding.left}
					y2={y}
					stroke="currentColor"
					stroke-width="0.5"
					opacity="0.2"
				/>
				<text
					x={padding.left - 5}
					y={y}
					text-anchor="end"
					dominant-baseline="middle"
					class="text-[8px] fill-slate-500"
				>
					{tick.toLocaleString()}
				</text>
			{/each}

			<!-- X Axis -->
			<line
				x1={padding.left}
				y1={height - padding.bottom}
				x2={100 - padding.right}
				y2={height - padding.bottom}
				stroke="currentColor"
				stroke-width="0.5"
				opacity="0.2"
			/>

			<!-- Grid Lines -->
			{#each yTicks as tick}
				{@const y = yScale(tick)}
				<line
					x1={padding.left}
					y1={y}
					x2={100 - padding.right}
					y2={y}
					stroke="currentColor"
					stroke-width="0.25"
					opacity="0.1"
				/>
			{/each}

			<!-- Bars -->
			{#each data as point, idx}
				{@const barWidth = xScale.bandwidth()}
				{@const barX = xScale(String(point.x))}
				{@const barY = yScale(point.y)}
				{@const barHeight = height - padding.bottom - barY}

				<rect
					x={barX}
					y={barY}
					width={barWidth}
					height={barHeight}
					fill={hoveredIndex === idx ? "#8b5cf6" : "rgba(139, 92, 246, 0.8)"}
					class="transition-all cursor-pointer dark:opacity-60"
					onmouseenter={() => handleBarHover(idx)}
					onmouseleave={handleBarLeave}
					role="button"
					tabindex="0"
					aria-label="Bar {idx}: {point.x} - {point.y}"
				/>
			{/each}
		</svg>

		<!-- X Axis Labels (Positioned Absolutely) -->
		<div class="absolute" style="left: {padding.left}%; bottom: 0; right: {padding.right}%; height: {padding.bottom}px;">
			<div class="flex h-full items-start justify-around text-center">
				{#each data as point}
					<span class="text-[10px] font-bold uppercase tracking-tight text-slate-400 leading-tight">
						{point.x}
					</span>
				{/each}
			</div>
		</div>

		<!-- Tooltip -->
		{#if hoveredIndex !== null}
			{@const point = data[hoveredIndex]}
			{@const barX = xScale(String(point.x)) || 0}
			{@const barWidth = xScale.bandwidth()}
			<div
				class="absolute pointer-events-none bg-white dark:bg-slate-900 p-2 px-3 rounded-lg border border-primary/20 shadow-xl z-10"
				style="left: {barX + barWidth / 2}%; top: {yScale(point.y)}%; transform: translate(-50%, -100%); margin-top: -8px;"
			>
				<p class="text-[10px] font-black uppercase tracking-widest text-primary">{point.x}</p>
				<p class="text-lg font-black tracking-tight text-slate-900 dark:text-white">{point.y.toLocaleString()}</p>
			</div>
		{/if}
	</div>
</div>
