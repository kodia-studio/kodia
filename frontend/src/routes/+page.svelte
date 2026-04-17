<script lang="ts">
	import Navbar from '$lib/components/ui/Navbar.svelte';
	import { 
		Zap, 
		ShieldCheck, 
		Box, 
		Terminal, 
		ChevronRight, 
		Code, 
		Cpu, 
		Globe, 
		Layers,
		ArrowRight 
	} from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';

	let visible = $state(false);

	onMount(() => {
		visible = true;
	});

	const features = [
		{
			icon: Cpu,
			title: 'Golang Performance',
			desc: 'Harness the speed of Go with Gin for ultra-low latency backend APIs.'
		},
		{
			icon: Box,
			title: 'Svelte 5 Runes',
			desc: 'The most reactive frontend experience using the latest Svelte 5 technology.'
		},
		{
			icon: ShieldCheck,
			title: 'Secure by Default',
			desc: 'JWT auth, role-based access, and security middleware ready out-of-the-box.'
		},
		{
			icon: Terminal,
			title: 'Instant Scaffolding',
			desc: 'Our powerful CLI helps you build features, handlers, and repos in seconds.'
		},
		{
			icon: Zap,
			title: 'Clean Architecture',
			desc: 'Maintainable and testable code following enterprise-grade architecture patterns.'
		},
		{
			icon: Globe,
			title: 'Deployment Ready',
			desc: 'Pre-configured Docker and Makefile for seamless production deployment.'
		}
	];
</script>

<svelte:head>
	<title>Kodia Framework | Professional Fullstack framework for Go & Svelte</title>
	<meta name="description" content="Build professional applications faster with Kodia Framework." />
</svelte:head>

<Navbar />

