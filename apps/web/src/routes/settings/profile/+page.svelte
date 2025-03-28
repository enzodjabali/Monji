<script lang="ts">
    /**
     * We receive "data" from the parent layout load function.
     * That includes:
     *    data.user        => { id, first_name, last_name, email, role }
     *    data.permissions => { environments: [...], databases: [...] }
     */
    export let data: {
      user: {
        id: number;
        first_name: string;
        last_name: string;
        email: string;
        role: string;
      };
      permissions: {
        environments: Array<{
          environment_id: number;
          environment_name: string;
          permission: string;
        }>;
        databases: Array<{
          environment_id: number;
          environment_name: string;
          db_name: string;
          permission: string;
        }>;
      };
    };
  </script>
  
  <h1 class="text-2xl font-bold mb-4">My Profile</h1>
  
  <!-- Basic user info -->
  <div class="bg-white shadow rounded p-4 mb-6">
    <h2 class="text-xl font-semibold mb-2">User Details</h2>
    <p><strong>First Name:</strong> {data.user.first_name}</p>
    <p><strong>Last Name:</strong> {data.user.last_name}</p>
    <p><strong>Email:</strong> {data.user.email}</p>
    <p><strong>Role:</strong> {data.user.role}</p>
  </div>
  
  <!-- Permissions block (only shown if NOT admin or superadmin) -->
  {#if data.user.role !== 'admin' && data.user.role !== 'superadmin'}
    <div class="bg-white shadow rounded p-4">
      <h2 class="text-xl font-semibold mb-4">My Permissions</h2>
  
      <!-- Environment permissions -->
      <section class="mb-4">
        <h3 class="text-lg font-bold mb-1">Environments</h3>
        {#if data.permissions.environments.length > 0}
          <ul class="list-disc list-inside space-y-1">
            {#each data.permissions.environments as env}
              <li>
                <strong>{env.environment_name}</strong> → <em>{env.permission}</em>
              </li>
            {/each}
          </ul>
        {:else}
          <p class="text-sm text-gray-500">No environment permissions found.</p>
        {/if}
      </section>
  
      <!-- Database permissions -->
      <section>
        <h3 class="text-lg font-bold mb-1">Databases</h3>
        {#if data.permissions.databases.length > 0}
          <ul class="list-disc list-inside space-y-1">
            {#each data.permissions.databases as db}
              <li>
                <strong>{db.environment_name}</strong> / {db.db_name} → <em>{db.permission}</em>
              </li>
            {/each}
          </ul>
        {:else}
          <p class="text-sm text-gray-500">No database permissions found.</p>
        {/if}
      </section>
    </div>
  {/if}