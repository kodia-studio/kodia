<script lang="ts">
  import AdminLayout from "$lib/components/layouts/AdminLayout.svelte";
  import Container from "$lib/components/shared/Container.svelte";
  import Grid from "$lib/components/shared/Grid.svelte";
  import ChartCard from "$lib/components/charts/ChartCard.svelte";
  import AreaChart from "$lib/components/charts/AreaChart.svelte";
  import { Users, Shield, Zap, Database, TrendingUp, ArrowUpRight } from "lucide-svelte";

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
  <div class="space-y-8">
    <div>
      <h1 class="text-3xl font-bold font-heading tracking-tight">System Overview</h1>
      <p class="text-muted-foreground mt-1.5">Real-time performance and user metrics for Kodia Framework.</p>
    </div>

    <!-- Stats Grid -->
    <Grid cols={1} md={2} lg={4} gap="md">
      {#each stats as stat}
        <div class="card p-6 flex items-start justify-between">
          <div>
            <p class="text-sm font-medium text-muted-foreground">{stat.name}</p>
            <h2 class="text-2xl font-bold mt-1">{stat.value}</h2>
            <div class="flex items-center gap-1 mt-2 text-xs font-semibold">
              <span class={stat.change.startsWith("+") ? "text-green-500" : "text-amber-500"}>
                {stat.change}
              </span>
              <span class="text-muted-foreground font-normal">vs last month</span>
            </div>
          </div>
          <div class="p-3 rounded-xl {stat.bg} {stat.color}">
            <stat.icon class="w-5 h-5" />
          </div>
        </div>
      {/each}
    </Grid>

    <!-- Charts Area -->
    <Grid cols={1} lg={12} gap="lg">
      <div class="lg:col-span-8">
        <ChartCard 
          title="Overall Performance" 
          subtitle="4.8k Requests" 
        >
          {#snippet actions()}
            <button class="flex items-center gap-2 px-3 py-1.5 bg-muted/50 hover:bg-muted rounded-lg text-xs font-bold transition-colors">
              <TrendingUp class="w-3.5 h-3.5" />
              Analyze
            </button>
          {/snippet}
          
          <AreaChart data={chartData} height={350} />
        </ChartCard>
      </div>

      <div class="lg:col-span-4 space-y-6">
        <div class="card bg-primary text-primary-foreground overflow-hidden relative group">
          <div class="absolute -top-[20%] -right-[10%] w-[60%] h-[60%] rounded-full bg-white/10 blur-[80px] group-hover:scale-125 transition-transform duration-700"></div>
          
          <div class="relative z-10">
            <h3 class="text-sm font-bold uppercase tracking-widest opacity-80">Security Audit</h3>
            <p class="text-2xl font-bold mt-2">All systems clear.</p>
            <div class="mt-6 p-4 rounded-xl bg-white/10 backdrop-blur-md border border-white/10">
              <div class="flex items-center justify-between text-xs font-bold mb-2">
                <span>Compliance</span>
                <span>100%</span>
              </div>
              <div class="w-full h-1.5 bg-white/20 rounded-full overflow-hidden">
                <div class="h-full bg-white w-full"></div>
              </div>
            </div>
            <button class="mt-8 w-full py-3 rounded-xl bg-white text-primary font-bold text-sm hover:bg-primary-50 transition-colors flex items-center justify-center gap-2">
              View Report
              <ArrowUpRight class="w-4 h-4" />
            </button>
          </div>
        </div>

        <div class="card p-6">
          <h3 class="text-sm font-bold text-muted-foreground uppercase tracking-widest leading-none mb-6">
            Recent Activities
          </h3>
          <div class="space-y-6">
            {#each Array(3) as _, i}
              <div class="flex gap-4">
                <div class="w-2 h-2 mt-1.5 rounded-full bg-primary shrink-0 ring-4 ring-primary/10"></div>
                <div>
                  <p class="text-sm font-semibold leading-tight">Config updated by Admin</p>
                  <p class="text-xs text-muted-foreground mt-1">2 hours ago • Project: Kodia-Cloud</p>
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </Grid>
  </div>
</AdminLayout>
