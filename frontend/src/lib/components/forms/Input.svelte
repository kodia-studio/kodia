<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import type { HTMLInputAttributes } from "svelte/elements";

  interface Props extends HTMLInputAttributes {
    error?: boolean;
    icon?: import('svelte').Snippet;
    class?: string;
  }

  let { 
    error = false, 
    icon, 
    class: className, 
    value = $bindable(),
    ...props 
  }: Props = $props();
</script>

<div class="relative w-full group">
  {#if icon}
    <div class="absolute left-3.5 top-1/2 -translate-y-1/2 text-muted-foreground group-focus-within:text-primary transition-colors">
      {@render icon()}
    </div>
  {/if}

  <input
    class={cn(
      "flex h-11 w-full rounded-xl border bg-slate-50/50 dark:bg-slate-900/50 px-4 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/20 focus-visible:border-primary disabled:cursor-not-allowed disabled:opacity-50 transition-all duration-200",
      icon && "pl-11",
      error ? "border-destructive focus-visible:border-destructive focus-visible:ring-destructive/20" : "border-input hover:border-muted-foreground/30",
      className
    )}
    {...props}
  />
</div>
