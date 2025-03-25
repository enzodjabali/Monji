import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) throw redirect(303, '/login');

  // 1) Fetch user info.
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) throw redirect(303, '/login');
  const userData = await userRes.json();

  // 2) Fetch environments for the Navbar.
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) throw redirect(303, '/login');
  const envData = await envRes.json();

  // 3) Fetch the document to edit.
  const { id, dbName, collectionName, docID } = params;
  const docRes = await fetch(
    `http://api:8080/environments/${id}/databases/${dbName}/collections/${collectionName}/documents/${docID}`,
    { headers: { Authorization: `Bearer ${token}` } }
  );
  if (!docRes.ok) {
    throw redirect(303, `/environments/${id}/databases/${dbName}/collections/${collectionName}/documents`);
  }
  const docData = await docRes.json();
  return {
    user: userData.user,
    environments: envData.environments || [],
    document: docData.document,
    database: docData.database,
    collection: docData.collection,
    currentEnvironmentId: id,
    currentDatabase: dbName,
    currentCollection: collectionName,
    docID
  };
};

export const actions: Actions = {
  default: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) throw redirect(303, '/login');

    const { id, dbName, collectionName, docID } = params;
    const formData = await request.formData();
    const jsonData = formData.get('document');
    if (typeof jsonData !== 'string') {
      return fail(400, { error: 'Invalid document data' });
    }

    // Validate JSON format.
    try {
      JSON.parse(jsonData);
    } catch (err) {
      return fail(400, { error: 'Invalid JSON format' });
    }

    // Call your API with a PUT request.
    const res = await fetch(
      `http://api:8080/environments/${id}/databases/${dbName}/collections/${collectionName}/documents/${docID}`,
      {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        // Pass the raw JSON string as expected by your API.
        body: jsonData
      }
    );
    if (!res.ok) {
      return fail(400, { error: 'Failed to update document' });
    }
    const result = await res.json();
    return { success: true, result };
  }
};