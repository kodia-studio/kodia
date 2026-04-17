import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
	id: string;
	message: string;
	type: ToastType;
	duration?: number;
}

function createToastStore() {
	const { subscribe, update } = writable<Toast[]>([]);

	return {
		subscribe,
		add: (message: string, type: ToastType = 'info', duration = 3000) => {
			const id = Math.random().toString(36).substring(2, 9);
			update((toasts) => [...toasts, { id, message, type, duration }]);

			if (duration > 0) {
				setTimeout(() => {
					update((toasts) => toasts.filter((t) => t.id !== id));
				}, duration);
			}
		},
		remove: (id: string) => {
			update((toasts) => toasts.filter((t) => t.id !== id));
		},
		success: (msg: string) => createToastStore().add(msg, 'success'),
		error: (msg: string) => createToastStore().add(msg, 'error')
	};
}

export const toastStore = createToastStore();
