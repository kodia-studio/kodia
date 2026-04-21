import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	// Check if user is authenticated
	const token = locals.token || null;

	if (!token) {
		throw redirect(302, '/login');
	}

	return {
		token
	};
};
