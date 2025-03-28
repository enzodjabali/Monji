<!-- src/routes/+layout.svelte -->
<script lang="ts">
    import '../app.css';
    import { navigating } from '$app/stores';
  
    // Reactive declaration: whenever $navigating is non-null,
    // isNavigating becomes true, and vice versa.
    $: isNavigating = $navigating !== null;
</script>
  
<!-- The loading bar at the top -->
<div class="loading-bar" class:active={isNavigating}></div>
  
<!-- Your actual app content -->
<slot />
  
<style>
    /* A thin bar fixed at the top of the screen */
    .loading-bar {
      position: fixed;
      top: 0;
      left: 0;
      width: 0%;              /* Start hidden */
      height: 4px;            /* Adjust thickness here */
      background-color: #1B6609; /* A green color, pick whatever you like */
      transition: width 0.3s ease-in-out;
      z-index: 9999;          /* Ensure it stays above other elements */
    }
    .loading-bar.active {
      width: 100%;            /* Expand to full width while navigating */
    }
</style>