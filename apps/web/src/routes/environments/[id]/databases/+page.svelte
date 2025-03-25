<!-- apps/web/src/routes/environments/[id]/databases/+page.svelte -->
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
    databases: {
      Name: string;
      SizeOnDisk: number;
      Empty: boolean;
    }[];
    totalSize: number;
    currentEnvironmentId: string;
  };

  // Manage modals for Create / Edit / Delete
  let showCreateModal = false;
  let showEditModal = false;
  let showDeleteModal = false;

  // For "Manage" dropdown on each database
  let manageDropdownOpen: string | null = null;

  // Fields for creating a DB
  let newDbName = '';

  // Fields for renaming a DB
  let oldDbName: string | null = null; // store the original name
  let renameDbNewName = '';

  // Fields for deleting a DB
  let deleteDbName: string | null = null;
  let typedDbName = ''; // confirm typed name

  // Toggle the "Manage" dropdown for a given dbName
  function toggleManageDropdown(dbName: string) {
    manageDropdownOpen = manageDropdownOpen === dbName ? null : dbName;
  }

  function closeManageDropdown() {
    manageDropdownOpen = null;
  }

  // Open the Create Database modal
  function openCreateModal() {
    newDbName = '';
    showCreateModal = true;
  }

  // Open the Rename modal
  function openEditModal(dbName: string) {
    oldDbName = dbName;
    renameDbNewName = dbName; // start with the current name
    closeManageDropdown();
    showEditModal = true;
  }

  // Open the Delete modal
  function openDeleteModal(dbName: string) {
    deleteDbName = dbName;
    typedDbName = '';
    closeManageDropdown();
    showDeleteModal = true;
  }

  // Close all modals
  function closeModals() {
    showCreateModal = false;
    showEditModal = false;
    showDeleteModal = false;
  }

  // Close modals if user clicks the overlay background
  function handleOverlayClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      closeModals();
    }
  }

  // If user clicks anywhere outside the open "Manage" dropdown, close it
  function handleWindowClick(e: MouseEvent) {
    if (!manageDropdownOpen) return;
    const container = document.getElementById(`db-manage-dropdown-${manageDropdownOpen}`);
    if (!container) return;

    if (!container.contains(e.target as Node)) {
      manageDropdownOpen = null;
    }
  }
</script>

<!-- We pass currentEnvironmentId to the Navbar so it’s preselected -->
<Navbar
  user={data.user}
  environments={data.environments}
  currentEnvironmentId={data.currentEnvironmentId}
/>

<!-- BREADCRUMB: environmentId => "Environments / Databases" -->
<Breadcrumb environmentId={data.currentEnvironmentId} />

<!-- Close "Manage" dropdown if user clicks outside -->
<svelte:window on:click={handleWindowClick} />

