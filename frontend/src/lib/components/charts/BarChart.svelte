<script lang="ts">
	import { Chart, Svg, Axis, Bars, Tooltip } from "layerchart";
	import { scaleBand } from "d3-scale";
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

	let { data, title, height = 300, class: className }: Props = $props();
</script>

<div class={cn("space-y-4", className)}>
	{#if title}
		<h3 class="text-xs font-black uppercase tracking-widest text-slate-500 ml-1">
			{title}
		</h3>
	{/if}

	<div class="card-premium p-6" style="height: {height}px">
		<Chart
			{data}
			x="x"
			xScale={scaleBand().padding(0.4)}
			y="y"
			yDomain={[0, null]}
			yNice
			padding={{ left: 16, bottom: 24, top: 8, right: 8 }}
		>
			<Svg>
				<Axis placement="bottom" grid={{ class: 'stroke-slate-100 dark:stroke-slate-800' }} rule={false} tickLabelProps={{ class: 'fill-slate-400 text-[10px] font-bold uppercase tracking-tight' }} />
				<Axis placement="left" grid={{ class: 'stroke-slate-100 dark:stroke-slate-800' }} rule={false} tickLabelProps={{ class: 'fill-slate-400 text-[10px] font-bold uppercase tracking-tight' }} />
				<Bars 
					class="fill-primary/80 dark:fill-primary/60 hover:fill-primary transition-colors cursor-pointer"
					radius={4}
				/>
			</Svg>
			<Tooltip.Root let:data>
				<div class="glass p-2 px-3 rounded-lg shadow-xl border border-primary/20">
					<p class="text-[10px] font-black uppercase tracking-widest text-primary">{data.x}</p>
					<p class="text-lg font-black tracking-tight text-slate-900 dark:text-white">{data.y.toLocaleString()}</p>
				</div>
			</Tooltip.Root>
		</Chart>
	</div>
</div>
