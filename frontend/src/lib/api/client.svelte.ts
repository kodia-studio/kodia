import { PUBLIC_API_URL } from '$env/static/public';
import type { paths } from '../types/generated/api';
import type { ApiError } from '../types/api.types';
import { authStore } from '../stores/auth.store';
import { devStore } from '../stores/dev.svelte';
import { get } from 'svelte/store';

type HttpMethod = 'get' | 'post' | 'put' | 'patch' | 'delete';

class ApiClient {
	private baseUrl: string;
	isLoading = $state(false);
	error = $state<ApiError | null>(null);

	constructor() {
		// Institutional standard: Strip /api if present as generated paths include it
		this.baseUrl = (PUBLIC_API_URL || 'http://localhost:8080/api').replace(/\/api$/, '');
	}

	/**
	 * Institutional-grade request handler with flexible typing
	 */
	async request<T = any>(
		path: string,
		method: HttpMethod,
		options: {
			params?: any;
			body?: any;
			headers?: Record<string, string>;
		} = {}
	): Promise<T> {
		this.isLoading = true;
		this.error = null;

		let url = `${this.baseUrl}${path}`;
		
		// Handle path parameters
		if (options.params?.path) {
			Object.entries(options.params.path).forEach(([key, value]) => {
				url = url.replace(`{${key}}`, String(value));
			});
		}

		// Handle query parameters
		if (options.params?.query || (method === 'get' && options.params)) {
			const searchParams = new URLSearchParams();
			const queryParams = options.params?.query || options.params;
			
			if (typeof queryParams === 'object') {
				Object.entries(queryParams).forEach(([key, value]) => {
					if (value !== undefined && value !== null) {
						searchParams.append(key, String(value));
					}
				});
				const queryString = searchParams.toString();
				if (queryString) {
					url += (url.includes('?') ? '&' : '?') + queryString;
				}
			}
		}

		const headers = new Headers(options.headers);
		if (!headers.has('Content-Type')) {
			headers.set('Content-Type', 'application/json');
		}

		const token = get(authStore).accessToken;
		if (token) {
			headers.set('Authorization', `Bearer ${token}`);
		}

		const logId = Math.random().toString(36).substring(7);
		const startTime = performance.now();

		devStore.logRequest({
			id: logId,
			method: method.toUpperCase(),
			path: url.replace(this.baseUrl, ''),
			requestBody: options.body
		});

		try {
			const response = await fetch(url, {
				method: method.toUpperCase(),
				headers,
				body: options.body ? JSON.stringify(options.body) : undefined
			});

			const result = await response.json();
			const duration = performance.now() - startTime;

			devStore.updateLog(logId, {
				status: response.status,
				duration,
				responseBody: result
			});

			if (!response.ok) {
				const apiError = result as ApiError;
				this.error = apiError;
				throw apiError;
			}

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

	// Helper methods with path autocomplete support
	get<T = any, P extends keyof paths | (string & {}) = string>(path: P, params?: any) {
		return this.request<T>(path as string, 'get', { params });
	}

	post<T = any, P extends keyof paths | (string & {}) = string>(path: P, body?: any, params?: any) {
		return this.request<T>(path as string, 'post', { body, params });
	}

	put<T = any, P extends keyof paths | (string & {}) = string>(path: P, body?: any, params?: any) {
		return this.request<T>(path as string, 'put', { body, params });
	}

	patch<T = any, P extends keyof paths | (string & {}) = string>(path: P, body?: any, params?: any) {
		return this.request<T>(path as string, 'patch', { body, params });
	}

	delete<T = any, P extends keyof paths | (string & {}) = string>(path: P, params?: any) {
		return this.request<T>(path as string, 'delete', { params });
	}
}

export const api = new ApiClient();
