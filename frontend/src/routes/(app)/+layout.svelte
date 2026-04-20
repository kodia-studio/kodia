<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.store';
	import { themeStore } from '$lib/stores/theme.store';

	let { children } = $props();

	onMount(() => {
		// Simple client-side auth guard
		if (!$authStore.isLoading && !$authStore.isAuthenticated) {
			goto('/login');
		}
	});

	function handleLogout() {
		authStore.logout();
		goto('/login');
	}

	const navItems = [
		{ name: 'Dashboard', href: '/dashboard', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6' },
		{ name: 'Users', href: '/users', icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z' },
		{ name: 'Settings', href: '/settings', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37a1.724 1.724 0 002.572-1.065z' }
	];
</script>

{#if $authStore.isAuthenticated}
	<div class="flex h-screen overflow-hidden">
		<!-- Sidebar -->
		<aside class="w-64 glass border-r h-full hidden md:flex flex-col">
			<div class="p-6">
				<div class="flex items-center gap-2">
					<div class="w-8 h-8 rounded-lg bg-primary-600 flex items-center justify-center text-white font-bold">K</div>
					<span class="text-xl font-bold tracking-tight">Kodia</span>
				</div>
			</div>

			<nav class="flex-1 px-4 space-y-1 mt-4">
				{#each navItems as item}
					<a
						href={item.href}
						class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors hover:bg-slate-100 dark:hover:bg-slate-800"
					>
						<svg class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={item.icon} />
						</svg>
						{item.name}
					</a>
				{/each}
			</nav>

			<div class="p-4 border-t border-slate-200 dark:border-slate-800">
				<div class="flex items-center gap-3 px-3 py-2">
					<div class="w-8 h-8 rounded-full bg-slate-200 dark:bg-slate-800 flex items-center justify-center overflow-hidden">
						{#if $authStore.user?.avatar_url}
							<img src={$authStore.user.avatar_url} alt="Avatar" class="w-full h-full object-cover">
						{:else}
							<span class="text-xs font-bold">{$authStore.user?.name.charAt(0)}</span>
						{/if}
					</div>
					<div class="flex-1 min-w-0">
						<p class="text-sm font-medium truncate">{$authStore.user?.name}</p>
						<p class="text-xs text-slate-500 truncate">{$authStore.user?.email}</p>
					</div>
				</div>
				<button
					onclick={handleLogout}
					class="w-full mt-2 flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
				>
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
					</svg>
					Log out
				</button>
			</div>
		</aside>

		<!-- Main Content -->
		<main class="flex-1 overflow-y-auto relative bg-slate-50/50 dark:bg-slate-950/50">
			<!-- Header -->
			<header class="sticky top-0 z-30 h-16 glass border-b px-8 flex items-center justify-between">
				<button class="md:hidden text-slate-500" aria-label="Toggle menu">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7" />
					</svg>
				</button>

				<div class="flex-1"></div>

				<div class="flex items-center gap-4">
					<button
						class="p-2 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
						onclick={() => themeStore.set($themeStore === 'dark' ? 'light' : 'dark')}
					>
						{#if $themeStore === 'dark'}
							<svg class="w-5 h-5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M12 5a7 7 0 100 14 7 7 0 000-14z" />
							</svg>
						{:else}
							<svg class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
							</svg>
						{/if}
					</button>
				</div>
			</header>

			<div class="p-8">
				{@render children()}
			</div>
		</main>
	</div>
{:else if !$authStore.isLoading}
	<!-- Not authenticated, logic in onMount will redirect -->
{/if}
