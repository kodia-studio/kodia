<script lang="ts">
	import { cn } from "$lib/utils";
	import type { HTMLButtonAttributes } from "svelte/elements";

	interface Props extends HTMLButtonAttributes {
		variant?: "primary" | "secondary" | "outline" | "ghost" | "premium" | "danger";
		size?: "sm" | "md" | "lg" | "icon";
		loading?: boolean;
	}

	let {
		children,
		class: className,
		variant = "primary",
		size = "md",
		loading = false,
		disabled,
		...rest
	}: Props = $props();

	const variants = {
		primary: "bg-primary text-white hover:bg-primary-dark shadow-lg shadow-primary/20",
		secondary: "bg-secondary text-white hover:bg-secondary/90 shadow-lg shadow-secondary/20",
		outline: "border-2 border-slate-200 dark:border-slate-800 hover:border-primary hover:text-primary bg-transparent",
		ghost: "hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400",
		premium: "btn-premium",
		danger: "bg-red-500 text-white hover:bg-red-600 shadow-lg shadow-red-500/20"
	};

	const sizes = {
		sm: "px-3 py-1.5 text-sm",
		md: "px-5 py-2.5",
		lg: "px-8 py-3.5 text-lg",
		icon: "p-2.5"
	};
</script>

<button
	class={cn(
		"inline-flex items-center justify-center gap-2 rounded-kodia font-bold transition-all duration-300 active:scale-95 disabled:opacity-50 disabled:pointer-events-none focus:outline-none focus-visible:ring-4 focus-visible:ring-primary/20",
		variants[variant],
		sizes[size],
		className
	)}
	{disabled}
	{...rest}
>
	{#if loading}
		<svg class="animate-spin h-4 w-4 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
		</svg>
	{/if}
	{@render children?.()}
</button>
