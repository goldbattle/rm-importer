<script lang="ts">
  import Router from 'svelte-spa-router';
  import Menu from './pages/Menu.svelte';
  import RmFiles from './pages/RmFiles.svelte';
  const routes = {
    '/': Menu,
    '/rm-files': RmFiles,
  };
</script>

<main>
  <Router {routes} />
</main>