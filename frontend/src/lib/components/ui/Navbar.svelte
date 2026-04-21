<script lang="ts">
  import { toggleMode, mode } from "mode-watcher";
  import { Sun, Moon, Layers, UserCircle } from "lucide-svelte";
  import { cn } from "$lib/utils/styles";
  import { authStore } from "$lib/stores/auth.store";

  let isScrolled = $state(false);

  if (typeof window !== "undefined") {
    window.addEventListener("scroll", () => {
      isScrolled = window.scrollY > 20;
    });
  }
</script>

<nav
  class={cn(
    "fixed top-0 left-0 right-0 z-40 transition-all duration-300",
    isScrolled ? "py-3 glass shadow-lg" : "py-6 bg-white dark:bg-transparent"
  )}
>
  <div class="container mx-auto px-6 flex items-center justify-between">
    <a href="/" class="flex items-center gap-2 group">
      <div
        class="w-10 h-10 bg-linear-to-br from-blue-500 to-blue-600 dark:from-blue-600 dark:to-cyan-600 rounded-xl flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform"
      >
        <Layers class="text-white w-6 h-6" />
      </div>
      <span class="text-2xl font-bold font-heading tracking-tight text-slate-900 dark:text-white">
        Kodia<span class="text-blue-500 dark:text-blue-400">.</span>
      </span>
    </a>

    <div class="hidden md:flex items-center gap-8 text-sm font-medium text-slate-700 dark:text-slate-300">
      <a href="/docs" class="hover:text-blue-600 dark:hover:text-blue-400 transition-colors">Documentation</a>
      <!-- Examples dan Showcase akan diimplementasikan nanti -->
    </div>

    <div class="flex items-center gap-3">
      <a
        href="https://github.com/kodia-studio/kodia"
        target="_blank"
        class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors hidden sm:block text-slate-700 dark:text-slate-300"
        aria-label="GitHub"
      >
        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v 3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
        </svg>
      </a>

      <button
        onclick={toggleMode}
        class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors overflow-hidden relative w-9 h-9 text-slate-700 dark:text-slate-300"
        aria-label="Toggle theme"
      >
        <div class="flex flex-col items-center justify-center h-full">
          {#if mode.current === "dark"}
            <Moon class="w-5 h-5" />
          {:else}
            <Sun class="w-5 h-5" />
          {/if}
        </div>
      </button>

      <div class="h-8 w-px bg-border/50 mx-1 hidden sm:block"></div>

      {#if $authStore.isAuthenticated}
        <a href="/dashboard" class="btn-primary py-2 px-5 text-sm flex items-center gap-2">
          <UserCircle class="w-4 h-4" />
          Dashboard
        </a>
      {:else}
        <div class="flex items-center gap-2">
          <a href="/login" class="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
            Sign In
          </a>
          <a href="/register" class="btn-primary py-2 px-5 text-sm shadow-xl shadow-primary/20">
            Get Started
          </a>
        </div>
      {/if}
    </div>
  </div>
</nav>
