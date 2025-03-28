<script lang="ts">
  /**
   * We show up to four breadcrumb steps:
   *
   * 1) "Environments" (always) - links to /environments
   * 2) environmentName (if environmentId && environmentName)
   * 3) databaseName    (if databaseName)
   * 4) collectionName  (if collectionName)
   *
   * In each case, we link to the relevant route.
   */

  export let environmentId: string | null = null;
  export let environmentName: string | null = null;
  export let databaseName: string | null = null;
  export let collectionName: string | null = null;
  export let documentId: string | null = null;
</script>

<nav class="flex items-center space-x-2 text-sm text-gray-700 my-4" aria-label="Breadcrumb">
  <ol class="inline-flex items-center space-x-1 md:space-x-3">
    <!-- 1) Always link to /environments -->
    <li>
      <div class="flex items-center">
        <!-- Icon - you can change it to any symbol/text you want -->
        <svg
          class="w-4 h-4 mr-1 text-gray-400"
          fill="currentColor"
          viewBox="0 0 20 20"
          aria-hidden="true"
        >
          <path d="M10 3.172l5.95 5.95-1.414 1.414L10 5.586l-4.536 4.95L4.05 9.121z" />
        </svg>
        <a
          href="/environments"
          class="text-blue-600 hover:underline font-medium"
        >
          Environments
        </a>
      </div>
    </li>

    <!-- 2) environmentName -->
    {#if environmentId && environmentName}
      <li>
        <div class="flex items-center">
          <svg
            class="w-6 h-6 text-gray-400 mx-1"
            fill="currentColor"
            viewBox="0 0 20 20"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M7.293 14.707a1 1 0 010-1.414L10.586 10 
                 7.293 6.707a1 1 0 011.414-1.414l4 4a1 
                 1 0 010 1.414l-4 4a1 1 0 
                 01-1.414 0z"
              clip-rule="evenodd"
            />
          </svg>
          <a
            href={`/environments/${environmentId}/databases`}
            class="text-blue-600 hover:underline font-medium"
          >
            {environmentName}
          </a>
        </div>
      </li>
    {/if}

    <!-- 3) databaseName -->
    {#if environmentId && databaseName}
      <li>
        <div class="flex items-center">
          <svg
            class="w-6 h-6 text-gray-400 mx-1"
            fill="currentColor"
            viewBox="0 0 20 20"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M7.293 14.707a1 1 0 010-1.414L10.586 10 
                 7.293 6.707a1 1 0 011.414-1.414l4 4a1 
                 1 0 010 1.414l-4 4a1 1 0 
                 01-1.414 0z"
              clip-rule="evenodd"
            />
          </svg>
          <a
            href={`/environments/${environmentId}/databases/${databaseName}/collections`}
            class="text-blue-600 hover:underline font-medium"
          >
            {databaseName}
          </a>
        </div>
      </li>
    {/if}

    <!-- 4) collectionName -->
    {#if environmentId && databaseName && collectionName}
      <li>
        <div class="flex items-center">
          <svg
            class="w-6 h-6 text-gray-400 mx-1"
            fill="currentColor"
            viewBox="0 0 20 20"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M7.293 14.707a1 1 0 010-1.414L10.586 10 
                 7.293 6.707a1 1 0 011.414-1.414l4 4a1 
                 1 0 010 1.414l-4 4a1 1 0 
                 01-1.414 0z"
              clip-rule="evenodd"
            />
          </svg>
          <a
            href={`/environments/${environmentId}/databases/${databaseName}/collections/${collectionName}/documents`}
            class="text-blue-600 hover:underline font-medium"
          >
            {collectionName}
          </a>
        </div>
      </li>
    {/if}
    <!-- 5) documentId -->
    {#if environmentId && databaseName && collectionName && documentId}
      <li>
        <div class="flex items-center">
          <svg
            class="w-6 h-6 text-gray-400 mx-1"
            fill="currentColor"
            viewBox="0 0 20 20"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M7.293 14.707a1 1 0 010-1.414L10.586 10 
                 7.293 6.707a1 1 0 011.414-1.414l4 4a1 
                 1 0 010 1.414l-4 4a1 1 0 
                 01-1.414 0z"
              clip-rule="evenodd"
            />
          </svg>
          <a
            href={`/environments/${environmentId}/databases/${databaseName}/collections/${collectionName}/documents/${documentId}`}
            class="text-blue-600 hover:underline font-medium"
          >
            {documentId}
          </a>
        </div>
      </li>
    {/if}
  </ol>
</nav>