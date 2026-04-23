<script lang="ts">
    import { Database, Mail, Zap, Cloud, Shield, Activity, Share2 } from 'lucide-svelte';
    import { fade, fly } from 'svelte/transition';

    const nodes = [
        { id: 'core', label: 'Kodia Core', icon: Zap, class: 'bg-primary text-white shadow-primary/40' },
        { id: 'sail', label: 'Kodia Sail', icon: Cloud, class: 'bg-blue-500/10 text-blue-500' },
        { id: 'mail', label: 'Kodia Mailer', icon: Mail, class: 'bg-amber-500/10 text-amber-500' },
        { id: 'db', label: 'GORM + Postgres', icon: Database, class: 'bg-emerald-500/10 text-emerald-500' },
        { id: 'auth', label: 'Identity Auth', icon: Shield, class: 'bg-indigo-500/10 text-indigo-500' },
        { id: 'worker', label: 'Queue Worker', icon: Activity, class: 'bg-cyan-500/10 text-cyan-500' }
    ];

    let mounted = $state(false);
    import { onMount } from 'svelte';
    onMount(() => mounted = true);
</script>

<section class="py-32 bg-slate-900 border-y border-white/5 relative overflow-hidden">
    <!-- Animated background grid -->
    <div class="absolute inset-0 bg-hive opacity-10 pointer-events-none"></div>
    
    <div class="container mx-auto px-6 relative z-10">
        <div class="text-center mb-24">
            <span class="text-xs font-black tracking-[0.4em] text-amber-500 uppercase mb-4 block">The Ecosystem</span>
            <h2 class="text-5xl md:text-6xl font-black tracking-tight text-white mb-6">
                Connect the <span class="text-gradient">Hive.</span>
            </h2>
            <p class="text-slate-400 max-w-xl mx-auto text-lg leading-relaxed">
                Kodia components are tightly integrated but loosely coupled, 
                forming a robust colony for your application.
            </p>
        </div>

        <div class="relative h-[600px] flex items-center justify-center">
            <!-- Connection Lines (SVG) -->
            <svg class="absolute inset-0 w-full h-full pointer-events-none opacity-20">
                <defs>
                    <linearGradient id="line-grad" x1="0%" y1="0%" x2="100%" y2="0%">
                        <stop offset="0%" stop-color="#3b82f6" />
                        <stop offset="100%" stop-color="#fbbf24" />
                    </linearGradient>
                </defs>
                {#if mounted}
                    <!-- Simplified connection lines representing the hive -->
                    <circle cx="50%" cy="50%" r="200" fill="none" stroke="url(#line-grad)" stroke-width="1" stroke-dasharray="10 20" class="animate-spin-slow" />
                    <line x1="50%" y1="50%" x2="30%" y2="20%" stroke="url(#line-grad)" stroke-width="2" />
                    <line x1="50%" y1="50%" x2="70%" y2="20%" stroke="url(#line-grad)" stroke-width="2" />
                    <line x1="50%" y1="50%" x2="80%" y2="50%" stroke="url(#line-grad)" stroke-width="2" />
                    <line x1="50%" y1="50%" x2="70%" y2="80%" stroke="url(#line-grad)" stroke-width="2" />
                    <line x1="50%" y1="50%" x2="30%" y2="80%" stroke="url(#line-grad)" stroke-width="2" />
                    <line x1="50%" y1="50%" x2="20%" y2="50%" stroke="url(#line-grad)" stroke-width="2" />
                {/if}
            </svg>

            <!-- Center Hub -->
            <div class="relative group">
                <div class="absolute inset-0 bg-primary/40 blur-[80px] rounded-full group-hover:bg-primary/60 transition-all duration-1000"></div>
                <div class="relative w-32 h-32 md:w-40 md:h-40 {nodes[0].class} rounded-full flex flex-col items-center justify-center text-center p-4 border border-white/20 shadow-2xl z-20 transition-transform duration-500 scale-100 group-hover:scale-105">
                    <img src="/logo.png" alt="Kodia Core" class="h-12 w-auto mb-2" />
                    <span class="text-xs font-black uppercase tracking-tighter">Kodia Core</span>
                </div>
            </div>

            <!-- Satellite Nodes (Fixed positions for art-direction) -->
            <div class="absolute top-[10%] left-[20%] lg:left-[25%] node group">
                <div class="glass-dark p-4 rounded-2xl border border-white/5 flex items-center gap-4 shadow-2xl transition-all group-hover:scale-110 group-hover:border-blue-500/50">
                    <div class="w-10 h-10 rounded-xl bg-blue-500/10 text-blue-500 flex items-center justify-center">
                        <Cloud size={20} />
                    </div>
                    <span class="text-sm font-bold text-white whitespace-nowrap">Kodia Sail</span>
                </div>
            </div>

            <div class="absolute top-[10%] right-[20%] lg:right-[25%] node group">
                <div class="glass-dark p-4 rounded-2xl border border-white/5 flex items-center gap-4 shadow-2xl transition-all group-hover:scale-110 group-hover:border-amber-500/50">
                    <div class="w-10 h-10 rounded-xl bg-amber-500/10 text-amber-500 flex items-center justify-center">
                        <Mail size={20} />
                    </div>
                    <span class="text-sm font-bold text-white whitespace-nowrap">Kodia Mailer</span>
                </div>
            </div>

            <div class="absolute bottom-[10%] left-[20%] lg:left-[25%] node group">
                <div class="glass-dark p-4 rounded-2xl border border-white/5 flex items-center gap-4 shadow-2xl transition-all group-hover:scale-110 group-hover:border-emerald-500/50">
                    <div class="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-500 flex items-center justify-center">
                        <Database size={20} />
                    </div>
                    <span class="text-sm font-bold text-white whitespace-nowrap">GORM + Postgres</span>
                </div>
            </div>

            <div class="absolute bottom-[10%] right-[20%] lg:right-[25%] node group">
                <div class="glass-dark p-4 rounded-2xl border border-white/5 flex items-center gap-4 shadow-2xl transition-all group-hover:scale-110 group-hover:border-indigo-500/50">
                    <div class="w-10 h-10 rounded-xl bg-indigo-500/10 text-indigo-500 flex items-center justify-center">
                        <Shield size={20} />
                    </div>
                    <span class="text-sm font-bold text-white whitespace-nowrap">Identity Auth</span>
                </div>
            </div>

            <div class="absolute left-[5%] lg:left-[10%] top-1/2 -translate-y-1/2 node group hidden md:block">
                <div class="glass-dark p-4 rounded-2xl border border-white/5 flex items-center gap-4 shadow-2xl transition-all group-hover:scale-110 group-hover:border-cyan-500/50">
                    <div class="w-10 h-10 rounded-xl bg-cyan-500/10 text-cyan-500 flex items-center justify-center">
                        <Activity size={20} />
                    </div>
                    <span class="text-sm font-bold text-white whitespace-nowrap">Queue Worker</span>
                </div>
            </div>

            <div class="absolute right-[5%] lg:right-[10%] top-1/2 -translate-y-1/2 node group hidden md:block">
                <div class="glass-dark p-4 rounded-2xl border border-white/5 flex items-center gap-4 shadow-2xl transition-all group-hover:scale-110 group-hover:border-primary/50">
                    <div class="w-10 h-10 rounded-xl bg-primary/10 text-primary flex items-center justify-center">
                        <Share2 size={20} />
                    </div>
                    <span class="text-sm font-bold text-white whitespace-nowrap">Event Dispatcher</span>
                </div>
            </div>
        </div>

        <div class="mt-20 glass-dark p-12 rounded-[40px] border border-white/5 text-center relative overflow-hidden">
            <div class="absolute inset-0 bg-primary/5 animate-pulse"></div>
            <h3 class="text-2xl font-black text-white mb-4 relative z-10">Fully Automated, Zero Glue.</h3>
            <p class="text-slate-400 max-w-2xl mx-auto text-sm leading-relaxed relative z-10">
                Stop wasting time wiring up services. Kodia detects your infrastructure 
                components and registers them automatically with the container, 
                letting you focus on writing code that matters.
            </p>
        </div>
    </div>
</section>

<style>
    .glass-dark {
        background: rgba(15, 23, 42, 0.8);
        backdrop-filter: blur(20px);
        -webkit-backdrop-filter: blur(20px);
    }
    .animate-spin-slow {
        animation: spin 30s linear infinite;
    }
    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
</style>
