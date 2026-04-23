<script lang="ts">
  import Navbar from "../Navbar.svelte";
  import { fade } from "svelte/transition";

  let { children } = $props();
</script>

<div class="relative min-h-screen flex flex-col">
  <Navbar />
  
  <main class="flex-1 pt-24">
    {#key children}
      <div in:fade={{ duration: 300, delay: 150 }} out:fade={{ duration: 150 }}>
        {@render children?.()}
      </div>
    {/key}
  </main>

  <footer class="py-12 border-t bg-muted/30">
    <div class="container mx-auto px-6 text-center text-muted-foreground text-sm">
      <p>&copy; {new Date().getFullYear()} Kodia Framework. Built like a user, code like a pro.</p>
    </div>
  </footer>
</div>
