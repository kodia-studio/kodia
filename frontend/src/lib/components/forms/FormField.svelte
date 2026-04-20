<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import { fade } from "svelte/transition";

  interface Props {
    label?: string;
    error?: string;
    description?: string;
    required?: boolean;
    class?: string;
    id?: string;
    children?: import('svelte').Snippet<[{ id: string }]>;
  }

  let { 
    label, 
    error, 
    description, 
    required = false, 
    class: className, 
    id = "field-" + Math.random().toString(36).substring(2, 9),
    children 
  }: Props = $props();
</script>

<div class={cn("space-y-2 w-full", className)}>
  {#if label}
    <div class="flex items-center justify-between">
      <label for={id} class="text-sm font-semibold tracking-tight text-foreground/90 cursor-pointer">
        {label}
        {#if required}
          <span class="text-destructive ml-0.5">*</span>
        {/if}
      </label>
    </div>
  {/if}

  <div class="relative">
    {@render children?.({ id })}
  </div>

  {#if description && !error}
    <p class="text-xs text-muted-foreground leading-relaxed">
      {description}
    </p>
  {/if}

  {#if error}
    <p 
      class="text-xs font-medium text-destructive animate-in fade-in slide-in-from-top-1 duration-200" 
      in:fade={{ duration: 150 }}
    >
      {error}
    </p>
  {/if}
</div>
