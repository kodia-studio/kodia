<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import AdminLayout from "$lib/components/layouts/AdminLayout.svelte";
  import Grid from "$lib/components/shared/Grid.svelte";
  import ChartCard from "$lib/components/charts/ChartCard.svelte";
  import AreaChart from "$lib/components/charts/AreaChart.svelte";
  import { cn } from "$lib/utils/styles";
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
  <div class="space-y-12 pb-10">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 pb-6 border-b border-slate-200/50 dark:border-white/5">
      <div class="flex flex-col">
        <div class="flex items-center gap-4">
          <h1 class="text-4xl font-black font-heading tracking-tight text-slate-900 dark:text-white leading-none">System Pulse.</h1>
          <div class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest flex items-center gap-2 
            {stats.status === 'up' ? 'bg-emerald-500/10 text-emerald-500 border border-emerald-500/20 shadow-[0_0_15px_rgba(16,185,129,0.1)]' : 
             stats.status === 'connecting' ? 'bg-amber-500/10 text-amber-500 border border-amber-500/20' : 
             'bg-rose-500/10 text-rose-500 border border-rose-500/20'}">
            <span class="w-1.5 h-1.5 rounded-full {stats.status === 'up' ? 'bg-emerald-500 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.8)]' : stats.status === 'connecting' ? 'bg-amber-500' : 'bg-rose-500'}"></span>
            {stats.status}
          </div>
        </div>
        <p class="text-xs font-black uppercase tracking-[0.3em] text-slate-400 mt-3 leading-none">High-Frequency Telemetry Stream</p>
      </div>
      
      <div class="flex items-center gap-3 px-4 py-2 bg-slate-100/50 dark:bg-white/5 rounded-xl border border-slate-200/50 dark:border-white/5 text-[10px] font-black uppercase tracking-widest text-slate-400">
        <Activity class="w-3.5 h-3.5" />
        Endpoint: /pulse/stream
      </div>
    </div>

    <!-- Elite Vitals Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
      {#each vitals as vital}
        <div class="group relative">
          <div class="absolute -inset-1 bg-linear-to-r from-primary to-secondary rounded-3xl blur opacity-0 group-hover:opacity-10 transition duration-500"></div>
          <div class="glass relative p-7 rounded-3xl border border-slate-200/50 dark:border-white/5 flex items-start justify-between hover:scale-[1.02] transition-all duration-300">
            <div>
              <p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mb-2 leading-none">{vital.name}</p>
              <h2 class="text-3xl font-black text-slate-900 dark:text-white tracking-tighter tabular-nums">{vital.value}</h2>
            </div>
            <div class={cn("p-4 rounded-2xl ring-1 ring-white/20 shadow-inner group-hover:scale-110 transition-transform", vital.bg, vital.color)}>
              <vital.icon class="w-6 h-6" />
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Chart Analytics Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- CPU Load Chart -->
      <div class="glass rounded-3xl border border-slate-200/50 dark:border-white/5 overflow-hidden flex flex-col bg-white/30 dark:bg-slate-900/40">
        <div class="px-8 py-7 border-b border-slate-200/50 dark:border-white/5 flex items-center justify-between">
          <div>
            <h3 class="text-lg font-black text-slate-900 dark:text-white leading-none uppercase tracking-tight">CPU Core Load.</h3>
            <p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mt-2">Active processing intensity</p>
          </div>
          <div class="w-2.5 h-2.5 rounded-full bg-primary/20 animate-pulse ring-4 ring-primary/5"></div>
        </div>
        <div class="p-8 pb-4">
          <AreaChart data={cpuHistory} height={250} />
        </div>
      </div>

      <!-- Memory Usage Chart -->
      <div class="glass rounded-3xl border border-slate-200/50 dark:border-white/5 overflow-hidden flex flex-col bg-white/30 dark:bg-slate-900/40">
        <div class="px-8 py-7 border-b border-slate-200/50 dark:border-white/5 flex items-center justify-between">
          <div>
            <h3 class="text-lg font-black text-slate-900 dark:text-white leading-none uppercase tracking-tight">Heap Statistics.</h3>
            <p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mt-2">Memory allocation & GC cycle</p>
          </div>
          <div class="w-2.5 h-2.5 rounded-full bg-secondary/20 animate-pulse ring-4 ring-secondary/5"></div>
        </div>
        <div class="p-8 pb-4">
          <AreaChart data={memHistory} height={250} />
        </div>
      </div>
    </div>

    <!-- High-Fidelity Terminal Logs -->
    <div class="relative group">
      <div class="absolute -inset-1 bg-linear-to-r from-primary to-secondary rounded-[32px] blur opacity-10 group-hover:opacity-20 transition duration-1000"></div>
      
      <div class="relative rounded-[32px] overflow-hidden border border-slate-800 shadow-2xl bg-[#0b1120]">
        <!-- Terminal Header -->
        <div class="bg-slate-900/80 px-7 py-4 flex items-center justify-between border-b border-white/5">
          <div class="flex items-center gap-6">
             <div class="flex gap-1.5">
                <div class="w-3 h-3 rounded-full bg-rose-500/50"></div>
                <div class="w-3 h-3 rounded-full bg-amber-500/50"></div>
                <div class="w-3 h-3 rounded-full bg-emerald-500/50"></div>
             </div>
             <div class="flex items-center gap-3">
               <Terminal size={14} class="text-primary" />
               <span class="text-[11px] font-black uppercase tracking-[0.2em] text-slate-500">Log Protocol Alpha</span>
             </div>
          </div>
          <div class="flex gap-4">
             <div class="flex items-center gap-2">
                <div class="w-2 h-2 rounded-full bg-amber-500"></div>
                <span class="text-[10px] font-black uppercase text-amber-500/80 tracking-widest">Warning</span>
             </div>
             <div class="flex items-center gap-2">
                <div class="w-2 h-2 rounded-full bg-rose-500"></div>
                <span class="text-[10px] font-black uppercase text-rose-500/80 tracking-widest">Critical</span>
             </div>
          </div>
        </div>
        
        <div class="h-[400px] overflow-y-auto p-8 font-mono text-[13px] space-y-3 bg-[#0b1120] custom-scrollbar selection:bg-primary/30">
          {#if logs.length === 0}
            <div class="h-full flex flex-col items-center justify-center text-slate-600 gap-4">
              <div class="w-12 h-12 rounded-full border-2 border-slate-800 border-t-primary animate-spin"></div>
              <p class="font-black uppercase tracking-widest text-[10px] animate-pulse">Syncing Event Logs...</p>
            </div>
          {:else}
            {#each logs as log (log.message + log.timestamp)}
              <div in:slide={{ duration: 300 }} class="flex gap-4 py-2 px-4 rounded-xl border border-transparent hover:border-white/5 hover:bg-white/5 transition-all group">
                <span class="text-slate-500 shrink-0 font-bold select-none tabular-nums opacity-60">
                  {new Date().toLocaleTimeString('en-GB', { hour12: false })}
                </span>
                <span class={cn(
                  "font-black shrink-0 uppercase w-16 text-center rounded-md px-1.5 py-0.5 text-[10px] tracking-widest",
                  log.level === 'error' ? 'bg-rose-500/20 text-rose-500 border border-rose-500/30' : 'bg-amber-500/20 text-amber-500 border border-amber-500/30'
                )}>
                  {log.level}
                </span>
                <span class="text-white/80 group-hover:text-white leading-relaxed">{log.message}</span>
              </div>
            {/each}
          {/if}
        </div>
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
