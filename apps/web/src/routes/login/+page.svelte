<script lang="ts">
    // This interface will be populated by server actions (e.g., for login errors)
    export let data: { error?: string };
    </script>
    
<!-- Outer container: two-column layout on larger screens -->
<div class="min-h-screen md:flex bg-white">
    <!-- Left column: login form -->
    <div class="flex flex-col justify-center px-8 py-12 md:w-1/2">
        <div class="max-w-sm mx-auto w-full">
          <!-- Monji logo -->
          <img
            src="https://monji-assets.fra1.cdn.digitaloceanspaces.com/images/monji-logo-black.png"
            alt="Monji logo"
            class="h-20 w-auto mb-6"
          />
    
          <!-- Heading -->
          <h1 class="text-2xl font-bold mb-2">Log in to your account</h1>
          <p class="mb-6">
            First time here?
            <a href="#" class="text-[#1B6609] hover:underline">
              Check out our newcomers documentation!
            </a>
          </p>
    
          <!-- Divider -->
          <div class="flex items-center my-6">
            <hr class="flex-grow border-t border-gray-300" />
            <span class="mx-2 text-gray-400">Login</span>
            <hr class="flex-grow border-t border-gray-300" />
          </div>
    
          <!-- Actual login form -->
          <form method="post" class="space-y-5">
            {#if data.error}
              <p class="text-red-500 font-semibold">{data.error}</p>
            {/if}
    
            <div>
              <label for="email" class="block font-semibold mb-1">Email</label>
              <input
                type="email"
                id="email"
                name="email"
                required
                class="w-full border border-gray-300 rounded px-3 py-2
                       focus:outline-none focus:ring-2 focus:ring-[#1B6609]"
              />
            </div>
    
            <div>
              <label for="password" class="block font-semibold mb-1">Password</label>
              <input
                type="password"
                id="password"
                name="password"
                required
                class="w-full border border-gray-300 rounded px-3 py-2
                       focus:outline-none focus:ring-2 focus:ring-[#1B6609]"
              />
            </div>
    
            <button
              type="submit"
              class="w-full bg-[#1B6609] text-white py-2 rounded
                     hover:bg-green-800 transition"
            >
              Login
            </button>
          </form>
        </div>
      </div>
    
      <!-- Right column: marketing / info panel with SVG background -->
      <div
        class="hidden md:flex md:w-1/2 items-center justify-center p-10"
        style="
            background-color: #1B6609;
            background-image: url('https://monji-assets.fra1.cdn.digitaloceanspaces.com/images/background-waves-white.svg');
            background-repeat: no-repeat;
            background-position: right bottom;
            background-size: contain;
        "
        >
        <div class="max-w-sm text-white">
          <h2 class="text-3xl font-bold mb-4">MongoDB 8.0 is here</h2>
          <p class="mb-4">
            Up to 32% higher throughput, improved horizontal scaling,
            expanded queryable encryption capabilities, and more.
          </p>
          <a href="#" class="underline hover:text-gray-200">
            See everything that's new &rarr;
          </a>
        </div>
    </div>
</div>