// apps/web/src/routes/environments/[id]/databases/[dbName]/collections/[collectionName]/documents/+page.server.ts
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

  // 2) Fetch environments for the navbar
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  // 3) Fetch documents for this collection
  const { id, dbName, collectionName } = params;
  const docsRes = await fetch(
    `http://api:8080/environments/${id}/databases/${dbName}/collections/${collectionName}/documents`,
    {
      headers: { Authorization: `Bearer ${token}` }
    }
  );
  if (!docsRes.ok) {
    // If API call fails, redirect back to the database page
    throw redirect(303, `/environments/${id}/databases/${dbName}`);
  }
  const docsData = await docsRes.json();
  // Example shape:
  // {
  //   "collection": "delete_me",
  //   "database": "ilovemongodb",
  //   "documents": [ {...}, {...} ]
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