<script lang="ts">
	import { Button, Input, Badge, Card } from '$lib/components/ui';
	import InfiniteList from '$lib/components/data/InfiniteList.svelte';
	import VirtualList from '$lib/components/data/VirtualList.svelte';
	import { createOptimistic } from '$lib/api/optimistic.svelte';
	import { Loader2, RefreshCcw, Wifi, WifiOff, Zap } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	// 1. Optimistic Update Demo
	const userStatus = createOptimistic('Offline', async (newStatus) => {
		// Simulate API call
		await new Promise(resolve => setTimeout(resolve, 2000));
		if (newStatus === 'Error') throw new Error('Failed to update status');
	});

	// 2. Infinite List Demo
	let infiniteItems = $state(Array.from({ length: 20 }, (_, i) => `Item ${i + 1}`));
	let hasMore = $state(true);
	let loadingInfinite = $state(false);

	async function loadMore() {
		loadingInfinite = true;
		await new Promise(resolve => setTimeout(resolve, 1500));
		const next = infiniteItems.length;
		infiniteItems = [...infiniteItems, ...Array.from({ length: 10 }, (_, i) => `Item ${next + i + 1}`)];
		if (infiniteItems.length >= 50) hasMore = false;
		loadingInfinite = false;
	}

	// 3. Virtual List Demo
	const virtualItems = Array.from({ length: 1000 }, (_, i) => ({
		id: i,
		title: `High Performance Data Row ${i}`,
		value: Math.floor(Math.random() * 1000)
	}));

	let activeTab = $state<'optimistic' | 'infinite' | 'virtual'>('optimistic');
</script>

