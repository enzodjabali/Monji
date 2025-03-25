import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // Fetch user info
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // Fetch environments for the Navbar
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  // Fetch collections for the selected database
  const envId = params.id;
  const dbName = params.dbName;
  const colRes = await fetch(
    `http://api:8080/environments/${envId}/databases/${dbName}/collections`,
    {
      headers: { Authorization: `Bearer ${token}` }
    }
  );
  if (!colRes.ok) {
    // Redirect back to the databases page if the API call fails
    throw redirect(303, `/environments/${envId}/databases`);
  }
  const colData = await colRes.json();
  // Example shape: { "database": "...", "collections": [ {...}, ... ] }

  return {
    user: userData.user,
    environments: envData.environments || [],
    collections: colData.collections || [],
    database: colData.database,
    currentEnvironmentId: envId,
    currentDatabase: dbName
  };
};