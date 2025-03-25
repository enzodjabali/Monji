import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // Fetch connected user info
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // Fetch the environments list for the navbar
  const envRes = await fetch('http://api:8080/environments', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  // Fetch databases for the selected environment (params.id)
  const envId = params.id;
  const dbRes = await fetch(`http://api:8080/environments/${envId}/databases`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  if (!dbRes.ok) {
    throw redirect(303, '/environments');
  }
  const dbData = await dbRes.json();

  return {
    user: userData.user,                // e.g. { id, first_name, last_name, email, role }
    environments: envData.environments,   // array of environment objects
    databases: dbData.Databases,          // array of { Name, SizeOnDisk, Empty }
    totalSize: dbData.TotalSize,          // total size across all databases
    currentEnvironmentId: envId
  };
};