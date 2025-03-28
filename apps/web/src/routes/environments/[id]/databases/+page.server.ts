import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // Fetch user
  const userRes = await fetch('http://api:8080/whoami', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // Fetch environment list for the Navbar
  const envRes = await fetch('http://api:8080/environments', {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  // Fetch specific environment to get its name
  const envId = params.id;
  const singleEnvRes = await fetch(`http://api:8080/environments/${envId}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!singleEnvRes.ok) {
    throw redirect(303, '/environments');
  }
  const singleEnvData = await singleEnvRes.json();
  const environmentName = singleEnvData.environment?.name ?? 'Unknown Env';

  // Fetch databases for that environment
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
    currentEnvironmentId: envId,

    // For the breadcrumb
    environmentId: envId,
    environmentName,
    databaseName: null,      // on databases page, no DB name yet
    collectionName: null     // no collection name
  };
};

export const actions: Actions = {
  // Create DB
  createDb: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const envId = params.id;
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

    throw redirect(303, `/environments/${envId}/databases`);
  },

  // Update (rename) DB
  updateDb: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const envId = params.id;
    const formData = await request.formData();
    const oldDbName = formData.get('oldDbName');
    const newDbName = formData.get('newDbName');

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

  // Delete DB
  deleteDb: async ({ request, params, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const envId = params.id;
    const formData = await request.formData();
    const dbName = formData.get('dbName');

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