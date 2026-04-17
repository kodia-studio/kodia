<script lang="ts">
	import { page } from '$app/stores';
	import { themeStore } from '$lib/stores/theme.store';
	import '../app.css';

	// The $page.status will give us the HTTP error code (e.g., 404, 500)
	// The $page.error will contain the error details
</script>

<svelte:head>
	<title>Error {$page.status} - Kodia</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=Outfit:wght@100..900&display=swap" rel="stylesheet">
</svelte:head>

<div class="min-h-screen flex items-center justify-center p-4 bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 font-sans">
	<div class="text-center max-w-lg w-full">
		<div class="mb-8 relative">
			<h1 class="text-9xl font-bold text-slate-200 dark:text-slate-800 tracking-tighter mix-blend-multiply dark:mix-blend-screen opacity-50">
				{$page.status}
			</h1>
			<div class="absolute inset-0 flex items-center justify-center">
				<div class="w-20 h-20 rounded-2xl bg-primary-600 flex items-center justify-center text-white shadow-xl shadow-primary-500/30">
					<svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				</div>
			</div>
		</div>

		<h2 class="text-2xl font-bold mb-4 tracking-tight">
			{#if $page.status === 404}
				Page not found
			{:else}
				Something went wrong
			{/if}
		</h2>
		
		<p class="text-slate-500 dark:text-slate-400 mb-8 max-w-sm mx-auto">
			{$page.error?.message || "We couldn't find the page you're looking for or an unexpected error occurred."}
		</p>

		<div class="flex items-center justify-center gap-4">
			<button onclick={() => window.history.back()} class="px-4 py-2 rounded-lg font-medium border border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors">
				Go back
			</button>
			<a href="/" class="btn-primary">
				Return home
			</a>
		</div>

		{#if process.env.NODE_ENV === 'development' && $page.error?.stack}
			<div class="mt-12 text-left bg-slate-100 dark:bg-slate-900 p-4 rounded-xl overflow-x-auto border border-red-200 dark:border-red-900/30">
				<p class="text-xs font-bold text-red-500 mb-2 uppercase tracking-wider">Error Stack (Dev Only)</p>
				<pre class="text-xs text-slate-600 dark:text-slate-400 font-mono">{$page.error.stack}</pre>
			</div>
		{/if}
	</div>
</div>
