import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // 1) Fetch user info
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // 2) Fetch environments (for Navbar)
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  // 3) Fetch documents for the specified collection
  const { id, dbName, collectionName } = params;
  const docsRes = await fetch(
    `http://api:8080/environments/${id}/databases/${dbName}/collections/${collectionName}/documents`,
    {
      headers: { Authorization: `Bearer ${token}` }
    }
  );
  if (!docsRes.ok) {
    // If call fails, redirect back to the collections page
    throw redirect(303, `/environments/${id}/databases/${dbName}/collections`);
  }
  const docsData = await docsRes.json();
  // Example shape:
  // {
  //   "collection": "system.users",
  //   "database": "admin",
  //   "documents": [ { "_id": "...", ... }, ... ]
  // }

  return {
    user: userData.user,
    environments: envData.environments || [],
    documents: docsData.documents || [],
    database: docsData.database,
    collection: docsData.collection,
    currentEnvironmentId: id,
    currentDatabase: dbName,
    currentCollection: collectionName
  };
};