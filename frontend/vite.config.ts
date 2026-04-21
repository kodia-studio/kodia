import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit(), tailwindcss()],
	build: {
		rollupOptions: {
			onwarn(warning, warn) {
				// Silence known circular dependency noise from node_modules
				if (
					warning.code === "CIRCULAR_DEPENDENCY" &&
					(warning.message.includes("node_modules") || warning.ids?.some((id) => id.includes("node_modules")))
				) {
					return;
				}
				warn(warning);
			}
		}
	}
});
