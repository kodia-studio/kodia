<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import { 
    LayoutDashboard, 
    Users, 
    Settings, 
    ShieldCheck, 
    BarChart3, 
    Layers,
    Activity,
    Zap,
    ChevronLeft,
    LogOut
  } from "lucide-svelte";
  import { page } from "$app/state";
  import { authStore } from "$lib/stores/auth.store";
  import { goto } from "$app/navigation";

  let { collapsed = $bindable(false) } = $props();

  const menuItems = [
    { name: "Dashboard", icon: LayoutDashboard, href: "/dashboard" },
    { name: "Pulse", icon: Activity, href: "/pulse" },
    { name: "Users", icon: Users, href: "/users" },
    { name: "Settings", icon: Settings, href: "/settings" },
  ];

  async function handleLogout() {
    authStore.logout();
    await goto("/login");
  }
</script>

<aside 
  class={cn(
    "fixed left-0 top-0 bottom-0 transition-all duration-300 z-50 border-r",
    "bg-white/40 dark:bg-slate-900/40 backdrop-blur-3xl border-slate-200/50 dark:border-white/5 shadow-2xl flex flex-col",
    collapsed ? "w-20" : "w-64"
  )}
>
  <!-- Sidebar Header / Brand -->
  <div class="h-24 flex items-center px-7 relative">
    <div class="absolute inset-x-0 bottom-0 h-px bg-linear-to-r from-transparent via-slate-200/50 dark:via-white/5 to-transparent"></div>
    
    <a href="/" class="flex items-center gap-4 overflow-hidden group">
      <div class="min-w-[44px] h-11 bg-linear-to-tr from-primary to-secondary rounded-2xl flex items-center justify-center shadow-lg shadow-primary/20 group-hover:scale-105 transition-transform">
        <Zap class="text-white w-6 h-6 fill-white/20" />
      </div>
      {#if !collapsed}
        <div class="flex flex-col">
          <span class="text-xl font-black font-heading tracking-tighter text-slate-900 dark:text-white leading-none">
            Kodia Console<span class="text-primary">.</span>
          </span>
          <span class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 mt-1">
            Colony v1.0
          </span>
        </div>
      {/if}
    </a>
  </div>

  <!-- Navigation -->
  <nav class="flex-1 px-4 py-8 space-y-3 overflow-y-auto custom-scrollbar">
    {#each menuItems as item}
      {@const isActive = page.url.pathname === item.href}
      <a 
        href={item.href}
        class={cn(
          "flex items-center gap-4 px-4 py-3.5 rounded-2xl transition-all duration-300 group relative overflow-hidden",
          isActive 
            ? "bg-primary text-white shadow-[0_8px_30px_rgb(59,130,246,0.3)] ring-1 ring-white/20" 
            : "hover:bg-slate-100/50 dark:hover:bg-white/5 text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white"
        )}
      >
        {#if isActive}
          <div class="absolute inset-0 bg-linear-to-r from-primary to-secondary opacity-100"></div>
          <div class="absolute -right-2 top-0 bottom-0 w-8 bg-white/20 blur-xl rotate-12"></div>
        {/if}

        <item.icon class={cn("w-5.5 h-5.5 relative z-10", isActive ? "text-white" : "group-hover:scale-110 transition-transform")} />
        
        {#if !collapsed}
          <span class="font-black text-xs uppercase tracking-widest relative z-10">{item.name}</span>
        {/if}

        {#if isActive && !collapsed}
          <div class="ml-auto w-1.5 h-1.5 bg-white rounded-full relative z-10 animate-pulse"></div>
        {/if}
      </a>
    {/each}
  </nav>

  <!-- Bottom Section -->
  <div class="px-4 py-8 relative">
    <div class="absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-slate-200/50 dark:via-white/5 to-transparent"></div>
    
    <div class="space-y-2">
      <button 
        onclick={() => (collapsed = !collapsed)}
        class="w-full flex items-center gap-4 px-4 py-3.5 rounded-2xl hover:bg-slate-100/50 dark:hover:bg-white/5 text-slate-500 dark:text-slate-400 transition-all active:scale-95 group"
      >
        <ChevronLeft class={cn("w-5.5 h-5.5 transition-transform duration-500", collapsed && "rotate-180")} />
        {#if !collapsed}
          <span class="text-xs font-black uppercase tracking-widest group-hover:text-slate-900 dark:group-hover:text-white">Minimize.</span>
        {/if}
      </button>
      
      <button 
        onclick={handleLogout}
        class="w-full flex items-center gap-4 px-4 py-3.5 rounded-2xl hover:bg-rose-500/10 text-slate-500 dark:text-slate-400 hover:text-rose-500 transition-all active:scale-95 group"
      >
        <LogOut class="w-5.5 h-5.5" />
        {#if !collapsed}
          <span class="text-xs font-black uppercase tracking-widest leading-none pt-0.5">Termination.</span>
        {/if}
      </button>
    </div>
  </div>
</aside>
