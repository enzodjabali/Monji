// apps/web/src/routes/environments/[id]/databases/[dbName]/collections/+page.server.ts
import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) throw redirect(303, '/login');

  // Fetch user info
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) throw redirect(303, '/login');
  const userData = await userRes.json();

  // Fetch environments for the Navbar
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) throw redirect(303, '/login');
  const envData = await envRes.json();

  // Fetch collections for the selected database
  const { id, dbName } = params;
  const colRes = await fetch(
    `http://api:8080/environments/${id}/databases/${dbName}/collections`,
    {
      headers: { Authorization: `Bearer ${token}` }
    }
  );
  if (!colRes.ok) {
    // Redirect back to the databases page if the API call fails
    throw redirect(303, `/environments/${id}/databases`);
  }
  const colData = await colRes.json();

  return {
    user: userData.user,
    environments: envData.environments || [],
    collections: colData.collections || [],
    database: colData.database,
    currentEnvironmentId: id,
    currentDatabase: dbName
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

    // Call API: POST /environments/:id/databases/:dbName/collections
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

    // Call API: PUT /environments/:id/databases/:dbName/collections/:collName
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

    // Call API: DELETE /environments/:id/databases/:dbName/collections/:collName
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