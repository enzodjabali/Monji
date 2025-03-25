// apps/web/src/routes/environments/+page.server.ts
import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  // Check for the authentication token
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // 1) Fetch the current user
  const userRes = await fetch('http://api:8080/whoami', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  if (!userRes.ok) {
    // If something’s wrong with the token, redirect to login
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // 2) Fetch environments
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

  // Return both user and environments to the Svelte page
  return {
    user: userData.user,               // e.g. { id, first_name, last_name, ... }
    environments: envData.environments // e.g. array of environment objects
  };
};