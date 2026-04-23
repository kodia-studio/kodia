<script lang="ts">
	import type { SuperForm } from "sveltekit-superforms";
	import { cn } from "$lib/utils";
	import Button from "../ui/Button.svelte";

	interface Props {
		form: SuperForm<any>;
		submitLabel?: string;
		loadingLabel?: string;
		children?: import('svelte').Snippet;
		class?: string;
		onsubmit?: (event: SubmitEvent) => void;
	}

	let {
		form,
		submitLabel = "Submit",
		loadingLabel = "Processing...",
		children,
		class: className,
		onsubmit
	}: Props = $props();

	const delayed = $derived(form.delayed);
</script>

<form
	method="POST"
	use:form.enhance
	class={cn("flex flex-col gap-6", className)}
	onsubmit={onsubmit}
>
	<div class="flex flex-col gap-4">
		{@render children?.()}
	</div>

	<div class="flex items-center justify-end pt-4">
		<Button type="submit" variant="premium" loading={$delayed}>
			{$delayed ? loadingLabel : submitLabel}
		</Button>
	</div>
</form>
