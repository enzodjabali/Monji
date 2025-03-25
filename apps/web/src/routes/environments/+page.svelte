<!-- apps/web/src/routes/environments/+page.svelte -->
<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import Breadcrumb from '$lib/components/Breadcrumb.svelte';
  import { fade, scale } from 'svelte/transition';

  export let data: {
    user: {
      id: number;
      first_name: string;
      last_name: string;
      email: string;
      role: string;
    };
    environments: {
      id: number;
      name: string;
      connection_string: string;
      created_by: number;
    }[];
  };

  // Modals
  let showCreateModal = false;
  let showEditModal = false;
  let showDeleteModal = false;

  // For "edit" modal
  let editEnvId: number | null = null;
  let editName = '';
  let editConnection = '';

  // For "delete" modal
  let deleteEnvId: number | null = null;
  let deleteEnvName = '';

  // For "create" modal
  let newName = '';
  let newConnection = '';

  function openCreateModal() {
    newName = '';
    newConnection = '';
    showCreateModal = true;
  }

  function openEditModal(envId: number, name: string, connection: string) {
    editEnvId = envId;
    editName = name;
    editConnection = connection;
    showEditModal = true;
  }

  function openDeleteModal(envId: number, name: string) {
    deleteEnvId = envId;
    deleteEnvName = name;
    showDeleteModal = true;
  }

  function closeModals() {
    showCreateModal = false;
    showEditModal = false;
    showDeleteModal = false;
  }

  // Close modal if user clicks on the overlay (but not on modal content)
  function handleOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      closeModals();
    }
  }
</script>

<Navbar user={data.user} environments={data.environments} />
<!-- No environmentId in breadcrumb => only shows "Environments" -->
<Breadcrumb />

