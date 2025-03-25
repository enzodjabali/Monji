<!-- apps/web/src/lib/components/Navbar.svelte -->
<script lang="ts">
    // We expect user and environments passed in as props
    export let user: {
      id: number;
      first_name: string;
      last_name: string;
      email: string;
      role: string;
    } | null = null;
  
    export let environments: {
      id: number;
      name: string;
      connection_string: string;
      created_by: number;
    }[] = [];
  
    // Dropdown toggle
    let showDropdown = false;
  
    // Function to toggle the dropdown menu
    function toggleDropdown() {
      showDropdown = !showDropdown;
    }
  
    // Derive the initial from the user’s first name
    $: userInitial = user?.first_name
      ? user.first_name[0].toUpperCase()
      : '?';
  </script>
  
  <nav class="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
    <!-- Left side: Logo + Environment select -->
    <div class="flex items-center space-x-4">
      <img
        src="https://monji-assets.fra1.cdn.digitaloceanspaces.com/images/monji-logo-black.png"
        alt="Monji logo"
        class="h-8 w-auto"
      />
  
      <!-- Environments dropdown -->
      <div class="relative">
        <select
          class="border border-gray-300 rounded px-3 py-1 focus:outline-none focus:ring-2 focus:ring-blue-600"
        >
          {#each environments as env}
            <option value={env.id}>{env.name}</option>
          {/each}
        </select>
      </div>
    </div>
  
    <!-- Right side: Avatar & dropdown -->
    <div class="relative">
      <!-- Avatar Circle -->
      <div
        class="bg-gray-700 text-white h-8 w-8 flex items-center justify-center rounded-full cursor-pointer"
        on:click={toggleDropdown}
      >
        {userInitial}
      </div>
  
      <!-- Dropdown -->
      {#if showDropdown}
        <div class="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded shadow-lg z-10">
          <!-- User info section -->
          <div class="px-4 py-2">
            <p class="font-semibold">
              {user?.first_name} {user?.last_name}
            </p>
            <p class="text-sm text-gray-600">
              {user?.email}
            </p>
          </div>
          <hr />
          <!-- Action links -->
          <ul class="py-1">
            <li>
              <a href="/settings" class="block px-4 py-2 hover:bg-gray-100">
                Settings
              </a>
            </li>
            <li>
              <form method="post" action="/logout">
                <!-- Or a link if you do not need a POST -->
                <button
                  type="submit"
                  class="w-full text-left px-4 py-2 hover:bg-gray-100"
                >
                  Logout
                </button>
              </form>
            </li>
          </ul>
        </div>
      {/if}
    </div>
  </nav>