<script lang="ts">
	import { cn } from "$lib/utils";

	interface Props {
		src?: string;
		alt?: string;
		fallback?: string;
		size?: "sm" | "md" | "lg" | "xl";
		class?: string;
	}

	let { src, alt = "Avatar", fallback = "K", size = "md", class: className }: Props = $props();

	const sizes = {
		sm: "w-8 h-8 text-[10px]",
		md: "w-10 h-10 text-xs",
		lg: "w-16 h-16 text-lg",
		xl: "w-24 h-24 text-2xl"
	};
</script>

<div
	class={cn(
		"relative inline-flex items-center justify-center rounded-full overflow-hidden border-2 border-white dark:border-slate-800 shadow-kodia bg-slate-100 dark:bg-slate-800",
		sizes[size],
		className
	)}
>
	{#if src}
		<img {src} {alt} class="w-full h-full object-cover" />
	{:else}
		<span class="font-black text-slate-500 dark:text-slate-400 uppercase">
			{fallback}
		</span>
	{/if}
</div>
