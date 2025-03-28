import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
  const token = cookies.get('token');

  if (!token) {
    // Not logged in → go to /login
    throw redirect(303, '/login');
  } else {
    // Logged in → go to /environments
    throw redirect(303, '/environments');
  }
};