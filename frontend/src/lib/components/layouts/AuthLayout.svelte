<script lang="ts">
  import { Zap } from "lucide-svelte";
  import { fade, fly } from "svelte/transition";

  let { children } = $props();
</script>

<div class="min-h-screen flex items-center justify-center relative overflow-hidden bg-white dark:bg-slate-950 transition-colors duration-500">
  <!-- Artisanal Background Architecture -->
  <div class="absolute inset-0 bg-hive opacity-[0.03] dark:opacity-[0.05] -z-10"></div>
  <div class="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full bg-linear-to-b from-primary/10 via-transparent to-transparent opacity-20 dark:opacity-40"></div>
  
  <!-- Floating Identity Blobs -->
  <div class="absolute -top-[20%] -left-[10%] w-[60%] h-[60%] rounded-full bg-primary/5 blur-[120px] animate-pulse"></div>
  <div class="absolute -bottom-[20%] -right-[10%] w-[60%] h-[60%] rounded-full bg-amber-500/5 blur-[120px] animate-pulse" style="animation-duration: 4s"></div>

  <div class="w-full max-w-[440px] px-6 relative z-10 py-12">
    <!-- Brand Identity -->
    <div class="flex flex-col items-center mb-10" in:fly={{ y: -20, duration: 800 }}>
      <a href="/" class="flex items-center gap-3 group">
        <div class="w-14 h-14 bg-primary rounded-2xl flex items-center justify-center text-white shadow-2xl shadow-primary/30 group-hover:rotate-12 transition-transform duration-500">
          <Zap size={32} fill="currentColor" />
        </div>
        <span class="text-4xl font-black tracking-tighter text-slate-900 dark:text-white uppercase">Kodia</span>
      </a>
    </div>

    <!-- Main Auth Stage -->
    <div 
      class="glass relative rounded-[40px] p-8 md:p-10 shadow-[0_32px_100px_-20px_rgba(0,0,0,0.1)] dark:shadow-[0_32px_100px_-20px_rgba(0,0,0,0.5)] border border-slate-200 dark:border-slate-800 bg-white/40 dark:bg-slate-900/40 backdrop-blur-3xl overflow-hidden"
      in:fly={{ y: 20, duration: 800, delay: 200 }}
    >
      <!-- Premium Accent Line -->
      <div class="absolute top-0 inset-x-0 h-px bg-linear-to-r from-transparent via-primary/50 to-transparent"></div>
      
      <div class="relative z-10">
        {@render children?.()}
      </div>
    </div>

    <!-- Footer Attribution -->
    <div class="text-center mt-10" in:fade={{ delay: 1000 }}>
      <p class="text-[10px] font-black uppercase tracking-[0.3em] text-slate-400 dark:text-slate-500">
        &copy; {new Date().getFullYear()} Kodia Studio • Professional Colony
      </p>
    </div>
  </div>
</div>

<style>
  :global(.glass) {
    background: rgba(255, 255, 255, 0.4);
  }
  :global(.dark .glass) {
    background: rgba(15, 23, 42, 0.4);
  }
</style>
