<script lang="ts">
    export let data: { error?: string } | undefined;
</script>
  
<h1>Login</h1>

{#if data?.error}
    <p style="color: red">{data.error}</p>
{/if}

<form method="post">
    <label>
        Email:
        <input type="email" name="email" required />
    </label>
    <br />
    <label>
        Password:
        <input type="password" name="password" required />
    </label>
    <br />
    <button type="submit">Login</button>
</form>