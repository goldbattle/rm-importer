<script lang="ts">
  import "./app.css";
  import Router from 'svelte-spa-router';
  import Menu from "./Menu.svelte";

  const routes = {
    '/': Menu,
  }

</script>

<Router {routes}/>