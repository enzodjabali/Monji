import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) throw redirect(303, '/login');

  // 1) Fetch user info
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) throw redirect(303, '/login');
  const userData = await userRes.json();

  // 2) Fetch all environments for the Navbar
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) throw redirect(303, '/login');
  const envData = await envRes.json();

  // 3) Fetch environment name
  const { id, dbName } = params;
  const singleEnvRes = await fetch(`http://api:8080/environments/${id}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!singleEnvRes.ok) {
    throw redirect(303, '/environments');
  }
  const singleEnvData = await singleEnvRes.json();
  const environmentName = singleEnvData.environment?.name ?? 'Unknown Env';

  // 4) If you want a “display name” for the DB, fetch it from your API. For now, just use dbName as is.
  const databaseDisplayName = dbName;

  // 5) Fetch collections
  const colRes = await fetch(
    `http://api:8080/environments/${id}/databases/${dbName}/collections`,
    {
      headers: { Authorization: `Bearer ${token}` }
    }
  );
  if (!colRes.ok) {
    throw redirect(303, `/environments/${id}/databases`);
  }
  const colData = await colRes.json();

  return {
    user: userData.user,
    environments: envData.environments || [],
    collections: colData.collections || [],
    database: colData.database, // might be "dbName" or something
    currentEnvironmentId: id,
    currentDatabase: dbName,

    // For breadcrumb
    environmentId: id,
    environmentName,
    databaseName: databaseDisplayName,
    collectionName: null
  };
};

export const actions: Actions = {
  // Create a new collection
  createCollection: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) throw redirect(303, '/login');

    const { id, dbName } = params;
    const formData = await request.formData();
    const collectionName = formData.get('collectionName');

    if (typeof collectionName !== 'string') {
      return fail(400, { error: 'Invalid collection name' });
    }

    // POST /environments/:id/databases/:dbName/collections
    const res = await fetch(`http://api:8080/environments/${id}/databases/${dbName}/collections`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ collectionName })
    });

    if (!res.ok) {
      return fail(400, { error: 'Failed to create collection' });
    }

    throw redirect(303, `/environments/${id}/databases/${dbName}/collections`);
  },

  // Rename an existing collection
  updateCollection: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) throw redirect(303, '/login');

    const { id, dbName } = params;
    const formData = await request.formData();
    const oldCollectionName = formData.get('oldCollectionName');
    const newCollectionName = formData.get('newCollectionName');

    if (typeof oldCollectionName !== 'string' || typeof newCollectionName !== 'string') {
      return fail(400, { error: 'Invalid form data' });
    }

    // PUT /environments/:id/databases/:dbName/collections/:collName
    const res = await fetch(
      `http://api:8080/environments/${id}/databases/${dbName}/collections/${oldCollectionName}`,
      {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ newCollectionName })
      }
    );

    if (!res.ok) {
      return fail(400, { error: 'Failed to rename collection' });
    }

    throw redirect(303, `/environments/${id}/databases/${dbName}/collections`);
  },

  // Delete a collection
  deleteCollection: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) throw redirect(303, '/login');

    const { id, dbName } = params;
    const formData = await request.formData();
    const collectionName = formData.get('collectionName');

    if (typeof collectionName !== 'string') {
      return fail(400, { error: 'Invalid collection name' });
    }

    // DELETE ...
    const res = await fetch(
      `http://api:8080/environments/${id}/databases/${dbName}/collections/${collectionName}`,
      {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      }
    );

    if (!res.ok) {
      return fail(400, { error: 'Failed to delete collection' });
    }

    throw redirect(303, `/environments/${id}/databases/${dbName}/collections`);
  }
};