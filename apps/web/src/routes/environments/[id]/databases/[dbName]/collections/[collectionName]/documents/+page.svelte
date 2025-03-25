<!-- apps/web/src/routes/environments/[id]/databases/[dbName]/collections/[collectionName]/documents/+page.svelte -->
<script lang="ts">
    import Navbar from '$lib/components/Navbar.svelte';
  
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
    };
  
    /** Helper to pretty-print a document as JSON. */
    function prettyPrint(doc: any) {
      return JSON.stringify(doc, null, 2);
    }
  </script>
  
  <!-- NAVBAR (with user & environments) -->
  <Navbar user={data.user} environments={data.environments} />
  
  <!-- PAGE LAYOUT -->
  <div class="bg-gray-100 min-h-screen p-8">
    <div class="max-w-7xl mx-auto">
      <div class="grid gap-6 md:grid-cols-[2fr_1fr]">
        
        <!-- LEFT BUBBLE: DOCUMENTS LIST -->
        <div class="bg-white rounded-lg shadow p-6 space-y-4">
          <h2 class="text-2xl font-bold text-gray-800">
            Documents in {data.collection} collection (DB: {data.database})
          </h2>
          
          {#if data.documents.length > 0}
            {#each data.documents as doc, i}
              <div class="border border-gray-200 rounded p-4 mb-4">
                <!-- Document Title -->
                <div class="text-gray-800 font-semibold mb-2">
                  Document {i + 1}
                </div>
                
                <!-- Example: show _id if present -->
                {#if doc._id}
                  <p class="text-sm text-gray-600 mb-2">
                    <span class="font-semibold">_id:</span> {doc._id}
                  </p>
                {/if}
  
                <!-- Two-column layout for top-level fields + full JSON -->
                <div class="flex flex-col md:flex-row gap-4">
                  <!-- Left: Key-Value fields (excluding _id) -->
                  <div class="md:w-1/2 space-y-1">
                    {#each Object.keys(doc).filter(k => k !== '_id') as key}
                      <p class="text-sm text-gray-600">
                        <span class="font-semibold">{key}:</span> {JSON.stringify(doc[key])}
                      </p>
                    {/each}
                  </div>
  
                  <!-- Right: Pretty-printed JSON -->
                  <div class="md:w-1/2 bg-gray-50 border border-gray-200 rounded p-2 overflow-auto">
                    <pre class="text-sm text-gray-700">
                      {prettyPrint(doc)}
                    </pre>
                  </div>
                </div>
              </div>
            {/each}
          {:else}
            <p class="text-gray-600">No documents found in this collection.</p>
          {/if}
        </div>
        
        <!-- RIGHT BUBBLE: TOOLBAR / LINKS -->
        <div class="bg-white rounded-lg shadow p-6 space-y-6">
          <h2 class="text-2xl font-bold text-gray-800">Toolbar</h2>
          <div>
            <h3 class="text-lg font-semibold text-gray-800 mb-2">Recommended Resources</h3>
            <ul class="list-disc list-inside space-y-1">
              <li>
                <a href="https://docs.mongodb.com" target="_blank" class="text-[#1B6609] hover:underline">
                  Documentation
                </a>
              </li>
              <li>
                <a href="https://university.mongodb.com" target="_blank" class="text-[#1B6609] hover:underline">
                  University
                </a>
              </li>
              <li>
                <a href="https://community.mongodb.com" target="_blank" class="text-[#1B6609] hover:underline">
                  Forums
                </a>
              </li>
              <li>
                <a href="https://support.mongodb.com" target="_blank" class="text-[#1B6609] hover:underline">
                  Support
                </a>
              </li>
            </ul>
          </div>
        </div>
  
      </div>
    </div>
  </div>