<div class="bg-gray-100 min-h-screen p-8">
  <div class="max-w-7xl mx-auto grid gap-6 md:grid-cols-[2fr_1fr]">

    <!-- LEFT COLUMN: Databases box -->
    <div class="bg-white rounded-lg shadow p-6 space-y-4">
      <!-- Heading + create button -->
      <div class="flex justify-between items-center mb-2">
        <h2 class="text-2xl font-bold text-gray-800">Databases</h2>
        <button
          on:click={openCreateModal}
          class="text-sm px-3 py-1 bg-[#1B6609] text-white rounded hover:bg-[#1B6609]/90 transition"
        >
          Create a database
        </button>
      </div>
      <!-- Total Size -->
      <p class="text-gray-700 mb-4">
        Total Size: {data.totalSize} bytes
      </p>

      <!-- Databases list -->
      {#if data.databases?.length > 0}
        <div class="grid gap-4">
          {#each data.databases as db}
            <div class="border border-gray-200 rounded p-4 hover:shadow transition relative">
              <!-- Top row: Name + Manage button -->
              <div class="flex items-center justify-between mb-1">
                <h3 class="font-semibold text-lg text-gray-800">
                  <!-- Link to the collections for this DB -->
                  <a
                    href={`/environments/${data.currentEnvironmentId}/databases/${db.Name}/collections`}
                    class="hover:underline"
                  >
                    {db.Name}
                  </a>
                </h3>
                <!-- Manage button -->
                <div
                  class="relative"
                  id={"db-manage-dropdown-" + db.Name}
                >
                  <button
                    on:click={() => toggleManageDropdown(db.Name)}
                    class="text-sm px-3 py-1 bg-[#1B6609] text-white rounded
                           hover:bg-[#1B6609]/90 transition"
                  >
                    Manage
                  </button>
                  {#if manageDropdownOpen === db.Name}
                    <div
                      class="absolute right-0 mt-1 w-32 bg-white border border-gray-200 rounded shadow z-10"
                      transition:scale
                    >
                      <ul class="py-1">
                        <li>
                          <button
                            class="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm"
                            on:click={() => openEditModal(db.Name)}
                          >
                            Rename
                          </button>
                        </li>
                        <li>
                          <button
                            class="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm text-red-600"
                            on:click={() => openDeleteModal(db.Name)}
                          >
                            Delete
                          </button>
                        </li>
                      </ul>
                    </div>
                  {/if}
                </div>
              </div>

              <!-- Some DB stats -->
              <p class="text-sm text-gray-600">
                Size on Disk: {db.SizeOnDisk} bytes
              </p>
              <p class="text-sm text-gray-600">
                Empty: {db.Empty ? 'Yes' : 'No'}
              </p>
            </div>
          {/each}
        </div>
      {:else}
        <p class="text-gray-600">No databases found.</p>
      {/if}
    </div>

    <!-- RIGHT COLUMN: "Toolbar" or other info -->
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

<!-- CREATE DATABASE MODAL -->
{#if showCreateModal}
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    transition:fade={{ duration: 150 }}
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      transition:scale={{ duration: 150 }}
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Create a database</h2>
      <!-- Action form => ?/createDb -->
      <form method="post" action="?/createDb" class="space-y-4">
        <!-- We won't ask for initialCollection; we pass "delete_me" behind the scenes -->
        <input type="hidden" name="initialCollection" value="delete_me" />

        <div>
          <label class="block font-semibold mb-1" for="newDbName">Database name</label>
          <input
            id="newDbName"
            name="dbName"
            type="text"
            bind:value={newDbName}
            placeholder="e.g. myNewDatabase"
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
            Create
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- EDIT (RENAME) DATABASE MODAL -->
{#if showEditModal && oldDbName !== null}
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    transition:fade={{ duration: 150 }}
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      transition:scale={{ duration: 150 }}
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4">Rename the database</h2>
      <!-- Action form => ?/updateDb -->
      <form method="post" action="?/updateDb" class="space-y-4">
        <!-- Hidden field for the old DB name -->
        <input type="hidden" name="oldDbName" value={oldDbName} />

        <div>
          <label class="block font-semibold mb-1" for="renameDbNewName">New name</label>
          <input
            id="renameDbNewName"
            name="newDbName"
            type="text"
            bind:value={renameDbNewName}
            placeholder="Enter new database name"
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

<!-- DELETE DATABASE MODAL -->
{#if showDeleteModal && deleteDbName !== null}
  <div
    class="fixed inset-0 flex items-center justify-center bg-black/20 z-50"
    transition:fade={{ duration: 150 }}
    on:click={handleOverlayClick}
  >
    <div
      class="bg-white rounded-md p-6 w-full max-w-md"
      transition:scale={{ duration: 150 }}
      on:click|stopPropagation
    >
      <h2 class="text-xl font-bold mb-4 text-red-600">Delete the database</h2>
      <p class="mb-4">
        To confirm, type the database name: <strong>"{deleteDbName}"</strong> below.
      </p>

      <!-- Action form => ?/deleteDb -->
      <form method="post" action="?/deleteDb" class="space-y-4">
        <!-- Hidden field for the DB name to delete -->
        <input type="hidden" name="dbName" value={deleteDbName} />

        <div>
          <label class="block font-semibold mb-1" for="typedDbName">
            Database name:
          </label>
          <input
            id="typedDbName"
            type="text"
            bind:value={typedDbName}
            placeholder="{deleteDbName}"
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
            disabled={typedDbName !== deleteDbName}
          >
            Delete
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}