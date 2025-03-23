import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  // Check for the authentication token
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // Fetch environments from your API using the token for authentication
  const res = await fetch('http://api:8080/environments', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });

  if (res.ok) {
    const data = await res.json();
    return data;
  } else {
    // If the API returns an error (e.g. token expired), redirect to login
    throw redirect(303, '/login');
  }
};