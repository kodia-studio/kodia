<script lang="ts">
  import AreaChart from "$lib/components/charts/AreaChart.svelte";
  import { cn } from "$lib/utils/styles";
  import { Users, Shield, Zap, Database, TrendingUp } from "lucide-svelte";

  const stats = [
    { name: "Total Users", value: "12,842", change: "+12%", trend: "up", icon: Users, color: "text-blue-500", bg: "bg-blue-500/10" },
    { name: "Security Score", value: "98.2%", change: "+0.5%", trend: "up", icon: Shield, color: "text-green-500", bg: "bg-green-500/10" },
    { name: "API Latency", value: "24ms", change: "-4ms", trend: "down", icon: Zap, color: "text-amber-500", bg: "bg-amber-500/10" },
    { name: "Storage Used", value: "1.2TB", change: "+8%", trend: "up", icon: Database, color: "text-purple-500", bg: "bg-purple-500/10" },
  ];

  const chartData = Array.from({ length: 30 }, (_, i) => ({
    date: new Date(Date.now() - (29 - i) * 24 * 60 * 60 * 1000),
    value: Math.floor(Math.random() * 5000) + 2000
  }));

  const activities = [
    { name: "Database Backup", time: "2 hours ago", category: "System", icon: Database },
    { name: "Security Update", time: "4 hours ago", category: "Security", icon: Shield },
    { name: "User Login", time: "10 minutes ago", category: "Auth", icon: Users },
  ];
</script>

<div class="space-y-8">
  <!-- Header -->
  <div class="flex flex-col gap-2">
    <h1 class="text-3xl font-bold text-slate-900 dark:text-white">Dashboard</h1>
    <p class="text-sm text-slate-600 dark:text-slate-400">Welcome back! Here's your system overview.</p>
  </div>

  <!-- Stats Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    {#each stats as stat}
      <div class="bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800 p-6 hover:border-slate-300 dark:hover:border-slate-700 transition-colors">
        <div class="flex items-start justify-between mb-4">
          <div class={cn("p-3 rounded-lg", stat.bg)}>
            <stat.icon class={cn("w-5 h-5", stat.color)} />
          </div>
          <div class={cn("text-xs font-semibold px-2 py-1 rounded-full", stat.trend === "up" ? "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400" : "bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400")}>
            {stat.change}
          </div>
        </div>
        <p class="text-xs text-slate-500 dark:text-slate-400 font-medium mb-1">{stat.name}</p>
        <p class="text-2xl font-bold text-slate-900 dark:text-white">{stat.value}</p>
      </div>
    {/each}
  </div>

  <!-- Main Content Grid -->
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Chart -->
    <div class="lg:col-span-2 bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800">
      <div class="flex items-center justify-between p-6 border-b border-slate-200 dark:border-slate-800">
        <div>
          <h2 class="text-lg font-bold text-slate-900 dark:text-white">Performance</h2>
          <p class="text-xs text-slate-600 dark:text-slate-400 mt-1">Last 30 days</p>
        </div>
        <button class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors">
          <TrendingUp class="w-5 h-5 text-slate-400" />
        </button>
      </div>
      <div class="p-6">
        <AreaChart data={chartData} height={320} />
      </div>
    </div>

    <!-- Quick Stats Card -->
    <div class="bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800 p-6">
      <div class="flex items-center gap-3 mb-6">
        <div class="p-3 bg-green-100 dark:bg-green-900/30 rounded-lg">
          <Shield class="w-5 h-5 text-green-600 dark:text-green-400" />
        </div>
        <div>
          <p class="text-xs text-slate-500 dark:text-slate-400 font-medium">System Status</p>
          <p class="text-lg font-bold text-green-600 dark:text-green-400">Healthy</p>
        </div>
      </div>
      <div class="space-y-4">
        <div>
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-slate-600 dark:text-slate-400">CPU Usage</span>
            <span class="text-xs font-bold text-slate-900 dark:text-white">45%</span>
          </div>
          <div class="h-2 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
            <div class="h-full bg-blue-500 w-[45%]"></div>
          </div>
        </div>
        <div>
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-slate-600 dark:text-slate-400">Memory Usage</span>
            <span class="text-xs font-bold text-slate-900 dark:text-white">62%</span>
          </div>
          <div class="h-2 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
            <div class="h-full bg-amber-500 w-[62%]"></div>
          </div>
        </div>
        <button class="w-full mt-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-semibold rounded-lg transition-colors">
          View Details
        </button>
      </div>
    </div>
  </div>

  <!-- Recent Activities -->
  <div class="bg-white dark:bg-slate-900/50 rounded-xl border border-slate-200 dark:border-slate-800">
    <div class="flex items-center justify-between p-6 border-b border-slate-200 dark:border-slate-800">
      <h2 class="text-lg font-bold text-slate-900 dark:text-white">Recent Activity</h2>
      <button class="text-xs font-semibold text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300">View All</button>
    </div>
    <div class="divide-y divide-slate-200 dark:divide-slate-800">
      {#each activities as activity}
        <div class="p-4 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors flex items-center justify-between">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-slate-100 dark:bg-slate-800 rounded-lg">
              <activity.icon class="w-4 h-4 text-slate-600 dark:text-slate-400" />
            </div>
            <div>
              <p class="text-sm font-medium text-slate-900 dark:text-white">{activity.name}</p>
              <p class="text-xs text-slate-500 dark:text-slate-400">{activity.time}</p>
            </div>
          </div>
          <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 px-2 py-1 rounded">{activity.category}</span>
        </div>
      {/each}
    </div>
  </div>
</div>
