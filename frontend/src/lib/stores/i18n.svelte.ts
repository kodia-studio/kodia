/**
 * Kodia i18n Engine 🐨🌍
 * High-performance, reactive translation layer using Svelte 5 Runes.
 */

class I18nStore {
	locale = $state('en');
	dictionary = $state<Record<string, string>>({});
	isLoading = $state(false);

	/**
	 * Reactive translation function
	 */
	t = (key: string, vars: Record<string, string | number> = {}) => {
		let text = this.dictionary[key] || key;
		
		// Interpolation
		Object.entries(vars).forEach(([k, v]) => {
			text = text.replace(`{${k}}`, String(v));
		});
		
		return text;
	};

	/**
	 * Set locale and load dictionary
	 */
	async setLocale(newLocale: string) {
		this.isLoading = true;
		try {
			// In a real app, this would fetch from a JSON file or API
			// For now, we'll use a local mock loader
			const dict = await this.loadDictionary(newLocale);
			this.dictionary = dict;
			this.locale = newLocale;
			document.documentElement.lang = newLocale;
		} catch (err) {
			console.error(`[i18n] Failed to load locale: ${newLocale}`, err);
		} finally {
			this.isLoading = false;
		}
	}

	private async loadDictionary(locale: string) {
		// Mock lazy loading
		const dictionaries: Record<string, any> = {
			en: {
				"app.title": "Kodia Framework",
				"app.tagline": "Build like a user, code like a pro.",
				"auth.login": "Login to Intelligence",
				"auth.register": "Join Ecosystem",
				"common.search": "Search...",
				"common.loading": "Processing..."
			},
			id: {
				"app.title": "Kodia Framework",
				"app.tagline": "Bangun seperti user, koding seperti pro.",
				"auth.login": "Masuk ke Sistem",
				"auth.register": "Gabung Ekosistem",
				"common.search": "Cari...",
				"common.loading": "Memproses..."
			}
		};
		
		return dictionaries[locale] || dictionaries.en;
	}
}

export const i18n = new I18nStore();
