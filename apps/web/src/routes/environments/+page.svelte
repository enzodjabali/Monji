<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import Breadcrumb from '$lib/components/Breadcrumb.svelte';
  import { goto } from '$app/navigation';

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
    environmentId: string | null;
    environmentName: string | null;
    databaseName: string | null;
    collectionName: string | null;
    documentId: string | null;
  };

  // State for modals & dropdown
  let showCreateModal = false;
  let showEditModal = false;
  let showDeleteModal = false;

  let manageDropdownOpen: number | null = null;

  let editEnvId: number | null = null;
  let editName = '';
  let editConnection = '';

  let deleteEnvId: number | null = null;
  let deleteEnvName = '';
  let deleteInputName = '';

  let newName = '';
  let newConnection = '';

  function openCreateModal() {
    newName = '';
    newConnection = '';
    showCreateModal = true;
  }

  function toggleManageDropdown(envId: number) {
    manageDropdownOpen = manageDropdownOpen === envId ? null : envId;
  }
  function closeManageDropdown() {
    manageDropdownOpen = null;
  }

  function openEditModal(envId: number, name: string, connection: string) {
    editEnvId = envId;
    editName = name;
    editConnection = connection;
    closeManageDropdown();
    showEditModal = true;
  }

  function openDeleteModal(envId: number, name: string) {
    deleteEnvId = envId;
    deleteEnvName = name;
    deleteInputName = '';
    closeManageDropdown();
    showDeleteModal = true;
  }

  function closeModals() {
    showCreateModal = false;
    showEditModal = false;
    showDeleteModal = false;
  }

  function handleOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      closeModals();
    }
  }

  function handleWindowClick(e: MouseEvent) {
    if (manageDropdownOpen === null) return;
    const container = document.getElementById(`env-manage-dropdown-${manageDropdownOpen}`);
    if (!container) return;
    if (!container.contains(e.target as Node)) {
      manageDropdownOpen = null;
    }
  }

  function goToDatabases(envId: number) {
    goto(`/environments/${envId}/databases`);
  }
</script>

<!-- Close "Manage" dropdown if user clicks outside -->
<svelte:window on:click={handleWindowClick} />

<Navbar user={data.user} environments={data.environments} currentEnvironmentId="" />

<!-- Enhanced breadcrumb: all null for this page, but pass anyway for consistency -->
<Breadcrumb
  environmentId={data.environmentId}
  environmentName={data.environmentName}
  databaseName={data.databaseName}
  collectionName={data.collectionName}
  documentId={data.documentId}
/>

<div class="bg-gray-100 min-h-screen p-8">
  <div class="max-w-7xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Environments</h2>
      <button
        on:click={openCreateModal}
        class="text-sm px-3 py-1 bg-[#1B6609] text-white rounded hover:bg-[#1B6609]/90 transition"
      >
        Add an environment
      </button>
    </div>

    <div class="grid gap-6 md:grid-cols-[2fr_1fr]">
      <!-- LEFT COLUMN: ENV LIST -->
      <div class="bg-white rounded-lg shadow p-6 space-y-4">
        {#if data.environments?.length > 0}
          <div class="grid gap-4">
            {#each data.environments as env}
              <div
                class="border border-gray-200 rounded-lg p-4 hover:shadow transition relative cursor-pointer"
                on:click={() => goToDatabases(env.id)}
              >
                <div class="flex items-center justify-between mb-1">
                  <h3 class="text-lg font-semibold text-gray-800">
                    {env.name}
                  </h3>
                  <!-- Manage button -->
                  <div
                    class="relative"
                    id={"env-manage-dropdown-" + env.id}
                    on:click|stopPropagation
                  >
                    <button
                      on:click={() => toggleManageDropdown(env.id)}
                      class="text-sm px-3 py-1 bg-[#1B6609] text-white rounded
                             hover:bg-[#1B6609]/90 transition border border-transparent"
                    >
                      Manage
                    </button>
                    {#if manageDropdownOpen === env.id}
                      <div
                        class="absolute right-0 mt-1 w-32 bg-white border border-gray-200 rounded shadow z-10"
                      >
                        <ul class="py-1">
                          <li>
                            <button
                              class="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm"
                              on:click={() =>
                                openEditModal(env.id, env.name, env.connection_string)
                              }
                            >
                              Edit
                            </button>
                          </li>
                          <li>
                            <button
                              class="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm text-red-600"
                              on:click={() => openDeleteModal(env.id, env.name)}
                            >
                              Remove
                            </button>
                          </li>
                        </ul>
                      </div>
                    {/if}
                  </div>
                </div>
                <p class="text-sm text-gray-500">
                  <span class="font-semibold">Connection:</span> {env.connection_string}
                </p>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-gray-600">No environments available.</p>
        {/if}
      </div>

      <!-- RIGHT COLUMN: TOOLBAR -->
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
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Add an environment</h2>
      <form method="post" action="?/createEnv" class="space-y-4">
        <div>
          <label class="block font-semibold mb-1" for="newName">Name</label>
          <input
            id="newName"
            name="name"
            type="text"
            bind:value={newName}
            placeholder="e.g. Mongo production environment"
            required
            class="w-full border border-gray-300 rounded px-3 py-2
                   focus:outline-none focus:ring-2 focus:ring-blue-600"
          />
        </div>
        <div>
          <label class="block font-semibold mb-1" for="newConnection">Connection string</label>
          <input
            id="newConnection"
            name="connection_string"
            type="text"
            bind:value={newConnection}
            placeholder="mongodb://root:password@host:27017"
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
            class="px-4 py-2 rounded bg-[#1B6609] text-white hover:bg-[#1B6609]/90"
          >
            Add
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- EDIT MODAL -->
{#if showEditModal && editEnvId !== null}
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Edit the environment</h2>
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
          <label class="block font-semibold mb-1" for="editConnection">Connection string</label>
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
            class="px-4 py-2 rounded bg-[#1B6609] text-white hover:bg-[#1B6609]/90"
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
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4 text-red-600">Remove the environment</h2>
      <p class="mb-4">
        To confirm, type the environment name:
        <strong>"{deleteEnvName}"</strong> below.
      </p>
      <form method="post" action="?/deleteEnv" class="space-y-4">
        <input type="hidden" name="id" value={deleteEnvId} />

        <div>
          <label class="block font-semibold mb-1" for="deleteInputName">
            Environment name:
          </label>
          <input
            id="deleteInputName"
            type="text"
            bind:value={deleteInputName}
            placeholder="{deleteEnvName}"
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
            class="px-4 py-2 rounded bg-red-600 text-white hover:bg-red-500"
            disabled={deleteInputName !== deleteEnvName}
          >
            Remove
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
