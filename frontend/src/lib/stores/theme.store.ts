import { browser } from '$app/environment';
import { writable } from 'svelte/store';

type Theme = 'light' | 'dark' | 'system';

function createThemeStore() {
	const { subscribe, set } = writable<Theme>('system');

	return {
		subscribe,
		set: (theme: Theme) => {
			if (browser) {
				localStorage.setItem('theme', theme);
				if (theme === 'dark' || (theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
					document.documentElement.classList.add('dark');
				} else {
					document.documentElement.classList.remove('dark');
				}
			}
			set(theme);
		},
		init: () => {
			if (browser) {
				const theme = (localStorage.getItem('theme') as Theme) || 'system';
				createThemeStore().set(theme);
			}
		}
	};
}

export const themeStore = createThemeStore();
