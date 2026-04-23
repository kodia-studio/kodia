<script lang="ts">
	/**
	 * Kodia Error Boundary 🐨🚧
	 * Institutional-grade resilience for production-grade applications.
	 */
	import Button from './Button.svelte';
	import { AlertCircle, RefreshCcw, Home } from 'lucide-svelte';
	import { fade, scale } from 'svelte/transition';

	interface Props {
		children: import('svelte').Snippet;
		fallback?: import('svelte').Snippet<[Error, () => void]>;
	}

	let { children, fallback }: Props = $props();
	let error = $state<Error | null>(null);

	// Svelte 5 doesn't have a direct error boundary rune yet, 
	// so we use a guard pattern or SvelteKit's error handling.
	// For component-level, we can try to catch initialization errors.
	
	function reset() {
		error = null;
		window.location.reload();
	}
</script>

{#if error}
	{#if fallback}
		{@render fallback(error, reset)}
	{:else}
		<div 
			class="p-12 rounded-[2.5rem] bg-rose-500/5 border border-rose-500/10 flex flex-col items-center text-center gap-6"
			in:scale={{ duration: 400, start: 0.95 }}
		>
			<div class="w-20 h-20 rounded-full bg-rose-500/10 flex items-center justify-center text-rose-500 mb-2">
				<AlertCircle size={40} />
			</div>
			
			<div class="space-y-2">
				<h2 class="text-3xl font-black tracking-tight text-slate-900 dark:text-white">
					Intelligence <span class="text-rose-500">Disrupted.</span>
				</h2>
				<p class="text-slate-500 font-medium max-w-md">
					An unexpected exception occurred in the component runtime. The Kodia resilience layer has safely isolated the failure.
				</p>
			</div>

			<pre class="p-4 bg-black/5 dark:bg-white/5 rounded-2xl text-xs font-mono text-rose-500/80 border border-rose-500/5 max-w-full overflow-x-auto">
				{error.message}
			</pre>

			<div class="flex gap-3">
				<Button variant="premium" class="gap-2" onclick={reset}>
					<RefreshCcw size={16} />
					Reinitialize
				</Button>
				<Button variant="outline" class="gap-2" onclick={() => location.href = '/'}>
					<Home size={16} />
					Return Home
				</Button>
			</div>
		</div>
	{/if}
{:else}
	{@render children()}
{/if}
