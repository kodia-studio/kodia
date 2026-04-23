<script lang="ts">
	/**
	 * Kodia Skeleton Loader 🐨🦴
	 * High-fidelity shimmer animations for smooth content transitions.
	 */
	import { cn } from '$lib/utils';

	interface Props {
		variant?: 'text' | 'circle' | 'rect' | 'avatar';
		width?: string | number;
		height?: string | number;
		class?: string;
	}

	let { 
		variant = 'rect', 
		width, 
		height, 
		class: className 
	}: Props = $props();

	const style = $derived(`
		${width ? `width: ${typeof width === 'number' ? width + 'px' : width};` : ''}
		${height ? `height: ${typeof height === 'number' ? height + 'px' : height};` : ''}
	`);
</script>

<div
	class={cn(
		"relative overflow-hidden bg-slate-200 dark:bg-slate-800 animate-pulse",
		variant === 'circle' && "rounded-full",
		variant === 'avatar' && "rounded-2xl w-12 h-12",
		variant === 'text' && "rounded-md h-3 w-full mb-2",
		variant === 'rect' && "rounded-2xl",
		className
	)}
	{style}
>
	<!-- Shimmer Effect -->
	<div class="absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/10 dark:via-white/5 to-transparent shimmer-animation"></div>
</div>

<style>
	@keyframes shimmer {
		100% { transform: translateX(100%); }
	}
	.shimmer-animation {
		animation: shimmer 2s infinite;
	}
</style>
