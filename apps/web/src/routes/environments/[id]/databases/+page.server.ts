// apps/web/src/routes/environments/[id]/databases/+page.server.ts
import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // Fetch connected user info
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // Fetch the environments list for the navbar
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  // Fetch databases for the selected environment (params.id)
  const envId = params.id;
  const dbRes = await fetch(`http://api:8080/environments/${envId}/databases`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!dbRes.ok) {
    throw redirect(303, '/environments');
  }
  const dbData = await dbRes.json();

  return {
    user: userData.user,
    environments: envData.environments,
    databases: dbData.Databases,
    totalSize: dbData.TotalSize,
    currentEnvironmentId: envId
  };
};

export const actions: Actions = {
  // Create a new database
  createDb: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const envId = params.id; // environment ID from the route
    const formData = await request.formData();
    const dbName = formData.get('dbName');
    const initialCollection = formData.get('initialCollection');

    if (typeof dbName !== 'string' || typeof initialCollection !== 'string') {
      return fail(400, { error: 'Invalid form data' });
    }

    // POST /environments/:id/databases
    const res = await fetch(`http://api:8080/environments/${envId}/databases`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        dbName,
        initialCollection
      })
    });

    if (!res.ok) {
      return fail(400, { error: 'Failed to create database' });
    }

    // Refresh the page
    throw redirect(303, `/environments/${envId}/databases`);
  },

  // Update (rename) a database
  updateDb: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const envId = params.id; // environment ID
    const formData = await request.formData();
    const oldDbName = formData.get('oldDbName'); // the current name
    const newDbName = formData.get('newDbName'); // the new name

    if (typeof oldDbName !== 'string' || typeof newDbName !== 'string') {
      return fail(400, { error: 'Invalid form data' });
    }

    // PUT /environments/:id/databases/:dbName
    const res = await fetch(
      `http://api:8080/environments/${envId}/databases/${oldDbName}`,
      {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({
          newDbName
        })
      }
    );

    if (!res.ok) {
      return fail(400, { error: 'Failed to rename database' });
    }

    throw redirect(303, `/environments/${envId}/databases`);
  },

  // Delete a database
  deleteDb: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const envId = params.id; // environment ID
    const formData = await request.formData();
    const dbName = formData.get('dbName'); // the database name to delete

    if (typeof dbName !== 'string') {
      return fail(400, { error: 'Invalid database name' });
    }

    // DELETE /environments/:id/databases/:dbName
    const res = await fetch(
      `http://api:8080/environments/${envId}/databases/${dbName}`,
      {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      }
    );

    if (!res.ok) {
      return fail(400, { error: 'Failed to delete database' });
    }

    throw redirect(303, `/environments/${envId}/databases`);
  }
};