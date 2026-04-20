<script lang="ts">
  import { cn } from "$lib/utils/styles";
  import type { HTMLAttributes } from "svelte/elements";

  interface Props extends HTMLAttributes<HTMLDivElement> {
    fluid?: boolean;
    size?: "sm" | "md" | "lg" | "xl" | "2xl";
  }

  let { fluid = false, size = "xl", class: className, children, ...props }: Props = $props();

  const sizes = {
    sm: "max-w-screen-sm",
    md: "max-w-screen-md",
    lg: "max-w-screen-lg",
    xl: "max-w-screen-xl",
    "2xl": "max-w-screen-2xl",
  };
</script>

<div
  class={cn(
    "mx-auto px-4 sm:px-6 lg:px-8 w-full",
    !fluid && sizes[size],
    className
  )}
  {...props}
>
  {@render children?.()}
</div>
