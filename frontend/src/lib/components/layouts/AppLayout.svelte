<script lang="ts">
  import Navbar from "../Navbar.svelte";
  import DevTools from "../dev/DevTools.svelte";
  import ErrorBoundary from "../ui/ErrorBoundary.svelte";
  import { i18n } from "$lib/stores/i18n.svelte";
  import { themeStore } from "$lib/stores/theme.svelte";
  import { Toaster } from "svelte-sonner";
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";

  let { children } = $props();

  onMount(() => {
    i18n.setLocale('en');
  });

  const theme = $derived(themeStore.dark ? 'dark' : 'light');
</script>

<div class="relative min-h-screen flex flex-col">
  <Navbar />
  
  <main class="flex-1 pt-24">
    <div class="container mx-auto px-6">
      <ErrorBoundary>
        {#key children}
          <div in:fade={{ duration: 300, delay: 150 }} out:fade={{ duration: 150 }}>
            {@render children?.()}
          </div>
        {/key}
      </ErrorBoundary>
    </div>
  </main>

  <footer class="py-12 border-t bg-muted/30">
    <div class="container mx-auto px-6 text-center text-muted-foreground text-sm">
      <p>&copy; {new Date().getFullYear()} Kodia Framework. Built like a user, code like a pro.</p>
    </div>
  </footer>

  <Toaster {theme} richColors position="top-right" />
  <DevTools />
</div>
