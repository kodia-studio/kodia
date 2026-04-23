<script lang="ts">
	import { cn } from "$lib/utils";

	interface Props {
		label?: string;
		error?: string;
		id?: string;
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { label, error, id, class: className, children }: Props = $props();
</script>

<div class={cn("flex flex-col gap-1.5", className)}>
	{#if label}
		<label 
			for={id}
			class="text-xs font-black uppercase tracking-widest text-slate-500 ml-1"
		>
			{label}
		</label>
	{/if}
	
	{@render children?.()}

	{#if error}
		<span class="text-[10px] font-bold text-red-500 ml-1 uppercase tracking-tight">
			{error}
		</span>
	{/if}
</div>
