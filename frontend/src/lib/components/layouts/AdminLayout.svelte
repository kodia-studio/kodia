<script lang="ts">
  import Sidebar from "../ui/Sidebar.svelte";
  import { cn } from "$lib/utils/styles";
  import { fade } from "svelte/transition";
  import { Search, Bell, User as UserIcon } from "lucide-svelte";

  let { children } = $props();
  let sidebarCollapsed = $state(false);
</script>

<div class="min-h-screen bg-muted/30">
  <Sidebar bind:collapsed={sidebarCollapsed} />

  <div 
    class={cn(
      "transition-all duration-300 min-h-screen flex flex-col",
      sidebarCollapsed ? "pl-20" : "pl-64"
    )}
  >
    <!-- Header -->
    <header class="h-20 border-b bg-background/80 backdrop-blur-md sticky top-0 z-40 px-8 flex items-center justify-between">
      <div class="relative w-96 hidden md:block">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <input 
          type="text" 
          placeholder="Search everything..." 
          class="w-full pl-10 pr-4 py-2 bg-muted/50 border-none rounded-xl text-sm focus:ring-2 focus:ring-primary/20 outline-none transition-all"
        />
      </div>

      <div class="flex items-center gap-4">
        <button class="p-2.5 hover:bg-muted rounded-xl relative transition-colors">
          <Bell class="w-5 h-5 text-muted-foreground" />
          <span class="absolute top-2.5 right-2.5 w-2 h-2 bg-primary rounded-full border-2 border-background"></span>
        </button>
        
        <div class="h-8 w-px bg-border mx-2"></div>

        <button class="flex items-center gap-3 p-1.5 hover:bg-muted rounded-xl transition-colors">
          <div class="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
            <UserIcon class="w-4 h-4" />
          </div>
          <div class="text-left hidden sm:block">
            <p class="text-sm font-semibold leading-tight">Admin User</p>
            <p class="text-xs text-muted-foreground">Level 4 Access</p>
          </div>
        </button>
      </div>
    </header>

    <!-- Main Content -->
    <main class="p-8 flex-1">
      {#key children}
        <div in:fade={{ duration: 200, delay: 100 }}>
          {@render children?.()}
        </div>
      {/key}
    </main>

    <footer class="px-8 py-6 border-t text-sm text-muted-foreground">
      <p>&copy; {new Date().getFullYear()} Kodia Admin Console</p>
    </footer>
  </div>
</div>
