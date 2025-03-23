import type { Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

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
        // For production, consider setting `secure: true` and an appropriate `maxAge`
      });
      // Redirect to the environments page upon successful login
      throw redirect(303, '/environments');
    } else {
      // Return an error message if login failed
      return fail(401, { error: 'Invalid email or password' });
    }
  }
};