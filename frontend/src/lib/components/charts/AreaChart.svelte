<script lang="ts">
  import { Chart, Svg, Axis, Area, Tooltip, LinearGradient } from "layerchart";
  import { scaleTime } from "d3-scale";
  import { curveMonotoneX } from "d3-shape";

  interface DataPoint {
    date: Date;
    value: number;
  }

  interface Props {
    data: DataPoint[];
    height?: number;
    color?: string;
  }

  let { 
    data, 
    height = 300,
    color = "hsl(var(--primary))" 
  }: Props = $props();
</script>

<div style:height="{height}px" class="w-full">
  <Chart
    {data}
    x="date"
    xScale={scaleTime()}
    y="value"
    yDomain={[0, null]}
    padding={{ top: 10, right: 10, bottom: 30, left: 40 }}
    tooltip={{ mode: "bisect-x" }}
  >
    <Svg>
      <LinearGradient {color} vertical let:gradient>
        <Area
          line={{ class: "stroke-2", stroke: color }}
          fill={gradient}
          curve={curveMonotoneX}
          class="fill-primary/20"
        />
      </LinearGradient>
      <Axis placement="left" grid={{ class: "stroke-muted/30" }} ticks={5} class="text-xs text-muted-foreground" />
      <Axis placement="bottom" grid={{ class: "stroke-muted/30" }} class="text-xs text-muted-foreground" />
    </Svg>
    <Tooltip.Root let:data>
      <div class="glass p-3 rounded-xl border-white/20 shadow-2xl">
        <div class="text-[10px] uppercase tracking-wider text-muted-foreground font-bold mb-1">
          {data.date.toLocaleDateString()}
        </div>
        <div class="text-sm font-bold text-foreground">
          {data.value.toLocaleString()}
        </div>
      </div>
    </Tooltip.Root>
  </Chart>
</div>
