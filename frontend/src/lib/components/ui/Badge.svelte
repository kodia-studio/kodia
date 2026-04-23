<script lang="ts">
	import { cn } from "$lib/utils";

	interface Props {
		variant?: "primary" | "secondary" | "success" | "warning" | "danger" | "ghost" | "premium";
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { variant = "primary", class: className, children }: Props = $props();

	const variants = {
		primary: "bg-primary/10 text-primary border-primary/20",
		secondary: "bg-secondary/10 text-secondary border-secondary/20",
		success: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20",
		warning: "bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20",
		danger: "bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/20",
		ghost: "bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border-transparent",
		premium: "bg-linear-to-r from-primary/10 to-secondary/10 text-primary border-primary/20 shadow-sm"
	};
</script>

<span
	class={cn(
		"inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-widest border",
		variants[variant],
		className
	)}
>
	{@render children?.()}
</span>
