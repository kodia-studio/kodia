<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import AdminLayout from "$lib/components/layouts/AdminLayout.svelte";
  import Grid from "$lib/components/shared/Grid.svelte";
  import ChartCard from "$lib/components/charts/ChartCard.svelte";
  import AreaChart from "$lib/components/charts/AreaChart.svelte";
  import { authStore } from "$lib/stores/auth.store";
  import { Activity, Server, Database, AlertCircle, Terminal, Cpu, HardDrive, Network } from "lucide-svelte";
  import { fade, slide } from "svelte/transition";

  // State
  let stats = $state({
    cpu_usage_percent: 0,
    memory_usage_percent: 0,
    disk_usage_percent: 0,
    goroutines: 0,
    status: "connecting"
  });

  let cpuHistory = $state<any[]>([]);
  let memHistory = $state<any[]>([]);
  let logs = $state<any[]>([]);
  let ws: WebSocket;

  // Max samples for charts
  const MAX_SAMPLES = 30;

  function connect() {
    const token = $authStore.accessToken;
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    // Adjust port/host as needed for your environment logic
    const host = window.location.hostname === "localhost" ? "localhost:8080" : window.location.host;
    
    ws = new WebSocket(`${protocol}//${host}/api/pulse/stream?token=${token}`);

    ws.onopen = () => {
      stats.status = "up";
    };

    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      
      if (msg.type === "stats") {
        stats = { ...msg.data, status: "up" };
        
        // Update History
        const now = new Date();
        
        cpuHistory = [...cpuHistory, { date: now, value: stats.cpu_usage_percent }].slice(-MAX_SAMPLES);
        memHistory = [...memHistory, { date: now, value: stats.memory_usage_percent }].slice(-MAX_SAMPLES);
      } else if (msg.type === "log") {
        logs = [msg.data, ...logs].slice(0, 100);
      }
    };

    ws.onclose = () => {
      stats.status = "down";
      // Auto reconnect after 3 seconds
      setTimeout(connect, 3000);
    };
  }

  onMount(() => {
    connect();
  });

  onDestroy(() => {
    if (ws) ws.close();
  });

  const vitals = $derived([
    { name: "CPU Usage", value: `${stats.cpu_usage_percent.toFixed(1)}%`, icon: Cpu, color: "text-blue-500", bg: "bg-blue-500/10" },
    { name: "RAM Used", value: `${stats.memory_usage_percent.toFixed(1)}%`, icon: Server, color: "text-purple-500", bg: "bg-purple-500/10" },
    { name: "Disk Space", value: `${stats.disk_usage_percent.toFixed(1)}%`, icon: HardDrive, color: "text-emerald-500", bg: "bg-emerald-500/10" },
    { name: "Goroutines", value: stats.goroutines.toString(), icon: Activity, color: "text-amber-500", bg: "bg-amber-500/10" },
  ]);
</script>

<AdminLayout>
  <div class="space-y-8">
    <div class="flex items-end justify-between">
      <div>
        <div class="flex items-center gap-3">
          <h1 class="text-3xl font-bold font-heading tracking-tight">Kodia Pulse</h1>
          <div class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider flex items-center gap-1.5 
            {stats.status === 'up' ? 'bg-emerald-500/10 text-emerald-500 border border-emerald-500/20' : 
             stats.status === 'connecting' ? 'bg-amber-500/10 text-amber-500 border border-amber-500/20' : 
             'bg-rose-500/10 text-rose-500 border border-rose-500/20'}">
            <span class="w-1.5 h-1.5 rounded-full {stats.status === 'up' ? 'bg-emerald-500 animate-pulse' : stats.status === 'connecting' ? 'bg-amber-500' : 'bg-rose-500'}"></span>
            {stats.status}
          </div>
        </div>
        <p class="text-muted-foreground mt-1.5">Real-time system vitals and professional telemetry.</p>
      </div>
      
      <div class="flex items-center gap-2 text-xs font-medium text-muted-foreground bg-muted/30 px-3 py-1.5 rounded-lg border border-border/50">
        <Network class="w-3.5 h-3.5" />
        Stream: /api/pulse/stream
      </div>
    </div>

    <!-- Vitals Grid -->
    <Grid cols={1} md={2} lg={4} gap="md">
      {#each vitals as vital}
        <div class="card p-6 flex items-start justify-between group hover:border-primary/50 transition-colors">
          <div>
            <p class="text-sm font-medium text-muted-foreground">{vital.name}</p>
            <h2 class="text-2xl font-bold mt-1 tabular-nums">{vital.value}</h2>
          </div>
          <div class="p-3 rounded-xl {vital.bg} {vital.color} group-hover:scale-110 transition-transform">
            <vital.icon class="w-5 h-5" />
          </div>
        </div>
      {/each}
    </Grid>

    <Grid cols={1} lg={2} gap="lg">
      <!-- CPU Chart -->
      <ChartCard title="CPU Load" subtitle="Real-time processor performance">
        <AreaChart data={cpuHistory} height={250} />
      </ChartCard>

      <!-- Memory Chart -->
      <ChartCard title="Memory Usage" subtitle="JVM/Heap utilization stats">
        <AreaChart data={memHistory} height={250} />
      </ChartCard>
    </Grid>

    <!-- Real-time Logs -->
    <div class="card overflow-hidden border-border/50 bg-black/5 dark:bg-black/20">
      <div class="px-6 py-4 border-b border-border/50 flex items-center justify-between">
        <div class="flex items-center gap-2 uppercase tracking-widest text-[11px] font-bold text-muted-foreground">
          <Terminal class="w-4 h-4 text-primary" />
          Live Console Feed (Warn/Error)
        </div>
        <div class="flex gap-2">
           <div class="px-2 py-1 rounded bg-amber-500/10 text-amber-500 text-[10px] font-bold border border-amber-500/20">WARN</div>
           <div class="px-2 py-1 rounded bg-rose-500/10 text-rose-500 text-[10px] font-bold border border-rose-500/20">ERROR</div>
        </div>
      </div>
      
      <div class="h-[300px] overflow-y-auto p-4 font-mono text-xs space-y-2 bg-[#09090b]">
        {#if logs.length === 0}
          <div class="h-full flex items-center justify-center text-muted-foreground italic opacity-50">
            Scanning for system alerts...
          </div>
        {:else}
          {#each logs as log (log.message + log.timestamp)}
            <div in:slide={{ duration: 200 }} class="flex gap-3 py-1 border-b border-white/5 last:border-0 hover:bg-white/5 px-2 rounded transition-colors group">
              <span class="text-white/30 shrink-0 select-none">[{new Date().toLocaleTimeString()}]</span>
              <span class={log.level === 'error' ? 'text-rose-500' : 'text-amber-500' + ' font-bold shrink-0 uppercase w-12'}>
                {log.level}
              </span>
              <span class="text-white/80 group-hover:text-white break-all">{log.message}</span>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</AdminLayout>

<style>
  /* Custom scrollbar for terminal */
  div::-webkit-scrollbar {
    width: 6px;
  }
  div::-webkit-scrollbar-track {
    background: transparent;
  }
  div::-webkit-scrollbar-thumb {
    background: #27272a;
    border-radius: 10px;
  }
</style>
