<script lang="ts">
    import { onMount } from 'svelte';
    import { themeStore } from '$lib/stores/theme.store';
    import { authStore } from '$lib/stores/auth.store';
    import { Sun, Moon, Menu, X } from 'lucide-svelte';

    let scrolled = $state(false);
    let isDark = $state(false);
    let mobileMenuOpen = $state(false);

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

    function closeMobileMenu() {
        mobileMenuOpen = false;
    }
</script>

<nav class="fixed top-0 left-0 right-0 z-50 transition-all duration-300 bg-white dark:bg-slate-950 border-b border-slate-200 dark:border-slate-800 shadow-md hover:shadow-lg {scrolled ? 'py-3' : 'py-4'}">
    <div class="container mx-auto px-4 md:px-6 flex items-center justify-between">
        <!-- Logo -->
        <a href="/" class="flex items-center gap-2 group" onclick={closeMobileMenu}>
            <img src="/logo.png" alt="Kodia Logo" class="h-8 w-auto group-hover:scale-105 transition-transform" />
            <span class="text-lg md:text-xl font-black tracking-tighter">KODIA</span>
        </a>

        <!-- Desktop Links -->
        <div class="hidden md:flex items-center gap-8">
            <a href="https://kodia.id/docs/prologue/getting-started" class="text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-primary transition-colors">Documentation</a>
            <a href="https://kodia.id/features" class="text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-primary transition-colors">Features</a>

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
                <a href="/login" class="text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-primary transition-colors">
                    Login
                </a>
                <a href="/register" class="btn-premium py-2 px-5 text-sm">
                    Register
                </a>
            {/if}
        </div>

        <!-- Mobile Menu Toggle -->
        <div class="flex md:hidden items-center gap-2">
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
            <button
                onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
                class="p-2 text-slate-600 dark:text-slate-400 hover:text-primary transition-colors rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800"
                aria-label="Toggle Menu"
            >
                {#if mobileMenuOpen}
                    <X size={24} />
                {:else}
                    <Menu size={24} />
                {/if}
            </button>
        </div>
    </div>

    <!-- Mobile Menu -->
    {#if mobileMenuOpen}
        <div class="md:hidden border-t border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 animate-in slide-in-from-top-2">
            <div class="container mx-auto px-4 py-4 space-y-3">
                <a
                    href="https://kodia.id/docs/prologue/getting-started"
                    class="block px-4 py-2.5 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-primary hover:bg-slate-50 dark:hover:bg-slate-800/50 rounded-lg transition-colors"
                    onclick={closeMobileMenu}
                >
                    Documentation
                </a>
                <a
                    href="https://kodia.id/features"
                    class="block px-4 py-2.5 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-primary hover:bg-slate-50 dark:hover:bg-slate-800/50 rounded-lg transition-colors"
                    onclick={closeMobileMenu}
                >
                    Features
                </a>

                <div class="h-px bg-slate-200 dark:bg-slate-800 my-2"></div>

                {#if $authStore.isAuthenticated}
                    <a
                        href="/dashboard"
                        class="block px-4 py-2.5 text-sm font-medium btn-primary text-center rounded-lg"
                        onclick={closeMobileMenu}
                    >
                        Dashboard
                    </a>
                {:else}
                    <a
                        href="/login"
                        class="block px-4 py-2.5 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-primary hover:bg-slate-50 dark:hover:bg-slate-800/50 rounded-lg transition-colors"
                        onclick={closeMobileMenu}
                    >
                        Login
                    </a>
                    <a
                        href="/register"
                        class="block px-4 py-2.5 text-sm font-medium btn-primary text-center rounded-lg"
                        onclick={closeMobileMenu}
                    >
                        Register
                    </a>
                {/if}
            </div>
        </div>
    {/if}
</nav>
