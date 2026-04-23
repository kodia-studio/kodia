import { PUBLIC_API_URL } from '$env/static/public';
import type { paths } from '../types/generated/api';
import type { ApiError } from '../types/api.types';
import { authStore } from '../stores/auth.store';
import { get } from 'svelte/store';

type HttpMethod = 'get' | 'post' | 'put' | 'patch' | 'delete';

class ApiClient {
	private baseUrl: string;
	isLoading = $state(false);
	error = $state<ApiError | null>(null);

	constructor() {
		this.baseUrl = PUBLIC_API_URL || 'http://localhost:8080/api';
	}

	/**
	 * Institutional-grade request handler with full type safety
	 */
	async request<
		P extends keyof paths,
		M extends HttpMethod & keyof paths[P]
	>(
		path: P,
		method: M,
		options: {
			params?: paths[P][M] extends { parameters: { query?: infer Q, path?: infer PH } } ? { query?: Q, path?: PH } : never;
			body?: paths[P][M] extends { requestBody: { content: { "application/json": infer B } } } ? B : never;
			headers?: Record<string, string>;
		} = {} as any
	): Promise<any> {
		this.isLoading = true;
		this.error = null;

		let url = `${this.baseUrl}${path as string}`;
		
		// Handle path parameters
		if (options.params?.path) {
			Object.entries(options.params.path).forEach(([key, value]) => {
				url = url.replace(`{${key}}`, String(value));
			});
		}

		// Handle query parameters
		if (options.params?.query) {
			const searchParams = new URLSearchParams();
			Object.entries(options.params.query).forEach(([key, value]) => {
				if (value !== undefined && value !== null) {
					searchParams.append(key, String(value));
				}
			});
			const queryString = searchParams.toString();
			if (queryString) url += `?${queryString}`;
		}

		const headers = new Headers(options.headers);
		if (!headers.has('Content-Type')) {
			headers.set('Content-Type', 'application/json');
		}

		const token = get(authStore).accessToken;
		if (token) {
			headers.set('Authorization', `Bearer ${token}`);
		}

		try {
			const response = await fetch(url, {
				method: method.toUpperCase(),
				headers,
				body: options.body ? JSON.stringify(options.body) : undefined
			});

			const result = await response.json();

			if (!response.ok) {
				const apiError = result as ApiError;
				this.error = apiError;
				throw apiError;
			}

			// Kodia Framework standard: always return the 'data' field if it exists
			return result.data !== undefined ? result.data : result;
		} catch (err) {
			if (!this.error) {
				this.error = {
					success: false,
					code: 'Network Error',
					message: err instanceof Error ? err.message : 'Unknown error occurred'
				};
			}
			throw this.error;
		} finally {
			this.isLoading = false;
		}
	}

	// Helper methods for cleaner syntax
	get<P extends keyof paths>(path: P, params?: any) {
		return this.request(path, 'get', { params });
	}

	post<P extends keyof paths>(path: P, body?: any, params?: any) {
		return this.request(path, 'post', { body, params });
	}

	put<P extends keyof paths>(path: P, body?: any, params?: any) {
		return this.request(path, 'put', { body, params });
	}

	patch<P extends keyof paths>(path: P, body?: any, params?: any) {
		return this.request(path, 'patch', { body, params });
	}

	delete<P extends keyof paths>(path: P, params?: any) {
		return this.request(path, 'delete', { params });
	}
}

export const api = new ApiClient();
