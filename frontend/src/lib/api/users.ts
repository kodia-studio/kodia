import { api } from './client.svelte';
import type { UserResponse, UpdateUserRequest, ChangePasswordRequest } from '../types/user.types';
import type { ApiResponse } from '../types/api.types';

export const userApi = {
	getAll: async (page = 1, perPage = 15) => {
		return await api.get<UserResponse[]>(`/api/users?page=${page}&per_page=${perPage}`);
	},

	getById: async (id: string) => {
		return await api.get<UserResponse>(`/api/users/${id}`);
	},

	update: async (id: string, data: UpdateUserRequest) => {
		return await api.patch<UserResponse>(`/api/users/${id}`, data);
	},

	delete: async (id: string) => {
		return await api.delete(`/api/users/${id}`);
	},

	changePassword: async (data: ChangePasswordRequest) => {
		return await api.post('/api/users/me/change-password', data);
	}
};
