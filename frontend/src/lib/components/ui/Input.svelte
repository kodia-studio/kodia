<script lang="ts">
	import { cn } from "$lib/utils";
	import type { HTMLInputAttributes } from "svelte/elements";

	interface Props extends HTMLInputAttributes {
		label?: string;
		error?: string;
		containerClass?: string;
	}

	let {
		label,
		error,
		value = $bindable(""),
		class: className,
		containerClass,
		...rest
	}: Props = $props();
</script>

<div class={cn("flex flex-col gap-1.5", containerClass)}>
	{#if label}
		<label class="text-xs font-black uppercase tracking-widest text-slate-500 ml-1" for={rest.id}>
			{label}
		</label>
	{/if}
	
	<div class="relative">
		<input
			class={cn(
				"input-kodia",
				error && "border-red-500 focus:border-red-500 focus:ring-red-500/20",
				className
			)}
			bind:value
			{...rest}
		/>
	</div>

	{#if error}
		<span class="text-[10px] font-bold text-red-500 ml-1 uppercase tracking-tight">
			{error}
		</span>
	{/if}
</div>
