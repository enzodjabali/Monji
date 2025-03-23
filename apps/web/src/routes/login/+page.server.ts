import type { Actions } from '@sveltejs/kit';
import { redirect, fail } from '@sveltejs/kit';

export const actions: Actions = {
	default: async ({ request, fetch, cookies }) => {
		const data = await request.formData();
		const email = data.get('email') as string;
		const password = data.get('password') as string;

		const res = await fetch('http://localhost:8080/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ email, password })
		});

		if (res.ok) {
			const result = await res.json();
			const token = result.token;
			// Store the token in an HTTP‑only cookie
			cookies.set('token', token, {
				path: '/',
				httpOnly: true,
				sameSite: 'strict',
				// secure: process.env.NODE_ENV === 'production',
				maxAge: 60 * 60 * 24 // 1 day
			});
			throw redirect(302, '/environments');
		} else {
			const errorData = await res.json();
			return fail(401, { error: errorData.error || 'Login failed' });
		}
	}
};