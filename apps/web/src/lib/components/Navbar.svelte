<script lang="ts">
    import { scale } from 'svelte/transition';
  
    // Props
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
  
    // Local state
    let showDropdown = false;
  
    // References to elements for click-outside detection
    let dropdownRef: HTMLElement;
    let avatarRef: HTMLElement;
  
    // Toggle dropdown
    function toggleDropdown() {
      showDropdown = !showDropdown;
    }
  
    // If user clicks anywhere outside the avatar or dropdown, close it
    function handleClickOutside(event: MouseEvent) {
      if (!dropdownRef || !avatarRef) return;
  
      // If the click is NOT inside the dropdown or the avatar, close the dropdown
      if (
        !dropdownRef.contains(event.target as Node) &&
        !avatarRef.contains(event.target as Node)
      ) {
        showDropdown = false;
      }
    }
  
    // Compute the user's initial (e.g. first letter of first name)
    $: userInitial = user?.first_name
      ? user.first_name[0].toUpperCase()
      : '?';
  </script>
  
  <!-- Listen for clicks on the window to detect outside clicks -->
  <svelte:window on:click={handleClickOutside} />
  
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
  
    <!-- Right side: Avatar & Dropdown -->
    <div class="relative">
      <!-- Avatar Circle -->
      <div
        class="bg-gray-700 text-white h-8 w-8 flex items-center justify-center rounded-full cursor-pointer"
        on:click={toggleDropdown}
        bind:this={avatarRef}
      >
        {userInitial}
      </div>
  
      <!-- Dropdown, shown if 'showDropdown' is true -->
      {#if showDropdown}
        <div
          class="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded shadow-lg z-10"
          transition:scale={{ duration: 150 }}
          bind:this={dropdownRef}
        >
          <!-- User info -->
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
              <!-- Example: your logout could be a link or a form POST -->
              <form method="post" action="/logout">
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