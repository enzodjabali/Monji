import { redirect, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
  const token = cookies.get('token');
  if (token) {
    // User is already logged in. Redirect to /environments
    throw redirect(303, '/environments');
  }
};

export const actions: Actions = {
  default: async ({ request, fetch, cookies }) => {
    const formData = await request.formData();
    const email = formData.get('email');
    const password = formData.get('password');

    // Call your API to log in
    const res = await fetch('http://api:8080/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });

    if (res.ok) {
      const result = await res.json();
      // Store the token in an HTTP-only cookie
      cookies.set('token', result.token, {
        path: '/',
        httpOnly: true,
        // For production, also consider secure: true, a maxAge, etc.
      });
      // Redirect to /environments on success
      throw redirect(303, '/environments');
    } else {
      // Return an error message if login failed
      return fail(401, { error: 'Invalid email or password' });
    }
  }
};