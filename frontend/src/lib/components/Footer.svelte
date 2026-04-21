<script lang="ts">
    import { Zap, Globe } from 'lucide-svelte';
</script>

<footer class="pt-24 pb-12 bg-white dark:bg-slate-950 border-t border-slate-200 dark:border-slate-800 transition-colors">
    <div class="container mx-auto px-6">
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-12 mb-20">
            <!-- Brand Column -->
            <div class="col-span-2">
                <a href="/" class="flex items-center gap-3 mb-6 group">
                    <div class="w-10 h-10 bg-primary rounded-xl flex items-center justify-center text-white font-black group-hover:rotate-12 transition-transform shadow-lg shadow-primary/20">
                        <Zap size={22} fill="currentColor" />
                    </div>
                    <span class="text-2xl font-black tracking-tighter">KODIA</span>
                </a>
                <p class="text-sm text-slate-500 dark:text-slate-400 max-w-xs leading-relaxed">
                    The world-class fullstack framework for Go. 
                    Designed for speed, built for reliability, 
                    and engineered for developer joy.
                </p>
                <div class="flex gap-5 mt-8 items-center text-slate-400">
                    <a href="https://x.com/kodia_studio" target="_blank" class="hover:text-primary transition-colors" aria-label="Twitter">
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-twitter"><path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"/></svg>
                    </a>
                    <a href="https://github.com/kodia-studio/kodia" target="_blank" class="hover:text-primary transition-colors" aria-label="GitHub">
                        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-github"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.28 1.15-.28 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/><path d="M9 18c-4.51 2-5-2-7-2"/></svg>
                    </a>
                    <a href="https://kodia.id" target="_blank" class="hover:text-primary transition-colors" aria-label="Official Website"><Globe size={18} /></a>
                </div>
            </div>

            <!-- Documentation -->
            <div>
                <h4 class="text-xs font-black uppercase tracking-widest text-slate-900 dark:text-white mb-6">Docs</h4>
                <ul class="space-y-4 text-sm font-medium text-slate-500 dark:text-slate-400">
                    <li><a href="https://kodia.id/docs" class="hover:text-primary transition-colors">Getting Started</a></li>
                    <li><a href="https://kodia.id/docs" class="hover:text-primary transition-colors">Architecture</a></li>
                    <li><a href="https://kodia.id/docs" class="hover:text-primary transition-colors">Basic Usage</a></li>
                    <li><a href="https://kodia.id/docs" class="hover:text-primary transition-colors">Security</a></li>
                </ul>
            </div>

            <!-- Ecosystem -->
            <div>
                <h4 class="text-xs font-black uppercase tracking-widest text-slate-900 dark:text-white mb-6">Ecosystem</h4>
                <ul class="space-y-4 text-sm font-medium text-slate-500 dark:text-slate-400">
                    <li><a href="https://kodia.id/ecosystem/sail" class="hover:text-primary transition-colors">Kodia Sail</a></li>
                    <li><a href="https://kodia.id/ecosystem/events" class="hover:text-primary transition-colors">Event Dispatcher</a></li>
                    <li><a href="https://kodia.id/ecosystem/mailer" class="hover:text-primary transition-colors">Kodia Mailer</a></li>
                    <li><a href="https://kodia.id/ecosystem/tasks" class="hover:text-primary transition-colors">Task Board</a></li>
                </ul>
            </div>

            <!-- Community -->
            <div>
                <h4 class="text-xs font-black uppercase tracking-widest text-slate-900 dark:text-white mb-6">Community</h4>
                <ul class="space-y-4 text-sm font-medium text-slate-500 dark:text-slate-400">
                    <li><a href="https://github.com/kodia-studio/kodia/issues" target="_blank" rel="noopener noreferrer" class="hover:text-primary transition-colors">GitHub Issues</a></li>
                    <li><a href="https://github.com/kodia-studio/kodia/discussions" target="_blank" rel="noopener noreferrer" class="hover:text-primary transition-colors">Discussions</a></li>
                    <li><a href="https://kodia.id/changelog" class="hover:text-primary transition-colors">Changelog</a></li>
                    <li><a href="https://discord.gg/kodia" target="_blank" rel="noopener noreferrer" class="hover:text-primary transition-colors">Discord</a></li>
                </ul>
            </div>

            <!-- Legal -->
            <div>
                <h4 class="text-xs font-black uppercase tracking-widest text-slate-900 dark:text-white mb-6">License</h4>
                <ul class="space-y-4 text-sm font-medium text-slate-500 dark:text-slate-400">
                    <li><a href="https://kodia.id/license" class="hover:text-primary transition-colors">MIT License</a></li>
                    <li><a href="https://kodia.id/privacy" class="hover:text-primary transition-colors">Privacy Policy</a></li>
                </ul>
            </div>
        </div>

        <div class="border-t border-slate-200 dark:border-slate-800 pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
            <p class="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">
                © 2026 Kodia Studio. All Rights Reserved.
            </p>
            <div class="flex items-center gap-2">
                <span class="text-[9px] font-black uppercase tracking-widest text-slate-400">Built with</span>
                <div class="flex items-center gap-1.5 px-2 py-1 bg-primary/5 rounded border border-primary/10">
                    <Zap size={10} class="text-primary fill-primary" />
                    <span class="text-[9px] font-black text-primary uppercase">Kodia Framework</span>
                </div>
            </div>
        </div>
    </div>
</footer>
