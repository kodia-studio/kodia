<script lang="ts">
  import { Search, Loader2, X } from "lucide-svelte";
  import { onMount } from "svelte";
  import { api } from "$lib/api/client";
  import { cn } from "$lib/utils/styles";

  interface Props {
    endpoint?: string;
    placeholder?: string;
    onresults?: (results: any) => void;
    class?: string;
    debounceMs?: number;
  }

  let { 
    endpoint = "/api/search", 
    placeholder = "Search everything...", 
    onresults, 
    class: className,
    debounceMs = 300 
  }: Props = $props();

  let query = $state("");
  let isLoading = $state(false);
  let timer: any;

  async function performSearch() {
    if (!query) {
      if (onresults) onresults([]);
      return;
    }

    isLoading = true;
    try {
      const results = await api.get(`${endpoint}?q=${encodeURIComponent(query)}`);
      if (onresults) onresults(results);
    } catch (error) {
      console.error("Search failed:", error);
    } finally {
      isLoading = false;
    }
  }

  function handleInput() {
    clearTimeout(timer);
    timer = setTimeout(() => {
      performSearch();
    }, debounceMs);
  }

  function clearSearch() {
    query = "";
    if (onresults) onresults([]);
  }
</script>

<div class={cn("relative group w-full", className)}>
  <div class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground group-focus-within:text-primary transition-colors">
    <Search class="w-4 h-4" />
  </div>

  <input
    type="text"
    bind:value={query}
    oninput={handleInput}
    {placeholder}
    class="w-full pl-10 pr-10 py-2.5 bg-muted/50 border border-border/50 rounded-xl text-sm focus:ring-2 focus:ring-primary/20 focus:bg-background outline-none transition-all"
  />

  <div class="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-2">
    {#if isLoading}
      <Loader2 class="w-4 h-4 animate-spin text-primary" />
    {:else if query}
      <button 
        onclick={clearSearch}
        class="p-0.5 hover:bg-muted rounded-md text-muted-foreground transition-colors"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    {/if}
  </div>
</div>
