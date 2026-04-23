<script lang="ts">
	/**
	 * Kodia Virtual List 🐨⚡
	 * High-performance list virtualization for institutional-grade applications.
	 */

	interface Props {
		items: any[];
		itemHeight: number;
		height: string | number;
		children: import('svelte').Snippet<[any, number]>;
		class?: string;
	}

	let { 
		items, 
		itemHeight, 
		height, 
		children, 
		class: className 
	}: Props = $props();

	let scrollTop = $state(0);
	let containerHeight = $state(0);

	let startIndex = $derived(Math.floor(scrollTop / itemHeight));
	let endIndex = $derived(Math.min(
		items.length,
		Math.ceil((scrollTop + containerHeight) / itemHeight) + 1
	));

	let visibleItems = $derived(items.slice(startIndex, endIndex).map((item, i) => ({
		item,
		index: startIndex + i,
		style: `position: absolute; top: 0; left: 0; width: 100%; transform: translateY(${(startIndex + i) * itemHeight}px);`
	})));

	let totalHeight = $derived(items.length * itemHeight);

	function handleScroll(e: Event) {
		scrollTop = (e.target as HTMLElement).scrollTop;
	}
</script>

<div
	class={className}
	style="height: {typeof height === 'number' ? height + 'px' : height}; overflow-y: auto; position: relative;"
	onscroll={handleScroll}
	bind:clientHeight={containerHeight}
>
	<div style="height: {totalHeight}px; width: 100%; pointer-events: none;"></div>
	
	{#each visibleItems as { item, index, style } (index)}
		<div {style}>
			{@render children(item, index)}
		</div>
	{/each}
</div>
