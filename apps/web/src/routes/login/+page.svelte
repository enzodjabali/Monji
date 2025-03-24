<script lang="ts">
// This interface will be populated by server actions (e.g., for login errors)
export let data: { error?: string };
</script>

<!-- Outer container: two-column layout on larger screens -->
<div class="min-h-screen md:flex bg-white">
    <!-- Left column: login form and social sign-in -->
    <div class="flex flex-col justify-center px-8 py-12 md:w-1/2">
        <div class="max-w-sm mx-auto w-full">
        <!-- Monji logo -->
        <img
            src="https://monji-assets.fra1.cdn.digitaloceanspaces.com/images/monji-logo-black.png"
            alt="Monji logo"
            class="h-20 w-auto mb-6"
        />
        <!-- Logo / heading -->
        <h1 class="text-2xl font-bold mb-2">Log in to your account</h1>
        <p class="mb-6">
            First time here?
            <a href="#" class="text-green-600 hover:underline">Check out our newcomers documentation!</a>
        </p>
        <!-- <p class="mb-6">
            Don’t have an account?
            <a href="/signup" class="text-green-600 hover:underline">Sign Up</a>
        </p> -->

        <!-- Social login buttons (example) -->
        <!-- <div class="space-y-3">
            <button
            type="button"
            class="w-full flex items-center justify-center gap-2 py-2 border border-gray-300 rounded hover:bg-gray-50 transition"
            > -->
            <!-- Example Google icon (SVG or from an icon library) -->
            <!-- <svg
                class="w-5 h-5"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
            >
                <path
                d="M12 11.988h10.88c.08.535.12 1.08.12 1.63 0 5.52-3.69 9.45-9 9.45-5.22 0-9.48-4.27-9.48-9.53s4.26-9.53 9.48-9.53c2.6 0 4.78.97 6.44 2.55l-2.61 2.51c-.71-.68-1.94-1.46-3.83-1.46-3.24 0-5.88 2.68-5.88 5.93s2.64 5.93 5.88 5.93c3.75 0 5.16-2.7 5.39-4.09H12v-3.43z"
                />
            </svg>
            <span>Sign in with Google</span>
            </button>
            <button
            type="button"
            class="w-full flex items-center justify-center gap-2 py-2 border border-gray-300 rounded hover:bg-gray-50 transition"
            > -->
            <!-- Example GitHub icon (SVG or from an icon library) -->
            <!-- <svg
                class="w-5 h-5"
                fill="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
            >
                <path
                d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.387.6.113.82-.258.82-.577 
                    0-.285-.01-1.04-.016-2.04-3.338.726-4.042-1.416-4.042-1.416-.546-1.387-1.333-1.757-1.333-1.757-1.088-.745.084-.73.084-.73 
                    1.205.084 1.84 1.236 1.84 1.236 1.07 1.835 2.807 1.304 3.492.997.107-.776.42-1.305.762-1.605-2.665-.303-5.466-1.334-5.466-5.93 
                    0-1.31.47-2.382 1.236-3.22-.124-.304-.536-1.527.116-3.176 0 0 1.008-.322 3.301 1.23a11.52 11.52 0 0 1 3.003-.403c1.018.005 
                    2.042.137 3.003.403 2.291-1.552 3.298-1.23 3.298-1.23.653 1.649.241 2.872.118 3.176.77.838 1.236 1.91 1.236 3.22 
                    0 4.61-2.803 5.625-5.476 5.92.43.37.81 1.1.81 2.22 0 1.604-.015 2.896-.015 3.286 0 .317.22.686.825.57C20.565 
                    22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"
                />
            </svg>
            <span>Sign in with GitHub</span>
            </button>
        </div> -->

        <!-- Divider -->
        <!-- <div class="flex items-center my-6">
            <hr class="flex-grow border-t border-gray-300" />
            <span class="mx-2 text-gray-400">or</span>
            <hr class="flex-grow border-t border-gray-300" />
        </div> -->
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
                class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            />
            </div>

            <div>
            <label for="password" class="block font-semibold mb-1">Password</label>
            <input
                type="password"
                id="password"
                name="password"
                required
                class="w-full border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-green-500"
            />
            </div>

            <button
            type="submit"
            class="w-full bg-green-600 text-white py-2 rounded hover:bg-green-700 transition"
            >
            Login
            </button>
        </form>
        </div>
    </div>

    <!-- Right column: can display an image, illustration, or marketing info -->
    <div class="hidden md:flex md:w-1/2 bg-gray-50 items-center justify-center p-10">
        <!-- Example marketing content -->
        <div class="max-w-sm">
        <h2 class="text-3xl font-bold mb-4">MongoDB 8.0 is here</h2>
        <p class="text-gray-700 mb-4">
            Up to 32% higher throughput, improved horizontal scaling, expanded
            queryable encryption capabilities, and more.
        </p>
        <a href="#" class="text-green-600 hover:underline">
            See everything that's new &rarr;
        </a>
        </div>
    </div>
</div>