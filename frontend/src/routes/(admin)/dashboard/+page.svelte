<script lang="ts">
  import AreaChart from "$lib/components/charts/AreaChart.svelte";
  import BarChart from "$lib/components/charts/BarChart.svelte";
  import PieChart from "$lib/components/charts/PieChart.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import { cn } from "$lib/utils/styles";
  import {
    DollarSign, ShoppingCart, TrendingUp, Target,
    Package, Tag, Zap, ArrowUpRight, ArrowDownRight, Download
  } from "lucide-svelte";

  // Types
  type ProductStatus = "active" | "low_stock" | "out_of_stock";
  type TxStatus = "completed" | "processing" | "failed" | "refunded";

  // Mock Data: 90 days of revenue with upward trend and weekly seasonality
  const ALL_REVENUE_DATA = Array.from({ length: 90 }, (_, i) => {
    const base = 12_000_000 + Math.sin(i / 7) * 3_000_000;
    const noise = (Math.random() - 0.5) * 2_000_000;
    const trend = i * 50_000;
    return {
      date: new Date(Date.now() - (89 - i) * 24 * 60 * 60 * 1000),
      value: Math.max(5_000_000, Math.floor(base + noise + trend))
    };
  });

  const KPI_CARDS = [
    {
      label: "Total Revenue",
      value: "Rp 847.2 Jt",
      change: "+14.2%",
      trend: "up" as const,
      icon: DollarSign,
      color: "text-blue-500",
      bg: "bg-blue-500/10"
    },
    {
      label: "Total Orders",
      value: "3,847",
      change: "+8.5%",
      trend: "up" as const,
      icon: ShoppingCart,
      color: "text-emerald-500",
      bg: "bg-emerald-500/10"
    },
    {
      label: "Avg. Order Value",
      value: "Rp 220.3 Rb",
      change: "+5.3%",
      trend: "up" as const,
      icon: TrendingUp,
      color: "text-amber-500",
      bg: "bg-amber-500/10"
    },
    {
      label: "Conversion Rate",
      value: "3.68%",
      change: "-0.4%",
      trend: "down" as const,
      icon: Target,
      color: "text-rose-500",
      bg: "bg-rose-500/10"
    }
  ];

  const CATEGORY_SALES = [
    { x: "Gadgets", y: 284_500_000 },
    { x: "Fashion", y: 198_300_000 },
    { x: "Home", y: 156_700_000 },
    { x: "Beauty", y: 124_100_000 },
    { x: "Sports", y: 83_600_000 }
  ];

  const TOP_PRODUCTS: Array<{ name: string; category: string; unitsSold: number; revenue: number; status: ProductStatus }> = [
    { name: "iPhone 15 Pro Max", category: "Gadgets", unitsSold: 142, revenue: 284_500_000, status: "active" },
    { name: "Nike Air Max 2024", category: "Fashion", unitsSold: 237, revenue: 118_500_000, status: "active" },
    { name: "Dyson V15 Detect", category: "Home", unitsSold: 48, revenue: 96_000_000, status: "low_stock" },
    { name: "La Mer Moisturiser", category: "Beauty", unitsSold: 89, revenue: 62_300_000, status: "active" },
    { name: "Garmin Fenix 7", category: "Sports", unitsSold: 33, revenue: 49_500_000, status: "out_of_stock" }
  ];

  const RECENT_TRANSACTIONS: Array<{ orderId: string; customer: string; items: number; total: number; status: TxStatus; date: string }> = [
    { orderId: "#ORD-7841", customer: "Budi Santoso", items: 3, total: 4_250_000, status: "completed", date: "20 Mei 2026" },
    { orderId: "#ORD-7840", customer: "Siti Rahayu", items: 1, total: 1_899_000, status: "processing", date: "20 Mei 2026" },
    { orderId: "#ORD-7839", customer: "Ahmad Fauzi", items: 5, total: 7_650_000, status: "completed", date: "19 Mei 2026" },
    { orderId: "#ORD-7838", customer: "Dewi Lestari", items: 2, total: 3_100_000, status: "failed", date: "19 Mei 2026" },
    { orderId: "#ORD-7837", customer: "Reza Pratama", items: 1, total: 890_000, status: "refunded", date: "18 Mei 2026" },
    { orderId: "#ORD-7836", customer: "Nurul Hidayah", items: 4, total: 5_420_000, status: "completed", date: "18 Mei 2026" }
  ];

  const QUICK_STATS = {
    salesTarget: { current: 847_200_000, target: 1_000_000_000 },
    topCategory: "Gadgets",
    activePromotions: 5
  };

  // Reactive State
  let selectedPeriod = $state<"7d" | "30d" | "90d">("30d");

  // Derived State
  const chartData = $derived.by(() => {
    const days = selectedPeriod === "7d" ? 7 : selectedPeriod === "30d" ? 30 : 90;
    return ALL_REVENUE_DATA.slice(-days);
  });

  const salesTargetPct = $derived(
    Math.round((QUICK_STATS.salesTarget.current / QUICK_STATS.salesTarget.target) * 100)
  );

  // Helper Functions
  function formatIDR(value: number): string {
    if (value >= 1_000_000_000) return `Rp ${(value / 1_000_000_000).toFixed(1)} M`;
    if (value >= 1_000_000) return `Rp ${(value / 1_000_000).toFixed(1)} Jt`;
    if (value >= 1_000) return `Rp ${(value / 1_000).toFixed(0)} Rb`;
    return `Rp ${value.toLocaleString("id-ID")}`;
  }

  function txStatusVariant(status: TxStatus): "success" | "warning" | "danger" | "ghost" {
    const map = { completed: "success", processing: "warning", failed: "danger", refunded: "ghost" } as const;
    return map[status];
  }

  function productStatusVariant(status: ProductStatus): "success" | "warning" | "danger" {
    const map = { active: "success", low_stock: "warning", out_of_stock: "danger" } as const;
    return map[status];
  }
