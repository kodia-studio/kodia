<svelte:head>
	<script>
		// Initialize theme before page renders to prevent flash
		const theme = localStorage.getItem('theme') || 'system';
		if (theme === 'dark' || (theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	</script>
</svelte:head>

<script lang="ts">
	import { onMount } from 'svelte';
	import { goto, afterNavigate } from '$app/navigation';
	import { Check } from 'lucide-svelte';
	import { authStore } from '$lib/stores/auth.store';
	import { themeStore } from '$lib/stores/theme.store';
	import Sidebar from '$lib/components/ui/Sidebar.svelte';
	import { Bell, LogOut, Settings, Search, Moon, Sun, Menu, X } from 'lucide-svelte';
	import { api } from '$lib/api/client.svelte';
	import { fade } from 'svelte/transition';
	import { websocketStore } from '$lib/stores/websocket.store.svelte';
	import {
		setupNotificationListener,
		teardownNotificationListener
	} from '$lib/websocket/listeners/notification-listener';

	let { children } = $props();
	let sidebarCollapsed = $state(false);
	let sidebarOpen = $state(false);
	let accountDropdownOpen = $state(false);
	let notificationDropdownOpen = $state(false);
	let notifications = $state<any[]>([]);
	let searchQuery = $state('');
	let isMobile = $state(false);

	let accountDropdownEl = $state<HTMLElement | null>(null);
	let notificationDropdownEl = $state<HTMLElement | null>(null);

	let unreadCount = $derived($authStore.notificationUnreadCount ?? 0);

	async function fetchUnreadCount() {
		try {
			const res = await api.get<any>('/api/notifications/unread/count');
			authStore.update((s) => ({ ...s, notificationUnreadCount: res?.count ?? 0 }));
		} catch {
			// Silently fail
		}
	}

	async function fetchNotifications() {
		try {
			const res = await api.get<any>('/api/notifications?per_page=5');
			notifications = res?.notifications || [];
		} catch {
			notifications = [];
		}
	}

	onMount(() => {
		if (!$authStore.isLoading && !$authStore.isAuthenticated) {
			goto('/login');
			return;
		}

		// Connect WebSocket once for the entire admin session
		websocketStore.connect();
		setupNotificationListener();
		fetchUnreadCount();
		fetchNotifications();

		// Detect mobile screen size
		const checkMobile = () => {
			isMobile = window.innerWidth < 768;
			if (isMobile) {
				sidebarOpen = false;
			}
		};
		checkMobile();
		window.addEventListener('resize', checkMobile);

		// Fallback polling every 30 seconds
		const interval = setInterval(() => {
			fetchUnreadCount();
			if (notificationDropdownOpen) {
				fetchNotifications();
			}
		}, 30000);

		// Close dropdowns when clicking outside
		const handleClickOutside = (event: MouseEvent) => {
			const target = event.target as Node;
			if (accountDropdownOpen && accountDropdownEl && !accountDropdownEl.contains(target)) {
				accountDropdownOpen = false;
			}
			if (
				notificationDropdownOpen &&
				notificationDropdownEl &&
				!notificationDropdownEl.contains(target)
			) {
				notificationDropdownOpen = false;
			}
		};

		document.addEventListener('mousedown', handleClickOutside);

		return () => {
			clearInterval(interval);
			teardownNotificationListener();
			websocketStore.disconnect();
			window.removeEventListener('resize', checkMobile);
			document.removeEventListener('mousedown', handleClickOutside);
		};
	});

	// Close dropdowns on navigation
	afterNavigate(() => {
		accountDropdownOpen = false;
		notificationDropdownOpen = false;
	});

	async function handleNotificationClick(notifId: string, isRead: boolean) {
		if (!isRead) {
			try {
				await api.post(`/api/notifications/${notifId}/read`, {});
				authStore.update((s) => ({
					...s,
					notificationUnreadCount: Math.max(0, (s.notificationUnreadCount ?? 0) - 1)
				}));
			} catch (err) {
				console.error('Failed to mark notification as read:', err);
			}
		}
		notificationDropdownOpen = false;
		goto(`/notifications/${notifId}`);
	}

	function handleLogout() {
		authStore.logout();
		goto('/login');
	}
</script>

{#if $authStore.isAuthenticated}
	<div class="flex h-screen overflow-hidden">
		<!-- Sidebar - Hidden on mobile, visible on md+ screens -->
		<div class="hidden md:block">
			<Sidebar bind:collapsed={sidebarCollapsed} />
		</div>

		<!-- Mobile Sidebar Overlay -->
		{#if sidebarOpen && isMobile}
			<button
				type="button"
				class="fixed inset-0 z-40 bg-black/50"
				onclick={() => (sidebarOpen = false)}
				aria-label="Close sidebar menu"
			></button>
			<div class="fixed inset-y-0 left-0 z-40 w-64 md:hidden">
				<Sidebar bind:collapsed={sidebarCollapsed} onNavigate={() => (sidebarOpen = false)} />
			</div>
		{/if}

		<!-- Main Content -->
		<main
			class="relative flex flex-1 flex-col overflow-y-auto bg-slate-50/50 transition-all duration-300 dark:bg-slate-950/50"
			style="margin-left: {isMobile ? '0' : sidebarCollapsed ? '80px' : '256px'}"
		>
			<!-- Header -->
			<header
				class="fixed top-0 z-30 flex h-16 items-center justify-between border-b border-slate-200/50 bg-white/60 backdrop-blur-3xl transition-all duration-300 dark:border-slate-800/50 dark:bg-slate-900/60 md:h-20"
				style="right: 0; left: {isMobile ? '0' : sidebarCollapsed ? '80px' : '256px'}; padding: 0 {isMobile ? '1rem' : '2rem'}"
			>
				<!-- Mobile Menu Button -->
				<button
					type="button"
					onclick={() => (sidebarOpen = !sidebarOpen)}
					class="md:hidden p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
					aria-label="Toggle sidebar menu"
				>
					{#if sidebarOpen}
						<X class="w-5 h-5 text-slate-600 dark:text-slate-400" />
					{:else}
						<Menu class="w-5 h-5 text-slate-600 dark:text-slate-400" />
					{/if}
				</button>

				<!-- Search Box -->
				<div class="group relative hidden flex-1 mx-4 md:block max-w-md lg:max-w-lg">
					<Search
						class="group-focus-within:text-primary absolute top-1/2 left-4 h-4 w-4 -translate-y-1/2 text-slate-400 transition-colors"
					/>
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search..."
						class="focus:ring-primary/10 focus:border-primary/30 w-full rounded-xl border border-slate-200/50 bg-slate-100/50 py-2.5 pr-4 pl-11 text-sm font-medium transition-all outline-none placeholder:text-slate-400 focus:ring-4 dark:border-slate-700/50 dark:bg-slate-800/50"
					/>
				</div>

				<!-- Right Actions -->
				<div class="flex items-center gap-2 md:gap-5">
					<!-- Notification Bell -->
					<div class="relative" bind:this={notificationDropdownEl}>
						<button
							onclick={async () => {
								notificationDropdownOpen = !notificationDropdownOpen;
								if (notificationDropdownOpen) {
									await fetchNotifications();
								}
							}}
							class="group relative rounded-lg p-2 md:p-2.5 md:rounded-xl transition-all hover:bg-slate-100 active:scale-95 dark:hover:bg-slate-800"
							aria-label="Notifications"
						>
							<Bell class="group-hover:text-primary h-4 w-4 md:h-5 md:w-5 text-slate-500 transition-colors" />
							{#if unreadCount > 0}
								<span
									class="absolute top-1 right-1 md:top-1.5 md:right-1.5 flex h-4 md:h-4.5 min-w-4 md:min-w-4.5 items-center justify-center rounded-full border-2 border-white bg-red-500 px-1 text-[8px] md:text-[10px] leading-none font-black text-white shadow-sm dark:border-slate-900"
								>
									{unreadCount > 9 ? '9+' : unreadCount}
								</span>
							{:else}
								<span
									class="bg-primary absolute top-2 right-2 md:top-2.5 md:right-2.5 h-2 w-2 animate-pulse rounded-full border-2 border-white shadow-sm dark:border-slate-900"
								></span>
							{/if}
						</button>

						<!-- Notifications Dropdown -->
						{#if notificationDropdownOpen}
							<div
								transition:fade={{ duration: 150 }}
								class="fixed top-16 left-4 right-4 z-50 overflow-hidden rounded-xl border border-slate-200 bg-white/90 shadow-2xl ring-1 ring-black/5 backdrop-blur-2xl dark:border-white/10 dark:bg-slate-900/90 md:absolute md:top-auto md:left-auto md:right-0 md:mt-2 md:w-96 md:rounded-2xl"
							>
								<!-- Header -->
								<div
									class="border-b border-slate-100 bg-slate-50/50 px-4 py-3 dark:border-white/5 dark:bg-white/5"
								>
									<p class="text-[10px] font-black tracking-[0.2em] text-slate-400 uppercase">
										Recent Notifications
									</p>
								</div>

								<!-- Notifications List -->
								<div class="custom-scrollbar max-h-96 overflow-y-auto">
									{#if notifications.length > 0}
										{#each notifications as notif (notif.id)}
											<button
												type="button"
												onclick={() => handleNotificationClick(notif.id, notif.is_read)}
												class={`w-full border-b border-slate-100 px-4 py-3 text-left transition-colors dark:border-white/5 ${notif.is_read ? 'hover:bg-slate-50 dark:hover:bg-slate-800/50' : 'bg-blue-50/30 hover:bg-blue-50/50 dark:bg-blue-900/20 dark:hover:bg-blue-900/30'}`}
											>
												<div class="flex items-start gap-3">
													<div
														class="mt-0.5 h-2 w-2 rounded-full {notif.is_read
															? 'bg-slate-300'
															: 'bg-blue-500'} shrink-0"
													></div>
													<div class="min-w-0 flex-1">
														<p
															class="line-clamp-1 text-sm font-semibold text-slate-900 dark:text-white"
														>
															{notif.title}
														</p>
														<p
															class="mt-0.5 line-clamp-2 text-xs text-slate-600 dark:text-slate-400"
														>
															{notif.message}
														</p>
														<p class="mt-1 text-[10px] text-slate-500 dark:text-slate-500">
															{new Date(notif.created_at).toLocaleString()}
														</p>
													</div>
													{#if !notif.is_read}
														<Check class="mt-0.5 h-4 w-4 shrink-0 text-blue-500" />
													{/if}
												</div>
											</button>
										{/each}
									{:else}
										<div class="px-4 py-8 text-center">
											<p class="text-sm text-slate-500 dark:text-slate-400">No notifications yet</p>
										</div>
									{/if}
								</div>

								<!-- Footer -->
								<a
									href="/notifications"
									onclick={() => (notificationDropdownOpen = false)}
									class="text-primary block border-t border-slate-100 px-4 py-3 text-center text-sm font-semibold transition-colors hover:bg-slate-100 dark:border-white/5 dark:hover:bg-slate-800/50"
								>
									View all notifications
								</a>
							</div>
						{/if}
					</div>

					<!-- Theme Toggle -->
					<button
						onclick={() => themeStore.set($themeStore === 'dark' ? 'light' : 'dark')}
						class="group rounded-lg p-2 md:p-2.5 md:rounded-xl transition-all hover:bg-slate-100 active:scale-95 dark:hover:bg-slate-800"
					>
						{#if $themeStore === 'dark'}
							<Sun class="h-4 w-4 md:h-5 md:w-5 text-yellow-500 transition-colors group-hover:text-yellow-600" />
						{:else}
							<Moon class="h-4 w-4 md:h-5 md:w-5 text-slate-500 transition-colors group-hover:text-slate-700" />
						{/if}
					</button>

					<div class="hidden h-6 w-px bg-slate-200 dark:bg-slate-800 sm:block"></div>

					<!-- Account Dropdown -->
					<div class="relative" bind:this={accountDropdownEl}>
						<button
							onclick={() => (accountDropdownOpen = !accountDropdownOpen)}
							aria-label="Account menu"
							class="flex items-center gap-2 md:gap-3 rounded-lg md:rounded-xl border border-transparent p-1.5 md:pr-3 transition-all hover:border-slate-200 hover:bg-slate-100 active:scale-95 dark:hover:border-white/5 dark:hover:bg-slate-800"
						>
							<div
								class="from-primary/20 to-secondary/20 text-primary border-primary/20 ring-primary/10 flex h-8 w-8 md:h-9 md:w-9 items-center justify-center overflow-hidden rounded-lg border bg-linear-to-tr shadow-inner ring-1"
							>
								{#if $authStore.user?.avatar_url}
									<img
										src={$authStore.user.avatar_url}
										alt="Avatar"
										class="h-full w-full object-cover"
									/>
								{:else}
									<span class="text-primary text-xs font-black uppercase">
										{$authStore.user?.name?.charAt(0).toUpperCase() || 'U'}
									</span>
								{/if}
							</div>
							<div class="hidden text-left md:block">
								<p
									class="text-xs leading-tight font-black tracking-tight text-slate-900 uppercase dark:text-white"
								>
									{$authStore.user?.name || 'Member'}
								</p>
								<div class="flex items-center gap-1.5">
									<div class="h-1.5 w-1.5 rounded-full bg-emerald-500"></div>
									<p
										class="text-[10px] font-black tracking-widest text-emerald-600 uppercase dark:text-emerald-400"
									>
										Online
									</p>
								</div>
							</div>
						</button>

						<!-- Dropdown Menu -->
						{#if accountDropdownOpen}
							<div
								transition:fade={{ duration: 150 }}
								class="absolute right-0 z-50 mt-2 w-48 sm:w-56 overflow-hidden rounded-xl md:rounded-2xl border border-slate-200 bg-white/90 py-2 shadow-2xl ring-1 ring-black/5 backdrop-blur-2xl dark:border-white/10 dark:bg-slate-900/90"
							>
								<div
									class="mb-1 border-b border-slate-100 bg-slate-50/50 px-4 py-3 dark:border-white/5 dark:bg-white/5"
								>
									<p class="text-[10px] font-black tracking-[0.2em] text-slate-400 uppercase">
										Settings & Profile
									</p>
								</div>
								<a
									href="/settings"
									class="hover:bg-primary/5 hover:text-primary flex items-center gap-3 px-4 py-3 text-sm font-bold text-slate-600 transition-all dark:text-slate-400"
									onclick={() => (accountDropdownOpen = false)}
								>
									<Settings class="h-4 w-4" />
									Settings
								</a>
								<button
									onclick={handleLogout}
									class="flex w-full items-center gap-3 px-4 py-3 text-sm font-bold text-rose-500 transition-all hover:bg-rose-500/5"
								>
									<LogOut class="h-4 w-4" />
									Sign Out
								</button>
							</div>
						{/if}
					</div>
				</div>
			</header>

			<div class="flex-1 px-4 pb-4 md:px-6 md:pb-6 lg:px-10 lg:pb-10 pt-20 md:pt-28">
				{@render children()}
			</div>

			<!-- Footer -->
			<footer
				class="border-t border-slate-200/50 bg-white/30 px-4 md:px-6 lg:px-10 py-6 md:py-8 text-[10px] md:text-xs font-black text-slate-400 backdrop-blur-sm dark:border-white/5 dark:bg-transparent dark:text-slate-600"
			>
				<div class="flex flex-col sm:flex-row sm:items-center gap-3 sm:gap-4 text-center sm:text-left">
					<span>&copy; {new Date().getFullYear()} Kodia. All rights reserved.</span>
					<div class="hidden sm:block h-1.5 w-1.5 rounded-full bg-slate-200 dark:bg-slate-800"></div>
					<span>Administration Platform</span>
				</div>
			</footer>
		</main>
	</div>
{:else if !$authStore.isLoading}
	<!-- Not authenticated, logic in onMount will redirect -->
{/if}
