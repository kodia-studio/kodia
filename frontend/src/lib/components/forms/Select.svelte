<script lang="ts">
  import { Select as SelectPrimitive } from "bits-ui";
  import { Check, ChevronDown } from "lucide-svelte";
  import { cn } from "$lib/utils/styles";

  interface Item {
    value: string;
    label: string;
    disabled?: boolean;
  }

  interface Props {
    items: Item[];
    placeholder?: string;
    value?: string;
    onValueChange?: (value: string) => void;
    error?: boolean;
    class?: string;
    disabled?: boolean;
    name?: string;
  }

  let { 
    items, 
    placeholder = "Select an option", 
    value = $bindable(""), 
    onValueChange,
    error = false,
    class: className,
    disabled = false,
    name
  }: Props = $props();
</script>

<SelectPrimitive.Root
  type="single"
  bind:value
  onValueChange={(v) => { if (v) onValueChange?.(v); }}
  {disabled}
  {name}
>
  <SelectPrimitive.Trigger
    class={cn(
      "flex h-11 w-full items-center justify-between rounded-xl border bg-slate-50/50 dark:bg-slate-900/50 px-4 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary disabled:cursor-not-allowed disabled:opacity-50 transition-all duration-200",
      error ? "border-destructive focus-visible:ring-destructive/20" : "border-input hover:border-muted-foreground/30",
      className
    )}
  >
    <SelectPrimitive.Value {placeholder} class="text-sm" />
    <ChevronDown class="w-4 h-4 opacity-50" />
  </SelectPrimitive.Trigger>

  <SelectPrimitive.Content
    class="z-50 min-w-32 overflow-hidden rounded-xl border bg-popover text-popover-foreground shadow-2xl animate-in fade-in zoom-in-95 duration-200"
    sideOffset={4}
  >
    <SelectPrimitive.Viewport class="p-1">
      {#each items as item}
        <SelectPrimitive.Item
          value={item.value}
          disabled={item.disabled}
          class="relative flex w-full cursor-default select-none items-center rounded-lg py-2 px-3 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-disabled:pointer-events-none data-disabled:opacity-50 transition-colors"
        >
          {#snippet children({ selected })}
            <span class="flex-1">{item.label}</span>
            {#if selected}
               <Check class="w-4 h-4 text-primary ml-auto" />
            {/if}
          {/snippet}
        </SelectPrimitive.Item>
      {/each}
    </SelectPrimitive.Viewport>
  </SelectPrimitive.Content>
</SelectPrimitive.Root>
