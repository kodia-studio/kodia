import { browser } from '$app/environment';

class ThemeStore {
	dark = $state(false);

	init() {
		if (!browser) return;

		// Check localStorage first
		const stored = localStorage.getItem('theme');
		if (stored) {
			this.dark = stored === 'dark';
			this.applyTheme(this.dark);
			return;
		}

		// Fall back to system preference
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		this.dark = prefersDark;
		this.applyTheme(this.dark);
	}

	private applyTheme(isDark: boolean) {
		const html = document.documentElement;
		if (isDark) {
			html.classList.add('dark');
		} else {
			html.classList.remove('dark');
		}
	}

	toggle() {
		if (!browser) return;
		this.dark = !this.dark;
		this.applyTheme(this.dark);
		localStorage.setItem('theme', this.dark ? 'dark' : 'light');
	}
}

export const themeStore = new ThemeStore();

export function initTheme() {
	themeStore.init();
}
