<!-- apps/web/src/routes/environments/[id]/databases/+page.svelte -->
<script lang="ts">
    import Navbar from '$lib/components/Navbar.svelte';
  
    // This is what we returned from +page.server.ts
    export let data: {
      databases: {
        Name: string;
        SizeOnDisk: number;
        Empty: boolean;
      }[];
      totalSize: number;
    };
  </script>
  
  <!-- Optional: If you want the same navbar, pass in user/environments data from a shared layout 
       or re-fetch user data here. For simplicity, we'll just show the content. -->
  <!-- <Navbar user={...} environments={...} /> -->
  
  <div class="bg-gray-100 min-h-screen p-8">
    <div class="max-w-4xl mx-auto bg-white rounded-lg shadow p-6">
      <h1 class="text-2xl font-bold mb-4">Databases</h1>
      
      <p class="text-gray-700 mb-6">
        Total Size: {data.totalSize} bytes
      </p>
  
      {#if data.databases.length > 0}
        <div class="space-y-4">
          {#each data.databases as db}
            <div class="border border-gray-200 rounded p-4">
              <h2 class="font-semibold text-lg text-gray-800">
                {db.Name}
              </h2>
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
  </div>