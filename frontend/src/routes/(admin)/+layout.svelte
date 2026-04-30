<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth.store';
	import { themeStore } from '$lib/stores/theme.store';
	import Sidebar from '$lib/components/ui/Sidebar.svelte';
	import { Bell, LogOut, Settings, Search, Moon, Sun } from 'lucide-svelte';
	import { fade } from 'svelte/transition';

	let { children } = $props();
	let sidebarCollapsed = $state(false);
	let accountDropdownOpen = $state(false);
	let searchQuery = $state("");

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
</script>

{#if $authStore.isAuthenticated}
	<div class="flex h-screen overflow-hidden">
		<!-- Sidebar -->
		<Sidebar bind:collapsed={sidebarCollapsed} />

		<!-- Main Content -->
		<main class="flex-1 overflow-y-auto relative bg-slate-50/50 dark:bg-slate-950/50 transition-all duration-300 flex flex-col" style="margin-left: {sidebarCollapsed ? '80px' : '256px'}">
			<!-- Header -->
			<header class="fixed top-0 z-30 h-20 border-b border-slate-200/50 dark:border-slate-800/50 bg-white/60 dark:bg-slate-900/60 backdrop-blur-3xl transition-all duration-300 px-8 flex items-center justify-between" style="right: 0; left: {sidebarCollapsed ? '80px' : '256px'}">
				<!-- Search Box -->
				<div class="relative w-96 hidden md:block group">
					<Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 group-focus-within:text-primary transition-colors" />
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search for metrics, users, or logs..."
						class="w-full pl-11 pr-4 py-2.5 bg-slate-100/50 dark:bg-slate-800/50 border border-slate-200/50 dark:border-slate-700/50 rounded-xl text-sm focus:ring-4 focus:ring-primary/10 focus:border-primary/30 outline-none transition-all placeholder:text-slate-400 font-medium"
					/>
				</div>

				<!-- Right Actions -->
				<div class="flex items-center gap-5">
					<!-- Notification Bell -->
					<button class="p-2.5 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl relative transition-all active:scale-95 group">
						<Bell class="w-5 h-5 text-slate-500 group-hover:text-primary transition-colors" />
						<span class="absolute top-2.5 right-2.5 w-2 h-2 bg-primary rounded-full border-2 border-white dark:border-slate-900 shadow-sm animate-pulse"></span>
					</button>

					<!-- Theme Toggle -->
					<button
						onclick={() => themeStore.set($themeStore === 'dark' ? 'light' : 'dark')}
						class="p-2.5 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition-all active:scale-95 group"
					>
						{#if $themeStore === 'dark'}
							<Sun class="w-5 h-5 text-yellow-500 group-hover:text-yellow-600 transition-colors" />
						{:else}
							<Moon class="w-5 h-5 text-slate-500 group-hover:text-slate-700 transition-colors" />
						{/if}
					</button>

					<div class="h-6 w-px bg-slate-200 dark:bg-slate-800"></div>

					<!-- Account Dropdown -->
					<div class="relative">
						<button
							onclick={() => (accountDropdownOpen = !accountDropdownOpen)}
							class="flex items-center gap-3 p-1.5 pr-3 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition-all active:scale-95 border border-transparent hover:border-slate-200 dark:hover:border-white/5"
						>
							<div class="w-9 h-9 rounded-lg bg-linear-to-tr from-primary/20 to-secondary/20 flex items-center justify-center text-primary shadow-inner border border-primary/20 ring-1 ring-primary/10 overflow-hidden">
								{#if $authStore.user?.avatar_url}
									<img src={$authStore.user.avatar_url} alt="Avatar" class="w-full h-full object-cover" />
								{:else}
									<span class="text-xs font-black text-primary uppercase">
										{$authStore.user?.name?.charAt(0).toUpperCase() || 'U'}
									</span>
								{/if}
							</div>
							<div class="text-left hidden sm:block">
								<p class="text-xs font-black tracking-tight leading-tight text-slate-900 dark:text-white uppercase">
									{$authStore.user?.name || "Member"}
								</p>
								<div class="flex items-center gap-1.5">
									<div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
									<p class="text-[10px] text-emerald-600 dark:text-emerald-400 font-black uppercase tracking-widest">
										Online
									</p>
								</div>
							</div>
						</button>

						<!-- Dropdown Menu -->
						{#if accountDropdownOpen}
							<div
								transition:fade={{ duration: 150 }}
								class="absolute right-0 mt-3 w-56 bg-white/90 dark:bg-slate-900/90 backdrop-blur-2xl border border-slate-200 dark:border-white/10 rounded-2xl shadow-2xl py-2 z-50 overflow-hidden ring-1 ring-black/5"
							>
								<div class="px-4 py-3 border-b border-slate-100 dark:border-white/5 mb-1 bg-slate-50/50 dark:bg-white/5">
									<p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Settings & Profile</p>
								</div>
								<a href="/settings" class="flex items-center gap-3 px-4 py-3 text-sm font-bold text-slate-600 dark:text-slate-400 hover:bg-primary/5 hover:text-primary transition-all" onclick={() => (accountDropdownOpen = false)}>
									<Settings class="w-4 h-4" />
									Console Settings
								</a>
								<button
									onclick={handleLogout}
									class="w-full flex items-center gap-3 px-4 py-3 text-sm font-bold text-rose-500 hover:bg-rose-500/5 transition-all"
								>
									<LogOut class="w-4 h-4" />
									Sign Out.
								</button>
							</div>
						{/if}
					</div>
				</div>
			</header>

			<div class="flex-1 p-10 pt-24">
				{@render children()}
			</div>

			<!-- Footer -->
			<footer class="px-10 py-8 border-t border-slate-200/50 dark:border-white/5 text-[10px] font-black uppercase tracking-[0.3em] text-slate-400 dark:text-slate-600 bg-white/30 dark:bg-transparent backdrop-blur-sm">
				<div class="flex items-center gap-4">
					<span>&copy; {new Date().getFullYear()} Kodia Elite Console.</span>
					<div class="w-1.5 h-1.5 rounded-full bg-slate-200 dark:bg-slate-800"></div>
					<span>High-Fidelity Framework</span>
				</div>
			</footer>
		</main>
	</div>
{:else if !$authStore.isLoading}
	<!-- Not authenticated, logic in onMount will redirect -->
{/if}
