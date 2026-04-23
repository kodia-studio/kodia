<script lang="ts">
	import { onMount } from 'svelte';
	import { Loader2 } from 'lucide-svelte';

	interface Props {
		items: any[];
		hasMore: boolean;
		loading: boolean;
		loadMore: () => Promise<void>;
		children: import('svelte').Snippet<[any]>;
		loader?: import('svelte').Snippet;
		empty?: import('svelte').Snippet;
		threshold?: number;
		class?: string;
	}

	let { 
		items, 
		hasMore, 
		loading, 
		loadMore, 
		children, 
		loader, 
		empty,
		threshold = 0.5,
		class: className 
	}: Props = $props();

	let observer: IntersectionObserver;
	let sentinel: HTMLElement;

	onMount(() => {
		observer = new IntersectionObserver((entries) => {
			if (entries[0].isIntersecting && hasMore && !loading) {
				loadMore();
			}
		}, { threshold });

		if (sentinel) observer.observe(sentinel);

		return () => observer.disconnect();
	});

	// Re-observe if sentinel changes
	$effect(() => {
		if (sentinel && observer) {
			observer.observe(sentinel);
		}
	});
</script>

<div class={className}>
	{#each items as item}
		{@render children(item)}
	{:else}
		{#if !loading && empty}
			{@render empty()}
		{/if}
	{/each}

	<div 
		bind:this={sentinel} 
		class="h-10 flex items-center justify-center w-full"
	>
		{#if loading}
			{#if loader}
				{@render loader()}
			{:else}
				<div class="flex items-center gap-2 text-primary font-black uppercase tracking-widest text-[10px]">
					<Loader2 class="w-4 h-4 animate-spin" />
					Loading Intelligence...
				</div>
			{/if}
		{/if}
	</div>
</div>
