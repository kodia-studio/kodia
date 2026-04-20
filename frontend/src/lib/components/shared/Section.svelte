<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import type { HTMLAttributes } from "svelte/elements";

  interface Props extends HTMLAttributes<HTMLElement> {
    spacing?: "none" | "sm" | "md" | "lg" | "xl";
    variant?: "default" | "muted" | "glass" | "primary";
  }

  let { 
    spacing = "md", 
    variant = "default", 
    class: className, 
    children, 
    ...props 
  }: Props = $props();

  const spacings = {
    none: "py-0",
    sm: "py-8",
    md: "py-16",
    lg: "py-24",
    xl: "py-32",
  };

  const variants = {
    default: "bg-background",
    muted: "bg-muted/50",
    glass: "glass",
    primary: "bg-primary text-primary-foreground",
  };
</script>

<section
  class={cn(
    spacings[spacing],
    variants[variant],
    className
  )}
  {...props}
>
  {@render children?.()}
</section>
