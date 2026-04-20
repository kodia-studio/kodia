import { PUBLIC_API_URL } from '$env/static/public';
import type { ApiResponse, ApiError } from '../types/api.types';
import { authStore } from '../stores/auth.store';
import { get } from 'svelte/store';

class ApiClient {
	private baseUrl: string;
	isLoading = $state(false);
	error = $state<ApiError | null>(null);

	constructor() {
		this.baseUrl = PUBLIC_API_URL || 'http://localhost:8080/api';
	}

	private async request<T>(
		path: string,
		options: RequestInit = {}
	): Promise<T> {
		this.isLoading = true;
		this.error = null;

		const url = `${this.baseUrl}${path}`;
		const headers = new Headers(options.headers);

		if (!headers.has('Content-Type')) {
			headers.set('Content-Type', 'application/json');
		}

		// Attach access token if exists
		const token = get(authStore).accessToken;
		if (token) {
			headers.set('Authorization', `Bearer ${token}`);
		}

		try {
			const response = await fetch(url, {
				...options,
				headers
			});

			const result = await response.json();

			if (!response.ok) {
				const apiError = result as ApiError;
				this.error = apiError;
				throw apiError;
			}

			return (result as ApiResponse<T>).data;
		} catch (err) {
			if (!this.error) {
				this.error = {
					error: 'Network Error',
					message: err instanceof Error ? err.message : 'Unknown error occurred'
				};
			}
			throw this.error;
		} finally {
			this.isLoading = false;
		}
	}

	get<T>(path: string, options?: RequestInit) {
		return this.request<T>(path, { ...options, method: 'GET' });
	}

	post<T>(path: string, body?: any, options?: RequestInit) {
		return this.request<T>(path, {
			...options,
			method: 'POST',
			body: JSON.stringify(body)
		});
	}

	put<T>(path: string, body?: any, options?: RequestInit) {
		return this.request<T>(path, {
			...options,
			method: 'PUT',
			body: JSON.stringify(body)
		});
	}

	patch<T>(path: string, body?: any, options?: RequestInit) {
		return this.request<T>(path, {
			...options,
			method: 'PATCH',
			body: JSON.stringify(body)
		});
	}

	delete<T>(path: string, options?: RequestInit) {
		return this.request<T>(path, { ...options, method: 'DELETE' });
	}
}

export const api = new ApiClient();
