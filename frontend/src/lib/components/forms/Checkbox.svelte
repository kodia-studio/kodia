<script lang="ts">
  import { Checkbox as CheckboxPrimitive } from "bits-ui";
  import { Check, Minus } from "lucide-svelte";
  import { cn } from "$lib/utils/styles";

  interface Props {
    checked?: boolean | "indeterminate";
    onCheckedChange?: (checked: boolean | "indeterminate") => void;
    label?: string;
    description?: string;
    disabled?: boolean;
    class?: string;
    id?: string;
  }

  let { 
    checked = $bindable(false), 
    onCheckedChange, 
    label, 
    description,
    disabled = false,
    class: className,
    id = "checkbox-" + Math.random().toString(36).substring(2, 9)
  }: Props = $props();

  function handleCheckedChange(c: boolean) {
    checked = c;
    onCheckedChange?.(c);
  }
</script>

<div class={cn("flex items-start gap-3", className)}>
  <CheckboxPrimitive.Root
    checked={checked === true}
    indeterminate={checked === "indeterminate"}
    {onCheckedChange}
    {disabled}
    {id}
    class="peer group h-5 w-5 shrink-0 rounded-md border border-input ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground data-[state=checked]:border-primary transition-all duration-200"
  >
    {#snippet children({ checked, indeterminate })}
      <div class="flex items-center justify-center text-current">
        {#if checked}
          <Check class="h-3.5 w-3.5 stroke-3" />
        {:else if indeterminate}
          <Minus class="h-3.5 w-3.5 stroke-3" />
        {/if}
      </div>
    {/snippet}
  </CheckboxPrimitive.Root>

  {#if label}
    <div class="grid gap-1.5 leading-none">
      <label
        for={id}
        class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
      >
        {label}
      </label>
      {#if description}
        <p class="text-xs text-muted-foreground leading-relaxed">
          {description}
        </p>
      {/if}
    </div>
  {/if}
</div>
