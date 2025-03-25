<script lang="ts">
    // We expect an array of environments passed in as a prop
    export let environments = [];
  </script>
  
  <nav class="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
    <!-- Left side: Logo + Select -->
    <div class="flex items-center space-x-4">
      <!-- Monji logo -->
      <img
        src="https://monji-assets.fra1.cdn.digitaloceanspaces.com/images/monji-logo-black.png"
        alt="Monji logo"
        class="h-8 w-auto"
      />
  
      <!-- The select element for environments -->
      <div class="relative">
        <select
          class="border border-gray-300 rounded px-3 py-1 focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          <!-- Render each environment in the dropdown -->
          {#each environments as env}
            <option value={env.id}>
              {env.name}
            </option>
          {/each}
        </select>
      </div>
    </div>
  
    <!-- Right side: Add any links, user menu, or “Logout” button here -->
    <div class="flex items-center space-x-4">
      <a
        href="/login"
        class="text-gray-700 hover:text-gray-900 font-medium"
      >
        Logout
      </a>
    </div>
  </nav>