</script>

<svelte:head>
  <title>Sales Dashboard | Kodia Console</title>
</svelte:head>

<div class="space-y-6 md:space-y-8 pb-4">
  <!-- Header -->
  <div class="flex flex-col gap-4 md:gap-6">
    <div>
      <h1 class="text-2xl md:text-3xl font-bold text-slate-900 dark:text-white">Sales Dashboard</h1>
      <p class="text-xs md:text-sm text-slate-500 dark:text-slate-400 mt-2 md:mt-1">Overview for May 2026 — all figures in Indonesian Rupiah (IDR)</p>
    </div>
    <div>
      <button class="btn-premium text-xs md:text-sm flex items-center justify-center md:justify-start gap-2 w-full md:w-fit px-4 md:px-6 py-3 md:py-2">
        <Download class="w-4 h-4" />
        <span>Export Report</span>
      </button>
    </div>
  </div>

  <!-- KPI Row -->
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-5">
    {#each KPI_CARDS as kpi}
      <div class="bg-white dark:bg-slate-900/50 rounded-lg md:rounded-xl border border-slate-200 dark:border-slate-800 p-4 md:p-6 hover:border-slate-300 dark:hover:border-slate-700 transition-all duration-200 group">
        <div class="flex items-start justify-between mb-3 md:mb-4">
          <div class={cn("p-2 md:p-2.5 rounded-lg md:rounded-xl", kpi.bg)}>
            <kpi.icon class={cn("w-4 h-4 md:w-5 md:h-5", kpi.color)} />
          </div>
          <span class={cn(
            "flex items-center gap-0.5 text-[10px] md:text-xs font-bold px-2 py-1 rounded-full",
            kpi.trend === "up"
              ? "bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400"
              : "bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400"
          )}>
            {#if kpi.trend === "up"}
              <ArrowUpRight class="w-3 h-3" />
            {:else}
              <ArrowDownRight class="w-3 h-3" />
            {/if}
            {kpi.change}
          </span>
        </div>
        <p class="text-[10px] md:text-xs font-semibold uppercase tracking-widest text-slate-500 dark:text-slate-400 mb-1 md:mb-2">
          {kpi.label}
        </p>
        <p class="text-xl md:text-2xl font-bold text-slate-900 dark:text-white tracking-tight">
          {kpi.value}
        </p>
      </div>
    {/each}
  </div>

  <!-- Sales by Category - 2 Column Grid -->
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 md:gap-6">
    <!-- Bar Chart Card -->
    <div class="bg-white dark:bg-slate-900/50 rounded-lg md:rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="flex items-center justify-between px-4 md:px-6 py-3 md:py-4 border-b border-slate-200 dark:border-slate-800">
        <div>
          <h2 class="text-sm md:text-base font-bold text-slate-900 dark:text-white">Sales by Category</h2>
          <p class="text-[10px] md:text-xs text-slate-500 dark:text-slate-400 mt-0.5">Revenue breakdown</p>
        </div>
      </div>
      <div class="p-3 md:p-5">
        <BarChart data={CATEGORY_SALES} height={220} />
      </div>
    </div>

    <!-- Pie Chart Card -->
    <div class="bg-white dark:bg-slate-900/50 rounded-lg md:rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="flex items-center justify-between px-4 md:px-6 py-3 md:py-4 border-b border-slate-200 dark:border-slate-800">
        <div>
          <h2 class="text-sm md:text-base font-bold text-slate-900 dark:text-white">Category Distribution</h2>
          <p class="text-[10px] md:text-xs text-slate-500 dark:text-slate-400 mt-0.5">Market share by category</p>
        </div>
      </div>
      <div class="p-3 md:p-5">
        <PieChart data={CATEGORY_SALES} />
      </div>
    </div>
  </div>

  <!-- Revenue Chart + Quick Stats -->
  <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-6">
    <!-- Revenue Trend Chart: col-span-2 -->
    <div class="lg:col-span-2 bg-white dark:bg-slate-900/50 rounded-lg md:rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden flex flex-col">
      <!-- Card Header -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 md:gap-4 px-4 md:px-6 py-3 md:py-4 border-b border-slate-200 dark:border-slate-800 shrink-0">
        <div>
          <h2 class="text-sm md:text-base font-bold text-slate-900 dark:text-white">Revenue Trend</h2>
          <p class="text-[10px] md:text-xs text-slate-500 dark:text-slate-400 mt-0.5">Daily revenue for the selected period</p>
        </div>

        <!-- Period Selector Tabs -->
        <div class="flex items-center gap-1 p-1 bg-slate-100 dark:bg-slate-800 rounded-lg w-fit">
          {#each (["7d", "30d", "90d"] as const) as period}
            <button
              onclick={() => selectedPeriod = period}
              class={cn(
                "px-2 md:px-3 py-1.5 rounded-md text-[10px] md:text-xs font-semibold transition-all duration-150",
                selectedPeriod === period
                  ? "bg-white dark:bg-slate-700 text-slate-900 dark:text-white shadow-sm"
                  : "text-slate-500 hover:text-slate-700 dark:hover:text-slate-300"
              )}
            >
              {period.toUpperCase()}
            </button>
          {/each}
        </div>
      </div>

      <!-- Chart Body -->
      <div class="p-3 md:p-4 flex-1">
        <AreaChart data={chartData} height={250} color="#3b82f6" />
      </div>
    </div>

    <!-- Quick Stats: col-span-1 -->
    <div class="bg-white dark:bg-slate-900/50 rounded-lg md:rounded-xl border border-slate-200 dark:border-slate-800 p-4 md:p-5 flex flex-col gap-4 md:gap-5">
      <h2 class="text-sm md:text-base font-bold text-slate-900 dark:text-white">Quick Stats</h2>

      <!-- Sales Target Progress -->
      <div class="space-y-2 md:space-y-3">
        <div class="flex items-center justify-between mb-1 md:mb-2">
          <span class="text-[10px] md:text-xs font-semibold uppercase tracking-widest text-slate-500 dark:text-slate-400">
            Sales Target
          </span>
          <span class="text-xs md:text-sm font-bold text-slate-900 dark:text-white">
            {salesTargetPct}%
          </span>
        </div>
        <div class="h-2 md:h-2.5 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
          <div
            class="h-full bg-linear-to-r from-primary to-secondary rounded-full transition-all duration-700"
            style="width: {salesTargetPct}%"
          ></div>
        </div>
        <div class="flex items-center justify-between mt-1 md:mt-1.5">
          <span class="text-[9px] md:text-[10px] text-slate-400">
            {formatIDR(QUICK_STATS.salesTarget.current)}
          </span>
          <span class="text-[9px] md:text-[10px] text-slate-400">
            {formatIDR(QUICK_STATS.salesTarget.target)}
          </span>
        </div>
      </div>

      <!-- Divider -->
      <div class="h-px bg-slate-100 dark:bg-slate-800"></div>

      <!-- Top Category -->
      <div class="flex items-center justify-between gap-2">
        <div class="flex items-center gap-2 md:gap-3 min-w-0">
          <div class="p-1.5 md:p-2 bg-amber-500/10 rounded-lg shrink-0">
            <Tag class="w-3 h-3 md:w-4 md:h-4 text-amber-500" />
          </div>
          <div class="min-w-0">
            <p class="text-[9px] md:text-[10px] font-semibold uppercase tracking-widest text-slate-500 dark:text-slate-400">
              Top Category
            </p>
            <p class="text-xs md:text-sm font-bold text-slate-900 dark:text-white truncate">
              {QUICK_STATS.topCategory}
            </p>
          </div>
        </div>
        <span class="text-[10px] md:text-xs font-bold text-amber-600 dark:text-amber-400 bg-amber-500/10 px-2 py-1 rounded-lg shrink-0">
          #1
        </span>
      </div>

      <!-- Active Promotions -->
      <div class="flex items-center justify-between gap-2">
        <div class="flex items-center gap-2 md:gap-3 min-w-0">
          <div class="p-1.5 md:p-2 bg-purple-500/10 rounded-lg shrink-0">
            <Zap class="w-3 h-3 md:w-4 md:h-4 text-purple-500" />
          </div>
          <div class="min-w-0">
            <p class="text-[9px] md:text-[10px] font-semibold uppercase tracking-widest text-slate-500 dark:text-slate-400">
              Active Promos
            </p>
            <p class="text-xs md:text-sm font-bold text-slate-900 dark:text-white">
              {QUICK_STATS.activePromotions} Running
            </p>
          </div>
        </div>
        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shrink-0"></span>
      </div>

      <!-- View Details Button -->
      <button class="w-full mt-auto py-2 md:py-2.5 rounded-lg border border-slate-200 dark:border-slate-800 text-[10px] md:text-xs font-semibold text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors">
        View Full Report
      </button>
    </div>
  </div>

  <!-- Top Products + Recent Transactions -->
  <div class="grid grid-cols-1 xl:grid-cols-2 gap-4 md:gap-6">
    <!-- Top Products Table -->
    <div class="bg-white dark:bg-slate-900/50 rounded-lg md:rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="flex items-center justify-between p-4 md:p-6 border-b border-slate-200 dark:border-slate-800">
        <h2 class="text-sm md:text-base font-bold text-slate-900 dark:text-white">Top Products</h2>
        <button class="text-[10px] md:text-xs font-semibold text-primary hover:text-primary/80 transition-colors">
          View All
        </button>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-slate-100 dark:border-slate-800 bg-slate-50/80 dark:bg-slate-800/30">
              <th class="px-3 md:px-5 py-2 md:py-3 text-left text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Product</th>
              <th class="hidden sm:table-cell px-3 md:px-5 py-2 md:py-3 text-left text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Category</th>
              <th class="px-3 md:px-5 py-2 md:py-3 text-right text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Units</th>
              <th class="hidden md:table-cell px-3 md:px-5 py-2 md:py-3 text-right text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Revenue</th>
              <th class="px-3 md:px-5 py-2 md:py-3 text-center text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-800/50">
            {#each TOP_PRODUCTS as product}
              <tr class="hover:bg-slate-50 dark:hover:bg-slate-800/40 transition-colors">
                <td class="px-3 md:px-5 py-2.5 md:py-3.5">
                  <div class="flex items-center gap-2">
                    <div class="w-6 h-6 md:w-7 md:h-7 rounded-lg bg-slate-100 dark:bg-slate-800 flex items-center justify-center shrink-0">
                      <Package class="w-3 h-3 md:w-3.5 md:h-3.5 text-slate-500" />
                    </div>
                    <span class="text-xs md:text-sm font-semibold text-slate-900 dark:text-white truncate max-w-24 md:max-w-30">
                      {product.name}
                    </span>
                  </div>
                </td>
                <td class="hidden sm:table-cell px-3 md:px-5 py-2.5 md:py-3.5 text-[10px] md:text-xs text-slate-500 dark:text-slate-400">
                  {product.category}
                </td>
                <td class="px-3 md:px-5 py-2.5 md:py-3.5 text-right text-xs md:text-sm font-bold text-slate-900 dark:text-white">
                  {product.unitsSold.toLocaleString()}
                </td>
                <td class="hidden md:table-cell px-3 md:px-5 py-2.5 md:py-3.5 text-right text-xs md:text-sm font-bold text-slate-900 dark:text-white">
                  {formatIDR(product.revenue)}
                </td>
                <td class="px-3 md:px-5 py-2.5 md:py-3.5 text-center">
                  <Badge variant={productStatusVariant(product.status)}>
                    {product.status.replace(/_/g, " ")}
                  </Badge>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- Recent Transactions Table -->
    <div class="bg-white dark:bg-slate-900/50 rounded-lg md:rounded-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="flex items-center justify-between p-4 md:p-6 border-b border-slate-200 dark:border-slate-800">
        <h2 class="text-sm md:text-base font-bold text-slate-900 dark:text-white">Recent Transactions</h2>
        <button class="text-[10px] md:text-xs font-semibold text-primary hover:text-primary/80 transition-colors">
          View All
        </button>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-slate-100 dark:border-slate-800 bg-slate-50/80 dark:bg-slate-800/30">
              <th class="px-3 md:px-5 py-2 md:py-3 text-left text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Order</th>
              <th class="px-3 md:px-5 py-2 md:py-3 text-left text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Customer</th>
              <th class="hidden sm:table-cell px-3 md:px-5 py-2 md:py-3 text-right text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Total</th>
              <th class="px-3 md:px-5 py-2 md:py-3 text-center text-[9px] md:text-[10px] font-black uppercase tracking-widest text-slate-500">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100 dark:divide-slate-800/50">
            {#each RECENT_TRANSACTIONS as tx}
              <tr class="hover:bg-slate-50 dark:hover:bg-slate-800/40 transition-colors">
                <td class="px-3 md:px-5 py-2.5 md:py-3.5">
                  <div>
                    <span class="text-xs md:text-sm font-bold text-primary">{tx.orderId}</span>
                    <p class="text-[9px] md:text-[10px] text-slate-400 mt-0.5">{tx.date}</p>
                  </div>
                </td>
                <td class="px-3 md:px-5 py-2.5 md:py-3.5">
                  <div>
                    <p class="text-xs md:text-sm font-semibold text-slate-900 dark:text-white">
                      {tx.customer}
                    </p>
                    <p class="text-[9px] md:text-[10px] text-slate-400 mt-0.5">{tx.items} item{tx.items > 1 ? "s" : ""}</p>
                  </div>
                </td>
                <td class="hidden sm:table-cell px-3 md:px-5 py-2.5 md:py-3.5 text-right text-xs md:text-sm font-bold text-slate-900 dark:text-white">
                  {formatIDR(tx.total)}
                </td>
                <td class="px-3 md:px-5 py-2.5 md:py-3.5 text-center">
                  <Badge variant={txStatusVariant(tx.status)}>
                    {tx.status}
                  </Badge>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
