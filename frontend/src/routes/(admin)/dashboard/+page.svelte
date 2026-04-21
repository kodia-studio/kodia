<script lang="ts">
  import AdminLayout from "$lib/components/layouts/AdminLayout.svelte";
  import Container from "$lib/components/shared/Container.svelte";
  import Grid from "$lib/components/shared/Grid.svelte";
  import ChartCard from "$lib/components/charts/ChartCard.svelte";
  import AreaChart from "$lib/components/charts/AreaChart.svelte";
  import { cn } from "$lib/utils/styles";
  import { Users, Shield, Zap, Database, TrendingUp, ArrowUpRight, Activity } from "lucide-svelte";

  const stats = [
    { name: "Total Users", value: "12,842", change: "+12%", icon: Users, color: "text-blue-500", bg: "bg-blue-500/10" },
    { name: "Security Score", value: "98.2%", change: "+0.5%", icon: Shield, color: "text-green-500", bg: "bg-green-500/10" },
    { name: "API Latency", value: "24ms", change: "-4ms", icon: Zap, color: "text-amber-500", bg: "bg-amber-500/10" },
    { name: "Storage Used", value: "1.2TB", change: "+8%", icon: Database, color: "text-purple-500", bg: "bg-purple-500/10" },
  ];

  // Mock data for the chart
  const chartData = Array.from({ length: 30 }, (_, i) => ({
    date: new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000),
    value: Math.floor(Math.random() * 5000) + 2000
  }));
</script>

