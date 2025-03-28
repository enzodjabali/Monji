import type { LayoutServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    // Not logged in? Go to /login
    throw redirect(303, '/login');
  }

  // Hit your API endpoint /whoami
  const whoamiRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });

  if (!whoamiRes.ok) {
    // If it fails, also go to /login
    throw redirect(303, '/login');
  }

  // This JSON should contain { permissions: {...}, user: {...} }
  const whoamiData = await whoamiRes.json();

  return {
    user: whoamiData.user,
    permissions: whoamiData.permissions
  };
};