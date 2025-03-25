<script lang="ts">
    import Navbar from '$lib/components/Navbar.svelte';
    // data.environments is the array returned by the load function
    export let data: {
      environments: {
        id: number;
        name: string;
        connection_string: string;
        created_by: number;
      }[];
    };
  </script>

  <Navbar title="Monji" />
  <div class="bg-gray-100 min-h-screen p-8 space-y-6">
    <!-- Page Header -->
    <header class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-800">
        Environments
      </h1>
      <div>
        <!-- Example 'Create environment' button -->
        <button
          class="bg-blue-600 text-white px-4 py-2 rounded shadow hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          Create Environment
        </button>
      </div>
    </header>
  
    {#if data.environments.length > 0}
      <!-- Grid of environment cards -->
      <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {#each data.environments as env}
          <div class="bg-white rounded shadow p-5 space-y-2">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-700">
                {env.name}
              </h2>
              <!-- Example "Manage" link or button -->
              <button class="text-blue-600 text-sm hover:text-blue-800">
                Manage
              </button>
            </div>
            <p class="text-sm text-gray-500">
              <span class="font-semibold">ID:</span> {env.id}
            </p>
            <p class="text-sm text-gray-500">
              <span class="font-semibold">Connection:</span> {env.connection_string}
            </p>
            <p class="text-sm text-gray-500">
              <span class="font-semibold">Created by:</span> {env.created_by}
            </p>
          </div>
        {/each}
      </div>
    {:else}
      <p class="text-gray-600">
        No environments available.
      </p>
    {/if}
  </div>