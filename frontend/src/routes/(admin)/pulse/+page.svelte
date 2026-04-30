<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import AreaChart from "$lib/components/charts/AreaChart.svelte";
  import { cn } from "$lib/utils/styles";
  import { authStore } from "$lib/stores/auth.store";
  import { Cpu, Server, HardDrive, Activity, Terminal, Zap } from "lucide-svelte";
  import { slide } from "svelte/transition";

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

  const MAX_SAMPLES = 30;

  function connect() {
    const token = $authStore.accessToken;
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const host = window.location.hostname === "localhost" ? "localhost:8080" : window.location.host;

    ws = new WebSocket(`${protocol}//${host}/api/pulse/stream?token=${token}`);

    ws.onopen = () => {
      stats.status = "up";
    };

    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data);

      if (msg.type === "stats") {
        stats = { ...msg.data, status: "up" };
        const now = new Date();
        cpuHistory = [...cpuHistory, { date: now, value: stats.cpu_usage_percent }].slice(-MAX_SAMPLES);
        memHistory = [...memHistory, { date: now, value: stats.memory_usage_percent }].slice(-MAX_SAMPLES);
      } else if (msg.type === "log") {
        logs = [msg.data, ...logs].slice(0, 100);
      }
    };

    ws.onclose = () => {
      stats.status = "down";
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

<div class="space-y-8">
  <!-- Header -->
  <div class="flex flex-col gap-2">
    <div class="flex items-center gap-3">
      <h1 class="text-3xl font-bold text-slate-900 dark:text-white">System Pulse</h1>
      <span class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold {stats.status === 'up' ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400' : stats.status === 'connecting' ? 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-400' : 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'}">
        <span class="w-2 h-2 rounded-full {stats.status === 'up' ? 'bg-green-600 dark:bg-green-400 animate-pulse' : stats.status === 'connecting' ? 'bg-amber-600 dark:bg-amber-400' : 'bg-red-600 dark:bg-red-400'}"></span>
        {stats.status}
      </span>
    </div>
    <p class="text-sm text-slate-600 dark:text-slate-400">Real-time system monitoring and performance metrics</p>
  </div>

  <!-- Vitals Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    {#each vitals as vital}
      <div class="bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800 p-6 hover:border-slate-300 dark:hover:border-slate-700 transition-colors">
        <div class="flex items-start justify-between mb-4">
          <div class={cn("p-3 rounded-lg", vital.bg)}>
            <vital.icon class={cn("w-5 h-5", vital.color)} />
          </div>
        </div>
        <p class="text-xs text-slate-500 dark:text-slate-400 font-medium mb-1">{vital.name}</p>
        <p class="text-3xl font-bold text-slate-900 dark:text-white font-mono">{vital.value}</p>
      </div>
    {/each}
  </div>

  <!-- Charts Row -->
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <!-- CPU Chart -->
    <div class="bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="p-6 border-b border-slate-200 dark:border-slate-800">
        <h2 class="text-lg font-bold text-slate-900 dark:text-white">CPU Usage</h2>
        <p class="text-xs text-slate-600 dark:text-slate-400 mt-1">Last 30 samples</p>
      </div>
      <div class="p-6">
        <AreaChart data={cpuHistory} height={250} />
      </div>
    </div>

    <!-- Memory Chart -->
    <div class="bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="p-6 border-b border-slate-200 dark:border-slate-800">
        <h2 class="text-lg font-bold text-slate-900 dark:text-white">Memory Usage</h2>
        <p class="text-xs text-slate-600 dark:text-slate-400 mt-1">Last 30 samples</p>
      </div>
      <div class="p-6">
        <AreaChart data={memHistory} height={250} />
      </div>
    </div>
  </div>

  <!-- Logs Section -->
  <div class="bg-slate-950 dark:bg-slate-900 rounded-xl border border-slate-800 overflow-hidden">
    <div class="bg-slate-900 px-6 py-4 border-b border-slate-800 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <Terminal class="w-5 h-5 text-blue-400" />
        <div>
          <h2 class="text-sm font-bold text-white">System Logs</h2>
          <p class="text-xs text-slate-400">Real-time event stream</p>
        </div>
      </div>
    </div>

    <div class="h-100 overflow-y-auto font-mono text-sm space-y-1 p-6 custom-scrollbar">
      {#if logs.length === 0}
        <div class="h-full flex items-center justify-center text-slate-600">
          <div class="text-center">
            <div class="w-10 h-10 rounded-full border-2 border-slate-700 border-t-blue-500 animate-spin mx-auto mb-3"></div>
            <p class="text-xs font-medium">Connecting to logs...</p>
          </div>
        </div>
      {:else}
        {#each logs as log (log.message + log.timestamp)}
          <div in:slide={{ duration: 300 }} class="flex gap-3 py-1 px-2 rounded hover:bg-slate-800/50 transition-colors group">
            <span class="text-slate-500 shrink-0 tabular-nums text-xs opacity-70">{new Date().toLocaleTimeString('en-GB', { hour12: false })}</span>
            <span class={cn("shrink-0 text-xs font-bold px-2 py-0.5 rounded", log.level === 'error' ? 'text-red-400 bg-red-500/20' : 'text-amber-400 bg-amber-500/20')}>
              {log.level?.toUpperCase()}
            </span>
            <span class="text-slate-300 group-hover:text-slate-100 text-xs">{log.message}</span>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>

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
