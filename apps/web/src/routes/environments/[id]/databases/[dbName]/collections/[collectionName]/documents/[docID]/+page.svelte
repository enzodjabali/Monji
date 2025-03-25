<script lang="ts">
    import Navbar from '$lib/components/Navbar.svelte';
    import Breadcrumb from '$lib/components/Breadcrumb.svelte';
    import { onMount } from 'svelte';
  
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
      document: Record<string, any>;
      database: string;
      collection: string;
      currentEnvironmentId: string;
      currentDatabase: string;
      currentCollection: string;
      docID: string;
    };
  
    let documentContent: string = JSON.stringify(data.document, null, 2);
    let errorMsg = '';
    let successMsg = '';
  
    async function handleSubmit(event: SubmitEvent) {
      event.preventDefault();
      const form = event.currentTarget as HTMLFormElement;
      const formData = new FormData(form);
      const response = await fetch(form.action, {
        method: form.method,
        body: formData
      });
      if (response.ok) {
        const result = await response.json();
        if (result.error) {
          errorMsg = result.error;
        } else {
          successMsg = 'Document updated successfully';
        }
      } else {
        errorMsg = 'Failed to update document';
      }
    }
  </script>
  
  <Navbar user={data.user} environments={data.environments} />
  
  <!-- BREADCRUMB: environmentId, databaseName, collectionName => "Environments / Databases / Collections / Documents" -->
  <Breadcrumb
    environmentId={data.currentEnvironmentId}
    databaseName={data.currentDatabase}
    collectionName={data.currentCollection}
  />
  
  <div class="bg-gray-100 min-h-screen p-8">
    <div class="max-w-4xl mx-auto bg-white rounded-lg shadow p-6 space-y-4">
      <h2 class="text-2xl font-bold text-gray-800">
        Editing Document: {data.docID}
      </h2>
  
      {#if errorMsg}
        <p class="text-red-600">{errorMsg}</p>
      {/if}
      {#if successMsg}
        <p class="text-green-600">{successMsg}</p>
      {/if}
  
      <!-- JSON Editor Form -->
      <form method="post" on:submit={handleSubmit} class="space-y-4">
        <textarea
          name="document"
          class="w-full h-80 border border-gray-300 rounded p-2 font-mono text-sm"
          bind:value={documentContent}
        ></textarea>
  
        <div class="flex space-x-2">
          <button
            type="submit"
            class="bg-[#1B6609] text-white px-4 py-2 rounded hover:bg-green-700 transition"
          >
            Save
          </button>
          <a
            href={`/environments/${data.currentEnvironmentId}/databases/${data.currentDatabase}/collections/${data.currentCollection}/documents`}
            class="bg-gray-300 text-gray-700 px-4 py-2 rounded hover:bg-gray-400 transition"
          >
            Cancel
          </a>
        </div>
      </form>
    </div>
  </div>