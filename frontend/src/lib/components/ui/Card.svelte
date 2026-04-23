<script lang="ts">
	/**
	 * Kodia Card 🐨🎴
	 * Institutional-grade container with premium aesthetics.
	 */
	import { cn } from '$lib/utils';

	interface Props {
		children: import('svelte').Snippet;
		class?: string;
		padding?: 'none' | 'sm' | 'md' | 'lg';
		hover?: boolean;
	}

	let { 
		children, 
		class: className,
		padding = 'md',
		hover = false
	}: Props = $props();

	const paddings = {
		none: 'p-0',
		sm: 'p-4',
		md: 'p-6',
		lg: 'p-8'
	};
</script>

<div class={cn(
	"glass rounded-[2rem] border border-slate-200 dark:border-slate-800 overflow-hidden transition-all duration-500",
	hover && "hover:shadow-kodia-xl hover:shadow-primary/5 hover:-translate-y-1",
	paddings[padding],
	className
)}>
	{@render children()}
</div>
