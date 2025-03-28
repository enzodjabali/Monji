<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';
  import Breadcrumb from '$lib/components/Breadcrumb.svelte';
  import JsonEditor from '$lib/components/JsonEditor.svelte';

  // Data provided by the server (via load function)
  export let data: {
    user: { id: number; first_name: string; last_name: string; email: string; role: string };
    environments: Array<{ id: number; name: string; connection_string: string; created_by: number }>;
    document: Record<string, any>;
    database: string;
    collection: string;
    currentEnvironmentId: string;
    currentDatabase: string;
    currentCollection: string;
    docID: string;
  };

  // Derive the breadcrumb properties from the provided data.
  // Use currentEnvironmentId and find the environment name from the environments array.
  const environmentId = data.currentEnvironmentId;
  const environmentName =
    data.environments.find(
      (env) => env.id === Number(data.currentEnvironmentId)
    )?.name ?? '';
  const databaseName = data.database;
  const collectionName = data.collection;
  const documentId = data.docID;

  // Initialize the editor content as a formatted JSON string.
  let documentContent: string = JSON.stringify(data.document, null, 2);

  let errorMsg = '';
  let successMsg = '';

  // Use FormData to submit the form (so it’s sent as form‑encoded data)
  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    const form = event.currentTarget as HTMLFormElement;
    const formData = new FormData(form);

    // Ensure the latest editor content is in the hidden field.
    formData.set('document', documentContent);

    // Submit to the same URL (which calls the SvelteKit action)
    const response = await fetch(form.action, {
      method: form.method,
      body: formData
      // Do not set Content-Type manually—let the browser set it.
    });

    if (response.ok) {
      const result = await response.json();
      if (result.error) {
        errorMsg = result.error;
        successMsg = '';
      } else {
        successMsg = 'Document updated successfully!';
        errorMsg = '';
      }
    } else {
      errorMsg = 'Failed to update document';
      successMsg = '';
    }
  }
</script>

<!-- Navbar and Breadcrumb remain unchanged -->
<Navbar user={data.user} environments={data.environments} currentEnvironmentId={data.currentEnvironmentId} />

<Breadcrumb
  environmentId={environmentId}
  environmentName={environmentName}
  databaseName={databaseName}
  collectionName={collectionName}
  documentId={documentId}
/>

<div class="bg-gray-100 min-h-screen p-8">
  <div class="max-w-4xl mx-auto bg-white rounded-lg shadow p-6 space-y-4">
    <h2 class="text-2xl font-bold text-gray-800">
      Editing Document: {documentId}
    </h2>

    {#if errorMsg}
      <p class="text-red-600">{errorMsg}</p>
    {/if}
    {#if successMsg}
      <p class="text-green-600">{successMsg}</p>
    {/if}

    <!-- The form uses the default SvelteKit action -->
    <form method="post" on:submit={handleSubmit} class="space-y-4">
      <!-- The CodeMirror-based JSON editor -->
      <JsonEditor bind:value={documentContent} />

      <!-- Hidden textarea to send the updated JSON to the action -->
      <textarea name="document" class="hidden" bind:value={documentContent}></textarea>

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