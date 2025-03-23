import type { PageServerLoad } from './$types';
import { redirect, error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const token = cookies.get('token');
	if (!token) {
		throw redirect(302, '/login');
	}

	const res = await fetch('http://localhost:8080/environments', {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json',
			'Authorization': `Bearer ${token}`
		}
	});

	if (res.ok) {
		const data = await res.json();
		return {
			environments: data.environments
		};
	} else {
		throw error(res.status, 'Failed to fetch environments');
	}
};