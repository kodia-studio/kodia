<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import { 
    LayoutDashboard, 
    Users, 
    Settings, 
    ShieldCheck, 
    BarChart3, 
    Layers,
    ChevronLeft,
    LogOut,
    Activity
  } from "lucide-svelte";
  import { page } from "$app/state";

  let { collapsed = $bindable(false) } = $props();

  const menuItems = [
    { name: "Dashboard", icon: LayoutDashboard, href: "/admin" },
    { name: "Pulse", icon: Activity, href: "/admin/pulse" },
    { name: "Users", icon: Users, href: "/admin/users" },
    { name: "Analytics", icon: BarChart3, href: "/admin/analytics" },
    { name: "Security", icon: ShieldCheck, href: "/admin/security" },
    { name: "Settings", icon: Settings, href: "/admin/settings" },
  ];
</script>

<aside 
  class={cn(
    "fixed left-0 top-0 h-screen transition-all duration-300 z-50 bg-card border-r flex flex-col",
    collapsed ? "w-20" : "w-64"
  )}
>
  <div class="h-20 flex items-center px-6">
    <a href="/" class="flex items-center gap-3 overflow-hidden">
      <div class="min-w-[40px] h-10 accent-gradient rounded-xl flex items-center justify-center shadow-lg">
        <Layers class="text-white w-6 h-6" />
      </div>
      {#if !collapsed}
        <span class="text-2xl font-bold font-heading tracking-tight whitespace-nowrap">
          Kodia<span class="text-primary">.</span>
        </span>
      {/if}
    </a>
  </div>

  <nav class="flex-1 px-4 py-6 space-y-2">
    {#each menuItems as item}
      {@const isActive = page.url.pathname === item.href}
      <a 
        href={item.href}
        class={cn(
          "flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200 group",
          isActive 
            ? "bg-primary text-primary-foreground shadow-lg shadow-primary/20" 
            : "hover:bg-muted text-muted-foreground hover:text-foreground"
        )}
      >
        <item.icon class={cn("w-5 h-5", !isActive && "group-hover:scale-110 transition-transform")} />
        {#if !collapsed}
          <span class="font-medium text-sm">{item.name}</span>
        {/if}
      </a>
    {/each}
  </nav>

  <div class="p-4 border-t space-y-2">
    <button 
      onclick={() => (collapsed = !collapsed)}
      class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-muted text-muted-foreground transition-colors"
    >
      <ChevronLeft class={cn("w-5 h-5 transition-transform duration-300", collapsed && "rotate-180")} />
      {#if !collapsed}
        <span class="text-sm font-medium">Collapse</span>
      {/if}
    </button>
    
    <button class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-destructive/10 text-muted-foreground hover:text-destructive transition-colors">
      <LogOut class="w-5 h-5" />
      {#if !collapsed}
        <span class="text-sm font-medium">Logout</span>
      {/if}
    </button>
  </div>
</aside>
