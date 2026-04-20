<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import { MoreHorizontal } from "lucide-svelte";
  import { fade } from "svelte/transition";

  interface Props {
    title: string;
    subtitle?: string;
    class?: string;
    children?: import('svelte').Snippet;
    actions?: import('svelte').Snippet;
  }

  let { 
    title, 
    subtitle, 
    class: className, 
    children, 
    actions 
  }: Props = $props();
</script>

<div 
  class={cn("card relative overflow-hidden group", className)}
  in:fade={{ duration: 400 }}
>
  <div class="flex items-start justify-between mb-6">
    <div>
      <h3 class="text-sm font-bold text-muted-foreground uppercase tracking-widest leading-none mb-1.5">
        {title}
      </h3>
      {#if subtitle}
        <p class="text-2xl font-bold font-heading tracking-tight">
          {subtitle}
        </p>
      {/if}
    </div>
    
    {#if actions}
      {@render actions()}
    {:else}
      <button class="p-2 hover:bg-muted rounded-xl transition-colors text-muted-foreground">
        <MoreHorizontal class="w-5 h-5" />
      </button>
    {/if}
  </div>

  <div class="relative w-full">
    {@render children?.()}
  </div>
</div>