<div class="bg-gray-100 min-h-screen p-8">
  <div class="max-w-7xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Environments</h2>
      <button
        on:click={openCreateModal}
        class="bg-green-600 text-white px-4 py-2 rounded shadow hover:bg-green-500 transition"
      >
        + Create New Environment
      </button>
    </div>

    <div class="grid gap-6 md:grid-cols-[2fr_1fr]">
      <!-- Left: Environments List -->
      <div class="bg-white rounded-lg shadow p-6 space-y-4">
        {#if data.environments?.length > 0}
          <div class="grid gap-4">
            {#each data.environments as env}
              <div class="border border-gray-200 rounded-lg p-4 hover:shadow transition">
                <div class="flex items-center justify-between mb-1">
                  <h3 class="text-lg font-semibold text-gray-800">{env.name}</h3>
                  <a
                    href={`/environments/${env.id}/databases`}
                    class="text-[#1B6609] text-sm hover:underline"
                  >
                    Manage
                  </a>
                </div>
                <p class="text-sm text-gray-500">
                  <span class="font-semibold">Connection:</span> {env.connection_string}
                </p>
                <div class="mt-3 flex space-x-3">
                  <button
                    on:click={() => openEditModal(env.id, env.name, env.connection_string)}
                    class="text-sm px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-500"
                  >
                    Edit
                  </button>
                  <button
                    on:click={() => openDeleteModal(env.id, env.name)}
                    class="text-sm px-3 py-1 bg-red-600 text-white rounded hover:bg-red-500"
                  >
                    Delete
                  </button>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-gray-600">No environments available.</p>
        {/if}
      </div>

      <!-- Right: Toolbar/Links -->
      <div class="bg-white rounded-lg shadow p-6 space-y-6">
        <h2 class="text-2xl font-bold text-gray-800">Toolbar</h2>
        <div>
          <h3 class="text-lg font-semibold text-gray-800 mb-2">Recommended Resources</h3>
          <ul class="list-disc list-inside space-y-1">
            <li>
              <a
                href="https://docs.mongodb.com"
                target="_blank"
                class="text-[#1B6609] hover:underline"
              >
                Documentation
              </a>
            </li>
            <li>
              <a
                href="https://university.mongodb.com"
                target="_blank"
                class="text-[#1B6609] hover:underline"
              >
                University
              </a>
            </li>
            <li>
              <a
                href="https://community.mongodb.com"
                target="_blank"
                class="text-[#1B6609] hover:underline"
              >
                Forums
              </a>
            </li>
            <li>
              <a
                href="https://support.mongodb.com"
                target="_blank"
                class="text-[#1B6609] hover:underline"
              >
                Support
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- CREATE MODAL -->
{#if showCreateModal}
  <!-- Overlay -->
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    transition:fade={{ duration: 150 }}
    on:click={handleOverlayClick}
  >
    <!-- Modal Content -->
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      transition:scale={{ duration: 150 }}
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Create New Environment</h2>
      <form method="post" action="?/createEnv" class="space-y-4">
        <div>
          <label class="block font-semibold mb-1" for="newName">Name</label>
          <input
            id="newName"
            name="name"
            type="text"
            bind:value={newName}
            required
            class="w-full border border-gray-300 rounded px-3 py-2
                   focus:outline-none focus:ring-2 focus:ring-blue-600"
          />
        </div>
        <div>
          <label class="block font-semibold mb-1" for="newConnection">Connection String</label>
          <input
            id="newConnection"
            name="connection_string"
            type="text"
            bind:value={newConnection}
            required
            class="w-full border border-gray-300 rounded px-3 py-2
                   focus:outline-none focus:ring-2 focus:ring-blue-600"
          />
        </div>
        <div class="flex justify-end space-x-2">
          <button
            type="button"
            on:click={closeModals}
            class="px-4 py-2 rounded bg-gray-300 hover:bg-gray-400 text-gray-700"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-4 py-2 rounded bg-green-600 text-white hover:bg-green-500"
          >
            Create
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- EDIT MODAL -->
{#if showEditModal && editEnvId !== null}
  <!-- Overlay -->
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    transition:fade={{ duration: 150 }}
    on:click={handleOverlayClick}
  >
    <!-- Modal Content -->
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      transition:scale={{ duration: 150 }}
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Edit Environment</h2>
      <form method="post" action="?/updateEnv" class="space-y-4">
        <input type="hidden" name="id" value={editEnvId} />

        <div>
          <label class="block font-semibold mb-1" for="editName">Name</label>
          <input
            id="editName"
            name="name"
            type="text"
            bind:value={editName}
            required
            class="w-full border border-gray-300 rounded px-3 py-2
                   focus:outline-none focus:ring-2 focus:ring-blue-600"
          />
        </div>
        <div>
          <label class="block font-semibold mb-1" for="editConnection">Connection String</label>
          <input
            id="editConnection"
            name="connection_string"
            type="text"
            bind:value={editConnection}
            required
            class="w-full border border-gray-300 rounded px-3 py-2
                   focus:outline-none focus:ring-2 focus:ring-blue-600"
          />
        </div>
        <div class="flex justify-end space-x-2">
          <button
            type="button"
            on:click={closeModals}
            class="px-4 py-2 rounded bg-gray-300 hover:bg-gray-400 text-gray-700"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-500"
          >
            Save
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- DELETE MODAL -->
{#if showDeleteModal && deleteEnvId !== null}
  <!-- Overlay -->
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    transition:fade={{ duration: 150 }}
    on:click={handleOverlayClick}
  >
    <!-- Modal Content -->
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      transition:scale={{ duration: 150 }}
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4 text-red-600">Delete Environment</h2>
      <p class="mb-4">
        Are you sure you want to delete "<strong>{deleteEnvName}</strong>"?
        This action cannot be undone.
      </p>
      <form method="post" action="?/deleteEnv" class="space-y-4">
        <input type="hidden" name="id" value={deleteEnvId} />
        <div class="flex justify-end space-x-2">
          <button
            type="button"
            on:click={closeModals}
            class="px-4 py-2 rounded bg-gray-300 hover:bg-gray-400 text-gray-700"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-4 py-2 rounded bg-red-600 text-white hover:bg-red-500"
          >
            Delete
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}