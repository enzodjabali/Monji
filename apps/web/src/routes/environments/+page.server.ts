// apps/web/src/routes/environments/+page.server.ts

import type { PageServerLoad, Actions } from './$types';
import { redirect, fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
  const token = cookies.get('token');
  if (!token) {
    throw redirect(303, '/login');
  }

  // Fetch current user
  const userRes = await fetch('http://api:8080/whoami', {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`
    }
  });
  if (!userRes.ok) {
    throw redirect(303, '/login');
  }
  const userData = await userRes.json();

  // Fetch environments
  const envRes = await fetch('http://api:8080/environments', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`
    }
  });
  if (!envRes.ok) {
    throw redirect(303, '/login');
  }
  const envData = await envRes.json();

  return {
    user: userData.user,
    environments: envData.environments || []
  };
};

export const actions: Actions = {
  // Create a new environment
  createEnv: async ({ request, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const formData = await request.formData();
    const name = formData.get('name');
    const connection_string = formData.get('connection_string');

    if (typeof name !== 'string' || typeof connection_string !== 'string') {
      return fail(400, { error: 'Invalid form data' });
    }

    // Call API: POST /environments
    const res = await fetch('http://api:8080/environments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        name,
        connection_string
      })
    });

    if (!res.ok) {
      return fail(400, { error: 'Failed to create environment' });
    }

    // If successful, reload the /environments page
    throw redirect(303, '/environments');
  },

  // Update an existing environment
  updateEnv: async ({ request, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const formData = await request.formData();
    const envId = formData.get('id');
    const name = formData.get('name');
    const connection_string = formData.get('connection_string');

    if (
      typeof envId !== 'string' ||
      typeof name !== 'string' ||
      typeof connection_string !== 'string'
    ) {
      return fail(400, { error: 'Invalid form data' });
    }

    // Call API: PUT /environments/:id
    const res = await fetch(`http://api:8080/environments/${envId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        name,
        connection_string
      })
    });

    if (!res.ok) {
      return fail(400, { error: 'Failed to update environment' });
    }

    throw redirect(303, '/environments');
  },

  // Delete an environment
  deleteEnv: async ({ request, cookies, fetch }) => {
    const token = cookies.get('token');
    if (!token) {
      throw redirect(303, '/login');
    }

    const formData = await request.formData();
    const envId = formData.get('id');

    if (typeof envId !== 'string') {
      return fail(400, { error: 'Invalid environment ID' });
    }

    // Call API: DELETE /environments/:id
    const res = await fetch(`http://api:8080/environments/${envId}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`
      }
    });

    if (!res.ok) {
      return fail(400, { error: 'Failed to delete environment' });
    }

    throw redirect(303, '/environments');
  }
};