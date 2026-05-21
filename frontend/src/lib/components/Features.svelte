<script lang="ts">
    import { Anchor, Radio, Mail, Sprout, Zap, Palette, Server, Shield, Cpu } from 'lucide-svelte';
    import { fade, fly } from 'svelte/transition';

    const features = [
        {
            title: "Kodia Sail",
            description: "One-command Docker infrastructure with Postgres, Redis, and MinIO. Zero configuration required.",
            icon: Anchor,
            span: "col-span-12 lg:col-span-8",
            color: "text-blue-500",
            bg: "bg-blue-500/5",
            tag: "Infrastructure"
        },
        {
            title: "Event Dispatcher",
            description: "Built-in Pub/Sub system with automatic async worker orchestration.",
            icon: Radio,
            span: "col-span-12 lg:col-span-4",
            color: "text-purple-500",
            bg: "bg-purple-500/5",
            tag: "Realtime"
        },
        {
            title: "Svelte 5 UI",
            description: "Next-gen reactive UI with Tailwind v4. The fastest way to build modern frontends for Go applications.",
            icon: Palette,
            span: "col-span-12 lg:col-span-7",
            color: "text-rose-500",
            bg: "bg-rose-500/5",
            tag: "Frontend"
        },
        {
            title: "Kodia Mailer",
            description: "Beautiful HTML emails with multi-provider support.",
            icon: Mail,
            span: "col-span-12 lg:col-span-5",
            color: "text-amber-500",
            bg: "bg-amber-500/5",
            tag: "Communication"
        },
        {
            title: "High Performance",
            description: "Built on Gin and GORM with advanced optimizations.",
            icon: Zap,
            span: "col-span-12 lg:col-span-4",
            color: "text-emerald-500",
            bg: "bg-emerald-500/5",
            tag: "Core"
        },
        {
            title: "Security First",
            description: "CSRF, JWT, and Session management out of the box.",
            icon: Shield,
            span: "col-span-12 lg:col-span-4",
            color: "text-indigo-500",
            bg: "bg-indigo-500/5",
            tag: "Security"
        },
        {
            title: "Background Tasks",
            description: "Redis-backed queue workers with retry logic.",
            icon: Cpu,
            span: "col-span-12 lg:col-span-4",
            color: "text-cyan-500",
            bg: "bg-cyan-500/5",
            tag: "Worker"
        }
    ];

    let mounted = $state(false);
    import { onMount } from 'svelte';
    onMount(() => mounted = true);
</script>

<section class="pt-12 md:pt-16 lg:pt-20 pb-16 md:pb-24 lg:pb-32 relative overflow-hidden transition-colors duration-500">
    <!-- Background patterns -->
    <div class="absolute inset-0 bg-hive opacity-5 -z-10"></div>

    <div class="container mx-auto px-4 md:px-6">
        <div class="flex flex-col gap-6 md:gap-8 lg:gap-0 lg:flex-row lg:items-end lg:justify-between mb-12 md:mb-16 lg:mb-20">
            <div class="text-left max-w-2xl">
                <span class="text-[10px] md:text-xs font-black tracking-[0.3em] text-primary uppercase mb-3 md:mb-4 block">Engineered for Velocity</span>
                <h2 class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-black tracking-tight text-slate-900 dark:text-white leading-[0.9]">
                    Key <br/>
                    <span class="text-gradient">Features.</span>
                </h2>
            </div>
            <p class="text-slate-500 dark:text-slate-400 max-w-sm text-xs md:text-sm font-medium leading-relaxed">
                Kodia isn't just a library. It's a cohesive ecosystem designed to solve everyday fullstack problems with Go.
            </p>
        </div>

        <div class="grid grid-cols-12 gap-4 md:gap-5 lg:gap-6">
            {#if mounted}
                {#each features as feature, i}
                    <div
                        class="{feature.span} group relative"
                        in:fly={{ y: 20, duration: 800, delay: i * 100 }}
                    >
                        <div class="h-full glass rounded-2xl md:rounded-3xl p-4 md:p-6 lg:p-8 border border-slate-200 dark:border-slate-800 transition-all duration-500 hover:shadow-2xl hover:shadow-primary/5 hover:-translate-y-1 group-hover:bg-white/50 dark:group-hover:bg-white/5 flex flex-col justify-between">
                            <div>
                                <div class="flex items-center justify-between mb-6 md:mb-8 lg:mb-10">
                                    <div class="w-10 h-10 md:w-12 md:h-12 lg:w-14 lg:h-14 {feature.bg} {feature.color} rounded-lg md:rounded-2xl flex items-center justify-center transition-transform duration-500 group-hover:scale-110 group-hover:rotate-6 shrink-0">
                                        <feature.icon size={20} class="md:hidden" />
                                        <feature.icon size={24} class="hidden md:block lg:hidden" />
                                        <feature.icon size={28} class="hidden lg:block" />
                                    </div>
                                    <span class="text-[9px] md:text-[10px] font-black uppercase tracking-widest opacity-30 group-hover:opacity-100 transition-opacity">
                                        {feature.tag}
                                    </span>
                                </div>
                                <h3 class="text-lg md:text-xl lg:text-2xl font-black mb-2 md:mb-3 lg:mb-4 text-slate-900 dark:text-white tracking-tight">
                                    {feature.title}
                                </h3>
                                <p class="text-slate-600 dark:text-slate-400 text-xs md:text-sm leading-relaxed max-w-xs">
                                    {feature.description}
                                </p>
                            </div>

                            <div class="mt-6 md:mt-7 lg:mt-8 flex items-center gap-2 text-[9px] md:text-[10px] font-bold uppercase tracking-tighter text-primary opacity-0 group-hover:opacity-100 transition-all -translate-x-2.5 group-hover:translate-x-0">
                                Learn more about {feature.title.split(' ')[1] || feature.title}
                                <Zap size={10} fill="currentColor" />
                            </div>
                        </div>
                    </div>
                {/each}
            {/if}
        </div>
    </div>
</section>
