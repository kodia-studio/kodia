<script lang="ts">
	import { themeStore } from '$lib/stores/theme.store';
	import { Sun, Moon, Globe, Layers } from 'lucide-svelte';
	import { fade } from 'svelte/transition';

	let isScrolled = $state(false);

	if (typeof window !== 'undefined') {
		window.addEventListener('scroll', () => {
			isScrolled = window.scrollY > 20;
		});
	}

	function toggleTheme() {
		const current = localStorage.getItem('theme') || 'system';
		const next = current === 'dark' ? 'light' : 'dark';
		themeStore.set(next);
	}
</script>

<nav
	class="fixed top-0 left-0 right-0 z-40 transition-all duration-300 {isScrolled
		? 'py-3 glass shadow-lg'
		: 'py-6 bg-transparent'}"
>
	<div class="container mx-auto px-6 flex items-center justify-between">
		<a href="/" class="flex items-center gap-2 group">
			<div
				class="w-10 h-10 accent-gradient rounded-xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform"
			>
				<Layers class="text-white w-6 h-6" />
			</div>
			<span class="text-2xl font-bold font-heading tracking-tight">
				Kodia<span class="text-primary-600">.</span>
			</span>
		</a>

		<div class="hidden md:flex items-center gap-8 text-sm font-medium">
			<a href="/docs" class="hover:text-primary-600 transition-colors">Documentation</a>
			<a href="/examples" class="hover:text-primary-600 transition-colors">Examples</a>
			<a href="/showcase" class="hover:text-primary-600 transition-colors">Showcase</a>
		</div>

		<div class="flex items-center gap-3">
			<a
				href="https://github.com/kodia-studio/kodia"
				target="_blank"
				class="p-2 hover:bg-slate-200 dark:hover:bg-slate-800 rounded-lg transition-colors"
			>
				<Globe class="w-5 h-5" />
			</a>

			<button
				onclick={toggleTheme}
				class="p-2 hover:bg-slate-200 dark:hover:bg-slate-800 rounded-lg transition-colors overflow-hidden relative w-9 h-9"
				aria-label="Toggle theme"
			>
				<div class="flex flex-col items-center justify-center h-full">
					<div class="dark:hidden">
						<Sun class="w-5 h-5" />
					</div>
					<div class="hidden dark:block">
						<Moon class="w-5 h-5" />
					</div>
				</div>
			</button>

			<a href="/docs/get-started" class="btn-primary py-2 px-5 text-sm"> Get Started </a>
		</div>
	</div>
</nav>
