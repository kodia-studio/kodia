<script lang="ts">
    import { onMount } from 'svelte';
    import { fade, fly } from 'svelte/transition';
    import { ArrowRight, BookOpen, Zap, CheckCircle2 } from 'lucide-svelte';

    let mounted = $state(false);
    onMount(() => {
        mounted = true;
    });
</script>

<section class="relative pt-32 pb-20 overflow-hidden min-h-screen flex items-center">
    <!-- Artisanal Background Elements -->
    <div class="absolute inset-0 bg-hive -z-10 animate-pulse-slow"></div>
    <div class="hero-glow"></div>
    
    <!-- Floating Hive Blobs -->
    <div class="absolute top-1/4 -right-20 w-96 h-96 bg-primary/20 rounded-full blur-[120px] animate-blob"></div>
    <div class="absolute bottom-1/4 -left-20 w-80 h-80 bg-amber-500/10 rounded-full blur-[100px] animate-blob animation-delay-2000"></div>

    <div class="container mx-auto px-6 relative z-10">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
            
            <!-- LEFT: CLI Terminal (Identity focus) -->
            {#if mounted}
            <div class="order-2 lg:order-1" in:fly={{ x: -40, duration: 1000, delay: 200 }}>
                <div class="relative group">
                    <!-- Layered perspective effect -->
                    <div class="absolute -inset-1 bg-linear-to-r from-primary to-amber-500 rounded-3xl blur opacity-25 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"></div>
                    
                    <div class="relative rounded-2xl overflow-hidden border border-slate-800 shadow-[0_22px_70px_4px_rgba(0,0,0,0.56)] bg-[#0b1120]">
                        <!-- Terminal Header -->
                        <div class="bg-slate-900/80 px-5 py-3 flex items-center justify-between border-b border-white/5">
                            <div class="flex gap-2">
                                <div class="w-3 h-3 rounded-full bg-rose-500/80"></div>
                                <div class="w-3 h-3 rounded-full bg-amber-500/80"></div>
                                <div class="w-3 h-3 rounded-full bg-emerald-500/80"></div>
                            </div>
                            <div class="flex items-center gap-2">
                                <img src="/logo.png" alt="Kodia CLI" class="h-3 w-auto" />
                                <span class="text-[10px] text-slate-500 font-bold font-mono uppercase tracking-tighter">Kodia CLI v1.0</span>
                            </div>
                        </div>

                        <!-- Terminal Content -->
                        <div class="p-8 text-left font-mono text-base leading-relaxed overflow-x-auto whitespace-pre">
                            <div class="flex gap-3">
                                <span class="text-primary font-bold">➜</span>
                                <span class="text-amber-400">~</span>
                                <span class="text-white">kodia new my-app</span>
                            </div>
                            <div class="mt-4 space-y-1">
                                <div class="flex gap-3 text-slate-500 italic">
                                    <span>[kodia]</span>
                                    <span>Building your colony...</span>
                                </div>
                                <div class="flex gap-3 text-emerald-400/90 items-center">
                                    <CheckCircle2 size={14} />
                                    <span>Backend: Go 1.22 + Gin + GORM</span>
                                </div>
                                <div class="flex gap-3 text-emerald-400/90 items-center">
                                    <CheckCircle2 size={14} />
                                    <span>Frontend: Svelte 5 + Tailwind v4</span>
                                </div>
                                <div class="flex gap-3 text-emerald-400/90 items-center">
                                    <CheckCircle2 size={14} />
                                    <span>Infra: Docker + Kodia Sail ready</span>
                                </div>
                                <div class="mt-4 flex gap-3">
                                    <span class="text-primary font-bold">➜</span>
                                    <span class="text-amber-400">my-app</span>
                                    <span class="text-white">kodia sail up</span>
                                </div>
                                <div class="text-amber-400/80 animate-pulse">⚡ Hive is active at http://localhost:8080</div>
                            </div>
                        </div>
                    </div>

                    <!-- Floating Badge -->
                    <div class="absolute -bottom-12 -right-6 glass p-4 rounded-2xl shadow-xl border border-white/10 animate-bounce-slow hidden sm:block">
                        <div class="flex items-center gap-3">
                            <div class="w-10 h-10 bg-amber-500/10 rounded-xl flex items-center justify-center text-amber-500">
                                <Zap size={20} fill="currentColor" />
                            </div>
                            <div>
                                <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Performance</p>
                                <p class="text-sm font-black text-slate-900 dark:text-white">0.4ms Response</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            {/if}

            <!-- RIGHT: Bold Content -->
            {#if mounted}
            <div class="order-1 lg:order-2 text-left" in:fly={{ x: 40, duration: 1000, delay: 400 }}>
                <h1 class="text-7xl md:text-8xl font-black tracking-tight leading-[0.85] mb-8 text-slate-900 dark:text-white">
                    Built for <br/>
                    <span class="text-gradient">The Sting.</span>
                </h1>
                
                <p class="text-lg md:text-xl text-slate-600 dark:text-slate-400 mb-12 leading-relaxed max-w-xl">
                    Kodia isn't just another framework. It's a high-velocity 
                    fullstack ecosystem for Go developers who prioritize 
                    architecture, speed, and developer joy.
                </p>

                <div class="flex flex-col sm:flex-row items-center gap-6">
                    <a href="https://kodia.id/docs/prologue/getting-started" class="btn-premium py-4 px-10 text-lg w-full sm:w-auto flex items-center justify-center gap-3 group">
                        Join the Colony
                        <ArrowRight size={20} class="group-hover:translate-x-1 transition-transform" />
                    </a>
                    <button class="flex items-center gap-3 font-black text-sm uppercase tracking-widest hover:text-primary transition-colors py-4">
                        <BookOpen size={18} class="text-amber-500" />
                        Explore Docs
                    </button>
                </div>

                <!-- Trust Badges / Stats -->
                <div class="mt-16 pt-12 border-t border-slate-200/50 dark:border-slate-800/50 flex flex-wrap gap-12">
                    <div>
                        <p class="text-3xl font-black mb-1">20x</p>
                        <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Development Speed</p>
                    </div>
                    <div>
                        <p class="text-3xl font-black mb-1">Go</p>
                        <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Modern 1.22 Runtime</p>
                    </div>
                    <div>
                        <p class="text-3xl font-black mb-1">Svelte</p>
                        <p class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">Next-gen Frontend</p>
                    </div>
                </div>
            </div>
            {/if}

        </div>
    </div>
</section>

<style>
    .animate-blob {
        animation: blob 7s infinite;
    }
    .animation-delay-2000 {
        animation-delay: 2s;
    }
    .animate-pulse-slow {
        animation: pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite;
    }
    .animate-bounce-slow {
        animation: bounce 3s infinite;
    }

    @keyframes blob {
        0% { transform: translate(0px, 0px) scale(1); }
        33% { transform: translate(30px, -50px) scale(1.1); }
        66% { transform: translate(-20px, 20px) scale(0.9); }
        100% { transform: translate(0px, 0px) scale(1); }
    }

    @keyframes bounce {
        0%, 100% { transform: translateY(-5%); animation-timing-function: cubic-bezier(0.8,0,1,1); }
        50% { transform: translateY(0); animation-timing-function: cubic-bezier(0,0,0.2,1); }
    }
</style>