<div class="min-h-screen pt-24 pb-12 px-6">
	<div class="max-w-5xl mx-auto space-y-12">
		<!-- Header -->
		<div class="space-y-4">
			<Badge variant="premium">Institutional State Management</Badge>
			<h1 class="text-5xl font-black tracking-tight text-slate-900 dark:text-white">
				Data <span class="text-primary">Intelligence</span> Layer.
			</h1>
			<p class="text-xl text-slate-500 font-medium max-w-2xl">
				Experience high-velocity UX with Svelte 5 Runes, Optimistic Updates, and Elite Data Virtualization.
			</p>
		</div>

		<!-- Tab Navigation -->
		<div class="flex gap-2 p-1 glass rounded-2xl w-fit">
			{#each ['optimistic', 'infinite', 'virtual'] as tab}
				<button
					onclick={() => activeTab = tab as any}
					class="px-6 py-2.5 rounded-xl text-xs font-black uppercase tracking-widest transition-all {activeTab === tab ? 'bg-primary text-white shadow-lg shadow-primary/20' : 'text-slate-500 hover:text-slate-900 dark:hover:text-white'}"
				>
					{tab}
				</button>
			{/each}
		</div>

		<!-- Content Sections -->
		{#if activeTab === 'optimistic'}
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8" in:fade>
				<div class="p-8 glass rounded-3xl border border-slate-200 dark:border-slate-800 space-y-6">
					<div class="flex items-center justify-between">
						<h2 class="text-2xl font-black tracking-tight">Status Sync</h2>
						<div class="flex items-center gap-2">
							{#if userStatus.isSyncing}
								<Loader2 class="w-4 h-4 animate-spin text-primary" />
							{/if}
							<Badge variant={userStatus.value === 'Online' ? 'success' : 'outline'}>
								{userStatus.value}
							</Badge>
						</div>
					</div>

					<p class="text-slate-500 text-sm leading-relaxed">
						Update your status instantly. If the server fails, the state will инновативно rollback to its previous value.
					</p>

					<div class="flex gap-3">
						<Button 
							variant="primary" 
							class="flex-1"
							onclick={() => userStatus.update('Online')}
						>
							Go Online
						</Button>
						<Button 
							variant="outline" 
							class="flex-1"
							onclick={() => userStatus.update('Offline')}
						>
							Go Offline
						</Button>
						<Button 
							variant="danger"
							onclick={() => userStatus.update('Error').catch(() => toast.error('Rollback Triggered!'))}
						>
							Trigger Error
						</Button>
					</div>

					{#if userStatus.error}
						<div class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-500 text-xs font-bold uppercase tracking-tight">
							Error: {userStatus.error.message} (State Restored)
						</div>
					{/if}
				</div>

				<div class="p-8 glass rounded-3xl border border-slate-200 dark:border-slate-800 bg-slate-900 text-white space-y-6">
					<div class="p-4 rounded-2xl bg-white/5 border border-white/10 space-y-2">
						<div class="flex items-center gap-2 text-[10px] font-black uppercase tracking-widest text-primary">
							<Zap class="w-3 h-3" />
							Real-time Intelligence
						</div>
						<p class="text-sm font-medium text-slate-400">
							The UI updates in <span class="text-white">~0ms</span>, while the server synchronization happens asynchronously in the background.
						</p>
					</div>
					<div class="flex items-center gap-4 text-xs font-bold text-slate-500">
						<div class="flex items-center gap-1"><Wifi class="w-3 h-3" /> Low Latency</div>
						<div class="flex items-center gap-1"><RefreshCcw class="w-3 h-3" /> Auto-Sync</div>
					</div>
				</div>
			</div>

		{:else if activeTab === 'infinite'}
			<div class="glass rounded-3xl border border-slate-200 dark:border-slate-800 overflow-hidden" in:fade>
				<div class="p-6 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between">
					<h2 class="text-xl font-black tracking-tight">Infinite Intelligence Stream</h2>
					<Badge>{infiniteItems.length} Items Loaded</Badge>
				</div>
				<div class="max-h-[500px] overflow-y-auto custom-scrollbar p-6">
					<InfiniteList 
						items={infiniteItems} 
						{hasMore} 
						loading={loadingInfinite} 
						{loadMore}
						class="grid grid-cols-1 sm:grid-cols-2 gap-4"
					>
						{#snippet children(item)}
							<div class="p-4 glass rounded-2xl border border-slate-100 dark:border-slate-800 hover:border-primary/50 transition-all group">
								<div class="text-xs font-black text-primary mb-1 uppercase tracking-widest">Kodia Unit</div>
								<div class="font-bold">{item}</div>
							</div>
						{/snippet}
					</InfiniteList>
				</div>
			</div>

		{:else if activeTab === 'virtual'}
			<div class="glass rounded-3xl border border-slate-200 dark:border-slate-800 overflow-hidden" in:fade>
				<div class="p-6 border-b border-slate-100 dark:border-slate-800">
					<h2 class="text-xl font-black tracking-tight">Institutional Data Virtualization</h2>
					<p class="text-xs text-slate-500 font-bold uppercase tracking-widest mt-1">Rendering 1,000 items with zero performance degradation</p>
				</div>
				<VirtualList 
					items={virtualItems} 
					itemHeight={80} 
					height={400}
					class="custom-scrollbar"
				>
					{#snippet children(item, index)}
						<div class="px-6 h-[80px] flex items-center justify-between border-b border-slate-50 dark:border-slate-900 hover:bg-slate-50/50 dark:hover:bg-slate-900/50 transition-colors">
							<div class="flex items-center gap-4">
								<div class="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-black text-xs">
									{index}
								</div>
								<div>
									<div class="font-bold text-sm">{item.title}</div>
									<div class="text-[10px] font-black text-slate-400 uppercase tracking-widest">ID: {item.id}</div>
								</div>
							</div>
							<div class="text-right">
								<div class="text-sm font-black text-emerald-500">+{item.value}</div>
								<div class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">KODIA COINS</div>
							</div>
						</div>
					{/snippet}
				</VirtualList>
			</div>
		{/if}
	</div>
</div>

<style>
	.glass {
		background: rgba(255, 255, 255, 0.02);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
	}
	:global(.dark) .glass {
		background: rgba(15, 23, 42, 0.6);
	}
</style>
