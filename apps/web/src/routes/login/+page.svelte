<script lang="ts">
    // `data.error` will be passed back if the login fails
    export let data: { error?: string };
</script>
  
<h1>Login</h1>
{#if data.error}
    <p style="color: red;">{data.error}</p>
{/if}
  
<form method="post">
    <div>
      <label for="email">Email:</label>
      <input type="email" name="email" id="email" required />
    </div>
    <div>
      <label for="password">Password:</label>
      <input type="password" name="password" id="password" required />
    </div>
    <button type="submit">Login</button>
</form>