<script lang="ts">
	import { mapKodiaErrors, isValidationError } from '../../utils/errors';
	import type { SuperForm } from 'sveltekit-superforms';
	import { toast } from 'svelte-sonner';

	interface Props {
		form: SuperForm<any>;
		onsuccess?: (data: any) => void;
		onerror?: (error: any) => void;
		class?: string;
		children?: import('svelte').Snippet;
	}

	let { form, onsuccess, onerror, class: className, children }: Props = $props();

	const enhance = $derived(form.enhance);
	const errors = $derived(form.errors);
	const message = $derived(form.message);
	const tainted = $derived(form.tainted);
	const submitting = $derived(form.submitting);

	// Enhanced form handling
	function handleResult({ result }: any) {
		if (result.type === 'success') {
			if (onsuccess) onsuccess(result.data);
			toast.success('Operation successful');
		} else if (result.type === 'failure' || result.type === 'error') {
			const errorData = result.data || result.error;
			
			if (isValidationError(errorData)) {
				// Map backend errors to Superforms
				const mapped = mapKodiaErrors(errorData);
				errors.set(mapped as any);
				toast.error('Please correct the errors below');
			} else {
				if (onerror) onerror(errorData);
				toast.error(errorData?.message || 'An unexpected error occurred');
			}
		}
	}
</script>

<form 
	method="POST" 
	use:enhance={{ 
		onUpdate: handleResult 
	}}
	class={className}
>
	{#if $message}
		<div class="p-4 mb-6 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-500 text-sm font-medium">
			{$message}
		</div>
	{/if}

	{@render children?.()}
</form>
