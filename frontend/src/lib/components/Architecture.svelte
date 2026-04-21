<script lang="ts">
    import { Code2, ArrowRight, ShieldCheck, Cpu, Database, Layout } from 'lucide-svelte';
    import { fly } from 'svelte/transition';

    const steps = [
        { title: 'Request', icon: Layout, desc: 'Incoming HTTP/gRPC', color: 'text-blue-500' },
        { title: 'Security', icon: ShieldCheck, desc: 'Auth & CSRF', color: 'text-indigo-500' },
        { title: 'Handler', icon: Code2, desc: 'DTO Validation', color: 'text-amber-500' },
        { title: 'Service', icon: Cpu, desc: 'Business Logic', color: 'text-primary' },
        { title: 'Repository', icon: Database, desc: 'GORM / Cache', color: 'text-emerald-500' }
    ];

    let mounted = $state(false);
    import { onMount } from 'svelte';
    onMount(() => mounted = true);
</script>

<section class="py-32 relative bg-white dark:bg-slate-950 transition-colors duration-500">
    <div class="container mx-auto px-6">
        <div class="flex flex-col lg:flex-row gap-20 items-center">
            
            <!-- Blueprint Diagram -->
            <div class="w-full lg:w-1/2 order-2 lg:order-1">
                <div class="relative glass p-4 sm:p-8 rounded-4xl border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/50">
                    <div class="absolute inset-x-0 -top-px h-px bg-linear-to-r from-transparent via-primary to-transparent opacity-20"></div>
                    
                    <div class="space-y-4">
                        {#if mounted}
                            {#each steps as step, i}
                                <div 
                                    class="flex items-center gap-4 group"
                                    in:fly={{ x: -20, duration: 800, delay: i * 150 }}
                                >
                                    <div class="w-12 h-12 shrink-0 rounded-xl {step.color} bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 shadow-sm flex items-center justify-center transition-transform group-hover:scale-110">
                                        <step.icon size={20} />
                                    </div>
                                    <div class="grow h-px bg-slate-200 dark:bg-slate-800 flex items-center px-4">
                                        {#if i < steps.length - 1}
                                            <div class="w-2 h-2 rounded-full bg-primary animate-ping"></div>
                                        {/if}
                                    </div>
                                    <div class="w-40 shrink-0 text-right">
                                        <h4 class="text-sm font-black text-slate-900 dark:text-white uppercase tracking-widest">{step.title}</h4>
                                        <p class="text-[10px] font-mono text-slate-500 italic mt-0.5">{step.desc}</p>
                                    </div>
                                </div>
                                {#if i < steps.length - 1}
                                    <div class="ml-6 w-px h-8 border-l border-dashed border-slate-300 dark:border-slate-700"></div>
                                {/if}
                            {/each}
                        {/if}
                    </div>

                    <!-- Side Annotations -->
                    <div class="hidden sm:block absolute -right-20 top-1/2 -translate-y-1/2 space-y-24">
                        <div class="glass p-3 rounded-lg border border-slate-200 dark:border-slate-800 text-[9px] uppercase font-bold tracking-widest text-primary rotate-90">
                            Middleware Engine
                        </div>
                        <div class="glass p-3 rounded-lg border border-slate-200 dark:border-slate-800 text-[9px] uppercase font-bold tracking-widest text-amber-500 rotate-90">
                            Dependency Injection
                        </div>
                    </div>
                </div>
            </div>

            <!-- Content -->
            <div class="w-full lg:w-1/2 order-1 lg:order-2">
                <span class="text-xs font-black tracking-[0.3em] text-primary uppercase mb-6 block">The Blueprint</span>
                <h2 class="text-5xl md:text-6xl font-black tracking-tight text-slate-900 dark:text-white mb-8 leading-[0.9]">
                    Artisanal <br/>
                    <span class="text-gradient">Core Logic.</span>
                </h2>
                <div class="space-y-6">
                    <p class="text-lg text-slate-600 dark:text-slate-400 leading-relaxed">
                        Kodia enforces a clean, modular architecture that protects your business 
                        logic from infrastructure changes. It’s the "Hive Way" of building 
                        software that lasts.
                    </p>
                    <ul class="space-y-4">
                        <li class="flex items-center gap-4 text-sm font-bold text-slate-800 dark:text-slate-200">
                            <div class="w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center">
                                <ArrowRight size={14} />
                            </div>
                            Automatic Dependency Injection
                        </li>
                        <li class="flex items-center gap-4 text-sm font-bold text-slate-800 dark:text-slate-200">
                            <div class="w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center">
                                <ArrowRight size={14} />
                            </div>
                            Strongly Typed DTO Protection
                        </li>
                        <li class="flex items-center gap-4 text-sm font-bold text-slate-800 dark:text-slate-200">
                            <div class="w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center">
                                <ArrowRight size={14} />
                            </div>
                            Centralized Service Architecture
                        </li>
                    </ul>
                </div>
            </div>

        </div>
    </div>
</section>
