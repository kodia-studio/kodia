<script lang="ts">
  import { scaleLinear, scaleTime } from "d3-scale";
  import { area, line } from "d3-shape";
  import { cn } from "$lib/utils";

  interface DataPoint {
    date: Date;
    value: number;
  }

  interface Props {
    data: DataPoint[];
    height?: number;
    color?: string;
    class?: string;
  }

  let {
    data = [],
    height = 300,
    color = "#8b5cf6",
    class: className
  }: Props = $props();

  const padding = { top: 20, right: 20, bottom: 40, left: 50 };

  let hoveredIndex = $state<number | null>(null);
  let tooltipPos = $state({ x: 0, y: 0 });

  const xScale = $derived.by(() => {
    if (data.length === 0) return scaleTime();
    const dates = data.map(d => d.date);
    return scaleTime()
      .domain([Math.min(...dates.map(d => d.getTime())), Math.max(...dates.map(d => d.getTime()))])
      .range([padding.left, 100 - padding.right]);
  });

  const yScale = $derived.by(() => {
    if (data.length === 0) return scaleLinear();
    const values = data.map(d => d.value);
    const max = Math.max(...values);
    return scaleLinear()
      .domain([0, max])
      .range([height - padding.bottom, padding.top]);
  });

  const pathData = $derived.by(() => {
    if (data.length === 0) return "";
    const areaGenerator = area<DataPoint>()
      .x(d => (xScale(d.date) as number) * (100 / (100 - padding.left - padding.right)) - (padding.left * 100 / (100 - padding.left - padding.right)))
      .y0(height - padding.bottom)
      .y1(d => yScale(d.value) as number);
    return areaGenerator(data) || "";
  });

  const lineData = $derived.by(() => {
    if (data.length === 0) return "";
    const lineGenerator = line<DataPoint>()
      .x(d => (xScale(d.date) as number) * (100 / (100 - padding.left - padding.right)) - (padding.left * 100 / (100 - padding.left - padding.right)))
      .y(d => yScale(d.value) as number);
    return lineGenerator(data) || "";
  });

  const yTicks = $derived.by(() => {
    if (data.length === 0) return [];
    const max = Math.max(...data.map(d => d.value));
    const step = Math.pow(10, Math.floor(Math.log10(max))) / 2;
    const ticks = [];
    for (let i = 0; i <= max; i += step) {
      ticks.push(i);
    }
    return ticks.slice(0, 5);
  });

  function handleMouseMove(e: MouseEvent) {
    const svg = e.currentTarget as SVGElement;
    const rect = svg.getBoundingClientRect();
    const x = ((e.clientX - rect.left) / rect.width) * 100;

    const xPixels = (x - padding.left) / ((100 - padding.left - padding.right) / 100);
    const closestIndex = data.reduce((closest, _, idx) => {
      const distance = Math.abs((xScale(data[idx].date) as number) - xPixels);
      const closestDistance = Math.abs((xScale(data[closest].date) as number) - xPixels);
      return distance < closestDistance ? idx : closest;
    }, 0);

    hoveredIndex = closestIndex;
    tooltipPos = { x: (xScale(data[closestIndex].date) as number), y: yScale(data[closestIndex].value) as number };
  }

  function handleMouseLeave() {
    hoveredIndex = null;
  }
</script>

<div style="height: {height}px" class={cn("w-full relative", className)}>
  <svg
    {height}
    viewBox="0 0 100 {height}"
    preserveAspectRatio="none"
    class="w-full absolute inset-0"
    role="img"
    aria-label="Area chart"
    onmousemove={handleMouseMove}
    onmouseleave={handleMouseLeave}
  >
    <!-- Area -->
    <path
      d={pathData}
      fill={color}
      fill-opacity="0.2"
      class="transition-opacity"
    />

    <!-- Line -->
    <path
      d={lineData}
      fill="none"
      stroke={color}
      stroke-width="2"
      class="transition-all"
      vector-effect="non-scaling-stroke"
    />

    <!-- Y Axis -->
    <line x1={padding.left} y1={padding.top} x2={padding.left} y2={height - padding.bottom} stroke="currentColor" stroke-width="0.5" opacity="0.2" />

    <!-- Y Axis Ticks -->
    {#each yTicks as tick}
      {@const y = yScale(tick)}
      <line x1={padding.left - 2} y1={y} x2={padding.left} y2={y} stroke="currentColor" stroke-width="0.5" opacity="0.2" />
      <text x={padding.left - 5} y={y} text-anchor="end" dominant-baseline="middle" class="text-[8px] fill-slate-500">
        {tick.toLocaleString()}
      </text>
    {/each}

    <!-- X Axis -->
    <line x1={padding.left} y1={height - padding.bottom} x2={100 - padding.right} y2={height - padding.bottom} stroke="currentColor" stroke-width="0.5" opacity="0.2" />

    <!-- Hover Indicator -->
    {#if hoveredIndex !== null}
      <circle
        cx={tooltipPos.x}
        cy={tooltipPos.y}
        r="2"
        fill={color}
        class="pointer-events-none"
      />
    {/if}
  </svg>

  <!-- Tooltip -->
  {#if hoveredIndex !== null}
    <div
      class="absolute pointer-events-none bg-white dark:bg-slate-900 p-2 rounded border border-slate-200 dark:border-slate-800 text-sm shadow-lg z-10"
      style="left: {(tooltipPos.x / 100) * 100}%; top: {(tooltipPos.y / height) * 100}%; transform: translate(-50%, -100%); margin-top: -8px;"
    >
      <div class="text-[10px] uppercase tracking-wider text-slate-500 font-bold mb-1">
        {data[hoveredIndex].date.toLocaleDateString()}
      </div>
      <div class="text-sm font-bold text-slate-900 dark:text-white">
        {data[hoveredIndex].value.toLocaleString()}
      </div>
    </div>
  {/if}
</div>
