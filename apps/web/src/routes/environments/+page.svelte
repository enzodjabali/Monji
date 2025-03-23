<script lang="ts">
    // `data.environments` is the array returned by the load function
    export let data: { environments: { id: number; name: string; connection_string: string; created_by: number }[] };
</script>
  
<h1>Environments</h1>
{#if data.environments.length > 0}
    <ul>
      {#each data.environments as env}
        <li>
          <strong>{env.name}</strong>
          <div>ID: {env.id}</div>
          <div>Connection: {env.connection_string}</div>
          <div>Created by: {env.created_by}</div>
        </li>
      {/each}
    </ul>
  {:else}
    <p>No environments available.</p>
{/if}