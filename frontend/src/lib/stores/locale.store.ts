import { browser } from "$app/environment";
import { api } from "$lib/api/client";

/**
 * Kodia Locale Store (Svelte 5)
 * Reactive store for handling internationalization on the frontend.
 */

// State using Svelte 5 runes
let currentLocale = $state(browser ? localStorage.getItem("kodia_locale") || "en" : "en");
let translations = $state<Record<string, string>>({});
let isLoading = $state(false);

/**
 * Loads translation bundle from the server/static folder.
 */
async function loadTranslations(locale: string) {
  isLoading = true;
  try {
    const response = await fetch(`/locales/${locale}.json`);
    if (response.ok) {
      translations = await response.json();
      currentLocale = locale;
      if (browser) localStorage.setItem("kodia_locale", locale);
    }
  } catch (error) {
    console.error("Failed to load translations:", error);
  } finally {
    isLoading = false;
  }
}

/**
 * T translates a key with optional replacements.
 */
export function t(key: string, replacements: Record<string, string> = {}) {
  let text = translations[key] || key;
  
  for (const [k, v] of Object.entries(replacements)) {
    text = text.replace(`{${k}}`, v);
  }
  
  return text;
}

export const localeStore = {
  get current() { return currentLocale; },
  get translations() { return translations; },
  get isLoading() { return isLoading; },
  set: loadTranslations,
};

// Initialize if in browser
if (browser) {
  loadTranslations(currentLocale);
}
