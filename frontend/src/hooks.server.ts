import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	// Extract auth token from cookies
	const token = event.cookies.get('access_token');

	if (token) {
		event.locals.token = token;
		event.locals.isAuthenticated = true;
	} else {
		event.locals.isAuthenticated = false;
	}

	return resolve(event);
};
