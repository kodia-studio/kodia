<script lang="ts">
	import { cn } from '$lib/utils/styles';
	import { LayoutDashboard, Users, Settings, ChevronLeft } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { authStore } from '$lib/stores/auth.store';

	let { collapsed = $bindable(false), onNavigate }: { collapsed: boolean; onNavigate?: () => void } = $props();

	const allMenuItems = [
		{ name: 'Dashboard', icon: LayoutDashboard, href: '/dashboard', requiresAdmin: false },
		{ name: 'Users', icon: Users, href: '/users', requiresAdmin: true },
		{ name: 'Settings', icon: Settings, href: '/settings', requiresAdmin: false }
	];

	const menuItems = $derived.by(() => {
		const isAdmin = $authStore.user?.role === 'admin';
		return allMenuItems.filter((item) => !item.requiresAdmin || isAdmin);
	});
</script>

<aside
	class={cn(
		'fixed top-0 bottom-0 left-0 z-50 border-r transition-all duration-300',
		'flex flex-col border-slate-200/50 bg-white shadow-2xl dark:border-white/5 dark:bg-slate-900',
		collapsed ? 'w-20' : 'w-64'
	)}
>
	<!-- Sidebar Header / Brand -->
	<div
		class={cn(
			'relative flex items-center justify-center transition-all duration-300',
			collapsed ? 'h-20 px-4' : 'h-24 px-7'
		)}
	>
		<div
			class="absolute inset-x-0 bottom-0 h-px bg-linear-to-r from-transparent via-slate-200/50 to-transparent dark:via-white/5"
		></div>

		<a href="/" class="group flex items-center gap-4 overflow-hidden">
			<img
				src="/logo.png"
				alt="Kodia Logo"
				class={cn(
					'shrink-0 transition-all duration-300 group-hover:scale-105',
					collapsed ? 'h-10 w-10' : 'h-11 w-11 min-w-11'
				)}
			/>
			{#if !collapsed}
				<div class="flex flex-col">
					<span
						class="font-heading text-xl leading-none font-black tracking-tighter text-slate-900 dark:text-white"
					>
						Kodia Console
					</span>
					<span class="mt-1 text-[10px] font-black tracking-[0.2em] text-slate-400 uppercase">
						Admin Dashboard
					</span>
				</div>
			{/if}
		</a>
	</div>

	<!-- Navigation -->
	<nav class="custom-scrollbar flex-1 space-y-3 overflow-y-auto px-4 py-8">
		{#each menuItems as item}
			{@const isActive = page.url.pathname === item.href}
			<a
				href={item.href}
				onclick={() => onNavigate?.()}
				class={cn(
					'group relative flex items-center gap-4 overflow-hidden rounded-2xl px-4 py-3.5 transition-all duration-300',
					isActive
						? 'bg-primary text-white shadow-[0_8px_30px_rgb(59,130,246,0.3)] ring-1 ring-white/20'
						: 'text-slate-500 hover:bg-slate-100/50 hover:text-slate-900 dark:text-slate-400 dark:hover:bg-white/5 dark:hover:text-white'
				)}
			>
				{#if isActive}
					<div class="from-primary to-secondary absolute inset-0 bg-linear-to-r opacity-100"></div>
					<div class="absolute top-0 -right-2 bottom-0 w-8 rotate-12 bg-white/20 blur-xl"></div>
				{/if}

				<item.icon
					class={cn(
						'relative z-10 h-5.5 w-5.5',
						isActive ? 'text-white' : 'transition-transform group-hover:scale-110'
					)}
				/>

				{#if !collapsed}
					<span class="relative z-10 text-sm font-black">{item.name}</span>
				{/if}

				{#if isActive && !collapsed}
					<div class="relative z-10 ml-auto h-1.5 w-1.5 animate-pulse rounded-full bg-white"></div>
				{/if}
			</a>
		{/each}
	</nav>

	<!-- Bottom Section -->
	<div class="relative space-y-3 px-4 py-8">
		<div
			class="absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-slate-200/50 to-transparent dark:via-white/5"
		></div>

		<!-- Minimize Button -->
		<button
			onclick={() => (collapsed = !collapsed)}
			class="group flex w-full items-center gap-4 rounded-2xl px-4 py-3.5 text-slate-500 transition-all hover:bg-slate-100/50 active:scale-95 dark:text-slate-400 dark:hover:bg-white/5"
		>
			<ChevronLeft
				class={cn('h-5.5 w-5.5 transition-transform duration-500', collapsed && 'rotate-180')}
			/>
			{#if !collapsed}
				<span class="text-sm font-black group-hover:text-slate-900 dark:group-hover:text-white"
					>Collapse</span
				>
			{/if}
		</button>
	</div>
</aside>
