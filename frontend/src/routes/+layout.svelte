<script lang="ts">
  import "../app.css";
  import { onMount } from "svelte";
  import { initTheme, themeStore } from "$lib/stores/theme.svelte";
  import { Toaster } from "svelte-sonner";
  import { authStore } from "$lib/stores/auth.store";

  let { children } = $props();

  onMount(() => {
    initTheme();
  });

  const theme = $derived(themeStore.dark ? 'dark' : 'light');
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
  <link
    href="https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=Outfit:wght@100..900&display=swap"
    rel="stylesheet"
  />
</svelte:head>

<Toaster {theme} position="top-right" richColors closeButton />

<div class="min-h-screen bg-background text-foreground font-sans selection:bg-primary/20">
  {@render children()}
</div>
