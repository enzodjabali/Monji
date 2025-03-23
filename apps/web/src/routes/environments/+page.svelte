<script lang="ts">
    export let data: {
      environments: Array<{
        id: number;
        name: string;
        connection_string: string;
        created_by: number;
      }>;
    };
</script>
  
<h1>Environments</h1>
  
{#if data.environments && data.environments.length > 0}
    <ul>
      {#each data.environments as env}
        <li>
          <strong>{env.name}</strong> – {env.connection_string}
        </li>
      {/each}
    </ul>
  {:else}
    <p>No environments found.</p>
{/if}