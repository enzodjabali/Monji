<script lang="ts">
    import Navbar from '$lib/components/Navbar.svelte';
  
    // Data returned by +page.server.ts
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
    };
  </script>
  
  <!-- NAVBAR -->
  <Navbar user={data.user} environments={data.environments} />
  
  <!-- PAGE BACKGROUND -->
  <div class="bg-gray-100 min-h-screen p-8">
    <!-- Optional container to center content -->
    <div class="max-w-7xl mx-auto">
      <!-- 2-column layout: left = 2/3, right = 1/3 -->
      <div class="grid gap-6 md:grid-cols-[2fr_1fr]">
        
        <!-- LEFT BUBBLE: ENVIRONMENTS LIST -->
        <div class="bg-white rounded-lg shadow p-6 space-y-4">
          <h2 class="text-2xl font-bold text-gray-800">Environments</h2>
          
          {#if data.environments.length > 0}
            <div class="grid gap-4">
              {#each data.environments as env}
                <div class="bg-white border border-gray-200 rounded-lg p-4">
                  <div class="flex items-center justify-between mb-1">
                    <h3 class="text-lg font-semibold text-gray-800">
                      {env.name}
                    </h3>
                    <!-- "Manage" button in #1B6609 -->
                    <button class="text-[#1B6609] text-sm hover:underline">
                      Manage
                    </button>
                  </div>
                  <!-- Only show the connection string, remove ID/Created By -->
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
  
        <!-- RIGHT BUBBLE: TOOLBAR / LINKS -->
        <div class="bg-white rounded-lg shadow p-6 space-y-6">
          <h2 class="text-2xl font-bold text-gray-800">Toolbar</h2>
          
          <!-- Example: "Recommended Resources" section -->
          <div>
            <h3 class="text-lg font-semibold text-gray-800 mb-2">
              Recommended Resources
            </h3>
            <!-- Make all links #1B6609 -->
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