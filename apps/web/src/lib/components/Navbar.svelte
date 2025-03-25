<script lang="ts">
  import { scale } from 'svelte/transition';
  import { goto } from '$app/navigation';

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

  /**
   * The ID of the currently selected environment, or ''/null if none is selected.
   * For example:
   *   <Navbar
   *     user={data.user}
   *     environments={data.environments}
   *     currentEnvironmentId={data.currentEnvironmentId ?? ''}
   *   />
   */
  export let currentEnvironmentId: string | null = null;

  // Local state
  let showDropdown = false;

  // References for click-outside detection
  let dropdownRef: HTMLElement;
  let avatarRef: HTMLElement;

  // Toggle the user dropdown
  function toggleDropdown() {
    showDropdown = !showDropdown;
  }

  // Close user dropdown if clicking outside it
  function handleClickOutside(event: MouseEvent) {
    if (!dropdownRef || !avatarRef) return;
    if (
      !dropdownRef.contains(event.target as Node) &&
      !avatarRef.contains(event.target as Node)
    ) {
      showDropdown = false;
    }
  }

  // Compute the user's initial (e.g., first letter of first name)
  $: userInitial = user?.first_name ? user.first_name[0].toUpperCase() : '?';

  /**
   * Handle environment selection.
   * - If user picks '', navigate to /environments
   * - Else navigate to /environments/<envId>/databases
   */
  function handleSelectEnv(e: Event) {
    const envId = (e.target as HTMLSelectElement).value;
    if (!envId) {
      // User chose "Select an environment"
      goto('/environments');
    } else {
      goto(`/environments/${envId}/databases`);
    }
  }

  // Clicking the Monji logo => /environments
  function handleLogoClick() {
    goto('/environments');
  }
</script>

<!-- Listen for clicks on the window to detect outside clicks -->
<svelte:window on:click={handleClickOutside} />

<nav class="bg-white border-b border-gray-200 px-4 py-3 flex items-center justify-between">
  <!-- Left side: Logo + Environment select -->
  <div class="flex items-center space-x-4">
    <!-- Clickable logo -->
    <img
      src="https://monji-assets.fra1.cdn.digitaloceanspaces.com/images/monji-logo-black.png"
      alt="Monji logo"
      class="h-8 w-auto cursor-pointer"
      on:click={handleLogoClick}
    />

    <!-- Environments dropdown -->
    <div class="relative">
      <select
        class="border border-gray-300 rounded px-3 py-1 focus:outline-none focus:ring-2 focus:ring-blue-600"
        bind:value={currentEnvironmentId}
        on:change={handleSelectEnv}
      >
        <!-- The first option: "Select an environment" -->
        <option value="">Select an environment</option>

        {#each environments as env}
          <option value={String(env.id)}>
            {env.name}
          </option>
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