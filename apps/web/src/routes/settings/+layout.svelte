<script lang="ts">
    /**
     * Data returned from +layout.server.ts
     *    user:        { id, first_name, last_name, email, role }
     *    permissions: { environments: [...], databases: [...] }
     */
    export let data: {
      user: {
        id: number;
        first_name: string;
        last_name: string;
        email: string;
        role: string;
      };
      permissions: {
        environments: Array<{
          environment_id: number;
          environment_name: string;
          permission: string;
        }>;
        databases: Array<{
          environment_id: number;
          environment_name: string;
          db_name: string;
          permission: string;
        }>;
      };
    };
  </script>
  
  <div class="min-h-screen flex bg-gray-100">
    <!-- LEFT SIDEBAR -->
    <aside class="w-64 bg-white shadow">
      <div class="p-6 border-b">
        <h2 class="text-2xl font-bold">Settings</h2>
      </div>
      <nav class="p-6 space-y-2">
        <!-- Link(s) for different settings pages -->
        <a
          href="/settings/profile"
          class="block px-3 py-2 rounded hover:bg-gray-100 text-gray-700 transition"
        >
          Profile
        </a>
        <!-- Add more links here for other Settings pages in the future -->
      </nav>
    </aside>
  
    <!-- MAIN CONTENT AREA -->
    <main class="flex-1 p-6">
      <!-- Render whichever page is active: e.g. /settings/profile -->
      <slot />
    </main>
  </div>