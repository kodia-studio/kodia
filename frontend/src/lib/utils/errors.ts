/**
 * Kodia Error Mapper
 * Maps backend validation errors (Record<string, string[]>) to frontend form structures.
 */

import type { ApiError } from '../types/api.types';

/**
 * Maps Kodia backend errors to a format compatible with Superforms.
 * Kodia returns errors as: { "email": ["already exists"], "password": ["too short"] }
 */
export function mapKodiaErrors(apiError: ApiError | any): Record<string, string | string[]> {
	if (!apiError || !apiError.errors) return {};

	const mapped: Record<string, string[]> = {};

	for (const [field, messages] of Object.entries(apiError.errors)) {
		mapped[field] = Array.isArray(messages) ? messages : [String(messages)];
	}

	return mapped;
}

/**
 * Checks if an error is a Kodia validation error (400 or 422).
 */
export function isValidationError(error: any): boolean {
	return (
		error &&
		typeof error === 'object' &&
		'errors' in error &&
		Object.keys(error.errors || {}).length > 0
	);
}
