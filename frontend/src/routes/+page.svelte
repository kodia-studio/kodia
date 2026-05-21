<script lang="ts">
    import Hero from '$lib/components/Hero.svelte';
    import Features from '$lib/components/Features.svelte';
    import Ecosystem from '$lib/components/Ecosystem.svelte';
    import Architecture from '$lib/components/Architecture.svelte';
    import { Zap, Book, Server, Shield, Globe } from 'lucide-svelte';
    import SEO from '$lib/components/SEO.svelte';
    import Navbar from '$lib/components/Navbar.svelte';
    import Footer from '$lib/components/Footer.svelte';
</script>

<SEO 
    title="Kodia Framework - Build World-Class Apps with Go & Svelte 5"
    description="The institutional-grade fullstack framework for Go. High-velocity developer experience, zero-config infrastructure, and real-time observability."
    keywords="Go, Svelte 5, Fullstack, Framework, Golang, Web Development, SEO, SGO, Performance"
/>

<Navbar />

<main class="min-h-screen bg-white dark:bg-slate-950 transition-all duration-500 overflow-hidden">
    <Hero />
    <Features />
    <Ecosystem />
    <Architecture />
    
    <!-- Unique 'Hive Core' CTA Section -->
    <section class="pt-12 md:pt-16 lg:pt-20 pb-16 md:pb-24 lg:pb-40 relative isolate">
        <!-- Background Identity Elements -->
        <div class="absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-slate-200 dark:via-slate-800 to-transparent"></div>
        <div class="absolute inset-0 bg-hive opacity-[0.03] -z-10"></div>

        <div class="container mx-auto px-4 md:px-6">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12 lg:gap-20 items-center">

                <!-- Left: Hive Core Visual -->
                <div class="relative order-2 lg:order-1 flex justify-center lg:justify-start">
                    <div class="relative w-48 h-48 sm:w-56 sm:h-56 md:w-72 md:h-72 lg:w-96 lg:h-96 group">
                        <!-- Holographic Orb -->
                        <div class="absolute inset-0 bg-linear-to-tr from-primary to-amber-500 rounded-full blur-[80px] opacity-20 group-hover:opacity-40 transition-opacity duration-1000 animate-pulse"></div>
                        <div class="absolute inset-4 bg-white dark:bg-slate-800 rounded-full border border-slate-200 dark:border-slate-700 shadow-2xl flex items-center justify-center p-4 md:p-6 lg:p-8 overflow-hidden group-hover:scale-105 transition-transform duration-700">
                            <!-- Inner Hive Pattern -->
                            <div class="absolute inset-0 bg-hive opacity-10 dark:opacity-5"></div>
                            <div class="relative z-10 w-16 h-16 sm:w-20 sm:h-20 md:w-24 md:h-24 lg:w-32 lg:h-32 bg-white dark:bg-slate-900 rounded-full flex items-center justify-center text-slate-900 dark:text-white font-black shadow-lg">
                                <img src="/logo.png" alt="Kodia Hive" class="w-10 sm:w-12 md:w-16 h-auto animate-pulse" />
                            </div>
                        </div>

                        <!-- Technical Tags Floating -->
                        <div class="absolute -top-2 sm:-top-3 md:-top-4 right-0 glass px-2 sm:px-3 md:px-4 py-1 md:py-2 rounded-lg md:rounded-xl text-[9px] md:text-[10px] font-black uppercase tracking-widest text-primary border border-primary/20 shadow-xl animate-bounce-slow">
                            Scalable Apps
                        </div>
                        <div class="absolute bottom-6 sm:bottom-8 md:bottom-10 -left-4 sm:-left-5 md:-left-6 glass px-2 sm:px-3 md:px-4 py-1 md:py-2 rounded-lg md:rounded-xl text-[9px] md:text-[10px] font-black uppercase tracking-widest text-amber-500 border border-amber-500/20 shadow-xl animate-bounce-slow" style="animation-delay: 1s">
                            High Performance
                        </div>
                    </div>
                </div>

                <!-- Right: Editorial Content -->
                <div class="order-1 lg:order-2">
                    <span class="inline-flex items-center gap-2 px-2 md:px-3 py-1 mb-4 md:mb-6 text-[10px] font-black tracking-[0.3em] text-primary uppercase bg-primary/5 border border-primary/10 rounded-md">
                        Get Started
                    </span>
                    <h2 class="text-3xl sm:text-5xl md:text-6xl lg:text-8xl font-black mb-6 md:mb-8 text-slate-900 dark:text-white leading-[0.85] tracking-tight">
                        Build Your <br/>
                        <span class="text-gradient">Application.</span>
                    </h2>
                    <p class="text-xs sm:text-sm md:text-lg lg:text-xl text-slate-600 dark:text-slate-400 mb-8 md:mb-12 leading-relaxed max-w-xl">
                        Kodia is designed for developers who value performance and clean code.
                        It's fast, it's opinionated, and it's ready for your next
                        production-grade application.
                    </p>

                    <div class="space-y-6 md:space-y-6">
                        <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 md:gap-6">
                            <a href="https://kodia.id/docs/prologue/getting-started" class="btn-premium py-3 md:py-4 lg:py-5 px-6 md:px-10 lg:px-12 text-xs md:text-base lg:text-xl w-full sm:w-auto shadow-[0_20px_50px_rgba(59,130,246,0.3)] group text-center">
                                Start Building
                            </a>
                            <a href="https://github.com/kodia-studio" target="_blank" class="flex items-center justify-center sm:justify-start gap-2 md:gap-3 font-black text-xs md:text-sm uppercase tracking-widest hover:text-primary transition-colors py-3 md:py-4 group">
                                <span class="p-1.5 md:p-2 glass rounded-lg group-hover:bg-primary/10 transition-colors">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-github md:w-5"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.28 1.15-.28 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/><path d="M9 18c-4.51 2-5-2-7-2"/></svg>
                                </span>
                                Source Code
                            </a>
                        </div>

                        <div class="flex items-center gap-4 md:gap-8 pt-6 md:pt-8 border-t border-slate-200 dark:border-slate-800">
                            <div class="flex flex-col">
                                <span class="text-lg md:text-2xl font-black text-slate-900 dark:text-white leading-none">MIT</span>
                                <span class="text-[9px] font-bold uppercase tracking-widest text-slate-500 mt-1">Open Source</span>
                            </div>
                            <div class="w-px h-6 md:h-8 bg-slate-200 dark:bg-slate-800"></div>
                            <div class="flex flex-col">
                                <span class="text-lg md:text-2xl font-black text-slate-900 dark:text-white leading-none">100%</span>
                                <span class="text-[9px] font-bold uppercase tracking-widest text-slate-500 mt-1">Free to Use</span>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    </section>

</main>

<Footer />
