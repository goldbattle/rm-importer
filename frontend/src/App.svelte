<script lang="ts">
  import Router from 'svelte-spa-router';
  import Menu from '../components/Menu.svelte';
  import RmFiles from '../components/RmFiles.svelte';
  const routes = {
    '/': Menu,
    '/rm-files': RmFiles,
  };
</script>

<main>
  <Router {routes} />
</main>