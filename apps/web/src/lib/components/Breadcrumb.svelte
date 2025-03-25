<script lang="ts">
    /**
     * We show up to four breadcrumb levels:
     * 1) Environments (always shown, links to /environments)
     * 2) Databases (requires environmentId)
     * 3) Collections (requires environmentId + databaseName)
     * 4) Documents  (requires environmentId + databaseName + collectionName)
     *
     * Pass the IDs/names you have. If a prop is omitted or null, that level is skipped.
     */
    export let environmentId: string | null = null;
    export let databaseName: string | null = null;
    export let collectionName: string | null = null;
  </script>
  
  <nav class="text-sm text-gray-500 mb-4">
    <ol class="flex items-center space-x-2">
      <!-- Always show "Environments" linking to /environments -->
      <li>
        <a href="/environments" class="text-blue-600 hover:underline">
          Environments
        </a>
      </li>
  
      <!-- If environmentId is present, show "Databases" link -->
      {#if environmentId}
        <li>/</li>
        <li>
          <a
            href={`/environments/${environmentId}/databases`}
            class="text-blue-600 hover:underline"
          >
            Databases
          </a>
        </li>
      {/if}
  
      <!-- If databaseName is present, show "Collections" link -->
      {#if environmentId && databaseName}
        <li>/</li>
        <li>
          <a
            href={`/environments/${environmentId}/databases/${databaseName}/collections`}
            class="text-blue-600 hover:underline"
          >
            Collections
          </a>
        </li>
      {/if}
  
      <!-- If collectionName is present, show "Documents" link -->
      {#if environmentId && databaseName && collectionName}
        <li>/</li>
        <li>
          <a
            href={`/environments/${environmentId}/databases/${databaseName}/collections/${collectionName}/documents`}
            class="text-blue-600 hover:underline"
          >
            Documents
          </a>
        </li>
      {/if}
    </ol>
  </nav>