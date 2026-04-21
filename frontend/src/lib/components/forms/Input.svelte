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
    <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-focus-within:text-primary transition-all duration-300">
      {@render icon()}
    </div>
  {/if}

  <input
    class={cn(
      "flex h-12 w-full rounded-2xl border transition-all duration-300",
      "bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm",
      "text-sm font-semibold tracking-tight",
      "ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-slate-400 dark:placeholder:text-slate-500 disabled:cursor-not-allowed disabled:opacity-50",
      icon && "pl-12",
      !icon && "px-5",
      error 
        ? "border-destructive focus-visible:border-destructive focus-visible:ring-2 focus-visible:ring-destructive/20" 
        : "border-slate-200 dark:border-slate-800 focus-visible:outline-none focus:border-primary focus:ring-4 focus:ring-primary/10 hover:border-slate-300 dark:hover:border-slate-700",
      className
    )}
    {...props}
    bind:value
  />
</div>
