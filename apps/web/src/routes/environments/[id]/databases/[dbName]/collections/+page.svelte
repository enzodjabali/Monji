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
    collections: {
      name: string;
      count: number;
      size: number;
      storageSize: number;
      totalIndexSize: number;
    }[];
    database: string;
    currentEnvironmentId: string;
    currentDatabase: string;

    // For the breadcrumb
    environmentId: string | null;
    environmentName: string | null;
    databaseName: string | null;
    collectionName: string | null;
    documentId: string | null;
  };

  let showCreateModal = false;
  let showEditModal = false;
  let showDeleteModal = false;

  let manageDropdownOpen: string | null = null;
  let newCollectionName = '';
  let oldCollectionName: string | null = null;
  let renameNewCollectionName = '';
  let deleteCollectionName: string | null = null;
  let typedCollectionName = '';

  function toggleManageDropdown(collName: string) {
    manageDropdownOpen = manageDropdownOpen === collName ? null : collName;
  }

  function closeManageDropdown() {
    manageDropdownOpen = null;
  }

  function openCreateModal() {
    newCollectionName = '';
    showCreateModal = true;
  }

  function openEditModal(collName: string) {
    oldCollectionName = collName;
    renameNewCollectionName = collName;
    closeManageDropdown();
    showEditModal = true;
  }

  function openDeleteModal(collName: string) {
    deleteCollectionName = collName;
    typedCollectionName = '';
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
    if (!manageDropdownOpen) return;
    const container = document.getElementById(`coll-manage-dropdown-${manageDropdownOpen}`);
    if (!container) return;
    if (!container.contains(e.target as Node)) {
      manageDropdownOpen = null;
    }
  }

  function goToDocuments(collName: string) {
    goto(
      `/environments/${data.currentEnvironmentId}/databases/${data.currentDatabase}/collections/${collName}/documents`
    );
  }
</script>

<svelte:window on:click={handleWindowClick} />

<Navbar
  user={data.user}
  environments={data.environments}
  currentEnvironmentId={data.currentEnvironmentId}
/>

<!-- Enhanced breadcrumb with environment & DB name -->
<Breadcrumb
  environmentId={data.environmentId}
  environmentName={data.environmentName}
  databaseName={data.databaseName}
  collectionName={data.collectionName}
  documentId={data.documentId}
/>

<div class="bg-gray-100 min-h-screen p-8">
  <div class="max-w-7xl mx-auto grid gap-6 md:grid-cols-[2fr_1fr]">
    <!-- LEFT: Collections -->
    <div class="bg-white rounded-lg shadow p-6 space-y-4">
      <div class="flex justify-between items-center mb-2">
        <h2 class="text-2xl font-bold text-gray-800">
          Collections in {data.database}
        </h2>
        <button
          on:click={openCreateModal}
          class="text-sm px-3 py-1 bg-[#1B6609] text-white rounded hover:bg-[#1B6609]/90 transition"
        >
          Create a collection
        </button>
      </div>

      {#if data.collections?.length > 0}
        <div class="grid gap-4">
          {#each data.collections as coll}
            <div
              class="border border-gray-200 rounded p-4 hover:shadow transition relative cursor-pointer"
              on:click={() => goToDocuments(coll.name)}
            >
              <div class="flex items-center justify-between mb-1">
                <h3 class="font-semibold text-lg text-gray-800">
                  {coll.name}
                </h3>
                <div
                  class="relative"
                  id={"coll-manage-dropdown-" + coll.name}
                  on:click|stopPropagation
                >
                  <button
                    on:click={() => toggleManageDropdown(coll.name)}
                    class="text-sm px-3 py-1 bg-[#1B6609] text-white rounded hover:bg-[#1B6609]/90 transition"
                  >
                    Manage
                  </button>
                  {#if manageDropdownOpen === coll.name}
                    <div
                      class="absolute right-0 mt-1 w-32 bg-white border border-gray-200 rounded shadow z-10"
                    >
                      <ul class="py-1">
                        <li>
                          <button
                            class="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm"
                            on:click={() => openEditModal(coll.name)}
                          >
                            Rename
                          </button>
                        </li>
                        <li>
                          <button
                            class="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm text-red-600"
                            on:click={() => openDeleteModal(coll.name)}
                          >
                            Delete
                          </button>
                        </li>
                      </ul>
                    </div>
                  {/if}
                </div>
              </div>
              <p class="text-sm text-gray-600">Count: {coll.count}</p>
              <p class="text-sm text-gray-600">Size: {coll.size} bytes</p>
              <p class="text-sm text-gray-600">Storage Size: {coll.storageSize} bytes</p>
              <p class="text-sm text-gray-600">Total Index Size: {coll.totalIndexSize} bytes</p>
            </div>
          {/each}
        </div>
      {:else}
        <p class="text-gray-600">No collections found.</p>
      {/if}
    </div>

    <!-- RIGHT: Toolbox / Resources -->
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

<!-- CREATE COLLECTION MODAL -->
{#if showCreateModal}
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Create a collection</h2>
      <form method="post" action="?/createCollection" class="space-y-4">
        <div>
          <label class="block font-semibold mb-1" for="newCollectionName">Collection name</label>
          <input
            id="newCollectionName"
            name="collectionName"
            type="text"
            bind:value={newCollectionName}
            required
            class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-600"
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
            Create
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- RENAME COLLECTION MODAL -->
{#if showEditModal && oldCollectionName !== null}
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Rename collection</h2>
      <form method="post" action="?/updateCollection" class="space-y-4">
        <input type="hidden" name="oldCollectionName" value={oldCollectionName} />
        <div>
          <label class="block font-semibold mb-1" for="renameNewCollectionName">New name</label>
          <input
            id="renameNewCollectionName"
            name="newCollectionName"
            type="text"
            bind:value={renameNewCollectionName}
            required
            class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-600"
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

<!-- DELETE COLLECTION MODAL -->
{#if showDeleteModal && deleteCollectionName !== null}
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4 text-red-600">Delete collection</h2>
      <p class="mb-4">
        To confirm, type the collection name: <strong>"{deleteCollectionName}"</strong> below.
      </p>
      <form method="post" action="?/deleteCollection" class="space-y-4">
        <input type="hidden" name="collectionName" value={deleteCollectionName} />
        <div>
          <label class="block font-semibold mb-1" for="typedCollectionName">
            Collection name:
          </label>
          <input
            id="typedCollectionName"
            type="text"
            bind:value={typedCollectionName}
            placeholder="{deleteCollectionName}"
            class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-600"
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
            disabled={typedCollectionName !== deleteCollectionName}
          >
            Delete
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}