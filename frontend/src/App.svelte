<script lang="ts">
  import "./app.css";
  import Menu from "./pages/Menu.svelte";
  import FileSelection from "./pages/FileSelection.svelte";
  import Router from "svelte-spa-router";

  const routes = {
    '/': Menu,
    '/files': FileSelection,
  };
</script>

<Router {routes}/>