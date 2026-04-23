<script lang="ts">
	import { DropdownMenu as BitsDropdown } from "bits-ui";
	import { cn } from "$lib/utils";
	import { fly } from "svelte/transition";

	interface Item {
		label: string;
		onClick?: () => void;
		icon?: any;
		variant?: "default" | "danger";
		disabled?: boolean;
	}

	interface Props {
		items?: Item[];
		children?: import('svelte').Snippet;
		class?: string;
	}

	let { items = [], children, class: className }: Props = $props();
</script>

<BitsDropdown.Root>
	<BitsDropdown.Trigger class="outline-none">
		{@render children?.()}
	</BitsDropdown.Trigger>
	<BitsDropdown.Content
		class={cn(
			"z-50 min-w-[12rem] rounded-kodia glass border p-1 shadow-kodia-lg outline-none",
			className
		)}
	>
		{#snippet child({ props, open })}
			{#if open}
				<div 
					{...props} 
					transition:fly={{ y: 8, duration: 200 }}
				>
					{#each items as item}
						<BitsDropdown.Item
							disabled={item.disabled}
							onSelect={() => item.onClick?.()}
							class={cn(
								"flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm font-bold outline-none transition-colors cursor-pointer",
								item.variant === "danger"
									? "text-red-500 hover:bg-red-500/10"
									: "text-slate-600 dark:text-slate-300 hover:bg-primary/10 hover:text-primary",
								item.disabled && "opacity-50 cursor-not-allowed"
							)}
						>
							{#if item.icon}
								<item.icon class="w-4 h-4" />
							{/if}
							{item.label}
						</BitsDropdown.Item>
					{/each}
				</div>
			{/if}
		{/snippet}
	</BitsDropdown.Content>
</BitsDropdown.Root>
