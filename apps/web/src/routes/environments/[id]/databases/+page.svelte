<script lang="ts">
  import Navbar from '$lib/components/Navbar.svelte';

  // Data returned from the load function
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
</script>

<!-- Reuse the same Navbar (with user and environments) -->
<Navbar user={data.user} environments={data.environments} />

<!-- Page background and layout -->
<div class="bg-gray-100 min-h-screen p-8">
  <div class="max-w-7xl mx-auto">
    <!-- Two-column grid: left = 2/3, right = 1/3 -->
    <div class="grid gap-6 md:grid-cols-[2fr_1fr]">
      
      <!-- LEFT BUBBLE: Databases List -->
      <div class="bg-white rounded-lg shadow p-6 space-y-4">
        <h2 class="text-2xl font-bold text-gray-800">Databases</h2>
        <p class="text-gray-700 mb-4">Total Size: {data.totalSize} bytes</p>
        {#if data.databases.length > 0}
          <div class="space-y-4">
            {#each data.databases as db}
              <div class="border border-gray-200 rounded p-4">
                <h3 class="font-semibold text-lg text-gray-800">{db.Name}</h3>
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

      <!-- RIGHT BUBBLE: Toolbar / Links -->
      <div class="bg-white rounded-lg shadow p-6 space-y-6">
        <h2 class="text-2xl font-bold text-gray-800">Toolbar</h2>
        <div>
          <h3 class="text-lg font-semibold text-gray-800 mb-2">
            Recommended Resources
          </h3>
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