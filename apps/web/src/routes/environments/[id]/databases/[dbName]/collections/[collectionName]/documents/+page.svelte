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
    documents: Record<string, any>[];
    database: string;
    collection: string;
    currentEnvironmentId: string;
    currentDatabase: string;
    currentCollection: string;

    // breadcrumb props
    environmentId: string;
    environmentName: string;
    databaseName: string | null;
    collectionName: string | null;
    documentId: string | null;
  };

  let allFields = new Set<string>();
  for (const doc of data.documents) {
    for (const key of Object.keys(doc)) {
      allFields.add(key);
    }
  }
  const fieldNames = Array.from(allFields);

  function previewValue(value: any) {
    if (typeof value === 'object' && value !== null) {
      return JSON.stringify(value);
    }
    return String(value);
  }

  function handleRowClick(doc: any) {
    const docID = doc._id;
    goto(
      `/environments/${data.currentEnvironmentId}/databases/${data.currentDatabase}/collections/${data.currentCollection}/documents/${docID}`
    );
  }
</script>

<Navbar
  user={data.user}
  environments={data.environments}
  currentEnvironmentId={data.currentEnvironmentId}
/>

<Breadcrumb
  environmentId={data.environmentId}
  environmentName={data.environmentName}
  databaseName={data.databaseName}
  collectionName={data.collectionName}
  documentId={data.documentId}
/>

<div class="bg-gray-100 min-h-screen p-8">
  <div class="max-w-7xl mx-auto bg-white rounded-lg shadow p-6 space-y-4">
    <h2 class="text-2xl font-bold text-gray-800">
      Viewing Collection: {data.collection} (DB: {data.database})
    </h2>

    <div class="flex items-center space-x-2">
      <button class="bg-green-600 text-white px-4 py-2 rounded shadow hover:bg-green-500 transition">
        + New Document
      </button>
      <button class="bg-green-600 text-white px-4 py-2 rounded shadow hover:bg-green-500 transition">
        + New Index
      </button>
    </div>

    {#if data.documents.length > 0}
      <div class="overflow-auto">
        <table class="min-w-full mt-4 border-collapse">
          <thead>
            <tr class="bg-gray-50 border-b">
              {#each fieldNames as field}
                <th class="py-2 px-3 text-left font-semibold text-gray-700 border-r last:border-r-0">
                  {field}
                </th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each data.documents as doc}
              <tr
                class="border-b hover:bg-gray-50 cursor-pointer"
                on:click={() => handleRowClick(doc)}
              >
                {#each fieldNames as field}
                  <td class="py-2 px-3 border-r last:border-r-0 align-top text-sm text-gray-700">
                    {#if doc[field] !== undefined}
                      {previewValue(doc[field])}
                    {:else}
                      <span class="text-gray-400">--</span>
                    {/if}
                  </td>
                {/each}
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {:else}
      <p class="text-gray-600">No documents found in this collection.</p>
    {/if}
  </div>
</div>