<AdminLayout>
  <div class="space-y-12 pb-10">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 pb-6 border-b border-slate-200/50 dark:border-white/5">
      <div>
        <h1 class="text-4xl font-black font-heading tracking-tight text-slate-900 dark:text-white leading-none">Console Overview.</h1>
        <p class="text-xs font-black uppercase tracking-[0.3em] text-slate-400 mt-3">Autonomous Framework Monitoring</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="flex -space-x-3">
          {#each Array(3) as _, i}
            <div class="w-10 h-10 rounded-full border-4 border-white dark:border-slate-900 bg-slate-100 dark:bg-slate-800 flex items-center justify-center overflow-hidden ring-1 ring-slate-200 dark:ring-white/10">
              <img src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${i + 10}`} alt="Admin" class="w-full h-full object-cover" />
            </div>
          {/each}
        </div>
        <button class="btn-premium py-2.5 px-6 text-xs uppercase tracking-widest group">
          Deployment.
        </button>
      </div>
    </div>

    <!-- Elite Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
      {#each stats as stat}
        <div class="group relative">
          <div class="absolute -inset-1 bg-linear-to-r from-primary to-secondary rounded-3xl blur opacity-0 group-hover:opacity-10 transition duration-500"></div>
          <div class="glass relative p-7 rounded-3xl border border-slate-200/50 dark:border-white/5 flex flex-col justify-between h-full hover:scale-[1.02] transition-transform duration-300">
            <div class="flex items-start justify-between">
              <div class={cn("p-4 rounded-2xl ring-1 ring-white/20 shadow-inner", stat.bg, stat.color)}>
                <stat.icon class="w-6 h-6" />
              </div>
              <div class="bg-slate-50 dark:bg-white/5 px-2.5 py-1 rounded-full border border-slate-200/50 dark:border-white/5">
                <span class={cn("text-[10px] font-black", stat.change.startsWith("+") ? "text-emerald-500" : "text-rose-500")}>
                  {stat.change}
                </span>
              </div>
            </div>
            
            <div class="mt-8">
              <p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mb-1 leading-none">{stat.name}</p>
              <h2 class="text-3xl font-black text-slate-900 dark:text-white tracking-tighter">{stat.value}</h2>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Secondary Console Area -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      <!-- High-Fidelity Performance Chart -->
      <div class="lg:col-span-8">
        <div class="glass rounded-3xl border border-slate-200/50 dark:border-white/5 overflow-hidden flex flex-col h-full bg-white/30 dark:bg-slate-900/40">
          <div class="px-8 py-7 border-b border-slate-200/50 dark:border-white/5 flex items-center justify-between bg-white/20 dark:bg-white/5">
            <div>
              <h3 class="text-lg font-black text-slate-900 dark:text-white leading-none">Global Performance.</h3>
              <p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mt-2">Requests per microservice (last 30d)</p>
            </div>
            <div class="flex items-center gap-2">
              <button class="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-white/5 text-slate-400 transition-colors">
                <TrendingUp size={16} />
              </button>
            </div>
          </div>
          <div class="p-8 pb-4">
            <AreaChart data={chartData} height={320} />
          </div>
        </div>
      </div>

      <!-- Holographic Security Panel -->
      <div class="lg:col-span-4 space-y-8">
        <div class="relative group h-full">
          <!-- Holographic Glow Base -->
          <div class="absolute -inset-1 bg-linear-to-r from-emerald-500 via-primary to-secondary rounded-[32px] blur opacity-25 group-hover:opacity-40 transition duration-1000"></div>
          
          <div class="relative h-full glass rounded-[32px] overflow-hidden border border-white/20 bg-slate-900/90 p-8 flex flex-col text-white">
            <!-- Hive pattern overlay for widget -->
            <div class="absolute inset-0 bg-hive opacity-10 pointer-events-none"></div>
            
            <div class="relative z-10 flex flex-col h-full">
              <div class="flex items-center justify-between mb-8">
                <div class="w-12 h-12 rounded-2xl bg-emerald-500/20 border border-emerald-500/30 flex items-center justify-center text-emerald-400">
                  <Shield class="w-6 h-6 fill-emerald-500/20" />
                </div>
                <div class="px-3 py-1 rounded-full bg-emerald-500/20 text-emerald-400 text-[10px] font-black uppercase tracking-widest border border-emerald-500/30">
                  Secured
                </div>
              </div>

              <h3 class="text-3xl font-black tracking-tight leading-none mb-4">Elite Audit.</h3>
              <p class="text-sm font-medium text-slate-400 leading-relaxed mb-10">All systems are currently operating under the Kodia Security Protocol v2.1.</p>
              
              <div class="mt-auto space-y-6">
                <div class="p-5 rounded-2xl bg-white/5 border border-white/10 backdrop-blur-md">
                  <div class="flex items-center justify-between text-[10px] font-black uppercase tracking-[0.2em] mb-3 opacity-60">
                    <span>Audit Score</span>
                    <span class="text-emerald-400">Perfect.</span>
                  </div>
                  <div class="w-full h-2 bg-white/10 rounded-full overflow-hidden p-0.5 border border-white/5">
                    <div class="h-full bg-linear-to-r from-emerald-500 to-primary rounded-full w-full"></div>
                  </div>
                </div>

                <button class="w-full py-4 rounded-2xl bg-white text-slate-950 font-black text-xs uppercase tracking-[0.2em] hover:scale-105 active:scale-95 transition-all flex items-center justify-center gap-3">
                  Initiate Scan.
                  <ArrowUpRight class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Final Activities Row -->
    <div class="glass rounded-3xl border border-slate-200/50 dark:border-white/5 p-10 bg-white/40 dark:bg-slate-900/20 backdrop-blur-xl">
      <div class="flex items-center justify-between mb-10">
        <div>
          <h3 class="text-xl font-black text-slate-900 dark:text-white leading-none">Event Stream.</h3>
          <p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mt-3">Live activity from the Colony.</p>
        </div>
        <button class="text-[10px] font-black uppercase tracking-widest text-primary hover:underline">View Archive.</button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-10">
        {#each Array(3) as _, i}
          <div class="flex gap-5 group cursor-pointer active:scale-95 transition-all">
            <div class="w-12 h-12 rounded-2xl bg-slate-100 dark:bg-white/5 border border-slate-200/50 dark:border-white/5 flex items-center justify-center shrink-0 group-hover:border-primary/30 transition-colors">
              <Activity size={20} class="text-slate-400 dark:text-slate-500 group-hover:text-primary transition-colors" />
            </div>
            <div>
              <p class="text-xs font-black uppercase tracking-widest text-slate-900 dark:text-white leading-tight">Config Update.</p>
              <p class="text-[10px] text-slate-400 font-bold mt-1.5 uppercase">Admin • 2 hours ago</p>
              <div class="mt-3 px-2 py-0.5 rounded-md bg-secondary/5 text-[9px] font-black text-secondary uppercase tracking-widest w-fit">
                Kodia-Cloud
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>
</AdminLayout>
