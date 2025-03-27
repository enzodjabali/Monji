import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
  default: async ({ cookies }) => {
    // Remove/delete the token cookie
    cookies.delete('token', { path: '/' });

    // Redirect the user to /login
    throw redirect(303, '/login');
  }
};