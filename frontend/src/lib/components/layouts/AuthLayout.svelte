<script lang="ts">
  import { Layers } from "lucide-svelte";
  import { fade } from "svelte/transition";

  let { children } = $props();
</script>

<div class="min-h-screen flex items-center justify-center relative overflow-hidden bg-background">
  <!-- Decorative background elements -->
  <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] rounded-full bg-primary/5 blur-[120px]"></div>
  <div class="absolute -bottom-[10%] -right-[10%] w-[40%] h-[40%] rounded-full bg-primary/10 blur-[120px]"></div>

  <div class="w-full max-w-md px-6 relative z-10">
    <div class="flex flex-col items-center mb-8">
      <a href="/" class="flex items-center gap-2 mb-6 group">
        <div class="w-12 h-12 accent-gradient rounded-2xl flex items-center justify-center shadow-xl group-hover:scale-110 transition-transform">
          <Layers class="text-white w-7 h-7" />
        </div>
        <span class="text-3xl font-bold font-heading tracking-tight">Kodia</span>
      </a>
    </div>

    <div 
      class="card glass p-8 shadow-2xl backdrop-blur-2xl border-white/20"
      in:fade={{ duration: 400, delay: 100 }}
    >
      {@render children?.()}
    </div>

    <p class="text-center mt-8 text-sm text-muted-foreground">
      &copy; {new Date().getFullYear()} Kodia Framework
    </p>
  </div>
</div>
