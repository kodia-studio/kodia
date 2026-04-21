<script lang="ts">
	import Navbar from '$lib/components/ui/Navbar.svelte';
	import {
		Zap,
		Shield,
		Terminal,
		Layers,
		ArrowRight,
		Sparkles,
		Rocket,
		Code,
		GitBranch
	} from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';

	let visible = $state(false);
	let commandIndex = $state(0);

	onMount(() => {
		visible = true;
		const interval = setInterval(() => {
			commandIndex = (commandIndex + 1) % commands.length;
		}, 3000);
		return () => clearInterval(interval);
	});

	const ecosystem = [
		{
			icon: Terminal,
			title: 'Powerful CLI',
			desc: 'Generate features, handlers, and migrations in one command. The fastest way to scaffold.'
		},
		{
			icon: Zap,
			title: 'Go Backend',
			desc: 'Lightning-fast APIs with Gin. Clean architecture that scales from MVP to millions of users.'
		},
		{
			icon: Sparkles,
			title: 'Svelte 5 Frontend',
			desc: 'Ultra-reactive UI with Runes. Write less, do more with the most intuitive framework.'
		},
		{
			icon: Shield,
			title: 'Security First',
			desc: 'JWT, RBAC, and enterprise-grade middleware. Ship confidently, audit easily.'
		},
		{
			icon: GitBranch,
			title: 'Git Native',
			desc: 'Full-stack generated from commands. Version control your architecture decisions.'
		},
		{
			icon: Rocket,
			title: 'Deploy Anywhere',
			desc: 'Docker, Kubernetes, or serverless. One codebase, infinite possibilities.'
		}
	];

	const commands = [
		'kodia new my-project',
		'kodia make:feature Users',
		'kodia make:model Product',
		'kodia migrate'
	];
</script>

<svelte:head>
	<title>Kodia Framework | One Command. Full Stack. Zero Friction.</title>
	<meta name="description" content="The fastest way to build full-stack applications with Go & Svelte. Kodia Framework delivers powerful CLI, clean architecture, and instant scaffolding." />
</svelte:head>

<Navbar />

