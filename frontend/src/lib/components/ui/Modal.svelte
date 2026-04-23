<script lang="ts">
	import { Dialog as BitsDialog } from "bits-ui";
	import { cn } from "$lib/utils";
	import { fade, scale } from "svelte/transition";
	import { X } from "lucide-svelte";
	import { backOut } from "svelte/easing";

	interface Props {
		open?: boolean;
		title?: string;
		description?: string;
		children?: import('svelte').Snippet;
		footer?: import('svelte').Snippet;
		size?: "sm" | "md" | "lg" | "xl" | "full";
		class?: string;
	}

	let {
		open = $bindable(false),
		title,
		description,
		children,
		footer,
		size = "md",
		class: className
	}: Props = $props();

	const sizes = {
		sm: "max-w-sm",
		md: "max-w-md",
		lg: "max-w-2xl",
		xl: "max-w-4xl",
		full: "max-w-[95vw]"
	};
</script>

<BitsDialog.Root bind:open>
	<BitsDialog.Portal>
		<BitsDialog.Overlay
			class="fixed inset-0 z-50 bg-slate-950/40 backdrop-blur-md"
		>
			{#snippet child({ props, open })}
				{#if open}
					<div {...props} transition:fade={{ duration: 300 }}></div>
				{/if}
			{/snippet}
		</BitsDialog.Overlay>
		<BitsDialog.Content
			class={cn(
				"fixed left-[50%] top-[50%] z-50 w-full translate-x-[-50%] translate-y-[-50%] p-0 outline-none",
				sizes[size]
			)}
		>
			{#snippet child({ props, open })}
				{#if open}
					<div 
						{...props} 
						transition:scale={{ duration: 400, start: 0.95, opacity: 0, easing: backOut }}
					>
						<div class={cn("card-premium m-4 flex flex-col gap-4 shadow-2xl", className)}>
							<div class="flex items-start justify-between">
								<div class="flex flex-col gap-1">
									{#if title}
										<BitsDialog.Title class="text-xl font-black tracking-tight text-slate-900 dark:text-white">
											{title}
										</BitsDialog.Title>
									{/if}
									{#if description}
										<BitsDialog.Description class="text-sm text-slate-500 dark:text-slate-400">
											{description}
										</BitsDialog.Description>
									{/if}
								</div>
								<BitsDialog.Close
									class="rounded-full p-2 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
								>
									<X class="w-4 h-4 text-slate-500" />
								</BitsDialog.Close>
							</div>

							<div class="py-2">
								{@render children?.()}
							</div>

							{#if footer}
								<div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-100 dark:border-slate-800">
									{@render footer()}
								</div>
							{/if}
						</div>
					</div>
				{/if}
			{/snippet}
		</BitsDialog.Content>
	</BitsDialog.Portal>
</BitsDialog.Root>
