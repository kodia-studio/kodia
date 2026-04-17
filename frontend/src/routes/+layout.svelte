<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth.store';
	import { themeStore } from '$lib/stores/theme.store';
	import { toastStore } from '$lib/stores/toast.store';

	let { children } = $props();

	onMount(() => {
		authStore.init();
		themeStore.init();
	});
</script>

<svelte:head>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=Outfit:wght@100..900&display=swap" rel="stylesheet">
</svelte:head>

<div class="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 font-sans">
	{@render children()}
</div>

<!-- Global Toasts -->
<div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
	{#each $toastStore as toast (toast.id)}
		<div
			class="glass px-6 py-3 rounded-xl shadow-lg flex items-center gap-3 animate-in slide-in-from-right fade-in duration-300"
			class:border-green-500={toast.type === 'success'}
			class:border-red-500={toast.type === 'error'}
		>
			{#if toast.type === 'success'}
				<div class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
			{:else if toast.type === 'error'}
				<div class="w-2 h-2 rounded-full bg-red-500 animate-pulse"></div>
			{/if}
			<p class="text-sm font-medium">{toast.message}</p>
			<button
				class="ml-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
				onclick={() => toastStore.remove(toast.id)}
			>
				×
			</button>
		</div>
	{/each}
</div>

<style>
	:global(.animate-in) {
		animation: enter 0.3s ease-out;
	}

	@keyframes enter {
		from {
			opacity: 0;
			transform: translateX(20px);
		}
		to {
			opacity: 1;
			transform: translateX(0);
		}
	}
</style>
