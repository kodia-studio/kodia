import { api } from './client.svelte';
import type { AuthResponse, LoginRequest, RegisterRequest } from '../types/auth.types';
import type { ApiResponse } from '../types/api.types';
import { authStore } from '../stores/auth.store';

export const authApi = {
	register: async (data: RegisterRequest) => {
		const res = await api.post<AuthResponse>('/auth/register', data);
		authStore.login(res.user, res.access_token);
		return res;
	},

	login: async (data: LoginRequest) => {
		const res = await api.post<AuthResponse>('/auth/login', data);
		authStore.login(res.user, res.access_token);
		return res;
	},

	logout: async (refreshToken: string) => {
		try {
			await api.post('/auth/logout', { refresh_token: refreshToken });
		} finally {
			authStore.logout();
		}
	},

	getMe: async () => {
		return await api.get('/auth/me');
	}
};
