<script lang="ts">
  import "./app.css";
  import Menu from "./pages/Menu.svelte";
  import FileSelection from "./pages/FileSelection.svelte";
  import Router from "svelte-spa-router";
  import Export from "./pages/Export.svelte";

  const routes = {
    '/': Menu,
    '/files': FileSelection,
    '/export': Export,
    '*': Menu,
  };
</script>

<Router {routes}/>