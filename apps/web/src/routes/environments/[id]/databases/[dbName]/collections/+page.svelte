<script lang="ts">
    import Navbar from '$lib/components/Navbar.svelte';
    import Breadcrumb from '$lib/components/Breadcrumb.svelte';
  
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
        count: number;
        name: string;
        size: number;
        storageSize: number;
        totalIndexSize: number;
      }[];
      database: string;
      currentEnvironmentId: string;
      currentDatabase: string;
    };
  </script>
  
  <Navbar user={data.user} environments={data.environments} />
  
  <!-- BREADCRUMB: pass environmentId + databaseName => "Environments / Databases / Collections" -->
  <Breadcrumb
    environmentId={data.currentEnvironmentId}
    databaseName={data.currentDatabase}
  />
  
  <div class="bg-gray-100 min-h-screen p-8">
    <div class="max-w-7xl mx-auto">
      <div class="grid gap-6 md:grid-cols-[2fr_1fr]">
        
        <!-- LEFT BUBBLE: Collections List -->
        <div class="bg-white rounded-lg shadow p-6 space-y-4">
          <h2 class="text-2xl font-bold text-gray-800">
            Collections in {data.database}
          </h2>
          
          {#if data.collections && data.collections.length > 0}
            <div class="space-y-4">
              {#each data.collections as col}
                <!-- CLICKABLE LINK to Documents page -->
                <a
                  href={`/environments/${data.currentEnvironmentId}/databases/${data.currentDatabase}/collections/${col.name}/documents`}
                  class="block border border-gray-200 rounded p-4 hover:shadow transition"
                >
                  <h3 class="font-semibold text-lg text-gray-800">{col.name}</h3>
                  <p class="text-sm text-gray-600">Count: {col.count}</p>
                  <p class="text-sm text-gray-600">Size: {col.size} bytes</p>
                  <p class="text-sm text-gray-600">Storage Size: {col.storageSize} bytes</p>
                  <p class="text-sm text-gray-600">Total Index Size: {col.totalIndexSize} bytes</p>
                </a>
              {/each}
            </div>
          {:else}
            <p class="text-gray-600">No collections found in this database.</p>
          {/if}
        </div>
        
        <!-- RIGHT BUBBLE: Toolbar / Links -->
        <div class="bg-white rounded-lg shadow p-6 space-y-6">
          <h2 class="text-2xl font-bold text-gray-800">Toolbar</h2>
          <div>
            <h3 class="text-lg font-semibold text-gray-800 mb-2">
              Recommended Resources
            </h3>
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