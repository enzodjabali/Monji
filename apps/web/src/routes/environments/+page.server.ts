// apps/web/src/routes/environments/+page.server.ts
import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // Fetch current user
  const userRes = await fetch('http://api:8080/whoami', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // Fetch environments
  const envRes = await fetch('http://api:8080/environments', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  return {
    user: userData.user,
    // If envData.environments is null, default to an empty array.
    environments: envData.environments || []
  };
};