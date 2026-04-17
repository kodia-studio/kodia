import { browser } from '$app/environment';
import { writable } from 'svelte/store';
import type { UserResponse } from '../types/user.types';

interface AuthState {
	user: UserResponse | null;
	accessToken: string | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}

const initialState: AuthState = {
	user: null,
	accessToken: null,
	isAuthenticated: false,
	isLoading: true
};

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);

	return {
		subscribe,
		set,
		update,
		login: (user: UserResponse, token: string) => {
			if (browser) {
				localStorage.setItem('access_token', token);
				localStorage.setItem('user', JSON.stringify(user));
			}
			set({ user, accessToken: token, isAuthenticated: true, isLoading: false });
		},
		logout: () => {
			if (browser) {
				localStorage.removeItem('access_token');
				localStorage.removeItem('user');
			}
			set({ ...initialState, isLoading: false });
		},
		init: () => {
			if (browser) {
				const token = localStorage.getItem('access_token');
				const userJson = localStorage.getItem('user');
				if (token && userJson) {
					try {
						const user = JSON.parse(userJson);
						set({ user, accessToken: token, isAuthenticated: true, isLoading: false });
						return;
					} catch {
						localStorage.removeItem('access_token');
						localStorage.removeItem('user');
					}
				}
			}
			set({ ...initialState, isLoading: false });
		}
	};
}

export const authStore = createAuthStore();
