<script lang="ts">
	import { devStore } from '$lib/stores/dev.svelte';
	import { onMount } from 'svelte';
	import { 
		Zap, X, Terminal, Globe, Wifi, Trash2, 
		ChevronRight, ChevronDown, CheckCircle2, 
		AlertCircle, Clock, Activity 
	} from 'lucide-svelte';
	import Button from '../ui/Button.svelte';
	import { fade, slide, scale } from 'svelte/transition';

	let activeTab = $state<'network' | 'socket' | 'state'>('network');
	let selectedLogId = $state<string | null>(null);

	onMount(() => {
		const handleKeyDown = (e: KeyboardEvent) => {
			if (e.ctrlKey && e.key === 'k') {
				e.preventDefault();
				devStore.toggle();
			}
		};
		window.addEventListener('keydown', handleKeyDown);
		return () => window.removeEventListener('keydown', handleKeyDown);
	});

	function formatDuration(ms?: number) {
		if (!ms) return '-';
		return ms < 1000 ? `${ms.toFixed(0)}ms` : `${(ms / 1000).toFixed(2)}s`;
	}

	const selectedLog = $derived(devStore.logs.find(l => l.id === selectedLogId));
</script>

{#if devStore.isOpen}
	<div 
		class="fixed inset-y-0 right-0 w-[500px] bg-slate-950 text-slate-300 shadow-2xl z-[9999] border-l border-white/10 flex flex-col font-mono text-xs"
		transition:slide={{ axis: 'x', duration: 300 }}
	>
		<!-- Header -->
		<div class="p-4 border-b border-white/10 flex items-center justify-between bg-slate-900/50">
			<div class="flex items-center gap-2">
				<div class="w-6 h-6 bg-primary rounded flex items-center justify-center text-white">
					<Zap size={14} fill="currentColor" />
				</div>
				<span class="font-black tracking-tight text-white uppercase text-[10px]">Kodia DevTools</span>
			</div>
			<div class="flex items-center gap-2">
				<button 
					onclick={() => devStore.clear()}
					class="p-1.5 hover:bg-white/5 rounded transition-colors text-slate-500 hover:text-white"
					title="Clear Logs"
				>
					<Trash2 size={14} />
				</button>
				<button 
					onclick={() => devStore.toggle()}
					class="p-1.5 hover:bg-white/5 rounded transition-colors text-slate-500 hover:text-white"
				>
					<X size={16} />
				</button>
			</div>
		</div>

		<!-- Tabs -->
		<div class="flex border-b border-white/5 px-2">
			{#each ['network', 'socket', 'state'] as tab}
				<button
					onclick={() => activeTab = tab as any}
					class="px-4 py-3 border-b-2 transition-all uppercase tracking-widest text-[9px] font-black {activeTab === tab ? 'border-primary text-white' : 'border-transparent text-slate-500 hover:text-slate-300'}"
				>
					{tab}
				</button>
			{/each}
		</div>

		<!-- Content -->
		<div class="flex-1 overflow-hidden flex flex-col">
			{#if activeTab === 'network'}
				<div class="flex-1 overflow-y-auto custom-scrollbar">
					{#each devStore.logs as log}
						<button
							onclick={() => selectedLogId = log.id === selectedLogId ? null : log.id}
							class="w-full text-left p-3 border-b border-white/5 hover:bg-white/5 transition-colors group {selectedLogId === log.id ? 'bg-primary/10' : ''}"
						>
							<div class="flex items-center justify-between mb-1">
								<div class="flex items-center gap-2">
									<span class="font-black text-[10px] {log.status && log.status < 400 ? 'text-emerald-500' : 'text-rose-500'}">
										{log.status || 'PENDING'}
									</span>
									<span class="text-white font-bold">{log.method}</span>
									<span class="text-slate-500 truncate max-w-[200px]">{log.path}</span>
								</div>
								<div class="text-[9px] text-slate-500">
									{formatDuration(log.duration)}
								</div>
							</div>
							
							{#if selectedLogId === log.id}
								<div class="mt-4 space-y-4" transition:slide>
									{#if log.requestBody}
										<div class="space-y-1">
											<div class="text-[9px] font-black uppercase text-slate-500">Payload</div>
											<pre class="p-2 bg-black/40 rounded border border-white/5 overflow-x-auto text-[10px]">{JSON.stringify(log.requestBody, null, 2)}</pre>
										</div>
									{/if}
									{#if log.responseBody}
										<div class="space-y-1">
											<div class="text-[9px] font-black uppercase text-slate-500">Response</div>
											<pre class="p-2 bg-black/40 rounded border border-white/5 overflow-x-auto text-[10px]">{JSON.stringify(log.responseBody, null, 2)}</pre>
										</div>
									{/if}
								</div>
							{/if}
						</button>
					{:else}
						<div class="h-full flex flex-col items-center justify-center text-slate-600 gap-2 opacity-50">
							<Activity size={32} />
							<p class="text-[10px] font-black uppercase tracking-widest">Awaiting Activity...</p>
						</div>
					{/each}
				</div>

			{:else if activeTab === 'socket'}
				<div class="flex-1 overflow-y-auto custom-scrollbar p-2 space-y-1">
					{#each devStore.socketLogs as log}
						<div class="p-2 rounded border border-white/5 bg-white/[0.02] flex items-start gap-3">
							<div class="mt-0.5">
								{#if log.type === 'send'}
									<ChevronRight size={12} class="text-primary" />
								{:else}
									<ChevronRight size={12} class="text-emerald-500 rotate-180" />
								{/if}
							</div>
							<div class="flex-1 overflow-hidden">
								<div class="flex items-center justify-between mb-1">
									<span class="text-[9px] font-black uppercase tracking-tighter {log.type === 'send' ? 'text-primary' : 'text-emerald-500'}">
										{log.type === 'send' ? 'Outgoing' : 'Incoming'}
									</span>
									<span class="text-[8px] text-slate-600">
										{log.timestamp.toLocaleTimeString()}
									</span>
								</div>
								<pre class="text-[10px] text-slate-400 overflow-x-auto">{JSON.stringify(log.data, null, 2)}</pre>
							</div>
						</div>
					{:else}
						<div class="h-full flex flex-col items-center justify-center text-slate-600 gap-2 opacity-50">
							<Wifi size={32} />
							<p class="text-[10px] font-black uppercase tracking-widest">No Socket Traffic</p>
						</div>
					{/each}
				</div>

			{:else if activeTab === 'state'}
				<div class="p-4 space-y-6 overflow-y-auto custom-scrollbar flex-1">
					<div class="space-y-2">
						<div class="flex items-center gap-2 text-primary">
							<Globe size={14} />
							<span class="text-[10px] font-black uppercase tracking-widest">Environment</span>
						</div>
						<div class="grid grid-cols-2 gap-2">
							<div class="p-2 bg-white/5 rounded border border-white/5">
								<div class="text-[8px] text-slate-500 uppercase font-black">Mode</div>
								<div class="text-white">Development</div>
							</div>
							<div class="p-2 bg-white/5 rounded border border-white/5">
								<div class="text-[8px] text-slate-500 uppercase font-black">Framework</div>
								<div class="text-white">Svelte 5 (Runes)</div>
							</div>
						</div>
					</div>
					
					<div class="space-y-2">
						<div class="flex items-center gap-2 text-emerald-500">
							<Terminal size={14} />
							<span class="text-[10px] font-black uppercase tracking-widest">Quick Commands</span>
						</div>
						<div class="flex flex-wrap gap-2">
							<Button variant="outline" class="h-7 text-[9px] border-white/10 hover:bg-white/10" onclick={() => location.reload()}>
								Hard Reload
							</Button>
							<Button variant="outline" class="h-7 text-[9px] border-white/10 hover:bg-white/10" onclick={() => devStore.clear()}>
								Purge Logs
							</Button>
						</div>
					</div>
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="p-3 border-t border-white/5 bg-black/20 flex items-center justify-between text-[9px] font-bold text-slate-500">
			<div class="flex items-center gap-4">
				<div class="flex items-center gap-1">
					<div class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></div>
					Engine Active
				</div>
				<div>Ctrl + K to toggle</div>
			</div>
			<div>v1.0.0-PRO</div>
		</div>
	</div>
{/if}

<style>
	.custom-scrollbar::-webkit-scrollbar {
		width: 4px;
	}
	.custom-scrollbar::-webkit-scrollbar-track {
		background: transparent;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.1);
		border-radius: 20px;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.2);
	}
</style>