<main class="overflow-hidden">
	<!-- Hero Section -->
	<section class="relative min-h-screen flex items-center justify-center pt-32 pb-20 md:pt-24 md:pb-32 px-6">
		<!-- Dynamic Grid Background -->
		<div class="absolute inset-0 -z-10">
			<div class="absolute inset-0 bg-linear-to-br from-white via-slate-50 to-white dark:from-slate-950 dark:via-slate-900 dark:to-slate-950"></div>
			<!-- Animated Hive Pattern -->
			<div class="absolute inset-0 opacity-30">
				<svg class="w-full h-full" viewBox="0 0 1200 600" preserveAspectRatio="xMidYMid slice">
					<defs>
						<pattern id="hexagons" x="0" y="0" width="120" height="104" patternUnits="userSpaceOnUse">
							<polygon points="30,0 60,17.32 60,51.96 30,69.28 0,51.96 0,17.32" fill="none" stroke="url(#hiveGradient)" stroke-width="1" opacity="0.3"/>
						</pattern>
						<linearGradient id="hiveGradient" x1="0%" y1="0%" x2="100%" y2="100%">
							<stop offset="0%" style="stop-color:#3b82f6;stop-opacity:1" />
							<stop offset="100%" style="stop-color:#06b6d4;stop-opacity:1" />
						</linearGradient>
					</defs>
					<rect width="1200" height="600" fill="url(#hexagons)"/>
				</svg>
			</div>
			<!-- Blue Glow -->
			<div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-blue-500/20 rounded-full blur-3xl"></div>
		</div>

		<div class="container mx-auto text-center relative z-10">
			{#if visible}
				<div in:fly={{ y: 30, duration: 1000, delay: 100 }}>
					<!-- Badge -->
					<div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-blue-500/10 border border-blue-500/30 text-blue-400 text-sm font-medium mb-8">
						<span class="relative flex h-2 w-2">
							<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
							<span class="relative inline-flex rounded-full h-2 w-2 bg-blue-400"></span>
						</span>
						The Hive is Live
					</div>

					<!-- Main Heading -->
					<h1 class="text-6xl md:text-7xl lg:text-8xl font-black mb-8 leading-tight text-slate-900 dark:text-white max-w-5xl mx-auto">
						One Command.<br>
						<span class="bg-linear-to-r from-blue-600 dark:from-blue-400 to-cyan-600 dark:to-cyan-400 bg-clip-text text-transparent">Full Stack.</span><br>
						<span class="text-slate-600 dark:text-slate-300">Zero Friction.</span>
					</h1>

					<!-- Subtitle -->
					<p class="text-lg md:text-xl text-slate-600 dark:text-slate-400 max-w-3xl mx-auto mb-12 leading-relaxed">
						Experience the future of full-stack development. Go backend, Svelte frontend, and a CLI so powerful you'll never go back.
					</p>

					<!-- CTA Buttons -->
					<div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16">
						<a href="/docs/prologue/getting-started" class="group relative px-8 py-4 bg-linear-to-r from-blue-500 to-blue-600 text-white font-bold rounded-xl hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 flex items-center gap-2 w-full sm:w-auto justify-center">
							Start Building Now
							<ArrowRight class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
						</a>
						<a href="https://github.com/kodia-studio/kodia" target="_blank" class="px-8 py-4 bg-slate-800/50 border border-slate-700 text-white font-bold rounded-xl hover:bg-slate-800 transition-all duration-300 w-full sm:w-auto">
							View on GitHub
						</a>
					</div>
				</div>

				<!-- Live Command Demo -->
				<div in:fade={{ duration: 1200, delay: 400 }} class="mt-20 relative px-4">
					<div class="max-w-3xl mx-auto">
						<!-- Terminal Window -->
						<div class="bg-slate-900/50 backdrop-blur-xl border border-slate-700/50 rounded-2xl p-6 shadow-2xl">
							<div class="flex items-center gap-2 mb-6">
								<div class="flex gap-1.5">
									<div class="w-3 h-3 rounded-full bg-red-500/70"></div>
									<div class="w-3 h-3 rounded-full bg-yellow-500/70"></div>
									<div class="w-3 h-3 rounded-full bg-green-500/70"></div>
								</div>
								<span class="text-xs text-slate-400 ml-3">Terminal</span>
							</div>

							<div class="font-mono text-left">
								<div class="flex items-center gap-3 text-sm">
									<span class="text-blue-400">$</span>
									<span class="text-slate-300 flex items-center">
										<span class="text-slate-300 font-semibold inline">kodia</span>
										<span class="ml-2 text-slate-500">
											{#key commandIndex}
												<span in:fly={{ y: 20, duration: 500 }}>
													{commands[commandIndex]}
												</span>
											{/key}
										</span>
									</span>
									<span class="text-slate-500 ml-auto animate-pulse">▌</span>
								</div>
								<div class="mt-6 space-y-2 text-xs text-slate-400">
									<p>🚀 Initializing hive structure...</p>
									<p class="text-green-400">✔ Backend generated (Go + Gin)</p>
									<p class="text-green-400">✔ Frontend generated (Svelte 5)</p>
									<p class="text-green-400">✔ Database configured</p>
									<p class="text-cyan-400">→ Ready to ship</p>
								</div>
							</div>
						</div>

						<p class="text-slate-600 dark:text-slate-400 text-sm mt-4 text-center">
							The entire project setup in seconds. No configuration hell.
						</p>
					</div>
				</div>
			{/if}
		</div>
	</section>

	<!-- The Hive Ecosystem -->
	<section class="py-32 px-6 relative">
		<div class="container mx-auto">
			<div class="text-center mb-24">
				<span class="inline-block px-4 py-2 rounded-full bg-blue-500/10 border border-blue-500/30 text-blue-400 text-sm font-semibold mb-6">
					The Kodia Hive
				</span>
				<h2 class="text-5xl md:text-6xl font-black mb-8">
					Everything you need,<br>
					<span class="bg-linear-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent">nothing you don't</span>
				</h2>
				<p class="text-lg text-slate-600 dark:text-slate-400 max-w-2xl mx-auto">
					A tightly integrated ecosystem designed for teams that move fast and ship quality. Every piece matters.
				</p>
			</div>

			<!-- 3x2 Grid -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
				{#each ecosystem as item, i}
					<div
						class="group relative p-8 bg-white dark:bg-slate-900/50 border border-slate-200 dark:border-slate-800 rounded-2xl hover:border-blue-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/10"
						in:fly={{ y: 30, duration: 500, delay: 100 + i * 50 }}
					>
						<!-- Icon -->
						<div class="mb-6 relative">
							<div class="absolute inset-0 bg-blue-500/20 blur-xl group-hover:blur-2xl transition-all rounded-xl"></div>
							<div class="relative w-14 h-14 bg-linear-to-br from-blue-500 to-cyan-500 rounded-xl flex items-center justify-center text-white">
								{#if item.icon === Terminal}
									<Terminal class="w-7 h-7" />
								{:else if item.icon === Zap}
									<Zap class="w-7 h-7" />
								{:else if item.icon === Sparkles}
									<Sparkles class="w-7 h-7" />
								{:else if item.icon === Shield}
									<Shield class="w-7 h-7" />
								{:else if item.icon === GitBranch}
									<GitBranch class="w-7 h-7" />
								{:else if item.icon === Rocket}
									<Rocket class="w-7 h-7" />
								{/if}
							</div>
						</div>

						<h3 class="text-xl font-bold mb-3 text-slate-900 dark:text-white">{item.title}</h3>
						<p class="text-slate-600 dark:text-slate-400 leading-relaxed">{item.desc}</p>

						<!-- Accent line -->
						<div class="absolute bottom-0 left-0 right-0 h-px bg-linear-to-r from-transparent via-blue-500 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- Why Kodia Section -->
	<section class="py-32 px-6 bg-linear-to-b from-slate-100 dark:from-slate-900/50 to-slate-50 dark:to-slate-950/50 border-y border-slate-200 dark:border-slate-800">
		<div class="container mx-auto">
			<div class="max-w-4xl mx-auto">
				<h2 class="text-5xl md:text-6xl font-black mb-12 text-center text-slate-900 dark:text-white">
					Built for <span class="text-blue-600 dark:text-blue-400">modern teams</span>
				</h2>

				<div class="grid md:grid-cols-2 gap-12">
					<div>
						<h3 class="text-2xl font-bold mb-6 text-slate-900 dark:text-white flex items-center gap-3">
							<Code class="w-8 h-8 text-blue-600 dark:text-blue-400" />
							Developer Experience
						</h3>
						<ul class="space-y-4 text-slate-700 dark:text-slate-300">
							<li class="flex gap-3">
								<span class="text-blue-600 dark:text-blue-400 font-bold">›</span>
								<span>Type-safe from database to UI</span>
							</li>
							<li class="flex gap-3">
								<span class="text-blue-600 dark:text-blue-400 font-bold">›</span>
								<span>Hot reload everything by default</span>
							</li>
							<li class="flex gap-3">
								<span class="text-blue-600 dark:text-blue-400 font-bold">›</span>
								<span>Zero-config development environment</span>
							</li>
							<li class="flex gap-3">
								<span class="text-blue-600 dark:text-blue-400 font-bold">›</span>
								<span>Instant scaffolding for common patterns</span>
							</li>
						</ul>
					</div>

					<div>
						<h3 class="text-2xl font-bold mb-6 text-slate-900 dark:text-white flex items-center gap-3">
							<Rocket class="w-8 h-8 text-cyan-600 dark:text-cyan-400" />
							Production Ready
						</h3>
						<ul class="space-y-4 text-slate-700 dark:text-slate-300">
							<li class="flex gap-3">
								<span class="text-cyan-600 dark:text-cyan-400 font-bold">›</span>
								<span>Sub-millisecond Go backend APIs</span>
							</li>
							<li class="flex gap-3">
								<span class="text-cyan-600 dark:text-cyan-400 font-bold">›</span>
								<span>Enterprise auth & security built-in</span>
							</li>
							<li class="flex gap-3">
								<span class="text-cyan-600 dark:text-cyan-400 font-bold">›</span>
								<span>Docker & deployment ready from day one</span>
							</li>
							<li class="flex gap-3">
								<span class="text-cyan-600 dark:text-cyan-400 font-bold">›</span>
								<span>Observable and audit-friendly architecture</span>
							</li>
						</ul>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="py-32 px-6">
		<div class="container mx-auto max-w-4xl">
			<div class="bg-linear-to-r from-blue-600/20 to-cyan-600/20 border border-blue-500/30 rounded-3xl p-12 md:p-20 text-center relative overflow-hidden">
				<!-- Decorative Elements -->
				<div class="absolute -top-20 -right-20 w-40 h-40 bg-blue-500/10 rounded-full blur-3xl"></div>
				<div class="absolute -bottom-20 -left-20 w-40 h-40 bg-cyan-500/10 rounded-full blur-3xl"></div>

				<div class="relative z-10">
					<h2 class="text-5xl md:text-6xl font-black mb-8 text-slate-900 dark:text-white">
						Ready to join the Hive?
					</h2>
					<p class="text-xl text-slate-700 dark:text-slate-300 mb-12 max-w-2xl mx-auto leading-relaxed">
						Stop wasting time on boilerplate. Start building what matters. The future of full-stack development starts here.
					</p>

					<div class="flex flex-col sm:flex-row items-center justify-center gap-4">
						<a href="/docs/prologue/getting-started" class="group px-10 py-4 bg-linear-to-r from-blue-500 to-blue-600 text-white font-bold rounded-xl hover:shadow-lg hover:shadow-blue-500/50 transition-all flex items-center gap-2 w-full sm:w-auto justify-center">
							Get Started Free
							<ArrowRight class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
						</a>
						<a href="https://github.com/kodia-studio/kodia" target="_blank" class="px-10 py-4 bg-slate-800 border border-slate-700 text-white font-bold rounded-xl hover:bg-slate-700 transition-all w-full sm:w-auto">
							GitHub Repository
						</a>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Footer -->
	<footer class="py-12 border-t border-slate-200 dark:border-slate-800 px-6">
		<div class="container mx-auto flex flex-col md:flex-row items-center justify-between gap-6">
			<div class="flex items-center gap-2">
				<Layers class="text-blue-500 dark:text-blue-400 w-5 h-5" />
				<span class="font-bold font-heading text-slate-900 dark:text-white">Kodia<span class="text-blue-500 dark:text-blue-400">.</span></span>
				<span class="text-xs text-slate-600 dark:text-slate-500 ml-4">© 2026 Kodia Studio. MIT License.</span>
			</div>

			<div class="flex items-center gap-8 text-sm text-slate-600 dark:text-slate-400">
				<a href="https://github.com/kodia-studio/kodia" class="hover:text-blue-500 dark:hover:text-blue-400 transition-colors">GitHub</a>
				<a href="https://kodia.id/privacy" class="hover:text-blue-500 dark:hover:text-blue-400 transition-colors">Privacy Policy</a>
			</div>
		</div>
	</footer>
</main>

