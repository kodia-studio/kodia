<script lang="ts">
    import { onMount } from 'svelte';
    import { themeStore } from '$lib/stores/theme.store';
    import { authStore } from '$lib/stores/auth.store';
    import { Sun, Moon, Menu, Zap } from 'lucide-svelte';

    let scrolled = $state(false);
    let isDark = $state(false);

    onMount(() => {
        const handleScroll = () => {
            scrolled = window.scrollY > 20;
        };
        window.addEventListener('scroll', handleScroll);

        // Track effective theme
        isDark = document.documentElement.classList.contains('dark');
        const observer = new MutationObserver(() => {
            isDark = document.documentElement.classList.contains('dark');
        });
        observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });

        return () => {
            window.removeEventListener('scroll', handleScroll);
            observer.disconnect();
        };
    });
</script>

<nav class="fixed top-0 left-0 right-0 z-50 transition-all duration-300 border-b border-t-0 border-slate-200/50 dark:border-slate-800/50 glass shadow-sm {scrolled ? 'py-3' : 'py-4'}">
    <div class="container mx-auto px-6 flex items-center justify-between">
            <!-- Logo -->
            <a href="/" class="flex items-center gap-2 group">
                <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-white font-black group-hover:rotate-12 transition-transform shadow-lg shadow-primary/20">
                    <Zap size={18} fill="currentColor" />
                </div>
                <span class="text-xl font-black tracking-tighter">KODIA</span>
            </a>

            <!-- Desktop Links -->
            <div class="hidden md:flex items-center gap-8">
                <a href="https://kodia.id/docs/prologue/getting-started" class="text-sm font-medium hover:text-primary transition-colors">Documentation</a>
                <a href="/state-management-demo" class="text-sm font-medium hover:text-primary transition-colors">Data Layer</a>
                <a href="https://kodia.id/features" class="text-sm font-medium hover:text-primary transition-colors">Features</a>
                <a href="https://github.com/kodia-studio/kodia" target="_blank" class="text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors p-2" aria-label="GitHub">
                    <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v 3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                </a>

                <div class="w-px h-4 bg-slate-200 dark:bg-slate-700"></div>

                <button
                    onclick={() => themeStore.toggle()}
                    class="p-2 rounded-xl hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors text-slate-500 dark:text-slate-400"
                    aria-label="Toggle Theme"
                >
                    {#if isDark}
                        <Sun size={20} />
                    {:else}
                        <Moon size={20} />
                    {/if}
                </button>

                {#if $authStore.isAuthenticated}
                    <a href="/dashboard" class="btn-premium py-2 px-5 text-sm">
                        Dashboard
                    </a>
                {:else}
                    <a href="/login" class="text-sm font-medium hover:text-primary transition-colors">
                        Login
                    </a>
                    <a href="/register" class="btn-premium py-2 px-5 text-sm">
                        Register
                    </a>
                {/if}
            </div>

            <!-- Mobile Menu Toggle (simplified) -->
            <button class="md:hidden p-2 text-slate-500" aria-label="Toggle Menu">
                <Menu size={24} />
            </button>
    </div>
</nav>

<style>
    nav :global(.glass) {
        border-color: rgba(255, 255, 255, 0.1);
    }
</style>