<main class="overflow-hidden">
	<!-- Hero Section -->
	<section class="relative pt-32 pb-20 md:pt-48 md:pb-32 px-6">
		<!-- Background Decorations -->
		<div
			class="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-6xl h-full -z-10 pointer-events-none"
		>
			<div
				class="absolute top-10 left-10 w-72 h-72 bg-primary-500/10 rounded-full blur-3xl animate-pulse"
			></div>
			<div
				class="absolute bottom-20 right-10 w-96 h-96 bg-primary-600/10 rounded-full blur-3xl"
			></div>
		</div>

		<div class="container mx-auto text-center">
			{#if visible}
				<div in:fly={{ y: 20, duration: 800, delay: 100 }}>
					<span
						class="inline-block px-4 py-1.5 rounded-full bg-primary-500/10 text-primary-600 text-xs font-bold uppercase tracking-wider mb-6 border border-primary-500/20 shadow-inner"
					>
						Version 1.0 is here
					</span>
					<h1 class="text-5xl md:text-7xl lg:text-8xl font-black mb-8 leading-[1.1]">
						Code like a <span class="text-gradient">PRO</span><br />
						Build at <span class="text-primary-600">Light</span> Speed.
					</h1>
					<p class="text-lg md:text-xl text-slate-500 dark:text-slate-400 max-w-2xl mx-auto mb-12">
						The most professional full-stack framework for Gopher and Svelte developers. 
						Clean architecture, lightning-fast performance, and premium aesthetics.
					</p>

					<div class="flex flex-col sm:flex-row items-center justify-center gap-4">
						<a href="/docs/get-started" class="btn-primary w-full sm:w-auto flex items-center justify-center gap-2">
							Start Building Now
							<ArrowRight class="w-4 h-4" />
						</a>
						<a href="https://github.com/kodia-studio/kodia" target="_blank" class="btn-glass w-full sm:w-auto flex items-center justify-center gap-2">
							<Globe class="w-4 h-4" />
							Star on GitHub
						</a>
					</div>
				</div>

				<div in:fade={{ duration: 1000, delay: 600 }} class="mt-20 relative px-4">
					<div class="glass dark:glass-dark max-w-5xl mx-auto rounded-2xl p-2 shadow-2xl relative overflow-hidden group">
                        <!-- Code Mockup Area -->
						<div class="bg-slate-950 rounded-xl p-6 md:p-10 text-left font-mono text-sm overflow-x-auto">
							<div class="flex gap-1.5 mb-6">
								<div class="w-3 h-3 rounded-full bg-red-500"></div>
								<div class="w-3 h-3 rounded-full bg-yellow-500"></div>
								<div class="w-3 h-3 rounded-full bg-green-500"></div>
							</div>
							<div class="space-y-2">
								<p class="text-slate-500">// Initialize your next big thing</p>
								<p><span class="text-primary-400">$</span> kodia new <span class="text-yellow-400">awesome-app</span></p>
								<p class="text-slate-500 italic">// Generate features in seconds</p>
								<p><span class="text-primary-400">$</span> kodia make:feature <span class="text-yellow-400">Inventory</span></p>
								<br />
								<p class="text-blue-400">🚀 Starting Kodia Development Servers...</p>
								<p class="text-green-400">✅ Backend running on port 8080</p>
								<p class="text-green-400">✅ Frontend running on port 5173</p>
								<p class="text-slate-500 cursor-blink">_</p>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</section>

	<!-- Features Grid -->
	<section class="py-24 bg-slate-100 dark:bg-slate-900/30 border-y border-slate-200 dark:border-slate-800">
		<div class="container mx-auto px-6">
			<div class="mb-16 text-center max-w-3xl mx-auto">
				<h2 class="text-3xl md:text-5xl font-bold mb-6">Everything you need to ship.</h2>
				<p class="text-slate-500 dark:text-slate-400">
					Stop wasting time on boilerplate. Kodia provides everything from auth to deployment 
					so you can focus on your core product.
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
				{#each features as feature, i}
					<div 
						class="card group hover:border-primary-500/50"
						in:fade={{ duration: 500, delay: 200 * i }}
					>
						<div class="w-12 h-12 rounded-xl bg-primary-500/10 text-primary-600 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-primary-500 transition-all group-hover:text-white">
							<feature.icon class="w-6 h-6" />
						</div>
						<h3 class="text-xl font-bold mb-3">{feature.title}</h3>
						<p class="text-slate-500 dark:text-slate-400 leading-relaxed">
							{feature.desc}
						</p>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="py-24 relative overflow-hidden">
		<div class="container mx-auto px-6">
			<div class="accent-gradient rounded-3xl p-12 text-center text-white relative overflow-hidden shadow-2xl shadow-primary-500/20">
				<!-- Shapes -->
				<div class="absolute -top-10 -right-10 w-64 h-64 bg-white/10 rounded-full blur-3xl"></div>
				<div class="absolute -bottom-10 -left-10 w-64 h-64 bg-black/10 rounded-full blur-3xl"></div>
				
				<div class="relative z-10">
					<h2 class="text-4xl md:text-5xl font-black mb-6">Ready to build the future?</h2>
					<p class="text-primary-100 text-lg mb-10 max-w-xl mx-auto leading-relaxed">
						Join thousands of developers building scalable, clean, and professional applications with Kodia Framework.
					</p>
					<div class="flex flex-col sm:flex-row items-center justify-center gap-4">
						<a href="/docs" class="bg-white text-primary-600 px-8 py-3.5 rounded-xl font-bold hover:bg-primary-50 transition-all shadow-xl active:scale-95">
							Get Started Free
						</a>
						<a href="/community" class="bg-primary-700/50 text-white border border-primary-500 px-8 py-3.5 rounded-xl font-bold hover:bg-primary-700 transition-all active:scale-95">
							Join the Community
						</a>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="py-12 border-t border-slate-200 dark:border-slate-800">
		<div class="container mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6">
			<div class="flex items-center gap-2">
				<Layers class="text-primary-600 w-5 h-5" />
				<span class="font-bold font-heading">Kodia<span class="text-primary-600">.</span></span>
				<span class="text-xs text-slate-500 ml-4">© 2026 Kodia Studio. MIT License.</span>
			</div>
			
			<div class="flex items-center gap-8 text-sm text-slate-500">
				<a href="https://twitter.com/kodia" class="hover:text-primary-600">Twitter</a>
				<a href="https://github.com/kodia-studio/kodia" class="hover:text-primary-600">GitHub</a>
				<a href="/docs/privacy" class="hover:text-primary-600">Privacy Policy</a>
			</div>
		</div>
	</footer>
</main>

<style>
	:global(.cursor-blink) {
		animation: blink 1s step-end infinite;
	}

	@keyframes blink {
		from, to { opacity: 1; }
		50% { opacity: 0; }
	}
</style>
