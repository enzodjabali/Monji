import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // 1) Fetch user
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // 2) Fetch environments (Navbar)
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  // 3) environment name, DB name, etc.
  const { id, dbName, collectionName } = params;

  // Fetch environment name
  const singleEnvRes = await fetch(`http://api:8080/environments/${id}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!singleEnvRes.ok) {
    throw redirect(303, '/environments');
  }
  const singleEnvData = await singleEnvRes.json();
  const environmentName = singleEnvData.environment?.name || 'Unknown Env';

  // For DB name, we just use dbName
  const databaseDisplayName = dbName;

  // For collection name, we can just use collectionName
  const collectionDisplayName = collectionName;

  // 4) Fetch documents
  const docsRes = await fetch(
    `http://api:8080/environments/${id}/databases/${dbName}/collections/${collectionName}/documents`,
    {
      headers: { Authorization: `Bearer ${token}` }
    }
  );
  if (!docsRes.ok) {
    throw redirect(303, `/environments/${id}/databases/${dbName}/collections`);
  }
  const docsData = await docsRes.json();

  return {
    user: userData.user,
    environments: envData.environments || [],
    documents: docsData.documents || [],
    database: docsData.database,
    collection: docsData.collection,
    currentEnvironmentId: id,
    currentDatabase: dbName,
    currentCollection: collectionName,

    // For the breadcrumb
    environmentId: id,
    environmentName,
    databaseName: databaseDisplayName,
    collectionName: collectionDisplayName
  };
};