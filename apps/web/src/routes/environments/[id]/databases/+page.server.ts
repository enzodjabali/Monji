// apps/web/src/routes/environments/[id]/databases/+page.server.ts
import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    // If no token, redirect to login
    throw redirect(303, '/login');
  }

  // 'params.id' is the environment ID from the URL
  const envId = params.id;

  // Fetch databases for this environment
  const res = await fetch(`http://api:8080/environments/${envId}/databases`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  // If the API call fails or returns an error, handle it
  if (!res.ok) {
    // For example, redirect back to /environments or show an error
    throw redirect(303, '/environments');
  }

  // Example response shape:
  // {
  //   "Databases": [
  //     { "Name": "admin", "SizeOnDisk": 151552, "Empty": false },
  //     ...
  //   ],
  //   "TotalSize": 425984
  // }
  const data = await res.json();

  // Return the relevant info to +page.svelte
  return {
    databases: data.Databases, // array of { Name, SizeOnDisk, Empty }
    totalSize: data.TotalSize  // total size across all DBs
  };
};